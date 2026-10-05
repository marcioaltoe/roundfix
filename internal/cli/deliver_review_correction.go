package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/spec"
)

var _ delivery.ReviewCorrectionProver = (*commandDeliveryWorkflow)(nil)

func (workflow *commandDeliveryWorkflow) ProveReviewCorrection(ctx context.Context, workDir string, archivedSpecs []string, candidate, head string) (delivery.ReviewCorrection, error) {
	refuse := func(reason string) (delivery.ReviewCorrection, error) {
		return delivery.ReviewCorrection{Reason: reason}, nil
	}
	if _, err := workflow.git.RunGit(ctx, workDir, "merge-base", "--is-ancestor", candidate, head); err != nil {
		var exitErr interface{ ExitCode() int }
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return refuse("item head does not descend from parked candidate")
		}
		return delivery.ReviewCorrection{}, fmt.Errorf("check review correction ancestry: %w", err)
	}
	artifactDir, err := roundconfig.ValidateArtifactDirectory(workflow.loaded.Config.Defaults.ArtifactDir, workDir, workflow.loaded.HomeDir)
	if err != nil {
		return delivery.ReviewCorrection{}, err
	}
	record, err := readReviewRecord(filepath.Join(reviewCheckoutDir(artifactDir, workDir), reviewRecordFileName))
	if err != nil {
		return delivery.ReviewCorrection{}, fmt.Errorf("read review correction record: %w", err)
	}
	if record.Repository != workDir || record.HeadCommit != candidate || record.Outcome != reviewOutcomeFindings {
		return refuse("review record must name the item worktree and parked candidate with outcome findings")
	}
	ledger, err := readReviewFindingDispositions(filepath.Join(artifactDir, reviewDispositionLedgerFileName))
	if err != nil {
		return delivery.ReviewCorrection{}, err
	}
	dispositions := reviewDispositionsAtHead(ledger, workDir, candidate)
	for _, finding := range record.FindingItems {
		if reviewFindingDismissedByValidation(finding) {
			continue
		}
		matches := 0
		for _, d := range dispositions {
			if d.Finding != finding.ID || d.Text != finding.Text {
				continue
			}
			valid, err := reviewDispositionValidAtHead(ctx, workflow.git, workDir, candidate, head, d)
			if err != nil {
				return delivery.ReviewCorrection{}, err
			}
			if valid {
				matches++
			}
		}
		if matches != 1 {
			return refuse(fmt.Sprintf("finding %s requires exactly one valid disposition with evidence or a fixing commit between the candidate and item head", finding.ID))
		}
	}
	root, err := roundconfig.ResolveSpecsRoot(workflow.loaded, workDir)
	if err != nil {
		return delivery.ReviewCorrection{}, err
	}
	if root.External {
		return refuse("archived Specs directory is outside the item repository")
	}
	archiveRoot, err := filepath.Rel(workDir, spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot))
	if err != nil {
		return delivery.ReviewCorrection{}, fmt.Errorf("resolve review correction archive directory: %w", err)
	}
	var prefixes []string
	for _, slug := range archivedSpecs {
		if slug == "" || slug == "." || slug == ".." || strings.ContainsAny(slug, "/\\") {
			return refuse("invalid archived Spec slug")
		}
		prefixes = append(prefixes, filepath.ToSlash(filepath.Join(archiveRoot, slug))+"/")
	}
	paths, err := workflow.git.RunGit(ctx, workDir, "diff", "--name-only", "--no-renames", "-z", candidate, head, "--")
	if err != nil {
		return delivery.ReviewCorrection{}, fmt.Errorf("list review correction paths: %w", err)
	}
	for _, name := range nulPaths(paths) {
		allowed := false
		for _, prefix := range prefixes {
			allowed = allowed || strings.HasPrefix(name, prefix)
		}
		if !allowed {
			return refuse(fmt.Sprintf("changed path %s is outside the named archived Specs", name))
		}
	}
	return delivery.ReviewCorrection{Accepted: true}, nil
}
