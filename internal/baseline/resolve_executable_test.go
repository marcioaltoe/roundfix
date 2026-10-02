package baseline

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/testfixture"
)

func TestResolveExecutableNeverRunsTheCandidate(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "executed")
	candidate := testfixture.FixtureBinary(t, "tool", fmt.Sprintf(`package main
import "os"
func main() {
	if err := os.WriteFile(%q, []byte("executed"), 0o600); err != nil { panic(err) }
}
`, marker))
	dir := filepath.Dir(candidate)
	if err := os.Symlink(candidate, filepath.Join(dir, "linked")); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"tool", "linked"} {
		path, reason := ResolveExecutable(name, []string{dir})
		if path != filepath.Join(dir, name) || reason != "" {
			t.Fatalf("resolve %s = %q, %q", name, path, reason)
		}
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("candidate was executed")
	}
	if _, reason := ResolveExecutable("absent", []string{dir}); reason == "" {
		t.Fatal("absent executable accepted")
	}
	if err := os.Chmod(candidate, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, reason := ResolveExecutable("linked", []string{dir}); reason == "" {
		t.Fatal("non-executable link target accepted")
	}
}
