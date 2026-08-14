//go:build repocontract

// Suite: package suite-guard installation contract.
// Invariant: this package installs suiteguard.Main whenever its tests spawn a process.
// Boundary IN: this package's Go test sources and TestMain wiring.
// Boundary OUT: the repository-wide guarded-package inventory in internal/suiteguard.
package preflight

import (
	"testing"

	"roundfix/internal/suiteguardcontract"
)

func TestEverySpawningPackageInstallsTheSuiteGuard(t *testing.T) {
	suiteguardcontract.CheckCurrentPackage(t)
}
