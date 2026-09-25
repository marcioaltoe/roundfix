// Suite: immutable skills-lock reconciliation.
// Invariant: only source-selected lock entries proven obsolete at one immutable commit are removed.
// Boundary IN: built-in-style Profile requirements, real local bare Git acquisition, ordered lock JSON, and the restore transaction.
// Boundary OUT: CLI parsing and rendering, owned by internal/cli/baseline_skills_reconcile_test.go.
package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries(t *testing.T) {
	t.Parallel()

	repo, source, revision, _ := newSkillsReconcileFixture(t, map[string]string{
		"skills/present/SKILL.md": "# present\n",
	})
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "present": {
      "source": "example/skills",
      "ref": "old",
      "sourceType": "github",
      "skillPath": "skills/present/SKILL.md",
      "computedHash": "present"
    },
    "obsolete": {
      "source": "example/skills",
      "ref": "old",
      "sourceType": "github",
      "skillPath": "skills/obsolete/SKILL.md",
      "computedHash": "obsolete"
    }
  }
}
`)
	request := SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        "go-cli-tui",
		SourceRepository: "example/skills",
		Commit:           revision,
		SourceDir:        source,
	}

	preview, err := ReconcileSkillsLock(context.Background(), request)
	assertReconcileErrorCode(t, err, "plan.confirmation.required")
	if preview.PlanDigest == nil || !lowercaseSHA256.MatchString(*preview.PlanDigest) {
		t.Fatalf("preview Plan Digest = %v", preview.PlanDigest)
	}
	assertReconcileDisposition(t, preview, "present", SkillsLockPresent)
	assertReconcileDisposition(t, preview, "obsolete", SkillsLockObsolete)
	if len(preview.PlannedChanges) != 1 ||
		preview.PlannedChanges[0].Action != "remove-lock-entry" ||
		preview.PlannedChanges[0].Skill != "obsolete" {
		t.Fatalf("planned changes = %+v", preview.PlannedChanges)
	}

	request.Confirmation = *preview.PlanDigest
	applied, err := ReconcileSkillsLock(context.Background(), request)
	if err != nil {
		t.Fatalf("apply lock reconciliation: %v", err)
	}
	if !applied.OK || !applied.Applied {
		t.Fatalf("applied payload = %+v", applied)
	}
	lock := readSkillsRestoreLock(t, repo)
	skills := lock["skills"].(map[string]any)
	if _, exists := skills["obsolete"]; exists {
		t.Fatalf("obsolete entry survived apply: %+v", skills)
	}
	if _, exists := skills["present"]; !exists {
		t.Fatalf("present entry was removed: %+v", skills)
	}
}

func TestReconcileSkillsLockBlocksARequiredRemovedSkill(t *testing.T) {
	t.Parallel()

	repo, source, revision, _ := newSkillsReconcileFixture(t, map[string]string{
		"skills/other/SKILL.md": "# other\n",
	})
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "coding-guidelines": {
      "source": "example/skills",
      "skillPath": "skills/coding-guidelines/SKILL.md",
      "future": {"keep": true}
    },
    "obsolete": {
      "source": "example/skills",
      "skillPath": "skills/obsolete/SKILL.md"
    }
  }
}
`)
	before := snapshotVisibleTree(t, repo)

	payload, err := ReconcileSkillsLock(context.Background(), SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        "go-cli-tui",
		SourceRepository: "example/skills",
		Commit:           revision,
		SourceDir:        source,
	})
	assertReconcileErrorCode(t, err, "reconcile.required-removed")
	assertReconcileDisposition(t, payload, "coding-guidelines", SkillsLockRequiredRemoved)
	assertReconcileDisposition(t, payload, "obsolete", SkillsLockObsolete)
	if payload.PlanDigest != nil || len(payload.PlannedChanges) != 0 {
		t.Fatalf("blocked payload contains an applicable plan: %+v", payload)
	}
	assertVisibleTree(t, repo, before)
}

func TestReconcileSkillsLockUnreachableSourceRemovesNothing(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "repository\n")
	commitInspectionRepository(t, repo, "seed repository")
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "obsolete": {
      "source": "example/skills",
      "skillPath": "skills/obsolete/SKILL.md"
    }
  }
}
`)
	before := snapshotVisibleTree(t, repo)
	dependencies := reconcileDependencies{profile: restoreProfile{
		ID: "fixture", Setup: "fixture", RequiredSkills: map[string]struct{}{},
	}}

	payload, err := reconcileSkillsLock(context.Background(), SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        dependencies.profile.ID,
		SourceRepository: "example/skills",
		Commit:           strings.Repeat("a", 40),
		SourceDir:        t.TempDir(),
	}, dependencies)
	assertReconcileErrorCode(t, err, "source.commit-unavailable")
	if len(payload.Skills) != 0 || len(payload.PlannedChanges) != 0 || payload.PlanDigest != nil {
		t.Fatalf("failed acquisition produced classifications or a plan: %+v", payload)
	}
	assertVisibleTree(t, repo, before)
}

func TestReconcileSkillsLockKeepsMovedAndUnrelatedEntries(t *testing.T) {
	t.Parallel()

	repo, source, revision, dependencies := newSkillsReconcileFixture(t, map[string]string{
		"new/location/moved/SKILL.md": "# moved\n",
	})
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "moved": {
      "source": "example/skills",
      "skillPath": "old/location/moved/SKILL.md"
    },
    "unrelated": {
      "source": "another/source",
      "skillPath": "skills/unrelated/SKILL.md"
    }
  }
}
`)
	before := snapshotVisibleTree(t, repo)

	payload, err := reconcileSkillsLock(context.Background(), SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        dependencies.profile.ID,
		SourceRepository: "example/skills",
		Commit:           revision,
		SourceDir:        source,
	}, dependencies)
	if err != nil {
		t.Fatalf("empty reconciliation: %v", err)
	}
	if !payload.OK || payload.Applied || len(payload.PlannedChanges) != 0 {
		t.Fatalf("empty reconciliation payload = %+v", payload)
	}
	assertReconcileDisposition(t, payload, "moved", SkillsLockMoved)
	assertReconcileDisposition(t, payload, "unrelated", SkillsLockUnrelated)
	assertVisibleTree(t, repo, before)
}

func TestReconcileSkillsLockPreservesInstalledTreesAndUnknownFields(t *testing.T) {
	t.Parallel()

	repo, source, revision, dependencies := newSkillsReconcileFixture(t, map[string]string{
		"skills/kept/SKILL.md": "# kept\n",
	})
	writeInspectionFile(t, repo, ".agents/skills/obsolete/SKILL.md", "# installed and retained\n")
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "futureRoot": {"enabled": true},
  "skills": {
    "first": {
      "source": "another/source",
      "future": [1, 2, 3]
    },
    "obsolete": {
      "source": "example/skills",
      "skillPath": "skills/obsolete/SKILL.md",
      "futureEntry": "removed with its obsolete entry"
    },
    "kept": {
      "source": "example/skills",
      "skillPath": "skills/kept/SKILL.md",
      "futureEntry": {"preserved": true}
    },
    "last": {
      "source": "another/source",
      "future": null
    }
  },
  "futureTail": "keep"
}
`)
	request := SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        dependencies.profile.ID,
		SourceRepository: "example/skills",
		Commit:           revision,
		SourceDir:        source,
	}
	preview, err := reconcileSkillsLock(context.Background(), request, dependencies)
	assertReconcileErrorCode(t, err, "plan.confirmation.required")
	entry := reconcileEntryByName(t, preview, "obsolete")
	if entry.InstalledTree != SkillsInstalledTreeRetained {
		t.Fatalf("obsolete installed tree = %q, want retained", entry.InstalledTree)
	}

	request.Confirmation = *preview.PlanDigest
	if _, err := reconcileSkillsLock(context.Background(), request, dependencies); err != nil {
		t.Fatalf("apply lock reconciliation: %v", err)
	}
	installed, err := os.ReadFile(filepath.Join(
		repo, ".agents", "skills", "obsolete", "SKILL.md",
	))
	if err != nil || string(installed) != "# installed and retained\n" {
		t.Fatalf("installed obsolete tree = %q, %v", installed, err)
	}
	lockBytes, err := os.ReadFile(filepath.Join(repo, skillsLockPath))
	if err != nil {
		t.Fatal(err)
	}
	var lock map[string]any
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatal(err)
	}
	if lock["futureRoot"].(map[string]any)["enabled"] != true || lock["futureTail"] != "keep" {
		t.Fatalf("unknown root fields were not preserved: %s", lockBytes)
	}
	kept := lock["skills"].(map[string]any)["kept"].(map[string]any)
	if kept["futureEntry"].(map[string]any)["preserved"] != true {
		t.Fatalf("unknown kept-entry fields were not preserved: %s", lockBytes)
	}
	wantOrder := []string{`"first"`, `"kept"`, `"last"`}
	position := -1
	for _, field := range wantOrder {
		next := strings.Index(string(lockBytes[position+1:]), field)
		if next < 0 {
			t.Fatalf("lock field %s missing after apply: %s", field, lockBytes)
		}
		position += next + 1
	}
}

func TestReconcileSkillsLockStalePreimageWritesNothing(t *testing.T) {
	t.Parallel()

	repo, source, revision, dependencies := newSkillsReconcileFixture(t, map[string]string{
		"skills/other/SKILL.md": "# other\n",
	})
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "obsolete": {
      "source": "example/skills",
      "skillPath": "skills/obsolete/SKILL.md"
    }
  }
}
`)
	request := SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        dependencies.profile.ID,
		SourceRepository: "example/skills",
		Commit:           revision,
		SourceDir:        source,
	}
	preview, err := reconcileSkillsLock(context.Background(), request, dependencies)
	assertReconcileErrorCode(t, err, "plan.confirmation.required")
	request.Confirmation = *preview.PlanDigest

	concurrent := []byte("{\n  \"version\": 1,\n  \"skills\": {},\n  \"concurrent\": true\n}\n")
	changed := false
	dependencies.transactionHook = func(point transactionFaultPoint) error {
		if point.Phase != transactionPhaseStaged || changed {
			return nil
		}
		changed = true
		return os.WriteFile(filepath.Join(repo, skillsLockPath), concurrent, 0o600)
	}
	payload, err := reconcileSkillsLock(context.Background(), request, dependencies)
	assertReconcileErrorCode(t, err, "plan.confirmation.stale")
	if payload.Applied {
		t.Fatalf("stale reconciliation reported applied: %+v", payload)
	}
	after, readErr := os.ReadFile(filepath.Join(repo, skillsLockPath))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(after) != string(concurrent) {
		t.Fatalf("stale apply changed concurrent preimage:\ngot:  %s\nwant: %s", after, concurrent)
	}
}

func TestReconcileSkillsLockRejectsMutableRevisionBeforeAcquisition(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "repository\n")
	commitInspectionRepository(t, repo, "seed repository")
	dependencies := reconcileDependencies{profile: restoreProfile{
		ID: "fixture", Setup: "fixture", RequiredSkills: map[string]struct{}{},
	}}

	_, err := reconcileSkillsLock(context.Background(), SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        dependencies.profile.ID,
		SourceRepository: "example/skills",
		Commit:           "main",
		SourceDir:        filepath.Join(t.TempDir(), "unreachable"),
	}, dependencies)
	assertReconcileErrorCode(t, err, "reconcile.commit-invalid")
}

func newSkillsReconcileFixture(
	t *testing.T,
	files map[string]string,
) (string, string, string, reconcileDependencies) {
	t.Helper()
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "repository\n")
	commitInspectionRepository(t, repo, "seed repository")

	sourceWorktree := newInspectionRepository(t)
	writeInspectionFile(t, sourceWorktree, "README.md", "source\n")
	for relative, content := range files {
		writeInspectionFile(t, sourceWorktree, relative, content)
	}
	commitInspectionRepository(t, sourceWorktree, "fixture source")
	revision, err := (ExecGitRunner{}).RunGit(
		context.Background(), sourceWorktree, "rev-parse", "HEAD",
	)
	if err != nil {
		t.Fatal(err)
	}
	bareParent := t.TempDir()
	bare := filepath.Join(bareParent, "source.git")
	if _, err := (ExecGitRunner{}).RunGit(
		context.Background(), bareParent, "clone", "--quiet", "--bare", sourceWorktree, bare,
	); err != nil {
		t.Fatal(err)
	}
	return repo, bare, revision, reconcileDependencies{profile: restoreProfile{
		ID: "fixture", Setup: "fixture", RequiredSkills: map[string]struct{}{},
	}}
}

func assertReconcileErrorCode(t *testing.T, err error, want string) {
	t.Helper()
	var restoreErr *SkillsRestoreError
	if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != want {
		t.Fatalf("error = %v, want %s", err, want)
	}
}

func assertReconcileDisposition(
	t *testing.T,
	payload SkillsReconcilePayload,
	name string,
	want SkillsLockDisposition,
) {
	t.Helper()
	entry := reconcileEntryByName(t, payload, name)
	if entry.Disposition != want {
		t.Fatalf("entry %s disposition = %q, want %q", name, entry.Disposition, want)
	}
}

func reconcileEntryByName(
	t *testing.T,
	payload SkillsReconcilePayload,
	name string,
) SkillsReconcileEntry {
	t.Helper()
	for _, entry := range payload.Skills {
		if entry.Skill == name {
			return entry
		}
	}
	t.Fatalf("entry %s missing from payload: %+v", name, payload.Skills)
	return SkillsReconcileEntry{}
}
