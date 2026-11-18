package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/baseline"
	"roundfix/internal/gittest"
	"roundfix/internal/phaserecord"
)

const phaseTranscript6 = `Usage:
  roundfix phase start <phase> [--spec <slug>] [--model <identifier>] [--effort <level>] [--author-subagent]
  roundfix phase model <phase> --spec <slug>
  roundfix phase fallback <phase> --spec <slug> --error <message>

Commands:
  start     Record in the Phase Record that this session starts a phase
  model     Print the model and effort a phase's Author Subagent starts on
  fallback  Decide and record an Authoring Fallback from the error an Author Subagent ended on

Phases:
  follow-up, idea, prd, techspec, tasks, delivery, release
  model and fallback take idea, prd, techspec or tasks

Options:
  --spec             Spec folder slug the phase works on; required for every phase
                     but follow-up and release
  --model            Model identifier the harness states for this session; recorded
                     for idea, prd, techspec and tasks
  --effort           Reasoning effort; used when CLAUDE_EFFORT is unset
  --author-subagent  The phase's Author Subagent runs the phase
  --error            Error the Author Subagent ended on; required for fallback
`

func phaseModelEnvironment(t *testing.T, session, effort string) {
	t.Helper()
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.getenv = func(key string) string {
			switch key {
			case "CLAUDE_CODE_SESSION_ID":
				return session
			case "CLAUDE_EFFORT":
				return effort
			// Configuration of a model is never evidence of the model running.
			case "ANTHROPIC_MODEL":
				return "configured-model"
			default:
				return ""
			}
		}
	})
}

func TestPhaseSurfaceTranscript6(t *testing.T) {
	t.Parallel()
	home, _ := phaseFixture(t, "")
	phaseModelEnvironment(t, "", "")
	phaseCall(t, []string{"phase", "--help"}, phaseTranscript6, "", 0)
	assertNoPhaseRecord(t, home)
}

func TestPhaseSurfaceTranscript7(t *testing.T) {
	t.Parallel()
	home, _ := phaseFixture(t, "3f9a1c20")
	phaseModelEnvironment(t, "3f9a1c20", "high")
	phaseCall(t, []string{"phase", "start", "prd", "--spec", "0413-a-ledger-export", "--model", "claude-fable-5-2", "--author-subagent"}, "Recorded prd of 0413-a-ledger-export for claude session 3f9a1c20: model claude-fable-5-2 (reported), effort high (environment), run by the Author Subagent\n", "", 0)
	if len(phaseDay(t, home)) != 1 {
		t.Fatal("expected one entry")
	}
}

type phaseModelValues struct {
	RepositoryID string `json:"repository_id"`
	Model        string `json:"model"`
	ModelSource  string `json:"model_source"`
	Effort       string `json:"effort"`
	EffortSource string `json:"effort_source"`
	Runner       string `json:"runner"`
}

func phaseModelEntry(t *testing.T, home string) phaseModelValues {
	t.Helper()
	lines := phaseDay(t, home)
	if len(lines) != 1 {
		t.Fatalf("lines = %v", lines)
	}
	var values phaseModelValues
	if err := json.Unmarshal([]byte(lines[0]), &values); err != nil {
		t.Fatal(err)
	}
	entries, skipped, err := phaserecord.Phases(home, time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC))
	if err != nil || skipped != 0 || len(entries) != 1 {
		t.Fatalf("read = %+v, %d, %v", entries, skipped, err)
	}
	entry := entries[0]
	read := phaseModelValues{entry.RepositoryID, entry.Model, entry.ModelSource, entry.Effort, entry.EffortSource, entry.Runner}
	if read != values {
		t.Fatalf("read = %+v, written = %+v", read, values)
	}
	return values
}

func TestPhaseStartRecordsTheModelAndEffortOfEachAuthoringPhase(t *testing.T) {
	t.Parallel()
	const invalidModel = "roundfix: warning: phase start: --model is not a model identifier of at most 200 characters; the model is recorded as unknown\n"
	const invalidEnvironment = "roundfix: warning: phase start: CLAUDE_EFFORT is not an effort level; the effort is recorded as unknown\n"
	const invalidEffort = "roundfix: warning: phase start: --effort is not an effort level; the effort is recorded as unknown\n"
	for _, phase := range []string{"idea", "prd", "techspec", "tasks"} {
		t.Run(phase, func(t *testing.T) {
			t.Parallel()
			for _, tc := range []struct {
				name                                                      string
				flags                                                     []string
				environment                                               string
				noSession                                                 bool
				model, modelSource, effort, effortSource, runner, warning string
			}{
				{"reported_model_and_environment", []string{"--model", "fixture-model", "--effort", "low"}, "high", false, "fixture-model", "reported", "high", "environment", "session", ""},
				{"effort_by_flag_only", []string{"--model", "fixture-model", "--effort", "xhigh"}, "", false, "fixture-model", "reported", "xhigh", "reported", "session", ""},
				{"no_reported_model", nil, "high", false, "", "unknown", "high", "environment", "session", ""},
				{"no_values", nil, "", false, "", "unknown", "", "unknown", "session", ""},
				{"no_session", []string{"--model", "fixture-model", "--effort", "low", "--author-subagent"}, "high", true, "", "unknown", "", "unknown", "unknown", ""},
				{"model_of_300_characters", []string{"--model", strings.Repeat("a", 300)}, "high", false, "", "unknown", "high", "environment", "session", invalidModel},
				{"malformed_model", []string{"--model", "bad model"}, "", false, "", "unknown", "", "unknown", "session", invalidModel},
				{"empty_model", []string{"--model="}, "", false, "", "unknown", "", "unknown", "session", invalidModel},
				{"model_at_limit", []string{"--model", strings.Repeat("a", 200)}, "", false, strings.Repeat("a", 200), "reported", "", "unknown", "session", ""},
				{"model_identifier_punctuation", []string{"--model", "A1._:/@[]-"}, "", false, "A1._:/@[]-", "reported", "", "unknown", "session", ""},
				{"malformed_environment_no_fallback", []string{"--model", "fixture-model", "--effort", "high"}, "Not an effort", false, "fixture-model", "reported", "", "unknown", "session", invalidEnvironment},
				{"malformed_effort_flag", []string{"--effort", "not an effort"}, "", false, "", "unknown", "", "unknown", "session", invalidEffort},
				{"effort_over_limit", []string{"--effort", strings.Repeat("a", 33)}, "", false, "", "unknown", "", "unknown", "session", invalidEffort},
				{"effort_at_limit", []string{"--effort", "a" + strings.Repeat("_-", 15) + "0"}, "", false, "", "unknown", "a" + strings.Repeat("_-", 15) + "0", "reported", "session", ""},
				{"author_subagent", []string{"--model", "fixture-model", "--author-subagent"}, "high", false, "fixture-model", "reported", "high", "environment", "author-subagent", ""},
			} {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					session := phaseTestSession
					if tc.noSession {
						session = ""
					}
					home, _ := phaseFixture(t, session)
					phaseModelEnvironment(t, session, tc.environment)
					args := append([]string{"phase", "start", phase, "--spec", "0412-example"}, tc.flags...)
					display := func(value, source string) string {
						if source == "unknown" {
							return "unknown"
						}
						return value + " (" + source + ")"
					}
					runner := "the session"
					if tc.runner == "author-subagent" {
						runner = "the Author Subagent"
					}
					wantOut := fmt.Sprintf("Recorded %s of 0412-example for claude session %s: model %s, effort %s, run by %s\n", phase, session, display(tc.model, tc.modelSource), display(tc.effort, tc.effortSource), runner)
					if tc.noSession {
						wantOut = "Recorded " + phase + " of 0412-example without a harness session: the environment names none\n"
					}
					phaseCall(t, args, wantOut, tc.warning, 0)
					want := phaseModelValues{"", tc.model, tc.modelSource, tc.effort, tc.effortSource, tc.runner}
					if got := phaseModelEntry(t, home); got != want {
						t.Fatalf("entry = %+v, want %+v", got, want)
					}
				})
			}
		})
	}
}

func TestPhaseStartIgnoresTheModelFlagsOnAnotherPhase(t *testing.T) {
	t.Parallel()
	for _, phase := range []string{"release", "follow-up", "delivery"} {
		t.Run(phase, func(t *testing.T) {
			t.Parallel()
			for _, tc := range []struct {
				name  string
				flags []string
			}{
				{"all", []string{"--model", "bad model", "--effort", "bad effort", "--author-subagent"}},
				{"empty_model", []string{"--model="}},
				{"effort", []string{"--effort", "high"}},
				{"false_subagent", []string{"--author-subagent=false"}},
				{"no_flags", nil},
			} {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					home, repo := phaseFixture(t, phaseTestSession)
					phaseModelEnvironment(t, phaseTestSession, "high")
					wantErr := ""
					if len(tc.flags) != 0 {
						wantErr = "roundfix: warning: phase start: --model, --effort and --author-subagent apply to idea, prd, techspec and tasks; ignored for " + phase + "\n"
					}
					args := append([]string{"phase", "start", phase, "--spec", "0412-example"}, tc.flags...)
					phaseCall(t, args, "Recorded "+phase+" of 0412-example for claude session "+phaseTestSession+"\n", wantErr, 0)
					lines := phaseDay(t, home)
					if len(lines) != 1 || lines[0] != phaseExpectedLine(t, phase, "0412-example", repo, phaseTestSession) {
						t.Fatalf("nine-key entry = %v", lines)
					}
				})
			}
		})
	}
}

func TestPhaseStartRecordsTheRepositoryIdentityOfItsClone(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.Run(t, repo, "commit", "--allow-empty", "-m", "root")
	identity, err := baseline.InspectRepositoryIdentity(t.Context(), repo, nil)
	if err != nil || identity.Digest == "" {
		t.Fatalf("identity = %+v, %v", identity, err)
	}
	first, second := filepath.Join(t.TempDir(), "first"), filepath.Join(t.TempDir(), "second")
	gittest.Run(t, repo, "clone", "--no-hardlinks", repo, first)
	gittest.Run(t, repo, "clone", "--no-hardlinks", repo, second)
	worktree := filepath.Join(t.TempDir(), "linked")
	gittest.Run(t, repo, "worktree", "add", "--detach", worktree, "HEAD")
	// This address is only local config data; nothing fetches it or records it.
	gittest.Run(t, second, "remote", "set-url", "origin", "https://fixture-user:fixture-pass@example.invalid/repo.git")
	for _, tc := range []struct{ name, path, session string }{
		{"original", repo, phaseTestSession}, {"first_clone", first, phaseTestSession},
		{"second_clone_with_credentials", second, phaseTestSession}, {"linked_worktree", worktree, phaseTestSession},
		{"without_session", first, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home, _ := phaseFixture(t, tc.session)
			phaseModelEnvironment(t, tc.session, "")
			setCommandWorkDirForTest(t, tc.path)
			wantOut := "Recorded prd of 0412-example for claude session " + tc.session + ": model unknown, effort unknown, run by the session\n"
			if tc.session == "" {
				wantOut = "Recorded prd of 0412-example without a harness session: the environment names none\n"
			}
			phaseCall(t, []string{"phase", "start", "prd", "--spec", "0412-example"}, wantOut, "", 0)
			got := phaseModelEntry(t, home)
			if got.RepositoryID != identity.Digest {
				t.Fatalf("identity = %q, want %q", got.RepositoryID, identity.Digest)
			}
			if text := strings.Join(phaseDay(t, home), ""); strings.Contains(text, "fixture-user") || strings.Contains(text, "fixture-pass") || strings.Contains(text, "example.invalid") {
				t.Fatalf("entry contains a remote address: %s", text)
			}
		})
	}
	// Exercise the root help usage line without changing other help expectations.
	var help strings.Builder
	var diagnostics strings.Builder
	if code := runCLI(t, []string{"--help"}, &help, &diagnostics); code != 0 || diagnostics.Len() != 0 || !strings.Contains(help.String(), strings.Split(phaseTranscript6, "\n")[1]+"\n") {
		t.Fatalf("root help = %s, %s, exit %d", &help, &diagnostics, code)
	}
}
