// Suite: repository Baseline Profile skill snapshot contracts
// Invariant: built-in contracts stay unchanged; repository profiles use agreeing required external contracts.
// Boundary IN: embedded catalog, repository profile files, installed trees and local Git roots.
// Boundary OUT: network acquisition, CLI rendering and daemon Verification.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func snapshotProfileFixture(t *testing.T) (string, *Catalog, string) {
	t.Helper()
	root := t.TempDir()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	const id = "repository-go"
	if _, err := InitCustomProfile(root, id, "go-cli-tui", catalog); err != nil {
		t.Fatal(err)
	}
	return root, catalog, id
}

func TestRestoreSkillsAcceptsARepositoryProfile(t *testing.T) {
	t.Parallel()
	root, _, id := snapshotProfileFixture(t)
	runApplyGit(t, root, "init")
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	payload, err := RestoreSkills(t.Context(), SkillsRestoreRequest{
		Repository: nested, ProfileID: id, Skills: []string{"context7-cli"}, SourceDir: filepath.Join(root, "missing"),
	})
	if err == nil || payload.Finding == nil || payload.Finding.Code != "restore.source-dir-invalid" {
		t.Fatalf("payload = %+v, error = %v", payload, err)
	}
	if payload.Profile != id || payload.Setup != nil {
		t.Fatalf("profile/setup = %q/%v", payload.Profile, payload.Setup)
	}
}

func TestReconcileSkillsLockAcceptsARepositoryProfile(t *testing.T) {
	t.Parallel()
	root, catalog, id := snapshotProfileFixture(t)
	runApplyGit(t, root, "init")
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	builtin, err := loadRestoreProfile(catalog, "go-cli-tui")
	if err != nil {
		t.Fatal(err)
	}
	source := builtin.Skills["context7-cli"].Source
	payload, err := ReconcileSkillsLock(t.Context(), SkillsReconcileRequest{
		Repository: nested, ProfileID: id, SourceRepository: source.Repository, Commit: source.Ref, SourceDir: filepath.Join(root, "missing"),
	})
	if err == nil || payload.Finding == nil || payload.Finding.Code != "restore.source-dir-invalid" {
		t.Fatalf("payload = %+v, error = %v", payload, err)
	}
	if payload.Profile != id || payload.Setup != nil {
		t.Fatalf("profile/setup = %q/%v", payload.Profile, payload.Setup)
	}
}

func TestSkillSnapshotProfileKeepsEveryBuiltInProfile(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range catalog.ProfileIDs() {
		t.Run(id, func(t *testing.T) {
			want, err := loadRestoreProfile(catalog, id)
			if err != nil {
				t.Fatal(err)
			}
			got, err := loadSkillSnapshotProfile("", id, catalog)
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("profile %q = %+v, want %+v, error = %v", id, got, want, err)
			}
		})
	}
}

func TestSkillSnapshotProfileResolvesARepositoryProfile(t *testing.T) {
	t.Parallel()
	root, catalog, id := snapshotProfileFixture(t)
	resolved, err := ResolveProfile(root, id, catalog)
	if err != nil {
		t.Fatal(err)
	}
	required := make(map[string]struct{})
	for _, moduleID := range resolved.Modules {
		entry, ok := catalog.Module(moduleID)
		if !ok {
			t.Fatalf("missing module %s", moduleID)
		}
		var module struct {
			RequiredSkills []string `json:"requiredSkills"`
		}
		if err := json.Unmarshal(entry.Data, &module); err != nil {
			t.Fatal(err)
		}
		for _, name := range module.RequiredSkills {
			required[name] = struct{}{}
		}
	}
	builtin, err := loadRestoreProfile(catalog, "go-cli-tui")
	if err != nil {
		t.Fatal(err)
	}
	if builtin.Setup != "go" {
		t.Fatalf("source setup = %q", builtin.Setup)
	}
	contracts := make(map[string]restoreSkillContract)
	for name := range required {
		if contract, ok := builtin.Skills[name]; ok {
			contracts[name] = contract
		}
	}
	got, err := loadSkillSnapshotProfile(root, id, catalog)
	want := restoreProfile{ID: id, RequiredSkills: required, Skills: contracts}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("profile = %+v, want %+v, error = %v", got, want, err)
	}
}

func TestSkillSnapshotProfileNamesTheRepositoryPath(t *testing.T) {
	t.Parallel()
	root, catalog, _ := snapshotProfileFixture(t)
	_, err := loadSkillSnapshotProfile(root, "missing-profile", catalog)
	var finding *SkillsRestoreError
	if !errors.As(err, &finding) || finding.Category != SkillsRestoreInvalid || finding.Finding.Code != "restore.profile-unresolved" {
		t.Fatalf("error = %v", err)
	}
	if !strings.Contains(finding.Finding.Message, ".roundfix/baseline/profiles/missing-profile.json") || strings.Contains(finding.Finding.Message, "built-in Baseline Profile") {
		t.Fatalf("message = %s", finding.Finding.Message)
	}
	var resolution *ProfileResolutionError
	if !errors.As(err, &resolution) {
		t.Fatalf("resolution error was not wrapped: %v", err)
	}
	_, compareErr := TrailingSetupSkills(root, "missing-profile", nil)
	if compareErr == nil || !strings.HasPrefix(compareErr.Error(), "load profile for snapshot comparison: ") || !strings.Contains(compareErr.Error(), ".roundfix/baseline/profiles/missing-profile.json") {
		t.Fatalf("comparison error = %v", compareErr)
	}
}

func TestRepositorySnapshotContractsRefusesDisagreement(t *testing.T) {
	t.Parallel()
	contract := restoreSkillContract{Name: "alpha", Source: RestoreSource{Provider: "github", Repository: "owner/repo", Ref: "commit", Path: "skills/alpha"}, TreeDigest: "digest"}
	only := restoreSkillContract{Name: "only", Source: contract.Source, TreeDigest: contract.TreeDigest}
	required := map[string]struct{}{"alpha": {}, "zeta": {}, "only": {}, "owned": {}}
	snapshots := []setupSnapshotSkills{
		{ID: "first", Skills: map[string]restoreSkillContract{"alpha": contract, "only": only, "zeta": contract}},
		{ID: "second", Skills: map[string]restoreSkillContract{"alpha": contract, "zeta": contract}},
	}
	got, err := repositorySnapshotContracts("repo-profile", required, snapshots)
	want := map[string]restoreSkillContract{"alpha": contract, "only": only, "zeta": contract}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("agreeing merge = %+v, error = %v", got, err)
	}
	mutations := map[string]func(*restoreSkillContract){
		"tree":       func(c *restoreSkillContract) { c.TreeDigest = "different" },
		"provider":   func(c *restoreSkillContract) { c.Source.Provider = "different" },
		"repository": func(c *restoreSkillContract) { c.Source.Repository = "different" },
		"ref":        func(c *restoreSkillContract) { c.Source.Ref = "different" },
		"path":       func(c *restoreSkillContract) { c.Source.Path = "different" },
	}
	for field, mutate := range mutations {
		t.Run(field, func(t *testing.T) {
			changed := contract
			mutate(&changed)
			snapshots[1].Skills["alpha"] = changed
			snapshots[1].Skills["zeta"] = changed
			_, err := repositorySnapshotContracts("repo-profile", required, snapshots)
			var finding *SkillsRestoreError
			if !errors.As(err, &finding) || finding.Finding.Code != "restore.snapshot-conflict" {
				t.Fatalf("error = %v", err)
			}
			message := finding.Finding.Message
			if !strings.Contains(message, `skill "alpha"`) || !strings.Contains(message, `"first" and "second"`) {
				t.Fatalf("message = %s", message)
			}
		})
	}
}

func TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	snapshots := embeddedSnapshotContracts(catalog)
	required := make(map[string]struct{})
	for _, snapshot := range snapshots {
		for name := range snapshot.Skills {
			required[name] = struct{}{}
		}
	}
	if len(required) == 0 {
		t.Fatal("embedded snapshots have no external skills")
	}
	if _, err := repositorySnapshotContracts("catalog-agreement", required, snapshots); err != nil {
		t.Fatal(err)
	}
}

func TestTrailingSetupSkillsComparesARepositoryProfile(t *testing.T) {
	t.Parallel()
	root, catalog, id := snapshotProfileFixture(t)
	profile, err := loadSkillSnapshotProfile(root, id, catalog)
	if err != nil {
		t.Fatal(err)
	}
	const matching = "golang-cli"
	files, exists, err := inspectRestoreTarget("../..", ".agents/skills/"+matching)
	if err != nil {
		t.Fatal(err)
	}
	if !exists || portableRestoreDigest(files) != profile.Skills[matching].TreeDigest {
		t.Fatal("control tree does not match the snapshot")
	}
	for _, file := range files {
		target := filepath.Join(root, ".agents", "skills", matching, filepath.FromSlash(file.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, file.Content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	const edited = "testing-boss"
	writeInspectionFile(t, root, ".agents/skills/"+edited+"/SKILL.md", "older installed upstream tree\n")
	got, err := TrailingSetupSkills(root, id, []string{edited, matching, edited, "golang-context", "roundfix"})
	if err != nil || !reflect.DeepEqual(got, []string{edited}) {
		t.Fatalf("trailing = %v, error = %v", got, err)
	}
}
