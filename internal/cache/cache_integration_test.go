package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/cache"
)

// TestCache_ConcurrentAccess ensures no data races under parallel reads/writes.
func TestCache_ConcurrentAccess(t *testing.T) {
	c := cache.New(5 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			key := "path"
			c.Set(key, map[string]string{"K": "v"})
		}(i)
		go func() {
			defer wg.Done()
			c.Get("path")
		}()
	}
	wg.Wait()

	if c.Size() < 0 {
		t.Error("unexpected negative size")
	}
}
