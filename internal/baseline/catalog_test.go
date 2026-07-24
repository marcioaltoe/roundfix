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
	"io/fs"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

//go:embed testdata
var compatibilityFixtures embed.FS

func TestEmbeddedCatalog(t *testing.T) {
	t.Chdir(t.TempDir())

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
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

func TestCatalogCompatibility(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}

	wantNormalized, err := fs.ReadFile(compatibilityFixtures, "testdata/catalog.normalized.json")
	if err != nil {
		t.Fatalf("read normalized compatibility fixture: %v", err)
	}
	if !bytes.Equal(catalog.Normalized(), wantNormalized) {
		t.Fatal("embedded catalog normalized bytes differ from the compatibility fixture")
	}

	wantDigestBytes, err := fs.ReadFile(compatibilityFixtures, "testdata/catalog.digest")
	if err != nil {
		t.Fatalf("read digest compatibility fixture: %v", err)
	}
	if got, want := catalog.Digest()+"\n", string(wantDigestBytes); got != want {
		t.Fatalf("Digest() = %q, want %q", got, want)
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

func TestCatalogMutation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		edit func(t *testing.T, assets fstest.MapFS)
	}{
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

	for _, test := range tests {
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
