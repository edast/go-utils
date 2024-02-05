package stream

import "container/heap"

// Item holds the details of the queue item, including its value, priority, and index in the queue.
type Item[T any] struct {
	value    T
	priority int
	index    int
}

// PriorityQueue represents a priority queue where the highest priority items are retrieved first.
// It is implemented as a heap to efficiently support priority-based retrieval.
type PriorityQueue[T any] []*Item[T]

// NewPriorityQueue creates a new instance of an empty PriorityQueue.
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{}
}

// Len returns the number of elements in the priority queue. It is part of heap.Interface.
func (pq PriorityQueue[T]) Len() int {
	return len(pq)
}

// Less returns true if the item at index i has a higher priority than the item at index j. It is part of heap.Interface.
func (pq PriorityQueue[T]) Less(i, j int) bool {
	// Note: Higher value means higher priority.
	return pq[i].priority > pq[j].priority
}

// Swap swaps the elements at indexes i and j. It updates their indexes to reflect their new positions. It is part of heap.Interface.
func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Pop removes and returns the highest priority item from the queue. It is part of heap.Interface.
func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	if n == 0 {
		// Handle empty heap case. It might be better to return a zero value or handle this case outside the method.
		return nil // Or consider a more appropriate way to handle an empty heap.
	}
	item := old[n-1]
	item.index = -1 // Mark as removed
	*pq = old[0 : n-1]
	return item
}

// Push adds an item to the priority queue. It is part of heap.Interface.
func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	item.index = len(*pq)
	*pq = append(*pq, item)
}

// Update modifies the priority and value of an Item in the queue and adjusts the queue to maintain heap invariant.
func (pq *PriorityQueue[T]) Update(item *Item[T], value T, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// Ensure PriorityQueue implements heap.Interface at compile time.
var _ heap.Interface = (*PriorityQueue[string])(nil)
