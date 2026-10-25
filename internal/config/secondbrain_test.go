// Suite: retired Secondbrain configuration.
// Invariant: retired keys warn once and are ignored; unknown keys remain errors.
// Boundary IN: temporary layered YAML files.
// Boundary OUT: no host configuration or external services.
package config

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheRetiredSecondbrainKeysLoadWithAWarning(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, user, project, want string }{
		{"User root", "secondbrain:\n  root: ~/brain\n", "", "root"},
		{"User invalid values", "secondbrain:\n  root: relative\n  inbox_namespace: Fluxus\n", "", "root,inbox_namespace"},
		{"Project namespace", "", "secondbrain:\n  inbox_namespace: _triaged\n", "inbox_namespace"},
		{"Project root", "", "secondbrain:\n  root: relative\n", "root"},
		{"Both layers", "secondbrain:\n  root: relative\n  inbox_namespace: Fluxus\n", "secondbrain:\n  root: /ignored\n  inbox_namespace: _triaged\n", "root,inbox_namespace"},
		{"Ignored value types", "secondbrain:\n  root: [invalid, path]\n  inbox_namespace: {invalid: value}\n", "", "root,inbox_namespace"},
		{"Empty mapping", "secondbrain: {}\n", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, repo := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustMkdir(t, filepath.Join(repo, ".git"))
			mustWrite(t, filepath.Join(home, ".roundfix/config.yml"), tc.user)
			mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), tc.project)
			var warnings bytes.Buffer
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: repo, Stderr: &warnings})
			if err != nil {
				t.Fatal(err)
			}
			want := ""
			if tc.want != "" {
				for _, key := range strings.Split(tc.want, ",") {
					want += "config: secondbrain." + key + " is deprecated and ignored; the Secondbrain inbox source was retired, so remove it\n"
				}
			}
			if warnings.String() != want {
				t.Fatalf("warnings=%q want=%q", warnings.String(), want)
			}
		})
	}
}

func TestAnUnknownSecondbrainKeyStillFailsStrictValidation(t *testing.T) {
	t.Parallel()
	for _, layer := range []string{"User", "Project"} {
		t.Run(layer, func(t *testing.T) {
			t.Parallel()
			home, repo := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustMkdir(t, filepath.Join(repo, ".git"))
			path := filepath.Join(home, ".roundfix/config.yml")
			if layer == "Project" {
				path = filepath.Join(repo, ".roundfixrc.yml")
			}
			mustWrite(t, path, "secondbrain:\n  root: relative\n  inbox_namespace: Fluxus\n  unknown: true\n")
			var warnings bytes.Buffer
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: repo, Stderr: &warnings})
			if err == nil || !strings.Contains(err.Error(), "field secondbrain not found") {
				t.Fatalf("error=%v", err)
			}
		})
	}
}
