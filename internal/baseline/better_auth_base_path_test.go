package baseline

import (
	"bytes"
	"context"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestTheBetterAuthBasePathIsAnOptionalDecision(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	declaration, exists := catalog.decisions["auth.base_path"]
	if !exists || !decisionOptional(declaration) || declaration["type"] != "string" || declaration["default"] != "/api/auth" {
		t.Fatalf("auth.base_path declaration = %#v, want optional string with /api/auth default", declaration)
	}
	if !slices.Equal(stringsOrEmpty(declaration["requiresCapabilities"]), []string{"capability.stack.better-auth"}) {
		t.Fatalf("auth.base_path capabilities = %v", declaration["requiresCapabilities"])
	}
	for _, id := range []string{"standard-typescript-monorepo", "go-cli-typescript-monorepo"} {
		t.Run(id, func(t *testing.T) {
			profile, err := ResolveProfile("", id, catalog)
			if err != nil {
				t.Fatal(err)
			}
			providerIndex := slices.Index(profile.Decisions, "auth.provider")
			if providerIndex < 0 || slices.Index(profile.Decisions, "auth.base_path") != providerIndex+1 {
				t.Fatalf("profile decisions = %v, want auth.base_path immediately after auth.provider", profile.Decisions)
			}
		})
	}
}

func TestARecordedBetterAuthBasePathMovesTheProviderException(t *testing.T) {
	t.Parallel()
	for _, scope := range []string{"/api/auth/*", "/v1/auth/*"} {
		t.Run(scope, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			provider := completeAuthProviderDecision()
			exception := provider["routeException"].(map[string]any)
			exception["scope"] = scope
			methods := "`GET` and `POST`"
			if scope == "/v1/auth/*" {
				exception["methods"] = []any{"GET"}
				exception["reason"] = "Keep the repository's recorded provider protocol."
				methods = "`GET`"
			}
			health := map[string]any{"scope": "/health", "methods": []any{"GET"}, "owner": "Operations", "reason": "Health checks."}
			source := map[string]any{"note": "Repository choice"}
			for i := range decisions {
				switch decisions[i].ID {
				case "auth.provider":
					decisions[i].Value = provider
				case "http.contract":
					decisions[i].Value = map[string]any{"mode": "Post-only", "exceptions": []any{cloneJSONValue(exception), health}, "source": source}
				}
			}
			decisions = append(decisions, DecisionValue{ID: "auth.base_path", Value: "/auth"})
			plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
			wantException := cloneJSONValue(exception).(map[string]any)
			wantException["scope"] = "/auth/*"
			for _, values := range [][]DecisionValue{plan.Decisions, betterAuthManifestDecisions(plan.SetupManifest)} {
				if got := planDecisionValue(t, values, "auth.base_path"); got != "/auth" {
					t.Errorf("base path = %v", got)
				}
				gotProvider := planDecisionValue(t, values, "auth.provider").(map[string]any)
				if !reflect.DeepEqual(gotProvider["routeException"], wantException) {
					t.Errorf("provider exception = %#v, want %#v", gotProvider["routeException"], wantException)
				}
				gotHTTP := planDecisionValue(t, values, "http.contract").(map[string]any)
				wantHTTP := map[string]any{"mode": "Post-only", "exceptions": []any{wantException, health}, "source": source}
				if !reflect.DeepEqual(gotHTTP, wantHTTP) {
					t.Errorf("HTTP contract = %#v, want %#v", gotHTTP, wantHTTP)
				}
			}
			backend := string(planPostimage(t, plan, "docs/agents/backend.md").Content)
			for _, sentence := range []string{
				"**Better Auth** owns " + methods + " for `/auth/*`:",
				"Better Auth owns the authentication protocol for `/auth/*`.",
				"\n\nThe recorded Better Auth base path is `/auth`, the value of Better Auth's `basePath` option.",
			} {
				if !strings.Contains(backend, sentence) {
					t.Errorf("backend missing %q:\n%s", sentence, backend)
				}
			}
			if strings.Count(backend, "`/auth/*`") != 2 || strings.Contains(backend, "`"+scope+"`") {
				t.Errorf("backend contains stale or duplicate scopes:\n%s", backend)
			}
			// Normalizing a recorded projection again must keep the one moved exception.
			again := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), betterAuthManifestDecisions(plan.SetupManifest))
			if !reflect.DeepEqual(again.Decisions, plan.Decisions) {
				t.Fatal("recorded decisions changed on replan")
			}
		})
	}
}

func TestAnUnrecordedBetterAuthBasePathKeepsTheRecordedScope(t *testing.T) {
	t.Parallel()
	for _, scope := range []string{"/api/auth/*", "/v1/auth/*"} {
		t.Run(scope, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			for i := range decisions {
				switch decisions[i].ID {
				case "auth.provider":
					provider := completeAuthProviderDecision()
					provider["routeException"].(map[string]any)["scope"] = scope
					decisions[i].Value = provider
				case "http.contract":
					decisions[i].Value = map[string]any{"mode": "Post-only"}
				}
			}
			plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
			for _, values := range [][]DecisionValue{plan.Decisions, betterAuthManifestDecisions(plan.SetupManifest)} {
				if slices.ContainsFunc(values, func(d DecisionValue) bool { return d.ID == "auth.base_path" }) {
					t.Fatal("unasked base path was recorded")
				}
				provider := planDecisionValue(t, values, "auth.provider").(map[string]any)
				if provider["routeException"].(map[string]any)["scope"] != scope {
					t.Fatal("recorded provider scope changed")
				}
				contract := planDecisionValue(t, values, "http.contract").(map[string]any)
				exceptions := contract["exceptions"].([]any)
				if len(exceptions) != 1 || exceptions[0].(map[string]any)["scope"] != scope {
					t.Fatalf("HTTP exceptions = %#v", exceptions)
				}
			}
			backend := planPostimage(t, plan, "docs/agents/backend.md").Content
			if strings.Contains(string(backend), "The recorded Better Auth base path") || strings.Count(string(backend), "`"+scope+"`") != 2 {
				t.Fatalf("unrecorded backend = %s", backend)
			}
			if scope == "/api/auth/*" {
				golden, ok := mustEmbeddedCatalog(t).Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md")
				if !ok || !bytes.Equal(backend, golden.Data) {
					t.Fatal("unrecorded backend differs from the existing golden")
				}
			}
		})
	}
}

func TestAnInvalidBetterAuthBasePathIsRefused(t *testing.T) {
	t.Parallel()
	const invalid = "auth.base_path must be an absolute path with at least one segment and no trailing slash, wildcard, query, fragment or whitespace"
	tests := []struct {
		name            string
		value           any
		withoutProvider bool
		message         string
	}{
		{"relative", "auth", false, invalid},
		{"root", "/", false, invalid},
		{"trailing slash", "/auth/", false, invalid},
		{"wildcard", "/auth/*", false, invalid},
		{"space", "/a b", false, invalid},
		{"double slash", "/a//b", false, invalid},
		{"query", "/a?x", false, invalid},
		{"fragment", "/auth#x", false, invalid},
		{"leading whitespace", " /auth", false, invalid},
		{"unicode whitespace", "/a\u00a0b", false, invalid},
		{"backtick", "/a`b", false, invalid},
		{"control", "/a\x00b", false, invalid},
		{"marker", "/<!--setup-context-driven:begin-->", false, invalid},
		{"token", "/{{auth}}", false, invalid},
		{"non-string", 42, false, invalid},
		{"without provider", "/auth", true, "auth.base_path requires auth.provider"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			if test.withoutProvider {
				decisions = slices.DeleteFunc(decisions, func(d DecisionValue) bool { return d.ID == "auth.provider" })
			}
			decisions = append(decisions, DecisionValue{ID: "auth.base_path", Value: test.value})
			// BuildPlan asks for missing required decisions and validates scalar types
			// before invoking the project normalizer; exercise its refusal directly.
			if test.withoutProvider || test.name == "non-string" {
				_, err := normalizeProjectDecisions(decisions, mustEmbeddedCatalog(t))
				if err == nil || !strings.Contains(err.Error(), test.message) {
					t.Fatalf("project normalization error = %v, want %q", err, test.message)
				}
				return
			}
			outcome, err := BuildPlan(context.Background(), PlanRequest{
				Repository: newProjectDecisionPlanRepository(t), ProfileID: "standard-typescript-monorepo", Decisions: decisions,
				Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
			})
			if err == nil || outcome.Plan != nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("invalid base path produced Plan=%v error=%v, want %q", outcome.Plan != nil, err, test.message)
			}
		})
	}
}

func betterAuthManifestDecisions(manifest SetupManifest) []DecisionValue {
	decisions := make([]DecisionValue, 0, len(manifest.Decisions))
	for id, decision := range manifest.Decisions {
		decisions = append(decisions, DecisionValue{ID: id, Value: decision.Value})
	}
	slices.SortFunc(decisions, func(a, b DecisionValue) int { return strings.Compare(a.ID, b.ID) })
	return decisions
}
