// Suite: replaced Inngest skill migration through Baseline update.
// Invariant: the seven replacements are required; approval deletes only an unedited old copy.
// Boundary IN: temporary adopted repositories, real preview/apply, lock hashes and skill links.
// Boundary OUT: network acquisition and the operator's installed skills.

package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/skills"
)

func TestTheBackendModuleRequiresTheSevenInngestSkills(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	writeBaselinePlanTestFile(t, repo, "docs/agents/setup-context.json", `{"modules":["typescript","backend"]}`)
	required, manifest, err := resolveExternalSkillRequirement(repo)
	if err != nil || !manifest {
		t.Fatalf("required skills = %v, manifest = %t, err = %v", required, manifest, err)
	}
	if slices.Contains(required, "inngest") {
		t.Error("backend still requires inngest")
	}
	for _, name := range []string{"inngest-cli", "inngest-durable-functions", "inngest-events", "inngest-flow-control", "inngest-middleware", "inngest-setup", "inngest-steps"} {
		if !slices.Contains(required, name) {
			t.Errorf("backend does not require %s", name)
		}
	}
}

func TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unedited", "edited"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo := newBaselineUpdateRepository(t)
			installRequiredExternalSkills(t, repo)
			lockPath := filepath.Join(repo, "skills-lock.json")
			var lock map[string]any
			if err := json.Unmarshal(mustReadBytes(t, lockPath), &lock); err != nil {
				t.Fatal(err)
			}
			writeBaselinePlanTestFile(t, repo, ".agents/skills/inngest/SKILL.md", "# Inngest\n")
			hash, err := skills.SkillFolderHash(t.Context(), filepath.Join(repo, ".agents/skills/inngest"))
			if err != nil {
				t.Fatal(err)
			}
			entries := lock["skills"].(map[string]any)
			entries["inngest"] = map[string]any{"computedHash": hash}
			data, err := json.Marshal(lock)
			if err != nil {
				t.Fatal(err)
			}
			mustWrite(t, lockPath, string(data))
			link := filepath.Join(repo, ".claude/skills/inngest")
			if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("../../.agents/skills/inngest", link); err != nil {
				t.Fatal(err)
			}
			commitBaselinePlanTestRepository(t, repo)
			want := []baseline.InstalledRetiredSkill{{
				Skill: "inngest", Paths: []string{".agents/skills/inngest", ".claude/skills/inngest"},
				LockEntry: true, LockHash: hash, Action: baseline.RetiredSkillDelete,
			}}
			if name == "edited" {
				writeBaselinePlanTestFile(t, repo, ".agents/skills/inngest/SKILL.md", "# Repository-specific Inngest guidance\n")
				commitBaselinePlanTestRepository(t, repo)
				want[0].Action, want[0].Reason = baseline.RetiredSkillKeep, baseline.RetiredSkillEdited
			}
			before := baselinePlanTestTree(t, repo)
			preview, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--format=json")
			if !reflect.DeepEqual(preview.Skills.Retired, want) || len(preview.Skills.Deleted) != 0 || stderr != "" {
				t.Fatalf("preview exit=%d stdout=%s stderr=%s; want retired %+v", code, out, stderr, want)
			}
			if name == "edited" {
				if code != exitOK || preview.State != "current" {
					t.Fatalf("edited preview exit=%d stdout=%s", code, out)
				}
				_, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--yes", "--format=text")
				const kept = "- retired inngest: kept, its copy differs from its skills-lock.json computedHash; delete .agents/skills/inngest, .claude/skills/inngest and its skills-lock.json entry by hand if the repository no longer uses it\n"
				if code != exitOK || stderr != "" || !strings.Contains(out, kept) {
					t.Fatalf("edited apply exit=%d stdout=%s stderr=%s", code, out, stderr)
				}
				if before != baselinePlanTestTree(t, repo) {
					t.Fatal("keeping edited inngest changed repository bytes")
				}
				return
			}
			if code != exitUnverified || preview.State != "plan_ready" || preview.PlanDigest == "" {
				t.Fatalf("unedited preview exit=%d stdout=%s", code, out)
			}
			_, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--format=text")
			const deletion = "- retired inngest: no longer required by the Baseline; the approved update deletes .agents/skills/inngest, .claude/skills/inngest and its skills-lock.json entry\n"
			if code != exitUnverified || stderr != "" || !strings.Contains(out, deletion) {
				t.Fatalf("text preview exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			if before != baselinePlanTestTree(t, repo) {
				t.Fatal("preview changed repository bytes")
			}
			applied, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--confirm-plan", preview.PlanDigest, "--format=json")
			if code != exitOK || stderr != "" || applied.State != "verified" || applied.ApprovedPlanDigest != preview.PlanDigest || !reflect.DeepEqual(applied.Skills.Deleted, want) || len(applied.Skills.Retired) != 0 {
				t.Fatalf("approved apply exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			for _, path := range want[0].Paths {
				if _, err := os.Lstat(filepath.Join(repo, path)); !os.IsNotExist(err) {
					t.Fatalf("deleted path %s remains: %v", path, err)
				}
			}
			var afterLock map[string]any
			if err := json.Unmarshal(mustReadBytes(t, lockPath), &afterLock); err != nil {
				t.Fatal(err)
			}
			delete(entries, "inngest")
			if !reflect.DeepEqual(afterLock, lock) {
				t.Fatalf("lock changed beyond inngest: %+v", afterLock)
			}
		})
	}
}
