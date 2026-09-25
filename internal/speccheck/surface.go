package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/spec"
)

var guideRoots = []string{
	".agents/skills",
	"skills",
	"docs/user-guide",
	"docs/agents",
}

func detectUndocumentedCLISurfaces(result *Result, specsRoot, repoRoot string, graph *spec.Graph) error {
	tasksByID := make(map[string]spec.Task, len(graph.Tasks))
	for _, task := range graph.Tasks {
		tasksByID[task.ID] = task
	}

	for _, task := range graph.Tasks {
		if task.Status != spec.StatusPending || task.Type == spec.TaskTypeQA {
			continue
		}
		surface, present := firstCLISurface(task.Context)
		if !present || taskOrDependenciesNameGuide(task, tasksByID, make(map[string]bool)) {
			continue
		}

		taskPath := filepath.Join(filepath.Clean(specsRoot), filepath.FromSlash(task.File))
		content, err := os.ReadFile(taskPath)
		if err != nil {
			return fmt.Errorf("read Task %q for undocumented CLI surfaces: %w", task.File, err)
		}
		displayPath := artifactDisplayPath(repoRoot, taskPath)
		result.Findings = append(result.Findings, Finding{
			Code:     CodeCLIUndocumented,
			Severity: SeverityGap,
			Summary:  displayPath + " names CLI surface " + surface + " without a guide in the Task or its dependencies",
			Where: []Location{{
				Path: displayPath,
				Line: sectionLineContaining(content, "Context", surface),
			}},
			Fix: "Name the matching guide under `.agents/skills/`, `skills/`, `docs/user-guide/`, or `docs/agents/` in this Task or one of its dependencies.",
		})
	}
	return nil
}

func firstCLISurface(refs []spec.TaskContextRef) (string, bool) {
	for _, ref := range refs {
		if (ref.Kind == spec.ContextKindInterface || ref.Kind == spec.ContextKindCreates) && isCLISurface(ref.Path) {
			return ref.Path, true
		}
	}
	return "", false
}

func isCLISurface(candidate string) bool {
	if candidate == "internal/cli/cli_test.go" {
		return true
	}
	if !strings.HasSuffix(candidate, ".go") || strings.HasSuffix(candidate, "_test.go") {
		return false
	}
	return pathUnder(candidate, "internal/cli") || pathUnder(candidate, "cmd/roundfix")
}

func taskOrDependenciesNameGuide(task spec.Task, tasksByID map[string]spec.Task, visited map[string]bool) bool {
	if visited[task.ID] {
		return false
	}
	visited[task.ID] = true
	if taskNamesGuide(task) {
		return true
	}
	for _, dependencyID := range task.Needs {
		if dependency, present := tasksByID[dependencyID]; present && taskOrDependenciesNameGuide(dependency, tasksByID, visited) {
			return true
		}
	}
	return false
}

func taskNamesGuide(task spec.Task) bool {
	for _, ref := range task.Context {
		if ref.Kind != spec.ContextKindInterface && ref.Kind != spec.ContextKindCreates {
			continue
		}
		for _, root := range guideRoots {
			if pathUnder(ref.Path, root) {
				return true
			}
		}
	}
	return false
}

func pathUnder(candidate, root string) bool {
	return strings.HasPrefix(candidate, root+"/")
}
