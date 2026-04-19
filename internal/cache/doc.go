// Package cache implements a thread-safe, TTL-based in-memory cache for
// Vault secret maps. It is used by the syncer and watcher to avoid
// redundant reads from Vault when secrets are unlikely to have changed
// within the configured time-to-live window.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	c.Set("secret/myapp", secrets)
//	if cached, ok := c.Get("secret/myapp"); ok {
//	    // use cached secrets
//	}
package cache
