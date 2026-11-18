package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/correction"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func TestDeliverStatusShowsUndeclaredCorrectionSigns(t *testing.T) {
	t.Parallel()
	for _, setting := range []string{"true", "false", ""} {
		name := setting
		if name == "" {
			name = "absent"
		}
		t.Run(name, func(t *testing.T) {
			f := newCorrectionSignsFixture(t, setting)
			head := correctionSignsCommit(t, f.repo, "fix: correction before merge")
			// A dirty owner setting cannot override the default branch's requirement.
			mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  require_defect_class: false\n")
			db, err := store.Open(t.Context(), f.home)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			queue, err := db.CreateDeliveryQueue(t.Context(), f.repo, []string{implementTestSlug})
			if err != nil {
				t.Fatal(err)
			}
			item := queue.Items[0]
			item.Branch = f.branch
			item.Worktree = f.repo
			item.Stage = store.DeliveryStageMerging
			item.CandidateCommits = []string{head}
			if _, _, _, err := db.RecordDeliveryQueueItemWorktree(t.Context(), f.repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
				t.Fatal(err)
			}
			if err := db.UpdateDeliveryQueueItem(t.Context(), f.repo, item); err != nil {
				t.Fatal(err)
			}
			call := func() string {
				t.Helper()
				var out, diag bytes.Buffer
				code := runCLI(t, []string{"deliver", "status"}, &out, &diag)
				if code != 0 || diag.Len() != 0 {
					t.Fatalf("exit=%d out=%q err=%q", code, out.String(), diag.String())
				}
				return out.String()
			}
			base := implementTestSlug + "\tmerging\t-\t" + f.repo + "\n" +
				"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\n" +
				"Usage: no Runs recorded\n"
			got := call()
			if setting != "true" {
				if got != base {
					t.Fatalf("default-off status=%q want Delivery Base %q", got, base)
				}
				return
			}
			line := "Corrections: " + implementTestSlug + " undeclared: commit " + head[:7] + ` "fix: correction before merge"; declare with roundfix corrections declare --spec ` + implementTestSlug + "\n"
			expected := strings.Replace(base, "Limits:", line+"Limits:", 1)
			if got != expected {
				t.Fatalf("status=%q want %q", got, expected)
			}
			path, err := f.workflow.correctionRecordPath(t.Context(), f.repo, implementTestSlug)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("status wrote the record: %v", err)
			}
			// An unreadable item's record is omitted without failing status.
			mustMkdir(t, path)
			if got := call(); got != base {
				t.Fatalf("unreadable record status=%q", got)
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			var out, diag bytes.Buffer
			code := runCLI(t, []string{"corrections", "declare", "--spec", implementTestSlug, "--class", "implementation", "--reason", "fixture fix", "--commit", head}, &out, &diag)
			if code != 0 {
				t.Fatalf("declare exit=%d out=%q err=%q", code, out.String(), diag.String())
			}
			if got := call(); got != base {
				t.Fatalf("covered status=%q", got)
			}
			record, err := correction.Read(path)
			if err != nil || len(record.Snapshots) != 0 {
				t.Fatalf("status appended signs: %+v err=%v", record, err)
			}
			gittest.Run(t, f.repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/unavailable")
			gittest.Run(t, f.repo, "update-ref", "-d", "refs/remotes/origin/main")
			if got := call(); got != base {
				t.Fatalf("unresolved default status=%q", got)
			}
		})
	}
}
