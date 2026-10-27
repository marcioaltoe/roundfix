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

func defaultDoctorDiskSpaceResults(ctx context.Context, loaded roundconfig.Loaded) []CheckResult {
	availableBytes := commandDependenciesForContext(ctx).availableBytes
	var findings []readinessFinding
	low := false
	inspect := func(label, path string, resolveErr error) {
		var available int64
		err := resolveErr
		if err == nil {
			var ancestor string
			ancestor, err = doctorDiskSpaceAncestor(path)
			if err == nil {
				available, err = availableBytes(ancestor)
			}
		}
		if err != nil {
			findings = append(findings, readinessFinding{
				Code: "DR-DISK-UNREAD", Status: CheckStatusWarn,
				Text: fmt.Sprintf("could not read free space for %s %s: %v", label, path, err),
			})
		} else if available < doctorDiskSpaceThreshold {
			low = true
			findings = append(findings, readinessFinding{
				Code: "DR-DISK-LOW", Status: CheckStatusWarn,
				Text: fmt.Sprintf("%s %s has %d MiB free, below 10240 MiB", label, path, available/(1024*1024)),
			})
		}
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
		if _, err := os.Stat(path); err == nil {
			return path, nil
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
