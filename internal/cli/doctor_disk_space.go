package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	roundconfig "roundfix/internal/config"
)

const doctorDiskSpaceThreshold int64 = 10 * 1024 * 1024 * 1024

type freeSpaceLocation struct {
	label, path string
	available   int64
	err         error
}

func inspectFreeSpace(ctx context.Context, loaded roundconfig.Loaded) []freeSpaceLocation {
	availableBytes := commandDependenciesForContext(ctx).availableBytes
	var locations []freeSpaceLocation
	inspect := func(label, path string, err error) {
		location := freeSpaceLocation{label: label, path: path, err: err}
		if location.err == nil {
			var ancestor string
			ancestor, location.err = doctorDiskSpaceAncestor(path)
			if location.err == nil {
				location.available, location.err = availableBytes(ancestor)
			}
		}
		locations = append(locations, location)
	}
	path, err := roundconfig.ResolveWorktreeLocation(loaded.Config.Worktree.Location, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		path = loaded.Config.Worktree.Location
	}
	inspect("worktree location", path, err)
	for _, cache := range loaded.Config.Worktree.SharedCaches {
		path, err := roundconfig.ResolveSharedCachePath(cache.Path, os.TempDir(), loaded.HomeDir, os.Getuid())
		if err != nil {
			path = cache.Path
		}
		inspect("Shared Build Cache", path, err)
	}
	inspect("temporary directory", os.TempDir(), nil)
	return locations
}

func defaultDoctorDiskSpaceResults(ctx context.Context, loaded roundconfig.Loaded) []CheckResult {
	var findings []readinessFinding
	low := false
	for _, location := range inspectFreeSpace(ctx, loaded) {
		if location.err != nil {
			findings = append(findings, readinessFinding{
				Code: "DR-DISK-UNREAD", Status: CheckStatusWarn,
				Text: fmt.Sprintf("could not read free space for %s %s: %v", location.label, location.path, location.err),
			})
		} else if location.available < doctorDiskSpaceThreshold {
			low = true
			findings = append(findings, readinessFinding{
				Code: "DR-DISK-LOW", Status: CheckStatusWarn,
				Text: fmt.Sprintf("%s %s has %d MiB free, below 10240 MiB", location.label, location.path, location.available/(1024*1024)),
			})
		}
	}
	if len(findings) == 0 {
		return nil
	}
	result := readinessResult(HealthCheckDiskSpace, "", findings)
	if low {
		result.NextAction = "roundfix gc"
	}
	return []CheckResult{result}
}

func doctorDiskSpaceAncestor(path string) (string, error) {
	for {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				return path, nil
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("stat disk-space location %q: %w", path, err)
		}
		parent := filepath.Dir(path)
		if parent == path {
			return "", fmt.Errorf("find existing ancestor of %q: %w", path, os.ErrNotExist)
		}
		path = parent
	}
}
