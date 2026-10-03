package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/config"
)

func loadCursorConfig(t *testing.T, fragment string) (config.Loaded, error) {
	t.Helper()
	home, work := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".roundfix"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".roundfix", "config.yml"), []byte(fragment), 0600); err != nil {
		t.Fatal(err)
	}
	return config.Load(config.LoadOptions{HomeDir: home, WorkDir: work})
}

// Before the runtime change this fragment was refused at preferred.runtime.
func TestCursorProfileLoadsWithAnEmptyEffort(t *testing.T) {
	for _, effort := range []string{"", "  "} {
		t.Run(fmt.Sprintf("effort_%q", effort), func(t *testing.T) {
			loaded, err := loadCursorConfig(t, fmt.Sprintf(`profiles:
  backend:
    preferred: {runtime: cursor, model: " grok-4-20[thinking=true] ", reasoning_effort: %q}
    fallbacks:
      - {runtime: cursor, model: "default[]", reasoning_effort: ""}
`, effort))
			if err != nil {
				t.Fatal(err)
			}
			profile, err := config.ResolveProfile(loaded.Config, config.CategoryBackend, nil)
			if err != nil {
				t.Fatal(err)
			}
			want := config.AgentSelection{Runtime: "cursor", Model: "grok-4-20[thinking=true]"}
			if profile.Profile.Preferred != want {
				t.Fatalf("preferred = %+v, want %+v", profile.Profile.Preferred, want)
			}
			if profile.Profile.Fallbacks[0] != (config.AgentSelection{Runtime: "cursor", Model: "default[]"}) {
				t.Fatalf("fallback = %+v", profile.Profile.Fallbacks[0])
			}
		})
	}
}

func TestCursorProfileRefusesAReasoningEffort(t *testing.T) {
	for _, location := range []string{"preferred", "fallback"} {
		t.Run(location, func(t *testing.T) {
			preferred, fallback := `""`, `""`
			path := "profiles.backend.preferred"
			if location == "preferred" {
				preferred = `" high "`
			} else {
				fallback = `"high"`
				path = "profiles.backend.fallbacks[0]"
			}
			_, err := loadCursorConfig(t, fmt.Sprintf(`profiles:
  backend:
    preferred: {runtime: cursor, model: "grok-4-20[thinking=true]", reasoning_effort: %s}
    fallbacks:
      - {runtime: cursor, model: "default[]", reasoning_effort: %s}
`, preferred, fallback))
			want := path + `.reasoning_effort must be "" for runtime cursor; Cursor states effort, thinking, context and speed inside the model value, for example "grok-4-20[thinking=true]"`
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("error = %v, want %s", err, want)
			}
		})
	}
}

func TestSupportedRuntimesIsTheOneList(t *testing.T) {
	want := []string{"codex", "claude", "cursor", "opencode"}
	if got := config.SupportedRuntimes(); !reflect.DeepEqual(got, want) {
		t.Fatalf("runtimes = %v, want %v", got, want)
	}
	copy := config.SupportedRuntimes()
	copy[0] = "changed"
	if !reflect.DeepEqual(config.SupportedRuntimes(), want) {
		t.Fatal("list aliases caller memory")
	}
	for _, runtime := range config.SupportedRuntimes() {
		t.Run(runtime, func(t *testing.T) {
			_, err := loadCursorConfig(t, fmt.Sprintf(`profiles:
  backend:
    preferred: {runtime: %s, model: preferred, reasoning_effort: ""}
    fallbacks:
      - {runtime: %s, model: fallback, reasoning_effort: ""}
`, runtime, runtime))
			if err != nil {
				t.Fatalf("refused %s: %v", runtime, err)
			}
		})
	}
	for _, location := range []string{"preferred", "fallback"} {
		t.Run(location, func(t *testing.T) {
			preferred, fallback, path := "codex", "codex", "profiles.backend.preferred"
			if location == "preferred" {
				preferred = "unknown"
			} else {
				fallback = "unknown"
				path = "profiles.backend.fallbacks[0]"
			}
			_, err := loadCursorConfig(t, fmt.Sprintf(`profiles:
  backend:
    preferred: {runtime: %s, model: preferred, reasoning_effort: ""}
    fallbacks:
      - {runtime: %s, model: fallback, reasoning_effort: ""}
`, preferred, fallback))
			want := path + `.runtime "unknown" is invalid; supported values: ` + strings.Join(config.SupportedRuntimes(), ", ")
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("error = %v, want %s", err, want)
			}
		})
	}
}

func TestLegacyRuntimesRefuseCursor(t *testing.T) {
	for _, tc := range []struct{ name, fragment, want string }{
		{"defaults", "defaults:\n  agent: cursor\n", `defaults.agent "cursor" is invalid; supported values: codex, claude, opencode; name cursor in profiles`},
		{"runtimes", "runtimes:\n  cursor:\n    model: example\n", "runtimes.cursor is not a supported config key"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loadCursorConfig(t, tc.fragment)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %s", err, tc.want)
			}
		})
	}
	for _, runtime := range []string{"codex", "claude", "opencode"} {
		t.Run(runtime, func(t *testing.T) {
			_, err := loadCursorConfig(t, "defaults:\n  agent: "+runtime+"\nruntimes:\n  "+runtime+":\n    model: example\n    reasoning_effort: \"\"\n")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestModelCatalogRetainsOfficialCodexIdentifiers(t *testing.T) {
	t.Parallel()
	catalog := agent.ModelCatalog("codex")
	for _, identifier := range []string{"gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna"} {
		found := false
		for _, choice := range catalog {
			if choice.Value == identifier {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Codex Model Catalog is missing official identifier %q: %#v", identifier, catalog)
		}
	}
}
