package sgr

import "io"

// NewWriter constructs a new Writer over the given target writer.
func NewWriter(target io.Writer) *Writer {
	p := &Writer{target: target}
	p.head = defaultState
	p.upstream = defaultState
	p.stack.bgc = make([]Color, 0, 8)
	p.stack.fgc = make([]Color, 0, 8)
	p.stack.ulc = make([]Color, 0, 8)
	p.stack.modes = make([]ModeSet, 0, 8)
	p.scratchCommands = make(Sequence, 0, 8)
	p.scratchBytes = make([]byte, 128)

	return p
}

// Writer extends functionality of the target writer
// by adding methods for changing color and style using SGR commands.
//
// Writer is using lazy approach for rendering SGR sequences that means
// it calculates an renders needed sequence only when its Write or Flush method is called.
// Until that only current values are remembered.
type Writer struct {
	target   io.Writer
	head     state
	upstream state
	stack    struct {
		bgc   []Color
		fgc   []Color
		ulc   []Color
		modes []ModeSet
	}
	scratchCommands Sequence
	scratchBytes    []byte
}

// Reset resets current SGR state to terminal defaults.
func (w *Writer) Reset() {
	w.head = defaultState
}

// SetBackgroundColor changes current background color.
func (w *Writer) SetBackgroundColor(color IntoColor) {
	w.head.bgc = color.Color()
}

// PushBackgroundColor changes background color and pushes old value to a stack
// so that it can be restored using PopBackgroundColor method.
func (w *Writer) PushBackgroundColor(color IntoColor) {
	w.stack.bgc = append(w.stack.bgc, w.head.bgc)
	w.SetBackgroundColor(color)
}

// PopBackgroundColor restores old background color that was saved at last PushBackgroundColor call.
func (w *Writer) PopBackgroundColor() {
	i := len(w.stack.bgc) - 1
	w.head.bgc = w.stack.bgc[i]
	w.stack.bgc = w.stack.bgc[:i]
}

// SetForegroundColor changes current foreground color.
func (w *Writer) SetForegroundColor(color IntoColor) {
	w.head.fgc = color.Color()
}

// PushForegroundColor changes foreground color and pushes old value to a stack
// so that it can be restored using PopForegroundColor method.
func (w *Writer) PushForegroundColor(color IntoColor) {
	w.stack.fgc = append(w.stack.fgc, w.head.fgc)
	w.SetForegroundColor(color)
}

// PopForegroundColor restores old foreground color that was saved at last PushForegroundColor call.
func (w *Writer) PopForegroundColor() {
	i := len(w.stack.fgc) - 1
	w.head.fgc = w.stack.fgc[i]
	w.stack.fgc = w.stack.fgc[:i]
}

// SetUnderlineColor changes current underline color.
func (w *Writer) SetUnderlineColor(color IntoColor) {
	w.head.ulc = color.Color()
}

// PushUnderlineColor changes underline color and pushes old value to a stack
// so that it can be restored using PopUnderlineColor method.
func (w *Writer) PushUnderlineColor(color IntoColor) {
	w.stack.ulc = append(w.stack.ulc, w.head.ulc)
	w.SetUnderlineColor(color)
}

// PopUnderlineColor restores old underline color that was saved at last PushUnderlineColor call.
func (w *Writer) PopUnderlineColor() {
	i := len(w.stack.ulc) - 1
	w.head.ulc = w.stack.ulc[i]
	w.stack.ulc = w.stack.ulc[:i]
}

// SetModes changes current modes using specified modes and action to calculate new modes.
func (w *Writer) SetModes(modes ModeSet, action ModeAction) {
	w.head.modes = w.head.modes.WithOther(modes, action)
}

// PushModes changes current modes using specified modes and action to calculate new modes
// and pushes old value to a stack so that it can be restored using PopModes method.
func (w *Writer) PushModes(modes ModeSet, action ModeAction) {
	w.stack.modes = append(w.stack.modes, w.head.modes)
	w.SetModes(modes, action)
}

// PopModes restores old modes that were saved at last PushModes call.
func (w *Writer) PopModes() {
	i := len(w.stack.modes) - 1
	w.head.modes = w.stack.modes[i]
	w.stack.modes = w.stack.modes[:i]
}

// Write flushes current style changes by generating CSI/SGR sequence and writing it
// to the target writer and then finally writes the given data to it.
func (w *Writer) Write(data []byte) (n int, err error) {
	err = w.sync()
	if err != nil {
		return 0, err
	}

	return w.target.Write(data)
}

// Flush just flushes current style changes by generating CSI/SGR sequence and writing it
// to the target writer.
// It is recommended to call Flush at the end of writing a line or a stream.
func (w *Writer) Flush() error {
	return w.sync()
}

func (w *Writer) sync() error {
	seq := w.scratchCommands[0:0]
	buf := w.scratchBytes[0:0]

	w.head.bgc = w.head.bgc.OrDefault()
	w.head.fgc = w.head.fgc.OrDefault()
	w.head.ulc = w.head.ulc.OrDefault()

	if w.head != w.upstream {
		if w.head == defaultState {
			seq = append(seq, ResetAll)
			w.upstream = w.head
		} else {
			if w.head.bgc != w.upstream.bgc {
				seq = append(seq, setBackgroundColor(w.head.bgc))
				w.upstream.bgc = w.head.bgc
			}
			if w.head.fgc != w.upstream.fgc {
				seq = append(seq, setForegroundColor(w.head.fgc))
				w.upstream.fgc = w.head.fgc
			}
			if w.head.ulc != w.upstream.ulc {
				seq = append(seq, setUnderlineColor(w.head.ulc))
				w.upstream.ulc = w.head.ulc
			}
			if w.head.modes != w.upstream.modes {
				seq = modeDiffToCommands(w.upstream.modes, w.head.modes, seq)
				w.upstream.modes = w.head.modes
			}
		}
	}

	if len(seq) != 0 {
		buf = seq.Render(buf)
		_, err := w.target.Write(buf)
		if err != nil {
			return err
		}
	}

	w.scratchCommands = seq[0:0]
	w.scratchBytes = buf[0:0]

	return nil
}

// ---

type state struct {
	bgc   Color
	fgc   Color
	ulc   Color
	modes ModeSet
}

var defaultState = state{
	bgc: Default.Color(),
	fgc: Default.Color(),
	ulc: Default.Color(),
}

func modeDiffToCommands(old, new ModeSet, seq Sequence) Sequence {
	changed := old ^ new

	for _, row := range modeSyncTableSingle {
		mask := row.mode.ModeSet()
		if changed&mask == 0 {
			continue
		}
		seq = append(seq, row.commands[(old&mask)>>row.mode])
	}

	for _, row := range modeSyncTableDual {
		mask := row.mask
		if changed&mask == 0 {
			continue
		}

		oi := int(old&mask) >> int(row.firstMode)
		ni := int(new&mask) >> int(row.firstMode)
		cmdMask := &modeSyncDualCommandMask[oi][ni]
		for i, cmd := range row.commands {
			if cmdMask[i] != 0 {
				seq = append(seq, cmd)
			}
		}
	}

	return seq
}

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
