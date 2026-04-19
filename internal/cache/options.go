package cache

import "time"

// Option is a functional option for configuring a Cache.
type Option func(*Cache)

// WithTTL sets a custom TTL for the cache.
func WithTTL(d time.Duration) Option {
	return func(c *Cache) {
		c.ttl = d
	}
}

// NewWithOptions creates a Cache applying the provided options.
// Defaults to a 60-second TTL if none is specified.
func NewWithOptions(opts ...Option) *Cache {
	c := &Cache{
		entries: make(map[string]*Entry),
		ttl:     60 * time.Second,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}
