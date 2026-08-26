package daemon

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
)

func TestExecVerifierRemovesSuccessfulOutputArtifact(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "runs", "run-test", "verification", "batch-001-attempt-1.log")
	result, err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    "printf verified",
		OutputPath: outputPath,
	})

	if err != nil {
		t.Fatalf("verify command: %v", err)
	}
	if result.OutputPath != outputPath {
		t.Fatalf("expected result output path %q, got %q", outputPath, result.OutputPath)
	}
	if _, statErr := os.Stat(outputPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected successful verification artifact removed, stat err=%v", statErr)
	}
}

func TestExecVerifierRetainsFailedOutputAsTypedCommandError(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "runs", "run-test", "verification", "batch-001-attempt-1.log")
	command := `printf '\163\164\144\157\165\164'; printf '\163\164\144\145\162\162' >&2; exit 7`
	_, err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    command,
		OutputPath: outputPath,
	})

	if err == nil {
		t.Fatal("expected verification failure")
	}
	var commandErr *VerificationCommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected typed command failure, got %T %[1]v", err)
	}
	if commandErr.Command != command || commandErr.OutputPath != outputPath {
		t.Fatalf("unexpected typed error metadata: %+v", commandErr)
	}
	content, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatalf("read failed verification artifact: %v", readErr)
	}
	if string(content) != "stdoutstderr" {
		t.Fatalf("expected combined output artifact, got %q", string(content))
	}
	if strings.Contains(err.Error(), "stdout") || strings.Contains(err.Error(), "stderr") {
		t.Fatalf("expected error to omit output body, got %v", err)
	}
}

func TestExecVerifierTemporaryExit75PreservesDiagnosticChain(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "runs", "run-test", "verification", "batch-001-attempt-1.log")
	command := `printf 'temporary stdout'; printf 'temporary stderr' >&2; exit 75`
	_, err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    command,
		OutputPath: outputPath,
	})

	if err == nil {
		t.Fatal("expected temporary verification failure")
	}
	var temporaryErr *TemporaryVerificationFailureError
	if !errors.As(err, &temporaryErr) {
		t.Fatalf("expected typed temporary verification failure, got %T %[1]v", err)
	}
	var commandErr *VerificationCommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected temporary failure to retain command failure chain, got %T %[1]v", err)
	}
	if temporaryErr.CommandFailure != commandErr {
		t.Fatalf("expected one preserved command failure, got temporary=%p command=%p", temporaryErr.CommandFailure, commandErr)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != TemporaryVerificationExitCode {
		t.Fatalf("expected child exit %d in error chain, got %v", TemporaryVerificationExitCode, err)
	}
	content, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatalf("read temporary verification artifact: %v", readErr)
	}
	if string(content) != "temporary stdouttemporary stderr" {
		t.Fatalf("expected combined temporary diagnostics, got %q", content)
	}
}

func TestExecVerifierExit1WithTimeoutTextRemainsDeterministic(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "verification.log")
	command := `printf 'timeout waiting for listener on database port'; exit 1`
	_, err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    command,
		OutputPath: outputPath,
	})

	if err == nil {
		t.Fatal("expected deterministic verification failure")
	}
	var temporaryErr *TemporaryVerificationFailureError
	if errors.As(err, &temporaryErr) {
		t.Fatalf("expected output text not to classify exit 1 as temporary, got %+v", temporaryErr)
	}
	var commandErr *VerificationCommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected deterministic command failure, got %T %[1]v", err)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected child exit 1 in deterministic chain, got %v", err)
	}
}

func TestExecVerifierTemporaryRetryDiagnosticPathsRemainDistinct(t *testing.T) {
	t.Parallel()
	artifactDir := t.TempDir()
	initialPath := VerificationOutputPath(artifactDir, "run-test", 1, 1)
	retryPath := VerificationRetryOutputPath(artifactDir, "run-test", 1, 1, 1)
	workDir := t.TempDir()

	_, initialErr := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    workDir,
		Command:    `printf 'initial temporary diagnostic'; exit 75`,
		OutputPath: initialPath,
	})
	_, retryErr := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    workDir,
		Command:    `printf 'exclusive deterministic diagnostic'; exit 1`,
		OutputPath: retryPath,
	})

	var temporaryErr *TemporaryVerificationFailureError
	if !errors.As(initialErr, &temporaryErr) {
		t.Fatalf("expected initial typed temporary failure, got %v", initialErr)
	}
	var retryCommandErr *VerificationCommandError
	if !errors.As(retryErr, &retryCommandErr) {
		t.Fatalf("expected retry deterministic command failure, got %v", retryErr)
	}
	if retryCommandErr.OutputPath != retryPath || retryPath == initialPath {
		t.Fatalf("expected distinct retry path %q from initial %q, got %+v", retryPath, initialPath, retryCommandErr)
	}
	initialDiagnostic, err := os.ReadFile(initialPath)
	if err != nil {
		t.Fatalf("read initial diagnostic: %v", err)
	}
	retryDiagnostic, err := os.ReadFile(retryPath)
	if err != nil {
		t.Fatalf("read retry diagnostic: %v", err)
	}
	if string(initialDiagnostic) != "initial temporary diagnostic" || string(retryDiagnostic) != "exclusive deterministic diagnostic" {
		t.Fatalf("unexpected retained diagnostics: initial=%q retry=%q", initialDiagnostic, retryDiagnostic)
	}
}

func TestExecVerifierClassifiesCancellationAsInfrastructureError(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "verification.log")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := ExecVerifier{}.Verify(ctx, VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    "printf canceled",
		OutputPath: outputPath,
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	var commandErr *VerificationCommandError
	if errors.As(err, &commandErr) {
		t.Fatalf("expected cancellation to stay infrastructure error, got %+v", commandErr)
	}
	if _, statErr := os.Stat(outputPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected no retained artifact on cancellation, stat err=%v", statErr)
	}
}

func TestExecVerifierClassifiesDeadlineAsUnknownVerdict(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "verification.log")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	_, err := ExecVerifier{}.Verify(ctx, VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    "printf never-started",
		OutputPath: outputPath,
	})

	var unknownErr *VerificationUnknownError
	if !errors.As(err, &unknownErr) {
		t.Fatalf("expected deadline to produce an unknown verdict, got %T %[1]v", err)
	}
	if unknownErr.Command != "printf never-started" || unknownErr.DiagnosticPath != outputPath {
		t.Fatalf("unexpected unknown-verdict metadata: %+v", unknownErr)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline identity preserved, got %v", err)
	}
}

func TestExecVerifierClassifiesProcessStartFailureAsUnknownVerdict(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "verification.log")
	_, err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    filepath.Join(t.TempDir(), "missing"),
		Command:    "printf never-started",
		OutputPath: outputPath,
	})

	if err == nil {
		t.Fatal("expected process-start failure")
	}
	var commandErr *VerificationCommandError
	if errors.As(err, &commandErr) {
		t.Fatalf("expected process-start failure not to become a command verdict, got %+v", commandErr)
	}
	var unknownErr *VerificationUnknownError
	if !errors.As(err, &unknownErr) {
		t.Fatalf("expected process-start failure to produce an unknown verdict, got %T %[1]v", err)
	}
	if unknownErr.Command != "printf never-started" || unknownErr.DiagnosticPath != outputPath {
		t.Fatalf("unexpected unknown-verdict metadata: %+v", unknownErr)
	}
	if _, statErr := os.Stat(outputPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected no retained artifact on process-start failure, stat err=%v", statErr)
	}
}

func TestExecVerifierClassifiesArtifactRetentionFailureAsInfrastructureError(t *testing.T) {
	t.Parallel()
	outputPath := filepath.Join(t.TempDir(), "verification.log")
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		t.Fatalf("create conflicting output path: %v", err)
	}

	_, err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir:    t.TempDir(),
		Command:    "printf broken; exit 7",
		OutputPath: outputPath,
	})

	if err == nil {
		t.Fatal("expected artifact retention failure")
	}
	var commandErr *VerificationCommandError
	if errors.As(err, &commandErr) {
		t.Fatalf("expected artifact failure to stay infrastructure error, got %+v", commandErr)
	}
	if !strings.Contains(err.Error(), "retain failed verification diagnostics") {
		t.Fatalf("expected artifact retention context, got %v", err)
	}
}

func TestGitCommitterValidatesRequest(t *testing.T) {
	t.Parallel()
	err := GitCommitter{}.Commit(context.Background(), CommitRequest{})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "git root is required") {
		t.Fatalf("expected git root validation error, got %v", err)
	}
}

func TestGitCommitterExcludesProjectConfigFromBatchCommit(t *testing.T) {
	t.Parallel()
	repoDir := t.TempDir()
	gittest.InitRepo(t, repoDir)
	gittest.AppendConfig(t, repoDir, "[user]\n\temail = test@example.com\n\tname = Test\n[commit]\n\tgpgsign = false\n")
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	runGitForTest(t, repoDir, "add", "tracked.txt")
	runGitForTest(t, repoDir, "commit", "-q", "-m", "init")

	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "changed\n")
	mustWriteForTest(t, filepath.Join(repoDir, "created.txt"), "created\n")
	mustWriteForTest(t, filepath.Join(repoDir, ".roundfixrc.yml"), "defaults:\n  agent: codex\n")

	err := GitCommitter{}.Commit(context.Background(), CommitRequest{
		WorkDir:      repoDir,
		Message:      "fix: test batch",
		ExcludePaths: []string{".roundfixrc.yml"},
	})
	if err != nil {
		t.Fatalf("expected Batch commit, got %v", err)
	}

	committed := runGitForTest(t, repoDir, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "tracked.txt") || !strings.Contains(committed, "created.txt") {
		t.Fatalf("expected code changes in commit, got %q", committed)
	}
	if strings.Contains(committed, ".roundfixrc.yml") {
		t.Fatalf("did not expect Project Config in commit, got %q", committed)
	}
	status := runGitForTest(t, repoDir, "status", "--porcelain=v1")
	if !strings.Contains(status, "?? .roundfixrc.yml") {
		t.Fatalf("expected Project Config to remain untracked, got %q", status)
	}
}

func TestGitCommitterStagesSelectedTrackedPathMatchedByGlobalIgnore(t *testing.T) {
	repoDir := t.TempDir()
	globalConfigDir := t.TempDir()
	globalIgnore := filepath.Join(globalConfigDir, "git", "ignore")
	if err := os.MkdirAll(filepath.Dir(globalIgnore), 0o755); err != nil {
		t.Fatalf("create global git config directory: %v", err)
	}
	mustWriteForTest(t, globalIgnore, "dist/\n")
	t.Setenv("XDG_CONFIG_HOME", globalConfigDir)
	gittest.InitRepo(t, repoDir)
	gittest.AppendConfig(t, repoDir, "[user]\n\temail = test@example.com\n\tname = Test\n[commit]\n\tgpgsign = false\n")
	trackedPath := filepath.Join(repoDir, "dist", "tracked.txt")
	ignoredPath := filepath.Join(repoDir, "dist", "untracked.txt")
	if err := os.MkdirAll(filepath.Dir(trackedPath), 0o755); err != nil {
		t.Fatalf("create dist directory: %v", err)
	}
	mustWriteForTest(t, trackedPath, "base\n")
	runGitForTest(t, repoDir, "add", "-f", "dist/tracked.txt")
	runGitForTest(t, repoDir, "commit", "-q", "-m", "init")

	mustWriteForTest(t, trackedPath, "changed\n")
	mustWriteForTest(t, ignoredPath, "not selected\n")

	err := GitCommitter{}.Commit(context.Background(), CommitRequest{
		WorkDir: repoDir,
		Message: "fix: test exact staging",
		Paths:   []string{"dist/tracked.txt"},
	})
	if err != nil {
		t.Fatalf("commit selected tracked path matched by global ignore: %v", err)
	}

	committed := runGitForTest(t, repoDir, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "dist/tracked.txt") {
		t.Fatalf("expected selected tracked path in commit, got %q", committed)
	}
	if strings.Contains(committed, "dist/untracked.txt") {
		t.Fatalf("did not expect unselected ignored path in commit, got %q", committed)
	}
	if _, err := os.Stat(ignoredPath); err != nil {
		t.Fatalf("expected unselected ignored file preserved: %v", err)
	}
}

func TestGitPusherValidatesRequest(t *testing.T) {
	t.Parallel()
	err := GitPusher{}.Push(context.Background(), PushRequest{})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "git root is required") {
		t.Fatalf("expected git root validation error, got %v", err)
	}
}

func runGitForTest(t *testing.T, workDir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-C", workDir}, gitConfigArgsForTest()...)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = isolatedGitEnvForTest()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func gitConfigArgsForTest() []string {
	return gittest.ConfigArgs()
}

func isolatedGitEnvForTest() []string {
	return gittest.IsolatedEnv()
}

func mustWriteForTest(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readFileForTest(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func TestSnapshotDiffCommitStagesOnlyAgentChangesInRealRepo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repoDir := t.TempDir()
	gittest.InitRepo(t, repoDir, "-b", "main")
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = test@example.com\n[commit]\n\tgpgsign = false\n")
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "original\n")
	runGitForTest(t, repoDir, "add", "tracked.txt")
	runGitForTest(t, repoDir, "commit", "-m", "initial")

	// Pre-existing user work that slipped past Preflight Validation.
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "user edit\n")
	mustWriteForTest(t, filepath.Join(repoDir, "user-notes.md"), "wip\n")

	snapshotter := GitWorktreeSnapshotter{}
	before, err := snapshotter.Snapshot(ctx, repoDir)
	if err != nil {
		t.Fatalf("before snapshot: %v", err)
	}

	// The Agent creates a fix.
	mustWriteForTest(t, filepath.Join(repoDir, "agent-fix.go"), "package fix\n")

	after, err := snapshotter.Snapshot(ctx, repoDir)
	if err != nil {
		t.Fatalf("after snapshot: %v", err)
	}
	changed := diffSnapshots(before, after)
	if len(changed) != 1 || changed[0] != "agent-fix.go" {
		t.Fatalf("expected only the Agent-made change in the diff, got %v", changed)
	}

	if err := (GitCommitter{}).Commit(ctx, CommitRequest{
		WorkDir: repoDir,
		Message: BatchCommitMessage(1),
		Paths:   changed,
	}); err != nil {
		t.Fatalf("snapshot-diff commit: %v", err)
	}

	committed := runGitForTest(t, repoDir, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "agent-fix.go") {
		t.Fatalf("expected Agent change committed, got %q", committed)
	}
	if strings.Contains(committed, "tracked.txt") || strings.Contains(committed, "user-notes.md") {
		t.Fatalf("expected user changes kept out of the Batch commit, got %q", committed)
	}
	status := runGitForTest(t, repoDir, "status", "--porcelain=v1")
	if !strings.Contains(status, " M tracked.txt") || !strings.Contains(status, "?? user-notes.md") {
		t.Fatalf("expected user work preserved in the worktree, got %q", status)
	}
}

func TestRunGitForTestIgnoresForcedSigningConfig(t *testing.T) {
	t.Setenv("GIT_CONFIG_COUNT", "2")
	t.Setenv("GIT_CONFIG_KEY_0", "commit.gpgsign")
	t.Setenv("GIT_CONFIG_VALUE_0", "true")
	t.Setenv("GIT_CONFIG_KEY_1", "gpg.program")
	t.Setenv("GIT_CONFIG_VALUE_1", "false")

	if got := mustReadUnisolatedGitConfigForCanary(t, "commit.gpgsign"); strings.TrimSpace(got) != "true" {
		t.Fatalf("expected forced signing visible to unisolated git, got %q", got)
	}

	repoDir := t.TempDir()
	gittest.InitRepo(t, repoDir)
	if got := runGitForTest(t, repoDir, "config", "--get", "commit.gpgsign"); strings.TrimSpace(got) != "false" {
		t.Fatalf("expected isolated helper to override forced signing, got %q", got)
	}
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = test@example.com\n[commit]\n\tgpgsign = false\n")
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	runGitForTest(t, repoDir, "add", "tracked.txt")
	runGitForTest(t, repoDir, "commit", "-q", "-m", "init")
}

func mustReadUnisolatedGitConfigForCanary(t *testing.T, key string) string {
	t.Helper()
	cmd := exec.Command("git", "config", "--get", key)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unisolated git config --get %s failed: %v\n%s", key, err, output)
	}
	return string(output)
}

func TestReviewArtifactsCommitMessage(t *testing.T) {
	t.Parallel()
	if got := ReviewArtifactsCommitMessage(1, "18"); got != "docs: review round 001 for pr 18" {
		t.Fatalf("expected single-round message, got %q", got)
	}
	if got := ReviewArtifactsCommitMessage(0, "18"); got != "docs: review rounds for pr 18" {
		t.Fatalf("expected all-rounds message, got %q", got)
	}
}
