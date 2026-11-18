package config

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

func TestRepositoryProfilesFollowTheRecommendedProfile(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile(filepath.Join("..", "..", projectConfigName))
	if err != nil {
		t.Fatal(err)
	}
	home, work := t.TempDir(), t.TempDir()
	gittest.MarkRepository(t, work)
	if err := os.WriteFile(filepath.Join(work, projectConfigName), content, 0600); err != nil {
		t.Fatal(err)
	}
	// Use the Project Config parser, including its strict profile-key checks,
	// with only a temporary home and repository; loading needs no runtime or network.
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work, Stderr: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	configured := 0
	for _, category := range AllWorkCategories() {
		entry, ok := loaded.Config.Profiles[category]
		if !ok || entry.Source != ProfileSourceProject {
			continue
		}
		configured++
		t.Run(string(category), func(t *testing.T) {
			want, ok := RecommendedProfile(category)
			if !ok || !reflect.DeepEqual(entry.Profile, want) {
				t.Errorf("profile = %+v, want Recommended Profile %+v (exists: %t)", entry.Profile, want, ok)
			}
			if entry.Deviation != nil {
				t.Errorf("repository profile carries a deviation: %+v", entry.Deviation)
			}
		})
	}
	if configured == 0 {
		t.Fatal("repository defines no profiles")
	}
}
