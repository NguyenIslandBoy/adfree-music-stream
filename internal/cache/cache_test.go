package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(5*time.Minute, 1*time.Minute)

	c.Set("foo", "bar", 5*time.Minute)

	val, ok := c.Get("foo")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if val.(string) != "bar" {
		t.Fatalf("expected 'bar', got '%v'", val)
	}

	_, ok = c.Get("missing")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}

	t.Logf("cache item count: %d", c.ItemCount())
}
