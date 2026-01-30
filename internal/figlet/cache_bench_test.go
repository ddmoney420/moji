package figlet

import (
	"testing"
)

const benchmarkFontContent = `flf2a$ 7 7 13 0 7 0 64 0
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

// BenchmarkParseUncached measures parsing without caching
func BenchmarkParseUncached(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseFont(benchmarkFontContent)
	}
}

// BenchmarkParseCached measures parsing with caching (cache hits)
func BenchmarkParseCached(b *testing.B) {
	cache := NewFontCache(10)
	cache.Set(benchmarkFontContent, nil) // Populate cache

	// Parse once to ensure it's cached
	font, _ := ParseFont(benchmarkFontContent)
	cache.Set(benchmarkFontContent, font)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(benchmarkFontContent)
	}
}

// BenchmarkParseFontCached measures the full ParseFontCached flow
func BenchmarkParseFontCached(b *testing.B) {
	ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseFontCached(benchmarkFontContent)
	}
}

// BenchmarkCacheSetAndGet measures cache set and get operations
func BenchmarkCacheSetAndGet(b *testing.B) {
	cache := NewFontCache(100)
	font := &Font{Header: Header{Height: 7}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(benchmarkFontContent, font)
		cache.Get(benchmarkFontContent)
	}
}

// BenchmarkCacheHashComputation measures hash computation speed
func BenchmarkHashContent(b *testing.B) {
	content := benchmarkFontContent

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hashContent(content)
	}
}

// BenchmarkConcurrentCacheAccess measures performance with concurrent access
func BenchmarkConcurrentCacheAccess(b *testing.B) {
	cache := NewFontCache(100)
	font := &Font{Header: Header{Height: 7}}
	cache.Set(benchmarkFontContent, font)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get(benchmarkFontContent)
		}
	})
}

// BenchmarkConcurrentCacheWrite measures concurrent writes to cache
func BenchmarkConcurrentCacheWrite(b *testing.B) {
	cache := NewFontCache(1000)
	font := &Font{Header: Header{Height: 7}}

	b.RunParallel(func(pb *testing.PB) {
		idx := 0
		for pb.Next() {
			content := benchmarkFontContent + string(rune(idx))
			cache.Set(content, font)
			idx++
		}
	})
}
