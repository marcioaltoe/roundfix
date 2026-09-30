package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const doctorStorageFreeBytesThreshold int64 = 64 * 1024 * 1024

func defaultDoctorStorageResults(ctx context.Context, loaded roundconfig.Loaded) []CheckResult {
	databasePath := store.DatabasePath(loaded.HomeDir)
	if _, err := os.Stat(databasePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []CheckResult{doctorStorageOKResult("no Run Database")}
		}
		return []CheckResult{doctorStoragePartialResult(fmt.Errorf("stat Run Database: %w", err))}
	}

	runStore, err := store.OpenReader(ctx, loaded.HomeDir)
	if err != nil {
		return []CheckResult{doctorStoragePartialResult(err)}
	}
	result := readDoctorStorageResult(ctx, runStore, loaded, time.Now().UTC())
	if closeErr := runStore.Close(); closeErr != nil {
		if result.Status == CheckStatusPartial {
			result.Detail += fmt.Sprintf("; could not close Run Database reader: %v", closeErr)
		} else {
			result = doctorStoragePartialResult(fmt.Errorf("close Run Database reader: %w", closeErr))
		}
	}
	return []CheckResult{result}
}

func readDoctorStorageResult(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, now time.Time) CheckResult {
	freeBytes, err := runStore.FreeBytes(ctx)
	if err != nil {
		return doctorStoragePartialResult(err)
	}

	retention := loaded.Config.Store.JournalRetention
	if retention == 0 {
		return doctorStorageResult(0, freeBytes, retention, true)
	}

	candidates, err := runStore.TerminalRunPruneCandidates(ctx, now.Add(-retention))
	if err != nil {
		return doctorStoragePartialResult(err)
	}

	artifactsInspected := strings.TrimSpace(loaded.GitRoot) != ""
	hasArtifactDir := func(string) (bool, error) { return false, nil }
	if artifactsInspected {
		artifactRoot, resolveErr := roundconfig.ResolveArtifactDirectory(
			loaded.Config.Defaults.ArtifactDir,
			loaded.GitRoot,
			loaded.HomeDir,
		)
		if resolveErr != nil {
			return doctorStoragePartialResult(fmt.Errorf("resolve Artifact Root: %w", resolveErr))
		}
		hasArtifactDir = func(runID string) (bool, error) {
			path, pathErr := gcRunArtifactPath(artifactRoot, runID)
			if pathErr != nil {
				return false, pathErr
			}
			info, statErr := os.Lstat(path)
			if errors.Is(statErr, os.ErrNotExist) {
				return false, nil
			}
			if statErr != nil {
				return false, fmt.Errorf("lstat Run artifact directory %q: %w", path, statErr)
			}
			if !info.IsDir() {
				return false, fmt.Errorf("Run artifact path %q is not a directory", path)
			}
			return true, nil
		}
	}

	reclaimableRunIDs, err := retentionReclaimable(candidates, hasArtifactDir)
	if err != nil {
		return doctorStoragePartialResult(err)
	}
	return doctorStorageResult(len(reclaimableRunIDs), freeBytes, retention, artifactsInspected)
}

func doctorStorageResult(reclaimableRuns int, freeBytes int64, retention time.Duration, artifactsInspected bool) CheckResult {
	runsClause := fmt.Sprintf("Runs reclaimable: %d", reclaimableRuns)
	if retention == 0 {
		runsClause = "Journal Retention 0 keeps every Run"
	} else if !artifactsInspected {
		runsClause += "; artifact directories not inspected outside a Git repository"
	}

	commands := make([]string, 0, 2)
	if reclaimableRuns > 0 {
		commands = append(commands, "roundfix gc")
	}
	if freeBytes >= doctorStorageFreeBytesThreshold {
		commands = append(commands, "roundfix gc compact")
	}
	if len(commands) > 0 {
		return CheckResult{
			Name:   HealthCheckStorage,
			Status: CheckStatusFound,
			Detail: fmt.Sprintf(
				"%s; Run Database free bytes: %d; next: %s",
				runsClause,
				freeBytes,
				strings.Join(commands, " and "),
			),
		}
	}
	return doctorStorageOKResult(fmt.Sprintf(
		"nothing to reclaim; %s; Run Database free bytes: %d",
		runsClause,
		freeBytes,
	))
}

func doctorStorageOKResult(detail string) CheckResult {
	return CheckResult{Name: HealthCheckStorage, Status: CheckStatusOK, Detail: detail}
}

func doctorStoragePartialResult(err error) CheckResult {
	return CheckResult{
		Name:   HealthCheckStorage,
		Status: CheckStatusPartial,
		Detail: fmt.Sprintf("could not read Run Database: %v", err),
	}
}
