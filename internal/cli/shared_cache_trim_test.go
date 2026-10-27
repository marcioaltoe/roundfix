// Suite: automatic Shared Build Cache trims at Run start and the repository gate.
// Invariant: trims preserve active entries and never decide the Run or gate outcome.
// Boundary IN: command harness, real cache files, SQLite events and gate Verification.
// Boundary OUT: Agent execution and Linux Lane tools use their existing fixtures.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/runevent"
	"roundfix/internal/sharedcache"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func sharedCacheImplementWorkspace(t *testing.T) (sharedCacheFixture, string, string) {
	t.Helper()
	f := newSharedCacheFixture(t, 1)
	env := commandEnvironmentForTest(t)
	home, repo := env.homeDir, env.workDir
	gittest.InitRepo(t, repo, "--initial-branch=main")
	writeImplementSpec(t, repo, implementTestSlug, []implementSeed{{id: "task_01", title: "Build the widget"}})
	// Give the gate an immutable default with Linux Lane disabled.
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), mustRead(t, filepath.Join(repo, ".roundfixrc.yml"))+"delivery:\n  linux_lane: false\n")
	gitImplement(t, repo, "add", "-A")
	gitImplement(t, repo, "commit", "-m", "seed cache and spec")
	gitImplement(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	gitImplement(t, repo, "checkout", "-b", "feat/widget")
	return f, home, repo
}

const sharedCacheTaskOutput = "task_01 completed — Build the widget\nClean: all 1 Task(s) completed.\nTokens: none reported by 1 prompt(s); cost not reported\n"

func runSharedCacheImplement(t *testing.T, home, repo string) (store.Run, string) {
	t.Helper()
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	withImplementCollaborators(t, runner)
	var out, diagnostics bytes.Buffer
	code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diagnostics)
	if code != exitOK || runner.calls != 1 || out.String() != sharedCacheTaskOutput {
		t.Fatalf("exit=%d calls=%d stdout=%q stderr=%q", code, runner.calls, out.String(), diagnostics.String())
	}
	run := onlyRunFromStore(t, home)
	if run.State != store.StateClean {
		t.Fatalf("Run state = %s", run.State)
	}
	return run, diagnostics.String()
}

type sharedCacheTrimPayload struct {
	Event   string `json:"event"`
	Trigger string `json:"trigger"`
	Path    string `json:"path"`
	Outcome string `json:"outcome"`
	Size    int64  `json:"size_bytes"`
	Limit   int64  `json:"limit_bytes"`
	Entries int    `json:"removed_entries"`
	Bytes   int64  `json:"removed_bytes"`
	Cause   string `json:"cause"`
}

func sharedCacheTrimEvents(t *testing.T, home string, run store.Run) ([]store.JournalEvent, []sharedCacheTrimPayload) {
	t.Helper()
	reader, err := store.OpenReader(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Error(err)
		}
	}()
	events, err := reader.RunEventsOfKinds(t.Context(), run.ID, runevent.KindDaemonStatus, runevent.KindDaemonTask)
	if err != nil {
		t.Fatal(err)
	}
	var trims []store.JournalEvent
	var payloads []sharedCacheTrimPayload
	for _, event := range events {
		var payload sharedCacheTrimPayload
		if err := json.Unmarshal(event.Event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Event == "shared_cache_trim" {
			if event.Event.Source != runevent.SourceDaemon || event.Event.Kind != runevent.KindDaemonStatus || event.Event.Summary != runevent.BoundSummary("Shared Build Cache "+payload.Path+": "+payload.Outcome) {
				t.Fatalf("unexpected event: %+v", event)
			}
			trims = append(trims, event)
			payloads = append(payloads, payload)
		}
	}
	if len(trims) > 0 {
		foundTask := false
		for _, event := range events {
			if event.Event.Kind == runevent.KindDaemonTask {
				foundTask = true
				for _, trim := range trims {
					if trim.Cursor >= event.Cursor {
						t.Fatal("trim event followed a Task event")
					}
				}
				break
			}
		}
		if !foundTask {
			t.Fatal("no Task event")
		}
	}
	return trims, payloads
}

func TestImplementStartTrimsTheDeclaredSharedBuildCache(t *testing.T) {
	t.Parallel()
	f, home, repo := sharedCacheImplementWorkspace(t)
	// Observe files at the worktree boundary, before any Task can start.
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	withImplementCollaborators(t, runner)
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		create := deps.createRunWorktree
		deps.createRunWorktree = func(ctx context.Context, opts runworktree.CreateOptions) (runworktree.Ref, error) {
			assertPathMissing(t, f.paths[0])
			assertPathMissing(t, f.paths[1])
			assertPathExists(t, f.paths[2])
			return create(ctx, opts)
		}
	})
	var out, diagnostics bytes.Buffer
	code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diagnostics)
	if code != exitOK || out.String() != sharedCacheTaskOutput {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
	run := onlyRunFromStore(t, home)
	if run.State != store.StateClean {
		t.Fatalf("Run state = %s", run.State)
	}
	_, payloads := sharedCacheTrimEvents(t, home, run)
	if len(payloads) != 1 {
		t.Fatalf("trim events = %d", len(payloads))
	}
	want := sharedCacheTrimPayload{Event: "shared_cache_trim", Trigger: "run_start", Path: f.root, Outcome: "trimmed", Size: f.size, Limit: 1048576, Entries: 2, Bytes: 2097152}
	if payloads[0] != want {
		t.Fatalf("payload=%+v want=%+v", payloads[0], want)
	}
	line := fmt.Sprintf("roundfix: Shared Build Cache %s: removed 2 idle entries (2097152 bytes) at Run start\n", f.root)
	if strings.Count(diagnostics.String(), line) != 1 {
		t.Fatalf("stderr=%q", diagnostics.String())
	}
	assertPathMissing(t, f.paths[0])
	assertPathMissing(t, f.paths[1])
	assertPathExists(t, f.paths[2])
	assertFileContent(t, filepath.Join(f.root, "README"), sharedcache.GoBuildReadme+"\n")
}

func TestImplementStartRunsOnWhenTheTrimDoesNot(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"refused cache", "budget ends"} {
		t.Run(name, func(t *testing.T) {
			f, home, repo := sharedCacheImplementWorkspace(t)
			if name == "refused cache" {
				if err := os.Remove(filepath.Join(f.root, "README")); err != nil {
					t.Fatal(err)
				}
			} else {
				// Zero expires synchronously, without waiting for a real-time timer.
				updateCommandDependenciesForTest(t, func(deps *commandDependencies) { deps.gc.sharedCacheTrimBudget = 0 })
			}
			run, diagnostics := runSharedCacheImplement(t, home, repo)
			_, payloads := sharedCacheTrimEvents(t, home, run)
			if len(payloads) != 1 {
				t.Fatalf("trim events = %d", len(payloads))
			}
			p := payloads[0]
			if p.Path != f.root || p.Trigger != "run_start" || p.Limit != 1048576 || p.Size != 0 || p.Entries != 0 || p.Bytes != 0 {
				t.Fatalf("payload=%+v", p)
			}
			if name == "refused cache" {
				if p.Outcome != "kept" || !strings.Contains(p.Cause, "README") || strings.Contains(diagnostics, "Shared Build Cache") {
					t.Fatalf("payload=%+v stderr=%q", p, diagnostics)
				}
			} else {
				line := fmt.Sprintf("roundfix: Shared Build Cache %s: trim stopped after its 0s budget at Run start; removed 0 idle entries (0 bytes); the next trim continues\n", f.root)
				if p.Outcome != "stopped" || p.Cause != "budget of 0s elapsed" || strings.Count(diagnostics, line) != 1 {
					t.Fatalf("payload=%+v stderr=%q", p, diagnostics)
				}
			}
			for _, path := range f.paths {
				assertPathExists(t, path)
			}
		})
	}
}

func TestImplementStartWithoutASharedBuildCacheRecordsNoTrim(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Build the widget"}})
	run, diagnostics := runSharedCacheImplement(t, home, repo)
	events, _ := sharedCacheTrimEvents(t, home, run)
	if len(events) != 0 || strings.Contains(diagnostics, "Shared Build Cache") {
		t.Fatalf("events=%v stderr=%q", events, diagnostics)
	}
}

func TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"idle entries removed first", "budget ends"} {
		t.Run(name, func(t *testing.T) {
			f, home, repo := sharedCacheImplementWorkspace(t)
			loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil {
				t.Fatal(err)
			}
			loaded.Config.Defaults.Verification = "test ! -e '" + f.paths[0] + "'"
			if name == "budget ends" {
				loaded.Config.Defaults.Verification = "true"
				// Zero expires synchronously, without waiting for a real-time timer.
				updateCommandDependenciesForTest(t, func(deps *commandDependencies) { deps.gc.sharedCacheTrimBudget = 0 })
			}
			var log bytes.Buffer
			workflow := commandDeliveryWorkflow{loaded: loaded, log: &log}
			ctx := commandContextForTest(t, t.Context())
			result, err := workflow.Gate(ctx, repo, implementTestSlug, "HEAD")
			if err != nil || !result.Passed || result.LaneSkip != "" {
				t.Fatalf("gate=%+v err=%v log=%q", result, err, log.String())
			}
			line := fmt.Sprintf("roundfix: Delivery Queue item %s: Shared Build Cache %s: removed 2 idle entries (2097152 bytes) before the repository gate\n", implementTestSlug, f.root)
			if name == "budget ends" {
				line = fmt.Sprintf("roundfix: Delivery Queue item %s: Shared Build Cache %s: trim stopped after its 0s budget before the repository gate; removed 0 idle entries (0 bytes); the next trim continues\n", implementTestSlug, f.root)
				for _, path := range f.paths {
					assertPathExists(t, path)
				}
			} else {
				assertPathMissing(t, f.paths[0])
				assertPathMissing(t, f.paths[1])
				assertPathExists(t, f.paths[2])
				assertFileContent(t, filepath.Join(f.root, "README"), sharedcache.GoBuildReadme+"\n")
			}
			if log.String() != line {
				t.Fatalf("log=%q want=%q", log.String(), line)
			}
		})
	}
}

func TestSharedCacheTrimReportsEveryDeclarationAfterCancellation(t *testing.T) {
	t.Parallel()
	f := newSharedCacheFixture(t, 1)
	env := commandEnvironmentForTest(t)
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: env.homeDir, WorkDir: env.workDir})
	if err != nil {
		t.Fatal(err)
	}
	// Include a later refused cache and a later path that cannot resolve.
	missing := loaded.Config.Worktree.SharedCaches[0]
	missing.Path = filepath.Join(env.homeDir, "missing-cache")
	invalid := missing
	invalid.Path = "relative-path"
	loaded.Config.Worktree.SharedCaches = append(loaded.Config.Worktree.SharedCaches, missing, invalid)
	ctx, cancel := context.WithCancel(commandContextForTest(t, t.Context()))
	cancel()
	reports := trimSharedCaches(ctx, loaded, false)
	if len(reports) != 3 {
		t.Fatalf("reports=%+v", reports)
	}
	for _, report := range reports {
		if report.outcome != sharedCacheStopped || !strings.Contains(report.cause, context.Canceled.Error()) || report.size != 0 || report.entries != 0 || report.bytes != 0 {
			t.Fatalf("report=%+v", report)
		}
	}
	for _, path := range f.paths {
		assertPathExists(t, path)
	}
}

func TestSharedCacheTrimClassifiesAnUnchangedOrFailedCache(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unchanged", "failed path"} {
		t.Run(name, func(t *testing.T) {
			f := newSharedCacheFixture(t, 9)
			env := commandEnvironmentForTest(t)
			loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: env.homeDir, WorkDir: env.workDir})
			if err != nil {
				t.Fatal(err)
			}
			if name == "failed path" {
				loaded.Config.Worktree.SharedCaches[0].Path = "relative-path"
			}
			reports := trimSharedCaches(commandContextForTest(t, t.Context()), loaded, false)
			if len(reports) != 1 {
				t.Fatalf("reports=%+v", reports)
			}
			report := reports[0]
			var diagnostics bytes.Buffer
			printSharedCacheTrim(&diagnostics, report, 30*time.Second, "", "at Run start", "the Run continues")
			if name == "unchanged" {
				if report.outcome != sharedCacheUnchanged || report.cause != "" || report.size != f.size || diagnostics.Len() != 0 {
					t.Fatalf("report=%+v stderr=%q", report, diagnostics.String())
				}
			} else {
				want := fmt.Sprintf("roundfix: warning: Shared Build Cache relative-path: trim failed at Run start: %s; the Run continues\n", report.cause)
				if report.outcome != sharedCacheFailed || report.cause == "" || diagnostics.String() != want {
					t.Fatalf("report=%+v stderr=%q", report, diagnostics.String())
				}
			}
			for _, path := range f.paths {
				assertPathExists(t, path)
			}
		})
	}
}

func TestImplementStartRunsOnWhenTheTrimEventCannotBeJournaled(t *testing.T) {
	t.Parallel()
	f, home, repo := sharedCacheImplementWorkspace(t)
	_, _ = retentionStartStore(t, home)
	db := openGCRetentionDB(t, home)
	gcRetentionExec(t, db, `CREATE TRIGGER refuse_trim_event BEFORE INSERT ON run_events
 WHEN json_extract(NEW.payload, '$.event') = 'shared_cache_trim'
 BEGIN SELECT RAISE(ABORT, 'trim event refused'); END`)
	run, diagnostics := runSharedCacheImplement(t, home, repo)
	line := "roundfix: warning: Shared Build Cache " + f.root + ": trim event not journaled: "
	if strings.Count(diagnostics, line) != 1 || !strings.Contains(diagnostics, "trim event refused") || !strings.Contains(diagnostics, "; the Run continues\n") {
		t.Fatalf("stderr=%q", diagnostics)
	}
	events, _ := sharedCacheTrimEvents(t, home, run)
	if len(events) != 0 {
		t.Fatalf("unexpected journaled trim events: %+v", events)
	}
	assertPathMissing(t, f.paths[0])
	assertPathMissing(t, f.paths[1])
	assertPathExists(t, f.paths[2])
}
