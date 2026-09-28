// Suite: skills-lock planning read ordering.
// Invariant: a skill-lock plan and its transaction preimage describe the same post-acquisition bytes.
// Boundary IN: real local Git acquisition, skills-lock parsing, plan digests, and transaction documents.
// Boundary OUT: CLI help and user documentation, owned by internal/cli and docs/user-guide.
package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func TestReconcileSkillsLockPlansFromTheLockPresentAfterTheFetch(t *testing.T) {
	repo, source, revision, dependencies := newSkillsReconcileFixture(t, map[string]string{
		"skills/present/SKILL.md": "# present\n",
	})
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "obsolete-one": {
      "source": "example/skills",
      "skillPath": "skills/obsolete-one/SKILL.md"
    }
  }
}
`)
	rewritten := []byte(`{
  "version": 1,
  "skills": {
    "obsolete-one": {
      "source": "example/skills",
      "skillPath": "skills/obsolete-one/SKILL.md"
    },
    "obsolete-two": {
      "source": "example/skills",
      "skillPath": "skills/obsolete-two/SKILL.md"
    }
  }
}
`)
	installFetchTimeSkillsLockRewrite(t, repo, rewritten)

	payload, err := reconcileSkillsLock(context.Background(), SkillsReconcileRequest{
		Repository:       repo,
		ProfileID:        dependencies.profile.ID,
		SourceRepository: "example/skills",
		Commit:           revision,
		SourceDir:        source,
	}, dependencies)
	assertReconcileErrorCode(t, err, "plan.confirmation.required")
	if len(payload.PlannedChanges) != 2 {
		t.Fatalf("planned changes = %+v, want both post-fetch obsolete entries", payload.PlannedChanges)
	}
	for _, skill := range []string{"obsolete-one", "obsolete-two"} {
		if !hasPlannedLockChange(payload.PlannedChanges, skill) {
			t.Fatalf("planned changes omit %s: %+v", skill, payload.PlannedChanges)
		}
	}
}

func TestReconcileSkillsLockConfirmedApplyKeepsAFetchTimeRewrite(t *testing.T) {
	repo, source, revision, dependencies := newSkillsReconcileFixture(t, map[string]string{
		"skills/present/SKILL.md": "# present\n",
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
	rewritten := []byte(`{
  "version": 1,
  "skills": {
    "obsolete": {
      "source": "example/skills",
      "skillPath": "skills/obsolete/SKILL.md"
    },
    "concurrent": {
      "source": "another/source",
      "skillPath": "skills/concurrent/SKILL.md"
    }
  }
}
`)
	installFetchTimeSkillsLockRewrite(t, repo, rewritten)

	_, err = reconcileSkillsLock(context.Background(), request, dependencies)
	assertReconcileErrorCode(t, err, "plan.confirmation.stale")
	assertFileBytes(t, filepath.Join(repo, skillsLockPath), rewritten)
}

func TestSkillsRestorePlansFromTheLockPresentAfterTheFetch(t *testing.T) {
	repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
		"agentic-cli-design": {"SKILL.md": "# restored\n"},
	})
	writeInspectionFile(t, repo, skillsLockPath, "{\n  \"version\": 1,\n  \"skills\": {}\n}\n")
	rewritten := []byte(`{
  "version": 1,
  "skills": {
    "unrelated": {
      "source": "another/source",
      "skillPath": "skills/unrelated/SKILL.md"
    }
  }
}
`)
	installFetchTimeSkillsLockRewrite(t, repo, rewritten)

	plan, err := buildSkillsRestorePlan(context.Background(), SkillsRestoreRequest{
		Repository: repo,
		ProfileID:  dependencies.profile.ID,
		Skills:     []string{"agentic-cli-design"},
		SourceDir:  source,
	}, dependencies)
	if err != nil {
		t.Fatalf("build restoration plan: %v", err)
	}
	postimage := transactionPostimageByPath(t, plan.document, skillsLockPath)
	var lock map[string]any
	if err := json.Unmarshal(postimage.Content, &lock); err != nil {
		t.Fatalf("decode planned lock postimage: %v", err)
	}
	if _, ok := lock["skills"].(map[string]any)["unrelated"]; !ok {
		t.Fatalf("planned lock postimage dropped post-fetch entry: %s", postimage.Content)
	}
}

func TestSkillsRestoreConfirmedApplyKeepsAFetchTimeRewrite(t *testing.T) {
	repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
		"agentic-cli-design": {"SKILL.md": "# restored\n"},
	})
	writeInspectionFile(t, repo, skillsLockPath, "{\n  \"version\": 1,\n  \"skills\": {}\n}\n")
	request := SkillsRestoreRequest{
		Repository: repo,
		ProfileID:  dependencies.profile.ID,
		Skills:     []string{"agentic-cli-design"},
		SourceDir:  source,
	}
	preview, err := restoreSkills(context.Background(), request, dependencies)
	assertSkillsRestoreErrorCode(t, err, "plan.confirmation.required")
	request.Confirmation = *preview.PlanDigest
	rewritten := []byte(`{
  "version": 1,
  "skills": {
    "concurrent": {
      "source": "another/source",
      "skillPath": "skills/concurrent/SKILL.md"
    }
  }
}
`)
	installFetchTimeSkillsLockRewrite(t, repo, rewritten)

	_, err = restoreSkills(context.Background(), request, dependencies)
	assertSkillsRestoreErrorCode(t, err, "plan.confirmation.stale")
	assertFileBytes(t, filepath.Join(repo, skillsLockPath), rewritten)
}

func TestReconcileTransactionDocumentRefusesAChangedLock(t *testing.T) {
	repo, identity := newTransactionDocumentRepository(t, []byte("disk lock\n"))
	document, err := buildSkillsReconcileTransactionDocument(
		repo, identity, "digest", []byte("planned lock\n"), []byte("postimage\n"),
	)
	assertLockChangedDuringPlan(t, document, err)
}

func TestRestoreTransactionDocumentRefusesAChangedLock(t *testing.T) {
	repo, identity := newTransactionDocumentRepository(t, []byte("disk lock\n"))
	document, err := buildRestoreTransactionDocument(
		repo,
		identity,
		"digest",
		lockEditingRestoreSkills(),
		map[string][]restoreFile{},
		[]byte("planned lock\n"),
		[]byte("postimage\n"),
	)
	assertLockChangedDuringPlan(t, document, err)
}

func TestReconcileTransactionDocumentAcceptsTheLockItPlanned(t *testing.T) {
	planned := []byte("planned lock\n")
	repo, identity := newTransactionDocumentRepository(t, planned)
	document, err := buildSkillsReconcileTransactionDocument(
		repo, identity, "digest", planned, []byte("postimage\n"),
	)
	if err != nil {
		t.Fatalf("build reconciliation transaction document: %v", err)
	}
	assertPlannedLockPreimage(t, document, planned)
}

func TestRestoreTransactionDocumentAcceptsTheLockItPlanned(t *testing.T) {
	planned := []byte("planned lock\n")
	repo, identity := newTransactionDocumentRepository(t, planned)
	document, err := buildRestoreTransactionDocument(
		repo,
		identity,
		"digest",
		lockEditingRestoreSkills(),
		map[string][]restoreFile{},
		planned,
		[]byte("postimage\n"),
	)
	if err != nil {
		t.Fatalf("build restoration transaction document: %v", err)
	}
	assertPlannedLockPreimage(t, document, planned)
}

func installFetchTimeSkillsLockRewrite(t *testing.T, repo string, rewritten []byte) {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("locate real git: %v", err)
	}
	bin := t.TempDir()
	rewritePath := filepath.Join(bin, "rewritten-skills-lock.json")
	if err := os.WriteFile(rewritePath, rewritten, 0o600); err != nil {
		t.Fatalf("write replacement skills lock: %v", err)
	}
	marker := filepath.Join(bin, "rewritten")
	script := fmt.Sprintf(`#!/bin/sh
if [ "$1" = "init" ] && [ "$2" = "--bare" ] && [ ! -e %s ]; then
  touch %s
  cp %s %s
fi
exec %s "$@"
`, shellTestArgument(marker), shellTestArgument(marker), shellTestArgument(rewritePath),
		shellTestArgument(filepath.Join(repo, skillsLockPath)), shellTestArgument(realGit))
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0o755); err != nil {
		t.Fatalf("write git wrapper: %v", err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func shellTestArgument(value string) string {
	return strconv.Quote(value)
}

func hasPlannedLockChange(changes []RestorePlannedChange, skill string) bool {
	for _, change := range changes {
		if change.Action == "remove-lock-entry" && change.Skill == skill {
			return true
		}
	}
	return false
}

func assertFileBytes(t *testing.T, filename string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read %s: %v", filename, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s bytes changed:\ngot:  %q\nwant: %q", filename, got, want)
	}
}

func assertSkillsRestoreErrorCode(t *testing.T, err error, want string) {
	t.Helper()
	var restoreErr *SkillsRestoreError
	if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != want {
		t.Fatalf("error = %v, want %s", err, want)
	}
}

func transactionPostimageByPath(t *testing.T, document PlanDocument, want string) Postimage {
	t.Helper()
	for _, postimage := range document.Postimages {
		if postimage.Path == want {
			return postimage
		}
	}
	t.Fatalf("transaction document has no postimage for %s: %+v", want, document.Postimages)
	return Postimage{}
}

func newTransactionDocumentRepository(t *testing.T, lock []byte) (string, RepositoryIdentity) {
	t.Helper()
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "repository\n")
	commitInspectionRepository(t, repo, "seed repository")
	if err := os.WriteFile(filepath.Join(repo, skillsLockPath), lock, 0o600); err != nil {
		t.Fatalf("write transaction lock fixture: %v", err)
	}
	root, identity, err := inspectRepositoryIdentity(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect transaction repository: %v", err)
	}
	return root, identity
}

func lockEditingRestoreSkills() []RestoreSkill {
	return []RestoreSkill{{
		Skill:      "example",
		TargetPath: ".agents/skills/example",
		LockEdit: &RestoreLockEdit{
			Action: "update-lock-entry",
			Path:   skillsLockPath,
			Skill:  "example",
		},
	}}
}

func assertLockChangedDuringPlan(t *testing.T, document PlanDocument, err error) {
	t.Helper()
	if !reflect.DeepEqual(document, PlanDocument{}) {
		t.Fatalf("changed lock produced a transaction document: %+v", document)
	}
	var restoreErr *SkillsRestoreError
	if !errors.As(err, &restoreErr) {
		t.Fatalf("error = %v, want SkillsRestoreError", err)
	}
	if restoreErr.Category != SkillsRestoreAction ||
		restoreErr.Finding.Code != "lock.changed-during-plan" ||
		restoreErr.Finding.Message != "skills-lock.json changed during planning; the plan no longer describes it." ||
		restoreErr.Finding.Action != "Rerun the preview and confirm its new Plan Digest." {
		t.Fatalf("changed-lock refusal = %+v", restoreErr)
	}
}

func assertPlannedLockPreimage(t *testing.T, document PlanDocument, planned []byte) {
	t.Helper()
	for _, preimage := range document.Preimages {
		if preimage.Path != skillsLockPath {
			continue
		}
		if preimage.Kind != PreimageRegular ||
			preimage.ContentIdentity != transactionContentIdentity(planned) {
			t.Fatalf("lock preimage = %+v, want planned identity", preimage)
		}
		return
	}
	t.Fatalf("transaction document has no %s preimage: %+v", skillsLockPath, document.Preimages)
}
