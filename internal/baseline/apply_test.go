// Suite: Exact portable Baseline apply
// Invariant: apply installs only one confirmed Plan Document or leaves the complete visible preimage unchanged.
// Boundary IN: strict plan approval, clone lineage, immutable backups, Baseline verification, transaction rollback, and idempotent reapply.
// Boundary OUT: public CLI flag parsing and stdout/stderr rendering.

package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestApplyCanonicalizesLegacyRepositoryRules(t *testing.T) {
	repo := newPlanRepository(t)
	const rules = "# Repository rules\n\nKeep the repository boundary explicit.\n"
	writeInspectionFile(t, repo, legacyRepositoryPath, rules)
	commitInspectionRepository(t, repo, "seed legacy repository rules")

	plan := buildPlanWithRepositoryExtension(t, repo)
	if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
		t.Fatalf("apply legacy repository rule migration: %v", err)
	}
	canonical, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(specificRepositoryPath)))
	if err != nil {
		t.Fatalf("read canonical repository rules: %v", err)
	}
	if !bytes.Equal(canonical, []byte(rules)) {
		t.Fatalf("canonical repository rules = %q, want exact legacy bytes %q", canonical, rules)
	}
	if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(legacyRepositoryPath))); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("legacy repository rule carrier remains after apply: %v", err)
	}
	agents, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read applied AGENTS.md: %v", err)
	}
	if !bytes.Contains(agents, []byte(specificRepositoryPath)) ||
		bytes.Contains(agents, []byte(legacyRepositoryPath)) {
		t.Fatalf("applied AGENTS.md does not point only to the canonical carrier:\n%s", agents)
	}
}

func TestRepositoryRuleBlockPreservesRepositoryEditAndEmptyReapply(t *testing.T) {
	repo := newPlanRepository(t)
	rule := []byte("Keep CLI diagnostics on stderr for this repository.\n")
	writeInspectionFile(t, repo, legacyRepositoryPath, string(rule))
	commitInspectionRepository(t, repo, "seed semantic repository rule")

	plan := buildSemanticCarrierPlan(t, repo, legacyRepositoryPath)
	if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
		t.Fatalf("apply semantic repository rule: %v", err)
	}
	beforeReapply := snapshotVisibleTree(t, repo)
	result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("empty semantic reapply: %v", err)
	}
	if result.State != "verified" || !strings.Contains(result.Message, "already applied") {
		t.Fatalf("empty semantic reapply result = %+v", result)
	}
	if afterReapply := snapshotVisibleTree(t, repo); !reflect.DeepEqual(afterReapply, beforeReapply) {
		t.Fatalf("empty semantic reapply changed managed files")
	}

	guidePath := "docs/agents/cli.md"
	guide, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(guidePath)))
	if err != nil {
		t.Fatal(err)
	}
	editedRule := []byte("Keep CLI diagnostics concise and actionable for this repository.\n")
	editedGuide := bytes.Replace(guide, rule, editedRule, 1)
	if bytes.Equal(editedGuide, guide) {
		t.Fatal("semantic repository-rule body was not found for edit")
	}
	writeTransactionFile(t, repo, guidePath, string(editedGuide), 0o644)
	commitInspectionRepository(t, repo, "edit repository-owned semantic rule")

	fresh := buildTestPlan(t, repo)
	freshGuide := planPostimage(t, fresh, guidePath)
	blocks, err := parseRepositoryRuleBlocks(guidePath, freshGuide.Content)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 || !bytes.Equal(blocks[0].Body, editedRule) {
		t.Fatalf("fresh plan semantic blocks = %+v, want edited bytes %q", blocks, editedRule)
	}
	for _, change := range fresh.FileChanges {
		if change.Path == guidePath {
			t.Fatalf("repository-owned semantic edit was treated as setup drift: %+v", change)
		}
	}
	var retained, inventoried bool
	for _, evidence := range fresh.Retention {
		if strings.HasPrefix(evidence.FromClause, "repository-rule.") &&
			evidence.Disposition == "repository-document" &&
			len(evidence.Targets) == 1 &&
			evidence.Targets[0] == guidePath {
			retained = true
		}
	}
	for _, entry := range fresh.ManagedEntries {
		if strings.HasPrefix(entry.ID, "repository-rule:") &&
			entry.Path == guidePath &&
			entry.Kind == "repository-owned" &&
			entry.ContentIdentity == planContentIdentity(editedRule) {
			inventoried = true
		}
	}
	if !retained || !inventoried {
		t.Fatalf("edited repository rule retention=%t inventory=%t", retained, inventoried)
	}
}

func TestRepositoryRuleBlockRollbackRestoresSemanticGuide(t *testing.T) {
	repo := newPlanRepository(t)
	const rule = "Keep repository CLI output stable.\n"
	writeInspectionFile(t, repo, legacyRepositoryRulesPath, rule)
	commitInspectionRepository(t, repo, "seed semantic rollback rule")
	plan := buildSemanticCarrierPlan(t, repo, legacyRepositoryRulesPath)
	before := snapshotVisibleTree(t, repo)

	transaction, err := beginTransaction(context.Background(), repo, plan, nil)
	if err != nil {
		t.Fatal(err)
	}
	tx := transaction.(*fileTransaction)
	tx.phaseHook = failTransactionOnce(
		transactionPhaseVerifying,
		"docs/agents/cli.md",
		errors.New("injected semantic guide verification failure"),
	)
	if _, err := tx.Apply(context.Background()); err == nil ||
		!strings.Contains(err.Error(), "injected semantic guide verification failure") {
		t.Fatalf("semantic rollback Apply() error = %v", err)
	}
	if err := tx.Close(); err != nil {
		t.Fatalf("close semantic rollback transaction: %v", err)
	}
	assertVisibleTree(t, repo, before)
}

func TestApplyExactDigest(t *testing.T) {
	repo := newPlanRepository(t)
	plan := buildTestPlan(t, repo)
	before := snapshotVisibleTree(t, repo)

	_, err := ApplyPlan(context.Background(), repo, plan, "sha256:"+strings.Repeat("0", 64))
	var applyErr *ApplyError
	if !errors.As(err, &applyErr) || applyErr.Kind != ApplyErrorApproval {
		t.Fatalf("ApplyPlan() error = %v, want approval refusal", err)
	}
	assertVisibleTree(t, repo, before)

	result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("ApplyPlan() exact digest: %v", err)
	}
	if result.State != "verified" || result.PlanDigest != plan.PlanDigest ||
		len(result.VerifiedPostimages) != len(plan.Postimages) {
		t.Fatalf("ApplyPlan() result = %+v", result)
	}
	for _, postimage := range plan.Postimages {
		assertAppliedPostimage(t, repo, postimage)
	}
}

func TestResultStatusMatrix(t *testing.T) {
	repo := newPlanRepository(t)
	plan := buildTestPlan(t, repo)

	first, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("first ApplyPlan(): %v", err)
	}
	assertResultStatusMatrix(t, first, ResultStatusMatrix{
		ApprovedPostimages:     EvidenceStatusVerified,
		SemanticRetention:      EvidenceStatusNotRun,
		ProfileAlignment:       EvidenceStatusVerified,
		RepositoryVerification: EvidenceStatusNotRun,
		Idempotence:            EvidenceStatusNotRun,
	})
	for _, axis := range []string{
		"Approved postimages: verified",
		"Semantic retention: not run",
		"Profile alignment: verified",
		"Repository Verification: not run",
		"Idempotence: not run",
	} {
		if !strings.Contains(first.Message, axis) {
			t.Errorf("human result lacks %q: %q", axis, first.Message)
		}
	}

	data, err := MarshalResult(first)
	if err != nil {
		t.Fatalf("MarshalResult(): %v", err)
	}
	var machine map[string]any
	if err := json.Unmarshal(data, &machine); err != nil {
		t.Fatalf("decode machine result: %v", err)
	}
	matrix, ok := machine["statusMatrix"].(map[string]any)
	wantMachineMatrix := map[string]any{
		"approvedPostimages":     "verified",
		"semanticRetention":      "not run",
		"profileAlignment":       "verified",
		"repositoryVerification": "not run",
		"idempotence":            "not run",
	}
	if !ok || !reflect.DeepEqual(matrix, wantMachineMatrix) {
		t.Fatalf("machine statusMatrix = %#v, want %#v", machine["statusMatrix"], wantMachineMatrix)
	}
	for _, field := range []string{
		"schemaVersion",
		"operation",
		"state",
		"message",
		"planDigest",
		"verifiedPostimages",
		"warnings",
		"recommendations",
	} {
		if _, exists := machine[field]; !exists {
			t.Errorf("machine result lost prior field %q", field)
		}
	}
	matrix["idempotence"] = "passed"
	invalid, err := json.Marshal(machine)
	if err != nil {
		t.Fatalf("encode invalid machine result: %v", err)
	}
	if _, err := ParseResult(invalid); err == nil || !strings.Contains(err.Error(), "invalid status") {
		t.Fatalf("ParseResult() invalid status error = %v", err)
	}

	second, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("second ApplyPlan(): %v", err)
	}
	assertResultStatusMatrix(t, second, ResultStatusMatrix{
		ApprovedPostimages:     EvidenceStatusVerified,
		SemanticRetention:      EvidenceStatusNotRun,
		ProfileAlignment:       EvidenceStatusVerified,
		RepositoryVerification: EvidenceStatusNotRun,
		Idempotence:            EvidenceStatusVerified,
	})
}

func TestCompletionLanguageRequiresRetention(t *testing.T) {
	t.Run("postimages and idempotence without retention", func(t *testing.T) {
		repo := newPlanRepository(t)
		plan := buildTestPlan(t, repo)
		if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
			t.Fatalf("first ApplyPlan(): %v", err)
		}
		result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
		if err != nil {
			t.Fatalf("second ApplyPlan(): %v", err)
		}
		if strings.Contains(strings.ToLower(result.Message), "complete") {
			t.Fatalf("result without retention reads as complete: %q", result.Message)
		}
	})

	t.Run("retention without idempotence", func(t *testing.T) {
		repo := newPlanRepository(t)
		plan := planWithVerifiedRetention(t, buildTestPlan(t, repo))
		result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
		if err != nil {
			t.Fatalf("first ApplyPlan(): %v", err)
		}
		if strings.Contains(strings.ToLower(result.Message), "complete") {
			t.Fatalf("result without idempotence reads as complete: %q", result.Message)
		}

		result, err = ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
		if err != nil {
			t.Fatalf("second ApplyPlan(): %v", err)
		}
		if !strings.Contains(strings.ToLower(result.Message), "complete") {
			t.Fatalf("result with retention and idempotence lacks completion language: %q", result.Message)
		}
	})
}

func assertResultStatusMatrix(t *testing.T, result Result, want ResultStatusMatrix) {
	t.Helper()
	if result.StatusMatrix == nil || *result.StatusMatrix != want {
		t.Fatalf("status matrix = %+v, want %+v", result.StatusMatrix, want)
	}
}

func planWithVerifiedRetention(t *testing.T, plan PlanDocument) PlanDocument {
	t.Helper()
	const clauseID = "clause.status-matrix-retained"
	delta := newClauseDelta()
	delta.Dispositions[clauseID] = ClauseRetained
	delta.Counts[ClauseRetained] = 1
	plan.Retention = append(plan.Retention, RetentionEvidence{
		FromClause:  clauseID,
		Enforcement: "mandatory",
		Disposition: string(ClauseRetained),
		Targets:     []string{clauseID},
		Reason:      "Stable clause identity and enforcement remain in the selected Baseline.",
	})
	plan.ClauseDelta = &delta
	var err error
	plan.PlanDigest, err = computePlanDigest(plan)
	if err != nil {
		t.Fatalf("compute retained Plan Digest: %v", err)
	}
	if err := ValidatePlanDocument(plan); err != nil {
		t.Fatalf("validate retained Plan: %v", err)
	}
	return plan
}

func TestApplyStalePreimage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		content  string
		wantPath string
	}{
		{
			name:     "consulted input",
			path:     "Makefile",
			content:  "verify:\n\t@echo changed\n",
			wantPath: "Makefile",
		},
		{
			name:     "mutation target",
			path:     "AGENTS.md",
			content:  "concurrent repository policy\n",
			wantPath: "AGENTS.md",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newPlanRepository(t)
			plan := buildTestPlan(t, repo)
			writeTransactionFile(t, repo, test.path, test.content, 0o644)
			before := snapshotVisibleTree(t, repo)

			_, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
			var applyErr *ApplyError
			if !errors.As(err, &applyErr) || applyErr.Kind != ApplyErrorStale ||
				!strings.Contains(err.Error(), test.wantPath) {
				t.Fatalf("ApplyPlan() error = %v, want stale %s refusal", err, test.wantPath)
			}
			assertVisibleTree(t, repo, before)
		})
	}
}

func TestApplyCrossClone(t *testing.T) {
	source := newPlanRepository(t)
	writeTransactionFile(t, source, "AGENTS.md", "repository policy\n", 0o644)
	commitInspectionRepository(t, source, "add root carrier")
	plan := buildTestPlan(t, source)

	cloneParent := t.TempDir()
	clone := filepath.Join(cloneParent, "matching")
	runApplyGit(t, cloneParent, "clone", "--quiet", source, clone)
	if _, err := ApplyPlan(context.Background(), clone, plan, plan.PlanDigest); err != nil {
		t.Fatalf("ApplyPlan() matching clone: %v", err)
	}

	unrelated := newUnrelatedPlanRepository(t)
	before := snapshotVisibleTree(t, unrelated)
	_, err := ApplyPlan(context.Background(), unrelated, plan, plan.PlanDigest)
	var applyErr *ApplyError
	if !errors.As(err, &applyErr) || applyErr.Kind != ApplyErrorStale ||
		!strings.Contains(err.Error(), "identity") {
		t.Fatalf("ApplyPlan() unrelated lineage error = %v", err)
	}
	assertVisibleTree(t, unrelated, before)
}

func TestImmutableRootBackup(t *testing.T) {
	const original = "repository policy\n"
	for _, existing := range []bool{false, true} {
		name := "creates exclusive backup"
		if existing {
			name = "accepts exact existing backup"
		}
		t.Run(name, func(t *testing.T) {
			repo := newPlanRepository(t)
			writeTransactionFile(t, repo, "AGENTS.md", original, 0o644)
			if existing {
				digest := strings.TrimPrefix(planContentIdentity([]byte(original)), "sha256:")
				writeTransactionFile(t, repo, "AGENTS."+digest+".md", original, 0o644)
			}
			commitInspectionRepository(t, repo, "add root carrier")
			plan := buildTestPlan(t, repo)

			var backup ManagedEntry
			for _, entry := range plan.ManagedEntries {
				if entry.Kind == "backup" {
					backup = entry
					break
				}
			}
			if backup.ID != "backup:AGENTS.md" ||
				!strings.Contains(backup.Path, strings.TrimPrefix(backup.ContentIdentity, "sha256:")) {
				t.Fatalf("immutable backup entry = %+v", backup)
			}
			if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
				t.Fatalf("ApplyPlan() with root backup: %v", err)
			}
			content, err := os.ReadFile(filepath.Join(repo, backup.Path))
			if err != nil {
				t.Fatalf("read immutable backup: %v", err)
			}
			if string(content) != original {
				t.Fatalf("immutable backup bytes = %q, want %q", content, original)
			}
		})
	}
}

func TestManagedRootFreshPlan(t *testing.T) {
	const managed = `<!-- setup-context-driven:begin id=root.core version=0.0.1 -->
managed root guidance
<!-- setup-context-driven:end id=root.core -->
`
	for _, test := range []struct {
		name       string
		root       string
		wantBackup bool
	}{
		{
			name: "setup-managed root needs no backup",
			root: managed,
		},
		{
			name:       "user-owned root keeps immutable backup",
			root:       "repository policy\n\n" + managed,
			wantBackup: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo := newPlanRepository(t)
			writeTransactionFile(t, repo, "AGENTS.md", test.root, 0o644)
			if err := os.Symlink("AGENTS.md", filepath.Join(repo, "CLAUDE.md")); err != nil {
				t.Fatalf("create root carrier alias: %v", err)
			}
			commitInspectionRepository(t, repo, "seed root carrier")

			initial := buildTestPlan(t, repo)
			if _, err := ApplyPlan(context.Background(), repo, initial, initial.PlanDigest); err != nil {
				t.Fatalf("ApplyPlan() initial Baseline: %v", err)
			}

			fresh := buildTestPlan(t, repo)
			var backups []ManagedEntry
			for _, entry := range fresh.ManagedEntries {
				if entry.Kind == "backup" {
					backups = append(backups, entry)
				}
			}
			if test.wantBackup {
				if len(backups) != 1 || backups[0].ID != "backup:AGENTS.md" {
					t.Fatalf("fresh Plan backups = %+v, want AGENTS.md backup", backups)
				}
			} else {
				if len(backups) != 0 {
					t.Fatalf("fresh Plan backups = %+v, want none", backups)
				}
				if len(fresh.FileChanges) != 0 {
					t.Fatalf("fresh Plan file changes = %+v, want none", fresh.FileChanges)
				}
			}
			if _, err := ApplyPlan(context.Background(), repo, fresh, fresh.PlanDigest); err != nil {
				t.Fatalf("ApplyPlan() fresh Baseline: %v", err)
			}
		})
	}
}

func TestApplyPostimageFailureRollsBack(t *testing.T) {
	repo := newPlanRepository(t)
	plan := buildTestPlan(t, repo)
	before := snapshotVisibleTree(t, repo)
	transaction, err := beginTransaction(context.Background(), repo, plan, nil)
	if err != nil {
		t.Fatal(err)
	}
	tx := transaction.(*fileTransaction)
	tx.phaseHook = failTransactionOnce(
		transactionPhaseVerifying,
		plan.Postimages[0].Path,
		errors.New("injected postimage verification failure"),
	)
	if _, err := tx.Apply(context.Background()); err == nil ||
		!strings.Contains(err.Error(), "injected postimage verification failure") {
		t.Fatalf("Apply() error = %v", err)
	}
	if err := tx.Close(); err != nil {
		t.Fatalf("close failed apply: %v", err)
	}
	assertVisibleTree(t, repo, before)
}

func TestProfileAdaptationApplyVerifiesProfileAndManifest(t *testing.T) {
	repo := newBackendProfileRepository(t, true)
	request, _ := backendProfileDraftPlanRequest(t, repo)
	outcome, err := BuildPlan(context.Background(), request)
	if err != nil || outcome.Plan == nil {
		t.Fatalf("BuildPlan() profile adaptation = plan=%v result=%+v error=%v",
			outcome.Plan != nil, outcome.Result, err)
	}
	plan := *outcome.Plan

	result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("ApplyPlan() profile adaptation error = %v", err)
	}
	const profilePath = ".roundfix/baseline/profiles/backend-only.json"
	profileBytes, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(profilePath)))
	if err != nil {
		t.Fatalf("read applied Profile: %v", err)
	}
	profilePostimage := planPostimage(t, plan, profilePath)
	if !bytes.Equal(profileBytes, profilePostimage.Content) {
		t.Fatal("applied Profile bytes differ from the approved postimage")
	}
	manifestBytes, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(manifestPath)))
	if err != nil {
		t.Fatalf("read applied Setup Manifest: %v", err)
	}
	var manifest SetupManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("decode applied Setup Manifest: %v", err)
	}
	if manifest.Profile != plan.Profile.ID || manifest.ProfileDigest != plan.Profile.Digest {
		t.Fatalf("applied Setup Manifest identity = %q %q, want %q %q",
			manifest.Profile, manifest.ProfileDigest, plan.Profile.ID, plan.Profile.Digest)
	}
	var verifiedProfile bool
	for _, postimage := range result.VerifiedPostimages {
		if postimage.Path == profilePath && bytes.Equal(postimage.Content, profileBytes) {
			verifiedProfile = true
		}
	}
	if !verifiedProfile {
		t.Fatal("apply result did not verify the planned Profile postimage")
	}
}

func TestProfileDraftRollbackRestoresMissingProfile(t *testing.T) {
	repo := newBackendProfileRepository(t, true)
	request, _ := backendProfileDraftPlanRequest(t, repo)
	outcome, err := BuildPlan(context.Background(), request)
	if err != nil || outcome.Plan == nil {
		t.Fatalf("BuildPlan() profile rollback = plan=%v result=%+v error=%v",
			outcome.Plan != nil, outcome.Result, err)
	}
	plan := *outcome.Plan
	before := snapshotVisibleTree(t, repo)
	transaction, err := beginTransaction(context.Background(), repo, plan, nil)
	if err != nil {
		t.Fatalf("begin profile rollback transaction: %v", err)
	}
	tx := transaction.(*fileTransaction)
	tx.phaseHook = failTransactionOnce(
		transactionPhaseVerifying,
		".roundfix/baseline/profiles/backend-only.json",
		errors.New("injected Profile verification failure"),
	)
	if _, err := tx.Apply(context.Background()); err == nil ||
		!strings.Contains(err.Error(), "injected Profile verification failure") {
		t.Fatalf("profile rollback Apply() error = %v", err)
	}
	if err := tx.Close(); err != nil {
		t.Fatalf("close profile rollback transaction: %v", err)
	}
	assertVisibleTree(t, repo, before)
}

func TestProfileDraftStaleTargetProducesNoMutation(t *testing.T) {
	repo := newBackendProfileRepository(t, true)
	request, _ := backendProfileDraftPlanRequest(t, repo)
	outcome, err := BuildPlan(context.Background(), request)
	if err != nil || outcome.Plan == nil {
		t.Fatalf("BuildPlan() stale profile = plan=%v result=%+v error=%v",
			outcome.Plan != nil, outcome.Result, err)
	}
	plan := *outcome.Plan
	const profilePath = ".roundfix/baseline/profiles/backend-only.json"
	writeTransactionFile(t, repo, profilePath, "{\"concurrent\":true}\n", 0o644)
	before := snapshotVisibleTree(t, repo)

	_, err = ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	var applyErr *ApplyError
	if !errors.As(err, &applyErr) ||
		applyErr.Kind != ApplyErrorStale ||
		!strings.Contains(err.Error(), profilePath) {
		t.Fatalf("ApplyPlan() stale profile error = %v, want stale refusal", err)
	}
	assertVisibleTree(t, repo, before)
}

func TestBaselineVerification(t *testing.T) {
	repo := newPlanRepository(t)
	plan := buildTestPlan(t, repo)
	marker := filepath.Join(repo, "verification-command-ran")
	for index := range plan.SetupManifest.Verification {
		plan.SetupManifest.Verification[index].Command = "touch " + marker
	}
	manifestBytes, err := marshalSetupManifest(plan.SetupManifest)
	if err != nil {
		t.Fatal(err)
	}
	for index := range plan.Postimages {
		if plan.Postimages[index].Path == manifestPath {
			plan.Postimages[index].Content = manifestBytes
			plan.Postimages[index].ContentIdentity = planContentIdentity(manifestBytes)
		}
	}
	for index := range plan.ManagedEntries {
		if plan.ManagedEntries[index].Path == manifestPath {
			plan.ManagedEntries[index].AfterIdentity = planContentIdentity(manifestBytes)
			plan.ManagedEntries[index].ContentIdentity = planContentIdentity(manifestBytes)
		}
	}
	plan.FileChanges, err = deriveFileChanges(plan.ManagedEntries, plan.Preimages, plan.Postimages)
	if err != nil {
		t.Fatal(err)
	}
	plan.PlanDigest, err = computePlanDigest(plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidatePlanDocument(plan); err != nil {
		t.Fatalf("validate command-recommendation plan: %v", err)
	}

	result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("ApplyPlan(): %v", err)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("repository Verification command ran: %v", err)
	}
	if len(result.Recommendations) == 0 || result.Recommendations[0] != "touch "+marker {
		t.Fatalf("recommendations = %v", result.Recommendations)
	}
}

func TestEmptyReapply(t *testing.T) {
	repo := newPlanRepository(t)
	const policy = "repository policy\n"
	writeTransactionFile(t, repo, "AGENTS.md", policy, 0o644)
	digest := strings.TrimPrefix(planContentIdentity([]byte(policy)), "sha256:")
	writeTransactionFile(t, repo, "AGENTS."+digest+".md", policy, 0o644)
	commitInspectionRepository(t, repo, "add root carrier and existing backup")
	plan := buildTestPlan(t, repo)
	if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
		t.Fatalf("first ApplyPlan(): %v", err)
	}
	before := snapshotVisibleTree(t, repo)
	result, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("empty reapply: %v", err)
	}
	if result.State != "verified" || !strings.Contains(result.Message, "already applied") {
		t.Fatalf("empty reapply result = %+v", result)
	}
	if after := snapshotVisibleTree(t, repo); !reflect.DeepEqual(after, before) {
		t.Fatalf("empty reapply changed managed files:\nbefore=%+v\nafter=%+v", before, after)
	}
}

func marshalSetupManifest(manifest SetupManifest) ([]byte, error) {
	data, err := jsonMarshalIndent(manifest)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func jsonMarshalIndent(value any) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}

func assertAppliedPostimage(t *testing.T, repo string, postimage Postimage) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(postimage.Path)))
	if err != nil {
		t.Fatalf("read applied postimage %s: %v", postimage.Path, err)
	}
	if !reflect.DeepEqual(content, postimage.Content) {
		t.Fatalf("applied postimage %s differs", postimage.Path)
	}
}

func runApplyGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func newUnrelatedPlanRepository(t *testing.T) string {
	t.Helper()
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeInspectionFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeInspectionFile(t, repo, "Makefile", "verify:\n\t@true\n")
	writeInspectionFile(t, repo, ".unrelated-lineage", "distinct root history\n")
	commitInspectionRepository(t, repo, "seed unrelated portable plan")
	return repo
}
