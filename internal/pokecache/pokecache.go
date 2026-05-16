package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c
}

func (c *Cache) reapLoop() {
	// create a ticker
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// do forever
	for range ticker.C {
		c.mu.Lock()

		currentTime := time.Now()

		for key, entry := range c.cache {
			// If entry was added longer ago than the interval allowed
			if currentTime.Sub(entry.createdAt) > c.interval {
				delete(c.cache, key)
			}
		}

		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, v []byte) {
	entry := cacheEntry{createdAt: time.Now(), val: v}
	c.cache[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	v, ok := c.cache[key]
	return v.val, ok
}
