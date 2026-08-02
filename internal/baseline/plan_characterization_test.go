// Suite: Baseline plan characterization
// Invariant: every named repository shape keeps its complete plan and alignment outcome until an explicit regeneration records an intended change.
// Boundary IN: profile alignment, plan assembly, verified apply setup, diagnostics, divergences, warnings, and deterministic golden comparison.
// Boundary OUT: CLI rendering, interactive decisions, repository Verification execution, and production behavior changes.

package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var updateBaselinePlanCharacterization = flag.Bool(
	"update-baseline-plan-characterization",
	false,
	"regenerate the Baseline plan characterization corpus",
)

const baselinePlanCharacterizationSchema = "roundfix/baseline-plan-characterization/v1"

type baselinePlanCharacterization struct {
	SchemaVersion string                            `json:"schemaVersion"`
	Shape         string                            `json:"shape"`
	Fixture       baselinePlanCharacterizationShape `json:"fixture"`
	Alignment     ProfileAlignment                  `json:"alignment"`
	Outcome       PlanOutcome                       `json:"outcome"`
}

type baselinePlanCharacterizationShape struct {
	ProfileID             string `json:"profileId"`
	VerifiedApplyState    string `json:"verifiedApplyState,omitempty"`
	ExistingBaseline      string `json:"existingBaseline,omitempty"`
	ExistingProfileDigest string `json:"existingProfileDigest,omitempty"`
	ExistingCatalogDigest string `json:"existingCatalogDigest,omitempty"`
}

type baselinePlanCharacterizationCase struct {
	name  string
	build func(*testing.T) (string, PlanRequest, baselinePlanCharacterizationShape)
}

func TestBaselinePlanCharacterization(t *testing.T) {
	characterizationBin := baselinePlanCharacterizationBin(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("locate git for characterization repositories: %v", err)
	}
	if err := os.Symlink(gitPath, filepath.Join(characterizationBin, "git")); err != nil {
		t.Fatalf("expose git in controlled characterization PATH: %v", err)
	}
	t.Setenv("PATH", characterizationBin)
	t.Setenv("GIT_AUTHOR_DATE", "2000-01-01T00:00:00Z")
	t.Setenv("GIT_COMMITTER_DATE", "2000-01-01T00:00:00Z")

	cases := []baselinePlanCharacterizationCase{
		{name: "clean-adoption", build: buildCleanAdoptionCharacterization},
		{name: "idempotent-replan-after-verified-apply", build: buildIdempotentReplanCharacterization},
		{name: "unsatisfied-blocking-capabilities", build: buildBlockingCapabilitiesCharacterization},
		{name: "advisory-only-divergences", build: buildAdvisoryDivergencesCharacterization},
		{name: "same-baseline-changed-profile-and-catalog-digests", build: buildSameBaselineDigestDriftCharacterization},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			repository, request, fixture := test.build(t)
			catalog, err := LoadEmbeddedCatalog()
			if err != nil {
				t.Fatalf("%s: load embedded catalog: %v", test.name, err)
			}
			profile, err := ResolveProfile(repository, request.ProfileID, catalog)
			if err != nil {
				t.Fatalf("%s: resolve profile: %v", test.name, err)
			}
			decisions, missing, err := ResolveDecisionInput(profile, request.Decisions, catalog)
			if err != nil {
				t.Fatalf("%s: resolve decisions: %v", test.name, err)
			}
			if len(missing) != 0 {
				t.Fatalf("%s: unresolved characterization decisions: %v", test.name, missing)
			}
			alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
				ProfileID: request.ProfileID,
				Decisions: profileAlignmentDecisions(profile, decisions),
				Profile:   &profile,
			}, catalog)
			if err != nil {
				t.Fatalf("%s: resolve profile alignment: %v", test.name, err)
			}
			outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
			if err != nil {
				t.Fatalf("%s: build plan: %v", test.name, err)
			}

			record := baselinePlanCharacterization{
				SchemaVersion: baselinePlanCharacterizationSchema,
				Shape:         test.name,
				Fixture:       fixture,
				Alignment:     alignment,
				Outcome:       outcome,
			}
			actual := marshalBaselinePlanCharacterization(t, record, repository, characterizationBin)
			goldenPath := filepath.Join("testdata", "plan-characterization", test.name+".golden.json")
			if *updateBaselinePlanCharacterization {
				writeBaselinePlanCharacterizationGolden(t, goldenPath, actual)
				return
			}
			expected, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf(
					"shape %q: read golden %s: %v; regenerate deliberately with -update-baseline-plan-characterization",
					test.name,
					goldenPath,
					err,
				)
			}
			if err := compareBaselinePlanCharacterization(test.name, expected, actual); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBaselinePlanCharacterizationDiffNamesShapeAndField(t *testing.T) {
	want := []byte("{\n  \"outcome\": {\n    \"Result\": {\n      \"state\": \"ready\"\n    }\n  }\n}\n")
	got := []byte("{\n  \"outcome\": {\n    \"Result\": {\n      \"state\": \"action_required\"\n    }\n  }\n}\n")
	err := compareBaselinePlanCharacterization("clean-adoption", want, got)
	if err == nil {
		t.Fatal("changed characterization did not fail comparison")
	}
	if !strings.Contains(err.Error(), `shape "clean-adoption"`) ||
		!strings.Contains(err.Error(), "$.outcome.Result.state") {
		t.Fatalf("comparison error = %q, want shape and changed field", err)
	}
}

func TestBaselinePlanCharacterizationPublicCommandJourneys(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}
	command := exec.Command(
		"go",
		"test",
		"./internal/cli",
		"-run",
		"^TestBaselineMacroJourneysPublicCLI$",
		"-count=1",
	)
	command.Dir = projectRoot
	command.Env = os.Environ()
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("public command characterization journeys changed: %v\n%s", err, output)
	}
}

func buildCleanAdoptionCharacterization(t *testing.T) (string, PlanRequest, baselinePlanCharacterizationShape) {
	t.Helper()
	repository := newBaselinePlanCharacterizationRepository(t, true, true, true)
	return repository, baselinePlanCharacterizationRequest(repository), baselinePlanCharacterizationShape{
		ProfileID: "go-cli-tui",
	}
}

func buildIdempotentReplanCharacterization(t *testing.T) (string, PlanRequest, baselinePlanCharacterizationShape) {
	t.Helper()
	repository := newBaselinePlanCharacterizationRepository(t, true, true, true)
	request := baselinePlanCharacterizationRequest(repository)
	outcome, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("build pre-apply characterization plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("pre-apply characterization plan returned result: %+v", outcome.Result)
	}
	result, err := ApplyPlan(context.Background(), repository, *outcome.Plan, outcome.Plan.PlanDigest)
	if err != nil {
		t.Fatalf("apply characterization plan: %v", err)
	}
	if result.State != "verified" {
		t.Fatalf("characterization apply state = %q, want verified", result.State)
	}
	return repository, request, baselinePlanCharacterizationShape{
		ProfileID:          "go-cli-tui",
		VerifiedApplyState: result.State,
	}
}

func buildBlockingCapabilitiesCharacterization(t *testing.T) (string, PlanRequest, baselinePlanCharacterizationShape) {
	t.Helper()
	repository := newBaselinePlanCharacterizationRepository(t, false, true, true)
	return repository, baselinePlanCharacterizationRequest(repository), baselinePlanCharacterizationShape{
		ProfileID: "go-cli-tui",
	}
}

func buildAdvisoryDivergencesCharacterization(t *testing.T) (string, PlanRequest, baselinePlanCharacterizationShape) {
	t.Helper()
	repository := newBaselinePlanCharacterizationRepository(t, true, true, false)
	return repository, baselinePlanCharacterizationRequest(repository), baselinePlanCharacterizationShape{
		ProfileID: "go-cli-tui",
	}
}

func buildSameBaselineDigestDriftCharacterization(t *testing.T) (string, PlanRequest, baselinePlanCharacterizationShape) {
	t.Helper()
	repository := newBaselinePlanCharacterizationRepository(t, true, true, true)
	request := baselinePlanCharacterizationRequest(repository)
	outcome, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("build digest-drift fixture plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("digest-drift fixture plan returned result: %+v", outcome.Result)
	}
	result, err := ApplyPlan(context.Background(), repository, *outcome.Plan, outcome.Plan.PlanDigest)
	if err != nil {
		t.Fatalf("apply digest-drift fixture plan: %v", err)
	}
	if result.State != "verified" {
		t.Fatalf("digest-drift fixture apply state = %q, want verified", result.State)
	}

	manifestFile := filepath.Join(repository, filepath.FromSlash(manifestPath))
	manifestBytes, err := os.ReadFile(manifestFile)
	if err != nil {
		t.Fatalf("read digest-drift Setup Manifest: %v", err)
	}
	var manifest SetupManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("decode digest-drift Setup Manifest: %v", err)
	}
	manifest.ProfileDigest = "sha256:" + strings.Repeat("1", 64)
	manifest.CatalogDigest = "sha256:" + strings.Repeat("2", 64)
	manifestBytes, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode digest-drift Setup Manifest: %v", err)
	}
	manifestBytes = append(manifestBytes, '\n')
	if err := os.WriteFile(manifestFile, manifestBytes, 0o644); err != nil {
		t.Fatalf("write digest-drift Setup Manifest: %v", err)
	}

	return repository, request, baselinePlanCharacterizationShape{
		ProfileID:             "go-cli-tui",
		VerifiedApplyState:    result.State,
		ExistingBaseline:      manifest.Generator.Baseline,
		ExistingProfileDigest: manifest.ProfileDigest,
		ExistingCatalogDigest: manifest.CatalogDigest,
	}
}

func baselinePlanCharacterizationRequest(repository string) PlanRequest {
	return PlanRequest{
		Repository:   repository,
		ProfileID:    "go-cli-tui",
		Decisions:    planTestDecisions(),
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	}
}

func newBaselinePlanCharacterizationRepository(
	t *testing.T,
	includeContext7 bool,
	includeExa bool,
	includeFirecrawl bool,
) string {
	t.Helper()
	repository := newInspectionRepository(t)
	if includeContext7 {
		writeInspectionFile(t, repository, ".agents/skills/context7/SKILL.md", "# Context7\n")
	}
	if includeExa {
		writeInspectionFile(t, repository, ".agents/skills/exa-web-search/SKILL.md", "# Exa\n")
	}
	if includeFirecrawl {
		writeInspectionFile(t, repository, ".agents/skills/firecrawl/SKILL.md", "# Firecrawl\n")
	}
	writeInspectionFile(t, repository, "Makefile", "verify:\n\t@true\n")
	commitInspectionRepository(t, repository, "seed Baseline plan characterization")
	return repository
}

func baselinePlanCharacterizationBin(t *testing.T) string {
	t.Helper()
	directory := t.TempDir()
	for _, name := range []string{"rg", "rtk"} {
		path := filepath.Join(directory, name)
		if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write controlled %s executable: %v", name, err)
		}
	}
	return directory
}

func marshalBaselinePlanCharacterization(
	t *testing.T,
	record baselinePlanCharacterization,
	repository string,
	characterizationBin string,
) []byte {
	t.Helper()
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s characterization: %v", record.Shape, err)
	}
	replacements := []struct {
		from string
		to   string
	}{
		{from: filepath.ToSlash(repository), to: "<repository>"},
		{from: filepath.ToSlash(characterizationBin), to: "<characterization-bin>"},
	}
	for _, replacement := range replacements {
		data = bytes.ReplaceAll(data, []byte(replacement.from), []byte(replacement.to))
	}
	return append(data, '\n')
}

func writeBaselinePlanCharacterizationGolden(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create characterization golden directory: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write characterization golden %s: %v", path, err)
	}
	t.Logf("regenerated %s", path)
}

func compareBaselinePlanCharacterization(shape string, expected, actual []byte) error {
	if bytes.Equal(expected, actual) {
		return nil
	}
	want, err := decodeBaselinePlanCharacterizationJSON(expected)
	if err != nil {
		return fmt.Errorf("shape %q: decode recorded characterization: %w", shape, err)
	}
	got, err := decodeBaselinePlanCharacterizationJSON(actual)
	if err != nil {
		return fmt.Errorf("shape %q: decode current characterization: %w", shape, err)
	}
	path, wantValue, gotValue, different := firstBaselinePlanCharacterizationDifference("$", want, got)
	if !different {
		path = "$ (JSON encoding)"
		wantValue = string(expected)
		gotValue = string(actual)
	}
	return fmt.Errorf(
		"shape %q changed at %s:\n- want %s\n+ got  %s\nregenerate deliberately with -update-baseline-plan-characterization",
		shape,
		path,
		formatBaselinePlanCharacterizationValue(wantValue),
		formatBaselinePlanCharacterizationValue(gotValue),
	)
}

func decodeBaselinePlanCharacterizationJSON(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

func firstBaselinePlanCharacterizationDifference(path string, want, got any) (string, any, any, bool) {
	if reflect.DeepEqual(want, got) {
		return "", nil, nil, false
	}
	wantMap, wantIsMap := want.(map[string]any)
	gotMap, gotIsMap := got.(map[string]any)
	if wantIsMap && gotIsMap {
		keys := make([]string, 0, len(wantMap)+len(gotMap))
		seen := make(map[string]struct{}, len(wantMap)+len(gotMap))
		for key := range wantMap {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
		for key := range gotMap {
			if _, ok := seen[key]; !ok {
				keys = append(keys, key)
			}
		}
		sort.Strings(keys)
		for _, key := range keys {
			wantValue, wantOK := wantMap[key]
			gotValue, gotOK := gotMap[key]
			fieldPath := path + "." + key
			if !wantOK || !gotOK {
				return fieldPath, wantValue, gotValue, true
			}
			if differencePath, differenceWant, differenceGot, different := firstBaselinePlanCharacterizationDifference(
				fieldPath,
				wantValue,
				gotValue,
			); different {
				return differencePath, differenceWant, differenceGot, true
			}
		}
	}
	wantSlice, wantIsSlice := want.([]any)
	gotSlice, gotIsSlice := got.([]any)
	if wantIsSlice && gotIsSlice {
		limit := min(len(wantSlice), len(gotSlice))
		for index := 0; index < limit; index++ {
			if differencePath, differenceWant, differenceGot, different := firstBaselinePlanCharacterizationDifference(
				fmt.Sprintf("%s[%d]", path, index),
				wantSlice[index],
				gotSlice[index],
			); different {
				return differencePath, differenceWant, differenceGot, true
			}
		}
		if len(wantSlice) != len(gotSlice) {
			return path + ".length", len(wantSlice), len(gotSlice), true
		}
	}
	return path, want, got, true
}

func formatBaselinePlanCharacterizationValue(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(data)
}
