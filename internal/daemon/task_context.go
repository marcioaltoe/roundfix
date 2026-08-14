package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/spec"
)

const (
	specContextBundlePathLimit = 200
	implementTaskSkillPath     = ".agents/skills/implement-task/SKILL.md"
	rootAgentInstructionsPath  = "AGENTS.md"
)

func (engine *Engine) buildTaskContextBundle(ctx context.Context, plan TaskPlan, task spec.Task) (agent.SpecContextBundle, error) {
	var prior []string
	if strings.TrimSpace(plan.HeadSHA) != "" {
		resolved, err := engine.deps.PriorChanges.PriorChangedFiles(ctx, plan.WorkDir, plan.HeadSHA)
		if err != nil {
			return agent.SpecContextBundle{}, fmt.Errorf("resolve prior changed files for Task %s: %w", task.ID, err)
		}
		prior = resolved
	}
	return assembleTaskContextBundle(plan, task, prior), nil
}

type qaPromptContext struct {
	Bundle             agent.SpecContextBundle
	PreviousReportPath string
	PreviousReportHead string
}

func (engine *Engine) buildQAPromptContext(ctx context.Context, plan TaskPlan, task spec.Task) (qaPromptContext, error) {
	bundle, err := engine.buildTaskContextBundle(ctx, plan, task)
	if err != nil {
		return qaPromptContext{}, err
	}
	previousReport, err := spec.NewestQAReport(plan.Spec.Dir)
	if errors.Is(err, spec.ErrNoQAReport) {
		return qaPromptContext{Bundle: bundle}, nil
	}
	if err != nil {
		return qaPromptContext{}, fmt.Errorf("resolve previous QA Report for Spec %s: %w", plan.Spec.Slug, err)
	}
	// The Run starts from HeadSHA, so any report already visible in this fresh
	// Run Worktree was the report available at that head. The shared prior-file
	// resolver uses the same head as its changed-path base.
	previousReportHead := strings.TrimSpace(plan.HeadSHA)
	if previousReportHead == "" {
		return qaPromptContext{}, fmt.Errorf("resolve previous QA Report for Spec %s: Run start head is required", plan.Spec.Slug)
	}
	return qaPromptContext{
		Bundle:             bundle,
		PreviousReportPath: taskPromptPath(plan, previousReport),
		PreviousReportHead: previousReportHead,
	}, nil
}

func assembleTaskContextBundle(plan TaskPlan, task spec.Task, priorFiles []string) agent.SpecContextBundle {
	remaining := specContextBundlePathLimit
	seen := map[string]bool{}
	add := func(path string) (string, bool) {
		if remaining <= 0 {
			return "", false
		}
		clean := cleanBundlePath(path)
		if clean == "" || seen[clean] {
			return "", false
		}
		seen[clean] = true
		remaining--
		return clean, true
	}
	var bundle agent.SpecContextBundle
	if path, ok := add(existingSpecPath(plan, "_prd.md")); ok {
		bundle.PRD = path
	}
	if path, ok := add(existingSpecPath(plan, "_techspec.md")); ok {
		bundle.TechSpec = path
	}
	if path, ok := add(existingSpecPath(plan, "_tasks.md")); ok {
		bundle.TaskGraph = path
	}
	if path, ok := add(existingRepoPath(plan.WorkDir, rootAgentInstructionsPath)); ok {
		bundle.Instructions = append(bundle.Instructions, path)
	}
	if path, ok := add(existingRepoPath(plan.WorkDir, implementTaskSkillPath)); ok {
		bundle.Instructions = append(bundle.Instructions, path)
	}
	for _, ref := range task.Context {
		switch ref.Kind {
		case spec.ContextKindInstruction:
			if path, ok := add(ref.Path); ok {
				bundle.Instructions = append(bundle.Instructions, path)
			}
		case spec.ContextKindInterface, spec.ContextKindCreates:
			if path, ok := add(ref.Path); ok {
				bundle.Interfaces = append(bundle.Interfaces, path)
			}
		}
	}
	uniquePrior := make([]string, 0, len(priorFiles))
	priorSeen := map[string]bool{}
	for _, path := range priorFiles {
		clean := cleanBundlePath(path)
		if clean == "" || seen[clean] || priorSeen[clean] {
			continue
		}
		priorSeen[clean] = true
		uniquePrior = append(uniquePrior, clean)
	}
	sort.Strings(uniquePrior)
	for _, path := range uniquePrior {
		if remaining <= 0 {
			bundle.OmittedPriorFiles++
			continue
		}
		seen[path] = true
		remaining--
		bundle.PriorChangedFiles = append(bundle.PriorChangedFiles, path)
	}
	return bundle
}

func existingSpecPath(plan TaskPlan, name string) string {
	path := filepath.Join(plan.Spec.Dir, name)
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return taskPromptPath(plan, path)
}

func existingRepoPath(workDir string, rel string) string {
	clean := cleanBundlePath(rel)
	if clean == "" {
		return ""
	}
	if _, err := os.Stat(filepath.Join(workDir, filepath.FromSlash(clean))); err != nil {
		return ""
	}
	return clean
}

func cleanBundlePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	clean := filepath.Clean(path)
	if clean == "." {
		return ""
	}
	return filepath.ToSlash(clean)
}
