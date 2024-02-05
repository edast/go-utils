package stream

import (
	"testing"
	"time"
)

func TestLatestItemQueue_ProduceConsume(t *testing.T) {
	queue := NewLatestItemQueue[int]()

	// Test producing and consuming a single item
	expected := 42
	queue.Produce(expected)
	item := <-queue.ConsumeChannel()
	if item != expected {
		t.Errorf("Expected %v, got %v", expected, item)
	}

	// Test overwriting an item and consuming the latest one
	newExpected := 43
	queue.Produce(newExpected)     // This item might get overwritten
	queue.Produce(newExpected + 1) // Latest item

	item = <-queue.ConsumeChannel()
	if item != newExpected+1 {
		t.Errorf("Expected %v, got %v", newExpected+1, item)
	}
}

func TestLatestItemQueue_EmptyConsume(t *testing.T) {
	queue := NewLatestItemQueue[int]()
	done := make(chan struct{})

	go func() {
		select {
		case item, ok := <-queue.ConsumeChannel():
			if ok { // This means the channel was not closed and we received an item, which shouldn't happen.
				t.Errorf("Received item %v from empty queue", item)
			}
		case <-time.After(100 * time.Millisecond): // Adjust the duration as needed
			// If no item is received in the given duration, assume success
			close(done)
		}
	}()

	<-done // Wait for the goroutine to finish or timeout
}

func TestLatestItemQueue_ProduceConsume2(t *testing.T) {
	queue := NewLatestItemQueue[int]()
	defer queue.Close()

	expected := 42
	queue.Produce(expected)

	select {
	case item := <-queue.ConsumeChannel():
		if item != expected {
			t.Errorf("Expected %d, got %d", expected, item)
		}
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for item")
	}
}

func TestLatestItemQueue_Overwrite(t *testing.T) {
	queue := NewLatestItemQueue[int]()
	defer queue.Close()

	queue.Produce(1) // This item should be overwritten
	expected := 2
	queue.Produce(expected)

	select {
	case item := <-queue.ConsumeChannel():
		if item != expected {
			t.Errorf("Expected %d, got %d", expected, item)
		}
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for item")
	}
}

func TestLatestItemQueue_Close(t *testing.T) {
	queue := NewLatestItemQueue[int]()
	queue.Produce(1) // Ensure there's something to potentially read.
	queue.Close()    // Close the queue, which should close the channel.

	// Attempt to receive from the channel; this might fetch the last item.
	<-queue.ConsumeChannel()
	// Attempt to receive again; this time, we expect it to be closed.
	_, ok := <-queue.ConsumeChannel()
	if ok {
		t.Fatal("Expected channel to be closed, but it was still open")
	}
}

func TestLatestItemQueue_ConcurrentProduce(t *testing.T) {
	queue := NewLatestItemQueue[int]()
	defer queue.Close()

	// Start multiple producers
	for i := 0; i < 10; i++ {
		go func(val int) {
			queue.Produce(val)
		}(i)
	}

	time.Sleep(time.Millisecond * 100) // Allow some time for goroutines to run

	// Ensure the queue is not blocked and can still receive
	queue.Produce(42)
	select {
	case item := <-queue.ConsumeChannel():
		if item != 42 {
			t.Errorf("Expected 42, got %d", item)
		}
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for item")
	}
}
