package stream

import (
	"container/heap"
	"testing"
)

// TestPriorityQueue_New tests the creation of a new PriorityQueue.
func TestPriorityQueue_New(t *testing.T) {
	pq := NewPriorityQueue[int]()
	if pq == nil {
		t.Fatal("NewPriorityQueue[int]() returned nil")
	}
	if pq.Len() != 0 {
		t.Errorf("Expected length of 0, got %d", pq.Len())
	}
}

// TestPriorityQueue_PushPop tests pushing items to and popping items from the PriorityQueue.
func TestPriorityQueue_PushPop(t *testing.T) {
	pq := NewPriorityQueue[int]()
	heap.Init(pq) // Initialize the priority queue to establish heap invariants
	pq.Push(&Item[int]{value: 1, priority: 3})
	pq.Push(&Item[int]{value: 2, priority: 1})
	pq.Push(&Item[int]{value: 3, priority: 2})

	if pq.Len() != 3 {
		t.Fatalf("Expected length of 3, got %d", pq.Len())
	}

	// Pop items and check if they are in correct order based on priority.
	expectedPriorities := []int{3, 2, 1}
	for _, expectedPriority := range expectedPriorities {
		item := heap.Pop(pq).(*Item[int])
		if item.priority != expectedPriority {
			t.Errorf("Expected priority %d, got %d", expectedPriority, item.priority)
		}
	}

	if pq.Len() != 0 {
		t.Errorf("Expected empty queue after popping all items, got length %d", pq.Len())
	}
}

// TestPriorityQueue_Update tests updating the priority and value of an item.
func TestPriorityQueue_Update(t *testing.T) {
	pq := NewPriorityQueue[int]()
	heap.Init(pq)
	item := &Item[int]{value: 1, priority: 1}
	pq.Push(item)

	// Update the item's priority and value.
	pq.Update(item, 2, 10)

	updatedItem := heap.Pop(pq).(*Item[int])
	if updatedItem.value != 2 || updatedItem.priority != 10 {
		t.Errorf("Expected value of 2 and priority of 10, got value %d and priority %d", updatedItem.value, updatedItem.priority)
	}
}

// TestPriorityQueue_EmptyPop tests popping from an empty queue.
func TestPriorityQueue_EmptyPop(t *testing.T) {
	pq := NewPriorityQueue[int]()
	heap.Init(pq)
	item := pq.Pop()

	if item != nil {
		t.Error("Expected nil when popping from an empty queue")
	}
}

// TestPriorityQueue_MixedOperations tests a sequence of operations to ensure queue integrity.
func TestPriorityQueue_MixedOperations(t *testing.T) {
	pq := NewPriorityQueue[int]()
	heap.Init(pq)
	pq.Push(&Item[int]{value: 1, priority: 5})
	pq.Push(&Item[int]{value: 2, priority: 2})
	item := &Item[int]{value: 3, priority: 3}
	pq.Push(item)

	pq.Update(item, 3, 10) // Update item to have the highest priority

	if pq.Len() != 3 {
		t.Fatalf("Expected length of 3, got %d", pq.Len())
	}

	// First popped item should be the one with updated priority
	firstItem := heap.Pop(pq).(*Item[int])
	if firstItem.value != 3 || firstItem.priority != 10 {
		t.Errorf("Expected value of 3 and priority of 10 for first popped item, got value %d and priority %d", firstItem.value, firstItem.priority)
	}
}
