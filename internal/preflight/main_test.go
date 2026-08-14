// Suite: preflight repository boundary.
// Invariant: preflight tests leave every non-ignored repository path unchanged.
// Boundary IN: the preflight package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package preflight

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
