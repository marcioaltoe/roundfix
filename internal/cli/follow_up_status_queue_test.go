// Suite: Spec Queue Table in the status view.
// Invariant: forecast times come from retained deliveries, conditions remain undated,
// and the advisory report preserves its existing fields and read-only behavior.
// Boundary IN: private homes, real temporary Git history, SQLite, clock and zone.
// Boundary OUT: forge and judge calls use per-test fakes; no Runs are started.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func newStatusQueueFixture(t *testing.T, now time.Time) (followUpFixture, followUpQueueInput) {
	t.Helper()
	in := newFollowUpReadInput(t, now)
	f := followUpFixture{repo: in.Loaded.GitRoot, home: in.Loaded.HomeDir, brain: t.TempDir(), judge: &followUpFakeJudge{t: t}}
	setCommandEnvironmentForTest(t, f.home, f.repo)
	updateCommandDependenciesForTest(t, func(d *commandDependencies) {
		d.environ = func() []string { return []string{} }
		d.getenv = func(string) string { return "" }
		d.followUpGHRunner = followUpFakeGH{t: t, root: f.repo, response: "[]"}
		d.judgeTransport = f.judge
		d.followUpNow = func() time.Time { return now }
		d.runsListLocation = func() *time.Location { return time.UTC }
	})
	return f, in
}

func statusQueueTableFixture(t *testing.T, now time.Time) (followUpFixture, followUpQueueInput) {
	t.Helper()
	f, in := newStatusQueueFixture(t, now)
	for _, slug := range []string{"0009-queued", "0002-planned"} {
		writeImplementSpecAtRoot(t, in.SpecRoot, slug, []implementSeed{{id: "task_01", status: "completed"}, {id: "task_02", status: "pending"}})
	}
	mustMkdir(t, filepath.Join(in.SpecRoot, "0003-authoring"))
	followUpQueuePhase(t, in, now.Add(-time.Hour), "prd", "0003-authoring", f.repo)
	followUpWrite(t, f.repo, "docs/backlog/open.md", "---\ntype: feat\nstatus: open\ncreated: "+now.AddDate(0, 0, -2).Format("2006-01-02")+"\n---\n# Open entry\n")
	followUpWrite(t, f.repo, "docs/backlog/waiting.md", "---\ntype: feat\nstatus: open\ncreated: "+now.AddDate(0, 0, -1).Format("2006-01-02")+"\nwaiting_on: the maintainer\nsince: "+now.AddDate(0, 0, -1).Format("2006-01-02")+"\n---\n# Waiting entry\n")
	followUpCommitAt(t, f.repo, now.Add(-2*time.Hour), "docs: plan queue")
	gitImplement(t, f.repo, "tag", "v0.4.1")
	followUpStatusSeedQueue(t, f, []string{"0009-queued"}, func(q *store.DeliveryQueue) {
		q.Items[0].Stage = store.DeliveryStageRunning
		q.Items[0].Blocker = "checks pending\nsecond line"
	})
	return f, in
}

func assertStatusQueueRows(t *testing.T, doc followUpStatusDocument, now time.Time) {
	t.Helper()
	if len(doc.Queue) != 5 {
		t.Fatalf("queue=%+v", doc.Queue)
	}
	kinds := []string{"delivery_queue", "planned", "authoring", "backlog", "backlog"}
	specs := []string{"0009-queued", "0002-planned", "0003-authoring", "", ""}
	statuses := []string{"running 1/2 Tasks", "not queued 1/2 Tasks", "authoring: prd", "not yet a Spec", "not yet a Spec"}
	notes := []string{"checks pending", "", "prd since " + now.Add(-time.Hour).Format("01-02 15:04"), "Open entry", "Waiting entry; waits on the maintainer since " + now.AddDate(0, 0, -1).Format("2006-01-02")}
	for i, r := range doc.Queue {
		if r.Row != i+1 || r.Kind != kinds[i] || r.Spec != specs[i] || r.Status != statuses[i] || r.Note != notes[i] {
			t.Fatalf("row %d=%+v", i, r)
		}
		assertQueueCellShape(t, r.Start)
		assertQueueCellShape(t, r.End)
		if i >= 3 && (r.Ref == "" || r.Title == "" || r.Start.State != "none" || r.End.State != "none" || r.Start.Reason != "" || r.End.Reason != "") {
			t.Fatalf("backlog=%+v", r)
		}
	}
}

func assertQueueCellShape(t *testing.T, c followUpQueueCell) {
	t.Helper()
	if c.State != "actual" && c.Actual != nil {
		t.Fatalf("non-actual=%+v", c)
	}
	if c.State != "forecast" && (c.Likely != nil || c.Limit != nil) {
		t.Fatalf("non-forecast=%+v", c)
	}
	if c.State != "condition" && (c.Condition != "" || c.LikelyAfterSeconds != nil || c.LimitAfterSeconds != nil) {
		t.Fatalf("non-condition=%+v", c)
	}
	if c.State != "none" && c.Reason != "" {
		t.Fatalf("non-none=%+v", c)
	}
	for _, v := range []*string{c.Actual, c.Likely, c.Limit} {
		if v != nil {
			at, err := time.Parse(time.RFC3339, *v)
			if err != nil || at.UTC().Format(time.RFC3339) != *v {
				t.Fatalf("UTC seconds timestamp=%q error=%v", *v, err)
			}
		}
	}
}

func TestTheStatusViewShowsTheSpecQueueTable(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 13, 36, 0, 0, time.UTC)
	for _, name := range []string{"text", "json", "state_is_unchanged"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if name == "state_is_unchanged" {
				f, in := newStatusQueueFixture(t, now)
				mustMkdir(t, filepath.Join(in.SpecRoot, "0003-authoring"))
				followUpWrite(t, f.repo, "docs/backlog/open.md", "---\nstatus: open\n---\n# Open\n")
				followUpCommitAt(t, f.repo, now.Add(-time.Hour), "docs: record future work")
				gitImplement(t, f.repo, "tag", "v0.4.1")
				doc := followUpStatusJSON(t)
				if doc.State != "nothing_to_do" || len(doc.Queue) != 2 || len(doc.PendingSpecs) != 0 {
					t.Fatalf("doc=%+v", doc)
				}
				return
			}
			f, _ := statusQueueTableFixture(t, now)
			doc := followUpStatusJSON(t)
			assertStatusQueueRows(t, doc, now)
			if name == "json" {
				return
			}
			out, stderr, code := followUpRunCLI(t, "--status")
			if code != 0 || stderr != "" {
				t.Fatalf("exit=%d stderr=%q", code, stderr)
			}
			pending := strings.Index(out, "Pending Specs:")
			start := strings.Index(out, "Spec queue:")
			end := strings.Index(out, "Next release:")
			if pending < 0 || start < pending || end < start {
				t.Fatalf("sections=%s", out)
			}
			table := out[start:end]
			previous := 0
			for _, want := range []string{"0009  running 1/2 Tasks", "0002  not queued 1/2 Tasks", "0003  authoring: prd", "Open entry", "Waiting entry; waits on the maintainer since 2026-11-15"} {
				at := strings.Index(table, want)
				if at < previous {
					t.Fatalf("missing/out-of-order %q in %s", want, table)
				}
				previous = at
			}
			// The legacy report is pinned in full after removing exactly the new section
			// and the three new source lines.
			legacy := out[:start] + out[end:]
			for _, line := range []string{"  backlog: checked 2\n", "  phase-record: checked 1\n", "  delivery-history: checked 0\n"} {
				legacy = strings.Replace(legacy, line, "", 1)
			}
			want := fmt.Sprintf("Status for widget in %s\nIn progress: 1\n  queue 0 0009-queued running checks pending\nsecond line\nPending Specs: 2\n  0 0009-queued running 1/2 Tasks\n  - 0002-planned not queued 1/2 Tasks\nNext release: no unreleased changes\nBlockers: 0\nWaiting on the maintainer: 0\nOpen authoring questions: wait (default)\nDelegated decisions: 0\nSources:\n  runs: checked 0\n  delivery-queue: checked 1\n  pull-requests: checked 0\n  git: checked 0\n  specs: checked 2\n  release-plan: checked 0\n", f.repo)
			if legacy != want {
				t.Fatalf("legacy report changed:\ngot:\n%s\nwant:\n%s", legacy, want)
			}
		})
	}
}

func TestTheStatusJSONCarriesTheQueueAndTheBasis(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 13, 36, 0, 0, time.UTC)
	f, _ := statusQueueTableFixture(t, now)
	doc := followUpStatusJSON(t)
	assertStatusQueueRows(t, doc, now)
	if doc.State != "work_found" || doc.Repository.Root != f.repo || doc.Repository.Namespace != "widget" || len(doc.InProgress.Runs) != 0 || len(doc.InProgress.DeliveryQueue) != 1 || len(doc.InProgress.PullRequests) != 0 || doc.InProgress.Git.Branch != "main" || doc.InProgress.Git.Head != strings.TrimSpace(gitImplementOutput(t, f.repo, "rev-parse", "--short=7", "HEAD")) || doc.InProgress.Git.Upstream != "" || doc.InProgress.Git.UnpushedCommits != 0 || len(doc.InProgress.Git.ChangedPaths) != 0 || len(doc.Blockers) != 0 || len(doc.Maintainer) != 0 {
		t.Fatalf("existing fields=%+v", doc)
	}
	if !reflect.DeepEqual(doc.PendingSpecs, []followUpPendingSpec{{Spec: "0009-queued", QueuePosition: statusQueueInt(0), Stage: followUpString("running"), TasksCompleted: 1, TasksTotal: 2}, {Spec: "0002-planned", TasksCompleted: 1, TasksTotal: 2}}) {
		t.Fatalf("pending=%+v", doc.PendingSpecs)
	}
	if !reflect.DeepEqual(doc.Release, followUpRelease{State: "no_unreleased_changes", Changes: []followUpReleaseChange{}}) {
		t.Fatalf("release=%+v", doc.Release)
	}
	wantBasis := followUpForecastBasis{Minimum: 11, LikelyPercentile: 50, LimitPercentile: 85}
	if !reflect.DeepEqual(doc.ForecastBasis, wantBasis) {
		t.Fatalf("basis=%+v", doc.ForecastBasis)
	}
	if !reflect.DeepEqual(doc.Sources[6:], []followUpSource{followUpSourceResult("backlog", 2, nil), followUpSourceResult("phase-record", 1, nil), followUpSourceResult("delivery-history", 0, nil)}) {
		t.Fatalf("sources=%+v", doc.Sources)
	}
	out, _, _ := followUpRunCLI(t, "--status", "--format", "json")
	expected := []string{`"pending_specs":`, `"queue":`, `"forecast_basis":`, `"release":`, `"blockers":`, `"maintainer":`, `"authoring_questions":`, `"delegated_decisions":`}
	previous := -1
	for _, field := range expected {
		at := strings.Index(out, field)
		if at <= previous {
			t.Fatalf("field order %q: %s", field, out)
		}
		previous = at
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 14 {
		t.Fatalf("fields=%v", fields)
	}
}
func statusQueueInt(n int) *int { return &n }

// Seed a real default-branch delivery history. All plans precede the first
// start; a retained Implement Run supplies each start and an archive commit
// supplies each end. S is the thirteenth start, three hours after the twelfth.
func statusQueueDeliveryHistory(t *testing.T, now, S time.Time, count int, active bool) (followUpFixture, followUpQueueInput) {
	t.Helper()
	f, in := newStatusQueueFixture(t, now)
	first := S.Add(-12 * 3 * time.Hour)
	for i := 1; i <= count; i++ {
		followUpWrite(t, in.SpecRoot, fmt.Sprintf("%04d-delivered/_prd.md", i), "---\nstatus: completed\n---\n# Delivered\n")
	}
	writeImplementSpecAtRoot(t, in.SpecRoot, "0013-running", []implementSeed{{id: "task_01", status: "pending"}})
	followUpCommitAt(t, f.repo, first.Add(-time.Hour), "docs: plan deliveries")
	for i := 1; i <= count; i++ {
		slug := fmt.Sprintf("%04d-delivered", i)
		start := first.Add(time.Duration(i-1) * 3 * time.Hour)
		followUpArchiveFixture(t, in, slug, spec.ArchiveNoQA)
		followUpCommitAt(t, f.repo, start.Add(time.Duration(i)*10*time.Minute), "docs: archive delivery")
		followUpSeedTimedRun(t, in, slug, store.StateClean, start)
	}
	if active {
		followUpSeedTimedRun(t, in, "0013-running", store.StateActive, S)
	}
	return f, in
}

func TestTheStatusViewForecastsFromTheDeliveryHistory(t *testing.T) {
	t.Parallel()
	S := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	for _, name := range []string{"in_progress", "longer_than_the_time_passed", "past_the_usual_time", "one_two_and_three_behind", "the_fourth_is_a_condition", "a_prerequisite_is_not_counted", "after_the_spec_in_progress_is_merged", "no_spec_in_progress", "a_run_slot", "a_park", "a_parked_spec_that_had_started", "a_spec_being_authored"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			elapsed := 5 * time.Minute
			switch name {
			case "longer_than_the_time_passed":
				elapsed = 75 * time.Minute
			case "past_the_usual_time":
				elapsed = 85 * time.Minute
			case "after_the_spec_in_progress_is_merged":
				elapsed = 185 * time.Minute
			}
			now := S.Add(elapsed)
			f, in := statusQueueDeliveryHistory(t, now, S, 12, name != "no_spec_in_progress")
			slugs := []string{"0013-running"}
			for i := 14; i <= 17; i++ {
				slug := fmt.Sprintf("%04d-planned", i)
				writeImplementSpecAtRoot(t, in.SpecRoot, slug, []implementSeed{{id: "task_01", status: "pending"}})
				slugs = append(slugs, slug)
			}
			if name == "no_spec_in_progress" {
				slugs = slugs[:4]
			}
			if name == "a_spec_being_authored" {
				mustMkdir(t, filepath.Join(in.SpecRoot, "0018-authoring"))
				followUpQueuePhase(t, in, now.Add(-time.Minute), "prd", "0018-authoring", f.repo)
			}
			followUpStatusSeedQueue(t, f, slugs, func(q *store.DeliveryQueue) {
				q.Items[0].Stage = store.DeliveryStageRunning
				switch name {
				case "a_prerequisite_is_not_counted":
					q.Items[2].Stage = store.DeliveryStageParked
					q.Items[2].Blocker = delivery.BlockerPrerequisiteUnmerged + ": 0001-delivered"
				case "a_run_slot":
					q.Items[1].Blocker = delivery.WaitRunSlot + ": holder"
				case "a_park":
					q.Items[1].Stage = store.DeliveryStageParked
					q.Items[1].Blocker = delivery.BlockerChecksTimeout
				case "a_parked_spec_that_had_started":
					q.Items[0].Stage = store.DeliveryStageParked
					q.Items[0].Blocker = delivery.BlockerChecksTimeout
				}
			})
			doc := followUpStatusJSON(t)
			out, stderr, code := followUpRunCLI(t, "--status")
			if code != 0 || stderr != "" {
				t.Fatalf("exit=%d stderr=%s", code, stderr)
			}
			b := doc.ForecastBasis

			gaps := 12
			if name == "no_spec_in_progress" {
				gaps = 11
			}
			if b.Minimum != 11 || b.LikelyPercentile != 50 || b.LimitPercentile != 85 || b.FirstEnd == nil || *b.FirstEnd != S.Add(-36*time.Hour+10*time.Minute).Format(time.RFC3339) || b.LastEnd == nil || *b.LastEnd != S.Add(-time.Hour).Format(time.RFC3339) || b.StartGap.Count != gaps || b.StartGap.LikelySeconds == nil || *b.StartGap.LikelySeconds != 10800 || b.StartGap.LimitSeconds == nil || *b.StartGap.LimitSeconds != 10800 {
				t.Fatalf("basis=%+v", b)
			}
			if b.DeliveredSpecs != 12 || b.StartToMerge.Count != 12 || b.StartToMerge.LikelySeconds == nil || *b.StartToMerge.LikelySeconds != 3600 || *b.StartToMerge.LimitSeconds != 6600 {
				t.Fatalf("basis=%+v", b)
			}
			for _, r := range doc.Queue {
				assertQueueCellShape(t, r.Start)
				assertQueueCellShape(t, r.End)
			}
			forecast := func(row int, start bool, likely, limit time.Time) {
				t.Helper()
				c := doc.Queue[row].End
				if start {
					c = doc.Queue[row].Start
				}
				want := followUpQueueCell{State: "forecast", Likely: followUpString(likely.UTC().Format(time.RFC3339)), Limit: followUpString(limit.UTC().Format(time.RFC3339))}
				if !reflect.DeepEqual(c, want) {
					t.Fatalf("row %d start=%v cell=%+v want=%+v", row, start, c, want)
				}
				text := "~" + likely.Format("01-02 15:04") + ", by " + limit.Format("01-02 15:04")
				if !strings.Contains(out, text) {
					t.Fatalf("missing %s in %s", text, out)
				}
			}
			condition := func(row int, start bool, want string, durations bool) {
				t.Helper()
				c := doc.Queue[row].End
				if start {
					c = doc.Queue[row].Start
				}
				if c.State != "condition" || c.Condition != want || c.Likely != nil || c.Limit != nil {
					t.Fatalf("condition=%+v want=%s", c, want)
				}
				text := want
				if durations {
					if c.LikelyAfterSeconds == nil || *c.LikelyAfterSeconds != 3600 || c.LimitAfterSeconds == nil || *c.LimitAfterSeconds != 6600 {
						t.Fatalf("durations=%+v", c)
					}
					text = "+1h00m, by +1h50m after its start"
				} else if c.LikelyAfterSeconds != nil || c.LimitAfterSeconds != nil {
					t.Fatalf("unexpected durations=%+v", c)
				}
				if !strings.Contains(out, text) {
					t.Fatalf("missing %s in %s", text, out)
				}
			}
			switch name {
			case "in_progress":
				forecast(0, false, S.Add(time.Hour), S.Add(110*time.Minute))
			case "longer_than_the_time_passed":
				forecast(0, false, S.Add(100*time.Minute), S.Add(120*time.Minute))
			case "past_the_usual_time":
				if doc.Queue[0].End.State != "none" || doc.Queue[0].End.Reason != "past the usual time" || !strings.Contains(out, "past the usual time") {
					t.Fatalf("row=%+v text=%s", doc.Queue[0], out)
				}
			case "one_two_and_three_behind":
				for i := 1; i <= 3; i++ {
					forecast(i, true, S.Add(time.Duration(i)*3*time.Hour), S.Add(time.Duration(i)*3*time.Hour))
					forecast(i, false, S.Add(time.Duration(i)*3*time.Hour+time.Hour), S.Add(time.Duration(i)*3*time.Hour+110*time.Minute))
				}
			case "the_fourth_is_a_condition":
				condition(4, true, "after 0016 starts", false)
				condition(4, false, "", true)
			case "a_prerequisite_is_not_counted":
				condition(2, true, "after prerequisite 0001 is merged", false)
				condition(2, false, "", true)
				for _, pair := range [][2]int{{1, 1}, {3, 2}, {4, 3}} {
					forecast(pair[0], true, S.Add(time.Duration(pair[1])*3*time.Hour), S.Add(time.Duration(pair[1])*3*time.Hour))
				}
			case "after_the_spec_in_progress_is_merged":
				for i := 1; i <= 3; i++ {
					condition(i, true, "after 0013 is merged", false)
					condition(i, false, "", true)
				}
			case "no_spec_in_progress":
				condition(0, true, "next in the queue", false)
				condition(1, true, "after 0013 starts", false)
				condition(4, true, "after it is queued", false)
			case "a_run_slot":
				condition(1, true, "after a Run slot frees", false)
				condition(1, false, "", true)
			case "a_park":
				condition(1, true, "after its park is answered", false)
				condition(1, false, "", true)
			case "a_parked_spec_that_had_started":
				condition(0, false, "after its park is answered", false)
			case "a_spec_being_authored":
				condition(5, true, "after its plan is merged", false)
				condition(5, false, "", true)
			}
			if name != "no_spec_in_progress" {
				if doc.Queue[0].Start.State != "actual" || doc.Queue[0].Start.Actual == nil || *doc.Queue[0].Start.Actual != S.Format(time.RFC3339) || !strings.Contains(out, S.Format("01-02 15:04")) {
					t.Fatalf("actual=%+v", doc.Queue[0].Start)
				}
			}
		})
	}
}

func TestTheStatusViewSaysWhenTheHistoryIsTooSmall(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 12, 5, 0, 0, time.UTC)
	S := now.Add(-5 * time.Minute)
	f, in := statusQueueDeliveryHistory(t, now, S, 3, false)
	followUpStatusSeedQueue(t, f, []string{"0013-running"}, func(q *store.DeliveryQueue) { q.Items[0].Stage = store.DeliveryStageRunning })
	writeImplementSpecAtRoot(t, in.SpecRoot, "0014-planned", []implementSeed{{id: "task_01", status: "pending"}})
	doc := followUpStatusJSON(t)
	out, stderr, code := followUpRunCLI(t, "--status")
	if code != 0 || stderr != "" || !strings.Contains(out, "Forecast basis: not enough history, 3 of 11 delivered Specs.\n") {
		t.Fatalf("exit=%d stderr=%q out=%s", code, stderr, out)
	}
	for _, row := range doc.Queue {
		if row.End.State != "none" || row.End.Reason != "not enough history (3 of 11)" {
			t.Fatalf("row=%+v", row)
		}
		for _, c := range []followUpQueueCell{row.Start, row.End} {
			assertQueueCellShape(t, c)
			if c.Actual != nil || c.Likely != nil || c.Limit != nil || c.LikelyAfterSeconds != nil || c.LimitAfterSeconds != nil {
				t.Fatalf("time without basis=%+v", c)
			}
		}
	}
	if f.judge.calls != 0 {
		t.Fatal("judge called")
	}
}

type statusQueueForge func(context.Context, string, ...string) (string, error)

func (f statusQueueForge) RunGH(ctx context.Context, dir string, args ...string) (string, error) {
	return f(ctx, dir, args...)
}

func statusQueueWithGit(t *testing.T, runner preflight.GitRunner) (followUpStatusDocument, string, string, int) {
	t.Helper()
	env := commandEnvironmentForTest(t)
	var out, stderr bytes.Buffer
	loaded, err := loadCommandConfig(env, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	code := runFollowUpStatusWithHistoryGit(t.Context(), followUpRequest{status: true, format: "json"}, &out, &stderr, loaded, env, runner)
	var doc followUpStatusDocument
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v; exit=%d stderr=%s stdout=%s", err, code, stderr.String(), out.String())
	}
	return doc, out.String(), stderr.String(), code
}

func TestTheStatusViewReadsTheHistoryBesideTheForgeCall(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 12, 5, 0, 0, time.UTC)
	f, _ := statusQueueDeliveryHistory(t, now, now.Add(-5*time.Minute), 12, true)
	followUpStatusSeedQueue(t, f, []string{"0013-running"}, func(q *store.DeliveryQueue) { q.Items[0].Stage = store.DeliveryStageRunning })
	logStarted := make(chan struct{})
	forgeRead := make(chan struct{})
	updateCommandDependenciesForTest(t, func(d *commandDependencies) {
		original := d.followUpGHRunner
		d.followUpGHRunner = statusQueueForge(func(ctx context.Context, dir string, args ...string) (string, error) {
			select {
			case <-logStarted:
			case <-ctx.Done():
				return "", ctx.Err()
			}
			close(forgeRead)
			return original.RunGH(ctx, dir, args...)
		})
	})
	calls := 0
	runner := followUpReadGit(func(ctx context.Context, dir string, args ...string) (string, error) {
		if args[0] == "log" {
			calls++
			close(logStarted)
			select {
			case <-forgeRead:
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
		return (preflight.ExecGitRunner{}).RunGit(ctx, dir, args...)
	})
	doc, out, stderr, code := statusQueueWithGit(t, runner)
	if code != 0 || stderr != "" || calls != 1 || doc.ForecastBasis.DeliveredSpecs != 12 || len(doc.Queue) != 1 || doc.Queue[0].End.State != "forecast" || doc.Queue[0].End.Likely == nil {
		t.Fatalf("exit=%d stderr=%q calls=%d out=%s", code, stderr, calls, out)
	}
}

func TestTheStatusViewContinuesWithoutAnUnreadableSource(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 13, 36, 0, 0, time.UTC)
	for _, name := range []string{"phase_record", "delivery_history", "backlog"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f, in := statusQueueDeliveryHistory(t, now, now.Add(-5*time.Minute), 12, true)
			followUpStatusSeedQueue(t, f, []string{"0013-running"}, func(q *store.DeliveryQueue) { q.Items[0].Stage = store.DeliveryStageRunning })
			writeImplementSpecAtRoot(t, in.SpecRoot, "0020-planned", []implementSeed{{id: "task_01", status: "pending"}})
			mustMkdir(t, filepath.Join(in.SpecRoot, "0021-folder"))
			followUpQueuePhase(t, in, now.Add(-time.Hour), "prd", "0022-phase", f.repo)
			followUpWrite(t, f.repo, "docs/backlog/open.md", "---\nstatus: open\n---\n# Open entry\n")
			before := followUpStatusJSON(t)
			if before.ForecastBasis.DeliveredSpecs != 12 || before.Queue[0].End.State != "forecast" || len(before.Queue) != 5 {
				t.Fatalf("fixture=%+v", before)
			}
			runner := preflight.GitRunner(preflight.ExecGitRunner{})
			switch name {
			case "phase_record":
				// A file in place of the directory fails on every host, including root.
				phaseDir := filepath.Join(f.home, ".roundfix/phase-record")
				if err := os.Rename(phaseDir, phaseDir+"-saved"); err != nil {
					t.Fatal(err)
				}
				mustWrite(t, phaseDir, "unreadable directory")
			case "backlog":
				followUpWrite(t, f.repo, "docs/backlog/broken.md", "---\nbroken: [\n---\n")
			case "delivery_history":
				runner = followUpReadGit(func(ctx context.Context, dir string, args ...string) (string, error) {
					if args[0] == "log" {
						return "", errors.New("history unavailable\nprivate second line")
					}
					return (preflight.ExecGitRunner{}).RunGit(ctx, dir, args...)
				})
			}
			doc, out, stderr, code := statusQueueWithGit(t, runner)
			sourceName := strings.ReplaceAll(name, "_", "-")
			found := false
			for _, s := range doc.Sources {
				if s.Name == sourceName {
					found = true
					if s.State != "unavailable" || s.Detail == "" || strings.Contains(s.Detail, "\n") || strings.Contains(s.Detail, "private second") {
						t.Fatalf("source=%+v", s)
					}
				}
			}
			want := 5
			if name == "phase_record" || name == "backlog" {
				want = 4
			}
			if !found || code != 0 || stderr != "" || len(doc.Queue) != want || doc.Queue[0].Spec != "0013-running" || doc.Queue[1].Spec != "0020-planned" || doc.Queue[2].Spec != "0021-folder" {
				t.Fatalf("exit=%d stderr=%q report=%s", code, stderr, out)
			}
			if name == "phase_record" && (doc.Queue[2].Status != "authoring" || doc.Queue[2].Note != "" || doc.Queue[3].Kind != "backlog") {
				t.Fatalf("authoring without phase=%+v", doc.Queue)
			}
			if name == "delivery_history" {
				if doc.ForecastBasis.DeliveredSpecs != 0 || doc.Queue[0].Start.State != "actual" || doc.Queue[0].Start.Actual == nil || *doc.Queue[0].Start.Actual != now.Add(-5*time.Minute).Format(time.RFC3339) {
					t.Fatalf("history fallback=%+v", doc)
				}
			} else if doc.ForecastBasis.DeliveredSpecs != 12 || doc.Queue[0].End.State != "forecast" {
				t.Fatalf("unaffected history=%+v", doc)
			}
			for _, r := range doc.Queue {
				for _, c := range []followUpQueueCell{r.Start, r.End} {
					assertQueueCellShape(t, c)
					if name == "delivery_history" && (c.Likely != nil || c.Limit != nil || c.LikelyAfterSeconds != nil || c.LimitAfterSeconds != nil) {
						t.Fatalf("forecast without history=%+v", c)
					}
				}
			}
		})
	}
}

func statusQueueSeconds(n int64) *int64  { return &n }
func statusQueueAt(at time.Time) *string { return followUpString(at.UTC().Format(time.RFC3339)) }
func statusQueueForecast(likely, limit time.Time) followUpQueueCell {
	return followUpQueueCell{State: "forecast", Likely: statusQueueAt(likely), Limit: statusQueueAt(limit)}
}
func statusQueueAfter(likely, limit int64) followUpQueueCell {
	return followUpQueueCell{State: "condition", LikelyAfterSeconds: statusQueueSeconds(likely), LimitAfterSeconds: statusQueueSeconds(limit)}
}

func TestTheSpecQueueTableRendersTimesAndDurations(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 13, 36, 0, 0, time.UTC)
	zone := time.FixedZone("-03", -3*60*60)
	base := func() followUpStatusDocument {
		return followUpStatusDocument{now: now, zone: zone, ForecastBasis: followUpForecastBasis{Minimum: 11, DeliveredSpecs: 109, FirstEnd: statusQueueAt(now.AddDate(0, 0, -30).Add(86*time.Minute + 53*time.Second)), LastEnd: statusQueueAt(now.AddDate(0, 0, -1).Add(-39*time.Minute - 8*time.Second)), LikelyPercentile: 50, LimitPercentile: 85, StartToMerge: followUpForecastMeasure{Count: 109, LikelySeconds: statusQueueSeconds(6990), LimitSeconds: statusQueueSeconds(15711)}, StartGap: followUpForecastMeasure{Count: 74, LikelySeconds: statusQueueSeconds(7365), LimitSeconds: statusQueueSeconds(13963)}}}
	}
	for _, name := range []string{"the_recorded_moment", "not_enough_history", "only_the_start_gaps_are_too_few", "only_backlog_rows_print_no_basis_line", "no_row_prints_only_the_count", "ten_rows_align", "a_likely_value_is_cut_and_a_limit_is_rounded_up"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			doc := base()
			switch name {
			case "the_recorded_moment":
				doc.Queue = []followUpQueueRow{
					{Row: 1, Kind: "delivery_queue", Spec: "0397-running", Status: "running 0/5 Tasks", TasksTotal: 5, Start: followUpQueueCell{State: "actual", Actual: statusQueueAt(now.Add(-126*time.Minute + 44*time.Second))}, End: statusQueueForecast(now.Add(83*time.Minute+9*time.Second), now.Add(216*time.Minute+44*time.Second))},
					{Row: 2, Kind: "planned", Spec: "0398-planned", Status: "not queued 0/4 Tasks", TasksTotal: 4, Start: statusQueueForecast(now.Add(44*time.Minute+15*time.Second), now.Add(266*time.Minute+58*time.Second)), End: statusQueueForecast(now.Add(211*time.Minute+24*time.Second), now.Add(544*time.Minute+19*time.Second))},
					{Row: 3, Kind: "authoring", Spec: "0399-authoring", Status: "authoring: prd", Start: followUpQueueCell{State: "condition", Condition: "after its plan is merged"}, End: statusQueueAfter(6990, 15711), Note: "prd since " + now.Add(-41*time.Minute).In(zone).Format("01-02 15:04")},
					{Row: 4, Kind: "backlog", Status: "not yet a Spec", Start: followUpQueueCell{State: "none"}, End: followUpQueueCell{State: "none"}, Note: "A report of late items; waits on the maintainer since " + now.AddDate(0, 0, -5).Format("2006-01-02")},
				}
				want := statusQueueRecordedTranscript
				if got := renderFollowUpSpecQueue(doc); got != want {
					t.Fatalf("got:\n%s\nwant:\n%s", got, want)
				}
			case "not_enough_history":
				doc.ForecastBasis = followUpForecastBasis{DeliveredSpecs: 3, Minimum: 11, LikelyPercentile: 50, LimitPercentile: 85, StartToMerge: followUpForecastMeasure{Count: 3}, StartGap: followUpForecastMeasure{Count: 2}}
				doc.Queue = []followUpQueueRow{{Row: 1, Kind: "planned", Spec: "0007-planned", Status: "not queued 0/2 Tasks", TasksTotal: 2, Start: followUpQueueCell{State: "condition", Condition: "after it is queued"}, End: followUpQueueCell{State: "none", Reason: "not enough history (3 of 11)"}}, {Row: 2, Kind: "backlog", Status: "not yet a Spec", Start: followUpQueueCell{State: "none"}, End: followUpQueueCell{State: "none"}, Note: "Export the ledger as a spreadsheet"}}
				if got := renderFollowUpSpecQueue(doc); got != statusQueueSmallTranscript {
					t.Fatalf("got:\n%s\nwant:\n%s", got, statusQueueSmallTranscript)
				}
			case "only_the_start_gaps_are_too_few":
				doc.ForecastBasis.StartGap = followUpForecastMeasure{Count: 10}
				doc.Queue = []followUpQueueRow{{Row: 1, Kind: "planned", Spec: "0001-planned", Status: "not queued 0/0 Tasks", Start: followUpQueueCell{State: "condition", Condition: "next in the queue"}, End: statusQueueAfter(6990, 15711)}}
				got := renderFollowUpSpecQueue(doc)
				want := "Forecast basis: 109 Specs delivered from 10-17 to 11-15 on this machine. Likely is the median and the limit the 85th percentile: 1h56m and 4h22m from start to merge.\n"
				if !strings.HasSuffix(got, want) || strings.Contains(got, "from one start") {
					t.Fatalf("table=%s", got)
				}
			case "only_backlog_rows_print_no_basis_line":
				doc.Queue = []followUpQueueRow{{Row: 1, Kind: "backlog", Status: "not yet a Spec", Start: followUpQueueCell{State: "none"}, End: followUpQueueCell{State: "none"}, Note: "Open entry"}}
				got := renderFollowUpSpecQueue(doc)
				if strings.Contains(got, "Forecast basis:") || !strings.HasPrefix(got, "Spec queue: 1\n") || !strings.Contains(got, "Open entry\n") {
					t.Fatalf("table=%s", got)
				}
			case "no_row_prints_only_the_count":
				if got := renderFollowUpSpecQueue(doc); got != "Spec queue: 0\n" {
					t.Fatalf("table=%q", got)
				}
			case "ten_rows_align":
				for i := 1; i <= 10; i++ {
					doc.Queue = append(doc.Queue, followUpQueueRow{Row: i, Kind: "planned", Spec: fmt.Sprintf("%04d-planned", i), Status: fmt.Sprintf("not queued %d/10 Tasks", i), Start: followUpQueueCell{State: "condition", Condition: "after it is queued"}, End: statusQueueAfter(6990, 15711), Note: "-"})
				}
				lines := strings.Split(renderFollowUpSpecQueue(doc), "\n")
				if len(lines) != 14 {
					t.Fatalf("lines=%q", lines)
				}
				header := lines[1]
				for i, line := range lines[2:12] {
					for _, cell := range []string{fmt.Sprintf("%04d", i+1), fmt.Sprintf("not queued %d/10 Tasks", i+1), "after it is queued", "+1h56m, by +4h22m after its start"} {
						title := "Spec"
						switch {
						case strings.HasPrefix(cell, "not queued"):
							title = "Status"
						case strings.HasPrefix(cell, "after"):
							title = "Start (-03)"
						case strings.HasPrefix(cell, "+"):
							title = "End (-03)"
						}
						if strings.Index(line, cell) != strings.Index(header, title) {
							t.Fatalf("misaligned line %q header %q", line, header)
						}
					}
					if !strings.HasSuffix(line, "  -") || strings.HasSuffix(line, " ") {
						t.Fatalf("last column padded %q", line)
					}
				}
			case "a_likely_value_is_cut_and_a_limit_is_rounded_up":
				for _, tc := range []struct {
					likely, limit int64
					want          string
				}{{6990, 15711, "+1h56m, by +4h22m after its start"}, {3601, 3601, "+1h00m, by +1h01m after its start"}, {2700, 2700, "+45m, by +45m after its start"}} {
					doc.Queue = []followUpQueueRow{{Row: 1, Kind: "planned", Spec: "0001-planned", Status: "not queued 0/0 Tasks", Start: statusQueueForecast(now.Add(time.Second), now.Add(time.Second)), End: statusQueueAfter(tc.likely, tc.limit)}}
					got := renderFollowUpSpecQueue(doc)
					if !strings.Contains(got, tc.want) || !strings.Contains(got, "~11-16 10:36, by 11-16 10:37") {
						t.Fatalf("table=%s want=%s", got, tc.want)
					}
				}
			}
		})
	}
}

const statusQueueRecordedTranscript = "Spec queue: 4\n  #  Spec  Status                Start (-03)                   End (-03)                          Note\n  1  0397  running 0/5 Tasks     11-16 08:30                   ~11-16 11:59, by 11-16 14:13       -\n  2  0398  not queued 0/4 Tasks  ~11-16 11:20, by 11-16 15:03  ~11-16 14:07, by 11-16 19:41       -\n  3  0399  authoring: prd        after its plan is merged      +1h56m, by +4h22m after its start  prd since 11-16 09:55\n  4  -     not yet a Spec        -                             -                                  A report of late items; waits on the maintainer since 2026-11-11\nForecast basis: 109 Specs delivered from 10-17 to 11-15 on this machine. Likely is the median and the limit the 85th percentile: 1h56m and 4h22m from start to merge, 2h02m and 3h53m from one start to the next.\n"

const statusQueueSmallTranscript = "Spec queue: 2\n  #  Spec  Status                Start (-03)         End (-03)                     Note\n  1  0007  not queued 0/2 Tasks  after it is queued  not enough history (3 of 11)  -\n  2  -     not yet a Spec        -                   -                             Export the ledger as a spreadsheet\nForecast basis: not enough history, 3 of 11 delivered Specs.\n"

const statusQueueSurface1 = "Status for <namespace> in <root>\nIn progress: 0\nPending Specs: 1\n  - 0002-export-the-ledger not queued 0/0 Tasks\nSpec queue: 3\n  #  Spec  Status                Start (UTC)               End (UTC)                     Note\n  1  0002  not queued 0/0 Tasks  after it is queued        not enough history (0 of 11)  -\n  2  0003  authoring             after its plan is merged  not enough history (0 of 11)  -\n  3  -     not yet a Spec        -                         -                             A table of late items; waits on the maintainer since 2026-11-18\nForecast basis: not enough history, 0 of 11 delivered Specs.\nNext release: unavailable: <reason>\nBlockers: 0\nWaiting on the maintainer: 0\nOpen authoring questions: wait (default)\nDelegated decisions: 0\nSources:\n  runs: checked 0\n  delivery-queue: checked 0\n  pull-requests: <state>\n  git: checked 0\n  specs: <state>\n  release-plan: unavailable: <reason>\n  backlog: checked 1\n  phase-record: checked 0\n  delivery-history: checked 0\n"

const statusQueueSurface2 = "Status for <namespace> in <root>\nIn progress: 0\nPending Specs: 0\nSpec queue: 0\nNext release: unavailable: <reason>\nBlockers: 0\nWaiting on the maintainer: 0\nOpen authoring questions: wait (default)\nDelegated decisions: 0\nSources:\n  runs: checked 0\n  delivery-queue: checked 0\n  pull-requests: <state>\n  git: checked 0\n  specs: checked 0\n  release-plan: unavailable: <reason>\n  backlog: checked 0\n  phase-record: checked 0\n  delivery-history: checked 0\n"

func TestTheStatusViewSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 18, 13, 36, 0, 0, time.UTC)
	for i, wantTemplate := range []string{statusQueueSurface1, statusQueueSurface2} {
		t.Run(fmt.Sprintf("transcript_%d", i+1), func(t *testing.T) {
			t.Parallel()
			f, in := newStatusQueueFixture(t, now)
			gitImplement(t, f.repo, "remote", "remove", "origin")
			if i == 0 {
				followUpWrite(t, in.SpecRoot, "0002-export-the-ledger/_prd.md", "---\nstatus: active\n---\n# Export the ledger\n")
				followUpWrite(t, in.SpecRoot, "0003-a-report-of-late-items/_idea.md", "# A report of late items\n")
				followUpWrite(t, f.repo, "docs/backlog/"+now.Format("2006-01-02")+"-a-table-of-late-items.md", "---\ntype: feat\nstatus: open\ncreated: "+now.Format("2006-01-02")+"\nwaiting_on: the maintainer\nsince: "+now.Format("2006-01-02")+"\n---\n# A table of late items\n")
			}
			// The transcript specifies one commit. Replace the seeded tree in that same
			// commit, at an instant derived from now, instead of adding history.
			gitImplement(t, f.repo, "add", ".")
			statusQueueAmendAt(t, f.repo, now.Add(-time.Hour))
			out, stderr, code := followUpRunCLI(t, "--status")
			doc := followUpStatusJSON(t)
			state := func(source followUpSource) string {
				if source.State == "checked" {
					return fmt.Sprintf("checked %d", source.Count)
				}
				if source.Detail != "" {
					return source.State + ": " + source.Detail
				}
				return source.State
			}
			want := strings.NewReplacer("<namespace>", filepath.Base(f.repo), "<root>", f.repo, "<reason>", *doc.Release.Reason).Replace(wantTemplate)
			want = strings.Replace(want, "pull-requests: <state>", "pull-requests: "+state(doc.Sources[2]), 1)
			want = strings.Replace(want, "specs: <state>", "specs: "+state(doc.Sources[4]), 1)
			if code != 0 || stderr != "" || out != want {
				t.Fatalf("exit=%d stderr=%q\ngot:\n%s\nwant:\n%s", code, stderr, out, want)
			}
			if count := strings.TrimSpace(gitImplementOutput(t, f.repo, "rev-list", "--count", "main")); count != "1" {
				t.Fatalf("commit count=%s", count)
			}
			if _, err := os.Stat(store.DatabasePath(f.home)); !os.IsNotExist(err) {
				t.Fatalf("Run Database: %v", err)
			}
			if _, err := os.Stat(filepath.Join(f.home, ".roundfix/phase-record")); !os.IsNotExist(err) {
				t.Fatalf("Phase Record: %v", err)
			}
		})
	}
}

func statusQueueAmendAt(t *testing.T, repo string, at time.Time) {
	t.Helper()
	args := append(gitConfigArgsForTest(), "commit", "--amend", "--no-edit", "--date="+at.UTC().Format(time.RFC3339))
	cmd := exec.CommandContext(t.Context(), "git", args...)
	cmd.Dir = repo
	cmd.Env = append(isolatedGitEnvForTest(), "GIT_COMMITTER_DATE="+at.UTC().Format(time.RFC3339), "GIT_AUTHOR_DATE="+at.UTC().Format(time.RFC3339))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("amend: %v\n%s", err, out)
	}
}
