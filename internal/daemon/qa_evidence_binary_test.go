// Suite: QA Report evidence binaries
// Invariant: image evidence commits and other binary evidence stays uncommitted.
// Boundary IN: TaskCycle with the existing fake QA Agent and real temporary Git repository
// Boundary OUT: ACP runtime and push scanning
package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func qaEvidenceBinaryFixture(t *testing.T, files map[string]string) (*taskCycleFixture, TaskCycleResult, string) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	root := filepath.ToSlash(filepath.Join("docs", "specs", plan.Spec.Slug, "qa", "evidence"))
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass),
		afterQA: func() {
			for name, content := range files {
				path := filepath.Join(plan.WorkDir, root, name)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				mustWriteForTest(t, path, content)
			}
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	return fixture, result, root
}

func TestQAReportCommitsImageEvidence(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"screen.png":  "\x89PNG\r\n\x1a\n\x00image",
		"screen.jpg":  "\xff\xd8\xff\xe0\x00image",
		"diagram.svg": "<svg xmlns=\"http://www.w3.org/2000/svg\"></svg>\n",
		"capture":     "\x89PNG\r\n\x1a\n\x00image",
	}
	fixture, result, root := qaEvidenceBinaryFixture(t, files)
	if !result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusCompleted) {
		t.Fatal("QA with image evidence was not accepted")
	}
	message := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "log", "-1", "--format=%B"))
	if message != QACommitMessage(taskCycleSlug, spec.VerdictPass) {
		t.Fatalf("QA commit message = %q", message)
	}
	for name, content := range files {
		if got := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+root+"/"+name); got != content {
			t.Fatalf("committed %s = %q, want original bytes", name, got)
		}
	}
}

func TestQAReportWithANonImageBinarySettlesFailedWithoutACommit(t *testing.T) {
	t.Parallel()
	for _, names := range [][]string{{"program.png"}, {"z-state.db", "a-program.png"}} {
		t.Run(strings.Join(names, "+"), func(t *testing.T) {
			t.Parallel()
			files := make(map[string]string)
			for _, name := range names {
				files[name] = "\x7fELF\x00private fixture bytes"
			}
			fixture, result, root := qaEvidenceBinaryFixture(t, files)
			if result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusFailed) {
				t.Fatal("QA with a non-image binary must settle failed")
			}
			var paths []string
			for _, name := range names {
				paths = append(paths, root+"/"+name)
			}
			sort.Strings(paths)
			wantReason := "QA evidence carries a non-image binary: " + strings.Join(paths, "; ")
			var reason string
			for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
				payload := eventPayloadMap(t, event)
				if event.ReviewIssue == fixture.graph.QATaskID && payload["phase"] == "settled" {
					reason, _ = payload["reason"].(string)
				}
			}
			if reason != wantReason {
				t.Fatalf("QA reason = %q, want %q", reason, wantReason)
			}
			message := fmt.Sprintf("QA Report refused: %d non-image binary file(s)", len(paths))
			wantPaths := make([]any, len(paths))
			for index, path := range paths {
				wantPaths[index] = path
			}
			refusals := 0
			for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
				payload := eventPayloadMap(t, event)
				if payload["phase"] == "evidence_binary" {
					refusals++
					want := map[string]any{"phase": "evidence_binary", "paths": wantPaths}
					if event.Summary != message || !reflect.DeepEqual(payload, want) {
						t.Fatalf("refusal event = %s %v, want %s %v", event.Summary, payload, message, want)
					}
				}
			}
			if refusals != 1 || !strings.Contains(fixture.progress.String(), message+"\n") {
				t.Fatalf("refusals = %d, progress = %s", refusals, fixture.progress.String())
			}
			if history := runGitForTest(t, fixture.gitRoot, "log", "--format=%s"); strings.Contains(history, "docs: qa report for "+taskCycleSlug) {
				t.Fatal("refused QA Report was committed")
			}
			for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit) {
				if strings.HasPrefix(event.Summary, "QA Report commit created:") {
					t.Fatal("refused QA Report published a commit event")
				}
			}
			for name, content := range files {
				got, err := os.ReadFile(filepath.Join(fixture.gitRoot, root, name))
				if err != nil || string(got) != content {
					t.Fatalf("refusal lost or altered %s: %v", name, err)
				}
			}
			if _, err := os.Stat(filepath.Join(fixture.gitRoot, result.QAReportPath)); err != nil {
				t.Fatalf("refusal lost the report: %v", err)
			}
			outputs := fixture.progress.String()
			for _, event := range fixture.sink.events {
				outputs += event.Summary + string(event.Payload)
			}
			if strings.Contains(outputs, "private fixture bytes") {
				t.Fatal("refusal exposed binary content")
			}
		})
	}
}

func TestQAReportBinaryCheckReadsOnlyPreparedEvidence(t *testing.T) {
	t.Parallel()
	workDir := t.TempDir()
	report := "custom/qa/report.md"
	root := "custom/qa/evidence/"
	files := map[string]string{
		root + "nested/image":               "\xff\xd8\xff\x00image",
		root + "a-last-byte.bin":            strings.Repeat("a", speccheck.EvidenceProbeBytes-1) + "\x00",
		root + "beyond-probe.txt":           strings.Repeat("a", speccheck.EvidenceProbeBytes) + "\x00",
		root + "z-program.png":              "\x7fELF\x00fixture",
		root + "unprepared.bin":             "\x00fixture",
		"custom/qa/evidence-other/state.db": "SQLite format 3\x00fixture",
		"other/qa/evidence/state.db":        "SQLite format 3\x00fixture",
	}
	for path, content := range files {
		fullPath := filepath.Join(workDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		mustWriteForTest(t, fullPath, content)
	}
	paths := []string{
		root + "z-program.png", root + "deleted.bin", root + "nested/image",
		root + "beyond-probe.txt", root + "a-last-byte.bin",
		"custom/qa/evidence-other/state.db", "other/qa/evidence/state.db",
	}
	got, err := nonImageQAReportEvidence(t.Context(), workDir, report, paths)
	want := []string{root + "a-last-byte.bin", root + "z-program.png"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("binary paths = %v, %v; want %v", got, err, want)
	}
	// Existing directories are infrastructure errors rather than silently skipped evidence.
	if _, err := nonImageQAReportEvidence(t.Context(), workDir, report, []string{root + "nested"}); err == nil {
		t.Fatal("directory read error was accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := nonImageQAReportEvidence(ctx, workDir, report, paths); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled scan = %v, want context.Canceled", err)
	}
}
