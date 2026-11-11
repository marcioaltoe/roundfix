package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func defaultDoctorVerificationCoresResult(ctx context.Context, loaded roundconfig.Loaded) CheckResult {
	budget := verificationCoreBudget(ctx, loaded.Config)
	result := CheckResult{
		Name:   HealthCheckVerificationCores,
		Status: CheckStatusOK,
		Detail: fmt.Sprintf("%d per Verification: %s", budget.Cores, budget.Describe()),
	}
	active, err := doctorActiveImplementRuns(ctx, loaded.HomeDir)
	if err != nil {
		result.Status = CheckStatusWarn
		result.Detail += fmt.Sprintf("; active Implement Runs unknown: %v", err)
		result.NextAction = "inspect the Run Database with roundfix runs list"
		return result
	}
	result.Detail += fmt.Sprintf("; %d active Implement Runs", active)
	concurrent := budget.Concurrent(active)
	total := int64(concurrent) * int64(budget.Cores)
	if total > int64(budget.LogicalCPUs) {
		result.Status = CheckStatusWarn
		result.Detail += fmt.Sprintf("; %d at once may take %d of %d logical CPUs", concurrent, total, budget.LogicalCPUs)
		result.NextAction = "lower GOMAXPROCS in the environment or set runs.max_verifying in User Config"
	}
	return result
}

func doctorActiveImplementRuns(ctx context.Context, homeDir string) (int, error) {
	if _, err := os.Stat(store.DatabasePath(homeDir)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, fmt.Errorf("stat Run Database: %w", err)
	}
	runStore, err := openGCStore(ctx, true, homeDir)
	if err != nil {
		return 0, err
	}
	runs, readErr := runStore.ActiveImplementRuns(ctx)
	closeErr := runStore.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("close Run Database reader: %w", closeErr)
	}
	return len(runs), errors.Join(readErr, closeErr)
}
