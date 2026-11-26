// Package testfixture provides compiled executables for tests that cross an
// operating-system process boundary.
package testfixture

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type binaryKey struct{ name, source string }
type binaryBuild struct {
	once   sync.Once
	path   string
	err    error
	output []byte
}

var binaryBuilds sync.Map

// ProductBinary builds the product once for this test binary.
func ProductBinary(t testing.TB) string {
	t.Helper()
	return cachedBinary(t, binaryKey{name: "roundfix"}, true)
}

// FixtureBinary builds each name and source pair once for this test binary.
// Cached executables are immutable: callers that change one must copy it first.
func FixtureBinary(t testing.TB, name string, source string) string {
	t.Helper()
	return cachedBinary(t, binaryKey{name: name, source: source}, false)
}

func cachedBinary(t testing.TB, key binaryKey, product bool) string {
	t.Helper()
	value, _ := binaryBuilds.LoadOrStore(struct {
		key     binaryKey
		product bool
	}{key, product}, &binaryBuild{})
	build := value.(*binaryBuild)
	build.once.Do(func() {
		// This directory belongs to the test process, not the first caller's test.
		// A t.TempDir would disappear while later tests still need the executable.
		executable, err := os.Executable()
		if err != nil {
			build.err = err
			return
		}
		dir, err := os.MkdirTemp(filepath.Dir(executable), "roundfix-test-binary-")
		if err != nil {
			build.err = err
			return
		}
		build.path = filepath.Join(dir, filepath.Base(key.name))
		if runtime.GOOS == "windows" && !strings.HasSuffix(build.path, ".exe") {
			build.path += ".exe"
		}
		target := filepath.Join(dir, "main.go")
		if product {
			target = "./cmd/roundfix"
		} else if err := os.WriteFile(target, []byte(key.source), 0o600); err != nil {
			build.err = err
			return
		}
		command := exec.Command("go", "build", "-buildvcs=false", "-o", build.path, target)
		if product {
			_, file, _, ok := runtime.Caller(0)
			if !ok {
				build.err = &os.PathError{Op: "locate product source", Path: target, Err: os.ErrNotExist}
				return
			}
			command.Dir = filepath.Join(filepath.Dir(file), "..", "..")
		}
		build.output, build.err = command.CombinedOutput()
	})
	if build.err != nil {
		t.Fatalf("compile %s fixture: %v\n%s", key.name, build.err, build.output)
	}
	return build.path
}
