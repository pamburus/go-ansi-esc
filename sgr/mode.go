package sgr

import (
	"fmt"
	"strings"

	"github.com/veggiemonk/strcase"
)

// ---

// All supported Mode values.
const (
	Bold Mode = iota
	Faint
	Italic
	SlowBlink
	RapidBlink
	Reversed
	Concealed
	CrossedOut
	Underlined
	DoublyUnderlined
	Framed
	Encircled
	Overlined
	Superscript
	Subscript
)

// ---

// Mode is flag number describing one of terminal modes supported by CSI/SGR sequences.
type Mode uint8

// ModeSet converts m to a ModeSet containing only m.
func (m Mode) ModeSet() ModeSet {
	return NewModeSet().With(m)
}

// ModeList converts m to a ModeList containing only m.
func (m Mode) ModeList() ModeList {
	return NewModeList(m)
}

// String returns textual description of m that can be used for debugging or logging purposes.
func (m Mode) String() string {
	if name, ok := modeNames[m]; ok {
		return name
	}

	return fmt.Sprintf("<!0x%02x>", uint8(m))
}

// Validate check that m has a valid value.
func (m Mode) Validate() error {
	if _, ok := modeNames[m]; !ok {
		return ErrInvalidModeValue{m}
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler interface
// that allows Mode to be used in any compatible marshaler like JSON, YAML, etc.
func (m Mode) MarshalText() ([]byte, error) {
	if name, ok := modeNames[m]; ok {
		return []byte(name), nil
	}

	return nil, ErrInvalidModeValue{m}
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows Mode to be used in any compatible unmarshaler like JSON, YAML, etc.
func (m *Mode) UnmarshalText(data []byte) error {
	return m.unmarshalText(string(data))
}

func (m *Mode) unmarshalText(text string) error {
	text = strings.TrimSpace(text)
	try := func(text string) bool {
		value, ok := textToMode[text]
		if ok {
			*m = value
		}

		return ok
	}

	switch {
	case try(text):
		return nil
	case try(strcase.Pascal(text)):
		return nil
	default:
		return ErrInvalidModeText{text}
	}
}

// ---

// EmptyModeSet constructs a new empty ModeSet.
func EmptyModeSet() ModeSet {
	return ModeSet(0)
}

// NewModeSet constructs a new empty ModeSet.
func NewModeSet() ModeSet {
	return ModeSet(0)
}

// ModeSetWith constructs a new ModeSet with the given Mode values.
func ModeSetWith(modes ...Mode) ModeSet {
	return ModeList(modes).ModeSet()
}

// ---

// ModeSet is a bit mask containing values for all Mode values.
type ModeSet uint16

// IsZero returns true if s is empty.
func (s ModeSet) IsZero() bool {
	return s == 0
}

// IsEmpty returns true if s is empty.
func (s ModeSet) IsEmpty() bool {
	return s == 0
}

// With returns a copy of ModeSet with specified mode bit set to 1.
func (s ModeSet) With(mode Mode) ModeSet {
	return s | (1 << mode)
}

// Without returns a copy of ModeSet with specified mode bit reset to 0.
func (s ModeSet) Without(mode Mode) ModeSet {
	return s & ^(1 << mode)
}

// WithToggled returns a copy of ModeSet with specified mode bit changed to the opposite value.
func (s ModeSet) WithToggled(mode Mode) ModeSet {
	return s ^ (1 << mode)
}

// Has returns true if s includes mode.
func (s ModeSet) Has(mode Mode) bool {
	return s&(1<<mode) != 0
}

// WithOther returns a new ModeSet that is combined with s according to the specified action.
func (s ModeSet) WithOther(other ModeSet, action ModeAction) ModeSet {
	switch action {
	case ModeReplace:
		return other
	case ModeAdd:
		return s | other
	case ModeRemove:
		return s &^ other
	case ModeToggle:
		return s ^ other
	default:
		return s
	}
}

// Diff returns a ModeSetDiff with Old set to s and New set to other.
func (s ModeSet) Diff(other ModeSet) ModeSetDiff {
	return ModeSetDiff{
		Old: s,
		New: other,
	}
}

// ModeList converts ModeSet to ModeList.
func (s ModeSet) ModeList() ModeList {
	result := make(ModeList, 0, 16)
	for i := 0; i != 16; i++ {
		if s&(1<<i) != 0 {
			result = append(result, Mode(i))
		}
	}

	return result
}

// ---

// NewModeList constructs a new ModeList with the given modes.
func NewModeList(modes ...Mode) ModeList {
	return ModeList(modes)
}

// ModeList is a slice of Mode values.
type ModeList []Mode

// ModeSet constructs a new ModeSet with the values in l.
func (l ModeList) ModeSet() ModeSet {
	result := EmptyModeSet()
	for _, mode := range l {
		result = result.With(mode)
	}

	return result
}

// ---

// ModeAction is an operation that can be applied to mode bits of a ModeSet when combining two ModeSet values.
type ModeAction uint

// Valid values for ModeAction.
const (
	ModeReplace ModeAction = iota
	ModeAdd
	ModeRemove
	ModeToggle
)

// ---

// ModeSetDiff contains Old and New mode sets and allows to generate a command sequence based on them.
type ModeSetDiff struct {
	Old ModeSet
	New ModeSet
}

// Added returns modes that present in New but are not present in Old.
func (d ModeSetDiff) Added() ModeSet {
	return d.New &^ d.Old
}

// Removed returns modes that present in Old but are not present in New.
func (d ModeSetDiff) Removed() ModeSet {
	return d.Old &^ d.New
}

// Changed returns modes that differ in Old and New.
func (d ModeSetDiff) Changed() ModeSet {
	return d.Old ^ d.New
}

// Reversed returns a reversed mode set where New and Old are swapped.
func (d ModeSetDiff) Reversed() ModeSetDiff {
	return ModeSetDiff{Old: d.New, New: d.Old}
}

// ToCommands appends commands needed to bring old mode set to new mode set to seq and returns modified seq.
func (d ModeSetDiff) ToCommands(seq Sequence) Sequence {
	changed := d.Old ^ d.New
	if changed == 0 {
		return seq
	}

	for _, row := range modeSyncTableSingle {
		mask := row.mode.ModeSet()
		if changed&mask == 0 {
			continue
		}
		seq = append(seq, row.commands[(d.Old&mask)>>row.mode])
	}

	for _, row := range modeSyncTableDual {
		mask := row.mask
		if changed&mask == 0 {
			continue
		}

		oi := int(d.Old&mask) >> int(row.firstMode)
		ni := int(d.New&mask) >> int(row.firstMode)
		cmdMask := &modeSyncDualCommandMask[oi][ni]
		for i, cmd := range row.commands {
			if cmdMask[i] != 0 {
				seq = append(seq, cmd)
			}
		}
	}

	return seq
}

// ---

// Mode sync table for independent modes.
var modeSyncTableSingle = []struct {
	mode     Mode
	commands [2]Command
}{
	{
		Italic,
		[2]Command{
			SetItalic,
			ResetItalic,
		},
	},
	{
		Reversed,
		[2]Command{
			SetReversed,
			ResetReversed,
		},
	},
	{
		Concealed,
		[2]Command{
			SetConcealed,
			ResetConcealed,
		},
	},
	{
		CrossedOut,
		[2]Command{
			SetCrossedOut,
			ResetCrossedOut,
		},
	},
	{
		Overlined,
		[2]Command{
			SetOverlined,
			ResetOverlined,
		},
	},
}

var modeSyncTableDual = []struct {
	firstMode Mode
	mask      ModeSet
	commands  [3]Command
}{
	{
		Bold,
		NewModeSet().
			With(Bold).
			With(Faint),
		[3]Command{
			ResetBoldAndFaint,
			SetBold,
			SetFaint,
		},
	},
	{
		SlowBlink,
		NewModeSet().
			With(SlowBlink).
			With(RapidBlink),
		[3]Command{
			ResetAllBlinks,
			SetSlowBlink,
			SetRapidBlink,
		},
	},
	{
		Framed,
		NewModeSet().
			With(Framed).
			With(Encircled),
		[3]Command{
			ResetFramedAndEncircled,
			SetFramed,
			SetEncircled,
		},
	},
	{
		Superscript,
		NewModeSet().
			With(Superscript).
			With(Subscript),
		[3]Command{
			ResetSuperscriptAndSubscript,
			SetSuperscript,
			SetSubscript,
		},
	},
	{
		Underlined,
		NewModeSet().
			With(Underlined).
			With(DoublyUnderlined),
		[3]Command{
			ResetAllUnderlines,
			SetUnderlined,
			SetDoublyUnderlined,
		},
	},
}

var modeSyncDualCommandMask = [4][4][3]int{
	{{0, 0, 0}, {0, 1, 0}, {0, 0, 1}, {0, 1, 1}},
	{{1, 0, 0}, {0, 0, 0}, {1, 0, 1}, {0, 0, 1}},
	{{1, 0, 0}, {1, 1, 0}, {0, 0, 0}, {0, 1, 0}},
	{{1, 0, 0}, {1, 1, 0}, {1, 0, 1}, {0, 0, 0}},
}

// ---

var modeNames = map[Mode]string{
	Bold:             "Bold",
	Faint:            "Faint",
	Italic:           "Italic",
	SlowBlink:        "SlowBlink",
	RapidBlink:       "RapidBlink",
	Reversed:         "Reversed",
	Concealed:        "Concealed",
	CrossedOut:       "CrossedOut",
	Underlined:       "Underlined",
	DoublyUnderlined: "DoublyUnderlined",
	Framed:           "Framed",
	Encircled:        "Encircled",
	Overlined:        "Overlined",
	Superscript:      "Superscript",
	Subscript:        "Subscript",
}

var textToMode = map[string]Mode{
	"Bold":             Bold,
	"Faint":            Faint,
	"Italic":           Italic,
	"SlowBlink":        SlowBlink,
	"RapidBlink":       RapidBlink,
	"Reversed":         Reversed,
	"Concealed":        Concealed,
	"CrossedOut":       CrossedOut,
	"Underlined":       Underlined,
	"DoublyUnderlined": DoublyUnderlined,
	"Framed":           Framed,
	"Encircled":        Encircled,
	"Overlined":        Overlined,
	"Superscript":      Superscript,
	"Subscript":        Subscript,
}
