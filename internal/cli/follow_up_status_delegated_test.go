// Boundary IN: private homes, temporary local repositories and command output.
// Boundary OUT: no remote, provider, judge, or authoring session is started.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
)

const statusDelegatedEntry = `- delegated: 2026-11-20, prd
  - question: which failure behavior does an expired import have?
  - options: retry automatically up to 3 times; fail fast and notify the user; park it for manual review
  - adopted: retry automatically up to 3 times — it preserves the sync retry behavior
  - under: Project Config
  - answer: none yet
`

func statusDelegatedRow(slug string) followUpDelegatedDecision {
	return followUpDelegatedDecision{Spec: slug, Line: 6, Date: "2026-11-20", Phase: "prd", Question: "which failure behavior does an expired import have?", Options: "retry automatically up to 3 times; fail fast and notify the user; park it for manual review", Adopted: "retry automatically up to 3 times — it preserves the sync retry behavior", Under: "Project Config", Answer: "none yet"}
}

func statusDelegatedSpec(t *testing.T, root, slug, entry string) {
	t.Helper()
	writeImplementSpecAtRoot(t, root, slug, []implementSeed{{id: "task_01", status: "pending"}})
	followUpWrite(t, root, slug+"/_prd.md", "---\nstatus: active\n---\n## Decisions\n\n"+entry)
}

func TestTheStatusViewListsTheDelegatedDecisionsOfActiveSpecs(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"two_specs_in_slug_order", "an_answered_entry", "an_unreadable_entry_with_its_line", "a_spec_without_an_entry", "an_archived_spec_is_not_listed", "an_unreadable_prd_names_its_spec", "none_prints_zero"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			now := time.Date(2026, 11, 19, 20, 0, 0, 0, time.UTC)
			f, in := newStatusQueueFixture(t, now)
			gitImplement(t, f.repo, "remote", "remove", "origin")
			want := []followUpDelegatedDecision{}
			block := "Delegated decisions: 0\n"
			const slug = "0001-expired-imports"
			const line = "  0001-expired-imports 2026-11-20 prd: which failure behavior does an expired import have? -> retry automatically up to 3 times (answer: none yet)\n"
			switch name {
			case "two_specs_in_slug_order":
				statusDelegatedSpec(t, in.SpecRoot, "0009-later", statusDelegatedEntry)
				statusDelegatedSpec(t, in.SpecRoot, slug, statusDelegatedEntry+"\n"+statusDelegatedEntry)
				second := statusDelegatedRow(slug)
				second.Line = 13
				want = []followUpDelegatedDecision{statusDelegatedRow(slug), second, statusDelegatedRow("0009-later")}
				block = "Delegated decisions: 3\n" + line + line + strings.Replace(line, slug, "0009-later", 1)
			case "an_answered_entry":
				answer := `2026-11-22, changed — "Fail fast"`
				statusDelegatedSpec(t, in.SpecRoot, slug, strings.Replace(statusDelegatedEntry, "none yet", answer, 1))
				row := statusDelegatedRow(slug)
				row.Answer = answer
				want = append(want, row)
				block = "Delegated decisions: 1\n" + strings.Replace(line, "none yet", "2026-11-22, changed", 1)
			case "an_unreadable_entry_with_its_line":
				statusDelegatedSpec(t, in.SpecRoot, slug, strings.Replace(statusDelegatedEntry, "2026-11-20", "2026-02-30", 1))
				row := statusDelegatedRow(slug)
				row.Date = "2026-02-30"
				row.Problem = "the head needs a YYYY-MM-DD date and the phase prd, techspec or tasks"
				want = append(want, row)
				block = "Delegated decisions: 1\n  " + slug + " line 6: unreadable: " + row.Problem + "\n"
			case "a_spec_without_an_entry":
				statusDelegatedSpec(t, in.SpecRoot, slug, "Autonomous route: no person present; decisions taken from the request.\n")
				mustMkdir(t, filepath.Join(in.SpecRoot, "0002-no-prd"))
			case "an_archived_spec_is_not_listed":
				statusDelegatedSpec(t, in.SpecRoot, slug, statusDelegatedEntry)
				if err := os.RemoveAll(filepath.Join(in.SpecRoot, slug)); err != nil {
					t.Fatal(err)
				}
				followUpArchiveFixture(t, in, slug, spec.ArchiveNoQA)
			case "an_unreadable_prd_names_its_spec":
				mustMkdir(t, filepath.Join(in.SpecRoot, slug, "_prd.md"))
			}
			followUpCommitAt(t, f.repo, now.Add(-time.Hour), "docs: plan imports")
			gitImplement(t, f.repo, "tag", "v0.4.1")
			out, stderr, code := followUpRunCLI(t, "--status")
			if code != 0 || stderr != "" {
				t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
			}
			start, end := strings.Index(out, "Delegated decisions:"), strings.Index(out, "Sources:\n")
			if start < 0 || end < start || out[start:end] != block {
				t.Fatalf("stdout=%s\nwant block=%s", out, block)
			}
			raw, stderr, code := followUpRunCLI(t, "--status", "--format", "json")
			if code != 0 || stderr != "" {
				t.Fatalf("JSON exit=%d stderr=%q", code, stderr)
			}
			var doc followUpStatusDocument
			if err := json.Unmarshal([]byte(raw), &doc); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(doc.DelegatedDecisions, want) {
				t.Fatalf("rows=%+v, want %+v", doc.DelegatedDecisions, want)
			}
			var fields struct {
				Rows []map[string]json.RawMessage `json:"delegated_decisions"`
			}
			if err := json.Unmarshal([]byte(raw), &fields); err != nil {
				t.Fatal(err)
			}
			if fields.Rows == nil {
				t.Fatalf("null or missing rows: %s", raw)
			}
			for _, row := range fields.Rows {
				if len(row) != 10 {
					t.Fatalf("row fields=%v", row)
				}
			}
			found := false
			for _, source := range doc.Sources {
				if source.Name == "specs" {
					found = true
					if name == "an_unreadable_prd_names_its_spec" {
						if source.State != "unavailable" || !strings.Contains(source.Detail, slug) || !strings.Contains(source.Detail, "_prd.md") {
							t.Fatalf("source=%+v", source)
						}
					} else if source.State != "checked" || source.Detail != "" {
						t.Fatalf("source=%+v", source)
					}
				}
			}
			if !found {
				t.Fatal("missing specs source")
			}
		})
	}
}

func TestTheStatusViewAuthoringSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"transcript_1", "transcript_2", "transcript_3"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			now := time.Date(2026, 11, 19, 20, 0, 0, 0, time.UTC)
			f, in := newStatusQueueFixture(t, now)
			gitImplement(t, f.repo, "remote", "remove", "origin")
			if name == "transcript_2" {
				followUpWrite(t, f.repo, ".roundfixrc.yml", "authoring:\n  open_questions: follow\n")
				statusDelegatedSpec(t, in.SpecRoot, "0001-expired-imports", statusDelegatedEntry)
			}
			gitImplement(t, f.repo, "add", ".")
			statusQueueAmendAt(t, f.repo, now.Add(-time.Hour))
			gitImplement(t, f.repo, "tag", "v0.4.1")
			if name == "transcript_3" {
				followUpWrite(t, f.repo, ".roundfixrc.yml", "authoring:\n  open_questions: maybe\n")
			}
			out, stderr, code := followUpRunCLI(t, "--status")
			wantCode, wantErr := 0, ""
			want := fmt.Sprintf("Status for %s in %s\nIn progress: 0\nPending Specs: 0\nSpec queue: 0\nNext release: no unreleased changes\nBlockers: 0\nWaiting on the maintainer: 0\nOpen authoring questions: wait (default)\nDelegated decisions: 0\nSources:\n  runs: checked 0\n  delivery-queue: checked 0\n  pull-requests: checked 0\n  git: checked 0\n  specs: checked 0\n  release-plan: checked 0\n  backlog: checked 0\n  phase-record: checked 0\n  delivery-history: checked 0\n", filepath.Base(f.repo), f.repo)
			if name == "transcript_2" {
				want = fmt.Sprintf("Status for %s in %s\n", filepath.Base(f.repo), f.repo) + statusDelegatedTranscript2
			}
			if name == "transcript_3" {
				wantCode, want = 2, ""
				wantErr = fmt.Sprintf("roundfix: follow-up failed: parse config %q: authoring.open_questions must be wait or follow\nRun 'roundfix follow-up --help' for usage.\n", filepath.Join(f.repo, ".roundfixrc.yml"))
			}
			if code != wantCode || out != want || stderr != wantErr {
				t.Fatalf("exit=%d want=%d\nstderr=%q want=%q\nstdout:\n%s\nwant:\n%s", code, wantCode, stderr, wantErr, out, want)
			}
		})
	}
}

const statusDelegatedTranscript2 = `In progress: 0
Pending Specs: 1
  - 0001-expired-imports not queued 0/1 Tasks
Spec queue: 1
  #  Spec  Status                Start (UTC)         End (UTC)                     Note
  1  0001  not queued 0/1 Tasks  after it is queued  not enough history (0 of 11)  -
Forecast basis: not enough history, 0 of 11 delivered Specs.
Next release: no unreleased changes
Blockers: 0
Waiting on the maintainer: 0
Open authoring questions: follow (Project Config)
Delegated decisions: 1
  0001-expired-imports 2026-11-20 prd: which failure behavior does an expired import have? -> retry automatically up to 3 times (answer: none yet)
Sources:
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: checked 0
  git: checked 0
  specs: checked 1
  release-plan: checked 0
  backlog: checked 0
  phase-record: checked 0
  delivery-history: checked 0
`

func TestTheStatusViewKeepsEveryUnreadablePRDInTheSpecsSource(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 19, 20, 0, 0, 0, time.UTC)
	f, in := newStatusQueueFixture(t, now)
	gitImplement(t, f.repo, "remote", "remove", "origin")
	followUpWrite(t, in.SpecRoot, "0001-missing-graph/_prd.md", "---\nstatus: active\n---\n")
	for _, slug := range []string{"0002-unreadable", "0003-unreadable"} {
		mustMkdir(t, filepath.Join(in.SpecRoot, slug, "_prd.md"))
	}
	followUpCommitAt(t, f.repo, now.Add(-time.Hour), "docs: plan imports")
	doc := followUpStatusJSON(t)
	for _, source := range doc.Sources {
		if source.Name == "specs" {
			if source.State != "unavailable" || strings.Contains(source.Detail, "\n") {
				t.Fatalf("source=%+v", source)
			}
			for _, slug := range []string{"0001-missing-graph", "0002-unreadable", "0003-unreadable"} {
				if !strings.Contains(source.Detail, slug) {
					t.Errorf("source loses %s: %+v", slug, source)
				}
			}
			return
		}
	}
	t.Fatal("missing specs source")
}
