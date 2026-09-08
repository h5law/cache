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
		t.Fatalf("expected %q, got %q", "one", value)
	}
}

func TestCacheMiss(t *testing.T) {
	c := New[int, string](2)

	value, ok := c.Get(1)

	if ok {
		t.Fatal("expected cache miss")
	}

	if value != "" {
		t.Fatalf("expected zero value, got %q", value)
	}
}

func TestCacheEviction(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")

	// Access key 1, making key 2 the least recently used.
	c.Get(1)

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

func TestCacheEvictsLeastRecentlyUsed(t *testing.T) {
	c := New[int, string](3)

	c.Put(1, "one")
	c.Put(2, "two")
	c.Put(3, "three")

	// Access 1 and 2.
	// LRU order is now:
	// 2 -> 1 -> 3
	c.Get(1)
	c.Get(2)

	// Key 3 should be evicted.
	c.Put(4, "four")

	if _, ok := c.Get(3); ok {
		t.Fatal("expected least recently used key 3 to be evicted")
	}

	for _, key := range []int{1, 2, 4} {
		if _, ok := c.Get(key); !ok {
			t.Fatalf("expected key %d to remain", key)
		}
	}
}

func TestCacheGetUpdatesLRUOrder(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")

	// Key 1 becomes most recently used.
	c.Get(1)

	// Key 2 should now be evicted.
	c.Put(3, "three")

	if _, ok := c.Get(2); ok {
		t.Fatal("expected key 2 to be evicted")
	}

	if _, ok := c.Get(1); !ok {
		t.Fatal("expected key 1 to remain")
	}
}

func TestCachePutUpdatesLRUOrder(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")

	// Updating an existing key should make it most recently used.
	c.Put(1, "updated")

	c.Put(3, "three")

	if _, ok := c.Get(2); ok {
		t.Fatal("expected key 2 to be evicted")
	}

	value, ok := c.Get(1)
	if !ok {
		t.Fatal("expected key 1 to remain")
	}

	if value != "updated" {
		t.Fatalf("expected %q, got %q", "updated", value)
	}
}

func TestCacheUpdate(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(1, "updated")

	value, ok := c.Get(1)

	if !ok {
		t.Fatal("expected cache hit")
	}

	if value != "updated" {
		t.Fatalf("expected %q, got %q", "updated", value)
	}

	if c.Len() != 1 {
		t.Fatalf("expected length 1, got %d", c.Len())
	}
}

func TestCacheDelete(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Delete(1)

	if _, ok := c.Get(1); ok {
		t.Fatal("expected key 1 to be deleted")
	}

	if c.Len() != 0 {
		t.Fatalf("expected length 0, got %d", c.Len())
	}
}

func TestCacheDeleteMissingKey(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Delete(2)

	if c.Len() != 1 {
		t.Fatalf("expected length 1, got %d", c.Len())
	}

	if _, ok := c.Get(1); !ok {
		t.Fatal("expected key 1 to remain")
	}
}

func TestCacheDeleteUpdatesLRUOrder(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")

	c.Delete(1)
	c.Put(3, "three")

	// Key 2 should remain because key 1 was explicitly deleted.
	if _, ok := c.Get(2); !ok {
		t.Fatal("expected key 2 to remain")
	}

	if _, ok := c.Get(3); !ok {
		t.Fatal("expected key 3 to remain")
	}
}

func TestCacheCapacityOne(t *testing.T) {
	c := New[int, string](1)

	c.Put(1, "one")

	if c.Len() != 1 {
		t.Fatalf("expected length 1, got %d", c.Len())
	}

	c.Put(2, "two")

	if c.Len() != 1 {
		t.Fatalf("expected length 1, got %d", c.Len())
	}

	if _, ok := c.Get(1); ok {
		t.Fatal("expected key 1 to be evicted")
	}

	value, ok := c.Get(2)
	if !ok {
		t.Fatal("expected key 2 to remain")
	}

	if value != "two" {
		t.Fatalf("expected %q, got %q", "two", value)
	}
}

func TestCacheLen(t *testing.T) {
	c := New[int, string](3)

	if c.Len() != 0 {
		t.Fatalf("expected initial length 0, got %d", c.Len())
	}

	c.Put(1, "one")

	if c.Len() != 1 {
		t.Fatalf("expected length 1, got %d", c.Len())
	}

	c.Put(2, "two")

	if c.Len() != 2 {
		t.Fatalf("expected length 2, got %d", c.Len())
	}

	c.Put(3, "three")

	if c.Len() != 3 {
		t.Fatalf("expected length 3, got %d", c.Len())
	}

	c.Put(4, "four")

	if c.Len() != 3 {
		t.Fatalf("expected length 3 after eviction, got %d", c.Len())
	}

	c.Delete(2)

	if c.Len() != 2 {
		t.Fatalf("expected length 2 after deletion, got %d", c.Len())
	}
}

func TestCacheEvictionMaintainsCapacity(t *testing.T) {
	const capacity = 10

	c := New[int, int](capacity)

	for i := 0; i < 1_000; i++ {
		c.Put(i, i)

		if c.Len() > capacity {
			t.Fatalf(
				"cache exceeded capacity: got %d, capacity %d",
				c.Len(),
				capacity,
			)
		}
	}

	if c.Len() != capacity {
		t.Fatalf(
			"expected final length %d, got %d",
			capacity,
			c.Len(),
		)
	}
}

func TestCacheDeleteAll(t *testing.T) {
	c := New[int, string](3)

	c.Put(1, "one")
	c.Put(2, "two")
	c.Put(3, "three")

	c.Delete(1)
	c.Delete(2)
	c.Delete(3)

	if c.Len() != 0 {
		t.Fatalf("expected empty cache, got length %d", c.Len())
	}

	for _, key := range []int{1, 2, 3} {
		if _, ok := c.Get(key); ok {
			t.Fatalf("expected key %d to be absent", key)
		}
	}
}

func TestCacheReuseAfterDelete(t *testing.T) {
	c := New[int, string](2)

	c.Put(1, "one")
	c.Delete(1)
	c.Put(1, "new one")

	value, ok := c.Get(1)

	if !ok {
		t.Fatal("expected key 1 to be present")
	}

	if value != "new one" {
		t.Fatalf("expected %q, got %q", "new one", value)
	}

	if c.Len() != 1 {
		t.Fatalf("expected length 1, got %d", c.Len())
	}
}

func TestCacheZeroValue(t *testing.T) {
	type value struct {
		ID   int
		Name string
	}

	c := New[int, value](2)

	c.Put(1, value{})

	got, ok := c.Get(1)

	if !ok {
		t.Fatal("expected cache hit")
	}

	if got != (value{}) {
		t.Fatalf("expected zero value, got %+v", got)
	}
}

func TestCacheGenericTypes(t *testing.T) {
	c := New[string, int](2)

	c.Put("one", 1)
	c.Put("two", 2)

	value, ok := c.Get("one")

	if !ok {
		t.Fatal("expected cache hit")
	}

	if value != 1 {
		t.Fatalf("expected 1, got %d", value)
	}
}

func TestCacheInvalidCapacity(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{
			name:     "zero",
			capacity: 0,
		},
		{
			name:     "negative",
			capacity: -1,
		},
		{
			name:     "large negative",
			capacity: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected New to panic")
				}
			}()

			New[int, string](tt.capacity)
		})
	}
}
