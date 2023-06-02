package sgr

import "fmt"

// ---

// ErrInvalidColorText is an error that occurs in case of parsing an invalid textual color representation.
type ErrInvalidColorText struct {
	Value   string
	details error
}

// Error returns the error message.
func (e ErrInvalidColorText) Error() string {
	return fmt.Sprintf("invalid color text %q", e.Value)
}

// Unwrap returns the underlying error.
func (e ErrInvalidColorText) Unwrap() error {
	return e.details
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidColorText) Is(err error) bool {
	if other, ok := err.(ErrInvalidColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidColorValue is an error that occurs in case explicit validation or marshaling discovers an invalid value.
type ErrInvalidColorValue struct {
	Value Color
}

// Error returns the error message.
func (e ErrInvalidColorValue) Error() string {
	return fmt.Sprintf("invalid basic color value %d", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidColorValue) Is(err error) bool {
	if other, ok := err.(ErrInvalidColorValue); ok {
		return other.Value == 0 || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidBasicColorValue is an error that occurs in case explicit validation or marshaling discovers an invalid value.
type ErrInvalidBasicColorValue struct {
	Value BasicColor
}

// Error returns the error message.
func (e ErrInvalidBasicColorValue) Error() string {
	return fmt.Sprintf("invalid basic color value %d", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidBasicColorValue) Is(err error) bool {
	if other, ok := err.(ErrInvalidBasicColorValue); ok {
		return other.Value == 0 || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidBasicColorText is an error that occurs in case of parsing an invalid textual representation of BasicColor.
type ErrInvalidBasicColorText struct {
	Value string
}

// Error returns the error message.
func (e ErrInvalidBasicColorText) Error() string {
	return fmt.Sprintf("invalid basic color text %q", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidBasicColorText) Is(err error) bool {
	if other, ok := err.(ErrInvalidBasicColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	if other, ok := err.(ErrInvalidColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidDefaultColorText is an error that occurs in case of parsing an invalid textual representation of DefaultColor.
type ErrInvalidDefaultColorText struct {
	Value string
}

// Error returns the error message.
func (e ErrInvalidDefaultColorText) Error() string {
	return fmt.Sprintf("invalid default color text %q", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidDefaultColorText) Is(err error) bool {
	if other, ok := err.(ErrInvalidDefaultColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	if other, ok := err.(ErrInvalidColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidBrightnessValue is an error that occurs in case explicit validation or marshaling discovers an invalid value.
type ErrInvalidBrightnessValue struct {
	Value Brightness
}

// Error returns the error message.
func (e ErrInvalidBrightnessValue) Error() string {
	return fmt.Sprintf("invalid basic color value %d", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidBrightnessValue) Is(err error) bool {
	if other, ok := err.(ErrInvalidBrightnessValue); ok {
		return other.Value == 0 || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidBrightnessText is an error that occurs in case of parsing an invalid textual representation of Brightness.
type ErrInvalidBrightnessText struct {
	Value string
}

// Error returns the error message.
func (e ErrInvalidBrightnessText) Error() string {
	return fmt.Sprintf("invalid brightness text %q", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidBrightnessText) Is(err error) bool {
	if other, ok := err.(ErrInvalidBrightnessText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidPaletteColorText is an error that occurs in case of parsing an invalid textual representation of PaletteColor.
type ErrInvalidPaletteColorText struct {
	Value string
}

// Error returns the error message.
func (e ErrInvalidPaletteColorText) Error() string {
	return fmt.Sprintf("invalid palette color text %q", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidPaletteColorText) Is(err error) bool {
	if other, ok := err.(ErrInvalidPaletteColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	if other, ok := err.(ErrInvalidColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	return false
}

// ---

// ErrInvalidRGBColorText is an error that occurs in case of parsing an invalid textual representation of RGBColor.
type ErrInvalidRGBColorText struct {
	Value string
}

// Error returns the error message.
func (e ErrInvalidRGBColorText) Error() string {
	return fmt.Sprintf("invalid rgb color text %q", e.Value)
}

// Is returns true if e is a sub-class of err.
func (e ErrInvalidRGBColorText) Is(err error) bool {
	if other, ok := err.(ErrInvalidRGBColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	if other, ok := err.(ErrInvalidColorText); ok {
		return other.Value == "" || other.Value == e.Value
	}

	return false
}
