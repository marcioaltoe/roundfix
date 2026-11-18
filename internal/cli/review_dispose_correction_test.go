package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/correction"
	"roundfix/internal/gittest"
)

func correctionDispositionFixture(t *testing.T, specs []string) (reviewCommandFixture, string) {
	t.Helper()
	f := newReviewDispositionFixture(t)
	mustMkdir(t, filepath.Join(f.repository, "docs/specs/0412-example"))
	mustWrite(t, filepath.Join(f.repository, "docs/specs/0412-example/_prd.md"), "# Fixture\n")
	record := dispositionReviewRecord(f.repository, f.headCommit, "Findings:\n- internal/cli/review.go:10: finding")
	record.Specs = specs
	writeDispositionRecord(t, f, record)
	mustWrite(t, filepath.Join(f.repository, "fix.txt"), "fixed\n")
	gittest.Run(t, f.repository, "add", "fix.txt")
	gittest.Run(t, f.repository, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "fix finding")
	updateCommandDependenciesForTest(t, func(d *commandDependencies) { d.getenv = func(string) string { return "" } })
	return f, strings.TrimSpace(gittest.Run(t, f.repository, "rev-parse", "HEAD"))
}
func TestReviewDisposeSurfaceTranscript8(t *testing.T) {
	t.Parallel()
	home, _ := correctionsFixture(t)
	before := correctionTree(t, home)
	code, out, diag := runReviewDispose(t, "F1", "--dismiss", "--evidence", "not a defect", "--class", "plan", "--reason", "The plan named the wrong file")
	if code != 2 || out != "" || diag != "roundfix: review dispose refused: --class cannot be combined with --dismiss\n" {
		t.Fatalf("%d %q %q", code, out, diag)
	}
	if !reflect.DeepEqual(before, correctionTree(t, home)) {
		t.Fatal("dismissal wrote home")
	}
}

type correctionAfterDispositionWriter struct {
	bytes.Buffer
	after func() error
}

func (w *correctionAfterDispositionWriter) Write(data []byte) (int, error) {
	n, err := w.Buffer.Write(data)
	if err != nil {
		return n, err
	}
	if w.after != nil {
		after := w.after
		w.after = nil
		if err := after(); err != nil {
			return n, err
		}
	}
	return n, nil
}

func TestReviewDisposeRecordsACorrectionWithTheFix(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"inferred spec", "explicit spec", "existing correction", "not a correction", "plan"} {
		t.Run(mode, func(t *testing.T) {
			specs := []string{"0412-example"}
			if mode == "explicit spec" {
				specs = []string{"0412-example", "0413-other"}
			}
			f, fix := correctionDispositionFixture(t, specs)
			args := []string{"F1", "--fixed-by", fix}
			switch mode {
			case "existing correction":
				assertCorrectionCall(t, "Correction C1 of 0412-example: implementation; covers Corrective Task task_01.", "declare", "--spec", "0412-example", "--class", "implementation", "--reason", "same cause", "--task", "task_01")
				args = append(args, "--correction", "C1")
			case "not a correction":
				args = append(args, "--class", "not-a-correction", "--reason", "maintainer decision")
			case "plan":
				args = append(args, "--class", "plan", "--artifact", "techspec", "--reason", "wrong plan")
			default:
				args = append(args, "--class", "implementation", "--reason", "wrong implementation")
			}
			if mode == "explicit spec" {
				args = append(args, "--spec", "0412-example")
			}
			code, out, diag := runReviewDispose(t, args...)
			ledger := readSingleDispositionLine(t, f.artifactDir)
			if code != 0 || diag != "" || !strings.HasPrefix(out, ledger+"Correction C1 of 0412-example: ") {
				t.Fatalf("%d %q %q", code, out, diag)
			}
			assertDispositionEntry(t, decodeDispositionLine(t, ledger), f, "F1", "internal/cli/review.go:10: finding", "fixed")
			record := readCorrectionFixture(t, f.homeDir, f.repository)
			standing := correction.Standing(record)
			if len(standing) != 1 {
				t.Fatalf("%+v", record)
			}
			signs := standing[0].Signs
			want := []correction.Sign{{Kind: "finding", ID: "F1@" + f.headCommit}, {Kind: "commit", ID: fix}}
			if mode == "existing correction" {
				want = append([]correction.Sign{{Kind: "task", ID: "task_01"}}, want...)
				if len(record.Declarations) != 2 {
					t.Fatal("did not append")
				}
			}
			if !reflect.DeepEqual(signs, want) {
				t.Fatalf("signs=%+v want %+v", signs, want)
			}
			// Standalone resolution finds the same fixed finding by its bare or qualified name.
			for _, name := range []string{"F1", "F1@" + f.headCommit[:9]} {
				code, _, diag := correctionCall(t, "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known", "--finding", name)
				if code != 0 || diag != "" {
					t.Fatalf("finding %s: %d %s", name, code, diag)
				}
			}
		})
	}
	t.Run("Correction write after disposition fails", func(t *testing.T) {
		f, fix := correctionDispositionFixture(t, []string{"0412-example"})
		var diagnostics bytes.Buffer
		out := &correctionAfterDispositionWriter{after: func() error {
			if err := os.MkdirAll(filepath.Join(f.homeDir, ".roundfix"), 0o700); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(f.homeDir, ".roundfix", "corrections"), []byte("blocks correction write"), 0o600)
		}}
		code := runWithContext(t.Context(), []string{"review", "dispose", "F1", "--fixed-by", fix, "--class", "unclear", "--reason", "a reason with $text and 'quotes'"}, out, &diagnostics, commandEnvironmentForTest(t))
		ledger := readSingleDispositionLine(t, f.artifactDir)
		if code != 1 || out.String() != ledger || !strings.Contains(diagnostics.String(), "disposition recorded but Correction could not be written") || !strings.Contains(diagnostics.String(), "roundfix corrections declare --spec 0412-example --finding F1@"+f.headCommit+" --commit "+fix) {
			t.Fatalf("%d %q %q", code, out.String(), diagnostics.String())
		}
		if !strings.Contains(diagnostics.String(), "--reason "+posixSingleQuote("a reason with $text and 'quotes'")) {
			t.Fatal(diagnostics.String())
		}
	})
	t.Run("ambiguous finding and explicit head", func(t *testing.T) {
		f, fix := correctionDispositionFixture(t, []string{"0412-example"})
		// Two independently fixed heads make the bare F1 ambiguous.
		path := filepath.Join(f.artifactDir, reviewDispositionLedgerFileName)
		entries := []reviewFindingDisposition{
			{Repository: f.repository, HeadCommit: f.headCommit, Finding: "F1", Disposition: "fixed", FixedBy: fix},
			{Repository: f.repository, HeadCommit: f.baseCommit, Finding: "F1", Disposition: "fixed", FixedBy: fix},
		}
		var ledger bytes.Buffer
		for _, d := range entries {
			if err := json.NewEncoder(&ledger).Encode(d); err != nil {
				t.Fatal(err)
			}
		}
		mustWrite(t, path, ledger.String())
		before := correctionTree(t, f.homeDir)
		code, out, diag := correctionCall(t, "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known", "--finding", "F1")
		if code != 2 || out != "" || !strings.Contains(diag, "F1@"+f.headCommit) || !strings.Contains(diag, "F1@"+f.baseCommit) {
			t.Fatalf("%d %q %q", code, out, diag)
		}
		if !reflect.DeepEqual(before, correctionTree(t, f.homeDir)) {
			t.Fatal("refusal wrote home")
		}
		code, _, diag = correctionCall(t, "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known", "--finding", "F1@"+f.headCommit)
		if code != 0 || diag != "" {
			t.Fatalf("%d %q", code, diag)
		}
	})
}
func TestReviewDisposeWithoutADeclarationKeepsItsBytes(t *testing.T) {
	t.Parallel()
	f, fix := correctionDispositionFixture(t, []string{"0412-example"})
	recordPath := filepath.Join(reviewCheckoutDir(f.artifactDir, f.repository), reviewRecordFileName)
	before, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	code, out, diag := runReviewDispose(t, "F1", "--fixed-by", fix)
	if code != 0 || diag != "" {
		t.Fatalf("%d %q", code, diag)
	}
	ledger := readSingleDispositionLine(t, f.artifactDir)
	d := decodeDispositionLine(t, ledger)
	if _, err := time.Parse(time.RFC3339, d.RecordedAt); err != nil {
		t.Fatal(err)
	}
	// Literal delivery-base key order, optional fields, newline and stdout.
	want := fmt.Sprintf(`{"repository":%q,"headCommit":%q,"finding":"F1","text":"internal/cli/review.go:10: finding","disposition":"fixed","fixedBy":%q,"recordedAt":%q}`+"\n", f.repository, f.headCommit, fix, d.RecordedAt)
	if ledger != want || out != want {
		t.Fatalf("ledger=%q stdout=%q want=%q", ledger, out, want)
	}
	after, err := os.ReadFile(recordPath)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("changed review record", err)
	}
	if _, err := os.Stat(filepath.Join(f.homeDir, ".roundfix", "corrections")); !os.IsNotExist(err) {
		t.Fatalf("read or created corrections: %v", err)
	}
}
func TestReviewDisposeRefusalLeavesTheLedgerAndTheRecord(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		flags []string
		specs []string
	}{
		{"class without reason", []string{"--class", "implementation"}, nil},
		{"reason without class", []string{"--reason", "known"}, nil},
		{"correction with class", []string{"--correction", "C1", "--class", "unclear", "--reason", "known"}, nil},
		{"correction with reason", []string{"--correction", "C1", "--reason", "known"}, nil},
		{"correction with artifact", []string{"--correction", "C1", "--artifact", "prd"}, nil},
		{"spec alone", []string{"--spec", "0412-example"}, nil},
		{"artifact alone", []string{"--artifact", "prd"}, nil},
		{"plan without artifact", []string{"--class", "plan", "--reason", "known"}, nil},
		{"unknown artifact", []string{"--class", "plan", "--artifact", "idea", "--reason", "known"}, nil},
		{"unknown class", []string{"--class", "unknown", "--reason", "known"}, nil},
		{"artifact on implementation", []string{"--class", "implementation", "--artifact", "prd", "--reason", "known"}, nil},
		{"artifact on not correction", []string{"--class", "not-a-correction", "--artifact", "prd", "--reason", "known"}, nil},
		{"blank reason", []string{"--class", "unclear", "--reason", " "}, nil},
		{"long reason", []string{"--class", "unclear", "--reason", strings.Repeat("é", 201)}, nil},
		{"newline reason", []string{"--class", "unclear", "--reason", "first\nsecond"}, nil},
		{"unknown correction", []string{"--correction", "C99"}, []string{"0412-example"}},
		{"no specs", []string{"--class", "unclear", "--reason", "known"}, []string{}},
		{"several specs", []string{"--class", "unclear", "--reason", "known"}, []string{"0412-example", "0413-other"}},
		{"absent spec", []string{"--class", "unclear", "--reason", "known", "--spec", "0413-absent"}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, fix := correctionDispositionFixture(t, c.specs)
			assertCorrectionCall(t, "Correction C1 of 0412-example: unclear; covers no Correction Sign.", "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known")
			beforeHome, beforeArtifacts := correctionTree(t, f.homeDir), correctionTree(t, f.artifactDir)
			args := append([]string{"F1", "--fixed-by", fix}, c.flags...)
			code, out, diag := runReviewDispose(t, args...)
			if code != 2 || out != "" || !strings.HasPrefix(diag, "roundfix: review dispose refused: ") {
				t.Fatalf("%d %q %q", code, out, diag)
			}
			if !reflect.DeepEqual(beforeHome, correctionTree(t, f.homeDir)) || !reflect.DeepEqual(beforeArtifacts, correctionTree(t, f.artifactDir)) {
				t.Fatal("refusal changed record or ledger")
			}
		})
	}
	t.Run("dismissal names the first actual flag", func(t *testing.T) {
		correctionsFixture(t)
		code, out, diag := runReviewDispose(t, "F1", "--dismiss", "--evidence", "--reason", "--class", "unclear", "--reason", "known")
		if code != 2 || out != "" || diag != "roundfix: review dispose refused: --class cannot be combined with --dismiss\n" {
			t.Fatalf("%d %q %q", code, out, diag)
		}
	})
	for _, name := range []string{"class", "reason", "artifact", "correction", "spec"} {
		t.Run("dismiss "+name, func(t *testing.T) {
			f, _ := correctionDispositionFixture(t, []string{"0412-example"})
			before := correctionTree(t, f.artifactDir)
			code, out, diag := runReviewDispose(t, "F1", "--dismiss", "--evidence", "not a defect", "--"+name, "value")
			if code != 2 || out != "" || diag != "roundfix: review dispose refused: --"+name+" cannot be combined with --dismiss\n" {
				t.Fatalf("%d %q %q", code, out, diag)
			}
			if !reflect.DeepEqual(before, correctionTree(t, f.artifactDir)) {
				t.Fatal("dismissal wrote ledger")
			}
		})
	}
}
