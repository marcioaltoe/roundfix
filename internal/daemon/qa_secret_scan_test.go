// Suite: QA Report secret refusal
// Invariant: evidence with a Secret Match stays uncommitted and reports only its location.
// Boundary IN: TaskCycle with the existing fake QA Agent and real temporary Git repository
// Boundary OUT: ACP runtime and push scanning (task_03)
package daemon

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func qaSecretScanFixture(t *testing.T, evidence string) (*taskCycleFixture, TaskCycleResult, string) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	path := filepath.ToSlash(filepath.Join("docs", "specs", plan.Spec.Slug, "qa", "evidence", "command.txt"))
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass),
		afterQA: func() {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(plan.WorkDir, path)), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, filepath.Join(plan.WorkDir, path), "command output\n"+evidence+"\n")
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	return fixture, result, path
}

func TestQAReportWithASecretMatchSettlesFailedWithoutACommit(t *testing.T) {
	t.Parallel()
	sample := "sk-or-" + "v1-" + strings.Repeat("b7", 32)
	fixture, result, path := qaSecretScanFixture(t, sample)
	if result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusFailed) {
		t.Fatalf("QA accepted = %t; QA Task must settle failed", result.QAAccepted)
	}
	wantReason := "QA evidence carries a secret: " + path + ":2: openrouter-key"
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
	refusals := 0
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "secret_scan" {
			continue
		}
		refusals++
		wantMatches := []any{map[string]any{"path": path, "line": float64(2), "kind": "openrouter-key"}}
		if event.Summary != "QA Report refused: 1 Secret Match(es)" || !reflect.DeepEqual(payload["matches"], wantMatches) {
			t.Fatalf("refusal event = %s %v", event.Summary, payload)
		}
	}
	if refusals != 1 {
		t.Fatalf("secret_scan events = %d, want 1", refusals)
	}
	history := runGitForTest(t, fixture.gitRoot, "log", "--format=%s")
	if strings.Contains(history, "docs: qa report for "+taskCycleSlug) {
		t.Fatal("refused QA Report was committed")
	}
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit) {
		if strings.HasPrefix(event.Summary, "QA Report commit created:") {
			t.Fatal("refused QA Report published a commit event")
		}
	}
	content, err := os.ReadFile(filepath.Join(fixture.gitRoot, path))
	if err != nil || string(content) != "command output\n"+sample+"\n" {
		t.Fatal("refusal lost or altered the evidence file")
	}
	if _, err := os.Stat(filepath.Join(fixture.gitRoot, result.QAReportPath)); err != nil {
		t.Fatalf("refusal lost the report: %v", err)
	}
	task, err := os.ReadFile(filepath.Join(fixture.gitRoot, taskFileRel(taskCycleSlug, fixture.graph.QATaskID)))
	if err != nil {
		t.Fatal(err)
	}
	outputs := []string{reason, fixture.progress.String(), string(task)}
	for _, event := range fixture.sink.events {
		outputs = append(outputs, event.Summary, string(event.Payload))
	}
	for _, output := range outputs {
		for index := 0; index+8 <= len(sample); index++ {
			if strings.Contains(output, sample[index:index+8]) {
				t.Fatal("QA output contains a sample byte run")
			}
		}
	}
}

func TestQAReportWithoutASecretMatchKeepsItsCommit(t *testing.T) {
	t.Parallel()
	fixture, result, path := qaSecretScanFixture(t, "safe command result")
	if !result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusCompleted) {
		t.Fatal("clean QA Task was not accepted")
	}
	message := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "log", "-1", "--format=%B"))
	if message != QACommitMessage(taskCycleSlug, spec.VerdictPass) {
		t.Fatalf("QA commit message = %q", message)
	}
	if result.QAReportPath != qaReportRelPathForTest() {
		t.Fatalf("QA report = %q, want fixture's current report", result.QAReportPath)
	}
	if committed := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+path); committed != "command output\nsafe command result\n" {
		t.Fatalf("clean evidence commit = %q", committed)
	}
	wantPhases := []string{"prior_report", "mechanical", "verdict", "evidence_snapshots"}
	var phases []string
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		phase, _ := eventPayloadMap(t, event)["phase"].(string)
		phases = append(phases, phase)
	}
	if !reflect.DeepEqual(phases, wantPhases) {
		t.Fatalf("clean QA phases = %v, want %v", phases, wantPhases)
	}
}

func TestQASecretScanUsesTheReportCommitPathSet(t *testing.T) {
	t.Parallel()
	fixture, plan, engine := priorFixture(t)
	report := filepath.ToSlash(filepath.Join("docs", "specs", taskCycleSlug, "qa", qaReportNameForTest()))
	task := taskFileRel(taskCycleSlug, fixture.graph.QATaskID)
	imported := filepath.ToSlash(filepath.Join("docs", "specs", taskCycleSlug, "qa", "prior.txt"))
	sample := "sk-or-" + "v1-" + strings.Repeat("d4", 32)
	mustWriteForTest(t, filepath.Join(plan.WorkDir, "deleted.txt"), "old text\n")
	runGitForTest(t, plan.WorkDir, "add", "deleted.txt")
	runGitForTest(t, plan.WorkDir, "commit", "-m", "test: establish deletion input")
	// These paths are already dirty before QA; only report, Task and imported
	// files are forced into the commit even when the snapshot diff excludes them.
	if err := os.MkdirAll(filepath.Dir(filepath.Join(plan.WorkDir, report)), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{report, imported, "preexisting.txt"} {
		mustWriteForTest(t, filepath.Join(plan.WorkDir, path), sample+"\n")
	}
	originalTask, err := os.ReadFile(filepath.Join(plan.WorkDir, task))
	if err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(plan.WorkDir, task), string(originalTask)+"\n"+sample+"\n")
	before, err := engine.snapshotQAPaths(t.Context(), plan.WorkDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(plan.WorkDir, "deleted.txt")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"changed.txt", "verification.txt", projectConfigPath, "executable.txt"} {
		mustWriteForTest(t, filepath.Join(plan.WorkDir, path), sample+"\n")
	}
	if err := os.Chmod(filepath.Join(plan.WorkDir, "executable.txt"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(plan.WorkDir, "binary.bin"), "\x00"+sample)
	prepared, err := engine.prepareQAReportCommit(t.Context(), plan, before, []string{"verification.txt"}, report, plan.Tasks[len(plan.Tasks)-1], imported)
	if err != nil {
		t.Fatal(err)
	}
	wantPaths := []string{"binary.bin", "changed.txt", "deleted.txt", imported, report, task}
	if !reflect.DeepEqual(prepared.paths, wantPaths) {
		t.Fatalf("prepared paths = %v, want %v", prepared.paths, wantPaths)
	}
	matches, err := scanQAReportPaths(t.Context(), plan.WorkDir, prepared.paths)
	if err != nil {
		t.Fatal(err)
	}
	var matchedPaths []string
	for _, match := range matches {
		matchedPaths = append(matchedPaths, match.Path)
	}
	if want := []string{"changed.txt", imported, report, task}; !reflect.DeepEqual(matchedPaths, want) {
		t.Fatalf("matched paths = %v, want %v", matchedPaths, want)
	}
	// A read error must refuse scanning rather than silently accepting the file.
	if _, err := scanQAReportPaths(t.Context(), plan.WorkDir, []string{"docs"}); err == nil {
		t.Fatal("directory read error was accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := scanQAReportPaths(ctx, plan.WorkDir, prepared.paths); err == nil {
		t.Fatal("cancelled scan was accepted")
	}
}

func TestQAFormatterOutputIsScannedBeforeTheQACommit(t *testing.T) {
	t.Parallel()
	sample := "sk-or-" + "v1-" + strings.Repeat("e5", 32)
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	path := filepath.ToSlash(filepath.Join("docs", "specs", plan.Spec.Slug, "qa", "evidence", "command.txt"))
	// The Agent writes clean evidence; only the formatter adds the sample.
	plan.FormatCommand = qaFormatterScript(t, "for path do case \"$path\" in *command.txt) printf '%s\\n' '"+sample+"' >> \"$path\";; esac; done\n")
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass),
		afterQA: func() {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(plan.WorkDir, path)), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, filepath.Join(plan.WorkDir, path), "command output\n")
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusFailed) {
		t.Fatalf("QA accepted = %t; formatter output with a secret must settle the QA Task failed", result.QAAccepted)
	}
	if history := runGitForTest(t, fixture.gitRoot, "log", "--format=%s"); strings.Contains(history, "docs: qa report for "+taskCycleSlug) {
		t.Fatal("formatter output with a secret was committed")
	}
	refusals := 0
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "secret_scan" {
			continue
		}
		refusals++
		wantMatches := []any{map[string]any{"path": path, "line": float64(2), "kind": "openrouter-key"}}
		if !reflect.DeepEqual(payload["matches"], wantMatches) {
			t.Fatalf("refusal event = %v", payload)
		}
	}
	if refusals != 1 {
		t.Fatalf("secret_scan events = %d, want 1", refusals)
	}
}

func TestQAFormatterNeverSeesEvidenceWithASecretMatch(t *testing.T) {
	t.Parallel()
	sample := "sk-or-" + "v1-" + strings.Repeat("f6", 32)
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	path := filepath.ToSlash(filepath.Join("docs", "specs", plan.Spec.Slug, "qa", "evidence", "command.txt"))
	// A failing formatter that echoes its input would put the sample in its diagnostic.
	plan.FormatCommand = qaFormatterScript(t, "cat \"$@\" >&2\nexit 1\n")
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass),
		afterQA: func() {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(plan.WorkDir, path)), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, filepath.Join(plan.WorkDir, path), "command output\n"+sample+"\n")
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusFailed) {
		t.Fatalf("QA accepted = %t; evidence with a secret must settle the QA Task failed", result.QAAccepted)
	}
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		if eventPayloadMap(t, event)["phase"] == "format" {
			t.Fatal("formatter ran over evidence with a Secret Match")
		}
	}
	outputs := []string{fixture.progress.String()}
	for _, event := range fixture.sink.events {
		outputs = append(outputs, event.Summary, string(event.Payload))
	}
	for _, output := range outputs {
		for index := 0; index+8 <= len(sample); index++ {
			if strings.Contains(output, sample[index:index+8]) {
				t.Fatal("QA output contains a sample byte run")
			}
		}
	}
}
