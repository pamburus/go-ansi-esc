package sgr_test

import (
	"testing"

	. "github.com/pamburus/go-tst/tst"

	"github.com/pamburus/go-ansi-esc/sgr"
)

func TestMode(tt *testing.T) {
	t := New(tt)

	t.Expect(sgr.NewModeSet().IsEmpty()).ToBeTrue()
	t.Expect(sgr.NewModeSet().IsZero()).ToBeTrue()
	t.Expect(sgr.EmptyModeSet().IsEmpty()).ToBeTrue()
	t.Expect(
		sgr.ModeSetWith(sgr.Overlined, sgr.Italic),
	).ToEqual(
		sgr.Overlined.ModeSet() | sgr.Italic.ModeSet(),
	)
	t.Expect(
		sgr.ModeSetWith(sgr.Overlined, sgr.Italic).ModeList(),
	).ToEqual(
		sgr.ModeList{sgr.Italic, sgr.Overlined},
	)
	t.Expect(
		sgr.NewModeList(sgr.Overlined, sgr.Italic),
	).ToEqual(
		sgr.ModeList{sgr.Overlined, sgr.Italic},
	)
	t.Expect(sgr.Italic.ModeList()).ToEqual(sgr.ModeList{sgr.Italic})

	t.Expect(sgr.DoublyUnderlined.String()).ToEqual("DoublyUnderlined")
	t.Expect(sgr.Mode(0x67).String()).ToEqual("<!0x67>")
	t.Expect(sgr.Mode(0x67).Validate()).To(MatchError(sgr.ErrInvalidModeValue{}))
	t.Expect(sgr.DoublyUnderlined.Validate()).ToSucceed()
	t.Expect(sgr.RapidBlink.MarshalText()).ToSucceed().AndResult().ToEqual([]byte("RapidBlink"))
	t.Expect(sgr.Mode(0x67).MarshalText()).ToFailWith(sgr.ErrInvalidModeValue{})

	var v sgr.Mode
	t.Expect(v.UnmarshalText([]byte("RapidBlink"))).ToSucceed()
	t.Expect(v).ToEqual(sgr.RapidBlink)
	t.Expect(v.UnmarshalText([]byte("rapid blink"))).ToSucceed()
	t.Expect(v).ToEqual(sgr.RapidBlink)
	t.Expect(v.UnmarshalText([]byte("rapid-blink"))).ToSucceed()
	t.Expect(v).ToEqual(sgr.RapidBlink)
	t.Expect(v.UnmarshalText([]byte("invalid"))).ToFailWith(sgr.ErrInvalidModeText{})

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

	t.Expect(set.Diff(other)).ToEqual(sgr.ModeSetDiff{Old: set, New: other})
	t.Expect(set.Diff(other).Reversed()).ToEqual(sgr.ModeSetDiff{Old: other, New: set})
	t.Expect(set.Diff(other).Added()).ToEqual(sgr.Subscript.ModeSet())
	t.Expect(set.Diff(other).Removed()).ToEqual(sgr.Overlined.ModeSet())
	t.Expect(set.Diff(other).Changed()).ToEqual(sgr.Overlined.ModeSet().With(sgr.Subscript))

	t.Expect(set.Diff(set).ToCommands(nil)).ToEqual(sgr.Sequence(nil))
}
