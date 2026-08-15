// Suite: Spec consistency-check repository boundary.
// Invariant: speccheck tests leave every non-ignored repository path unchanged.
// Boundary IN: the speccheck package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package speccheck

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
