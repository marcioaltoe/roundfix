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
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

func TestBaselineUpdateUnrecordedManagedRegionTextOutput(t *testing.T) {
	tests := []struct {
		name      string
		approval  []string
		wantState string
		wantExit  int
	}{
		{
			name:      "presented plan",
			wantState: "plan_ready",
			wantExit:  exitUnverified,
		},
		{
			name:      "applied plan",
			approval:  []string{"--yes"},
			wantState: "verified",
			wantExit:  exitOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := unrecordedBaselineUpdateRepository(t, []string{"line removed by refresh"})
			args := []string{"baseline", "update", "--repo", repository, "--format=text"}
			args = append(args, test.approval...)

			_, stdout, stderr, code := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				args...,
			)
			if code != test.wantExit || stderr != "" {
				t.Fatalf("%s update exit=%d stdout=%s stderr=%s", test.name, code, stdout, stderr)
			}
			for _, want := range []string{
				"Baseline update: " + strings.ReplaceAll(test.wantState, "_", " "),
				"Unrecorded managed regions: 1",
				"- Path: docs/agents/agent-instructions.md",
				"  Managed identity: guide.agent-instructions",
				"  Reason: digest-mismatch",
				"  Removed lines:\n    line removed by refresh",
			} {
				if !strings.Contains(stdout, want) {
					t.Fatalf("%s update output missing %q:\n%s", test.name, want, stdout)
				}
			}
		})
	}
}

func TestBaselineUpdateUnrecordedManagedRegionJSONOutput(t *testing.T) {
	tests := []struct {
		name      string
		approval  []string
		wantState string
		wantExit  int
	}{
		{
			name:      "presented plan",
			wantState: "plan_ready",
			wantExit:  exitUnverified,
		},
		{
			name:      "applied plan",
			approval:  []string{"--yes"},
			wantState: "verified",
			wantExit:  exitOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := unrecordedBaselineUpdateRepository(t, []string{"line removed by refresh"})
			args := []string{"baseline", "update", "--repo", repository, "--format=json"}
			args = append(args, test.approval...)

			result, stdout, stderr, code := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				args...,
			)
			if code != test.wantExit || stderr != "" {
				t.Fatalf("%s update exit=%d stdout=%s stderr=%s", test.name, code, stdout, stderr)
			}
			if result.State != test.wantState || result.SchemaVersion != baselineUpdateResultSchema {
				t.Fatalf("%s update result = %+v", test.name, result)
			}
			if len(result.UnrecordedManagedRegions) != 1 {
				t.Fatalf("%s unrecorded managed regions = %+v", test.name, result.UnrecordedManagedRegions)
			}
			region := result.UnrecordedManagedRegions[0]
			if region.Path != "docs/agents/agent-instructions.md" ||
				region.ManagedID != "guide.agent-instructions" ||
				region.Reason != baseline.UnrecordedManagedRegionReasonDigestMismatch ||
				!slices.Equal(region.RemovedLines, []string{"line removed by refresh"}) {
				t.Fatalf("%s unrecorded managed region = %+v", test.name, region)
			}
			if !strings.Contains(stdout, `"unrecordedManagedRegions"`) {
				t.Fatalf("%s JSON omitted unrecorded managed regions: %s", test.name, stdout)
			}
		})
	}
}

func TestBaselineUpdateUnrecordedManagedRegionWithoutRemovedLinesTextOutput(t *testing.T) {
	repository := unrecordedBaselineUpdateRepository(t, nil)

	_, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=text",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("no-removal update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "  Removed lines: no lines removed") {
		t.Fatalf("no-removal update lacks explicit statement:\n%s", stdout)
	}
}

func TestBaselineUpdateNoUnrecordedManagedRegionOmitsOutputs(t *testing.T) {
	repository := newBaselineUpdateRepository(t)

	_, textOutput, textError, textCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=text",
	)
	if textCode != exitOK || textError != "" {
		t.Fatalf("current text update exit=%d stdout=%s stderr=%s", textCode, textOutput, textError)
	}
	if strings.Contains(textOutput, "Unrecorded managed regions:") {
		t.Fatalf("current text update reported an unrecorded-region block:\n%s", textOutput)
	}

	result, jsonOutput, jsonError, jsonCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if jsonCode != exitOK || jsonError != "" {
		t.Fatalf("current JSON update exit=%d stdout=%s stderr=%s", jsonCode, jsonOutput, jsonError)
	}
	if result.UnrecordedManagedRegions != nil {
		t.Fatalf("current JSON update unrecorded regions = %#v, want nil", result.UnrecordedManagedRegions)
	}
	if strings.Contains(jsonOutput, `"unrecordedManagedRegions"`) {
		t.Fatalf("current JSON update emitted the optional field: %s", jsonOutput)
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

func TestBaselineUpdateUnresolvedProfileDiagnosis(t *testing.T) {
	tests := []struct {
		name          string
		profileID     string
		wantKind      baseline.UnresolvedProfileKind
		wantLocations []string
		wantAction    string
	}{
		{
			name:          "missing repository-owned Profile",
			profileID:     "repository-backend",
			wantKind:      baseline.UnresolvedProfileRepositoryMissing,
			wantLocations: []string{".roundfix/baseline/profiles/repository-backend.json"},
			wantAction:    "restore",
		},
		{
			name:          "unknown catalog identity",
			profileID:     "retired/profile",
			wantKind:      baseline.UnresolvedProfileCatalogUnknown,
			wantLocations: []string{},
			wantAction:    "adopt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			manifest := ReadBaselineSetupManifest(t, repository)
			manifest.Profile = test.profileID
			manifest.ProfileDigest = "sha256:" + strings.Repeat("0", 64)
			manifest.Generator.Baseline = "baseline." + test.profileID + "-" + baseline.ManifestVersion
			WriteBaselineSetupManifest(t, repository, manifest)

			result, stdout, stderr, code := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				"baseline", "update", "--repo", repository, "--format=json",
			)
			if code != exitPreflight || result.State != "failed" || result.Category != "manifest" {
				t.Fatalf("JSON update exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
			}
			if result.UnresolvedProfile == nil {
				t.Fatalf("JSON update unresolved Profile diagnosis is nil: %s", stdout)
			}
			diagnosis := *result.UnresolvedProfile
			if diagnosis.Identity != test.profileID || diagnosis.Kind != test.wantKind ||
				!slices.Equal(diagnosis.SearchedLocations, test.wantLocations) ||
				!strings.Contains(strings.ToLower(diagnosis.Action), test.wantAction) {
				t.Fatalf("JSON update unresolved Profile diagnosis = %+v", diagnosis)
			}
			if !strings.Contains(result.Message, diagnosis.Identity) ||
				!strings.Contains(result.Message, diagnosis.Action) ||
				strings.Contains(result.Message, "lstat") || strings.Contains(result.Message, "open ") ||
				strings.Contains(stderr, "lstat") || strings.Contains(stderr, "open ") {
				t.Fatalf("JSON update message=%q stderr=%q", result.Message, stderr)
			}

			_, text, textErr, textCode := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				"baseline", "update", "--repo", repository, "--format=text",
			)
			if textCode != exitPreflight ||
				!strings.Contains(text, diagnosis.Identity) ||
				!strings.Contains(text, diagnosis.Action) ||
				strings.Contains(text, "lstat") || strings.Contains(text, "open ") ||
				strings.Contains(textErr, "lstat") || strings.Contains(textErr, "open ") {
				t.Fatalf("text update exit=%d stdout=%q stderr=%q", textCode, text, textErr)
			}
			for _, location := range diagnosis.SearchedLocations {
				if !strings.Contains(text, location) {
					t.Fatalf("text update output %q lacks %q", text, location)
				}
			}
		})
	}
}

func TestBaselineUpdateFleetSweep(t *testing.T) {
	// Recorded fleet pattern: roundfix and fiscus had Setup Manifests whose
	// recorded digests predated their otherwise untouched Managed Regions.
	const manifestPredatesRegions = "roundfix and fiscus: manifest predates managed regions"
	// Recorded fleet pattern: conexus, tax-poc, and vortex were missing the same
	// fourteen structural clauses that the current catalog emits again.
	const structuralClausesMissing = "conexus, tax-poc, and vortex: structural clauses missing"
	// Recorded fleet pattern: fluxus named a repository-owned Baseline Profile
	// that was absent from the checkout.
	const unresolvedProfile = "fluxus: recorded Baseline Profile does not resolve"
	// Recorded fleet cohort: gss and oraculum were the only measured copies that
	// reached planning. This already-current copy exercises Task 09's required
	// zero-change endpoint for that non-blocking cohort.
	const alreadyCurrent = "gss and oraculum: current catalog has no proposed changes"

	type fleetCopy struct {
		name    string
		pattern string
		build   func(*testing.T) string
		apply   bool
	}
	corpus := []fleetCopy{
		{
			name:    "manifest-predates-managed-regions",
			pattern: manifestPredatesRegions,
			build: func(t *testing.T) string {
				return unrecordedBaselineUpdateRepository(t, nil)
			},
			apply: true,
		},
		{
			name:    "structural-clauses-missing",
			pattern: structuralClausesMissing,
			build:   newFleetStructuralClauseRepository,
			apply:   true,
		},
		{
			name:    "recorded-profile-does-not-resolve",
			pattern: unresolvedProfile,
			build:   newFleetUnresolvedProfileRepository,
		},
		{
			name:    "already-current",
			pattern: alreadyCurrent,
			build:   newBaselineUpdateRepository,
		},
	}

	corpusRoot := t.TempDir()
	for _, copy := range corpus {
		t.Run(copy.name, func(t *testing.T) {
			built := copy.build(t)
			repository := filepath.Join(corpusRoot, copy.name)
			if err := os.Rename(built, repository); err != nil {
				t.Fatalf("%s: move copy into fleet corpus: %v", copy.pattern, err)
			}
			relative, err := filepath.Rel(corpusRoot, repository)
			if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
				t.Fatalf("%s: repository %q escapes corpus root %q", copy.pattern, repository, corpusRoot)
			}

			result, stdout, stderr, code := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				"baseline", "update", "--repo", repository, "--format=json",
			)
			switch copy.pattern {
			case manifestPredatesRegions:
				if code != exitUnverified || stderr != "" || result.State != "plan_ready" ||
					result.Category != "approval" || result.PlanDigest == "" {
					t.Fatalf("%s: exit=%d result=%+v stdout=%s stderr=%s", copy.pattern, code, result, stdout, stderr)
				}
				if err := fleetManifestPredatesRegionsOracle(result); err != nil {
					t.Fatalf("%s: %v", copy.pattern, err)
				}
				withoutClassification := result
				withoutClassification.UnrecordedManagedRegions = nil
				if err := fleetManifestPredatesRegionsOracle(withoutClassification); err == nil {
					t.Fatalf("%s: removing the preservation classification did not break the sweep oracle", copy.pattern)
				}
			case structuralClausesMissing:
				if code != exitUnverified || stderr != "" || result.State != "plan_ready" ||
					result.Category != "approval" || result.PlanDigest == "" ||
					len(result.UnrecordedManagedRegions) != 0 {
					t.Fatalf("%s: exit=%d result=%+v stdout=%s stderr=%s", copy.pattern, code, result, stdout, stderr)
				}
				assertFleetStructuralClauseChanges(t, copy.pattern, result.FileChanges)
			case unresolvedProfile:
				if code != exitPreflight || result.State != "failed" || result.Category != "manifest" ||
					result.UnresolvedProfile == nil || result.NextAction == "" {
					t.Fatalf("%s: exit=%d result=%+v stdout=%s stderr=%s", copy.pattern, code, result, stdout, stderr)
				}
				diagnosis := result.UnresolvedProfile
				if diagnosis.Identity != "oraculum-backend" ||
					diagnosis.Kind != baseline.UnresolvedProfileRepositoryMissing ||
					!slices.Equal(diagnosis.SearchedLocations, []string{".roundfix/baseline/profiles/oraculum-backend.json"}) ||
					!strings.Contains(strings.ToLower(diagnosis.Action), "restore") ||
					strings.Contains(result.Message, "lstat") || strings.Contains(result.Message, "open ") {
					t.Fatalf("%s: diagnosis=%+v message=%q", copy.pattern, diagnosis, result.Message)
				}
			case alreadyCurrent:
				if code != exitOK || stderr != "" || result.State != "current" ||
					len(result.FileChanges) != 0 || result.PlanDigest == "" ||
					!strings.Contains(result.Message, "already matches the current Baseline catalog") {
					t.Fatalf("%s: exit=%d result=%+v stdout=%s stderr=%s", copy.pattern, code, result, stdout, stderr)
				}
			default:
				t.Fatalf("unhandled fleet pattern %q", copy.pattern)
			}

			if result.State != "plan_ready" && result.State != "current" && result.NextAction == "" {
				t.Fatalf("%s: state %q blocks before planning without a named human action", copy.pattern, result.State)
			}
			if !copy.apply {
				return
			}

			applied, applyOut, applyErr, applyCode := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				"baseline", "update", "--repo", repository, "--yes", "--format=json",
			)
			if applyCode != exitOK || applyErr != "" || applied.State != "verified" {
				t.Fatalf("%s: apply exit=%d result=%+v stdout=%s stderr=%s", copy.pattern, applyCode, applied, applyOut, applyErr)
			}
			if copy.pattern == structuralClausesMissing {
				assertFleetStructuralClausesPresent(t, repository)
			}

			current, currentOut, currentErr, currentCode := runBaselineUpdateTestCommand(
				t,
				context.Background(),
				"baseline", "update", "--repo", repository, "--format=json",
			)
			if currentCode != exitOK || currentErr != "" || current.State != "current" || len(current.FileChanges) != 0 {
				t.Fatalf("%s: next run exit=%d result=%+v stdout=%s stderr=%s", copy.pattern, currentCode, current, currentOut, currentErr)
			}
		})
	}
}

func fleetManifestPredatesRegionsOracle(result baselineUpdateResult) error {
	if len(result.UnrecordedManagedRegions) != 1 {
		return fmt.Errorf("unrecorded Managed Regions = %d, want 1", len(result.UnrecordedManagedRegions))
	}
	region := result.UnrecordedManagedRegions[0]
	if region.Path != "docs/agents/agent-instructions.md" ||
		region.ManagedID != "guide.agent-instructions" ||
		region.Reason != baseline.UnrecordedManagedRegionReasonDigestMismatch {
		return fmt.Errorf("unrecorded Managed Region = %+v, want agent-instructions digest mismatch", region)
	}
	return nil
}

type fleetStructuralClause struct {
	path string
	ids  []string
	line string
}

var fleetStructuralClauses = []fleetStructuralClause{
	{
		path: "docs/agents/backend.md",
		ids:  []string{"clause.backend.boundary-contracts", "rule.backend.boundary-contracts"},
		line: "- **mandatory**: Keep blocking, network, process, database, and daemon boundaries explicit about ownership, cancellation, timeouts, and error reporting. Test the lowest real boundary that proves the repository-authored contract; do not invent authentication, database, or transport policy.",
	},
	{
		path: "docs/agents/backend.md",
		ids:  []string{"clause.backend.http-independent-use-cases"},
		line: "- **mandatory**: Keep application use cases independent of HTTP request, response, router, and middleware types.",
	},
	{
		path: "docs/agents/backend.md",
		ids:  []string{"clause.backend.layered-architecture"},
		line: "- **mandatory**: Organize backend behavior through domain, application, and infrastructure layers. Dependencies point inward toward domain behavior.",
	},
	{
		path: "docs/agents/backend.md",
		ids:  []string{"clause.backend.persistence-owner"},
		line: "- **mandatory**: Keep persistence implementation in infrastructure and behind application-owned boundaries; schema and query definitions belong to the selected persistence capability.",
	},
	{
		path: "docs/agents/backend.md",
		ids:  []string{"clause.backend.prohibit-generic-layers"},
		line: "- **prohibited**: Do not introduce generic `modules` or `services` buckets as the normative backend architecture.",
	},
	{
		path: "docs/agents/backend.md",
		ids:  []string{"clause.backend.thin-http-handlers"},
		line: "- **mandatory**: Keep HTTP handlers thin: validate and translate transport input, invoke one application use case, and translate the result into the repository's HTTP Contract.",
	},
	{
		path: "docs/agents/domain.md",
		ids:  []string{"clause.domain.canonical-language"},
		line: "- **mandatory**: Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym.",
	},
	{
		path: "docs/agents/domain.md",
		ids:  []string{"clause.domain.layout-decision"},
		line: "- **mandatory**: Follow the repository's declared single-context or multi-context layout. Setup can require that decision but cannot infer bounded contexts from directory names.",
	},
	{
		path: "docs/agents/frontend.md",
		ids:  []string{"clause.frontend.organize-by-system"},
		line: "- **mandatory**: Organize frontend feature code by domain system. Each system exposes one public boundary while its internal components, hooks, queries, routes, and state import each other directly.",
	},
	{
		path: "docs/agents/frontend.md",
		ids:  []string{"clause.frontend.public-system-boundary"},
		line: "- **mandatory**: Import another system through that system's public boundary instead of reaching into its internal modules.",
	},
	{
		path: "docs/agents/issue-tracker.md",
		ids:  []string{"clause.spec.local-task-tracker-only"},
		line: "- **mandatory**: Use the local Spec folder, Task Graph, and Task files as the implementation issue tracker. Do not introduce external triage labels or external issue status as Task state.",
	},
	{
		path: "docs/agents/issue-tracker.md",
		ids:  []string{"clause.spec.status-only-in-task"},
		line: "- **mandatory**: Keep Task status only in the assigned Task file frontmatter. The Task Graph records topology and dependencies, not progress.",
	},
	{
		path: "docs/agents/monorepo.md",
		ids:  []string{"rule.monorepo.context-boundaries"},
		line: "- **mandatory**: Identify the owning package and bounded context before editing. Read its context and local instructions, keep changes inside the Task slice, and require an explicit owning contract plus boundary Verification for cross-package changes.",
	},
}

func newFleetStructuralClauseRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineReleaseRepository(t, "standard-typescript-monorepo")
	args := []string{"baseline", "plan", "--repo", repository, "--profile", "standard-typescript-monorepo"}
	for _, decision := range baselineReleaseDecisionArgs("standard-typescript-monorepo", "greenfield") {
		args = append(args, "--decision", decision)
	}
	args = append(args, "--format=json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunContext(context.Background(), args, &stdout, &stderr); code != exitOK {
		t.Fatalf("build structural-clause adoption plan exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	plan, err := baseline.ParsePlanDocument(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse structural-clause adoption plan: %v\n%s", err, stdout.String())
	}
	if _, err := baseline.ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatalf("apply structural-clause adoption plan: %v", err)
	}

	manifest := ReadBaselineSetupManifest(t, repository)
	manifest.CatalogDigest = "sha256:" + strings.Repeat("0", 64)
	removed := 0
	touched := make(map[string]struct{})
	for _, clause := range fleetStructuralClauses {
		path := filepath.Join(repository, filepath.FromSlash(clause.path))
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s structural-clause carrier: %v", clause.path, err)
		}
		needle := []byte(clause.line + "\n\n")
		if count := bytes.Count(content, needle); count != len(clause.ids) {
			t.Fatalf("structural clauses %v occur %d times in %s, want %d", clause.ids, count, clause.path, len(clause.ids))
		}
		content = bytes.ReplaceAll(content, needle, nil)
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatalf("write %s without structural clauses: %v", clause.path, err)
		}
		removed += len(clause.ids)
		touched[clause.path] = struct{}{}
	}
	if removed != 14 {
		t.Fatalf("removed structural clauses = %d, want 14", removed)
	}
	for path := range touched {
		updateFleetManifestArtifactDigest(t, repository, &manifest, path)
	}
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func updateFleetManifestArtifactDigest(
	t *testing.T,
	repository string,
	manifest *baseline.SetupManifest,
	carrierPath string,
) {
	t.Helper()
	for index := range manifest.ManagedArtifacts {
		artifact := &manifest.ManagedArtifacts[index]
		if artifact.Path != carrierPath {
			continue
		}
		content, err := os.ReadFile(filepath.Join(repository, filepath.FromSlash(carrierPath)))
		if err != nil {
			t.Fatalf("read managed carrier %s: %v", carrierPath, err)
		}
		begin := []byte("<!-- setup-context-driven:begin id=" + artifact.ID + " version=" + artifact.Version + " -->")
		end := []byte("<!-- setup-context-driven:end id=" + artifact.ID + " -->")
		start := bytes.Index(content, begin)
		finish := bytes.Index(content, end)
		if start < 0 || finish < start {
			t.Fatalf("managed carrier %s lacks markers for %s", carrierPath, artifact.ID)
		}
		body := content[start+len(begin) : finish]
		body = bytes.TrimPrefix(body, []byte("\n\n"))
		body = bytes.TrimSuffix(body, []byte("\n"))
		sum := sha256.Sum256(body)
		artifact.Digest = hex.EncodeToString(sum[:])
		return
	}
	t.Fatalf("Setup Manifest has no managed artifact for %s", carrierPath)
}

func newFleetUnresolvedProfileRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	manifest := ReadBaselineSetupManifest(t, repository)
	manifest.Profile = "oraculum-backend"
	manifest.ProfileDigest = "sha256:" + strings.Repeat("0", 64)
	manifest.Generator.Baseline = "baseline.oraculum-backend-" + baseline.ManifestVersion
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func assertFleetStructuralClauseChanges(t *testing.T, pattern string, changes []baseline.FileChange) {
	t.Helper()
	changed := make(map[string]struct{}, len(changes))
	for _, change := range changes {
		changed[change.Path] = struct{}{}
	}
	for _, clause := range fleetStructuralClauses {
		if _, ok := changed[clause.path]; !ok {
			t.Errorf("%s: plan does not restore %v in %s", pattern, clause.ids, clause.path)
		}
	}
}

func assertFleetStructuralClausesPresent(t *testing.T, repository string) {
	t.Helper()
	for _, clause := range fleetStructuralClauses {
		content, err := os.ReadFile(filepath.Join(repository, filepath.FromSlash(clause.path)))
		if err != nil {
			t.Fatalf("read restored structural-clause carrier %s: %v", clause.path, err)
		}
		if count := strings.Count(string(content), clause.line); count != len(clause.ids) {
			t.Errorf("restored structural clauses %v occur %d times in %s, want %d", clause.ids, count, clause.path, len(clause.ids))
		}
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
	manifest := ReadBaselineSetupManifest(t, repository)
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
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func unrecordedBaselineUpdateRepository(t *testing.T, removedLines []string) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	const (
		carrierPath = "docs/agents/agent-instructions.md"
		managedID   = "guide.agent-instructions"
	)

	if len(removedLines) == 0 {
		manifest := ReadBaselineSetupManifest(t, repository)
		found := false
		for index := range manifest.ManagedArtifacts {
			artifact := &manifest.ManagedArtifacts[index]
			if artifact.Path == carrierPath && artifact.ID == managedID {
				artifact.Digest = strings.Repeat("0", 64)
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Setup Manifest has no managed region %s:%s", carrierPath, managedID)
		}
		WriteBaselineSetupManifest(t, repository, manifest)
	} else {
		path := filepath.Join(repository, filepath.FromSlash(carrierPath))
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read managed guide %q: %v", carrierPath, err)
		}
		endMarker := []byte("<!-- setup-context-driven:end id=" + managedID + " -->")
		if !bytes.Contains(content, endMarker) {
			t.Fatalf("managed guide %q lacks end marker for %q", carrierPath, managedID)
		}
		inserted := []byte(strings.Join(removedLines, "\n") + "\n")
		content = bytes.Replace(content, endMarker, append(inserted, endMarker...), 1)
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatalf("write unrecorded managed guide %q: %v", carrierPath, err)
		}
	}
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func incompleteBaselineUpdateRepository(t *testing.T, decisionID string) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	manifest := ReadBaselineSetupManifest(t, repository)
	if _, exists := manifest.Decisions[decisionID]; !exists {
		t.Fatalf("adopted fixture manifest has no decision %q", decisionID)
	}
	delete(manifest.Decisions, decisionID)
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func ReadBaselineSetupManifest(t *testing.T, repository string) baseline.SetupManifest {
	t.Helper()
	path := filepath.Join(repository, filepath.FromSlash(baselineSetupManifestPath))
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

func WriteBaselineSetupManifest(t *testing.T, repository string, manifest baseline.SetupManifest) {
	t.Helper()
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode Setup Manifest: %v", err)
	}
	data = append(data, '\n')
	path := filepath.Join(repository, filepath.FromSlash(baselineSetupManifestPath))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write Setup Manifest: %v", err)
	}
}
