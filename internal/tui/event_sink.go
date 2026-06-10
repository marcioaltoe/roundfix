package tui

import (
	"context"
	"sync"
	"sync/atomic"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
)

// agentEventBufferSize bounds the Live Run View's event queue; rendering
// pressure beyond it drops events instead of blocking producers.
const agentEventBufferSize = 512

// EventBuffer adapts Run Event publication to rendering as a best-effort,
// non-blocking consumer: events queue into a bounded buffer served by one
// delivery goroutine, and overflow is counted, never blocking.
type EventBuffer struct {
	updates chan agent.StreamUpdate
	deliver func(agent.StreamUpdate)
	done    chan struct{}
	drops   atomic.Uint64

	mu     sync.Mutex
	closed bool
}

// NewEventBuffer starts the delivery goroutine. Callers own the lifecycle
// and must call Close to stop it.
func NewEventBuffer(capacity int, deliver func(agent.StreamUpdate)) *EventBuffer {
	if capacity <= 0 {
		capacity = agentEventBufferSize
	}
	buffer := &EventBuffer{
		updates: make(chan agent.StreamUpdate, capacity),
		deliver: deliver,
		done:    make(chan struct{}),
	}
	go buffer.drain()
	return buffer
}

func (buffer *EventBuffer) drain() {
	defer close(buffer.done)
	for update := range buffer.updates {
		buffer.deliver(update)
	}
}

// Publish enqueues the event without ever blocking: a full buffer or a
// closed view counts a drop and reports success, per the best-effort policy.
func (buffer *EventBuffer) Publish(_ context.Context, event runevent.RunEvent) error {
	update, ok := agent.StreamUpdateFromEvent(event)
	if !ok {
		return nil
	}
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	if buffer.closed {
		buffer.drops.Add(1)
		return nil
	}
	select {
	case buffer.updates <- update:
	default:
		buffer.drops.Add(1)
	}
	return nil
}

// Close stops accepting events, drains the queue, and waits for the
// delivery goroutine to exit. It is safe to call more than once.
func (buffer *EventBuffer) Close() {
	buffer.mu.Lock()
	if buffer.closed {
		buffer.mu.Unlock()
		return
	}
	buffer.closed = true
	close(buffer.updates)
	buffer.mu.Unlock()
	<-buffer.done
}

// DroppedEvents reports events discarded under rendering pressure.
func (buffer *EventBuffer) DroppedEvents() uint64 {
	return buffer.drops.Load()
}
