// Package cache provides a simple TTL-based in-memory cache for Vault secrets
// to reduce redundant API calls during watch intervals.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached set of secrets and its expiry time.
type Entry struct {
	Secrets   map[string]string
	FetchedAt time.Time
	ExpiresAt time.Time
}

// Cache is a TTL-based in-memory store for secret maps.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a new Cache with the given TTL duration.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Set stores secrets under the given key.
func (c *Cache) Set(key string, secrets map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.entries[key] = &Entry{
		Secrets:   secrets,
		FetchedAt: now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get retrieves secrets for the given key if present and not expired.
// Returns nil, false if missing or stale.
func (c *Cache) Get(key string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		return nil, false
	}
	return e.Secrets, true
}

// Invalidate removes the entry for the given key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*Entry)
}

// Size returns the number of entries currently in the cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
