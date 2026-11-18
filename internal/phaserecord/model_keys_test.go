package phaserecord_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"roundfix/internal/phaserecord"
)

var authoringKeys = append(append([]string(nil), phaseKeys...), "repository_id", "model", "model_source", "effort", "effort_source", "runner")

func TestAnAuthoringPhaseEntryHoldsTheSixKeysInOrder(t *testing.T) {
	t.Parallel()
	for _, phase := range []phaserecord.Phase{"idea", "prd", "techspec", "tasks", "release"} {
		t.Run(string(phase), func(t *testing.T) {
			t.Parallel()
			for _, source := range []string{"environment", "reported", "unknown"} {
				t.Run(source, func(t *testing.T) {
					home := t.TempDir()
					entry := phaserecord.PhaseEntry{Phase: phase, Spec: "0412-example", Repository: "/repo", Harness: "claude", Session: "s", RepositoryID: "sha256:fixture", Model: "fixture-model", ModelSource: "reported", Effort: "high", EffortSource: source, Runner: "author-subagent"}
					if source == "reported" {
						entry.Runner = "session"
					}
					if source == "unknown" {
						entry.Model, entry.ModelSource, entry.Effort, entry.Runner = "", "unknown", "", "unknown"
						entry.Harness, entry.Session = "", ""
					}
					if err := phaserecord.New(home, func() time.Time { return fixedTime }).StartPhase(t.Context(), entry); err != nil {
						t.Fatal(err)
					}
					lines := readLines(t, home, fixedTime)
					if len(lines) != 1 {
						t.Fatalf("lines = %d", len(lines))
					}
					keys := authoringKeys
					if phase == "release" {
						keys = phaseKeys
					}
					decodeLine(t, lines[0], keys)
					entries, skipped, err := phaserecord.Phases(home, fixedTime.Add(-time.Hour))
					if err != nil || skipped != 0 || len(entries) != 1 {
						t.Fatalf("read = %+v, %d, %v", entries, skipped, err)
					}
					want := phaserecord.RecordedPhase{Time: fixedTime.UTC().Truncate(time.Second), Phase: phase, Spec: entry.Spec, Repository: entry.Repository, RepositoryID: entry.RepositoryID, Model: entry.Model, ModelSource: entry.ModelSource, Effort: entry.Effort, EffortSource: entry.EffortSource, Runner: entry.Runner}
					if phase == "release" {
						want.RepositoryID, want.Model, want.Effort = "", "", ""
						want.ModelSource, want.EffortSource, want.Runner = "unknown", "unknown", "unknown"
					}
					if !reflect.DeepEqual(entries[0], want) {
						t.Fatalf("entry = %+v, want %+v", entries[0], want)
					}
				})
			}
		})
	}
}

func TestAnEntryWithoutTheSixKeysReadsBackAsUnknown(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	dir := filepath.Join(home, ".roundfix", "phase-record")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	line := `{"schema":"roundfix/phase-record/v1","kind":"phase","time":"2026-11-21T02:05:11Z","phase":"prd","spec":"0412-example","repository":"/repo","harness":"claude","session":"s","boundary":false}` + "\n"
	if err := os.WriteFile(filepath.Join(dir, "2026-11-21.jsonl"), []byte(line), 0600); err != nil {
		t.Fatal(err)
	}
	entries, skipped, err := phaserecord.Phases(home, fixedTime.Add(-time.Hour))
	if err != nil || skipped != 0 || len(entries) != 1 {
		t.Fatalf("read = %+v, %d, %v", entries, skipped, err)
	}
	want := phaserecord.RecordedPhase{Time: fixedTime.UTC().Truncate(time.Second), Phase: "prd", Spec: "0412-example", Repository: "/repo", ModelSource: "unknown", EffortSource: "unknown", Runner: "unknown"}
	if !reflect.DeepEqual(entries[0], want) {
		t.Fatalf("entry = %+v, want %+v", entries[0], want)
	}
}
