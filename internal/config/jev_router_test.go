package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// The config loader owns scope and effort validation; no runtime is launched.
func jevRouterProfileYAML(fallback bool, effort string) string {
	router := fmt.Sprintf("{runtime: opencode, model: %s, reasoning_effort: %q}", JevRouterModel, effort)
	normal := `{runtime: codex, model: gpt-6.1-sol, reasoning_effort: high}`
	if fallback {
		router, normal = normal, router
	}
	return fmt.Sprintf("profiles:\n  docs:\n    preferred: %s\n    fallbacks: [%s]\n", router, normal)
}

func loadJevRouterConfig(t *testing.T, user, project string) (Loaded, error) {
	t.Helper()
	home, work := t.TempDir(), t.TempDir()
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustMkdir(t, filepath.Join(work, ".git"))
	if user != "" {
		mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), user)
	}
	if project != "" {
		mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), project)
	}
	return Load(LoadOptions{HomeDir: home, WorkDir: work})
}

func TestJevRouterLoadsFromProjectConfig(t *testing.T) {
	t.Parallel()
	for _, fallback := range []bool{false, true} {
		t.Run(fmt.Sprintf("fallback=%t", fallback), func(t *testing.T) {
			loaded, err := loadJevRouterConfig(t, "", jevRouterProfileYAML(fallback, ""))
			if err != nil {
				t.Fatal(err)
			}
			resolved, err := ResolveProfile(loaded.Config, CategoryDocs, nil)
			if err != nil {
				t.Fatal(err)
			}
			selection := resolved.Profile.Preferred
			if fallback {
				selection = resolved.Profile.Fallbacks[0]
			}
			if resolved.Source != ProfileSourceProject || selection.Model != JevRouterModel || selection.ReasoningEffort != "" {
				t.Fatalf("unexpected profile: %+v", resolved)
			}
		})
	}
}

func TestJevRouterRefusedFromUserConfig(t *testing.T) {
	t.Parallel()
	for _, fallback := range []bool{false, true} {
		t.Run(fmt.Sprintf("fallback=%t", fallback), func(t *testing.T) {
			for _, shadowed := range []bool{false, true} {
				t.Run(fmt.Sprintf("shadowed=%t", shadowed), func(t *testing.T) {
					project := ""
					if shadowed {
						project = `profiles:
  docs:
    preferred: {runtime: codex, model: gpt-6.1-sol, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: opus, reasoning_effort: high}]
`
					}
					_, err := loadJevRouterConfig(t, jevRouterProfileYAML(fallback, ""), project)
					path := "profiles.docs.preferred"
					if fallback {
						path = "profiles.docs.fallbacks[0]"
					}
					want := path + " names the Jev Router (" + JevRouterModel + "), which only Project Config may select"
					if err == nil || !strings.Contains(err.Error(), want) || !strings.Contains(err.Error(), "parse config") {
						t.Fatalf("error = %v, want %s", err, want)
					}
				})
			}
		})
	}
}

func TestJevRouterRefusesAReasoningEffort(t *testing.T) {
	t.Parallel()
	for _, fallback := range []bool{false, true} {
		t.Run(fmt.Sprintf("fallback=%t", fallback), func(t *testing.T) {
			_, err := loadJevRouterConfig(t, "", jevRouterProfileYAML(fallback, " high "))
			path := "profiles.docs.preferred"
			if fallback {
				path = "profiles.docs.fallbacks[0]"
			}
			want := path + `.reasoning_effort must be "" for the Jev Router; the router chooses the effort`
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("error = %v, want %s", err, want)
			}
		})
	}
}

func TestJevRouterRefusedAsOneRunOverride(t *testing.T) {
	t.Parallel()
	override := AgentSelection{Runtime: " opencode ", Model: " " + JevRouterModel + " "}
	_, err := ResolveProfile(Config{}, CategoryDocs, &override)
	want := "the Jev Router cannot be a one-Run override; name it in Project Config"
	if err == nil || err.Error() != want {
		t.Fatalf("error = %v, want %s", err, want)
	}
}

func TestJevRouterRefusedFromNonProjectSources(t *testing.T) {
	t.Parallel()
	for _, source := range []ProfileSource{ProfileSourceBuiltIn, ProfileSourceUser, ProfileSourceInvocation, ""} {
		t.Run(string(source), func(t *testing.T) {
			entries := builtinProfiles()
			entries[CategoryDocs] = ProfileEntry{Source: source, Profile: AgentSelectionProfile{
				Preferred: AgentSelection{Runtime: "opencode", Model: JevRouterModel},
				Fallbacks: []AgentSelection{{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"}},
			}}
			err := validateProfiles(entries)
			if err == nil || !strings.Contains(err.Error(), "which only Project Config may select") {
				t.Fatalf("source %q error = %v", source, err)
			}
		})
	}
}

func TestJevRouterRefusedFromLegacyRuntimeConfig(t *testing.T) {
	t.Parallel()
	for _, runtime := range []string{"opencode", "codex"} {
		t.Run(runtime, func(t *testing.T) {
			document := fmt.Sprintf(`defaults:
  agent: %s
runtimes:
  opencode:
    model: %s
    reasoning_effort: ""
`, runtime, JevRouterModel)
			_, err := loadJevRouterConfig(t, "", document)
			if err == nil || !strings.Contains(err.Error(), "the Jev Router cannot be selected through legacy runtimes; name it in Project Config profiles") {
				t.Fatalf("error = %v", err)
			}
		})
	}
}
