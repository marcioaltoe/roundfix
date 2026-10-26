package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/sharedcache"
)

type sharedCacheReport struct {
	path, kept         string
	size, limit, bytes int64
	entries            int
}

func sweepSharedCaches(ctx context.Context, loaded roundconfig.Loaded, dryRun bool) ([]sharedCacheReport, error) {
	var reports []sharedCacheReport
	var errs []error
	now := commandDependenciesForContext(ctx).gc.now
	for _, cache := range loaded.Config.Worktree.SharedCaches {
		if err := ctx.Err(); err != nil {
			return reports, errors.Join(append(errs, err)...)
		}
		path, err := roundconfig.ResolveSharedCachePath(cache.Path, os.TempDir(), loaded.HomeDir, os.Getuid())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		report := sharedCacheReport{path: path, limit: cache.MaxSizeMB * 1024 * 1024}
		plan, err := sharedcache.Measure(path, report.limit, cache.Idle, now())
		if err != nil {
			report.kept = err.Error()
			if !errors.Is(err, sharedcache.ErrNotGoBuildCache) {
				errs = append(errs, err)
			}
		} else {
			report.size, report.entries, report.bytes = plan.Bytes, len(plan.Remove), plan.RemoveSize
			if !dryRun {
				if err := ctx.Err(); err != nil {
					return reports, errors.Join(append(errs, err)...)
				}
				report.entries, report.bytes, err = sharedcache.Apply(plan, cache.Idle, now)
				if errors.Is(err, sharedcache.ErrNotGoBuildCache) {
					report.kept = err.Error()
				} else if err != nil {
					errs = append(errs, err)
				}
			}
		}
		reports = append(reports, report)
	}
	return reports, errors.Join(errs...)
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
