// Suite: Retired Skill reporting in Baseline update.
// Invariant: approval deletes only proven unedited Retired Skills and keeps every other copy.
// Boundary IN: adopted repository, real preview/apply, output, and exit codes.
// Boundary OUT: network acquisition (the existing skills-stage seam is injected).

package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/skills"
)

const retiredSkillsTranscript = "Skills retired: 2\n" +
	"- retired council: kept, skills-lock.json records no computedHash to prove its copy unedited; delete .agents/skills/council by hand if the repository no longer uses it\n" +
	"- retired the-fool: kept, skills-lock.json records no computedHash to prove its copy unedited; delete .agents/skills/the-fool and its skills-lock.json entry by hand if the repository no longer uses it\n"

func TestBaselineUpdateRetiredSkillDeletionTargets(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		paths []string
		lock  bool
		want  string
	}{
		{name: "lock only", lock: true, want: "its skills-lock.json entry"},
		{name: "two paths", paths: []string{".agents/skills/council", ".claude/skills/council"}, want: ".agents/skills/council and .claude/skills/council"},
		{name: "both paths and lock", paths: []string{".agents/skills/council", ".claude/skills/council"}, lock: true, want: ".agents/skills/council, .claude/skills/council and its skills-lock.json entry"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := newBaselineUpdateResult()
			result.Skills.Retired = []baseline.InstalledRetiredSkill{{Skill: "council", Paths: test.paths, LockEntry: test.lock, Action: baseline.RetiredSkillDelete}}
			var out bytes.Buffer
			if err := writeBaselineUpdateResult(result, false, &out); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), "Skills retired: 1\n- retired council: no longer required by the Baseline; the approved update deletes "+test.want+"\n") {
				t.Fatalf("retired deletion targets = %s", out.String())
			}
		})
	}
	t.Run("digest binds only planned deletions", func(t *testing.T) {
		digest := "original plan digest"
		kept := baseline.InstalledRetiredSkill{Skill: "council", Action: baseline.RetiredSkillKeep, Reason: baseline.RetiredSkillUnrecorded}
		deleted := baseline.InstalledRetiredSkill{Skill: "tdd", Paths: []string{".agents/skills/tdd"}, LockEntry: true, Action: baseline.RetiredSkillDelete, LockHash: "recorded hash"}
		if got := retiredSkillsPlanDigest(digest, []baseline.InstalledRetiredSkill{kept}); got != digest {
			t.Fatalf("keep changed digest: %s", got)
		}
		data, err := json.Marshal(map[string]any{"planDigest": digest, "retiredSkills": []baseline.InstalledRetiredSkill{deleted}})
		if err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(data)
		want := "sha256:" + hex.EncodeToString(hash[:])
		if got := retiredSkillsPlanDigest(digest, []baseline.InstalledRetiredSkill{kept, deleted}); got != want {
			t.Fatalf("deletion digest = %s; want %s", got, want)
		}
	})
	t.Run("kept unsafe copy", func(t *testing.T) {
		result := newBaselineUpdateResult()
		result.Skills.Retired = []baseline.InstalledRetiredSkill{{Skill: "council", Paths: []string{".agents/skills/council"}, Action: baseline.RetiredSkillKeep, Reason: baseline.RetiredSkillUnsafe}}
		var out bytes.Buffer
		if err := writeBaselineUpdateResult(result, false, &out); err != nil {
			t.Fatal(err)
		}
		want := "- retired council: kept, an entry is not a skill folder or a link to .agents/skills; delete .agents/skills/council by hand if the repository no longer uses it\n"
		if !strings.Contains(out.String(), want) {
			t.Fatalf("unsafe keep text = %s", out.String())
		}
	})
}

func TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState(t *testing.T) {
	t.Parallel()
	repo := adoptedRepositoryWithRetiredSkills(t)
	before := baselinePlanTestTree(t, repo)
	_, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(),
		"baseline", "update", "--repo", repo, "--format=text")
	if code != exitOK || stderr != "" || !strings.Contains(stdout, "Baseline update: current\n") ||
		!strings.Contains(stdout, "Skills drifted: 0\n"+retiredSkillsTranscript) {
		t.Fatalf("retired preview exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(),
		"baseline", "update", "--repo", repo, "--format=json")
	if code != exitOK || stderr != "" || result.State != "current" || result.Category != "" || result.NextAction != "" ||
		result.Message != "the repository already matches the current Baseline catalog" ||
		!reflect.DeepEqual(result.Skills.Retired, expectedInstalledRetiredSkills()) {
		t.Fatalf("retired JSON preview = %+v exit=%d stdout=%s stderr=%s", result, code, stdout, stderr)
	}
	if before != baselinePlanTestTree(t, repo) {
		t.Fatal("retired preview changed repository bytes")
	}
	t.Run("no retired copies", func(t *testing.T) {
		root := newBaselineUpdateRepository(t)
		for _, format := range []string{"text", "json"} {
			_, out, errOut, exit := runBaselineUpdateTestCommand(t, context.Background(),
				"baseline", "update", "--repo", root, "--format="+format)
			if exit != exitOK || errOut != "" || strings.Contains(out, "Skills retired") || strings.Contains(out, `"retired"`) {
				t.Fatalf("empty retired report exit=%d stdout=%s stderr=%s", exit, out, errOut)
			}
		}
	})
}

func TestBaselineUpdateWithoutSkillsOmitsRetiredSkills(t *testing.T) {
	t.Parallel()
	for _, approval := range []string{"", "--yes"} {
		for _, format := range []string{"text", "json"} {
			repo := deletableRetiredSkillsRepository(t)
			before := baselinePlanTestTree(t, repo)
			args := []string{"baseline", "update", "--repo", repo, "--no-skills", "--format=" + format}
			if approval != "" {
				args = append(args, approval)
			}
			stage := func(context.Context, baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
				t.Fatal("--no-skills ran the skills stage")
				return baselineUpdateSkillsResult{}, nil
			}
			result, out, errOut, code := runBaselineUpdateTestCommandWithSkillsStage(t, context.Background(), stage, args...)
			if code != exitOK || errOut != "" || len(result.Skills.Retired) != 0 || len(result.Skills.Deleted) != 0 || strings.Contains(out, "Skills retired") || strings.Contains(out, `"retired"`) {
				t.Fatalf("skipped retired report exit=%d stdout=%s stderr=%s", code, out, errOut)
			}
			if before != baselinePlanTestTree(t, repo) {
				t.Fatal("--no-skills changed current repository bytes")
			}
		}
	}
}

func TestBaselineUpdateAppliesAndListsRetiredSkills(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"text", "json"} {
		repo := adoptedRepositoryWithRetiredSkills(t)
		beforeSkills := baselinePlanTestTree(t, filepath.Join(repo, ".agents/skills"))
		beforeLock := string(mustReadBytes(t, filepath.Join(repo, "skills-lock.json")))
		manifest := ReadBaselineSetupManifest(t, repo)
		manifest.CatalogDigest = "sha256:" + strings.Repeat("0", 64)
		WriteBaselineSetupManifest(t, repo, manifest)
		commitBaselinePlanTestRepository(t, repo)
		result, out, errOut, code := runBaselineUpdateTestCommand(t, context.Background(),
			"baseline", "update", "--repo", repo, "--yes", "--format="+format)
		if code != exitOK || errOut != "" {
			t.Fatalf("applied retired report exit=%d stdout=%s stderr=%s", code, out, errOut)
		}
		if format == "text" {
			if !strings.Contains(out, "Baseline update: verified\n") || !strings.Contains(out, retiredSkillsTranscript) {
				t.Fatalf("applied text = %s", out)
			}
		} else if result.State != "verified" || result.Category != "" || !reflect.DeepEqual(result.Skills.Retired, expectedInstalledRetiredSkills()) {
			t.Fatalf("applied JSON = %+v", result)
		}
		if beforeSkills != baselinePlanTestTree(t, filepath.Join(repo, ".agents/skills")) || beforeLock != string(mustReadBytes(t, filepath.Join(repo, "skills-lock.json"))) {
			t.Fatal("applied update changed retired skills or lock")
		}
	}
}

func TestBaselineUpdateReportsRetiredSkillInspectionFailure(t *testing.T) {
	t.Parallel()
	repo := newBaselineUpdateRepository(t)
	mustWrite(t, filepath.Join(repo, "skills-lock.json"), "{")
	commitBaselinePlanTestRepository(t, repo)
	result, out, errOut, code := runBaselineUpdateTestCommand(t, context.Background(),
		"baseline", "update", "--repo", repo, "--format=json")
	if code != exitRunFailed || result.Category != "execution" || !strings.Contains(errOut, "skills-lock.json") ||
		result.NextAction != "repair skills-lock.json and rerun roundfix baseline update" {
		t.Fatalf("inspection failure = %+v exit=%d stdout=%s stderr=%s", result, code, out, errOut)
	}
}

func adoptedRepositoryWithRetiredSkills(t *testing.T) string {
	t.Helper()
	repo := newBaselineApplyTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/council/SKILL.md", "# council\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/the-fool/SKILL.md", "# the-fool\n")
	lock := map[string]any{"version": 1, "skills": map[string]any{"the-fool": map[string]any{}}}
	data, err := json.Marshal(lock)
	if err != nil {
		t.Fatal(err)
	}
	writeBaselinePlanTestFile(t, repo, "skills-lock.json", string(data))
	commitBaselinePlanTestRepository(t, repo)
	plan, _ := baselineApplyTestPlan(t, repo)
	if _, err := baseline.ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
		t.Fatalf("adopt retired skills fixture: %v", err)
	}
	commitBaselinePlanTestRepository(t, repo)
	installRequiredExternalSkills(t, repo)
	return repo
}

func expectedInstalledRetiredSkills() []baseline.InstalledRetiredSkill {
	return []baseline.InstalledRetiredSkill{
		{Skill: "council", Paths: []string{".agents/skills/council"}, Action: baseline.RetiredSkillKeep, Reason: baseline.RetiredSkillUnrecorded},
		{Skill: "the-fool", Paths: []string{".agents/skills/the-fool"}, LockEntry: true, Action: baseline.RetiredSkillKeep, Reason: baseline.RetiredSkillUnrecorded},
	}
}

func deletableRetiredSkillsRepository(t *testing.T) string {
	t.Helper()
	repo := newBaselineUpdateRepository(t)
	installRequiredExternalSkills(t, repo)
	lockPath := filepath.Join(repo, "skills-lock.json")
	var lock map[string]any
	if err := json.Unmarshal(mustReadBytes(t, lockPath), &lock); err != nil {
		t.Fatal(err)
	}
	entries := lock["skills"].(map[string]any)
	for _, name := range []string{"tdd", "design-patterns", "local-only"} {
		writeBaselinePlanTestFile(t, repo, ".agents/skills/"+name+"/SKILL.md", "# "+name+"\n")
		hash, err := skills.SkillFolderHash(t.Context(), filepath.Join(repo, ".agents", "skills", name))
		if err != nil {
			t.Fatal(err)
		}
		entries[name] = map[string]any{"computedHash": hash}
	}
	writeBaselinePlanTestFile(t, repo, ".agents/skills/design-patterns/SKILL.md", "# repository edit\n")
	for _, name := range []string{"tdd", "design-patterns", "context7-cli", "exa-web-search", "local-only"} {
		link := filepath.Join(repo, ".claude", "skills", name)
		if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("../../.agents/skills/"+name, link); err != nil {
			t.Fatal(err)
		}
	}
	data, err := json.Marshal(lock)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, lockPath, string(data))
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			repo := deletableRetiredSkillsRepository(t)
			before := baselinePlanTestTree(t, repo)
			lockPath := filepath.Join(repo, "skills-lock.json")
			var beforeLock map[string]any
			if err := json.Unmarshal(mustReadBytes(t, lockPath), &beforeLock); err != nil {
				t.Fatal(err)
			}
			wantRetired := []baseline.InstalledRetiredSkill{
				{Skill: "design-patterns", Paths: []string{".agents/skills/design-patterns", ".claude/skills/design-patterns"}, LockEntry: true, Action: baseline.RetiredSkillKeep, Reason: baseline.RetiredSkillEdited, LockHash: beforeLock["skills"].(map[string]any)["design-patterns"].(map[string]any)["computedHash"].(string)},
				{Skill: "tdd", Paths: []string{".agents/skills/tdd", ".claude/skills/tdd"}, LockEntry: true, Action: baseline.RetiredSkillDelete, LockHash: beforeLock["skills"].(map[string]any)["tdd"].(map[string]any)["computedHash"].(string)},
			}
			survivors := map[string]string{}
			for _, name := range []string{"design-patterns", "context7-cli", "exa-web-search", "local-only"} {
				survivors[name] = baselinePlanTestTree(t, filepath.Join(repo, ".agents", "skills", name))
			}
			preview, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--format="+format)
			message := "guidance matches the current Baseline catalog; the approved update deletes 1 Retired Skill(s)"
			if code != exitUnverified || stderr != "" || preview.PlanDigest == "" {
				t.Fatalf("preview exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			if format == "json" {
				if preview.State != "plan_ready" || preview.Category != "approval" || preview.Message != message || !reflect.DeepEqual(preview.Skills.Retired, wantRetired) || len(preview.Skills.Deleted) != 0 {
					t.Fatalf("preview = %+v", preview)
				}
			} else {
				for _, line := range []string{
					"Baseline update: plan ready\nCategory: approval\nResult: " + message + "\n",
					"- retired tdd: no longer required by the Baseline; the approved update deletes .agents/skills/tdd, .claude/skills/tdd and its skills-lock.json entry\n",
					"- retired design-patterns: kept, its copy differs from its skills-lock.json computedHash; delete .agents/skills/design-patterns, .claude/skills/design-patterns and its skills-lock.json entry by hand if the repository no longer uses it\n",
				} {
					if !strings.Contains(out, line) {
						t.Fatalf("preview omits %q:\n%s", line, out)
					}
				}
			}
			skipped, skipOut, skipErr, skipCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--no-skills", "--format=json")
			if skipCode != exitOK || skipErr != "" || skipped.PlanDigest == preview.PlanDigest || strings.Contains(skipOut, preview.PlanDigest) {
				t.Fatalf("no-skills digest exit=%d stdout=%s stderr=%s", skipCode, skipOut, skipErr)
			}
			if before != baselinePlanTestTree(t, repo) {
				t.Fatal("preview changed repository bytes")
			}
			applied, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--confirm-plan", preview.PlanDigest, "--format="+format)
			if code != exitOK || stderr != "" {
				t.Fatalf("apply exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			if format == "text" {
				want := "Skills deleted: 1\n- deleted tdd: .agents/skills/tdd, .claude/skills/tdd and its skills-lock.json entry\nSkills retired: 1\n"
				if !strings.Contains(out, want) {
					t.Fatalf("apply text = %s", out)
				}
			} else {
				if applied.State != "verified" || applied.ApprovedPlanDigest != preview.PlanDigest || !reflect.DeepEqual(applied.Skills.Deleted, wantRetired[1:]) || !reflect.DeepEqual(applied.Skills.Retired, wantRetired[:1]) {
					t.Fatalf("apply JSON = %s", out)
				}
			}
			for _, root := range []string{".agents", ".claude"} {
				if _, err := os.Lstat(filepath.Join(repo, root, "skills", "tdd")); !os.IsNotExist(err) {
					t.Fatalf("tdd %s remains: %v", root, err)
				}
			}
			var afterLock map[string]any
			if err := json.Unmarshal(mustReadBytes(t, lockPath), &afterLock); err != nil {
				t.Fatal(err)
			}
			delete(beforeLock["skills"].(map[string]any), "tdd")
			if !reflect.DeepEqual(afterLock, beforeLock) {
				t.Fatalf("lock changed beyond tdd: %+v", afterLock)
			}
			for name, tree := range survivors {
				if tree != baselinePlanTestTree(t, filepath.Join(repo, ".agents", "skills", name)) {
					t.Fatalf("survivor %s changed", name)
				}
				target, err := os.Readlink(filepath.Join(repo, ".claude", "skills", name))
				if err != nil || target != "../../.agents/skills/"+name {
					t.Fatalf("survivor link %s = %q, %v", name, target, err)
				}
			}
		})
	}
	t.Run("yes", func(t *testing.T) {
		repo := deletableRetiredSkillsRepository(t)
		result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--yes", "--format=json")
		if code != exitOK || stderr != "" || len(result.Skills.Deleted) != 1 || result.Skills.Deleted[0].Skill != "tdd" {
			t.Fatalf("--yes exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
		if _, err := os.Lstat(filepath.Join(repo, ".agents", "skills", "tdd")); !os.IsNotExist(err) {
			t.Fatalf("--yes left tdd: %v", err)
		}
	})
	t.Run("reinspects after the skills stage", func(t *testing.T) {
		repo := deletableRetiredSkillsRepository(t)
		lockBefore := mustReadBytes(t, filepath.Join(repo, "skills-lock.json"))
		stage := func(context.Context, baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
			mustWrite(t, filepath.Join(repo, ".agents", "skills", "tdd", "SKILL.md"), "# edited during skills stage\n")
			return successfulBaselineUpdateSkillsResult(), nil
		}
		result, out, stderr, code := runBaselineUpdateTestCommandWithSkillsStage(t, t.Context(), stage, "baseline", "update", "--repo", repo, "--yes", "--format=json")
		if code != exitRunFailed || result.Category != "skills" || result.NextAction != "delete the remaining Retired Skills by hand or rerun roundfix baseline update" || !strings.Contains(stderr, `retired skill "tdd" changed after the plan`) || len(result.Skills.Deleted) != 0 {
			t.Fatalf("post-stage refusal exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
		if !bytes.Equal(lockBefore, mustReadBytes(t, filepath.Join(repo, "skills-lock.json"))) {
			t.Fatal("post-stage refusal changed lock")
		}
		target, err := os.Readlink(filepath.Join(repo, ".claude", "skills", "tdd"))
		if err != nil || target != "../../.agents/skills/tdd" {
			t.Fatalf("post-stage refusal removed link: %q, %v", target, err)
		}
	})
}

func TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan(t *testing.T) {
	t.Parallel()
	repo := deletableRetiredSkillsRepository(t)
	preview, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--format=json")
	if code != exitUnverified || stderr != "" {
		t.Fatalf("preview exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	mustWrite(t, filepath.Join(repo, ".agents", "skills", "tdd", "SKILL.md"), "# edited after preview\n")
	before := baselinePlanTestTree(t, repo)
	result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--confirm-plan", preview.PlanDigest, "--format=json")
	if code != exitUnverified || result.State != "action_required" || result.Category != "approval" || !strings.Contains(stderr, "does not match supplied plan") || strings.Contains(out, `"deleted"`) {
		t.Fatalf("stale confirmation exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	if before != baselinePlanTestTree(t, repo) {
		t.Fatal("stale confirmation changed repository bytes")
	}
}
