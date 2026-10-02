// Suite: Declared workspace locations
// Invariant: guides name only safe workspaces bound by built-in profiles.
// Boundary IN: embedded profiles and rendered managed guides.
// Boundary OUT: applying a Plan and executing repository Verification.

package baseline

import (
	"strings"
	"testing"
)

func TestTheBackendAndFrontendGuidesNameTheirDeclaredWorkspace(t *testing.T) {
	t.Parallel()
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
	for _, test := range []struct{ path, scope string }{
		{"docs/agents/backend.md", "These rules govern the repository's TypeScript backend workspace at `packages/backend`. A service or\ncommand written in another language follows its own guide."},
		{"docs/agents/frontend.md", "These rules govern the repository's web frontend workspace at `packages/frontend`. A terminal\ninterface follows its own guide."},
	} {
		t.Run(test.path, func(t *testing.T) {
			if text := string(planPostimage(t, plan, test.path).Content); !strings.Contains(text, test.scope) {
				t.Errorf("guide lacks scope %q; rendered guide:\n%s", test.scope, text)
			}
		})
	}
}

func TestAGuideWithoutADeclaredWorkspaceKeepsItsGenericScope(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name       string
		custom     bool
		workspaces []any
		location   string
	}{
		{name: "no binding"},
		{name: "repository-owned profile", custom: true},
		{name: "two sorted paths", workspaces: []any{
			map[string]any{"guide": "guide.backend", "path": "packages/z"},
			map[string]any{"guide": "guide.backend", "path": "packages/a"},
		}, location: " at `packages/a` and `packages/z`"},
		{name: "three sorted paths", workspaces: []any{
			map[string]any{"guide": "guide.backend", "path": "packages/z"},
			map[string]any{"guide": "guide.backend", "path": "packages/m"},
			map[string]any{"guide": "guide.backend", "path": "packages/a"},
		}, location: " at `packages/a`, `packages/m` and `packages/z`"},
		{name: "unbound and unsafe paths", workspaces: []any{
			map[string]any{"path": "packages/unbound"},
			map[string]any{"guide": "guide.backend", "path": "../outside"},
			map[string]any{"guide": "guide.backend", "path": "/absolute"},
			map[string]any{"guide": "guide.backend", "path": "C:/drive"},
			map[string]any{"guide": "guide.backend", "path": "packages/../backend"},
			map[string]any{"guide": "guide.backend", "path": "packages\\backend"},
			map[string]any{"guide": "guide.backend", "path": ""},
			map[string]any{"guide": "guide.backend", "path": "packages/safe"},
		}, location: " at `packages/safe`"},
	} {
		t.Run(test.name, func(t *testing.T) {
			catalog := cloneCatalogForRetentionDrift(t, mustEmbeddedCatalog(t))
			profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
			if err != nil {
				t.Fatal(err)
			}
			if test.custom {
				profile.ID = "repository-owned"
				profile.Source = ProfileSourceRepository
			} else {
				catalog.profiles[profile.ID]["workspaces"] = test.workspaces
			}
			_, artifacts, err := resolveManagedArtifacts(catalog, profile, standardTypeScriptDecisions("make verify"), false)
			if err != nil {
				t.Fatal(err)
			}
			found := 0
			for _, artifact := range artifacts {
				var scope string
				switch artifact.ID {
				case "guide.backend":
					scope = "These rules govern the repository's TypeScript backend workspace" + test.location + ". A service or\ncommand written in another language follows its own guide."
				case "guide.frontend":
					scope = "These rules govern the repository's web frontend workspace. A terminal\ninterface follows its own guide."
				default:
					continue
				}
				found++
				if !strings.Contains(artifact.Body, scope) {
					t.Errorf("%s lacks scope %q", artifact.ID, scope)
				}
			}
			if found != 2 {
				t.Fatalf("rendered %d workspace guides, want 2", found)
			}
		})
	}
}

func TestEveryBuiltInWorkspaceBindsAGuideItsProfileRenders(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for id, profile := range catalog.profiles {
		t.Run(id, func(t *testing.T) {
			artifacts := profileModuleArtifacts(stringsOrEmpty(profile["modules"]), catalog)
			for _, workspace := range objectsOrEmpty(profile["workspaces"]) {
				guide, ok := stringValue(workspace, "guide")
				if !ok || !strings.HasPrefix(guide, "guide.") {
					t.Errorf("workspace %v lacks a guide binding", workspace["path"])
				}
				if _, exists := artifacts[guide]; !exists {
					t.Errorf("workspace %v binds %q, which the profile does not render", workspace["path"], guide)
				}
			}
		})
	}
}
