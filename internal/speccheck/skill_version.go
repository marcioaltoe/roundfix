package speccheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/spec"
)

// CodeSkillVersionUntasked identifies an owned skill edit without its version work.
const CodeSkillVersionUntasked = "SC-SKILL-VERSION-UNTASKED"

const ownedSkillVersionRecordPath = "skills/testdata/owned-skill-versions.json"

func detectSkillVersion(result *Result, repoRoot, specDir string, graph *spec.Graph) error {
	recordPath := filepath.Join(repoRoot, ownedSkillVersionRecordPath)
	info, err := os.Stat(recordPath)
	if errors.Is(err, os.ErrNotExist) || (err == nil && !info.Mode().IsRegular()) {
		addSkip(result, CodeSkillVersionUntasked, ownedSkillVersionRecordPath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat owned skill version record %q: %w", recordPath, err)
	}
	if graph == nil {
		addSkip(result, CodeSkillVersionUntasked, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md")))
		return nil
	}
	content, err := os.ReadFile(recordPath)
	if err != nil {
		return fmt.Errorf("read owned skill version record %q: %w", recordPath, err)
	}
	var record struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(content, &record); err != nil {
		return fmt.Errorf("parse owned skill version record %q: %w", recordPath, err)
	}
	for _, task := range graph.Tasks {
		if task.Status == spec.StatusCompleted || task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
			continue
		}
		declarations := make(map[string]spec.ContextKind)
		edits := make(map[string][]spec.TaskContextRef)
		var names []string
		for _, ref := range task.Context {
			if ref.Kind == spec.ContextKindInstruction {
				continue
			}
			declarations[ref.Path] = ref.Kind
			relative, canonical := strings.CutPrefix(ref.Path, ".agents/skills/")
			name, file, nested := strings.Cut(relative, "/")
			if !canonical || !nested || file == "" {
				continue
			}
			if _, owned := record.Skills[name]; !owned {
				continue
			}
			if len(edits[name]) == 0 {
				names = append(names, name)
			}
			edits[name] = append(edits[name], ref)
		}
		step := false
		for _, requirement := range task.Requirements {
			step = step || strings.Contains(requirement.Text, "-record-skill-versions")
		}
		writable := func(path string) bool {
			kind := declarations[path]
			return kind == spec.ContextKindInterface || kind == spec.ContextKindCreates
		}
		for _, name := range names {
			var missing []string
			seen := make(map[string]bool)
			addMissing := func(path string) {
				if !seen[path] {
					missing = append(missing, "interface: "+path)
					seen[path] = true
				}
			}
			for _, ref := range edits[name] {
				mirror := strings.TrimPrefix(ref.Path, ".agents/")
				if ref.Kind == spec.ContextKindDeletes {
					if declarations[mirror] != spec.ContextKindDeletes {
						addMissing(mirror)
					}
				} else if !writable(mirror) {
					addMissing(mirror)
				}
			}
			for _, path := range []string{".agents/skills/" + name + "/SKILL.md", "skills/" + name + "/SKILL.md", ownedSkillVersionRecordPath} {
				if !writable(path) {
					addMissing(path)
				}
			}
			if !step {
				missing = append(missing, "a requirement naming -record-skill-versions")
			}
			if len(missing) == 0 {
				continue
			}
			path := filepath.Join(filepath.Dir(graph.Spec.Dir), filepath.FromSlash(task.File))
			taskContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read owned skill Task %q: %w", path, err)
			}
			result.Findings = append(result.Findings, Finding{
				Code: CodeSkillVersionUntasked, Severity: SeverityGap,
				Summary: fmt.Sprintf("Task %s edits owned skill %s without declaring its mirror, both SKILL.md files and %s and naming its -record-skill-versions step", task.ID, name, ownedSkillVersionRecordPath),
				Where:   []Location{{Path: artifactDisplayPath(repoRoot, path), Line: skillVersionContextLine(taskContent, edits[name][0])}},
				Fix:     "In this Task, add " + strings.Join(missing, "; ") + ", following the write-tasks rule \"Owned skill versions\".",
			})
		}
	}
	return nil
}

// Locate the mutation itself, even if an instruction cites the same path first.
func skillVersionContextLine(content []byte, ref spec.TaskContextRef) int {
	inContext := false
	for index, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			inContext = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Context"
			continue
		}
		value, matches := strings.CutPrefix(trimmed, "- "+string(ref.Kind)+":")
		if inContext && matches && strings.Trim(strings.TrimSpace(value), "`") == ref.Path {
			return index + 1
		}
	}
	return 1
}
