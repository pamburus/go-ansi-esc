package sgr

import (
	"encoding"
	"fmt"
	"strings"
)

// ---

// Complete set of valid BasicColor values.
const (
	Black         BasicColor = 0
	Red           BasicColor = 1
	Green         BasicColor = 2
	Yellow        BasicColor = 3
	Blue          BasicColor = 4
	Magenta       BasicColor = 5
	Cyan          BasicColor = 6
	White         BasicColor = 7
	BrightBlack   BasicColor = 8
	BrightRed     BasicColor = 9
	BrightGreen   BasicColor = 10
	BrightYellow  BasicColor = 11
	BrightBlue    BasicColor = 12
	BrightMagenta BasicColor = 13
	BrightCyan    BasicColor = 14
	BrightWhite   BasicColor = 15
)

// BasicColor is a value from the most widely supported color set that was originally defined in ECMA-48 standard.
type BasicColor uint8

// Color converts c to a uniform Color value that can be used in all functions dealing with colors.
func (c BasicColor) Color() Color {
	return Color(c) | colorKindBasic
}

// Normal returns a counter-part color that has Normal brightness
// or just c if it already has Normal brightness.
func (c BasicColor) Normal() BasicColor {
	if c >= 8 {
		c -= 8
	}

	return c
}

// Bright returns a counter-part color that has Bright brightness
// or just c if it already has Bright brightness.
func (c BasicColor) Bright() BasicColor {
	if c < 8 {
		c += 8
	}

	return c
}

// WithBrightness ensures c has requested brightness and returns it
// or returns a counter-part color that has the requested brightness.
func (c BasicColor) WithBrightness(b Brightness) BasicColor {
	if b == Bright {
		return c.Bright()
	}

	return c.Normal()
}

// PaletteColor converts c to PaletteColor.
// Palette colors are backward-compatible with basic color because
// its first 16 colors have the same codes and meaning.
func (c BasicColor) PaletteColor() PaletteColor {
	return PaletteColor(c)
}

// String returns textual description of c that can be used for debugging or logging purposes.
func (c BasicColor) String() string {
	if name, ok := basicColorNames[c]; ok {
		return name
	}

	return fmt.Sprintf("<!0x%02x>", uint8(c))
}

// Validate check that c has a valid value that is one of the values defined by the ECMA-48 standard.
func (c BasicColor) Validate() error {
	if _, ok := basicColorNames[c]; !ok {
		return ErrInvalidBasicColorValue{c}
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler interface
// that allows BasicColor to be used in any compatible marshaler like JSON, YAML, etc.
func (c BasicColor) MarshalText() ([]byte, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	return []byte(c.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows BasicColor to be used in any compatible unmarshaler like JSON, YAML, etc.
func (c *BasicColor) UnmarshalText(data []byte) error {
	return c.unmarshalText(string(data))
}

func (c *BasicColor) unmarshalText(text string) error {
	brightness := Normal
	t := text

	if len(t) > len(textBright) && strings.EqualFold(t[:len(textBright)], textBright) {
		t = strings.TrimLeft(t[len(textBright):], " _-")
		brightness = Bright
	}

	colors := []BasicColor{
		Black,
		Red,
		Green,
		Yellow,
		Blue,
		Magenta,
		Cyan,
		White,
	}

	for _, color := range colors {
		if strings.EqualFold(t, color.String()) {
			*c = color.WithBrightness(brightness)

			return nil
		}
	}

	return ErrInvalidBasicColorText{text}
}

// ---

// Complete set of valid Brightness values .
const (
	Normal Brightness = iota
	Bright
)

// Brightness defines brightness of a BasicColor value.
type Brightness uint8

// String returns textual description of b that can be used for debugging or logging purposes.
func (b Brightness) String() string {
	switch b {
	case Normal:
		return textNormal
	case Bright:
		return textBright
	default:
		return fmt.Sprintf("<!0x%02x>", uint8(b))
	}
}

// Validate check that b has a valid value that is one of the valid values.
func (b Brightness) Validate() error {
	switch b {
	case Normal, Bright:
		return nil
	default:
		return ErrInvalidBrightnessValue{b}
	}
}

// MarshalText implements encoding.TextMarshaler interface
// that allows Brightness to be used in any compatible marshaler like JSON, YAML, etc.
func (b Brightness) MarshalText() ([]byte, error) {
	err := b.Validate()
	if err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows Brightness to be used in any compatible unmarshaler like JSON, YAML, etc.
func (b *Brightness) UnmarshalText(text []byte) error {
	s := string(text)
	switch {
	case strings.EqualFold(s, textBright):
		*b = Bright
	case strings.EqualFold(s, textNormal):
		*b = Normal
	default:
		return ErrInvalidBrightnessText{s}
	}

	return nil
}

// ---

// Default is a value of DefaultColor type.
// It can be used to reset color to the default color defined by the terminal.
var Default DefaultColor

// DefaultColor is the default color defined by the terminal.
type DefaultColor struct{}

// Color returns a uniform Color value having the same meaning that can be used in all functions dealing with colors.
func (DefaultColor) Color() Color {
	return Color(colorKindDefault)
}

// String returns textual description of c that can be used for debugging or logging purposes.
func (DefaultColor) String() string {
	return "Default"
}

// Validate always returns true.
func (DefaultColor) Validate() error {
	return nil
}

// MarshalText implements encoding.TextMarshaler interface
// that allows DefaultColor to be used in any compatible marshaler like JSON, YAML, etc.
func (DefaultColor) MarshalText() ([]byte, error) {
	return []byte(Default.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows DefaultColor to be used in any compatible unmarshaler like JSON, YAML, etc.
func (DefaultColor) UnmarshalText(data []byte) error {
	return Default.unmarshalText(string(data))
}

func (DefaultColor) unmarshalText(text string) error {
	if strings.EqualFold(text, Default.String()) {
		return nil
	}

	return ErrInvalidDefaultColorText{text}
}

// ---

// PaletteColor is a 256-color palette color code.
// First 16 colors a compatible with BasicColor codes.
type PaletteColor uint8

// Color converts c to a uniform Color value having the same meaning that can be used in all functions dealing with colors.
func (c PaletteColor) Color() Color {
	return Color(c) | colorKindPalette
}

// String returns textual description of c that can be used for debugging or logging purposes.
func (c PaletteColor) String() string {
	return fmt.Sprintf("#%02x", uint8(c))
}

// Validate always returns true.
func (PaletteColor) Validate() error {
	return nil
}

// MarshalText implements encoding.TextMarshaler interface
// that allows PaletteColor to be used in any compatible marshaler like JSON, YAML, etc.
func (c PaletteColor) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows PaletteColor to be used in any compatible unmarshaler like JSON, YAML, etc.
func (c *PaletteColor) UnmarshalText(data []byte) error {
	return c.unmarshalText(string(data))
}

func (c *PaletteColor) unmarshalText(text string) error {
	text = strings.TrimSpace(text)
	var v PaletteColor

	if len(text) != 3 {
		return ErrInvalidPaletteColorText{text}
	}

	n, err := fmt.Sscanf(text, "#%02x\n", &v)
	if err != nil || n != 1 {
		return ErrInvalidPaletteColorText{text}
	}

	*c = v

	return nil
}

// ---

// RGB constructs a RGBColor with the given red, green and blue values.
func RGB(red, green, blue uint8) RGBColor {
	return RGBColor(uint32(red)<<16 | uint32(green)<<8 | uint32(blue))
}

// RGBColor is a color represented by a combination of 8-bit brightness values of red, green and blue components.
type RGBColor uint32

// R returns red component brightness value.
func (c RGBColor) R() uint8 {
	return uint8((c >> 16) & 0xFF)
}

// G returns green component brightness value.
func (c RGBColor) G() uint8 {
	return uint8((c >> 8) & 0xFF)
}

// B returns blue component brightness value.
func (c RGBColor) B() uint8 {
	return uint8((c >> 0) & 0xFF)
}

// Color converts c to a uniform Color value having the same meaning that can be used in all functions dealing with colors.
func (c RGBColor) Color() Color {
	return Color(c) | colorKindRGB
}

// String returns textual description of c that can be used for debugging or logging purposes.
func (c RGBColor) String() string {
	return fmt.Sprintf("#%06x", uint32(c))
}

// Validate returns nil.
func (RGBColor) Validate() error {
	return nil
}

// MarshalText implements encoding.TextMarshaler interface
// that allows RGBColor to be used in any compatible marshaler like JSON, YAML, etc.
func (c RGBColor) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows RGBColor to be used in any compatible unmarshaler like JSON, YAML, etc.
func (c *RGBColor) UnmarshalText(data []byte) error {
	return c.unmarshalText(string(data))
}

func (c *RGBColor) unmarshalText(text string) error {
	text = strings.TrimSpace(text)
	var v RGBColor

	if len(text) != 7 {
		return ErrInvalidRGBColorText{text}
	}

	n, err := fmt.Sscanf(text, "#%06x\n", &v)
	if err != nil || n != 1 {
		return ErrInvalidRGBColorText{text}
	}

	*c = v

	return nil
}

// ---

// Color is a uniform color value that can represent a BasicColor, DefaultColor, PaletteColor or RGBColor.
type Color uint32

// Color implements IntoColor interface that is needed to avoid explicit type conversions.
func (c Color) Color() Color {
	return c
}

// IsZero returns true if c is zero-initialized value that can be treated as nil.
func (c Color) IsZero() bool {
	return c == 0
}

// OrDefault returns default color if c is zero otherwise it returns c.
func (c Color) OrDefault() Color {
	if c.IsZero() {
		return Default.Color()
	}

	return c
}

// IsBasicColor returns true if c represents a BasicColor value.
func (c Color) IsBasicColor() bool {
	return c.kind() == colorKindBasic
}

// BasicColor returns a BasicColor and true if c represents a BasicColor value or zero value and false otherwise.
func (c Color) BasicColor() (BasicColor, bool) {
	if c.IsBasicColor() {
		return c.AsBasicColor(), true
	}

	return 0, false
}

// AsBasicColor converts c to BasicColor or panics if c does not represent a BasicColor.
func (c Color) AsBasicColor() BasicColor {
	if !c.IsBasicColor() {
		if !c.IsPaletteColor() || (c&0xF0) != 0 {
			return BasicColor(0)
		}
	}

	return BasicColor(c & 0xFF)
}

// IsDefaultColor returns true if c represents a DefaultColor value.
func (c Color) IsDefaultColor() bool {
	return c.kind() == colorKindDefault
}

// DefaultColor returns a DefaultColor and true if c represents a DefaultColor value or zero value and false otherwise.
func (c Color) DefaultColor() (DefaultColor, bool) {
	if c.IsDefaultColor() {
		return Default, true
	}

	return Default, false
}

// AsDefaultColor converts c to DefaultColor or panics if c does not represent a DefaultColor.
func (c Color) AsDefaultColor() DefaultColor {
	return Default
}

// IsPaletteColor returns true if c represents a PaletteColor value.
func (c Color) IsPaletteColor() bool {
	return c.kind() == colorKindPalette
}

// PaletteColor returns a PaletteColor and true if c represents a PaletteColor value or zero value and false otherwise.
func (c Color) PaletteColor() (PaletteColor, bool) {
	if c.IsPaletteColor() {
		return c.AsPaletteColor(), true
	}

	return 0, false
}

// AsPaletteColor converts c to PaletteColor or panics if c does not represent a PaletteColor.
func (c Color) AsPaletteColor() PaletteColor {
	if !c.IsPaletteColor() && !c.IsBasicColor() {
		return PaletteColor(0)
	}

	return PaletteColor(c & 0xFF)
}

// IsRGBColor returns true if c represents a RGBColor value.
func (c Color) IsRGBColor() bool {
	return c.kind() == colorKindRGB
}

// RGBColor returns a RGBColor and true if c represents a RGBColor value or zero value and false otherwise.
func (c Color) RGBColor() (RGBColor, bool) {
	if c.IsRGBColor() {
		return c.AsRGBColor(), true
	}

	return 0, false
}

// AsRGBColor converts c to RGBColor or panics if c does not represent a RGBColor.
func (c Color) AsRGBColor() RGBColor {
	if !c.IsRGBColor() {
		return RGBColor(0)
	}

	return RGBColor(c & 0xFFFFFF)
}

// String returns textual description of c that can be used for debugging or logging purposes.
func (c Color) String() string {
	if c.IsZero() {
		return ""
	}

	switch c.kind() {
	case colorKindDefault:
		return c.AsDefaultColor().String()
	case colorKindBasic:
		return c.AsBasicColor().String()
	case colorKindPalette:
		return c.AsPaletteColor().String()
	case colorKindRGB:
		return c.AsRGBColor().String()
	default:
		return fmt.Sprintf("<!0x%08x>", uint32(c))
	}
}

// Validate check if c has a valid value.
// Zero value is not considered an error, checks for zero value should be done explicitly.
func (c Color) Validate() error {
	if c.IsZero() {
		return nil
	}

	switch c.kind() {
	case colorKindDefault:
		return c.AsDefaultColor().Validate()
	case colorKindBasic:
		return c.AsBasicColor().Validate()
	case colorKindPalette:
		return c.AsPaletteColor().Validate()
	case colorKindRGB:
		return c.AsRGBColor().Validate()
	default:
		return ErrInvalidColorValue{c}
	}
}

// MarshalText implements encoding.TextMarshaler interface
// that allows Color to be used in any compatible marshaler like JSON, YAML, etc.
func (c Color) MarshalText() ([]byte, error) {
	if c.IsZero() {
		return nil, nil
	}

	err := c.Validate()
	if err != nil {
		return nil, err
	}

	return []byte(c.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface
// that allows Color to be used in any compatible unmarshaler like JSON, YAML, etc.
func (c *Color) UnmarshalText(data []byte) error {
	text := strings.TrimSpace(string(data))

	switch {
	case len(text) == 0:
		*c = 0
	case len(text) == 3 && text[0] == '#':
		var v PaletteColor
		err := v.unmarshalText(text)
		if err != nil {
			return ErrInvalidColorText{text, err}
		}
		*c = Color(v) | colorKindPalette
	case len(text) == 7 && text[0] == '#':
		var v RGBColor
		err := v.unmarshalText(text)
		if err != nil {
			return ErrInvalidColorText{text, err}
		}
		*c = Color(v) | colorKindRGB
	case len(text) == 7:
		if Default.unmarshalText(text) == nil {
			*c = Default.Color()

			break
		}

		fallthrough
	default:
		var v BasicColor
		err := v.unmarshalText(text)
		if err != nil {
			return ErrInvalidColorText{text, nil}
		}
		*c = Color(v) | colorKindBasic
	}

	return nil
}

func (c Color) kind() Color {
	return c & colorKindMask
}

// ---

// IntoColor is an interface that can be used to uniformly accessing any color type.
type IntoColor interface {
	Color() Color
}

// ---

var basicColorNames = map[BasicColor]string{
	Black:         "Black",
	Red:           "Red",
	Green:         "Green",
	Yellow:        "Yellow",
	Blue:          "Blue",
	Magenta:       "Magenta",
	Cyan:          "Cyan",
	White:         "White",
	BrightBlack:   "BrightBlack",
	BrightRed:     "BrightRed",
	BrightGreen:   "BrightGreen",
	BrightYellow:  "BrightYellow",
	BrightBlue:    "BrightBlue",
	BrightMagenta: "BrightMagenta",
	BrightCyan:    "BrightCyan",
	BrightWhite:   "BrightWhite",
}

var textNormal = "Normal"
var textBright = "Bright"

// ---

const (
	colorKindMask    = 0xFF000000
	colorKindNone    = 0x00000000
	colorKindDefault = 0x01000000
	colorKindBasic   = 0x02000000
	colorKindPalette = 0x03000000
	colorKindRGB     = 0x04000000
)

// ---

var _ = []encoding.TextUnmarshaler{
	new(DefaultColor),
	new(BasicColor),
	new(PaletteColor),
	new(RGBColor),
	new(Color),
	Default,
}
