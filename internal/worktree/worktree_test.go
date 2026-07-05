package worktree

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Suite: Run Worktree git lifecycle.
// Invariant: Run Worktrees use named Run Branches and integrate through safe git porcelain.
// Boundary IN: internal/worktree behavior and real git temp repositories.
// Boundary OUT: daemon, store, and CLI wiring that will call this package in later tasks.

func TestCreateUsesNamedRunBranchUnderRoundfixHomeAndCopiesFiles(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	location := filepath.Join(homeDir, "configured-worktrees")
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, ".gitignore"), ".env\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", ".gitignore", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	mustWriteWorktreeTest(t, filepath.Join(repoDir, ".env"), "SECRET=1\n")

	var ref Ref
	stderr := captureStderr(t, func() {
		var err error
		ref, err = Create(ctx, CreateOptions{
			UserRoot: repoDir,
			Location: location,
			RunID:    "task01",
			HeadSHA:  headSHA,
			CopyList: []string{".env", "missing.env"},
		})
		if err != nil {
			t.Fatalf("create Run Worktree: %v", err)
		}
	})

	if ref.RunID != "task01" {
		t.Fatalf("expected Run ID task01, got %q", ref.RunID)
	}
	if ref.UserRoot != filepath.Clean(repoDir) {
		t.Fatalf("expected user root %q, got %q", repoDir, ref.UserRoot)
	}
	if ref.Branch != "roundfix/run-task01" {
		t.Fatalf("expected named Run Branch, got %q", ref.Branch)
	}
	if filepath.Base(ref.Path) != "task01" {
		t.Fatalf("expected path to end in Run ID, got %q", ref.Path)
	}
	if filepath.Dir(filepath.Dir(ref.Path)) != location {
		t.Fatalf("expected path under configured worktree location %q, got %q", location, ref.Path)
	}
	repoSlug := filepath.Base(filepath.Dir(ref.Path))
	if !strings.HasPrefix(repoSlug, filepath.Base(repoDir)+"-") || !hasLowerHexSuffix(repoSlug, 8) {
		t.Fatalf("expected readable repo slug with 8 hex suffix, got %q", repoSlug)
	}
	if got := strings.TrimSpace(gitWorktreeTest(t, ref.Path, "branch", "--show-current")); got != ref.Branch {
		t.Fatalf("expected worktree branch %q, got %q", ref.Branch, got)
	}
	if got := strings.TrimSpace(gitWorktreeTest(t, repoDir, "branch", "--show-current")); got != "main" {
		t.Fatalf("expected source checkout to stay on main, got %q", got)
	}
	if got := mustReadWorktreeTest(t, filepath.Join(ref.Path, ".env")); got != "SECRET=1\n" {
		t.Fatalf("expected copied env file, got %q", got)
	}
	if !strings.Contains(stderr, "missing.env") || !strings.Contains(stderr, "missing") {
		t.Fatalf("expected missing copy-list note on stderr, got %q", stderr)
	}
	if status := gitStatus(t, ref.Path); status != "" {
		t.Fatalf("expected clean Run Worktree, got %q", status)
	}
}

func TestDeriveRootPathUsesReadableUniqueRepoSlug(t *testing.T) {
	location := t.TempDir()
	parentOne := t.TempDir()
	parentTwo := t.TempDir()
	repoOne := filepath.Join(parentOne, "same-name")
	repoTwo := filepath.Join(parentTwo, "same-name")
	mustMkdirWorktreeTest(t, repoOne)
	mustMkdirWorktreeTest(t, repoTwo)

	pathOne, err := deriveRootPath(location, repoOne, "run-one")
	if err != nil {
		t.Fatalf("derive first path: %v", err)
	}
	pathTwo, err := deriveRootPath(location, repoTwo, "run-one")
	if err != nil {
		t.Fatalf("derive second path: %v", err)
	}

	slugOne := filepath.Base(filepath.Dir(pathOne))
	slugTwo := filepath.Base(filepath.Dir(pathTwo))
	if !strings.HasPrefix(slugOne, "same-name-") || !hasLowerHexSuffix(slugOne, 8) {
		t.Fatalf("expected first readable repo slug with 8 hex suffix, got %q", slugOne)
	}
	if !strings.HasPrefix(slugTwo, "same-name-") || !hasLowerHexSuffix(slugTwo, 8) {
		t.Fatalf("expected second readable repo slug with 8 hex suffix, got %q", slugTwo)
	}
	if slugOne == slugTwo {
		t.Fatalf("expected same-named repos at different paths to get unique slugs, both got %q", slugOne)
	}
	if filepath.Dir(filepath.Dir(pathOne)) != location || filepath.Base(pathOne) != "run-one" {
		t.Fatalf("expected path under location/repo-slug/run-id, got %q", pathOne)
	}
}

func TestIntegrateFastForwardsCleanCheckout(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "ff-clean")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModeFastForwardMerge || result.Reason != "" {
		t.Fatalf("expected ff merge result, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at run SHA %s, got %s", runSHA, head)
	}
	if status := gitStatus(t, fixture.repoDir); status != "" {
		t.Fatalf("expected clean user checkout, got %q", status)
	}
}

func TestIntegrateFastForwardsPreservingNonOverlappingDirt(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "ff-dirty")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "user.txt"), "user dirt\n")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModeFastForwardMerge || result.Reason != "" {
		t.Fatalf("expected ff merge result, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at run SHA %s, got %s", runSHA, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if !strings.Contains(status, " M user.txt") {
		t.Fatalf("expected non-overlapping user dirt preserved, got %q", status)
	}
}

func TestIntegrateReturnsPendingOnOverlapAndLeavesBranchAndDirt(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "overlap")
	before := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))
	runSHA := fixture.commitRunChange(t, "shared.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "shared.txt"), "user dirt\n")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModePending || result.Reason != ReasonOverlap {
		t.Fatalf("expected pending overlap, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != before {
		t.Fatalf("expected main to stay at %s, got %s", before, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if !strings.Contains(status, " M shared.txt") {
		t.Fatalf("expected overlapping user dirt preserved, got %q", status)
	}
	assertNoPhantomStagedEntries(t, status)
}

func TestIntegrateReturnsPendingOnDivergenceAndLeavesBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "diverged")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "user.txt"), "user commit\n")
	gitWorktreeTest(t, fixture.repoDir, "add", "user.txt")
	gitWorktreeTest(t, fixture.repoDir, "commit", "-m", "user commit")
	userSHA := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModePending || result.Reason != ReasonDiverged {
		t.Fatalf("expected pending divergence, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != userSHA {
		t.Fatalf("expected main to stay at %s, got %s", userSHA, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if status != "" {
		t.Fatalf("expected clean user checkout after divergence refusal, got %q", status)
	}
	assertNoPhantomStagedEntries(t, status)
}

func TestIntegrateMovesBranchWhenTargetCheckedOutNowhere(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "branch-move")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	gitWorktreeTest(t, fixture.repoDir, "checkout", "-b", "other")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModeBranchMove || result.Reason != "" {
		t.Fatalf("expected branch move result, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at run SHA %s, got %s", runSHA, head)
	}
	if branch := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "branch", "--show-current")); branch != "other" {
		t.Fatalf("expected user checkout to remain on other branch, got %q", branch)
	}
	if status := gitStatus(t, fixture.repoDir); status != "" {
		t.Fatalf("expected clean user checkout, got %q", status)
	}
}

func TestIntegrateReturnsPendingOnNonAncestryWhenTargetCheckedOutNowhere(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "non-ancestry")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "user.txt"), "user commit\n")
	gitWorktreeTest(t, fixture.repoDir, "add", "user.txt")
	gitWorktreeTest(t, fixture.repoDir, "commit", "-m", "user commit")
	userSHA := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))
	gitWorktreeTest(t, fixture.repoDir, "checkout", "-b", "other")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModePending || result.Reason != ReasonNonAncestry {
		t.Fatalf("expected pending non-ancestry, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != userSHA {
		t.Fatalf("expected main to stay at %s, got %s", userSHA, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if status != "" {
		t.Fatalf("expected clean user checkout after non-ancestry refusal, got %q", status)
	}
	assertNoPhantomStagedEntries(t, status)
}

func TestCleanupCleanRemovesCleanRunWorktreeAndBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "cleanup")

	if err := CleanupClean(ctx, fixture.ref); err != nil {
		t.Fatalf("cleanup Clean: %v", err)
	}

	if _, err := os.Stat(fixture.ref.Path); !os.IsNotExist(err) {
		t.Fatalf("expected Run Worktree removed, stat err=%v", err)
	}
	if branchExists(t, fixture.repoDir, fixture.ref.Branch) {
		t.Fatalf("expected Run Branch %q deleted", fixture.ref.Branch)
	}
}

func TestCleanupCleanKeepsDirtyRunWorktreeAndBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "dirty-cleanup")
	mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, "run.txt"), "dirty\n")

	err := CleanupClean(ctx, fixture.ref)

	if err == nil {
		t.Fatal("expected dirty Run Worktree removal to fail")
	}
	if _, statErr := os.Stat(fixture.ref.Path); statErr != nil {
		t.Fatalf("expected dirty Run Worktree kept, stat err=%v", statErr)
	}
	if !branchExists(t, fixture.repoDir, fixture.ref.Branch) {
		t.Fatalf("expected Run Branch %q kept", fixture.ref.Branch)
	}
}

func TestPruneTerminalReapsCrashedTerminalCleanRunOnly(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	location := filepath.Join(homeDir, ".roundfix", "worktrees")
	cleanRef, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "terminal-clean", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create terminal clean ref: %v", err)
	}
	keptRef, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "kept", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create kept ref: %v", err)
	}
	if err := os.RemoveAll(cleanRef.Path); err != nil {
		t.Fatalf("simulate crashed worktree removal: %v", err)
	}

	err = PruneTerminal(ctx, repoDir, location, func(runID string) bool {
		return runID == "terminal-clean"
	})
	if err != nil {
		t.Fatalf("prune terminal: %v", err)
	}

	if branchExists(t, repoDir, cleanRef.Branch) {
		t.Fatalf("expected terminal Clean Run Branch %q deleted", cleanRef.Branch)
	}
	if _, err := os.Stat(keptRef.Path); err != nil {
		t.Fatalf("expected non-terminal Run Worktree to survive, stat err=%v", err)
	}
	if !branchExists(t, repoDir, keptRef.Branch) {
		t.Fatalf("expected non-terminal Run Branch %q to survive", keptRef.Branch)
	}
}

type integrationFixture struct {
	repoDir string
	ref     Ref
}

func newIntegrationFixture(t *testing.T, runID string) integrationFixture {
	t.Helper()
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "shared.txt"), "base\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "run.txt"), "base\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "user.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "shared.txt", "run.txt", "user.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	location := filepath.Join(homeDir, ".roundfix", "worktrees")
	ref, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: runID, HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create Run Worktree: %v", err)
	}
	return integrationFixture{repoDir: repoDir, ref: ref}
}

func (fixture integrationFixture) commitRunChange(t *testing.T, path string, content string) string {
	t.Helper()
	mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, path), content)
	gitWorktreeTest(t, fixture.ref.Path, "add", path)
	gitWorktreeTest(t, fixture.ref.Path, "commit", "-m", "run change")
	return strings.TrimSpace(gitWorktreeTest(t, fixture.ref.Path, "rev-parse", "HEAD"))
}

func initWorktreeRepo(t *testing.T) string {
	t.Helper()
	repoDir := t.TempDir()
	gitWorktreeTest(t, repoDir, "init", "-b", "main")
	gitWorktreeTest(t, repoDir, "config", "user.name", "Roundfix Test")
	gitWorktreeTest(t, repoDir, "config", "user.email", "test@example.com")
	gitWorktreeTest(t, repoDir, "config", "commit.gpgsign", "false")
	return repoDir
}

func gitWorktreeTest(t *testing.T, workDir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-C", workDir}, gitConfigArgsForWorktreeTest()...)
	cmdArgs = append(cmdArgs, "-c", "core.fsmonitor=false")
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = isolatedGitEnvForWorktreeTest()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func gitConfigArgsForWorktreeTest() []string {
	return []string{
		"-c", "user.name=Roundfix Test",
		"-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
	}
}

func isolatedGitEnvForWorktreeTest() []string {
	env := make([]string, 0, len(os.Environ())+2)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_CONFIG_") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
}

func gitStatus(t *testing.T, workDir string) string {
	t.Helper()
	return gitWorktreeTest(t, workDir, "status", "--porcelain=v1")
}

func assertNoPhantomStagedEntries(t *testing.T, status string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimRight(status, "\n"), "\n") {
		if line == "" || strings.HasPrefix(line, "?? ") {
			continue
		}
		if len(line) >= 2 && line[0] != ' ' {
			t.Fatalf("expected no staged phantom entry, got status %q", status)
		}
	}
}

func branchExists(t *testing.T, repoDir string, branch string) bool {
	t.Helper()
	return strings.TrimSpace(gitWorktreeTest(t, repoDir, "branch", "--list", branch)) != ""
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stderr
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("capture stderr pipe: %v", err)
	}
	os.Stderr = writePipe
	output := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		_, _ = buffer.ReadFrom(readPipe)
		output <- buffer.String()
	}()

	fn()

	if err := writePipe.Close(); err != nil {
		t.Fatalf("close captured stderr writer: %v", err)
	}
	os.Stderr = original
	got := <-output
	if err := readPipe.Close(); err != nil {
		t.Fatalf("close captured stderr reader: %v", err)
	}
	return got
}

func mustWriteWorktreeTest(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustMkdirWorktreeTest(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustReadWorktreeTest(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func hasLowerHexSuffix(value string, length int) bool {
	if len(value) < length {
		return false
	}
	suffix := value[len(value)-length:]
	for _, char := range suffix {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}
