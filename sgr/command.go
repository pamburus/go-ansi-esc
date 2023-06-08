package sgr

import (
	"fmt"
)

// ---

// SetForegroundColor returns a command that will change foreground color
// when sent to a terminal in as part of a Sequence.
func SetForegroundColor(color IntoColor) Command {
	return setForegroundColor(color.Color())
}

// SetBackgroundColor returns a command that will change background color
// when sent to a terminal in as part of a Sequence.
func SetBackgroundColor(color IntoColor) Command {
	return setBackgroundColor(color.Color())
}

// SetUnderlineColor returns a command that will change underline color
// when sent to a terminal in as part of a Sequence.
func SetUnderlineColor(color IntoColor) Command {
	return setUnderlineColor(color.Color())
}

// ---

// Complete set of simple SGR commands without arguments.
const (
	ResetAll                        = commandValid | Command(CodeResetAll)
	SetBold                         = commandValid | Command(CodeSetBold)
	SetFaint                        = commandValid | Command(CodeSetFaint)
	SetItalic                       = commandValid | Command(CodeSetItalic)
	SetUnderlined                   = commandValid | Command(CodeSetUnderlined)
	SetSlowBlink                    = commandValid | Command(CodeSetSlowBlink)
	SetRapidBlink                   = commandValid | Command(CodeSetRapidBlink)
	SetReversed                     = commandValid | Command(CodeSetReversed)
	SetConcealed                    = commandValid | Command(CodeSetConcealed)
	SetCrossedOut                   = commandValid | Command(CodeSetCrossedOut)
	SetDoublyUnderlined             = commandValid | Command(CodeSetDoublyUnderlined)
	ResetBoldAndFaint               = commandValid | Command(CodeResetBoldAndFaint)
	ResetItalic                     = commandValid | Command(CodeResetItalic)
	ResetAllUnderlines              = commandValid | Command(CodeResetAllUnderlines)
	ResetAllBlinks                  = commandValid | Command(CodeResetAllBlinks)
	ResetReversed                   = commandValid | Command(CodeResetReversed)
	ResetConcealed                  = commandValid | Command(CodeResetConcealed)
	ResetCrossedOut                 = commandValid | Command(CodeResetCrossedOut)
	SetForegroundColorBlack         = commandValid | Command(CodeSetForegroundColorBlack)
	SetForegroundColorRed           = commandValid | Command(CodeSetForegroundColorRed)
	SetForegroundColorGreen         = commandValid | Command(CodeSetForegroundColorGreen)
	SetForegroundColorYellow        = commandValid | Command(CodeSetForegroundColorYellow)
	SetForegroundColorBlue          = commandValid | Command(CodeSetForegroundColorBlue)
	SetForegroundColorMagenta       = commandValid | Command(CodeSetForegroundColorMagenta)
	SetForegroundColorCyan          = commandValid | Command(CodeSetForegroundColorCyan)
	SetForegroundColorWhite         = commandValid | Command(CodeSetForegroundColorWhite)
	ResetForegroundColor            = commandValid | Command(CodeResetForegroundColor)
	SetBackgroundColorBlack         = commandValid | Command(CodeSetBackgroundColorBlack)
	SetBackgroundColorRed           = commandValid | Command(CodeSetBackgroundColorRed)
	SetBackgroundColorGreen         = commandValid | Command(CodeSetBackgroundColorGreen)
	SetBackgroundColorYellow        = commandValid | Command(CodeSetBackgroundColorYellow)
	SetBackgroundColorBlue          = commandValid | Command(CodeSetBackgroundColorBlue)
	SetBackgroundColorMagenta       = commandValid | Command(CodeSetBackgroundColorMagenta)
	SetBackgroundColorCyan          = commandValid | Command(CodeSetBackgroundColorCyan)
	SetBackgroundColorWhite         = commandValid | Command(CodeSetBackgroundColorWhite)
	ResetBackgroundColor            = commandValid | Command(CodeResetBackgroundColor)
	SetFramed                       = commandValid | Command(CodeSetFramed)
	SetEncircled                    = commandValid | Command(CodeSetEncircled)
	SetOverlined                    = commandValid | Command(CodeSetOverlined)
	ResetFramedAndEncircled         = commandValid | Command(CodeResetFramedAndEncircled)
	ResetOverlined                  = commandValid | Command(CodeResetOverlined)
	ResetUnderlineColor             = commandValid | Command(CodeResetUnderlineColor)
	SetSuperscript                  = commandValid | Command(CodeSetSuperscript)
	SetSubscript                    = commandValid | Command(CodeSetSubscript)
	ResetSuperscriptAndSubscript    = commandValid | Command(CodeResetSuperscriptAndSubscript)
	SetForegroundColorBrightBlack   = commandValid | Command(CodeSetForegroundColorBrightBlack)
	SetForegroundColorBrightRed     = commandValid | Command(CodeSetForegroundColorBrightRed)
	SetForegroundColorBrightGreen   = commandValid | Command(CodeSetForegroundColorBrightGreen)
	SetForegroundColorBrightYellow  = commandValid | Command(CodeSetForegroundColorBrightYellow)
	SetForegroundColorBrightBlue    = commandValid | Command(CodeSetForegroundColorBrightBlue)
	SetForegroundColorBrightMagenta = commandValid | Command(CodeSetForegroundColorBrightMagenta)
	SetForegroundColorBrightCyan    = commandValid | Command(CodeSetForegroundColorBrightCyan)
	SetForegroundColorBrightWhite   = commandValid | Command(CodeSetForegroundColorBrightWhite)
	SetBackgroundColorBrightBlack   = commandValid | Command(CodeSetBackgroundColorBrightBlack)
	SetBackgroundColorBrightRed     = commandValid | Command(CodeSetBackgroundColorBrightRed)
	SetBackgroundColorBrightGreen   = commandValid | Command(CodeSetBackgroundColorBrightGreen)
	SetBackgroundColorBrightYellow  = commandValid | Command(CodeSetBackgroundColorBrightYellow)
	SetBackgroundColorBrightBlue    = commandValid | Command(CodeSetBackgroundColorBrightBlue)
	SetBackgroundColorBrightMagenta = commandValid | Command(CodeSetBackgroundColorBrightMagenta)
	SetBackgroundColorBrightCyan    = commandValid | Command(CodeSetBackgroundColorBrightCyan)
	SetBackgroundColorBrightWhite   = commandValid | Command(CodeSetBackgroundColorBrightWhite)
)

// Command is an SGR command that can be used to control terminal's SGR (Select Graphic Rendition) parameters.
type Command uint64

// String returns textual description of c that can be used for debugging or logging purposes.
func (c Command) String() string {
	if c.IsZero() {
		return ""
	}
	if !c.valid() {
		return fmt.Sprintf("<!0x%08x>", uint64(c))
	}

	code := CommandCode(c & 0xFF)
	switch code {
	case CodeSetBackgroundColor, CodeSetForegroundColor, CodeSetUnderlineColor:
		return fmt.Sprintf("%s(%s)", code, commandToColor(c))
	}

	return code.String()
}

// IsZero returns true if command is zero-initialized that actually also means the command is invalid.
// Zero commands are excluded from sequence during its rendition.
func (c Command) IsZero() bool {
	return c == 0
}

// Code returns the code of the command leaving aside its possible arguments.
func (c Command) Code() CommandCode {
	return CommandCode(c & 0xFF)
}

func (c Command) render(buf []byte) []byte {
	if !c.valid() {
		return buf
	}

	buf = append(buf, btoa(byte(c.Code()))...)

	for i := 0; i != c.argCount(); i++ {
		buf = append(buf, ';')
		buf = append(buf, btoa(c.arg(i))...)
	}

	return buf
}

func (c Command) valid() bool {
	return (c & commandValid) != 0
}

func (c Command) argCount() int {
	return int((c & commandArgCountMask) >> commandArgCountShift)
}

func (c Command) arg(i int) uint8 {
	return uint8((c >> ((i + 1) << 3)) & 0xFF)
}

// ---

// Complete set of supported CommandCode values.
const (
	CodeResetAll                        CommandCode = 0
	CodeSetBold                         CommandCode = 1
	CodeSetFaint                        CommandCode = 2
	CodeSetItalic                       CommandCode = 3
	CodeSetUnderlined                   CommandCode = 4
	CodeSetSlowBlink                    CommandCode = 5
	CodeSetRapidBlink                   CommandCode = 6
	CodeSetReversed                     CommandCode = 7
	CodeSetConcealed                    CommandCode = 8
	CodeSetCrossedOut                   CommandCode = 9
	CodeSetDoublyUnderlined             CommandCode = 21
	CodeResetBoldAndFaint               CommandCode = 22
	CodeResetItalic                     CommandCode = 23
	CodeResetAllUnderlines              CommandCode = 24
	CodeResetAllBlinks                  CommandCode = 25
	CodeResetReversed                   CommandCode = 27
	CodeResetConcealed                  CommandCode = 28
	CodeResetCrossedOut                 CommandCode = 29
	CodeSetForegroundColorBlack         CommandCode = 30
	CodeSetForegroundColorRed           CommandCode = 31
	CodeSetForegroundColorGreen         CommandCode = 32
	CodeSetForegroundColorYellow        CommandCode = 33
	CodeSetForegroundColorBlue          CommandCode = 34
	CodeSetForegroundColorMagenta       CommandCode = 35
	CodeSetForegroundColorCyan          CommandCode = 36
	CodeSetForegroundColorWhite         CommandCode = 37
	CodeSetForegroundColor              CommandCode = 38
	CodeResetForegroundColor            CommandCode = 39
	CodeSetBackgroundColorBlack         CommandCode = 40
	CodeSetBackgroundColorRed           CommandCode = 41
	CodeSetBackgroundColorGreen         CommandCode = 42
	CodeSetBackgroundColorYellow        CommandCode = 43
	CodeSetBackgroundColorBlue          CommandCode = 44
	CodeSetBackgroundColorMagenta       CommandCode = 45
	CodeSetBackgroundColorCyan          CommandCode = 46
	CodeSetBackgroundColorWhite         CommandCode = 47
	CodeSetBackgroundColor              CommandCode = 48
	CodeResetBackgroundColor            CommandCode = 49
	CodeSetFramed                       CommandCode = 51
	CodeSetEncircled                    CommandCode = 52
	CodeSetOverlined                    CommandCode = 53
	CodeResetFramedAndEncircled         CommandCode = 54
	CodeResetOverlined                  CommandCode = 55
	CodeSetUnderlineColor               CommandCode = 58
	CodeResetUnderlineColor             CommandCode = 59
	CodeSetSuperscript                  CommandCode = 73
	CodeSetSubscript                    CommandCode = 74
	CodeResetSuperscriptAndSubscript    CommandCode = 75
	CodeSetForegroundColorBrightBlack   CommandCode = 90
	CodeSetForegroundColorBrightRed     CommandCode = 91
	CodeSetForegroundColorBrightGreen   CommandCode = 92
	CodeSetForegroundColorBrightYellow  CommandCode = 93
	CodeSetForegroundColorBrightBlue    CommandCode = 94
	CodeSetForegroundColorBrightMagenta CommandCode = 95
	CodeSetForegroundColorBrightCyan    CommandCode = 96
	CodeSetForegroundColorBrightWhite   CommandCode = 97
	CodeSetBackgroundColorBrightBlack   CommandCode = 100
	CodeSetBackgroundColorBrightRed     CommandCode = 101
	CodeSetBackgroundColorBrightGreen   CommandCode = 102
	CodeSetBackgroundColorBrightYellow  CommandCode = 103
	CodeSetBackgroundColorBrightBlue    CommandCode = 104
	CodeSetBackgroundColorBrightMagenta CommandCode = 105
	CodeSetBackgroundColorBrightCyan    CommandCode = 106
	CodeSetBackgroundColorBrightWhite   CommandCode = 107
)

// CommandCode is the value of first octet of an SGR command.
type CommandCode uint8

// String returns textual description of c that can be used for debugging or logging purposes.
func (c CommandCode) String() string {
	if name, ok := commandNames[c]; ok {
		return name
	}

	return fmt.Sprintf("<!%d>", c)
}

// ---

var commandNames = map[CommandCode]string{
	CodeSetBold:                         "SetBold",
	CodeSetFaint:                        "SetFaint",
	CodeSetItalic:                       "SetItalic",
	CodeSetUnderlined:                   "SetUnderlined",
	CodeSetSlowBlink:                    "SetSlowBlink",
	CodeSetRapidBlink:                   "SetRapidBlink",
	CodeSetReversed:                     "SetReversed",
	CodeSetConcealed:                    "SetConcealed",
	CodeSetCrossedOut:                   "SetCrossedOut",
	CodeSetDoublyUnderlined:             "SetDoublyUnderlined",
	CodeResetBoldAndFaint:               "ResetBoldAndFaint",
	CodeResetItalic:                     "ResetItalic",
	CodeResetAllUnderlines:              "ResetAllUnderlines",
	CodeResetAllBlinks:                  "ResetAllBlinks",
	CodeResetReversed:                   "ResetReversed",
	CodeResetConcealed:                  "ResetConcealed",
	CodeResetCrossedOut:                 "ResetCrossedOut",
	CodeSetForegroundColorBlack:         "SetForegroundColorBlack",
	CodeSetForegroundColorRed:           "SetForegroundColorRed",
	CodeSetForegroundColorGreen:         "SetForegroundColorGreen",
	CodeSetForegroundColorYellow:        "SetForegroundColorYellow",
	CodeSetForegroundColorBlue:          "SetForegroundColorBlue",
	CodeSetForegroundColorMagenta:       "SetForegroundColorMagenta",
	CodeSetForegroundColorCyan:          "SetForegroundColorCyan",
	CodeSetForegroundColorWhite:         "SetForegroundColorWhite",
	CodeSetForegroundColor:              "SetForegroundColor",
	CodeResetForegroundColor:            "ResetForegroundColor",
	CodeSetBackgroundColorBlack:         "SetBackgroundColorBlack",
	CodeSetBackgroundColorRed:           "SetBackgroundColorRed",
	CodeSetBackgroundColorGreen:         "SetBackgroundColorGreen",
	CodeSetBackgroundColorYellow:        "SetBackgroundColorYellow",
	CodeSetBackgroundColorBlue:          "SetBackgroundColorBlue",
	CodeSetBackgroundColorMagenta:       "SetBackgroundColorMagenta",
	CodeSetBackgroundColorCyan:          "SetBackgroundColorCyan",
	CodeSetBackgroundColorWhite:         "SetBackgroundColorWhite",
	CodeSetBackgroundColor:              "SetBackgroundColor",
	CodeResetBackgroundColor:            "ResetBackgroundColor",
	CodeSetFramed:                       "SetFramed",
	CodeSetEncircled:                    "SetEncircled",
	CodeSetOverlined:                    "SetOverlined",
	CodeSetUnderlineColor:               "SetUnderlineColor",
	CodeResetFramedAndEncircled:         "ResetFramedAndEncircled",
	CodeResetOverlined:                  "ResetOverlined",
	CodeResetUnderlineColor:             "ResetUnderlineColor",
	CodeSetSuperscript:                  "SetSuperscript",
	CodeSetSubscript:                    "SetSubscript",
	CodeResetSuperscriptAndSubscript:    "ResetSuperscriptAndSubscript",
	CodeSetForegroundColorBrightBlack:   "SetForegroundColorBrightBlack",
	CodeSetForegroundColorBrightRed:     "SetForegroundColorBrightRed",
	CodeSetForegroundColorBrightGreen:   "SetForegroundColorBrightGreen",
	CodeSetForegroundColorBrightYellow:  "SetForegroundColorBrightYellow",
	CodeSetForegroundColorBrightBlue:    "SetForegroundColorBrightBlue",
	CodeSetForegroundColorBrightMagenta: "SetForegroundColorBrightMagenta",
	CodeSetForegroundColorBrightCyan:    "SetForegroundColorBrightCyan",
	CodeSetForegroundColorBrightWhite:   "SetForegroundColorBrightWhite",
	CodeSetBackgroundColorBrightBlack:   "SetBackgroundColorBrightBlack",
	CodeSetBackgroundColorBrightRed:     "SetBackgroundColorBrightRed",
	CodeSetBackgroundColorBrightGreen:   "SetBackgroundColorBrightGreen",
	CodeSetBackgroundColorBrightYellow:  "SetBackgroundColorBrightYellow",
	CodeSetBackgroundColorBrightBlue:    "SetBackgroundColorBrightBlue",
	CodeSetBackgroundColorBrightMagenta: "SetBackgroundColorBrightMagenta",
	CodeSetBackgroundColorBrightCyan:    "SetBackgroundColorBrightCyan",
	CodeSetBackgroundColorBrightWhite:   "SetBackgroundColorBrightWhite",
}

// ---

func setForegroundColor(color Color) Command {
	switch color.kind() {
	case colorKindDefault:
		return ResetForegroundColor
	case colorKindBasic:
		{
			bc := color.AsBasicColor()
			if bc >= BrightBlack {
				return SetForegroundColorBrightBlack + Command(bc-BrightBlack)
			}

			return SetForegroundColorBlack + Command(bc)
		}
	case colorKindPalette:
		return paletteColorToCommand(color.AsPaletteColor(), CodeSetForegroundColor)
	case colorKindRGB:
		return rgbColorToCommand(color.AsRGBColor(), CodeSetForegroundColor)
	default:
		return 0
	}
}

func setBackgroundColor(color Color) Command {
	switch color.kind() {
	case colorKindDefault:
		return ResetBackgroundColor
	case colorKindBasic:
		{
			bc := color.AsBasicColor()
			if bc >= BrightBlack {
				return SetBackgroundColorBrightBlack + Command(bc-BrightBlack)
			}

			return SetBackgroundColorBlack + Command(bc)
		}
	case colorKindPalette:
		return paletteColorToCommand(color.AsPaletteColor(), CodeSetBackgroundColor)
	case colorKindRGB:
		return rgbColorToCommand(color.AsRGBColor(), CodeSetBackgroundColor)
	default:
		return 0
	}
}

func setUnderlineColor(color Color) Command {
	switch color.kind() {
	case colorKindDefault:
		return ResetUnderlineColor
	case colorKindBasic:
		return paletteColorToCommand(color.AsBasicColor().PaletteColor(), CodeSetUnderlineColor)
	case colorKindPalette:
		return paletteColorToCommand(color.AsPaletteColor(), CodeSetUnderlineColor)
	case colorKindRGB:
		return rgbColorToCommand(color.AsRGBColor(), CodeSetUnderlineColor)
	default:
		return 0
	}
}

// ---

func paletteColorToCommand(color PaletteColor, code CommandCode) Command {
	return Command(code) | commandValid | commandArgCount2 | cmdArg(0, 5) | cmdArg(1, uint8(color))
}

func rgbColorToCommand(color RGBColor, code CommandCode) Command {
	return Command(code) | commandValid | commandArgCount4 | cmdArg(0, 2) | cmdArg(1, color.R()) | cmdArg(2, color.G()) | cmdArg(3, color.B())
}

func commandToColor(command Command) Color {
	if command.valid() {
		switch command.Code() {
		case CodeSetForegroundColor, CodeSetBackgroundColor, CodeSetUnderlineColor:
			switch command.arg(0) {
			case 5:
				return PaletteColor(command.arg(1)).Color()
			case 2:
				return RGB(command.arg(1), command.arg(2), command.arg(3)).Color()
			}
		}
	}

	return Color(0)
}

// ---

func cmdArg(i int, val uint8) Command {
	return Command(val) << ((i + 1) << 3)
}

// ---

const (
	commandValid         = 0x1000000000000000
	commandArgCountMask  = 0x0F00000000000000
	commandArgCount1     = 0x0100000000000000
	commandArgCount2     = 0x0200000000000000
	commandArgCount3     = 0x0300000000000000
	commandArgCount4     = 0x0400000000000000
	commandArgCountShift = 64 - 8
)
