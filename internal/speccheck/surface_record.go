package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/skillcoverage"
	"roundfix/internal/spec"
)

// CodeSurfaceRecordUntasked identifies a guide Task without its own record step.
const CodeSurfaceRecordUntasked = "SC-SURFACE-RECORD-UNTASKED"

func detectSurfaceRecord(result *Result, repoRoot, specDir string, graph *spec.Graph) error {
	if !regularRepositoryFile(repoRoot, skillcoverage.RecordPath) {
		addSkip(result, CodeSurfaceRecordUntasked, skillcoverage.RecordPath)
		return nil
	}
	if graph == nil {
		addSkip(result, CodeSurfaceRecordUntasked, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md")))
		return nil
	}
	const rule = `following the write-tasks rule "Skills Tasks for changed Behavior Surfaces".`
	for _, task := range graph.Tasks {
		if task.Status == spec.StatusCompleted || task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
			continue
		}
		guide, record, step := "", false, false
		for _, ref := range task.Context {
			if ref.Kind != spec.ContextKindInterface && ref.Kind != spec.ContextKindCreates && ref.Kind != spec.ContextKindDeletes {
				continue
			}
			if guide == "" && strings.HasPrefix(ref.Path, "docs/user-guide/") && strings.HasSuffix(ref.Path, ".md") {
				guide = ref.Path
			}
			if ref.Path == skillcoverage.RecordPath && ref.Kind != spec.ContextKindDeletes {
				record = true
			}
		}
		if guide == "" {
			continue
		}
		for _, requirement := range task.Requirements {
			step = step || strings.Contains(requirement.Text, "-record-skill-coverage")
		}
		if record && step {
			continue
		}
		path := filepath.Join(filepath.Dir(graph.Spec.Dir), filepath.FromSlash(task.File))
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read guide Task %q for surface record: %w", path, err)
		}
		result.Findings = append(result.Findings, Finding{
			Code: CodeSurfaceRecordUntasked, Severity: SeverityGap,
			Summary: fmt.Sprintf("Task %s declares user guide %s without declaring %s and naming its -record-skill-coverage step", task.ID, guide, skillcoverage.RecordPath),
			Where:   []Location{{Path: artifactDisplayPath(repoRoot, path), Line: sectionLineContaining(content, "Context", guide)}},
			Fix:     "In this guide Task, declare " + skillcoverage.RecordPath + " under interface: or creates: and name -record-skill-coverage in a requirement, " + rule,
		})
	}
	path := filepath.Join(specDir, "_authorization.md")
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read surface record authorization %q: %w", path, err)
	}
	inLimits := false
	for index, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			inLimits = trimmed == "## Limits"
			continue
		}
		if !inLimits || !strings.Contains(line, skillcoverage.RecordPath) {
			continue
		}
		result.Findings = append(result.Findings, Finding{
			Code: CodeSurfaceRecordUntasked, Severity: SeverityGap,
			Summary: artifactDisplayPath(repoRoot, path) + " Limits name " + skillcoverage.RecordPath,
			Where:   []Location{{Path: artifactDisplayPath(repoRoot, path), Line: index + 1}},
			Fix:     "Remove " + skillcoverage.RecordPath + " from the authorization Limits so each guide Task can record its own fingerprint, " + rule,
		})
	}
	return nil
}
