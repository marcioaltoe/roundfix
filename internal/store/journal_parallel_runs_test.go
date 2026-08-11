package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

// Suite: parallel Run Event Journal writers
// Invariant: six Runs sharing one Run Database append without SQLITE_BUSY while preserving each Run's cursor and publisher order.
// Boundary IN: independent Store connections, SQLite, journal batching, and the machine-wide advisory writer lock.
// Boundary OUT: Run scheduling and consumer projections, which have their own package and corpus suites.

const (
	parallelRunsPreRaiseBusyTimeoutMillis = 5_000
	parallelRunsCount                     = 6
	parallelRunEvents                     = journalBatchSize * 2
)

type parallelRunFixture struct {
	store     *Store
	runID     string
	summaries []string
}

type parallelRunWriteResult struct {
	runIndex int
	attempts int
	latency  time.Duration
	err      error
}

type parallelRunMeasurement struct {
	completedRuns         int
	totalEvents           int
	sqliteBusy            int
	elapsed               time.Duration
	writerLatencyP50      time.Duration
	writerLatencyP95      time.Duration
	concurrentPerThousand time.Duration
}

type parallelRunScenario struct {
	runs        []parallelRunFixture
	measurement parallelRunMeasurement
}

func TestParallelRuns(t *testing.T) {
	t.Run("six concurrent Runs append events at the pre-raise timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		defer cancel()

		scenario := runParallelRunScenario(t, ctx)
		measurement := scenario.measurement
		if measurement.sqliteBusy != 0 {
			t.Fatalf("expected zero SQLITE_BUSY occurrences, got %d", measurement.sqliteBusy)
		}
		if measurement.completedRuns != parallelRunsCount {
			t.Fatalf("expected all %d Runs to complete appends, got %d", parallelRunsCount, measurement.completedRuns)
		}
		if measurement.totalEvents != parallelRunsCount*parallelRunEvents {
			t.Fatalf("expected %d appended events, got %d", parallelRunsCount*parallelRunEvents, measurement.totalEvents)
		}
		t.Logf(
			"PARALLEL_RUNS runs=%d events_per_run=%d total_events=%d busy_timeout_ms=%d sqlite_busy=%d completed_runs=%d elapsed_us=%d writer_latency_p50_us=%d writer_latency_p95_us=%d concurrent_wall_per_1000_events_us=%d",
			parallelRunsCount,
			parallelRunEvents,
			measurement.totalEvents,
			parallelRunsPreRaiseBusyTimeoutMillis,
			measurement.sqliteBusy,
			measurement.completedRuns,
			measurement.elapsed.Microseconds(),
			measurement.writerLatencyP50.Microseconds(),
			measurement.writerLatencyP95.Microseconds(),
			measurement.concurrentPerThousand.Microseconds(),
		)
	})

	t.Run("events read back per Run keep cursor and publisher order", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		defer cancel()

		scenario := runParallelRunScenario(t, ctx)
		if scenario.measurement.sqliteBusy != 0 {
			t.Fatalf("expected zero SQLITE_BUSY occurrences, got %d", scenario.measurement.sqliteBusy)
		}
		if scenario.measurement.completedRuns != parallelRunsCount {
			t.Fatalf("expected all %d Runs to complete appends, got %d", parallelRunsCount, scenario.measurement.completedRuns)
		}
		for runIndex, fixture := range scenario.runs {
			events, err := fixture.store.RunEventsAfter(ctx, fixture.runID, 0, parallelRunEvents+1)
			if err != nil {
				t.Fatalf("read events for Run %d: %v", runIndex+1, err)
			}
			if len(events) != len(fixture.summaries) {
				t.Fatalf("Run %d: expected %d events, got %d", runIndex+1, len(fixture.summaries), len(events))
			}
			for eventIndex, event := range events {
				wantCursor := int64(eventIndex + 1)
				if event.Cursor != wantCursor {
					t.Fatalf("Run %d event %d: cursor = %d, want contiguous cursor %d", runIndex+1, eventIndex+1, event.Cursor, wantCursor)
				}
				if event.Event.Summary != fixture.summaries[eventIndex] {
					t.Fatalf("Run %d cursor %d: summary = %q, want publisher event %q", runIndex+1, event.Cursor, event.Event.Summary, fixture.summaries[eventIndex])
				}
			}
		}
	})

	t.Run("cancelling one Run releases the writer lock", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		defer cancel()

		fixtures := newParallelRunFixtures(t, ctx)
		for runIndex := 1; runIndex < len(fixtures); runIndex++ {
			// This rehearsal controls the flush boundary explicitly. A long
			// test-only linger keeps the production timer from racing that
			// boundary; the production batch size and write path stay intact.
			fixtures[runIndex].store.journal.maxLinger = time.Hour
			event := parallelRunEvent(fixtures[runIndex].runID, runIndex, 0)
			fixtures[runIndex].summaries = append(fixtures[runIndex].summaries, event.Summary)
			if err := fixtures[runIndex].store.JournalSink().Publish(ctx, event); err != nil {
				t.Fatalf("queue event for remaining Run %d: %v", runIndex+1, err)
			}
		}

		cancelledCtx, cancelRun := context.WithCancel(ctx)
		held := make(chan struct{})
		cancelledRunDone := make(chan error, 1)
		go func() {
			cancelledRunDone <- fixtures[0].store.withWriteTx(cancelledCtx, "cancelled Run rehearsal", func(_ *sql.Tx) error {
				close(held)
				<-cancelledCtx.Done()
				return cancelledCtx.Err()
			})
		}()
		select {
		case <-held:
		case err := <-cancelledRunDone:
			t.Fatalf("cancelled Run failed before acquiring the writer lock: %v", err)
		case <-ctx.Done():
			t.Fatalf("wait for cancelled Run to acquire the writer lock: %v", ctx.Err())
		}

		startFlush := make(chan struct{})
		flushAttempted := make(chan struct{}, parallelRunsCount-1)
		results := make(chan parallelRunWriteResult, parallelRunsCount-1)
		var ready sync.WaitGroup
		ready.Add(parallelRunsCount - 1)
		for runIndex := 1; runIndex < len(fixtures); runIndex++ {
			go func() {
				ready.Done()
				<-startFlush
				flushAttempted <- struct{}{}
				results <- parallelRunWriteResult{
					runIndex: runIndex,
					attempts: 1,
					err:      fixtures[runIndex].store.FlushJournal(ctx),
				}
			}()
		}
		ready.Wait()
		close(startFlush)
		for range parallelRunsCount - 1 {
			<-flushAttempted
		}
		cancelRun()

		if err := <-cancelledRunDone; !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled Run transaction error = %v, want context.Canceled", err)
		}

		busyCount := 0
		for range parallelRunsCount - 1 {
			result := <-results
			switch {
			case result.err == nil:
			case isSQLiteBusy(result.err):
				busyCount++
			default:
				t.Fatalf("remaining Run %d after cancellation: %v", result.runIndex+1, result.err)
			}
		}
		if busyCount != 0 {
			t.Fatalf("expected zero SQLITE_BUSY occurrences after cancellation, got %d", busyCount)
		}

		for runIndex := 1; runIndex < len(fixtures); runIndex++ {
			events, err := fixtures[runIndex].store.RunEventsAfter(ctx, fixtures[runIndex].runID, 0, 2)
			if err != nil {
				t.Fatalf("read remaining Run %d after cancellation: %v", runIndex+1, err)
			}
			if len(events) != 1 || events[0].Cursor != 1 || events[0].Event.Summary != fixtures[runIndex].summaries[0] {
				t.Fatalf("remaining Run %d did not preserve its append after cancellation: %#v", runIndex+1, events)
			}
		}
	})
}

func runParallelRunScenario(t *testing.T, ctx context.Context) parallelRunScenario {
	t.Helper()
	fixtures := newParallelRunFixtures(t, ctx)

	start := make(chan struct{})
	results := make(chan parallelRunWriteResult, parallelRunsCount)
	var ready sync.WaitGroup
	ready.Add(parallelRunsCount)
	for runIndex := range fixtures {
		go func() {
			ready.Done()
			<-start
			started := time.Now()
			result := parallelRunWriteResult{runIndex: runIndex}
			sink := fixtures[runIndex].store.JournalSink()
			for eventIndex := range parallelRunEvents {
				event := parallelRunEvent(fixtures[runIndex].runID, runIndex, eventIndex)
				fixtures[runIndex].summaries = append(fixtures[runIndex].summaries, event.Summary)
				result.attempts++
				if err := sink.Publish(ctx, event); err != nil {
					result.latency = time.Since(started)
					result.err = err
					results <- result
					return
				}
			}
			result.err = fixtures[runIndex].store.FlushJournal(ctx)
			result.latency = time.Since(started)
			results <- result
		}()
	}
	ready.Wait()
	started := time.Now()
	close(start)

	measurement := parallelRunMeasurement{}
	latencies := make([]time.Duration, 0, parallelRunsCount)
	var unexpected []string
	for range parallelRunsCount {
		result := <-results
		measurement.totalEvents += result.attempts
		switch {
		case result.err == nil:
			measurement.completedRuns++
			latencies = append(latencies, result.latency)
		case isSQLiteBusy(result.err):
			measurement.sqliteBusy++
		default:
			unexpected = append(unexpected, fmt.Sprintf("Run %d: %v", result.runIndex+1, result.err))
		}
	}
	measurement.elapsed = time.Since(started)
	if len(unexpected) > 0 {
		t.Fatalf("parallel Run appends failed: %v", unexpected)
	}
	measurement.writerLatencyP50 = journalHarnessPercentile(latencies, 50)
	measurement.writerLatencyP95 = journalHarnessPercentile(latencies, 95)
	if measurement.totalEvents > 0 {
		measurement.concurrentPerThousand = measurement.elapsed * 1_000 / time.Duration(measurement.totalEvents)
	}
	return parallelRunScenario{runs: fixtures, measurement: measurement}
}

func newParallelRunFixtures(t *testing.T, ctx context.Context) []parallelRunFixture {
	t.Helper()
	homeDir := t.TempDir()
	fixtures := make([]parallelRunFixture, parallelRunsCount)
	for runIndex := range fixtures {
		runStore := openTestStore(t, ctx, homeDir)
		t.Cleanup(func() {
			if err := runStore.Close(); err != nil {
				t.Errorf("close parallel Run %d store: %v", runIndex+1, err)
			}
		})
		setParallelRunsBusyTimeout(t, ctx, runStore)

		req := sampleCreateRunRequest()
		req.HeadBranch = fmt.Sprintf("parallel-run-%d", runIndex+1)
		req.LocalBranch = req.HeadBranch
		req.PRNumber = fmt.Sprintf("parallel-%d", runIndex+1)
		run, err := runStore.CreateRunSkippingActiveLock(ctx, req)
		if err != nil {
			t.Fatalf("create parallel Run %d: %v", runIndex+1, err)
		}
		fixtures[runIndex] = parallelRunFixture{store: runStore, runID: run.ID}
	}
	return fixtures
}

func setParallelRunsBusyTimeout(t *testing.T, ctx context.Context, runStore *Store) {
	t.Helper()
	//nolint:sql-injection-exec-sprintf-go // PRAGMA does not accept bind parameters; the value is an internal integer constant.
	if _, err := runStore.db.ExecContext(ctx, fmt.Sprintf("PRAGMA busy_timeout = %d", parallelRunsPreRaiseBusyTimeoutMillis)); err != nil {
		t.Fatalf("set pre-raise busy timeout: %v", err)
	}
	var got int
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA busy_timeout`).Scan(&got); err != nil {
		t.Fatalf("read pre-raise busy timeout: %v", err)
	}
	if got != parallelRunsPreRaiseBusyTimeoutMillis {
		t.Fatalf("busy timeout = %d ms, want pre-raise %d ms", got, parallelRunsPreRaiseBusyTimeoutMillis)
	}
}

func parallelRunEvent(runID string, runIndex int, eventIndex int) runevent.RunEvent {
	sequence := eventIndex + 1
	summary := fmt.Sprintf("parallel Run %d publisher event %d", runIndex+1, sequence)
	return runevent.RunEvent{
		RunID:   runID,
		Batch:   eventIndex/journalBatchSize + 1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentMessage,
		Summary: summary,
		Time:    time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC).Add(time.Duration(eventIndex) * time.Millisecond),
		Payload: []byte(fmt.Sprintf(`{"run":%d,"sequence":%d}`, runIndex+1, sequence)),
	}
}
