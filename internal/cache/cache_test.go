package cache_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/cache"
)

func TestSet_AndGet_ReturnsSecrets(t *testing.T) {
	c := cache.New(5 * time.Second)
	secrets := map[string]string{"KEY": "value"}
	c.Set("mypath", secrets)

	got, ok := c.Get("mypath")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", got["KEY"])
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	c := cache.New(5 * time.Second)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("path", map[string]string{"A": "1"})
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("path")
	if ok {
		t.Fatal("expected expired cache miss")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("path", map[string]string{"X": "y"})
	c.Invalidate("path")
	_, ok := c.Get("path")
	if ok {
		t.Fatal("expected cache miss after invalidate")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("a", map[string]string{"K": "v"})
	c.Set("b", map[string]string{"K": "v"})
	c.Flush()
	if c.Size() != 0 {
		t.Errorf("expected size 0, got %d", c.Size())
	}
}

func TestSize_ReturnsCount(t *testing.T) {
	c := cache.New(5 * time.Second)
	if c.Size() != 0 {
		t.Fatalf("expected 0, got %d", c.Size())
	}
	c.Set("p1", map[string]string{})
	c.Set("p2", map[string]string{})
	if c.Size() != 2 {
		t.Errorf("expected 2, got %d", c.Size())
	}
}
