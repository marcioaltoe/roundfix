package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	roundskills "roundfix/skills"
)

const authoringInstallAction = "roundfix skills install --target project"

type updateCheckAuthoring struct {
	State          string   `json:"state"`
	Harness        string   `json:"harness"`
	Declared       []string `json:"declared"`
	Problems       []string `json:"problems"`
	Next           string   `json:"next,omitempty"`
	Installed      []string `json:"installed,omitempty"`
	Recommendation string   `json:"recommendation,omitempty"`
}

func authoringNext(state roundskills.AuthoringState) string {
	for _, problem := range state.Problems {
		if problem.Installable {
			return authoringInstallAction
		}
	}
	return ""
}

func authoringProblemTexts(state roundskills.AuthoringState) []string {
	texts := make([]string, 0, len(state.Problems))
	for _, problem := range state.Problems {
		texts = append(texts, problem.Text)
	}
	return texts
}

// authoringReportState adds the installed, unpublished skills to the existing
// inspection without changing what skills check reports.
type authoringReportState struct {
	roundskills.AuthoringState
	Installed []string
}

func loadedAuthoringState(loaded roundconfig.Loaded) (authoringReportState, bool) {
	if loaded.GitRoot == "" || loaded.RepositoryVersionErr != nil {
		return authoringReportState{}, false
	}
	declaration := authoringDeclaration(loaded.Config.Skills)
	state := authoringReportState{AuthoringState: roundskills.InspectAuthoring(loaded.GitRoot, declaration)}
	if declaration == nil && len(state.Problems) == 0 {
		state.Installed = installedUnpublishedAuthoringSkills(loaded.GitRoot)
	}
	return state, len(state.Declared) > 0 || len(state.Problems) > 0 || len(state.Installed) > 0
}

func installedUnpublishedAuthoringSkills(root string) []string {
	var names []string
	for _, phase := range roundskills.AuthoringPhases() {
		info, err := os.Stat(filepath.Join(root, ".agents", "skills", phase.Skill, "SKILL.md"))
		if err == nil && info.Mode().IsRegular() && !roundskills.PublishesAuthoringSkill(root, phase.Skill) {
			names = append(names, phase.Skill)
		}
	}
	return names
}

func authoringRecommendationLines(block string) []string {
	return strings.Split(strings.TrimSuffix(block, "\n"), "\n")
}

func printSkillsAuthoringWarning(ctx context.Context, environment commandEnvironment, stderr io.Writer) []string {
	workDir, err := environment.resolveWorkDir("resolve skills check working directory")
	if err != nil {
		return nil
	}
	root, err := environment.dependencies.resolveSkillsProjectRoot(ctx, workDir)
	if err != nil {
		return nil
	}
	declaration, err := loadAuthoringDeclaration(environment, root, io.Discard)
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: warning: Authoring Model Declaration not read: %v\n", err)
		return nil
	}
	state := roundskills.InspectAuthoring(root, declaration)
	if len(state.Problems) == 0 {
		return declaredAuthoringSkills(root, declaration)
	}
	fmt.Fprintf(stderr, "roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.%s):\n", state.Harness)
	for _, problem := range state.Problems {
		fmt.Fprintf(stderr, "  %s\n", problem.Text)
	}
	if next := authoringNext(state); next != "" {
		fmt.Fprintf(stderr, "  next: %s\n", next)
	}
	return declaredAuthoringSkills(root, declaration)
}

func doctorAuthoringModelsResult(state authoringReportState) CheckResult {
	result := CheckResult{Name: HealthCheckAuthoringModels, Status: CheckStatusOK, Detail: state.Harness + ": " + strings.Join(state.Declared, ", ") + " applied"}
	if len(state.Problems) > 0 {
		result.Status = CheckStatusWarn
		result.Detail = "the installed files do not match the Authoring Model Declaration: " + strings.Join(authoringProblemTexts(state.AuthoringState), "; ")
		result.NextAction = authoringNext(state.AuthoringState)
	} else if len(state.Installed) > 0 {
		result.Status = CheckStatusWarn
		result.Detail = "the authoring skills are installed without an Authoring Model Declaration: " + strings.Join(state.Installed, ", ")
		result.NextAction = "add the Recommended Authoring Declaration below to .roundfixrc.yml, then run " + authoringInstallAction
		result.Lines = authoringRecommendationLines(roundconfig.RecommendedAuthoringYAML())
	}
	return result
}

func updateCheckAuthoringModels(loaded roundconfig.Loaded) *updateCheckAuthoring {
	state, report := loadedAuthoringState(loaded)
	if !report {
		return nil
	}
	result := &updateCheckAuthoring{State: "applied", Harness: state.Harness, Declared: append([]string{}, state.Declared...), Problems: authoringProblemTexts(state.AuthoringState), Next: authoringNext(state.AuthoringState)}
	if len(state.Problems) > 0 {
		result.State = "unapplied"
	} else if len(state.Installed) > 0 {
		result.State = "undeclared"
		result.Installed = state.Installed
		result.Recommendation = roundconfig.RecommendedAuthoringYAML()
	}
	return result
}

func printUpdateCheckAuthoringModels(stdout io.Writer, report *updateCheckAuthoring) {
	if report == nil {
		return
	}
	if report.State == "undeclared" {
		fmt.Fprintf(stdout, "Authoring models: the authoring skills are installed without a declaration: %s; add the Recommended Authoring Declaration below to .roundfixrc.yml, then apply with %s.\n", strings.Join(report.Installed, ", "), authoringInstallAction)
		for _, line := range authoringRecommendationLines(report.Recommendation) {
			fmt.Fprintf(stdout, "  %s\n", line)
		}
		return
	}
	if report.State != "unapplied" {
		return
	}
	fmt.Fprintf(stdout, "Authoring models: the installed files do not match the declaration: %s", strings.Join(report.Problems, "; "))
	if report.Next != "" {
		fmt.Fprintf(stdout, "; apply with %s", report.Next)
	}
	fmt.Fprintln(stdout, ".")
}

func declaredAuthoringSkills(root string, declaration *roundskills.AuthoringDeclaration) []string {
	if declaration == nil {
		return nil
	}
	var names []string
	for _, phase := range roundskills.AuthoringPhases() {
		if _, ok := declaration.Entries[phase.Key]; ok && !roundskills.PublishesAuthoringSkill(root, phase.Skill) {
			names = append(names, phase.Skill)
		}
	}
	return names
}

func loadedDeclaredAuthoringSkills(loaded roundconfig.Loaded) []string {
	if loaded.GitRoot == "" || loaded.RepositoryVersionErr != nil {
		return nil
	}
	return declaredAuthoringSkills(loaded.GitRoot, authoringDeclaration(loaded.Config.Skills))
}
