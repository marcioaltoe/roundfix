package agent

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"roundfix/internal/testfixture"
)

func TestAgentNodeOptionsDropsAMissingPreload(t *testing.T) {
	t.Parallel()
	const missing = "/missing.cjs"
	const memory = "--max-old-space-size=4096"
	cases := []struct {
		name, value, want string
		dropped           []string
	}{
		{"require equals", "--require=" + missing + " " + memory, memory, []string{missing}},
		{"require space", "--require " + missing + " " + memory, memory, []string{missing}},
		{"short require", "-r " + missing + " " + memory, memory, []string{missing}},
		{"import URL", "--import file:///missing.cjs " + memory, memory, []string{missing}},
		{"import equals URL", "--import=file:///missing.cjs " + memory, memory, []string{missing}},
		{"import path", "--import " + missing + " " + memory, memory, []string{missing}},
		{"quoted missing path", `--require "/missing file.cjs" ` + memory, memory, []string{"/missing file.cjs"}},
		{"encoded URL", "--import=file:///missing%20file.cjs " + memory, memory, []string{"/missing file.cjs"}},
		{"all removed", "--require=" + missing + " -r /other.cjs", "", []string{missing, "/other.cjs"}},
		{"kept order", "--trace-warnings --require=" + missing + " " + memory, "--trace-warnings " + memory, []string{missing}},
		{"existing", "--require=/existing.cjs", "--require=/existing.cjs", nil},
		{"existing URL", "--import=file:///existing.cjs", "--import=file:///existing.cjs", nil},
		{"relative", "--require ./missing.cjs", "--require ./missing.cjs", nil},
		{"package", "-r package-name", "-r package-name", nil},
		{"relative file URL", "--import=file:relative.cjs", "--import=file:relative.cjs", nil},
		{"other URL", "--import=https://example.com/missing.cjs", "--import=https://example.com/missing.cjs", nil},
		{"quoted kept path", `--require "/existing file.cjs"`, `--require "/existing file.cjs"`, nil},
		{"escaped quotes", `--require "/existing \"file\".cjs"`, `--require "/existing \"file\".cjs"`, nil},
		{"escaped backslash", `--require "/existing \\file.cjs"`, `--require "/existing \\file.cjs"`, nil},
		{"spaces normalized", "  --trace-warnings   " + memory + " ", "--trace-warnings " + memory, nil},
		{"unknown option", "--unrelated=" + missing, "--unrelated=" + missing, nil},
		{"missing argument", "--require", "--require", nil},
		{"empty word", `--require ""`, `--require ""`, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kept, dropped, ok := agentNodeOptions(tc.value, func(path string) bool {
				if !filepath.IsAbs(path) {
					t.Fatalf("existence check for nonabsolute path %q", path)
				}
				return strings.HasPrefix(path, "/existing")
			})
			if !ok || kept != tc.want || !reflect.DeepEqual(dropped, tc.dropped) {
				t.Fatalf("filter(%q) = (%q, %v, %v), want (%q, %v, true)", tc.value, kept, dropped, ok, tc.want, tc.dropped)
			}
		})
	}
}

func TestAgentNodeOptionsKeepsWhatItCannotSplit(t *testing.T) {
	t.Parallel()
	for _, value := range []string{`--require "/missing.cjs`, `--require=/missing.cjs --title="unfinished`, `"trailing\`} {
		t.Run(value, func(t *testing.T) {
			kept, dropped, ok := agentNodeOptions(value, func(string) bool {
				t.Fatal("cannot inspect preloads before splitting succeeds")
				return false
			})
			if ok || kept != value || len(dropped) != 0 {
				t.Fatalf("unsplittable value = (%q, %v, %v), want (%q, [], false)", kept, dropped, ok, value)
			}
		})
	}
}

func TestACPXRunnerDropsAMissingPreloadWithOneNotice(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.cjs")
	output := filepath.Join(dir, "environment.txt")
	scriptPath := filepath.Join(dir, "acpx.sh")
	script := "#!/bin/sh\nprintf '%s|%s|%s|%s\\n' \"$NODE_OPTIONS\" \"${CODEX_PATH-unset}\" \"${CLAUDECODE-unset}\" \"$KEEP\" >> \"$OUTPUT\"\nprintf '%s\\n' '" + MinimumACPXVersion + "'\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	// Compile the launcher through the shared fixture helper. The shell
	// script stays data; the Go tool owns writing the executable.
	command := testfixture.FixtureBinary(t, "acpx", `package main
import (
	"fmt"
	"os"
	"os/exec"
)
func main() {
	args := append([]string{os.Getenv("ACPX_TEST_SCRIPT")}, os.Args[1:]...)
	command := exec.Command("/bin/sh", args...)
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`)
	environment := []string{"NODE_OPTIONS=--require=" + missing + " --max-old-space-size=4096", "CODEX_PATH=/dirty", "CLAUDECODE=1", "KEEP=kept", "OUTPUT=" + output, "ACPX_TEST_SCRIPT=" + scriptPath}
	original := append([]string(nil), environment...)
	var notices bytes.Buffer
	runner := &ACPXRunner{Command: command, Environment: environment, Notices: &notices}
	for range 2 {
		if err := runner.Probe(context.Background(), ProbeRequest{}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Repeat("--max-old-space-size=4096|unset|unset|kept\n", 2)
	if string(got) != want {
		t.Fatalf("child environment = %q, want %q", got, want)
	}
	wantNotice := fmt.Sprintf("roundfix: notice: NODE_OPTIONS preload %q does not exist; Roundfix left it out of the agent environment\n", missing)
	if notices.String() != wantNotice {
		t.Fatalf("notices = %q, want %q", notices.String(), wantNotice)
	}
	if !reflect.DeepEqual(environment, original) {
		t.Fatalf("explicit environment mutated: %v", environment)
	}
	// A second runner in this process must not repeat the path's notice.
	(&ACPXRunner{Environment: environment, Notices: &notices}).commandEnv(nil)
	if notices.String() != wantNotice {
		t.Fatalf("second runner repeated notice: %q", notices.String())
	}
}

func TestACPXRunnerNodeOptionsEnvironment(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.cjs")
	existing := filepath.Join(dir, "existing file.cjs")
	if err := os.WriteFile(existing, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		values  []string
		want    string
		present bool
	}{
		{"last entry wins", []string{"--require=" + missing, "--trace-warnings"}, "--trace-warnings", true},
		{"all removed including shadowed entry", []string{"--trace-warnings", "--require=" + missing}, "", false},
		{"existing quoted path", []string{`--require "` + existing + `"`}, `--require "` + existing + `"`, true},
		{"unbalanced untouched", []string{`--require "` + missing}, `--require "` + missing, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			env := []string{"KEEP=value"}
			for _, value := range tc.values {
				env = append(env, "NODE_OPTIONS="+value)
			}
			var notices bytes.Buffer
			got := (&ACPXRunner{Environment: env, Notices: &notices}).commandEnv(nil)
			if value := environmentValue(got, "NODE_OPTIONS"); value != tc.want {
				t.Fatalf("NODE_OPTIONS = %q, want %q", value, tc.want)
			}
			present := false
			for _, entry := range got {
				if strings.HasPrefix(entry, "NODE_OPTIONS=") {
					present = true
				}
			}
			if present != tc.present || environmentValue(got, "KEEP") != "value" {
				t.Fatalf("environment = %v", got)
			}
		})
	}
}

func TestACPXRunnerNodeOptionsProcessEnvironment(t *testing.T) {
	// Sequential: changes the process environment.
	value := "--require=" + filepath.Join(t.TempDir(), "missing.cjs") + " --trace-warnings"
	t.Setenv("NODE_OPTIONS", value)
	var notices bytes.Buffer
	got := (&ACPXRunner{Notices: &notices}).commandEnv(nil)
	if environmentValue(got, "NODE_OPTIONS") != "--trace-warnings" || os.Getenv("NODE_OPTIONS") != value {
		t.Fatal("process environment was not filtered independently")
	}
}

func TestACPXRunnerNodeOptionsConcurrentNotices(t *testing.T) {
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "missing.cjs")
	var notices bytes.Buffer
	runner := &ACPXRunner{Environment: []string{"NODE_OPTIONS=--require=" + missing + " --import=" + missing}, Notices: &notices}
	var group sync.WaitGroup
	for range 8 {
		group.Go(func() { runner.commandEnv(nil) })
	}
	group.Wait()
	if strings.Count(notices.String(), "\n") != 1 {
		t.Fatalf("concurrent repeated notices: %q", notices.String())
	}
}

func TestACPXEnvironmentDropsAMissingPreloadForAnyProcess(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "gone.cjs")
	var notices bytes.Buffer
	got := ACPXEnvironment([]string{"KEEP=1", "NODE_OPTIONS=--require=" + missing + " --max-old-space-size=4096"}, &notices)
	want := []string{"KEEP=1", "NODE_OPTIONS=--max-old-space-size=4096"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("environment = %q, want %q", got, want)
	}
	if !strings.Contains(notices.String(), missing) {
		t.Fatalf("notice = %q, want it to name %s", notices.String(), missing)
	}
}
