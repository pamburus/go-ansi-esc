package sgr_test

import (
	"errors"
	"testing"

	. "github.com/pamburus/go-tst/tst"

	"github.com/pamburus/go-ansi-esc/sgr"
)

func TestError(tt *testing.T) {
	t := New(tt)
	t.Expect(sgr.ErrInvalidColorText{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidDefaultColorText{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidBasicColorText{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidPaletteColorText{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidRGBColorText{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidBasicColorValue{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidColorValue{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidBrightnessValue{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidBrightnessText{}.Error()).ToNotEqual("")
	t.Expect(sgr.ErrInvalidColorValue{}).To(MatchError(sgr.ErrInvalidColorValue{}))
	t.Expect(sgr.ErrInvalidBrightnessValue{}).To(MatchError(sgr.ErrInvalidBrightnessValue{}))
	t.Expect(sgr.ErrInvalidBrightnessValue{sgr.Bright}).To(MatchError(sgr.ErrInvalidBrightnessValue{}))
	t.Expect(sgr.ErrInvalidBrightnessValue{}).ToNot(MatchError(sgr.ErrInvalidColorValue{}))
	t.Expect(sgr.ErrInvalidColorValue{sgr.Color(123435345)}).To(MatchError(sgr.ErrInvalidColorValue{}))
	t.Expect(sgr.ErrInvalidColorText{}).To(MatchError(sgr.ErrInvalidColorText{}))
	t.Expect(sgr.ErrInvalidBrightnessText{}).To(MatchError(sgr.ErrInvalidBrightnessText{}))
	t.Expect(sgr.ErrInvalidBrightnessText{"text"}).To(MatchError(sgr.ErrInvalidBrightnessText{}))
	t.Expect(sgr.ErrInvalidDefaultColorText{"text"}).To(MatchError(sgr.ErrInvalidDefaultColorText{}))
	t.Expect(sgr.ErrInvalidDefaultColorText{}).To(MatchError(sgr.ErrInvalidColorText{}))
	t.Expect(sgr.ErrInvalidBasicColorText{"text"}).To(MatchError(sgr.ErrInvalidBasicColorText{}))
	t.Expect(sgr.ErrInvalidBasicColorText{}).To(MatchError(sgr.ErrInvalidBasicColorText{}))
	t.Expect(sgr.ErrInvalidBasicColorText{}).To(MatchError(sgr.ErrInvalidColorText{}))
	t.Expect(sgr.ErrInvalidPaletteColorText{"text"}).To(MatchError(sgr.ErrInvalidPaletteColorText{}))
	t.Expect(sgr.ErrInvalidPaletteColorText{}).To(MatchError(sgr.ErrInvalidPaletteColorText{}))
	t.Expect(sgr.ErrInvalidPaletteColorText{}).To(MatchError(sgr.ErrInvalidColorText{}))
	t.Expect(sgr.ErrInvalidRGBColorText{"text"}).To(MatchError(sgr.ErrInvalidRGBColorText{}))
	t.Expect(sgr.ErrInvalidRGBColorText{}).To(MatchError(sgr.ErrInvalidRGBColorText{}))
	t.Expect(sgr.ErrInvalidRGBColorText{}).To(MatchError(sgr.ErrInvalidColorText{}))

	t.Expect(errors.Is(sgr.ErrInvalidColorText{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidBasicColorText{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidBrightnessText{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidDefaultColorText{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidPaletteColorText{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidRGBColorText{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidBasicColorValue{}, errors.New("some"))).ToEqual(false)
	t.Expect(errors.Is(sgr.ErrInvalidColorValue{}, errors.New("some"))).ToEqual(false)
}
