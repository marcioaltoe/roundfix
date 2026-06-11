package runevent

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type recordingSink struct {
	mu     sync.Mutex
	events []RunEvent
	err    error
}

func (sink *recordingSink) Publish(_ context.Context, event RunEvent) error {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	sink.events = append(sink.events, event)
	return sink.err
}

func (sink *recordingSink) count() int {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	return len(sink.events)
}

type blockingSink struct {
	release chan struct{}
}

func (sink *blockingSink) Publish(context.Context, RunEvent) error {
	<-sink.release
	return nil
}

func TestFanoutPropagatesCriticalSinkError(t *testing.T) {
	critical := &recordingSink{err: errors.New("journal write failed")}
	fanout := NewFanout([]Sink{critical}, nil)
	defer fanout.Close()

	err := fanout.Publish(context.Background(), RunEvent{Kind: KindAgentMessage})

	if err == nil || !errors.Is(err, critical.err) {
		t.Fatalf("expected critical sink error to propagate, got %v", err)
	}
}

func TestFanoutSwallowsAndCountsBestEffortErrors(t *testing.T) {
	broken := &recordingSink{err: errors.New("ui broke")}
	fanout := NewFanout(nil, []Sink{broken})

	if err := fanout.Publish(context.Background(), RunEvent{Kind: KindAgentMessage}); err != nil {
		t.Fatalf("expected best-effort error swallowed, got %v", err)
	}
	fanout.Close()

	if broken.count() != 1 {
		t.Fatalf("expected delivery attempt, got %d", broken.count())
	}
	if fanout.BestEffortFailures() != 1 {
		t.Fatalf("expected one recorded failure, got %d", fanout.BestEffortFailures())
	}
}

func TestFanoutBlockedBestEffortSinkNeverStallsPublication(t *testing.T) {
	blocked := &blockingSink{release: make(chan struct{})}
	healthy := &recordingSink{}
	fanout := NewFanout([]Sink{healthy}, []Sink{blocked})

	published := make(chan struct{})
	go func() {
		for index := 0; index < bestEffortBuffer+10; index++ {
			if err := fanout.Publish(context.Background(), RunEvent{Kind: KindAgentMessage}); err != nil {
				t.Errorf("publish: %v", err)
			}
		}
		close(published)
	}()

	select {
	case <-published:
	case <-time.After(5 * time.Second):
		t.Fatal("publication stalled behind a blocked best-effort sink")
	}
	if healthy.count() != bestEffortBuffer+10 {
		t.Fatalf("expected critical sink to receive every event, got %d", healthy.count())
	}
	if fanout.DroppedEvents() == 0 {
		t.Fatal("expected drops counted while the best-effort sink was blocked")
	}
	close(blocked.release)
	fanout.Close()
}

func TestFanoutDeliversEverythingUnderConcurrentPublishers(t *testing.T) {
	var received atomic.Int64
	counter := sinkFunc(func(context.Context, RunEvent) error {
		received.Add(1)
		return nil
	})
	critical := &recordingSink{}
	fanout := NewFanout([]Sink{critical}, []Sink{counter})

	const publishers = 8
	const perPublisher = 50
	var wg sync.WaitGroup
	for index := 0; index < publishers; index++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < perPublisher; n++ {
				_ = fanout.Publish(context.Background(), RunEvent{Kind: KindAgentMessage})
			}
		}()
	}
	wg.Wait()
	fanout.Close()

	total := int64(publishers * perPublisher)
	if critical.count() != int(total) {
		t.Fatalf("expected %d critical deliveries, got %d", total, critical.count())
	}
	if received.Load()+int64(fanout.DroppedEvents()) != total {
		t.Fatalf("expected delivered+dropped=%d, got %d+%d", total, received.Load(), fanout.DroppedEvents())
	}
}

func TestFanoutCloseIsIdempotentAndDropsLatePublishes(t *testing.T) {
	sink := &recordingSink{}
	fanout := NewFanout(nil, []Sink{sink})
	fanout.Close()
	fanout.Close()

	if err := fanout.Publish(context.Background(), RunEvent{Kind: KindAgentMessage}); err != nil {
		t.Fatalf("expected publish after close to stay silent, got %v", err)
	}
	if fanout.DroppedEvents() != 1 {
		t.Fatalf("expected late publish counted as dropped, got %d", fanout.DroppedEvents())
	}
}

type sinkFunc func(ctx context.Context, event RunEvent) error

func (fn sinkFunc) Publish(ctx context.Context, event RunEvent) error { return fn(ctx, event) }
