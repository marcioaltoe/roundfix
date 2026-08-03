package cli

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/releaseplan"
)

// Suite: Release Plan GitSource
// Invariant: local Git range resolution is read-only and returns committed evidence only.
// Boundary IN: temporary Git repositories through preflight.GitRunner.
// Boundary OUT: Release Plan classification, CLI rendering, network services, release mutation.

func TestReleasePlanGitSourceDefaultRangeResolvesLatestReachableStableTag(t *testing.T) {
	t.Parallel()
	repo := newReleasePlanGitSourceRepo(t)
	source := newReleasePlanGitSource(repo.dir, preflight.ExecGitRunner{})
	before := snapshotReleasePlanRepo(t, repo.dir)

	releaseRange, err := source.ResolveRange(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ResolveRange: %v", err)
	}
	commits, err := source.Commits(context.Background(), releaseRange)
	if err != nil {
		t.Fatalf("Commits: %v", err)
	}

	assertReleasePlanRepoUnchanged(t, repo.dir, before)
	if releaseRange.Base.Tag != "v0.1.1" {
		t.Fatalf("Base.Tag = %q, want latest reachable stable tag v0.1.1", releaseRange.Base.Tag)
	}
	if releaseRange.Base.CommitSHA != repo.patchTagSHA {
		t.Fatalf("Base.CommitSHA = %q, want %q", releaseRange.Base.CommitSHA, repo.patchTagSHA)
	}
	if releaseRange.Target.Name != "HEAD" {
		t.Fatalf("Target.Name = %q, want HEAD", releaseRange.Target.Name)
	}
	if releaseRange.Target.CommitSHA != repo.fixSHA {
		t.Fatalf("Target.CommitSHA = %q, want %q", releaseRange.Target.CommitSHA, repo.fixSHA)
	}
	if gotSubjects := commitSubjects(commits); !reflect.DeepEqual(gotSubjects, []string{"feat: add release plan", "fix: keep range stable"}) {
		t.Fatalf("commit subjects = %v, want feature then fix", gotSubjects)
	}
	if !strings.Contains(commits[0].Body, "Release Plan Command reads committed history.") {
		t.Fatalf("first commit body = %q, want full body text", commits[0].Body)
	}
	if got := commits[0].ChangedPaths; !reflect.DeepEqual(got, []string{"cmd/roundfix/main.go", "internal/cli/release.go"}) {
		t.Fatalf("first commit paths = %v, want deterministic changed paths", got)
	}
}

func TestReleasePlanGitSourceExplicitEndpointsReturnOrderedCommitsAndPaths(t *testing.T) {
	t.Parallel()
	repo := newReleasePlanGitSourceRepo(t)
	source := newReleasePlanGitSource(repo.dir, preflight.ExecGitRunner{})
	before := snapshotReleasePlanRepo(t, repo.dir)

	releaseRange, err := source.ResolveRange(context.Background(), "v0.1.1", repo.featureSHA)
	if err != nil {
		t.Fatalf("ResolveRange: %v", err)
	}
	commits, err := source.Commits(context.Background(), releaseRange)
	if err != nil {
		t.Fatalf("Commits: %v", err)
	}

	assertReleasePlanRepoUnchanged(t, repo.dir, before)
	if releaseRange.Base.Tag != "v0.1.1" || releaseRange.Base.CommitSHA != repo.patchTagSHA {
		t.Fatalf("base = %+v, want v0.1.1 at %s", releaseRange.Base, repo.patchTagSHA)
	}
	if releaseRange.Target.Name != repo.featureSHA || releaseRange.Target.CommitSHA != repo.featureSHA {
		t.Fatalf("target = %+v, want explicit feature commit", releaseRange.Target)
	}
	if len(commits) != 1 {
		t.Fatalf("commits = %+v, want one feature commit", commits)
	}
	if commits[0].SHA != repo.featureSHA {
		t.Fatalf("commit SHA = %q, want %q", commits[0].SHA, repo.featureSHA)
	}
	if got := commits[0].ChangedPaths; !reflect.DeepEqual(got, []string{"cmd/roundfix/main.go", "internal/cli/release.go"}) {
		t.Fatalf("changed paths = %v, want feature commit paths", got)
	}
}

func TestReleasePlanGitSourceDirtyWorktreeReportsPathsAndPreservesRepo(t *testing.T) {
	t.Parallel()
	repo := newReleasePlanGitSourceRepo(t)
	writeReleasePlanFile(t, repo.dir, "internal/cli/release.go", "dirty tracked change\n")
	writeReleasePlanFile(t, repo.dir, "scratch.txt", "untracked change\n")
	source := newReleasePlanGitSource(repo.dir, preflight.ExecGitRunner{})
	before := snapshotReleasePlanRepo(t, repo.dir)

	_, err := source.ResolveRange(context.Background(), "", "")

	assertReleasePlanRepoUnchanged(t, repo.dir, before)
	if !errors.Is(err, releaseplan.ErrDirtyWorktree) {
		t.Fatalf("ResolveRange error = %v, want dirty worktree", err)
	}
	var sourceErr releaseplan.GitSourceError
	if !errors.As(err, &sourceErr) {
		t.Fatalf("ResolveRange error = %T, want GitSourceError", err)
	}
	for _, path := range []string{"internal/cli/release.go", "scratch.txt"} {
		if !containsString(sourceErr.Paths, path) {
			t.Fatalf("dirty paths = %v, want %q", sourceErr.Paths, path)
		}
	}
	if !strings.Contains(sourceErr.NextAction, "commit, stash, or remove") {
		t.Fatalf("NextAction = %q, want commit, stash, or remove guidance", sourceErr.NextAction)
	}
}

func TestReleasePlanGitSourceInvalidInputsFailReadOnly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		prepare func(t *testing.T) (string, string, string)
		want    error
	}{
		{
			name: "omitted base without stable tag",
			prepare: func(t *testing.T) (string, string, string) {
				repoDir := newEmptyReleasePlanGitRepo(t)
				writeReleasePlanFile(t, repoDir, "README.md", "seed\n")
				gitReleasePlan(t, repoDir, "add", "-A")
				commitReleasePlan(t, repoDir, "chore: seed")
				return repoDir, "", ""
			},
			want: releaseplan.ErrNoStableReleaseTag,
		},
		{
			name: "malformed explicit base",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				return repo.dir, "release-1", ""
			},
			want: releaseplan.ErrMalformedStableVersion,
		},
		{
			name: "pre-release explicit base",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				return repo.dir, "v0.2.0-rc.1", ""
			},
			want: releaseplan.ErrPrereleaseVersion,
		},
		{
			name: "missing explicit base tag",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				return repo.dir, "v9.9.9", ""
			},
			want: releaseplan.ErrUnresolvedRevision,
		},
		{
			name: "unresolved target",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				return repo.dir, "v0.1.1", "missing-target"
			},
			want: releaseplan.ErrUnresolvedRevision,
		},
		{
			name: "non commit target",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				treeSHA := strings.TrimSpace(gitReleasePlanOutput(t, repo.dir, "rev-parse", "HEAD^{tree}"))
				return repo.dir, "v0.1.1", treeSHA
			},
			want: releaseplan.ErrNonCommitRevision,
		},
		{
			name: "empty range",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				return repo.dir, "v0.1.1", "v0.1.1"
			},
			want: releaseplan.ErrInvalidReleaseRange,
		},
		{
			name: "reversed range",
			prepare: func(t *testing.T) (string, string, string) {
				repo := newReleasePlanGitSourceRepo(t)
				return repo.dir, "v0.1.1", "v0.1.0"
			},
			want: releaseplan.ErrInvalidReleaseRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir, from, to := tt.prepare(t)
			source := newReleasePlanGitSource(repoDir, preflight.ExecGitRunner{})
			before := snapshotReleasePlanRepo(t, repoDir)

			releaseRange, err := source.ResolveRange(context.Background(), from, to)

			assertReleasePlanRepoUnchanged(t, repoDir, before)
			if err == nil {
				t.Fatalf("ResolveRange succeeded with range %+v, want error", releaseRange)
			}
			if !errors.Is(err, tt.want) {
				t.Fatalf("ResolveRange error = %v, want errors.Is(..., %v)", err, tt.want)
			}
			if !reflect.DeepEqual(releaseRange, releaseplan.Range{}) {
				t.Fatalf("ResolveRange returned partial range %+v, want zero range", releaseRange)
			}
		})
	}
}

func TestReleasePlanGitSourceHonorsContextCancellation(t *testing.T) {
	t.Parallel()
	repo := newReleasePlanGitSourceRepo(t)
	source := newReleasePlanGitSource(repo.dir, preflight.ExecGitRunner{})
	before := snapshotReleasePlanRepo(t, repo.dir)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, resolveErr := source.ResolveRange(ctx, "", "")
	_, commitsErr := source.Commits(ctx, releaseplan.Range{})

	assertReleasePlanRepoUnchanged(t, repo.dir, before)
	if !errors.Is(resolveErr, context.Canceled) {
		t.Fatalf("ResolveRange error = %v, want context.Canceled", resolveErr)
	}
	if !errors.Is(commitsErr, context.Canceled) {
		t.Fatalf("Commits error = %v, want context.Canceled", commitsErr)
	}
}

type releasePlanGitRepo struct {
	dir         string
	seedSHA     string
	patchTagSHA string
	featureSHA  string
	fixSHA      string
}

func newReleasePlanGitSourceRepo(t *testing.T) releasePlanGitRepo {
	t.Helper()
	repoDir := newEmptyReleasePlanGitRepo(t)

	writeReleasePlanFile(t, repoDir, "README.md", "seed\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	seedSHA := commitReleasePlan(t, repoDir, "chore: seed")
	gitReleasePlan(t, repoDir, "tag", "v0.1.0")

	gitReleasePlan(t, repoDir, "checkout", "-b", "future-release")
	writeReleasePlanFile(t, repoDir, "future.txt", "future\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	commitReleasePlan(t, repoDir, "feat: future release")
	gitReleasePlan(t, repoDir, "tag", "v9.0.0")
	gitReleasePlan(t, repoDir, "checkout", "main")

	writeReleasePlanFile(t, repoDir, "internal/releaseplan/version.go", "package releaseplan\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	patchTagSHA := commitReleasePlan(t, repoDir, "fix: prior release")
	gitReleasePlan(t, repoDir, "tag", "v0.1.1")

	writeReleasePlanFile(t, repoDir, "cmd/roundfix/main.go", "package main\n")
	writeReleasePlanFile(t, repoDir, "internal/cli/release.go", "package cli\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	featureSHA := commitReleasePlan(t, repoDir, "feat: add release plan", "Release Plan Command reads committed history.")

	writeReleasePlanFile(t, repoDir, "internal/releaseplan/source.go", "package releaseplan\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	fixSHA := commitReleasePlan(t, repoDir, "fix: keep range stable")
	gitReleasePlan(t, repoDir, "tag", "v0.2.0-rc.1")
	gitReleasePlan(t, repoDir, "remote", "add", "origin", "https://example.invalid/roundfix.git")

	return releasePlanGitRepo{
		dir:         repoDir,
		seedSHA:     seedSHA,
		patchTagSHA: patchTagSHA,
		featureSHA:  featureSHA,
		fixSHA:      fixSHA,
	}
}

func newEmptyReleasePlanGitRepo(t *testing.T) string {
	t.Helper()
	repoDir := t.TempDir()
	gittest.InitRepo(t, repoDir, "--initial-branch=main")
	gitReleasePlan(t, repoDir, "config", "user.name", "Roundfix Test")
	gitReleasePlan(t, repoDir, "config", "user.email", "roundfix-test@example.com")
	gitReleasePlan(t, repoDir, "config", "commit.gpgsign", "false")
	return repoDir
}

func writeReleasePlanFile(t *testing.T, repoDir string, name string, content string) {
	t.Helper()
	path := filepath.Join(repoDir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func commitReleasePlan(t *testing.T, repoDir string, subject string, bodyParagraphs ...string) string {
	t.Helper()
	args := []string{"commit", "-m", subject}
	for _, paragraph := range bodyParagraphs {
		args = append(args, "-m", paragraph)
	}
	gitReleasePlan(t, repoDir, args...)
	return strings.TrimSpace(gitReleasePlanOutput(t, repoDir, "rev-parse", "HEAD"))
}

func gitReleasePlan(t *testing.T, dir string, args ...string) {
	t.Helper()
	_ = gitReleasePlanOutput(t, dir, args...)
}

func gitReleasePlanOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	return gitImplementOutput(t, dir, args...)
}

type releasePlanRepoSnapshot struct {
	Files   map[string]string
	Refs    string
	Remotes string
	Config  string
	Status  string
}

func snapshotReleasePlanRepo(t *testing.T, repoDir string) releasePlanRepoSnapshot {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(repoDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(repoDir, path)
		if err != nil {
			return err
		}
		if relative == ".git" {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(relative)] = string(content)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot repository files: %v", err)
	}
	return releasePlanRepoSnapshot{
		Files:   files,
		Refs:    gitReleasePlanOutput(t, repoDir, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads", "refs/tags"),
		Remotes: gitReleasePlanOutput(t, repoDir, "remote", "-v"),
		Config:  gitReleasePlanOutput(t, repoDir, "config", "--local", "--list"),
		Status:  gitReleasePlanOutput(t, repoDir, "--no-optional-locks", "status", "--porcelain=v1", "-z"),
	}
}

func assertReleasePlanRepoUnchanged(t *testing.T, repoDir string, before releasePlanRepoSnapshot) {
	t.Helper()
	after := snapshotReleasePlanRepo(t, repoDir)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("repository changed\nbefore: %+v\nafter:  %+v", before, after)
	}
}

func commitSubjects(commits []releaseplan.Commit) []string {
	subjects := make([]string, 0, len(commits))
	for _, commit := range commits {
		subjects = append(subjects, commit.Subject)
	}
	return subjects
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
