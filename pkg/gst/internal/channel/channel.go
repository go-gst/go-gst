package channel

import (
	"sync"
)

// Growable is a channel-like structure that does not have a fixed size buffer
//
// Compared to a regular channel this also does not panic on a send if the channel is closed.
type Growable[T any] struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	buffer   []T
	closed   bool
}

func NewGrowable[T any]() *Growable[T] {
	gc := &Growable[T]{}
	gc.notEmpty = sync.NewCond(&gc.mu)
	return gc
}

// Send sends a value to the channel. It returns false if the channel is closed.
func (gc *Growable[T]) Send(val T) bool {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	if gc.closed {
		return false
	}
	gc.buffer = append(gc.buffer, val)
	gc.notEmpty.Signal()

	return true
}

// Receive receives a value from the channel. It returns the value and true if successful,
// or a zero value and false if the channel is closed.
func (gc *Growable[T]) Receive() (T, bool) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	for len(gc.buffer) == 0 && !gc.closed {
		gc.notEmpty.Wait()
	}

	var zero T
	if gc.closed {
		return zero, false
	}

	val := gc.buffer[0]
	gc.buffer = gc.buffer[1:]
	return val, true
}

func (gc *Growable[T]) Close() {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	gc.closed = true

	// drop all buffered values
	gc.buffer = nil

	gc.notEmpty.Broadcast()
}
