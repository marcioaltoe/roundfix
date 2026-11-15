// Suite: Spec Check Worktree Bootstrap
// Invariant: the authoring probe prepares its tree once as the Daemon does and never probes after failed preparation.
// Boundary IN: CLI dispatch, temporary Git repositories, configured copy and bootstrap, and text/JSON streams and exit codes.
// Boundary OUT: Daemon scheduling and bootstrap internals, covered by their owning suites.
package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes(t *testing.T) {
	t.Parallel()
	const bootstrap = "test -f local.env && touch .deps-installed"
	for _, tc := range []struct{ name, format string }{
		{"surface_transcript_1", "text"},
		{"json_names_the_bootstrap", "json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, repo := newSpecCheckVerificationWorkspace(t, []string{
				"test -f .deps-installed", "test -f local.env", "test -f output-the-task-writes",
			})
			mustWrite(t, filepath.Join(repo, ".gitignore"), ".deps-installed\nlocal.env\n")
			mustWrite(t, filepath.Join(repo, "local.env"), "local fixture\n")
			commitSpecCheckBootstrapConfig(t, repo, bootstrap, []string{"local.env"})
			code, output, document := probeProvenanceCLI(t, tc.format)
			if code != exitRunFailed {
				t.Fatalf("exit=%d, want 1; stdout=%s", code, output)
			}
			if tc.format == "text" {
				assertSpecCheckBootstrapText(t, output, "Spec clean\nNo findings. Authored Verification commands executed: 3.\n",
					"Verification tree: HEAD\n"+
						"Verification tree bootstrap: \"test -f local.env && touch .deps-installed\"\n"+
						"- task_01: vacuous — \"test -f .deps-installed\" (exited zero before work)\n"+
						"- task_01: vacuous — \"test -f local.env\" (exited zero before work)\n"+
						"- task_01: honest — \"test -f output-the-task-writes\" (exited non-zero before work)\n")
				return
			}
			var got struct {
				Verification struct {
					Tree      string `json:"tree"`
					Bootstrap string `json:"bootstrap"`
				} `json:"verification"`
			}
			if err := json.Unmarshal([]byte(output), &got); err != nil {
				t.Fatal(err)
			}
			if got.Verification.Tree != "HEAD" || got.Verification.Bootstrap != bootstrap {
				t.Fatalf("verification=%+v", got.Verification)
			}
			commands := document.Verification.Commands
			if len(commands) != 3 || commands[0].Verdict != "vacuous" || commands[1].Verdict != "vacuous" || commands[2].Verdict != "honest" {
				t.Fatalf("commands=%+v", commands)
			}
		})
	}
}

func TestSpecCheckStopsWhenTheWorktreeBootstrapFails(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"text", "json"} {
		t.Run("surface_transcript_2_"+format, func(t *testing.T) {
			marker := filepath.Join(t.TempDir(), "command-ran")
			_, repo := newSpecCheckVerificationWorkspace(t, []string{fmt.Sprintf("touch %q; exit 1", marker)})
			commitSpecCheckBootstrapConfig(t, repo, "echo bootstrap-broke >&2; exit 7", nil)
			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, t.Context(), []string{"spec", "check", "clean", "--run-verification", "--format", format}, &stdout, &stderr)
			const wantStderr = "roundfix: spec check failed: prepare Verification tree at HEAD: worktree bootstrap failed: echo bootstrap-broke >&2; exit 7: exit status 7; bootstrap work may have been applied\n" +
				"No Verification command was probed: a Run stops at the same Worktree Bootstrap before any Agent work.\n" +
				"Worktree Bootstrap output tail:\n" +
				"bootstrap-broke\n"
			if code != exitRunFailed || stdout.String() != "" || stderr.String() != wantStderr {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=1 stdout=empty stderr=%q", code, stdout.String(), stderr.String(), wantStderr)
			}
			if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("Verification command ran after failed bootstrap: %v", err)
			}
		})
	}
}

func TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand(t *testing.T) {
	t.Parallel()
	t.Run("nothing_to_probe", func(t *testing.T) {
		for _, scenario := range []string{"completed_Task_and_QA_gate", "Task_without_commands", "Spec_without_Task_Graph"} {
			t.Run(scenario, func(t *testing.T) {
				counter := filepath.Join(t.TempDir(), "bootstrap-count")
				var repo string
				switch scenario {
				case "completed_Task_and_QA_gate":
					_, repo = newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
						{status: spec.StatusCompleted, commands: []string{"true"}},
						{status: spec.StatusPending, taskType: spec.TaskTypeQA},
					})
				case "Task_without_commands":
					_, repo = newSpecCheckVerificationWorkspace(t, nil)
				case "Spec_without_Task_Graph":
					_, repo = newSpecCheckWorkspace(t, "clean")
				}
				commitSpecCheckBootstrapConfig(t, repo, fmt.Sprintf("echo ran >> %q", counter), nil)
				var stdout, stderr bytes.Buffer
				code := runCLIContext(t, t.Context(), []string{"spec", "check", "clean", "--run-verification"}, &stdout, &stderr)
				if scenario == "Task_without_commands" {
					// Loading rejects this graph before any command can be probed.
					if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), "has no Verification commands") {
						t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
					}
				} else if code != exitOK || stderr.Len() != 0 || strings.Contains(stdout.String(), "Verification tree bootstrap:") {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
				}
				if _, err := os.Stat(counter); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("bootstrap ran without a probed command: %v", err)
				}
			})
		}
	})
	t.Run("two_Specs_and_two_Tasks", func(t *testing.T) {
		counter := filepath.Join(t.TempDir(), "bootstrap-count")
		_, repo := newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
			{status: spec.StatusPending, commands: []string{"test -f output-the-task-writes"}},
			{status: spec.StatusPending, commands: []string{"test -f another-output-the-task-writes"}},
		})
		second := filepath.Join(repo, "docs/specs/second")
		if err := copyDir(filepath.Join(repo, "docs/specs/clean"), second); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(second)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			path := filepath.Join(second, entry.Name())
			mustWrite(t, path, strings.ReplaceAll(strings.ReplaceAll(mustRead(t, path), "spec: clean", "spec: second"), "authorizations/clean.md", "authorizations/second.md"))
		}
		authorization := filepath.Join(repo, "docs/workflow/authorizations/clean.md")
		mustWrite(t, filepath.Join(repo, "docs/workflow/authorizations/second.md"), strings.ReplaceAll(mustRead(t, authorization), "clean", "second"))
		commitSpecCheckBootstrapConfig(t, repo, fmt.Sprintf("echo ran >> %q", counter), nil)
		var stdout, stderr bytes.Buffer
		code := runCLIContext(t, t.Context(), []string{"spec", "check", "clean", "second", "--run-verification"}, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 || strings.Count(stdout.String(), "Verification tree bootstrap:") != 2 || strings.Count(stdout.String(), ": honest — ") != 4 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if got := mustRead(t, counter); got != "ran\n" {
			t.Fatalf("bootstrap counter=%q, want one line", got)
		}
	})
}

func commitSpecCheckBootstrapConfig(t *testing.T, repo, bootstrap string, copyList []string) {
	t.Helper()
	var config strings.Builder
	fmt.Fprintf(&config, "worktree:\n  bootstrap: %q\n", bootstrap)
	if len(copyList) > 0 {
		config.WriteString("  copy:\n")
		for _, path := range copyList {
			fmt.Fprintf(&config, "    - %q\n", path)
		}
	}
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), config.String())
	gitImplement(t, repo, "add", "-A")
	gitImplement(t, repo, "commit", "-m", "test: configure Worktree Bootstrap")
}

func TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			_, _ = newSpecCheckVerificationWorkspace(t, []string{"test -f output-the-task-writes", "true"})
			code, output, _ := probeProvenanceCLI(t, format)
			if code != exitRunFailed {
				t.Fatalf("exit=%d, want 1; stdout=%s", code, output)
			}
			if format == "text" {
				assertSpecCheckBootstrapText(t, output, "Spec clean\nNo findings. Authored Verification commands executed: 2.\n",
					"Verification tree: HEAD\n"+
						"- task_01: honest — \"test -f output-the-task-writes\" (exited non-zero before work)\n"+
						"- task_01: vacuous — \"true\" (exited zero before work)\n")
				return
			}
			const suffix = `"verification":{"executed":true,"tree":"HEAD","commands":[{"task":"task_01","command":"test -f output-the-task-writes","verdict":"honest"},{"task":"task_01","command":"true","verdict":"vacuous"}],"uncommitted":[]}}` + "\n"
			if !strings.HasSuffix(output, suffix) {
				t.Fatalf("JSON changed; stdout=%s; want suffix=%s", output, suffix)
			}
		})
	}
}

// Only the fixture's Skipped block may vary in these Surface Transcripts.
func assertSpecCheckBootstrapText(t *testing.T, output, prefix, tail string) {
	t.Helper()
	if !strings.HasPrefix(output, prefix) {
		t.Fatalf("stdout=%s; want prefix=%s", output, prefix)
	}
	skipped, gotTail, found := strings.Cut(strings.TrimPrefix(output, prefix), "Verification tree: HEAD\n")
	if !found || "Verification tree: HEAD\n"+gotTail != tail {
		t.Fatalf("stdout=%s; want tail=%s", output, tail)
	}
	if skipped == "" {
		return
	}
	if !strings.HasPrefix(skipped, "Skipped:\n") || !strings.HasSuffix(skipped, "\n") {
		t.Fatalf("unexpected transcript block: %q", skipped)
	}
	for _, line := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(skipped, "Skipped:\n"), "\n"), "\n") {
		if !strings.HasPrefix(line, "  ") {
			t.Fatalf("unexpected line in fixture's Skipped block: %q", line)
		}
	}
}
