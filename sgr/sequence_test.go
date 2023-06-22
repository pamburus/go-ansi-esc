package sgr_test

import (
	"testing"

	. "github.com/pamburus/go-tst/tst"

	"github.com/pamburus/go-ansi-esc/sgr"
)

func TestSequence(tt *testing.T) {
	t := New(tt)

	t.Run("Basic/Basic/Italic", func(t Test) {
		t.Expect(
			string(sgr.Sequence{
				sgr.SetBackgroundColor(sgr.Black.Bright()),
				sgr.SetForegroundColor(sgr.Cyan),
				sgr.SetItalic,
			}.Bytes()),
		).ToEqual(
			"\x1b[100;36;3m",
		)
	})

	t.Run("Default/Palette/Italic", func(t Test) {
		t.Expect(
			string(sgr.Sequence{
				sgr.SetBackgroundColor(sgr.Default),
				sgr.SetForegroundColor(sgr.PaletteColor(12)),
				sgr.SetItalic,
			}.Bytes()),
		).ToEqual(
			"\x1b[49;38;5;12;3m",
		)
	})

	t.Run("Palette/RGB/Bold", func(t Test) {
		t.Expect(
			string(sgr.Sequence{
				sgr.SetBackgroundColor(sgr.PaletteColor(12)),
				sgr.SetForegroundColor(sgr.RGB(128, 140, 192)),
				sgr.SetBold,
			}.Bytes()),
		).ToEqual(
			"\x1b[48;5;12;38;2;128;140;192;1m",
		)
	})

	t.Run("RGB/RGB/Bold", func(t Test) {
		t.Expect(
			string(sgr.Sequence{
				sgr.SetBackgroundColor(sgr.RGB(12, 14, 19).Color()),
				sgr.SetForegroundColor(sgr.RGB(128, 140, 192).Color()),
				sgr.SetBold,
			}.Bytes()),
		).ToEqual(
			"\x1b[48;2;12;14;19;38;2;128;140;192;1m",
		)
	})

	t.Run("Zero/Bold", func(t Test) {
		t.Expect(
			string(sgr.Sequence{
				sgr.SetBackgroundColor(sgr.Color(0)),
				sgr.SetBold,
			}.Bytes()),
		).ToEqual(
			"\x1b[1m",
		)
	})

	t.Run("UnderlineColor/Basic", func(t Test) {
		t.Expect(
			string(sgr.Sequence{
				sgr.SetUnderlineColor(sgr.Green),
			}.Bytes()),
		).ToEqual(
			"\x1b[58;5;2m",
		)
	})
}

func BenchmarkSequence(b *testing.B) {
	run := func(b *testing.B, seq sgr.Sequence) {
		buf := make([]byte, 0, 100)

		for i := 0; i < b.N; i++ {
			seq.Render(buf)
		}
	}

	b.Run("None/Basic/Italic", func(b *testing.B) {
		run(b, sgr.Sequence{
			sgr.SetForegroundColor(sgr.BrightWhite),
			sgr.SetItalic,
		})
	})

	b.Run("Basic/RGB/Italic", func(b *testing.B) {
		run(b, sgr.Sequence{
			sgr.SetBackgroundColor(sgr.Black),
			sgr.SetForegroundColor(sgr.RGB(10, 45, 34)),
			sgr.SetItalic,
		})
	})
}
