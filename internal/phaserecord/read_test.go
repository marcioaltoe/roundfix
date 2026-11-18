package phaserecord

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPhasesReadsBackThePhaseEntriesOfTheWindow(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 14, 0, 0, 0, time.UTC)
	for _, name := range []string{"oldest_first", "an_entry_an_earlier_version_marked", "prompt_entries_are_skipped", "a_malformed_line_is_counted", "another_schema_is_counted", "before_since_is_left_out", "a_missing_directory_is_an_empty_record", "a_file_that_is_no_day_file_is_ignored", "an_unreadable_record_is_an_error", "nothing_is_written", "written_by_start_phase", "utc_day_and_instant"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			dir := filepath.Join(home, ".roundfix", "phase-record")
			since := now.Add(-24 * time.Hour)
			want := []RecordedPhase{{Time: now, Phase: "prd", Spec: "0001-widget", Repository: "fixture", ModelSource: "unknown", EffortSource: "unknown", Runner: "unknown"}}
			line := func(at time.Time, phase Phase) string {
				t.Helper()
				data, err := json.Marshal(phaseLine{Schema: Schema, Kind: "phase", Time: at.Format(time.RFC3339), Phase: phase, Spec: "0001-widget", Repository: "fixture"})
				if err != nil {
					t.Fatal(err)
				}
				return string(data) + "\n"
			}
			if name == "a_missing_directory_is_an_empty_record" {
				entries, skipped, err := Phases(home, since)
				if err != nil || len(entries) != 0 || skipped != 0 {
					t.Fatalf("entries=%+v skipped=%d err=%v", entries, skipped, err)
				}
				if _, err := os.Stat(filepath.Join(home, ".roundfix")); !os.IsNotExist(err) {
					t.Fatalf("reader created directory: %v", err)
				}
				return
			}
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, now.Format("2006-01-02")+".jsonl")
			content := line(now, "prd")
			skip := 0
			switch name {
			case "utc_day_and_instant":
				since = now.Truncate(24 * time.Hour).Add(-time.Hour).In(time.FixedZone("east", 2*60*60))
				if err := os.WriteFile(filepath.Join(dir, since.UTC().Format("2006-01-02")+".jsonl"), []byte(line(since, "idea")), 0600); err != nil {
					t.Fatal(err)
				}
				want = append([]RecordedPhase{{Time: since, Phase: "idea", Spec: "0001-widget", Repository: "fixture", ModelSource: "unknown", EffortSource: "unknown", Runner: "unknown"}}, want...)
			case "an_entry_an_earlier_version_marked":
				content += `{"schema":"roundfix/phase-record/v1","kind":"phase","time":"` + since.Format(time.RFC3339) + `","phase":"idea","spec":"0001-widget","repository":"fixture","harness":"claude","session":"s","boundary":true,"previous_phase":"prd","previous_spec":"0001-widget"}` + "\n"
				want = append([]RecordedPhase{{Time: since, Phase: "idea", Spec: "0001-widget", Repository: "fixture", ModelSource: "unknown", EffortSource: "unknown", Runner: "unknown"}}, want...)
			case "oldest_first":
				content += line(since, "idea")
				want = append([]RecordedPhase{{Time: since, Phase: "idea", Spec: "0001-widget", Repository: "fixture", ModelSource: "unknown", EffortSource: "unknown", Runner: "unknown"}}, want...)
			case "prompt_entries_are_skipped":
				content += `{"schema":"` + Schema + `","kind":"prompt","time":42,"phase":false}` + "\n"
			case "a_malformed_line_is_counted":
				content += "broken\n[]\nnull\n" + line(now, "invalid") + strings.Replace(line(now, "prd"), now.Format(time.RFC3339), "bad-time", 1)
				skip = 5
			case "another_schema_is_counted":
				content += strings.Replace(line(now, "prd"), Schema, "other/schema", 1)
				skip = 1
			case "before_since_is_left_out":
				content += line(since.Add(-time.Second), "idea")
				if err := os.WriteFile(filepath.Join(dir, since.AddDate(0, 0, -1).Format("2006-01-02")+".jsonl"), []byte("broken\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "a_file_that_is_no_day_file_is_ignored":
				for _, filename := range []string{"notes.jsonl", "2026-02-30.jsonl", "2026-1-1.jsonl"} {
					if err := os.WriteFile(filepath.Join(dir, filename), []byte("broken\n"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.Mkdir(filepath.Join(dir, now.AddDate(0, 0, 1).Format("2006-01-02")+".jsonl"), 0700); err != nil {
					t.Fatal(err)
				}
			case "an_unreadable_record_is_an_error":
				// A regular file in place of the directory fails listing on every host,
				// including a privileged test process and Windows.
				if err := os.Remove(dir); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(dir, []byte("not a directory"), 0600); err != nil {
					t.Fatal(err)
				}
				if _, _, err := Phases(home, since); err == nil {
					t.Fatal("unreadable directory was accepted")
				}
				return
			case "written_by_start_phase":
				if err := New(home, func() time.Time { return now }).StartPhase(context.Background(), PhaseEntry{Phase: "prd", Spec: "0001-widget", Repository: "fixture"}); err != nil {
					t.Fatal(err)
				}
			}
			if !strings.HasPrefix(name, "written_by_start_phase") {
				if err := os.WriteFile(path, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			entries, skipped, err := Phases(home, since)
			if err != nil || skipped != skip || len(entries) != len(want) {
				t.Fatalf("entries=%+v want=%+v skipped=%d want=%d err=%v", entries, want, skipped, skip, err)
			}
			for i, entry := range entries {
				if !entry.Time.Equal(want[i].Time) || entry.Phase != want[i].Phase || entry.Spec != want[i].Spec || entry.Repository != want[i].Repository || entry.RepositoryID != want[i].RepositoryID || entry.Model != want[i].Model || entry.ModelSource != want[i].ModelSource || entry.Effort != want[i].Effort || entry.EffortSource != want[i].EffortSource || entry.Runner != want[i].Runner {
					t.Fatalf("entry[%d]=%+v want=%+v", i, entry, want[i])
				}
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(before) != string(after) {
				t.Fatal("reader changed record")
			}
		})
	}
}
