package figlet

import (
	"hash/fnv"
	"sync"
	"testing"
)

// Test sample font content
const sampleFontContent = `flf2a$ 7 7 13 0 7 0 64 0
Font Author: Test
This is a test font

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
   @@`

const sampleFontContent2 = `flf2a$ 7 7 13 0 7 0 64 0
Font Author: Test2
Different font

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
   @@`

func TestNewFontCache(t *testing.T) {
	cache := NewFontCache(5)
	if cache.capacity != 5 {
		t.Errorf("expected capacity 5, got %d", cache.capacity)
	}
	if cache.Size() != 0 {
		t.Errorf("expected size 0, got %d", cache.Size())
	}
}

func TestNewFontCacheDefaults(t *testing.T) {
	cache := NewFontCache(0)
	if cache.capacity != 10 {
		t.Errorf("expected default capacity 10, got %d", cache.capacity)
	}

	cache = NewFontCache(-1)
	if cache.capacity != 10 {
		t.Errorf("expected default capacity 10 for negative, got %d", cache.capacity)
	}
}

func TestHashContent(t *testing.T) {
	content1 := "test content"
	content2 := "test content"
	content3 := "different content"

	hash1 := hashContent(content1)
	hash2 := hashContent(content2)
	hash3 := hashContent(content3)

	if hash1 != hash2 {
		t.Errorf("same content should produce same hash")
	}
	if hash1 == hash3 {
		t.Errorf("different content should produce different hash")
	}

	// Verify it uses fnv-1a algorithm
	h := fnv.New64a()
	h.Write([]byte(content1))
	expected := h.Sum64()
	if hash1 != expected {
		t.Errorf("expected hash %d, got %d", expected, hash1)
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cache := NewFontCache(10)
	font := &Font{
		Header: Header{Height: 7},
		Characters: map[rune][]string{
			'A': {"line1", "line2"},
		},
	}

	cache.Set(sampleFontContent, font)
	if cache.Size() != 1 {
		t.Errorf("expected size 1, got %d", cache.Size())
	}

	retrieved, found := cache.Get(sampleFontContent)
	if !found {
		t.Errorf("expected font to be found in cache")
	}
	if retrieved != font {
		t.Errorf("expected same font object")
	}
}

func TestCacheMiss(t *testing.T) {
	cache := NewFontCache(10)
	_, found := cache.Get("nonexistent content")
	if found {
		t.Errorf("expected cache miss for nonexistent content")
	}
}

func TestCacheStats(t *testing.T) {
	cache := NewFontCache(10)
	font := &Font{Header: Header{Height: 7}}

	cache.Set(sampleFontContent, font)

	// Miss
	cache.Get("nonexistent")
	stats := cache.Stats()
	if stats.Misses != 1 {
		t.Errorf("expected 1 miss, got %d", stats.Misses)
	}

	// Hit
	cache.Get(sampleFontContent)
	stats = cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("expected 1 hit, got %d", stats.Hits)
	}
}

func TestCacheEviction(t *testing.T) {
	cache := NewFontCache(3)
	font1 := &Font{Header: Header{Height: 1}}
	font2 := &Font{Header: Header{Height: 2}}
	font3 := &Font{Header: Header{Height: 3}}
	font4 := &Font{Header: Header{Height: 4}}

	cache.Set("content1", font1)
	cache.Set("content2", font2)
	cache.Set("content3", font3)

	if cache.Size() != 3 {
		t.Errorf("expected size 3, got %d", cache.Size())
	}

	// Adding 4th font should evict the least recently used (content1)
	cache.Set("content4", font4)

	if cache.Size() != 3 {
		t.Errorf("expected size 3 after eviction, got %d", cache.Size())
	}

	_, found := cache.Get("content1")
	if found {
		t.Errorf("expected content1 to be evicted")
	}

	_, found = cache.Get("content2")
	if !found {
		t.Errorf("expected content2 to still be in cache")
	}

	stats := cache.Stats()
	if stats.Evictions != 1 {
		t.Errorf("expected 1 eviction, got %d", stats.Evictions)
	}
}

func TestCacheLRUOrdering(t *testing.T) {
	cache := NewFontCache(2)
	font1 := &Font{Header: Header{Height: 1}}
	font2 := &Font{Header: Header{Height: 2}}
	font3 := &Font{Header: Header{Height: 3}}

	cache.Set("content1", font1)
	cache.Set("content2", font2)

	// Access content1 to make it recently used
	cache.Get("content1")

	// Add content3, which should evict content2 (least recently used)
	cache.Set("content3", font3)

	_, found := cache.Get("content1")
	if !found {
		t.Errorf("expected content1 to still be in cache")
	}

	_, found = cache.Get("content2")
	if found {
		t.Errorf("expected content2 to be evicted (LRU)")
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewFontCache(10)
	font := &Font{Header: Header{Height: 7}}

	cache.Set(sampleFontContent, font)
	if cache.Size() != 1 {
		t.Errorf("expected size 1 before clear")
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", cache.Size())
	}

	_, found := cache.Get(sampleFontContent)
	if found {
		t.Errorf("expected cache to be empty after clear")
	}
}

func TestSetCapacity(t *testing.T) {
	cache := NewFontCache(3)
	font := &Font{Header: Header{Height: 1}}

	cache.Set("content1", font)
	cache.Set("content2", font)
	cache.Set("content3", font)

	if cache.Size() != 3 {
		t.Errorf("expected size 3")
	}

	// Reduce capacity to 2
	cache.SetCapacity(2)

	if cache.Size() != 2 {
		t.Errorf("expected size 2 after capacity reduction, got %d", cache.Size())
	}

	if cache.capacity != 2 {
		t.Errorf("expected capacity 2, got %d", cache.capacity)
	}
}

func TestSetCapacityDefaults(t *testing.T) {
	cache := NewFontCache(5)
	cache.SetCapacity(0)
	if cache.capacity != 10 {
		t.Errorf("expected default capacity 10, got %d", cache.capacity)
	}

	cache.SetCapacity(-5)
	if cache.capacity != 10 {
		t.Errorf("expected default capacity 10 for negative, got %d", cache.capacity)
	}
}

func TestGlobalCache(t *testing.T) {
	// Clear global cache first
	ClearCache()

	font, err := ParseFontCached(sampleFontContent)
	if err != nil {
		t.Errorf("expected no error parsing font, got %v", err)
	}
	if font == nil {
		t.Errorf("expected valid font")
	}

	// Second call should hit cache
	stats := GetCacheStats()
	initialHits := stats.Hits

	_, err = ParseFontCached(sampleFontContent)
	if err != nil {
		t.Errorf("expected no error on second parse")
	}

	// Verify cache hit
	stats = GetCacheStats()
	if stats.Hits <= initialHits {
		t.Errorf("expected cache hit, hits before: %d, after: %d", initialHits, stats.Hits)
	}
}

func TestGlobalCacheSetCapacity(t *testing.T) {
	ClearCache()
	SetCacheCapacity(2)

	font := &Font{Header: Header{Height: 1}}

	// Use a local cache to set up the state
	// The global functions use defaultCache
	defaultCache.Set("content1", font)
	defaultCache.Set("content2", font)

	if defaultCache.Size() != 2 {
		t.Errorf("expected size 2")
	}

	SetCacheCapacity(1)

	if defaultCache.Size() != 1 {
		t.Errorf("expected size 1 after capacity reduction, got %d", defaultCache.Size())
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewFontCache(100)
	font := &Font{Header: Header{Height: 7}}

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Concurrent writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			content := sampleFontContent + " version " + string(rune(idx))
			cache.Set(content, font)
		}(i)
	}

	// Concurrent readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = cache.Get(sampleFontContent)
			_, _ = cache.Get(sampleFontContent2)
		}()
	}

	wg.Wait()
	close(errors)

	if cache.Size() < 0 {
		t.Errorf("cache size should not be negative")
	}
}

func TestConcurrentParseAndCache(t *testing.T) {
	ClearCache()
	SetCacheCapacity(50)

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Multiple goroutines parsing the same font
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			font, err := ParseFontCached(sampleFontContent)
			if err == nil && font != nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if successCount != 10 {
		t.Errorf("expected 10 successful parses, got %d", successCount)
	}

	stats := GetCacheStats()
	if stats.Hits == 0 {
		t.Errorf("expected cache hits with concurrent access")
	}
}

func TestNilFontNotCached(t *testing.T) {
	cache := NewFontCache(10)
	cache.Set(sampleFontContent, nil)

	if cache.Size() != 0 {
		t.Errorf("expected nil fonts not to be cached, size: %d", cache.Size())
	}
}

func TestDuplicateSetUpdates(t *testing.T) {
	cache := NewFontCache(10)
	font1 := &Font{Header: Header{Height: 1}}
	font2 := &Font{Header: Header{Height: 2}}

	cache.Set(sampleFontContent, font1)
	if cache.Size() != 1 {
		t.Errorf("expected size 1")
	}

	// Set same content with different font
	cache.Set(sampleFontContent, font2)
	if cache.Size() != 1 {
		t.Errorf("expected size 1 after update")
	}

	retrieved, _ := cache.Get(sampleFontContent)
	if retrieved != font2 {
		t.Errorf("expected updated font")
	}
}

func TestParseError(t *testing.T) {
	ClearCache()

	invalidFont := "not a valid font"
	font, err := ParseFontCached(invalidFont)
	if err == nil {
		t.Errorf("expected error parsing invalid font")
	}
	if font != nil {
		t.Errorf("expected nil font on error")
	}

	// Error result should not be cached
	// Verify by calling again - should get same error
	_, err2 := ParseFontCached(invalidFont)
	if err2 == nil {
		t.Errorf("expected error on second parse of invalid font")
	}
}

func TestCacheEdgeCases(t *testing.T) {
	cache := NewFontCache(1)
	font := &Font{Header: Header{Height: 1}}

	// Fill cache
	cache.Set("content1", font)

	// Access it (should remain)
	cache.Get("content1")

	// Add new entry
	cache.Set("content2", font)

	// First should be evicted
	_, found := cache.Get("content1")
	if found {
		t.Errorf("expected content1 to be evicted")
	}

	// Second should still exist
	_, found = cache.Get("content2")
	if !found {
		t.Errorf("expected content2 to exist")
	}
}

func TestStatsAccuracy(t *testing.T) {
	cache := NewFontCache(10)
	font := &Font{Header: Header{Height: 1}}

	cache.Set("content1", font)

	// 5 misses
	for i := 0; i < 5; i++ {
		cache.Get("nonexistent" + string(rune(i)))
	}

	// 3 hits
	for i := 0; i < 3; i++ {
		cache.Get("content1")
	}

	stats := cache.Stats()
	if stats.Hits != 3 {
		t.Errorf("expected 3 hits, got %d", stats.Hits)
	}
	if stats.Misses != 5 {
		t.Errorf("expected 5 misses, got %d", stats.Misses)
	}
}

func TestCacheSize(t *testing.T) {
	cache := NewFontCache(10)
	if cache.Size() != 0 {
		t.Errorf("expected initial size 0, got %d", cache.Size())
	}

	font := &Font{Header: Header{Height: 1}}
	cache.Set("content1", font)
	if cache.Size() != 1 {
		t.Errorf("expected size 1, got %d", cache.Size())
	}

	cache.Set("content2", font)
	if cache.Size() != 2 {
		t.Errorf("expected size 2, got %d", cache.Size())
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", cache.Size())
	}
}

func TestCacheGetRefreshesLRU(t *testing.T) {
	cache := NewFontCache(2)
	font := &Font{Header: Header{Height: 1}}

	cache.Set("content1", font)
	cache.Set("content2", font)

	// Access content1 to refresh it
	cache.Get("content1")

	// Add content3, should evict content2 not content1
	cache.Set("content3", font)

	_, found := cache.Get("content1")
	if !found {
		t.Errorf("expected content1 to remain after being accessed")
	}

	_, found = cache.Get("content2")
	if found {
		t.Errorf("expected content2 to be evicted")
	}
}

func TestNewFontCacheWithNegativeCapacity(t *testing.T) {
	cache := NewFontCache(-100)
	if cache.capacity != 10 {
		t.Errorf("expected default capacity 10 for negative input, got %d", cache.capacity)
	}
}

func TestMultipleMisses(t *testing.T) {
	cache := NewFontCache(10)

	for i := 0; i < 20; i++ {
		cache.Get("nonexistent" + string(rune(i)))
	}

	stats := cache.Stats()
	if stats.Misses != 20 {
		t.Errorf("expected 20 misses, got %d", stats.Misses)
	}
}

func TestEvictionCountWithCapacityChange(t *testing.T) {
	cache := NewFontCache(5)
	font := &Font{Header: Header{Height: 1}}

	// Fill cache
	for i := 0; i < 5; i++ {
		cache.Set("content"+string(rune(i)), font)
	}

	stats := cache.Stats()
	if stats.Evictions != 0 {
		t.Errorf("expected 0 evictions initially, got %d", stats.Evictions)
	}

	// Add one more to trigger eviction
	cache.Set("content5", font)
	stats = cache.Stats()
	if stats.Evictions != 1 {
		t.Errorf("expected 1 eviction, got %d", stats.Evictions)
	}

	// Reduce capacity to trigger more evictions
	cache.SetCapacity(2)
	stats = cache.Stats()
	if stats.Evictions != 4 {
		t.Errorf("expected 4 evictions after capacity reduction, got %d", stats.Evictions)
	}
}

func TestParseAndCacheSameFontMultipleTimes(t *testing.T) {
	ClearCache()

	// Parse same font 5 times
	var fonts []*Font
	var errs []error
	uniqueContent := sampleFontContent + "unique_test_12345"
	for i := 0; i < 5; i++ {
		font, err := ParseFontCached(uniqueContent)
		fonts = append(fonts, font)
		errs = append(errs, err)
	}

	// Check all parses succeeded
	for i, err := range errs {
		if err != nil {
			t.Errorf("parse %d failed: %v", i, err)
		}
	}

	// All should return same cached object
	for i := 1; i < len(fonts); i++ {
		if fonts[i] != fonts[0] {
			t.Errorf("expected same cached object at parse %d", i)
		}
	}

	// Check cache statistics show cache hits
	stats := GetCacheStats()
	if stats.Hits < 4 {
		t.Errorf("expected at least 4 cache hits (first is miss), got %d", stats.Hits)
	}
}
