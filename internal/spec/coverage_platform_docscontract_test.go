//go:build docscontract

//verify:relevant dist/npm/platforms.json go.mod go.sum

package spec

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestCoverageCollectionMatchesTheToolchainOnThisPlatform(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..")
	actual, err := collectCoverageRecordFor(root, []string{runtime.GOOS})
	if err != nil {
		t.Fatal(err)
	}
	toolchain, err := collectToolchainCoverage(root)
	if err != nil {
		t.Fatal(err)
	}
	// The unchanged toolchain parser describes this host's one platform.
	toolchain.Platforms = []string{runtime.GOOS}
	comparison := compareCoverageRecords(toolchain, actual)
	for _, difference := range append(comparison.Regressions, comparison.Additions...) {
		t.Error(difference)
	}
	t.Logf("static collection and go test -list compared on %s", runtime.GOOS)
}
