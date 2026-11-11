// Package corebudget computes and exports the Verification Core Budget.
package corebudget

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Variable   = "ROUNDFIX_VERIFICATION_CORES"
	GoVariable = "GOMAXPROCS"
)

const (
	SourceEnvironment  = "environment"
	SourceMaxVerifying = "runs.max_verifying"
	SourceMaxActive    = "runs.max_active"
	SourceUndivided    = "undivided"
)

type Inputs struct {
	LogicalCPUs, MaxVerifying, MaxActive int
	Environ                              []string
}

type Budget struct {
	Cores, LogicalCPUs, Share, Divisor int
	Source                             string
	MaxVerifying, MaxActive            int
}

// Compute uses only the supplied CPU count, capacities and environment.
func Compute(in Inputs) Budget {
	budget := Budget{
		LogicalCPUs:  max(1, in.LogicalCPUs),
		Divisor:      1,
		Source:       SourceUndivided,
		MaxVerifying: in.MaxVerifying,
		MaxActive:    in.MaxActive,
	}
	budget.Share = max(1, budget.LogicalCPUs/2)

	var explicit string
	for _, entry := range in.Environ {
		if value, ok := strings.CutPrefix(entry, GoVariable+"="); ok {
			explicit = value
		}
	}
	digits := explicit != ""
	for i := 0; i < len(explicit); i++ {
		if explicit[i] < '0' || explicit[i] > '9' {
			digits = false
			break
		}
	}
	if digits {
		if cores, err := strconv.ParseInt(explicit, 10, 32); err == nil && cores > 0 {
			budget.Cores = int(cores)
			budget.Source = SourceEnvironment
			return budget
		}
	}

	if in.MaxVerifying > 0 {
		budget.Divisor = in.MaxVerifying
		budget.Source = SourceMaxVerifying
	} else if in.MaxActive > 0 {
		budget.Divisor = in.MaxActive
		budget.Source = SourceMaxActive
	}
	budget.Cores = max(1, budget.Share/budget.Divisor)
	return budget
}

// Describe identifies the source of the budget.
func (budget Budget) Describe() string {
	switch budget.Source {
	case SourceEnvironment:
		return fmt.Sprintf("GOMAXPROCS %d in the environment", budget.Cores)
	case SourceMaxVerifying:
		return fmt.Sprintf("half of %d logical CPUs divided by runs.max_verifying %d", budget.LogicalCPUs, budget.Divisor)
	case SourceMaxActive:
		return fmt.Sprintf("half of %d logical CPUs divided by runs.max_active %d while runs.max_verifying is 0", budget.LogicalCPUs, budget.Divisor)
	default:
		return fmt.Sprintf("half of %d logical CPUs, undivided while runs.max_verifying and runs.max_active are 0", budget.LogicalCPUs)
	}
}

// Concurrent returns the configured capacity, falling back to active Runs.
func (budget Budget) Concurrent(activeRuns int) int {
	if budget.MaxVerifying > 0 {
		return budget.MaxVerifying
	}
	if budget.MaxActive > 0 {
		return budget.MaxActive
	}
	return max(1, activeRuns)
}

// Environment copies environ and, for a positive budget, replaces both variables.
func Environment(environ []string, cores int) []string {
	result := make([]string, 0, len(environ)+2)
	for _, entry := range environ {
		name, _, _ := strings.Cut(entry, "=")
		if cores > 0 && (name == GoVariable || name == Variable) {
			continue
		}
		result = append(result, entry)
	}
	if cores > 0 {
		value := strconv.Itoa(cores)
		result = append(result, GoVariable+"="+value, Variable+"="+value)
	}
	return result
}
