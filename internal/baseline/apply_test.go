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
