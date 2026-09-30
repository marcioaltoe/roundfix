// Suite: Workspace scope and HTTP suggestions
// Invariant: guides name their workspace and HTTP suggestions preserve repository decisions.
// Boundary IN: embedded catalog, Plan postimages, and Setup Manifest resolution.
// Boundary OUT: applying a Plan and executing repository Verification.

package baseline

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func profileHTTPDefaultFindings(profile, decision document) []string {
	contract, _ := objectValue(profile["httpContract"])
	suggestion, _ := objectValue(decision["default"])
	mode, _ := suggestion["mode"].(string)
	if mode == "" || contract["default"] != mode {
		return []string{fmt.Sprintf("%v: HTTP default %v disagrees with catalog default %q", profile["id"], contract["default"], mode)}
	}
	return nil
}

func TestTheBackendAndFrontendGuidesSayWhatTheyGovern(t *testing.T) {
	t.Parallel()
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
	rules := []stackWordingRule{
		{guide: "docs/agents/backend.md", must: []string{"These rules govern the repository's TypeScript backend workspace. A service or command written in another language follows its own guide."}},
		{guide: "docs/agents/frontend.md", must: []string{"These rules govern the repository's web frontend workspace. A terminal interface follows its own guide."}},
	}
	rendered := make(map[string]string)
	for _, rule := range rules {
		rendered[rule.guide] = string(planPostimage(t, plan, rule.guide).Content)
		text := strings.Join(strings.Fields(rendered[rule.guide]), " ")
		scopeIndex := strings.Index(text, rule.must[0])
		ruleIndex := strings.Index(text, "- **mandatory**:")
		if scopeIndex < 0 || ruleIndex < 0 || scopeIndex >= ruleIndex {
			t.Errorf("%s: scope must precede rules", rule.guide)
		}
	}
	if findings := stackWordingFindings(rendered, rules); len(findings) != 0 {
		t.Fatalf("scope findings: %v", findings)
	}
}

func TestTheProfileHTTPDefaultMatchesTheDecisionCatalog(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	if findings := profileHTTPDefaultFindings(catalog.profiles["standard-typescript-monorepo"], catalog.decisions[httpContractDecisionID]); len(findings) != 0 {
		t.Fatalf("HTTP default findings: %v", findings)
	}
}

func TestAProfileHTTPDefaultThatDisagreesIsReported(t *testing.T) {
	t.Parallel()
	profile := document{"id": "disagreeing-profile", "httpContract": map[string]any{"default": "Post-only"}}
	decision := document{"default": map[string]any{"mode": "REST"}}
	want := []string{`disagreeing-profile: HTTP default Post-only disagrees with catalog default "REST"`}
	if findings := profileHTTPDefaultFindings(profile, decision); !reflect.DeepEqual(findings, want) {
		t.Fatalf("findings = %v, want %v", findings, want)
	}
}

func TestARecordedHTTPContractDecisionSurvivesAnUpdate(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, mode := range []string{"Post-only", "REST"} {
		t.Run(mode, func(t *testing.T) {
			repository := t.TempDir()
			decisions := manifestInputDecisionsForProfile(t, repository, catalog, "standard-typescript-monorepo")
			contract, ok := objectValue(decisions[httpContractDecisionID])
			if !ok {
				t.Fatal("fixture HTTP contract is not an object")
			}
			contract["mode"] = mode
			exceptions, ok := contract["exceptions"].([]any)
			if !ok {
				t.Fatal("fixture HTTP exceptions are not an array")
			}
			contract["exceptions"] = append(exceptions, map[string]any{
				"scope": "/health", "methods": []any{"GET"},
				"owner": "Operations", "reason": "Expose repository health checks.",
			})
			manifest := newManifestInputFixture(t, repository, catalog, "standard-typescript-monorepo", decisions)
			writeManifestInputFixture(t, repository, manifest)
			before := snapshotManifestInputRepository(t, repository)
			input, err := ResolveManifestInput(repository, catalog)
			if err != nil {
				t.Fatal(err)
			}
			if input.State != ManifestInputResolved || len(input.NewDecisions) != 0 {
				t.Fatalf("state = %q, new decisions = %v", input.State, input.NewDecisions)
			}
			assertManifestInputDecisionsEqual(t, input.Decisions, manifest.Decisions)
			index := decisionValueIndex(input.Decisions, httpContractDecisionID)
			if index < 0 {
				t.Fatal("recorded HTTP decision is absent")
			}
			resolved, ok := objectValue(input.Decisions[index].Value)
			if !ok || resolved["mode"] != mode || !reflect.DeepEqual(resolved["exceptions"], contract["exceptions"]) {
				t.Fatalf("resolved HTTP decision = %#v, want mode %q and exceptions %#v", resolved, mode, contract["exceptions"])
			}
			if after := snapshotManifestInputRepository(t, repository); !reflect.DeepEqual(after, before) {
				t.Fatal("manifest resolution wrote repository files")
			}
		})
	}
}

func TestAnAbsentHTTPContractDecisionIsSuggestedAsREST(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	repository := t.TempDir()
	decisions := manifestInputDecisionsForProfile(t, repository, catalog, "standard-typescript-monorepo")
	delete(decisions, httpContractDecisionID)
	manifest := newManifestInputFixture(t, repository, catalog, "standard-typescript-monorepo", decisions)
	writeManifestInputFixture(t, repository, manifest)
	before := snapshotManifestInputRepository(t, repository)
	input, err := ResolveManifestInput(repository, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if input.State != ManifestInputIncomplete || len(input.NewDecisions) != 1 || input.NewDecisions[0].ID != httpContractDecisionID {
		t.Fatalf("state = %q, new decisions = %v", input.State, input.NewDecisions)
	}
	suggestion, ok := objectValue(input.NewDecisions[0].SuggestedValue)
	if !ok || suggestion["mode"] != "REST" {
		t.Fatalf("HTTP suggestion = %#v, want REST", suggestion)
	}
	if decisionValueIndex(input.Decisions, httpContractDecisionID) >= 0 {
		t.Fatal("manifest resolution adopted the HTTP suggestion")
	}
	if after := snapshotManifestInputRepository(t, repository); !reflect.DeepEqual(after, before) {
		t.Fatal("manifest resolution wrote repository files")
	}
}
