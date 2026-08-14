// Package testfixture provides compiled executables for tests that cross an
// operating-system process boundary.
package testfixture

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// FixtureBinary compiles source into name and returns the executable path.
// The Go tool owns the output descriptor, so the test process never forks
// while it holds the executable open for writing.
func FixtureBinary(t testing.TB, name string, source string) string {
	t.Helper()
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write %s fixture source: %v", name, err)
	}
	binaryPath := filepath.Join(dir, name)
	command := exec.Command("go", "build", "-buildvcs=false", "-o", binaryPath, sourcePath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("compile %s fixture: %v\n%s", name, err, output)
	}
	return binaryPath
}
