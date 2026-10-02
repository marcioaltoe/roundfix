// Suite: composed Setup Snapshots
// Invariant: a materialized composition is the ordered union of named,
// non-composed components, and sync validates it before writing any assets.
// Boundary IN: embedded catalog overlays and private Git/asset fixtures.
// Boundary OUT: profiles, managed refresh, and network access.
package baseline

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

const composedSetupID = "go-cli-typescript-bun"

func TestTheComposedSetupIsTheUnionOfItsComponents(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	composed := catalog.setups[composedSetupID]
	source, _ := objectValue(composed["source"])
	if source["type"] != "composed" || !reflect.DeepEqual(stringsOrEmpty(source["setups"]), []string{"go", "typescript"}) {
		t.Fatalf("composition source = %v", source)
	}
	for _, field := range []struct{ name, key string }{{"skills", "name"}, {"activationBundles", "id"}} {
		seen := map[string]document{}
		want := []document{}
		for _, id := range []string{"go", "typescript"} {
			for _, entry := range objectsOrEmpty(catalog.setups[id][field.name]) {
				name, _ := stringValue(entry, field.key)
				if prior, exists := seen[name]; exists {
					if !reflect.DeepEqual(prior, entry) {
						t.Fatalf("components disagree on %s %s", field.name, name)
					}
					continue
				}
				seen[name] = entry
				want = append(want, entry)
			}
		}
		if got := objectsOrEmpty(composed[field.name]); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s is not the component union in component order", field.name)
		}
	}
	// Small explicit entries exercise equal deduplication without depending on
	// the current upstream membership or a second call to the implementation.
	components := []document{
		{"id": "first", "skills": []any{map[string]any{"name": "a"}, map[string]any{"name": "b"}}, "activationBundles": []any{map[string]any{"id": "shared", "skills": []any{"a"}}}},
		{"id": "second", "skills": []any{map[string]any{"name": "b"}, map[string]any{"name": "c"}}, "activationBundles": []any{map[string]any{"id": "shared", "skills": []any{"a"}}, map[string]any{"id": "last", "skills": []any{"c"}}}},
	}
	result, err := composeSetupSnapshot("example", components)
	if err != nil {
		t.Fatal(err)
	}
	var names, bundles []string
	for _, entry := range objectsOrEmpty(result["skills"]) {
		names = append(names, entry["name"].(string))
	}
	for _, entry := range objectsOrEmpty(result["activationBundles"]) {
		bundles = append(bundles, entry["id"].(string))
	}
	if !slices.Equal(names, []string{"a", "b", "c"}) || !slices.Equal(bundles, []string{"shared", "last"}) {
		t.Fatalf("ordered union skills=%v bundles=%v", names, bundles)
	}
}

func TestAComposedSetupThatDriftsFromItsComponentsIsRefused(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"skills", "activationBundles"} {
		t.Run(field, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			setup := compositionFixtureDocument(t, assets, composedSetupID)
			entries := setup[field].([]any)
			setup[field] = entries[1:]
			writeCompositionFixture(t, assets, composedSetupID, setup)
			_, err := LoadCatalog(assets)
			requireCatalogDiagnostic(t, err, "catalog.setup.composition.drift")
		})
	}
}

func TestAComposedSetupWithAnUnknownOrComposedComponentIsRefused(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		ids  any
	}{
		{"unknown", []any{"go", "missing"}},
		{"composed", []any{"go", composedSetupID}},
		{"one component", []any{"go"}},
		{"duplicate", []any{"go", "go"}},
		{"malformed list", "go"},
		{"malformed identifier", []any{"go", 42}},
	} {
		t.Run(test.name, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			setup := compositionFixtureDocument(t, assets, composedSetupID)
			source, _ := objectValue(setup["source"])
			source["setups"] = test.ids
			writeCompositionFixture(t, assets, composedSetupID, setup)
			_, err := LoadCatalog(assets)
			requireCatalogDiagnostic(t, err, "catalog.setup.composition.invalid")
		})
	}
}

func TestComponentsThatDisagreeOnASkillAreAConflict(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"skills", "activationBundles"} {
		t.Run(field, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			component := compositionFixtureDocument(t, assets, "typescript")
			var entry document
			for _, candidate := range objectsOrEmpty(component[field]) {
				if candidate["name"] == "knowledge-workspace" || candidate["id"] == "bundle.delivery" {
					entry = candidate
					break
				}
			}
			if entry == nil {
				t.Fatal("fixture shared entry missing")
			}
			if field == "skills" {
				entry["path"] = "skills/conflicting-path"
			} else {
				entry["skills"] = []any{"roundfix"}
			}
			writeCompositionFixture(t, assets, "typescript", component)
			_, err := LoadCatalog(assets)
			requireCatalogDiagnostic(t, err, "catalog.setup.composition.conflict")
		})
	}
	t.Run("sync conflict writes nothing", func(t *testing.T) {
		repository, root := newAssetsSyncTarget(t)
		source, _ := newAssetsSyncSource(t, root)
		var component assetsSyncSnapshot
		readAssetsSyncJSON(t, filepath.Join(root, "setups", "typescript.json"), &component)
		for index := range component.ActivationBundles {
			if component.ActivationBundles[index].ID == "bundle.delivery" {
				component.ActivationBundles[index].Skills = []string{"roundfix"}
			}
		}
		writeAssetsSyncJSON(t, filepath.Join(root, "setups", "typescript.json"), component)
		before := captureAssetsSyncTree(t, root)
		payload, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source}, assetsSyncDependencies{repository: repository, assetRoot: root})
		var syncErr *AssetsSyncError
		if !errors.As(err, &syncErr) || syncErr.Category != AssetsSyncInvalid || payload.OK {
			t.Fatalf("conflict result = %+v, error = %v", payload, err)
		}
		found := false
		for _, finding := range payload.Findings {
			found = found || (finding.Code == "skills.setup-snapshot.drift" && strings.Contains(finding.Message, "disagree"))
		}
		if !found {
			t.Fatalf("missing invalid-assets conflict finding: %+v", payload.Findings)
		}
		if !reflect.DeepEqual(before, captureAssetsSyncTree(t, root)) {
			t.Fatal("composition conflict wrote assets")
		}
	})
}

func TestAnAssetSyncRewritesTheComposedSetupWithItsComponents(t *testing.T) {
	t.Parallel()
	repository, root := newAssetsSyncTarget(t)
	source, _ := newAssetsSyncSource(t, root)
	dependencies := assetsSyncDependencies{repository: repository, assetRoot: root}
	if _, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source}, dependencies); err != nil {
		t.Fatal(err)
	}
	before := captureAssetsSyncTree(t, root)
	var goSetup assetsSyncSnapshot
	readAssetsSyncJSON(t, filepath.Join(root, "setups", "go.json"), &goSetup)
	var uniquePath string
	for _, skill := range goSetup.Skills {
		if skill.Name == "golang-testing" {
			uniquePath = skill.Path
		}
	}
	if uniquePath == "" {
		t.Fatal("fixture has no golang-testing")
	}
	target := filepath.Join(filepath.Dir(source), filepath.FromSlash(uniquePath), "SKILL.md")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, append(data, []byte("\nchanged component content\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	runAssetsSyncGit(t, filepath.Dir(source), "add", ".")
	runAssetsSyncGit(t, filepath.Dir(source), "commit", "--quiet", "-m", "change component content")
	payload, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source, Check: true}, dependencies)
	if err == nil || payload.OK {
		t.Fatalf("check missed drift: %+v, %v", payload, err)
	}
	found := false
	for _, finding := range payload.Findings {
		found = found || finding.ManagedID == "setup."+composedSetupID
	}
	if !found {
		t.Fatalf("check missed composed drift: %+v", payload.Findings)
	}
	if !reflect.DeepEqual(before, captureAssetsSyncTree(t, root)) {
		t.Fatal("check wrote assets")
	}
	if _, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source}, dependencies); err != nil {
		t.Fatal(err)
	}
	after := captureAssetsSyncTree(t, root)
	for _, id := range []string{"go", composedSetupID} {
		if after["setups/"+id+".json"] == before["setups/"+id+".json"] {
			t.Fatalf("%s was not rewritten in the component sync", id)
		}
	}
	catalog, err := LoadCatalog(os.DirFS(root))
	if err != nil {
		t.Fatal(err)
	}
	var componentSkill, composedSkill document
	for _, skill := range objectsOrEmpty(catalog.setups["go"]["skills"]) {
		if skill["name"] == "golang-testing" {
			componentSkill = skill
		}
	}
	for _, skill := range objectsOrEmpty(catalog.setups[composedSetupID]["skills"]) {
		if skill["name"] == "golang-testing" {
			composedSkill = skill
		}
	}
	if componentSkill == nil || !reflect.DeepEqual(componentSkill, composedSkill) {
		t.Fatal("composed snapshot did not inherit the changed skill")
	}
}

func TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup(t *testing.T) {
	t.Parallel()
	repository, root := newAssetsSyncTarget(t)
	source, _ := newAssetsSyncSource(t, root)
	dependencies := assetsSyncDependencies{repository: repository, assetRoot: root}
	if _, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source}, dependencies); err != nil {
		t.Fatal(err)
	}
	before := captureAssetsSyncTree(t, root)
	payload, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source}, dependencies)
	if err != nil || !payload.OK || len(payload.Findings) != 0 {
		t.Fatalf("unchanged sync = %+v, %v", payload, err)
	}
	if !reflect.DeepEqual(before, captureAssetsSyncTree(t, root)) {
		t.Fatal("unchanged sync changed snapshot bytes")
	}
	// Drift in the materialized composition alone is also visible to --check.
	var composed assetsSyncSnapshot
	target := filepath.Join(root, "setups", composedSetupID+".json")
	readAssetsSyncJSON(t, target, &composed)
	composed.Skills = composed.Skills[1:]
	writeAssetsSyncJSON(t, target, composed)
	drifted := captureAssetsSyncTree(t, root)
	payload, err = syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source, Check: true}, dependencies)
	if err == nil || len(payload.Findings) != 1 || payload.Findings[0].ManagedID != "setup."+composedSetupID {
		t.Fatalf("composed-only check = %+v, %v", payload, err)
	}
	if !reflect.DeepEqual(drifted, captureAssetsSyncTree(t, root)) {
		t.Fatal("composed-only check wrote assets")
	}
	if _, err := syncAssets(t.Context(), AssetsSyncRequest{SourceDir: source}, dependencies); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(before["setups/"+composedSetupID+".json"]), []byte(captureAssetsSyncTree(t, root)["setups/"+composedSetupID+".json"])) {
		t.Fatal("composed-only repair did not restore exact bytes")
	}
}

func compositionFixtureDocument(t *testing.T, assets fstest.MapFS, id string) document {
	t.Helper()
	doc, diagnostics := decodeDocument(assets["setups/"+id+".json"].Data, id)
	if len(diagnostics) != 0 {
		t.Fatal(diagnostics)
	}
	return doc
}

func writeCompositionFixture(t *testing.T, assets fstest.MapFS, id string, setup document) {
	t.Helper()
	var payload any = setup["skills"]
	if bundles, exists := setup["activationBundles"]; exists {
		payload = document{"skills": setup["skills"], "activationBundles": bundles}
	}
	digest, err := canonicalSHA256(payload)
	if err != nil {
		t.Fatal(err)
	}
	setup["digest"] = digest
	data, err := json.MarshalIndent(setup, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	assets["setups/"+id+".json"].Data = append(data, '\n')
}
