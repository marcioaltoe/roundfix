package runevent

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// bestEffortBuffer bounds each best-effort sink's queue; events beyond it
// are dropped rather than blocking the producer.
const bestEffortBuffer = 256

// Fanout multiplexes one publication across critical and best-effort sinks.
// Critical sink errors propagate to the producer. Best-effort sinks receive
// events through a bounded queue served by one goroutine per sink, so a slow
// or broken consumer never blocks or fails event production.
type Fanout struct {
	critical []Sink
	workers  []*bestEffortWorker

	mu     sync.Mutex
	closed bool
	wg     sync.WaitGroup

	dropped  atomic.Uint64
	failures atomic.Uint64
}

type bestEffortWorker struct {
	sink   Sink
	events chan RunEvent
}

// NewFanout starts one delivery goroutine per best-effort sink. Callers own
// the Fanout lifecycle and must call Close to stop delivery.
func NewFanout(critical []Sink, bestEffort []Sink) *Fanout {
	fanout := &Fanout{}
	for _, sink := range critical {
		if sink != nil {
			fanout.critical = append(fanout.critical, sink)
		}
	}
	for _, sink := range bestEffort {
		if sink == nil {
			continue
		}
		worker := &bestEffortWorker{sink: sink, events: make(chan RunEvent, bestEffortBuffer)}
		fanout.workers = append(fanout.workers, worker)
		fanout.wg.Add(1)
		go fanout.deliver(worker)
	}
	return fanout
}

func (fanout *Fanout) deliver(worker *bestEffortWorker) {
	defer fanout.wg.Done()
	for event := range worker.events {
		// Best-effort delivery is decoupled from producer cancellation.
		if err := worker.sink.Publish(context.Background(), event); err != nil {
			fanout.failures.Add(1)
		}
	}
}

// Publish sends the event to every sink. It returns the joined critical sink
// errors; best-effort failures are counted and swallowed.
func (fanout *Fanout) Publish(ctx context.Context, event RunEvent) error {
	var errs []error
	for _, sink := range fanout.critical {
		if err := sink.Publish(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}
	fanout.mu.Lock()
	for _, worker := range fanout.workers {
		if fanout.closed {
			fanout.dropped.Add(1)
			continue
		}
		select {
		case worker.events <- event:
		default:
			fanout.dropped.Add(1)
		}
	}
	fanout.mu.Unlock()
	return errors.Join(errs...)
}

// Close stops accepting best-effort deliveries, drains queued events, and
// waits for delivery goroutines to exit. It is safe to call more than once.
func (fanout *Fanout) Close() {
	fanout.mu.Lock()
	if fanout.closed {
		fanout.mu.Unlock()
		return
	}
	fanout.closed = true
	for _, worker := range fanout.workers {
		close(worker.events)
	}
	fanout.mu.Unlock()
	fanout.wg.Wait()
}

// DroppedEvents reports best-effort events discarded because a queue was
// full or the fanout was closed.
func (fanout *Fanout) DroppedEvents() uint64 { return fanout.dropped.Load() }

// BestEffortFailures reports swallowed best-effort publish errors.
func (fanout *Fanout) BestEffortFailures() uint64 { return fanout.failures.Load() }
