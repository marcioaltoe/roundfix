// Suite: Claude Skill Links in offline skill restoration
// Invariant: links are created only in an established link layout and share the tree transaction.
// Boundary IN: offline Git source, preview digest, filesystem entries, transaction faults.
// Boundary OUT: CLI rendering.
package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const restoredClaudeSkill = "agentic-cli-design"
const restoredClaudePath = ".claude/skills/agentic-cli-design"
const restoredClaudeTarget = "../../.agents/skills/agentic-cli-design"

func claudeRestoreFixture(t *testing.T) (string, SkillsRestoreRequest, restoreDependencies) {
	t.Helper()
	repo, source, deps := newSkillsRestoreFixture(t, map[string]map[string]string{restoredClaudeSkill: {"SKILL.md": "# restored\n"}})
	return repo, SkillsRestoreRequest{Repository: repo, ProfileID: deps.profile.ID, Skills: []string{restoredClaudeSkill}, SourceDir: source}, deps
}

func keepClaudeLinks(t *testing.T, repo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(repo, ".claude/skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../../.agents/skills/seed", filepath.Join(repo, ".claude/skills/seed")); err != nil {
		t.Fatal(err)
	}
}

func confirmClaudeRestore(t *testing.T, request SkillsRestoreRequest, deps restoreDependencies) SkillsRestorePayload {
	t.Helper()
	preview, err := restoreSkills(context.Background(), request, deps)
	var action *SkillsRestoreError
	if !errors.As(err, &action) || action.Category != SkillsRestoreAction || preview.PlanDigest == nil {
		t.Fatalf("preview=%+v error=%v", preview, err)
	}
	request.Confirmation = *preview.PlanDigest
	applied, err := restoreSkills(context.Background(), request, deps)
	if err != nil || !applied.Applied {
		t.Fatalf("apply=%+v error=%v", applied, err)
	}
	return preview
}

func assertClaudeLink(t *testing.T, repo string) {
	t.Helper()
	name := filepath.Join(repo, restoredClaudePath)
	info, err := os.Lstat(name)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("link info=%v error=%v", info, err)
	}
	target, err := os.Readlink(name)
	if err != nil || target != restoredClaudeTarget {
		t.Fatalf("target=%q error=%v", target, err)
	}
}

func TestSkillsRestoreLinksARestoredSkillWhereClaudeLinksAreKept(t *testing.T) {
	t.Parallel()
	repo, request, deps := claudeRestoreFixture(t)
	keepClaudeLinks(t, repo)
	preview := confirmClaudeRestore(t, request, deps)
	assertClaudeLink(t, repo)
	assertRestoreFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# restored\n")
	found := false
	for _, change := range preview.PlannedChanges {
		if change.Action != "create-link" {
			continue
		}
		found = true
		data, err := json.Marshal(change)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"action":"create-link","path":".claude/skills/agentic-cli-design","skill":"agentic-cli-design","target":"../../.agents/skills/agentic-cli-design"}`
		if string(data) != want {
			t.Fatalf("change=%s", data)
		}
	}
	if !found {
		t.Fatal("preview lacks create-link")
	}
	second, err := restoreSkills(context.Background(), request, deps)
	if err != nil || len(second.PlannedChanges) != 0 {
		t.Fatalf("second=%+v error=%v", second, err)
	}
}

func TestSkillsRestoreAddsNoLinkWithoutClaudeLinks(t *testing.T) {
	t.Parallel()
	for _, layout := range []string{"absent", "symbolic-link", "copies-only", "claude-symbolic-link", "wrong-target"} {
		t.Run(layout, func(t *testing.T) {
			repo, request, deps := claudeRestoreFixture(t)
			switch layout {
			case "symbolic-link":
				outside := t.TempDir()
				keepClaudeLinks(t, outside)
				if err := os.Mkdir(filepath.Join(repo, ".claude"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(outside, ".claude/skills"), filepath.Join(repo, ".claude/skills")); err != nil {
					t.Fatal(err)
				}
			case "claude-symbolic-link":
				outside := t.TempDir()
				keepClaudeLinks(t, outside)
				if err := os.Symlink(filepath.Join(outside, ".claude"), filepath.Join(repo, ".claude")); err != nil {
					t.Fatal(err)
				}
			case "copies-only":
				writeInspectionFile(t, repo, ".claude/skills/seed/SKILL.md", "copy\n")
			case "wrong-target":
				keepClaudeLinks(t, repo)
				if err := os.Remove(filepath.Join(repo, ".claude/skills/seed")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("../../.agents/skills/other", filepath.Join(repo, ".claude/skills/seed")); err != nil {
					t.Fatal(err)
				}
			}
			preview := confirmClaudeRestore(t, request, deps)
			for _, change := range preview.PlannedChanges {
				if change.Action == "create-link" {
					t.Fatalf("unexpected link: %+v", change)
				}
			}
			if _, err := os.Lstat(filepath.Join(repo, restoredClaudePath)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("unexpected entry: %v", err)
			}
		})
	}
}

func TestSkillsRestoreLeavesAnExistingClaudeEntry(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"file", "directory", "dangling-link"} {
		t.Run(kind, func(t *testing.T) {
			repo, request, deps := claudeRestoreFixture(t)
			keepClaudeLinks(t, repo)
			name := filepath.Join(repo, restoredClaudePath)
			switch kind {
			case "file":
				writeInspectionFile(t, repo, restoredClaudePath, "owned\n")
			case "directory":
				writeInspectionFile(t, repo, restoredClaudePath+"/SKILL.md", "owned\n")
			case "dangling-link":
				if err := os.Symlink("missing", name); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.Lstat(name)
			if err != nil {
				t.Fatal(err)
			}
			preview := confirmClaudeRestore(t, request, deps)
			for _, change := range preview.PlannedChanges {
				if change.Action == "create-link" {
					t.Fatalf("unexpected link: %+v", change)
				}
			}
			after, err := os.Lstat(name)
			if err != nil || !os.SameFile(before, after) {
				t.Fatalf("entry changed: %v", err)
			}
			if kind == "dangling-link" {
				target, err := os.Readlink(name)
				if err != nil || target != "missing" {
					t.Fatalf("target=%q error=%v", target, err)
				}
			}
		})
	}
}

func TestSkillsRestoreRollsBackTheClaudeLinkWithTheTree(t *testing.T) {
	t.Parallel()
	repo, request, deps := claudeRestoreFixture(t)
	keepClaudeLinks(t, repo)
	preview, err := restoreSkills(context.Background(), request, deps)
	if err == nil || preview.PlanDigest == nil {
		t.Fatalf("preview=%+v error=%v", preview, err)
	}
	injected := errors.New("after Claude link replacement")
	hit := false
	deps.transactionHook = func(point transactionFaultPoint) error {
		if point.Phase == transactionPhaseReplaced && point.Path == restoredClaudePath {
			hit = true
			assertClaudeLink(t, repo)
			assertRestoreFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# restored\n")
			return injected
		}
		return nil
	}
	request.Confirmation = *preview.PlanDigest
	_, err = restoreSkills(context.Background(), request, deps)
	if !hit || !errors.Is(err, injected) {
		t.Fatalf("fault hit=%v error=%v", hit, err)
	}
	for _, name := range []string{restoredClaudePath, ".agents/skills/agentic-cli-design", skillsLockPath} {
		if _, err := os.Lstat(filepath.Join(repo, name)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("rollback left %s: %v", name, err)
		}
	}
}

func TestSkillsRestoreLinksAnInstalledSkillThatLacksItsLink(t *testing.T) {
	t.Parallel()
	repo, request, deps := claudeRestoreFixture(t)
	confirmClaudeRestore(t, request, deps)
	keepClaudeLinks(t, repo)
	skillFile := filepath.Join(repo, ".agents/skills/agentic-cli-design/SKILL.md")
	lockFile := filepath.Join(repo, skillsLockPath)
	skillBefore, err := os.Lstat(skillFile)
	if err != nil {
		t.Fatal(err)
	}
	lockBefore, err := os.Lstat(lockFile)
	if err != nil {
		t.Fatal(err)
	}
	preview := confirmClaudeRestore(t, request, deps)
	if len(preview.PlannedChanges) != 1 || preview.PlannedChanges[0].Action != "create-link" {
		t.Fatalf("changes=%+v", preview.PlannedChanges)
	}
	lock, err := os.ReadFile(lockFile)
	if err != nil {
		t.Fatal(err)
	}
	withoutLink, err := restorePlanDigest(deps.profile, preview.Acquisitions, preview.Skills, nil, lock, lock)
	if err != nil {
		t.Fatal(err)
	}
	if *preview.PlanDigest == withoutLink {
		t.Fatal("link is absent from Plan Digest")
	}
	skillAfter, err := os.Lstat(skillFile)
	if err != nil {
		t.Fatal(err)
	}
	lockAfter, err := os.Lstat(lockFile)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(skillBefore, skillAfter) || !os.SameFile(lockBefore, lockAfter) {
		t.Fatal("link-only restore replaced skill bytes or lock")
	}
	assertClaudeLink(t, repo)
}
