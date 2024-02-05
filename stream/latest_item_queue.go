package stream

// LatestItemQueue is a generic type-safe queue that ensures the consumer always receives the most recent item.
// It is particularly useful in scenarios where processing speed varies and only the latest data is relevant,
// such as real-time data processing or event handling systems.
type LatestItemQueue[T any] struct {
	channel chan T        // A channel that holds the latest item.
	closed  chan struct{} // Indicator for closing the consume channel
}

// NewLatestItemQueue creates a new instance of LatestItemQueue with a predefined buffer.
// The buffer size is set to 1 to hold only the most recent item.
func NewLatestItemQueue[T any]() *LatestItemQueue[T] {
	return &LatestItemQueue[T]{
		channel: make(chan T, 1),
		closed:  make(chan struct{}),
	}
}

// Produce attempts to send an item to the queue.
// If the queue is full (already holding an item), it discards the oldest item and enqueues the new one,
// ensuring that the queue always contains the most recent item.
func (q *LatestItemQueue[T]) Produce(item T) {
	select {
	case <-q.closed: // Check if the queue is closed to prevent sending on closed channel
		return
	case q.channel <- item: // Try sending the item to the channel.
	default:
		// The channel is full, discard the oldest item and send the new one.
		<-q.channel       // Discard the oldest item.
		q.channel <- item // Send the new item.
	}
}

// ConsumeChannel provides access to the underlying channel for consuming items.
// Consumers can read from this channel to receive the most recent item available.
func (q *LatestItemQueue[T]) ConsumeChannel() <-chan T {
	return q.channel
}

// Close safely closes the consume channel, ensuring no more items can be sent.
func (q *LatestItemQueue[T]) Close() {
	select {
	case <-q.closed: // Prevent closing more than once
		return
	default:
		close(q.closed)  // Close the closed channel to signal closure
		close(q.channel) // Close the channel to signal no more sends
	}
}
