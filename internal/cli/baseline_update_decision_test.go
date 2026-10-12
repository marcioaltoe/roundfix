// Suite: declared decisions in a manifest-driven Baseline update
// Invariant: declarations are normalized and digest-bound before a confirmed managed refresh; invalid declarations never write.
// Boundary IN: real temporary adopters, CLI parsing, previews, confirmation, guide and manifest bytes, and refusal categories.
// Boundary OUT: external skill acquisition and Pending History conversion.
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineUpdateRecordsADeclaredDecision(t *testing.T) {
	t.Parallel()
	for _, missing := range []bool{false, true} {
		name := "changed recorded decision"
		if missing {
			name = "new required decision"
		}
		t.Run(name, func(t *testing.T) {
			repo := newBaselineUpdateRepository(t)
			manifest := ReadBaselineSetupManifest(t, repo)
			recorded := manifest.Decisions["secondbrain.enabled"].Value
			if missing {
				delete(manifest.Decisions, "secondbrain.enabled")
				WriteBaselineSetupManifest(t, repo, manifest)
				commitBaselinePlanTestRepository(t, repo)
				recorded = nil
			}
			before := baselinePlanTestTree(t, repo)
			args := []string{"baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--decision", "secondbrain.enabled=true", "--decision", "branch.prefix=<type>/"}
			preview, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), append(args, "--format=json")...)
			if code != exitUnverified || stderr != "" || preview.State != "plan_ready" || preview.PlanDigest == "" {
				t.Fatalf("preview exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			var payload struct {
				DeclaredDecisions []struct {
					ID       string
					Value    any
					Recorded any
				}
			}
			if err := json.Unmarshal([]byte(out), &payload); err != nil {
				t.Fatal(err)
			}
			if len(payload.DeclaredDecisions) != 2 {
				t.Fatalf("declarations = %+v", payload.DeclaredDecisions)
			}
			found := false
			for _, decision := range payload.DeclaredDecisions {
				if decision.ID == "secondbrain.enabled" {
					found = true
					if decision.Value != true || !reflect.DeepEqual(decision.Recorded, recorded) {
						t.Fatalf("declaration = %+v, recorded=%v", decision, recorded)
					}
				}
			}
			if !found || len(preview.NewDecisions) != 0 || len(preview.AdoptedSuggestions) != 0 {
				t.Fatalf("preview = %+v", preview)
			}
			_, text, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), append(args, "--format=text")...)
			prior := "false"
			if missing {
				prior = "none"
			}
			for _, want := range []string{"Declared decisions: 2", "- secondbrain.enabled=true (recorded: " + prior + ")", `- branch.prefix="\u003ctype\u003e/" (recorded: "ma/")`} {
				if !strings.Contains(text, want) {
					t.Errorf("text missing %q:\n%s", want, text)
				}
			}
			if code != exitUnverified || stderr != "" || before != baselinePlanTestTree(t, repo) {
				t.Fatal("preview wrote or did not present the plan")
			}

			changedArgs := []string{"baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--decision", "secondbrain.enabled=false", "--decision", "branch.prefix=<type>/", "--format=json"}
			changed, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), changedArgs...)
			if stderr != "" || code != exitUnverified || changed.PlanDigest == preview.PlanDigest {
				t.Fatalf("different declaration digest exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			refused, out, _, code := runBaselineUpdateTestCommand(t, t.Context(), append(changedArgs, "--confirm-plan", preview.PlanDigest)...)
			if code != exitUnverified || refused.Category != "approval" || before != baselinePlanTestTree(t, repo) {
				t.Fatalf("changed declaration accepted old digest: exit=%d stdout=%s", code, out)
			}
			applied, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), append(args, "--format=json", "--confirm-plan", preview.PlanDigest)...)
			if code != exitOK || stderr != "" || applied.State != "verified" || applied.ApprovedPlanDigest != preview.PlanDigest {
				t.Fatalf("confirmed declaration exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			after := ReadBaselineSetupManifest(t, repo)
			if after.Decisions["secondbrain.enabled"].Value != true || after.Decisions["branch.prefix"].Value != "<type>/" || after.Profile != manifest.Profile {
				t.Fatalf("written manifest = %+v", after)
			}
		})
	}
	t.Run("remaining suggestions", func(t *testing.T) {
		repo := incompleteBaselineUpdateRepository(t, "secondbrain.enabled")
		manifest := ReadBaselineSetupManifest(t, repo)
		delete(manifest.Decisions, "triage.external")
		WriteBaselineSetupManifest(t, repo, manifest)
		commitBaselinePlanTestRepository(t, repo)
		before := baselinePlanTestTree(t, repo)
		args := []string{"baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--decision", "triage.external=false", "--format=json"}
		result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), args...)
		if code != exitUnverified || stderr != "" || result.Category != "decision" || len(result.NewDecisions) != 1 || result.NewDecisions[0].ID != "secondbrain.enabled" || len(result.DeclaredDecisions) != 1 || result.DeclaredDecisions[0].Recorded != nil || before != baselinePlanTestTree(t, repo) {
			t.Fatalf("remaining decisions exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
		result, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), append(args, "--adopt-suggested", "--yes")...)
		if code != exitOK || stderr != "" || len(result.AdoptedSuggestions) != 1 || result.AdoptedSuggestions[0].ID != "secondbrain.enabled" || ReadBaselineSetupManifest(t, repo).Decisions["triage.external"].Value != false {
			t.Fatalf("remaining suggestion adoption exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
	})

	t.Run("repository profile", func(t *testing.T) {
		repo := newBaselineRepositoryProfileFixture(t, false)
		result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--decision", "branch.prefix=<type>/", "--yes", "--format=json")
		if code != exitOK || stderr != "" || result.State != "verified" || ReadBaselineSetupManifest(t, repo).Profile != baselineRepositoryProfileTestID {
			t.Fatalf("repository profile declaration exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
	})
	t.Run("empty array without declarations", func(t *testing.T) {
		repo := newBaselineUpdateRepository(t)
		_, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--format=json")
		var payload map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &payload); err != nil {
			t.Fatal(err)
		}
		if code != exitOK || stderr != "" || string(payload["declaredDecisions"]) != "[]" {
			t.Fatalf("empty declarations exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
	})
	t.Run("help", func(t *testing.T) {
		for _, args := range [][]string{{"--help"}, {"baseline", "--help"}, {"baseline", "update", "--help"}} {
			var out, stderr bytes.Buffer
			if code := RunContext(t.Context(), args, &out, &stderr); code != exitOK || stderr.Len() != 0 || !strings.Contains(out.String(), "[--adopt-suggested] [--decision <id=value> ...]") {
				t.Fatalf("help %v: stdout=%s stderr=%s", args, out.String(), stderr.String())
			}
		}
		if !strings.Contains(commandUsage("baseline update"), "Record or change a decision as id=value; repeat for several decisions") {
			t.Fatal("help lacks decision option")
		}
	})
}

func TestBaselineUpdateDeclaresTheBetterAuthBasePath(t *testing.T) {
	t.Parallel()
	repo := newBaselineDeclaredTypeScriptRepository(t)
	before := baselinePlanTestTree(t, repo)
	updateArgs := []string{"baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--decision", "auth.base_path=/auth", "--format=json"}
	preview, previewOut, previewErr, code := runBaselineUpdateTestCommand(t, t.Context(), updateArgs...)
	if code != exitUnverified || previewErr != "" || preview.State != "plan_ready" || before != baselinePlanTestTree(t, repo) {
		t.Fatalf("base path preview exit=%d stdout=%s stderr=%s", code, previewOut, previewErr)
	}
	if !strings.Contains(previewOut, `"id":"auth.base_path","value":"/auth","recorded":null`) {
		t.Fatalf("base path declaration missing: %s", previewOut)
	}
	result, updateOut, updateErr, code := runBaselineUpdateTestCommand(t, t.Context(), append(updateArgs, "--yes")...)
	if code != exitOK || updateErr != "" || result.State != "verified" || result.PlanDigest != preview.PlanDigest {
		t.Fatalf("base path apply exit=%d stdout=%s stderr=%s", code, updateOut, updateErr)
	}
	manifest := ReadBaselineSetupManifest(t, repo)
	if manifest.Profile != "standard-typescript-monorepo" || manifest.Decisions["auth.base_path"].Value != "/auth" {
		t.Fatalf("base path manifest=%+v", manifest)
	}
	provider, ok := manifest.Decisions["auth.provider"].Value.(map[string]any)
	if !ok {
		t.Fatalf("provider=%v", manifest.Decisions["auth.provider"].Value)
	}
	exception, ok := provider["routeException"].(map[string]any)
	if !ok || exception["scope"] != "/auth/*" {
		t.Fatalf("provider exception=%v", provider)
	}
	contract, ok := manifest.Decisions["http.contract"].Value.(map[string]any)
	if !ok {
		t.Fatalf("contract=%v", manifest.Decisions["http.contract"].Value)
	}
	exceptions, ok := contract["exceptions"].([]any)
	if !ok || len(exceptions) != 1 || !reflect.DeepEqual(exceptions[0], exception) {
		t.Fatalf("HTTP exceptions=%v", contract)
	}
	guide, err := os.ReadFile(filepath.Join(repo, "docs/agents/backend.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"`/auth/*`", "The recorded Better Auth base path is `/auth`, the value of Better Auth's `basePath` option."} {
		if !bytes.Contains(guide, []byte(want)) {
			t.Errorf("backend guide missing %q:\n%s", want, guide)
		}
	}
	if bytes.Contains(guide, []byte("/api/auth/*")) || bytes.Count(guide, []byte("`/auth/*`")) != 2 {
		t.Fatalf("backend has stale or duplicate scopes:\n%s", guide)
	}
	current, currentOut, currentErr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--format=json")
	if code != exitOK || currentErr != "" || current.State != "current" || len(current.FileChanges) != 0 {
		t.Fatalf("second update exit=%d stdout=%s stderr=%s", code, currentOut, currentErr)
	}
}

func newBaselineDeclaredTypeScriptRepository(t *testing.T) string {
	t.Helper()
	repo := newBaselineReleaseRepository(t, "standard-typescript-monorepo")
	args := []string{"baseline", "plan", "--repo", repo, "--profile", "standard-typescript-monorepo", "--format=json"}
	for _, decision := range baselineReleaseDecisionArgs("standard-typescript-monorepo", "greenfield") {
		args = append(args, "--decision", decision)
	}
	var out, stderr bytes.Buffer
	if code := RunContext(t.Context(), args, &out, &stderr); code != exitOK {
		t.Fatalf("adoption plan exit=%d stdout=%s stderr=%s", code, out.String(), stderr.String())
	}
	plan, err := baseline.ParsePlanDocument(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := baseline.ApplyPlan(t.Context(), repo, plan, plan.PlanDigest); err != nil {
		t.Fatal(err)
	}
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func TestBaselineUpdateRefusesAnUndeclarableDecision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		decisions []string
		fixed     bool
	}{
		{name: "malformed pair", decisions: []string{"secondbrain.enabled"}},
		{name: "empty id", decisions: []string{"=true"}},
		{name: "empty value", decisions: []string{"secondbrain.enabled="}},
		{name: "unknown decision", decisions: []string{"unknown.decision=true"}},
		{name: "unselected decision", decisions: []string{"auth.base_path=/auth"}},
		{name: "profile-fixed decision", decisions: []string{"branch.prefix=<type>/"}, fixed: true},
		{name: "same profile-fixed value", decisions: []string{"branch.prefix=ma/"}, fixed: true},
		{name: "preservation mode", decisions: []string{"preservation.mode=greenfield"}},
		{name: "invalid typed value", decisions: []string{"secondbrain.enabled=perhaps"}},
		{name: "duplicate declaration", decisions: []string{"secondbrain.enabled=true", "secondbrain.enabled=false"}},
		{name: "atomic refusal", decisions: []string{"branch.prefix=<type>/", "unknown.decision=true"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newBaselineUpdateRepository(t)
			if test.fixed {
				catalog, err := baseline.LoadEmbeddedCatalog()
				if err != nil {
					t.Fatal(err)
				}
				initialized, err := baseline.InitCustomProfile(repo, "repository-fixed", "go-cli-tui", catalog)
				if err != nil {
					t.Fatal(err)
				}
				var profile map[string]any
				data, err := os.ReadFile(initialized.Path)
				if err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(data, &profile); err != nil {
					t.Fatal(err)
				}
				profile["values"] = map[string]any{"branch.prefix": "ma/"}
				data, err = json.Marshal(profile)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(initialized.Path, data, 0o644); err != nil {
					t.Fatal(err)
				}
				manifest := ReadBaselineSetupManifest(t, repo)
				manifest.Profile = initialized.Profile.ID
				manifest.Generator.Baseline = "baseline." + initialized.Profile.ID + "-" + baseline.ManifestVersion
				WriteBaselineSetupManifest(t, repo, manifest)
				commitBaselinePlanTestRepository(t, repo)
			}
			before := baselinePlanTestTree(t, repo)
			args := []string{"baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--yes", "--format=json"}
			for _, decision := range test.decisions {
				args = append(args, "--decision", decision)
			}
			result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), args...)
			if code != exitPreflight || result.Category != "invalid" || result.OK || stderr == "" {
				t.Fatalf("refusal exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
			var payload map[string]json.RawMessage
			if err := json.Unmarshal([]byte(out), &payload); err != nil {
				t.Fatal(err)
			}
			if string(payload["declaredDecisions"]) != "[]" {
				t.Fatalf("refusal declaration array=%s", payload["declaredDecisions"])
			}
			if before != baselinePlanTestTree(t, repo) {
				t.Fatal("invalid declaration wrote repository bytes")
			}
		})
	}
	for _, value := range []string{"auth", "/", "/auth/", "/auth/*", "/a b", "/a//b", "/a?x"} {
		t.Run("invalid base path "+value, func(t *testing.T) {
			repo := newBaselineDeclaredTypeScriptRepository(t)
			before := baselinePlanTestTree(t, repo)
			result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--no-skills", "--no-history", "--decision", "auth.base_path="+value, "--yes", "--format=json")
			if code != exitPreflight || result.Category != "invalid" || !strings.Contains(stderr, "auth.base_path must be an absolute path") || before != baselinePlanTestTree(t, repo) {
				t.Fatalf("invalid base path exit=%d stdout=%s stderr=%s", code, out, stderr)
			}
		})
	}

}
