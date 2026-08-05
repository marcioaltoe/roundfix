// Suite: Spec survivor classification
// Invariant: every surviving Spec branch and worktree is classified from local evidence without mutation.
// Boundary IN: disposable Git repositories and the real read-only Run Database.
// Boundary OUT: CLI rendering and delivery verification, owned by later Spec 0068 tasks.
package specaudit

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

const auditFixtureSlug = "0068-spec-close-audit"

type auditFixture struct {
	t       *testing.T
	ctx     context.Context
	repoDir string
	homeDir string
	store   *store.Store
}

func TestAuditClassifiesPullRequestBranch(t *testing.T) {
	t.Parallel()
	fixture := newAuditFixture(t)
	branch := "ma/spec-close-pull-request"
	fixture.commitBranch(branch, "pull-request.txt")
	fixture.createRun(store.CreateRunRequest{
		Kind:           store.KindFetch,
		HeadRepository: "owner/repository",
		HeadBranch:     branch,
		BaseRepository: "owner/repository",
		PRNumber:       "42",
		GitRoot:        fixture.repoDir,
		LocalBranch:    branch,
		HeadSHA:        fixture.git("rev-parse", branch),
		ArtifactDir:    filepath.Join(fixture.homeDir, "artifacts"),
		SpecSlug:       auditFixtureSlug,
	}, store.StateClean)

	result := fixture.audit()
	survivor := requireSurvivor(t, result, branch, false)
	if survivor.Kind != KindPullRequest {
		t.Fatalf("branch kind = %q, want %q", survivor.Kind, KindPullRequest)
	}
	if !strings.Contains(survivor.Evidence, "Pull Request #42") {
		t.Fatalf("branch evidence = %q, want Pull Request #42", survivor.Evidence)
	}
	assertEverySurvivorHasEvidence(t, result)
}

func TestAuditClassifiesPendingBranch(t *testing.T) {
	t.Parallel()
	fixture := newAuditFixture(t)
	branch := "ma/spec-close-pending"
	fixture.commitBranch(branch, "pending.txt")
	fixture.createImplementRun(branch, "", store.StateClean)

	result := fixture.audit()
	survivor := requireSurvivor(t, result, branch, false)
	if survivor.Kind != KindPending {
		t.Fatalf("branch kind = %q, want %q", survivor.Kind, KindPending)
	}
	if !strings.Contains(survivor.Evidence, "1 commit") ||
		!strings.Contains(survivor.Evidence, "not represented") {
		t.Fatalf("branch evidence = %q, want the unintegrated commit evidence", survivor.Evidence)
	}
	assertEverySurvivorHasEvidence(t, result)
}

func TestAuditClassifiesResidueBranch(t *testing.T) {
	t.Parallel()
	fixture := newAuditFixture(t)
	branch := "ma/spec-close-residue"
	fixture.commitBranch(branch, "residue.txt")
	fixture.git("merge", "--squash", branch)
	fixture.git("commit", "-m", "feat: merge residue fixture")
	fixture.createImplementRun(branch, "", store.StateClean)

	result := fixture.audit()
	survivor := requireSurvivor(t, result, branch, false)
	if survivor.Kind != KindResidue {
		t.Fatalf("branch kind = %q, want %q", survivor.Kind, KindResidue)
	}
	wantReclaim := "git branch -d -- 'ma/spec-close-residue'"
	if survivor.Reclaim != wantReclaim {
		t.Fatalf("branch reclaim = %q, want %q", survivor.Reclaim, wantReclaim)
	}
	if !strings.Contains(survivor.Evidence, "content is fully represented") {
		t.Fatalf("branch evidence = %q, want content integration proof", survivor.Evidence)
	}
	assertEverySurvivorHasEvidence(t, result)
}

func TestAuditPreservesUnmatchedWorktree(t *testing.T) {
	t.Parallel()
	fixture := newAuditFixture(t)
	branch := "ma/spec-close-scratch"
	worktreePath := fixture.addWorktree(branch, "scratch.txt")
	fixture.createRun(store.CreateRunRequest{
		Kind:           store.KindFetch,
		HeadRepository: "owner/repository",
		HeadBranch:     branch,
		BaseRepository: "owner/repository",
		PRNumber:       "43",
		GitRoot:        fixture.repoDir,
		LocalBranch:    branch,
		HeadSHA:        fixture.git("rev-parse", branch),
		ArtifactDir:    filepath.Join(fixture.homeDir, "artifacts"),
		SpecSlug:       auditFixtureSlug,
	}, store.StateClean)

	result := fixture.audit()
	survivor := requireSurvivor(t, result, worktreePath, true)
	if survivor.Kind != KindPreserved {
		t.Fatalf("worktree kind = %q, want %q", survivor.Kind, KindPreserved)
	}
	if !strings.Contains(survivor.Evidence, "no matching Run") {
		t.Fatalf("worktree evidence = %q, want missing Run evidence", survivor.Evidence)
	}
	if survivor.Reclaim != "" {
		t.Fatalf("preserved worktree reclaim = %q, want empty", survivor.Reclaim)
	}
	assertEverySurvivorHasEvidence(t, result)
}

func TestAuditPreservesActiveRunSurvivors(t *testing.T) {
	t.Parallel()
	fixture := newAuditFixture(t)
	branch := "ma/spec-close-active"
	worktreePath := fixture.addWorktree(branch, "active.txt")
	run := fixture.createImplementRun(branch, worktreePath, store.StateActive)

	result := fixture.audit()
	for _, target := range []struct {
		name       string
		isWorktree bool
	}{
		{name: branch},
		{name: worktreePath, isWorktree: true},
	} {
		survivor := requireSurvivor(t, result, target.name, target.isWorktree)
		if survivor.Kind == KindResidue {
			t.Fatalf("Active Run survivor %q classified as residue", target.name)
		}
		if !strings.Contains(survivor.Evidence, "Active Run "+run.ID) {
			t.Fatalf("Active Run survivor evidence = %q, want Run %s", survivor.Evidence, run.ID)
		}
		if survivor.Reclaim != "" {
			t.Fatalf("Active Run survivor reclaim = %q, want empty", survivor.Reclaim)
		}
	}
	assertEverySurvivorHasEvidence(t, result)
}

func newAuditFixture(t *testing.T) *auditFixture {
	t.Helper()
	root := t.TempDir()
	repoDir := filepath.Join(root, "repository")
	homeDir := filepath.Join(root, "home")
	gittest.InitRepo(t, repoDir, "-b", "main")
	writeAuditFixtureFile(t, filepath.Join(repoDir, "README.md"), "fixture\n")
	gittest.Run(t, repoDir, "add", "README.md")
	gittest.Run(t, repoDir, "commit", "-m", "chore: initialize fixture")

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open fixture Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close fixture Run Database: %v", err)
		}
	})
	return &auditFixture{
		t:       t,
		ctx:     ctx,
		repoDir: repoDir,
		homeDir: homeDir,
		store:   runStore,
	}
}

func (fixture *auditFixture) commitBranch(branch, filename string) {
	fixture.t.Helper()
	fixture.git("switch", "-c", branch)
	writeAuditFixtureFile(fixture.t, filepath.Join(fixture.repoDir, filename), branch+"\n")
	fixture.git("add", filename)
	fixture.git(
		"commit",
		"-m", "feat: add "+filename,
		"-m", "Roundfix-Spec: "+auditFixtureSlug,
	)
	fixture.git("switch", "main")
}

func (fixture *auditFixture) addWorktree(branch, filename string) string {
	fixture.t.Helper()
	worktreePath := filepath.Join(filepath.Dir(fixture.repoDir), strings.ReplaceAll(branch, "/", "-"))
	fixture.git("worktree", "add", "-b", branch, worktreePath, "main")
	writeAuditFixtureFile(fixture.t, filepath.Join(worktreePath, filename), branch+"\n")
	gittest.Run(fixture.t, worktreePath, "add", filename)
	gittest.Run(
		fixture.t,
		worktreePath,
		"commit",
		"-m", "feat: add "+filename,
		"-m", "Roundfix-Spec: "+auditFixtureSlug,
	)
	resolved, err := filepath.EvalSymlinks(worktreePath)
	if err != nil {
		fixture.t.Fatalf("resolve fixture worktree path: %v", err)
	}
	return resolved
}

func (fixture *auditFixture) createImplementRun(branch, workDir, state string) store.Run {
	fixture.t.Helper()
	return fixture.createRun(store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     fixture.repoDir,
		LocalBranch: branch,
		WorkDir:     workDir,
		SpecSlug:    auditFixtureSlug,
	}, state)
}

func (fixture *auditFixture) createRun(req store.CreateRunRequest, state string) store.Run {
	fixture.t.Helper()
	run, err := fixture.store.CreateRun(fixture.ctx, req)
	if err != nil {
		fixture.t.Fatalf("create fixture Run: %v", err)
	}
	if state == "" || state == store.StateActive {
		return run
	}
	completed, err := fixture.store.CompleteRun(fixture.ctx, run.ID, state)
	if err != nil {
		fixture.t.Fatalf("complete fixture Run: %v", err)
	}
	return completed.Run
}

func (fixture *auditFixture) audit() Result {
	fixture.t.Helper()
	result, err := Audit(fixture.ctx, fixture.repoDir, fixture.homeDir, auditFixtureSlug)
	if err != nil {
		fixture.t.Fatalf("audit fixture: %v", err)
	}
	return result
}

func (fixture *auditFixture) git(args ...string) string {
	fixture.t.Helper()
	return strings.TrimSpace(gittest.Run(fixture.t, fixture.repoDir, args...))
}

func writeAuditFixtureFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture file %s: %v", path, err)
	}
}

func requireSurvivor(t *testing.T, result Result, name string, isWorktree bool) Survivor {
	t.Helper()
	for _, survivor := range result.Survivors {
		if survivor.Name == name && survivor.IsWorktree == isWorktree {
			return survivor
		}
	}
	t.Fatalf("survivor %q (worktree=%t) not found in %#v", name, isWorktree, result.Survivors)
	return Survivor{}
}

func assertEverySurvivorHasEvidence(t *testing.T, result Result) {
	t.Helper()
	if len(result.Survivors) == 0 {
		t.Fatal("audit returned no survivors")
	}
	for _, survivor := range result.Survivors {
		if strings.TrimSpace(survivor.Evidence) == "" {
			t.Fatalf("survivor %#v has empty evidence", survivor)
		}
	}
}
