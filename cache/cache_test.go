package cache

import "testing"

func TestCachePutGet(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")

	value, ok := c.Get(1)

	if !ok {
		t.Fatal("expected cache hit")
	}

	if value != "one" {
		t.Fatalf("expected one, got %q", value)
	}
}

func TestCacheEviction(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")

	// 1 is now the least recently used item.
	c.Get(1)

	// 2 becomes the LRU item.
	c.Put(3, "three")

	if _, ok := c.Get(2); ok {
		t.Fatal("expected key 2 to be evicted")
	}

	if _, ok := c.Get(1); !ok {
		t.Fatal("expected key 1 to remain")
	}

	if _, ok := c.Get(3); !ok {
		t.Fatal("expected key 3 to remain")
	}
}

func TestCacheDelete(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Delete(1)

	if _, ok := c.Get(1); ok {
		t.Fatal("expected key to be deleted")
	}

	if c.Len() != 0 {
		t.Fatalf("expected empty cache, got %d", c.Len())
	}
}
