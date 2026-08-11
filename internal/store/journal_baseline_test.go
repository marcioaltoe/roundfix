package store

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"roundfix/internal/runevent"
)

// Suite: Run Event Journal measurement harness
// Invariant: a fresh self-seeded journal reports comparable write and Run-start costs without using the operator database.
// Boundary IN: Store appends, SQLite writer contention, and the retention sweep run at Run start.
// Boundary OUT: later batching, writer-lock, and retention-query changes measured by the same harness.

type journalHarnessParameters struct {
	JournalSizes    []int
	WriterCounts    []int
	WritesPerWriter int
	PayloadBytes    int
	RunStartSamples int
}

type journalHarnessWriteResult struct {
	JournalEvents int
	Writers       int
	Attempts      int
	Successes     int
	SQLiteBusy    int
	LatencyP50    time.Duration
	LatencyP95    time.Duration
	LockWaitP50   time.Duration
	LockWaitP95   time.Duration
}

type journalHarnessRunStartResult struct {
	JournalEvents int
	DatabaseBytes int64
	Samples       int
	LatencyP50    time.Duration
	LatencyP95    time.Duration
}

type journalHarnessWriteSample struct {
	latency time.Duration
	err     error
}

func TestJournalMeasurementHarness(t *testing.T) {
	if testing.Short() {
		t.Skip("journal measurement harness seeds large journals; run without -short")
	}
	params := journalHarnessParameters{
		JournalSizes:    []int{0, 1_000, 10_000},
		WriterCounts:    []int{1, 2, 4},
		WritesPerWriter: 8,
		PayloadBytes:    1_800,
		RunStartSamples: 5,
	}
	if err := validateJournalHarnessParameters(params); err != nil {
		t.Fatalf("validate harness parameters: %v", err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	writeResults, runStartResults, facts := runJournalMeasurementHarness(t, ctx, params)
	t.Logf(
		"JOURNAL_HARNESS parameters journal_sizes=%v writer_counts=%v writes_per_writer=%d payload_bytes=%d run_start_samples=%d busy_timeout_ms=%d",
		params.JournalSizes,
		params.WriterCounts,
		params.WritesPerWriter,
		params.PayloadBytes,
		params.RunStartSamples,
		busyTimeoutMillis,
	)
	t.Logf(
		"JOURNAL_HARNESS machine go=%s os=%s arch=%s cpus=%d sqlite=%s journal_mode=%s synchronous=%d",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		facts.sqliteVersion,
		facts.journalMode,
		facts.synchronous,
	)
	for _, result := range writeResults {
		t.Logf(
			"JOURNAL_HARNESS write journal_events=%d writers=%d attempts=%d successes=%d sqlite_busy=%d latency_p50_us=%d latency_p95_us=%d lock_wait_p50_us=%d lock_wait_p95_us=%d",
			result.JournalEvents,
			result.Writers,
			result.Attempts,
			result.Successes,
			result.SQLiteBusy,
			result.LatencyP50.Microseconds(),
			result.LatencyP95.Microseconds(),
			result.LockWaitP50.Microseconds(),
			result.LockWaitP95.Microseconds(),
		)
	}
	for _, result := range runStartResults {
		t.Logf(
			"JOURNAL_HARNESS run_start_retention journal_events=%d database_bytes=%d samples=%d latency_p50_us=%d latency_p95_us=%d",
			result.JournalEvents,
			result.DatabaseBytes,
			result.Samples,
			result.LatencyP50.Microseconds(),
			result.LatencyP95.Microseconds(),
		)
	}
}

func TestJournalMeasurementHarnessRejectsInvalidParameters(t *testing.T) {
	tests := []struct {
		name   string
		params journalHarnessParameters
		want   string
	}{
		{
			name: "missing empty-journal baseline",
			params: journalHarnessParameters{
				JournalSizes: []int{1_000}, WriterCounts: []int{1}, WritesPerWriter: 1, PayloadBytes: 100, RunStartSamples: 1,
			},
			want: "journal sizes must start at zero",
		},
		{
			name: "missing uncontended writer baseline",
			params: journalHarnessParameters{
				JournalSizes: []int{0}, WriterCounts: []int{2}, WritesPerWriter: 1, PayloadBytes: 100, RunStartSamples: 1,
			},
			want: "writer counts must start at one",
		},
		{
			name: "payload cannot hold raw JSON envelope",
			params: journalHarnessParameters{
				JournalSizes: []int{0}, WriterCounts: []int{1}, WritesPerWriter: 1, PayloadBytes: 8, RunStartSamples: 1,
			},
			want: "payload bytes must be at least",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJournalHarnessParameters(tt.params)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("validateJournalHarnessParameters() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

type journalHarnessMachineFacts struct {
	sqliteVersion string
	journalMode   string
	synchronous   int
}

func runJournalMeasurementHarness(
	t *testing.T,
	ctx context.Context,
	params journalHarnessParameters,
) ([]journalHarnessWriteResult, []journalHarnessRunStartResult, journalHarnessMachineFacts) {
	t.Helper()

	writeResults := make([]journalHarnessWriteResult, 0, len(params.JournalSizes)*len(params.WriterCounts))
	runStartResults := make([]journalHarnessRunStartResult, 0, len(params.JournalSizes))
	var facts journalHarnessMachineFacts
	for sizeIndex, journalSize := range params.JournalSizes {
		runStartHomeDir := t.TempDir()
		seedStore := openTestStore(t, ctx, runStartHomeDir)
		seedJournalHarness(t, ctx, seedStore, journalSize, params.PayloadBytes)
		if sizeIndex == 0 {
			facts = readJournalHarnessMachineFacts(t, ctx, seedStore)
		}

		runStartResults = append(runStartResults, measureJournalHarnessRunStart(
			t,
			ctx,
			seedStore,
			journalSize,
			params.RunStartSamples,
		))
		closeStore(t, seedStore)

		var uncontendedMedian time.Duration
		for _, writerCount := range params.WriterCounts {
			homeDir := t.TempDir()
			writeSeedStore := openTestStore(t, ctx, homeDir)
			seedJournalHarness(t, ctx, writeSeedStore, journalSize, params.PayloadBytes)
			closeStore(t, writeSeedStore)

			result, latencies := measureJournalHarnessWrites(
				t,
				ctx,
				homeDir,
				journalSize,
				writerCount,
				params.WritesPerWriter,
				params.PayloadBytes,
			)
			if writerCount == 1 {
				uncontendedMedian = result.LatencyP50
			}
			lockWaits := make([]time.Duration, 0, len(latencies))
			for _, latency := range latencies {
				if writerCount == 1 || latency <= uncontendedMedian {
					lockWaits = append(lockWaits, 0)
					continue
				}
				lockWaits = append(lockWaits, latency-uncontendedMedian)
			}
			result.LockWaitP50 = journalHarnessPercentile(lockWaits, 50)
			result.LockWaitP95 = journalHarnessPercentile(lockWaits, 95)
			writeResults = append(writeResults, result)
		}
	}
	return writeResults, runStartResults, facts
}

func validateJournalHarnessParameters(params journalHarnessParameters) error {
	const rawJSONEnvelopeBytes = len(`{"chunk":""}`)
	if len(params.JournalSizes) == 0 || params.JournalSizes[0] != 0 {
		return errors.New("journal sizes must start at zero")
	}
	if len(params.WriterCounts) == 0 || params.WriterCounts[0] != 1 {
		return errors.New("writer counts must start at one")
	}
	if !slices.IsSorted(params.JournalSizes) || !slices.IsSorted(params.WriterCounts) {
		return errors.New("journal sizes and writer counts must be sorted")
	}
	if params.WritesPerWriter <= 0 {
		return errors.New("writes per writer must be positive")
	}
	if params.PayloadBytes < rawJSONEnvelopeBytes {
		return fmt.Errorf("payload bytes must be at least %d", rawJSONEnvelopeBytes)
	}
	if params.RunStartSamples <= 0 {
		return errors.New("Run-start samples must be positive")
	}
	return nil
}

func seedJournalHarness(t *testing.T, ctx context.Context, runStore *Store, journalSize int, payloadBytes int) {
	t.Helper()
	runStore.now = func() time.Time { return journalHarnessCutoff().Add(time.Hour) }
	req := sampleCreateRunRequest()
	req.HeadBranch = fmt.Sprintf("measurement-seed-%d", journalSize)
	req.LocalBranch = req.HeadBranch
	req.PRNumber = fmt.Sprintf("seed-%d", journalSize)
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create %d-event journal seed Run: %v", journalSize, err)
	}

	const seedBatchSize = 1_000
	payload := journalHarnessPayload(payloadBytes)
	for seeded := 0; seeded < journalSize; seeded += seedBatchSize {
		batchSize := min(seedBatchSize, journalSize-seeded)
		events := make([]runevent.RunEvent, batchSize)
		for index := range events {
			events[index] = journalHarnessEvent(run.ID, payload)
		}
		if _, err := runStore.AppendRunEvents(ctx, events); err != nil {
			t.Fatalf("seed journal events %d..%d of %d: %v", seeded, seeded+batchSize, journalSize, err)
		}
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatalf("complete %d-event journal seed Run: %v", journalSize, err)
	}
}

func measureJournalHarnessWrites(
	t *testing.T,
	ctx context.Context,
	homeDir string,
	journalSize int,
	writerCount int,
	writesPerWriter int,
	payloadBytes int,
) (journalHarnessWriteResult, []time.Duration) {
	t.Helper()

	writers := make([]*Store, writerCount)
	runIDs := make([]string, writerCount)
	for writerIndex := range writerCount {
		writers[writerIndex] = openTestStore(t, ctx, homeDir)
		req := sampleCreateRunRequest()
		req.HeadBranch = fmt.Sprintf("measurement-writers-%d-%d", writerCount, writerIndex)
		req.LocalBranch = req.HeadBranch
		req.PRNumber = fmt.Sprintf("writers-%d-%d", writerCount, writerIndex)
		run, err := writers[writerIndex].CreateRunSkippingActiveLock(ctx, req)
		if err != nil {
			t.Fatalf("create measurement Run for writer %d of %d: %v", writerIndex+1, writerCount, err)
		}
		runIDs[writerIndex] = run.ID
	}

	payload := journalHarnessPayload(payloadBytes)
	samples := make(chan journalHarnessWriteSample, writerCount*writesPerWriter)
	for range writesPerWriter {
		start := make(chan struct{})
		var ready sync.WaitGroup
		ready.Add(writerCount)
		var writersDone sync.WaitGroup
		for writerIndex := range writerCount {
			writersDone.Go(func() {
				ready.Done()
				<-start
				started := time.Now()
				_, err := writers[writerIndex].AppendRunEvent(ctx, journalHarnessEvent(runIDs[writerIndex], payload))
				samples <- journalHarnessWriteSample{latency: time.Since(started), err: err}
			})
		}
		ready.Wait()
		close(start)
		writersDone.Wait()
	}
	close(samples)

	for _, writer := range writers {
		closeStore(t, writer)
	}

	result := journalHarnessWriteResult{
		JournalEvents: journalSize,
		Writers:       writerCount,
		Attempts:      writerCount * writesPerWriter,
	}
	latencies := make([]time.Duration, 0, result.Attempts)
	for sample := range samples {
		switch {
		case sample.err == nil:
			result.Successes++
			latencies = append(latencies, sample.latency)
		case isSQLiteBusy(sample.err):
			result.SQLiteBusy++
		default:
			t.Fatalf("append measurement event with %d writers: %v", writerCount, sample.err)
		}
	}
	if result.Successes+result.SQLiteBusy != result.Attempts {
		t.Fatalf("measurement accounted for %d of %d write attempts", result.Successes+result.SQLiteBusy, result.Attempts)
	}
	result.LatencyP50 = journalHarnessPercentile(latencies, 50)
	result.LatencyP95 = journalHarnessPercentile(latencies, 95)
	return result, latencies
}

func measureJournalHarnessRunStart(
	t *testing.T,
	ctx context.Context,
	runStore *Store,
	journalSize int,
	sampleCount int,
) journalHarnessRunStartResult {
	t.Helper()

	latencies := make([]time.Duration, 0, sampleCount)
	cutoff := journalHarnessCutoff()
	for range sampleCount {
		started := time.Now()
		candidates, err := runStore.TerminalRunPruneCandidates(ctx, cutoff)
		if err != nil {
			t.Fatalf("measure Run-start retention eligibility at %d events: %v", journalSize, err)
		}
		if len(candidates) != 0 {
			t.Fatalf("expected self-seeded Active Run to be ineligible, got %d candidates", len(candidates))
		}
		pruned, err := runStore.PruneTerminalRuns(ctx, cutoff)
		if err != nil {
			t.Fatalf("measure locked Run-start retention scan at %d events: %v", journalSize, err)
		}
		if len(pruned.RunIDs) != 0 || pruned.Events != 0 {
			t.Fatalf("expected Run-start retention scan to preserve fixture, got %#v", pruned)
		}
		latencies = append(latencies, time.Since(started))
	}

	return journalHarnessRunStartResult{
		JournalEvents: journalSize,
		DatabaseBytes: journalHarnessDatabaseBytes(t, ctx, runStore),
		Samples:       sampleCount,
		LatencyP50:    journalHarnessPercentile(latencies, 50),
		LatencyP95:    journalHarnessPercentile(latencies, 95),
	}
}

func journalHarnessCutoff() time.Time {
	return time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC)
}

func journalHarnessEvent(runID string, payload []byte) runevent.RunEvent {
	return runevent.RunEvent{
		RunID:   runID,
		Batch:   1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentMessage,
		Summary: "measurement event",
		Time:    time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC),
		Payload: payload,
	}
}

func journalHarnessPayload(size int) []byte {
	const prefix = `{"chunk":"`
	const suffix = `"}`
	return []byte(prefix + strings.Repeat("x", size-len(prefix)-len(suffix)) + suffix)
}

func journalHarnessPercentile(samples []time.Duration, percentile int) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	sorted := slices.Clone(samples)
	slices.Sort(sorted)
	index := (len(sorted)*percentile + 99) / 100
	return sorted[index-1]
}

func isSQLiteBusy(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code()&0xff == sqlite3.SQLITE_BUSY
}

func readJournalHarnessMachineFacts(t *testing.T, ctx context.Context, runStore *Store) journalHarnessMachineFacts {
	t.Helper()
	var facts journalHarnessMachineFacts
	if err := runStore.db.QueryRowContext(ctx, `SELECT sqlite_version()`).Scan(&facts.sqliteVersion); err != nil {
		t.Fatalf("read SQLite version: %v", err)
	}
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA journal_mode`).Scan(&facts.journalMode); err != nil {
		t.Fatalf("read SQLite journal mode: %v", err)
	}
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA synchronous`).Scan(&facts.synchronous); err != nil {
		t.Fatalf("read SQLite synchronous mode: %v", err)
	}
	return facts
}

func journalHarnessDatabaseBytes(t *testing.T, ctx context.Context, runStore *Store) int64 {
	t.Helper()
	var pageCount int64
	var pageSize int64
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA page_count`).Scan(&pageCount); err != nil {
		t.Fatalf("read journal fixture page count: %v", err)
	}
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA page_size`).Scan(&pageSize); err != nil {
		t.Fatalf("read journal fixture page size: %v", err)
	}
	return pageCount * pageSize
}
