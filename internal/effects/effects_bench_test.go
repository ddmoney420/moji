package effects

import (
	"testing"
)

// Baseline test strings of various sizes
const (
	shortText   = "Hello"
	mediumText  = "The Quick Brown Fox"
	longText    = "The Quick Brown Fox Jumps Over The Lazy Dog The Quick Brown Fox Jumps Over The Lazy Dog"
	repeating   = "AAAAAAAAAA"
	mixed       = "Hello123!@#World456$%^Test"
)

// BenchmarkIndividualEffects measures performance of each text effect
func BenchmarkIndividualEffects(b *testing.B) {
	effects := []struct {
		name string
		fn   func(string) string
	}{
		{"Flip", Flip},
		{"Reverse", Reverse},
		{"Mirror", Mirror},
		{"Wave", Wave},
		{"Zalgo", func(s string) string { return Zalgo(s, 3) }},
		{"Bubble", Bubble},
		{"Square", Square},
		{"Bold", Bold},
		{"Italic", Italic},
		{"Strikethrough", Strikethrough},
		{"Underline", Underline},
		{"SmallCaps", SmallCaps},
		{"Fullwidth", Fullwidth},
		{"Monospace", Monospace},
		{"Script", Script},
		{"Fraktur", Fraktur},
		{"DoubleStruck", DoubleStruck},
		{"Sparkle", Sparkle},
	}

	for _, e := range effects {
		b.Run(e.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = e.fn(mediumText)
			}
		})
	}
}

// BenchmarkEffectByTextSize measures effect performance with different text sizes
func BenchmarkEffectByTextSize(b *testing.B) {
	texts := []struct {
		name string
		text string
	}{
		{"short", shortText},
		{"medium", mediumText},
		{"long", longText},
		{"repeating", repeating},
		{"mixed", mixed},
	}

	for _, txt := range texts {
		b.Run(txt.name, func(b *testing.B) {
			b.Run("Flip", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Flip(txt.text)
				}
			})

			b.Run("Zalgo", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Zalgo(txt.text, 3)
				}
			})

			b.Run("Fullwidth", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Fullwidth(txt.text)
				}
			})

			b.Run("Script", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Script(txt.text)
				}
			})
		})
	}
}

// BenchmarkTransformationEffects measures text transformation effects
func BenchmarkTransformationEffects(b *testing.B) {
	b.Run("Flip", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Flip(longText)
		}
	})

	b.Run("Reverse", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Reverse(longText)
		}
	})

	b.Run("Mirror", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Mirror(mediumText)
		}
	})

	b.Run("Wave", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Wave(mediumText)
		}
	})
}

// BenchmarkDecorativeEffects measures decorative text effects
func BenchmarkDecorativeEffects(b *testing.B) {
	b.Run("Bubble", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Bubble(mediumText)
		}
	})

	b.Run("Square", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Square(mediumText)
		}
	})

	b.Run("Sparkle", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Sparkle(mediumText)
		}
	})
}

// BenchmarkZalgoVariations measures Zalgo effect at different intensities
func BenchmarkZalgoVariations(b *testing.B) {
	b.Run("intensity_1", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Zalgo(mediumText, 1)
		}
	})

	b.Run("intensity_3", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Zalgo(mediumText, 3)
		}
	})

	b.Run("intensity_6", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Zalgo(mediumText, 6)
		}
	})

	b.Run("longText_intensity_3", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Zalgo(longText, 3)
		}
	})
}

// BenchmarkStyleEffects measures text styling effects
func BenchmarkStyleEffects(b *testing.B) {
	b.Run("Bold", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Bold(mediumText)
		}
	})

	b.Run("Italic", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Italic(mediumText)
		}
	})

	b.Run("Strikethrough", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Strikethrough(mediumText)
		}
	})

	b.Run("Underline", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Underline(mediumText)
		}
	})

	b.Run("SmallCaps", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = SmallCaps(mediumText)
		}
	})
}

// BenchmarkFontVariationEffects measures font variation effects
func BenchmarkFontVariationEffects(b *testing.B) {
	b.Run("Fullwidth", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Fullwidth(mediumText)
		}
	})

	b.Run("Monospace", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Monospace(mediumText)
		}
	})

	b.Run("Script", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Script(mediumText)
		}
	})

	b.Run("Fraktur", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Fraktur(mediumText)
		}
	})

	b.Run("DoubleStruck", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = DoubleStruck(mediumText)
		}
	})
}

// BenchmarkApplyEffect measures the Apply dispatcher function
func BenchmarkApplyEffect(b *testing.B) {
	effects := []struct {
		name string
		effect string
	}{
		{"flip", "flip"},
		{"reverse", "reverse"},
		{"zalgo", "zalgo"},
		{"bubble", "bubble"},
		{"bold", "bold"},
		{"fullwidth", "fullwidth"},
		{"script", "script"},
		{"sparkle", "sparkle"},
	}

	for _, e := range effects {
		b.Run(e.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Apply(e.effect, mediumText)
			}
		})
	}
}

// BenchmarkComplexWorkflow measures typical effect workflows
func BenchmarkComplexWorkflow(b *testing.B) {
	b.Run("flip_then_uppercase", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Flip(mediumText)
		}
	})

	b.Run("bubble_then_sparkle", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			text := Bubble(mediumText)
			_ = Sparkle(text)
		}
	})

	b.Run("bold_italic_underline", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			text := Bold(mediumText)
			text = Italic(text)
			_ = Underline(text)
		}
	})

	b.Run("zalgo_intense", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Zalgo(mediumText, 6)
		}
	})
}

// BenchmarkConcurrentEffects measures concurrent effect application
func BenchmarkConcurrentEffects(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Flip(mediumText)
		}
	})
}

// BenchmarkMemoryAllocations measures memory characteristics of effects
func BenchmarkMemoryAllocations(b *testing.B) {
	b.Run("Flip", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Flip(longText)
		}
	})

	b.Run("Zalgo", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Zalgo(longText, 3)
		}
	})

	b.Run("Fullwidth", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Fullwidth(longText)
		}
	})

	b.Run("Mirror", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Mirror(longText)
		}
	})
}
