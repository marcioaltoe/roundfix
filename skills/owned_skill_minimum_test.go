package skills

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// minimumVersionDifferences compares the contract map with actual embedded
// declarations. Keeping the comparison shared exercises its negative case too.
func minimumVersionDifferences(t *testing.T, minimums map[string]string) []string {
	t.Helper()
	var differences []string
	for _, name := range Names() {
		data, err := embedded.ReadFile(name + "/SKILL.md")
		if err != nil {
			t.Fatal(err)
		}
		metadata, ok := parseSkillFrontmatter(string(data))
		if !ok || !ValidVersion(metadata.Version) {
			t.Fatalf("%s has no valid embedded version", name)
		}
		if minimums[name] != metadata.Version {
			differences = append(differences, name)
		}
	}
	sort.Strings(differences)
	return differences
}

func TestTheOwnedSkillMinimumIsTheEmbeddedVersion(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := Install(t.Context(), InstallRequest{Target: "project", ProjectDir: root}); err != nil {
		t.Fatal(err)
	}
	if got := minimumVersionDifferences(t, ownedSkillMinimumVersions); len(got) != 0 {
		t.Fatalf("minimum differs from embedded version: %v", got)
	}
	if len(ownedSkillMinimumVersions) != len(Names()) {
		t.Fatal("minimum map does not match owned set")
	}
}

func TestAMinimumThatDiffersFromTheEmbeddedVersionIsReported(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := Install(t.Context(), InstallRequest{Target: "project", ProjectDir: root}); err != nil {
		t.Fatal(err)
	}
	minimums := make(map[string]string, len(ownedSkillMinimumVersions))
	for name, version := range ownedSkillMinimumVersions {
		minimums[name] = version
	}
	minimums["qa-gate"] = skillVersionBelow(t, minimums["qa-gate"])
	if got := minimumVersionDifferences(t, minimums); !reflect.DeepEqual(got, []string{"qa-gate"}) {
		t.Fatalf("differences = %v, want qa-gate", got)
	}
}

func TestAnInstalledOwnedSkillOlderThanTheBundleIsBelow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := Install(t.Context(), InstallRequest{Target: "project", ProjectDir: root}); err != nil {
		t.Fatal(err)
	}
	const name = "qa-gate"
	minimum := ownedSkillMinimumVersions[name]
	below := skillVersionBelow(t, minimum)
	path := filepath.Join(root, ".agents", "skills", name, "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("\nversion: "+minimum+"\n"), []byte("\nversion: "+below+"\n"), 1)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	readiness, err := CheckRepositoryWithExternal(t.Context(), root, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, owned := range readiness.Owned {
		if owned.Skill == name {
			if owned.State != ReadinessBelow || owned.Minimum != minimum || owned.Found != below {
				t.Fatalf("owned readiness = %+v", owned)
			}
			if readiness.Ready() {
				t.Fatal("older installed skill is ready")
			}
			return
		}
	}
	t.Fatal("qa-gate not reported")
}
