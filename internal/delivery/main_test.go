// Suite: Delivery Queue repository boundary.
// Invariant: delivery tests leave every non-ignored repository path unchanged.
// Boundary IN: the delivery package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package delivery

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
