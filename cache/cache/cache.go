package cache

import (
	lru "cache/cache/rlu"
	"fmt"
	"sync"
)

type cache struct {
	mu    sync.Mutex
	lru   *lru.LRU
	bytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.NewLRU(c.bytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return ByteView{}, fmt.Errorf("mismatch")
	}

	value, ok := c.lru.Get(key)
	if ok {
		return value.(ByteView), nil
	}
	return NewByteView(nil), fmt.Errorf("mismatch")
}
