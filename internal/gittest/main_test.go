// Suite: Git test-helper repository boundary.
// Invariant: gittest tests leave every non-ignored repository path unchanged.
// Boundary IN: the gittest package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package gittest

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
