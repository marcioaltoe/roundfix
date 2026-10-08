package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/baseline"
	"roundfix/internal/releaseplan"
	"roundfix/internal/skillcoverage"
)

const releasePlanCheckNextAction = "complete the skills and guides check in the release runbook before the release Pull Request"

type releasePlanCheck struct {
	Status     string `json:"status"`
	Detail     string `json:"detail"`
	NextAction string `json:"nextAction,omitempty"`
}

type releasePlanChecks struct {
	Skills        releasePlanCheck              `json:"skills"`
	Baseline      releasePlanCheck              `json:"baseline"`
	SkillCoverage releasePlanSkillCoverageCheck `json:"skillCoverage"`
}

func newReleasePlanCheck(status, detail string) releasePlanCheck {
	result := releasePlanCheck{Status: status, Detail: detail}
	if status != "ok" && status != "current" {
		result.NextAction = releasePlanCheckNextAction
	}
	return result
}

func collectReleasePlanChecks(ctx context.Context, source releasePlanGitSource, environment commandEnvironment, plan releaseplan.Plan) releasePlanChecks {
	root, err := source.git(ctx, source.workDir, "rev-parse", "--show-toplevel")
	root = strings.TrimSpace(root)
	if err == nil && root == "" {
		err = errors.New("git returned an empty root")
	}
	if err != nil {
		check := newReleasePlanCheck("failed", "resolve repository root: "+err.Error())
		return releasePlanChecks{Skills: check, Baseline: check, SkillCoverage: failedReleasePlanSkillCoverageCheck(plan, err)}
	}
	skills := repositorySkillsCheck(ctx, commandDependenciesForContext(ctx).doctor, root)
	return releasePlanChecks{
		Skills:        newReleasePlanCheck(string(skills.Status), skills.Detail),
		Baseline:      collectReleasePlanBaselineCheck(ctx, root, environment),
		SkillCoverage: collectReleasePlanSkillCoverageCheck(ctx, source, plan),
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
	check := checks.SkillCoverage
	fmt.Fprintf(stdout, "skill-coverage: %s: %s", check.Status, check.Detail)
	if check.NextAction != "" {
		fmt.Fprintf(stdout, "; next: %s", check.NextAction)
	}
	fmt.Fprintln(stdout)
	for _, lagging := range check.Lagging {
		detail := "no covering skill"
		if len(lagging.Skills) != 0 {
			detail = "covering skills unchanged: " + strings.Join(lagging.Skills, ", ")
		}
		fmt.Fprintf(stdout, "- lagging: %s (%s; %s)\n", lagging.Surface, lagging.Change, detail)
	}

}

func releasePlanJSONWithChecks(plan releaseplan.Plan, checks releasePlanChecks) releasePlanJSON {
	result := releasePlanJSONFromPlan(plan)
	result.Checks = checks
	return result
}

const releasePlanSkillCoverageNextAction = "update a covering skill or record a Coverage Review in docs/references/skill-coverage.json for each lagging surface, then rerun roundfix release plan"
const releasePlanBlockedNextAction = "resolve the blocking skill-coverage check, then rerun roundfix release plan before any release mutation."

type releasePlanLaggingSurface struct {
	Surface string   `json:"surface"`
	Change  string   `json:"change"`
	Skills  []string `json:"skills"`
}
type releasePlanSkillCoverageCheck struct {
	releasePlanCheck
	Blocking bool                        `json:"blocking"`
	Lagging  []releasePlanLaggingSurface `json:"lagging,omitempty"`
}

func failedReleasePlanSkillCoverageCheck(plan releaseplan.Plan, err error) releasePlanSkillCoverageCheck {
	return releasePlanSkillCoverageCheck{releasePlanCheck: releasePlanCheck{Status: "failed", Detail: err.Error(), NextAction: releasePlanSkillCoverageNextAction}, Blocking: plan.State != releaseplan.StateNoRelease}
}
func collectReleasePlanSkillCoverageCheck(ctx context.Context, source releasePlanGitSource, plan releaseplan.Plan) releasePlanSkillCoverageCheck {
	target, err := source.skillCoverageSnapshot(ctx, plan.Target.CommitSHA)
	if err != nil {
		return failedReleasePlanSkillCoverageCheck(plan, err)
	}
	if target.Map == nil {
		return releasePlanSkillCoverageCheck{releasePlanCheck: releasePlanCheck{Status: "not_declared", Detail: "the repository has no Skill Coverage Map"}}
	}
	base, err := source.skillCoverageSnapshot(ctx, plan.Base.CommitSHA)
	if err != nil {
		return failedReleasePlanSkillCoverageCheck(plan, err)
	}
	if base.Map == nil {
		return releasePlanSkillCoverageCheck{releasePlanCheck: releasePlanCheck{Status: "introduced", Detail: fmt.Sprintf("the Skill Coverage Map is new since %s; no surface is compared", plan.Base.Tag)}}
	}
	paths, err := source.git(ctx, source.workDir, "diff", "--name-only", "--no-renames", plan.Base.CommitSHA, plan.Target.CommitSHA)
	if err != nil {
		return failedReleasePlanSkillCoverageCheck(plan, fmt.Errorf("read skill coverage changed paths: %w", err))
	}
	changed := map[string]bool{}
	for _, path := range splitNonEmptyLines(paths) {
		changed[path] = true
	}
	changes := skillcoverage.Compare(base, target, changed)
	result := releasePlanSkillCoverageCheck{}
	counts := map[string]int{}
	for _, change := range changes {
		counts[change.Outcome]++
		if change.Outcome == "lagging" {
			result.Lagging = append(result.Lagging, releasePlanLaggingSurface{Surface: change.ID, Change: change.Kind, Skills: change.Skills})
		}
	}
	result.Status = "current"
	result.Detail = fmt.Sprintf("%d changed surface(s) since %s: %d described, %d reviewed, %d uncovered", len(changes), plan.Base.Tag, counts["described"], counts["reviewed"], counts["uncovered"])
	if len(result.Lagging) != 0 {
		result.Status = "behind"
		result.Detail = fmt.Sprintf("%d lagging surface(s) since %s", len(result.Lagging), plan.Base.Tag)
		result.NextAction = releasePlanSkillCoverageNextAction
		result.Blocking = plan.State != releaseplan.StateNoRelease
	}
	return result
}
func (source releasePlanGitSource) skillCoverageSnapshot(ctx context.Context, revision string) (skillcoverage.Snapshot, error) {
	var result skillcoverage.Snapshot
	data, exists, err := source.skillCoverageFile(ctx, revision, skillcoverage.MapPath)
	if err != nil || !exists {
		return result, err
	}
	m, err := skillcoverage.ParseMap([]byte(data))
	if err != nil {
		return result, fmt.Errorf("parse %s at %s: %w", skillcoverage.MapPath, revision, err)
	}
	result.Map = &m
	data, exists, err = source.skillCoverageFile(ctx, revision, skillcoverage.RecordPath)
	if err != nil {
		return result, err
	}
	if !exists {
		return result, fmt.Errorf("missing %s at %s", skillcoverage.RecordPath, revision)
	}
	record, err := skillcoverage.ParseRecord([]byte(data))
	if err != nil {
		return result, fmt.Errorf("parse %s at %s: %w", skillcoverage.RecordPath, revision, err)
	}
	result.Record = &record
	return result, nil
}
func (source releasePlanGitSource) skillCoverageFile(ctx context.Context, revision, path string) (string, bool, error) {
	object := revision + ":" + path
	if _, err := source.git(ctx, source.workDir, "cat-file", "-e", object); err != nil {
		// Distinguish an absent path from an unreadable object or a Git failure.
		listing, listErr := source.git(ctx, source.workDir, "ls-tree", "-z", revision, "--", path)
		if listErr == nil && listing == "" {
			return "", false, nil
		}
		return "", false, fmt.Errorf("inspect %s: %w", object, err)
	}
	data, err := source.git(ctx, source.workDir, "show", object)
	if err != nil {
		return "", true, fmt.Errorf("read %s: %w", object, err)
	}
	return data, true, nil
}
