package runevent

import "context"

type sourceFilterSink struct {
	next Sink
	drop Source
}

// NewSourceFilterSink returns a Sink decorator that drops events from the
// configured source and forwards every other event to the wrapped sink.
func NewSourceFilterSink(next Sink, drop Source) Sink {
	if next == nil {
		return Discard
	}
	return sourceFilterSink{next: next, drop: drop}
}

func (sink sourceFilterSink) Publish(ctx context.Context, event RunEvent) error {
	if event.Source == sink.drop {
		return nil
	}
	return sink.next.Publish(ctx, event)
}
