package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/config"
	roundskills "roundfix/skills"
)

func loadAuthoringDeclaration(environment commandEnvironment, root string, stderr io.Writer) (*roundskills.AuthoringDeclaration, error) {
	environment.workDir, environment.workDirErr = root, nil
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return nil, err
	}
	return authoringDeclaration(loaded.Config.Skills), nil
}

func authoringDeclaration(settings config.Skills) *roundskills.AuthoringDeclaration {
	h := settings.Authoring.Claude
	if h == nil {
		return nil
	}
	declaration := &roundskills.AuthoringDeclaration{Harness: "claude", Entries: map[string]roundskills.AuthoringEntry{}}
	phases := []*config.AuthoringEntry{h.Idea, h.PRD, h.TechSpec, h.Tasks}
	for i, phase := range roundskills.AuthoringPhases() {
		if entry := phases[i]; entry != nil {
			declaration.Entries[phase.Key] = roundskills.AuthoringEntry{Model: entry.Model, Effort: entry.ReasoningEffort}
		}
	}
	if h.Research != nil {
		declaration.Entries["research"] = roundskills.AuthoringEntry{Model: h.Research.Model, Effort: h.Research.ReasoningEffort, Tools: append([]string{}, h.Research.Tools...)}
	}
	if len(declaration.Entries) == 0 {
		return nil
	}
	return declaration
}

func writeAuthoringInstall(result roundskills.AuthoringInstall, declaration *roundskills.AuthoringDeclaration, stdout, stderr io.Writer) {
	for _, skill := range result.Headers {
		for _, phase := range roundskills.AuthoringPhases() {
			if skill != phase.Skill || declaration == nil {
				continue
			}
			entry := declaration.Entries[phase.Key]
			var values []string
			if entry.Model != "" {
				values = append(values, "model "+entry.Model)
			}
			if entry.Effort != "" {
				values = append(values, "effort "+entry.Effort)
			}
			fmt.Fprintf(stdout, "Applied declared header: %s (%s)\n", skill, strings.Join(values, ", "))
		}
	}
	for _, path := range result.Subagents {
		fmt.Fprintf(stdout, "Generated subagent: %s\n", path)
	}
	for _, path := range result.Removed {
		fmt.Fprintf(stdout, "Removed generated subagent: %s\n", path)
	}
	writeUncoveredAuthoringNotices(result, stderr)
}

func writeUncoveredAuthoringNotices(result roundskills.AuthoringInstall, stderr io.Writer) {
	for _, skill := range result.Uncovered {
		for _, phase := range roundskills.AuthoringPhases() {
			if phase.Skill == skill {
				fmt.Fprintf(stderr, "roundfix: notice: .agents/skills/%s/SKILL.md held model or effort lines that no declaration covers; declare skills.authoring.claude.%s in .roundfixrc.yml to keep them\n", skill, phase.Key)
			}
		}
	}
}

func authoringInstallReport(result roundskills.AuthoringInstall) *roundskills.AuthoringInstall {
	if len(result.Headers)+len(result.Subagents)+len(result.Removed) == 0 {
		return nil
	}
	return &result
}

func baselineAuthoringLoader(environment commandEnvironment, stderr io.Writer) func(context.Context, string) (*roundskills.AuthoringDeclaration, error) {
	return func(ctx context.Context, root string) (*roundskills.AuthoringDeclaration, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return loadAuthoringDeclaration(environment, root, stderr)
	}
}
