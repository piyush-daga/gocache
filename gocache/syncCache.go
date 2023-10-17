package gocache

import (
	"sync"
	"time"
)

type CacheItem struct {
	V  interface{}
	Ex int64
}

type CacheStore struct {
	store sync.Map
}

func NewCacheWithSyncMap() *CacheStore {
	return &CacheStore{}
}

func (c *CacheStore) Set(k string, v interface{}, ttl time.Duration) {
	ex := time.Now().Add(ttl).UnixNano()
	c.store.Store(k, CacheItem{v, ex})

	// Evict the entry when the time duration is complete
	time.AfterFunc(ttl, func() {
		c.store.Delete(k)
	})
}

func (c *CacheStore) Get(k string) (CacheItem, bool) {
	item, found := c.store.Load(k)
	if !found {
		return CacheItem{}, found
	}

	return CacheItem{item.(CacheItem).V, (item.(CacheItem).Ex - time.Now().UnixNano()) / int64(time.Millisecond)}, found
}

func (c *CacheStore) Delete(k string) {
	c.store.Delete(k)
}
