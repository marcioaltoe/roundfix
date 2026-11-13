// Suite: Build Cache Bound Makefile wiring.
// Invariant: each top-level build or test bounds only its chosen cache once and fails open.
// Boundary IN: the real Makefile, compiled Go stub, and isolated temporary directories.
// Boundary OUT: real cache entries and real Go command execution by make.
package verifyselect_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func TestMakeBoundsItsOwnGoBuildCache(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	uid, err := exec.CommandContext(t.Context(), "id", "-u").Output()
	if err != nil {
		t.Fatalf("id -u: %v", err)
	}
	cacheName := "roundfix-go-build-" + strings.TrimSpace(string(uid))
	goStub := testfixture.FixtureBinary(t, "bound-go-stub", `package main
import (
 "encoding/json"
 "fmt"
 "os"
)
func main() {
 if len(os.Args) > 1 && os.Args[1] == "env" { return }
 record, err := os.OpenFile(os.Getenv("ROUNDFIX_BOUND_RECORD"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
 if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(2) }
 err = json.NewEncoder(record).Encode(os.Args[1:])
 closeErr := record.Close()
 if err != nil || closeErr != nil { fmt.Fprintln(os.Stderr, err, closeErr); os.Exit(2) }
 if len(os.Args) > 3 && os.Args[1] == "run" && os.Args[3] == "./cmd/gocache-bound" && os.Getenv("ROUNDFIX_BOUND_FAIL") == "1" { os.Exit(1) }
}
`)
	var environment []string
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "GOCACHE", "TMPDIR", "MAKEFLAGS", "MFLAGS", "MAKEOVERRIDES", "MAKELEVEL", "GOCACHE_BOUND_MB", "GOCACHE_BOUND_IDLE", "ROUNDFIX_BOUND_RECORD", "ROUNDFIX_BOUND_FAIL":
			continue
		}
		environment = append(environment, entry)
	}
	for _, name := range []string{"default", "once per invocation", "fallback cache", "sub-make", "explicit cache", "explicit cache with a bound", "turned off", "idle override", "failing bound"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			runDir, err := filepath.EvalSymlinks(dir)
			if err != nil {
				t.Fatal(err)
			}
			tmpRoot := t.TempDir()
			record := filepath.Join(dir, "calls")
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "GO="+goStub, "GO_TEST_PARALLEL=16")
			command.Dir = dir
			command.Env = append(append([]string{}, environment...), "TMPDIR="+tmpRoot+"/", "ROUNDFIX_BOUND_RECORD="+record)
			targets := []string{"test"}
			cache, maxMB, idle := filepath.Join(tmpRoot, cacheName), "10240", "2h"
			skip := false
			switch name {
			case "once per invocation":
				targets = []string{"vet", "test"}
			case "fallback cache":
				// A regular file is unwritable as a directory even when the test runs as root.
				if err := os.WriteFile(filepath.Join(tmpRoot, cacheName), []byte("not a directory"), 0600); err != nil {
					t.Fatal(err)
				}
				cache = filepath.Join(runDir, ".gocache")
			case "sub-make":
				command.Env = append(command.Env, "MAKELEVEL=1")
				skip = true
			case "explicit cache", "explicit cache with a bound":
				cache = filepath.Join(dir, "explicit cache")
				command.Env = append(command.Env, "GOCACHE="+cache)
				if name == "explicit cache" {
					skip = true
				} else {
					maxMB = "4096"
					command.Env = append(command.Env, "GOCACHE_BOUND_MB="+maxMB)
				}
			case "turned off":
				command.Env = append(command.Env, "GOCACHE_BOUND_MB=0")
				skip = true
			case "idle override":
				idle = "3h"
				command.Env = append(command.Env, "GOCACHE_BOUND_IDLE="+idle)
			case "failing bound":
				command.Env = append(command.Env, "ROUNDFIX_BOUND_FAIL=1")
			}
			command.Args = append(command.Args, targets...)
			var stdout, stderr bytes.Buffer
			command.Stdout, command.Stderr = &stdout, &stderr
			if err := command.Run(); err != nil {
				t.Fatalf("make %v: %v\nstdout: %s\nstderr: %s", targets, err, &stdout, &stderr)
			}
			contents, err := os.ReadFile(record)
			if err != nil {
				t.Fatal(err)
			}
			var got [][]string
			decoder := json.NewDecoder(bytes.NewReader(contents))
			for decoder.More() {
				var args []string
				if err := decoder.Decode(&args); err != nil {
					t.Fatalf("decode invocation: %v\n%s", err, contents)
				}
				got = append(got, args)
			}
			var want [][]string
			if !skip {
				want = append(want, []string{"run", "-buildvcs=false", "./cmd/gocache-bound", "-dir", cache, "-max-mb", maxMB, "-idle", idle})
			}
			if name == "once per invocation" {
				want = append(want, []string{"vet", "./..."})
			}
			want = append(want, []string{"test", "-parallel", "16", "./..."})
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("invocations = %q, want %q\nstdout: %s\nstderr: %s", got, want, &stdout, &stderr)
			}
			wantStderr := ""
			if name == "failing bound" {
				wantStderr = "gocache-bound: warning: the Build Cache Bound did not run (exit 1); the build continues\n"
			}
			if stderr.String() != wantStderr {
				t.Fatalf("stderr = %q, want %q", stderr.String(), wantStderr)
			}
		})
	}
}
