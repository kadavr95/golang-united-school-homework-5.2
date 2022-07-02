package cache

import (
	"time"
)

type CacheEntry struct {
	value      string
	expiration time.Time
}

type Cache struct {
	entries map[string]CacheEntry
}

func NewCache() Cache {
	return Cache{entries: make(map[string]CacheEntry)}
}

func (c Cache) Get(key string) (string, bool) {
	v, ok := c.entries[key]
	if v.expiration.IsZero() || time.Now().Before(v.expiration) {
		return v.value, ok
	}
	return "", false
}

func (c *Cache) Put(key, value string) {
	c.entries[key] = CacheEntry{value: value}
}

func (c Cache) Keys() []string {
	keys := make([]string, 0)
	for k := range c.entries {
		v := c.entries[k]
		if v.expiration.IsZero() || time.Now().Before(v.expiration) {
			keys = append(keys, k)
		}
	}
	// Some clean up to prevent bloating of a cache with expired values
	if float32(len(keys)) > float32(len(c.entries))*0.3 {
		c.CleanUp()
	}
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.entries[key] = CacheEntry{value: value, expiration: deadline}
}

func (c *Cache) CleanUp() {
	for k := range c.entries {
		v := c.entries[k]
		if !v.expiration.IsZero() && time.Now().After(v.expiration) {
			delete(c.entries, k)
		}
	}
}
