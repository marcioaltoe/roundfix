package phaserecord_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/phaserecord"
)

var fixedTime = time.Date(2026, 11, 20, 23, 5, 11, 987, time.FixedZone("fixture", -3*60*60))
var phaseKeys = []string{"schema", "kind", "time", "phase", "spec", "repository", "harness", "session", "boundary"}

func readLines(t *testing.T, home string, now time.Time) [][]byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(home, ".roundfix", "phase-record", now.UTC().Format("2006-01-02")+".jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasSuffix(data, []byte("\n")) {
		t.Fatal("entry lacks newline")
	}
	return bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
}

func decodeLine(t *testing.T, line []byte, keys []string) map[string]json.RawMessage {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(line))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		t.Fatalf("object: %v, %v", token, err)
	}
	var gotKeys []string
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		gotKeys = append(gotKeys, token.(string))
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(gotKeys, keys) {
		t.Fatalf("keys = %v, want %v", gotKeys, keys)
	}
	var row map[string]json.RawMessage
	if err := json.Unmarshal(line, &row); err != nil {
		t.Fatal(err)
	}
	return row
}

func TestAPhaseEntryIsOneLineOfItsDeclaredKeys(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"keys_in_order", "no_boundary_keys_after_a_release", "modes", "one_line_per_entry", "oversized_field"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			recorder := phaserecord.New(home, func() time.Time { return fixedTime })
			entry := phaserecord.PhaseEntry{Phase: "prd", Spec: "0412-example", Repository: "/repo", Harness: "claude", Session: "session"}
			if name == "oversized_field" {
				entry.Repository = strings.Repeat("a", 1023) + "é"
				entry.Spec = strings.Repeat("s", 1100)
				entry.Harness = strings.Repeat("h", 1100)
				entry.Session = strings.Repeat("x", 1100)
			}
			if err := recorder.StartPhase(t.Context(), entry); err != nil {
				t.Fatal(err)
			}
			lines := readLines(t, home, fixedTime)
			if len(lines) != 1 {
				t.Fatalf("lines = %d", len(lines))
			}
			row := decodeLine(t, lines[0], phaseKeys)
			if string(row["schema"]) != `"roundfix/phase-record/v1"` || string(row["kind"]) != `"phase"` || string(row["time"]) != `"2026-11-21T02:05:11Z"` || string(row["boundary"]) != "false" {
				t.Fatalf("row = %s", lines[0])
			}
			switch name {
			case "keys_in_order":
				want := `{"schema":"roundfix/phase-record/v1","kind":"phase","time":"2026-11-21T02:05:11Z","phase":"prd","spec":"0412-example","repository":"/repo","harness":"claude","session":"session","boundary":false}`
				if string(lines[0]) != want {
					t.Fatalf("line = %s, want %s", lines[0], want)
				}
			case "no_boundary_keys_after_a_release":
				entry.Phase, entry.Spec = "release", ""
				// A release without a Spec must not add boundary keys to the next entry.
				entry.Session = "empty-spec"
				if err := recorder.StartPhase(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
				entry.Phase = "follow-up"
				if err := recorder.StartPhase(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
				lines = readLines(t, home, fixedTime)
				got := decodeLine(t, lines[2], phaseKeys)
				if string(got["boundary"]) != "false" {
					t.Fatalf("boundary = %s", lines[2])
				}
			case "modes":
				if runtime.GOOS != "windows" {
					for path, want := range map[string]os.FileMode{filepath.Join(home, ".roundfix", "phase-record"): 0700, filepath.Join(home, ".roundfix", "phase-record", "2026-11-21.jsonl"): 0600} {
						info, err := os.Stat(path)
						if err != nil {
							t.Fatal(err)
						}
						if info.Mode().Perm() != want {
							t.Errorf("%s: mode = %o, want %o", path, info.Mode().Perm(), want)
						}
					}
				}
			case "one_line_per_entry":
				entry.Repository = "line\nbreak\"\\"
				if err := recorder.StartPhase(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
				lines = readLines(t, home, fixedTime)
				if len(lines) != 2 {
					t.Fatalf("lines = %d", len(lines))
				}
				got := decodeLine(t, lines[1], phaseKeys)
				var repository string
				if err := json.Unmarshal(got["repository"], &repository); err != nil {
					t.Fatal(err)
				}
				if repository != entry.Repository {
					t.Fatalf("repository = %q", repository)
				}
			case "oversized_field":
				for _, key := range []string{"spec", "repository", "harness", "session"} {
					var value string
					if err := json.Unmarshal(row[key], &value); err != nil {
						t.Fatal(err)
					}
					wantLength := 1024
					if key == "repository" {
						wantLength = 1023
					}
					if len(value) != wantLength || !utf8.ValidString(value) {
						t.Errorf("%s: length = %d, valid UTF-8 = %t", key, len(value), utf8.ValidString(value))
					}
				}
			}
		})
	}
}

func TestStartPhaseMarksNoBoundaryWhateverCameBefore(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                   string
		previous, next         phaserecord.Phase
		previousSpec, nextSpec string
	}{
		{"after_follow_up", "follow-up", "techspec", "", "0412-example"},
		{"after_idea", "idea", "techspec", "0412-example", "0412-example"},
		{"after_prd", "prd", "techspec", "0412-example", "0412-example"},
		{"after_techspec", "techspec", "techspec", "0412-example", "0412-example"},
		{"after_tasks", "tasks", "techspec", "0412-example", "0412-example"},
		{"after_delivery", "delivery", "techspec", "0412-example", "0412-example"},
		{"after_release", "release", "follow-up", "0412-example", ""},
		{"after_prd_of_another_spec", "prd", "techspec", "0413-another", "0412-example"},
		{"release_of_another_spec_after_release", "release", "release", "0412-example", "0413-another"},
		{"release_of_the_same_spec_again", "release", "release", "0412-example", "0412-example"},
		{"release_without_a_spec_again", "release", "release", "", ""},
		{"after_a_marked_entry", "prd", "techspec", "0412-example", "0412-example"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			earlier := fixedTime.Add(-time.Minute)
			if tc.name == "after_a_marked_entry" {
				directory := filepath.Join(home, ".roundfix", "phase-record")
				if err := os.MkdirAll(directory, 0700); err != nil {
					t.Fatal(err)
				}
				line := `{"schema":"roundfix/phase-record/v1","kind":"phase","time":"` + earlier.UTC().Format(time.RFC3339) + `","phase":"prd","spec":"0412-example","repository":"/repo","harness":"claude","session":"s","boundary":true,"previous_phase":"idea","previous_spec":"0412-example"}` + "\n"
				if err := os.WriteFile(filepath.Join(directory, earlier.UTC().Format("2006-01-02")+".jsonl"), []byte(line), 0600); err != nil {
					t.Fatal(err)
				}
			} else if err := phaserecord.New(home, func() time.Time { return earlier }).StartPhase(t.Context(), phaserecord.PhaseEntry{Phase: tc.previous, Spec: tc.previousSpec, Session: "s"}); err != nil {
				t.Fatal(err)
			}
			before := readLines(t, home, fixedTime)
			if err := phaserecord.New(home, func() time.Time { return fixedTime }).StartPhase(t.Context(), phaserecord.PhaseEntry{Phase: tc.next, Spec: tc.nextSpec, Session: "s"}); err != nil {
				t.Fatal(err)
			}
			lines := readLines(t, home, fixedTime)
			if len(lines) != len(before)+1 {
				t.Fatalf("lines = %d, want %d", len(lines), len(before)+1)
			}
			for i, line := range before {
				if !bytes.Equal(line, lines[i]) {
					t.Fatalf("earlier line %d changed: %s", i, lines[i])
				}
			}
			row := decodeLine(t, lines[len(lines)-1], phaseKeys)
			if string(row["boundary"]) != "false" || string(row["phase"]) != `"`+string(tc.next)+`"` || string(row["session"]) != `"s"` || string(row["spec"]) != `"`+tc.nextSpec+`"` {
				t.Fatalf("written = %s", lines[len(lines)-1])
			}
		})
	}
}

func TestParsePhaseKnowsTheSevenPhases(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"follow-up", "idea", "prd", "techspec", "tasks", "delivery", "release"} {
		t.Run(name, func(t *testing.T) {
			got, err := phaserecord.ParsePhase(name)
			if err != nil || string(got) != name {
				t.Fatalf("ParsePhase = %q, %v", got, err)
			}
		})
	}
	for _, name := range []string{"", "plan", "PRD"} {
		if _, err := phaserecord.ParsePhase(name); err == nil {
			t.Errorf("accepted unknown phase %q", name)
		}
	}
}

func TestStartPhaseReadsNoEarlierEntry(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"an_unreadable_earlier_day", "a_malformed_day_file"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			directory := filepath.Join(home, ".roundfix", "phase-record")
			if err := os.MkdirAll(directory, 0700); err != nil {
				t.Fatal(err)
			}
			var before []byte
			if name == "an_unreadable_earlier_day" {
				yesterday := fixedTime.UTC().AddDate(0, 0, -1)
				if err := os.Mkdir(filepath.Join(directory, yesterday.Format("2006-01-02")+".jsonl"), 0700); err != nil {
					t.Fatal(err)
				}
			} else {
				earlier := fixedTime.Add(-time.Minute)
				if err := phaserecord.New(home, func() time.Time { return earlier }).StartPhase(t.Context(), phaserecord.PhaseEntry{Phase: "release", Session: "s"}); err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(directory, fixedTime.UTC().Format("2006-01-02")+".jsonl")
				valid, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				before = append([]byte(strings.Repeat("x", 1<<20)+"\n"), valid...)
				if err := os.WriteFile(path, before, 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := phaserecord.New(home, func() time.Time { return fixedTime }).StartPhase(t.Context(), phaserecord.PhaseEntry{Phase: "techspec", Spec: "0412-example", Session: "s"}); err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(filepath.Join(directory, fixedTime.UTC().Format("2006-01-02")+".jsonl"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(after, before) {
				t.Fatal("earlier bytes changed")
			}
			appended := after[len(before):]
			if bytes.Count(appended, []byte("\n")) != 1 || !bytes.HasSuffix(appended, []byte("\n")) {
				t.Fatalf("expected exactly one appended line: %q", appended)
			}
			row := decodeLine(t, bytes.TrimSuffix(appended, []byte("\n")), phaseKeys)
			if string(row["boundary"]) != "false" {
				t.Fatalf("written = %s", appended)
			}
		})
	}
}
