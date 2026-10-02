// Suite: recorded causes.
// Invariant: terminal-window membership, attempt deduplication and corrective
// graph membership determine counts without changing the temporary home.
package runcause

import (
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

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func causeFixture(t *testing.T) (Table, Request, string, store.Run) {
	t.Helper()
	home, repo := t.TempDir(), t.TempDir()
	var err error
	home, err = filepath.EvalSymlinks(home)
	if err != nil {
		t.Fatal(err)
	}
	table, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	req := Request{RepositoryRoot: repo, SpecsRoot: filepath.Join(repo, "specs"), ArchiveRoot: filepath.Join(repo, "archive")}
	dir := filepath.Join(req.SpecsRoot, "0300-example")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	writeCauseFile(t, filepath.Join(dir, "_tasks.md"), "---\nschema: spec-tasks/v1\nqa: task_03\ngraph:\n  nodes:\n    - {id: task_01, file: task_01.md}\n    - {id: task_02, file: task_02.md}\n    - {id: task_03, file: task_03.md}\n    - {id: task_04, file: task_04.md}\n---\n")
	writeCauseFile(t, filepath.Join(dir, "task_04.md"), "# Repair golden output\n\n## Overview\nPre-PR review found that the golden output was not re-recorded\n\n## Result\nauthorization must not classify this Task\n")
	s, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	run, err := s.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: "main", SpecSlug: "0300-example", ArtifactDir: filepath.Join(home, "artifacts")})
	if err != nil {
		t.Fatal(err)
	}
	logDir := filepath.Join(run.ArtifactDir, "runs", run.ID, "verification")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatal(err)
	}
	for i, diagnostic := range []string{"database is locked", "--- FAIL: TestCausesFixture", "unmatched diagnostic"} {
		writeCauseFile(t, filepath.Join(logDir, fmt.Sprintf("%d.log", i)), diagnostic)
	}
	events := []struct {
		kind    runevent.Kind
		payload string
	}{
		{runevent.KindDaemonTask, `{"task":"task_01","phase":"started"}`},
		{runevent.KindDaemonTask, `{"task":"task_02","phase":"started"}`},
		{runevent.KindDaemonTask, `{"task":"task_03","phase":"started"}`},
		{runevent.KindDaemonTask, `{"task":"task_04","phase":"started"}`},
		{runevent.KindDaemonVerification, fmt.Sprintf(`{"task":"task_01","attempt":1,"phase":"failed","command":"go test ./internal/store","diagnostic_path":%q}`, filepath.Join(logDir, "0.log"))},
		{runevent.KindDaemonVerification, `{"task":"task_01","attempt":1,"phase":"verdict","verdict":"failed"}`},
		{runevent.KindDaemonVerification, fmt.Sprintf(`{"task":"task_02","attempt":1,"phase":"failed","command":"go test ./internal/cli","diagnostic_path":%q}`, filepath.Join(logDir, "1.log"))},
		{runevent.KindDaemonVerification, `{"task":"task_02","attempt":1,"phase":"verdict","verdict":"failed"}`},
		{runevent.KindDaemonTask, `{"task":"task_01","phase":"verification_feedback"}`},
		{runevent.KindDaemonVerification, `{"task":"task_01","attempt":2,"phase":"verdict","verdict":"passed"}`},
		{runevent.KindDaemonVerification, fmt.Sprintf(`{"task":"task_02","attempt":2,"phase":"failed","commands":["make verify"],"diagnostic_path":%q}`, filepath.Join(logDir, "2.log"))},
		{runevent.KindDaemonVerification, `{"task":"task_02","attempt":2,"phase":"verdict","verdict":"failed"}`},
	}
	for _, e := range events {
		if _, err := s.AppendRunEvent(t.Context(), runevent.RunEvent{RunID: run.ID, Kind: e.kind, Source: runevent.SourceDaemon, Time: time.Now().UTC(), Payload: json.RawMessage(e.payload)}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.CompleteRun(t.Context(), run.ID, store.StateClean); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	reader, err := store.OpenReader(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := reader.Close(); err != nil {
			t.Error(err)
		}
	})
	req.Reader = reader
	return table, req, home, run
}

func writeCauseFile(t *testing.T, path, text string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestBuildReportsAttemptsTasksAndCorrectives(t *testing.T) {
	t.Parallel()
	table, req, _, run := causeFixture(t)
	report, err := Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Runs != 1 || len(report.Items) != 4 || len(report.Tasks) != 4 {
		t.Fatalf("report=%+v", report)
	}
	for i, want := range []string{"environment", "implementation-defect", "unclassified", "repository-convention"} {
		if report.Items[i].Class != want || report.Items[i].RunID != run.ID {
			t.Fatalf("item=%+v", report.Items[i])
		}
	}
	if report.Items[3].Check != "pre-pr-review" || report.Items[3].Attempt != nil || report.Items[2].Signature != nil || report.Items[2].Evidence != nil {
		t.Fatalf("nullable/trigger fields=%+v", report.Items)
	}
	want := map[string]int{"scope-or-authorization": 0, "shared-section-contract": 0, "repository-convention": 1, "implementation-defect": 1, "environment": 1, "unclassified": 1, "items": 4, "classified": 3, "repository_knowledge": 1}
	if !reflect.DeepEqual(report.Summary, want) {
		t.Fatalf("summary=%v", report.Summary)
	}
	for i, c := range report.Tasks {
		if c.Corrective == nil || *c.Corrective != (i == 3) || c.QA != (i == 2) || c.Runs != 1 {
			t.Fatalf("task=%+v", c)
		}
	}
	if c := report.Tasks[0]; c.VerdictsPassed != 1 || c.VerdictsFailed != 1 || c.FeedbackRounds != 1 {
		t.Fatalf("counts=%+v", c)
	}
	if report.Tasks[1].VerdictsFailed != 2 {
		t.Fatalf("counts=%+v", report.Tasks[1])
	}
	// The archive is consulted without an active PRD or task execution eligibility.
	if err := os.Rename(req.SpecsRoot, req.ArchiveRoot); err != nil {
		t.Fatal(err)
	}
	archived, err := Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, archived) {
		t.Fatal("archive changed report")
	}
	if err := os.Remove(filepath.Join(req.ArchiveRoot, "0300-example", "_tasks.md")); err != nil {
		t.Fatal(err)
	}
	missing, err := Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(missing.SpecsNotFound, []string{"0300-example"}) || len(missing.Items) != 3 {
		t.Fatalf("missing=%+v", missing)
	}
	for _, c := range missing.Tasks {
		if c.Corrective != nil {
			t.Fatal("missing graph implied corrective false")
		}
	}
}

func TestBuildIgnoresALogOutsideTheArtifactDirectory(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"leaf symlink", "parent symlink", "outside", "other run", "directory", "missing"} {
		t.Run(mode, func(t *testing.T) {
			table, req, home, run := causeFixture(t)
			path := filepath.Join(run.ArtifactDir, "runs", run.ID, "verification", "0.log")
			outside := filepath.Join(home, "outside.log")
			writeCauseFile(t, outside, "database is locked")
			switch mode {
			case "leaf symlink":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, path); err != nil {
					t.Fatal(err)
				}
			case "parent symlink":
				dir := filepath.Dir(path)
				if err := os.Rename(dir, dir+"-real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(dir+"-real", dir); err != nil {
					t.Fatal(err)
				}
			case "outside":
				path = outside
			case "other run":
				path = filepath.Join(run.ArtifactDir, "runs", "another", "0.log")
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				writeCauseFile(t, path, "database is locked")
			case "directory":
				path = filepath.Dir(path)
			case "missing":
				path = filepath.Join(filepath.Dir(path), "absent.log")
			}
			if mode != "leaf symlink" && mode != "parent symlink" {
				db, err := sql.Open("sqlite", store.DatabasePath(home))
				if err != nil {
					t.Fatal(err)
				}
				payload := fmt.Sprintf(`{"task":"task_01","attempt":1,"phase":"failed","command":"go test ./internal/store","diagnostic_path":%q}`, path)
				if _, err := db.Exec(`UPDATE run_events SET payload = ? WHERE run_id = ? AND cursor = 5`, payload, run.ID); err != nil {
					t.Fatal(err)
				}
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}
			report, err := Build(t.Context(), table, req)
			if err != nil {
				t.Fatal(err)
			}
			if report.Items[0].Class != "unclassified" || report.Items[0].Evidence != nil {
				t.Fatalf("unsafe log used: %+v", report.Items[0])
			}
		})
	}
}

func causeHashes(t *testing.T, home string) map[string][32]byte {
	t.Helper()
	hashes := map[string][32]byte{}
	err := filepath.WalkDir(home, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hashes[path] = sha256.Sum256(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hashes
}

func TestBuildLeavesTheRunDatabaseUnchanged(t *testing.T) {
	t.Parallel()
	table, req, home, _ := causeFixture(t)
	before := causeHashes(t, home)
	if _, err := Build(t.Context(), table, req); err != nil {
		t.Fatal(err)
	}
	if after := causeHashes(t, home); !reflect.DeepEqual(before, after) {
		t.Fatal("Build changed home files")
	}
}

func TestBuildWindowAndMultipleRuns(t *testing.T) {
	t.Parallel()
	table, req, home, run := causeFixture(t)
	s, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	next, err := s.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: req.RepositoryRoot, LocalBranch: "main", SpecSlug: run.SpecSlug, ArtifactDir: run.ArtifactDir})
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"task_01", "task_04", "task_99"} {
		payload := fmt.Sprintf(`{"task":%q,"phase":"started"}`, id)
		if _, err := s.AppendRunEvent(t.Context(), runevent.RunEvent{RunID: next.ID, Kind: runevent.KindDaemonTask, Time: time.Now().UTC(), Payload: []byte(payload)}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.CompleteRun(t.Context(), next.ID, store.StateClean); err != nil {
		t.Fatal(err)
	}
	// Active Runs and other repositories do not enter the window.
	if _, err := s.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: req.RepositoryRoot, LocalBranch: "main", SpecSlug: "active"}); err != nil {
		t.Fatal(err)
	}
	other, err := s.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: t.TempDir(), LocalBranch: "main", SpecSlug: "other"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.CompleteRun(t.Context(), other.ID, store.StateClean); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	report, err := Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Runs != 2 || len(report.Items) != 4 || report.Tasks[0].Runs != 2 || report.Tasks[3].Runs != 2 {
		t.Fatalf("multi-run=%+v", report)
	}
	if report.Tasks[4].Corrective == nil || *report.Tasks[4].Corrective {
		t.Fatal("nonmember became corrective")
	}
	req.Since = run.CreatedAt
	req.Until = next.CreatedAt
	report, err = Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Runs != 1 {
		t.Fatalf("half-open window: %d Runs", report.Runs)
	}
	req.Since = next.CreatedAt
	req.Until = time.Time{}
	report, err = Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Runs != 1 || len(report.Items) != 1 || report.Items[0].Kind != "corrective" {
		t.Fatalf("inclusive window=%+v", report)
	}
}

func TestBuildBoundsEvidenceAndExportsNoPrivateCommandText(t *testing.T) {
	t.Parallel()
	table, req, home, run := causeFixture(t)
	log := filepath.Join(run.ArtifactDir, "runs", run.ID, "verification", "0.log")
	// An early signature beyond the tail cannot override the retained failure.
	writeCauseFile(t, log, "database is locked\n"+strings.Repeat("x", int(table.DiagnosticTailBytes))+"\n--- FAIL: TestTail")
	// Only the bounded title/Overview is classified, never later sections.
	writeCauseFile(t, filepath.Join(req.SpecsRoot, run.SpecSlug, "task_04.md"), "# Follow-up\n## Overview\n"+strings.Repeat("x", table.TaskTextBytes)+" golden\n## Result\nauthorization\n")
	db, err := sql.Open("sqlite", store.DatabasePath(home))
	if err != nil {
		t.Fatal(err)
	}
	command := `API_KEY="credential-sentinel" go test /private/repository/package --token secret-sentinel`
	payload := fmt.Sprintf(`{"task":"task_01","attempt":1,"phase":"failed","command":%q,"diagnostic_path":%q}`, command, log)
	if _, err := db.Exec(`UPDATE run_events SET payload = ? WHERE run_id = ? AND cursor = 5`, payload, run.ID); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	report, err := Build(t.Context(), table, req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Items[0].Class != "implementation-defect" || report.Items[3].Class != "unclassified" {
		t.Fatalf("bounds ignored: %+v", report.Items)
	}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, private := range []string{"credential-sentinel", "secret-sentinel", "/private/repository", home, "database is locked", "TestTail"} {
		if strings.Contains(string(data), private) {
			t.Fatalf("private text exported: %q", private)
		}
	}
}
