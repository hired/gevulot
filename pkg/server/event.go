package server

import (
	"sync"
	"sync/atomic"
)

// Event represents a one-time event that may occur in the future. Borrowed from gRPC.
type Event struct {
	doneChan chan struct{}
	fireOnce sync.Once
	fired    int32
}

// NewEvent returns a new, ready-to-use Event.
func NewEvent() *Event {
	return &Event{doneChan: make(chan struct{})}
}

// Fire causes event to complete. It is safe to call multiple times, and
// concurrently. It returns true if this call to Fire caused the signaling
// channel returned by Done to close.
func (e *Event) Fire() bool {
	ret := false

	e.fireOnce.Do(func() {
		atomic.StoreInt32(&e.fired, 1)
		close(e.doneChan)
		ret = true
	})

	return ret
}

// Done returns a channel that will be closed when Fire is called.
func (e *Event) Done() <-chan struct{} {
	return e.doneChan
}

// HasFired returns true if Fire has been called.
func (e *Event) HasFired() bool {
	return atomic.LoadInt32(&e.fired) == 1
}
