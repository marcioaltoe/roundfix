package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	roundconfig "roundfix/internal/config"
	roundskills "roundfix/skills"
)

const authoringInstallAction = "roundfix skills install --target project"

type updateCheckAuthoring struct {
	State    string   `json:"state"`
	Harness  string   `json:"harness"`
	Declared []string `json:"declared"`
	Problems []string `json:"problems"`
	Next     string   `json:"next,omitempty"`
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

func loadedAuthoringState(loaded roundconfig.Loaded) (roundskills.AuthoringState, bool) {
	if loaded.GitRoot == "" || loaded.RepositoryVersionErr != nil {
		return roundskills.AuthoringState{}, false
	}
	state := roundskills.InspectAuthoring(loaded.GitRoot, authoringDeclaration(loaded.Config.Skills))
	return state, len(state.Declared) > 0 || len(state.Problems) > 0
}

func printSkillsAuthoringWarning(ctx context.Context, environment commandEnvironment, stderr io.Writer) {
	workDir, err := environment.resolveWorkDir("resolve skills check working directory")
	if err != nil {
		return
	}
	root, err := environment.dependencies.resolveSkillsProjectRoot(ctx, workDir)
	if err != nil {
		return
	}
	declaration, err := loadAuthoringDeclaration(environment, root, io.Discard)
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: warning: Authoring Model Declaration not read: %v\n", err)
		return
	}
	state := roundskills.InspectAuthoring(root, declaration)
	if len(state.Problems) == 0 {
		return
	}
	fmt.Fprintf(stderr, "roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.%s):\n", state.Harness)
	for _, problem := range state.Problems {
		fmt.Fprintf(stderr, "  %s\n", problem.Text)
	}
	if next := authoringNext(state); next != "" {
		fmt.Fprintf(stderr, "  next: %s\n", next)
	}
}

func doctorAuthoringModelsResult(state roundskills.AuthoringState) CheckResult {
	result := CheckResult{Name: HealthCheckAuthoringModels, Status: CheckStatusOK, Detail: state.Harness + ": " + strings.Join(state.Declared, ", ") + " applied"}
	if len(state.Problems) > 0 {
		result.Status = CheckStatusWarn
		result.Detail = "the installed files do not match the Authoring Model Declaration: " + strings.Join(authoringProblemTexts(state), "; ")
		result.NextAction = authoringNext(state)
	}
	return result
}

func updateCheckAuthoringModels(loaded roundconfig.Loaded) *updateCheckAuthoring {
	state, report := loadedAuthoringState(loaded)
	if !report {
		return nil
	}
	result := &updateCheckAuthoring{State: "applied", Harness: state.Harness, Declared: append([]string{}, state.Declared...), Problems: authoringProblemTexts(state), Next: authoringNext(state)}
	if len(state.Problems) > 0 {
		result.State = "unapplied"
	}
	return result
}

func printUpdateCheckAuthoringModels(stdout io.Writer, report *updateCheckAuthoring) {
	if report == nil || report.State != "unapplied" {
		return
	}
	fmt.Fprintf(stdout, "Authoring models: the installed files do not match the declaration: %s", strings.Join(report.Problems, "; "))
	if report.Next != "" {
		fmt.Fprintf(stdout, "; apply with %s", report.Next)
	}
	fmt.Fprintln(stdout, ".")
}
