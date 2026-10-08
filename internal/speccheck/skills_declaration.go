package speccheck

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"roundfix/internal/skillcoverage"
	"roundfix/internal/spec"
)

const (
	// CodeSkillsUntasked identifies a covered source with no skill-writing Task.
	CodeSkillsUntasked = "SC-SKILLS-UNTASKED"
	// CodeSkillsMalformed identifies an invalid PRD Skills Declaration entry.
	CodeSkillsMalformed = "SC-SKILLS-MALFORMED"
)

func skillsHorizon(repoRoot, prdPath string) contractHorizon {
	commit, readable := prdAddingCommit(repoRoot, prdPath)
	if !readable || commit == "" {
		return contractHorizon{held: true}
	}
	output, err := adrHorizonGitOutput(repoRoot, "log", "--diff-filter=A", "--format=%H", "--", skillcoverage.MapPath)
	if err != nil {
		return contractHorizon{held: true}
	}
	commits := strings.Fields(string(output))
	missing := "a PRD committed at or after " + skillcoverage.MapPath
	if len(commits) == 0 {
		return contractHorizon{missing: missing}
	}
	_, err = adrHorizonGitOutput(repoRoot, "merge-base", "--is-ancestor", commits[len(commits)-1], commit)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return contractHorizon{missing: missing}
	}
	return contractHorizon{held: true}
}

func parseSkillsDeclaration(result *Result, content, display string, coverage skillcoverage.Map) map[string]bool {
	known := make(map[string]bool, len(coverage.Surfaces))
	for _, surface := range coverage.Surfaces {
		known[surface.ID] = true
	}
	unchanged := map[string]bool{}
	section := false
	for i, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "## ") {
			section = line == "## Skills"
			continue
		}
		if !section || strings.TrimSpace(line) == "" {
			continue
		}
		value, prefix := strings.CutPrefix(line, "- unchanged: ")
		id, reason, separator := strings.Cut(value, " — ")
		id = strings.TrimSpace(id)
		if !prefix || !separator || !known[id] || strings.TrimSpace(reason) == "" {
			result.Findings = append(result.Findings, Finding{
				Code: CodeSkillsMalformed, Severity: SeverityError,
				Summary: display + " has a malformed Skills Declaration",
				Where:   []Location{{Path: display, Line: i + 1}},
				Fix:     "Use - unchanged: <surface id> — <reason> with an id in " + skillcoverage.MapPath + " and a non-blank reason.",
			})
			continue
		}
		unchanged[id] = true
	}
	return unchanged
}

func detectSkillsDeclaration(result *Result, repoRoot, prdPath string, stage Stage, graph *spec.Graph) error {
	skipAll := func(missing string) {
		for _, code := range []string{CodeSkillsMalformed, CodeSkillsUntasked} {
			addSkip(result, code, missing)
		}
	}
	content, err := os.ReadFile(filepath.Join(repoRoot, skillcoverage.MapPath))
	if errors.Is(err, os.ErrNotExist) {
		skipAll(skillcoverage.MapPath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read Skill Coverage Map: %w", err)
	}
	horizon := skillsHorizon(repoRoot, prdPath)
	if !horizon.held {
		skipAll(horizon.missing)
		return nil
	}
	coverage, err := skillcoverage.ParseMap(content)
	if err != nil {
		return fmt.Errorf("read Skill Coverage Map: %w", err)
	}
	prd, err := os.ReadFile(prdPath)
	if err != nil {
		return fmt.Errorf("read Skills Declaration: %w", err)
	}
	unchanged := parseSkillsDeclaration(result, string(prd), artifactDisplayPath(repoRoot, prdPath), coverage)
	if stage == StagePRD || stage == StageTechSpec {
		return nil
	}
	if graph == nil {
		skipAll(artifactDisplayPath(repoRoot, filepath.Join(filepath.Dir(prdPath), "_tasks.md")))
		return nil
	}
	type binding struct {
		task spec.Task
		ref  spec.TaskContextRef
	}
	changed := map[string]binding{}
	written := map[string]bool{}
	for _, task := range graph.Tasks {
		if task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
			continue
		}
		for _, ref := range task.Context {
			if ref.Kind == spec.ContextKindInterface || ref.Kind == spec.ContextKindCreates {
				written[ref.Path] = true
			}
			if ref.Kind != spec.ContextKindInterface && ref.Kind != spec.ContextKindCreates && ref.Kind != spec.ContextKindDeletes {
				continue
			}
			for _, surface := range coverage.SurfacesForPath(ref.Path) {
				if _, seen := changed[surface.ID]; !seen {
					changed[surface.ID] = binding{task, ref}
				}
			}
		}
	}
	for _, surface := range coverage.Surfaces {
		binding, matched := changed[surface.ID]
		if !matched || surface.Uncovered != "" {
			continue
		}
		covered := unchanged[surface.ID] && written[skillcoverage.MapPath]
		for _, file := range surface.Skills {
			covered = covered || written[file]
		}
		if covered {
			continue
		}
		taskPath := filepath.Join(filepath.Dir(graph.Spec.Dir), filepath.FromSlash(binding.task.File))
		taskContent, err := os.ReadFile(taskPath)
		if err != nil {
			return fmt.Errorf("read skills source Task %q: %w", taskPath, err)
		}
		result.Findings = append(result.Findings, Finding{
			Code: CodeSkillsUntasked, Severity: SeverityError,
			Summary: fmt.Sprintf("Behavior Surface %s is changed by Task %s declaring %s without a Task for covering files: %s", surface.ID, binding.task.ID, binding.ref.Path, strings.Join(surface.Skills, ", ")),
			Where:   []Location{{Path: artifactDisplayPath(repoRoot, taskPath), Line: sectionLineContaining(taskContent, "Context", binding.ref.Path)}},
			Fix:     "Add a non-QA skills Task declaring a covering file under interface: or creates:, or a - unchanged: entry in the PRD and a non-QA Task declaring " + skillcoverage.MapPath + " to record the Coverage Review.",
		})
	}
	return nil
}
