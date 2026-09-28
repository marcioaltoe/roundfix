// Suite: verification selection repository boundary.
// Invariant: verification-selection tests leave every non-ignored repository path unchanged.
// Boundary IN: the verifyselect package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package verifyselect_test

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
