package daemon

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func qaFormatterScript(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "formatter.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body), 0o600); err != nil {
		t.Fatal(err)
	}
	return "sh '" + strings.ReplaceAll(path, "'", "'\\''") + "'"
}

func qaFormatCommitFixture(t *testing.T, script string) (*taskCycleFixture, TaskPlan, *Engine, string, string) {
	t.Helper()
	f, p, e := priorFixture(t)
	before, err := e.snapshotQAPaths(context.Background(), p.WorkDir)
	if err != nil {
		t.Fatal(err)
	}
	report := filepath.ToSlash(filepath.Join("docs", "specs", p.Spec.Slug, "qa", priorReportName(1)))
	evidence := filepath.ToSlash(filepath.Join("docs", "specs", p.Spec.Slug, "qa", "evidence", "a file;$.txt"))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(p.WorkDir, evidence)), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(p.WorkDir, report), priorReportContent("fail"))
	mustWriteForTest(t, filepath.Join(p.WorkDir, evidence), "original evidence\n")
	p.FormatCommand = script
	if err := e.commitQAReport(context.Background(), p, 1, before, nil, "fail", report, p.Tasks[len(p.Tasks)-1]); err != nil {
		t.Fatal(err)
	}
	return f, p, e, report, evidence
}

func assertQAFormatOutcome(t *testing.T, f *taskCycleFixture, stage, outcome string, count int) map[string]any {
	t.Helper()
	var found []map[string]any
	for _, event := range taskEventsOfKind(f.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == "format" {
			found = append(found, payload)
		}
	}
	if len(found) != count {
		t.Fatalf("format events = %v, want %d", found, count)
	}
	if count == 0 {
		return nil
	}
	payload := found[count-1]
	if payload["stage"] != stage || payload["outcome"] != outcome {
		t.Fatalf("format event = %v", payload)
	}
	if outcome == "failed" || outcome == "reverted" {
		if payload["reason"] == "" || !strings.Contains(f.progress.String(), "roundfix: QA format "+stage+" "+outcome+":") {
			t.Fatalf("missing diagnostic: %v; %s", payload, f.progress)
		}
	}
	return payload
}

func qaCommittedBytes(t *testing.T, p TaskPlan, path string) string {
	t.Helper()
	content, err := priorQAGit(context.Background(), p.WorkDir, "show", "HEAD:"+path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestQAFormatCommitsTheFormattedReport(t *testing.T) {
	script := qaFormatterScript(t, "for path do printf '\\nformatted\\n' >> \"$path\"; done\n")
	f, p, _, report, evidence := qaFormatCommitFixture(t, script)
	if got := qaCommittedBytes(t, p, report); got != priorReportContent("fail")+"\nformatted\n" {
		t.Fatalf("committed report = %q", got)
	}
	if got := qaCommittedBytes(t, p, evidence); got != "original evidence\n\nformatted\n" {
		t.Fatalf("committed evidence = %q", got)
	}
	payload := assertQAFormatOutcome(t, f, "commit", "formatted", 1)
	if payload["paths"] != float64(2) {
		t.Fatalf("paths = %v", payload)
	}
	want := []any{evidence, report}
	if !reflect.DeepEqual(payload["changed"], want) {
		t.Fatalf("changed = %v, want %v", payload["changed"], want)
	}
}

func TestQAFormatFormatsAnImportedPassBeforeThePrecondition(t *testing.T) {
	f, p, _ := priorFixture(t)
	prior := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("fail"), QACommitMessage(p.Spec.Slug, "fail"), true)
	p.FormatCommand = qaFormatterScript(t, "for path do printf '\\nformatted\\n' >> \"$path\"; done\n")
	p.RepositoryVerification = `find docs/specs/` + p.Spec.Slug + `/qa -type f -exec sh -c 'for path do grep -aq formatted "$path" || exit 1; done' check {} +`
	runner := &taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot, qaReport: qaReportForTest(spec.VerdictPass)}
	e := f.engine(t, runner, ExecVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})
	if _, err := e.TaskCycle(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("gate never reached Agent: %s", f.progress)
	}
	for _, path := range prior.Files {
		if got := qaCommittedBytes(t, p, path); !strings.Contains(got, "formatted") {
			t.Fatalf("import not formatted: %s", path)
		}
	}
	events := taskEventsOfKind(f.sink, runevent.KindDaemonQA)
	imported := false
	for _, event := range events {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == "format" && payload["stage"] == "imported" && payload["outcome"] == "formatted" {
			imported = true
		}
	}
	if !imported {
		t.Fatal("missing imported format event")
	}
	assertQAFormatOutcome(t, f, "commit", "formatted", 2)
}

func TestQAFormatRevertsAFailingFormatter(t *testing.T) {
	script := qaFormatterScript(t, "for path do printf 'broken' > \"$path\"; chmod 700 \"$path\"; done\nprintf 'formatter error' >&2\nexit 1\n")
	f, p, _, report, evidence := qaFormatCommitFixture(t, script)
	if got := qaCommittedBytes(t, p, report); got != priorReportContent("fail") {
		t.Fatalf("report = %q", got)
	}
	if got := qaCommittedBytes(t, p, evidence); got != "original evidence\n" {
		t.Fatalf("evidence = %q", got)
	}
	for _, path := range []string{report, evidence} {
		info, err := os.Stat(filepath.Join(p.WorkDir, path))
		if err != nil || info.Mode().Perm() != 0o644 {
			t.Fatalf("mode of %s = %v %v", path, info, err)
		}
	}
	payload := assertQAFormatOutcome(t, f, "commit", "failed", 1)
	if payload["diagnostics"] != "formatter error" {
		t.Fatalf("diagnostics = %v", payload)
	}
}

func TestQAFormatRevertsAVerdictChange(t *testing.T) {
	script := qaFormatterScript(t, "for path do sed 's/verdict: fail/verdict: pass/' \"$path\" > \"$path.tmp\"; mv \"$path.tmp\" \"$path\"; printf '\\nchanged\\n' >> \"$path\"; done\n")
	f, p, _, report, evidence := qaFormatCommitFixture(t, script)
	if got := qaCommittedBytes(t, p, report); got != priorReportContent("fail") {
		t.Fatalf("report = %q", got)
	}
	if got := qaCommittedBytes(t, p, evidence); got != "original evidence\n" {
		t.Fatalf("evidence = %q", got)
	}
	assertQAFormatOutcome(t, f, "commit", "reverted", 1)
}

func TestQAFormatLeavesFilesOutsideTheQADirectory(t *testing.T) {
	log := filepath.Join(t.TempDir(), "arguments")
	script := qaFormatterScript(t, "printf '%s\\n' \"$@\" > '"+log+"'\n")
	f, p, e := priorFixture(t)
	before, err := e.snapshotQAPaths(context.Background(), p.WorkDir)
	if err != nil {
		t.Fatal(err)
	}
	report := filepath.ToSlash(filepath.Join("docs", "specs", p.Spec.Slug, "qa", priorReportName(1)))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(p.WorkDir, report)), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(p.WorkDir, report), priorReportContent("fail"))
	mustWriteForTest(t, filepath.Join(p.WorkDir, "outside.txt"), "outside\n")
	// Symlink leaves and symlink ancestors must not supply outside files.
	if err := os.Symlink(filepath.Join(p.WorkDir, "outside.txt"), filepath.Join(p.Spec.Dir, "qa", "link")); err != nil {
		t.Fatal(err)
	}
	p.FormatCommand = script
	if err := e.commitQAReport(context.Background(), p, 1, before, nil, "fail", report, p.Tasks[len(p.Tasks)-1]); err != nil {
		t.Fatal(err)
	}
	arguments, err := os.ReadFile(log)
	if err != nil || string(arguments) != report+"\n" {
		t.Fatalf("arguments = %q %v", arguments, err)
	}
	if got := qaCommittedBytes(t, p, "outside.txt"); got != "outside\n" {
		t.Fatalf("outside = %q", got)
	}
	assertQAFormatOutcome(t, f, "commit", "unchanged", 1)
}

func TestQAFormatEmptyCommandRunsNothing(t *testing.T) {
	f, p, _, report, evidence := qaFormatCommitFixture(t, "")
	if qaCommittedBytes(t, p, report) != priorReportContent("fail") || qaCommittedBytes(t, p, evidence) != "original evidence\n" {
		t.Fatal("empty formatter changed bytes")
	}
	assertQAFormatOutcome(t, f, "", "", 0)
	if strings.Contains(f.progress.String(), "roundfix: QA format") {
		t.Fatal("empty formatter wrote diagnostic")
	}
}

func TestQAFormatCancellationRestoresFiles(t *testing.T) {
	for _, name := range []string{"context", "stop request"} {
		t.Run(name, func(t *testing.T) {
			f, p, e := priorFixture(t)
			directory := t.TempDir()
			ready := filepath.Join(directory, "ready")
			report := filepath.Join(p.Spec.Dir, "qa", priorReportName(1))
			if err := os.MkdirAll(filepath.Dir(report), 0o755); err != nil {
				t.Fatal(err)
			}
			original := priorReportContent("fail")
			mustWriteForTest(t, report, original)
			// The formatter shell and its child both remain alive until cancellation.
			p.FormatCommand = qaFormatterScript(t, "printf 'broken' > \"$1\"\nprintf ready > '"+ready+"'\nsleep 30\nprintf late > \"$1\"\n")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			done := make(chan error, 1)
			go func() {
				ticker := time.NewTicker(10 * time.Millisecond)
				defer ticker.Stop()
				for {
					if _, err := os.Stat(ready); err == nil {
						if name == "context" {
							cancel()
							done <- nil
						} else {
							done <- f.store.RequestStop(ctx, p.RunID)
						}
						return
					}
					select {
					case <-ctx.Done():
						done <- ctx.Err()
						return
					case <-ticker.C:
					}
				}
			}()
			err := e.formatQADirectory(ctx, p, 1, "commit", []string{report}, "fail")
			triggerErr := <-done
			if triggerErr != nil {
				t.Fatal(triggerErr)
			}
			expected := context.Canceled
			if name == "stop request" {
				expected = ErrStopRequested
			}
			if !errors.Is(err, expected) {
				t.Fatalf("format error = %v, want %v", err, expected)
			}
			got, err := os.ReadFile(report)
			if err != nil || string(got) != original {
				t.Fatalf("restored report = %q %v", got, err)
			}
			assertQAFormatOutcome(t, f, "commit", "reverted", 1)
		})
	}
}

func TestQAFormatStartFailureRestoresFiles(t *testing.T) {
	f, p, e := priorFixture(t)
	report := filepath.Join(p.Spec.Dir, "qa", priorReportName(1))
	if err := os.MkdirAll(filepath.Dir(report), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, report, priorReportContent("fail"))
	p.FormatCommand = "formatter"
	t.Setenv("PATH", t.TempDir()) // sh itself cannot start.
	if err := e.formatQADirectory(context.Background(), p, 1, "commit", []string{report}, "fail"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(report)
	if err != nil || string(got) != priorReportContent("fail") {
		t.Fatalf("report = %q %v", got, err)
	}
	assertQAFormatOutcome(t, f, "commit", "failed", 1)
}

func TestQAFormatBoundsDiagnostics(t *testing.T) {
	script := qaFormatterScript(t, "i=0; while [ \"$i\" -lt 3000 ]; do printf x; i=$((i+1)); done; printf tail >&2; exit 1\n")
	f, _, _, _, _ := qaFormatCommitFixture(t, script)
	payload := assertQAFormatOutcome(t, f, "commit", "failed", 1)
	diagnostics := payload["diagnostics"].(string)
	if len(diagnostics) != 2048 || !strings.HasSuffix(diagnostics, "tail") {
		t.Fatalf("diagnostics length=%d tail=%q", len(diagnostics), diagnostics)
	}
}
