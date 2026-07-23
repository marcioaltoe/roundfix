// Suite: repository-owned Baseline Profiles
// Invariant: custom profiles resolve only embedded catalog entries from their exact repository-owned path.
// Boundary IN: strict profile JSON, catalog references, repository discovery, paths, and profile digests.
// Boundary OUT: CLI rendering, Agent Selection Profiles, Baseline planning, and apply.

package baseline

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestCustomProfileInitAndLoad(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}

	created, err := InitCustomProfile(repo, "team-go", "go-cli-tui", catalog)
	if err != nil {
		t.Fatalf("InitCustomProfile() error = %v", err)
	}
	wantPath := filepath.Join(repo, ".roundfix", "baseline", "profiles", "team-go.json")
	if created.Path != wantPath {
		t.Fatalf("InitCustomProfile() path = %q, want %q", created.Path, wantPath)
	}
	if created.Profile.Source != ProfileSourceRepository {
		t.Fatalf("InitCustomProfile() source = %q, want %q", created.Profile.Source, ProfileSourceRepository)
	}
	if created.Profile.ID != "team-go" || created.Profile.Digest == "" {
		t.Fatalf("InitCustomProfile() profile = %#v", created.Profile)
	}

	loaded, err := LoadRepositoryProfile(repo, "team-go", catalog)
	if err != nil {
		t.Fatalf("LoadRepositoryProfile() error = %v", err)
	}
	if loaded.Digest != created.Profile.Digest {
		t.Fatalf("loaded digest = %q, want %q", loaded.Digest, created.Profile.Digest)
	}
	if !slices.Equal(loaded.Modules, created.Profile.Modules) ||
		!slices.Equal(loaded.Decisions, created.Profile.Decisions) ||
		!slices.Equal(loaded.Templates, created.Profile.Templates) {
		t.Fatalf("loaded profile differs from initialized profile:\nloaded=%#v\ncreated=%#v", loaded, created.Profile)
	}

	discovered, err := DiscoverRepositoryProfiles(repo, catalog)
	if err != nil {
		t.Fatalf("DiscoverRepositoryProfiles() error = %v", err)
	}
	if len(discovered) != 1 || discovered[0].ID != "team-go" {
		t.Fatalf("DiscoverRepositoryProfiles() = %#v, want team-go", discovered)
	}
}

func TestCustomProfileRejectsInvalidDeclarations(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}
	valid := `{
  "schemaVersion": "roundfix/custom-baseline-profile/v1",
  "catalogSchema": "roundfix/baseline-catalog/v1",
  "id": "team",
  "modules": ["core"],
  "decisions": ["language.generated"],
  "capabilities": [],
  "templates": ["template.root.core"],
  "values": {"language.generated": "English"}
}`
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "unknown module",
			content: strings.Replace(valid, `"core"`, `"missing-module"`, 1),
			want:    "custom.profile.module.unknown",
		},
		{
			name:    "remote reference",
			content: strings.Replace(valid, `"values":`, `"remote": "https://example.test/profile.json", "values":`, 1),
			want:    "custom.profile.field.unknown",
		},
		{
			name:    "custom asset",
			content: strings.Replace(valid, `"values":`, `"assets": ["local.md"], "values":`, 1),
			want:    "custom.profile.field.unknown",
		},
		{
			name:    "executable content",
			content: strings.Replace(valid, `"values":`, `"command": "make verify", "values":`, 1),
			want:    "custom.profile.field.unknown",
		},
		{
			name:    "profile composition",
			content: strings.Replace(valid, `"values":`, `"profiles": ["go-cli-tui"], "values":`, 1),
			want:    "custom.profile.field.unknown",
		},
		{
			name:    "duplicate key",
			content: strings.Replace(valid, `"id": "team",`, `"id": "team", "id": "other",`, 1),
			want:    "custom.profile.json.key.duplicate",
		},
		{
			name:    "invalid decision value",
			content: strings.Replace(valid, `"English"`, `"Portuguese"`, 1),
			want:    "custom.profile.value.invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParseCustomProfile([]byte(tt.content), "team.json", catalog)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ParseCustomProfile() error = %v, want code %q", err, tt.want)
			}
		})
	}
}

func TestCustomProfileRejectsUnsafePathsAndNonRepositorySources(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}
	outside := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(outside, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write outside profile: %v", err)
	}
	if _, err := LoadCustomProfilePath(repo, outside, catalog); !errors.Is(err, ErrUnsafeCustomProfilePath) {
		t.Fatalf("LoadCustomProfilePath(outside) error = %v, want ErrUnsafeCustomProfilePath", err)
	}

	profileDir := filepath.Join(repo, ".roundfix", "baseline")
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		t.Fatalf("create profile parent: %v", err)
	}
	if err := os.Symlink(filepath.Dir(outside), filepath.Join(profileDir, "profiles")); err != nil {
		t.Fatalf("create escaping profiles symlink: %v", err)
	}
	if _, err := LoadCustomProfilePath(repo, filepath.Join(profileDir, "profiles", "outside.json"), catalog); !errors.Is(err, ErrUnsafeCustomProfilePath) {
		t.Fatalf("LoadCustomProfilePath(symlink) error = %v, want ErrUnsafeCustomProfilePath", err)
	}
}
