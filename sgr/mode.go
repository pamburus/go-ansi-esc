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
