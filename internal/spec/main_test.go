// Suite: Spec repository boundary.
// Invariant: spec tests leave every non-ignored repository path unchanged.
// Boundary IN: the spec package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package spec

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
