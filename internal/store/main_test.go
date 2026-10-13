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
	// A helper belongs to the spawning package's mark, not to a new
	// package run. In particular, tree-middle hands off a waiting child.
	if os.Getenv("ROUNDFIX_OWNER_PROCESS_HELPER") != "" {
		os.Exit(m.Run())
	}
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
