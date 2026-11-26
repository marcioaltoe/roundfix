// Suite: process-wide compiled fixture reuse.
// Invariant: each product or name/source key has one immutable executable.
// Boundary IN: repeated helper calls and executable results.
// Boundary OUT: toolchain caching and product command behavior.
package testfixture

import (
	"os"
	"os/exec"
	"testing"
)

func TestProductBinaryIsBuiltOncePerTestBinary(t *testing.T) {
	t.Parallel()
	first := ProductBinary(t)
	if second := ProductBinary(t); second != first {
		t.Fatalf("product paths differ: %s, %s", first, second)
	}
	if info, err := os.Stat(first); err != nil || info.Size() == 0 {
		t.Fatalf("product executable: %v, %v", info, err)
	}
	if output, err := exec.Command(first, "--version").CombinedOutput(); err != nil {
		t.Fatalf("execute product: %v\n%s", err, output)
	}
}

func TestFixtureBinaryIsBuiltOncePerNameAndSource(t *testing.T) {
	t.Parallel()
	const source = "package main\nimport \"fmt\"\nfunc main(){fmt.Print(\"first\")}"
	first := FixtureBinary(t, "stub", source)
	if second := FixtureBinary(t, "stub", source); second != first {
		t.Fatalf("fixture paths differ: %s, %s", first, second)
	}
	other := FixtureBinary(t, "stub", "package main\nimport \"fmt\"\nfunc main(){fmt.Print(\"other\")}")
	renamed := FixtureBinary(t, "renamed", source)
	if other == first || renamed == first {
		t.Fatal("different fixture keys shared an executable")
	}
	for path, want := range map[string]string{first: "first", other: "other", renamed: "first"} {
		if output, err := exec.Command(path).CombinedOutput(); err != nil || string(output) != want {
			t.Fatalf("execute %s: %q, %v; want %q", path, output, err, want)
		}
	}
}
