// Suite: immutable Repository Skill Set restoration
// Invariant: only an exact current Plan Digest may atomically install verified external skill bytes and lock evidence.
// Boundary IN: embedded-style profile contracts, real offline Git acquisition, safe repository targets, and the recoverable file transaction.
// Boundary OUT: public GitHub availability and CLI text rendering, which internal/cli/baseline_skills_restore_test.go owns.
package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSkillsRestoreOfflinePreviewApplyAndIdempotence(t *testing.T) {
	repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
		"agentic-cli-design": {
			"SKILL.md":            "# restored\n",
			"references/guide.md": "guide\n",
		},
	})
	writeInspectionFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# drifted\n")
	writeInspectionFile(t, repo, ".agents/skills/agentic-cli-design/references/removed.md", "remove\n")
	writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "unrelated": {
      "computedHash": "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
    }
  }
}
`)
	before := snapshotVisibleTree(t, repo)
	request := SkillsRestoreRequest{
		Repository: repo,
		ProfileID:  dependencies.profile.ID,
		Skills:     []string{"agentic-cli-design"},
		SourceDir:  source,
	}

	preview, err := restoreSkills(context.Background(), request, dependencies)
	var restoreErr *SkillsRestoreError
	if !errors.As(err, &restoreErr) ||
		restoreErr.Category != SkillsRestoreAction ||
		restoreErr.Finding.Code != "plan.confirmation.required" {
		t.Fatalf("preview error = %v, want confirmation required", err)
	}
	if preview.PlanDigest == nil || !lowercaseSHA256.MatchString(*preview.PlanDigest) {
		t.Fatalf("preview Plan Digest = %v", preview.PlanDigest)
	}
	if len(preview.PlannedChanges) != 4 {
		t.Fatalf("preview planned changes = %+v", preview.PlannedChanges)
	}
	assertVisibleTree(t, repo, before)

	request.Confirmation = *preview.PlanDigest
	applied, err := restoreSkills(context.Background(), request, dependencies)
	if err != nil {
		t.Fatalf("apply offline restoration: %v", err)
	}
	if !applied.OK || !applied.Applied || applied.Finding == nil ||
		applied.Finding.Code != "restore.completed" {
		t.Fatalf("applied payload = %+v", applied)
	}
	assertRestoreFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# restored\n")
	assertRestoreFile(t, repo, ".agents/skills/agentic-cli-design/references/guide.md", "guide\n")
	if _, err := os.Lstat(filepath.Join(
		repo,
		".agents",
		"skills",
		"agentic-cli-design",
		"references",
		"removed.md",
	)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("removed source file still exists: %v", err)
	}
	lock := readSkillsRestoreLock(t, repo)
	entry := lock["skills"].(map[string]any)["agentic-cli-design"].(map[string]any)
	if entry["source"] != "example/skills" ||
		entry["ref"] != dependencies.profile.Skills["agentic-cli-design"].Source.Ref ||
		entry["computedHash"] == "" {
		t.Fatalf("portable lock entry = %+v", entry)
	}
	if _, ok := lock["skills"].(map[string]any)["unrelated"]; !ok {
		t.Fatalf("unrelated lock entry was not preserved: %+v", lock)
	}
	contract := dependencies.profile.Skills["agentic-cli-design"]
	wantLock := fmt.Sprintf(`{
  "version": 1,
  "skills": {
    "unrelated": {
      "computedHash": "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
    },
    "agentic-cli-design": {
      "source": "example/skills",
      "ref": %q,
      "sourceType": "github",
      "skillPath": "skills/agentic-cli-design/SKILL.md",
      "computedHash": %q
    }
  }
}
`, contract.Source.Ref, externalSkillsLockDigest([]restoreFile{
		{Path: "SKILL.md", Content: []byte("# restored\n")},
		{Path: "references/guide.md", Content: []byte("guide\n")},
	}))
	lockBytes, err := os.ReadFile(filepath.Join(repo, skillsLockPath))
	if err != nil {
		t.Fatal(err)
	}
	if string(lockBytes) != wantLock {
		t.Fatalf("portable lock bytes differ:\ngot:\n%s\nwant:\n%s", lockBytes, wantLock)
	}
	lockInfo, err := os.Stat(filepath.Join(repo, skillsLockPath))
	if err != nil {
		t.Fatal(err)
	}
	if lockInfo.Mode().Perm() != 0o600 {
		t.Fatalf("skills-lock.json mode = %#o, want 0600", lockInfo.Mode().Perm())
	}

	request.Confirmation = ""
	second, err := restoreSkills(context.Background(), request, dependencies)
	if err != nil {
		t.Fatalf("empty restoration: %v", err)
	}
	if !second.OK || second.Applied || len(second.PlannedChanges) != 0 {
		t.Fatalf("empty restoration payload = %+v", second)
	}
}

func TestSkillsRestoreProvenanceAndPreMutationRefusals(t *testing.T) {
	t.Run("groups exact provenance", func(t *testing.T) {
		repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
			"agentic-cli-design": {"SKILL.md": "# agentic\n"},
			"domain-modeling": {
				"SKILL.md":            "# domain\n",
				"references/guide.md": "guide\n",
			},
		})
		payload, err := restoreSkills(context.Background(), SkillsRestoreRequest{
			Repository: repo,
			ProfileID:  dependencies.profile.ID,
			Skills:     []string{"domain-modeling", "agentic-cli-design"},
			SourceDir:  source,
		}, dependencies)
		if err == nil {
			t.Fatal("preview unexpectedly applied without confirmation")
		}
		if len(payload.Acquisitions) != 1 {
			t.Fatalf("acquisitions = %+v, want one exact provenance group", payload.Acquisitions)
		}
		if got := []string{payload.Skills[0].Skill, payload.Skills[1].Skill}; !reflect.DeepEqual(
			got,
			[]string{"agentic-cli-design", "domain-modeling"},
		) {
			t.Fatalf("skill order = %v", got)
		}
	})

	tests := []struct {
		name    string
		prepare func(t *testing.T, repo string, dependencies *restoreDependencies)
		code    string
	}{
		{
			name: "provenance commit mismatch",
			prepare: func(t *testing.T, _ string, dependencies *restoreDependencies) {
				contract := dependencies.profile.Skills["agentic-cli-design"]
				contract.Source.Ref = "HEAD"
				dependencies.profile.Skills["agentic-cli-design"] = contract
			},
			code: "source.commit-mismatch",
		},
		{
			name: "source digest mismatch",
			prepare: func(t *testing.T, _ string, dependencies *restoreDependencies) {
				contract := dependencies.profile.Skills["agentic-cli-design"]
				contract.TreeDigest = strings.Repeat("0", 64)
				dependencies.profile.Skills["agentic-cli-design"] = contract
			},
			code: "source.digest-mismatch",
		},
		{
			name: "lock adapter mismatch",
			prepare: func(t *testing.T, _ string, dependencies *restoreDependencies) {
				dependencies.adapterFixture = []byte(`{
  "schemaVersion": "setup-context-driven/external-lock-hash-compatibility-v1",
  "version": 1,
  "files": [{"path": "SKILL.md", "content": "# fixture\n"}],
  "expectedSha256": "0000000000000000000000000000000000000000000000000000000000000000"
}`)
			},
			code: "lock.adapter-incompatible",
		},
		{
			name: "unsafe target",
			prepare: func(t *testing.T, repo string, _ *restoreDependencies) {
				parent := filepath.Join(repo, ".agents", "skills")
				if err := os.MkdirAll(parent, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(
					t.TempDir(),
					filepath.Join(parent, "agentic-cli-design"),
				); err != nil {
					t.Fatal(err)
				}
			},
			code: "target.unsafe-tree",
		},
		{
			name: "malformed lock",
			prepare: func(t *testing.T, repo string, _ *restoreDependencies) {
				writeInspectionFile(t, repo, skillsLockPath, "{not json")
			},
			code: "lock.invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
				"agentic-cli-design": {"SKILL.md": "# restored\n"},
			})
			test.prepare(t, repo, &dependencies)
			before := snapshotVisibleTree(t, repo)
			_, err := restoreSkills(context.Background(), SkillsRestoreRequest{
				Repository: repo,
				ProfileID:  dependencies.profile.ID,
				Skills:     []string{"agentic-cli-design"},
				SourceDir:  source,
			}, dependencies)
			var restoreErr *SkillsRestoreError
			if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != test.code {
				t.Fatalf("error = %v, want %s", err, test.code)
			}
			assertVisibleTree(t, repo, before)
		})
	}
}

func TestSkillsRestoreStalePlanDoesNotMutate(t *testing.T) {
	repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
		"agentic-cli-design": {"SKILL.md": "# restored\n"},
	})
	writeInspectionFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# first drift\n")
	request := SkillsRestoreRequest{
		Repository: repo,
		ProfileID:  dependencies.profile.ID,
		Skills:     []string{"agentic-cli-design"},
		SourceDir:  source,
	}
	preview, err := restoreSkills(context.Background(), request, dependencies)
	if err == nil || preview.PlanDigest == nil {
		t.Fatalf("preview = %+v error=%v", preview, err)
	}
	writeInspectionFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# changed after preview\n")
	before := snapshotVisibleTree(t, repo)
	request.Confirmation = *preview.PlanDigest
	payload, err := restoreSkills(context.Background(), request, dependencies)
	var restoreErr *SkillsRestoreError
	if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != "plan.confirmation.stale" {
		t.Fatalf("stale error = %v payload=%+v", err, payload)
	}
	assertVisibleTree(t, repo, before)
}

func TestSkillsRestoreRollbackRestoresSkillAndLockPreimage(t *testing.T) {
	tests := []struct {
		name  string
		phase transactionPhase
		path  string
	}{
		{
			name:  "partial skill write",
			phase: transactionPhaseReplacing,
			path:  ".agents/skills/agentic-cli-design/references/guide.md",
		},
		{
			name:  "lock update",
			phase: transactionPhaseReplacing,
			path:  skillsLockPath,
		},
		{
			name:  "postwrite verification",
			phase: transactionPhaseVerifying,
			path:  skillsLockPath,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
				"agentic-cli-design": {
					"SKILL.md":            "# restored\n",
					"references/guide.md": "guide\n",
				},
			})
			writeInspectionFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# drifted\n")
			writeInspectionFile(t, repo, skillsLockPath, `{
  "version": 1,
  "skills": {
    "unrelated": {
      "computedHash": "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
    }
  }
}
`)
			request := SkillsRestoreRequest{
				Repository: repo,
				ProfileID:  dependencies.profile.ID,
				Skills:     []string{"agentic-cli-design"},
				SourceDir:  source,
			}
			preview, err := restoreSkills(context.Background(), request, dependencies)
			if err == nil || preview.PlanDigest == nil {
				t.Fatalf("preview = %+v error=%v", preview, err)
			}
			before := snapshotVisibleTree(t, repo)
			request.Confirmation = *preview.PlanDigest
			dependencies.transactionHook = failTransactionOnce(
				test.phase,
				test.path,
				errors.New("injected restoration failure"),
			)
			payload, err := restoreSkills(context.Background(), request, dependencies)
			var restoreErr *SkillsRestoreError
			if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != "restore.apply-failed" {
				t.Fatalf("rollback error = %v payload=%+v", err, payload)
			}
			assertVisibleTree(t, repo, before)
		})
	}
}

func TestSkillsRestoreCompatibilityMatchesMaintainedPythonShape(t *testing.T) {
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	fixture, ok := catalog.Asset(lockHashCompatibilityPath)
	if !ok {
		t.Fatal("embedded lock compatibility fixture is missing")
	}
	if err := validateLockAdapter(fixture.Data); err != nil {
		t.Fatalf("maintained lock adapter fixture: %v", err)
	}

	repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
		"agentic-cli-design": {
			"SKILL.md":            "# restored\n",
			"references/guide.md": "guide\n",
		},
	})
	payload, err := restoreSkills(context.Background(), SkillsRestoreRequest{
		Repository: repo,
		ProfileID:  dependencies.profile.ID,
		Skills:     []string{"agentic-cli-design"},
		SourceDir:  source,
	}, dependencies)
	if err == nil {
		t.Fatal("preview unexpectedly succeeded without confirmation")
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{
		"schemaVersion",
		"ok",
		"applied",
		"profile",
		"setup",
		"acquisitions",
		"skills",
		"plannedChanges",
		"planDigest",
		"finding",
	} {
		if _, exists := document[field]; !exists {
			t.Fatalf("maintained restore-v1 field %q is missing: %s", field, data)
		}
	}
	if document["schemaVersion"] != SkillsRestoreSchemaVersion {
		t.Fatalf("schemaVersion = %v", document["schemaVersion"])
	}
	changes := document["plannedChanges"].([]any)
	first := changes[0].(map[string]any)
	for _, field := range []string{"action", "path", "skill", "beforeDigest", "afterDigest"} {
		if _, exists := first[field]; !exists {
			t.Fatalf("maintained file-change field %q is missing: %s", field, data)
		}
	}
}

func newSkillsRestoreFixture(
	t *testing.T,
	skills map[string]map[string]string,
) (string, string, restoreDependencies) {
	t.Helper()
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "repository\n")
	commitInspectionRepository(t, repo, "seed repository")

	source := newInspectionRepository(t)
	for skill, files := range skills {
		for relative, content := range files {
			writeInspectionFile(t, source, filepath.ToSlash(filepath.Join("skills", skill, relative)), content)
		}
	}
	commitInspectionRepository(t, source, "fixture skill source")
	revision, err := (ExecGitRunner{}).RunGit(context.Background(), source, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	contracts := make(map[string]restoreSkillContract, len(skills))
	for skill, files := range skills {
		tree := make([]restoreFile, 0, len(files))
		for relative, content := range files {
			tree = append(tree, restoreFile{Path: relative, Content: []byte(content)})
		}
		contracts[skill] = restoreSkillContract{
			Name: skill,
			Source: RestoreSource{
				Provider:   "github",
				Repository: "example/skills",
				Ref:        revision,
				Path:       pathForRestoreFixture(skill),
			},
			TreeDigest: portableRestoreDigest(tree),
		}
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	adapter, ok := catalog.Asset(lockHashCompatibilityPath)
	if !ok {
		t.Fatal("embedded lock compatibility fixture is missing")
	}
	return repo, source, restoreDependencies{
		profile: restoreProfile{
			ID: "rust-cli", Setup: "rust-cli", Skills: contracts,
		},
		adapterFixture: adapter.Data,
	}
}

func pathForRestoreFixture(skill string) string {
	return "skills/" + skill
}

func assertRestoreFile(t *testing.T, repo, relative, want string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatalf("read restored %s: %v", relative, err)
	}
	if string(data) != want {
		t.Fatalf("restored %s = %q, want %q", relative, data, want)
	}
}

func readSkillsRestoreLock(t *testing.T, repo string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repo, skillsLockPath))
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	return document
}
