package gocache

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	K  string
	V  interface{}
	Ex int64
}

type Cache struct {
	items map[string]CacheEntry
	size  int
	mu    sync.RWMutex
}

func NewCacheWithMutex(size int) *Cache {
	c := Cache{}
	c.items = make(map[string]CacheEntry, size)
	c.size = size

	return &c
}

func (c *Cache) Set(k string, v interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we can insert new elements
	if len(c.items) == c.size {
		return fmt.Errorf("Cache size is full, delete elements to add more")
	}

	ex := time.Now().Add(ttl).UnixNano()
	c.items[k] = CacheEntry{
		K:  k,
		V:  v,
		Ex: ex,
	}

	// Delete item after ttl
	time.AfterFunc(ttl, func() {
		c.Delete(k)
	})

	return nil
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.items[k]
	if !ok {
		return nil, false
	}

	return val.V, true
}

func (c *Cache) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, k)
}

func (c *Cache) List() map[string]CacheEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.items
}
