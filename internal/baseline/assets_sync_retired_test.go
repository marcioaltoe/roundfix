// Suite: retired asset-sync skills
// Invariant: upstream and recorded repo entries cannot restore Retired Skills; other entries survive.
// Boundary IN: committed upstream lists and recorded snapshots in temporary Git repositories.
// Boundary OUT: CLI output and installed skill trees.
package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAssetSyncDropsARetiredSkillTheUpstreamListNames(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		list string
	}{
		{"upstream names both", "skills/council\nskills/the-fool\n"},
		{"duplicate retired names", "skills/council\nskills/council\nskills/the-fool\nskills/the-fool\n"},
		{"upstream omits recorded council", "skills/the-fool\n"},
		{"non-retired control", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			sourceDir := filepath.Join(root, "setups")
			if err := os.MkdirAll(sourceDir, 0o755); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"council", "the-fool", "retained"} {
				tree := filepath.Join(root, "skills", name)
				if err := os.MkdirAll(tree, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(tree, "SKILL.md"), []byte("# "+name+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			list := "skills/qa-gate\n" + tc.list + "skills/retained\n"
			if err := os.WriteFile(filepath.Join(sourceDir, "go.txt"), []byte(list), 0o644); err != nil {
				t.Fatal(err)
			}
			runAssetsSyncGit(t, root, "init", "--quiet")
			runAssetsSyncGit(t, root, "add", ".")
			runAssetsSyncGit(t, root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.com", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "upstream skills")
			revision := strings.TrimSpace(runAssetsSyncGit(t, root, "rev-parse", "HEAD"))
			owned := assetsSyncSource{Type: "repo", Name: "roundfix"}
			qa := assetsSyncSkill{Name: "qa-gate", Path: "skills/qa-gate", Source: owned, MinimumVersion: "0.0.2"}
			council := assetsSyncSkill{Name: "council", Path: "skills/council", Source: owned, MinimumVersion: "0.0.2"}
			roundfix := assetsSyncSkill{Name: "roundfix", Path: "skills/roundfix", Source: owned, MinimumVersion: "0.0.3"}
			snapshot, findings := buildAssetsSyncSnapshot(t.Context(), "go", sourceDir,
				assetsSyncSnapshot{Skills: []assetsSyncSkill{council, qa, roundfix}},
				assetsSyncCheckout{root: root, repository: "example/skills", revision: revision})
			if len(findings) != 0 || snapshot == nil {
				t.Fatalf("snapshot = %+v, findings = %+v", snapshot, findings)
			}
			digest, err := assetsSyncFilesystemTreeDigest(filepath.Join(root, "skills", "retained"))
			if err != nil {
				t.Fatal(err)
			}
			retained := assetsSyncSkill{Name: "retained", Path: "skills/retained", TreeDigest: digest,
				Source: assetsSyncSource{Type: "github", Repository: "example/skills", Ref: revision, Path: "skills/retained"}}
			if want := []assetsSyncSkill{qa, retained, roundfix}; !reflect.DeepEqual(snapshot.Skills, want) {
				t.Fatalf("synced skills = %+v, want %+v", snapshot.Skills, want)
			}
		})
	}
}
