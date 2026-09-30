// Suite: Delivery Queue owner staleness.
// Invariant: an item warns only when its starting main contains Roundfix source changes newer than the queue owner.
// Boundary IN: real disposable Git repositories, Auditor Staleness evidence, and Delivery Revalidation.
// Boundary OUT: the delivery engine records the returned warning in internal/delivery.
package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestRevalidateWarnsWhenTheOwnerPredatesSourceOnMain(t *testing.T) {
	fixture := newDeliveryOwnerStalenessFixture(t, filepath.Join("internal", "owner.go"))
	setOwnerBuildCommit(t, fixture.buildCommit[:8])

	result, err := fixture.workflow.Revalidate(t.Context(), fixture.repository, "clean", []string{})
	if err != nil {
		t.Fatalf("revalidate item start: %v", err)
	}

	want := "owner-older-than-main: owner build " + fixture.buildCommit[:12] +
		" predates starting main " + fixture.startingMain[:12]
	if result.OwnerWarning != want {
		t.Fatalf("owner warning = %q, want %q", result.OwnerWarning, want)
	}
}

func TestRevalidateStaysQuietForDocsOnlyMain(t *testing.T) {
	fixture := newDeliveryOwnerStalenessFixture(t, filepath.Join("docs", "owner.md"))
	setOwnerBuildCommit(t, fixture.buildCommit)

	result, err := fixture.workflow.Revalidate(t.Context(), fixture.repository, "clean", []string{})
	if err != nil {
		t.Fatalf("revalidate docs-only item start: %v", err)
	}
	if result.OwnerWarning != "" {
		t.Fatalf("docs-only owner warning = %q, want empty", result.OwnerWarning)
	}
}

func TestRevalidateStaysQuietForACurrentOwner(t *testing.T) {
	fixture := newDeliveryOwnerStalenessFixture(t, filepath.Join("internal", "owner.go"))
	setOwnerBuildCommit(t, fixture.startingMain)

	result, err := fixture.workflow.Revalidate(t.Context(), fixture.repository, "clean", []string{})
	if err != nil {
		t.Fatalf("revalidate current-owner item start: %v", err)
	}
	if result.OwnerWarning != "" {
		t.Fatalf("current owner warning = %q, want empty", result.OwnerWarning)
	}
}

func TestRevalidateStaysQuietWithoutTheBuildCommit(t *testing.T) {
	fixture := newDeliveryOwnerStalenessFixture(t, filepath.Join("internal", "owner.go"))
	setOwnerBuildCommit(t, strings.Repeat("f", 40))

	result, err := fixture.workflow.Revalidate(t.Context(), fixture.repository, "clean", []string{})
	if err != nil {
		t.Fatalf("revalidate item start without owner commit: %v", err)
	}
	if result.OwnerWarning != "" {
		t.Fatalf("missing-build owner warning = %q, want empty", result.OwnerWarning)
	}
}

func TestRevalidateDoesNotRecomputeTheOwnerWarningOnRetry(t *testing.T) {
	fixture := newDeliveryOwnerStalenessFixture(t, filepath.Join("internal", "owner.go"))
	setOwnerBuildCommit(t, fixture.buildCommit)

	result, err := fixture.workflow.Revalidate(t.Context(), fixture.repository, "clean", nil)
	if err != nil {
		t.Fatalf("revalidate retry: %v", err)
	}
	if result.OwnerWarning != "" {
		t.Fatalf("retry owner warning = %q, want no recomputation", result.OwnerWarning)
	}
}

func TestRevalidateStaysQuietWhenTheOwnerSourceDiffFails(t *testing.T) {
	fixture := newDeliveryOwnerStalenessFixture(t, filepath.Join("internal", "owner.go"))
	setOwnerBuildCommit(t, fixture.buildCommit)
	fixture.workflow.git = ownerDiffFailingGitRunner{delegate: preflight.ExecGitRunner{}}

	result, err := fixture.workflow.Revalidate(t.Context(), fixture.repository, "clean", []string{})
	if err != nil {
		t.Fatalf("revalidate item start with unreadable owner diff: %v", err)
	}
	if result.OwnerWarning != "" {
		t.Fatalf("failed-diff owner warning = %q, want empty", result.OwnerWarning)
	}
}

type deliveryOwnerStalenessFixture struct {
	repository   string
	buildCommit  string
	startingMain string
	workflow     *commandDeliveryWorkflow
}

func newDeliveryOwnerStalenessFixture(t *testing.T, changedPath string) deliveryOwnerStalenessFixture {
	t.Helper()
	revalidationFixture := newDeliveryRevalidationGitFixture(t)
	repository := revalidationFixture.before
	buildCommit := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(repository, changedPath)), 0o755); err != nil {
		t.Fatalf("create owner staleness fixture directory: %v", err)
	}
	mustWrite(t, filepath.Join(repository, changedPath), "owner staleness fixture\n")
	gittest.Run(t, repository, "add", "--", changedPath)
	gittest.Run(t, repository, "commit", "-m", "test: advance starting main")
	startingMain := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))

	return deliveryOwnerStalenessFixture{
		repository:   repository,
		buildCommit:  buildCommit,
		startingMain: startingMain,
		workflow: &commandDeliveryWorkflow{
			loaded: roundconfig.Loaded{
				GitRoot: repository,
				Config:  roundconfig.Config{Specs: roundconfig.Specs{Root: "docs/specs"}},
			},
			git: preflight.ExecGitRunner{},
		},
	}
}

func setOwnerBuildCommit(t *testing.T, commit string) {
	t.Helper()
	previous := app.BuildCommit
	app.BuildCommit = commit
	t.Cleanup(func() {
		app.BuildCommit = previous
	})
}

type ownerDiffFailingGitRunner struct {
	delegate preflight.GitRunner
}

func (runner ownerDiffFailingGitRunner) RunGit(
	ctx context.Context,
	workDir string,
	args ...string,
) (string, error) {
	if len(args) >= 2 && args[0] == "diff" && args[1] == "--name-only" {
		return "", errors.New("owner source diff unavailable")
	}
	return runner.delegate.RunGit(ctx, workDir, args...)
}
