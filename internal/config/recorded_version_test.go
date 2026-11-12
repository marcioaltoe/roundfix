// Suite: repository Recorded Roundfix Version compatibility
// Invariant: configuration refuses a newer record before reading Project Config.
// Boundary IN: local Setup Manifest, layered config, and release comparison.
// Boundary OUT: Baseline writing, owned by internal/baseline.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/app"
	"roundfix/internal/gittest"
)

const futureRecordedRoundfixVersion = "999999.0.0"

func TestARepositoryRecordedByANewerRoundfixRefusesLoad(t *testing.T) {
	t.Parallel()
	home, repo := recordedVersionConfigRepository(t, `{"roundfixVersion":"`+futureRecordedRoundfixVersion+`"}`, "future_section:\n  enabled: true\n")
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
	var refusal *NewerRepositoryError
	if !errors.As(err, &refusal) {
		t.Fatalf("expected newer repository refusal before YAML decode, got %v", err)
	}
	path := filepath.Join(repo, SetupManifestPath)
	want := fmt.Sprintf("Setup Manifest %q records Roundfix %s, newer than this binary's %s; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'", path, futureRecordedRoundfixVersion, app.NormalizeVersion(app.Version))
	if err != refusal || err.Error() != want {
		t.Fatalf("unwrapped refusal = %v, want %s", err, want)
	}
	if refusal.Path != path || refusal.Recorded != futureRecordedRoundfixVersion || refusal.Running != app.NormalizeVersion(app.Version) {
		t.Fatalf("refusal fields = %#v", refusal)
	}
	if loaded.ProjectConfigPath != filepath.Join(repo, projectConfigName) {
		t.Fatalf("Project Config path = %q", loaded.ProjectConfigPath)
	}
}

func TestAnEqualOrOlderRecordedRoundfixVersionLoadsProjectConfig(t *testing.T) {
	t.Parallel()
	for _, version := range []string{app.NormalizeVersion(app.Version), "0.0.1"} {
		t.Run(version, func(t *testing.T) {
			t.Parallel()
			home, repo := recordedVersionConfigRepository(t, `{"roundfixVersion":"`+version+`"}`, "update:\n  mode: auto\n")
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil || loaded.RepositoryVersionErr != nil || loaded.Config.Update.Mode != "auto" {
				t.Fatalf("loaded = %#v, err = %v", loaded, err)
			}
		})
	}
}

func TestARepositoryWithoutARecordedRoundfixVersionLoads(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, manifest string }{
		{"absent manifest", ""}, {"no field", `{}`}, {"not JSON", `not JSON`},
		{"array", `[]`}, {"null", `null`}, {"string", `"release"`}, {"unreadable file", "directory"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home, repo := recordedVersionConfigRepository(t, tt.manifest, "update:\n  mode: auto\n")
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil || loaded.RepositoryVersionErr != nil || loaded.Config.Update.Mode != "auto" {
				t.Fatalf("mode = %q, err = %v, record = %v", loaded.Config.Update.Mode, err, loaded.RepositoryVersionErr)
			}
		})
	}
	if version, err := RecordedRoundfixVersion(" "); version != "" || err != nil {
		t.Fatalf("blank root = %q, %v", version, err)
	}
}

func TestAMalformedRecordedRoundfixVersionRefusesLoad(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{`"0.140"`, `"v0.140.0"`, `"0.140.0-rc.1"`, `140`, `null`, `true`, `{}`, `[]`, `""`, `" 0.140.0"`, `"0.140.0\n"`} {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()
			home, repo := recordedVersionConfigRepository(t, `{"roundfixVersion":`+raw+`}`, "future_section:\n  enabled: true\n")
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			want := fmt.Sprintf("Setup Manifest %q records roundfixVersion %s, which is not a Roundfix release version such as 0.140.0; restore the value the Roundfix that wrote it recorded", filepath.Join(repo, SetupManifestPath), raw)
			if err == nil || err.Error() != want {
				t.Fatalf("refusal = %v, want %s", err, want)
			}
		})
	}
}

func TestAllowNewerRepositoryLoadsWithoutProjectConfig(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, project string }{
		{"valid project config", "update:\n  mode: auto\n"},
		{"unknown project key", "update:\n  mode: auto\nfuture_section:\n  enabled: true\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home, repo := recordedVersionConfigRepository(t, `{"roundfixVersion":"`+futureRecordedRoundfixVersion+`"}`, tt.project)
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo, AllowNewerRepository: true})
			var refusal *NewerRepositoryError
			if err != nil || !errors.As(loaded.RepositoryVersionErr, &refusal) {
				t.Fatalf("load error = %v, repository refusal = %v", err, loaded.RepositoryVersionErr)
			}
			if loaded.Config.Update.Mode != "notify" {
				t.Fatalf("Project Config was read: update.mode = %q", loaded.Config.Update.Mode)
			}
			if loaded.ProjectConfigPath != filepath.Join(repo, projectConfigName) {
				t.Fatalf("Project Config path = %q", loaded.ProjectConfigPath)
			}
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, userConfigRelPath), "update:\n  mode: auto\n")
			loaded, err = Load(LoadOptions{HomeDir: home, WorkDir: repo, AllowNewerRepository: true})
			if err != nil || loaded.Config.Update.Mode != "auto" || loaded.RepositoryVersionErr == nil {
				t.Fatalf("User Config not retained: mode = %q, err = %v, refusal = %v", loaded.Config.Update.Mode, err, loaded.RepositoryVersionErr)
			}
		})
	}
	t.Run("malformed record", func(t *testing.T) {
		t.Parallel()
		home, repo := recordedVersionConfigRepository(t, `{"roundfixVersion":140}`, "future_section:\n  enabled: true\n")
		loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo, AllowNewerRepository: true})
		if err != nil || loaded.RepositoryVersionErr == nil || loaded.Config.Update.Mode != "notify" {
			t.Fatalf("mode = %q, err = %v, refusal = %v", loaded.Config.Update.Mode, err, loaded.RepositoryVersionErr)
		}
	})
}

func TestADevelopmentBuildIgnoresTheRecordedRoundfixVersion(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{`"` + futureRecordedRoundfixVersion + `"`, `"malformed"`} {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()
			_, repo := recordedVersionConfigRepository(t, `{"roundfixVersion":`+raw+`}`, "")
			for _, running := range []string{"0.0.0-dev", ""} {
				if err := CheckRecordedRoundfixVersion(repo, running); err != nil {
					t.Fatalf("development %q refused: %v", running, err)
				}
			}
			if err := CheckRecordedRoundfixVersion(repo, app.Version); err == nil {
				t.Fatal("running release accepted newer or malformed record")
			}
		})
	}
}

func recordedVersionConfigRepository(t *testing.T, manifest, project string) (home, repo string) {
	t.Helper()
	home, repo = t.TempDir(), t.TempDir()
	gittest.MarkRepository(t, repo)
	path := filepath.Join(repo, SetupManifestPath)
	mustMkdir(t, filepath.Dir(path))
	if manifest == "directory" {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	} else if manifest != "" {
		mustWrite(t, path, manifest)
	}
	mustWrite(t, filepath.Join(repo, projectConfigName), project)
	return home, repo
}
