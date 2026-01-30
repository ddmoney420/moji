package filters

import (
	"testing"
)

// Baseline test strings of various sizes
const (
	shortText   = "Hello"
	mediumText  = "The Quick Brown Fox Jumps Over The Lazy Dog"
	longText    = "The Quick Brown Fox Jumps Over The Lazy Dog. The Quick Brown Fox Jumps Over The Lazy Dog. The Quick Brown Fox Jumps Over The Lazy Dog."
	multiline   = "Line One\nLine Two\nLine Three\nLine Four\nLine Five"
	largeText   = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris."
)

// BenchmarkIndividualFilters measures performance of each filter
func BenchmarkIndividualFilters(b *testing.B) {
	filters := []struct {
		name string
		fn   Filter
	}{
		{"Metal", Metal},
		{"Rainbow", Rainbow},
		{"Crop", Crop},
		{"Flip", Flip},
		{"Flop", Flop},
		{"Rotate180", Rotate180},
		{"Border", Border},
		{"Shadow", Shadow},
		{"Shadow3D", Shadow3D},
		{"Glitch", Glitch},
		{"Matrix", Matrix},
		{"Fire", Fire},
		{"Ice", Ice},
		{"Neon", Neon},
		{"RetroGreen", RetroGreen},
		{"Bold", Bold},
		{"Italic", Italic},
		{"Underline", Underline},
		{"Strikethrough", Strikethrough},
		{"Blink", Blink},
		{"Dim", Dim},
		{"Invert", Invert},
	}

	for _, f := range filters {
		b.Run(f.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = f.fn(mediumText)
			}
		})
	}
}

// BenchmarkFilterByTextSize measures filter performance with different text sizes
func BenchmarkFilterByTextSize(b *testing.B) {
	texts := []struct {
		name string
		text string
	}{
		{"short", shortText},
		{"medium", mediumText},
		{"long", longText},
		{"multiline", multiline},
		{"large", largeText},
	}

	for _, txt := range texts {
		b.Run(txt.name, func(b *testing.B) {
			b.Run("Metal", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Metal(txt.text)
				}
			})

			b.Run("Rainbow", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Rainbow(txt.text)
				}
			})

			b.Run("Glitch", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Glitch(txt.text)
				}
			})

			b.Run("Border", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Border(txt.text)
				}
			})
		})
	}
}

// BenchmarkColorFilters measures color-intensive filters
func BenchmarkColorFilters(b *testing.B) {
	b.Run("Metal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Metal(largeText)
		}
	})

	b.Run("Rainbow", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Rainbow(largeText)
		}
	})

	b.Run("Neon", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Neon(largeText)
		}
	})

	b.Run("Fire", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Fire(largeText)
		}
	})

	b.Run("Ice", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Ice(largeText)
		}
	})

	b.Run("RetroGreen", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = RetroGreen(largeText)
		}
	})
}

// BenchmarkStructuralFilters measures layout-changing filters
func BenchmarkStructuralFilters(b *testing.B) {
	b.Run("Crop", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Crop(multiline)
		}
	})

	b.Run("Flip", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Flip(mediumText)
		}
	})

	b.Run("Flop", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Flop(multiline)
		}
	})

	b.Run("Rotate180", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Rotate180(mediumText)
		}
	})

	b.Run("Border", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Border(multiline)
		}
	})

	b.Run("Shadow", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Shadow(mediumText)
		}
	})

	b.Run("Shadow3D", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Shadow3D(mediumText)
		}
	})
}

// BenchmarkEffectFilters measures effect-based filters
func BenchmarkEffectFilters(b *testing.B) {
	b.Run("Glitch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Glitch(mediumText)
		}
	})

	b.Run("Matrix", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Matrix(mediumText)
		}
	})

	b.Run("Blink", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Blink(mediumText)
		}
	})
}

// BenchmarkTextStyleFilters measures text styling filters
func BenchmarkTextStyleFilters(b *testing.B) {
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

	b.Run("Underline", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Underline(mediumText)
		}
	})

	b.Run("Strikethrough", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Strikethrough(mediumText)
		}
	})

	b.Run("Dim", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Dim(mediumText)
		}
	})

	b.Run("Invert", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Invert(mediumText)
		}
	})
}

// BenchmarkChainFilters measures chaining multiple filters
func BenchmarkChainFilters(b *testing.B) {
	b.Run("chain_2", func(b *testing.B) {
		filterNames := []string{"bold", "underline"}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Chain(mediumText, filterNames)
		}
	})

	b.Run("chain_3", func(b *testing.B) {
		filterNames := []string{"bold", "italic", "rainbow"}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Chain(mediumText, filterNames)
		}
	})

	b.Run("chain_5", func(b *testing.B) {
		filterNames := []string{"bold", "italic", "underline", "rainbow", "invert"}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Chain(mediumText, filterNames)
		}
	})
}

// BenchmarkParseChain measures filter chain parsing
func BenchmarkParseChain(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		spec := "bold,italic"

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ParseChain(spec)
		}
	})

	b.Run("complex", func(b *testing.B) {
		spec := "bold, italic, underline, rainbow, shadow, border"

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ParseChain(spec)
		}
	})
}

// BenchmarkComplexWorkflow measures typical filter workflows
func BenchmarkComplexWorkflow(b *testing.B) {
	b.Run("styling_workflow", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			text := mediumText
			text = Bold(text)
			text = Rainbow(text)
			_ = Border(text)
		}
	})

	b.Run("effect_workflow", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			text := mediumText
			text = Shadow(text)
			text = Glitch(text)
			_ = Border(text)
		}
	})
}
