// Suite: Go environment build cache Makefile wiring.
// Invariant: make follows a usable configured cache and warns once before a small fallback.
// Boundary IN: the real Makefile, isolated Go environment files, and a compiled recording stub.
// Boundary OUT: host Go environment files, real cache entries, network, and cache trimming.
package verifyselect_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func TestMakeFollowsTheGoBuildCacheSetWithGoEnv(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	realGo, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := exec.CommandContext(t.Context(), "id", "-u").Output()
	if err != nil {
		t.Fatalf("id -u: %v", err)
	}
	cacheName := "roundfix-go-build-" + strings.TrimSpace(string(uid))
	goStub := testfixture.FixtureBinary(t, "go-env-stub", `package main
import (
 "encoding/json"
 "fmt"
 "os"
 "os/exec"
)
func main() {
 if len(os.Args) > 1 && os.Args[1] == "env" {
  command := exec.Command(os.Getenv("ROUNDFIX_CACHE_REAL_GO"), os.Args[1:]...)
  command.Stdout, command.Stderr = os.Stdout, os.Stderr
  if err := command.Run(); err != nil { os.Exit(1) }
  return
 }
 record, err := os.OpenFile(os.Getenv("ROUNDFIX_CACHE_RECORD"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
 if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(2) }
 args := append(os.Args[1:], "GOCACHE="+os.Getenv("GOCACHE"))
 err = json.NewEncoder(record).Encode(args)
 closeErr := record.Close()
 if err != nil || closeErr != nil { fmt.Fprintln(os.Stderr, err, closeErr); os.Exit(2) }
}
`)
	var environment []string
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "GOCACHE", "GOENV", "TMPDIR", "MAKEFLAGS", "MFLAGS", "MAKEOVERRIDES", "MAKELEVEL", "GOCACHE_BOUND_MB", "GOCACHE_BOUND_IDLE", "ROUNDFIX_CACHE_REAL_GO", "ROUNDFIX_CACHE_RECORD":
			continue
		}
		environment = append(environment, entry)
	}
	for _, name := range []string{"set and writable", "set and not mounted", "set and read-only", "set to a file", "set with a bound", "environment wins", "nothing set"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			tmpRoot := t.TempDir()
			setPath := filepath.Join(dir, "external disk", "cache", "go-build")
			fallback := filepath.Join(tmpRoot, cacheName)
			cache, maxMB, wantStderr := setPath, "10240", ""
			if name == "set and not mounted" {
				if err := os.Mkdir(filepath.Join(dir, "Volumes"), 0700); err != nil {
					t.Fatal(err)
				}
				setPath = filepath.Join(dir, "Volumes", "missing disk", "cache", "go-build")
			} else if name == "set to a file" {
				if err := os.MkdirAll(filepath.Dir(setPath), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(setPath, []byte("not a directory"), 0600); err != nil {
					t.Fatal(err)
				}
			} else if name != "nothing set" {
				if err := os.MkdirAll(setPath, 0700); err != nil {
					t.Fatal(err)
				}
				if name == "set and read-only" {
					t.Cleanup(func() {
						if err := os.Chmod(setPath, 0700); err != nil {
							t.Error(err)
						}
					})
					if err := os.Chmod(setPath, 0555); err != nil {
						t.Fatal(err)
					}
				}
			}
			goenv := filepath.Join(dir, "go-env")
			setting := "GOCACHE=" + setPath + "\n"
			if name == "nothing set" {
				setting = ""
				cache = fallback
			}
			if err := os.WriteFile(goenv, []byte(setting), 0600); err != nil {
				t.Fatal(err)
			}
			record := filepath.Join(dir, "calls")
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "test", "GO="+goStub, "GO_TEST_PARALLEL=16")
			command.Dir = dir
			command.Env = append(append([]string{}, environment...), "TMPDIR="+tmpRoot+"/", "GOENV="+goenv, "ROUNDFIX_CACHE_REAL_GO="+realGo, "ROUNDFIX_CACHE_RECORD="+record)
			switch name {
			case "set and not mounted", "set and read-only", "set to a file":
				cache, maxMB = fallback, "2048"
				wantStderr = fmt.Sprintf("gocache: warning: GOCACHE set with go env (%s) is not a writable directory; using %s with a 2048 MiB default bound\n", setPath, fallback)
			case "set with a bound":
				maxMB = "4096"
				command.Env = append(command.Env, "GOCACHE_BOUND_MB="+maxMB)
			case "environment wins":
				cache = filepath.Join(dir, "explicit cache")
				command.Env = append(command.Env, "GOCACHE="+cache)
			}
			var stdout, stderr bytes.Buffer
			command.Stdout, command.Stderr = &stdout, &stderr
			if err := command.Run(); err != nil {
				t.Fatalf("make test: %v\nstdout: %s\nstderr: %s", err, &stdout, &stderr)
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
			if name != "environment wins" {
				want = append(want, []string{"run", "-buildvcs=false", "./cmd/gocache-bound", "-dir", cache, "-max-mb", maxMB, "-idle", "2h", "GOCACHE=" + cache})
			}
			want = append(want, []string{"test", "-parallel", "16", "./...", "GOCACHE=" + cache})
			if !reflect.DeepEqual(got, want) {
				t.Errorf("invocations = %q, want %q\nstdout: %s\nstderr: %s", got, want, &stdout, &stderr)
			}
			if stderr.String() != wantStderr {
				t.Errorf("stderr = %q, want %q", stderr.String(), wantStderr)
			}
			if name == "set and writable" {
				entries, err := os.ReadDir(setPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(entries) != 0 {
					t.Errorf("configured cache retains write probes: %v", entries)
				}
			}
			if name == "set and not mounted" {
				if _, err := os.Stat(filepath.Join(dir, "Volumes", "missing disk")); !os.IsNotExist(err) {
					t.Errorf("missing disk path must stay absent, stat error = %v", err)
				}
			}
		})
	}
}
