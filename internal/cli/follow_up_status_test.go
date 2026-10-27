// Suite: Follow-up Command status view.
// Invariant: status reports local work and release decisions without writes or judge calls.
// Boundary IN: public CLI dispatch, temporary Git repositories, Specs and SQLite.
// Boundary OUT: GitHub and judge HTTP use per-test fakes only.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestFollowUpHelpTranscript(t *testing.T) {
	t.Parallel()
	newFollowUpFixture(t)
	const want = `Usage:
  roundfix follow-up [--status] [--messages <path>] [--format <text|json>]

List this repository's pending work in Follow-up Priority order, or show its status.
It is read-only: it never writes a file, commits, pushes or prompts.
Options:
  --status    Show work in progress, pending Specs, the next release and blockers
  --messages  JSON file of ai-memory messages to list with the pending work
  --format    text (default) or json

Exits 0 for a report, including one with nothing to do; 2 for invalid usage,
configuration or messages file; 1 for an unreadable Run Database.
`
	out, stderr, code := followUpRunCLI(t, "--help")
	if code != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr, out, want)
	}
	var rootOut, rootErr bytes.Buffer
	code = runCLI(t, []string{"--help"}, &rootOut, &rootErr)
	if code != 0 || rootErr.Len() != 0 || !strings.Contains(rootOut.String(), "  roundfix follow-up [--status] [--messages <path>] [--format <text|json>]\n") {
		t.Fatalf("root help exit=%d stderr=%q stdout=%q", code, rootErr.String(), rootOut.String())
	}
}

func followUpStatusJSON(t *testing.T) followUpStatusDocument {
	t.Helper()
	out, stderr, code := followUpRunCLI(t, "--status", "--format", "json")
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
	}
	var doc followUpStatusDocument
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Schema != "roundfix/follow-up/v1" || doc.View != "status" {
		t.Fatalf("document=%+v", doc)
	}
	names := []string{}
	for _, source := range doc.Sources {
		names = append(names, source.Name)
	}
	if !reflect.DeepEqual(names, []string{"runs", "delivery-queue", "pull-requests", "git", "specs", "release-plan"}) {
		t.Fatalf("sources=%v", names)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &fields); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"judge", "items"} {
		if _, found := fields[field]; found {
			t.Fatalf("status contains list field %s", field)
		}
	}
	return doc
}

func followUpStatusSeedQueue(t *testing.T, f followUpFixture, slugs []string, update func(*store.DeliveryQueue)) store.DeliveryQueue {
	t.Helper()
	writer, err := store.Open(context.Background(), f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}()
	queue, err := writer.CreateDeliveryQueue(context.Background(), f.repo, slugs)
	if err != nil {
		t.Fatal(err)
	}
	update(&queue)
	for _, item := range queue.Items {
		if err := writer.UpdateDeliveryQueueItem(context.Background(), f.repo, item); err != nil {
			t.Fatal(err)
		}
	}
	return queue
}

func TestFollowUpStatusListsPendingSpecsInDeliveryOrder(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	for _, slug := range []string{"0003-first", "0009-second", "0005-unqueued", "0010-merged"} {
		writeImplementSpecAtRoot(t, filepath.Join(f.repo, "docs/specs"), slug, []implementSeed{
			{id: "task_01", status: string(spec.StatusCompleted)},
			{id: "task_02", status: string(spec.StatusPending)},
		})
	}
	queue := followUpStatusSeedQueue(t, f, []string{"0009-second", "0003-first", "0010-merged"}, func(queue *store.DeliveryQueue) {
		queue.Items[0].Stage = store.DeliveryStageRunning
		queue.Items[2].Stage = store.DeliveryStageMerged
	})
	doc := followUpStatusJSON(t)
	got := []string{}
	for _, row := range doc.PendingSpecs {
		got = append(got, row.Spec)
		if row.TasksCompleted != 1 || row.TasksTotal != 2 {
			t.Fatalf("Task counts=%+v", row)
		}
	}
	if !reflect.DeepEqual(got, []string{"0009-second", "0003-first", "0005-unqueued"}) {
		t.Fatalf("Specs=%v", got)
	}
	for i := range 2 {
		row := doc.PendingSpecs[i]
		if row.QueuePosition == nil || *row.QueuePosition != queue.Items[i].Position || row.Stage == nil || *row.Stage != string(queue.Items[i].Stage) {
			t.Fatalf("queued row=%+v want=%+v", row, queue.Items[i])
		}
	}
	if doc.PendingSpecs[2].QueuePosition != nil || doc.PendingSpecs[2].Stage != nil {
		t.Fatalf("unqueued=%+v", doc.PendingSpecs[2])
	}
	out, stderr, code := followUpRunCLI(t, "--status")
	want := fmt.Sprintf("Pending Specs: 3\n  %d 0009-second running 1/2 Tasks\n  %d 0003-first %s 1/2 Tasks\n  - 0005-unqueued not queued 1/2 Tasks\n", queue.Items[0].Position, queue.Items[1].Position, queue.Items[1].Stage)
	if code != 0 || stderr != "" || !strings.Contains(out, want) {
		t.Fatalf("exit=%d stderr=%q stdout=%s want=%s", code, stderr, out, want)
	}
}

func TestFollowUpStatusReportsTheNextRelease(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"HEAD at release", "feature", "dirty", "patch", "maintenance", "manual", "no tag"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			f := newFollowUpFixture(t)
			if mode != "no tag" {
				gitImplement(t, f.repo, "tag", "v0.4.1")
			}
			subject := ""
			switch mode {
			case "feature":
				subject = "feat: expose widget status"
			case "patch":
				subject = "fix: repair widget status"
			case "maintenance":
				subject = "docs: describe widget"
			case "manual":
				subject = "chore: change widget behavior"
			}
			if subject != "" {
				path := "internal/widget.go"
				if mode == "maintenance" {
					path = "docs/specs/0001-fixture/task_01.md"
				}
				followUpWrite(t, f.repo, path, "widget change\n")
				gitImplement(t, f.repo, "add", ".")
				gitImplement(t, f.repo, "commit", "-m", subject)
			}
			if mode == "dirty" {
				followUpWrite(t, f.repo, "README.md", "# Modified\n")
			}
			doc := followUpStatusJSON(t)
			out, stderr, code := followUpRunCLI(t, "--status")
			if code != 0 || stderr != "" {
				t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
			}
			switch mode {
			case "HEAD at release":
				if doc.Release.State != "no_unreleased_changes" || doc.State != "nothing_to_do" || !strings.Contains(out, "Next release: no unreleased changes\n") || doc.Release.Reason != nil || len(doc.Maintainer) != 0 {
					t.Fatalf("doc=%+v text=%s", doc, out)
				}
			case "dirty", "no tag":
				reason := "dirty"
				if mode == "no tag" {
					reason = "tag"
				}
				if doc.Release.State != "unavailable" || doc.Release.Reason == nil || !strings.Contains(*doc.Release.Reason, reason) || doc.State != "work_found" || !strings.Contains(out, "Next release: unavailable: "+*doc.Release.Reason+"\n") {
					t.Fatalf("doc=%+v text=%s", doc, out)
				}
			default:
				wantState, wantVersion := "approval_required", "v0.5.0"
				switch mode {
				case "patch":
					wantState, wantVersion = "ready", "v0.4.2"
				case "maintenance":
					wantState, wantVersion = "no_release", ""
				case "manual":
					wantState, wantVersion = "manual_classification_required", ""
				}
				sha := strings.TrimSpace(gitImplementOutput(t, f.repo, "rev-parse", "HEAD"))
				if doc.Release.State != wantState || doc.Release.Base != "v0.4.1" || doc.Release.ProposedVersion != wantVersion || (doc.Release.Reason != nil) != (mode == "manual") || len(doc.Release.Changes) != 1 || doc.Release.Changes[0].Subject != subject || doc.Release.Changes[0].Commit != sha || !strings.Contains(out, " "+subject+"\n") {
					t.Fatalf("release=%+v text=%s", doc.Release, out)
				}
				switch mode {
				case "feature":
					if doc.State != "work_found" || len(doc.Maintainer) != 1 || doc.Maintainer[0].Kind != "release_approval" || doc.Maintainer[0].Detail != "Approve the minor increment to v0.5.0?" || !strings.Contains(out, "Next release: v0.5.0 from v0.4.1 (approval_required)\n") || !strings.Contains(out, "  release v0.5.0: Approve the minor increment to v0.5.0?\n") {
						t.Fatalf("doc=%+v text=%s", doc, out)
					}
				case "manual":
					if doc.State != "work_found" || len(doc.Maintainer) != 1 || doc.Maintainer[0].Kind != "release_classification" || !strings.Contains(out, "manual release classification") {
						t.Fatalf("doc=%+v text=%s", doc, out)
					}
				case "maintenance":
					if doc.State != "nothing_to_do" || len(doc.Maintainer) != 0 || !strings.Contains(out, "Next release: none (no_release)\n") {
						t.Fatalf("doc=%+v text=%s", doc, out)
					}
				case "patch":
					if doc.State != "work_found" || len(doc.Maintainer) != 0 {
						t.Fatalf("doc=%+v", doc)
					}
				}
			}
		})
	}
}

func TestFollowUpStatusNamesTheManualClassification(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	gitImplement(t, f.repo, "tag", "v0.4.1")
	for i, subject := range []string{"feat: expose widget status", "chore: change widget behavior", "chore: change widget storage"} {
		followUpWrite(t, f.repo, fmt.Sprintf("internal/widget%d.go", i), "widget change\n")
		gitImplement(t, f.repo, "add", ".")
		gitImplement(t, f.repo, "commit", "-m", subject)
	}

	doc := followUpStatusJSON(t)
	const reason = "2 commit(s) since v0.4.1 need a manual release classification"
	const detail = "classify them with roundfix release plan --from v0.4.1 --to HEAD --impact <none|patch|minor|major> --reason <text>"
	if doc.Release.State != "manual_classification_required" || doc.Release.Base != "v0.4.1" || doc.Release.ProposedVersion != "" || doc.Release.Reason == nil || *doc.Release.Reason != reason || len(doc.Release.Changes) != 3 {
		t.Errorf("release=%+v, want manual classification with reason %q and three changes", doc.Release, reason)
	}
	wantMaintainer := []followUpMaintainer{{Kind: "release_classification", Detail: detail}}
	if !reflect.DeepEqual(doc.Maintainer, wantMaintainer) {
		t.Errorf("maintainer=%+v, want %+v", doc.Maintainer, wantMaintainer)
	}
	out, stderr, code := followUpRunCLI(t, "--status")
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
	}
	for _, want := range []string{
		"\nNext release: not proposed (manual_classification_required): " + reason + "\n",
		"\nWaiting on the maintainer: 1\n  release classification: " + detail + "\nSources:\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout=%s, want %q", out, want)
		}
	}
	for _, blank := range []string{"Next release:  ", "  release : "} {
		if strings.Contains(out, blank) {
			t.Errorf("stdout contains blank version %q: %s", blank, out)
		}
	}
}

func TestFollowUpStatusListsBlockersAndTheMaintainersQuestion(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	queue := followUpStatusSeedQueue(t, f, []string{"0009-second", "0003-first"}, func(queue *store.DeliveryQueue) {
		queue.Items[0].Stage = store.DeliveryStageParked
		queue.Items[0].Blocker = delivery.BlockerChecksTimeout
		queue.Items[1].Stage = store.DeliveryStageParked
		queue.Items[1].Blocker = delivery.BlockerReviewFindings
	})
	doc := followUpStatusJSON(t)
	if len(doc.Blockers) != 2 || len(doc.Maintainer) != 1 || doc.State != "work_found" {
		t.Fatalf("doc=%+v", doc)
	}
	out, stderr, code := followUpRunCLI(t, "--status")
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
	}
	for i, item := range queue.Items {
		park := delivery.ClassifyPark(queue, item)
		want := followUpBlocker{Spec: item.SpecSlug, Blocker: item.Blocker, ParkClass: string(park.Class), Next: park.Next}
		if doc.Blockers[i] != want || !strings.Contains(out, fmt.Sprintf("  %s %s: %s\n", item.SpecSlug, park.Class, park.Next)) {
			t.Fatalf("blocker=%+v want=%+v text=%s", doc.Blockers[i], want, out)
		}
	}
	question, found := delivery.PendingQuestionFor(queue)
	want := followUpMaintainer{Kind: "pending_question", Spec: "0009-second", Detail: question.Blocker, Answer: question.Answer}
	if !found || doc.Maintainer[0] != want || !strings.Contains(out, fmt.Sprintf("  question %s parked %s: %s\n", question.SpecSlug, question.Blocker, question.Answer)) {
		t.Fatalf("maintainer=%+v text=%s", doc.Maintainer, out)
	}
	// Compare against the existing public Delivery Queue status answer.
	var deliverOut, deliverErr bytes.Buffer
	code = runCLI(t, []string{"deliver", "status"}, &deliverOut, &deliverErr)
	if code != 0 || deliverErr.Len() != 0 || !strings.Contains(deliverOut.String(), "Answer: "+question.Answer+"\n") {
		t.Fatalf("deliver status exit=%d stderr=%s stdout=%s", code, deliverErr.String(), deliverOut.String())
	}
}

func TestFollowUpStatusNeverAsksTheJudge(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
	followUpSeedPending(t, f)
	// If status traversed repository notes, this malformed entry would mark them unavailable.
	followUpWrite(t, f.repo, "docs/_inbox/unreadable.md", "---\nbroken: [\n---\n")
	doc := followUpStatusJSON(t)
	out, stderr, code := followUpRunCLI(t, "--status")
	if f.judge.calls != 0 || code != 0 || stderr != "" || strings.Contains(out, "Judge:") {
		t.Fatalf("calls=%d exit=%d stderr=%s stdout=%s", f.judge.calls, code, stderr, out)
	}
	for _, source := range doc.Sources {
		if source.Name == "repository-inbox" || source.Name == "ai-memory-messages" {
			t.Fatalf("forbidden source=%+v", source)
		}
	}
}

func TestFollowUpStatusWritesNothing(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"no database", "at tag", "release", "dirty", "live WAL"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			f := newFollowUpFixture(t)
			gitImplement(t, f.repo, "tag", "v0.4.1")
			if mode == "release" {
				followUpWrite(t, f.repo, "internal/widget.go", "widget\n")
				gitImplement(t, f.repo, "add", ".")
				gitImplement(t, f.repo, "commit", "-m", "feat: add widget")
			}
			if mode == "dirty" {
				followUpWrite(t, f.repo, "README.md", "# Dirty\n")
			}
			if mode != "no database" {
				followUpSeedDatabase(t, f, true, true)
			}
			if mode == "live WAL" {
				writer, err := store.Open(context.Background(), f.home)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := writer.Close(); err != nil {
						t.Error(err)
					}
				}()
				if _, err := writer.CreateRun(context.Background(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.repo, LocalBranch: "main", SpecSlug: "0003-live"}); err != nil {
					t.Fatal(err)
				}
			}
			before := followUpSnapshot(t, f.repo, f.home, f.brain)
			status := gitImplementOutput(t, f.repo, "--no-optional-locks", "status", "--porcelain")
			for _, format := range []string{"text", "json"} {
				out, stderr, code := followUpRunCLI(t, "--status", "--format", format)
				if code != 0 || stderr != "" || out == "" {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
				}
				if after := followUpSnapshot(t, f.repo, f.home, f.brain); !reflect.DeepEqual(before, after) {
					t.Fatalf("%s report changed input bytes", format)
				}
				if after := gitImplementOutput(t, f.repo, "--no-optional-locks", "status", "--porcelain"); after != status {
					t.Fatalf("git status changed: before=%q after=%q", status, after)
				}
			}
		})
	}
}

func TestFollowUpStatusRefusesMessages(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"--status", "--messages", "messages.json"}, {"--messages=messages.json", "--status"}, {"--status", "--messages="}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			t.Parallel()
			f := newFollowUpFixture(t)
			// This configuration would fail if validation performed any reads first.
			followUpWrite(t, f.home, ".roundfix/config.yml", "invalid: [\n")
			if _, err := os.Stat(filepath.Join(f.repo, "messages.json")); !os.IsNotExist(err) {
				t.Fatalf("messages fixture exists: %v", err)
			}
			out, stderr, code := followUpRunCLI(t, args...)
			const want = "roundfix: follow-up failed: --messages applies only to the follow-up list; drop it with --status\nRun 'roundfix follow-up --help' for usage.\n"
			if code != 2 || out != "" || stderr != want {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, stderr)
			}
		})
	}
}
