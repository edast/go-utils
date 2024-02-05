// lru.go contains the implementation of the LRUCache type, applying a
// Least Recently Used (LRU) caching strategy. The LRUCache is safe for
// concurrent use by multiple goroutines, leveraging a mutex to protect
// shared state and a sync.Pool to minimize allocation overhead.

package cache

import (
	"container/list"
	"sync"
)

// entry holds a key-value pair for the cache. It is used internally by the LRUCache
// to store cache items in a linked list.
type entry[K comparable, V any] struct {
	key   K
	value V
}

// LRUCache implements a generic Least Recently Used (LRU) cache. It automatically
// evicts the least recently accessed items to maintain a fixed size. The cache is
// thread-safe, supporting concurrent access by multiple goroutines.
type LRUCache[K comparable, V any] struct {
	capacity int                 // Maximum number of items the cache can hold.
	list     *list.List          // Ordered list to track the least recently used items.
	dict     map[K]*list.Element // Map for quick access to list elements.
	pool     sync.Pool           // Pool to reuse entry objects.
	mu       sync.Mutex          // Mutex to protect concurrent access to the cache.
}

// NewLRUCache creates a new instance of an LRUCache with the given capacity.
// It initializes the internal data structures and prepares the cache for use.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("cache: capacity must be greater than zero")
	}

	return &LRUCache[K, V]{
		capacity: capacity,
		list:     list.New(),
		dict:     make(map[K]*list.Element, capacity),
		pool: sync.Pool{
			New: func() interface{} {
				return new(entry[K, V])
			},
		},
	}
}

// Get retrieves the value associated with the given key from the cache.
// If the key is found in the cache, Get returns the value and true.
// Otherwise, it returns the zero value for V and false.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.dict[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

// Put adds a key-value pair to the cache. If the key already exists, its value
// is updated. If adding a new key exceeds the cache's capacity, the least recently
// used item is evicted. Put is safe to call from multiple goroutines.
func (c *LRUCache[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.dict[key]; ok {
		elem.Value.(*entry[K, V]).value = val
		c.list.MoveToFront(elem)
		return
	}

	e := c.pool.Get().(*entry[K, V])
	e.key = key
	e.value = val

	if c.list.Len() >= c.capacity {
		c.evict()
	}

	elem := c.list.PushFront(e)
	c.dict[key] = elem
}

// evict removes the least recently used item from the cache.
// It is called internally by Put when adding a new item would exceed
// the cache's capacity. evict is safe to call from multiple goroutines.
func (c *LRUCache[K, V]) evict() {
	oldest := c.list.Back()
	if oldest != nil {
		oldEntry := oldest.Value.(*entry[K, V])
		delete(c.dict, oldEntry.key)
		c.list.Remove(oldest)
		c.pool.Put(oldEntry)
	}
}
