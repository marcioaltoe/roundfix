package corebudget_test

import (
	"slices"
	"testing"

	"roundfix/internal/corebudget"
)

func TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name        string
		inputs      corebudget.Inputs
		want        corebudget.Budget
		description string
	}{
		{"default on twelve", corebudget.Inputs{LogicalCPUs: 12, MaxVerifying: 1, MaxActive: 4},
			corebudget.Budget{Cores: 6, LogicalCPUs: 12, Share: 6, Divisor: 1, Source: corebudget.SourceMaxVerifying, MaxVerifying: 1, MaxActive: 4},
			"half of 12 logical CPUs divided by runs.max_verifying 1"},
		{"two at once", corebudget.Inputs{LogicalCPUs: 12, MaxVerifying: 2, MaxActive: 4},
			corebudget.Budget{Cores: 3, LogicalCPUs: 12, Share: 6, Divisor: 2, Source: corebudget.SourceMaxVerifying, MaxVerifying: 2, MaxActive: 4},
			"half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"odd count rounds down", corebudget.Inputs{LogicalCPUs: 13, MaxVerifying: 1},
			corebudget.Budget{Cores: 6, LogicalCPUs: 13, Share: 6, Divisor: 1, Source: corebudget.SourceMaxVerifying, MaxVerifying: 1},
			"half of 13 logical CPUs divided by runs.max_verifying 1"},
		{"capacity off divides by max active", corebudget.Inputs{LogicalCPUs: 12, MaxActive: 4},
			corebudget.Budget{Cores: 1, LogicalCPUs: 12, Share: 6, Divisor: 4, Source: corebudget.SourceMaxActive, MaxActive: 4},
			"half of 12 logical CPUs divided by runs.max_active 4 while runs.max_verifying is 0"},
		{"both off is undivided", corebudget.Inputs{LogicalCPUs: 12},
			corebudget.Budget{Cores: 6, LogicalCPUs: 12, Share: 6, Divisor: 1, Source: corebudget.SourceUndivided},
			"half of 12 logical CPUs, undivided while runs.max_verifying and runs.max_active are 0"},
		{"one CPU keeps one", corebudget.Inputs{LogicalCPUs: 1, MaxVerifying: 1},
			corebudget.Budget{Cores: 1, LogicalCPUs: 1, Share: 1, Divisor: 1, Source: corebudget.SourceMaxVerifying, MaxVerifying: 1},
			"half of 1 logical CPUs divided by runs.max_verifying 1"},
		{"never below one", corebudget.Inputs{LogicalCPUs: 4, MaxVerifying: 5},
			corebudget.Budget{Cores: 1, LogicalCPUs: 4, Share: 2, Divisor: 5, Source: corebudget.SourceMaxVerifying, MaxVerifying: 5},
			"half of 4 logical CPUs divided by runs.max_verifying 5"},
		{"zero CPUs counts as one", corebudget.Inputs{},
			corebudget.Budget{Cores: 1, LogicalCPUs: 1, Share: 1, Divisor: 1, Source: corebudget.SourceUndivided},
			"half of 1 logical CPUs, undivided while runs.max_verifying and runs.max_active are 0"},
		{"negative CPUs counts as one", corebudget.Inputs{LogicalCPUs: -2},
			corebudget.Budget{Cores: 1, LogicalCPUs: 1, Share: 1, Divisor: 1, Source: corebudget.SourceUndivided},
			"half of 1 logical CPUs, undivided while runs.max_verifying and runs.max_active are 0"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := corebudget.Compute(test.inputs)
			if got != test.want {
				t.Errorf("Compute(%+v) = %+v, want %+v", test.inputs, got, test.want)
			}
			if got.Describe() != test.description {
				t.Errorf("Describe() = %q, want %q", got.Describe(), test.description)
			}
		})
	}
}

func TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name        string
		environ     []string
		cores       int
		source      string
		divisor     int
		description string
	}{
		{"positive whole number", []string{"GOMAXPROCS=4"}, 4, corebudget.SourceEnvironment, 1, "GOMAXPROCS 4 in the environment"},
		{"larger than the share", []string{"GOMAXPROCS=10"}, 10, corebudget.SourceEnvironment, 1, "GOMAXPROCS 10 in the environment"},
		{"last entry wins", []string{"GOMAXPROCS=10", "GOMAXPROCS=4"}, 4, corebudget.SourceEnvironment, 1, "GOMAXPROCS 4 in the environment"},
		{"zero is ignored", []string{"GOMAXPROCS=0"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"negative is ignored", []string{"GOMAXPROCS=-2"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"signed is ignored", []string{"GOMAXPROCS=+4"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"text is ignored", []string{"GOMAXPROCS=four"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"empty is ignored", []string{"GOMAXPROCS="}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"the Roundfix variable alone is ignored", []string{"ROUNDFIX_VERIFICATION_CORES=4"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"invalid last entry wins", []string{"GOMAXPROCS=4", "GOMAXPROCS=0"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"maximum whole number", []string{"GOMAXPROCS=2147483647"}, 2147483647, corebudget.SourceEnvironment, 1, "GOMAXPROCS 2147483647 in the environment"},
		{"overflow is ignored", []string{"GOMAXPROCS=2147483648"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"non ASCII digits are ignored", []string{"GOMAXPROCS=４"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"whitespace is ignored", []string{"GOMAXPROCS= 4"}, 3, corebudget.SourceMaxVerifying, 2, "half of 12 logical CPUs divided by runs.max_verifying 2"},
		{"leading zeros are digits", []string{"GOMAXPROCS=004"}, 4, corebudget.SourceEnvironment, 1, "GOMAXPROCS 4 in the environment"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := corebudget.Compute(corebudget.Inputs{LogicalCPUs: 12, MaxVerifying: 2, MaxActive: 4, Environ: test.environ})
			want := corebudget.Budget{Cores: test.cores, LogicalCPUs: 12, Share: 6, Divisor: test.divisor, Source: test.source, MaxVerifying: 2, MaxActive: 4}
			if got != want {
				t.Errorf("Compute with %q = %+v, want %+v", test.environ, got, want)
			}
			if got.Describe() != test.description {
				t.Errorf("Describe() = %q, want %q", got.Describe(), test.description)
			}
		})
	}
}

func TestTheCoreBudgetEnvironmentReplacesBothVariables(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		cores int
		want  []string
	}{
		{"budget four", 4, []string{"FIRST=1", "MIDDLE=a=b", "GOMAXPROCS_EXTRA=keep", "LAST=3", "GOMAXPROCS=4", "ROUNDFIX_VERIFICATION_CORES=4"}},
		{"no budget", 0, nil},
		{"negative budget", -1, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			environ := []string{"FIRST=1", "GOMAXPROCS=12", "MIDDLE=a=b", "ROUNDFIX_VERIFICATION_CORES=7", "GOMAXPROCS_EXTRA=keep", "GOMAXPROCS=8", "LAST=3"}
			original := slices.Clone(environ)
			want := test.want
			if test.cores < 1 {
				want = original
			}
			got := corebudget.Environment(environ, test.cores)
			if !slices.Equal(got, want) {
				t.Errorf("Environment() = %q, want %q", got, want)
			}
			if !slices.Equal(environ, original) {
				t.Fatalf("Environment changed input to %q, want %q", environ, original)
			}
			got[0] = "CHANGED=1"
			if !slices.Equal(environ, original) {
				t.Errorf("returned environment aliases input: %q", environ)
			}
		})
	}
}

func TestTheCoreBudgetConcurrencyUsesConfiguredCapacity(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name       string
		inputs     corebudget.Inputs
		activeRuns int
		want       int
	}{
		{"machine capacity wins", corebudget.Inputs{MaxVerifying: 2, MaxActive: 4}, 8, 2},
		{"max active fallback", corebudget.Inputs{MaxActive: 4}, 8, 4},
		{"active Runs fallback", corebudget.Inputs{}, 8, 8},
		{"no active Runs keeps one", corebudget.Inputs{}, 0, 1},
		{"negative active Runs keeps one", corebudget.Inputs{}, -2, 1},
		{"negative capacities use active Runs", corebudget.Inputs{MaxVerifying: -1, MaxActive: -2}, 3, 3},
		{"explicit budget keeps machine capacity", corebudget.Inputs{MaxVerifying: 2, MaxActive: 4, Environ: []string{"GOMAXPROCS=10"}}, 8, 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := corebudget.Compute(test.inputs).Concurrent(test.activeRuns)
			if got != test.want {
				t.Errorf("Concurrent(%d) = %d, want %d", test.activeRuns, got, test.want)
			}
		})
	}
}
