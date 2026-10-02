// Suite: causes CLI transcripts.
// Invariant: stdout, stderr, exits and JSON reproduce the public contract;
// all data live in a temporary repository and Roundfix Home.
package cli

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/runcause"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

const causesFixtureID = "run_20261001T120000Z_0000000000000001"

func causesCLIWorkspace(t *testing.T) commandEnvironment {
	t.Helper()
	home, repo := t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo, "--initial-branch=main")
	var err error
	home, err = filepath.EvalSymlinks(home)
	if err != nil {
		t.Fatal(err)
	}
	repo, err = filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	return commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
}

func seedCausesCLI(t *testing.T) commandEnvironment {
	t.Helper()
	env := causesCLIWorkspace(t)
	dir := filepath.Join(env.workDir, "docs/specs/0300-example")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "_tasks.md"), "---\nschema: spec-tasks/v1\nqa: task_03\ngraph:\n  nodes:\n    - {id: task_01, file: task_01.md}\n    - {id: task_02, file: task_02.md}\n    - {id: task_03, file: task_03.md}\n    - {id: task_04, file: task_04.md}\n---\n")
	mustWrite(t, filepath.Join(dir, "task_04.md"), "# Corrective Task\n\n## Overview\nPre-PR review found that the golden output was not re-recorded\n")
	s, err := store.Open(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	run, err := s.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: env.workDir, LocalBranch: "main", SpecSlug: "0300-example", ArtifactDir: filepath.Join(env.homeDir, "artifacts")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.CompleteRun(t.Context(), run.ID, store.StateClean); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	// Pin the transcript's ID and time in this isolated fixture database only.
	db, err := sql.Open("sqlite", store.DatabasePath(env.homeDir))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE runs SET id = ?, created_at = ? WHERE id = ?`, causesFixtureID, "2026-10-01T12:00:00Z", run.ID); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = store.Open(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	logDir := filepath.Join(run.ArtifactDir, "runs", causesFixtureID, "verification")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatal(err)
	}
	appendEvent := func(kind runevent.Kind, payload string) {
		t.Helper()
		if _, err := s.AppendRunEvent(t.Context(), runevent.RunEvent{RunID: causesFixtureID, Kind: kind, Source: runevent.SourceDaemon, Time: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC), Payload: []byte(payload)}); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"task_01", "task_02", "task_03", "task_04"} {
		appendEvent(runevent.KindDaemonTask, fmt.Sprintf(`{"task":%q,"phase":"started"}`, id))
	}
	for i, c := range []struct {
		task, command, diagnostic string
		attempt                   int
	}{{"task_01", "go test ./internal/store", "database is locked", 1}, {"task_02", "go test ./internal/cli", "--- FAIL: TestCausesFixture", 1}, {"task_02", "make verify", "no matching words", 2}} {
		path := filepath.Join(logDir, fmt.Sprintf("%d.log", i))
		mustWrite(t, path, c.diagnostic)
		appendEvent(runevent.KindDaemonVerification, fmt.Sprintf(`{"task":%q,"phase":"failed","attempt":%d,"command":%q,"diagnostic_path":%q}`, c.task, c.attempt, c.command, path))
		appendEvent(runevent.KindDaemonVerification, fmt.Sprintf(`{"task":%q,"phase":"verdict","attempt":%d,"verdict":"failed"}`, c.task, c.attempt))
	}
	appendEvent(runevent.KindDaemonTask, `{"task":"task_01","phase":"verification_feedback"}`)
	appendEvent(runevent.KindDaemonVerification, `{"task":"task_01","phase":"verdict","attempt":2,"verdict":"passed"}`)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	return env
}

func runCausesCLI(t *testing.T, env commandEnvironment, args ...string) (int, string, string) {
	t.Helper()
	var out, diag bytes.Buffer
	code := runWithContext(t.Context(), append([]string{"runs", "causes"}, args...), &out, &diag, env)
	return code, out.String(), diag.String()
}

func causesDigest(t *testing.T) string {
	t.Helper()
	table, err := runcause.Load()
	if err != nil {
		t.Fatal(err)
	}
	return table.SHA256[:12]
}

func TestRunsCausesPrintsAClassifiedWindow(t *testing.T) {
	t.Parallel()
	env := seedCausesCLI(t)
	code, out, diag := runCausesCLI(t, env, "--since", "2026-10-01")
	want := "verification 0300-example/task_01 " + causesFixtureID + " attempt-1 environment database-locked: go test ./internal/store\n" +
		"verification 0300-example/task_02 " + causesFixtureID + " attempt-1 implementation-defect go-test-failure: go test ./internal/cli\n" +
		"verification 0300-example/task_02 " + causesFixtureID + " attempt-2 unclassified -: make verify\n" +
		"corrective 0300-example/task_04 " + causesFixtureID + " - repository-convention record-or-golden: pre-pr-review\n" +
		"Causes: 4 item(s) from 1 terminal Spec Run(s) in window 2026-10-01..*; scope-or-authorization 0, shared-section-contract 0, repository-convention 1, implementation-defect 1, environment 1, unclassified 1; repository knowledge 1 of 3 classified; signatures " + causesDigest(t) + "\n"
	if code != 0 || diag != "" || out != want {
		t.Fatalf("exit=%d stderr=%q\nout=%q\nwant=%q", code, diag, out, want)
	}
}

func TestRunsCausesPrintsAnEmptyWindow(t *testing.T) {
	t.Parallel()
	for _, database := range []bool{false, true} {
		t.Run(fmt.Sprint(database), func(t *testing.T) {
			env := causesCLIWorkspace(t)
			if database {
				env = seedCausesCLI(t)
			}
			before := causesHomeHashes(t, env.homeDir)
			code, out, diag := runCausesCLI(t, env, "--until", "2026-09-01")
			want := "Causes: 0 item(s) from 0 terminal Spec Run(s) in window *..2026-09-01; scope-or-authorization 0, shared-section-contract 0, repository-convention 0, implementation-defect 0, environment 0, unclassified 0; repository knowledge 0 of 0 classified; signatures " + causesDigest(t) + "\n"
			if code != 0 || diag != "" || out != want {
				t.Fatalf("exit=%d out=%q stderr=%q", code, out, diag)
			}
			assertCausesHomeUnchanged(t, before, causesHomeHashes(t, env.homeDir))
		})
	}
}

func TestRunsCausesRefusesAMalformedDate(t *testing.T) {
	t.Parallel()
	env := causesCLIWorkspace(t)
	code, out, diag := runCausesCLI(t, env, "--since", "2026-13-01")
	want := "roundfix: runs causes failed: --since must be a date in YYYY-MM-DD form, got \"2026-13-01\"\nRun 'roundfix runs causes --help' for usage.\n"
	if code != 2 || out != "" || diag != want {
		t.Fatalf("exit=%d out=%q stderr=%q", code, out, diag)
	}
}

func TestRunsCausesRequiresAGitRepository(t *testing.T) {
	t.Parallel()
	env := causesCLIWorkspace(t)
	env.workDir = t.TempDir()
	code, out, diag := runCausesCLI(t, env)
	want := "roundfix: runs causes failed: runs causes requires a Git repository\nRun 'roundfix runs causes --help' for usage.\n"
	if code != 2 || out != "" || diag != want {
		t.Fatalf("exit=%d out=%q stderr=%q", code, out, diag)
	}
}

func TestRunsCausesPrintsJSON(t *testing.T) {
	t.Parallel()
	env := seedCausesCLI(t)
	code, out, diag := runCausesCLI(t, env, "--since", "2026-10-01", "--until", "2026-10-02", "--format", "json")
	if code != 0 || diag != "" {
		t.Fatalf("exit=%d stderr=%s", code, diag)
	}
	var got runcause.Report
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got.Schema != "roundfix/runs-causes/v1" || got.Runs != 1 || len(got.Items) != 4 || len(got.Tasks) != 4 || got.Summary["classified"] != 3 || len(got.SpecsNotFound) != 0 {
		t.Fatalf("JSON=%s", out)
	}
	if got.Items[2].Signature != nil || got.Items[2].Evidence != nil || got.Items[3].Attempt != nil || got.Tasks[0].VerdictsPassed != 1 || got.Tasks[0].FeedbackRounds != 1 || got.Tasks[1].VerdictsFailed != 2 || got.Tasks[3].Corrective == nil || !*got.Tasks[3].Corrective {
		t.Fatalf("JSON=%s", out)
	}
	if *got.Window.Since != "2026-10-01" || *got.Window.Until != "2026-10-02" {
		t.Fatalf("window=%+v", got.Window)
	}
	for _, secret := range []string{env.homeDir, env.workDir, "database is locked", "--- FAIL: TestCausesFixture", "diagnostic_path", "API_KEY"} {
		if strings.Contains(out, secret) {
			t.Fatalf("JSON contains private evidence %q", secret)
		}
	}
	// Require every top-level and nullable field, rather than accepting omissions.
	var document map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &document); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"schema", "window", "signatures_sha256", "runs", "items", "tasks", "summary", "specs_not_found"} {
		if _, ok := document[field]; !ok {
			t.Fatalf("missing %s", field)
		}
	}
	for _, field := range []string{`"signature":null`, `"evidence":null`, `"attempt":null`} {
		if !strings.Contains(out, field) {
			t.Fatalf("missing nullable field %s", field)
		}
	}
}

func TestRunsCausesHelp(t *testing.T) {
	t.Parallel()
	env := causesCLIWorkspace(t)
	for _, args := range [][]string{{"runs", "causes", "--help"}, {"runs", "--help"}, {"--help"}} {
		var out, diag bytes.Buffer
		code := runWithContext(t.Context(), args, &out, &diag, env)
		if code != 0 || diag.Len() != 0 || !strings.Contains(out.String(), runsCausesSynopsis) {
			t.Fatalf("args=%v exit=%d out=%s stderr=%s", args, code, &out, &diag)
		}
	}
}

func TestRunsCausesRejectsInvalidOptionsAndUnreadableDatabase(t *testing.T) {
	t.Parallel()
	env := causesCLIWorkspace(t)
	for _, args := range [][]string{{"--unknown"}, {"extra"}, {"--until", "bad"}, {"--since", ""}, {"--since", "2026-10-01", "--until", "2026-10-01"}, {"--format", "yaml"}} {
		code, out, diag := runCausesCLI(t, env, args...)
		if code != 2 || out != "" || diag == "" {
			t.Fatalf("args=%v exit=%d out=%s stderr=%s", args, code, out, diag)
		}
	}
	if err := os.MkdirAll(filepath.Dir(store.DatabasePath(env.homeDir)), 0755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, store.DatabasePath(env.homeDir), "corrupt sqlite database")
	code, out, diag := runCausesCLI(t, env)
	if code != 1 || out != "" || diag == "" {
		t.Fatalf("exit=%d out=%s stderr=%s", code, out, diag)
	}
}

func causesHomeHashes(t *testing.T, home string) map[string][32]byte {
	t.Helper()
	hashes := map[string][32]byte{}
	err := filepath.WalkDir(home, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(home, path)
		if err != nil {
			return err
		}
		hashes[rel] = sha256.Sum256(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hashes
}

func assertCausesHomeUnchanged(t *testing.T, before, after map[string][32]byte) {
	t.Helper()
	for path, hash := range before {
		if got, ok := after[path]; !ok || got != hash {
			t.Fatalf("home file changed: %s", path)
		}
	}
	for path := range after {
		if _, ok := before[path]; !ok && path != ".roundfix/roundfix.db-wal" && path != ".roundfix/roundfix.db-shm" {
			t.Fatalf("unexpected new home file: %s", path)
		}
	}
}

func TestBuildLeavesTheRunDatabaseUnchanged(t *testing.T) {
	t.Parallel()
	env := seedCausesCLI(t)
	before := causesHomeHashes(t, env.homeDir)
	reader, err := store.OpenReader(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	table, err := runcause.Load()
	if err != nil {
		t.Fatal(err)
	}
	report, err := runcause.Build(t.Context(), table, runcause.Request{Reader: reader, RepositoryRoot: env.workDir, SpecsRoot: filepath.Join(env.workDir, "docs/specs"), ArchiveRoot: filepath.Join(env.workDir, "docs/history/specs")})
	if err != nil {
		t.Fatal(err)
	}
	assertCausesHomeUnchanged(t, before, causesHomeHashes(t, env.homeDir))
	code, out, diag := runCausesCLI(t, env, "--format", "json")
	if code != 0 || diag != "" {
		t.Fatalf("exit=%d stderr=%s", code, diag)
	}
	var fromCommand runcause.Report
	if err := json.Unmarshal([]byte(out), &fromCommand); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, fromCommand) {
		t.Fatal("Build and command reports disagree")
	}
	assertCausesHomeUnchanged(t, before, causesHomeHashes(t, env.homeDir))
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	assertCausesHomeUnchanged(t, before, causesHomeHashes(t, env.homeDir))
}
