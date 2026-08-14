// Suite: Run Database repository boundary.
// Invariant: store tests leave every non-ignored repository path unchanged.
// Boundary IN: the store package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package store

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
