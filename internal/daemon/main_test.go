// Suite: daemon repository boundary.
// Invariant: daemon tests leave every non-ignored repository path unchanged.
// Boundary IN: the daemon package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	if active, err := runVerificationEnvironmentFixture(); active {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
