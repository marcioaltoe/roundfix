// Suite: embedded Baseline catalog
// Invariant: the embedded catalog has one deterministic identity and rejects every invalid relationship.
// Boundary IN: embedded assets, strict JSON loading, normalization, references, and catalog digests.
// Boundary OUT: repository inspection, decisions, planning, apply, CLI, Agent Skills, and network access.

package baseline

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
)

//go:embed testdata
var compatibilityFixtures embed.FS

var updateBaselineDigests = flag.Bool("update", false, "regenerate derived digest artifacts")

var updateCatalogDiagnosticCharacterization = flag.Bool(
	"update-catalog-diagnostics",
	false,
	"regenerate the catalog diagnostic characterization corpus",
)

const baselineDigestRegenerationHint = "run 'make baseline-digests'"

const catalogDiagnosticGoldenPath = "testdata/catalog.diagnostics.golden.json"

func TestEmbeddedCatalog(t *testing.T) {
	t.Chdir(t.TempDir())

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf(
			"LoadEmbeddedCatalog() error = %v; %s to regenerate stale derived artifacts",
			err,
			baselineDigestRegenerationHint,
		)
	}

	wantProfiles := []string{
		"go-cli-tui",
		"rust-cli",
		"standard-typescript-monorepo",
	}
	if got := catalog.ProfileIDs(); !slices.Equal(got, wantProfiles) {
		t.Fatalf("ProfileIDs() = %v, want %v", got, wantProfiles)
	}

	wantModules := []string{
		"core",
		"context-workflow",
		"go",
		"cli-surface",
		"tui-surface",
		"autonomous-work",
		"spec-workflow",
		"external-triage",
		"secondbrain",
		"repository-extension",
	}
	gotModules, ok := catalog.OrderedModules("go-cli-tui")
	if !ok {
		t.Fatal("OrderedModules(go-cli-tui) missing")
	}
	if !slices.Equal(gotModules, wantModules) {
		t.Fatalf("OrderedModules(go-cli-tui) = %v, want %v", gotModules, wantModules)
	}

	for _, path := range []string{
		"profiles/go-cli-tui.json",
		"modules/core.json",
		"decisions.json",
		"templates/index.json",
		"retention/transition.managed-v2-to-portable-v3.json",
		"setups/go-cli.json",
	} {
		if _, ok := catalog.Asset(path); !ok {
			t.Errorf("Asset(%q) missing", path)
		}
	}

	profile, ok := catalog.Profile("rust-cli")
	if !ok {
		t.Fatal("Profile(rust-cli) missing")
	}
	profile.Data[0] = 'x'
	unchanged, _ := catalog.Profile("rust-cli")
	if unchanged.Data[0] == 'x' {
		t.Fatal("Profile(rust-cli) exposed mutable catalog bytes")
	}
}

func TestCatalogRegenerationMode(t *testing.T) {
	t.Parallel()

	const formatterDigestMismatch = "catalog.profile.formatter.goldenDigest.mismatch"
	if len(deferredDuringRegeneration) != 1 {
		t.Fatalf("deferredDuringRegeneration has %d codes, want 1", len(deferredDuringRegeneration))
	}
	if _, ok := deferredDuringRegeneration[formatterDigestMismatch]; !ok {
		t.Fatalf("deferredDuringRegeneration does not contain %q", formatterDigestMismatch)
	}

	var _ func(fs.FS) (*Catalog, error) = LoadCatalog
	var _ func() (*Catalog, error) = loadEmbeddedCatalogForRegeneration

	t.Run("strict load records allowlisted mismatch", func(t *testing.T) {
		t.Parallel()
		assets := cloneEmbeddedAssets(t)
		addFormatterGoldenDrift(t, assets)

		_, err := LoadCatalog(assets)
		requireCatalogDiagnostic(t, err, formatterDigestMismatch)
	})

	t.Run("regeneration load defers allowlisted mismatch", func(t *testing.T) {
		t.Parallel()
		assets := cloneEmbeddedAssets(t)
		addFormatterGoldenDrift(t, assets)

		if _, err := loadCatalog(assets, true); err != nil {
			t.Fatalf("load catalog in regeneration mode: %v", err)
		}
	})

	t.Run("regeneration load records non-allowlisted diagnostic", func(t *testing.T) {
		t.Parallel()
		assets := cloneEmbeddedAssets(t)
		replaceAsset(
			t,
			assets,
			"profiles/rust-cli.json",
			`"rust",`,
			`"rust", "missing-module",`,
		)

		_, err := loadCatalog(assets, true)
		requireCatalogDiagnostic(t, err, "catalog.profile.module.unknown")
	})
}

func TestRegenerationBreaksGoldenDigestCycle(t *testing.T) {
	t.Parallel()

	assets := newGoldenDigestCycleFixture(t)

	t.Run("strict load rejects the changed generated guide", func(t *testing.T) {
		t.Parallel()
		_, err := LoadCatalog(assets)
		requireCatalogDiagnostic(t, err, "catalog.profile.formatter.goldenDigest.mismatch")
	})

	t.Run("regeneration load accepts the changed generated guide", func(t *testing.T) {
		t.Parallel()
		catalog, err := loadCatalogForRegeneration(assets)
		if err != nil {
			t.Fatalf("load cycle fixture for regeneration: %v", err)
		}
		if _, ok := catalog.Profile("standard-typescript-monorepo"); !ok {
			t.Fatal("regeneration catalog has no standard TypeScript Profile")
		}
	})
}

func TestSourceMachinePathRecognizesMarkdownAndFileURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		match bool
	}{
		{
			name:  "plain macOS user path",
			value: "/Users/alice/work",
			match: true,
		},
		{
			name:  "Markdown link destination",
			value: "[repo](/Users/alice/work)",
			match: true,
		},
		{
			name:  "file URL",
			value: "file:///Users/alice/work",
			match: true,
		},
		{
			name:  "web URL home segment",
			value: "https://example.test/home/guide",
			match: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := sourceMachinePath.MatchString(test.value); got != test.match {
				t.Fatalf("sourceMachinePath.MatchString(%q) = %t, want %t", test.value, got, test.match)
			}
		})
	}
}

func TestCatalogDigest(t *testing.T) {
	t.Parallel()

	first, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("first LoadEmbeddedCatalog() error = %v", err)
	}
	second, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("second LoadEmbeddedCatalog() error = %v", err)
	}

	if first.Digest() == "" {
		t.Fatal("Digest() is empty")
	}
	if first.Digest() != second.Digest() {
		t.Fatalf("Digest() is nondeterministic: %q != %q", first.Digest(), second.Digest())
	}
	if !bytes.Equal(first.Normalized(), second.Normalized()) {
		t.Fatal("Normalized() is nondeterministic")
	}
	if !json.Valid(first.Normalized()) {
		t.Fatal("Normalized() is not valid JSON")
	}
}

func TestGuidanceCompositionAssets(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}
	if catalog.Digest() == "" {
		t.Fatal("embedded catalog digest is empty")
	}

	profile, ok := catalog.Profile("standard-typescript-monorepo")
	if !ok {
		t.Fatal("standard TypeScript Profile is missing")
	}
	var formatterProfile struct {
		Formatter struct {
			FixturePaths []string `json:"fixturePaths"`
			GoldenDigest string   `json:"goldenDigest"`
		} `json:"formatter"`
	}
	if err := json.Unmarshal(profile.Data, &formatterProfile); err != nil {
		t.Fatalf("decode formatter Profile: %v", err)
	}
	if len(formatterProfile.Formatter.FixturePaths) != 14 ||
		formatterProfile.Formatter.GoldenDigest == "" {
		t.Fatalf("formatter contract = %+v, want 14 digest-bound fixtures", formatterProfile.Formatter)
	}
	for _, fixturePath := range formatterProfile.Formatter.FixturePaths {
		if strings.Contains(fixturePath, specificRepositoryPath) {
			t.Fatalf("greenfield formatter fixture includes repository-specific carrier %q", fixturePath)
		}
		if _, ok := catalog.Asset(fixturePath); !ok {
			t.Errorf("formatter fixture %q is missing", fixturePath)
		}
	}

	for assetPath, asset := range cloneEmbeddedAssets(t) {
		if !strings.HasPrefix(assetPath, "modules/") &&
			!strings.HasPrefix(assetPath, "templates/") &&
			!strings.HasPrefix(assetPath, "formatter-fixtures/") &&
			!strings.Contains(assetPath, "/corpus/") {
			continue
		}
		folded := strings.ToLower(string(asset.Data))
		for _, brand := range []string{"fluxus", "oraculum"} {
			if strings.Contains(folded, brand) {
				t.Errorf("portable asset %q contains project brand %q", assetPath, brand)
			}
		}
	}
}

func TestProjectDecisionAssets(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	for _, profileID := range catalog.ProfileIDs() {
		profile, err := ResolveProfile("", profileID, catalog)
		if err != nil {
			t.Fatalf("resolve Profile %q: %v", profileID, err)
		}
		hasIdentifier := slices.Contains(profile.Decisions, "identifier.strategy")
		hasAuthProvider := slices.Contains(profile.Decisions, "auth.provider")
		switch profileID {
		case "standard-typescript-monorepo":
			if !hasIdentifier || !hasAuthProvider {
				t.Errorf(
					"Profile %q decisions = %v, want identifier.strategy and auth.provider",
					profileID,
					profile.Decisions,
				)
			}
		default:
			if hasIdentifier || hasAuthProvider {
				t.Errorf(
					"Profile %q unexpectedly selects project decisions: %v",
					profileID,
					profile.Decisions,
				)
			}
		}
	}

	for assetPath, required := range map[string][]string{
		"formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md": {
			"Use UUID version 7 for new project-owned Internal Identifiers only.",
			"external provider identifiers",
			"natural keys",
		},
		"formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md": {
			"Application HTTP mode: **Post-only**.",
			"**Better Auth** owns `GET` and `POST` for `/api/auth/*`",
			"Session, OAuth redirect, callback, and related provider protocol routes",
		},
		"formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md": {
			"without express maintainer authorization",
			"plugin declaration, or version pin",
		},
		"formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md": {
			"Project Constraints",
			"bounded repository-relative files",
			"completed or archived legacy Specs byte-identical",
		},
	} {
		asset, ok := catalog.Asset(assetPath)
		if !ok {
			t.Errorf("asset %q is missing", assetPath)
			continue
		}
		for _, clause := range required {
			if !strings.Contains(string(asset.Data), clause) {
				t.Errorf("asset %q is missing %q", assetPath, clause)
			}
		}
	}

	const toolingClauseID = "clause.core.require-tooling-authorization"
	source, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatalf("load maintained Source Baseline: %v", err)
	}
	if !sourceBaselineHasEntry(source, toolingClauseID) {
		t.Errorf("maintained Source Baseline does not account for %q", toolingClauseID)
	}
	for _, transitionID := range catalog.TransitionIDs() {
		transition, err := catalog.UpgradeRetentionContract(transitionID)
		if err != nil {
			t.Fatalf("load retention contract %q: %v", transitionID, err)
		}
		if !retentionTargetsClause(transition, toolingClauseID) {
			t.Errorf("retention contract %q does not retain %q", transitionID, toolingClauseID)
		}
	}

	for assetPath, asset := range cloneEmbeddedAssets(t) {
		if !strings.HasPrefix(assetPath, "modules/") &&
			!strings.HasPrefix(assetPath, "profiles/") &&
			!strings.HasPrefix(assetPath, "templates/") &&
			!strings.HasPrefix(assetPath, "formatter-fixtures/") &&
			!strings.Contains(assetPath, "/corpus/") {
			continue
		}
		folded := strings.ToLower(string(asset.Data))
		for _, brand := range []string{"fluxus", "oraculum"} {
			if strings.Contains(folded, brand) {
				t.Errorf("project-agnostic asset %q contains %q", assetPath, brand)
			}
		}
	}

	tests := []struct {
		name       string
		decisionID string
	}{
		{name: "identifier strategy", decisionID: "identifier.strategy"},
		{name: "Better Auth provider", decisionID: "auth.provider"},
	}
	for _, test := range tests {
		t.Run("catalog rejects missing "+test.name, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			replaceAsset(
				t,
				assets,
				"profiles/standard-typescript-monorepo.json",
				`    "`+test.decisionID+`",`+"\n",
				"",
			)
			_, err := LoadCatalog(assets)
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) ||
				!validationErr.Has("catalog.project-decision.profile.missing") {
				t.Fatalf("LoadCatalog() error = %v, want missing project-decision diagnostic", err)
			}
		})
	}
}

func TestCatalogCompatibility(t *testing.T) {
	if !*updateBaselineDigests {
		t.Parallel()
	}

	var catalog *Catalog
	var err error
	if *updateBaselineDigests {
		catalog, err = loadEmbeddedCatalogForRegeneration()
	} else {
		catalog, err = LoadEmbeddedCatalog()
	}
	if err != nil {
		t.Fatalf(
			"LoadEmbeddedCatalog() error = %v; %s to regenerate stale derived artifacts",
			err,
			baselineDigestRegenerationHint,
		)
	}

	wantNormalized := catalog.Normalized()
	wantDigestBytes := []byte(catalog.Digest() + "\n")
	if *updateBaselineDigests {
		regenerateCatalogCompatibility(t, catalog)
	} else {
		wantNormalized, err = fs.ReadFile(
			compatibilityFixtures,
			"testdata/catalog.normalized.json",
		)
		if err != nil {
			t.Fatalf("read normalized compatibility fixture: %v", err)
		}
		wantDigestBytes, err = fs.ReadFile(compatibilityFixtures, "testdata/catalog.digest")
		if err != nil {
			t.Fatalf("read digest compatibility fixture: %v", err)
		}
	}
	if !bytes.Equal(catalog.Normalized(), wantNormalized) {
		t.Fatalf(
			"embedded catalog normalized bytes differ from the compatibility fixture; %s",
			baselineDigestRegenerationHint,
		)
	}
	if got, want := catalog.Digest()+"\n", string(wantDigestBytes); got != want {
		t.Fatalf("Digest() = %q, want %q; %s", got, want, baselineDigestRegenerationHint)
	}

	legacyAssets, err := fs.Sub(compatibilityFixtures, "testdata/legacy-v2/assets")
	if err != nil {
		t.Fatalf("open maintained v2 compatibility assets: %v", err)
	}
	legacy, err := LoadCatalog(legacyAssets)
	if err != nil {
		t.Fatalf("LoadCatalog(legacy-v2) error = %v", err)
	}
	if got, want := legacy.ProfileIDs(), []string{"fixture"}; !slices.Equal(got, want) {
		t.Fatalf("legacy ProfileIDs() = %v, want %v", got, want)
	}
}

func regenerateCatalogCompatibility(t *testing.T, catalog *Catalog) {
	t.Helper()
	writeBaselineDerivedArtifact(
		t,
		"testdata/catalog.normalized.json",
		catalog.Normalized(),
	)
	writeBaselineDerivedArtifact(t, "testdata/catalog.digest", []byte(catalog.Digest()+"\n"))
}

func regenerateCatalogCompatibilityFromAssets(t *testing.T) {
	t.Helper()
	catalog, err := loadCatalogForRegeneration(os.DirFS("assets"))
	if err != nil {
		t.Fatalf("load regenerated Baseline assets: %v", err)
	}
	regenerateCatalogCompatibility(t, catalog)
}

func writeBaselineDerivedArtifact(t *testing.T, filePath string, data []byte) {
	t.Helper()
	current, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read derived artifact %s: %v", filePath, err)
	}
	if bytes.Equal(current, data) {
		return
	}
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		t.Fatalf("write derived artifact %s: %v", filePath, err)
	}
	t.Logf("regenerated %s", filePath)
}

type catalogDiagnosticRecord struct {
	Code    string `json:"code"`
	Subject string `json:"subject"`
	Detail  string `json:"detail"`
}

type catalogDiagnosticResult struct {
	Catalog     string                    `json:"catalog"`
	Diagnostics []catalogDiagnosticRecord `json:"diagnostics"`
}

type catalogDiagnosticCorpus struct {
	SchemaVersion string                    `json:"schemaVersion"`
	Catalogs      []catalogDiagnosticResult `json:"catalogs"`
}

type catalogDiagnosticFixture struct {
	name   string
	assets func(t *testing.T) fs.FS
}

func TestCatalogDiagnosticCharacterization(t *testing.T) {
	actual := loadCatalogDiagnosticCorpus(t)
	actualBytes := marshalCatalogDiagnosticCorpus(t, actual)

	if *updateCatalogDiagnosticCharacterization {
		writeBaselineDerivedArtifact(t, catalogDiagnosticGoldenPath, actualBytes)
	}

	wantBytes, err := os.ReadFile(catalogDiagnosticGoldenPath)
	if err != nil {
		t.Fatalf("read catalog diagnostic characterization: %v", err)
	}
	var want catalogDiagnosticCorpus
	if err := json.Unmarshal(wantBytes, &want); err != nil {
		t.Fatalf("decode catalog diagnostic characterization: %v", err)
	}
	if want.SchemaVersion != actual.SchemaVersion {
		t.Fatalf(
			"catalog diagnostic characterization schema = %q, want %q",
			want.SchemaVersion,
			actual.SchemaVersion,
		)
	}
	if diff := catalogDiagnosticDiff(want.Catalogs, actual.Catalogs); diff != "" {
		t.Fatalf(
			"catalog diagnostic characterization differs; run with -update-catalog-diagnostics to re-record an intentional change:\n%s",
			diff,
		)
	}

	t.Run("failure diff names the changed code", func(t *testing.T) {
		mutated := cloneCatalogDiagnosticResults(actual.Catalogs)
		var changedCode string
		for catalogIndex := range mutated {
			if len(mutated[catalogIndex].Diagnostics) == 0 {
				continue
			}
			changedCode = mutated[catalogIndex].Diagnostics[0].Code
			mutated[catalogIndex].Diagnostics[0].Detail += " deliberately changed"
			break
		}
		if changedCode == "" {
			t.Fatal("catalog diagnostic characterization contains no diagnostic to mutation-test")
		}
		diff := catalogDiagnosticDiff(actual.Catalogs, mutated)
		if !strings.Contains(diff, changedCode) {
			t.Fatalf("diagnostic diff does not name changed code %q:\n%s", changedCode, diff)
		}
		if !strings.Contains(diff, "- ") || !strings.Contains(diff, "+ ") {
			t.Fatalf("diagnostic diff does not show removed and added records:\n%s", diff)
		}
	})
}

func loadCatalogDiagnosticCorpus(t *testing.T) catalogDiagnosticCorpus {
	t.Helper()

	fixtures := catalogDiagnosticFixtures()
	sort.Slice(fixtures, func(i, j int) bool { return fixtures[i].name < fixtures[j].name })
	results := make([]catalogDiagnosticResult, 0, len(fixtures))
	for index, fixture := range fixtures {
		if index != 0 && fixtures[index-1].name == fixture.name {
			t.Fatalf("duplicate catalog diagnostic fixture %q", fixture.name)
		}
		_, err := LoadCatalog(fixture.assets(t))
		diagnostics := make([]catalogDiagnosticRecord, 0)
		if err != nil {
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("load catalog diagnostic fixture %q: %v", fixture.name, err)
			}
			for _, diagnostic := range validationErr.Diagnostics {
				diagnostics = append(diagnostics, catalogDiagnosticRecord{
					Code:    diagnostic.Code,
					Subject: diagnostic.Path,
					Detail:  diagnostic.Info,
				})
			}
			sortCatalogDiagnosticRecords(diagnostics)
		}
		results = append(results, catalogDiagnosticResult{
			Catalog:     fixture.name,
			Diagnostics: diagnostics,
		})
	}
	return catalogDiagnosticCorpus{
		SchemaVersion: "roundfix/baseline-catalog-diagnostics/v1",
		Catalogs:      results,
	}
}

func catalogDiagnosticFixtures() []catalogDiagnosticFixture {
	fixtures := []catalogDiagnosticFixture{
		{
			name: "embedded",
			assets: func(*testing.T) fs.FS {
				return embeddedAssets
			},
		},
		{
			name: "legacy-v2",
			assets: func(t *testing.T) fs.FS {
				t.Helper()
				assets, err := fs.Sub(compatibilityFixtures, "testdata/legacy-v2/assets")
				if err != nil {
					t.Fatalf("open maintained v2 compatibility assets: %v", err)
				}
				return assets
			},
		},
		{
			name: "repository-assets",
			assets: func(*testing.T) fs.FS {
				return os.DirFS("assets")
			},
		},
	}

	for _, mutation := range catalogMutationTests() {
		mutation := mutation
		fixtures = append(fixtures, editedCatalogDiagnosticFixture(
			"mutation/"+mutation.name,
			mutation.edit,
		))
	}
	for _, decisionID := range []string{"auth.provider", "identifier.strategy"} {
		decisionID := decisionID
		fixtures = append(fixtures, editedCatalogDiagnosticFixture(
			"project-decision/missing "+decisionID,
			func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"profiles/standard-typescript-monorepo.json",
					`    "`+decisionID+`",`+"\n",
					"",
				)
			},
		))
	}
	for _, authority := range catalogToolingAuthorityTests() {
		authority := authority
		fixtures = append(fixtures, editedCatalogDiagnosticFixture(
			"tooling-authority/"+authority.name,
			authority.edit,
		))
	}

	const toolingClauseID = "clause.core.require-tooling-authorization"
	fixtures = append(fixtures,
		editedCatalogDiagnosticFixture(
			"source-baseline/new required clause",
			func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/core.json",
					`"clauses": [
        {
          "id": "clause.core.fix-root-causes",`,
					`"clauses": [
        {
          "id": "clause.core.characterization-new-clause",
          "enforcement": "mandatory",
          "guidance": "Characterize a newly required clause."
        },
        {
          "id": "clause.core.fix-root-causes",`,
				)
			},
		),
		editedCatalogDiagnosticFixture(
			"tooling-authority/source accounting removal",
			func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json",
					toolingClauseID,
					"clause.missing-tooling-authority-accounting",
				)
			},
		),
		editedCatalogDiagnosticFixture(
			"tooling-authority/retention accounting removal",
			func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"retention/transition.managed-v2-to-portable-v3.json",
					`, "`+toolingClauseID+`"`,
					``,
				)
			},
		),
	)
	return fixtures
}

func editedCatalogDiagnosticFixture(
	name string,
	edit func(t *testing.T, assets fstest.MapFS),
) catalogDiagnosticFixture {
	return catalogDiagnosticFixture{
		name: name,
		assets: func(t *testing.T) fs.FS {
			t.Helper()
			assets := cloneEmbeddedAssets(t)
			edit(t, assets)
			return assets
		},
	}
}

func marshalCatalogDiagnosticCorpus(t *testing.T, corpus catalogDiagnosticCorpus) []byte {
	t.Helper()
	data, err := json.MarshalIndent(corpus, "", "  ")
	if err != nil {
		t.Fatalf("encode catalog diagnostic characterization: %v", err)
	}
	return append(data, '\n')
}

func cloneCatalogDiagnosticResults(results []catalogDiagnosticResult) []catalogDiagnosticResult {
	cloned := make([]catalogDiagnosticResult, len(results))
	for index, result := range results {
		cloned[index] = result
		cloned[index].Diagnostics = slices.Clone(result.Diagnostics)
	}
	return cloned
}

func catalogDiagnosticDiff(
	want []catalogDiagnosticResult,
	got []catalogDiagnosticResult,
) string {
	wantByCatalog := make(map[string][]catalogDiagnosticRecord, len(want))
	gotByCatalog := make(map[string][]catalogDiagnosticRecord, len(got))
	names := make(map[string]struct{}, len(want)+len(got))
	for _, result := range want {
		wantByCatalog[result.Catalog] = result.Diagnostics
		names[result.Catalog] = struct{}{}
	}
	for _, result := range got {
		gotByCatalog[result.Catalog] = result.Diagnostics
		names[result.Catalog] = struct{}{}
	}

	changedCodes := make(map[string]struct{})
	var details strings.Builder
	for _, name := range sortedKeys(names) {
		wantDiagnostics, wantOK := wantByCatalog[name]
		gotDiagnostics, gotOK := gotByCatalog[name]
		if wantOK && gotOK && slices.Equal(wantDiagnostics, gotDiagnostics) {
			continue
		}
		fmt.Fprintf(&details, "catalog %q:\n", name)
		if !wantOK {
			details.WriteString("+ catalog added\n")
		}
		if !gotOK {
			details.WriteString("- catalog removed\n")
		}
		removed, added := changedCatalogDiagnostics(wantDiagnostics, gotDiagnostics)
		for _, diagnostic := range removed {
			changedCodes[diagnostic.Code] = struct{}{}
			fmt.Fprintf(&details, "- %s\n", formatCatalogDiagnostic(diagnostic))
		}
		for _, diagnostic := range added {
			changedCodes[diagnostic.Code] = struct{}{}
			fmt.Fprintf(&details, "+ %s\n", formatCatalogDiagnostic(diagnostic))
		}
	}
	if details.Len() == 0 {
		return ""
	}
	codes := sortedKeys(changedCodes)
	if len(codes) == 0 {
		return "changed diagnostic codes: none (catalog set changed)\n" + details.String()
	}
	return "changed diagnostic codes: " + strings.Join(codes, ", ") + "\n" + details.String()
}

func changedCatalogDiagnostics(
	want []catalogDiagnosticRecord,
	got []catalogDiagnosticRecord,
) (removed []catalogDiagnosticRecord, added []catalogDiagnosticRecord) {
	remaining := make(map[catalogDiagnosticRecord]int, len(want))
	for _, diagnostic := range want {
		remaining[diagnostic]++
	}
	for _, diagnostic := range got {
		if remaining[diagnostic] != 0 {
			remaining[diagnostic]--
			continue
		}
		added = append(added, diagnostic)
	}
	for diagnostic, count := range remaining {
		for range count {
			removed = append(removed, diagnostic)
		}
	}
	sortCatalogDiagnosticRecords(removed)
	sortCatalogDiagnosticRecords(added)
	return removed, added
}

func sortCatalogDiagnosticRecords(diagnostics []catalogDiagnosticRecord) {
	sort.Slice(diagnostics, func(i, j int) bool {
		left, right := diagnostics[i], diagnostics[j]
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		if left.Subject != right.Subject {
			return left.Subject < right.Subject
		}
		return left.Detail < right.Detail
	})
}

func formatCatalogDiagnostic(diagnostic catalogDiagnosticRecord) string {
	return fmt.Sprintf(
		"%s subject=%q detail=%q",
		diagnostic.Code,
		diagnostic.Subject,
		diagnostic.Detail,
	)
}

func TestInstructionHierarchy(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}

	levels := catalog.InstructionHierarchy()
	want := []string{
		"universal",
		"context",
		"spec",
		"autonomous",
		"stack",
		"surface",
		"optional-knowledge",
		"repository-specific",
	}
	got := make([]string, len(levels))
	for index, level := range levels {
		got[index] = level.ID
	}
	if !slices.Equal(got, want) {
		t.Fatalf("InstructionHierarchy() = %v, want %v", got, want)
	}
	if got := levels[5].RootBlocks; !slices.Equal(got, []string{
		"root.cli-surface",
		"root.tui-surface",
	}) {
		t.Fatalf("surface RootBlocks = %v", got)
	}
}

func TestSemanticOwnerRegistry(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}
	profile, err := ResolveProfile("", "go-cli-tui", catalog)
	if err != nil {
		t.Fatalf("ResolveProfile() error = %v", err)
	}
	activeModules, artifacts, err := resolveManagedArtifacts(
		catalog,
		profile,
		planTestDecisions(),
		false,
	)
	if err != nil {
		t.Fatalf("resolveManagedArtifacts() error = %v", err)
	}
	activeArtifacts := make([]string, len(artifacts))
	for index, artifact := range artifacts {
		activeArtifacts[index] = artifact.ID
	}

	registry := catalog.SemanticOwnerRegistry(activeModules, activeArtifacts)
	for _, artifact := range artifacts {
		if artifact.Kind != "guide" {
			continue
		}
		owner, ok := registry[artifact.ID]
		if !ok {
			t.Errorf("active guide %q has no semantic owner", artifact.ID)
			continue
		}
		if owner.ManagedID != artifact.ID || owner.Module != artifact.Module ||
			owner.Path != artifact.Path || owner.Title == "" ||
			len(owner.Classifications) == 0 {
			t.Errorf("semantic owner %q = %+v", artifact.ID, owner)
		}
	}
	for _, inactive := range []string{
		"guide.external-triage",
		"guide.secondbrain",
	} {
		if _, ok := registry[inactive]; ok {
			t.Errorf("inactive guide %q has a semantic destination", inactive)
		}
	}
}

type catalogMutationTest struct {
	name string
	code string
	edit func(t *testing.T, assets fstest.MapFS)
}

func catalogMutationTests() []catalogMutationTest {
	return []catalogMutationTest{
		{
			name: "duplicate JSON key",
			code: "catalog.json.key.duplicate",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"decisions.json",
					`"schemaVersion": "setup-context-driven/decisions-v1",`,
					`"schemaVersion": "setup-context-driven/decisions-v1", "schemaVersion": "setup-context-driven/decisions-v1",`,
				)
			},
		},
		{
			name: "unknown schema field",
			code: "catalog.profile.field.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"profiles/rust-cli.json",
					`"schemaVersion": "setup-context-driven/profile/0.0.1",`,
					`"schemaVersion": "setup-context-driven/profile/0.0.1", "unexpected": true,`,
				)
			},
		},
		{
			name: "duplicate catalog identifier",
			code: "catalog.module.id.duplicate",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/bun.json",
					"\"schemaVersion\": \"setup-context-driven/module-v3\",\n  \"id\": \"bun\"",
					"\"schemaVersion\": \"setup-context-driven/module-v3\",\n  \"id\": \"core\"",
				)
			},
		},
		{
			name: "unknown profile module",
			code: "catalog.profile.module.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(t, assets, "profiles/rust-cli.json", `"rust",`, `"rust", "missing-module",`)
			},
		},
		{
			name: "unknown setup",
			code: "catalog.profile.setup.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(t, assets, "profiles/rust-cli.json", `"setup": "rust-cli"`, `"setup": "missing-setup"`)
			},
		},
		{
			name: "unknown template reference",
			code: "catalog.template.reference.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/rust.json",
					`"template": "template.root.rust"`,
					`"template": "template.missing"`,
				)
			},
		},
		{
			name: "module dependency cycle",
			code: "catalog.module.dependency.cycle",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(t, assets, "modules/core.json", `"dependsOn": []`, `"dependsOn": ["context-workflow"]`)
			},
		},
		{
			name: "decision dependency cycle",
			code: "catalog.decision.dependency.cycle",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"decisions.json",
					`"activateModules": ["spec-workflow"]`,
					`"activateModules": ["spec-workflow"], "requireDecisions": ["spec.scaffold"]`,
				)
			},
		},
		{
			name: "invalid decision default",
			code: "catalog.decision.default.invalid",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"decisions.json",
					`"id": "spec.scaffold",
      "version": 1,
      "type": "boolean",
      "default": true`,
					`"id": "spec.scaffold",
      "version": 1,
      "type": "boolean",
      "default": "yes"`,
				)
			},
		},
		{
			name: "unknown decision effect module",
			code: "catalog.decision.effect.module.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(t, assets, "decisions.json", `"activateModules": ["spec-workflow"]`, `"activateModules": ["missing-module"]`)
			},
		},
		{
			name: "unknown retention target",
			code: "catalog.transition.target.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"retention/transition.managed-v2-to-portable-v3.json",
					`"fromClause": "clause.managed-v2.keep-root-compact", "disposition": "retained", "targets": ["clause.core.keep-root-compact"]`,
					`"fromClause": "clause.managed-v2.keep-root-compact", "disposition": "retained", "targets": ["clause.missing"]`,
				)
			},
		},
		{
			name: "undeclared template token",
			code: "catalog.template.token.undeclared",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				asset := assets["templates/root/core.md"]
				asset.Data = append(append([]byte(nil), asset.Data...), []byte("\n{{decision.missing}}\n")...)
				assets["templates/root/core.md"] = asset
			},
		},
		{
			name: "unknown Source Baseline profile",
			code: "catalog.sourceBaseline.profile.unknown",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"source-baselines/index.json",
					`"profile": "standard-typescript-monorepo"`,
					`"profile": "missing-profile"`,
				)
			},
		},
		{
			name: "missing formatter fixture",
			code: "catalog.profile.formatter.fixture.missing",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				delete(
					assets,
					"formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md",
				)
			},
		},
		{
			name: "formatter golden drift",
			code: "catalog.profile.formatter.goldenDigest.mismatch",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				addFormatterGoldenDrift(t, assets)
			},
		},
		{
			name: "stale Source Baseline accounting target",
			code: "catalog.sourceBaseline.integrity.invalid",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"source-baselines/baseline.standard-typescript-monorepo-0.0.1/accounting.json",
					`"clause.core.review-new-dependencies"`,
					`"clause.missing"`,
				)
			},
		},
		{
			name: "setup snapshot digest drift",
			code: "catalog.setup.digest.mismatch",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(t, assets, "setups/rust-cli.json", `"digest": "`, `"digest": "0`)
			},
		},
		{
			name: "invalid instruction precedence",
			code: "catalog.instruction-hierarchy.order.invalid",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/core.json",
					`"id": "context", "title": "Context and documentation"`,
					`"id": "spec", "title": "Context and documentation"`,
				)
			},
		},
		{
			name: "instruction dependency after dependent",
			code: "catalog.instruction-hierarchy.dependency.order",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/core.json",
					`"rootBlocks": ["root.cli-surface", "root.tui-surface"]`,
					`"rootBlocks": ["root.tui-surface", "root.cli-surface"]`,
				)
			},
		},
		{
			name: "duplicate semantic ownership",
			code: "catalog.semantic-owner.classification.duplicate",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/core.json",
					`"classifications": ["domain-language", "identifier-policy"]`,
					`"classifications": ["universal-execution", "identifier-policy"]`,
				)
			},
		},
		{
			name: "narrower clause weakens universal policy",
			code: "catalog.clause.weakening.prohibited",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"modules/autonomous-work.json",
					`"guidance": "The Supervisor must not write feature code or tests."`,
					`"guidance": "The Supervisor must not write feature code or tests.", "weakens": ["clause.core.fix-root-causes"]`,
				)
			},
		},
	}
}

func TestCatalogMutation(t *testing.T) {
	t.Parallel()

	for _, test := range catalogMutationTests() {
		t.Run(test.name, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			test.edit(t, assets)

			_, err := LoadCatalog(assets)
			if err == nil {
				t.Fatalf("LoadCatalog() error = nil, want diagnostic %q", test.code)
			}
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("LoadCatalog() error = %T, want *ValidationError: %v", err, err)
			}
			if !validationErr.Has(test.code) {
				t.Fatalf("LoadCatalog() diagnostics = %v, want %q", validationErr.Diagnostics, test.code)
			}
		})
	}
}

type catalogToolingAuthorityTest struct {
	name string
	code string
	edit func(t *testing.T, assets fstest.MapFS)
}

func catalogToolingAuthorityTests() []catalogToolingAuthorityTest {
	return []catalogToolingAuthorityTest{
		{
			name: "Profile omits the universal rule",
			code: "catalog.tooling-authority.profile.rule.missing",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				removeProfileRule(t, assets, "profiles/rust-cli.json", "rule.core.tooling-authority")
			},
		},
		{
			name: "decision excludes the universal guide",
			code: "catalog.tooling-authority.effect.prohibited",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"decisions.json",
					`"activateModules": ["spec-workflow"]`,
					`"activateModules": ["spec-workflow"], "excludeArtifacts": ["guide.agent-instructions"]`,
				)
			},
		},
		{
			name: "decision controls the core module",
			code: "catalog.tooling-authority.effect.prohibited",
			edit: func(t *testing.T, assets fstest.MapFS) {
				t.Helper()
				replaceAsset(
					t,
					assets,
					"decisions.json",
					`"activateModules": ["spec-workflow"]`,
					`"activateModules": ["spec-workflow", "core"]`,
				)
			},
		},
	}
}

func TestToolingAuthorityCannotBeDisabled(t *testing.T) {
	t.Parallel()

	for _, test := range catalogToolingAuthorityTests() {
		t.Run(test.name, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			test.edit(t, assets)

			_, err := LoadCatalog(assets)
			if err == nil {
				t.Fatalf("LoadCatalog() error = nil, want diagnostic %q", test.code)
			}
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("LoadCatalog() error = %T, want *ValidationError: %v", err, err)
			}
			if !validationErr.Has(test.code) {
				t.Fatalf("LoadCatalog() diagnostics = %v, want %q", validationErr.Diagnostics, test.code)
			}
		})
	}
}

func TestToolingAuthorityAccounting(t *testing.T) {
	t.Parallel()

	const clauseID = "clause.core.require-tooling-authorization"
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	source, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatalf("load maintained Source Baseline: %v", err)
	}
	if !sourceBaselineHasEntry(source, clauseID) {
		t.Fatalf("maintained Source Baseline does not account for %q", clauseID)
	}
	for _, transitionID := range catalog.TransitionIDs() {
		transition, err := catalog.UpgradeRetentionContract(transitionID)
		if err != nil {
			t.Fatalf("load retention contract %q: %v", transitionID, err)
		}
		if !retentionTargetsClause(transition, clauseID) {
			t.Errorf("retention contract %q does not retain %q", transitionID, clauseID)
		}
	}

	t.Run("source accounting removal invalidates the catalog", func(t *testing.T) {
		assets := cloneEmbeddedAssets(t)
		replaceAsset(
			t,
			assets,
			"source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json",
			clauseID,
			"clause.missing-tooling-authority-accounting",
		)
		if _, err := LoadCatalog(assets); err == nil {
			t.Fatal("LoadCatalog() error = nil after removing tooling-authority source accounting")
		}
	})

	t.Run("retention accounting removal invalidates the catalog", func(t *testing.T) {
		assets := cloneEmbeddedAssets(t)
		replaceAsset(
			t,
			assets,
			"retention/transition.managed-v2-to-portable-v3.json",
			`, "`+clauseID+`"`,
			``,
		)
		_, err := LoadCatalog(assets)
		if err == nil {
			t.Fatal("LoadCatalog() error = nil after removing tooling-authority retention accounting")
		}
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) ||
			!validationErr.Has("catalog.tooling-authority.retention.missing") {
			t.Fatalf("LoadCatalog() error = %v, want tooling-authority retention diagnostic", err)
		}
	})
}

func sourceBaselineHasEntry(source SourceBaseline, id string) bool {
	for _, entry := range source.Entries {
		if entry.ID == id {
			return true
		}
	}
	return false
}

func retentionTargetsClause(contract UpgradeRetentionContract, clauseID string) bool {
	for _, accounting := range contract.Accounting {
		if slices.Contains(accounting.Targets, clauseID) {
			return true
		}
	}
	return false
}

func cloneEmbeddedAssets(t *testing.T) fstest.MapFS {
	t.Helper()

	assets := fstest.MapFS{}
	err := fs.WalkDir(embeddedAssets, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(embeddedAssets, path)
		if err != nil {
			return err
		}
		assets[path] = &fstest.MapFile{Data: append([]byte(nil), data...), Mode: 0o444}
		return nil
	})
	if err != nil {
		t.Fatalf("clone embedded assets: %v", err)
	}
	return assets
}

func addFormatterGoldenDrift(t *testing.T, assets fstest.MapFS) {
	t.Helper()

	const fixturePath = "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md"
	asset := assets[fixturePath]
	asset.Data = append(append([]byte(nil), asset.Data...), []byte("\ndrift\n")...)
	assets[fixturePath] = asset
}

func newGoldenDigestCycleFixture(t *testing.T) fstest.MapFS {
	t.Helper()

	assets := cloneEmbeddedAssets(t)
	addFormatterGoldenDrift(t, assets)
	return assets
}

func mustEmbeddedCatalogForRegeneration(t *testing.T) *Catalog {
	t.Helper()

	catalog, err := loadEmbeddedCatalogForRegeneration()
	if err != nil {
		t.Fatalf("load embedded Baseline catalog for regeneration: %v", err)
	}
	return catalog
}

func requireCatalogDiagnostic(t *testing.T, err error, code string) {
	t.Helper()

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("load catalog error = %v, want *ValidationError", err)
	}
	if !validationErr.Has(code) {
		t.Fatalf("load catalog diagnostics = %v, want %q", validationErr.Diagnostics, code)
	}
}

// removeProfileRule drops one rule from a profile's requiredRules by identity.
// A textual edit here would depend on how the asset happens to wrap its rule
// list, and a reformat would turn the mutation into a silent no-op whose only
// symptom is the catalog loading successfully.
func removeProfileRule(t *testing.T, assets fstest.MapFS, path, rule string) {
	t.Helper()

	asset, ok := assets[path]
	if !ok {
		t.Fatalf("asset %q missing", path)
	}
	var profile map[string]any
	if err := json.Unmarshal(asset.Data, &profile); err != nil {
		t.Fatalf("decode profile %q: %v", path, err)
	}
	rules, ok := profile["requiredRules"].([]any)
	if !ok {
		t.Fatalf("profile %q has no requiredRules array", path)
	}
	kept := slices.DeleteFunc(slices.Clone(rules), func(entry any) bool {
		return entry == rule
	})
	if len(kept) != len(rules)-1 {
		t.Fatalf("profile %q requiredRules removed %d entries for %q, want 1", path, len(rules)-len(kept), rule)
	}
	profile["requiredRules"] = kept
	mutated, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("encode profile %q: %v", path, err)
	}
	asset.Data = mutated
	assets[path] = asset
}

func replaceAsset(t *testing.T, assets fstest.MapFS, path, old, replacement string) {
	t.Helper()

	asset, ok := assets[path]
	if !ok {
		t.Fatalf("asset %q missing", path)
	}
	content := string(asset.Data)
	if strings.Count(content, old) != 1 {
		t.Fatalf("asset %q has %d occurrences of %q, want 1", path, strings.Count(content, old), old)
	}
	asset.Data = []byte(strings.Replace(content, old, replacement, 1))
	assets[path] = asset
}
