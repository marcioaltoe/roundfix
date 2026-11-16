package phaserecord_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/phaserecord"
)

var promptKeys = []string{"schema", "kind", "time", "started", "ended", "run_id", "spec", "scope_kind", "scope_id", "repository", "work_dir", "harness", "session"}

func TestAPromptEntryIsOneLineOfItsDeclaredKeys(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"keys_in_order", "absent_values_are_empty_strings", "oversized_field"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			recorder := phaserecord.New(home, func() time.Time { return fixedTime })
			entry := phaserecord.PromptEntry{Started: fixedTime.Add(-time.Second), Ended: fixedTime, RunID: "run", Spec: "0412-example", ScopeKind: "task", ScopeID: "task_02", Repository: "/repo", WorkDir: "/work", Harness: "codex", Session: "adapter-session"}
			switch name {
			case "absent_values_are_empty_strings":
				entry = phaserecord.PromptEntry{Started: entry.Started, Ended: entry.Ended}
			case "oversized_field":
				value := strings.Repeat("a", 1023) + "é" + strings.Repeat("b", 100)
				entry.RunID, entry.Spec, entry.ScopeKind, entry.ScopeID = value, value, value, value
				entry.Repository, entry.WorkDir, entry.Harness, entry.Session = value, value, value, value
			}
			for range 2 {
				if err := recorder.RecordPrompt(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
			}
			lines := readLines(t, home, fixedTime)
			if len(lines) != 2 || !bytes.Equal(lines[0], lines[1]) {
				t.Fatalf("lines = %q", lines)
			}
			row := decodeLine(t, lines[0], promptKeys)
			if name == "keys_in_order" {
				want := `{"schema":"roundfix/phase-record/v1","kind":"prompt","time":"2026-11-21T02:05:11Z","started":"2026-11-21T02:05:10Z","ended":"2026-11-21T02:05:11Z","run_id":"run","spec":"0412-example","scope_kind":"task","scope_id":"task_02","repository":"/repo","work_dir":"/work","harness":"codex","session":"adapter-session"}`
				if string(lines[0]) != want {
					t.Fatalf("line = %s, want %s", lines[0], want)
				}
				if runtime.GOOS != "windows" {
					for path, want := range map[string]os.FileMode{filepath.Join(home, ".roundfix", "phase-record"): 0700, filepath.Join(home, ".roundfix", "phase-record", "2026-11-21.jsonl"): 0600} {
						info, err := os.Stat(path)
						if err != nil {
							t.Fatal(err)
						}
						if info.Mode().Perm() != want {
							t.Fatalf("mode %s = %o, want %o", path, info.Mode().Perm(), want)
						}
					}
				}
			}
			for _, key := range promptKeys[5:] {
				switch name {
				case "absent_values_are_empty_strings":
					if string(row[key]) != `""` {
						t.Fatalf("%s = %s", key, row[key])
					}
				case "oversized_field":
					if !utf8.Valid(row[key]) || string(row[key]) != `"`+strings.Repeat("a", 1023)+`"` {
						t.Fatalf("%s = %s", key, row[key])
					}
				}
			}
		})
	}
}

func TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"live", "dry_run", "cutoff_date_is_kept", "not_a_date", "missing_directory"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			dir := filepath.Join(home, ".roundfix", "phase-record")
			names := []string{"2026-11-18.jsonl", "2026-11-19.jsonl", "2026-11-20.jsonl", "2026-11-21.jsonl"}
			if name == "not_a_date" {
				names = append(names, "notes.jsonl", "2026-02-30.jsonl", "2026-1-01.jsonl", "2026-11-17.txt", "prefix-2026-11-17.jsonl")
			}
			if name != "missing_directory" {
				if err := os.MkdirAll(dir, 0700); err != nil {
					t.Fatal(err)
				}
				for _, file := range names {
					if err := os.WriteFile(filepath.Join(dir, file), []byte("entry\n"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				// A directory with a date name is not a day file.
				if err := os.Mkdir(filepath.Join(dir, "2026-11-01.jsonl"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			// Local date is the 19th; the cutoff's UTC date is the 20th.
			cutoff := time.Date(2026, 11, 19, 23, 0, 0, 0, time.FixedZone("fixture", -3*3600))
			count, err := phaserecord.New(home, func() time.Time { return fixedTime }).RemoveBefore(t.Context(), cutoff, name == "dry_run")
			wantCount := 2
			if name == "missing_directory" {
				wantCount = 0
			}
			if err != nil || count != wantCount {
				t.Fatalf("count = %d, err = %v", count, err)
			}
			if name == "missing_directory" {
				if _, err := os.Stat(dir); !os.IsNotExist(err) {
					t.Fatalf("directory created: %v", err)
				}
				return
			}
			got, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			var gotNames []string
			for _, item := range got {
				if !item.IsDir() {
					gotNames = append(gotNames, item.Name())
				}
			}
			var wantNames []string
			for _, file := range names {
				if name != "dry_run" && (file == "2026-11-18.jsonl" || file == "2026-11-19.jsonl") {
					continue
				}
				wantNames = append(wantNames, file)
				data, err := os.ReadFile(filepath.Join(dir, file))
				if err != nil || string(data) != "entry\n" {
					t.Fatalf("kept %s: %s, %v", file, data, err)
				}
			}
			// ReadDir sorts its result; the fixture's order is not significant.
			if len(gotNames) != len(wantNames) {
				t.Fatalf("files = %v, want %v", gotNames, wantNames)
			}
			for _, file := range wantNames {
				if !slices.Contains(gotNames, file) {
					t.Fatalf("missing %s in %v", file, gotNames)
				}
			}
		})
	}
}
