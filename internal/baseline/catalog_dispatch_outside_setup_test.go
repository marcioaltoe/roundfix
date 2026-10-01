// Suite: profile dispatch membership
// Invariant: selected module dispatches and owned activation bundles name only setup skills.
// Boundary IN: cloned embedded catalogs and catalog diagnostics.
// Boundary OUT: skill installation, network, and managed refresh.
package baseline

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	for id, profile := range catalog.profiles {
		setupID, _ := stringValue(profile, "setup")
		var names []string
		for _, skill := range objectsOrEmpty(catalog.setups[setupID]["skills"]) {
			name, _ := stringValue(skill, "name")
			names = append(names, name)
		}
		for _, name := range profileNamedSkills(catalog, activations, stringsOrEmpty(profile["modules"])) {
			if !slices.Contains(names, name) {
				t.Errorf("profile %s setup %s omits %s", id, setupID, name)
			}
		}
	}
}

func TestADispatchedSkillOutsideTheProfileSetupIsReported(t *testing.T) {
	t.Parallel()
	assets := cloneEmbeddedAssets(t)
	const skill = "golang-testing"
	var setup document
	if err := json.Unmarshal(assets["setups/go-tui.json"].Data, &setup); err != nil {
		t.Fatal(err)
	}
	var retained []any
	removed := false
	for _, row := range objectsOrEmpty(setup["skills"]) {
		name, _ := stringValue(row, "name")
		if name == skill {
			removed = true
			continue
		}
		retained = append(retained, row)
	}
	if !removed {
		t.Fatal("fixture skill missing before removal")
	}
	setup["skills"] = retained
	data, err := json.Marshal(setup)
	if err != nil {
		t.Fatal(err)
	}
	assets["setups/go-tui.json"].Data = data
	_, err = LoadCatalog(assets)
	diagnostic := catalogDiagnosticByCode(t, err, "catalog.profile.skill.dispatch-outside-setup")
	if diagnostic.Path != "go-cli-tui" || diagnostic.Info != skill {
		t.Fatalf("diagnostic = %+v", diagnostic)
	}
}

func TestABundledSkillOutsideTheProfileSetupIsReported(t *testing.T) {
	t.Parallel()
	assets := cloneEmbeddedAssets(t)
	replaceAsset(t, assets, "skill-activations.json", `"skills": ["qa-gate", "evidence-gate"]`, `"skills": ["qa-gate", "evidence-gate", "unknown-qa-skill"]`)
	_, err := LoadCatalog(assets)
	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("catalog error = %v", err)
	}
	for _, profileID := range []string{"go-cli-tui", "rust-cli", "standard-typescript-monorepo"} {
		found := false
		for _, diagnostic := range validation.Diagnostics {
			if diagnostic.Code == "catalog.profile.skill.dispatch-outside-setup" && diagnostic.Path == profileID && diagnostic.Info == "unknown-qa-skill" {
				found = true
			}
		}
		if !found {
			t.Errorf("profile %s has no bundle diagnostic: %+v", profileID, validation.Diagnostics)
		}
	}
}

func TestAModuleTheProfileDoesNotSelectIsNotChecked(t *testing.T) {
	t.Parallel()
	assets := cloneEmbeddedAssets(t)
	catalog, err := LoadCatalog(assets)
	if err != nil {
		t.Fatal(err)
	}
	profile := catalog.profiles["go-cli-tui"]
	if slices.Contains(stringsOrEmpty(profile["modules"]), "typescript") {
		t.Fatal("fixture selects typescript")
	}
	setupID, _ := stringValue(profile, "setup")
	for _, row := range objectsOrEmpty(catalog.setups[setupID]["skills"]) {
		if row["name"] == "vitest" {
			t.Fatal("fixture setup lists vitest")
		}
	}
	if slices.Contains(profileNamedSkills(catalog, readUpstreamSkillDocument(t, "skill-activations.json"), stringsOrEmpty(profile["modules"])), "vitest") {
		t.Fatal("unselected TypeScript bundle was checked")
	}
}
