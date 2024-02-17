package cache

import (
	"sync"
	"time"
)

type Cache struct {
	val      map[string]cachedEntry
	mu       sync.Mutex
	interval time.Duration
}

type cachedEntry struct {
	createdAt time.Time
	val       []byte
}

type Ticker struct {
	C <-chan time.Time
}

func (c *Cache) NewCache(dur time.Duration) {
	c.val = make(map[string]cachedEntry)
	c.interval = dur
	t := time.NewTicker(dur)
	go c.reapLoop(t)
}

func (c *Cache) reapLoop(t *time.Ticker) {
	for {
		select {
		case <-t.C:
			c.mu.Lock()
			now := time.Now()
			for entry := range c.val {
				if c.interval < time.Time.Sub(now, c.val[entry].createdAt) {
					delete(c.val, entry)
				}
			}
			c.mu.Unlock()
		default:
			continue
		}
	}
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
	if !ok {
		return make([]byte, 0), false
	}
	return entry.val, true
}
