package cache

import "sync"

type Config struct {
	Space int
}

type Cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewCache[K comparable, V any](cfg Config) *Cache[K, V] {
	return &Cache[K, V]{
		mu:   sync.RWMutex{},
		data: make(map[K]V, cfg.Space),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()

	c.data[key] = value

	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()

	item, ok := c.data[key]

	c.mu.RUnlock()

	if !ok {
		var zero V

		return zero, false
	}

	return item, true
}

func (c *Cache[K, V]) GetAll() map[K]V {
	// Требований к сторогой консистентности нет, поэтому для оптимизации использовал RLock()
	c.mu.RLock()

	newMap := make(map[K]V, len(c.data))
	for k, v := range c.data {
		newMap[k] = v
	}

	c.mu.RUnlock()

	return newMap
}
