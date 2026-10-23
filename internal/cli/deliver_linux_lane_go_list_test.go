// Boundary: the real shell script tail and gate with compiled Go/Docker stand-ins.
// No container or external service is used.
package cli

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/testfixture"
)

const linuxLaneGoListToolSource = `package main
import ("fmt"; "os"; "strconv")
func main() {
 switch os.Args[1] {
 case "list":
  fmt.Print(os.Getenv("LIST_OUTPUT"))
  code, err := strconv.Atoi(os.Getenv("LIST_EXIT")); if err != nil { panic(err) }; os.Exit(code)
 case "test":
  if err := os.WriteFile(os.Getenv("TEST_RECORD"), []byte("ran"), 0600); err != nil { panic(err) }
  code, err := strconv.Atoi(os.Getenv("TEST_EXIT")); if err != nil { panic(err) }; os.Exit(code)
 default: panic("unexpected go command")
 }
}`

func TestLinuxLaneScriptNamesAFailedGoList(t *testing.T) {
	t.Parallel()
	start := strings.Index(linuxLaneScript, `pkgs="$(go list`)
	if start < 0 {
		t.Fatal("lane script has no package-selection go list")
	}
	shell := testfixture.RequireTool(t, "sh")
	name := "go"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	goBinary := testfixture.FixtureBinary(t, name, linuxLaneGoListToolSource)
	for _, tc := range []struct {
		name, listOutput, line       string
		listExit, testExit, wantExit int
		testRan                      bool
	}{
		{"go list failure", "", "roundfix-linux-lane: go list: exited 1", 1, 0, 4, false},
		{"go list signal", "", "roundfix-linux-lane: go list: exited 137", 137, 0, 4, false},
		{"go test failure", "fixture/pkg\n", "roundfix-linux-lane: packages: 1", 0, 1, 1, true},
		{"no package", "", "no package builds on linux", 0, 0, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := t.TempDir()
			gittest.InitRepo(t, repo, "-b", "main")
			record := filepath.Join(repo, "test-ran")
			cmd := exec.CommandContext(t.Context(), shell, "-ceu", linuxLaneScript[start:], "linux-lane", "./...")
			cmd.Dir = repo
			cmd.Env = []string{"PATH=" + filepath.Dir(goBinary), "LIST_OUTPUT=" + tc.listOutput, "LIST_EXIT=" + strconv.Itoa(tc.listExit), "TEST_EXIT=" + strconv.Itoa(tc.testExit), "TEST_RECORD=" + record}
			output, err := cmd.CombinedOutput()
			code := 0
			if err != nil {
				var exit *exec.ExitError
				if !errors.As(err, &exit) {
					t.Fatal(err)
				}
				code = exit.ExitCode()
			}
			if code != tc.wantExit || !slices.Contains(strings.Split(string(output), "\n"), tc.line) {
				t.Fatalf("exit=%d output=%q; want exit=%d whole line=%q", code, output, tc.wantExit, tc.line)
			}
			if tc.listExit == 0 && strings.Contains(string(output), "roundfix-linux-lane: go list:") {
				t.Fatalf("successful go list reported failure: %q", output)
			}
			_, err = os.Stat(record)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
			if ran := err == nil; ran != tc.testRan {
				t.Fatalf("go test ran=%v; want %v", ran, tc.testRan)
			}
		})
	}
}

func TestLinuxLaneGoListFailureIsNamed(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, tc := range []struct{ name, code, summary, blocker string }{
		{"change", "1", "Linux Lane change failure: go list exited 1", delivery.BlockerGateFailed + ": Linux Lane change failure: go list exited 1"},
		{"signal", "137", "Linux Lane container failure: go list exited 137: the process ended by signal 9", delivery.BlockerLinuxLaneContainer + ": go list exited 137: the process ended by signal 9"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26\n", "pkg/code.go": "package pkg\n"})
			f.commit(t, "pkg/testdata/input", "data")
			f.environment = withEnvValue(f.environment, "LANE_EXIT", "4")
			marker := "roundfix-linux-lane: go list: exited " + tc.code + "\n"
			f.environment = withEnvValue(f.environment, "LANE_OUTPUT", marker)
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			suffix := "; see " + f.logPath()
			if err != nil || result.Passed || result.Reason != tc.summary+suffix || result.Blocker != tc.blocker+suffix {
				t.Fatalf("gate=%+v err=%v; want reason=%q blocker=%q", result, err, tc.summary+suffix, tc.blocker+suffix)
			}
			if log := mustRead(t, f.logPath()); !strings.Contains(log, marker) || !strings.HasSuffix(log, "\n"+tc.summary+"\n") {
				t.Fatalf("lane log=%q; want marker and final summary=%q", log, tc.summary)
			}
		})
	}
}
