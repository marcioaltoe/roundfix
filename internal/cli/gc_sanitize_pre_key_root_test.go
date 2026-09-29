// Suite: GC sanitation of pre-repository-key Artifact Roots
// Invariant: sanitation accepts only a recorded Run's key-derived or checkout-derived default root.
// Boundary IN: public GC Command, Run Database metadata, real bare Git repositories, and Artifact Directories.
// Boundary OUT: Artifact Directory path derivation, owned by internal/config/artifact_directory_for_path_test.go.
package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func TestGCSanitizeReclaimsAPreKeyDefaultRootInABareLayout(t *testing.T) {
	ctx := context.Background()
	fixture := newBareGCSanitationFixture(t)
	now := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	run := seedBareGCSanitationRun(t, ctx, fixture, fixture.checkoutDefault, "pre-key", now)
	runDir := writeRunArtifact(t, fixture.checkoutDefault, run.ID, "pre-key artifact")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"gc", "sanitize", "--apply"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("gc sanitize --apply exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("gc sanitize --apply stderr = %q, want empty", stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{fixture.checkoutDefault, "Classification: orphaned", fixture.checkoutRoot, "removed: " + runDir} {
		if !strings.Contains(output, want) {
			t.Fatalf("gc sanitize --apply output = %q, want %q", output, want)
		}
	}
	if strings.Contains(output, "Classification: overridden") {
		t.Fatalf("gc sanitize --apply output = %q, pre-key default must not be overridden", output)
	}
	assertPathMissing(t, runDir)
}

func TestGCSanitizeStillPreservesARootEqualToNeitherDefault(t *testing.T) {
	ctx := context.Background()
	fixture := newBareGCSanitationFixture(t)
	now := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	overriddenRoot := filepath.Join(fixture.homeDir, ".roundfix", "artifacts", "operator-override")
	run := seedBareGCSanitationRun(t, ctx, fixture, overriddenRoot, "overridden", now)
	runDir := writeRunArtifact(t, overriddenRoot, run.ID, "overridden artifact")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"gc", "sanitize"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("gc sanitize exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("gc sanitize stderr = %q, want empty", stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{"Classification: overridden", fixture.keyDefault, fixture.checkoutDefault} {
		if !strings.Contains(output, want) {
			t.Fatalf("gc sanitize output = %q, want evidence containing %q", output, want)
		}
	}
	assertPathExists(t, runDir)
}

func TestGCSanitizeKeepsAKeyDerivedRootClassification(t *testing.T) {
	ctx := context.Background()
	fixture := newBareGCSanitationFixture(t)
	now := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	run := seedBareGCSanitationRun(t, ctx, fixture, fixture.keyDefault, "key-derived", now)
	runDir := writeRunArtifact(t, fixture.keyDefault, run.ID, "key-derived artifact")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"gc", "sanitize"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("gc sanitize exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("gc sanitize stderr = %q, want empty", stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, fixture.keyDefault) || !strings.Contains(output, "Classification: orphaned") {
		t.Fatalf("gc sanitize output = %q, want key-derived root classified orphaned", output)
	}
	if strings.Contains(output, "Classification: overridden") {
		t.Fatalf("gc sanitize output = %q, key-derived default must not be overridden", output)
	}
	assertPathExists(t, runDir)
}

func TestGCSanitizeAcceptsTheCheckoutDefaultWhenTheKeyDefaultCannotBeDerived(t *testing.T) {
	homeDir := t.TempDir()
	checkoutRoot := t.TempDir()
	checkoutDefault, err := roundconfig.DefaultArtifactDirectoryForPath(checkoutRoot, homeDir)
	if err != nil {
		t.Fatalf("derive checkout default Artifact Root: %v", err)
	}
	mustMkdir(t, filepath.Join(checkoutDefault, "runs"))

	report := classifyGCSanitationRoot(store.ArtifactRoot{
		Path: checkoutDefault,
		Runs: []store.ArtifactRootRun{{
			ID:         "run_checkout_default",
			GitRoot:    checkoutRoot,
			Repository: "",
			State:      store.StateClean,
		}},
	}, homeDir)

	if report.classification != gcSanitationOrphaned {
		t.Fatalf("classification = %q, want %q; evidence=%q", report.classification, gcSanitationOrphaned, report.evidence)
	}
	if !strings.Contains(report.evidence, checkoutRoot) {
		t.Fatalf("evidence = %q, want recorded checkout %q", report.evidence, checkoutRoot)
	}
}

func TestGCSanitizePreservesUnsafeWhenTheKeyDefaultCannotBeDerived(t *testing.T) {
	homeDir := t.TempDir()
	overriddenRoot := filepath.Join(homeDir, ".roundfix", "artifacts", "operator-override")
	mustMkdir(t, filepath.Join(overriddenRoot, "runs"))

	report := classifyGCSanitationRoot(store.ArtifactRoot{
		Path: overriddenRoot,
		Runs: []store.ArtifactRootRun{{
			ID:         "run_unresolved_key",
			GitRoot:    t.TempDir(),
			Repository: "",
			State:      store.StateClean,
		}},
	}, homeDir)

	if report.classification != gcSanitationUnsafe {
		t.Fatalf("classification = %q, want %q; evidence=%q", report.classification, gcSanitationUnsafe, report.evidence)
	}
	if !strings.Contains(report.evidence, "repository-key default Artifact Root") {
		t.Fatalf("evidence = %q, want repository-key derivation failure", report.evidence)
	}
}

type bareGCSanitationFixture struct {
	homeDir         string
	bareRoot        string
	checkoutRoot    string
	keyDefault      string
	checkoutDefault string
}

func newBareGCSanitationFixture(t *testing.T) bareGCSanitationFixture {
	t.Helper()
	fixtureRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("resolve fixture root: %v", err)
	}
	fixture := bareGCSanitationFixture{
		homeDir:      filepath.Join(fixtureRoot, "home"),
		bareRoot:     filepath.Join(fixtureRoot, "repository.git"),
		checkoutRoot: filepath.Join(fixtureRoot, "checkout"),
	}
	seedRoot := filepath.Join(fixtureRoot, "seed")
	gittest.InitRepo(t, seedRoot, "--initial-branch=main")
	gittest.Run(t, seedRoot, "commit", "--allow-empty", "-m", "seed repository")
	gittest.Run(t, seedRoot, "clone", "--bare", seedRoot, fixture.bareRoot)
	gittest.Harden(t, fixture.bareRoot)
	gittest.Run(t, fixture.bareRoot, "worktree", "add", "-b", "feature/gc-sanitize", fixture.checkoutRoot, "main")
	setCommandEnvironmentForTest(t, fixture.homeDir, fixture.checkoutRoot)
	mustWrite(t, filepath.Join(fixture.checkoutRoot, ".roundfixrc.yml"), "store:\n  journal_retention: 336h\n")

	fixture.keyDefault, err = roundconfig.ResolveArtifactDirectory("", fixture.checkoutRoot, fixture.homeDir)
	if err != nil {
		t.Fatalf("resolve repository-key default Artifact Root: %v", err)
	}
	fixture.checkoutDefault, err = roundconfig.DefaultArtifactDirectoryForPath(fixture.checkoutRoot, fixture.homeDir)
	if err != nil {
		t.Fatalf("derive checkout default Artifact Root: %v", err)
	}
	if fixture.keyDefault == fixture.checkoutDefault {
		t.Fatalf("bare layout defaults unexpectedly match: %q", fixture.keyDefault)
	}
	return fixture
}

func seedBareGCSanitationRun(
	t *testing.T,
	ctx context.Context,
	fixture bareGCSanitationFixture,
	artifactRoot string,
	name string,
	now time.Time,
) store.Run {
	t.Helper()
	runStore, err := store.Open(ctx, fixture.homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	run := createGCSanitationRun(t, ctx, runStore, fixture.checkoutRoot, artifactRoot, name, true)
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run store: %v", err)
	}
	setRunTimestamps(t, fixture.homeDir, run.ID, now.Add(-400*time.Hour), now.Add(-400*time.Hour))
	return run
}
