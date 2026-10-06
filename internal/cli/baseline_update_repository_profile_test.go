// Suite: repository-profile Baseline commands
// Invariant: repository profiles use their modules' embedded skill contracts.
// Boundary IN: real profile adoption, update preview, owned install, Doctor skills check, and restore refusal streams.
// Boundary OUT: source acquisition during refresh is replaced by a restore fake; machine Doctor checks are not exercised.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/skills"
)

const baselineRepositoryProfileTestID = "repository-go-cli-tui"

func TestBaselineUpdateReachesCurrentWithARepositoryProfile(t *testing.T) {
	for _, copied := range []bool{true, false} {
		t.Run(fmt.Sprintf("copied skills=%t", copied), func(t *testing.T) {
			root := newBaselineRepositoryProfileFixture(t, copied)
			before := baselinePlanTestTree(t, root)
			out, stderr, code := runRepositoryProfileCommand(t, "baseline", "update", "--repo", root, "--format", "json")
			var result baselineUpdateResult
			if err := json.Unmarshal([]byte(out), &result); err != nil {
				t.Fatal(err)
			}
			if code != exitOK || stderr != "" || result.State != "current" {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			out, stderr, code = runRepositoryProfileCommand(t, "baseline", "update", "--repo", root, "--format", "text")
			const prefix = "Baseline update: current\nResult: the repository already matches the current Baseline catalog\n"
			if code != exitOK || stderr != "" || !strings.HasPrefix(out, prefix) {
				t.Fatalf("Surface Transcript 1: exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			if before != baselinePlanTestTree(t, root) {
				t.Fatal("current update changed repository bytes")
			}
		})
	}
}

func TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile(t *testing.T) {
	root := newBaselineRepositoryProfileFixture(t, true)
	editRepositoryProfileTestSkill(t, root)
	before := baselinePlanTestTree(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "update", "--repo", root, "--format", "json")
	var result baselineUpdateResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	if code != exitUnverified || stderr != "" || result.State != "plan_ready" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	found := false
	for _, drift := range result.Skills.Drifted {
		if drift.Skill == "testing-boss" && strings.Contains(drift.Reason, "trails its Setup Snapshot") {
			found = true
		}
	}
	if !found {
		t.Fatalf("trailing skill absent from preview: %s", out)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("preview changed repository bytes")
	}
}

func TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile(t *testing.T) {
	root := newBaselineRepositoryProfileFixture(t, true)
	editRepositoryProfileTestSkill(t, root)
	// Force the real install to replace every owned skill, rather than just report them.
	for _, name := range skills.Names() {
		if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", name)); err != nil {
			t.Fatal(err)
		}
	}
	var requests []baseline.SkillsRestoreRequest
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{
		Repository: root, ProfileID: baselineRepositoryProfileTestID,
	}, baselineUpdateSkillsDependencies{
		resolveProjectRoot: defaultResolveSkillsProjectRoot,
		install:            skills.Install,
		ownedNames:         skills.Names,
		resolveExternal:    resolveExternalSkillRequirement,
		checkRepository:    skills.CheckRepositoryWithExternal,
		trailingSkills:     baseline.TrailingSetupSkills,
		restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
			requests = append(requests, request)
			return baseline.SkillsRestorePayload{Applied: true}, nil
		},
	})
	wantInstalled := slices.Clone(skills.Names())
	slices.Sort(wantInstalled)
	if err != nil || result.Status != baselineUpdateSkillsVerified ||
		result.InstalledCount != len(wantInstalled) || !slices.Equal(result.Installed, wantInstalled) ||
		!slices.Equal(result.Restored, []string{"testing-boss"}) || len(result.Drifted) != 0 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if len(requests) != 1 || requests[0].Repository != root || requests[0].ProfileID != baselineRepositoryProfileTestID ||
		!slices.Equal(requests[0].Skills, []string{"testing-boss"}) {
		t.Fatalf("restore requests=%+v", requests)
	}
	readiness, err := skills.CheckRepository(t.Context(), root)
	if err != nil || !readiness.Ready() {
		t.Fatalf("owned install readiness=%+v err=%v", readiness, err)
	}
}

func TestDoctorComparesARepositoryProfileWithItsSnapshot(t *testing.T) {
	root := newBaselineRepositoryProfileFixture(t, true)
	editRepositoryProfileTestSkill(t, root)
	before := baselinePlanTestTree(t, root)
	result := repositorySkillsCheck(t.Context(), defaultDoctorDependencies(), root)
	if result.Status != CheckStatusWarn || !strings.Contains(result.Detail, "DR-SKILL-TRAILS-SNAPSHOT: trails the Setup Snapshot: testing-boss") ||
		strings.Contains(result.Detail, "snapshot comparison unavailable") {
		t.Fatalf("Doctor skills=%+v", result)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("Doctor changed repository bytes")
	}
	path := filepath.Join(root, filepath.FromSlash(doctorSetupManifestPath))
	var manifest map[string]any
	if err := json.Unmarshal(mustReadBytes(t, path), &manifest); err != nil {
		t.Fatal(err)
	}
	manifest["profile"] = "missing-profile"
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, path, string(data))
	before = baselinePlanTestTree(t, root)
	result = repositorySkillsCheck(t.Context(), defaultDoctorDependencies(), root)
	if result.Status != CheckStatusOK || !strings.Contains(result.Detail, "snapshot comparison unavailable") ||
		!strings.Contains(result.Detail, "neither built-in nor resolvable at .roundfix/baseline/profiles/missing-profile.json") ||
		strings.Contains(result.Detail, "built-in Baseline Profile") {
		t.Fatalf("Doctor missing-profile skills=%+v", result)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("unavailable comparison changed repository bytes")
	}
}

func TestBaselineSkillsRestoreNamesTheProfilePath(t *testing.T) {
	root := newBaselineRepositoryProfileFixture(t, false)
	before := baselinePlanTestTree(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "skills", "restore", "--profile", "missing-profile", "--repo", root, "--format", "text")
	const message = `Baseline Profile "missing-profile" is neither built-in nor resolvable at .roundfix/baseline/profiles/missing-profile.json.`
	wantOut := "Baseline skills restore: blocked\nrestore.profile-unresolved: " + message + "\nNext action: Run roundfix baseline profile validate missing-profile to see why, restore that file, or choose a built-in profile id.\n"
	wantErr := "roundfix: baseline skills restore failed: " + message + "\n"
	if code != exitPreflight || out != wantOut || stderr != wantErr {
		t.Fatalf("Surface Transcript 2: exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("blocked restore changed repository bytes")
	}
}

func TestBaselineSkillsRestoreAcceptsARepositoryProfile(t *testing.T) {
	root := newBaselineRepositoryProfileFixture(t, false)
	source := filepath.Join(root, "missing-source")
	before := baselinePlanTestTree(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "skills", "restore", "--profile", baselineRepositoryProfileTestID,
		"--skill", "context7-cli", "--source-dir", source, "--repo", root, "--format", "text")
	message := "Offline Git object store is not a directory: " + source + "."
	wantOut := "Baseline skills restore: blocked\nrestore.source-dir-invalid: " + message + "\nNext action: Pass an existing Git checkout or bare object store to --source-dir.\n"
	wantErr := "roundfix: baseline skills restore failed: " + message + "\n"
	if code != exitPreflight || out != wantOut || stderr != wantErr {
		t.Fatalf("Surface Transcript 3: exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("blocked restore changed repository bytes")
	}
}

func newBaselineRepositoryProfileFixture(t *testing.T, copySkills bool) string {
	t.Helper()
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	root := newBaselineApplyTestRepository(t)
	t.Chdir(root) // profile init intentionally resolves the process working directory.
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "profile", "init", "--id", baselineRepositoryProfileTestID, "--from", "go-cli-tui")
	if code != exitOK || stderr != "" {
		t.Fatalf("profile init exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	commitBaselinePlanTestRepository(t, root)
	args := []string{"baseline", "plan", "--repo", root, "--profile", baselineRepositoryProfileTestID,
		"--decision", "preservation.mode=greenfield",
		"--decision", "language.generated=English",
		"--decision", "verification.gate=make verify",
		"--decision", "verification.incremental=make verify-incremental",
		"--decision", "branch.prefix=ma/",
		"--decision", "spec.scaffold=true",
		"--decision", "domain.layout=single-context",
		"--decision", "triage.external=false",
		"--decision", "autonomous.enabled=true",
		"--decision", "runtime.backend=codex gpt-5.5 xhigh",
		"--decision", "runtime.design=claude opus xhigh",
		"--decision", "secondbrain.enabled=false",
		"--decision", "repository.extension.enabled=false", "--format", "json"}
	out, stderr, code = runRepositoryProfileCommand(t, args...)
	if code != exitOK || stderr != "" {
		t.Fatalf("plan exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	plan, err := baseline.ParsePlanDocument([]byte(out))
	if err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(t.TempDir(), "plan.json")
	mustWrite(t, planPath, out)
	out, stderr, code = runRepositoryProfileCommand(t, "baseline", "apply", "--repo", root, "--plan", planPath, "--confirm-plan", plan.PlanDigest, "--format", "json")
	if code != exitOK || stderr != "" {
		t.Fatalf("apply exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	if copySkills {
		for _, relative := range []string{".agents/skills", "skills-lock.json"} {
			copyThisRepositorySkillSetPath(t, filepath.Join(source, relative), filepath.Join(root, relative))
		}
	}
	commitBaselinePlanTestRepository(t, root)
	return root
}

func runRepositoryProfileCommand(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := RunContext(t.Context(), args, &stdout, &stderr)
	return stdout.String(), stderr.String(), code
}

func editRepositoryProfileTestSkill(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md")
	mustWrite(t, path, string(mustReadBytes(t, path))+"\nOld upstream installation fixture.\n")
	lockPath := filepath.Join(root, "skills-lock.json")
	var lock map[string]any
	if err := json.Unmarshal(mustReadBytes(t, lockPath), &lock); err != nil {
		t.Fatal(err)
	}
	hash, err := skills.SkillFolderHash(t.Context(), filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	lock["skills"].(map[string]any)["testing-boss"].(map[string]any)["computedHash"] = hash
	data, err := json.Marshal(lock)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, lockPath, string(data))
	readiness, err := skills.CheckRepositoryWithExternal(t.Context(), root, []string{"testing-boss"})
	if err != nil || !readiness.Ready() {
		t.Fatalf("edited skill must still match its lock: %+v %v", readiness, err)
	}
}
