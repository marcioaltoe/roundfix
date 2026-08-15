// Suite: Spec commands
// Invariant: the CLI exposes Spec checks and close audits through stable streams, schemas, and exit codes.
// Boundary IN: public Run dispatch, configured Spec Root resolution, command renderers, and real Git and Run Database fixtures
// Boundary OUT: checker and audit classification correctness, owned by internal/speccheck and internal/specaudit
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/specaudit"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

func TestRunSpecCheckCleanText(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "clean"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), "Spec clean\nNo findings.\n") {
		t.Fatalf("stdout = %q, want clean report", stdout.String())
	}
	if strings.Contains(stdout.String(), "[error]") || strings.Contains(stdout.String(), "[gap]") {
		t.Fatalf("clean stdout contains a finding:\n%s", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestSpecCheckRunVerification(t *testing.T) {
	t.Run("executes commands for every requested Spec", func(t *testing.T) {
		_, _ = newSpecCheckWorkspace(t, "coverage-range", "coverage-untasked")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"spec", "check", "coverage-range", "coverage-untasked", "--run-verification",
		}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("multi-Spec run-verification exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
		}
		if got := strings.Count(stdout.String(), "Verification tree: HEAD"); got != 2 {
			t.Fatalf("HEAD report count = %d, want 2:\n%s", got, stdout.String())
		}
		if got := strings.Count(stdout.String(), `- task_01: vacuous — "true"`); got != 2 {
			t.Fatalf("executed command report count = %d, want 2:\n%s", got, stdout.String())
		}
	})

	t.Run("reports vacuous and honest commands against HEAD", func(t *testing.T) {
		_, repoDir := newSpecCheckVerificationWorkspace(t, []string{
			"test -f committed-marker.txt",
			"test -f task-output.txt",
		})
		mustWrite(t, filepath.Join(repoDir, "committed-marker.txt"), "present at HEAD\n")
		gitImplement(t, repoDir, "add", "committed-marker.txt")
		gitImplement(t, repoDir, "commit", "-m", "test: seed pre-work marker")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"spec", "check", "clean", "--run-verification",
		}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("run-verification exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
		}
		for _, want := range []string{
			"Verification tree: HEAD",
			`- task_01: vacuous — "test -f committed-marker.txt" (exited zero before work)`,
			`- task_01: honest — "test -f task-output.txt" (exited non-zero before work)`,
		} {
			if !strings.Contains(stdout.String(), want) {
				t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
			}
		}
		if stderr.String() != "" {
			t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
		}
	})

	t.Run("exits zero when every command honestly fails", func(t *testing.T) {
		_, _ = newSpecCheckVerificationWorkspace(t, []string{
			"test -f task-output-one.txt",
			"test -f task-output-two.txt",
		})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"spec", "check", "clean", "--run-verification",
		}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("honest run-verification exit = %d, want %d; stderr=%q\nstdout=%s", code, exitOK, stderr.String(), stdout.String())
		}
		if got := strings.Count(stdout.String(), ": honest — "); got != 2 {
			t.Fatalf("honest verdict count = %d, want 2:\n%s", got, stdout.String())
		}
	})

	t.Run("does not execute commands without the flag", func(t *testing.T) {
		marker := filepath.Join(t.TempDir(), "verification-ran")
		_, _ = newSpecCheckVerificationWorkspace(t, []string{"touch " + marker})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{"spec", "check", "clean"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
		}
		if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Verification command ran without opt-in: %v", err)
		}
		if !strings.Contains(stdout.String(), "Verification: not run (use --run-verification).") {
			t.Fatalf("stdout does not report unexecuted Verification:\n%s", stdout.String())
		}
	})

	t.Run("reports a command that cannot run as unknown", func(t *testing.T) {
		const command = "roundfix-verification-binary-that-does-not-exist"
		_, _ = newSpecCheckVerificationWorkspace(t, []string{command})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"spec", "check", "clean", "--run-verification",
		}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("unknown run-verification exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
		}
		want := `- task_01: unknown — "` + command + `" (command could not be executed (exit 127)`
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout does not report unknown command and cause %q:\n%s", want, stdout.String())
		}
		if strings.Contains(stdout.String(), ": vacuous —") || strings.Contains(stdout.String(), ": honest —") {
			t.Fatalf("unknown command received another verdict:\n%s", stdout.String())
		}
	})
}

func TestRunSpecCheckErrorText(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "tooling-unauthorized"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, want := range []string{
		"[error] " + speccheck.CodeToolingUnauthorized,
		"docs/specs/tooling-unauthorized/_prd.md:",
		"docs/workflow/authorizations/other-spec.md:",
		"  fix: ",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestSpecCheckStageExitsNonZeroOnAFinding(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "check", "tooling-unauthorized", "--stage", "prd",
	}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("stage-scoped spec check exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[error] "+speccheck.CodeToolingUnauthorized) {
		t.Fatalf("stage-scoped stdout does not contain the error finding:\n%s", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestSpecCheckStageExitsZeroWithoutAFinding(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "check", "clean", "--stage=prd",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("stage-scoped spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), "Spec clean\nNo findings.\n") {
		t.Fatalf("stdout = %q, want clean stage-scoped report", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestSpecCheckStageRejectsAnUnknownValue(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "check", "--stage", "design",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unknown stage exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no report", stdout.String())
	}
	for _, want := range []string{"design", "prd", "techspec", "tasks"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr does not contain %q: %q", want, stderr.String())
		}
	}
}

func TestSpecCheckWithoutStageIsUnchanged(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecCheckWorkspace(t, "tooling-unauthorized")
	wantResult, err := speccheck.Check(
		filepath.Join(repoDir, "docs", "specs"),
		repoDir,
		"tooling-unauthorized",
	)
	if err != nil {
		t.Fatalf("build unscoped comparison result: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "check", "tooling-unauthorized",
	}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("unscoped spec check exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	wantOutcome := specCheckOutcome{result: wantResult, repairInputs: []specCheckRepairInput{}}
	if want := renderSpecCheckText(wantOutcome, specCheckVerificationReport{Commands: []specCheckVerificationCommandReport{}}); stdout.String() != want {
		t.Fatalf("unscoped stdout = %q, want existing report %q", stdout.String(), want)
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecCheckGapStrictPromotion(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecCheckWorkspace(t, "citation-dirty")
	techSpecPath := filepath.Join(repoDir, "docs", "specs", "citation-dirty", "_techspec.md")
	techSpec := strings.ReplaceAll(mustRead(t, techSpecPath), "\n## Decisions\n\nThe implementation also follows ADR-0103.\n", "")
	mustWrite(t, techSpecPath, techSpec)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"spec", "check", "citation-dirty"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("gap-only spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[gap] "+speccheck.CodeADRRelated) {
		t.Fatalf("gap-only stdout does not contain gap finding:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "check", "citation-dirty", "--strict"}, &stdout, &stderr)
	if code != exitRunFailed {
		t.Fatalf("strict gap exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[error] "+speccheck.CodeADRRelated) {
		t.Fatalf("strict stdout does not contain promoted error:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "[gap] "+speccheck.CodeADRRelated) {
		t.Fatalf("strict stdout retained gap severity:\n%s", stdout.String())
	}
}

func TestSpecCheckClassifiesTheGateBoundary(t *testing.T) {
	t.Run("authoring stage keeps a declared term as an error", func(t *testing.T) {
		_, _ = newSpecCheckWorkspace(t, "vocabulary-missing")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"spec", "check", "vocabulary-missing", "--stage", "techspec", "--strict",
		}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("authoring-stage exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
		}
		if !strings.Contains(stdout.String(), "[error] "+speccheck.CodeVocabularyUndocumented) {
			t.Fatalf("authoring-stage stdout does not contain the vocabulary error:\n%s", stdout.String())
		}
		if strings.Contains(stdout.String(), "[repair input]") {
			t.Fatalf("authoring-stage stdout contains a gate repair input:\n%s", stdout.String())
		}
	})

	t.Run("strict full sweep reports a declared term as repair input", func(t *testing.T) {
		_, _ = newSpecCheckWorkspace(t, "vocabulary-missing")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"spec", "check", "vocabulary-missing", "--strict",
		}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("strict full-sweep exit = %d, want %d; stderr=%q\nstdout=%s", code, exitOK, stderr.String(), stdout.String())
		}
		for _, want := range []string{"[repair input] " + speccheck.CodeVocabularyUndocumented, "publish:"} {
			if !strings.Contains(stdout.String(), want) {
				t.Errorf("strict full-sweep stdout does not contain %q:\n%s", want, stdout.String())
			}
		}
		if strings.Contains(stdout.String(), "[error] "+speccheck.CodeVocabularyUndocumented) {
			t.Fatalf("strict full-sweep stdout retained the declared-term error:\n%s", stdout.String())
		}

		stdout.Reset()
		stderr.Reset()
		code = runCLIContext(t, context.Background(), []string{
			"spec", "check", "vocabulary-missing", "--strict", "--format", "json",
		}, &stdout, &stderr)
		if code != exitOK {
			t.Fatalf("strict full-sweep JSON exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
		}
		var document struct {
			Findings     []speccheck.Finding `json:"findings"`
			RepairInputs []struct {
				Code     string  `json:"code"`
				Severity *string `json:"severity"`
			} `json:"repairInputs"`
		}
		if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &document); err != nil {
			t.Fatalf("decode strict full-sweep JSON: %v; stdout=%q", err, stdout.String())
		}
		if len(document.Findings) != 0 || len(document.RepairInputs) != 1 ||
			document.RepairInputs[0].Code != speccheck.CodeVocabularyUndocumented ||
			document.RepairInputs[0].Severity != nil {
			t.Fatalf("strict full-sweep JSON = %#v, want one non-error repair input", document)
		}
	})

	t.Run("term with no Spec declaration remains an error in both modes", func(t *testing.T) {
		undeclared := speccheck.Result{
			Slug: "current-spec",
			Findings: []speccheck.Finding{{
				Code:     speccheck.CodeVocabularyUndocumented,
				Severity: speccheck.SeverityError,
				Summary:  `internal/example/emitter.go emits undocumented token "orphan:" absent from CONTEXT.md`,
				Where:    []speccheck.Location{{Path: "internal/example/emitter.go", Line: 7}},
				Fix:      `Document "orphan:" in CONTEXT.md.`,
			}},
		}

		for _, gateSweep := range []bool{false, true} {
			outcome := classifySpecCheckBoundary(undeclared, gateSweep)
			if len(outcome.result.Findings) != 1 || outcome.result.Findings[0].Severity != speccheck.SeverityError {
				t.Errorf("classifySpecCheckBoundary(gateSweep=%t) findings = %#v, want the undeclared-term error", gateSweep, outcome.result.Findings)
			}
			if len(outcome.repairInputs) != 0 {
				t.Errorf("classifySpecCheckBoundary(gateSweep=%t) repair inputs = %#v, want none", gateSweep, outcome.repairInputs)
			}
		}
	})
}

func TestRunSpecCheckJSONWritesOneObjectPerSpec(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean", "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "check", "clean", "tooling-unauthorized", "--format", "json",
	}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("multi-Spec JSON exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("JSON lines = %d, want 2; stdout=%q", len(lines), stdout.String())
	}
	for index, wantSlug := range []string{"clean", "tooling-unauthorized"} {
		var document struct {
			Schema       string `json:"schema"`
			Slug         string `json:"slug"`
			Verification struct {
				Executed bool  `json:"executed"`
				Commands []any `json:"commands"`
			} `json:"verification"`
		}
		if err := json.Unmarshal([]byte(lines[index]), &document); err != nil {
			t.Fatalf("line %d is not JSON: %v; line=%q", index+1, err, lines[index])
		}
		if document.Schema != speccheck.SchemaVersion || document.Slug != wantSlug {
			t.Errorf("line %d document = %#v, want schema %q slug %q", index+1, document, speccheck.SchemaVersion, wantSlug)
		}
		if document.Verification.Executed || document.Verification.Commands == nil || len(document.Verification.Commands) != 0 {
			t.Errorf("line %d Verification = %#v, want explicit not-run state", index+1, document.Verification)
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecCheckUnknownSlugIsUsageError(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "no-such-slug"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unknown slug exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no partial report", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no-such-slug") {
		t.Fatalf("stderr does not name unknown slug: %q", stderr.String())
	}
}

func TestRunSpecCheckWithoutSlugChecksEveryActiveSpec(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean", "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("all-active exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, slug := range []string{"clean", "tooling-unauthorized"} {
		count := 0
		for _, line := range strings.Split(stdout.String(), "\n") {
			if line == "Spec "+slug {
				count++
			}
		}
		if count != 1 {
			t.Errorf("stdout does not contain exactly one report for %q:\n%s", slug, stdout.String())
		}
	}
}

func TestRunSpecCheckUnreadableRootIsUsageError(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecCheckWorkspace(t, "clean")
	missingRoot := filepath.Join(repoDir, "missing-spec-root")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "specs:\n  root: "+missingRoot+"\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unreadable root exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no report", stdout.String())
	}
	if !strings.Contains(stderr.String(), missingRoot) || !strings.Contains(stderr.String(), "does not exist") {
		t.Fatalf("stderr does not name unreadable Spec Root: %q", stderr.String())
	}
}

func TestRunSpecCheckHelpAppearsInTopLevelUsageAndCommandList(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("top-level help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{
		"roundfix spec check [<slug> ...] [--format <text|json>] [--strict] [--run-verification]",
		"spec       Check Spec artifact consistency",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("top-level help does not contain %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "check", "--help"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("spec check help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{"--stage", "prd, techspec, or tasks", "--format", "--strict", "--run-verification", "HEAD", "0  no errors", "1  at least one error", "2  usage error"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("spec check help does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("help stderr = %q, want empty", stderr.String())
	}
}

func TestRunSpecAuditCleanText(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecAuditWorkspace(t)
	archiveSpecAuditFixture(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "audit", specAuditFixtureSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	want := "Spec audit " + specAuditFixtureSlug + "\nNo residue or undelivered work.\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want clean report %q", stdout.String(), want)
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecAuditResidueText(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecAuditWorkspace(t)
	const branch = "ma/spec-audit-residue"
	gitImplement(t, repoDir, "checkout", "-b", branch)
	mustWrite(t, filepath.Join(repoDir, "residue.txt"), "residue\n")
	gitImplement(t, repoDir, "add", "residue.txt")
	gitImplement(
		t,
		repoDir,
		"commit",
		"-m", "feat: add residue fixture",
		"-m", "Roundfix-Spec: "+specAuditFixtureSlug,
	)
	gitImplement(t, repoDir, "checkout", "main")
	gitImplement(t, repoDir, "merge", "--squash", branch)
	gitImplement(t, repoDir, "commit", "-m", "feat: merge residue fixture")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "audit", specAuditFixtureSlug}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, want := range []string{
		"residue branch " + branch,
		"evidence: survivor content is fully represented",
		"reclaim: git branch -D -- '" + branch + "'",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecAuditUndeliveredTextNamesHoldingBranch(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecAuditWorkspace(t)
	const branch = "ma/spec-audit-archive"
	gitImplement(t, repoDir, "checkout", "-b", branch)
	archivedPath := archiveTestPath(spec.ArchiveKindSpec, specAuditFixtureSlug)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(repoDir, archivedPath)), 0o755); err != nil {
		t.Fatalf("create archived Spec parent: %v", err)
	}
	gitImplement(t, repoDir, "mv", filepath.ToSlash(filepath.Join("docs", "specs", specAuditFixtureSlug)), archivedPath)
	gitImplement(
		t,
		repoDir,
		"commit",
		"-m", "docs: archive Spec audit fixture",
		"-m", "Roundfix-Spec: "+specAuditFixtureSlug,
	)
	gitImplement(t, repoDir, "checkout", "main")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "audit", specAuditFixtureSlug}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, want := range []string{archivedPath, "held by: " + branch} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecAuditJSONWritesOneObject(t *testing.T) {
	t.Parallel()
	_, _ = newSpecAuditWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "audit", specAuditFixtureSlug, "--format", "json",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	decoder := json.NewDecoder(strings.NewReader(stdout.String()))
	var document struct {
		Schema        string `json:"schema"`
		SchemaVersion string `json:"schemaVersion"`
		Type          string `json:"type"`
		OK            bool   `json:"ok"`
		Slug          string `json:"slug"`
	}
	if err := decoder.Decode(&document); err != nil {
		t.Fatalf("stdout is not JSON: %v; stdout=%q", err, stdout.String())
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("stdout contains more than one JSON object: %v; stdout=%q", err, stdout.String())
	}
	if document.Schema != specAuditSchemaVersion ||
		document.SchemaVersion != specAuditSchemaVersion ||
		document.Type != specAuditDocumentType ||
		!document.OK ||
		document.Slug != specAuditFixtureSlug {
		t.Fatalf(
			"document = %#v, want schema %q type %q ok true slug %q",
			document,
			specAuditSchemaVersion,
			specAuditDocumentType,
			specAuditFixtureSlug,
		)
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRenderSpecAuditJSONReportsAttention(t *testing.T) {
	t.Parallel()
	data, err := renderSpecAuditJSON(specaudit.Result{
		Slug: specAuditFixtureSlug,
		Survivors: []specaudit.Survivor{{
			Name: "ma/spec-audit-residue",
			Kind: specaudit.KindResidue,
		}},
	})
	if err != nil {
		t.Fatalf("render Spec audit JSON: %v", err)
	}
	var document struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatalf("decode Spec audit JSON: %v", err)
	}
	if document.OK {
		t.Fatal("attention-requiring Spec audit JSON ok = true, want false")
	}
}

func TestRunSpecAuditUnknownSlugIsUsageError(t *testing.T) {
	t.Parallel()
	_, _ = newSpecAuditWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "audit", "no-such-slug", "--format", "json",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unknown slug exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no partial JSON", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no-such-slug") {
		t.Fatalf("stderr does not name unknown slug: %q", stderr.String())
	}
}

func TestRunSpecAuditHelpAppearsInUsageAndCommandList(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("top-level help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{
		"roundfix spec audit <slug> [--format <text|json>]",
		"spec       Check Spec artifact consistency; audit Spec delivery",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("top-level help does not contain %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "audit", "--help"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("spec audit help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{"--format", "0  no residue", "1  residue or undelivered work", "2  usage error or unknown Spec slug"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("spec audit help does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("help stderr = %q, want empty", stderr.String())
	}
}

func TestRunSpecAuditPreservesSpecCheckBehavior(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "clean"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), "Spec clean\nNo findings.\n") {
		t.Fatalf("spec check stdout = %q, want unchanged clean report prefix", stdout.String())
	}
	if strings.Contains(stdout.String(), "[error]") || strings.Contains(stdout.String(), "[gap]") {
		t.Fatalf("spec check clean stdout contains a finding:\n%s", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("spec check stderr = %q, want empty", stderr.String())
	}
}

const specAuditFixtureSlug = "0068-spec-close-audit"

func newSpecAuditWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	gittest.InitRepo(t, repoDir, "--initial-branch=main")
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = roundfix-test@example.com\n[commit]\n\tgpgsign = false\n")
	if err := os.MkdirAll(filepath.Join(repoDir, "docs", "specs", specAuditFixtureSlug), 0o755); err != nil {
		t.Fatalf("create fixture Spec directory: %v", err)
	}
	mustWrite(
		t,
		filepath.Join(repoDir, "docs", "specs", specAuditFixtureSlug, "_prd.md"),
		"---\nspec: "+specAuditFixtureSlug+"\nstatus: active\n---\n",
	)
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "docs: seed Spec audit fixture")
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open fixture Run Database: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close fixture Run Database: %v", err)
	}
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve repo dir: %v", err)
	}
	setCommandEnvironmentForTest(t, homeDir, resolved)
	return homeDir, resolved
}

func archiveSpecAuditFixture(t *testing.T, repoDir string) {
	t.Helper()
	activePath := filepath.ToSlash(filepath.Join("docs", "specs", specAuditFixtureSlug))
	archivedPath := archiveTestPath(spec.ArchiveKindSpec, specAuditFixtureSlug)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(repoDir, archivedPath)), 0o755); err != nil {
		t.Fatalf("create archived Spec parent: %v", err)
	}
	gitImplement(t, repoDir, "mv", activePath, archivedPath)
	gitImplement(t, repoDir, "commit", "-m", "docs: archive Spec audit fixture")
}

func newSpecCheckWorkspace(t *testing.T, slugs ...string) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	fixtureRepo := filepath.Join(cliTestRepoRoot(t), "internal", "speccheck", "testdata", "repo")
	for _, relative := range []string{
		filepath.Join("docs", "agents"),
		filepath.Join("docs", "adr"),
		filepath.Join("docs", "workflow"),
	} {
		if err := copyDir(filepath.Join(fixtureRepo, relative), filepath.Join(repoDir, relative)); err != nil {
			t.Fatalf("copy fixture %s: %v", relative, err)
		}
	}
	for _, slug := range slugs {
		relative := filepath.Join("docs", "specs", slug)
		if err := copyDir(filepath.Join(fixtureRepo, relative), filepath.Join(repoDir, relative)); err != nil {
			t.Fatalf("copy fixture Spec %s: %v", slug, err)
		}
	}

	gittest.InitRepo(t, repoDir, "--initial-branch=main")
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = roundfix-test@example.com\n[commit]\n\tgpgsign = false\n")
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed Spec Consistency Check fixtures")
	gitImplement(t, repoDir, "checkout", "-b", "ma/spec-check")
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve repo dir: %v", err)
	}
	setCommandEnvironmentForTest(t, homeDir, resolved)
	return homeDir, resolved
}

func newSpecCheckVerificationWorkspace(t *testing.T, commands []string) (string, string) {
	t.Helper()
	homeDir, repoDir := newSpecCheckWorkspace(t, "clean")
	specDir := filepath.Join(repoDir, "docs", "specs", "clean")
	mustWrite(t, filepath.Join(specDir, "_tasks.md"), `---
schema: spec-tasks/v1
spec: clean
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---

# Clean fixture Task Graph
`)
	var task strings.Builder
	task.WriteString(`---
task: task_01
spec: clean
status: pending
type: backend
complexity: low
---

# Task 01: Exercise Verification reporting

## Verification

`)
	for _, command := range commands {
		fmt.Fprintf(&task, "- `%s`\n", command)
	}
	mustWrite(t, filepath.Join(specDir, "task_01.md"), task.String())
	gitImplement(t, repoDir, "add", "docs/specs/clean/_tasks.md", "docs/specs/clean/task_01.md")
	gitImplement(t, repoDir, "commit", "-m", "test: seed Verification fixture")
	return homeDir, repoDir
}
