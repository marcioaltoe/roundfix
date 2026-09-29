package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"

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
	graph, err := spec.Load(specsRoot.Path, specSlug)
	if err != nil {
		return delivery.Revalidation{}, fmt.Errorf("load item Spec %q: %w", specSlug, err)
	}
	result := delivery.Revalidation{Findings: findingCodes(findings)}
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

func sortedStringsContain(values []string, target string) bool {
	_, found := sort.Find(len(values), func(index int) int {
		return strings.Compare(target, values[index])
	})
	return found
}
