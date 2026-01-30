package figlet

import (
	"container/list"
	"hash/fnv"
	"sync"
)

// cacheEntry represents a single cache entry with LRU metadata
type cacheEntry struct {
	hash  uint64
	font  *Font
	bytes int
}

// FontCache provides thread-safe LRU caching for parsed fonts
type FontCache struct {
	mu       sync.RWMutex
	capacity int
	cache    map[uint64]*cacheEntry
	lru      *list.List // doubly-linked list for LRU ordering
	stats    CacheStats
}

// CacheStats contains cache statistics
type CacheStats struct {
	Hits      int
	Misses    int
	Evictions int
}

var defaultCache = &FontCache{
	capacity: 10,
	cache:    make(map[uint64]*cacheEntry),
	lru:      list.New(),
	stats:    CacheStats{},
}

// NewFontCache creates a new font cache with specified capacity
func NewFontCache(capacity int) *FontCache {
	if capacity <= 0 {
		capacity = 10
	}
	return &FontCache{
		capacity: capacity,
		cache:    make(map[uint64]*cacheEntry),
		lru:      list.New(),
		stats:    CacheStats{},
	}
}

// hashContent returns the FNV-1a hash of the content
func hashContent(content string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(content))
	return h.Sum64()
}

// Set stores a font in the cache
func (c *FontCache) Set(content string, font *Font) {
	if font == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	hash := hashContent(content)

	// If entry already exists, update it
	if entry, exists := c.cache[hash]; exists {
		entry.font = font
		c.lru.MoveToFront(c.lru.Front()) // Refresh position
		return
	}

	// Add new entry
	entry := &cacheEntry{
		hash:  hash,
		font:  font,
		bytes: len(content),
	}
	c.cache[hash] = entry
	c.lru.PushFront(hash)

	// Evict least recently used if over capacity
	if len(c.cache) > c.capacity {
		c.evictLRU()
	}
}

// Get retrieves a font from the cache if it exists
func (c *FontCache) Get(content string) (*Font, bool) {
	c.mu.RLock()
	hash := hashContent(content)
	entry, exists := c.cache[hash]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		c.stats.Misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	c.stats.Hits++
	// Move to front (most recently used)
	for elem := c.lru.Front(); elem != nil; elem = elem.Next() {
		if elem.Value == hash {
			c.lru.MoveToFront(elem)
			break
		}
	}
	c.mu.Unlock()

	return entry.font, true
}

// evictLRU removes the least recently used entry (must be called with lock held)
func (c *FontCache) evictLRU() {
	if elem := c.lru.Back(); elem != nil {
		hash := elem.Value.(uint64)
		delete(c.cache, hash)
		c.lru.Remove(elem)
		c.stats.Evictions++
	}
}

// Clear removes all entries from the cache
func (c *FontCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[uint64]*cacheEntry)
	c.lru = list.New()
}

// Stats returns current cache statistics
func (c *FontCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		Hits:      c.stats.Hits,
		Misses:    c.stats.Misses,
		Evictions: c.stats.Evictions,
	}
}

// Size returns the number of fonts currently in the cache
func (c *FontCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

// SetCapacity updates the cache capacity
func (c *FontCache) SetCapacity(capacity int) {
	if capacity <= 0 {
		capacity = 10
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = capacity

	// Evict entries if over new capacity
	for len(c.cache) > c.capacity {
		c.evictLRU()
	}
}

// Global cache management functions

// SetCacheCapacity sets the capacity of the default global cache
func SetCacheCapacity(capacity int) {
	defaultCache.SetCapacity(capacity)
}

// ClearCache clears the default global cache
func ClearCache() {
	defaultCache.Clear()
}

// CacheStats returns statistics from the default global cache
func GetCacheStats() CacheStats {
	return defaultCache.Stats()
}

// ParseFontCached parses a font with LRU caching using the default global cache
func ParseFontCached(content string) (*Font, error) {
	// Check cache first
	if font, found := defaultCache.Get(content); found {
		return font, nil
	}

	// Parse if not in cache
	font, err := ParseFont(content)
	if err != nil {
		return nil, err
	}

	// Store in cache
	defaultCache.Set(content, font)

	return font, nil
}
