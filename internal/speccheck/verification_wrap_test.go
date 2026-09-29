// Suite: wrap-fragile Task Verification classification
// Invariant: Pending non-QA Tasks cannot use line-bound multi-word grep patterns against Markdown.
// Boundary IN: parsed Task Verification commands and the public Spec Consistency Check
// Boundary OUT: Daemon Verification execution and general shell validation
package speccheck_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const (
	wrapFragilePhrase = "records no QA row"
	wrapFragileFile   = "docs/user-guide/commands.md"
)

func TestWrapFragileVerificationReportsALineBoundPhraseAgainstMarkdown(t *testing.T) {
	t.Parallel()

	findings := wrapFragileFindings(`grep -q "records no QA row" docs/user-guide/commands.md`)
	if len(findings) != 1 {
		t.Fatalf("WrapFragileVerification() = %#v, want one finding", findings)
	}
	if findings[0].Code != speccheck.CodeVerifyWrapFragile || findings[0].Severity != speccheck.SeverityError {
		t.Fatalf("finding identity = %s/%s, want %s/%s", findings[0].Code, findings[0].Severity, speccheck.CodeVerifyWrapFragile, speccheck.SeverityError)
	}
}

func TestWrapFragileVerificationNamesThePhraseTheFileAndTheFix(t *testing.T) {
	t.Parallel()

	findings := wrapFragileFindings(`grep -q "records no QA row" docs/user-guide/commands.md`)
	if len(findings) != 1 {
		t.Fatalf("WrapFragileVerification() = %#v, want one finding", findings)
	}
	finding := findings[0]
	for _, want := range []string{wrapFragilePhrase, wrapFragileFile} {
		if !strings.Contains(finding.Summary, want) {
			t.Errorf("finding summary = %q, want %q", finding.Summary, want)
		}
	}
	const wantFix = `tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "records no QA row" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands.md "records no QA row" >&2; exit 1; }`
	if finding.Fix != wantFix {
		t.Errorf("finding fix = %q, want %q", finding.Fix, wantFix)
	}
}

func TestWrapFragileVerificationReportsANegatedLineBoundPhrase(t *testing.T) {
	t.Parallel()

	findings := wrapFragileFindings(`! rtk grep -q "records no QA row" docs/user-guide/commands.md`)
	if len(findings) != 1 {
		t.Fatalf("WrapFragileVerification() = %#v, want one finding", findings)
	}
	const wantFix = `! { tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "records no QA row"; }`
	if findings[0].Fix != wantFix {
		t.Errorf("finding fix = %q, want %q", findings[0].Fix, wantFix)
	}
}

func TestWrapFragileVerificationReportsEachTopLevelCommand(t *testing.T) {
	t.Parallel()

	command := `grep -q "first wrapped phrase" docs/first.md && rtk grep -q "second wrapped phrase" docs/second.md || ! grep -q "third wrapped phrase" docs/third.md; grep -q "fourth wrapped phrase" docs/fourth.md | cat`
	findings := wrapFragileFindings(command)
	if len(findings) != 4 {
		t.Fatalf("WrapFragileVerification() = %#v, want one finding per top-level grep", findings)
	}
}

func TestWrapFragileVerificationUsesTheFirstRegexpOption(t *testing.T) {
	t.Parallel()

	findings := wrapFragileFindings(`grep -q -e "records no QA row" -e fallback docs/user-guide/commands.md`)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, wrapFragilePhrase) {
		t.Fatalf("WrapFragileVerification() = %#v, want first -e phrase reported", findings)
	}
}

func TestWrapFragileVerificationAcceptsTheWrapTolerantForm(t *testing.T) {
	t.Parallel()

	command := `tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "records no QA row" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands.md "records no QA row" >&2; exit 1; }`
	requireNoWrapFragileFinding(t, command)
}

func TestWrapFragileVerificationAcceptsAnAnchoredPattern(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "^records no QA row" docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsAnEndAnchoredPattern(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "records no QA row$" docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsAHeadingPattern(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "# Records no QA row" docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsATableRowPattern(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "| records no QA row" docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsASingleWord(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q records docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsStandardInput(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `printf '%s\n' "records no QA row" | grep -q "records no QA row"`)
}

func TestWrapFragileVerificationAcceptsAnExplicitStandardInputOperand(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "records no QA row" - docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsRedirectedStandardInput(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "records no QA row" < docs/user-guide/commands.md`)
}

func TestWrapFragileVerificationAcceptsANonMarkdownFile(t *testing.T) {
	t.Parallel()

	requireNoWrapFragileFinding(t, `grep -q "records no QA row" internal/docscontract/testdata/corpus-golden.json`)
}

func TestWrapFragileCheckSkipsACompletedTask(t *testing.T) {
	t.Parallel()

	result := checkWrapFragileFixture(t, spec.StatusCompleted, spec.TaskTypeBackend)
	findings := findingsWithCode(result, speccheck.CodeVerifyWrapFragile)
	if len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want completed Task skipped", speccheck.CodeVerifyWrapFragile, findings)
	}
}

func TestWrapFragileCheckSkipsTheQATask(t *testing.T) {
	t.Parallel()

	result := checkWrapFragileFixture(t, spec.StatusPending, spec.TaskTypeQA)
	findings := findingsWithCode(result, speccheck.CodeVerifyWrapFragile)
	if len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want qa Task skipped", speccheck.CodeVerifyWrapFragile, findings)
	}
}

func TestStrictKeepsAWrapFragilePhraseAnError(t *testing.T) {
	t.Parallel()

	result := checkWrapFragileFixture(t, spec.StatusPending, spec.TaskTypeBackend)
	findings := findingsWithCode(result, speccheck.CodeVerifyWrapFragile)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeVerifyWrapFragile, findings)
	}
	if findings[0].Severity != speccheck.SeverityError {
		t.Fatalf("severity = %q, want %q", findings[0].Severity, speccheck.SeverityError)
	}
	if len(findings[0].Where) != 1 || findings[0].Where[0].Path != "docs/specs/wrap-fragile/task_01.md" || findings[0].Where[0].Line != 13 {
		t.Errorf("finding locations = %#v, want Verification command at task_01.md:13", findings[0].Where)
	}

	speccheck.PromoteGaps(&result)
	findings = findingsWithCode(result, speccheck.CodeVerifyWrapFragile)
	if len(findings) != 1 || findings[0].Severity != speccheck.SeverityError {
		t.Fatalf("strict %s findings = %#v, want one error", speccheck.CodeVerifyWrapFragile, findings)
	}
}

func TestMissingTaskGraphListsTheWrapFragileSkip(t *testing.T) {
	t.Parallel()

	result, err := speccheck.Check(fixtureSpecRoot, "testdata/repo", "no-taskgraph")
	if err != nil {
		t.Fatalf("Check(no-taskgraph): %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeVerifyWrapFragile); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want detector skipped", speccheck.CodeVerifyWrapFragile, findings)
	}
	if !hasSkip(result, speccheck.CodeVerifyWrapFragile, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeVerifyWrapFragile)
	}
}

func TestWrapFragileVerificationRegistersAtTasksStage(t *testing.T) {
	t.Parallel()

	repoRoot, specsRoot, slug := writeStageScopeFixture(t)
	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
	if err != nil {
		t.Fatalf("CheckStage(StagePRD): %v", err)
	}
	if !hasSkip(result, speccheck.CodeVerifyWrapFragile, "stage prd") {
		t.Fatalf("StagePRD Skipped = %#v, want named %s detector", result.Skipped, speccheck.CodeVerifyWrapFragile)
	}
}

func TestWrapTolerantFormPassesOnAWrappedPhrase(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "wrapped.md")
	if err := os.WriteFile(path, []byte("records no\nQA row\n"), 0o644); err != nil {
		t.Fatalf("write wrapped Markdown: %v", err)
	}
	fix := wrapFragilePresenceFix(t, path)
	if output, err := exec.Command("sh", "-c", fix).CombinedOutput(); err != nil {
		t.Fatalf("sh -c suggested presence form: %v; output=%q", err, output)
	}
}

func TestWrapTolerantFormNamesTheMissingPhraseAndFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing.md")
	if err := os.WriteFile(path, []byte("another paragraph\n"), 0o644); err != nil {
		t.Fatalf("write missing-phrase Markdown: %v", err)
	}
	fix := wrapFragilePresenceFix(t, path)
	output, err := exec.Command("sh", "-c", fix).CombinedOutput()
	if err == nil {
		t.Fatalf("sh -c suggested presence form exited 0, want 1; output=%q", output)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("sh -c suggested presence form error = %v, want exit 1; output=%q", err, output)
	}
	for _, want := range []string{path, wrapFragilePhrase} {
		if !strings.Contains(string(output), want) {
			t.Errorf("stderr = %q, want %q", output, want)
		}
	}
}

func wrapFragileFindings(command string) []speccheck.Finding {
	return speccheck.WrapFragileVerification(spec.Task{
		File:         "fixture/task_01.md",
		Status:       spec.StatusPending,
		Type:         spec.TaskTypeBackend,
		Verification: []string{command},
	})
}

func requireNoWrapFragileFinding(t *testing.T, command string) {
	t.Helper()
	if findings := wrapFragileFindings(command); len(findings) != 0 {
		t.Fatalf("WrapFragileVerification(%q) = %#v, want no finding", command, findings)
	}
}

func wrapFragilePresenceFix(t *testing.T, path string) string {
	t.Helper()
	findings := wrapFragileFindings(fmt.Sprintf(`grep -q "records no QA row" %s`, path))
	if len(findings) != 1 {
		t.Fatalf("WrapFragileVerification() = %#v, want suggested presence form", findings)
	}
	return findings[0].Fix
}

func checkWrapFragileFixture(t *testing.T, status spec.Status, taskType spec.TaskType) speccheck.Result {
	t.Helper()

	repoRoot := t.TempDir()
	const slug = "wrap-fragile"
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: wrap-fragile
status: active
created: 2026-09-28
surfaces: [backend]
---

# Wrap-fragile fixture
`)
	qaDeclaration := "qa: declined\nqa_reason: fixture has no terminal QA gate"
	if taskType == spec.TaskTypeQA {
		qaDeclaration = "qa: task_01"
	}
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", fmt.Sprintf(`---
schema: spec-tasks/v1
spec: wrap-fragile
%s
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`, qaDeclaration))
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/task_01.md", fmt.Sprintf(`---
task: task_01
spec: wrap-fragile
status: %s
type: %s
complexity: low
---

# Task 01: Check a phrase

## Verification

- `+"`grep -q \"records no QA row\" docs/user-guide/commands.md`"+`
`, status, taskType))

	result, err := speccheck.Check(specsRoot, repoRoot, slug)
	if err != nil {
		t.Fatalf("Check(): %v", err)
	}
	return result
}
