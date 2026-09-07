package cache

import "container/list"

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	lru      *list.List
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("cache capacity must be positive")
	}

	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element, capacity),
		lru:      list.New(),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	element, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	c.lru.MoveToFront(element)

	return element.Value.(entry[K, V]).value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	if element, ok := c.items[key]; ok {
		element.Value = entry[K, V]{
			key:   key,
			value: value,
		}

		c.lru.MoveToFront(element)
		return
	}

	element := c.lru.PushFront(entry[K, V]{
		key:   key,
		value: value,
	})

	c.items[key] = element

	if c.lru.Len() > c.capacity {
		c.evict()
	}
}

func (c *Cache[K, V]) Delete(key K) {
	element, ok := c.items[key]
	if !ok {
		return
	}

	delete(c.items, key)
	c.lru.Remove(element)
}

func (c *Cache[K, V]) Len() int {
	return c.lru.Len()
}

func (c *Cache[K, V]) evict() {
	element := c.lru.Back()
	if element == nil {
		return
	}

	item := element.Value.(entry[K, V])

	delete(c.items, item.key)
	c.lru.Remove(element)
}
