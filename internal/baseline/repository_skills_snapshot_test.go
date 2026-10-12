package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	var manifest struct {
		Profile string `json:"profile"`
	}
	manifestBytes, err := os.ReadFile(filepath.Join(repoRoot, "docs", "agents", "setup-context.json"))
	if err != nil {
		t.Fatalf("read Setup Manifest: %v", err)
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("decode Setup Manifest: %v", err)
	}
	if manifest.Profile == "" {
		t.Fatal("Setup Manifest does not name a Baseline Profile")
	}

	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	lockBytes, err := os.ReadFile(filepath.Join(repoRoot, "skills-lock.json"))
	if err != nil {
		t.Fatalf("read skills lock: %v", err)
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	required := make([]string, 0, len(lock.Skills))
	for name := range lock.Skills {
		required = append(required, name)
	}
	sort.Strings(required)
	differing, err := SetupSnapshotDifferences(repoRoot, manifest.Profile, required)
	if err != nil {
		t.Fatalf("compare required skills with the Setup Snapshot: %v", err)
	}
	for _, name := range differing {
		t.Errorf("required skill %q differs from the Setup Snapshot; restore deliberately with roundfix baseline skills restore --profile %s --skill %s", name, manifest.Profile, name)
	}
}
