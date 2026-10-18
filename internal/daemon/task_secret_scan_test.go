// Suite: Task commit Secret Scan
// Invariant: Only newly added secret lines refuse settlement and no value reaches diagnostics.
// Boundary IN: TaskCycle, existing fake Agent, and real temporary Git repository.
// Boundary OUT: QA Report scanning and push scanning.
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

func taskSecretScanFixture(t *testing.T, initial, content string) (*taskCycleFixture, TaskCycleResult, string, string) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Write a fixture"}})
	path := "fixture.txt"
	if initial != "" {
		mustWriteForTest(t, filepath.Join(fixture.gitRoot, path), initial)
	}
	commitTaskFixtureSource(t, fixture.gitRoot, "test: establish Task scan input")
	head := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "rev-parse", "HEAD"))
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot,
		afterTask: func(string) { mustWriteForTest(t, filepath.Join(fixture.gitRoot, path), content) },
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), fixture.plan())
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	return fixture, result, path, head
}

func TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit(t *testing.T) {
	t.Parallel()
	value := strings.Repeat("K8", 15)
	content := "clean\n" + "SERVICE_API_KEY" + " = " + value + "\n"
	// Cover both paths absent from HEAD and paths already committed there.
	for _, test := range []struct{ name, initial string }{
		{"new file", ""},
		{"tracked file", "clean\nold fixture\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture, result, path, head := taskSecretScanFixture(t, test.initial, content)
			if result.Failed != 1 || result.Completed != 0 || taskStatusOnDisk(t, fixture.gitRoot, "task_01") != string(spec.StatusFailed) {
				t.Fatalf("Task result = %+v, want failed without completion", result)
			}
			wantReason := "Task commit carries a secret: " + path + ":2: credential-assignment"
			refusals := 0
			var reason string
			for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
				payload := eventPayloadMap(t, event)
				if payload["phase"] == "settled" {
					reason, _ = payload["reason"].(string)
				}
				if payload["phase"] != "secret_scan" {
					continue
				}
				refusals++
				wantMatches := []any{map[string]any{"path": path, "line": float64(2), "kind": "credential-assignment"}}
				if event.ReviewIssue != "task_01" || event.Summary != "Task task_01 commit refused: 1 Secret Match(es)" || !reflect.DeepEqual(payload["matches"], wantMatches) {
					t.Fatalf("refusal event = %s %v", event.Summary, payload)
				}
			}
			if refusals != 1 || reason != wantReason {
				t.Fatalf("refusals = %d, reason = %q, want one refusal with %q", refusals, reason, wantReason)
			}
			if got := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "rev-parse", "HEAD")); got != head {
				t.Fatal("refused Task created a commit")
			}
			if events := taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit); len(events) != 0 {
				t.Fatal("refused Task published a commit event")
			}
			preserved, err := os.ReadFile(filepath.Join(fixture.gitRoot, path))
			if err != nil || string(preserved) != content {
				t.Fatal("refusal lost or changed the Task output")
			}
			taskContent, err := os.ReadFile(filepath.Join(fixture.gitRoot, taskFileRel(taskCycleSlug, "task_01")))
			if err != nil {
				t.Fatal(err)
			}
			outputs := []string{reason, fixture.progress.String(), string(taskContent)}
			for _, event := range fixture.sink.events {
				outputs = append(outputs, event.Summary, string(event.Payload))
			}
			for _, output := range outputs {
				for index := 0; index+8 <= len(value); index++ {
					if strings.Contains(output, value[index:index+8]) {
						t.Fatal("Task diagnostic contains a sample byte run")
					}
				}
			}
		})
	}
}

func TestTaskCommitScansOnlyTheLinesItAdds(t *testing.T) {
	t.Parallel()
	sample := "SERVICE_API_KEY" + " = " + strings.Repeat("K8", 15) + "\n"
	content := sample + "changed neighboring line\n"
	fixture, result, path, head := taskSecretScanFixture(t, sample+"old neighboring line\n", content)
	assertTaskSecretScanCommitted(t, fixture, result, path, head, content)
}

func TestTaskCommitAcceptsALettersOnlyFixture(t *testing.T) {
	t.Parallel()
	content := "SERVICE_API_KEY" + " = " + strings.Repeat("z", 30) + "\n"
	fixture, result, path, head := taskSecretScanFixture(t, "", content)
	assertTaskSecretScanCommitted(t, fixture, result, path, head, content)
}

func assertTaskSecretScanCommitted(t *testing.T, fixture *taskCycleFixture, result TaskCycleResult, path, head, content string) {
	t.Helper()
	if result.Completed != 1 || result.Failed != 0 || taskStatusOnDisk(t, fixture.gitRoot, "task_01") != string(spec.StatusCompleted) {
		t.Fatalf("Task result = %+v, want completed without refusal", result)
	}
	if got := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "rev-parse", "HEAD")); got == head {
		t.Fatal("accepted Task did not commit")
	}
	if got := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+path); got != content {
		t.Fatal("accepted Task did not commit the expected content")
	}
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if eventPayloadMap(t, event)["phase"] == "secret_scan" {
			t.Fatal("accepted Task published a secret refusal")
		}
	}
}

func TestTaskSecretScanUsesHeadAndPreparedPaths(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	sample := "SERVICE_API_KEY" + " = " + strings.Repeat("K8", 15) + "\n"
	for _, path := range []string{"tracked.txt", "deleted.txt"} {
		mustWriteForTest(t, filepath.Join(fixture.gitRoot, path), sample+"old\n")
	}
	commitTaskFixtureSource(t, fixture.gitRoot, "test: seed Task scan paths")
	if err := os.Remove(filepath.Join(fixture.gitRoot, "deleted.txt")); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, "tracked.txt"), sample+"new\n"+sample)
	for _, path := range []string{"staged.txt", "untracked.txt", "excluded.txt"} {
		mustWriteForTest(t, filepath.Join(fixture.gitRoot, path), sample)
	}
	runGitForTest(t, fixture.gitRoot, "add", "staged.txt")
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, "binary.bin"), "\x00"+sample)
	paths := []string{"tracked.txt", "deleted.txt", "staged.txt", "untracked.txt", "binary.bin"}
	matches, err := scanTaskCommitPaths(t.Context(), fixture.gitRoot, paths)
	if err != nil {
		t.Fatal(err)
	}
	locations := make(map[string]int)
	for _, match := range matches {
		if match.Kind != "credential-assignment" || match.Commit != "" {
			t.Fatalf("unexpected match metadata: %+v", match)
		}
		locations[match.Path] = match.Line
	}
	want := map[string]int{"tracked.txt": 3, "staged.txt": 1, "untracked.txt": 1}
	if len(matches) != len(want) || !reflect.DeepEqual(locations, want) {
		t.Fatalf("locations = %v, want %v", locations, want)
	}
	if _, err := scanTaskCommitPaths(t.Context(), fixture.gitRoot, []string{"docs"}); err == nil {
		t.Fatal("directory read error was accepted")
	}
	if _, err := scanTaskCommitPaths(t.Context(), t.TempDir(), []string{"missing.txt"}); err == nil {
		t.Fatal("Git error was accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := scanTaskCommitPaths(ctx, fixture.gitRoot, paths); err == nil {
		t.Fatal("cancelled scan was accepted")
	}
}
