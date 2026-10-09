package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/spec"
)

func (workflow *commandDeliveryWorkflow) Revalidate(
	ctx context.Context,
	workDir string,
	specSlug string,
	priorMerges []string,
) (delivery.Revalidation, error) {
	specsRoot, err := roundconfig.ResolveSpecsRoot(workflow.loaded, workDir)
	if err != nil {
		return delivery.Revalidation{}, fmt.Errorf("resolve item Specs Root: %w", err)
	}
	findings, err := strictSpecFindings(specsRoot.Path, workDir, specSlug)
	if err != nil {
		return delivery.Revalidation{}, fmt.Errorf("run strict Spec Consistency Check: %w", err)
	}
	result := delivery.Revalidation{Findings: findingCodes(findings)}
	if priorMerges != nil {
		result.OwnerWarning = workflow.deliveryOwnerWarning(ctx, workDir)
	}
	graph, err := spec.Load(specsRoot.Path, specSlug)
	if err != nil {
		return delivery.Revalidation{}, fmt.Errorf("load item Spec %q: %w", specSlug, err)
	}
	premises := productionPremises(graph)
	changed := make(map[string]struct{})
	for _, mergeCommit := range priorMerges {
		mergeCommit = strings.TrimSpace(mergeCommit)
		paths, err := workflow.git.RunGit(ctx, workDir, "diff", "--name-only", mergeCommit+"^", mergeCommit)
		if err != nil {
			return delivery.Revalidation{}, fmt.Errorf("read prior merge %q: %w", mergeCommit, err)
		}
		mergeChanged := false
		for _, path := range strings.Split(strings.TrimSpace(paths), "\n") {
			path = strings.TrimSpace(path)
			if path == "" || !sortedStringsContain(premises, path) {
				continue
			}
			changed[path] = struct{}{}
			mergeChanged = true
		}
		if mergeChanged {
			result.ChangedBy = append(result.ChangedBy, mergeCommit)
		}
	}
	for path := range changed {
		result.ChangedPremises = append(result.ChangedPremises, path)
	}
	sort.Strings(result.ChangedPremises)
	return result, nil
}

func (workflow *commandDeliveryWorkflow) deliveryOwnerWarning(ctx context.Context, workDir string) string {
	startingMain, err := workflow.git.RunGit(ctx, workDir, "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	startingMain = strings.TrimSpace(startingMain)
	binary := commandDependenciesForContext(ctx).auditor()
	evidence := spec.ResolveAuditorEvidence(ctx, workDir, startingMain, binary)
	if !evidence.SelfAudit || evidence.Ancestry != app.AncestryOlder {
		return ""
	}

	buildCommit := strings.TrimSuffix(strings.TrimSpace(binary.Commit), "-dirty")
	resolvedBuild, err := workflow.git.RunGit(ctx, workDir, "rev-parse", buildCommit+"^{commit}")
	if err != nil {
		return ""
	}
	changed, err := workflow.git.RunGit(
		ctx,
		workDir,
		"diff",
		"--name-only",
		buildCommit,
		startingMain,
		"--",
		"cmd",
		"internal",
		"go.mod",
		"go.sum",
	)
	if err != nil || strings.TrimSpace(changed) == "" {
		return ""
	}
	return fmt.Sprintf(
		"owner-older-than-main: owner build %s predates starting main %s",
		deliveryCommitName(resolvedBuild),
		deliveryCommitName(startingMain),
	)
}

func deliveryCommitName(commit string) string {
	commit = strings.TrimSpace(commit)
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}

func sortedStringsContain(values []string, target string) bool {
	_, found := sort.Find(len(values), func(index int) int {
		return strings.Compare(target, values[index])
	})
	return found
}
