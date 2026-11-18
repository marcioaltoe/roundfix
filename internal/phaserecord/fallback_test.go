package phaserecord

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func fallbackFixtureEntry() FallbackEntry {
	return FallbackEntry{Phase: "techspec", Spec: "0412-example", Repository: "/repo", RepositoryID: "sha256:fixture", Harness: "claude", Session: "3f9a1c20", DeclaredModel: "fable", DeclaredEffort: "xhigh", Model: "opus", Effort: "xhigh", Limit: "Fable", Message: "You've hit your Fable limit · resets Nov 22 at 8pm"}
}

func TestAFallbackEntryIsOneLineOfItsDeclaredKeys(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	now := time.Date(2026, 11, 20, 23, 5, 11, 987, time.FixedZone("west", -3*60*60))
	recorder := New(home, func() time.Time { return now })
	entry := fallbackFixtureEntry()
	want := `{"schema":"roundfix/phase-record/v1","kind":"fallback","time":"2026-11-21T02:05:11Z","phase":"techspec","spec":"0412-example","repository":"/repo","repository_id":"sha256:fixture","harness":"claude","session":"3f9a1c20","declared_model":"fable","declared_effort":"xhigh","model":"opus","effort":"xhigh","limit":"Fable","message":"You've hit your Fable limit · resets Nov 22 at 8pm"}` + "\n"
	if err := recorder.RecordFallback(t.Context(), entry); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, ".roundfix", "phase-record", "2026-11-21.jsonl")
	data, err := os.ReadFile(path)
	if err != nil || string(data) != want {
		t.Fatalf("line = %s, err %v, want %s", data, err, want)
	}
	if runtime.GOOS != "windows" {
		for path, mode := range map[string]os.FileMode{filepath.Dir(path): 0700, path: 0600} {
			info, err := os.Stat(path)
			if err != nil || info.Mode().Perm() != mode {
				t.Fatalf("mode for %s: %v, %v", path, info, err)
			}
		}
	}
	entry.Message = "a quote \" and\na newline"
	entry.Effort, entry.DeclaredEffort, entry.Harness, entry.Session = "", "", "", ""
	if err := recorder.RecordFallback(t.Context(), entry); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(path)
	if err != nil || !strings.HasPrefix(string(data), want) || strings.Count(string(data), "\n") != 2 {
		t.Fatalf("append changed prior line: %s, %v", data, err)
	}
	rows, skipped, err := Fallbacks(home)
	if err != nil || skipped != 0 || len(rows) != 2 || rows[1].FallbackEntry != entry {
		t.Fatalf("read = %+v, %d, %v", rows, skipped, err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := recorder.RecordFallback(ctx, entry); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel = %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(data) {
		t.Fatalf("cancel changed record: %s, %v", after, err)
	}
}

func TestFallbacksReadsBackTheFallbackEntries(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"missing", "oldest_first", "malformed", "non_day_files", "unreadable"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			dir := filepath.Join(home, ".roundfix", "phase-record")
			if name == "missing" {
				rows, n, err := Fallbacks(home)
				if err != nil || n != 0 || len(rows) != 0 {
					t.Fatalf("read = %+v %d %v", rows, n, err)
				}
				if _, err := os.Stat(filepath.Join(home, ".roundfix")); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("reader wrote: %v", err)
				}
				return
			}
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			if name == "unreadable" {
				if err := os.Remove(dir); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(dir, []byte("obstruction"), 0600); err != nil {
					t.Fatal(err)
				}
				if _, _, err := Fallbacks(home); err == nil {
					t.Fatal("accepted obstruction")
				}
				return
			}
			now := time.Date(2026, 11, 21, 2, 5, 11, 0, time.UTC)
			entry := fallbackFixtureEntry()
			// Append out of timestamp order within one day and across days.
			for _, at := range []time.Time{now, now.Add(-time.Hour), now.Add(-24 * time.Hour)} {
				if err := New(home, func() time.Time { return at }).RecordFallback(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
			}
			path := filepath.Join(dir, "2026-11-21.jsonl")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			extra := `{"schema":"roundfix/phase-record/v1","kind":"phase","time":false}` + "\n" + `{"schema":"roundfix/phase-record/v1","kind":"prompt"}` + "\n"
			wantSkipped := 0
			if name == "malformed" {
				extra += "broken\n[]\nnull\n\n" + strings.Replace(string(data), `"techspec"`, `"release"`, 1) + strings.Replace(string(data), `"schema":"roundfix/phase-record/v1"`, `"schema":"other"`, 1) + strings.Replace(string(data), `"phase":"techspec"`, `"phase":false`, 1) + strings.Replace(string(data), `"time":"2026-11-21T02:05:11Z"`, `"time":"bad"`, 1)
				// Each transformed two-line block corrupts only its first line.
				wantSkipped = 8
			}
			if name == "non_day_files" {
				for _, file := range []string{"notes.jsonl", "2026-02-30.jsonl", "2026-1-1.jsonl"} {
					if err := os.WriteFile(filepath.Join(dir, file), []byte("broken"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.Mkdir(filepath.Join(dir, "2026-11-22.jsonl"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			data = append(data, []byte(extra)...)
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			rows, n, err := Fallbacks(home)
			wantCount := 3
			if name == "malformed" {
				wantCount += 4
			}
			if err != nil || n != wantSkipped || len(rows) != wantCount {
				t.Fatalf("read = %+v %d %v; want count %d skipped %d", rows, n, err, wantCount, wantSkipped)
			}
			for i, row := range rows {
				if row.FallbackEntry != entry || i > 0 && row.Time.Before(rows[i-1].Time) {
					t.Fatalf("row %d = %+v", i, row)
				}
			}
			after, err := os.ReadFile(path)
			if err != nil || !reflect.DeepEqual(data, after) {
				t.Fatalf("reader changed record: %v", err)
			}
		})
	}
}

func TestPhasesSkipsAFallbackEntry(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	now := time.Date(2026, 11, 21, 2, 5, 11, 0, time.UTC)
	recorder := New(home, func() time.Time { return now })
	if err := recorder.StartPhase(t.Context(), PhaseEntry{Phase: "idea", Spec: "0412-example"}); err != nil {
		t.Fatal(err)
	}
	if err := recorder.RecordFallback(t.Context(), fallbackFixtureEntry()); err != nil {
		t.Fatal(err)
	}
	rows, n, err := Phases(home, time.Time{})
	if err != nil || n != 0 || len(rows) != 1 || rows[0].Phase != "idea" {
		t.Fatalf("phases = %+v %d %v", rows, n, err)
	}
	// Even an unrecognised payload in a fallback line is not a phase error.
	path := filepath.Join(home, ".roundfix", "phase-record", "2026-11-21.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte(`{"schema":"roundfix/phase-record/v1","kind":"fallback","phase":false,"time":42}`+"\n")...)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	rows, n, err = Phases(home, time.Time{})
	if err != nil || n != 0 || len(rows) != 1 {
		t.Fatalf("phases = %+v %d %v", rows, n, err)
	}
}
