// Suite: External Test Tools.
// Invariant: tests resolve tools outside the Base Test Tools through RequireTool.
// Boundary IN: literal os/exec tool lookups in internal, cmd and skills test files.
// Boundary OUT: paths, variable programs and non-test files.
package testfixture

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// Sequential: sets PATH
func TestRequireToolReturnsAToolOnPath(t *testing.T) {
	name := "stand-in"
	if runtime.GOOS == "windows" {
		name += ".exe"
		t.Setenv("PATHEXT", ".EXE")
	}
	binary := FixtureBinary(t, name, "package main\nfunc main() {}\n")
	t.Setenv("PATH", filepath.Dir(binary))
	if got := RequireTool(t, name); got != binary {
		t.Fatalf("RequireTool = %q, want %q", got, binary)
	}
}

// Sequential: sets PATH
func TestRequireToolSkipsNamingAnAbsentTool(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	var skipped *testing.T
	var message string
	continued := false
	t.Run("absent", func(t *testing.T) {
		skipped = t
		RequireTool(&toolSkipRecorder{TB: t, message: &message}, "rustfmt")
		continued = true
	})
	if !skipped.Skipped() || continued || message != "requires rustfmt, which is not on PATH" {
		t.Fatalf("skipped = %v, continued = %v, message = %q", skipped.Skipped(), continued, message)
	}
}

type toolSkipRecorder struct {
	testing.TB
	message *string
}

func (t *toolSkipRecorder) Skipf(format string, args ...any) {
	t.Helper()
	*t.message = fmt.Sprintf(format, args...)
	t.TB.Skipf(format, args...)
}
