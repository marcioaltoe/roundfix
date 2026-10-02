// Suite: prior QA pass compiled-source exclusion
// Invariant: compiled evidence is never imported or compared and every skip is journaled.
// Boundary IN: real Git, Run journal, prior-pass import and worktree files
// Boundary OUT: ACP runtime and evidence carry proof (existing suites)
package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func recordPriorSourceEvidence(t *testing.T, f *taskCycleFixture, p TaskPlan, names []string) priorQAPass {
	t.Helper()
	base := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), false)
	runGitForTest(t, f.gitRoot, "checkout", "--detach", pass.Commit)
	for _, name := range names {
		path := filepath.ToSlash(filepath.Join(filepath.Dir(pass.Report), "evidence", name))
		mustWriteForTest(t, filepath.Join(f.gitRoot, path), "evidence bytes: "+name+"\r\n")
		runGitForTest(t, f.gitRoot, "add", path)
		pass.Files = append(pass.Files, path)
	}
	runGitForTest(t, f.gitRoot, "commit", "--amend", "--no-edit")
	pass.Commit = strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	runGitForTest(t, f.gitRoot, "branch", "-f", store.RunBranchPrefix+f.run.ID, pass.Commit)
	payload, err := json.Marshal(map[string]any{"decision": "created", "task": "task_02", "commit": pass.Commit, "report": pass.Report})
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.store.AppendRunEvent(context.Background(), runevent.RunEvent{RunID: f.run.ID, Batch: 2, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonCommit, Time: taskCycleNowForTest(), Payload: payload})
	if err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, f.gitRoot, "checkout", "--detach", base)
	return pass
}

func assertPriorImportedFilesAndSkips(t *testing.T, f *taskCycleFixture, wantFiles, wantSkipped []string) {
	t.Helper()
	events := taskEventsOfKind(f.sink, runevent.KindDaemonQA)
	if len(events) != 1 {
		t.Fatalf("prior events = %d, want 1", len(events))
	}
	var payload struct {
		Files   []string `json:"files"`
		Skipped []struct {
			Path   string `json:"path"`
			Reason string `json:"reason"`
		} `json:"skipped"`
	}
	if err := json.Unmarshal(events[0].Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(payload.Files, wantFiles) {
		t.Fatalf("imported files = %v, want %v", payload.Files, wantFiles)
	}
	if payload.Skipped == nil || len(payload.Skipped) != len(wantSkipped) {
		t.Fatalf("skipped = %+v, want %v", payload.Skipped, wantSkipped)
	}
	for i, path := range wantSkipped {
		if payload.Skipped[i].Path != path || payload.Skipped[i].Reason != "compiled source" {
			t.Fatalf("skip %d = %+v, want %s: compiled source", i, payload.Skipped[i], path)
		}
	}
	assertPriorOutcome(t, f, "imported", "")
}

func TestPriorQAPassSkipsCompiledSourceEvidence(t *testing.T) {
	f, p, e := priorFixture(t)
	prior := recordPriorSourceEvidence(t, f, p, []string{"replay_test.go", "ledger.txt", "replay_test.go.txt"})
	pass, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || !imported || pass.Commit != prior.Commit {
		t.Fatalf("import = %+v %v %v", pass, imported, err)
	}
	source := filepath.ToSlash(filepath.Join(filepath.Dir(prior.Report), "evidence", "replay_test.go"))
	if _, err := os.Stat(filepath.Join(f.gitRoot, source)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("compiled source imported: %v", err)
	}
	wantFiles := []string{prior.Files[0], prior.Files[3], prior.Files[4], prior.Report}
	sort.Strings(wantFiles)
	for _, path := range wantFiles {
		got, err := os.ReadFile(filepath.Join(f.gitRoot, path))
		if err != nil {
			t.Fatal(err)
		}
		want, err := priorQAGit(context.Background(), f.gitRoot, "cat-file", "blob", prior.Commit+":"+path)
		if err != nil || string(got) != string(want) {
			t.Fatalf("imported bytes differ for %s: %v", path, err)
		}
	}
	if !reflect.DeepEqual(pass.Files, wantFiles) {
		t.Fatalf("returned files = %v, want %v", pass.Files, wantFiles)
	}
	assertPriorImportedFilesAndSkips(t, f, wantFiles, []string{source})
}

func TestPriorQAPassSkipsEveryCompiledExtension(t *testing.T) {
	for _, extension := range []string{".go", ".rs", ".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs"} {
		t.Run(extension, func(t *testing.T) {
			f, p, e := priorFixture(t)
			prior := recordPriorSourceEvidence(t, f, p, []string{"replay" + extension})
			_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
			if err != nil || !imported {
				t.Fatalf("import = %v %v", imported, err)
			}
			source := prior.Files[2]
			if _, err := os.Stat(filepath.Join(f.gitRoot, source)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("compiled source imported: %v", err)
			}
			assertPriorImportedFilesAndSkips(t, f, []string{prior.Files[0], prior.Report}, []string{source})
		})
	}
}

func TestPriorQAPassSkipNeverReportsPathDiffers(t *testing.T) {
	for _, extension := range []string{".go", ".rs", ".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs"} {
		t.Run(extension, func(t *testing.T) {
			f, p, e := priorFixture(t)
			prior := recordPriorSourceEvidence(t, f, p, []string{"replay" + extension})
			source := prior.Files[2]
			full := filepath.Join(f.gitRoot, source)
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, full, "different local bytes")
			_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
			if err != nil || !imported {
				t.Fatalf("import = %v %v", imported, err)
			}
			got, err := os.ReadFile(full)
			if err != nil || string(got) != "different local bytes" {
				t.Fatalf("local source changed: %q %v", got, err)
			}
			assertPriorImportedFilesAndSkips(t, f, []string{prior.Files[0], prior.Report}, []string{source})
		})
	}
}

func TestPriorQAPassWithoutRecordedPassPreservesNullFiles(t *testing.T) {
	f, p, e := priorFixture(t)
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("import without prior pass = %v %v", imported, err)
	}
	events := taskEventsOfKind(f.sink, runevent.KindDaemonQA)
	if len(events) != 1 {
		t.Fatalf("prior events = %d, want 1", len(events))
	}
	payload := eventPayloadMap(t, events[0])
	files, present := payload["files"]
	if !present || files != nil {
		t.Fatalf("files without an import = %v (present: %v), want null", files, present)
	}
	skipped, ok := payload["skipped"].([]any)
	if !ok || len(skipped) != 0 {
		t.Fatalf("skipped without a prior pass = %v, want empty list", payload["skipped"])
	}
	assertPriorOutcome(t, f, "none", "")
}
