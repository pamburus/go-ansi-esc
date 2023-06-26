package sgr_test

import (
	"testing"

	. "github.com/pamburus/go-tst/tst"

	"github.com/pamburus/go-ansi-esc/sgr"
)

func TestColor(tt *testing.T) {
	t := New(tt)

	t.Run("Color", func(t Test) {
		t.Run("Zero", func(t Test) {
			var zero sgr.Color
			t.Expect(zero.Color()).ToEqual(zero)
			t.Expect(zero.IsZero()).ToEqual(true)
			t.Expect(zero.IsDefaultColor()).ToEqual(false)
			t.Expect(zero.IsBasicColor()).ToEqual(false)
			t.Expect(zero.IsPaletteColor()).ToEqual(false)
			t.Expect(zero.IsRGBColor()).ToEqual(false)
			t.Expect(zero.DefaultColor()).ToEqual(sgr.Default, false)
			t.Expect(zero.BasicColor()).ToEqual(sgr.Black, false)
			t.Expect(zero.PaletteColor()).ToEqual(sgr.PaletteColor(0), false)
			t.Expect(zero.RGBColor()).ToEqual(sgr.RGBColor(0), false)
			t.Expect(zero.AsDefaultColor()).ToEqual(sgr.Default)
			t.Expect(zero.AsBasicColor()).ToEqual(sgr.BasicColor(0))
			t.Expect(zero.AsPaletteColor()).ToEqual(sgr.PaletteColor(0))
			t.Expect(zero.AsRGBColor()).ToEqual(sgr.RGBColor(0))
			t.Expect(zero.String()).ToEqual("")
			t.Expect(zero.Validate()).ToSucceed()
			t.Expect(zero.MarshalText()).ToSucceed().AndResult().ToEqual([]byte(nil))
			t.Expect(zero.OrDefault()).ToEqual(sgr.Default.Color())
			var color sgr.Color
			t.Expect(color.UnmarshalText(nil)).ToSucceed()
			t.Expect(color).ToEqual(zero)
		})
		t.Run("Basic", func(t Test) {
			color := sgr.Yellow.Color()
			t.Expect(color.IsZero()).ToEqual(false)
			t.Expect(color.IsDefaultColor()).ToEqual(false)
			t.Expect(color.IsBasicColor()).ToEqual(true)
			t.Expect(color.AsBasicColor()).ToEqual(sgr.Yellow)
			t.Expect(color.IsPaletteColor()).ToEqual(false)
			t.Expect(color.IsRGBColor()).ToEqual(false)
			t.Expect(color.DefaultColor()).ToEqual(sgr.Default, false)
			t.Expect(color.BasicColor()).ToEqual(sgr.Yellow, true)
			t.Expect(color.PaletteColor()).ToEqual(sgr.PaletteColor(0), false)
			t.Expect(color.RGBColor()).ToEqual(sgr.RGBColor(0), false)
			t.Expect(color.String()).ToEqual("Yellow")
			t.Expect(color.Validate()).ToSucceed()
			t.Expect(color.MarshalText()).ToSucceed().AndResult().ToEqual([]byte("Yellow"))
			t.Expect(color.OrDefault()).ToEqual(color)
			var otherColor sgr.Color
			t.Expect(otherColor.UnmarshalText([]byte("Magenta"))).ToSucceed()
			t.Expect(otherColor).ToEqual(sgr.Magenta.Color())
			t.Run("Invalid", func(t Test) {
				color := sgr.BasicColor(127).Color()
				t.Expect(color.String()).ToEqual("<!0x7f>")
				t.Expect(color.Validate()).ToFailWith(sgr.ErrInvalidBasicColorValue{})
				t.Expect(color.MarshalText()).ToFailWith(sgr.ErrInvalidBasicColorValue{})
				var otherColor sgr.Color
				t.Expect(otherColor.UnmarshalText([]byte("<!0x7f>"))).ToFailWith(sgr.ErrInvalidColorText{})
			})
		})
		t.Run("Default", func(t Test) {
			color := sgr.Default.Color()
			t.Expect(color.IsZero()).ToEqual(false)
			t.Expect(color.IsDefaultColor()).ToEqual(true)
			t.Expect(color.AsDefaultColor()).ToEqual(sgr.Default)
			t.Expect(color.IsBasicColor()).ToEqual(false)
			t.Expect(color.IsPaletteColor()).ToEqual(false)
			t.Expect(color.IsRGBColor()).ToEqual(false)
			t.Expect(color.DefaultColor()).ToEqual(sgr.Default, true)
			t.Expect(color.BasicColor()).ToEqual(sgr.Black, false)
			t.Expect(color.PaletteColor()).ToEqual(sgr.PaletteColor(0), false)
			t.Expect(color.RGBColor()).ToEqual(sgr.RGBColor(0), false)
			t.Expect(color.String()).ToEqual("Default")
			t.Expect(color.Validate()).ToSucceed()
			t.Expect(color.MarshalText()).ToSucceed().AndResult().ToEqual([]byte("Default"))
			var otherColor sgr.Color
			t.Expect(otherColor.UnmarshalText([]byte("Default"))).ToSucceed()
			t.Expect(otherColor).ToEqual(color)
		})
		t.Run("Palette", func(t Test) {
			color := sgr.PaletteColor(255).Color()
			t.Expect(color.IsZero()).ToEqual(false)
			t.Expect(color.IsDefaultColor()).ToEqual(false)
			t.Expect(color.IsBasicColor()).ToEqual(false)
			t.Expect(color.IsPaletteColor()).ToEqual(true)
			t.Expect(color.AsPaletteColor()).ToEqual(sgr.PaletteColor(255))
			t.Expect(color.IsRGBColor()).ToEqual(false)
			t.Expect(color.DefaultColor()).ToEqual(sgr.Default, false)
			t.Expect(color.BasicColor()).ToEqual(sgr.Black, false)
			t.Expect(color.PaletteColor()).ToEqual(sgr.PaletteColor(255), true)
			t.Expect(color.RGBColor()).ToEqual(sgr.RGBColor(0), false)
			t.Expect(color.String()).ToEqual("#ff")
			t.Expect(color.Validate()).ToSucceed()
			t.Expect(color.MarshalText()).ToSucceed().AndResult().ToEqual([]byte("#ff"))
			t.Expect(color.OrDefault()).ToEqual(color)
			var otherColor sgr.Color
			t.Expect(otherColor.UnmarshalText([]byte("#ff"))).ToSucceed()
			t.Expect(otherColor).ToEqual(color)
			t.Run("Invalid", func(t Test) {
				var otherColor sgr.Color
				t.Expect(otherColor.UnmarshalText([]byte("#0q"))).ToFailWith(sgr.ErrInvalidPaletteColorText{})
			})
		})
		t.Run("RGB", func(t Test) {
			color := sgr.RGB(0, 128, 255).Color()
			t.Expect(color.IsZero()).ToEqual(false)
			t.Expect(color.IsDefaultColor()).ToEqual(false)
			t.Expect(color.IsBasicColor()).ToEqual(false)
			t.Expect(color.IsPaletteColor()).ToEqual(false)
			t.Expect(color.IsRGBColor()).ToEqual(true)
			t.Expect(color.AsRGBColor()).ToEqual(sgr.RGB(0, 128, 255))
			t.Expect(color.DefaultColor()).ToEqual(sgr.Default, false)
			t.Expect(color.BasicColor()).ToEqual(sgr.Black, false)
			t.Expect(color.PaletteColor()).ToEqual(sgr.PaletteColor(0), false)
			t.Expect(color.RGBColor()).ToEqual(sgr.RGB(0, 128, 255), true)
			t.Expect(color.String()).ToEqual("#0080ff")
			t.Expect(color.Validate()).ToSucceed()
			t.Expect(color.MarshalText()).ToSucceed().AndResult().ToEqual([]byte("#0080ff"))
			t.Expect(color.OrDefault()).ToEqual(color)
			var otherColor sgr.Color
			t.Expect(otherColor.UnmarshalText([]byte("#0080ff"))).ToSucceed()
			t.Expect(otherColor).ToEqual(color)
			t.Run("Invalid", func(t Test) {
				var otherColor sgr.Color
				t.Expect(otherColor.UnmarshalText([]byte("#0080fq"))).ToFailWith(sgr.ErrInvalidRGBColorText{})
			})
		})

		t.Run("Invalid", func(t Test) {
			t.Expect(sgr.Color(918239182).MarshalText()).ToFailWith(sgr.ErrInvalidColorValue{})
			t.Expect(sgr.Color(0x36bb37ce).String()).ToEqual("<!0x36bb37ce>")
		})
	})

	t.Run("Basic", func(t Test) {
		t.Expect(sgr.Black.Bright()).ToEqual(sgr.BrightBlack)
		t.Expect(sgr.BrightBlack.Normal()).ToEqual(sgr.Black)
		t.Expect(sgr.Black.Normal()).ToEqual(sgr.Black)
		t.Expect(sgr.BrightBlack.Bright()).ToEqual(sgr.BrightBlack)
		t.Expect(sgr.BrightBlack.WithBrightness(sgr.Normal)).ToEqual(sgr.Black)
		t.Expect(sgr.Black.WithBrightness(sgr.Bright)).ToEqual(sgr.BrightBlack)

		t.Run("String", func(t Test) {
			t.Run("Valid", func(t Test) {
				t.Expect(
					sgr.BrightGreen.String(),
				).ToEqual(
					"BrightGreen",
				)
			})
			t.Run("Invalid", func(t Test) {
				t.Expect(
					sgr.BasicColor(128).String(),
				).ToEqual(
					"<!0x80>",
				)
			})
		})

		t.Run("Validate", func(t Test) {
			t.Run("Valid", func(t Test) {
				t.Expect(sgr.BrightGreen.Validate()).ToSucceed()
			})
			t.Run("Invalid", func(t Test) {
				t.Expect(sgr.BasicColor(127).Validate()).ToFailWith(sgr.ErrInvalidBasicColorValue{127})
			})
		})

		t.Run("MarshalText", func(t Test) {
			t.Run("Valid", func(t Test) {
				text, err := sgr.BrightGreen.MarshalText()
				t.Expect(err).ToNot(HaveOccurred())
				t.Expect(text).ToEqual([]byte("BrightGreen"))
			})
			t.Run("Invalid", func(t Test) {
				t.Expect(sgr.BasicColor(127).MarshalText()).ToFailWith(sgr.ErrInvalidBasicColorValue{})
			})
		})

		t.Run("UnmarshalText", func(t Test) {
			t.Run("Valid/Normal", func(t Test) {
				var c sgr.BasicColor
				t.Expect(c.UnmarshalText([]byte("yellow"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.Yellow)
			})
			t.Run("Valid/Bright", func(t Test) {
				var c sgr.BasicColor
				t.Expect(c.UnmarshalText([]byte("bright red"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.BrightRed)
			})
			t.Run("Valid/Color", func(t Test) {
				var c sgr.Color
				t.Expect(c.UnmarshalText([]byte("yellow"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.Yellow.Color())
			})
			t.Run("Invalid", func(t Test) {
				t.Run("Direct", func(t Test) {
					var c sgr.BasicColor
					t.Expect(c.UnmarshalText([]byte("bright pink"))).ToFailWith(sgr.ErrInvalidColorText{Value: "bright pink"})
				})
				t.Run("Color", func(t Test) {
					var c sgr.Color
					t.Expect(c.UnmarshalText([]byte("bright pink"))).ToFailWith(sgr.ErrInvalidColorText{})
				})
			})
		})
	})

	t.Run("Brightness", func(t Test) {
		t.Run("String", func(t Test) {
			t.Run("Valid", func(t Test) {
				t.Expect(sgr.Normal.String()).ToEqual("Normal")
				t.Expect(sgr.Bright.String()).ToEqual("Bright")
			})
			t.Run("Invalid", func(t Test) {
				t.Expect(sgr.Brightness(5).String()).ToEqual("<!0x05>")
			})
		})

		t.Run("Validate", func(t Test) {
			t.Run("Valid", func(t Test) {
				t.Expect(sgr.Normal.Validate()).ToSucceed()
				t.Expect(sgr.Bright.Validate()).ToSucceed()
			})
			t.Run("Invalid", func(t Test) {
				t.Expect(sgr.Brightness(127).Validate()).ToFailWith(sgr.ErrInvalidBrightnessValue{127})
			})
		})

		t.Run("MarshalText", func(t Test) {
			t.Run("Valid", func(t Test) {
				text, err := sgr.Normal.MarshalText()
				t.Expect(err).ToNot(HaveOccurred())
				t.Expect(text).ToEqual([]byte("Normal"))
			})
			t.Run("Invalid", func(t Test) {
				t.Expect(sgr.Brightness(127).MarshalText()).ToFailWith(sgr.ErrInvalidBrightnessValue{})
			})
		})

		t.Run("UnmarshalText", func(t Test) {
			t.Run("Valid/Normal", func(t Test) {
				var b sgr.Brightness
				t.Expect(b.UnmarshalText([]byte("normal"))).ToSucceed()
				t.Expect(b).ToEqual(sgr.Normal)
			})
			t.Run("Valid/Bright", func(t Test) {
				var b sgr.Brightness
				t.Expect(b.UnmarshalText([]byte("bright"))).ToSucceed()
				t.Expect(b).ToEqual(sgr.Bright)
			})
			t.Run("Invalid", func(t Test) {
				var b sgr.Brightness
				t.Expect(b.UnmarshalText([]byte("extreme"))).ToFailWith(sgr.ErrInvalidBrightnessText{Value: "extreme"})
			})
		})
	})

	t.Run("Default", func(t Test) {
		t.Run("String", func(t Test) {
			t.Expect(
				sgr.Default.String(),
			).ToEqual(
				"Default",
			)
		})

		t.Run("Validate", func(t Test) {
			t.Expect(sgr.Default.Validate()).ToSucceed()
		})

		t.Run("MarshalText", func(t Test) {
			text, err := sgr.Default.MarshalText()
			t.Expect(err).ToNot(HaveOccurred())
			t.Expect(text).ToEqual([]byte("Default"))
		})

		t.Run("UnmarshalText", func(t Test) {
			t.Run("Valid/Direct", func(t Test) {
				var c sgr.DefaultColor
				t.Expect(c.UnmarshalText([]byte("default"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.Default)
			})
			t.Run("Valid/Color", func(t Test) {
				var c sgr.Color
				t.Expect(c.UnmarshalText([]byte("default"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.Default.Color())
			})
			t.Run("Invalid/Direct", func(t Test) {
				var c sgr.DefaultColor
				t.Expect(c.UnmarshalText([]byte("pink"))).ToFailWith(sgr.ErrInvalidDefaultColorText{"pink"})
			})
		})
	})

	t.Run("Palette", func(t Test) {
		t.Run("String", func(t Test) {
			t.Expect(
				sgr.PaletteColor(192).String(),
			).ToEqual(
				"#c0",
			)
		})

		t.Run("Validate", func(t Test) {
			t.Expect(sgr.PaletteColor(192).Validate()).ToSucceed()
		})

		t.Run("MarshalText", func(t Test) {
			text, err := sgr.PaletteColor(48).MarshalText()
			t.Expect(err).ToNot(HaveOccurred())
			t.Expect(text).ToEqual([]byte("#30"))
		})

		t.Run("UnmarshalText", func(t Test) {
			t.Run("Valid/Direct", func(t Test) {
				var c sgr.PaletteColor
				t.Expect(c.UnmarshalText([]byte("#40"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.PaletteColor(64))
			})
			t.Run("Valid/Color", func(t Test) {
				var c sgr.Color
				t.Expect(c.UnmarshalText([]byte("#50"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.PaletteColor(80).Color())
			})
			t.Run("Invalid/Direct", func(t Test) {
				var c sgr.PaletteColor
				t.Expect(c.UnmarshalText([]byte("pink"))).ToFailWith(sgr.ErrInvalidPaletteColorText{"pink"})
			})
			t.Run("Invalid/Color", func(t Test) {
				var c sgr.Color
				t.Expect(c.UnmarshalText([]byte("#zz"))).ToFailWith(sgr.ErrInvalidColorText{})
			})
		})
	})

	t.Run("RGB", func(t Test) {
		t.Run("String", func(t Test) {
			t.Expect(
				sgr.RGB(15, 30, 45).String(),
			).ToEqual(
				"#0f1e2d",
			)
		})

		t.Run("Validate", func(t Test) {
			t.Expect(sgr.RGB(15, 30, 45).Validate()).ToSucceed()
		})

		t.Run("MarshalText", func(t Test) {
			text, err := sgr.RGB(16, 32, 48).MarshalText()
			t.Expect(err).ToNot(HaveOccurred())
			t.Expect(text).ToEqual([]byte("#102030"))
		})

		t.Run("UnmarshalText", func(t Test) {
			t.Run("Valid/Direct", func(t Test) {
				var c sgr.RGBColor
				t.Expect(c.UnmarshalText([]byte("#405060"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.RGB(64, 80, 96))
			})
			t.Run("Valid/Color", func(t Test) {
				var c sgr.Color
				t.Expect(c.UnmarshalText([]byte("#506070"))).ToSucceed()
				t.Expect(c).ToEqual(sgr.RGB(80, 96, 112).Color())
			})
			t.Run("Invalid/Direct", func(t Test) {
				const text = "pink"
				var c sgr.RGBColor
				t.Expect(c.UnmarshalText([]byte(text))).ToFailWith(sgr.ErrInvalidRGBColorText{text})
			})
			t.Run("Invalid/Color", func(t Test) {
				const text = "#6940fz"
				var c sgr.Color
				t.Expect(c.UnmarshalText([]byte(text))).ToFailWith(sgr.ErrInvalidRGBColorText{text})
			})
		})
	})
}
