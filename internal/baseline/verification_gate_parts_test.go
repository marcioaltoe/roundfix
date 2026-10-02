package baseline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAGateThatReachesBothPartsReportsNoDivergence(t *testing.T) {
	t.Parallel()
	for _, makefile := range []string{
		"verify: verify-go workspace\nverify-go:\n\tgo test ./...\nworkspace:\n\t@-rtk bun run verify\n",
		"verify: all\nall: verify-go | workspace\nverify-go:\n\tgo test ./...\nworkspace: ; bun run verify\n",
	} {
		_, divergences := gatePartsProjection(t, makefile, "rtk make verify", nil)
		if len(divergences) != 0 {
			t.Fatalf("reached parts reported: %+v", divergences)
		}
	}
}

func TestAGateThatSkipsAPartReportsItOnce(t *testing.T) {
	t.Parallel()
	_, divergences := gatePartsProjection(t, "verify: verify-go verify-go\nverify-go:\n\tgo test ./...\nunused:\n\tbun run verify\n", "make verify", nil)
	if len(divergences) != 1 {
		t.Fatalf("divergences = %+v, want one", divergences)
	}
	got := divergences[0]
	if got.ID != "verification.workspace" || got.Blocking || got.Requirement != CapabilityRecommended ||
		got.Code != "verification.gate.part.missing" ||
		got.Message != `the repository gate "make verify" does not run "bun run verify"` ||
		got.NextAction != "make the repository gate run the named Verification, or map its role to a command the gate runs" {
		t.Fatalf("missing workspace divergence = %+v", got)
	}
}

func TestAPartReachedThroughARecipeInvocationCounts(t *testing.T) {
	t.Parallel()
	for _, invocation := range []string{"$(MAKE)", "${MAKE}", "make"} {
		t.Run(invocation, func(t *testing.T) {
			_, divergences := gatePartsProjection(t, "verify:\n\t@-rtk "+invocation+" verify-go\n\tbun run verify\nverify-go:\n\tgo test ./...\n", "make verify", nil)
			if len(divergences) != 0 {
				t.Fatalf("recipe reach = %+v", divergences)
			}
		})
	}
	t.Run("chained invocations", func(t *testing.T) {
		_, divergences := gatePartsProjection(t, "verify:\n\t$(MAKE) prepare;$(MAKE) verify-go workspace\nprepare:\n\ttrue\nverify-go:\n\tgo test ./...\nworkspace:\n\tbun run verify\n", "make verify", nil)
		if len(divergences) != 0 {
			t.Fatalf("chained recipe reach = %+v", divergences)
		}
	})
	t.Run("named but not invoked", func(t *testing.T) {
		_, divergences := gatePartsProjection(t, "verify:\n\techo make verify-go\n\t$(OTHER) verify-go\n\tbun run verify\nverify-go:\n\tgo test ./...\n", "make verify", nil)
		if len(divergences) != 1 || divergences[0].Code != "verification.gate.part.missing" || divergences[0].ID != "verification.go" {
			t.Fatalf("uninvoked part = %+v", divergences)
		}
	})
	t.Run("mapped role", func(t *testing.T) {
		projections, divergences := gatePartsProjection(t, "verify:\n\t$(MAKE) check-go\n\tbun run verify\ncheck-go:\n\tgo test ./...\n", "make verify", map[string]string{"go": "make check-go"})
		part, found := findVerificationProjection(projections, "verification.go")
		if len(divergences) != 0 || !found || part.SatisfiedByCommand != "make check-go" {
			t.Fatalf("mapped part = %+v, divergences = %+v", part, divergences)
		}
	})
}

func TestAGateThatIsNotAMakeTargetIsNotChecked(t *testing.T) {
	t.Parallel()
	for _, gate := range []string{"bun run verify", "make absent"} {
		t.Run(gate, func(t *testing.T) {
			_, divergences := gatePartsProjection(t, "verify:\n\ttrue\nverify-go:\n\tgo test ./...\n", gate, nil)
			for _, divergence := range divergences {
				if divergence.Code == "verification.gate.part.missing" {
					t.Fatalf("non-Make or undeclared gate checked: %+v", divergences)
				}
			}
		})
	}
}

func TestMakeTargetReachStopsAtACycle(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	writeProfileAlignmentFile(t, repository, "Makefile", "first: second\n\tgo test ./...\nsecond: first\n\tbun run verify\nunused:\n\tfalse\n")
	data, err := os.ReadFile(filepath.Join(repository, "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	targets, recipes := makeTargetReach(data, "first")
	if len(targets) != 2 || len(recipes) != 2 || !gatePartReached("make second", targets, recipes) || !gatePartReached("rtk bun run verify", targets, recipes) {
		t.Fatalf("cycle reach: targets=%v recipes=%v", targets, recipes)
	}
}

func TestANonBooleanPartOfGateIsRefused(t *testing.T) {
	t.Parallel()
	for _, value := range []string{`"true"`, "null", "1", "{}"} {
		t.Run(value, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			replaceAsset(t, assets, "profiles/go-cli-typescript-monorepo.json", `"command": "bun run verify", "partOfGate": true`, `"command": "bun run verify", "partOfGate": `+value)
			_, err := LoadCatalog(assets)
			assertFrontendDiagnostic(t, err, "catalog.profile.verification.invalid")
		})
	}
}

func gatePartsProjection(t *testing.T, makefile, gate string, mappings map[string]string) ([]VerificationProjection, []ProfileDivergence) {
	t.Helper()
	repository := t.TempDir()
	writeProfileAlignmentFile(t, repository, "Makefile", makefile)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(true, true))
	root, err := os.OpenRoot(repository)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := root.Close(); err != nil {
			t.Error(err)
		}
	})
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile(repository, "go-cli-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	projections, divergences, err := resolveVerificationProjection(root, profile, []DecisionValue{{ID: "verification.gate", Value: gate}}, mappings, catalog)
	if err != nil {
		t.Fatal(err)
	}
	// Incremental is unrelated to gate parts and absent from these minimal fixtures.
	var relevant []ProfileDivergence
	for _, divergence := range divergences {
		if !strings.Contains(divergence.ID, "incremental") {
			relevant = append(relevant, divergence)
		}
	}
	return projections, relevant
}
