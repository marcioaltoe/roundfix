// Suite: Public Baseline release journeys
// Invariant: a built Roundfix binary applies only one exact reviewed Baseline Plan and composes with repository-owned tools without hidden execution.
// Boundary IN: real process plan/apply handoff, greenfield, preservation, update, profile change, stale and cross-clone plans, unsafe carriers, external composition, and empty reapply.
// Boundary OUT: transaction fault injection (internal/baseline/release_gate_test.go), sealed ACP supervision (internal/baselineacp/release_gate_test.go), and the separately authorized live Fluxus repository.

package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/internal/gittest"
)

func TestGuidanceCompositionJourney(t *testing.T) {
	t.Parallel()
	binary := buildBaselineReleaseBinary(t)
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded Baseline catalog: %v", err)
	}
	profiles := []string{"go-cli-tui", "rust-cli", "standard-typescript-monorepo"}
	if got := fmt.Sprint(catalog.ProfileIDs()); got != fmt.Sprint(profiles) {
		t.Fatalf("maintained Profiles = %s, want %s", got, fmt.Sprint(profiles))
	}

	for _, profile := range profiles {
		t.Run(profile, func(t *testing.T) {
			testBaselineReleaseProfileJourney(t, binary, profile)
		})
	}
}

func TestProjectDecisionJourney(t *testing.T) {
	t.Parallel()
	t.Run("human and automation answers produce one Plan", TestProjectDecisionParity)
	t.Run("compatible decisions are reused on update", TestProjectDecisionReuse)
	t.Run("Fluxus-style defaults retain the persisted HTTP exception", TestBetterAuthSuggestionReusesFullHTTPException)

	binary := buildBaselineReleaseBinary(t)
	t.Run("affected Profile apply audit and empty reapply", func(t *testing.T) {
		testBaselineReleaseProfileJourney(t, binary, "standard-typescript-monorepo")
	})
	t.Run("missing automation input stops without mutation", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "standard-typescript-monorepo")
		before := baselineReleaseVisibleState(t, repo)
		args := []string{
			"baseline", "plan",
			"--repo", repo,
			"--profile", "standard-typescript-monorepo",
		}
		for _, decision := range baselineReleaseDecisionArgs(
			"standard-typescript-monorepo",
			"greenfield",
		) {
			if strings.HasPrefix(decision, "identifier.strategy=") {
				continue
			}
			args = append(args, "--decision", decision)
		}
		args = append(args, "--format=json")
		code, stdout, stderr := runBaselineReleaseCLI(t, binary, args...)
		if code != exitUnverified || len(stderr) != 0 {
			t.Fatalf("missing decision exit=%d stdout=%s stderr=%s", code, stdout, stderr)
		}
		result, err := baseline.ParseResult(stdout)
		if err != nil ||
			result.State != "action_required" ||
			result.Category != "decision" ||
			!strings.Contains(result.Message, "identifier.strategy") {
			t.Fatalf("missing identifier result = %+v error=%v", result, err)
		}
		if after := baselineReleaseVisibleState(t, repo); after != before {
			t.Fatal("missing automation input changed repository state")
		}
	})
}

func testBaselineReleaseProfileJourney(t *testing.T, binary, profile string) {
	t.Helper()
	repo := newBaselineReleaseRepository(t, profile)
	plan, planPath := baselineReleasePlan(t, binary, repo, profile)
	if plan.Profile.ID != profile {
		t.Fatalf("planned Profile = %q, want %q", plan.Profile.ID, profile)
	}
	assertBaselineReleaseCompleteManagedLedger(t, plan)
	assertBaselineReleaseNoRepositoryCarrier(t, repo, plan)

	first := baselineReleaseApply(t, binary, repo, plan, planPath)
	assertBaselineReleaseRecommendation(t, first, "make verify")
	if profile == "standard-typescript-monorepo" {
		assertBaselineReleaseRecommendation(t, first, "bun run fmt")
		assertBaselineReleaseNoRecommendation(t, first, "bun run format")
	}
	assertBaselineReleaseNoRepositoryCarrier(t, repo, plan)
	for _, marker := range []string{".qa-verification-ran", ".qa-profile-command-ran"} {
		if _, err := os.Stat(filepath.Join(repo, marker)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Baseline executed repository command marker %s: %v", marker, err)
		}
	}

	update, updatePath := baselineReleasePlan(t, binary, repo, profile)
	if len(update.FileChanges) != 0 {
		t.Fatalf(
			"first fresh Plan after greenfield apply has %d file changes: %+v",
			len(update.FileChanges),
			update.FileChanges,
		)
	}
	assertBaselineReleaseCompleteManagedLedger(t, update)
	baselineReleaseApply(t, binary, repo, update, updatePath)
	assertBaselineReleaseNoRepositoryCarrier(t, repo, update)

	managedBefore := baselineReleaseManagedState(t, repo, update)
	runBaselineReleaseFormatter(t, repo, profile)
	runBaselineReleaseExternal(t, repo, "make", "verify")
	if _, err := os.Stat(filepath.Join(repo, ".qa-verification-ran")); err != nil {
		t.Fatalf("external repository Verification did not run: %v", err)
	}
	if managedAfter := baselineReleaseManagedState(t, repo, update); managedAfter != managedBefore {
		t.Fatalf(
			"formatter or repository Verification changed managed state:\nbefore=%s\nafter=%s",
			managedBefore,
			managedAfter,
		)
	}

	fresh, freshPath := baselineReleasePlan(t, binary, repo, profile)
	if len(fresh.FileChanges) != 0 {
		t.Fatalf("fresh audit Plan has %d file changes: %+v", len(fresh.FileChanges), fresh.FileChanges)
	}
	beforeReapply := baselineReleaseVisibleState(t, repo)
	second := baselineReleaseApply(t, binary, repo, fresh, freshPath)
	if second.State != "verified" || !strings.Contains(second.Message, "already applied") {
		t.Fatalf("empty update apply result = %+v", second)
	}
	if afterReapply := baselineReleaseVisibleState(t, repo); afterReapply != beforeReapply {
		t.Fatalf("empty update apply changed repository state")
	}
}

func TestSemanticRedistributionJourney(t *testing.T) {
	t.Parallel()
	binary := buildBaselineReleaseBinary(t)
	for _, carrier := range []string{
		"docs/agents/repository.md",
		"docs/agents/repository-rules.md",
	} {
		t.Run("zero residual from "+filepath.Base(carrier), func(t *testing.T) {
			repo := newBaselineReleaseRepository(t, "go-cli-tui")
			rule := []byte("Keep requested CLI output on stdout.\n")
			writeBaselinePlanTestFile(t, repo, carrier, string(rule))
			commitBaselinePlanTestRepository(t, repo)

			decisionPath, sourceIDs := baselineReleaseReadoptionDecisionFile(
				t,
				repo,
				"go-cli-tui",
				func(entry baseline.ReadoptionSourceEntry) baseline.ReadoptionDisposition {
					return baseline.ReadoptionDisposition{
						EntryID:        entry.ID,
						EntryDigest:    entry.Digest,
						Classification: "normative-clause",
						Disposition:    "repository-document",
						Destination: &baseline.ReadoptionDestination{
							DocumentType: "agent-guide",
							Path:         "docs/agents/cli.md",
							Digest:       entry.Digest,
						},
						Reason: "The active CLI guide owns this repository policy.",
					}
				},
			)
			plan, planPath := baselineReleasePlanWithDecisionFile(
				t,
				binary,
				repo,
				"go-cli-tui",
				decisionPath,
			)
			assertBaselineReleaseCompleteManagedLedger(t, plan)
			assertBaselineReleaseRetentionLedger(
				t,
				plan,
				sourceIDs,
				"repository-document",
				"docs/agents/cli.md",
			)
			baselineReleaseApply(t, binary, repo, plan, planPath)
			assertBaselineReleaseNoRepositoryCarrier(t, repo, plan)

			guide, err := os.ReadFile(filepath.Join(repo, "docs", "agents", "cli.md"))
			if err != nil || !bytes.Contains(guide, rule) {
				t.Fatalf("semantic guide does not contain exact redistributed bytes: error=%v\n%s", err, guide)
			}
		})
	}

	t.Run("non-empty residual is retained canonically", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		const legacyCarrier = "docs/agents/repository.md"
		rule := []byte("Keep the repository-specific release name stable.\n")
		writeBaselinePlanTestFile(t, repo, legacyCarrier, string(rule))
		commitBaselinePlanTestRepository(t, repo)

		decisionPath, sourceIDs := baselineReleaseReadoptionDecisionFile(
			t,
			repo,
			"go-cli-tui",
			func(entry baseline.ReadoptionSourceEntry) baseline.ReadoptionDisposition {
				return baseline.ReadoptionDisposition{
					EntryID:        entry.ID,
					EntryDigest:    entry.Digest,
					Classification: "normative-clause",
					Disposition:    "repository-rules",
					Destination: &baseline.ReadoptionDestination{
						DocumentType:  "repository-rules",
						Path:          "docs/agents/specific-repository.md",
						Digest:        entry.Digest,
						ProposedBytes: base64.StdEncoding.EncodeToString(entry.SourceBytes),
					},
					Reason: "No active semantic guide owns this repository-specific policy.",
				}
			},
		)
		plan, planPath := baselineReleasePlanWithDecisionFile(
			t,
			binary,
			repo,
			"go-cli-tui",
			decisionPath,
		)
		assertBaselineReleaseCompleteManagedLedger(t, plan)
		assertBaselineReleaseRetentionLedger(
			t,
			plan,
			sourceIDs,
			"repository-rules",
			"docs/agents/specific-repository.md",
		)
		baselineReleaseApply(t, binary, repo, plan, planPath)

		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(legacyCarrier))); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("legacy carrier remains after residual migration: %v", err)
		}
		residual, err := os.ReadFile(filepath.Join(repo, "docs", "agents", "specific-repository.md"))
		if err != nil || !bytes.Equal(residual, rule) {
			t.Fatalf("canonical residual bytes = %q error=%v, want %q", residual, err, rule)
		}
		root, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
		if err != nil || !bytes.Contains(root, []byte("docs/agents/specific-repository.md")) {
			t.Fatalf("root residual pointer missing: error=%v\n%s", err, root)
		}
	})
}

func TestProfileAdaptationJourney(t *testing.T) {
	t.Parallel()
	binary := buildBaselineReleaseBinary(t)
	repo, input, decisions := baselinePlanProfileFileFixture(t)
	seedBaselineReleaseSourceManifest(t, repo, input.SourceProfileID)
	draftPath := filepath.Join(t.TempDir(), "guided-backend.json")
	if err := os.WriteFile(draftPath, input.Document, 0o600); err != nil {
		t.Fatalf("write reviewed Profile adaptation: %v", err)
	}

	code, stdout, stderr := runBaselineReleaseCLI(
		t,
		binary,
		baselinePlanProfileFileArgs(repo, draftPath, decisions)...,
	)
	if code != exitOK {
		t.Fatalf("Profile adaptation plan exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	plan, planPath := baselineReleasePlanFile(t, stdout)
	if plan.Profile.ID != "guided-backend" ||
		plan.Profile.Source != baseline.ProfileSourceRepository {
		t.Fatalf("reviewed repository-owned Profile = %+v", plan.Profile)
	}
	assertBaselineReleaseCompleteManagedLedger(t, plan)
	assertBaselineReleaseProfileLedger(t, plan)
	assertBaselineReleaseNoRepositoryCarrier(t, repo, plan)

	first := baselineReleaseApply(t, binary, repo, plan, planPath)
	assertBaselineReleaseRecommendation(t, first, "make verify")
	profileBytes, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(plan.Profile.Path)))
	if err != nil || len(profileBytes) == 0 {
		t.Fatalf("applied repository-owned Profile = %q error=%v", profileBytes, err)
	}

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	alignment, err := baseline.ResolveProfileAlignment(
		context.Background(),
		repo,
		baseline.ProfileAlignmentRequest{
			ProfileID:            plan.Profile.ID,
			Decisions:            plan.Decisions,
			Profile:              &plan.Profile,
			RemediationProfileID: input.SourceProfileID,
		},
		catalog,
	)
	if err != nil {
		t.Fatalf("audit applied Profile adaptation: %v", err)
	}
	if !alignment.Ready {
		t.Fatalf("applied Profile adaptation remains divergent: %+v", alignment.Divergences)
	}
	for _, capabilityID := range []string{"capability.context7", "capability.exa"} {
		assertBaselineReleaseUniversalCapability(t, alignment, capabilityID)
	}

	code, stdout, stderr = runBaselineReleaseCLI(
		t,
		binary,
		baselinePlanProfileFileArgs(repo, draftPath, decisions)...,
	)
	if code != exitOK {
		t.Fatalf("Profile update plan exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	update, updatePath := baselineReleasePlanFile(t, stdout)
	assertBaselineReleaseCompleteManagedLedger(t, update)
	baselineReleaseApply(t, binary, repo, update, updatePath)

	runBaselineReleaseExternal(t, repo, "make", "verify")
	code, stdout, stderr = runBaselineReleaseCLI(
		t,
		binary,
		baselinePlanProfileFileArgs(repo, draftPath, decisions)...,
	)
	if code != exitOK {
		t.Fatalf("fresh Profile plan exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	fresh, freshPath := baselineReleasePlanFile(t, stdout)
	if len(fresh.FileChanges) != 0 {
		t.Fatalf("fresh Profile Plan has %d file changes: %+v", len(fresh.FileChanges), fresh.FileChanges)
	}
	beforeReapply := baselineReleaseVisibleState(t, repo)
	second := baselineReleaseApply(t, binary, repo, fresh, freshPath)
	if second.State != "verified" || !strings.Contains(second.Message, "already applied") {
		t.Fatalf("empty Profile reapply result = %+v", second)
	}
	if afterReapply := baselineReleaseVisibleState(t, repo); afterReapply != beforeReapply {
		t.Fatal("empty Profile reapply changed repository state")
	}
}

func TestBaselineMacroJourneysPublicCLI(t *testing.T) {
	t.Parallel()
	binary := buildBaselineReleaseBinary(t)

	t.Run("greenfield automation apply and empty reapply", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		plan, planPath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		first := baselineReleaseApply(t, binary, repo, plan, planPath)
		if first.State != "verified" {
			t.Fatalf("first apply state = %q, want verified", first.State)
		}
		before := baselineReleaseManagedState(t, repo, plan)
		second := baselineReleaseApply(t, binary, repo, plan, planPath)
		if second.State != "verified" || !strings.Contains(second.Message, "already applied") {
			t.Fatalf("empty reapply result = %+v", second)
		}
		if after := baselineReleaseManagedState(t, repo, plan); after != before {
			t.Fatalf("empty reapply changed managed state:\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("preservation", func(t *testing.T) {
		const repositoryRule = "Keep this repository-specific release rule.\n"
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		writeBaselinePlanTestFile(t, repo, "AGENTS.md", repositoryRule)
		commitBaselinePlanTestRepository(t, repo)
		decisionPath := baselineReleasePreservationDecisionFile(t, repo)
		plan, planPath := baselineReleasePlanWithDecisionFile(
			t, binary, repo, "go-cli-tui", decisionPath,
		)
		baselineReleaseApply(t, binary, repo, plan, planPath)

		rules, err := os.ReadFile(filepath.Join(repo, "docs", "agents", "specific-repository.md"))
		if err != nil || !bytes.Contains(rules, []byte(strings.TrimSpace(repositoryRule))) {
			t.Fatalf("preserved repository rule = %q error=%v", rules, err)
		}
		var backup baseline.ManagedEntry
		for _, entry := range plan.ManagedEntries {
			if entry.Kind == "backup" {
				backup = entry
				break
			}
		}
		content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(backup.Path)))
		if err != nil || string(content) != repositoryRule {
			t.Fatalf("immutable preservation backup %q = %q error=%v", backup.Path, content, err)
		}
	})

	t.Run("update and profile change", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		goPlan, goPlanPath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		baselineReleaseApply(t, binary, repo, goPlan, goPlanPath)

		updatePlan, updatePath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		baselineReleaseApply(t, binary, repo, updatePlan, updatePath)

		rustPlan, rustPlanPath := baselineReleasePlan(t, binary, repo, "rust-cli")
		if rustPlan.Profile.ID != "rust-cli" || rustPlan.PlanDigest == goPlan.PlanDigest {
			t.Fatalf("profile-change plan = profile %q digest %q, original %q",
				rustPlan.Profile.ID, rustPlan.PlanDigest, goPlan.PlanDigest)
		}
		baselineReleaseApply(t, binary, repo, rustPlan, rustPlanPath)
		manifest, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(baselineSetupManifestPath)))
		if err != nil || !bytes.Contains(manifest, []byte(`"profile": "rust-cli"`)) {
			t.Fatalf("profile-change manifest = %q error=%v", manifest, err)
		}
	})

	t.Run("stale plan", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		plan, planPath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@echo changed\n")
		before := baselineReleaseVisibleState(t, repo)
		code, stdout, _ := runBaselineReleaseCLI(
			t, binary,
			"baseline", "apply",
			"--repo", repo,
			"--plan", planPath,
			"--confirm-plan", plan.PlanDigest,
			"--format=json",
		)
		if code != exitUnverified {
			t.Fatalf("stale apply exit = %d stdout=%s, want %d", code, stdout, exitUnverified)
		}
		result, err := baseline.ParseResult(stdout)
		if err != nil || result.Category != "stale" {
			t.Fatalf("stale apply result = %+v error=%v", result, err)
		}
		if after := baselineReleaseVisibleState(t, repo); after != before {
			t.Fatalf("stale apply changed repository:\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("cross-clone apply", func(t *testing.T) {
		source := newBaselineReleaseRepository(t, "go-cli-tui")
		plan, planPath := baselineReleasePlan(t, binary, source, "go-cli-tui")
		parent := t.TempDir()
		clone := filepath.Join(parent, "matching")
		runBaselinePlanTestCommand(t, parent, "git", "clone", "--quiet", source, clone)
		gittest.Harden(t, clone)
		result := baselineReleaseApply(t, binary, clone, plan, planPath)
		if result.State != "verified" {
			t.Fatalf("cross-clone apply result = %+v", result)
		}
	})

	t.Run("unsafe carrier", func(t *testing.T) {
		repo := newBaselinePlanTestRepository(t)
		if err := os.Symlink("../outside.md", filepath.Join(repo, "AGENTS.md")); err != nil {
			t.Fatal(err)
		}
		commitBaselinePlanTestRepository(t, repo)
		before := baselineReleaseVisibleState(t, repo)
		code, stdout, _ := runBaselineReleaseCLI(
			t, binary, "baseline", "plan", "--repo", repo, "--format=json",
		)
		if code != exitPreflight {
			t.Fatalf("unsafe plan exit = %d stdout=%s, want %d", code, stdout, exitPreflight)
		}
		result, err := baseline.ParseResult(stdout)
		if err != nil || result.Category != "preflight" {
			t.Fatalf("unsafe carrier result = %+v error=%v", result, err)
		}
		if after := baselineReleaseVisibleState(t, repo); after != before {
			t.Fatalf("unsafe planning changed repository:\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("interactive and automation parity", TestHumanAutomationPlanParity)
	t.Run("consolidated preservation review", TestConsolidatedReview)
	t.Run("rejected-plan revision", TestRejectedPlanRevision)
}

func TestBaselineFindingRegressionsHumanReview(t *testing.T) {
	t.Parallel()
	t.Run("decisions use one consolidated review", TestConsolidatedReview)
	t.Run("file projection precedes canonical ledgers", TestRejectedPlanRevision)
}

func TestBaselineDocumentationContractExamples(t *testing.T) {
	t.Parallel()
	t.Run("public examples parse", TestBaselineExamplesParse)
	t.Run("Decision Documents parse", TestBaselineDecisionExamples)
}

// coldBuiltBinary compiles the Roundfix CLI once per package run, from an
// empty build cache, and hands every caller the same path.
//
// The empty cache is deliberate: these tests exercise the binary a release
// would ship, so it has to compile from scratch rather than inherit whatever
// the developer's cache happens to hold. What was not deliberate was paying
// that cold compile seven times. Seven callers each built the whole project
// into their own t.TempDir() with their own empty GOCACHE — roughly 160s of
// the package's 437s of serial work — and because they run in parallel they
// saturated every core available. That is why the package took the same wall
// clock on two cores as on twelve, and why raising -parallel changed nothing:
// the package was not waiting on scheduling, it was compiling itself seven
// times over.
//
// Callers only exec the binary, never modify it, so one copy serves them all.
var coldBuiltBinary = struct {
	once sync.Once
	dir  string
	path string
	err  error
}{}

func buildBaselineReleaseBinary(t *testing.T) string {
	t.Helper()
	coldBuiltBinary.once.Do(func() {
		projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			coldBuiltBinary.err = err
			return
		}
		dir, err := os.MkdirTemp("", "roundfix-cold-build")
		if err != nil {
			coldBuiltBinary.err = err
			return
		}
		coldBuiltBinary.dir = dir
		binary := filepath.Join(dir, "roundfix")
		command := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
		command.Dir = projectRoot
		command.Env = append(os.Environ(), "GOCACHE="+filepath.Join(dir, "go-cache"))
		if output, err := command.CombinedOutput(); err != nil {
			coldBuiltBinary.err = fmt.Errorf("build Roundfix binary from an empty cache: %w\n%s", err, output)
			return
		}
		coldBuiltBinary.path = binary
	})
	if coldBuiltBinary.err != nil {
		t.Fatal(coldBuiltBinary.err)
	}
	return coldBuiltBinary.path
}

// removeColdBuiltBinary drops the shared build directory after the package's
// tests finish. It is called from TestMain, because the directory outlives
// every individual test that used it.
func removeColdBuiltBinary() {
	if coldBuiltBinary.dir != "" {
		_ = os.RemoveAll(coldBuiltBinary.dir)
	}
}

func baselineReleasePlan(
	t *testing.T,
	binary string,
	repo string,
	profile string,
) (baseline.PlanDocument, string) {
	t.Helper()
	args := []string{"baseline", "plan", "--repo", repo, "--profile", profile}
	for _, decision := range baselineReleaseDecisionArgs(profile, "greenfield") {
		args = append(args, "--decision", decision)
	}
	args = append(args, "--format=json")
	code, stdout, stderr := runBaselineReleaseCLI(t, binary, args...)
	if code != exitOK {
		t.Fatalf("release plan exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	return baselineReleasePlanFile(t, stdout)
}

func baselineReleasePlanWithDecisionFile(
	t *testing.T,
	binary string,
	repo string,
	profile string,
	decisionPath string,
) (baseline.PlanDocument, string) {
	t.Helper()
	code, stdout, stderr := runBaselineReleaseCLI(
		t,
		binary,
		"baseline", "plan",
		"--repo", repo,
		"--profile", profile,
		"--decision-file", decisionPath,
		"--format=json",
	)
	if code != exitOK {
		t.Fatalf("preservation plan exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	return baselineReleasePlanFile(t, stdout)
}

func baselineReleasePlanFile(t *testing.T, data []byte) (baseline.PlanDocument, string) {
	t.Helper()
	plan, err := baseline.ParsePlanDocument(data)
	if err != nil {
		t.Fatalf("parse release-gate Plan: %v\n%s", err, data)
	}
	planPath := filepath.Join(t.TempDir(), "baseline-plan.json")
	if err := os.WriteFile(planPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return plan, planPath
}

func baselineReleaseApply(
	t *testing.T,
	binary string,
	repo string,
	plan baseline.PlanDocument,
	planPath string,
) baseline.Result {
	t.Helper()
	code, stdout, stderr := runBaselineReleaseCLI(
		t,
		binary,
		"baseline", "apply",
		"--repo", repo,
		"--plan", planPath,
		"--confirm-plan", plan.PlanDigest,
		"--format=json",
	)
	if code != exitOK {
		t.Fatalf("release apply exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result, err := baseline.ParseResult(stdout)
	if err != nil {
		t.Fatalf("parse release apply result: %v\n%s", err, stdout)
	}
	return result
}

func baselineReleasePreservationDecisionFile(t *testing.T, repo string) string {
	t.Helper()
	inspection, err := baseline.InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	preservation, err := baseline.PlanRootPreservation(
		inspection,
		baseline.RootPreservationRequest{Mode: baseline.PreservationModePreservation},
	)
	if err != nil || preservation.DecisionSkeleton == nil {
		t.Fatalf("build preservation Decision Document: plan=%+v error=%v", preservation, err)
	}
	document := preservation.DecisionSkeleton.Document
	for _, raw := range baselineReleaseDecisionArgs("go-cli-tui", "preservation") {
		id, value, ok := strings.Cut(raw, "=")
		if !ok {
			t.Fatalf("invalid release decision %q", raw)
		}
		document.Decisions = append(document.Decisions, baseline.DecisionValue{
			ID: id, Value: baselineReleaseDecisionValue(value),
		})
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "preservation-decisions.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func baselineReleaseDecisionValue(raw string) any {
	switch raw {
	case "true":
		return true
	case "false":
		return false
	}
	if strings.HasPrefix(raw, "{") {
		var value any
		if json.Unmarshal([]byte(raw), &value) == nil {
			return value
		}
	}
	return raw
}

func baselineReleaseReadoptionDecisionFile(
	t *testing.T,
	repo string,
	profile string,
	disposition func(baseline.ReadoptionSourceEntry) baseline.ReadoptionDisposition,
) (string, []string) {
	t.Helper()
	inspection, err := baseline.InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect Readoption repository: %v", err)
	}
	preservation, err := baseline.PlanRootPreservation(
		inspection,
		baseline.RootPreservationRequest{Mode: baseline.PreservationModePreservation},
	)
	if err != nil {
		t.Fatalf("inventory Readoption source: %v", err)
	}
	if preservation.DecisionSkeleton == nil || len(preservation.SourceBaseline.Entries) == 0 {
		t.Fatalf("Readoption inventory has no editable source: %+v", preservation)
	}
	document := preservation.DecisionSkeleton.Document
	document.Readoption.Dispositions = make(
		[]baseline.ReadoptionDisposition,
		0,
		len(preservation.SourceBaseline.Entries),
	)
	sourceIDs := make([]string, 0, len(preservation.SourceBaseline.Entries))
	for _, entry := range preservation.SourceBaseline.Entries {
		document.Readoption.Dispositions = append(
			document.Readoption.Dispositions,
			disposition(entry),
		)
		sourceIDs = append(sourceIDs, entry.ID)
	}
	for _, raw := range baselineReleaseDecisionArgs(profile, "preservation") {
		id, value, ok := strings.Cut(raw, "=")
		if !ok {
			t.Fatalf("invalid release decision %q", raw)
		}
		document.Decisions = append(document.Decisions, baseline.DecisionValue{
			ID: id, Value: baselineReleaseDecisionValue(value),
		})
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatalf("marshal Readoption Decision Document: %v", err)
	}
	decisionPath := filepath.Join(t.TempDir(), "readoption-decisions.json")
	if err := os.WriteFile(decisionPath, append(data, '\n'), 0o600); err != nil {
		t.Fatalf("write Readoption Decision Document: %v", err)
	}
	return decisionPath, sourceIDs
}

func baselineReleaseDecisionArgs(profile, preservation string) []string {
	repositoryExtension := "false"
	if preservation == "preservation" {
		repositoryExtension = "true"
	}
	decisions := []string{
		"preservation.mode=" + preservation,
		"language.generated=English",
		"verification.gate=make verify",
		"spec.scaffold=true",
		"domain.layout=single-context",
		"triage.external=false",
		"autonomous.enabled=false",
		"secondbrain.enabled=false",
		"repository.extension.enabled=" + repositoryExtension,
	}
	if profile == "standard-typescript-monorepo" {
		decisions = append(
			decisions,
			`identifier.strategy={"kind":"uuid-v7"}`,
			`http.contract={"mode":"REST"}`,
			`auth.provider={"kind":"better-auth","routeException":{"scope":"/api/auth/*","methods":["GET","POST"],"owner":"Better Auth","reason":"Provider protocol routes require GET and POST semantics."}}`,
		)
	}
	return decisions
}

func assertBaselineReleaseCompleteManagedLedger(t *testing.T, plan baseline.PlanDocument) {
	t.Helper()
	if len(plan.ManagedEntries) == 0 {
		t.Fatalf(
			"incomplete Change Plan ledgers: managed=%d fileChanges=%d",
			len(plan.ManagedEntries),
			len(plan.FileChanges),
		)
	}
	entries := make(map[string]baseline.ManagedEntry, len(plan.ManagedEntries))
	for index, entry := range plan.ManagedEntries {
		if entry.Ordinal != index || entry.ID == "" || entry.Path == "" {
			t.Fatalf("managed-entry ledger item %d = %+v", index, entry)
		}
		if _, duplicate := entries[entry.ID]; duplicate {
			t.Fatalf("duplicate managed-entry ledger ID %q", entry.ID)
		}
		entries[entry.ID] = entry
	}
	projected := make(map[string]int, len(entries))
	for _, change := range plan.FileChanges {
		if change.Path == "" || len(change.ManagedEntries) == 0 {
			t.Fatalf("file-change projection is incomplete: %+v", change)
		}
		for _, id := range change.ManagedEntries {
			entry, ok := entries[id]
			if !ok || entry.Path != change.Path {
				t.Fatalf("file-change projection %q references invalid entry %q", change.Path, id)
			}
			projected[id]++
		}
	}
	for id := range entries {
		switch projected[id] {
		case 1:
			continue
		case 0:
			entry := entries[id]
			if !baselineReleasePathIsUnchanged(plan, entry.Path) {
				t.Fatalf("unprojected managed entry %q at %q is not unchanged", id, entry.Path)
			}
		default:
			t.Fatalf("managed entry %q appears %d times in file-change projection", id, projected[id])
		}
	}
}

func baselineReleasePathIsUnchanged(plan baseline.PlanDocument, path string) bool {
	var before baseline.Preimage
	for _, candidate := range plan.Preimages {
		if candidate.Path == path {
			before = candidate
			break
		}
	}
	var after baseline.Postimage
	for _, candidate := range plan.Postimages {
		if candidate.Path == path {
			after = candidate
			break
		}
	}
	if before.Exists {
		return before.Kind == after.Kind && before.ContentIdentity == after.ContentIdentity
	}
	return after.Kind == baseline.PreimageMissing
}

func assertBaselineReleaseRetentionLedger(
	t *testing.T,
	plan baseline.PlanDocument,
	sourceIDs []string,
	disposition string,
	target string,
) {
	t.Helper()
	if len(plan.Retention) != len(sourceIDs) {
		t.Fatalf(
			"Upgrade Retention Contract has %d entries, want %d: %+v",
			len(plan.Retention),
			len(sourceIDs),
			plan.Retention,
		)
	}
	evidence := make(map[string]baseline.RetentionEvidence, len(plan.Retention))
	for _, entry := range plan.Retention {
		if _, duplicate := evidence[entry.FromClause]; duplicate {
			t.Fatalf("duplicate Upgrade Retention Contract source %q", entry.FromClause)
		}
		evidence[entry.FromClause] = entry
	}
	for _, sourceID := range sourceIDs {
		entry, ok := evidence[sourceID]
		if !ok ||
			entry.Disposition != disposition ||
			len(entry.Targets) != 1 ||
			entry.Targets[0] != target ||
			strings.TrimSpace(entry.Reason) == "" {
			t.Fatalf("Upgrade Retention Contract source %q = %+v", sourceID, entry)
		}
	}
}

func assertBaselineReleaseNoRepositoryCarrier(
	t *testing.T,
	repo string,
	plan baseline.PlanDocument,
) {
	t.Helper()
	for _, carrier := range []string{
		"docs/agents/specific-repository.md",
		"docs/agents/repository.md",
		"docs/agents/repository-rules.md",
	} {
		for _, postimage := range plan.Postimages {
			if postimage.Path == carrier && postimage.Kind == baseline.PreimageRegular {
				t.Fatalf("zero-residual Plan contains repository carrier %q", carrier)
			}
		}
		if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(carrier))); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("zero-residual repository contains carrier %q: %v", carrier, err)
		}
	}
	for _, postimage := range plan.Postimages {
		if postimage.Path == "AGENTS.md" {
			for _, carrier := range []string{
				"docs/agents/specific-repository.md",
				"docs/agents/repository.md",
				"docs/agents/repository-rules.md",
			} {
				if bytes.Contains(postimage.Content, []byte(carrier)) {
					t.Fatalf("zero-residual root postimage contains pointer %q", carrier)
				}
			}
		}
	}
}

func assertBaselineReleaseProfileLedger(t *testing.T, plan baseline.PlanDocument) {
	t.Helper()
	entryID := "profile:" + plan.Profile.ID
	for _, entry := range plan.ManagedEntries {
		if entry.ID != entryID {
			continue
		}
		if entry.Kind != "profile" ||
			entry.Path != plan.Profile.Path ||
			entry.AfterIdentity == "" ||
			entry.AfterIdentity != entry.ContentIdentity {
			t.Fatalf("repository-owned Profile ledger entry = %+v", entry)
		}
		for _, change := range plan.FileChanges {
			if change.Path == entry.Path {
				for _, managedID := range change.ManagedEntries {
					if managedID == entryID {
						return
					}
				}
			}
		}
		t.Fatalf("repository-owned Profile %q is absent from file-change projection", entryID)
	}
	t.Fatalf("repository-owned Profile %q is absent from managed-entry ledger", entryID)
}

func assertBaselineReleaseUniversalCapability(
	t *testing.T,
	alignment baseline.ProfileAlignment,
	capabilityID string,
) {
	t.Helper()
	for _, outcome := range alignment.Capabilities {
		if outcome.ID != capabilityID {
			continue
		}
		if outcome.Requirement != baseline.CapabilityRequired ||
			outcome.Status != baseline.CapabilitySatisfied ||
			outcome.Blocking {
			t.Fatalf("universal capability %q became a waiver: %+v", capabilityID, outcome)
		}
		return
	}
	t.Fatalf("universal capability %q is absent from the applied audit", capabilityID)
}

func assertBaselineReleaseRecommendation(t *testing.T, result baseline.Result, command string) {
	t.Helper()
	for _, recommendation := range result.Recommendations {
		if recommendation == command {
			return
		}
	}
	t.Fatalf("apply recommendations = %v, want %q", result.Recommendations, command)
}

func assertBaselineReleaseNoRecommendation(t *testing.T, result baseline.Result, command string) {
	t.Helper()
	for _, recommendation := range result.Recommendations {
		if recommendation == command {
			t.Fatalf("apply recommendations include unsupported command %q: %v", command, result.Recommendations)
		}
	}
}

func seedBaselineReleaseSourceManifest(t *testing.T, repo, profileID string) {
	t.Helper()
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := baseline.ResolveProfile("", profileID, catalog)
	if err != nil {
		t.Fatal(err)
	}
	manifest := baseline.SetupManifest{
		SchemaVersion: baseline.ManifestSchema,
		Version:       baseline.ManifestVersion,
		Generator: baseline.ManifestGenerator{
			Skill:    "setup-context-driven",
			Version:  baseline.ManifestVersion,
			Baseline: "baseline." + profile.ID + "-" + baseline.ManifestVersion,
		},
		Profile:          profile.ID,
		ProfileDigest:    profile.Digest,
		CatalogDigest:    "sha256:" + strings.Repeat("0", 64),
		Modules:          []string{},
		Decisions:        map[string]baseline.ManifestDecision{},
		ManagedArtifacts: []baseline.ManifestArtifact{},
		LocalSkills:      []string{},
		Verification:     []baseline.VerificationProjection{},
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeBaselinePlanTestFile(
		t,
		repo,
		"docs/agents/setup-context.json",
		string(append(data, '\n')),
	)
	commitBaselinePlanTestRepository(t, repo)
}

func runBaselineReleaseFormatter(t *testing.T, repo, profile string) {
	t.Helper()
	switch profile {
	case "go-cli-tui":
		runBaselineReleaseExternal(t, repo, "gofmt", "-w", "main.go")
	case "rust-cli":
		runBaselineReleaseExternal(t, repo, "rustfmt", "src/main.rs")
	case "standard-typescript-monorepo":
		formatterBin, formatterVersion := provisionBaselineReleaseFormatter(t, repo)
		versionOutput := runBaselineReleaseFixtureTool(
			t,
			repo,
			formatterBin,
			formatterVersion,
			"oxfmt",
			"--version",
		)
		if got := strings.TrimSpace(string(versionOutput)); got != formatterVersion {
			t.Fatalf("fixture formatter version = %q, want %q", got, formatterVersion)
		}
		runBaselineReleaseFixtureTool(
			t,
			repo,
			formatterBin,
			formatterVersion,
			"oxfmt",
			"--check",
			"AGENTS.md",
			"docs/agents",
		)
		if _, err := os.Stat(filepath.Join(repo, ".qa-formatter-ran")); err != nil {
			t.Fatalf("fixture formatter did not run: %v", err)
		}
	default:
		t.Fatalf("unknown maintained Profile %q", profile)
	}
}

func provisionBaselineReleaseFormatter(t *testing.T, repo string) (string, string) {
	t.Helper()
	packageData, err := os.ReadFile(filepath.Join(repo, "package.json"))
	if err != nil {
		t.Fatalf("read disposable repository package.json: %v", err)
	}
	var repositoryPackage struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal(packageData, &repositoryPackage); err != nil {
		t.Fatalf("decode disposable repository package.json: %v", err)
	}
	formatterVersion := repositoryPackage.Dependencies["oxfmt"]
	if strings.TrimSpace(formatterVersion) == "" {
		t.Fatal("disposable repository does not own an Oxfmt version")
	}

	binDir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create fixture-owned formatter directory: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(binDir, "oxfmt"),
		[]byte(baselineReleaseFormatterFixture),
		0o755,
	); err != nil {
		t.Fatalf("provision repository-local formatter fixture: %v", err)
	}
	return binDir, formatterVersion
}

const baselineReleaseFormatterFixture = `#!/bin/sh
set -eu

version="${ROUNDFIX_FIXTURE_FORMATTER_VERSION:?missing fixture formatter version}"
if [ "${1:-}" = "--version" ] && [ "$#" -eq 1 ]; then
	printf '%s\n' "$version"
	exit 0
fi
if [ "${BUN_CONFIG_REGISTRY:-}" != "http://127.0.0.1:1" ]; then
	printf '%s\n' "formatter fixture requires the unreachable registry" >&2
	exit 1
fi
cache="${BUN_INSTALL_CACHE_DIR:?missing fixture-owned cache}"
marker="$cache/.roundfix-formatter-fixture-used"
if [ -e "$marker" ]; then
	printf '%s\n' "formatter fixture cache was not fresh" >&2
	exit 1
fi
mkdir -p "$cache"
: > "$marker"
if [ "$#" -ne 3 ] || [ "$1" != "--check" ] || [ "$2" != "AGENTS.md" ] || [ "$3" != "docs/agents" ]; then
	printf '%s\n' "unexpected formatter invocation" >&2
	exit 1
fi
if [ ! -s AGENTS.md ] || [ ! -d docs/agents ]; then
	printf '%s\n' "managed formatter targets are missing" >&2
	exit 1
fi
: > .qa-formatter-ran
`

func newBaselineReleaseRepository(t *testing.T, profile string) string {
	t.Helper()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", strings.Join([]string{
		"verify:",
		"\t@touch .qa-verification-ran",
		"",
	}, "\n"))

	switch profile {
	case "go-cli-tui":
		writeBaselinePlanTestFile(t, repo, "main.go", "package main\nfunc main( ){ }\n")
	case "rust-cli":
		writeBaselinePlanTestFile(t, repo, "src/main.rs", "fn main( ) {println!(\"ok\");}\n")
	case "standard-typescript-monorepo":
		writeBaselinePlanTestFile(t, repo, "package.json", baselineReleaseTypeScriptPackageJSON(t))
		writeBaselinePlanTestFile(t, repo, "packages/frontend/package.json", `{"name":"frontend"}`)
		writeBaselinePlanTestFile(
			t,
			repo,
			"packages/backend/package.json",
			`{"name":"backend","dependencies":{"postgres":"latest","drizzle-orm":"latest"}}`,
		)
		writeBaselinePlanTestFile(
			t,
			repo,
			"packages/backend/src/infra/controllers/http/app.ts",
			"app.get('/health', handler)\napp.post('/orders', handler)\n",
		)
		writeBaselinePlanTestFile(
			t,
			repo,
			"DATABASE.md",
			"# Database\n\nPostgreSQL is the repository database contract.\n",
		)
	default:
		t.Fatalf("unknown maintained profile %q", profile)
	}
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func baselineReleaseTypeScriptPackageJSON(t *testing.T) string {
	t.Helper()
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load maintained Profile formatter: %v", err)
	}
	entry, ok := catalog.Profile("standard-typescript-monorepo")
	if !ok {
		t.Fatal("standard TypeScript Profile is missing")
	}
	var maintained struct {
		Formatter struct {
			Version string `json:"version"`
		} `json:"formatter"`
	}
	if err := json.Unmarshal(entry.Data, &maintained); err != nil {
		t.Fatalf("decode maintained Profile formatter: %v", err)
	}
	if strings.TrimSpace(maintained.Formatter.Version) == "" {
		t.Fatal("maintained Profile formatter version is empty")
	}
	dependencies := []string{
		`"@logtape/logtape":"latest"`,
		`"@tanstack/react-query":"latest"`,
		`"@tanstack/react-router":"latest"`,
		`"better-auth":"latest"`,
		`"drizzle-orm":"latest"`,
		`"hono":"latest"`,
		`"oxfmt":"` + maintained.Formatter.Version + `"`,
		`"oxlint":"latest"`,
		`"postgres":"latest"`,
		`"react":"latest"`,
		`"shadcn":"latest"`,
		`"tailwindcss":"latest"`,
		`"turbo":"latest"`,
		`"typescript":"latest"`,
		`"vite":"latest"`,
		`"vitest":"latest"`,
		`"zod":"latest"`,
	}
	return `{"name":"root","packageManager":"bun@1.3.0","scripts":{` +
		`"fmt":"touch .qa-profile-command-ran","lint":"touch .qa-profile-command-ran",` +
		`"test":"touch .qa-profile-command-ran","build":"touch .qa-profile-command-ran",` +
		`"verify":"touch .qa-profile-command-ran"},` +
		`"dependencies":{` + strings.Join(dependencies, ",") + `}}`
}

func runBaselineReleaseCLI(
	t *testing.T,
	binary string,
	args ...string,
) (int, []byte, []byte) {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err == nil {
		return 0, stdout.Bytes(), stderr.Bytes()
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("run Roundfix release binary: %v", err)
	}
	return exitErr.ExitCode(), stdout.Bytes(), stderr.Bytes()
}

func runBaselineReleaseExternal(t *testing.T, repo, name string, args ...string) {
	t.Helper()
	runBaselineReleaseExternalWithEnv(t, repo, nil, name, args...)
}

func runBaselineReleaseExternalWithEnv(
	t *testing.T,
	repo string,
	environmentOverrides []string,
	name string,
	args ...string,
) {
	t.Helper()
	command := exec.Command(name, args...)
	command.Dir = repo
	command.Env = baselineReleaseEnvironment(environmentOverrides...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("run external %s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
}

func baselineReleaseEnvironment(overrides ...string) []string {
	replaced := map[string]struct{}{
		"GIT_CONFIG_NOSYSTEM": {},
		"GIT_OPTIONAL_LOCKS":  {},
	}
	for _, override := range overrides {
		key, _, _ := strings.Cut(override, "=")
		replaced[key] = struct{}{}
	}
	environment := make([]string, 0, len(os.Environ())+len(overrides)+2)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if _, ok := replaced[key]; !ok {
			environment = append(environment, entry)
		}
	}
	environment = append(
		environment,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_OPTIONAL_LOCKS=0",
	)
	return append(environment, overrides...)
}

func runBaselineReleaseFixtureTool(
	t *testing.T,
	repo, fixtureBin, formatterVersion, name string,
	args ...string,
) []byte {
	t.Helper()
	executable := filepath.Join(fixtureBin, name)
	if _, err := os.Stat(executable); err != nil {
		t.Fatalf("find fixture-owned %s: %v", name, err)
	}
	command := exec.Command(executable, args...)
	command.Dir = repo
	isolatedBunHome := t.TempDir()
	command.Env = baselineReleaseEnvironment(
		"BUN_CONFIG_REGISTRY=http://127.0.0.1:1",
		"BUN_INSTALL="+isolatedBunHome,
		"BUN_INSTALL_CACHE_DIR="+filepath.Join(isolatedBunHome, "cache"),
		"BUN_INSTALL_GLOBAL_BIN_DIR="+filepath.Join(isolatedBunHome, "bin"),
		"BUN_INSTALL_GLOBAL_DIR="+filepath.Join(isolatedBunHome, "global"),
		"BUN_RUNTIME_TRANSPILER_CACHE_PATH="+filepath.Join(isolatedBunHome, "transpiler-cache"),
		"ROUNDFIX_FIXTURE_FORMATTER_VERSION="+formatterVersion,
		"PATH="+strings.Join([]string{
			fixtureBin,
			"/usr/bin",
			"/bin",
			"/usr/sbin",
			"/sbin",
		}, string(os.PathListSeparator)),
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run fixture tool %s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return output
}

func baselineReleaseManagedState(
	t *testing.T,
	repo string,
	plan baseline.PlanDocument,
) string {
	t.Helper()
	digest := sha256.New()
	for _, postimage := range plan.Postimages {
		content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(postimage.Path)))
		if err != nil {
			t.Fatalf("read managed path %s: %v", postimage.Path, err)
		}
		fmt.Fprintf(digest, "%s\x00%x\x00", postimage.Path, sha256.Sum256(content))
	}
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func baselineReleaseVisibleState(t *testing.T, root string) string {
	t.Helper()
	var rows []string
	err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == ".git" {
			return filepath.SkipDir
		}
		if relative == "." {
			return nil
		}
		info, err := os.Lstat(name)
		if err != nil {
			return err
		}
		row := fmt.Sprintf("%s:%s:%#o", relative, info.Mode().Type(), info.Mode().Perm())
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(name)
			if err != nil {
				return err
			}
			row += ":" + target
		case info.Mode().IsRegular():
			content, err := os.ReadFile(name)
			if err != nil {
				return err
			}
			row += fmt.Sprintf(":%x", sha256.Sum256(content))
		}
		rows = append(rows, row)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(rows)
	return strings.Join(rows, "\n")
}
