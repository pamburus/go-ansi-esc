package sgr_test

import (
	"testing"

	. "github.com/pamburus/go-tst/tst"

	"github.com/pamburus/go-ansi-esc/sgr"
)

func TestMode(tt *testing.T) {
	t := New(tt)
	set := sgr.NewModeSet().
		With(sgr.Bold).
		With(sgr.Italic).
		With(sgr.Overlined)
	t.Expect(set.Has(sgr.Bold)).ToBeTrue()
	t.Expect(set.Has(sgr.Italic)).ToBeTrue()
	t.Expect(set.Has(sgr.Overlined)).ToBeTrue()
	t.Expect(set.Has(sgr.Concealed)).ToBeFalse()
	t.Expect(set.Has(sgr.CrossedOut)).ToBeFalse()
	t.Expect(set.Has(sgr.Subscript)).ToBeFalse()

	t.Expect(set.Without(sgr.Italic)).ToEqual(
		sgr.NewModeSet().
			With(sgr.Bold).
			With(sgr.Overlined),
	)
	t.Expect(set.WithToggled(sgr.Italic)).ToEqual(
		sgr.NewModeSet().
			With(sgr.Bold).
			With(sgr.Overlined),
	)
	t.Expect(set.WithToggled(sgr.CrossedOut)).ToEqual(
		sgr.NewModeSet().
			With(sgr.Bold).
			With(sgr.Overlined).
			With(sgr.Italic).
			With(sgr.CrossedOut),
	)

	other := sgr.NewModeSet().
		With(sgr.Bold).
		With(sgr.Italic).
		With(sgr.Subscript)
	t.Expect(set.WithOther(other, sgr.ModeReplace)).ToEqual(other)
	t.Expect(set.WithOther(other, sgr.ModeAdd)).ToEqual(
		sgr.NewModeSet().
			With(sgr.Bold).
			With(sgr.Italic).
			With(sgr.Overlined).
			With(sgr.Subscript),
	)
	t.Expect(set.WithOther(other, sgr.ModeRemove)).ToEqual(
		sgr.NewModeSet().
			With(sgr.Overlined),
	)
	t.Expect(set.WithOther(other, sgr.ModeToggle)).ToEqual(
		sgr.NewModeSet().
			With(sgr.Overlined).
			With(sgr.Subscript),
	)
	t.Expect(set.WithOther(other, sgr.ModeAction(345))).ToEqual(set)
}
