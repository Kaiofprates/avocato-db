package ledger

import (
	"sync"
)

type Cache struct {
	mu       sync.RWMutex
	items    map[string]*Block
	capacity int
	order    []string
}

func NewCache(capacity int) *Cache {
	return &Cache{
		items:    make(map[string]*Block),
		capacity: capacity,
		order:    make([]string, 0, capacity),
	}
}

func (c *Cache) Get(key string) (*Block, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	b, ok := c.items[key]
	return b, ok
}

func (c *Cache) Put(key string, b *Block) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		c.items[key] = b
		return
	}

	if len(c.order) >= c.capacity {
		// Evict oldest (FIFO for simplicity in this POC)
		oldest := c.order[0]
		delete(c.items, oldest)
		c.order = c.order[1:]
	}

	c.items[key] = b
	c.order = append(c.order, key)
}
