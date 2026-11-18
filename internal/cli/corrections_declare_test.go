package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/baseline"
	"roundfix/internal/correction"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func correctionsFixture(t *testing.T) (home, repo string) {
	t.Helper()
	home, repo = t.TempDir(), t.TempDir()
	gittest.Run(t, repo, "init", "-b", "main")
	mustMkdir(t, filepath.Join(repo, "docs/specs/0412-example"))
	mustWrite(t, filepath.Join(repo, "docs/specs/0412-example/_prd.md"), "# Fixture\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "root")
	physical, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	repo = physical
	setCommandEnvironmentForTest(t, home, repo)
	updateCommandDependenciesForTest(t, func(d *commandDependencies) {
		d.getenv = func(key string) string {
			if key == "CLAUDE_CODE_SESSION_ID" {
				return "correction-session"
			}
			return ""
		}
		d.phaseNow = func() time.Time { return time.Date(2026, 11, 20, 14, 2, 11, 0, time.UTC) }
	})
	return home, repo
}
func correctionCall(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, diagnostics bytes.Buffer
	code := runCLI(t, append([]string{"corrections"}, args...), &out, &diagnostics)
	return code, out.String(), diagnostics.String()
}
func correctionFixturePath(t *testing.T, home, repo string) string {
	t.Helper()
	identity, err := baseline.InspectRepositoryIdentity(t.Context(), repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	return correction.Path(home, identity.Digest, "0412-example")
}
func readCorrectionFixture(t *testing.T, home, repo string) correction.File {
	t.Helper()
	f, err := correction.Read(correctionFixturePath(t, home, repo))
	if err != nil {
		t.Fatal(err)
	}
	return f
}
func assertCorrectionCall(t *testing.T, want string, args ...string) {
	t.Helper()
	code, out, diag := correctionCall(t, args...)
	if code != 0 || out != want+"\n" || diag != "" {
		t.Fatalf("%v: exit=%d stdout=%q stderr=%q; want %q", args, code, out, diag, want)
	}
}
func correctionTree(t *testing.T, root string) map[string]string {
	t.Helper()
	result := make(map[string]string)
	err := filepath.WalkDir(root, func(p string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if e.IsDir() {
			result[rel] = "directory"
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		result[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestCorrectionsSurfaceTranscript2(t *testing.T) {
	t.Parallel()
	home, repo := correctionsFixture(t)
	assertCorrectionCall(t, "Correction C1 of 0412-example: implementation; covers no Correction Sign.", "declare", "--spec", "0412-example", "--class", "implementation", "--reason", "The note was cut in the middle of a multi-byte character")
	f := readCorrectionFixture(t, home, repo)
	if len(f.Declarations) != 1 || f.Declarations[0].Harness != "claude" || f.Declarations[0].Session != "correction-session" {
		t.Fatalf("%+v", f)
	}
	for _, p := range []string{correctionFixturePath(t, home, repo), filepath.Dir(correctionFixturePath(t, home, repo))} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		want := os.FileMode(0o600)
		if info.IsDir() {
			want = 0o700
		}
		if info.Mode().Perm() != want {
			t.Fatalf("mode %s = %v", p, info.Mode())
		}
	}
}
func TestCorrectionsSurfaceTranscript3(t *testing.T) {
	t.Parallel()
	home, _ := correctionsFixture(t)
	before := correctionTree(t, home)
	code, out, diag := correctionCall(t, "declare", "--spec", "0412-example", "--class", "plan", "--reason", "The TechSpec named a file that does not exist")
	want := "roundfix: corrections declare refused: class plan requires --artifact prd, techspec or tasks\nRun 'roundfix corrections --help' for usage.\n"
	if code != 2 || out != "" || diag != want {
		t.Fatalf("%d %q %q", code, out, diag)
	}
	if !reflect.DeepEqual(before, correctionTree(t, home)) {
		t.Fatal("refusal wrote home")
	}
}
func TestCorrectionsDeclareRefusalsWriteNothing(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"missing group", nil, "missing subcommand; use declare or report"},
		{"unknown group", []string{"unknown"}, `unknown subcommand "unknown"; use declare or report`},
		{"missing spec", []string{"declare", "--class", "unclear", "--reason", "known"}, "--spec must match"},
		{"slug", []string{"declare", "--spec", "../escape", "--class", "unclear", "--reason", "known"}, "--spec must match"},
		{"unknown class", []string{"--class", "wrong"}, "plan, implementation, environment, operator, unclear or not-a-correction"},
		{"blank class", []string{"--class", ""}, "plan, implementation, environment, operator, unclear or not-a-correction"},
		{"missing artifact", []string{"--class", "plan"}, "prd, techspec or tasks"},
		{"unknown artifact", []string{"--class", "plan", "--artifact", "idea"}, "prd, techspec or tasks"},
		{"unexpected artifact", []string{"--artifact", "prd"}, "only with class plan"},
		{"not correction artifact", []string{"--class", "not-a-correction", "--artifact", "prd"}, "only with class plan"},
		{"blank reason", []string{"--reason", " "}, "1 to 200 characters"},
		{"long reason", []string{"--reason", strings.Repeat("é", 201)}, "1 to 200 characters"},
		{"line reason", []string{"--reason", "one\ntwo"}, "one line"},
		{"control reason", []string{"--reason", "one\x00two"}, "control character"},
		{"separator reason", []string{"--reason", "one\u2028two"}, "line or paragraph separator"},
		{"paragraph reason", []string{"--reason", "one\u2029two"}, "line or paragraph separator"},
		{"blank by", []string{"--by", " "}, "1 to 64 characters"},
		{"long by", []string{"--by", strings.Repeat("é", 65)}, "1 to 64 characters"},
		{"control by", []string{"--by", "one\ttwo"}, "control character"},
		{"task", []string{"--task", "task_bad"}, "task_<digits>"},
		{"id", []string{"--correction", "C0"}, "C<n>"},
		{"missing correction", []string{"--correction", "C99"}, "does not exist"},
		{"commit", []string{"--commit", "missing"}, "does not resolve"},
		{"finding", []string{"--finding", "F1"}, "one fixed disposition"},
		{"unknown flag", []string{"--unknown"}, "flag provided but not defined"},
		{"not held", []string{"--spec", "0413-absent"}, "holds no Spec"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			home, repo := correctionsFixture(t)
			args := c.args
			if c.name != "missing group" && c.name != "unknown group" && c.name != "missing spec" && c.name != "slug" {
				args = append([]string{"declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known"}, args...)
			}
			beforeHome, beforeRepo := correctionTree(t, home), correctionTree(t, repo)
			code, out, diag := correctionCall(t, args...)
			if code != 2 || out != "" || !strings.Contains(diag, c.want) || !strings.HasSuffix(diag, "Run 'roundfix corrections --help' for usage.\n") {
				t.Fatalf("%v: %d %q %q", args, code, out, diag)
			}
			if !reflect.DeepEqual(beforeHome, correctionTree(t, home)) || !reflect.DeepEqual(beforeRepo, correctionTree(t, repo)) {
				t.Fatal("refusal wrote files")
			}
		})
	}
	t.Run("help", func(t *testing.T) {
		home, _ := correctionsFixture(t)
		before := correctionTree(t, home)
		for _, args := range [][]string{{"--help"}, {"-h"}, {"declare", "--help"}, {"declare", "-h"}} {
			code, out, diag := correctionCall(t, args...)
			if code != 0 || out != correctionsUsage || diag != "" {
				t.Fatalf("%d %q %q", code, out, diag)
			}
		}
		if !reflect.DeepEqual(before, correctionTree(t, home)) {
			t.Fatal("help wrote")
		}
	})
	t.Run("outside a Git working tree", func(t *testing.T) {
		home := t.TempDir()
		outside := t.TempDir()
		setCommandEnvironmentForTest(t, home, outside)
		before := correctionTree(t, home)
		code, out, diag := correctionCall(t, "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known")
		if code != 2 || out != "" || !strings.Contains(diag, "requires a git repository working tree") {
			t.Fatalf("%d %q %q", code, out, diag)
		}
		if !reflect.DeepEqual(before, correctionTree(t, home)) {
			t.Fatal("non-repository refusal wrote")
		}
	})
	t.Run("syntax before reads", func(t *testing.T) {
		home, repo := correctionsFixture(t)
		mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "malformed: [")
		before := correctionTree(t, home)
		code, _, diag := correctionCall(t, "declare", "--spec", "0412-example", "--class", "plan", "--reason", "known")
		if code != 2 || !strings.Contains(diag, "class plan requires") {
			t.Fatalf("%d %s", code, diag)
		}
		if !reflect.DeepEqual(before, correctionTree(t, home)) {
			t.Fatal("refusal wrote")
		}
	})
	t.Run("record read and write failures", func(t *testing.T) {
		for _, blockRecord := range []bool{true, false} {
			t.Run(fmt.Sprint(blockRecord), func(t *testing.T) {
				home, repo := correctionsFixture(t)
				path := correctionFixturePath(t, home, repo)
				if blockRecord {
					if err := os.MkdirAll(path, 0o700); err != nil {
						t.Fatal(err)
					}
				} else {
					mustMkdir(t, filepath.Join(home, ".roundfix"))
					mustWrite(t, filepath.Join(home, ".roundfix", "corrections"), "blocks directory")
				}
				before := correctionTree(t, home)
				code, out, _ := correctionCall(t, "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known")
				if code != 1 || out != "" {
					t.Fatalf("record failure exit=%d out=%q", code, out)
				}
				if !reflect.DeepEqual(before, correctionTree(t, home)) {
					t.Fatal("failure wrote")
				}
			})
		}
	})
}
func TestCorrectionsDeclareAgainAndWithdraw(t *testing.T) {
	t.Parallel()
	home, repo := correctionsFixture(t)
	assertCorrectionCall(t, "Correction C1 of 0412-example: implementation; covers commit "+strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "--short=7", "HEAD"))+".", "declare", "--spec", "0412-example", "--class", "implementation", "--reason", "known", "--commit", "HEAD")
	assertCorrectionCall(t, "Correction C1 of 0412-example: unclear; covers commit "+strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "--short=7", "HEAD"))+"; declared 2 times, class changed from implementation.", "declare", "--spec", "0412-example", "--correction", "C1", "--class", "unclear")
	assertCorrectionCall(t, "Correction C1 of 0412-example: unclear; covers commit "+strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "--short=7", "HEAD"))+" and Corrective Task task_02; declared 3 times, class unchanged.", "declare", "--spec", "0412-example", "--correction", "C1", "--task", "task_02")
	assertCorrectionCall(t, "Correction C1 of 0412-example: not a correction; covers commit "+strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "--short=7", "HEAD"))+" and Corrective Task task_02; declared 4 times, class changed from unclear.", "declare", "--spec", "0412-example", "--correction", "C1", "--class", "not-a-correction", "--reason", "maintainer decision")
	path := correctionFixturePath(t, home, repo)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	assertCorrectionCall(t, "Correction C1 of 0412-example: withdrawn; declared 5 times.", "declare", "--spec", "0412-example", "--correction", "C1", "--withdraw", "--reason", "declared in error")
	after, err := os.ReadFile(path)
	if err != nil || !bytes.HasPrefix(after, before) {
		t.Fatal("rewrote history", err)
	}
	f := readCorrectionFixture(t, home, repo)
	if len(f.Declarations) != 5 || len(correction.Standing(f)[0].Signs) != 0 {
		t.Fatalf("%+v", f)
	}
	assertCorrectionCall(t, "Correction C1 of 0412-example: implementation; covers Corrective Task task_03; declared 6 times, class changed from not-a-correction.", "declare", "--spec", "0412-example", "--correction", "C1", "--class", "implementation", "--task", "task_03")
	for _, class := range correction.Answers() {
		args := []string{"declare", "--spec", "0412-example", "--class", string(class), "--reason", "known"}
		if class == "plan" {
			args = append(args, "--artifact", "prd")
		}
		code, _, diag := correctionCall(t, args...)
		if code != 0 || diag != "" {
			t.Fatalf("%s: %d %s", class, code, diag)
		}
	}
	// A plan class can change when its retained artifact is explicitly cleared.
	code, out, diag := correctionCall(t, "declare", "--spec", "0412-example", "--correction", "C2", "--class", "implementation", "--artifact", "")
	if code != 0 || diag != "" || !strings.Contains(out, "class changed from plan") {
		t.Fatalf("%d %q %q", code, out, diag)
	}
	// Withdrawal and redeclaration refusals preserve every line and the ledger.
	for _, args := range [][]string{
		{"--correction", "C1"}, {"--correction", "C1", "--by", "name"},
		{"--withdraw", "--reason", "known"},
		{"--correction", "C1", "--withdraw"},
		{"--correction", "C1", "--withdraw", "--reason", "known", "--class", "unclear"},
		{"--correction", "C1", "--withdraw", "--reason", "known", "--artifact", "prd"},
		{"--correction", "C1", "--withdraw", "--reason", "known", "--task", "task_01"},
		{"--correction", "C1", "--withdraw", "--reason", "known", "--commit", "HEAD"},
		{"--correction", "C1", "--withdraw", "--reason", "known", "--finding", "F1"},
	} {
		before := correctionTree(t, home)
		code, out, _ := correctionCall(t, append([]string{"declare", "--spec", "0412-example"}, args...)...)
		if code != 2 || out != "" || !reflect.DeepEqual(before, correctionTree(t, home)) {
			t.Fatalf("refusal %v: %d %q", args, code, out)
		}
	}
}
func TestCorrectionsDeclareReachesOneRecordFromAWorktreeAndASecondClone(t *testing.T) {
	t.Parallel()
	home, repo := correctionsFixture(t)
	worktree := filepath.Join(t.TempDir(), "worktree")
	clone := filepath.Join(t.TempDir(), "clone")
	gittest.Run(t, repo, "worktree", "add", "--detach", worktree, "HEAD")
	gittest.Run(t, repo, "clone", "--no-hardlinks", repo, clone)
	for i, dir := range []string{repo, worktree, clone} {
		setCommandEnvironmentForTest(t, home, dir)
		assertCorrectionCall(t, fmt.Sprintf("Correction C%d of 0412-example: unclear; covers no Correction Sign.", i+1), "declare", "--spec", "0412-example", "--class", "unclear", "--reason", "known")
	}
	if f := readCorrectionFixture(t, home, repo); len(f.Declarations) != 3 {
		t.Fatalf("%+v", f)
	}
	matches, err := filepath.Glob(filepath.Join(home, ".roundfix", "corrections", "*", "*.jsonl"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("%v %v", matches, err)
	}
	// An Archive Record in the working tree is also an accepted Spec location.
	archive := filepath.Join(clone, "docs/history/specs/0413-archived.md")
	mustMkdir(t, filepath.Dir(archive))
	mustWrite(t, archive, "# Archive Record\n")
	code, out, diag := correctionCall(t, "declare", "--spec", "0413-archived", "--class", "operator", "--reason", "known")
	if code != 0 || diag != "" || !strings.Contains(out, "C1 of 0413-archived") {
		t.Fatalf("%d %q %q", code, out, diag)
	}
}
func TestGCRemovesCorrectionsRecordsPastTheCutoff(t *testing.T) {
	t.Parallel()
	for _, database := range []bool{false, true} {
		t.Run(fmt.Sprint(database), func(t *testing.T) {
			home, repo := correctionsFixture(t)
			now := time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)
			withGCNow(t, now)
			if database {
				s, err := store.Open(t.Context(), home)
				if err != nil {
					t.Fatal(err)
				}
				if err := s.Close(); err != nil {
					t.Fatal(err)
				}
			}
			path := correctionFixturePath(t, home, repo)
			if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
				t.Fatal(err)
			}
			line, err := correction.EncodeDeclaration("0412-example", correction.Declaration{Time: now.Add(-31 * 24 * time.Hour), Correction: "C1", Class: "unclear", Reason: "known"})
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, line, 0o600); err != nil {
				t.Fatal(err)
			}
			for _, dry := range []bool{true, false} {
				args := []string{"gc"}
				action := "removed"
				if dry {
					args = append(args, "--dry-run")
					action = "removable"
				}
				var out, diag bytes.Buffer
				code := runCLI(t, args, &out, &diag)
				if code != 0 || diag.Len() != 0 || !strings.Contains(out.String(), "  Phase Record day files "+action+": 0\n  Corrections records "+action+": 1\n") {
					t.Fatalf("%d %s %s", code, out.String(), diag.String())
				}
				_, err := os.Stat(path)
				if dry && err != nil {
					t.Fatal(err)
				}
				if !dry && !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("record remains: %v", err)
				}
			}
		})
	}
}

func TestConcurrentCorrectionDeclarationsMintDistinctIDs(t *testing.T) {
	t.Parallel()
	home, repo := correctionsFixture(t)
	environment := commandEnvironmentForTest(t)
	request, err := parseCorrectionDeclare([]string{"--spec", "0412-example", "--class", "unclear", "--reason", "known"})
	if err != nil {
		t.Fatal(err)
	}
	// Prepare two writers before either appends, reproducing the minting race.
	identity, err := baseline.InspectRepositoryIdentity(t.Context(), repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	p := preparedCorrection{path: correction.Path(home, identity.Digest, request.slug), slug: request.slug, request: request}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Go(func() {
			_, err := appendCorrection(context.Background(), p, environment.dependencies.phaseNow())
			errs <- err
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	f := readCorrectionFixture(t, home, repo)
	if len(f.Declarations) != 2 || f.Declarations[0].Correction != "C1" || f.Declarations[1].Correction != "C2" {
		t.Fatalf("%+v", f)
	}
}
