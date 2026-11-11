package cli

import (
	"context"
	"runtime"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/corebudget"
)

func verificationCoreBudget(ctx context.Context, config roundconfig.Config) corebudget.Budget {
	deps := commandDependenciesForContext(ctx)
	logicalCPUs := deps.logicalCPUs
	if logicalCPUs == nil {
		logicalCPUs = runtime.NumCPU
	}
	return corebudget.Compute(corebudget.Inputs{
		LogicalCPUs:  logicalCPUs(),
		MaxVerifying: config.Runs.MaxVerifying,
		MaxActive:    config.Runs.MaxActive,
		Environ:      deps.environ(),
	})
}
