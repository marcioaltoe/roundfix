// Suite: QA evidence snapshot settlement
// Invariant: the QA Report commit holds Daemon-recorded snapshots at the audited head.
// Boundary IN: TaskCycle, task-cycle fixture, real Git settlement and Run events
// Boundary OUT: ACP runtime and recorder qualification (internal/speccheck/evidence_record_test.go)
package daemon

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead(t *testing.T) {
	testQAEvidenceSnapshot(t, true)
}

func TestQAGateStripsAnAgentWrittenSnapshotWhenNoRowQualifies(t *testing.T) {
	testQAEvidenceSnapshot(t, false)
}

func testQAEvidenceSnapshot(t *testing.T, qualifies bool) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	mustWriteForTest(t, fixture.gitRoot+"/snapshot-input.txt", "audited input\n")
	runGitForTest(t, fixture.gitRoot, "add", "snapshot-input.txt")
	runGitForTest(t, fixture.gitRoot, "commit", "-m", "test: establish snapshot input")
	head := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "rev-parse", "HEAD"))
	status := "fail"
	if qualifies {
		status = "pass"
	}
	report := fmt.Sprintf("---\nverdict: pass\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\nevidence_snapshots: agent-written\n---\n\n# QA Report\n\n## Results\n\n| # | Status | Provenance | Evidence |\n| --- | --- | --- | --- |\n| 3 | %s | criterion | observed |\n\n### 3 evidence\n\n```yaml\ninputs:\n  - kind: repository_path\n    ref: snapshot-input.txt\n```\n", status)
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: report}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}
	committed := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+result.QAReportPath)
	if strings.Contains(committed, "agent-written") {
		t.Fatalf("Agent snapshot survives: %s", committed)
	}
	outcome := "none"
	rows := float64(0)
	if qualifies {
		outcome = "recorded"
		rows = 1
		if !strings.Contains(committed, "evidence_snapshots:") || !strings.Contains(committed, "head: "+head) || !strings.Contains(committed, fmt.Sprintf("{ref: snapshot-input.txt, count: 1, sha256: %x}", sha256.Sum256([]byte(fmt.Sprintf("%x  snapshot-input.txt\n", sha256.Sum256([]byte("audited input\n"))))))) {
			t.Fatalf("committed snapshot missing: %s", committed)
		}
	} else if strings.Contains(committed, "evidence_snapshots:") {
		t.Fatal(committed)
	}
	found := false
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "evidence_snapshots" {
			continue
		}
		found = true
		if payload["outcome"] != outcome || payload["head"] != head || payload["rows"] != rows || payload["report"] != result.QAReportPath || payload["error"] != "" {
			t.Fatalf("event=%+v", payload)
		}
	}
	if !found {
		t.Fatal("missing evidence_snapshots event")
	}
}

func TestQAEvidenceSnapshotGitErrorDoesNotPreventSettlement(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	path := qaReportRelPathForTest()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(fixture.gitRoot, path)), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, fixture.gitRoot+"/"+path, "---\nverdict: pass\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\nevidence_snapshots: untrusted\n---\n\n## Results\n\n| # | Status |\n| --- | --- |\n| 1 | pass |\n\n### 1 evidence\n```yaml\ninputs:\n  - kind: repository_path\n    ref: missing.txt\n```\n")
	engine := fixture.engine(t, &taskFakeRunner{}, &taskFakeVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})
	if err := engine.recordQAEvidenceSnapshots(context.Background(), plan, 1, path, "invalid-head"); err != nil {
		t.Fatalf("Git error blocked settlement: %v", err)
	}
	events := taskEventsOfKind(fixture.sink, runevent.KindDaemonQA)
	payload := eventPayloadMap(t, events[len(events)-1])
	if payload["outcome"] != "error" || payload["error"] == "" {
		t.Fatalf("event=%+v", payload)
	}
	content, err := os.ReadFile(fixture.gitRoot + "/" + path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "evidence_snapshots:") {
		t.Fatal(string(content))
	}
}

func TestQAEvidenceSnapshotSkipsAnExternalReport(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	path := t.TempDir() + "/report.md"
	const content = "---\nevidence_snapshots: external\n---\n"
	mustWriteForTest(t, path, content)
	engine := fixture.engine(t, &taskFakeRunner{}, &taskFakeVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})
	if err := engine.recordQAEvidenceSnapshots(context.Background(), plan, 1, path, "head"); err != nil {
		t.Fatal(err)
	}
	events := taskEventsOfKind(fixture.sink, runevent.KindDaemonQA)
	if eventPayloadMap(t, events[len(events)-1])["outcome"] != "skipped" {
		t.Fatal(events)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != content {
		t.Fatalf("external report changed: %q %v", after, err)
	}
}

func TestQAEvidenceSnapshotFileErrorIsInfrastructure(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	engine := fixture.engine(t, &taskFakeRunner{}, &taskFakeVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})
	path := qaReportRelPathForTest()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(fixture.gitRoot, path)), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, fixture.gitRoot+"/"+path, qaReportForTest(spec.VerdictPass))
	if err := os.Chmod(fixture.gitRoot+"/"+path, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(fixture.gitRoot+"/"+path, 0o644) })
	if err := engine.recordQAEvidenceSnapshots(context.Background(), fixture.qaPlan(), 1, path, "head"); err == nil {
		t.Fatal("report write failure did not stop QA")
	}
}

func TestQAEvidenceSnapshotCancellationPublishesStop(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	engine := fixture.engine(t, &taskFakeRunner{}, &taskFakeVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := engine.recordQAEvidenceSnapshots(ctx, fixture.qaPlan(), 1, "", "head")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
	if len(taskEventsOfKind(fixture.sink, runevent.KindDaemonStatus)) == 0 {
		t.Fatal("missing stop event")
	}
}
