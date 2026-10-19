package testfixture

import (
	"os/exec"
	"testing"
)

// RequireTool resolves an External Test Tool on PATH or skips naming the tool.
func RequireTool(t testing.TB, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("requires %s, which is not on PATH", name)
	}
	return path
}
