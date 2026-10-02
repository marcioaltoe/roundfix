package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/skills"
)

func TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot(t *testing.T) {
	root := t.TempDir()
	if _, err := skills.Install(t.Context(), skills.InstallRequest{Target: "project", ProjectDir: root}); err != nil {
		t.Fatal(err)
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := loadRestoreProfile(catalog, "go-cli-tui")
	if err != nil {
		t.Fatal(err)
	}
	// Read a repository-held upstream tree and prove it matches the embedded
	// snapshot before using it as the unchanged control.
	const matching = "golang-cli"
	files, exists, err := inspectRestoreTarget("../..", ".agents/skills/"+matching)
	if err != nil {
		t.Fatal(err)
	}
	if !exists || portableRestoreDigest(files) != profile.Skills[matching].TreeDigest {
		t.Fatal("control tree does not match the embedded Setup Snapshot")
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
	const trailing = "testing-boss"
	dir := filepath.Join(root, ".agents", "skills", trailing)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("older installed upstream tree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries := map[string]any{}
	for _, name := range []string{matching, trailing} {
		hash, err := skills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", name))
		if err != nil {
			t.Fatal(err)
		}
		entries[name] = map[string]any{"computedHash": hash}
	}
	data, err := json.Marshal(map[string]any{"version": 1, "skills": entries})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills-lock.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	readiness, err := skills.CheckRepositoryWithExternal(t.Context(), root, []string{matching, trailing})
	if err != nil || !readiness.Ready() {
		t.Fatalf("lock readiness = %+v, error = %v", readiness, err)
	}
	got, err := TrailingSetupSkills(root, profile.ID, []string{trailing, "golang-context", matching, "roundfix"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{trailing}) {
		t.Fatalf("trailing = %v", got)
	}
	// A second installed difference proves sorting and duplicate elimination.
	secondDir := filepath.Join(root, ".agents", "skills", "golang-context")
	if err := os.MkdirAll(secondDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(secondDir, "SKILL.md"), []byte("old context skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = TrailingSetupSkills(root, profile.ID, []string{trailing, "golang-context", trailing, matching})
	if err != nil || !reflect.DeepEqual(got, []string{"golang-context", trailing}) {
		t.Fatalf("sorted trailing = %v, error = %v", got, err)
	}
	if _, err := TrailingSetupSkills(root, "unknown-profile", []string{trailing}); err == nil {
		t.Fatal("unknown profile must make snapshot comparison unavailable")
	}
}
