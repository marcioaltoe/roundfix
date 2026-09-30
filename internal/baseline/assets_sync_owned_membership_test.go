// Suite: asset sync owned membership
// Invariant: omitted repo entries survive in recorded order while omitted external entries are dropped.
// Boundary IN: recorded snapshots and committed upstream lists in temporary Git repositories.
// Boundary OUT: asset transactions, catalog loading, and CLI output.

package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func buildOmittedSkillSnapshot(t *testing.T, recorded []assetsSyncSkill) *assetsSyncSnapshot {
	t.Helper()
	root := t.TempDir()
	sourceDir := filepath.Join(root, "setups")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "go-cli.txt"), []byte("skills/00-setup/qa-gate\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runAssetsSyncGit(t, root, "init", "--quiet")
	runAssetsSyncGit(t, root, "add", ".")
	runAssetsSyncGit(t, root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.com", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "upstream omits skills")
	revision := strings.TrimSpace(runAssetsSyncGit(t, root, "rev-parse", "HEAD"))
	snapshot, findings := buildAssetsSyncSnapshot(t.Context(), "go-cli", sourceDir,
		assetsSyncSnapshot{Skills: recorded},
		assetsSyncCheckout{root: root, repository: "example/skills", revision: revision})
	if len(findings) != 0 || snapshot == nil {
		t.Fatalf("snapshot = %+v, findings = %+v", snapshot, findings)
	}
	return snapshot
}

func TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits(t *testing.T) {
	qa := assetsSyncSkill{Name: "qa-gate", Path: "skills/00-setup/qa-gate", Source: assetsSyncSource{Type: "repo", Name: "roundfix"}, MinimumVersion: "0.0.2"}
	roundfix := assetsSyncSkill{Name: "roundfix", Path: "skills/06-review-repair/roundfix", Source: qa.Source, MinimumVersion: "0.0.2"}
	other := assetsSyncSkill{Name: "custom-owned", Path: "skills/custom-owned", Source: qa.Source, MinimumVersion: "0.1.0"}
	snapshot := buildOmittedSkillSnapshot(t, []assetsSyncSkill{roundfix, qa, other})
	if want := []assetsSyncSkill{qa, roundfix, other}; !reflect.DeepEqual(snapshot.Skills, want) {
		t.Fatalf("synced skills = %+v, want %+v", snapshot.Skills, want)
	}
}

func TestAssetSyncStillDropsAnExternalSkillTheUpstreamListOmits(t *testing.T) {
	qa := assetsSyncSkill{Name: "qa-gate", Path: "skills/00-setup/qa-gate", Source: assetsSyncSource{Type: "repo", Name: "roundfix"}, MinimumVersion: "0.0.2"}
	for _, sourceType := range []string{"github", "local"} {
		t.Run(sourceType, func(t *testing.T) {
			external := assetsSyncSkill{Name: "external", Path: "skills/external", Source: assetsSyncSource{Type: sourceType}}
			snapshot := buildOmittedSkillSnapshot(t, []assetsSyncSkill{external, qa})
			if want := []assetsSyncSkill{qa}; !reflect.DeepEqual(snapshot.Skills, want) {
				t.Fatalf("synced skills = %+v, want %+v", snapshot.Skills, want)
			}
		})
	}
}
