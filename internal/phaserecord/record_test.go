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
	for _, name := range []string{"keys_in_order", "boundary_keys_only_at_a_boundary", "modes", "one_line_per_entry", "oversized_field"} {
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
			if boundary, err := recorder.StartPhase(t.Context(), entry); err != nil || boundary != (phaserecord.Boundary{}) {
				t.Fatalf("first append: %+v, %v", boundary, err)
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
			case "boundary_keys_only_at_a_boundary":
				entry.Phase, entry.Spec = "release", ""
				// A boundary must retain an empty previous_spec as a present key.
				entry.Session = "empty-spec"
				if _, err := recorder.StartPhase(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
				entry.Phase = "follow-up"
				if _, err := recorder.StartPhase(t.Context(), entry); err != nil {
					t.Fatal(err)
				}
				lines = readLines(t, home, fixedTime)
				got := decodeLine(t, lines[2], append(append([]string{}, phaseKeys...), "previous_phase", "previous_spec"))
				if string(got["previous_spec"]) != `""` || string(got["previous_phase"]) != `"release"` {
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
				if _, err := recorder.StartPhase(t.Context(), entry); err != nil {
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

func TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                                                 string
		previous, next                                       phaserecord.Phase
		previousSpec, nextSpec, previousSession, nextSession string
		age                                                  int
		crossed                                              bool
	}{
		{"another_phase", "prd", "techspec", "0412-example", "0412-example", "s", "s", 0, true},
		{"another_spec", "prd", "prd", "0412-example", "0413-another", "s", "s", 0, true},
		{"after_follow_up", "follow-up", "prd", "", "0412-example", "s", "s", 0, false},
		{"delivery_then_release_of_the_spec", "delivery", "release", "0412-example", "0412-example", "s", "s", 0, false},
		{"delivery_then_release_without_a_spec", "delivery", "release", "0412-example", "", "s", "s", 0, false},
		{"delivery_then_release_of_another_spec", "delivery", "release", "0412-example", "0413-another", "s", "s", 0, true},
		{"same_phase_again", "prd", "prd", "0412-example", "0412-example", "s", "s", 0, false},
		{"another_session", "prd", "techspec", "0412-example", "0412-example", "a", "b", 0, false},
		{"no_session", "prd", "techspec", "0412-example", "0412-example", "", "", 0, false},
		{"six_days_old", "prd", "techspec", "0412-example", "0412-example", "s", "s", 6, true},
		{"eight_days_old", "prd", "techspec", "0412-example", "0412-example", "s", "s", 8, false},
		{"latest_entry_decides", "prd", "techspec", "0412-example", "0412-example", "s", "s", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			earlier := fixedTime.AddDate(0, 0, -tc.age).Add(-time.Hour)
			previous := phaserecord.PhaseEntry{Phase: tc.previous, Spec: tc.previousSpec, Session: tc.previousSession}
			if _, err := phaserecord.New(home, func() time.Time { return earlier }).StartPhase(t.Context(), previous); err != nil {
				t.Fatal(err)
			}
			if tc.name == "latest_entry_decides" {
				previous.Phase = "follow-up"
				if _, err := phaserecord.New(home, func() time.Time { return fixedTime.Add(-time.Minute) }).StartPhase(t.Context(), previous); err != nil {
					t.Fatal(err)
				}
				// Append an older timestamp later: latest means the recorded time, not file order.
				previous.Phase = "prd"
				if _, err := phaserecord.New(home, func() time.Time { return earlier }).StartPhase(t.Context(), previous); err != nil {
					t.Fatal(err)
				}
			}
			got, err := phaserecord.New(home, func() time.Time { return fixedTime }).StartPhase(t.Context(), phaserecord.PhaseEntry{Phase: tc.next, Spec: tc.nextSpec, Session: tc.nextSession})
			want := phaserecord.Boundary{}
			keys := append([]string{}, phaseKeys...)
			if tc.crossed {
				want = phaserecord.Boundary{Crossed: true, PreviousPhase: tc.previous, PreviousSpec: tc.previousSpec}
				keys = append(keys, "previous_phase", "previous_spec")
			}
			if err != nil || got != want {
				t.Fatalf("boundary = %+v, %v; want %+v", got, err, want)
			}
			lines := readLines(t, home, fixedTime)
			row := decodeLine(t, lines[len(lines)-1], keys)
			if string(row["boundary"]) != map[bool]string{true: "true", false: "false"}[tc.crossed] || string(row["phase"]) != `"`+string(tc.next)+`"` || string(row["session"]) != `"`+tc.nextSession+`"` || string(row["spec"]) != `"`+tc.nextSpec+`"` {
				t.Fatalf("written = %s", lines[len(lines)-1])
			}
			if tc.crossed && (string(row["previous_phase"]) != `"`+string(tc.previous)+`"` || string(row["previous_spec"]) != `"`+tc.previousSpec+`"`) {
				t.Fatalf("previous = %s", lines[len(lines)-1])
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

func TestStartPhaseSkipsUnrelatedAndMalformedHistory(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	directory := filepath.Join(home, ".roundfix", "phase-record")
	if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	// Malformed lines larger than a scanner's usual token limit must not hide
	// the valid phase entry that follows them.
	history := strings.Repeat("x", 1<<20) + "\nnot json\n" +
		`{"kind":"prompt","phase":"release","session":"s","time":"2026-11-21T02:05:10Z"}` + "\n" +
		`{"kind":"phase","phase":"release","session":"s","time":"bad"}` + "\n" +
		`{"kind":"phase","phase":"prd","spec":"0412-example","session":"s","time":"2026-11-21T01:00:00Z"}` + "\n" +
		`{"kind":"phase","phase":"follow-up","session":"other","time":"2026-11-21T02:05:10Z"}` + "\n"
	if err := os.WriteFile(filepath.Join(directory, "2026-11-21.jsonl"), []byte(history), 0600); err != nil {
		t.Fatal(err)
	}
	ignored := `{"kind":"phase","phase":"follow-up","session":"s","time":"2026-11-21T02:05:10Z"}` + "\n"
	if err := os.WriteFile(filepath.Join(directory, "not-a-date.jsonl"), []byte(ignored), 0600); err != nil {
		t.Fatal(err)
	}
	boundary, err := phaserecord.New(home, func() time.Time { return fixedTime }).StartPhase(t.Context(), phaserecord.PhaseEntry{Phase: "techspec", Spec: "0412-example", Session: "s"})
	want := phaserecord.Boundary{Crossed: true, PreviousPhase: "prd", PreviousSpec: "0412-example"}
	if err != nil || boundary != want {
		t.Fatalf("boundary = %+v, %v, want %+v", boundary, err, want)
	}
	lines := readLines(t, home, fixedTime)
	row := decodeLine(t, lines[len(lines)-1], append(append([]string{}, phaseKeys...), "previous_phase", "previous_spec"))
	if string(row["boundary"]) != "true" {
		t.Fatalf("written = %s", lines[len(lines)-1])
	}
}
