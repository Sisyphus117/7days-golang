package cache

import (
	lru "cache/rlu"
	"sync"
)

type cache struct {
	mu    sync.Mutex
	lru   *lru.LRU
	bytes int64
}
