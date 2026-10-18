//go:build unix

// Suite: Baseline repository inspection of Unix special files.
// Invariant: unsafe aliases block inspection, including aliases to named pipes.
// Boundary IN: temporary Git repositories, instruction aliases and Unix special files.
// Boundary OUT: Baseline apply transactions and Windows filesystem behavior.
package baseline

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestInstructionAliasUnsafeTargetsBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(*testing.T, string)
	}{
		{
			name: "absolute external",
			setup: func(t *testing.T, repo string) {
				outside := filepath.Join(t.TempDir(), "outside.md")
				writeInspectionFile(t, filepath.Dir(outside), filepath.Base(outside), "outside\n")
				mustInspectionSymlink(t, outside, filepath.Join(repo, "AGENTS.md"))
			},
		},
		{
			name: "relative escape",
			setup: func(t *testing.T, repo string) {
				mustInspectionSymlink(t, "../outside.md", filepath.Join(repo, "AGENTS.md"))
			},
		},
		{
			name: "cycle",
			setup: func(t *testing.T, repo string) {
				mustInspectionSymlink(t, "CLAUDE.md", filepath.Join(repo, "AGENTS.md"))
				mustInspectionSymlink(t, "AGENTS.md", filepath.Join(repo, "CLAUDE.md"))
			},
		},
		{
			name: "directory",
			setup: func(t *testing.T, repo string) {
				if err := os.Mkdir(filepath.Join(repo, "policy"), 0o755); err != nil {
					t.Fatalf("create policy directory: %v", err)
				}
				mustInspectionSymlink(t, "policy", filepath.Join(repo, "AGENTS.md"))
			},
		},
		{
			name: "special file",
			setup: func(t *testing.T, repo string) {
				special := filepath.Join(repo, "policy.pipe")
				if err := syscall.Mkfifo(special, 0o600); err != nil {
					t.Fatalf("create named pipe: %v", err)
				}
				mustInspectionSymlink(t, "policy.pipe", filepath.Join(repo, "AGENTS.md"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInspectionRepository(t)
			tt.setup(t, repo)
			commitInspectionRepository(t, repo, "seed")

			inspection, err := InspectRepository(context.Background(), repo, nil)
			if err != nil {
				t.Fatalf("unsafe carrier must be reported in snapshot, got error: %v", err)
			}
			if !hasRepositoryFinding(inspection.Snapshot.Blocking, "baseline.inventory.unsafe-alias", "AGENTS.md") {
				t.Fatalf("unsafe alias was not apply-blocking: %+v", inspection.Snapshot.Blocking)
			}
		})
	}
}

func TestBaselineReleaseGate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "greenfield", run: TestGreenfieldPlanBacksUpWithoutImport},
		{name: "preservation", run: TestPreservationPlanAcceptsCompleteDecisionDocument},
		{name: "profile selection", run: TestProfileAlignmentResolvesExactlyOneProfile},
		{name: "stale plan", run: TestApplyStalePreimage},
		{name: "cross-clone apply", run: TestApplyCrossClone},
		{name: "unsafe carrier", run: TestInstructionAliasUnsafeTargetsBlock},
		{name: "rollback", run: TestTransactionRollback},
		{name: "semantic redistribution rollback", run: TestRepositoryRuleBlockRollbackRestoresSemanticGuide},
		{name: "Profile adaptation rollback", run: TestProfileDraftRollbackRestoresMissingProfile},
		{name: "recovery", run: TestTransactionRecovery},
		{name: "empty reapply", run: TestEmptyReapply},
		{name: "Profile divergence adaptation", run: TestProfileDivergenceResolution},
		{name: "universal capability remediation", run: TestUniversalCapabilityRemediation},
		{name: "rejected-plan revision", run: TestRejectedPlanRevision},
		{name: "renewed approval", run: TestRevisionRequiresNewApproval},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
