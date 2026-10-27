// Suite: repository-profile Baseline and Doctor documentation contracts.
// Invariant: repository profiles use their modules' embedded skill contracts and Setup Snapshots.
// Boundary IN: live repository guides, skills, profile adoption, update preview, and Doctor skills checks.
// Boundary OUT: source acquisition during refresh is replaced by a restore fake; machine Doctor checks are not exercised.
// These tests read the live documentation trees and run under the docscontract tag per ADR-0326.
//go:build docscontract

//verify:always

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

func TestBaselineUpdateReachesCurrentWithARepositoryProfile(t *testing.T) {
	t.Parallel()
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

func TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile(t *testing.T) {
	t.Parallel()
	root := newBaselineRepositoryProfileFixture(t, true)
	editRepositoryProfileTestSkill(t, root)
	before := baselinePlanTestTree(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "update", "--repo", root, "--format", "json")
	var result baselineUpdateResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	if code != exitOK || stderr != "" || result.State != "current" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	found := false
	for _, difference := range result.Skills.Differing {
		if difference.Skill == "testing-boss" && difference.Reason == "differs from its Setup Snapshot, newer or older unknown; kept; restore deliberately with roundfix baseline skills restore --profile "+baselineRepositoryProfileTestID+" --skill testing-boss" {
			found = true
		}
	}
	if !found || len(result.Skills.Drifted) != 0 {
		t.Fatalf("differing skill absent from preview: %s", out)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("preview changed repository bytes")
	}
}

func TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile(t *testing.T) {
	t.Parallel()
	root := newBaselineRepositoryProfileFixture(t, true)
	editRepositoryProfileTestSkill(t, root)
	// Force the real install to replace every owned skill, rather than just report them.
	for _, name := range skills.Names() {
		if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", name)); err != nil {
			t.Fatal(err)
		}
	}
	beforeSkill := mustReadBytes(t, filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md"))
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{
		Repository: root, ProfileID: baselineRepositoryProfileTestID,
	}, baselineUpdateSkillsDependencies{
		recordLock:          baseline.RecordInstalledOwnedSkills,
		resolveProjectRoot:  defaultResolveSkillsProjectRoot,
		install:             skills.Install,
		ownedNames:          skills.Names,
		resolveExternal:     resolveExternalSkillRequirement,
		checkRepository:     skills.CheckRepositoryWithExternal,
		snapshotDifferences: baseline.SetupSnapshotDifferences,
		restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
			t.Fatalf("restore called for a lock-matching difference: %+v", request)
			return baseline.SkillsRestorePayload{}, nil
		},
	})
	wantInstalled := slices.Clone(skills.Names())
	slices.Sort(wantInstalled)
	if err != nil || result.Status != baselineUpdateSkillsVerified ||
		result.InstalledCount != len(wantInstalled) || !slices.Equal(result.Installed, wantInstalled) ||
		len(result.Restored) != 0 || len(result.Differing) != 1 || result.Differing[0].Skill != "testing-boss" || len(result.Drifted) != 0 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if !bytes.Equal(beforeSkill, mustReadBytes(t, filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md"))) {
		t.Fatal("skills refresh changed the differing skill's bytes")
	}
	readiness, err := skills.CheckRepository(t.Context(), root)
	if err != nil || !readiness.Ready() {
		t.Fatalf("owned install readiness=%+v err=%v", readiness, err)
	}
}

func TestDoctorComparesARepositoryProfileWithItsSnapshot(t *testing.T) {
	t.Parallel()
	root := newBaselineRepositoryProfileFixture(t, true)
	editRepositoryProfileTestSkill(t, root)
	before := baselinePlanTestTree(t, root)
	result := repositorySkillsCheck(t.Context(), defaultDoctorDependencies(), root)
	if result.Status != CheckStatusOK || result.NextAction != "" || !strings.Contains(result.Detail, "DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer or older unknown: testing-boss") ||
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

func TestBaselineUpdatePreviewListsASetupSnapshotDifference(t *testing.T) {
	t.Parallel()
	root := newBaselineUpdateRepository(t)
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{".agents/skills", "skills-lock.json"} {
		copyThisRepositorySkillSetPath(t, filepath.Join(source, relative), filepath.Join(root, relative))
	}
	// Make the installed tree intentionally different while keeping its lock intact.
	path := filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md")
	if err := os.WriteFile(path, append(mustReadBytes(t, path), []byte("\nDifferent upstream installation fixture.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
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
	if err := os.WriteFile(lockPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	before := baselinePlanTestTree(t, root)
	result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--format=json")
	if code != exitOK || result.State != "current" || stderr != "" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	found := false
	for _, difference := range result.Skills.Differing {
		if difference.Skill == "testing-boss" && difference.Reason == "differs from its Setup Snapshot, newer or older unknown; kept; restore deliberately with roundfix baseline skills restore --profile go-cli-tui --skill testing-boss" {
			found = true
		}
	}
	if !found || len(result.Skills.Drifted) != 0 {
		t.Fatalf("differing skill absent from preview: %s", out)
	}
	_, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--format=text")
	if code != exitOK || stderr != "" || !strings.Contains(out, "Skills drifted: 0\nSkills differing: 1\n- differs testing-boss: differs from its Setup Snapshot, newer or older unknown; kept; restore deliberately with roundfix baseline skills restore --profile go-cli-tui --skill testing-boss\n") {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	if baselinePlanTestTree(t, root) != before {
		t.Fatal("preview changed repository bytes")
	}
	result, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--no-skills", "--format=json")
	if code != exitOK || stderr != "" || result.Skills.Differing == nil || len(result.Skills.Differing) != 0 || !strings.Contains(out, `"differing":[]`) {
		t.Fatalf("--no-skills exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	beforeSkill := mustReadBytes(t, path)
	result, out, stderr, code = runBaselineUpdateTestCommandWithSkillsStage(t, t.Context(), runBaselineUpdateSkillsStage,
		"baseline", "update", "--repo", root, "--yes", "--format=json")
	if code != exitOK || stderr != "" || result.Skills.Status != baselineUpdateSkillsVerified || len(result.Skills.Restored) != 0 || len(result.Skills.Differing) != 1 {
		t.Fatalf("--yes exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	if !bytes.Equal(beforeSkill, mustReadBytes(t, path)) {
		t.Fatal("--yes changed the differing skill's bytes")
	}
	afterHash, err := skills.SkillFolderHash(t.Context(), filepath.Dir(path))
	if err != nil || afterHash != hash {
		t.Fatalf("--yes skill tree hash=%q, want %q; error=%v", afterHash, hash, err)
	}
	t.Run("difference with outdated and retired skills", func(t *testing.T) {
		found, required := lowerPreviewOwnedSkill(t, root, "qa-gate")
		retiredDir := filepath.Join(root, ".agents", "skills", "council")
		if err := os.MkdirAll(retiredDir, 0o755); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(retiredDir, "SKILL.md"), "retired skill fixture\n")
		beforePreview := baselinePlanTestTree(t, root)
		_, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--format=text")
		wantOwned := fmt.Sprintf("Skills drifted: 0\nSkills outdated: 1\n- outdated qa-gate: found %s, requires %s\n", found, required)
		wantRetired := "Skills retired: 1\n- retired council: no longer required by the Baseline; delete .agents/skills/council\n"
		wantDifference := "Skills differing: 1\n- differs testing-boss: differs from its Setup Snapshot, newer or older unknown; kept; restore deliberately with roundfix baseline skills restore --profile go-cli-tui --skill testing-boss\n"
		if code != exitUnverified || stderr != "" || !strings.Contains(out, wantOwned+wantRetired+wantDifference) {
			t.Fatalf("combined preview exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
		if baselinePlanTestTree(t, root) != beforePreview {
			t.Fatal("combined preview changed repository bytes")
		}
	})
}

func editRepositoryProfileTestSkill(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md")
	mustWrite(t, path, string(mustReadBytes(t, path))+"\nDifferent upstream installation fixture.\n")
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
