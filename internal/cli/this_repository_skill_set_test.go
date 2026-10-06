package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"roundfix/skills"
)

func TestThisRepositoryHoldsEveryRequiredExternalSkill(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	readiness := checkThisRepositorySkillSet(t, root, true)
	if !readiness.Ready() {
		t.Fatalf("repository skill set is not ready: missing=%v outdated=%v owned=%v", readiness.MissingExternal, readiness.OutdatedExternal, readiness.MissingOwned)
	}
}

func TestARepositoryMissingARequiredExternalSkillIsReported(t *testing.T) {
	root := copyThisRepositorySkillSetFixture(t)
	external, ok, err := resolveExternalSkillRequirement(root)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || len(external) == 0 {
		t.Fatalf("fixture external requirement = %v, manifest readable = %t", external, ok)
	}
	missing := external[0]
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", missing)); err != nil {
		t.Fatalf("remove fixture skill %q: %v", missing, err)
	}

	readiness := checkThisRepositorySkillSet(t, root, false)
	if readiness.Ready() {
		t.Fatal("repository with a missing required external skill was reported ready")
	}
	found := false
	for _, name := range readiness.MissingExternal {
		if name == missing {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing external skills = %v, want %q", readiness.MissingExternal, missing)
	}
}

func checkThisRepositorySkillSet(t *testing.T, root string, wantReady bool) skills.RepositoryReadiness {
	t.Helper()
	external, ok, err := resolveExternalSkillRequirement(root)
	if err != nil {
		t.Fatalf("resolve external skill requirement: %v", err)
	}
	if !ok {
		t.Fatal("Setup Manifest did not resolve an external skill requirement")
	}
	readiness, err := skills.CheckRepositoryWithExternal(t.Context(), root, external)
	if err != nil {
		t.Fatalf("check repository skill set: %v", err)
	}
	for _, name := range readiness.MissingExternal {
		t.Logf("missing external skill: %s", name)
	}
	for _, name := range readiness.OutdatedExternal {
		t.Logf("outdated external skill: %s", name)
	}
	if wantReady && !readiness.Ready() {
		t.Errorf("repository skill set is not ready: missing=%v outdated=%v", readiness.MissingExternal, readiness.OutdatedExternal)
	}

	lockBytes, err := os.ReadFile(filepath.Join(root, "skills-lock.json"))
	if err != nil {
		t.Fatalf("read skills lock: %v", err)
	}
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	for _, name := range []string{"context7", "feature-systems-pattern", "rust", "review", "triage", "the-fool", "autoresearch", "council"} {
		if _, exists := lock.Skills[name]; exists {
			t.Errorf("obsolete skill remains in skills-lock.json: %s", name)
		}
		if _, err := os.Stat(filepath.Join(root, ".agents", "skills", name)); err == nil {
			t.Errorf("obsolete skill directory remains: %s", name)
		} else if !os.IsNotExist(err) {
			t.Errorf("inspect obsolete skill directory %s: %v", name, err)
		}
	}
	return readiness
}

func copyThisRepositorySkillSetFixture(t *testing.T) string {
	t.Helper()
	source, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir()
	for _, relative := range []string{
		"docs/agents/setup-context.json",
		"skills-lock.json",
		".agents/skills",
	} {
		copyThisRepositorySkillSetPath(t, filepath.Join(source, relative), filepath.Join(destination, relative))
	}
	return destination
}

func copyThisRepositorySkillSetPath(t *testing.T, source, destination string) {
	t.Helper()
	info, err := os.Stat(source)
	if err != nil {
		t.Fatalf("stat fixture source %s: %v", source, err)
	}
	if info.IsDir() {
		if err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			relative, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}
			target := filepath.Join(destination, relative)
			if info.IsDir() {
				return os.MkdirAll(target, info.Mode().Perm())
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, data, info.Mode().Perm())
		}); err != nil {
			t.Fatalf("copy fixture directory %s: %v", source, err)
		}
		return
	}
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read fixture source %s: %v", source, err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, data, info.Mode().Perm()); err != nil {
		t.Fatalf("write fixture destination %s: %v", destination, err)
	}
}
