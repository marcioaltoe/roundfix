// Suite: Skill Activation ownership and required-skill composition.
// Invariant: an owner and its dependencies require every skill the rendered dispatch names.
// Boundary IN: embedded catalog, copied asset mutations, module compositions and rendered guide text.
// Boundary OUT: network, installed skills, host home and CLI output.
package baseline

import (
	"errors"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestAnActivationOwnerRequiresEverySkillOfItsBundle(t *testing.T) {
	t.Parallel()
	t.Run("the_embedded_catalog_loads", func(t *testing.T) {
		if _, err := LoadEmbeddedCatalog(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("a_sibling_module_does_not_satisfy_the_owner", func(t *testing.T) {
		assets := cloneEmbeddedAssets(t)
		replaceAsset(t, assets, "modules/typescript.json", "    \"solid\",\n", "")
		if !strings.Contains(string(assets["modules/frontend.json"].Data), "    \"solid\",\n") {
			t.Fatal("sibling frontend no longer requires solid")
		}
		_, err := LoadCatalog(assets)
		want := Diagnostic{Code: "catalog.skill.activation.skill.unrequired", Path: "trigger.production-code", Info: "solid"}
		if got := catalogDiagnosticByCode(t, err, want.Code); got != want {
			t.Fatalf("diagnostic = %+v, want %+v", got, want)
		}
		if got := activationRequirementDiagnostics(t, err); !reflect.DeepEqual(got, []Diagnostic{want}) {
			t.Fatalf("requirement diagnostics = %+v, want only %+v", got, want)
		}
	})
	t.Run("a_dependency_satisfies_the_owner", func(t *testing.T) {
		assets := cloneEmbeddedAssets(t)
		if _, err := LoadCatalog(assets); err != nil {
			t.Fatalf("owner with core dependency: %v", err)
		}
		replaceAsset(t, assets, "modules/core.json", "    \"systematic-debugging\",\n", "")
		_, err := LoadCatalog(assets)
		want := Diagnostic{Code: "catalog.skill.activation.skill.unrequired", Path: "trigger.debugging", Info: "systematic-debugging"}
		if got := catalogDiagnosticByCode(t, err, want.Code); got != want {
			t.Fatalf("diagnostic = %+v, want %+v", got, want)
		}
		if got := activationRequirementDiagnostics(t, err); !reflect.DeepEqual(got, []Diagnostic{want}) {
			t.Fatalf("requirement diagnostics = %+v, want only %+v", got, want)
		}
	})
	for _, field := range []string{"owner", "bundle"} {
		name := "an_unknown_owner_reports_only_the_owner"
		if field == "bundle" {
			name = "an_unknown_bundle_reports_only_the_bundle"
		}
		t.Run(name, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			old := "      \"owner\": \"typescript\",\n      \"when\": \"Diagnosing or fixing unexpected behavior.\","
			replacement := strings.Replace(old, "typescript", "unknown-owner", 1)
			if field == "bundle" {
				old = "      \"bundle\": \"bundle.debugging\""
				replacement = "      \"bundle\": \"bundle.unknown\""
			}
			replaceAsset(t, assets, "skill-activations.json", old, replacement)
			_, err := LoadCatalog(assets)
			code := "catalog.skill.activation." + field + ".unknown"
			diagnostic := catalogDiagnosticByCode(t, err, code)
			if diagnostic.Path != "trigger.debugging" {
				t.Fatalf("unknown %s diagnostic = %+v", field, diagnostic)
			}
			for _, diagnostic := range activationRequirementDiagnostics(t, err) {
				if diagnostic.Path == "trigger.debugging" {
					t.Fatalf("unknown %s also reported unrequired skill: %+v", field, diagnostic)
				}
			}
		})
	}
	t.Run("dependency_cycles_do_not_recurse_forever", func(t *testing.T) {
		assets := cloneEmbeddedAssets(t)
		replaceAsset(t, assets, "modules/typescript.json", "    \"solid\",\n", "")
		replaceAsset(t, assets, "modules/typescript.json", "    \"core\"\n", "    \"core\",\n    \"frontend\"\n")
		_, err := LoadCatalog(assets)
		catalogDiagnosticByCode(t, err, "catalog.module.dependency.cycle")
		if got := activationRequirementDiagnostics(t, err); len(got) != 0 {
			t.Fatalf("frontend dependency in a cycle did not supply solid: %+v", got)
		}
	})
	t.Run("diagnostics_follow_activation_and_bundle_order", func(t *testing.T) {
		assets := cloneEmbeddedAssets(t)
		// Put missing skills in two bundles whose activation and skill order
		// both differ from lexical order. Assert only the rule under test.
		replaceAsset(t, assets, "skill-activations.json", `"skills": ["coding-guidelines", "clean-code", "solid"]`, `"skills": ["missing-z", "missing-a"]`)
		replaceAsset(t, assets, "skill-activations.json", `"skills": ["testing-boss", "vitest"]`, `"skills": ["missing-y", "missing-b"]`)
		replaceAsset(t, assets, "skill-activations.json", `"id": "trigger.production-code"`, `"id": "trigger.z-production-code"`)
		_, err := LoadCatalog(assets)
		want := []Diagnostic{
			{Code: "catalog.skill.activation.skill.unrequired", Path: "trigger.z-production-code", Info: "missing-z"},
			{Code: "catalog.skill.activation.skill.unrequired", Path: "trigger.z-production-code", Info: "missing-a"},
			{Code: "catalog.skill.activation.skill.unrequired", Path: "trigger.testing", Info: "missing-y"},
			{Code: "catalog.skill.activation.skill.unrequired", Path: "trigger.testing", Info: "missing-b"},
		}
		if got := activationRequirementDiagnostics(t, err); !reflect.DeepEqual(got, want) {
			t.Fatalf("requirement diagnostics = %+v, want authored order %+v", got, want)
		}
	})
}

// Copied module mutations can also invalidate unrelated pins. Read only the
// activation requirement diagnostics so the test witnesses this rule alone.
func activationRequirementDiagnostics(t *testing.T, err error) []Diagnostic {
	t.Helper()
	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("catalog error = %v, want *ValidationError", err)
	}
	var diagnostics []Diagnostic
	for _, diagnostic := range validation.Diagnostics {
		if diagnostic.Code == "catalog.skill.activation.skill.unrequired" {
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	return diagnostics
}

func TestTheRequiredSkillsCoverTheRenderedDispatch(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("every_module_with_its_dependencies", func(t *testing.T) {
		for _, moduleID := range catalog.ModuleIDs() {
			t.Run(moduleID, func(t *testing.T) {
				visited := make(map[string]struct{})
				var visit func(string)
				visit = func(id string) {
					if _, seen := visited[id]; seen {
						return
					}
					visited[id] = struct{}{}
					for _, dependency := range stringsOrEmpty(catalog.modules[id]["dependsOn"]) {
						visit(dependency)
					}
				}
				visit(moduleID)
				assertRenderedSkillsAreRequired(t, catalog, sortedKeys(visited))
			})
		}
	})
	t.Run("every_built_in_profile", func(t *testing.T) {
		for _, profileID := range catalog.ProfileIDs() {
			t.Run(profileID, func(t *testing.T) {
				modules, ok := catalog.OrderedModules(profileID)
				if !ok {
					t.Fatalf("profile %q has no resolved modules", profileID)
				}
				assertRenderedSkillsAreRequired(t, catalog, modules)
			})
		}
	})
	t.Run("the_built_in_profiles_keep_their_names", func(t *testing.T) {
		for _, test := range []struct {
			id    string
			count int
			solid bool
		}{
			{"go-cli-tui", 41, false},
			{"go-cli-typescript-monorepo", 105, true},
			{"python", 38, false},
			{"rust-cli", 34, false},
			{"standard-typescript-monorepo", 95, true},
		} {
			t.Run(test.id, func(t *testing.T) {
				modules, ok := catalog.OrderedModules(test.id)
				if !ok {
					t.Fatalf("profile %q has no resolved modules", test.id)
				}
				required, err := catalog.RequiredSkills(modules)
				if err != nil {
					t.Fatal(err)
				}
				if len(required) != test.count || slices.Contains(required, "solid") != test.solid {
					t.Fatalf("required skills = %v; want %d names, solid=%t", required, test.count, test.solid)
				}
			})
		}
	})
}

func assertRenderedSkillsAreRequired(t *testing.T, catalog *Catalog, modules []string) {
	t.Helper()
	required, err := catalog.RequiredSkills(modules)
	if err != nil {
		t.Fatal(err)
	}
	// Parse the text an Agent reads, independently of profileNamedSkills and
	// the activation data structure used by the catalog validator.
	bundleLine := regexp.MustCompile("^  [-*+] `bundle\\.[^`]+`: (.*)$")
	individualLine := regexp.MustCompile("^[-*+] `([^`]+)`:$")
	token := regexp.MustCompile("`([^`]+)`")
	var names []string
	for _, line := range strings.Split(renderSkillDispatch(catalog, modules), "\n") {
		if match := bundleLine.FindStringSubmatch(line); match != nil {
			for _, skill := range token.FindAllStringSubmatch(match[1], -1) {
				names = append(names, skill[1])
			}
		} else if match := individualLine.FindStringSubmatch(line); match != nil {
			names = append(names, match[1])
		}
	}
	if len(names) == 0 {
		t.Fatal("rendered dispatch has no parsed skill names")
	}
	for _, name := range names {
		if !slices.Contains(required, name) {
			t.Errorf("rendered skill %q is absent from required skills for modules %v", name, modules)
		}
	}
}

func TestRequiredSkills(t *testing.T) {
	t.Parallel()
	catalog := &Catalog{modules: map[string]document{
		"owner":      {"requiredSkills": []any{"zeta", "alpha"}, "dependsOn": []any{"dependency"}},
		"sibling":    {"requiredSkills": []any{"alpha", "beta"}},
		"dependency": {"requiredSkills": []any{"dependency-only"}},
	}}
	for _, test := range []struct {
		name    string
		modules []string
		want    []string
		unknown string
	}{
		{name: "no_modules", want: []string{}},
		{name: "sorted_union_with_repeated_modules", modules: []string{"sibling", "owner", "sibling"}, want: []string{"alpha", "beta", "zeta"}},
		{name: "dependencies_are_not_implicit", modules: []string{"owner"}, want: []string{"alpha", "zeta"}},
		{name: "unknown_module_returns_no_partial_list", modules: []string{"owner", "unknown"}, unknown: "unknown"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := catalog.RequiredSkills(test.modules)
			if test.unknown != "" {
				var unknown *UnknownModuleError
				if got != nil || !errors.As(err, &unknown) || unknown.ModuleID != test.unknown || !strings.Contains(err.Error(), test.unknown) {
					t.Fatalf("required skills = %v, error = %v; want no list and unknown module %q", got, err, test.unknown)
				}
				return
			}
			if err != nil || !reflect.DeepEqual(got, test.want) {
				t.Fatalf("required skills = %v, error = %v; want %v", got, err, test.want)
			}
			if len(got) != 0 {
				got[0] = "caller-mutation"
				again, err := catalog.RequiredSkills(test.modules)
				if err != nil || !reflect.DeepEqual(again, test.want) {
					t.Fatalf("caller changed catalog: skills = %v, error = %v", again, err)
				}
			}
		})
	}
}
