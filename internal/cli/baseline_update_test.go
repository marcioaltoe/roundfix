// Suite: manifest-driven Baseline update command
// Invariant: update either presents a digest-bound managed refresh without writes or applies that exact plan.
// Boundary IN: CLI parsing, manifest projection, managed-refresh planning, apply, injected skill refresh, structured output, and exit categories.
// Boundary OUT: external skill acquisition and interactive Baseline adoption.

package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	roundskills "roundfix/skills"
)

func TestBaselineUpdateAppliesManifestPlanAndReportsJSON(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repository)

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("baseline update exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.SchemaVersion != baselineUpdateResultSchema ||
		result.Operation != "update" ||
		result.State != "verified" ||
		result.ApprovedPlanDigest == "" ||
		result.ApprovedPlanDigest != result.PlanDigest {
		t.Fatalf("baseline update result = %+v", result)
	}
	if result.PriorCatalog.Digest == "" ||
		result.CurrentCatalog.Digest == "" ||
		result.PriorCatalog.Digest == result.CurrentCatalog.Digest {
		t.Fatalf("baseline update catalog identities = prior %+v current %+v", result.PriorCatalog, result.CurrentCatalog)
	}
	if len(result.FileChanges) < 2 {
		t.Fatalf("baseline update file changes = %+v, want stale guide and manifest", result.FileChanges)
	}
	if result.Retention == nil || result.Warnings == nil || result.AdoptedSuggestions == nil {
		t.Fatalf("baseline update JSON arrays must be present: %+v", result)
	}
	if before == baselinePlanTestTree(t, repository) {
		t.Fatal("baseline update did not rewrite stale managed artifacts")
	}
}

func TestBaselineUpdateSkillStageReportsInstalledAndRestoredSkills(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	var gotRequest baselineUpdateSkillsRequest
	stage := func(_ context.Context, request baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
		gotRequest = request
		return baselineUpdateSkillsResult{
			Status:         baselineUpdateSkillsVerified,
			InstalledCount: 1,
			Installed:      []string{"roundfix"},
			Restored:       []string{"context7"},
			Drifted:        []baselineUpdateSkillDrift{},
		}, nil
	}

	result, stdout, stderr, code := runBaselineUpdateTestCommandWithSkillsStage(
		t,
		context.Background(),
		stage,
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("skill-refreshing update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if gotRequest.Repository != repository || gotRequest.ProfileID == "" || gotRequest.SourceDir != "" {
		t.Fatalf("skills stage request = %+v", gotRequest)
	}
	if result.Skills.Status != baselineUpdateSkillsVerified ||
		result.Skills.InstalledCount != 1 ||
		len(result.Skills.Installed) != 1 || result.Skills.Installed[0] != "roundfix" ||
		len(result.Skills.Restored) != 1 || result.Skills.Restored[0] != "context7" {
		t.Fatalf("skills result = %+v", result.Skills)
	}
}

func TestBaselineUpdateSkillWarningKeepsApplyAxisVerified(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	stage := func(_ context.Context, _ baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
		return baselineUpdateSkillsResult{
			Status:         baselineUpdateSkillsWarning,
			InstalledCount: 1,
			Installed:      []string{"roundfix"},
			Restored:       []string{},
			Drifted: []baselineUpdateSkillDrift{{
				Skill:  "context7",
				Reason: "immutable upstream is unreachable",
			}},
		}, nil
	}

	result, stdout, stderr, code := runBaselineUpdateTestCommandWithSkillsStage(
		t,
		context.Background(),
		stage,
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("warning update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "verified" || result.StatusMatrix == nil ||
		result.StatusMatrix.ApprovedPostimages != baseline.EvidenceStatusVerified {
		t.Fatalf("apply axis = state %q matrix %+v", result.State, result.StatusMatrix)
	}
	if result.Skills.Status != baselineUpdateSkillsWarning || len(result.Skills.Drifted) != 1 ||
		result.Skills.Drifted[0].Skill != "context7" ||
		!strings.Contains(result.Skills.Drifted[0].Reason, "unreachable") {
		t.Fatalf("skills warning axis = %+v", result.Skills)
	}
}

func TestBaselineUpdateSkipsSkillStageAndPreservesSkillDirectory(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	sentinel := filepath.Join(repository, ".agents", "skills", "maintainer", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(sentinel), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sentinel, []byte("maintainer skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	skillsRoot := filepath.Join(repository, ".agents", "skills")
	before := baselinePlanTestTree(t, skillsRoot)
	stage := func(context.Context, baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
		t.Fatal("suppressed skills stage was called")
		return baselineUpdateSkillsResult{}, nil
	}

	result, stdout, stderr, code := runBaselineUpdateTestCommandWithSkillsStage(
		t,
		context.Background(),
		stage,
		"baseline", "update", "--repo", repository, "--yes", "--no-skills", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("skills-suppressed update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.Skills.Status != baselineUpdateSkillsSkipped {
		t.Fatalf("suppressed skills result = %+v", result.Skills)
	}
	if after := baselinePlanTestTree(t, skillsRoot); after != before {
		t.Fatalf("suppressed skills stage changed the project skill directory: got %s want %s", after, before)
	}
}

func TestBaselineUpdatePassesOfflineSourceToSkillStage(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	source := t.TempDir()
	var gotSource string
	stage := func(_ context.Context, request baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
		gotSource = request.SourceDir
		return successfulBaselineUpdateSkillsResult(), nil
	}

	_, stdout, stderr, code := runBaselineUpdateTestCommandWithSkillsStage(
		t,
		context.Background(),
		stage,
		"baseline", "update", "--repo", repository, "--yes",
		"--skills-source-dir", source, "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("offline-source update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	wantSource, err := filepath.Abs(source)
	if err != nil {
		t.Fatal(err)
	}
	if gotSource != wantSource {
		t.Fatalf("skills source = %q, want %q", gotSource, wantSource)
	}
}

func TestBaselineUpdateSkillsStageUsesPreviewThenConfirmation(t *testing.T) {
	const externalSkill = "context7"
	const planDigest = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	var restoreRequests []baseline.SkillsRestoreRequest
	dependencies := baselineUpdateSkillsDependencies{
		resolveProjectRoot: func(context.Context, string) (string, error) { return "/repo", nil },
		install: func(_ context.Context, request roundskills.InstallRequest) (roundskills.InstallResult, error) {
			if request.Target != "project" || request.ProjectDir != "/repo" {
				t.Fatalf("install request = %+v", request)
			}
			return roundskills.InstallResult{Targets: []roundskills.InstalledTarget{{Target: "project", Dir: "/repo/.agents/skills", Files: 3}}}, nil
		},
		ownedNames: func() []string { return []string{"roundfix"} },
		resolveExternal: func(string) ([]string, bool, error) {
			return []string{externalSkill}, true, nil
		},
		checkRepository: func(context.Context, string, []string) (roundskills.RepositoryReadiness, error) {
			return roundskills.RepositoryReadiness{MissingExternal: []string{externalSkill}}, nil
		},
		restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
			restoreRequests = append(restoreRequests, request)
			payload := baseline.SkillsRestorePayload{
				SchemaVersion: baseline.SkillsRestoreSchemaVersion,
				Profile:       request.ProfileID,
				Skills:        []baseline.RestoreSkill{{Skill: externalSkill}},
				PlanDigest:    pointerToString(planDigest),
			}
			if request.Confirmation == "" {
				restoreErr := &baseline.SkillsRestoreError{
					Category: baseline.SkillsRestoreAction,
					Finding: baseline.RestoreFinding{
						Code: "plan.confirmation.required",
					},
					Err: errors.New("restoration plan is not confirmed"),
				}
				payload.Finding = &restoreErr.Finding
				return payload, restoreErr
			}
			payload.OK = true
			payload.Applied = true
			return payload, nil
		},
	}

	result, err := runBaselineUpdateSkillsStageWith(
		context.Background(),
		baselineUpdateSkillsRequest{
			Repository: "/repo",
			ProfileID:  "go-cli-tui",
			SourceDir:  "/offline",
		},
		dependencies,
	)
	if err != nil {
		t.Fatalf("skills stage: %v", err)
	}
	if len(restoreRequests) != 2 || restoreRequests[0].Confirmation != "" ||
		restoreRequests[1].Confirmation != planDigest ||
		restoreRequests[0].SourceDir != "/offline" || restoreRequests[1].SourceDir != "/offline" {
		t.Fatalf("restore requests = %+v", restoreRequests)
	}
	if result.Status != baselineUpdateSkillsVerified || result.InstalledCount != 1 ||
		len(result.Restored) != 1 || result.Restored[0] != externalSkill {
		t.Fatalf("skills stage result = %+v", result)
	}
}

func TestBaselineUpdateSkillsStageDegradesUnreachableSourcePerSkill(t *testing.T) {
	const unreachableSkill = "context7"
	const restorableSkill = "testing-boss"
	const planDigest = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	dependencies := baselineUpdateSkillsDependencies{
		resolveProjectRoot: func(context.Context, string) (string, error) { return "/repo", nil },
		install: func(context.Context, roundskills.InstallRequest) (roundskills.InstallResult, error) {
			return roundskills.InstallResult{}, nil
		},
		ownedNames: func() []string { return []string{"roundfix"} },
		resolveExternal: func(string) ([]string, bool, error) {
			return []string{unreachableSkill, restorableSkill}, true, nil
		},
		checkRepository: func(context.Context, string, []string) (roundskills.RepositoryReadiness, error) {
			return roundskills.RepositoryReadiness{
				OutdatedExternal: []string{unreachableSkill, restorableSkill},
			}, nil
		},
		restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
			if request.Skills[0] == unreachableSkill {
				restoreErr := &baseline.SkillsRestoreError{
					Category: baseline.SkillsRestoreExecution,
					Finding: baseline.RestoreFinding{
						Code:    "source.commit-unavailable",
						Message: "immutable upstream is unreachable",
					},
					Err: errors.New("git fetch failed"),
				}
				return baseline.SkillsRestorePayload{Finding: &restoreErr.Finding}, restoreErr
			}
			payload := baseline.SkillsRestorePayload{
				Skills:     []baseline.RestoreSkill{{Skill: restorableSkill}},
				PlanDigest: pointerToString(planDigest),
			}
			if request.Confirmation == "" {
				restoreErr := &baseline.SkillsRestoreError{
					Category: baseline.SkillsRestoreAction,
					Finding:  baseline.RestoreFinding{Code: "plan.confirmation.required"},
					Err:      errors.New("restoration plan is not confirmed"),
				}
				payload.Finding = &restoreErr.Finding
				return payload, restoreErr
			}
			payload.OK = true
			payload.Applied = true
			return payload, nil
		},
	}

	result, err := runBaselineUpdateSkillsStageWith(
		context.Background(),
		baselineUpdateSkillsRequest{Repository: "/repo", ProfileID: "go-cli-tui"},
		dependencies,
	)
	if err != nil {
		t.Fatalf("unreachable source must be a warning: %v", err)
	}
	if result.Status != baselineUpdateSkillsWarning || len(result.Drifted) != 1 ||
		result.Drifted[0].Skill != unreachableSkill ||
		!strings.Contains(result.Drifted[0].Reason, "unreachable") ||
		len(result.Restored) != 1 || result.Restored[0] != restorableSkill {
		t.Fatalf("unreachable skills result = %+v", result)
	}
}

func TestBaselineUpdateIdempotenceReportsZeroFileChanges(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)

	first, firstOut, firstErr, firstCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if firstCode != exitOK || firstErr != "" || len(first.FileChanges) == 0 {
		t.Fatalf("first baseline update exit=%d result=%+v stdout=%s stderr=%s", firstCode, first, firstOut, firstErr)
	}
	afterFirst := baselinePlanTestTree(t, repository)

	second, secondOut, secondErr, secondCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if secondCode != exitOK || secondErr != "" {
		t.Fatalf("second baseline update exit=%d stdout=%s stderr=%s", secondCode, secondOut, secondErr)
	}
	if len(second.FileChanges) != 0 || second.State != "verified" {
		t.Fatalf("second baseline update result = %+v, want verified zero-change result", second)
	}
	if second.StatusMatrix == nil || second.StatusMatrix.Idempotence != baseline.EvidenceStatusVerified {
		t.Fatalf("second baseline update idempotence evidence = %+v", second.StatusMatrix)
	}
	if afterSecond := baselinePlanTestTree(t, repository); afterSecond != afterFirst {
		t.Fatal("second baseline update changed repository bytes")
	}
}

func TestBaselineUpdateNoManifestRequiresAdoptionWithoutWrites(t *testing.T) {
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, "README.md", "unadopted repository\n")
	commitBaselinePlanTestRepository(t, repository)
	before := baselinePlanTestTree(t, repository)

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("no-manifest update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "adoption" ||
		!strings.Contains(result.NextAction, "roundfix baseline") ||
		!strings.Contains(strings.ToLower(result.NextAction), "adoption") {
		t.Fatalf("no-manifest update result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("no-manifest update changed repository bytes")
	}
}

func TestBaselineUpdateNewDecisionRequiresActionWithoutWrites(t *testing.T) {
	repository := incompleteBaselineUpdateRepository(t, "secondbrain.enabled")
	before := baselinePlanTestTree(t, repository)

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("new-decision update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "decision" ||
		len(result.NewDecisions) != 1 || result.NewDecisions[0].ID != "secondbrain.enabled" ||
		!strings.Contains(result.Message, "secondbrain.enabled") {
		t.Fatalf("new-decision update result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("new-decision refusal changed repository bytes")
	}
}

func TestBaselineUpdateAdoptsEverySuggestedDecision(t *testing.T) {
	repository := incompleteBaselineUpdateRepository(t, "secondbrain.enabled")

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--adopt-suggested", "--yes", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("suggestion-adopting update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "verified" || len(result.AdoptedSuggestions) != 1 ||
		result.AdoptedSuggestions[0].ID != "secondbrain.enabled" ||
		result.AdoptedSuggestions[0].SuggestedValue != true {
		t.Fatalf("suggestion-adopting result = %+v", result)
	}
}

func TestBaselineUpdatePresentsPlanAndConfirmsPreviousDigest(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repository)

	planned, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=text",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("unconfirmed update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	for _, want := range []string{
		"Baseline update: plan ready",
		"Prior catalog:",
		"Current catalog:",
		"File changes:",
		"Retention evidence:",
		"Plan Digest:",
		"Next action:",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("unconfirmed update output missing %q:\n%s", want, stdout)
		}
	}
	if planned.PlanDigest == "" || planned.ApprovedPlanDigest != "" {
		t.Fatalf("unconfirmed update result = %+v", planned)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("unconfirmed update changed repository bytes")
	}

	confirmed, confirmOut, confirmErr, confirmCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--confirm-plan", planned.PlanDigest, "--format=json",
	)
	if confirmCode != exitOK || confirmErr != "" || confirmed.ApprovedPlanDigest != planned.PlanDigest {
		t.Fatalf("previous-digest confirmation exit=%d result=%+v stdout=%s stderr=%s", confirmCode, confirmed, confirmOut, confirmErr)
	}
}

func TestBaselineUpdateRejectsMutuallyExclusiveConfirmationForms(t *testing.T) {
	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--yes", "--confirm-plan", "sha256:reviewed", "--format=json",
	)
	if code != exitPreflight || !strings.Contains(stderr, "mutually exclusive") {
		t.Fatalf("mutually exclusive update exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if result.State != "failed" || result.Category != "invalid" {
		t.Fatalf("mutually exclusive update result = %+v", result)
	}
}

func TestBaselineUpdateRejectsDigestOtherThanCurrentPlanWithoutWrites(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repository)
	const staleDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--confirm-plan", staleDigest, "--format=json",
	)
	if code != exitUnverified || !strings.Contains(stderr, "does not match") {
		t.Fatalf("stale confirmation exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "approval" ||
		result.PlanDigest == "" || result.PlanDigest == staleDigest || result.ApprovedPlanDigest != "" {
		t.Fatalf("stale confirmation result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("stale confirmation changed repository bytes")
	}
}

func TestBaselineUpdateHelpNamesNonInteractiveContract(t *testing.T) {
	_, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--help",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("baseline update help exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	for _, want := range []string{
		"roundfix baseline update",
		"--yes | --confirm-plan <digest>",
		"--adopt-suggested",
		"--no-skills",
		"--skills-source-dir",
		baselineUpdateResultSchema,
		"without prompting or invoking a semantic analyzer",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("baseline update help missing %q:\n%s", want, stdout)
		}
	}
}

func TestBaselineUpdateExitCategoriesAndNoACPRuntimeDependency(t *testing.T) {
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, "README.md", "unadopted repository\n")
	commitBaselinePlanTestRepository(t, repository)

	var stderr bytes.Buffer
	code := RunContext(
		context.Background(),
		[]string{"baseline", "update", "--repo", repository, "--format=json"},
		failingWriter{err: errors.New("injected baseline update output failure")},
		&stderr,
	)
	if code != exitRunFailed || !strings.Contains(stderr.String(), "injected baseline update output failure") {
		t.Fatalf("output failure exit=%d stderr=%s", code, stderr.String())
	}

	adopted := newBaselineUpdateRepository(t)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, stdout, canceledErr, canceledCode := runBaselineUpdateTestCommand(
		t,
		canceled,
		"baseline", "update", "--repo", adopted, "--yes", "--format=json",
	)
	if canceledCode != exitSIGINT || !strings.Contains(canceledErr, "context canceled") {
		t.Fatalf("canceled update exit=%d stdout=%s stderr=%s", canceledCode, stdout, canceledErr)
	}
}

func runBaselineUpdateTestCommand(
	t *testing.T,
	ctx context.Context,
	args ...string,
) (baselineUpdateResult, string, string, int) {
	t.Helper()
	return runBaselineUpdateTestCommandWithSkillsStage(
		t,
		ctx,
		func(context.Context, baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
			return successfulBaselineUpdateSkillsResult(), nil
		},
		args...,
	)
}

func runBaselineUpdateTestCommandWithSkillsStage(
	t *testing.T,
	ctx context.Context,
	stage baselineUpdateSkillsStage,
	args ...string,
) (baselineUpdateResult, string, string, int) {
	t.Helper()
	if len(args) < 2 || args[0] != "baseline" || args[1] != "update" {
		t.Fatalf("baseline update test arguments = %v", args)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineUpdateCommandWithSkillsStage(
		ctx,
		args[2:],
		&stdout,
		&stderr,
		commandEnvironmentFromProcess(),
		stage,
	)
	var result baselineUpdateResult
	if index := strings.Index(stdout.String(), "{"); index >= 0 {
		if err := json.Unmarshal(stdout.Bytes()[index:], &result); err != nil {
			t.Fatalf("decode baseline update result: %v\n%s", err, stdout.String())
		}
	} else if digest := baselineUpdateTextValue(stdout.String(), "Plan Digest:"); digest != "" {
		result.PlanDigest = digest
	}
	return result, stdout.String(), stderr.String(), code
}

func successfulBaselineUpdateSkillsResult() baselineUpdateSkillsResult {
	return baselineUpdateSkillsResult{
		Status:         baselineUpdateSkillsVerified,
		InstalledCount: 1,
		Installed:      []string{"roundfix"},
		Restored:       []string{},
		Drifted:        []baselineUpdateSkillDrift{},
	}
}

func pointerToString(value string) *string {
	return &value
}

func baselineUpdateTextValue(output, label string) string {
	for _, line := range strings.Split(output, "\n") {
		if value, found := strings.CutPrefix(line, label); found {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func newBaselineUpdateRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineApplyTestRepository(t)
	plan, _ := baselineApplyTestPlan(t, repository)
	if _, err := baseline.ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatalf("adopt Baseline update fixture: %v", err)
	}
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func staleBaselineUpdateRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	manifest := readBaselineUpdateManifest(t, repository)
	manifest.CatalogDigest = "sha256:" + strings.Repeat("0", 64)

	artifactIndex := -1
	for index, artifact := range manifest.ManagedArtifacts {
		if artifact.Kind == "guide" {
			artifactIndex = index
			break
		}
	}
	if artifactIndex < 0 {
		t.Fatal("adopted fixture has no managed guide")
	}
	artifact := manifest.ManagedArtifacts[artifactIndex]
	const staleBody = "Stale managed guidance from the prior catalog.\n"
	path := filepath.Join(repository, filepath.FromSlash(artifact.Path))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read managed guide %q: %v", artifact.Path, err)
	}
	begin := "<!-- setup-context-driven:begin id=" + artifact.ID + " version=" + artifact.Version + " -->"
	end := "<!-- setup-context-driven:end id=" + artifact.ID + " -->"
	start := bytes.Index(content, []byte(begin))
	finish := bytes.Index(content, []byte(end))
	if start < 0 || finish < start {
		t.Fatalf("managed guide %q lacks artifact markers for %q", artifact.Path, artifact.ID)
	}
	finish += len(end)
	replacement := []byte(begin + "\n\n" + staleBody + "\n" + end)
	content = append(append(append([]byte(nil), content[:start]...), replacement...), content[finish:]...)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write stale managed guide %q: %v", artifact.Path, err)
	}
	sum := sha256.Sum256([]byte(staleBody))
	manifest.ManagedArtifacts[artifactIndex].Digest = hex.EncodeToString(sum[:])
	writeBaselineUpdateManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func incompleteBaselineUpdateRepository(t *testing.T, decisionID string) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	manifest := readBaselineUpdateManifest(t, repository)
	if _, exists := manifest.Decisions[decisionID]; !exists {
		t.Fatalf("adopted fixture manifest has no decision %q", decisionID)
	}
	delete(manifest.Decisions, decisionID)
	writeBaselineUpdateManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func readBaselineUpdateManifest(t *testing.T, repository string) baseline.SetupManifest {
	t.Helper()
	path := filepath.Join(repository, "docs", "agents", "setup-context.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Setup Manifest: %v", err)
	}
	var manifest baseline.SetupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("decode Setup Manifest: %v", err)
	}
	return manifest
}

func writeBaselineUpdateManifest(t *testing.T, repository string, manifest baseline.SetupManifest) {
	t.Helper()
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode Setup Manifest: %v", err)
	}
	data = append(data, '\n')
	path := filepath.Join(repository, "docs", "agents", "setup-context.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write Setup Manifest: %v", err)
	}
}
