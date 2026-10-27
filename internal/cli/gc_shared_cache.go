package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/sharedcache"
	"roundfix/internal/store"
)

type sharedCacheOutcome string

const (
	sharedCacheTrimmed   sharedCacheOutcome = "trimmed"
	sharedCacheUnchanged sharedCacheOutcome = "unchanged"
	sharedCacheKept      sharedCacheOutcome = "kept"
	sharedCacheFailed    sharedCacheOutcome = "failed"
	sharedCacheStopped   sharedCacheOutcome = "stopped"
)

type sharedCacheReport struct {
	path, kept         string
	size, limit, bytes int64
	entries            int
	outcome            sharedCacheOutcome
	cause              string
	err                error
	unresolved         bool
}

func trimSharedCaches(ctx context.Context, loaded roundconfig.Loaded, dryRun bool) []sharedCacheReport {
	var reports []sharedCacheReport
	deps := commandDependenciesForContext(ctx).gc
	for _, cache := range loaded.Config.Worktree.SharedCaches {
		report := sharedCacheReport{path: cache.Path, limit: cache.MaxSizeMB * 1024 * 1024}
		path, err := roundconfig.ResolveSharedCachePath(cache.Path, os.TempDir(), loaded.HomeDir, os.Getuid())
		report.unresolved = err != nil
		if path != "" {
			report.path = path
		}
		if stopErr := ctx.Err(); stopErr != nil {
			err = stopErr
		}
		if err == nil {
			var plan sharedcache.Plan
			plan, err = sharedcache.MeasureContext(ctx, path, report.limit, cache.Idle, deps.now())
			if err == nil {
				report.size = plan.Bytes
				if dryRun {
					report.entries, report.bytes = len(plan.Remove), plan.RemoveSize
				} else {
					report.entries, report.bytes, err = sharedcache.ApplyContext(ctx, plan, cache.Idle, deps.now)
				}
			} else {
				// GC historically prints a measurement error as Kept, even when fatal.
				report.kept = err.Error()
			}
		}
		report.err = err
		switch {
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			report.outcome = sharedCacheStopped
			report.cause = err.Error()
			if errors.Is(err, context.DeadlineExceeded) {
				report.cause = fmt.Sprintf("budget of %s elapsed", deps.sharedCacheTrimBudget)
			}
		case errors.Is(err, sharedcache.ErrNotGoBuildCache):
			report.outcome, report.cause, report.kept = sharedCacheKept, err.Error(), err.Error()
		case err != nil:
			report.outcome, report.cause = sharedCacheFailed, err.Error()
		case report.entries > 0:
			report.outcome = sharedCacheTrimmed
		default:
			report.outcome = sharedCacheUnchanged
		}
		reports = append(reports, report)
	}
	return reports
}

func sweepSharedCaches(ctx context.Context, loaded roundconfig.Loaded, dryRun bool) ([]sharedCacheReport, error) {
	var reports []sharedCacheReport
	var errs []error
	for _, report := range trimSharedCaches(ctx, loaded, dryRun) {
		if report.err != nil && report.outcome != sharedCacheKept {
			errs = append(errs, report.err)
		}
		if !report.unresolved && report.outcome != sharedCacheStopped {
			reports = append(reports, report)
		}
	}
	return reports, errors.Join(errs...)
}

func printSharedCacheTrim(writer io.Writer, report sharedCacheReport, budget time.Duration, prefix, phase, continuation string) {
	if writer == nil {
		return
	}
	switch report.outcome {
	case sharedCacheTrimmed:
		fmt.Fprintf(writer, "roundfix: %sShared Build Cache %s: removed %d idle entries (%d bytes) %s\n", prefix, report.path, report.entries, report.bytes, phase)
	case sharedCacheStopped:
		fmt.Fprintf(writer, "roundfix: %sShared Build Cache %s: trim stopped after its %s budget %s; removed %d idle entries (%d bytes); the next trim continues\n", prefix, report.path, budget, phase, report.entries, report.bytes)
	case sharedCacheFailed:
		fmt.Fprintf(writer, "roundfix: warning: %sShared Build Cache %s: trim failed %s: %s; %s\n", prefix, report.path, phase, report.cause, continuation)
	}
}

func trimSharedCachesAtRunStart(ctx context.Context, loaded roundconfig.Loaded, runStore *store.Store, runID string, stderr io.Writer) {
	if len(loaded.Config.Worktree.SharedCaches) == 0 {
		return
	}
	budget := commandDependenciesForContext(ctx).gc.sharedCacheTrimBudget
	trimCtx, cancel := context.WithTimeout(ctx, budget)
	reports := trimSharedCaches(trimCtx, loaded, false)
	cancel()
	for _, report := range reports {
		printSharedCacheTrim(stderr, report, budget, "", "at Run start", "the Run continues")
		payload, err := json.Marshal(struct {
			Event   string             `json:"event"`
			Trigger string             `json:"trigger"`
			Path    string             `json:"path"`
			Outcome sharedCacheOutcome `json:"outcome"`
			Size    int64              `json:"size_bytes"`
			Limit   int64              `json:"limit_bytes"`
			Entries int                `json:"removed_entries"`
			Bytes   int64              `json:"removed_bytes"`
			Cause   string             `json:"cause,omitempty"`
		}{Event: "shared_cache_trim", Trigger: "run_start", Path: report.path, Outcome: report.outcome,
			Size: report.size, Limit: report.limit, Entries: report.entries, Bytes: report.bytes, Cause: report.cause})
		if err == nil {
			_, err = runStore.AppendRunEvent(ctx, runevent.RunEvent{
				RunID: runID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonStatus,
				Summary: runevent.BoundSummary(fmt.Sprintf("Shared Build Cache %s: %s", report.path, report.outcome)), Payload: payload,
			})
		}
		if err != nil {
			fmt.Fprintf(stderr, "roundfix: warning: Shared Build Cache %s: trim event not journaled: %v; the Run continues\n", report.path, err)
		}
	}
}

func printSharedCacheReport(stdout io.Writer, reports []sharedCacheReport, dryRun bool) {
	action := "removed"
	if dryRun {
		action = "reclaimable"
	}
	for _, report := range reports {
		fmt.Fprintf(stdout, "  Shared Build Cache: %s\n", report.path)
		if report.kept != "" {
			fmt.Fprintf(stdout, "    Kept: %s\n", report.kept)
			continue
		}
		fmt.Fprintf(stdout, "    Size: %d bytes of %d allowed\n", report.size, report.limit)
		fmt.Fprintf(stdout, "    Entries %s: %d (%d bytes)\n", action, report.entries, report.bytes)
	}
}
