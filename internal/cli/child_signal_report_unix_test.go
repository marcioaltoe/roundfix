//go:build darwin || linux

// Boundary: the fixture signals only its own PID in a temporary home and directory.
package cli

import (
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild(t *testing.T) {
	t.Parallel()
	binary := testfixture.FixtureBinary(t, "self-killing-roundfix", `package main
import ("fmt"; "os"; "syscall")
func main() {
 fmt.Fprint(os.Stderr, "fixture diagnostic")
 if err := syscall.Kill(os.Getpid(), syscall.SIGKILL); err != nil { panic(err) }
 panic("self signal returned")
}`)
	_, stderr, code := runRoundfixBinaryMacro(t, binary, t.TempDir(), t.TempDir(), t.TempDir(), "", nil)
	if code != -1 {
		t.Fatalf("child exit = %d, want -1; stderr=%q", code, stderr)
	}
	lines := strings.Split(strings.TrimSuffix(stderr, "\n"), "\n")
	const want = "child was killed by signal SIGKILL (9), no core dump"
	if got := lines[len(lines)-1]; got != want {
		t.Fatalf("last stderr line = %q, want %q; stderr=%q", got, want, stderr)
	}
	if stderr != "fixture diagnostic\n"+want+"\n" {
		t.Fatalf("stderr = %q, want fixture diagnostic followed by the signal note", stderr)
	}
}
