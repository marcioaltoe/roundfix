package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

// Suite: reopen over Git-proven Late Dependencies.
// Boundary IN: built CLI, working-tree graph and report-addition history.
// Boundary OUT: QA evidence and non-gate Task bytes remain unchanged.
func TestReopenReopensALateDependencyThroughTheBuiltBinary(t *testing.T) {
	home, repo, report := lateDependencyWorkspace(t, true, true)
	binary := filepath.Join(t.TempDir(), "roundfix")
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = root
	build.Env = isolatedGitEnvForTest()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	beforeReport := mustRead(t, report)
	beforeDependency := mustRead(t, implementTaskPath(repo, "task_03"))
	command := exec.Command(binary, "reopen", "--spec", implementTestSlug)
	command.Dir = repo
	command.Env = withEnvValue(isolatedGitEnvForTest(), "HOME", home)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	code := exitCodeFromWait(command.Run())
	if code != exitOK || stderr.Len() != 0 || stdout.String() != "reopened task_qa pending — invalidated qa/qa-report-2026-10-06.md\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	task := mustRead(t, implementTaskPath(repo, "task_qa"))
	if !strings.Contains(task, "status: pending") || !strings.Contains(task, "- Dependencies added after the QA Report: `task_03`\n") {
		t.Fatalf("QA Task: %s", task)
	}
	if mustRead(t, report) != beforeReport || mustRead(t, implementTaskPath(repo, "task_03")) != beforeDependency {
		t.Fatal("reopen changed evidence or dependency")
	}
	assertNoRunDatabase(t, home)
}

func TestReopenRefusesWhenTheNewestReportIsUncommitted(t *testing.T) {
	_, repo, report := lateDependencyWorkspace(t, true, true)
	newest := filepath.Join(filepath.Dir(report), "qa-report-2026-10-07.md")
	mustWrite(t, newest, mustRead(t, report))
	assertLateDependencyRefusal(t, repo, newest)
}

func TestReopenRefusesWhenTheRecordedClosureMatches(t *testing.T) {
	_, repo, report := lateDependencyWorkspace(t, false, true)
	assertLateDependencyRefusal(t, repo, report)
}

func lateDependencyWorkspace(t *testing.T, addLate, commitReport bool) (string, string, string) {
	t.Helper()
	seeds := []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}, implementQAGateSeed(string(spec.StatusCompleted), "task_01")}
	home, repo := newImplementWorkspace(t, seeds)
	report := filepath.Join(repo, "docs", "specs", implementTestSlug, "qa", "qa-report-2026-10-06.md")
	mustMkdir(t, filepath.Dir(report))
	mustWrite(t, report, "---\nverdict: pass\n---\n\nPrior evidence.\n")
	if commitReport {
		gitImplement(t, repo, "add", "-A")
		gitImplement(t, repo, "commit", "-m", "record QA report")
	}
	if addLate {
		writeImplementSpec(t, repo, implementTestSlug, []implementSeed{seeds[0], {id: "task_03", status: string(spec.StatusCompleted)}, implementQAGateSeed(string(spec.StatusCompleted), "task_01", "task_03")})
	}
	return home, repo, report
}

func assertLateDependencyRefusal(t *testing.T, repo, report string) {
	t.Helper()
	before := mustRead(t, implementTaskPath(repo, "task_qa"))
	evidence := mustRead(t, report)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"reopen", "--spec", implementTestSlug}, &stdout, &stderr)
	want := "terminal QA Task \"task_qa\" is not stale; every dependency is completed"
	if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), want) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if mustRead(t, implementTaskPath(repo, "task_qa")) != before || mustRead(t, report) != evidence {
		t.Fatal("refusal mutated Task or report")
	}
}

func TestReopenRefusesWhenLateDependencyIDsChangeBeforeWrite(t *testing.T) {
	_, repo, _ := lateDependencyWorkspace(t, true, true)
	var preflightStderr bytes.Buffer
	plan, err := preflightReopen(context.Background(), implementTestSlug, &preflightStderr, commandEnvironmentForTest(t))
	if err != nil {
		t.Fatalf("preflight: %v", err)
	}
	writeImplementSpec(t, repo, implementTestSlug, []implementSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
		{id: "task_03", status: string(spec.StatusCompleted)},
		{id: "task_04", status: string(spec.StatusCompleted)},
		implementQAGateSeed(string(spec.StatusCompleted), "task_01", "task_03", "task_04"),
	})
	before := mustRead(t, implementTaskPath(repo, "task_qa"))
	var stdout, stderr bytes.Buffer
	code := reopenFromPlan(context.Background(), implementTestSlug, plan, &stdout, &stderr)
	if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), "gate changed after preflight") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if mustRead(t, implementTaskPath(repo, "task_qa")) != before {
		t.Fatal("changed gate was mutated")
	}
}

func TestReopenUsesOldestReportAddition(t *testing.T) {
	_, repo, report := lateDependencyWorkspace(t, true, true)
	before := mustRead(t, report)
	if err := os.Remove(report); err != nil {
		t.Fatal(err)
	}
	gitImplement(t, repo, "add", "-A")
	gitImplement(t, repo, "commit", "-m", "remove report after adding dependency")
	mustWrite(t, report, before)
	gitImplement(t, repo, "add", "-A")
	gitImplement(t, repo, "commit", "-m", "restore report")
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"reopen", "--spec", implementTestSlug}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(mustRead(t, implementTaskPath(repo, "task_qa")), "- Dependencies added after the QA Report: `task_03`") {
		t.Fatal("oldest report addition was not used")
	}
}

func TestReopenRefusesWithoutRecordedManifest(t *testing.T) {
	for _, tc := range []struct {
		name      string
		malformed bool
	}{{"absent", false}, {"unparsable", true}} {
		t.Run(tc.name, func(t *testing.T) {
			_, repo, report := lateDependencyWorkspace(t, false, false)
			manifest := filepath.Join(repo, "docs", "specs", implementTestSlug, "_tasks.md")
			if tc.malformed {
				mustWrite(t, manifest, "not a manifest")
			} else if err := os.Remove(manifest); err != nil {
				t.Fatal(err)
			}
			gitImplement(t, repo, "add", "-A")
			gitImplement(t, repo, "commit", "-m", "record report without usable manifest")
			writeImplementSpec(t, repo, implementTestSlug, []implementSeed{
				{id: "task_01", status: string(spec.StatusCompleted)},
				{id: "task_03", status: string(spec.StatusCompleted)},
				implementQAGateSeed(string(spec.StatusCompleted), "task_01", "task_03"),
			})
			assertLateDependencyRefusal(t, repo, report)
		})
	}
}

func TestReopenRefusesWhenReportResolvesOutsideGitRoot(t *testing.T) {
	_, repo, report := lateDependencyWorkspace(t, true, true)
	outside := filepath.Join(t.TempDir(), "qa-report-2026-10-06.md")
	mustWrite(t, outside, mustRead(t, report))
	if err := os.Remove(report); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, report); err != nil {
		t.Fatal(err)
	}
	assertLateDependencyRefusal(t, repo, report)
}
