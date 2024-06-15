package cache

import (
	"sync"
	"time"
)

type Cache[I comparable, T any] interface {
	Put(time.Duration, I, T)
	Get(I) (T, bool)
	Remove(I) (T, bool)
	Clean()
}

type cacheItem[T any] struct {
	t time.Time
	v T
}

type cache[I comparable, T any] struct {
	mp map[I]cacheItem[T]

	mu sync.RWMutex
}

func New[I comparable, T any]() Cache[I, T] {
	return &cache[I, T]{
		mp: make(map[I]cacheItem[T]),
	}
}

func (c *cache[I, T]) Put(d time.Duration, i I, v T) {
	c.mu.Lock()
	c.mp[i] = cacheItem[T]{
		t: time.Now().Add(d),
		v: v,
	}
	c.mu.Unlock()
}

func (c *cache[I, T]) Get(i I) (T, bool) {
	var zero T

	c.mu.RLock()
	v, ok := c.mp[i]
	c.mu.RUnlock()

	if !ok {
		return zero, false
	}

	if v.t.Before(time.Now()) {
		c.mu.Lock()
		delete(c.mp, i)
		c.mu.Unlock()

		return zero, false
	}

	return v.v, true
}

func (c *cache[I, T]) Remove(i I) (T, bool) {
	var zero T

	c.mu.Lock()
	v, ok := c.mp[i]
	if ok {
		delete(c.mp, i)
	}
	c.mu.Unlock()

	if !ok {
		return zero, false
	}

	if v.t.Before(time.Now()) {
		return zero, false
	}

	return v.v, true
}

func (c *cache[I, T]) Clean() {
	var now = time.Now()

	c.mu.Lock()
	for k, v := range c.mp {
		if v.t.Before(now) {
			delete(c.mp, k)
		}
	}
	c.mu.Unlock()
}
