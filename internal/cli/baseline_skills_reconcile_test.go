// Suite: public Baseline skills-lock reconciliation CLI.
// Invariant: lock removal is previewed, digest-confirmed, and restricted to immutable source commits.
// Boundary IN: public dispatch, flag parsing, help, output, exit categories, local Git acquisition, and lock mutation.
// Boundary OUT: entry classification and transaction fault handling, owned by internal/baseline/skills_reconcile_test.go.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineSkillsReconcilePreviewThenConfirm(t *testing.T) {
	t.Parallel()

	repository := newBaselineSkillsReconcileGitRepository(t, map[string]string{
		"README.md": "repository\n",
	})
	source := newBaselineSkillsReconcileGitRepository(t, map[string]string{
		"README.md": "source\n",
	})
	revision := baselineSkillsReconcileGitOutput(t, source, "rev-parse", "HEAD")
	lockPath := filepath.Join(repository, "skills-lock.json")
	writeBaselineSkillsReconcileFile(t, lockPath, `{
  "version": 1,
  "skills": {
    "obsolete-task05": {
      "source": "example/skills",
      "skillPath": "skills/obsolete-task05/SKILL.md",
      "future": {"preserve-until-removal": true}
    }
  }
}
`)

	args := []string{
		"baseline", "skills", "reconcile",
		"--profile", "go-cli-tui",
		"--source", "example/skills",
		"--revision", revision,
		"--source-dir", source,
		"--repo", repository,
		"--format", "json",
	}
	var previewStdout bytes.Buffer
	var previewStderr bytes.Buffer
	previewCode := RunContext(context.Background(), args, &previewStdout, &previewStderr)
	if previewCode != exitUnverified {
		t.Fatalf(
			"preview exit = %d, want %d; stdout=%s stderr=%s",
			previewCode,
			exitUnverified,
			previewStdout.String(),
			previewStderr.String(),
		)
	}
	var preview baseline.SkillsReconcilePayload
	if err := json.Unmarshal(previewStdout.Bytes(), &preview); err != nil {
		t.Fatalf("decode preview JSON: %v\n%s", err, previewStdout.String())
	}
	if preview.PlanDigest == nil || len(*preview.PlanDigest) != 64 {
		t.Fatalf("preview Plan Digest = %v", preview.PlanDigest)
	}
	if len(preview.PlannedChanges) != 1 ||
		preview.PlannedChanges[0].Action != "remove-lock-entry" ||
		preview.PlannedChanges[0].Skill != "obsolete-task05" {
		t.Fatalf("preview planned changes = %+v", preview.PlannedChanges)
	}
	if !strings.Contains(previewStderr.String(), "baseline skills reconcile failed") {
		t.Fatalf("preview stderr = %q", previewStderr.String())
	}

	confirmedArgs := append(append([]string(nil), args...), "--confirm-plan", *preview.PlanDigest)
	var confirmedStdout bytes.Buffer
	var confirmedStderr bytes.Buffer
	confirmedCode := RunContext(
		context.Background(),
		confirmedArgs,
		&confirmedStdout,
		&confirmedStderr,
	)
	if confirmedCode != exitOK || confirmedStderr.Len() != 0 {
		t.Fatalf(
			"confirmed exit = %d, want %d; stdout=%s stderr=%s",
			confirmedCode,
			exitOK,
			confirmedStdout.String(),
			confirmedStderr.String(),
		)
	}
	var confirmed baseline.SkillsReconcilePayload
	if err := json.Unmarshal(confirmedStdout.Bytes(), &confirmed); err != nil {
		t.Fatalf("decode confirmed JSON: %v\n%s", err, confirmedStdout.String())
	}
	if !confirmed.OK || !confirmed.Applied {
		t.Fatalf("confirmed payload = %+v", confirmed)
	}
	lockBytes, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode applied skills lock: %v\n%s", err, lockBytes)
	}
	if _, exists := lock.Skills["obsolete-task05"]; exists {
		t.Fatalf("obsolete entry survived confirmed reconciliation: %s", lockBytes)
	}

	var currentStdout bytes.Buffer
	var currentStderr bytes.Buffer
	currentCode := RunContext(context.Background(), args, &currentStdout, &currentStderr)
	if currentCode != exitOK || currentStderr.Len() != 0 {
		t.Fatalf(
			"current exit = %d, want %d; stdout=%s stderr=%s",
			currentCode,
			exitOK,
			currentStdout.String(),
			currentStderr.String(),
		)
	}
	var current baseline.SkillsReconcilePayload
	if err := json.Unmarshal(currentStdout.Bytes(), &current); err != nil {
		t.Fatalf("decode current JSON: %v\n%s", err, currentStdout.String())
	}
	if !current.OK || current.Applied || len(current.PlannedChanges) != 0 {
		t.Fatalf("current payload = %+v", current)
	}
}

func TestBaselineSkillsReconcileRejectsAMutableRevision(t *testing.T) {
	t.Parallel()

	missingSource := filepath.Join(t.TempDir(), "fetch-must-not-run")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "skills", "reconcile",
		"--profile", "go-cli-tui",
		"--source", "example/skills",
		"--revision", "main",
		"--source-dir", missingSource,
		"--format", "json",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("mutable revision exit = %d, want %d; stdout=%s stderr=%s", code, exitPreflight, stdout.String(), stderr.String())
	}
	var payload baseline.SkillsReconcilePayload
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("decode refusal JSON: %v\n%s", err, stdout.String())
	}
	if payload.Finding == nil || payload.Finding.Code != "reconcile.commit-invalid" {
		t.Fatalf("mutable revision finding = %+v", payload.Finding)
	}
	if strings.Contains(stdout.String(), "source-dir-invalid") || strings.Contains(stderr.String(), "source-dir-invalid") {
		t.Fatalf("mutable revision reached source acquisition: stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
}

func TestBaselineSkillsReconcileHelpNamesTheConfirmationContract(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "skills", "reconcile", "--help",
	}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("help exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"baseline skills reconcile --profile <id> --source <owner/repo> --revision <commit>",
		"A non-empty preview exits 3",
		"--confirm-plan",
		"Exact lowercase Plan Digest returned by the current preview",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("help missing %q:\n%s", want, stdout.String())
		}
	}
}

func newBaselineSkillsReconcileGitRepository(
	t *testing.T,
	files map[string]string,
) string {
	t.Helper()
	repository := t.TempDir()
	baselineSkillsReconcileGit(t, repository, "init", "--quiet")
	baselineSkillsReconcileGit(t, repository, "config", "user.email", "task05@example.invalid")
	baselineSkillsReconcileGit(t, repository, "config", "user.name", "Task 05")
	baselineSkillsReconcileGit(t, repository, "config", "commit.gpgsign", "false")
	for relative, content := range files {
		writeBaselineSkillsReconcileFile(t, filepath.Join(repository, relative), content)
	}
	baselineSkillsReconcileGit(t, repository, "add", ".")
	baselineSkillsReconcileGit(t, repository, "commit", "--quiet", "-m", "fixture")
	return repository
}

func writeBaselineSkillsReconcileFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func baselineSkillsReconcileGit(t *testing.T, directory string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func baselineSkillsReconcileGitOutput(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}
