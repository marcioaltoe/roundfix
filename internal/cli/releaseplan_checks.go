package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/baseline"
	"roundfix/internal/releaseplan"
)

const releasePlanCheckNextAction = "complete the skills and guides check in the release runbook before the release Pull Request"

type releasePlanCheck struct {
	Status     string `json:"status"`
	Detail     string `json:"detail"`
	NextAction string `json:"nextAction,omitempty"`
}

type releasePlanChecks struct {
	Skills   releasePlanCheck `json:"skills"`
	Baseline releasePlanCheck `json:"baseline"`
}

func newReleasePlanCheck(status, detail string) releasePlanCheck {
	result := releasePlanCheck{Status: status, Detail: detail}
	if status != "ok" && status != "current" {
		result.NextAction = releasePlanCheckNextAction
	}
	return result
}

func collectReleasePlanChecks(ctx context.Context, source releasePlanGitSource, environment commandEnvironment) releasePlanChecks {
	root, err := source.git(ctx, source.workDir, "rev-parse", "--show-toplevel")
	root = strings.TrimSpace(root)
	if err == nil && root == "" {
		err = errors.New("git returned an empty root")
	}
	if err != nil {
		check := newReleasePlanCheck("failed", "resolve repository root: "+err.Error())
		return releasePlanChecks{Skills: check, Baseline: check}
	}
	skills := repositorySkillsCheck(ctx, commandDependenciesForContext(ctx).doctor, root)
	return releasePlanChecks{
		Skills:   newReleasePlanCheck(string(skills.Status), skills.Detail),
		Baseline: collectReleasePlanBaselineCheck(ctx, root, environment),
	}
}

func collectReleasePlanBaselineCheck(ctx context.Context, root string, environment commandEnvironment) releasePlanCheck {
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		return newReleasePlanCheck("failed", err.Error())
	}
	input, err := baseline.ResolveManifestInput(root, catalog)
	if errors.Is(err, baseline.ErrNoManifest) {
		return newReleasePlanCheck("action_required", "the repository has no Setup Manifest")
	}
	if err != nil {
		return newReleasePlanCheck("failed", err.Error())
	}
	if len(input.NewDecisions) != 0 {
		return newReleasePlanCheck("action_required", "the current Baseline catalog requires new decisions: "+strings.Join(baselineUpdateDecisionIDs(input.NewDecisions), ", "))
	}
	directories, err := environment.executableDirectories("resolve Baseline executable search path")
	if err != nil {
		return newReleasePlanCheck("failed", err.Error())
	}
	outcome, err := baseline.BuildPlan(ctx, baseline.PlanRequest{
		Repository:            root,
		ProfileID:             input.ProfileID,
		Decisions:             input.Decisions,
		Preservation:          baseline.RootPreservationRequest{Mode: baseline.PreservationModeManagedRefresh},
		ExecutableDirectories: directories,
	})
	if err != nil {
		return newReleasePlanCheck("failed", err.Error())
	}
	if outcome.Plan == nil {
		return newReleasePlanCheck(outcome.Result.State, outcome.Result.Message)
	}
	files, moves := len(outcome.Plan.FileChanges), len(outcome.Plan.HistoryMoves)
	if files == 0 && moves == 0 {
		return newReleasePlanCheck("current", "the repository already matches the current Baseline catalog")
	}
	detail := fmt.Sprintf("%d file change(s)", files)
	if moves != 0 {
		detail += fmt.Sprintf(", %d history move(s)", moves)
	}
	return newReleasePlanCheck("plan_ready", detail)
}

func printReleasePlanChecksText(checks releasePlanChecks, stdout io.Writer) {
	for _, entry := range []struct {
		name  string
		check releasePlanCheck
	}{{"skills", checks.Skills}, {"baseline", checks.Baseline}} {
		fmt.Fprintf(stdout, "%s: %s: %s", entry.name, entry.check.Status, entry.check.Detail)
		if entry.check.NextAction != "" {
			fmt.Fprintf(stdout, "; next: %s", entry.check.NextAction)
		}
		fmt.Fprintln(stdout)
	}
}

func releasePlanJSONWithChecks(plan releaseplan.Plan, checks releasePlanChecks) releasePlanJSON {
	result := releasePlanJSONFromPlan(plan)
	result.Checks = checks
	return result
}
