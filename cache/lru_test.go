package cache

import (
	"sync"
	"testing"
)

// TestLRUCache_PutGet tests basic put and get operations.
func TestLRUCache_PutGet(t *testing.T) {
	cache := NewLRUCache[string, string](2)

	// Test insertion
	cache.Put("key1", "val1")
	if v, ok := cache.Get("key1"); !ok || v != "val1" {
		t.Fatalf("cache.Get(\"key1\") = %v, %v; want %v, %v", v, ok, "val1", true)
	}

	// Test update
	cache.Put("key1", "val1-updated")
	if v, ok := cache.Get("key1"); !ok || v != "val1-updated" {
		t.Fatalf("cache.Get(\"key1\") after update = %v, %v; want %v, %v", v, ok, "val1-updated", true)
	}

	// Test eviction
	cache.Put("key2", "val2")
	cache.Put("key3", "val3") // This should evict "key1"
	if _, ok := cache.Get("key1"); ok {
		t.Fatal("Expected \"key1\" to be evicted")
	}
}

// TestLRUCache_EvictionOrder tests the LRU eviction policy.
func TestLRUCache_EvictionOrder(t *testing.T) {
	cache := NewLRUCache[int, int](2)

	cache.Put(1, 1)
	cache.Put(2, 2)
	cache.Put(3, 3) // Evicts key 1

	if _, ok := cache.Get(1); ok {
		t.Fatal("Expected key 1 to be evicted")
	}

	cache.Get(2)    // This access should make key 2 the most recently used
	cache.Put(4, 4) // Evicts key 3

	if _, ok := cache.Get(3); ok {
		t.Fatal("Expected key 3 to be evicted")
	}
}

// TestLRUCache_Concurrency tests the cache's thread-safety by performing parallel reads and writes.
func TestLRUCache_Concurrency(t *testing.T) {
	cache := NewLRUCache[int, int](100)
	var wg sync.WaitGroup

	// Populate the cache with initial values
	for i := 0; i < 50; i++ {
		cache.Put(i, i)
	}

	// Perform concurrent reads and writes
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(key int) {
			defer wg.Done()
			cache.Put(key, key*10)
		}(i)

		go func(key int) {
			defer wg.Done()
			cache.Get(key)
		}(i)
	}

	wg.Wait()

	// Verify updated values
	for i := 0; i < 50; i++ {
		if v, ok := cache.Get(i); !ok || v != i*10 {
			t.Fatalf("cache.Get(%d) = %v, %v; want %v, %v", i, v, ok, i*10, true)
		}
	}
}
