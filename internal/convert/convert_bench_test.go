package convert

import (
	"image"
	"image/color"
	"testing"
)

// createGradientTestImage creates a simple test image with gradient
func createGradientTestImage(width, height int) image.Image {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a gradient from black to white
			val := uint8((x * 255) / width)
			img.Set(x, y, color.RGBA{val, val, val, 255})
		}
	}

	return img
}

// BenchmarkFromImage measures image to ASCII conversion performance
func BenchmarkFromImage(b *testing.B) {
	b.Run("50x50", func(b *testing.B) {
		img := createGradientTestImage(50, 50)
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})

	b.Run("100x100", func(b *testing.B) {
		img := createGradientTestImage(100, 100)
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})

	b.Run("200x200", func(b *testing.B) {
		img := createGradientTestImage(200, 200)
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})

	b.Run("500x500", func(b *testing.B) {
		img := createGradientTestImage(500, 500)
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})
}

// BenchmarkSequentialProcessing measures sequential processing performance
func BenchmarkSequentialProcessing(b *testing.B) {
	b.Run("100x100", func(b *testing.B) {
		img := createGradientTestImage(100, 100)
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = fromImageSequential(img, opts)
		}
	})

	b.Run("300x300", func(b *testing.B) {
		img := createGradientTestImage(300, 300)
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = fromImageSequential(img, opts)
		}
	})
}

// BenchmarkSampleRegion measures sampling performance
func BenchmarkSampleRegion(b *testing.B) {
	img := createGradientTestImage(100, 100)

	b.Run("small_region", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _, _ = sampleRegion(img, 25, 25, 5, 5)
		}
	})

	b.Run("medium_region", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _, _ = sampleRegion(img, 25, 25, 20, 20)
		}
	})

	b.Run("large_region", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _, _ = sampleRegion(img, 10, 10, 50, 50)
		}
	})
}

// BenchmarkDetectEdges measures edge detection performance
func BenchmarkDetectEdges(b *testing.B) {
	b.Run("50x50", func(b *testing.B) {
		img := createGradientTestImage(50, 50)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = detectEdges(img, 50, 50, 1.0, 1.0)
		}
	})

	b.Run("100x100", func(b *testing.B) {
		img := createGradientTestImage(100, 100)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = detectEdges(img, 100, 100, 1.0, 1.0)
		}
	})

	b.Run("200x200", func(b *testing.B) {
		img := createGradientTestImage(200, 200)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = detectEdges(img, 200, 200, 1.0, 1.0)
		}
	})
}

// BenchmarkDifferentCharsets measures performance with different charsets
func BenchmarkDifferentCharsets(b *testing.B) {
	img := createGradientTestImage(100, 100)

	for _, charset := range []string{"standard", "simple", "detailed", "braille", "blocks"} {
		b.Run(charset, func(b *testing.B) {
			opts := Options{
				Width:   80,
				Charset: GetCharset(charset),
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = FromImage(img, opts)
			}
		})
	}
}

// BenchmarkDifferentSizes measures performance with different output sizes
func BenchmarkDifferentSizes(b *testing.B) {
	img := createGradientTestImage(200, 200)

	for _, width := range []int{40, 80, 120, 160} {
		b.Run("width_"+string(rune(width)), func(b *testing.B) {
			opts := Options{
				Width:   width,
				Charset: CharSets["standard"],
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = FromImage(img, opts)
			}
		})
	}
}

// BenchmarkConvertWithOptions measures conversion with various options
func BenchmarkConvertWithOptions(b *testing.B) {
	img := createGradientTestImage(100, 100)

	b.Run("default", func(b *testing.B) {
		opts := DefaultOptions()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})

	b.Run("inverted", func(b *testing.B) {
		opts := Options{
			Width:   80,
			Charset: CharSets["standard"],
			Invert:  true,
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})

	b.Run("edge_detect", func(b *testing.B) {
		opts := Options{
			Width:      80,
			Charset:    CharSets["standard"],
			EdgeDetect: true,
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = FromImage(img, opts)
		}
	})
}

// BenchmarkConcurrentConversion measures concurrent image conversion
func BenchmarkConcurrentConversion(b *testing.B) {
	img := createGradientTestImage(100, 100)
	opts := DefaultOptions()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = FromImage(img, opts)
		}
	})
}
