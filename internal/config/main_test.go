// Suite: Project Config repository boundary.
// Invariant: config tests leave every non-ignored repository path unchanged.
// Boundary IN: the config package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package config

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
