package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

type itemBinaryReport struct {
	paths   []string
	bytes   int64
	kept    int
	present bool
}

// Containment is proven before measuring and again before removing. Neither a
// linked Item Binary nor a linked item directory is followed.
func safeItemBinaryPath(root, slug string) (string, error) {
	if slug == "." || slug == ".." || filepath.Base(slug) != slug || slug == "" {
		return "", fmt.Errorf("unsafe Item Binary slug %q", slug)
	}
	item := filepath.Join(root, slug)
	for _, path := range []string{item, filepath.Join(item, "item-binary")} {
		info, err := os.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("inspect Item Binary path %q: %w", path, err)
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return "", nil
		}
	}
	physicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve delivery artifact root %q: %w", root, err)
	}
	path := filepath.Join(item, "item-binary")
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve Item Binary path %q: %w", path, err)
	}
	within, err := gcPathWithin(physicalRoot, physical)
	if err != nil {
		return "", fmt.Errorf("check Item Binary containment: %w", err)
	}
	if !within || physical == physicalRoot {
		return "", nil
	}
	return path, nil
}

func removeMergedItemBinaries(artifactRoot, slug string) error {
	path, err := safeItemBinaryPath(filepath.Join(artifactRoot, "delivery"), slug)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove merged Item Binary directory %q: %w", path, err)
	}
	return nil
}

func itemBinariesReclaimable(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, slug string) (eligible bool, err error) {
	if runStore == nil {
		// A queue may have been created since GC first checked the database.
		var exists bool
		exists, err = gcDatabaseExists(loaded.HomeDir)
		if err != nil || !exists {
			return !exists, err
		}
		runStore, err = openGCStore(ctx, true, loaded.HomeDir)
		if err != nil {
			return false, err
		}
		defer func() { err = errors.Join(err, runStore.Close()) }()
	}
	queue, _, err := runStore.DeliveryQueue(ctx, loaded.GitRoot)
	if err != nil {
		return false, err
	}
	for _, item := range queue.Items {
		if item.SpecSlug == slug && item.Stage != store.DeliveryStageMerged {
			return false, nil
		}
	}
	return true, nil
}

func sweepItemBinaries(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, dryRun bool) (itemBinaryReport, error) {
	report, scanErr := scanItemBinaries(ctx, runStore, loaded)
	if dryRun {
		return report, scanErr
	}
	removed, removalErr := removeItemBinaries(ctx, runStore, loaded, report)
	return removed, errors.Join(scanErr, removalErr)
}

func scanItemBinaries(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded) (itemBinaryReport, error) {
	var report itemBinaryReport
	if loaded.GitRoot == "" {
		return report, nil
	}
	artifactRoot, err := roundconfig.ResolveArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		return report, err
	}
	root := filepath.Join(artifactRoot, "delivery")
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return report, nil
	}
	if err != nil {
		return report, fmt.Errorf("read delivery artifact root %q: %w", root, err)
	}
	repository, err := roundconfig.RepositoryRoot(loaded.GitRoot)
	if err != nil {
		return report, fmt.Errorf("resolve Item Binary repository: %w", err)
	}
	loaded.GitRoot = repository
	var errs []error
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return report, errors.Join(append(errs, err)...)
		}
		if !entry.IsDir() && entry.Type()&os.ModeSymlink == 0 {
			continue
		}
		candidate := filepath.Join(root, entry.Name(), "item-binary")
		info, err := os.Lstat(candidate)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("inspect Item Binary %q: %w", candidate, err))
			continue
		}
		if !info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		report.present = true
		path, err := safeItemBinaryPath(root, entry.Name())
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if path == "" {
			report.kept++
			continue
		}
		eligible, err := itemBinariesReclaimable(ctx, runStore, loaded, entry.Name())
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if !eligible {
			report.kept++
			continue
		}
		size, exists, err := gcDirectoryBytes(path)
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if !exists {
			continue
		}
		report.paths = append(report.paths, path)
		report.bytes += size
	}
	return report, errors.Join(errs...)
}

// Remeasure each candidate and reread the queue immediately before removal.
func removeItemBinaries(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, scan itemBinaryReport) (itemBinaryReport, error) {
	report := itemBinaryReport{kept: scan.kept, present: scan.present}
	if len(scan.paths) == 0 {
		return report, nil
	}
	repository, err := roundconfig.RepositoryRoot(loaded.GitRoot)
	if err != nil {
		return report, fmt.Errorf("resolve Item Binary repository: %w", err)
	}
	loaded.GitRoot = repository
	var errs []error
	for _, candidate := range scan.paths {
		if err := ctx.Err(); err != nil {
			return report, errors.Join(append(errs, err)...)
		}
		root, slug := filepath.Dir(filepath.Dir(candidate)), filepath.Base(filepath.Dir(candidate))
		path, err := safeItemBinaryPath(root, slug)
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if path == "" {
			report.kept++
			continue
		}
		size, exists, err := gcDirectoryBytes(path)
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if !exists {
			continue
		}
		// A long byte walk cannot authorize deletion of a newly linked directory.
		path, err = safeItemBinaryPath(root, slug)
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if path == "" {
			report.kept++
			continue
		}
		eligible, err := itemBinariesReclaimable(ctx, runStore, loaded, slug)
		if err != nil {
			errs = append(errs, err)
			report.kept++
			continue
		}
		if !eligible {
			report.kept++
			continue
		}
		if err := ctx.Err(); err != nil {
			return report, errors.Join(append(errs, err)...)
		}
		if err := os.RemoveAll(path); err != nil {
			errs = append(errs, fmt.Errorf("remove Item Binary directory %q: %w", path, err))
			report.kept++
			continue
		}
		report.paths = append(report.paths, path)
		report.bytes += size
	}
	return report, errors.Join(errs...)
}

func printItemBinaryReport(stdout io.Writer, report itemBinaryReport, dryRun bool) {
	if !report.present {
		return
	}
	action, bytes, title := "removed", "reclaimed", "Removed"
	if dryRun {
		action, bytes, title = "reclaimable", "reclaimable", "Reclaimable"
	}
	fmt.Fprintf(stdout, "  Item Binaries %s: %d\n", action, len(report.paths))
	fmt.Fprintf(stdout, "  Item Binary bytes %s: %d\n", bytes, report.bytes)
	fmt.Fprintf(stdout, "  Item Binaries kept: %d\n", report.kept)
	fmt.Fprintf(stdout, "  %s Item Binaries:\n", title)
	for _, path := range report.paths {
		fmt.Fprintf(stdout, "    %s\n", path)
	}
}
