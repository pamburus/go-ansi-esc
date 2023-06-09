package sgr

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

// ---

// NewModeSet constructs a new empty ModeSet.
func NewModeSet() ModeSet {
	return ModeSet(0)
}

// ---

// ModeSet is a bit mask containing values for all Mode values.
type ModeSet uint16

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

// ToCommands appends commands needed to bring old mode set to new mode set to seq and returns modified seq.
func (d ModeSetDiff) ToCommands(seq Sequence) Sequence {
	changed := d.Old ^ d.New

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
