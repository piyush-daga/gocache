package gocache

import (
	"math"
	"reflect"
	"sync"
	"time"
)

const mutexLocked = 1

func MutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func RWMutexWriteLocked(rw *sync.RWMutex) bool {
	// RWMutex has a "w" sync.Mutex field for write lock
	state := reflect.ValueOf(rw).Elem().FieldByName("w").FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func RWMutexReadLocked(rw *sync.RWMutex) bool {
	state := reflect.ValueOf(rw).Elem().FieldByName("readerCount").FieldByName("v")
	return state.Int() > 0
}

type CacheEntryWithEviction struct {
	CacheEntry
	LastRecentlyAccessed time.Time
}

type CacheStoreWithEviction struct {
	items map[string]CacheEntryWithEviction
	size  int
	mu    sync.Mutex
}

// Default policy - LRU
func NewCacheWithEvictionPolicy(size int) *CacheStoreWithEviction {
	return &CacheStoreWithEviction{
		items: make(map[string]CacheEntryWithEviction, size),
		size:  size,
	}
}

func (c *CacheStoreWithEviction) Set(k string, v interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) == c.size {
		c.DeleteOldestAccessedItem()
	}

	ex := time.Now().Add(ttl).UnixNano()
	c.items[k] = CacheEntryWithEviction{
		CacheEntry{
			K:  k,
			V:  v,
			Ex: ex,
		},
		time.Now(),
	}

	time.AfterFunc(ttl, func() {
		c.Delete(k)
	})
}

func (c *CacheStoreWithEviction) Get(k string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.items[k]
	if !ok {
		return nil, false
	}

	// Access has succeeded - weird and counterintutive
	val.LastRecentlyAccessed = time.Now()
	c.items[k] = val

	return val, true
}

func (c *CacheStoreWithEviction) List() map[string]CacheEntryWithEviction {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		v.LastRecentlyAccessed = time.Now()
		c.items[k] = v
	}

	return c.items
}

// Since, this may or may not be called with a lock held
func deleteWithoutLocking(items map[string]CacheEntryWithEviction, k string) {
	delete(items, k)
}

func (c *CacheStoreWithEviction) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	deleteWithoutLocking(c.items, k)
}

func (c *CacheStoreWithEviction) DeleteOldestAccessedItem() {
	// Assuming this is always called with a lock held

	max := math.MinInt64
	oldest := ""
	for k, v := range c.items {
		if v.LastRecentlyAccessed.Nanosecond() > max {
			max = v.LastRecentlyAccessed.Nanosecond()
			oldest = k
		}
	}

	deleteWithoutLocking(c.items, oldest)
}
