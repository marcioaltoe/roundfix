package baseline

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest(t *testing.T) {
	t.Parallel()

	const guide = "docs/agents/agent-instructions.md"
	const body = "Managed guidance."
	artifact := func(relative, id string) ManifestArtifact {
		kind := "guide"
		if relative == "AGENTS.md" || relative == "CLAUDE.md" {
			kind = "root-block"
		}
		return ManifestArtifact{
			ID: id, Path: relative, Kind: kind, Module: "core", Template: id,
			Version: ManifestVersion, Digest: managedRegionBodyDigest(body),
		}
	}
	rule := func(relative string, start, end int, first string) MisplacedRepositoryRule {
		return MisplacedRepositoryRule{
			Path: relative, StartLine: start, EndLine: end, FirstLine: first,
			MoveTo: SpecificRepositoryPath,
		}
	}
	block := managedRegionFixture("guide.core", body)
	rootBlock := managedRegionFixture("root.core", body)
	foreign := "<!-- BEGIN:tool-rules -->\nForeign guidance.\n<!-- END:tool-rules -->\n"
	longLine := strings.Repeat("a", 119) + "界" + strings.Repeat("b", 8)
	mismatch := rule(guide, 1, 5, "Local guidance.")
	mismatch.ManagedID = "guide.core"
	tests := []struct {
		name       string
		files      map[string]string
		artifacts  []ManifestArtifact
		claudeLink bool
		noManifest bool
		invalid    bool
		extension  any
		readOnly   bool
		wantRules  []MisplacedRepositoryRule
		wantError  string
	}{
		{
			name:      "after_a_guide_managed_region",
			files:     map[string]string{guide: block + "\n\n  # Local rules  \n- Keep local policy.\n- Keep local conventions.\n\n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 8, 10, "# Local rules")},
		},
		{
			name:      "before_a_guide_managed_region",
			files:     map[string]string{guide: "# Local rules\nFirst rule.\nSecond rule.\nThird rule.\n\n" + block},
			wantRules: []MisplacedRepositoryRule{rule(guide, 1, 4, "# Local rules")},
		},
		{
			name:      "between_root_managed_regions_beside_a_foreign_block",
			files:     map[string]string{"AGENTS.md": rootBlock + "\n# Local rules\nKeep local policy.\n\n" + foreign + "\n" + managedRegionFixture("root.workflow", body)},
			artifacts: []ManifestArtifact{artifact("AGENTS.md", "root.core"), artifact("AGENTS.md", "root.workflow")},
			wantRules: []MisplacedRepositoryRule{rule("AGENTS.md", 7, 8, "# Local rules")},
		},
		{
			name:      "several_sections_after_one_managed_region",
			files:     map[string]string{guide: block + "\n<!-- local policy -->\n# Local policy\nFirst rule.\n\n<!-- local conventions -->\n# Local conventions\nSecond rule.\n\n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 7, 13, "<!-- local policy -->")},
		},
		{
			name:      "digest_mismatch_region",
			files:     map[string]string{guide: managedRegionFixture("guide.core", "Local guidance.")},
			wantRules: []MisplacedRepositoryRule{mismatch},
		},
		{
			name:  "marker_absent_region_is_not_an_entry",
			files: map[string]string{guide: "\n \n"},
		},
		{
			name:      "unclosed_begin_stays_in_the_passage",
			files:     map[string]string{guide: block + "\n<!-- BEGIN:tool-rules -->\nKeep this passage.\n\n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 7, 8, "<!-- BEGIN:tool-rules -->")},
		},
		{
			name:      "first_line_is_cut_at_120_bytes",
			files:     map[string]string{guide: block + "\n  " + longLine + " \n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 7, 7, strings.Repeat("a", 119)+"...")},
		},
		{name: "clean_carrier", files: map[string]string{guide: " \n\n" + block + "\n\t\n"}},
		{
			name:  "symbolic_link_claude_md",
			files: map[string]string{"AGENTS.md": rootBlock + "\n# Local rules\n"}, claudeLink: true,
			artifacts: []ManifestArtifact{artifact("AGENTS.md", "root.core"), artifact("CLAUDE.md", "claude.core")},
			wantRules: []MisplacedRepositoryRule{rule("AGENTS.md", 7, 7, "# Local rules")},
		},
		{
			name: "regular_claude_md", files: map[string]string{"CLAUDE.md": "\n# Local rules\nKeep this rule.\n\n"},
			artifacts: []ManifestArtifact{artifact("CLAUDE.md", "root.core")},
			wantRules: []MisplacedRepositoryRule{rule("CLAUDE.md", 2, 3, "# Local rules")},
		},
		{
			name:      "specific_repository_is_not_read",
			files:     map[string]string{guide: block, SpecificRepositoryPath: "Local rules live here.\n"},
			artifacts: []ManifestArtifact{artifact(guide, "guide.core"), artifact(SpecificRepositoryPath, "guide.local")},
		},
		{name: "no_manifest", noManifest: true, files: map[string]string{"AGENTS.md": "Local rules.\n"}},
		{name: "invalid_manifest", invalid: true, wantError: "read misplaced-rule Setup Manifest: baseline.preservation.managed-refresh.manifest.invalid"},
		{name: "extension_decision_off", extension: false, files: map[string]string{guide: block}},
		{name: "extension_decision_on", extension: true, files: map[string]string{guide: block}},
		{
			name: "reads_only", readOnly: true,
			files:     map[string]string{guide: block + "\n# Local policy\n", "unlisted.md": "Leave this alone.\n", SpecificRepositoryPath: "Local rules.\n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 7, 7, "# Local policy")},
		},
		{
			name:      "foreign_block_splits_a_passage",
			files:     map[string]string{guide: block + "\nFirst passage.\n\n" + foreign + "\nSecond passage.\n\n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 7, 7, "First passage."), rule(guide, 13, 13, "Second passage.")},
		},
		{
			name:      "foreign_block_matches_only_the_same_name",
			files:     map[string]string{guide: block + "\n  <!-- BEGIN:tool-rules --> \nForeign.\n<!-- END:other -->\n<!-- END:tool-rules -->\n\nLocal.\n"},
			wantRules: []MisplacedRepositoryRule{rule(guide, 12, 12, "Local.")},
		},
		{
			name:      "excluded_carriers_are_not_opened",
			files:     map[string]string{guide: managedRegionFixture("guide.core", "Local guidance.")},
			artifacts: []ManifestArtifact{artifact(guide, "guide.core"), artifact("nested/AGENTS.md", "root.nested"), artifact("nested/CLAUDE.md", "claude.nested"), artifact(SpecificRepositoryPath, "guide.local"), artifact("instructions.txt", "guide.text")},
			wantRules: []MisplacedRepositoryRule{mismatch},
		},
		{
			name:  "unlisted_carriers_are_not_read",
			files: map[string]string{guide: block, "AGENTS.md": "Unlisted root.\n", "nested/AGENTS.md": "Nested.\n", "nested/CLAUDE.md": "Nested alias.\n", "docs/agents/unlisted.md": "Unlisted guide.\n"},
		},
		{
			name:      "unlisted_regular_claude_md_is_read",
			files:     map[string]string{guide: block, "AGENTS.md": "Unlisted root.\n", "CLAUDE.md": "\n# Root repository policy\nKeep this rule.\n"},
			wantRules: []MisplacedRepositoryRule{rule("CLAUDE.md", 2, 3, "# Root repository policy")},
		},
		{
			name:  "unlisted_symbolic_link_claude_md_is_not_read",
			files: map[string]string{guide: block, "AGENTS.md": "Unlisted root.\n"}, claudeLink: true,
		},
		{
			name:      "missing_listed_carrier",
			files:     map[string]string{},
			wantError: fmt.Sprintf("read misplaced-rule carrier %q: %s", guide, fs.ErrNotExist),
		},
		{
			name:      "paths_then_lines",
			files:     map[string]string{guide: "First.\n\n" + block + "\nLast.\n", "AGENTS.md": "Root.\n" + rootBlock},
			artifacts: []ManifestArtifact{artifact(guide, "guide.core"), artifact("AGENTS.md", "root.core")},
			wantRules: []MisplacedRepositoryRule{rule("AGENTS.md", 1, 1, "Root."), rule(guide, 1, 1, "First."), rule(guide, 9, 9, "Last.")},
		},
		{
			name:      "crlf_and_no_final_newline",
			files:     map[string]string{guide: strings.ReplaceAll(block, "\n", "\r\n") + "\r\n \r\n  Local.  \r\n\r\nLast."},
			wantRules: []MisplacedRepositoryRule{rule(guide, 8, 10, "Local.")},
		},
		{name: "extension_string_is_not_true", extension: "true", files: map[string]string{guide: block}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := t.TempDir()
			for relative, data := range tt.files {
				writeInspectionFile(t, repo, relative, data)
			}
			if tt.claudeLink {
				if err := os.Symlink("AGENTS.md", filepath.Join(repo, "CLAUDE.md")); err != nil {
					t.Fatalf("create instruction alias: %v", err)
				}
			}
			if !tt.noManifest {
				artifacts := tt.artifacts
				if artifacts == nil {
					artifacts = []ManifestArtifact{artifact(guide, "guide.core")}
				}
				writeManagedRefreshManifest(t, repo, artifacts[0])
				manifestBytes, err := os.ReadFile(filepath.Join(repo, manifestPath))
				if err != nil {
					t.Fatalf("read fixture manifest: %v", err)
				}
				var manifest SetupManifest
				if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
					t.Fatalf("parse fixture manifest: %v", err)
				}
				manifest.ManagedArtifacts = artifacts
				if tt.extension != nil {
					manifest.Decisions["repository.extension.enabled"] = ManifestDecision{Value: tt.extension}
				}
				encoded, err := json.Marshal(manifest)
				if err != nil {
					t.Fatalf("encode fixture manifest: %v", err)
				}
				if tt.invalid {
					encoded = []byte("{}")
				}
				writeInspectionFile(t, repo, manifestPath, string(encoded))
			}
			if tt.readOnly {
				for _, relative := range []string{guide, "unlisted.md", SpecificRepositoryPath, manifestPath} {
					if err := os.Chmod(filepath.Join(repo, relative), 0o440); err != nil {
						t.Fatalf("set fixture mode: %v", err)
					}
				}
			}
			before := snapshotInspectionTree(t, repo)
			got, err := MisplacedRepositoryRules(repo)
			after := snapshotInspectionTree(t, repo)
			if !reflect.DeepEqual(before, after) {
				t.Errorf("repository changed: before = %#v, after = %#v", before, after)
			}
			gotError := ""
			if err != nil {
				gotError = err.Error()
			}
			if gotError != tt.wantError {
				t.Fatalf("error = %q, want %q", gotError, tt.wantError)
			}
			want := MisplacedRepositoryRuleReport{
				ManifestFound: !tt.noManifest, ExtensionEnabled: tt.extension == true, Rules: tt.wantRules,
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("report = %#v, want %#v", got, want)
			}
		})
	}
}
