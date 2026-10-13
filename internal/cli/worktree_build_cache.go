package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const doctorWorktreeBuildCacheBudget = 2 * time.Second

// Discovery merges all owners of a worktree: any unfinished owner protects it.
type worktreeBuildCache struct {
	Worktree, Path string
	Bytes          int64
}
type worktreeBuildCacheScan struct {
	Reclaimable []worktreeBuildCache
	Kept        int
	Bytes       int64
	Bounded     bool
}
type worktreeBuildCacheOwner struct {
	path  string
	ended bool
}

var worktreeBuildCacheTaskID = regexp.MustCompile(`^task_[0-9]+$`)

func discoverWorktreeBuildCacheOwners(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded) ([]worktreeBuildCacheOwner, error) {
	if strings.TrimSpace(loaded.GitRoot) == "" {
		return nil, nil
	}
	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{GitRoot: loaded.GitRoot, States: store.StatesAll})
	if err != nil {
		return nil, err
	}
	owners := map[string]worktreeBuildCacheOwner{}
	add := func(path string, ended bool) {
		if path == "" {
			return
		}
		path = filepath.Clean(path)
		key := path
		if physical, err := filepath.EvalSymlinks(path); err == nil {
			key = physical
		}
		previous, found := owners[key]
		if found {
			ended = ended && previous.ended
			path = previous.path
		}
		owners[key] = worktreeBuildCacheOwner{path: path, ended: ended}
	}
	location, err := roundconfig.ResolveWorktreeLocation(loaded.Config.Worktree.Location, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		return nil, err
	}
	for _, run := range runs {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if run.WorkDir == "" || !worktreeBuildCacheWithin(location, run.WorkDir) {
			continue
		}
		ended := store.IsTerminalState(run.State)
		add(run.WorkDir, ended)
		// Task Worktrees survive as siblings even when the Run Worktree is gone.
		entries, err := os.ReadDir(filepath.Dir(run.WorkDir))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read Task Worktrees of Run %q: %w", run.ID, err)
		}
		prefix := run.ID + "."
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && worktreeBuildCacheTaskID.MatchString(strings.TrimPrefix(entry.Name(), prefix)) {
				add(filepath.Join(filepath.Dir(run.WorkDir), entry.Name()), ended)
			}
		}
	}
	queue, _, err := runStore.DeliveryQueue(ctx, loaded.GitRoot)
	if err != nil {
		return nil, err
	}
	for _, item := range queue.Items {
		add(item.Worktree, item.Stage == store.DeliveryStageMerged)
	}
	result := make([]worktreeBuildCacheOwner, 0, len(owners))
	for _, owner := range owners {
		result = append(result, owner)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].path < result[j].path })
	return result, nil
}

func worktreeBuildCacheWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != "." && rel != ".." && !filepath.IsAbs(rel) && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// Prove lexical and physical containment before either measurement or removal.
func safeWorktreeBuildCachePath(loaded roundconfig.Loaded, worktree, entry string) (string, error) {
	location, err := roundconfig.ResolveWorktreeLocation(loaded.Config.Worktree.Location, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		return "", err
	}
	if !worktreeBuildCacheWithin(location, worktree) {
		return "", nil
	}
	physicalLocation, err := filepath.EvalSymlinks(location)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("resolve worktree location: %w", err)
	}
	physicalWorktree, err := filepath.EvalSymlinks(worktree)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("resolve worktree %q: %w", worktree, err)
	}
	if !worktreeBuildCacheWithin(physicalLocation, physicalWorktree) {
		return "", nil
	}
	path := filepath.Join(worktree, filepath.FromSlash(entry))
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("inspect Worktree Build Cache %q: %w", path, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", nil
	}
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve Worktree Build Cache %q: %w", path, err)
	}
	gitDir := filepath.Join(physicalWorktree, ".git")
	if !worktreeBuildCacheWithin(physicalWorktree, physical) || physical == gitDir || worktreeBuildCacheWithin(gitDir, physical) {
		return "", nil
	}
	return path, nil
}

func scanWorktreeBuildCaches(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, budget time.Duration) (worktreeBuildCacheScan, error) {
	var scan worktreeBuildCacheScan
	if len(loaded.Config.Worktree.Caches) == 0 {
		return scan, nil
	}
	owners, err := discoverWorktreeBuildCacheOwners(ctx, runStore, loaded)
	if err != nil {
		return scan, err
	}
	var errs []error
	now := commandDependenciesForContext(ctx).gc.now
	deadline := now().Add(budget)
	expired := func() bool { return budget > 0 && !now().Before(deadline) }
	for _, owner := range owners {
		for _, entry := range loaded.Config.Worktree.Caches {
			if err := ctx.Err(); err != nil {
				return scan, err
			}
			path, err := safeWorktreeBuildCachePath(loaded, owner.path, entry)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if path == "" {
				continue
			}
			if !owner.ended {
				scan.Kept++
				continue
			}
			cache := worktreeBuildCache{Worktree: owner.path, Path: path}
			if !scan.Bounded {
				cache.Bytes, scan.Bounded, err = measureWorktreeBuildCache(ctx, path, expired)
				if err != nil {
					errs = append(errs, err)
				}
			}
			scan.Bytes += cache.Bytes
			scan.Reclaimable = append(scan.Reclaimable, cache)
		}
	}
	return scan, errors.Join(errs...)
}

// Read directories in batches and check the shared deadline between entries.
// No background walk is abandoned when Doctor reaches its budget.
func measureWorktreeBuildCache(ctx context.Context, root string, expired func() bool) (int64, bool, error) {
	var bytes int64
	pending := []string{root}
	for len(pending) > 0 {
		if err := ctx.Err(); err != nil {
			return bytes, false, err
		}
		if expired() {
			return bytes, true, nil
		}
		dir := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		file, err := os.Open(dir)
		if err != nil {
			return bytes, false, fmt.Errorf("open Worktree Build Cache directory %q: %w", dir, err)
		}
		walkErr := func() error {
			for {
				if err := ctx.Err(); err != nil {
					return err
				}
				if expired() {
					return nil
				}
				entries, readErr := file.ReadDir(128)
				if readErr != nil && !errors.Is(readErr, io.EOF) {
					return readErr
				}
				for _, entry := range entries {
					if err := ctx.Err(); err != nil {
						return err
					}
					if expired() {
						return nil
					}
					if entry.Type()&os.ModeSymlink != 0 {
						continue
					}
					if entry.IsDir() {
						pending = append(pending, filepath.Join(dir, entry.Name()))
						continue
					}
					info, err := entry.Info()
					if err != nil {
						return err
					}
					if info.Mode().IsRegular() {
						bytes += info.Size()
					}
				}
				if errors.Is(readErr, io.EOF) {
					return nil
				}
			}
		}()
		closeErr := file.Close()
		if err := errors.Join(walkErr, closeErr); err != nil {
			return bytes, false, fmt.Errorf("measure Worktree Build Cache %q: %w", dir, err)
		}
		if expired() {
			return bytes, true, nil
		}
	}
	return bytes, false, nil
}

func removeWorktreeBuildCaches(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, scan worktreeBuildCacheScan) (worktreeBuildCacheScan, error) {
	report := worktreeBuildCacheScan{Kept: scan.Kept}
	var errs []error
	for _, cache := range scan.Reclaimable {
		entry, err := filepath.Rel(cache.Worktree, cache.Path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		path, err := safeWorktreeBuildCachePath(loaded, cache.Worktree, filepath.ToSlash(entry))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if path == "" {
			continue
		}
		bytes, _, err := measureWorktreeBuildCache(ctx, path, func() bool { return false })
		if err != nil {
			errs = append(errs, fmt.Errorf("measure Worktree Build Cache %q before removal: %w", path, err))
			continue
		}
		// Rediscovery rereads Runs and queue items, including competing owners.
		owners, err := discoverWorktreeBuildCacheOwners(ctx, runStore, loaded)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		ended := false
		for _, owner := range owners {
			physicalOwner, ownerErr := filepath.EvalSymlinks(owner.path)
			physicalCache, cacheErr := filepath.EvalSymlinks(cache.Worktree)
			if owner.path == cache.Worktree || (ownerErr == nil && cacheErr == nil && physicalOwner == physicalCache) {
				ended = owner.ended
				break
			}
		}
		if !ended {
			report.Kept++
			continue
		}
		// Reprove containment too, after the potentially long byte walk.
		path, err = safeWorktreeBuildCachePath(loaded, cache.Worktree, filepath.ToSlash(entry))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if path == "" {
			continue
		}
		err = os.RemoveAll(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("remove Worktree Build Cache %q: %w", path, err))
			continue
		}
		cache.Bytes = bytes
		report.Reclaimable = append(report.Reclaimable, cache)
		report.Bytes += bytes
	}
	return report, errors.Join(errs...)
}

func printWorktreeBuildCacheReport(stdout io.Writer, declared []string, report worktreeBuildCacheScan, dryRun bool) {
	if len(declared) == 0 {
		return
	}
	action, bytes, title := "removed", "reclaimed", "Removed"
	if dryRun {
		action, bytes, title = "reclaimable", "reclaimable", "Reclaimable"
	}
	fmt.Fprintf(stdout, "  Worktree Build Caches declared: %s\n", strings.Join(declared, ", "))
	fmt.Fprintf(stdout, "  Worktree Build Caches %s: %d\n", action, len(report.Reclaimable))
	fmt.Fprintf(stdout, "  Worktree Build Cache bytes %s: %d\n", bytes, report.Bytes)
	fmt.Fprintf(stdout, "  Worktree Build Caches kept: %d\n", report.Kept)
	paths := make([]string, 0, len(report.Reclaimable))
	for _, cache := range report.Reclaimable {
		paths = append(paths, cache.Path)
	}
	sort.Strings(paths)
	printGCIDList(stdout, title+" Worktree Build Caches", paths)
}
