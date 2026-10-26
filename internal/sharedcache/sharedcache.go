// Package sharedcache plans and applies conservative trims of Go build caches.
package sharedcache

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// MinimumIdle is twice the go command's hourly refresh of a used entry.
const MinimumIdle = 2 * time.Hour

// GoBuildReadme is the prefix identifying the go command's cache README.
const GoBuildReadme = "This directory holds cached build artifacts from the Go build system."

// ErrNotGoBuildCache identifies a directory refused for trimming.
var ErrNotGoBuildCache = errors.New("not a Go build cache")

// Entry is a regular cache file observed during measurement.
type Entry struct {
	Path    string
	Bytes   int64
	ModTime time.Time
}

// Plan describes the measured cache and its oldest idle entries to remove.
type Plan struct {
	Path         string
	Bytes, Limit int64
	Remove       []Entry
	RemoveSize   int64
}

var entryName = regexp.MustCompile(`^[0-9a-f]{2}/[0-9a-f]{64}-[ad]$`)

func openCache(dir string) (*os.Root, error) {
	info, err := os.Lstat(filepath.Clean(dir))
	if err != nil {
		return nil, fmt.Errorf("%w: inspect %q: %w", ErrNotGoBuildCache, dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%w: %q is not a directory or is a symbolic link", ErrNotGoBuildCache, dir)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: open %q: %w", ErrNotGoBuildCache, dir, err)
	}
	info, err = root.Lstat("README")
	if err == nil && !info.Mode().IsRegular() {
		err = errors.New("README is not a regular file")
	}
	var readme []byte
	if err == nil {
		readme, err = root.ReadFile("README")
	}
	if err == nil && !strings.HasPrefix(string(readme), GoBuildReadme) {
		err = errors.New("unrecognized README")
	}
	if err != nil {
		return nil, fmt.Errorf("%w: validate %q: %w", ErrNotGoBuildCache, dir, errors.Join(err, root.Close()))
	}
	return root, nil
}

// Measure reads all regular file sizes and plans only idle cache entries.
// It leaves the directory unchanged, even when it exceeds limit.
func Measure(dir string, limit int64, idle time.Duration, now time.Time) (plan Plan, err error) {
	if idle < MinimumIdle {
		return plan, fmt.Errorf("%w: idle time must be at least %s", ErrNotGoBuildCache, MinimumIdle)
	}
	root, err := openCache(dir)
	if err != nil {
		return plan, err
	}
	defer func() { err = errors.Join(err, root.Close()) }()
	plan.Path, plan.Limit = dir, limit
	var candidates []Entry
	err = fs.WalkDir(root.FS(), ".", func(path string, d fs.DirEntry, walkErr error) error {
		if errors.Is(walkErr, fs.ErrNotExist) {
			return nil
		}
		if walkErr != nil {
			return walkErr
		}
		info, statErr := root.Lstat(filepath.FromSlash(path))
		if errors.Is(statErr, fs.ErrNotExist) {
			return nil
		}
		if statErr != nil {
			return statErr
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		plan.Bytes += info.Size()
		if entryName.MatchString(path) && info.ModTime().Before(now.Add(-idle)) {
			candidates = append(candidates, Entry{Path: filepath.Join(dir, filepath.FromSlash(path)), Bytes: info.Size(), ModTime: info.ModTime()})
		}
		return nil
	})
	if err != nil {
		return plan, fmt.Errorf("measure Go build cache %q: %w", dir, err)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ModTime.Equal(candidates[j].ModTime) {
			return candidates[i].Path < candidates[j].Path
		}
		return candidates[i].ModTime.Before(candidates[j].ModTime)
	})
	for _, entry := range candidates {
		if plan.Bytes-plan.RemoveSize <= limit {
			break
		}
		plan.Remove = append(plan.Remove, entry)
		plan.RemoveSize += entry.Bytes
	}
	return plan, nil
}

// Apply rereads each planned entry immediately before removal, keeping entries
// used since measurement. Missing files are harmless; other errors are joined.
func Apply(plan Plan, idle time.Duration, now func() time.Time) (removed int, bytes int64, err error) {
	if idle < MinimumIdle {
		return 0, 0, fmt.Errorf("%w: idle time must be at least %s", ErrNotGoBuildCache, MinimumIdle)
	}
	root, err := openCache(plan.Path)
	if err != nil {
		return 0, 0, err
	}
	defer func() { err = errors.Join(err, root.Close()) }()
	var errs []error
	for _, entry := range plan.Remove {
		relative, relErr := filepath.Rel(plan.Path, entry.Path)
		if relErr != nil || !entryName.MatchString(filepath.ToSlash(relative)) {
			errs = append(errs, fmt.Errorf("refuse cache entry %q outside the entry naming contract", entry.Path))
			continue
		}
		// Do not follow a shard replaced by a symlink after measurement.
		shard, statErr := root.Lstat(filepath.Dir(relative))
		if statErr == nil && !shard.IsDir() {
			continue
		}
		if statErr == nil {
			var info fs.FileInfo
			info, statErr = root.Lstat(relative)
			if statErr == nil {
				if !info.Mode().IsRegular() || !info.ModTime().Before(now().Add(-idle)) {
					continue
				}
				removeErr := root.Remove(relative)
				if errors.Is(removeErr, fs.ErrNotExist) {
					continue
				}
				if removeErr != nil {
					errs = append(errs, fmt.Errorf("remove cache entry %q: %w", entry.Path, removeErr))
					continue
				}
				removed++
				bytes += info.Size()
				continue
			}
		}
		if !errors.Is(statErr, fs.ErrNotExist) {
			errs = append(errs, fmt.Errorf("inspect cache entry %q: %w", entry.Path, statErr))
		}
	}
	return removed, bytes, errors.Join(errs...)
}
