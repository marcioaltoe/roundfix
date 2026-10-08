package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

type runRetentionOptions struct {
	dryRun bool
	budget time.Duration
}

type runRetentionReport struct {
	retentionDays                                                int
	cutoff                                                       time.Time
	runIDs                                                       []string
	kept, queueReferenced, worktreePresent, artifactRootUnproven int
	rows                                                         store.RunRetentionRows
	estimatedBytes, artifactBytes, bytesBefore, bytesAfter       int64
	compaction                                                   string
	paused                                                       bool
}

func sweepRunDatabase(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, opts runRetentionOptions) (runRetentionReport, error) {
	deps := commandDependenciesForContext(ctx).gc
	started := deps.now().UTC()
	report := runRetentionReport{retentionDays: loaded.Config.Store.RunRetentionDays, compaction: "not needed"}
	report.cutoff = started.Add(-time.Duration(report.retentionDays) * 24 * time.Hour)
	expired := func() bool { return opts.budget > 0 && deps.now().Sub(started) >= opts.budget }
	var err error
	report.bytesBefore, _, _, err = runStore.RunRetentionStorage(ctx)
	if err != nil {
		return report, err
	}
	candidates, err := runStore.RunRetentionCandidates(ctx, report.cutoff)
	if err != nil {
		return report, err
	}
	for _, candidate := range candidates {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		worktree := false
		if candidate.WorkDir != "" {
			worktree, err = runRetentionPathExists(candidate.WorkDir)
			if err != nil {
				return report, err
			}
		}
		dir, unproven, err := runRetentionArtifact(candidate, loaded.HomeDir)
		if err != nil {
			return report, err
		}
		if candidate.QueueReferenced || worktree || unproven {
			report.kept++
			if candidate.QueueReferenced {
				report.queueReferenced++
			}
			if worktree {
				report.worktreePresent++
			}
			if unproven {
				report.artifactRootUnproven++
			}
			continue
		}
		if opts.dryRun {
			report.runIDs = append(report.runIDs, candidate.RunID)
			addRunRetentionRows(&report.rows, candidate.Rows)
			report.estimatedBytes += candidate.EstimatedBytes
			continue
		}
		if expired() {
			report.paused = true
			break
		}
		if dir.path != "" {
			if err := gcRemoveArtifactDirs([]gcArtifactDir{dir}); err != nil {
				return report, err
			}
			report.artifactBytes += dir.bytes
		}
		rows, err := runStore.RemoveRetainedRun(ctx, candidate.RunID, report.cutoff)
		if err != nil {
			var kept store.RunRetentionKeptError
			if errors.As(err, &kept) {
				report.kept++
				if kept.Reason == "Run is queue-referenced" {
					report.queueReferenced++
				}
				continue
			}
			return report, err
		}
		if rows.Runs > 0 {
			report.runIDs = append(report.runIDs, candidate.RunID)
			addRunRetentionRows(&report.rows, rows)
		}
	}
	if opts.dryRun {
		return report, nil
	}
	_, freePages, mode, err := runStore.RunRetentionStorage(ctx)
	if err != nil {
		return report, err
	}
	var pagesReleased int64
	if !report.paused && mode == 2 {
		for freePages > 0 {
			if expired() {
				report.paused = true
				break
			}
			slice, err := runStore.CompactIncrementally(ctx, store.IncrementalCompactionSlicePages)
			if err != nil {
				return report, err
			}
			pagesReleased += slice.PagesReleased
			freePages = slice.FreePages
		}
		if pagesReleased > 0 {
			report.compaction = fmt.Sprintf("incremental (%d pages)", pagesReleased)
		}
		if err := runStore.CheckpointRunRetention(ctx); err != nil {
			return report, err
		}
	} else if !report.paused && opts.budget == 0 && (freePages > 0 || mode != 2) {
		preview, err := deps.previewCompaction(ctx, runStore)
		if err == nil {
			_, err = deps.compact(ctx, runStore, preview)
		}
		if err != nil {
			var active store.ActiveRunCompactionError
			var writer store.WriterPresentCompactionError
			var capacity store.CompactionCapacityError
			var stale store.CompactionPreviewStaleError
			if !errors.As(err, &active) && !errors.As(err, &writer) && !errors.As(err, &capacity) && !errors.As(err, &stale) {
				return report, err
			}
			report.compaction = fmt.Sprintf("skipped (%v)", err)
		} else if mode == 0 {
			report.compaction = "full (converted to incremental)"
		} else {
			report.compaction = "full"
		}
	}
	if !report.paused {
		if expired() {
			report.paused = true
		} else if err := runStore.RecordRunRetentionSweep(ctx, store.RunRetentionSweep{CompletedAt: deps.now().UTC(), RetentionDays: report.retentionDays}); err != nil {
			return report, err
		}
	}
	report.bytesAfter, _, _, err = runStore.RunRetentionStorage(ctx)
	return report, err
}

func runRetentionPathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect Run Retention path %q: %w", path, err)
	}
	return true, nil
}

func runRetentionArtifact(candidate store.RunRetentionCandidate, homeDir string) (gcArtifactDir, bool, error) {
	root := candidate.ArtifactDir
	if root == "" {
		var err error
		root, err = roundconfig.ResolveArtifactDirectory("", candidate.Repository, homeDir)
		checkoutRoot, checkoutErr := roundconfig.DefaultArtifactDirectoryForPath(candidate.GitRoot, homeDir)
		if err != nil {
			if checkoutErr != nil {
				return gcArtifactDir{}, false, errors.Join(err, checkoutErr)
			}
			root = checkoutRoot
		} else if checkoutErr == nil && checkoutRoot != root {
			keyPath, pathErr := gcRunArtifactPath(root, candidate.RunID)
			if pathErr != nil {
				return gcArtifactDir{}, true, nil
			}
			keyExists, existsErr := runRetentionPathExists(keyPath)
			if existsErr != nil {
				return gcArtifactDir{}, false, existsErr
			}
			if !keyExists {
				root = checkoutRoot
			}
		}
	}
	path, err := gcRunArtifactPath(root, candidate.RunID)
	if err != nil {
		return gcArtifactDir{}, true, nil
	}
	exists, err := runRetentionPathExists(path)
	if err != nil || !exists {
		return gcArtifactDir{}, false, err
	}
	// Sanitation's physical-path and Home proof also applies to recorded roots;
	// an explicit root inside Home is valid for whole-Run retention.
	proof := classifyGCSanitationRoot(store.ArtifactRoot{Path: root}, homeDir)
	if proof.classification != gcSanitationOrphaned {
		return gcArtifactDir{}, true, nil
	}
	physical, err := filepath.EvalSymlinks(filepath.Join(root, "runs"))
	if err != nil {
		return gcArtifactDir{}, false, fmt.Errorf("resolve Run artifact root: %w", err)
	}
	home, err := filepath.EvalSymlinks(filepath.Dir(store.DatabasePath(homeDir)))
	if err != nil {
		return gcArtifactDir{}, false, fmt.Errorf("resolve Roundfix Home: %w", err)
	}
	inside, err := gcPathWithin(home, physical)
	if err != nil || !inside {
		return gcArtifactDir{}, true, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return gcArtifactDir{}, false, err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return gcArtifactDir{}, true, nil
	}
	dir, err := gcRunArtifactDir(root, candidate.RunID)
	return dir, false, err
}

func addRunRetentionRows(total *store.RunRetentionRows, rows store.RunRetentionRows) {
	total.Runs += rows.Runs
	total.RunEvents += rows.RunEvents
	total.AgentSelections += rows.AgentSelections
	total.TokenUsage += rows.TokenUsage
	total.ActiveRunLocks += rows.ActiveRunLocks
}

func printRunRetentionReport(stdout io.Writer, report runRetentionReport, dryRun bool) {
	fmt.Fprintf(stdout, "  Run Retention: %d days\n  Run Retention cutoff: %s\n", report.retentionDays, report.cutoff.Format(time.RFC3339Nano))
	action, list := "removed", "Removed Runs"
	if dryRun {
		action, list = "removable", "Removable Runs"
	}
	fmt.Fprintf(stdout, "  Runs %s: %d\n", action, len(report.runIDs))
	fmt.Fprintf(stdout, "  Runs kept past the cutoff: %d (queue-referenced %d, worktree present %d, artifact root unproven %d)\n", report.kept, report.queueReferenced, report.worktreePresent, report.artifactRootUnproven)
	fmt.Fprintf(stdout, "  Rows %s: runs=%d run_events=%d run_agent_selections=%d run_token_usage=%d active_run_locks=%d\n", action, report.rows.Runs, report.rows.RunEvents, report.rows.AgentSelections, report.rows.TokenUsage, report.rows.ActiveRunLocks)
	if dryRun {
		fmt.Fprintf(stdout, "  Database bytes reclaimable (estimated): %d\n", report.estimatedBytes)
	} else {
		fmt.Fprintf(stdout, "  Run artifact bytes reclaimed: %d\n  Compaction: %s\n  Database bytes before: %d\n  Database bytes after: %d\n", report.artifactBytes, report.compaction, report.bytesBefore, report.bytesAfter)
	}
	printGCIDList(stdout, list, report.runIDs)
}
