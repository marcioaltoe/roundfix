// Suite: Retired Skill reporting in Baseline update.
// Invariant: retired copies are advisory output and remain untouched.
// Boundary IN: adopted repository, real preview/apply, output, and exit codes.
// Boundary OUT: network acquisition (the existing skills-stage seam is injected).

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

const retiredSkillsTranscript = "Skills retired: 2\n" +
	"- retired council: no longer required by the Baseline; delete .agents/skills/council\n" +
	"- retired the-fool: no longer required by the Baseline; delete .agents/skills/the-fool and its skills-lock.json entry\n"

func TestBaselineUpdateRetiredSkillDeletionTargets(t *testing.T) {
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
			result.Skills.Retired = []baseline.InstalledRetiredSkill{{Skill: "council", Paths: test.paths, LockEntry: test.lock}}
			var out bytes.Buffer
			if err := writeBaselineUpdateResult(result, false, &out); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), "Skills retired: 1\n- retired council: no longer required by the Baseline; delete "+test.want+"\n") {
				t.Fatalf("retired deletion targets = %s", out.String())
			}
		})
	}
}

func TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState(t *testing.T) {
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
	for _, approval := range []string{"", "--yes"} {
		for _, format := range []string{"text", "json"} {
			repo := adoptedRepositoryWithRetiredSkills(t)
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
			if code != exitOK || errOut != "" || len(result.Skills.Retired) != 0 || strings.Contains(out, "Skills retired") || strings.Contains(out, `"retired"`) {
				t.Fatalf("skipped retired report exit=%d stdout=%s stderr=%s", code, out, errOut)
			}
			if before != baselinePlanTestTree(t, repo) {
				t.Fatal("--no-skills changed current repository bytes")
			}
		}
	}
}

func TestBaselineUpdateAppliesAndListsRetiredSkills(t *testing.T) {
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
	return repo
}

func expectedInstalledRetiredSkills() []baseline.InstalledRetiredSkill {
	return []baseline.InstalledRetiredSkill{
		{Skill: "council", Paths: []string{".agents/skills/council"}},
		{Skill: "the-fool", Paths: []string{".agents/skills/the-fool"}, LockEntry: true},
	}
}
