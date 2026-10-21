// Boundary: rendered preflight stderr, compiled item binary, real Run Database and status.
package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
	"roundfix/internal/testfixture"
)

func TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal(t *testing.T) {
	t.Parallel()
	ceiling := store.ActiveImplementRunCeilingError{MaxActive: 1}
	old := &store.Run{ID: "old"}
	for _, tc := range []struct {
		name, color   string
		reason        error
		code          int
		before, after *store.Run
		want          bool
	}{
		{name: "plain", reason: ceiling, code: exitPreflight, want: true},
		{name: "colored", color: "always", reason: ceiling, code: exitPreflight, want: true},
		{name: "existing unchanged Run", reason: ceiling, code: exitPreflight, before: old, after: old, want: true},
		{name: "working tree refusal", reason: store.ActiveRunError{Existing: store.Run{Kind: store.KindImplement}}, code: exitPreflight},
		{name: "new Run", reason: ceiling, code: exitPreflight, after: &store.Run{ID: "new"}},
		{name: "replaced Run", reason: ceiling, code: exitPreflight, before: old, after: &store.Run{ID: "new"}},
		{name: "exit one", reason: ceiling, code: exitRunFailed},
		{name: "later reason line", reason: fmt.Errorf("other refusal\n%s", ceiling), code: exitPreflight},
		{name: "prefix without colon", reason: errors.New(store.ActiveImplementRunCeilingReason + " later"), code: exitPreflight},
		{name: "blank reason lines", reason: fmt.Errorf("\n\n%s", ceiling), code: exitPreflight, want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stderr bytes.Buffer
			printPreflightFailure("implement", tc.reason, commandWriter(&stderr, commandEnvironment{colorMode: tc.color}))
			if tc.color != "" && !strings.Contains(stderr.String(), "\x1b[") {
				t.Fatal("fixture has no ANSI color")
			}
			command := roundfixCommandResult{exitCode: tc.code, stderr: stderr.String()}
			_, err := deliveryRunResult(command, "", tc.before, tc.after)
			var refusal *delivery.RunSlotRefusal
			if errors.As(err, &refusal) != tc.want {
				t.Fatalf("error=%v, want refusal=%v", err, tc.want)
			}
			if err == nil || err.Error() != command.failure("roundfix implement").Error() {
				t.Fatalf("failure text changed: %v", err)
			}
		})
	}
}

func TestRunSpecReturnsARunSlotRefusalWithItsHolders(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	printPreflightFailure("implement", store.ActiveImplementRunCeilingError{MaxActive: 1}, &stderr)
	binary := testfixture.FixtureBinary(t, "slot-refusal", fmt.Sprintf(`package main
import ("fmt"; "os")
func main() {
 if len(os.Args)==3 && os.Args[1]=="migrate" && os.Args[2]=="--check" { return }
 if len(os.Args)==4 && os.Args[1]=="implement" && os.Args[2]=="--spec" && os.Args[3]=="widget" { fmt.Fprint(os.Stderr, %q); os.Exit(2) }
 os.Exit(7)
}`, stderr.String()))
	home, repo := t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.PersistIdentity(t, repo)
	if err := os.WriteFile(filepath.Join(repo, ".gitignore"), []byte("bin/\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "fixture")
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	holderRepo := t.TempDir()
	gittest.InitRepo(t, holderRepo, "-b", "main")
	holder, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: holderRepo, LocalBranch: "feat/holder", SpecSlug: "0001-holder", OwnerPID: os.Getpid()})
	if err != nil {
		t.Fatal(err)
	}
	workflow := &commandDeliveryWorkflow{store: db, git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{GitRoot: repo, HomeDir: home}}
	workflow.loaded.Config.Runs.MaxActive = 1
	workflow.loaded.Config.Defaults.ArtifactDir = filepath.Join(home, "artifacts")
	workflow.loaded.Config.Delivery.ItemBinary = roundconfig.ItemBinaryDeclaration{Path: "bin/roundfix", Build: deliveryItemBuildCommand(binary)}
	_, err = workflow.RunSpec(commandContextForTest(t, t.Context()), repo, "widget")
	var refusal *delivery.RunSlotRefusal
	if !errors.As(err, &refusal) {
		t.Fatalf("error=%v, want RunSlotRefusal", err)
	}
	wantError := (roundfixCommandResult{exitCode: exitPreflight, stderr: stderr.String()}).failure("roundfix implement").Error()
	wantSlots := delivery.RunSlots{MaxActive: 1, Holders: []delivery.RunSlotHolder{{RunID: holder.ID, Repository: holder.GitRoot, SpecSlug: holder.SpecSlug}}}
	if refusal.Err.Error() != wantError || !reflect.DeepEqual(refusal.Slots, wantSlots) {
		t.Fatalf("refusal=%+v slots=%+v, want %s %+v", refusal, refusal.Slots, wantError, wantSlots)
	}
	// A closed database must report a slot-read failure.
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := workflow.RunSlots(t.Context()); err == nil {
		t.Fatal("closed database yielded slots")
	}
}

func TestDeliverStatusShowsARunSlotWaitTranscript(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	queue, err := db.CreateDeliveryQueue(t.Context(), repo, []string{"0291-example"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageRunning
	item.Worktree = "<worktree>"
	item.Branch = "roundfix/deliver-0291-example"
	if _, _, _, err := db.RecordDeliveryQueueItemWorktree(t.Context(), repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
		t.Fatal(err)
	}
	item.Blocker = delivery.RunSlotWaitText(time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC), delivery.RunSlots{MaxActive: 1, Holders: []delivery.RunSlotHolder{{RunID: "run_holder", Repository: "<holder repository>", SpecSlug: "0001-holder"}}})
	if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0291-example\trunning\trun-slot-wait: until 2026-10-21T00:00:00Z; 1 Run(s) hold 1 slot(s): run_holder in <holder repository> for 0001-holder\t<worktree>\nLimits: deadline none, retries per item none, concurrency 1, tokens none\nUsage: no Runs recorded\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr.String(), stdout.String(), want)
	}
}
