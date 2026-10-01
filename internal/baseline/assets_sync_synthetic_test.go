package baseline

import (
	"fmt"
	"path"
	"strings"
	"testing"
)

func assetsSyncSyntheticSkillFile(skillPath string) []byte {
	return []byte("# " + path.Base(skillPath) + "\n")
}

func assetsSyncSyntheticTreeDigest(skillPath string) string {
	return portableRestoreDigest([]restoreFile{{Path: "SKILL.md", Content: assetsSyncSyntheticSkillFile(skillPath)}})
}

func parityFixtureDigestFindings(fixture map[string]any) []string {
	findings := []string{}
	manifest, _ := fixture["manifest"].(map[string]any)
	setups, _ := manifest["setups"].([]any)
	for _, rawSetup := range setups {
		setup, _ := rawSetup.(map[string]any)
		skills, _ := setup["skills"].([]any)
		for _, rawSkill := range skills {
			skill, _ := rawSkill.(map[string]any)
			source, _ := skill["source"].(map[string]any)
			if source["type"] != "github" {
				continue
			}
			skillPath, _ := skill["path"].(string)
			want := assetsSyncSyntheticTreeDigest(skillPath)
			if skill["treeDigest"] != want {
				findings = append(findings, fmt.Sprintf("setup %s skill %s: treeDigest = %v, want %s", setup["id"], skillPath, skill["treeDigest"], want))
			}
		}
	}
	return findings
}

func TestTheParityFixtureDigestsFollowTheSyntheticSource(t *testing.T) {
	t.Parallel()
	fixture := readBaselineCompatibilityJSON[map[string]any](t, "fixtures/asset-sync.json")
	if findings := parityFixtureDigestFindings(fixture); len(findings) != 0 {
		t.Fatalf("synthetic digest mismatches: %v", findings)
	}
}

func TestASyntheticDigestThatDiffersIsReported(t *testing.T) {
	t.Parallel()
	fixture := map[string]any{"manifest": map[string]any{"setups": []any{
		map[string]any{"id": "synthetic", "skills": []any{
			map[string]any{"path": "skills/altered", "source": map[string]any{"type": "github"}, "treeDigest": strings.Repeat("0", 64)},
		}},
	}}}
	findings := parityFixtureDigestFindings(fixture)
	if len(findings) != 1 || !strings.Contains(findings[0], "setup synthetic skill skills/altered") {
		t.Fatalf("findings = %v, want one altered skill", findings)
	}
}
