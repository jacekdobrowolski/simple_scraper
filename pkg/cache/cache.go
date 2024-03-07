package cache

import (
	"sync"
)

type Cache[T any] struct {
	cacheMap map[string]T
	mux      sync.RWMutex
}

func New[T any]() *Cache[T] {
	return &Cache[T]{cacheMap: make(map[string]T)}
}

func (c *Cache[T]) Set(url string, value T) {
	defer c.mux.Unlock()
	c.mux.Lock()
	c.cacheMap[url] = value
}

func (c *Cache[T]) Get(url string) (T, bool) {
	defer c.mux.RUnlock()
	c.mux.RLock()
	value, ok := c.cacheMap[url]
	return value, ok
}
