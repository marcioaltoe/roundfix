// Suite: daemon repository boundary.
// Invariant: daemon tests leave every non-ignored repository path unchanged.
// Boundary IN: the daemon package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package daemon

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
