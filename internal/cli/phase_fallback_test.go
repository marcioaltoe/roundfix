package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/baseline"
	"roundfix/internal/gittest"
	"roundfix/internal/phaserecord"
)

const fallbackProjectConfig = `skills:
  authoring:
    claude:
      idea:
        model: fable
        reasoning_effort: xhigh
        fallback: {model: opus, reasoning_effort: xhigh}
      prd:
        model: fable
        reasoning_effort: xhigh
        fallback: {model: opus, reasoning_effort: xhigh}
      techspec:
        model: fable
        reasoning_effort: xhigh
        fallback: {model: opus, reasoning_effort: xhigh}
      tasks:
        model: opus
        reasoning_effort: high
        fallback: {model: sonnet, reasoning_effort: high}
`
const fallbackError = "You've hit your Fable limit · resets Nov 22 at 8pm"

func fallbackFixture(t *testing.T, configuration string) (string, string) {
	t.Helper()
	home, repo := phaseFixture(t, phaseTestSession)
	writeFallbackConfig(t, repo, configuration)
	return home, repo
}
func writeFallbackConfig(t *testing.T, repo, configuration string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(repo, ".roundfixrc.yml"), []byte(configuration), 0600); err != nil {
		t.Fatal(err)
	}
}
func fallbackArgs(command, phase, spec string) []string {
	args := []string{"phase", command, phase, "--spec", spec}
	if command == "fallback" {
		args = append(args, "--error", fallbackError)
	}
	return args
}
func fallbackRecordSnapshot(t *testing.T, home string) map[string]string {
	t.Helper()
	dir := filepath.Join(home, ".roundfix", "phase-record")
	files, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	result := map[string]string{}
	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		result[file.Name()] = string(data)
	}
	return result
}

func TestPhaseFallbackSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	home, repo := fallbackFixture(t, fallbackProjectConfig)
	for _, tc := range []struct {
		number          int
		args            []string
		out, diagnostic string
		code            int
	}{
		{1, []string{"phase", "--help"}, phaseTranscript6, "", 0},
		{2, fallbackArgs("fallback", "idea", "0412-example"), "Start roundfix-idea-author with model opus and effort xhigh, the fallback of fable: \"You've hit your Fable limit · resets Nov 22 at 8pm\"\n", "", 0},
		{3, fallbackArgs("model", "prd", "0412-example"), "Start roundfix-prd-author with model opus and effort xhigh, the fallback of fable: idea of 0412-example fell back when the Fable limit ran out.\n", "", 0},
		{4, fallbackArgs("model", "tasks", "0412-example"), "Start roundfix-tasks-author with no model and no effort: the declared model opus.\n", "", 0},
		{5, fallbackArgs("model", "prd", "0413-another"), "Start roundfix-prd-author with no model and no effort: the declared model fable.\n", "", 0},
		{6, []string{"phase", "fallback", "techspec", "--spec", "0412-example", "--error", "You've hit your weekly limit · resets Mon 12:00am"}, "No fallback for techspec of 0412-example: a weekly limit is shared by every model.\n", "", 3},
		{7, fallbackArgs("fallback", "tasks", "0412-example"), "No fallback for tasks of 0412-example: the Fable limit is not that of the declared model opus.\n", "", 3},
		{8, fallbackArgs("model", "delivery", "0412-example"), "", "roundfix: phase model failed: phase \"delivery\" has no Author Subagent; use idea, prd, techspec or tasks\nRun 'roundfix phase --help' for usage.\n", 2},
	} {
		t.Run(fmt.Sprintf("transcript_%d", tc.number), func(t *testing.T) {
			before := fallbackRecordSnapshot(t, home)
			phaseCall(t, tc.args, tc.out, tc.diagnostic, tc.code)
			if tc.number != 2 {
				if after := fallbackRecordSnapshot(t, home); !reflect.DeepEqual(before, after) {
					t.Fatalf("read/refusal changed day files: %v", after)
				}
				return
			}
			lines := phaseDay(t, home)
			if len(lines) != 1 {
				t.Fatalf("lines = %v", lines)
			}
			rows, n, err := phaserecord.Fallbacks(home)
			want := phaserecord.FallbackEntry{Phase: "idea", Spec: "0412-example", Repository: repo, Harness: "claude", Session: phaseTestSession, DeclaredModel: "fable", DeclaredEffort: "xhigh", Model: "opus", Effort: "xhigh", Limit: "Fable", Message: fallbackError}
			if err != nil || n != 0 || len(rows) != 1 || rows[0].FallbackEntry != want || !rows[0].Time.Equal(time.Date(2026, 11, 20, 13, 5, 11, 0, time.UTC)) {
				t.Fatalf("fallback = %+v, %d, %v", rows, n, err)
			}
		})
	}
}

func TestPhaseFallbackRefusesWithoutWriting(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, config, message, phase, reason string }{
		{"not_limit", fallbackProjectConfig, "server unavailable", "prd", "the error is not a usage limit"},
		{"session", fallbackProjectConfig, "You've hit your session limit", "prd", "a session limit is shared by every model"},
		{"weekly", fallbackProjectConfig, "You've hit your weekly limit", "prd", "a weekly limit is shared by every model"},
		{"spend", fallbackProjectConfig, "You've hit your monthly spend limit", "prd", "a monthly spend limit is not a model's usage limit"},
		{"no_model", "{}\n", fallbackError, "prd", "prd declares no model"},
		{"wrong_family", fallbackProjectConfig, fallbackError, "tasks", "the Fable limit is not that of the declared model opus"},
		{"no_fallback", "skills: {authoring: {claude: {prd: {model: fable}}}}\n", fallbackError, "prd", "prd declares no fallback"},
		{"fallback_same_family", "skills: {authoring: {claude: {prd: {model: fable, fallback: {model: fable-backup}}}}}\n", fallbackError, "prd", "the fallback fable-backup is a Fable model too"},
		{"limit_before_model", "{}\n", "not a limit", "prd", "the error is not a usage limit"},
		{"shared_before_model", "{}\n", "You've hit your weekly limit", "prd", "a weekly limit is shared by every model"},
		{"other_before_model", "{}\n", "You've hit your individual usage limit", "prd", "a individual usage limit is not a model's usage limit"},
		{"model_before_fallback", "skills: {authoring: {claude: {prd: {model: opus}}}}\n", fallbackError, "prd", "the Fable limit is not that of the declared model opus"},
		{"fallback_own_limit", fallbackProjectConfig, "You've hit your Opus limit", "prd", "the Opus limit is not that of the declared model fable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, _ := fallbackFixture(t, tc.config)
			args := fallbackArgs("fallback", tc.phase, "0412-example")
			args[len(args)-1] = tc.message
			phaseCall(t, args, "No fallback for "+tc.phase+" of 0412-example: "+tc.reason+".\n", "", 3)
			assertNoPhaseRecord(t, home)
			// A refusal also preserves existing day files byte for byte.
			if err := phaserecord.New(home, func() time.Time { return time.Date(2026, 11, 20, 1, 0, 0, 0, time.UTC) }).RecordFallback(t.Context(), phaserecord.FallbackEntry{Phase: "idea", Spec: "0400-existing"}); err != nil {
				t.Fatal(err)
			}
			before := fallbackRecordSnapshot(t, home)
			phaseCall(t, args, "No fallback for "+tc.phase+" of 0412-example: "+tc.reason+".\n", "", 3)
			if after := fallbackRecordSnapshot(t, home); !reflect.DeepEqual(before, after) {
				t.Fatal("refusal changed day files")
			}
		})
	}
	for _, command := range []string{"model", "fallback"} {
		t.Run(command+"_usage", func(t *testing.T) {
			t.Parallel()
			for _, tc := range []struct {
				name   string
				args   []string
				reason string
			}{
				{"missing_phase", []string{command}, "missing phase; use idea, prd, techspec or tasks"},
				{"non_author", []string{command, "release"}, `phase "release" has no Author Subagent; use idea, prd, techspec or tasks`},
				{"unknown_phase", []string{command, "unknown"}, `phase "unknown" has no Author Subagent; use idea, prd, techspec or tasks`},
				{"missing_spec", []string{command, "prd"}, `phase "prd" requires --spec <slug>`},
				{"bad_spec", []string{command, "prd", "--spec", "../0412-example"}, `invalid Spec slug "../0412-example"; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$`},
				{"empty_spec", []string{command, "prd", "--spec="}, `phase "prd" requires --spec <slug>`},
				{"unknown_flag", []string{command, "prd", "--unknown"}, "flag provided but not defined: -unknown"},
				{"flag_value", []string{command, "prd", "--spec"}, "flag needs an argument: -spec"},
				{"extra", []string{command, "prd", "--spec", "0412-example", "extra"}, `unexpected argument "extra"`},
				{"outside_repo", fallbackArgs(command, "prd", "0412-example")[1:], "phase " + command + " requires a git repository working tree"},
			} {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					home, _ := fallbackFixture(t, fallbackProjectConfig)
					if tc.name == "outside_repo" {
						setCommandWorkDirForTest(t, t.TempDir())
					}
					phaseCall(t, append([]string{"phase"}, tc.args...), "", "roundfix: phase "+command+" failed: "+tc.reason+"\nRun 'roundfix phase --help' for usage.\n", 2)
					assertNoPhaseRecord(t, home)
				})
			}
		})
	}
	for _, tc := range []struct {
		name, message string
		omit          bool
	}{
		{"missing_error", "", true}, {"empty_error", "", false}, {"whitespace_error", " \t\n", false}, {"oversize_error", strings.Repeat("é", 2049), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, _ := fallbackFixture(t, fallbackProjectConfig)
			args := []string{"phase", "fallback", "prd", "--spec", "0412-example"}
			if !tc.omit {
				args = append(args, "--error", tc.message)
			}
			phaseCall(t, args, "", "roundfix: phase fallback failed: --error must be non-empty and at most 4096 bytes\nRun 'roundfix phase --help' for usage.\n", 2)
			assertNoPhaseRecord(t, home)
		})
	}
	t.Run("configuration_error", func(t *testing.T) {
		t.Parallel()
		for _, command := range []string{"model", "fallback"} {
			t.Run(command, func(t *testing.T) {
				t.Parallel()
				home, repo := fallbackFixture(t, "skills: {authoring: {claude: {prd: {model: fable, fallback: {model: fable}}}}}\n")
				var out, diagnostic bytes.Buffer
				code := runCLI(t, fallbackArgs(command, "prd", "0412-example"), &out, &diagnostic)
				want := fmt.Sprintf("roundfix: phase %s failed: parse config %q: skills.authoring.claude.prd.fallback.model \"fable\" is the declared model; a fallback must name another model\nRun 'roundfix phase --help' for usage.\n", command, filepath.Join(repo, ".roundfixrc.yml"))
				if code != 2 || out.Len() != 0 || diagnostic.String() != want {
					t.Fatalf("exit %d stdout %q stderr %q want %q", code, out.String(), diagnostic.String(), want)
				}
				assertNoPhaseRecord(t, home)
			})
		}
	})
}

func TestPhaseModelKeepsTheFallbackForItsSpecAndModel(t *testing.T) {
	t.Parallel()
	home, repo := fallbackFixture(t, fallbackProjectConfig)
	// Give repositories distinct local root commits. No tag, remote ref or network.
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.Run(t, repo, "commit", "--allow-empty", "-m", "first repository")
	identity, err := baseline.InspectRepositoryIdentity(t.Context(), repo, nil)
	if err != nil || identity.Digest == "" {
		t.Fatalf("identity %+v %v", identity, err)
	}
	other := t.TempDir()
	gittest.InitRepo(t, other, "-b", "main")
	gittest.Run(t, other, "commit", "--allow-empty", "-m", "second repository")
	writeFallbackConfig(t, other, fallbackProjectConfig)
	clone := filepath.Join(t.TempDir(), "clone")
	gittest.Run(t, repo, "clone", "--no-hardlinks", repo, clone)
	writeFallbackConfig(t, clone, fallbackProjectConfig)
	now := time.Date(2026, 11, 20, 13, 5, 11, 0, time.UTC)
	// Newer entry first in file: selection must use the oldest applicable entry.
	for i, phase := range []phaserecord.Phase{"techspec", "idea"} {
		row := phaserecord.FallbackEntry{Phase: phase, Spec: "0412-example", Repository: repo, RepositoryID: identity.Digest, Limit: "FaBlE", DeclaredModel: "fable", Model: "historical-fallback"}
		at := now.Add(-time.Duration(i) * time.Hour)
		if err := phaserecord.New(home, func() time.Time { return at }).RecordFallback(t.Context(), row); err != nil {
			t.Fatal(err)
		}
	}
	before := fallbackRecordSnapshot(t, home)
	for _, tc := range []struct{ name, path, phase, spec, config, out string }{
		{"same_spec_and_identity", repo, "prd", "0412-example", fallbackProjectConfig, "Start roundfix-prd-author with model opus and effort xhigh, the fallback of fable: idea of 0412-example fell back when the FaBlE limit ran out.\n"},
		{"another_spec", repo, "prd", "0413-another", fallbackProjectConfig, "Start roundfix-prd-author with no model and no effort: the declared model fable.\n"},
		{"another_repository", other, "prd", "0412-example", fallbackProjectConfig, "Start roundfix-prd-author with no model and no effort: the declared model fable.\n"},
		{"clone_identity", clone, "prd", "0412-example", fallbackProjectConfig, "Start roundfix-prd-author with model opus and effort xhigh, the fallback of fable: idea of 0412-example fell back when the FaBlE limit ran out.\n"},
		{"different_declared_model", repo, "tasks", "0412-example", fallbackProjectConfig, "Start roundfix-tasks-author with no model and no effort: the declared model opus.\n"},
		{"no_fallback", repo, "prd", "0412-example", "skills: {authoring: {claude: {prd: {model: fable}}}}\n", "Start roundfix-prd-author with no model and no effort: the declared model fable.\n"},
		{"fallback_same_family", repo, "prd", "0412-example", "skills: {authoring: {claude: {prd: {model: fable, fallback: {model: fable-backup}}}}}\n", "Start roundfix-prd-author with no model and no effort: the declared model fable.\n"},
		{"substring_model_no_effort", repo, "prd", "0412-example", "skills: {authoring: {claude: {prd: {model: claude-fable-5-2, fallback: {model: opus}}}}}\n", "Start roundfix-prd-author with model opus and no effort, the fallback of claude-fable-5-2: idea of 0412-example fell back when the FaBlE limit ran out.\n"},
		{"no_model", repo, "prd", "0412-example", "skills: {authoring: {claude: {prd: {reasoning_effort: high}}}}\n", "No Author Subagent model for prd: the Authoring Model Declaration declares no model for it.\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setCommandEnvironmentForTest(t, home, tc.path)
			writeFallbackConfig(t, tc.path, tc.config)
			phaseCall(t, fallbackArgs("model", tc.phase, tc.spec), tc.out, "", 0)
			if after := fallbackRecordSnapshot(t, home); !reflect.DeepEqual(before, after) {
				t.Fatal("model changed record")
			}
		})
	}
	t.Run("identity_and_path_rules", func(t *testing.T) {
		for _, tc := range []struct {
			name, recordID, recordPath string
			emptyCurrent, wantFallback bool
		}{
			{"old_record_path_match", "", repo, false, true},
			{"old_record_path_differs", "", other, false, false},
			{"nonempty_identity_mismatch", "sha256:other", repo, false, false},
			{"empty_current_identity", identity.Digest, "", true, true},
			{"empty_current_path_differs", identity.Digest, repo, true, false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				freshHome, freshRepo := fallbackFixture(t, fallbackProjectConfig)
				root := freshRepo
				if !tc.emptyCurrent {
					gittest.InitRepo(t, root, "-b", "main")
					gittest.Run(t, root, "commit", "--allow-empty", "-m", "isolated repository")
				}
				path := tc.recordPath
				if path == "" || path == repo && !tc.emptyCurrent {
					path = root
				}
				setCommandWorkDirForTest(t, root)
				writeFallbackConfig(t, root, fallbackProjectConfig)
				row := phaserecord.FallbackEntry{Phase: "idea", Spec: "0412-example", Repository: path, RepositoryID: tc.recordID, Limit: "Fable"}
				if err := phaserecord.New(freshHome, func() time.Time { return now }).RecordFallback(t.Context(), row); err != nil {
					t.Fatal(err)
				}
				snapshot := fallbackRecordSnapshot(t, freshHome)
				out := "Start roundfix-prd-author with no model and no effort: the declared model fable.\n"
				if tc.wantFallback {
					out = "Start roundfix-prd-author with model opus and effort xhigh, the fallback of fable: idea of 0412-example fell back when the Fable limit ran out.\n"
				}
				phaseCall(t, fallbackArgs("model", "prd", "0412-example"), out, "", 0)
				if !reflect.DeepEqual(snapshot, fallbackRecordSnapshot(t, freshHome)) {
					t.Fatal("model wrote")
				}
			})
		}
	})
}

func TestPhaseFallbackRuntimeFailuresAndBoundedOutput(t *testing.T) {
	t.Parallel()
	t.Run("record_identity_and_error_byte_limit", func(t *testing.T) {
		t.Parallel()
		home, repo := fallbackFixture(t, fallbackProjectConfig)
		gittest.InitRepo(t, repo, "-b", "main")
		gittest.Run(t, repo, "commit", "--allow-empty", "-m", "fallback repository")
		identity, err := baseline.InspectRepositoryIdentity(t.Context(), repo, nil)
		if err != nil || identity.Digest == "" {
			t.Fatalf("identity %+v %v", identity, err)
		}
		message := fallbackError + strings.Repeat("x", 4096-len(fallbackError))
		args := fallbackArgs("fallback", "idea", "0412-example")
		args[len(args)-1] = message
		line := string([]rune(message)[:200])
		phaseCall(t, args, "Start roundfix-idea-author with model opus and effort xhigh, the fallback of fable: \""+line+"\"\n", "", 0)
		rows, n, err := phaserecord.Fallbacks(home)
		if err != nil || n != 0 || len(rows) != 1 || rows[0].RepositoryID != identity.Digest || rows[0].Message != line {
			t.Fatalf("rows %+v %d %v", rows, n, err)
		}
	})
	t.Run("model_without_a_record", func(t *testing.T) {
		t.Parallel()
		home, _ := fallbackFixture(t, fallbackProjectConfig)
		phaseCall(t, fallbackArgs("model", "prd", "0412-example"), "Start roundfix-prd-author with no model and no effort: the declared model fable.\n", "", 0)
		assertNoPhaseRecord(t, home)
	})
	for _, command := range []string{"model", "fallback"} {
		t.Run(command+"_obstructed_record", func(t *testing.T) {
			t.Parallel()
			home, _ := fallbackFixture(t, fallbackProjectConfig)
			dir := filepath.Join(home, ".roundfix", "phase-record")
			if err := os.MkdirAll(filepath.Dir(dir), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(dir, []byte("occupied"), 0600); err != nil {
				t.Fatal(err)
			}
			var out, diagnostic bytes.Buffer
			code := runCLI(t, fallbackArgs(command, "prd", "0412-example"), &out, &diagnostic)
			if code != 1 || out.Len() != 0 || !strings.HasPrefix(diagnostic.String(), "roundfix: phase "+command+" failed: ") || strings.Contains(diagnostic.String(), "for usage") {
				t.Fatalf("exit %d out %q err %q", code, out.String(), diagnostic.String())
			}
			data, err := os.ReadFile(dir)
			if err != nil || string(data) != "occupied" {
				t.Fatalf("obstruction changed: %q %v", data, err)
			}
		})
	}
	t.Run("no_effort_or_session", func(t *testing.T) {
		t.Parallel()
		home, _ := fallbackFixture(t, "skills: {authoring: {claude: {prd: {model: fable, fallback: {model: opus}}}}}\n")
		phaseModelEnvironment(t, "", "")
		message := "noise\n\tYou've hit your Fable limit\x00 · " + strings.Repeat("é", 300)
		line := "You've hit your Fable limit  · " + strings.Repeat("é", 171)
		line = string([]rune(line)[:200])
		args := fallbackArgs("fallback", "prd", "0412-example")
		args[len(args)-1] = message
		phaseCall(t, args, "Start roundfix-prd-author with model opus and no effort, the fallback of fable: \""+line+"\"\n", "", 0)
		var row map[string]any
		if err := json.Unmarshal([]byte(phaseDay(t, home)[0]), &row); err != nil {
			t.Fatal(err)
		}
		if row["message"] != line || row["harness"] != "" || row["session"] != "" || row["effort"] != "" {
			t.Fatalf("row = %+v", row)
		}
	})
}
