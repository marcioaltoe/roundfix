package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
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
	record, err := readReviewRecordAtHead(artifactDir, workDir, candidate)
	if errors.Is(err, os.ErrNotExist) {
		return refuse("review record must name the item worktree and parked candidate with outcome findings")
	}
	if err != nil {
		return delivery.ReviewCorrection{}, fmt.Errorf("read review correction record: %w", err)
	}
	if !sameReviewCheckout(record.Repository, workDir) || record.HeadCommit != candidate || record.Outcome != reviewOutcomeFindings {
		return refuse("review record must name the item worktree and parked candidate with outcome findings")
	}
	ledger, err := readReviewFindingDispositions(filepath.Join(artifactDir, reviewDispositionLedgerFileName))
	if err != nil {
		return delivery.ReviewCorrection{}, err
	}
	dispositions := reviewDispositionsAtHead(ledger, workDir, candidate)
	allFixed := true
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
				allFixed = allFixed && d.Disposition == "fixed"
			}
		}
		if matches != 1 {
			return refuse(fmt.Sprintf("finding %s requires exactly one valid disposition with evidence or a fixing commit between the candidate and item head", finding.ID))
		}
	}
	if record.Lineage != nil && record.Lineage.Round == 2 && !record.Lineage.Correction && allFixed {
		return delivery.ReviewCorrection{Accepted: true}, nil
	}
	root, err := resolveHistoricalSpecsRoot(workflow.loaded, workDir)
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
	var declared map[string]bool
	for _, name := range nulPaths(paths) {
		allowed := false
		for _, prefix := range prefixes {
			allowed = allowed || strings.HasPrefix(name, prefix) || name == strings.TrimSuffix(prefix, "/")+".md"
		}
		if allowed {
			continue
		}
		// Archive paths have their own proof. Read Task declarations only when
		// an outside path needs the additional declared-path proof.
		if declared == nil {
			declared = make(map[string]bool)
			for _, slug := range archivedSpecs {
				paths, err := workflow.archivedSpecDeclaredPaths(ctx, workDir, candidate, filepath.ToSlash(archiveRoot), slug)
				if err != nil {
					return refuse(fmt.Sprintf("declared paths of archived Spec %s are unavailable: changed path %s is outside the named archived Specs: %v", slug, name, err))
				}
				for _, path := range paths {
					declared[path] = true
				}
			}
		}
		if !declared[name] {
			return refuse(fmt.Sprintf("changed path %s is outside the named archived Specs and the paths their Tasks declared", name))
		}
	}
	return delivery.ReviewCorrection{Accepted: true}, nil
}

// archivedSpecDeclaredPaths reads Task declarations from the parked candidate
// for a folder, or the Archive Record's source_revision for a record.
func (workflow *commandDeliveryWorkflow) archivedSpecDeclaredPaths(ctx context.Context, workDir, candidate, archiveRoot, slug string) ([]string, error) {
	archived, err := readArchivedSpecAt(ctx, workflow.git, workDir, candidate, archiveRoot, slug)
	if err != nil {
		return nil, fmt.Errorf("read archived Spec: %w", err)
	}
	revision, folder := candidate, archived.Path
	if archived.Form == spec.ArchivedRecord {
		revision, folder = archived.Record.SourceRevision, archived.Record.Source
		if _, err := spec.EnsureSourceRevision(ctx, workDir, workflow.loaded.Config.Watch.PushRemote, archived.Record, candidate, archived.Path); err != nil {
			return nil, fmt.Errorf("recover Task source revision: %w", err)
		}
	}
	listing, err := workflow.git.RunGit(ctx, workDir, "ls-tree", "-r", "--name-only", "-z", revision, "--", folder)
	if err != nil {
		return nil, fmt.Errorf("list archived Task files: %w", err)
	}
	var declared []string
	for _, file := range nulPaths(listing) {
		matches, err := filepath.Match("task_*.md", filepath.Base(file))
		if err != nil {
			return nil, fmt.Errorf("match archived Task file: %w", err)
		}
		if !matches || filepath.ToSlash(filepath.Dir(file)) != folder {
			continue
		}
		content, err := workflow.git.RunGit(ctx, workDir, "show", revision+":"+file)
		if err != nil {
			return nil, fmt.Errorf("read archived Task %s: %w", file, err)
		}
		paths, err := spec.DeclaredTaskPaths([]byte(content), file)
		if err != nil {
			return nil, err
		}
		declared = append(declared, paths...)
	}
	return declared, nil
}
