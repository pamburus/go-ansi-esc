package sgr_test

import (
	"testing"

	"github.com/pamburus/go-ansi-esc/sgr"
	. "github.com/pamburus/go-tst/tst"
)

func TestCommand(tt *testing.T) {
	t := New(tt)

	t.Run("Zero", func(t Test) {
		var zero sgr.Command
		t.Expect(zero.IsZero()).ToBeTrue()
		t.Expect(zero.String()).ToEqual("")
		t.Expect(sgr.ResetAll.IsZero()).ToBeFalse()
		t.Expect(sgr.SetBackgroundColor(sgr.Color(0))).ToEqual(zero)
		t.Expect(sgr.SetForegroundColor(sgr.Color(0))).ToEqual(zero)
		t.Expect(sgr.SetUnderlineColor(sgr.Color(0))).ToEqual(zero)
	})

	t.Run("Invalid", func(t Test) {
		t.Expect(sgr.Command(500).String()).ToEqual("<!0x000001f4>")
		t.Expect((sgr.ResetAll | sgr.Command(sgr.CodeSetBackgroundColor)).String()).ToEqual("SetBackgroundColor()")
	})

	t.Run("String", func(t Test) {
		t.Run("Valid", func(t Test) {
			t.Expect(
				sgr.SetItalic.String(),
			).To(Equal(
				"SetItalic",
			))
		})
		t.Run("Invalid", func(t Test) {
			t.Expect(
				sgr.CommandCode(108).String(),
			).To(Equal(
				"<!108>",
			))
		})
	})

	t.Run("SetBackgroundColor", func(t Test) {
		t.Expect(
			sgr.SetBackgroundColor(sgr.BrightCyan).String(),
		).To(Equal(
			"SetBackgroundColorBrightCyan",
		))
		t.Expect(
			sgr.SetBackgroundColor(sgr.Red).String(),
		).To(Equal(
			"SetBackgroundColorRed",
		))
	})
	t.Run("SetForegroundColor", func(t Test) {
		t.Expect(
			sgr.SetForegroundColor(sgr.RGB(32, 48, 64)).String(),
		).To(Equal(
			"SetForegroundColor(#203040)",
		))
	})
	t.Run("SetUnderlineColor", func(t Test) {
		t.Expect(
			sgr.SetUnderlineColor(sgr.Green).String(),
		).To(Equal(
			"SetUnderlineColor(#02)",
		))
		t.Expect(
			sgr.SetUnderlineColor(sgr.PaletteColor(56)).String(),
		).To(Equal(
			"SetUnderlineColor(#38)",
		))
		t.Expect(
			sgr.SetUnderlineColor(sgr.RGB(12, 35, 65)).String(),
		).To(Equal(
			"SetUnderlineColor(#0c2341)",
		))
		t.Expect(
			sgr.SetUnderlineColor(sgr.Default).String(),
		).To(Equal(
			"ResetUnderlineColor",
		))
	})
}
