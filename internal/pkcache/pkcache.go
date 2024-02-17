package pkcache

import (
	"sync"
	"time"
)

type Cache struct {
	val map[string]cachedEntry
	mu  *sync.Mutex
}

type cachedEntry struct {
	createdAt time.Time
	val       []byte
}

type Ticker struct {
	C <-chan time.Time
}

func NewCache(dur time.Duration) Cache {
	c := Cache{
		val: make(map[string]cachedEntry),
		mu:  &sync.Mutex{},
	}
	go c.reapLoop(dur)
	return c
}

func (c *Cache) Add(url string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	newEntry := cachedEntry{time.Now(), val}
	c.val[url] = newEntry
}

func (c *Cache) Get(url string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.val[url]
	return entry.val, ok
}

func (c *Cache) reapLoop(dur time.Duration) {
	ticker := time.NewTicker(dur)
	for range ticker.C {
		c.reap(time.Now().UTC(), dur)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.val {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.val, k)
		}
	}
}
