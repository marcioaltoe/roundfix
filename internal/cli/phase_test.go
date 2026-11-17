package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
)

const phaseTestSession = "11111111-2222-4333-8444-555555555555"

func phaseFixture(t *testing.T, session string) (string, string) {
	t.Helper()
	home, repo := t.TempDir(), t.TempDir()
	gittest.MarkRepository(t, repo)
	setCommandEnvironmentForTest(t, home, repo)
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.getenv = func(key string) string {
			if key == "CLAUDE_CODE_SESSION_ID" {
				return session
			}
			return ""
		}
		deps.phaseNow = func() time.Time { return time.Date(2026, 11, 20, 13, 5, 11, 987, time.UTC) }
	})
	return home, repo
}

func phaseCall(t *testing.T, args []string, wantOut, wantErr string, wantCode int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, args, &stdout, &stderr)
	if code != wantCode || stdout.String() != wantOut || stderr.String() != wantErr {
		t.Fatalf("%v: exit = %d, want %d; stdout = %q, want %q; stderr = %q, want %q", args, code, wantCode, stdout.String(), wantOut, stderr.String(), wantErr)
	}
}

func phaseDay(t *testing.T, home string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(home, ".roundfix", "phase-record", "2026-11-20.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasSuffix(data, []byte("\n")) {
		t.Fatal("missing final newline")
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

func phaseExpectedLine(t *testing.T, phase, spec, repo, session string) string {
	t.Helper()
	quote := func(s string) string {
		data, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	harness := ""
	if session != "" {
		harness = "claude"
	}
	line := fmt.Sprintf(`{"schema":"roundfix/phase-record/v1","kind":"phase","time":"2026-11-20T13:05:11Z","phase":%s,"spec":%s,"repository":%s,"harness":%s,"session":%s,"boundary":false`, quote(phase), quote(spec), quote(repo), quote(harness), quote(session))
	return line + "}"
}

func assertNoPhaseRecord(t *testing.T, home string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(home, ".roundfix")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refusal wrote to home: %v", err)
	}
}

func TestPhaseStartSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	const help = `Usage:
  roundfix phase start <phase> [--spec <slug>]

Commands:
  start  Record in the Phase Record that this session starts a phase

Phases:
  follow-up, idea, prd, techspec, tasks, delivery, release

Options:
  --spec  Spec folder slug the phase works on; required for every phase
          but follow-up and release
`
	t.Run("transcript_1", func(t *testing.T) {
		home, _ := phaseFixture(t, phaseTestSession)
		for _, args := range [][]string{{"phase", "--help"}, {"phase", "-h"}, {"phase", "start", "--help"}, {"phase", "start", "-h"}} {
			phaseCall(t, args, help, "", 0)
		}
		assertNoPhaseRecord(t, home)
	})
	for _, number := range []int{2, 3, 4} {
		t.Run(fmt.Sprintf("transcript_%d", number), func(t *testing.T) {
			session := phaseTestSession
			if number == 4 {
				session = ""
			}
			home, repo := phaseFixture(t, session)
			phase, spec := "prd", "0412-example"
			wantErr := ""
			if number == 3 {
				phaseCall(t, []string{"phase", "start", "prd", "--spec", spec}, "Recorded prd of 0412-example for claude session "+session+"\n", "", 0)
				phase = "techspec"
			}
			args := []string{"phase", "start", phase, "--spec", spec}
			wantOut := "Recorded " + phase + " of " + spec + " for claude session " + session + "\n"
			if number == 4 {
				phase, spec = "follow-up", ""
				args = []string{"phase", "start", phase}
				wantOut = "Recorded follow-up without a harness session: the environment names none\n"
			}
			phaseCall(t, args, wantOut, wantErr, 0)
			lines := phaseDay(t, home)
			wantCount := 1
			if number == 3 {
				wantCount = 2
			}
			if len(lines) != wantCount {
				t.Fatalf("lines = %d, want %d", len(lines), wantCount)
			}
			if want := phaseExpectedLine(t, phase, spec, repo, session); lines[len(lines)-1] != want {
				t.Fatalf("line = %s, want %s", lines[len(lines)-1], want)
			}
			if number == 3 && lines[0] != phaseExpectedLine(t, "prd", spec, repo, session) {
				t.Fatalf("previous line was changed: %s", lines[0])
			}
		})
	}
	t.Run("transcript_5", func(t *testing.T) {
		home, _ := phaseFixture(t, phaseTestSession)
		phaseCall(t, []string{"phase", "start", "plan", "--spec", "0412-example"}, "", "roundfix: phase start failed: unknown phase \"plan\"; use follow-up, idea, prd, techspec, tasks, delivery or release\nRun 'roundfix phase --help' for usage.\n", 2)
		assertNoPhaseRecord(t, home)
	})
}

func TestPhaseStartRecordsEachOfTheSevenPhases(t *testing.T) {
	t.Parallel()
	for _, phase := range []string{"follow-up", "idea", "prd", "techspec", "tasks", "delivery", "release"} {
		t.Run(phase, func(t *testing.T) {
			home, repo := phaseFixture(t, phaseTestSession)
			spec := "0412-example"
			args := []string{"phase", "start", phase, "--spec", spec}
			if phase == "follow-up" || phase == "release" {
				spec = ""
				args = args[:3]
			}
			wantOut := "Recorded " + phase
			if spec != "" {
				wantOut += " of " + spec
			}
			wantOut += " for claude session " + phaseTestSession + "\n"
			phaseCall(t, args, wantOut, "", 0)
			lines := phaseDay(t, home)
			if len(lines) != 1 || lines[0] != phaseExpectedLine(t, phase, spec, repo, phaseTestSession) {
				t.Fatalf("lines = %v", lines)
			}
			decoder := json.NewDecoder(strings.NewReader(lines[0]))
			if _, err := decoder.Token(); err != nil {
				t.Fatal(err)
			}
			var keys []string
			for decoder.More() {
				key, err := decoder.Token()
				if err != nil {
					t.Fatal(err)
				}
				keys = append(keys, key.(string))
				var value json.RawMessage
				if err := decoder.Decode(&value); err != nil {
					t.Fatal(err)
				}
			}
			wantKeys := []string{"schema", "kind", "time", "phase", "spec", "repository", "harness", "session", "boundary"}
			if !reflect.DeepEqual(keys, wantKeys) {
				t.Fatalf("keys = %v", keys)
			}
		})
	}
}

func TestPhaseStartRefusalsWriteNothing(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		args   []string
		reason string
	}{
		{"missing_spec", []string{"start", "prd"}, `phase "prd" requires --spec <slug>`},
		{"malformed_spec", []string{"start", "prd", "--spec", "../0412-example"}, `invalid Spec slug "../0412-example"; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$`},
		{"extra_argument", []string{"start", "prd", "--spec", "0412-example", "extra"}, `unexpected argument "extra"`},
		{"unknown_subcommand", []string{"stop"}, `unknown subcommand "stop"; use start`},
		{"missing_subcommand", nil, `missing subcommand; use start`},
		{"missing_phase", []string{"start"}, `missing phase; use follow-up, idea, prd, techspec, tasks, delivery or release`},
		{"unknown_option", []string{"start", "release", "--unknown"}, `flag provided but not defined: -unknown`},
		{"empty_optional_spec", []string{"start", "release", "--spec="}, `invalid Spec slug ""; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$`},
		{"outside_a_repository", []string{"start", "follow-up"}, `phase start requires a git repository working tree`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home, _ := phaseFixture(t, phaseTestSession)
			if tc.name == "outside_a_repository" {
				setCommandWorkDirForTest(t, t.TempDir())
			}
			phaseCall(t, append([]string{"phase"}, tc.args...), "", "roundfix: phase start failed: "+tc.reason+"\nRun 'roundfix phase --help' for usage.\n", 2)
			assertNoPhaseRecord(t, home)
		})
	}
}

func TestPhaseStartExitsOneWhenTheRecordCannotBeWritten(t *testing.T) {
	t.Parallel()
	home, _ := phaseFixture(t, phaseTestSession)
	directory := filepath.Join(home, ".roundfix", "phase-record")
	if err := os.MkdirAll(filepath.Dir(directory), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(directory, []byte("occupied"), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"phase", "start", "follow-up"}, &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || !strings.HasPrefix(stderr.String(), "roundfix: phase start failed: ") || strings.Contains(stderr.String(), "for usage") {
		t.Fatalf("exit = %d; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(directory)
	if err != nil || string(data) != "occupied" {
		t.Fatalf("record obstruction changed: %q, %v", data, err)
	}
}

func TestPhaseStartLoadsNoConfiguration(t *testing.T) {
	t.Parallel()
	home, repo := phaseFixture(t, phaseTestSession)
	invalid := []byte("invalid: [unterminated\n")
	for _, path := range []string{filepath.Join(home, ".roundfix", "config.yml"), filepath.Join(repo, ".roundfixrc.yml")} {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, invalid, 0600); err != nil {
			t.Fatal(err)
		}
	}
	work := filepath.Join(repo, "nested")
	if err := os.Mkdir(work, 0700); err != nil {
		t.Fatal(err)
	}
	setCommandWorkDirForTest(t, work)
	phaseCall(t, []string{"phase", "start", "prd", "--spec", "0412-example"}, "Recorded prd of 0412-example for claude session "+phaseTestSession+"\n", "", 0)
	lines := phaseDay(t, home)
	if len(lines) != 1 || lines[0] != phaseExpectedLine(t, "prd", "0412-example", repo, phaseTestSession) {
		t.Fatalf("line = %v", lines)
	}
}

func TestPhaseStartPrintsNoNoticeAfterAnyPhase(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, previous, previousSpec, next, nextSpec string
		previousOut, nextOut                         string
	}{
		{"surface_transcript_1", "prd", "0412-example", "techspec", "0412-example",
			"Recorded prd of 0412-example for claude session <session>\n",
			"Recorded techspec of 0412-example for claude session <session>\n"},
		{"surface_transcript_2", "release", "0412-example", "follow-up", "",
			"Recorded release of 0412-example for claude session <session>\n",
			"Recorded follow-up for claude session <session>\n"},
		{"after_a_release_without_a_spec", "release", "", "follow-up", "",
			"Recorded release for claude session <session>\n",
			"Recorded follow-up for claude session <session>\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home, repo := phaseFixture(t, phaseTestSession)
			session := commandEnvironmentForTest(t).dependencies.getenv("CLAUDE_CODE_SESSION_ID")
			for _, entry := range []struct{ phase, spec, stdout string }{{tc.previous, tc.previousSpec, tc.previousOut}, {tc.next, tc.nextSpec, tc.nextOut}} {
				args := []string{"phase", "start", entry.phase}
				if entry.spec != "" {
					args = append(args, "--spec", entry.spec)
				}
				phaseCall(t, args, strings.ReplaceAll(entry.stdout, "<session>", session), "", 0)
			}
			lines := phaseDay(t, home)
			want := []string{phaseExpectedLine(t, tc.previous, tc.previousSpec, repo, session), phaseExpectedLine(t, tc.next, tc.nextSpec, repo, session)}
			if !reflect.DeepEqual(lines, want) {
				t.Fatalf("lines = %v, want %v", lines, want)
			}
		})
	}
}
