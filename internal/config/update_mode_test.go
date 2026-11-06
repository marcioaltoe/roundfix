// Suite: Update mode configuration.
// Invariant: only notify and auto are accepted, with Project Config winning.
// Boundary IN: layered YAML in temporary homes and repositories.
// Boundary OUT: no host configuration or external tools.
package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestUpdateModeIsNotifyOrAuto(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		user         string
		project      string
		want         string
		invalidScope string
	}{
		{name: "absent", want: "notify"},
		{name: "project auto", project: "auto", want: "auto"},
		{name: "user auto", user: "auto", want: "auto"},
		{name: "project over user", user: "notify", project: "auto", want: "auto"},
		{name: "project notify over user", user: "auto", project: "notify", want: "notify"},
		{name: "invalid project", user: "notify", invalidScope: "project"},
		{name: "invalid user", project: "auto", invalidScope: "user"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			check := func(t *testing.T, invalidMode string) {
				t.Helper()
				home, work := t.TempDir(), t.TempDir()
				gittest.MarkRepository(t, work)
				userPath := filepath.Join(home, ".roundfix", "config.yml")
				projectPath := filepath.Join(work, ".roundfixrc.yml")
				user, project := tc.user, tc.project
				invalidPath := ""
				switch tc.invalidScope {
				case "user":
					user, invalidPath = invalidMode, userPath
				case "project":
					project, invalidPath = invalidMode, projectPath
				}
				if user != "" || tc.invalidScope == "user" {
					mustMkdir(t, filepath.Dir(userPath))
					mustWrite(t, userPath, fmt.Sprintf("update:\n  mode: %q\n", user))
				}
				if project != "" || tc.invalidScope == "project" {
					mustWrite(t, projectPath, fmt.Sprintf("update:\n  mode: %q\n", project))
				}
				loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
				if invalidPath != "" {
					refusal := fmt.Sprintf("update.mode must be notify or auto, got %q", invalidMode)
					config := Builtin()
					config.Update.Mode = invalidMode
					if validateErr := Validate(config); validateErr == nil || validateErr.Error() != refusal {
						t.Errorf("Validate error = %v, want %s", validateErr, refusal)
					}
					want := fmt.Sprintf("parse config %q: %s", invalidPath, refusal)
					if err == nil || err.Error() != want {
						t.Fatalf("Load error = %v, want %s", err, want)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if loaded.Config.Update.Mode != tc.want {
					t.Fatalf("Update.Mode = %q, want %q", loaded.Config.Update.Mode, tc.want)
				}
			}
			if tc.invalidScope == "" {
				check(t, "")
				return
			}
			for _, mode := range []string{"always", "", "AUTO", " auto "} {
				t.Run(fmt.Sprintf("mode=%q", mode), func(t *testing.T) {
					check(t, mode)
				})
			}
		})
	}
	t.Run("templates unchanged", func(t *testing.T) {
		for _, template := range []struct {
			name    string
			content string
		}{
			{name: "user", content: DefaultConfigYAML()},
			{name: "project", content: DefaultProjectConfigYAML()},
		} {
			for _, line := range strings.Split(template.content, "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "update:") {
					t.Errorf("%s init template contains update mapping: %q", template.name, line)
				}
			}
		}
	})
}
