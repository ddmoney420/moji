package dither

import (
	"image"
	"image/color"
	"testing"
)

// createBenchGradientImage creates a grayscale gradient image for benchmarking
func createBenchGradientImage(width, height int) image.Image {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewGray(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a gradient from black to white
			val := uint8((x * 255) / width)
			img.Set(x, y, color.Gray{Y: val})
		}
	}

	return img
}

// BenchmarkDitheringAlgorithms measures all dithering algorithms
func BenchmarkDitheringAlgorithms(b *testing.B) {
	algorithms := []struct {
		name string
		algo Algorithm
	}{
		{"FloydSteinberg", FloydSteinberg},
		{"Bayer2x2", Bayer2x2},
		{"Bayer4x4", Bayer4x4},
		{"Bayer8x8", Bayer8x8},
		{"Atkinson", Atkinson},
		{"Sierra", Sierra},
		{"SierraLite", SierraLite},
		{"Stucki", Stucki},
		{"Burkes", Burkes},
		{"JarvisJudice", JarvisJudice},
	}

	for _, algo := range algorithms {
		b.Run(algo.name, func(b *testing.B) {
			img := createBenchGradientImage(100, 100)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Apply(img, algo.algo)
			}
		})
	}
}

// BenchmarkDitheringBySizes measures dithering performance on different image sizes
func BenchmarkDitheringBySizes(b *testing.B) {
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"small_50x50", 50, 50},
		{"medium_100x100", 100, 100},
		{"large_200x200", 200, 200},
		{"xlarge_500x500", 500, 500},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			b.Run("floyd_steinberg", func(b *testing.B) {
				img := createBenchGradientImage(size.width, size.height)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Apply(img, FloydSteinberg)
				}
			})

			b.Run("bayer4x4", func(b *testing.B) {
				img := createBenchGradientImage(size.width, size.height)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Apply(img, Bayer4x4)
				}
			})

			b.Run("sierra", func(b *testing.B) {
				img := createBenchGradientImage(size.width, size.height)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Apply(img, Sierra)
				}
			})
		})
	}
}

// BenchmarkOrderedDithering measures ordered dithering variants
func BenchmarkOrderedDithering(b *testing.B) {
	img := createBenchGradientImage(200, 200)

	b.Run("Bayer2x2", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Bayer2x2)
		}
	})

	b.Run("Bayer4x4", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Bayer4x4)
		}
	})

	b.Run("Bayer8x8", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Bayer8x8)
		}
	})
}

// BenchmarkErrorDiffusion measures error diffusion variants
func BenchmarkErrorDiffusion(b *testing.B) {
	img := createBenchGradientImage(200, 200)

	b.Run("FloydSteinberg", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, FloydSteinberg)
		}
	})

	b.Run("Sierra", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Sierra)
		}
	})

	b.Run("SierraLite", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, SierraLite)
		}
	})

	b.Run("Atkinson", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Atkinson)
		}
	})

	b.Run("Stucki", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Stucki)
		}
	})

	b.Run("Burkes", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, Burkes)
		}
	})

	b.Run("JarvisJudice", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Apply(img, JarvisJudice)
		}
	})
}

// BenchmarkToGrayscale measures grayscale conversion with different levels
func BenchmarkToGrayscale(b *testing.B) {
	img := createBenchGradientImage(200, 200).(*image.Gray)

	for _, levels := range []int{2, 4, 8, 16, 64, 256} {
		b.Run("levels_"+string(rune(levels/10))+"_"+string(rune(levels%10)), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ToGrayscale(img, levels)
			}
		})
	}
}

// BenchmarkContrastStretch measures contrast stretching performance
func BenchmarkContrastStretch(b *testing.B) {
	img := createBenchGradientImage(200, 200).(*image.Gray)

	b.Run("50x50_output", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ContrastStretch(img, 50, 200)
		}
	})

	b.Run("0x255_output", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ContrastStretch(img, 0, 255)
		}
	})
}

// BenchmarkConcurrentDithering measures concurrent dithering
func BenchmarkConcurrentDithering(b *testing.B) {
	img := createBenchGradientImage(100, 100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Apply(img, FloydSteinberg)
		}
	})
}

// BenchmarkComplexWorkflow measures typical dithering workflow
func BenchmarkComplexWorkflow(b *testing.B) {
	img := createBenchGradientImage(150, 150).(*image.Gray)

	b.Run("dither_then_grayscale", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dithered := Apply(img, FloydSteinberg)
			_ = ToGrayscale(dithered, 8)
		}
	})

	b.Run("stretch_then_dither", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stretched := ContrastStretch(img, 0, 255)
			_ = Apply(stretched, FloydSteinberg)
		}
	})

	b.Run("all_processing", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stretched := ContrastStretch(img, 0, 255)
			dithered := Apply(stretched, FloydSteinberg)
			_ = ToGrayscale(dithered, 4)
		}
	})
}
