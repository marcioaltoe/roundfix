// Suite: Baseline Profile decision gaps
// Invariant: a missing template decision is named before rendering and never mutates the repository.
// Boundary IN: profile validate, plan and update streams, exit codes, JSON and repository bytes.
// Boundary OUT: catalog reachability and binding selection, owned by internal/baseline.

package cli

import (
	"encoding/json"
	"os"
	"slices"
	"testing"

	"roundfix/internal/baseline"
)

const profileGapMessage = "profile sqllib does not list decision verification.incremental, which template template.guide.agent-instructions requires"
const profileGapRepair = `add "verification.incremental" to "decisions" in .roundfix/baseline/profiles/sqllib.json, then rerun with --decision 'verification.incremental=make verify-incremental'`
const profileGapValidationMessage = `custom.profile.decision.template-required: ` + profileGapMessage + `; add "verification.incremental" to "decisions" in .roundfix/baseline/profiles/sqllib.json, then run roundfix baseline update --decision 'verification.incremental=make verify-incremental'`

func TestBaselineProfileValidateNamesADecisionGap(t *testing.T) {
	t.Parallel()
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, "Makefile", "verify:\n\t@true\nverify-incremental:\n\t@true\n")
	writeProfileDecisionGapFixture(t, repository, false)
	commitBaselinePlanTestRepository(t, repository)
	setCommandWorkDirForTest(t, repository)
	before := baselinePlanTestTree(t, repository)
	wantErr := "roundfix: baseline profile.validate failed: " + profileGapValidationMessage + "\nRun 'roundfix baseline profile --help' for usage.\n"
	for _, target := range []string{"sqllib", "", ".roundfix/baseline/profiles/sqllib.json"} {
		args := []string{"baseline", "profile", "validate"}
		if target != "" {
			args = append(args, target)
		}
		stdout, stderr, code := runRepositoryProfileCommand(t, args...)
		if code != exitPreflight || stdout != "" || stderr != wantErr {
			t.Fatalf("Surface Transcript 1 (%q): exit=%d stdout=%q stderr=%q; want exit=2 stdout=empty stderr=%q", target, code, stdout, stderr, wantErr)
		}
		stdout, stderr, code = runRepositoryProfileCommand(t, append(args, "--format", "json")...)
		var result baselineProfileResult
		if err := json.Unmarshal([]byte(stdout), &result); err != nil {
			t.Fatalf("decode validate failure: %v; %s", err, stdout)
		}
		if code != exitPreflight || stderr != wantErr || result.State != "failed" || result.Category != "preflight" || result.Message != profileGapValidationMessage {
			t.Fatalf("validate JSON: exit=%d result=%+v stderr=%q", code, result, stderr)
		}
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("profile validation changed repository bytes")
	}
}

func TestBaselinePlanRefusesAProfileDecisionGap(t *testing.T) {
	t.Parallel()
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, "Makefile", "verify:\n\t@true\nverify-incremental:\n\t@true\n")
	writeProfileDecisionGapFixture(t, repository, false)
	commitBaselinePlanTestRepository(t, repository)
	setCommandWorkDirForTest(t, repository)
	before := baselinePlanTestTree(t, repository)
	stdout, stderr, code := runRepositoryProfileCommand(t, "baseline", "plan", "--profile", "sqllib", "--repo", repository, "--format", "text", "--decision", "preservation.mode=greenfield")
	wantOut := "Baseline plan: action required\nCategory: decision\nResult: " + profileGapMessage + "\nNext action: " + profileGapRepair + "\n"
	if code != exitUnverified || stdout != wantOut || stderr != "" {
		t.Fatalf("Surface Transcript 2: exit=%d stdout=%q stderr=%q; want exit=3 stdout=%q stderr=empty", code, stdout, stderr, wantOut)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("refused plan changed repository bytes")
	}
}

func TestBaselineUpdateRefusesAProfileDecisionGap(t *testing.T) {
	t.Parallel()
	repository := profileDecisionGapUpdateRepository(t, false)
	before := baselinePlanTestTree(t, repository)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format", "json")
	if code != exitUnverified || stderr != "" || result.State != "action_required" || result.Category != "decision" || result.Message != profileGapMessage || result.NextAction != profileGapRepair || len(result.NewDecisions) != 0 {
		t.Fatalf("gap update: exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("refused update changed repository bytes")
	}
}

func TestBaselineUpdateAsksForTheDecisionOnceTheProfileListsIt(t *testing.T) {
	t.Parallel()
	repository := profileDecisionGapUpdateRepository(t, true)
	before := baselinePlanTestTree(t, repository)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format", "json")
	if code != exitUnverified || stderr != "" || result.State != "action_required" || result.Category != "decision" || len(result.NewDecisions) != 1 || result.NewDecisions[0].ID != "verification.incremental" {
		t.Fatalf("listed decision update: exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("new-decision update changed repository bytes")
	}
}

func writeProfileDecisionGapFixture(t *testing.T, repository string, listIncremental bool) baseline.ResolvedProfile {
	t.Helper()
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	created, err := baseline.InitCustomProfile(repository, "sqllib", "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(created.Path)
	if err != nil {
		t.Fatal(err)
	}
	var profile map[string]json.RawMessage
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatal(err)
	}
	if !listIncremental {
		decisions := slices.DeleteFunc(slices.Clone(created.Profile.Decisions), func(id string) bool { return id == "verification.incremental" })
		profile["decisions"], err = json.Marshal(decisions)
		if err != nil {
			t.Fatal(err)
		}
	}
	data, err = json.MarshalIndent(profile, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(created.Path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	resolved, err := baseline.ResolveProfile(repository, "sqllib", catalog)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func profileDecisionGapUpdateRepository(t *testing.T, listIncremental bool) string {
	t.Helper()
	repository := singleGateBaselineUpdateRepository(t, true)
	profile := writeProfileDecisionGapFixture(t, repository, listIncremental)
	manifest := ReadBaselineSetupManifest(t, repository)
	manifest.Profile = profile.ID
	manifest.Generator.Baseline = "baseline." + profile.ID + "-" + baseline.ManifestVersion
	manifest.ProfileDigest = profile.Digest
	WriteBaselineSetupManifest(t, repository, manifest)
	setCommandWorkDirForTest(t, repository)
	return repository
}
