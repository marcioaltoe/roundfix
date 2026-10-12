package spec

// Suite: Pull request recording preserves every other Archive Record field.
// Boundary: parsed and rendered record files, including idempotence and refusals.

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRecordArchivePullRequestChangesOnlyThatField(t *testing.T) {
	t.Parallel()
	for _, disposition := range []ArchiveDisposition{ArchivePass, ArchivePartial, ArchiveQAOverride, ArchiveSuperseded, ArchiveFailedQA, ArchiveNoQA} {
		for _, long := range []bool{false, true} {
			t.Run(string(disposition)+"/"+map[bool]string{false: "short", true: "long"}[long], func(t *testing.T) {
				r := ArchiveRecord{
					Spec: "0256-fixture", Title: "Archive fixture", Created: "2026-10-01", Archived: "2026-10-12",
					Disposition: disposition, Source: "docs/specs/0256-fixture", SourceRevision: strings.Repeat("a", 40),
					QATask: "task_qa", QAReport: "qa/qa-report-2026-10-12.md", QAVerdict: VerdictPass,
					Unproven: []string{"published check"}, ADRs: []string{"ADR-0261"}, Sources: []string{"adopted.md"},
					Promoted: []string{"docs/references/lesson.md"}, Relinked: []string{"README.md"},
					Regeneration:   []ArchiveRegeneration{{Command: "make example", Outputs: []string{"generated.json"}}, {Command: "make empty", Outputs: []string{}}},
					DeliveryCommit: strings.Repeat("b", 40), Outcome: "The outcome stays unchanged.",
				}
				switch disposition {
				case ArchiveQAOverride:
					r.QAOverride = &QAArchiveOverrideRecord{"maintainer", "approved", "fail", "failed", r.SourceRevision}
				case ArchiveSuperseded:
					r.SupersededBy = "0257-successor"
				case ArchiveFailedQA:
					r.QAVerdict = VerdictFail
				}
				if long {
					r.Outcome = strings.Repeat("é", 2048)
				}
				original, err := RenderArchiveRecord(r)
				if err != nil {
					t.Fatal(err)
				}
				before, err := ParseArchiveRecord(original)
				if err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(t.TempDir(), r.Spec+".md")
				writeFile(t, path, string(original))
				changed, err := RecordArchivePullRequest(path, "7")
				if err != nil || !changed {
					t.Fatalf("changed=%v err=%v", changed, err)
				}
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				after, err := ParseArchiveRecord(content)
				if err != nil {
					t.Fatal(err)
				}
				expected := before
				expected.PullRequest = "7"
				rendered, err := RenderArchiveRecord(expected)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(content, rendered) {
					t.Fatal("writer differs from RenderArchiveRecord")
				}
				if len(content) > ArchiveRecordTargetBytes || !utf8.Valid(content) {
					t.Fatal("record exceeded byte target or broke UTF-8")
				}
				if !long && after.Outcome != before.Outcome {
					t.Fatal("short outcome changed")
				}
				if long && (len(after.Outcome) >= len(before.Outcome) || !strings.HasPrefix(before.Outcome, strings.TrimSuffix(after.Outcome, "…"))) {
					t.Fatal("outcome was not shortened")
				}
				after.PullRequest = before.PullRequest
				after.Outcome = before.Outcome
				if !reflect.DeepEqual(before, after) {
					t.Fatalf("other fields changed:\nbefore=%+v\nafter=%+v", before, after)
				}
				// A replay must preserve even noncanonical formatting byte for byte.
				same := bytes.Replace(content, []byte("schema: "+ArchiveRecordSchema), []byte("schema: '"+ArchiveRecordSchema+"'"), 1)
				writeFile(t, path, string(same))
				changed, err = RecordArchivePullRequest(path, "7")
				if err != nil || changed {
					t.Fatalf("replay changed=%v err=%v", changed, err)
				}
				got, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(got, same) {
					t.Fatalf("replay rewrote bytes: %v", err)
				}
				for _, number := range []string{"8", "0", "-7", "seven", "07", "", "12345678901"} {
					if changed, err := RecordArchivePullRequest(path, number); err == nil || changed {
						t.Fatalf("accepted %q: changed=%v err=%v", number, changed, err)
					}
					got, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(got, same) {
						t.Fatalf("refusal rewrote bytes: %v", err)
					}
				}
			})
		}
	}
	t.Run("malformed", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "malformed.md")
		writeFile(t, path, "not a record\n")
		if changed, err := RecordArchivePullRequest(path, "7"); err == nil || changed {
			t.Fatal("accepted malformed record")
		}
		if got, err := os.ReadFile(path); err != nil || string(got) != "not a record\n" {
			t.Fatal("malformed record changed")
		}
	})
	t.Run("missing", func(t *testing.T) {
		if changed, err := RecordArchivePullRequest(filepath.Join(t.TempDir(), "missing.md"), "7"); err == nil || changed {
			t.Fatal("accepted missing record")
		}
	})
}
