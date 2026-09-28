// Suite: tooling authorization repository boundary.
// Invariant: authorization tests leave every non-ignored repository path unchanged.
// Boundary IN: the authorization package test lifecycle.
// Boundary OUT: repository fingerprinting, owned by internal/suiteguard.
package authorization_test

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/suiteguard"
)

func TestMain(m *testing.M) {
	os.Exit(suiteguard.Main(m, filepath.Join("..", "..")))
}
