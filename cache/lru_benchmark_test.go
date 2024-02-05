package cache

import (
	"strconv"
	"sync/atomic"
	"testing"
)

// BenchmarkLRUCache_Put benchmarks the performance of the Put operation.
func BenchmarkLRUCache_Put(b *testing.B) {
	cache := NewLRUCache[int, string](b.N) // Use b.N as the capacity to avoid evictions

	b.ResetTimer() // Reset the timer to exclude setup time from the benchmark

	for i := 0; i < b.N; i++ {
		cache.Put(i, strconv.Itoa(i))
	}
}

// BenchmarkLRUCache_Get benchmarks the performance of the Get operation.
func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache[int, string](b.N)
	for i := 0; i < b.N; i++ {
		cache.Put(i, strconv.Itoa(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(i)
	}
}

// BenchmarkLRUCache_PutGet benchmarks the combined performance of Put and Get operations.
func BenchmarkLRUCache_PutGet(b *testing.B) {
	cache := NewLRUCache[int, string](b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Put(i, strconv.Itoa(i))
		cache.Get(i)
	}
}

// BenchmarkLRUCache_Concurrent benchmarks the performance of the cache under concurrent access.
func BenchmarkLRUCache_Concurrent(b *testing.B) {
	cache := NewLRUCache[int, string](1000) // Fixed size to force some evictions

	var keyCounter int64 // Atomic counter to generate unique keys

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var localKey int64
		for pb.Next() {
			// Use atomic.AddInt64 to ensure unique key generation across goroutines
			localKey = atomic.AddInt64(&keyCounter, 1)
			key := int(localKey) % 1000 // Keep keys within a realistic range to ensure some cache hits
			val := strconv.Itoa(key)
			cache.Put(key, val)
			_, _ = cache.Get(key)
		}
	})
}
