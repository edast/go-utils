# Generic Cache Implementations in Go

This repository provides thread-safe, generic implementations of various caching strategies, 
designed for high-performance concurrent access in Go. 
The initial offering includes an LRU (Least Recently Used) cache, with plans to expand 
to other strategies based on usage and demand.

## Features

- **Generic Implementation**: Works with any key and value types that are comparable in Go, thanks to Go's generics.
- **Concurrency Safe**: Designed to be safe for concurrent use by multiple goroutines without the need for external synchronization.
- **Efficient Operations**: Optimized for low memory overhead and fast operations, suitable for high-load environments.
- **Eviction Policies**: Each cache implementation comes with its own eviction policy, starting with LRU.

## Getting Started

### Installation

To include the cache package in your Go project, use the following `go get` command:

    go get -u github.com/edast/cache

### Usage

Here's a quick example of how to use the `LRUCache`:

    package main

    import (
        "fmt"

        "github.com/edast/cache"
    )

    func main() {
        // Create a new LRU cache with a capacity for 2 items
        lru := cache.NewLRUCache[string, int](2)

        // Add items to the cache
        lru.Put("Alice", 1)
        lru.Put("Bob", 2)

        // Retrieve and print an item
        if val, found := lru.Get("Alice"); found {
            fmt.Printf("Alice's value: %d\n", val)
        }

        // Add another item, causing `Bob` item to be evicted
        lru.Put("Charlie", 3)

        // Try to retrieve the evicted item
        if _, found := lru.Get("Bob"); !found {
            fmt.Println("Bob was evicted")
        }
    }


## Benchmarks

Benchmark tests are included to evaluate the performance of cache operations. Run the benchmarks using the Go tool:

    go test -bench=. github.com/edast/cache

## Contributing

Contributions are welcome, whether they're for new cache strategies, optimizations, bug reports, or documentation improvements. 
Please feel free to submit issues or pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
