package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

// SettlementChecker reads the same facts as the QA gate before a Task settles.
type SettlementChecker interface {
	RefusingFindings(specsRoot, workDir, slug string) ([]speccheck.Finding, error)
	AuditCommit(context.Context, speccheck.MechanicalRequest, speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error)
}

type SpecCheckSettlementChecker struct{}

func (SpecCheckSettlementChecker) RefusingFindings(specsRoot, workDir, slug string) ([]speccheck.Finding, error) {
	result, err := specGatePrecondition(specsRoot, workDir, slug)
	return result.Findings, err
}

func (SpecCheckSettlementChecker) AuditCommit(ctx context.Context, request speccheck.MechanicalRequest, commit speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
	return speccheck.AuditProspectiveTaskCommit(ctx, request, commit)
}

func settlementMechanicalRequest(ctx context.Context, plan TaskPlan) (speccheck.MechanicalRequest, error) {
	prdPath := filepath.Join(plan.Spec.Dir, "_prd.md")
	var authorizationPath string
	var authorizationReference speccheck.MechanicalAuthorizationReference
	deliveryTargetRevision := plan.HeadSHA
	taskCommitsFromRunStart := false
	if strings.TrimSpace(plan.HeadSHA) == "" {
		var err error
		authorizationPath, _, err = speccheck.MechanicalAuthorization(plan.WorkDir, prdPath)
		if err != nil {
			return speccheck.MechanicalRequest{}, err
		}
	} else {
		base, resolvedBase, err := qaDeliveryBase(ctx, plan)
		if err != nil {
			return speccheck.MechanicalRequest{}, err
		}
		if resolvedBase {
			deliveryTargetRevision = base
		} else {
			taskCommitsFromRunStart = true
		}
		resolved, _, err := speccheck.ResolveMechanicalAuthorization(ctx, plan.WorkDir, prdPath, deliveryTargetRevision)
		if err != nil {
			return speccheck.MechanicalRequest{}, err
		}
		authorizationReference = resolved
		authorizationPath = resolved.Path
	}

	return speccheck.MechanicalRequest{RepoRoot: plan.WorkDir, ConsumingSpec: plan.Spec.Slug, AuthorizationPath: authorizationPath, AuthorizationReference: authorizationReference, DeliveryTargetRevision: deliveryTargetRevision, TaskCommitsFromRunStart: taskCommitsFromRunStart}, nil
}

func (engine *Engine) authorizationSettlementCheck(plan TaskPlan, task spec.Task, before []string) verificationCheck {
	return verificationCheck{Label: "settlement check: authorization", Run: func(ctx context.Context, _ string) (string, error) {
		prepared, err := engine.prepareTaskCommit(ctx, plan, task, before)
		if err != nil {
			return "", err
		}
		taskFile := artifactCommitPath(plan, filepath.Join(plan.SpecsRoot, task.File))
		governed := false
		for _, path := range prepared.stageable {
			if path != taskFile && speccheck.GovernedPath(path) {
				governed = true
				break
			}
		}
		if !governed {
			return "", nil
		}
		parent, err := runGitOutput(ctx, plan.WorkDir, "rev-parse", "HEAD")
		if err != nil {
			return "", fmt.Errorf("resolve prospective Task commit parent: %w", err)
		}
		request, err := settlementMechanicalRequest(ctx, plan)
		if err != nil {
			return "", err
		}
		findings, err := engine.deps.SettlementChecker.AuditCommit(ctx, request, speccheck.ProspectiveTaskCommit{TaskID: task.ID, TaskFile: taskFile, Parent: strings.TrimSpace(parent), Changed: prepared.stageable})
		if err != nil {
			return "", err
		}
		var diagnostics strings.Builder
		for _, finding := range findings {
			fmt.Fprintf(&diagnostics, "%s: %s\n%s:%d\nFix: %s\n", finding.Code, finding.Detail, finding.File, finding.Line, finding.Fix)
		}
		return strings.TrimSpace(diagnostics.String()), nil
	}}
}

func settlementDiagnosticPath(req verificationAttemptRequest, label string) string {
	name := strings.ReplaceAll(strings.TrimPrefix(label, "settlement check: "), " ", "-")
	path := VerificationOutputPath(req.ArtifactDir, req.RunID, req.BatchNumber, req.Attempt)
	return strings.TrimSuffix(path, filepath.Ext(path)) + "-settlement-" + name + ".log"
}

func writeSettlementDiagnostics(path, diagnostic string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create settlement diagnostics directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(diagnostic+"\n"), 0o644); err != nil {
		return fmt.Errorf("write settlement diagnostics: %w", err)
	}
	return nil
}
