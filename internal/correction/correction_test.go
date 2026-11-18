package correction

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidateNamesTheAcceptedValues(t *testing.T) {
	t.Parallel()
	for _, class := range Answers() {
		t.Run(string(class), func(t *testing.T) {
			d := Declaration{Class: class, Reason: "what is known"}
			if class == "plan" {
				d.Artifact = "tasks"
			}
			if err := Validate(d); err != nil {
				t.Fatal(err)
			}
		})
	}
	cases := []struct {
		name string
		d    Declaration
		want string
	}{
		{"plan without artifact", Declaration{Class: "plan", Reason: "known"}, "prd, techspec or tasks"},
		{"unknown artifact", Declaration{Class: "plan", Artifact: "idea", Reason: "known"}, "prd, techspec or tasks"},
		{"implementation artifact", Declaration{Class: "implementation", Artifact: "prd", Reason: "known"}, "only with class plan"},
		{"not a correction artifact", Declaration{Class: NotACorrection, Artifact: "prd", Reason: "known"}, "only with class plan"},
		{"unknown class", Declaration{Class: "guess", Reason: "known"}, "plan, implementation, environment, operator, unclear or not-a-correction"},
		{"blank reason", Declaration{Class: "unclear", Reason: "  "}, "1 to 200 characters on one line"},
		{"long reason", Declaration{Class: "unclear", Reason: strings.Repeat("é", 201)}, "1 to 200 characters on one line"},
		{"newline", Declaration{Class: "unclear", Reason: "first\nsecond"}, "1 to 200 characters on one line"},
		{"control", Declaration{Class: "unclear", Reason: "first\x00second"}, "no control character"},
		{"line separator", Declaration{Class: "unclear", Reason: "first\u2028second"}, "line or paragraph separator"},
		{"paragraph separator", Declaration{Class: "unclear", Reason: "first\u2029second"}, "line or paragraph separator"},
		{"by", Declaration{Class: "unclear", Reason: "known", By: strings.Repeat("x", 65)}, "1 to 64 characters"},
		{"task", Declaration{Class: "unclear", Reason: "known", Signs: []Sign{{Kind: "task", ID: "task_no"}}}, "task_<digits>"},
		{"id", Declaration{Class: "unclear", Reason: "known", Correction: "C0"}, "C<n>"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := Validate(c.d); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("Validate = %v, want %q", err, c.want)
			}
		})
	}
	if err := Validate(Declaration{Class: "unclear", Reason: strings.Repeat("é", 200)}); err != nil {
		t.Fatal(err)
	}
}
func TestNextIDCountsFromOneAndNeverReuses(t *testing.T) {
	t.Parallel()
	if got := NextID(File{}); got != "C1" {
		t.Fatal(got)
	}
	f := File{Declarations: []Declaration{{Correction: "C4"}, {Correction: "C2"}, {Correction: "C4", Withdrawn: true}}}
	if got := NextID(f); got != "C5" {
		t.Fatal(got)
	}
	f.Declarations = append(f.Declarations, Declaration{Correction: "C999999999999999999999999"})
	if got := NextID(f); got != "C1000000000000000000000000" {
		t.Fatal(got)
	}
}
func TestStandingTakesTheLastDeclarationOfEachCorrection(t *testing.T) {
	t.Parallel()
	sign := Sign{Kind: "commit", ID: "a"}
	first := Declaration{Correction: "C1", Class: "plan", Artifact: "tasks", Reason: "first", Signs: []Sign{sign}}
	later := first
	later.Class = "implementation"
	later.Artifact = ""
	later.Reason = "later"
	withdrawn := later
	withdrawn.Withdrawn = true
	revived := later
	revived.Reason = "again"
	noDefect := Declaration{Correction: "C2", Class: NotACorrection, Reason: "decision", Signs: []Sign{{Kind: "task", ID: "task_02"}}}
	f := File{Declarations: []Declaration{first, later}}
	if got := Standing(f); !reflect.DeepEqual(got, []Declaration{later}) {
		t.Fatalf("%+v", got)
	}
	f.Declarations = append(f.Declarations, withdrawn)
	if got := Standing(f); !got[0].Withdrawn || len(got[0].Signs) != 0 {
		t.Fatalf("%+v", got)
	}
	f.Declarations = append(f.Declarations, revived, noDefect)
	if got := Standing(f); !reflect.DeepEqual(got, []Declaration{revived, noDefect}) {
		t.Fatalf("%+v", got)
	}
	if first.Signs[0] != sign {
		t.Fatal("mutated input")
	}
}
func TestUndeclaredCountsEverySignNoStandingDeclarationNames(t *testing.T) {
	t.Parallel()
	signs := []Sign{{Kind: "commit", ID: "a", Subject: "subject"}, {Kind: "finding", ID: "F1@a"}, {Kind: "task", ID: "task_01"}, {Kind: "commit", ID: "b"}}
	standing := []Declaration{{Correction: "C1", Signs: []Sign{{Kind: "commit", ID: "a"}}}, {Correction: "C2", Class: NotACorrection, Signs: signs[1:2]}, {Correction: "C3", Withdrawn: true, Signs: signs[2:3]}}
	if got := Undeclared(signs, standing); !reflect.DeepEqual(got, signs[2:]) {
		t.Fatalf("%+v", got)
	}
}
func TestTheTwoLinesKeepTheirKeyOrder(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 20, 14, 2, 11, 0, time.UTC)
	d := Declaration{Time: now, Correction: "C1", Class: "plan", Artifact: "techspec", Reason: "known", Signs: []Sign{{Kind: "finding", ID: "F1@abc", Subject: "must not be written"}}, Harness: "claude", Session: "session"}
	line, err := EncodeDeclaration("0412-example", d)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"schema":"roundfix/corrections/v1","kind":"declaration","time":"2026-11-20T14:02:11Z","spec":"0412-example","correction":"C1","class":"plan","artifact":"techspec","reason":"known","signs":[{"kind":"finding","id":"F1@abc"}],"withdrawn":false,"by":"","harness":"claude","session":"session"}` + "\n"
	if string(line) != want {
		t.Fatalf("got %s want %s", line, want)
	}
	snapshot, err := EncodeSnapshot("0412-example", Snapshot{Time: now, Head: "abc", Signs: []Sign{{Kind: "commit", ID: "def", Subject: "fix"}}, Dismissed: 1})
	if err != nil {
		t.Fatal(err)
	}
	want = `{"schema":"roundfix/corrections/v1","kind":"signs","time":"2026-11-20T14:02:11Z","spec":"0412-example","head":"abc","signs":[{"kind":"commit","id":"def","subject":"fix"}],"dismissed":1,"unavailable":""}` + "\n"
	if string(snapshot) != want {
		t.Fatalf("got %s want %s", snapshot, want)
	}
	empty, err := EncodeDeclaration("0412-example", Declaration{Time: now, Withdrawn: true, Signs: d.Signs})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(empty, []byte(`"signs":[]`)) {
		t.Fatal(string(empty))
	}
}
func TestReadSkipsAndCountsADamagedLine(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "record.jsonl")
	if f, err := Read(path); err != nil || len(f.Declarations) != 0 {
		t.Fatalf("%+v %v", f, err)
	}
	now := time.Date(2026, 11, 20, 14, 2, 11, 0, time.UTC)
	d, _ := EncodeDeclaration("0412-example", Declaration{Time: now, Correction: "C1", Class: "unclear", Reason: "known"})
	s, _ := EncodeSnapshot("0412-example", Snapshot{Time: now, Head: "head"})
	// Unknown fields remain lenient, while malformed kinds, schemas and types count.
	d = bytes.Replace(d, []byte(`"kind":"declaration"`), []byte(`"extra":true,"kind":"declaration"`), 1)
	content := append(append(d, []byte("broken\n{\"schema\":\"other\",\"kind\":\"declaration\"}\n{\"schema\":\"roundfix/corrections/v1\",\"kind\":\"other\"}\n{\"schema\":\"roundfix/corrections/v1\",\"kind\":\"signs\",\"time\":4}\n")...), s...)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	f, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if f.Skipped != 4 || len(f.Declarations) != 1 || len(f.Snapshots) != 1 {
		t.Fatalf("%+v", f)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(after, content) {
		t.Fatal("reader wrote", err)
	}
}
func TestRemoveBeforeRemovesAWholeFilePastTheCutoff(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	cutoff := time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)
	old, _ := EncodeDeclaration("0412-example", Declaration{Time: cutoff.Add(-time.Hour)})
	newer, _ := EncodeSnapshot("0412-example", Snapshot{Time: cutoff})
	write := func(id, slug string, data []byte) string {
		p := Path(home, id, slug)
		if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, data, 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	oldPath := write("aaaaaaaaaaaaaaaa", "0412-example", old)
	kept := write("bbbbbbbbbbbbbbbb", "0413-example", append(append([]byte(nil), old...), newer...))
	damaged := write("cccccccccccccccc", "0414-example", []byte("broken\n"))
	if err := os.Chtimes(damaged, cutoff.Add(-time.Hour), cutoff.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	recent := write("dddddddddddddddd", "0415-example", []byte("broken\n"))
	if err := os.Chtimes(recent, cutoff, cutoff); err != nil {
		t.Fatal(err)
	}
	n, err := RemoveBefore(t.Context(), home, cutoff, true)
	if err != nil || n != 2 {
		t.Fatalf("%d %v", n, err)
	}
	for _, p := range []string{oldPath, damaged, kept, recent} {
		if _, err := os.Stat(p); err != nil {
			t.Fatal(err)
		}
	}
	n, err = RemoveBefore(t.Context(), home, cutoff, false)
	if err != nil || n != 2 {
		t.Fatalf("%d %v", n, err)
	}
	for _, p := range []string{oldPath, damaged} {
		if _, err := os.Stat(filepath.Dir(p)); !os.IsNotExist(err) {
			t.Fatalf("directory remains: %v", err)
		}
	}
	for _, p := range []string{kept, recent} {
		if _, err := os.Stat(p); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := RemoveBefore(ctx, home, cutoff, false); err == nil {
		t.Fatal("ignored cancellation")
	}
}
