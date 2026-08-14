// Suite: Spec audit repository boundary.
// Invariant: specaudit tests leave every non-ignored repository path unchanged.
// Boundary IN: the specaudit package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package specaudit

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
