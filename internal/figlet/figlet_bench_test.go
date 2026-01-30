package figlet

import (
	"testing"
)

// Baseline font content for render benchmarks
const benchFontSmall = `flf2a$ 7 7 13 0 7 0 64 0
Font Author: Benchmark Font
This is a benchmark test font

$  $@
$  $@
$  $@
$  $@
$  $@
$  $@
$  $@@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@
██ @
██ @
██ @
   @
██ @
   @
   @@`

// BenchmarkFontParse measures font parsing performance
func BenchmarkFontParse(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ParseFont(benchFontSmall)
		}
	})

	b.Run("medium", func(b *testing.B) {
		// Create a larger font by repeating characters
		content := benchFontSmall
		for j := 0; j < 3; j++ {
			content += benchFontSmall[90:]
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ParseFont(content)
		}
	})
}

// BenchmarkFontRender measures rendering performance
func BenchmarkFontRender(b *testing.B) {
	font, _ := ParseFont(benchFontSmall)

	b.Run("shortText", func(b *testing.B) {
		text := "Hi"
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = font.Render(text)
		}
	})

	b.Run("mediumText", func(b *testing.B) {
		text := "Hello World"
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = font.Render(text)
		}
	})

	b.Run("longText", func(b *testing.B) {
		text := "The Quick Brown Fox Jumps Over The Lazy Dog"
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = font.Render(text)
		}
	})

	b.Run("repeatingText", func(b *testing.B) {
		text := "AAAAAAAAAA"
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = font.Render(text)
		}
	})

	b.Run("unicodeText", func(b *testing.B) {
		text := "こんにちは"
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = font.Render(text)
		}
	})
}

// BenchmarkParseHeader measures header parsing performance
func BenchmarkParseHeader(b *testing.B) {
	headerLine := "flf2a$ 7 7 13 0 7 0 64 0"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseHeader(headerLine)
	}
}

// BenchmarkParseAndRender measures combined parse and render
func BenchmarkParseAndRender(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			font, _ := ParseFont(benchFontSmall)
			_ = font.Render("Hi")
		}
	})

	b.Run("medium", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			font, _ := ParseFont(benchFontSmall)
			_ = font.Render("Hello World")
		}
	})
}

// BenchmarkConcurrentRender measures concurrent rendering
func BenchmarkConcurrentRender(b *testing.B) {
	font, _ := ParseFont(benchFontSmall)
	text := "Test"

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = font.Render(text)
		}
	})
}
