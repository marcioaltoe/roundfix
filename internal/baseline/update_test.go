// Suite: Setup Manifest plan inputs
// Invariant: manifest resolution returns typed, read-only plan inputs without adopting catalog suggestions.
// Boundary IN: manifest parsing, profile resolution, decision validation, and new-decision discovery.
// Boundary OUT: Baseline plan assembly, CLI rendering, prompting, and repository mutation.

package baseline

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestResolveManifestInputAbsent(t *testing.T) {
	t.Parallel()

	catalog := loadManifestInputCatalog(t)
	input, err := ResolveManifestInput(t.TempDir(), catalog)
	if !errors.Is(err, ErrNoManifest) {
		t.Fatalf("ResolveManifestInput() error = %v, want ErrNoManifest", err)
	}
	if input.State != ManifestInputAbsent {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputAbsent)
	}
}

func TestResolveManifestInputCurrent(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "go-cli-tui")
	manifest := newManifestInputFixture(t, repository, catalog, "go-cli-tui", decisions)
	writeManifestInputFixture(t, repository, manifest)

	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	if input.State != ManifestInputResolved {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputResolved)
	}
	if input.ProfileID != manifest.Profile {
		t.Fatalf("ResolveManifestInput() profile = %q, want %q", input.ProfileID, manifest.Profile)
	}
	if input.ProfileDigestChanged {
		t.Fatal("ResolveManifestInput() reported profile digest drift for a current manifest")
	}
	assertManifestInputDecisionsEqual(t, input.Decisions, manifest.Decisions)
	if len(input.NewDecisions) != 0 {
		t.Fatalf("ResolveManifestInput() new decisions = %#v, want none", input.NewDecisions)
	}
}

func TestResolveManifestInputIncompleteDoesNotAdoptSuggestion(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "go-cli-tui")
	delete(decisions, "secondbrain.enabled")
	manifest := newManifestInputFixture(t, repository, catalog, "go-cli-tui", decisions)
	writeManifestInputFixture(t, repository, manifest)

	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	if input.State != ManifestInputIncomplete {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputIncomplete)
	}
	if len(input.NewDecisions) != 1 || input.NewDecisions[0].ID != "secondbrain.enabled" {
		t.Fatalf("ResolveManifestInput() new decisions = %#v, want secondbrain.enabled", input.NewDecisions)
	}
	if input.NewDecisions[0].SuggestedValue != true {
		t.Fatalf("ResolveManifestInput() suggestion = %#v, want true", input.NewDecisions[0].SuggestedValue)
	}
	if input.NewDecisions[0].Summary == "" {
		t.Fatal("ResolveManifestInput() suggestion summary is empty")
	}
	if decisionValueIndex(input.Decisions, "secondbrain.enabled") >= 0 {
		t.Fatal("ResolveManifestInput() adopted a suggestion into the resolved decisions")
	}
}

func TestResolveManifestInputNamesEveryNewDecision(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "go-cli-tui")
	delete(decisions, "domain.layout")
	delete(decisions, "triage.external")
	manifest := newManifestInputFixture(t, repository, catalog, "go-cli-tui", decisions)
	writeManifestInputFixture(t, repository, manifest)

	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	if input.State != ManifestInputIncomplete {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputIncomplete)
	}
	want := []string{"domain.layout", "triage.external"}
	if len(input.NewDecisions) != len(want) {
		t.Fatalf("ResolveManifestInput() new decisions = %#v, want %v", input.NewDecisions, want)
	}
	for index, id := range want {
		if input.NewDecisions[index].ID != id || input.NewDecisions[index].Summary == "" {
			t.Fatalf("ResolveManifestInput() new decision %d = %#v, want %q with summary", index, input.NewDecisions[index], id)
		}
	}
}

func TestResolveManifestInput_standard_typescript_monorepoStructuredSuggestion(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	decisions := manifestInputDecisionsForProfile(
		t,
		repository,
		catalog,
		"standard-typescript-monorepo",
	)
	delete(decisions, authProviderDecisionID)
	manifest := newManifestInputFixture(
		t,
		repository,
		catalog,
		"standard-typescript-monorepo",
		decisions,
	)
	writeManifestInputFixture(t, repository, manifest)

	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	if input.State != ManifestInputIncomplete {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputIncomplete)
	}
	if len(input.NewDecisions) != 1 || input.NewDecisions[0].ID != authProviderDecisionID {
		t.Fatalf("ResolveManifestInput() new decisions = %#v, want %q", input.NewDecisions, authProviderDecisionID)
	}
	if err := ValidateDecisionValue(catalog, authProviderDecisionID, input.NewDecisions[0].SuggestedValue); err != nil {
		t.Fatalf("structured TypeScript suggestion is invalid: %v", err)
	}
	if decisionValueIndex(input.Decisions, authProviderDecisionID) >= 0 {
		t.Fatal("ResolveManifestInput() adopted the structured TypeScript suggestion")
	}
}

func TestResolveManifestInputProfileDigestDriftIsResolved(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "go-cli-tui")
	manifest := newManifestInputFixture(t, repository, catalog, "go-cli-tui", decisions)
	manifest.ProfileDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	writeManifestInputFixture(t, repository, manifest)

	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	if input.State != ManifestInputResolved {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputResolved)
	}
	if !input.ProfileDigestChanged {
		t.Fatal("ResolveManifestInput() did not report the profile catalog move")
	}
	assertManifestInputDecisionsEqual(t, input.Decisions, manifest.Decisions)
}

func TestResolveManifestInputUnresolvedProfileDiagnosis(t *testing.T) {
	t.Parallel()

	catalog := loadManifestInputCatalog(t)
	tests := []struct {
		name          string
		profileID     string
		wantKind      UnresolvedProfileKind
		wantLocations []string
		wantAction    string
	}{
		{
			name:          "missing repository-owned Profile",
			profileID:     "repository-backend",
			wantKind:      UnresolvedProfileRepositoryMissing,
			wantLocations: []string{".roundfix/baseline/profiles/repository-backend.json"},
			wantAction:    "restore",
		},
		{
			name:          "unknown catalog identity",
			profileID:     "retired/profile",
			wantKind:      UnresolvedProfileCatalogUnknown,
			wantLocations: []string{},
			wantAction:    "adopt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := t.TempDir()
			manifest := SetupManifest{
				SchemaVersion: ManifestSchema,
				Version:       ManifestVersion,
				Generator: ManifestGenerator{
					Skill:    "setup-context-driven",
					Version:  ManifestVersion,
					Baseline: "baseline." + test.profileID + "-" + ManifestVersion,
				},
				Profile:          test.profileID,
				ProfileDigest:    "sha256:" + strings.Repeat("0", 64),
				CatalogDigest:    catalog.Digest(),
				Modules:          []string{},
				Decisions:        map[string]ManifestDecision{},
				ManagedArtifacts: []ManifestArtifact{},
				LocalSkills:      []string{},
				Verification:     []VerificationProjection{},
			}
			writeManifestInputFixture(t, repository, manifest)

			input, err := ResolveManifestInput(repository, catalog)
			if !errors.Is(err, ErrManifestIncompatible) {
				t.Fatalf("ResolveManifestInput() error = %v, want ErrManifestIncompatible", err)
			}
			if input.State != ManifestInputIncompatible ||
				input.Incompatibility != ManifestInputProfileUnresolved {
				t.Fatalf("ResolveManifestInput() state = %q reason = %q", input.State, input.Incompatibility)
			}
			if input.UnresolvedProfile == nil {
				t.Fatal("ResolveManifestInput() unresolved Profile diagnosis is nil")
			}
			diagnosis := *input.UnresolvedProfile
			if diagnosis.Identity != test.profileID || diagnosis.Kind != test.wantKind ||
				!reflect.DeepEqual(diagnosis.SearchedLocations, test.wantLocations) ||
				!strings.Contains(strings.ToLower(diagnosis.Action), test.wantAction) {
				t.Fatalf("ResolveManifestInput() unresolved Profile diagnosis = %+v", diagnosis)
			}
			message := err.Error()
			if !strings.Contains(message, test.profileID) ||
				!strings.Contains(message, diagnosis.Action) ||
				strings.Contains(message, "lstat") || strings.Contains(message, "open ") {
				t.Fatalf("ResolveManifestInput() diagnosis message = %q", message)
			}
			for _, location := range test.wantLocations {
				if !strings.Contains(message, location) {
					t.Fatalf("ResolveManifestInput() diagnosis message %q lacks %q", message, location)
				}
			}
		})
	}
}

func TestResolveManifestInputDistinguishesAdoptionReasons(t *testing.T) {
	t.Parallel()

	catalog := loadManifestInputCatalog(t)
	tests := []struct {
		name       string
		manifest   func(*testing.T, string) SetupManifest
		wantReason ManifestInputIncompatibility
	}{
		{
			name: "profile no longer resolves",
			manifest: func(*testing.T, string) SetupManifest {
				return SetupManifest{
					SchemaVersion: ManifestSchema,
					Version:       ManifestVersion,
					Generator: ManifestGenerator{
						Skill:    "setup-context-driven",
						Version:  ManifestVersion,
						Baseline: "baseline.retired-profile-" + ManifestVersion,
					},
					Profile:          "retired-profile",
					ProfileDigest:    "sha256:0000000000000000000000000000000000000000000000000000000000000000",
					CatalogDigest:    catalog.Digest(),
					Modules:          []string{},
					Decisions:        map[string]ManifestDecision{},
					ManagedArtifacts: []ManifestArtifact{},
					LocalSkills:      []string{},
					Verification:     []VerificationProjection{},
				}
			},
			wantReason: ManifestInputProfileUnresolved,
		},
		{
			name: "stored decisions no longer validate",
			manifest: func(t *testing.T, repository string) SetupManifest {
				t.Helper()

				decisions := manifestInputDecisionsForProfile(t, repository, catalog, "go-cli-tui")
				decisions["language.generated"] = false
				return newManifestInputFixture(t, repository, catalog, "go-cli-tui", decisions)
			},
			wantReason: ManifestInputDecisionsInvalid,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := t.TempDir()
			writeManifestInputFixture(t, repository, test.manifest(t, repository))
			input, err := ResolveManifestInput(repository, catalog)
			if !errors.Is(err, ErrManifestIncompatible) {
				t.Fatalf("ResolveManifestInput() error = %v, want ErrManifestIncompatible", err)
			}
			if input.State != ManifestInputIncompatible || input.Incompatibility != test.wantReason {
				t.Fatalf(
					"ResolveManifestInput() = state %q reason %q, want state %q reason %q",
					input.State,
					input.Incompatibility,
					ManifestInputIncompatible,
					test.wantReason,
				)
			}
		})
	}
}

func TestResolveManifestInputRejectsUnreadableOrInvalidManifest(t *testing.T) {
	t.Parallel()

	catalog := loadManifestInputCatalog(t)
	tests := []struct {
		name       string
		prepare    func(*testing.T, string)
		wantReason ManifestInputIncompatibility
	}{
		{
			name: "invalid strict JSON",
			prepare: func(t *testing.T, repository string) {
				t.Helper()
				writeManifestInputBytes(t, repository, []byte(`{"schemaVersion":"wrong"}`))
			},
			wantReason: ManifestInputManifestInvalid,
		},
		{
			name: "manifest path is not a regular file",
			prepare: func(t *testing.T, repository string) {
				t.Helper()
				if err := os.MkdirAll(filepath.Join(repository, filepath.FromSlash(manifestPath)), 0o755); err != nil {
					t.Fatalf("create manifest directory: %v", err)
				}
			},
			wantReason: ManifestInputManifestUnreadable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := t.TempDir()
			test.prepare(t, repository)
			input, err := ResolveManifestInput(repository, catalog)
			if !errors.Is(err, ErrManifestIncompatible) {
				t.Fatalf("ResolveManifestInput() error = %v, want ErrManifestIncompatible", err)
			}
			if input.State != ManifestInputIncompatible || input.Incompatibility != test.wantReason {
				t.Fatalf(
					"ResolveManifestInput() = state %q reason %q, want state %q reason %q",
					input.State,
					input.Incompatibility,
					ManifestInputIncompatible,
					test.wantReason,
				)
			}
		})
	}
}

func TestResolveManifestInputProfileFixedDecisionIsNotNew(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	created, err := InitCustomProfile(repository, "fixed-go", "go-cli-tui", catalog)
	if err != nil {
		t.Fatalf("InitCustomProfile() error = %v", err)
	}
	data, err := os.ReadFile(created.Path)
	if err != nil {
		t.Fatalf("read custom profile: %v", err)
	}
	var profileDocument customProfileDocument
	if err := strictJSON(data, &profileDocument); err != nil {
		t.Fatalf("decode custom profile: %v", err)
	}
	profileDocument.Values["secondbrain.enabled"] = true
	data, err = marshalCustomProfileDocument(profileDocument)
	if err != nil {
		t.Fatalf("serialize custom profile: %v", err)
	}
	if err := os.WriteFile(created.Path, data, 0o644); err != nil {
		t.Fatalf("write custom profile: %v", err)
	}

	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "fixed-go")
	delete(decisions, "secondbrain.enabled")
	manifest := newManifestInputFixture(t, repository, catalog, "fixed-go", decisions)
	writeManifestInputFixture(t, repository, manifest)

	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	if input.State != ManifestInputResolved {
		t.Fatalf("ResolveManifestInput() state = %q, want %q", input.State, ManifestInputResolved)
	}
	if len(input.NewDecisions) != 0 {
		t.Fatalf("ResolveManifestInput() new decisions = %#v, want none", input.NewDecisions)
	}
	fixedIndex := decisionValueIndex(input.Decisions, "secondbrain.enabled")
	if fixedIndex < 0 || input.Decisions[fixedIndex].Value != true {
		t.Fatalf("ResolveManifestInput() fixed decision = %#v, want true", input.Decisions)
	}
}

func TestResolveManifestInputWritesNothing(t *testing.T) {
	t.Parallel()

	repository := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "go-cli-tui")
	manifest := newManifestInputFixture(t, repository, catalog, "go-cli-tui", decisions)
	writeManifestInputFixture(t, repository, manifest)
	writeManifestInputBytesAt(t, repository, "authored.txt", []byte("repository-authored bytes\n"))
	before := snapshotManifestInputRepository(t, repository)

	if _, err := ResolveManifestInput(repository, catalog); err != nil {
		t.Fatalf("ResolveManifestInput() error = %v", err)
	}
	after := snapshotManifestInputRepository(t, repository)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("repository changed during manifest resolution:\nbefore=%#v\nafter=%#v", before, after)
	}
}

func loadManifestInputCatalog(t *testing.T) *Catalog {
	t.Helper()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}
	return catalog
}

func manifestInputDecisionsForProfile(
	t *testing.T,
	repository string,
	catalog *Catalog,
	profileID string,
) map[string]any {
	t.Helper()

	profile, err := ResolveProfile(repository, profileID, catalog)
	if err != nil {
		t.Fatalf("ResolveProfile(%q) error = %v", profileID, err)
	}
	input := make([]DecisionValue, 0, len(profile.Decisions))
	selected := make(map[string]struct{}, len(profile.Decisions))
	for _, id := range profile.Decisions {
		if _, fixed := profile.Values[id]; fixed {
			continue
		}
		input = append(input, DecisionValue{ID: id, Value: manifestInputCatalogSuggestion(t, catalog, id)})
		selected[id] = struct{}{}
	}
	for {
		resolved, missing, err := ResolveDecisionInput(profile, input, catalog)
		if err != nil {
			t.Fatalf("ResolveDecisionInput(%q) error = %v", profileID, err)
		}
		if len(missing) == 0 {
			values := make(map[string]any, len(resolved))
			for _, decision := range resolved {
				values[decision.ID] = cloneJSONValue(decision.Value)
			}
			return values
		}
		for _, id := range missing {
			if _, alreadySelected := selected[id]; alreadySelected {
				t.Fatalf("ResolveDecisionInput(%q) repeatedly reported missing %q", profileID, id)
			}
			input = append(input, DecisionValue{ID: id, Value: manifestInputCatalogSuggestion(t, catalog, id)})
			selected[id] = struct{}{}
		}
	}
}

func manifestInputCatalogSuggestion(t *testing.T, catalog *Catalog, id string) any {
	t.Helper()

	declaration, ok := catalog.decisions[id]
	if !ok {
		t.Fatalf("catalog decision %q is missing", id)
	}
	if suggestion, ok := declaration["suggestion"]; ok {
		return cloneJSONValue(suggestion)
	}
	if suggestion, ok := declaration["default"]; ok {
		return cloneJSONValue(suggestion)
	}
	t.Fatalf("catalog decision %q has no suggestion or default", id)
	return nil
}

func newManifestInputFixture(
	t *testing.T,
	repository string,
	catalog *Catalog,
	profileID string,
	decisions map[string]any,
) SetupManifest {
	t.Helper()

	profile, err := ResolveProfile(repository, profileID, catalog)
	if err != nil {
		t.Fatalf("ResolveProfile(%q) error = %v", profileID, err)
	}
	stored := make(map[string]ManifestDecision, len(decisions))
	for id, value := range decisions {
		stored[id] = ManifestDecision{Value: cloneJSONValue(value)}
	}
	return SetupManifest{
		SchemaVersion: ManifestSchema,
		Version:       ManifestVersion,
		Generator: ManifestGenerator{
			Skill:    "setup-context-driven",
			Version:  ManifestVersion,
			Baseline: "baseline." + profileID + "-" + ManifestVersion,
		},
		Profile:          profileID,
		ProfileDigest:    profile.Digest,
		CatalogDigest:    catalog.Digest(),
		Modules:          []string{},
		Decisions:        stored,
		ManagedArtifacts: []ManifestArtifact{},
		LocalSkills:      []string{},
		Verification:     []VerificationProjection{},
	}
}

func writeManifestInputFixture(t *testing.T, repository string, manifest SetupManifest) {
	t.Helper()

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("serialize Setup Manifest fixture: %v", err)
	}
	writeManifestInputBytes(t, repository, append(data, '\n'))
}

func writeManifestInputBytes(t *testing.T, repository string, data []byte) {
	t.Helper()

	writeManifestInputBytesAt(t, repository, manifestPath, data)
}

func writeManifestInputBytesAt(t *testing.T, repository, relative string, data []byte) {
	t.Helper()

	target := filepath.Join(repository, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("create fixture directory for %q: %v", relative, err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		t.Fatalf("write fixture %q: %v", relative, err)
	}
}

func assertManifestInputDecisionsEqual(
	t *testing.T,
	resolved []DecisionValue,
	stored map[string]ManifestDecision,
) {
	t.Helper()

	if len(resolved) != len(stored) {
		t.Fatalf("resolved decision count = %d, want %d", len(resolved), len(stored))
	}
	for _, decision := range resolved {
		manifestDecision, ok := stored[decision.ID]
		if !ok || !valuesEqual(decision.Value, manifestDecision.Value) {
			t.Fatalf("resolved decision %q = %#v, want %#v", decision.ID, decision.Value, manifestDecision.Value)
		}
	}
}

type manifestInputFileSnapshot struct {
	Mode fs.FileMode
	Data []byte
}

func snapshotManifestInputRepository(t *testing.T, repository string) map[string]manifestInputFileSnapshot {
	t.Helper()

	snapshot := make(map[string]manifestInputFileSnapshot)
	err := filepath.WalkDir(repository, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(repository, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(relative)] = manifestInputFileSnapshot{
			Mode: info.Mode(),
			Data: data,
		}
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot fixture repository: %v", err)
	}
	return snapshot
}
