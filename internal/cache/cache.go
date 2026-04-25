package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	store *gocache.Cache
}

func New(defaultTTL, cleanupInterval time.Duration) *Cache {
	return &Cache{
		store: gocache.New(defaultTTL, cleanupInterval),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	return c.store.Get(key)
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.store.Set(key, value, ttl)
}

func (c *Cache) ItemCount() int {
	return c.store.ItemCount()
}
