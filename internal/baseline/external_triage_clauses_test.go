package baseline

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

var externalTriageClauses = []struct{ id, force, guidance string }{
	{"clause.external-triage.classify-before-labelling", "mandatory", "Use this workflow only for issues and pull requests managed in an external forge. Classify each item as a bug or an enhancement, and state its user-visible problem and next action in English, before changing its labels or status."},
	{"clause.external-triage.move-through-mapped-states", "mandatory", "Move each item through the triage states needs-triage, needs-info, ready, and wontfix, applying the forge label the repository maps to each state. The repository records that mapping in this guide, outside its setup markers, when it enables external triage."},
	{"clause.external-triage.ask-for-an-unmapped-label", "stop-and-ask", "Stop and ask the maintainer for the forge label of a triage state the repository has not mapped; never invent a label."},
	{"clause.external-triage.needs-info-asks-and-stops", "mandatory", "When an item lacks what triage needs, mark it needs-info, ask the reporter specific questions that name each missing fact, and stop triaging it until the reporter answers."},
	{"clause.external-triage.disclose-ai-authorship", "mandatory", "Start every comment posted on the forge with a sentence stating that an AI agent generated it."},
	{"clause.external-triage.record-wontfix-decisions", "mandatory", "Record each wontfix decision as a declined Backlog Entry whose reason cites the forge item, and answer a repeated request with that recorded decision instead of deciding it again."},
	{"clause.external-triage.pull-request-is-an-issue-with-code", "mandatory", "Triage an external pull request as an issue with code attached: classify it and move it through the same states before reviewing or merging its code."},
	{"clause.external-triage.route-accepted-work-to-specs", "mandatory", "Route accepted work into the repository's local Spec workflow; a forge label or status never stands in for a Task's status."},
}

func externalTriageClauseFindings(module document, guide string) []string {
	var findings []string
	var clauses []document
	for _, rule := range objectsOrEmpty(module["rules"]) {
		if id, _ := stringValue(rule, "id"); id == "rule.external-triage" {
			clauses = objectsOrEmpty(rule["clauses"])
			if guidance, _ := stringValue(rule, "guidance"); guidance != "" {
				findings = append(findings, "rule.external-triage has unlabelled guidance")
			}
		}
	}
	if len(clauses) != len(externalTriageClauses) {
		findings = append(findings, fmt.Sprintf("external triage has %d clauses, want %d", len(clauses), len(externalTriageClauses)))
	}
	for index, want := range externalTriageClauses {
		if index >= len(clauses) {
			findings = append(findings, "missing clause "+want.id)
		} else {
			clause := clauses[index]
			id, _ := stringValue(clause, "id")
			force, _ := stringValue(clause, "enforcement")
			guidance, _ := stringValue(clause, "guidance")
			if id != want.id || force != want.force || guidance != want.guidance {
				findings = append(findings, "clause differs at "+want.id)
			}
			if index == 0 && !reflect.DeepEqual(clause["replaces"], []any{"rule.external-triage"}) {
				findings = append(findings, want.id+" lacks replacement declaration")
			}
		}
		line := "- **" + want.force + "**: " + want.guidance
		if !strings.Contains("\n"+guide+"\n", "\n"+line+"\n") {
			findings = append(findings, "guide lacks forced clause "+want.id)
		}
	}
	for _, line := range strings.Split(guide, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") &&
			!strings.HasPrefix(line, "- **mandatory**: ") &&
			!strings.HasPrefix(line, "- **prohibited**: ") &&
			!strings.HasPrefix(line, "- **stop-and-ask**: ") {
			findings = append(findings, "unlabelled bullet: "+line)
		}
	}
	return findings
}

func TestTheExternalTriageGuideStatesItsClausesWithForce(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("assets/modules/external-triage.json")
	if err != nil {
		t.Fatal(err)
	}
	var module document
	if err := json.Unmarshal(data, &module); err != nil {
		t.Fatal(err)
	}
	guide, err := os.ReadFile("assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/external-triage.md")
	if err != nil {
		t.Fatal(err)
	}
	if findings := externalTriageClauseFindings(module, string(guide)); len(findings) != 0 {
		t.Fatalf("external triage findings: %v", findings)
	}
}

func TestAMissingOrUnlabelledExternalTriageClauseIsReported(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		mutate func([]any, string) ([]any, string)
		want   string
	}{
		{"missing clause", func(clauses []any, guide string) ([]any, string) { return clauses[:len(clauses)-1], guide }, "missing clause clause.external-triage.route-accepted-work-to-specs"},
		{"changed force", func(clauses []any, guide string) ([]any, string) {
			clauses[1].(map[string]any)["enforcement"] = "prohibited"
			return clauses, guide
		}, "clause differs at clause.external-triage.move-through-mapped-states"},
		{"changed identifier", func(clauses []any, guide string) ([]any, string) {
			clauses[1].(map[string]any)["id"] = "clause.external-triage.wrong"
			return clauses, guide
		}, "clause differs at clause.external-triage.move-through-mapped-states"},
		{"changed order", func(clauses []any, guide string) ([]any, string) {
			clauses[1], clauses[2] = clauses[2], clauses[1]
			return clauses, guide
		}, "clause differs at clause.external-triage.move-through-mapped-states"},
		{"changed guidance", func(clauses []any, guide string) ([]any, string) {
			clauses[1].(map[string]any)["guidance"] = "Invent forge labels."
			return clauses, guide
		}, "clause differs at clause.external-triage.move-through-mapped-states"},
		{"missing replacement", func(clauses []any, guide string) ([]any, string) {
			delete(clauses[0].(map[string]any), "replaces")
			return clauses, guide
		}, "clause.external-triage.classify-before-labelling lacks replacement declaration"},
		{"missing guide clause", func(clauses []any, guide string) ([]any, string) {
			return clauses, strings.ReplaceAll(guide, "- **mandatory**: "+externalTriageClauses[0].guidance+"\n", "")
		}, "guide lacks forced clause clause.external-triage.classify-before-labelling"},
		{"unlabelled bullet", func(clauses []any, guide string) ([]any, string) { return clauses, guide + "- Invent forge labels.\n" }, "unlabelled bullet: - Invent forge labels."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var clauses []any
			var guide string
			for _, want := range externalTriageClauses {
				clauses = append(clauses, map[string]any{"id": want.id, "enforcement": want.force, "guidance": want.guidance})
				guide += "- **" + want.force + "**: " + want.guidance + "\n"
			}
			clauses[0].(map[string]any)["replaces"] = []any{"rule.external-triage"}
			module := document{"rules": []any{map[string]any{"id": "rule.external-triage", "clauses": clauses}}}
			if findings := externalTriageClauseFindings(module, guide); len(findings) != 0 {
				t.Fatalf("valid literal has findings: %v", findings)
			}
			clauses, guide = tc.mutate(clauses, guide)
			module["rules"] = []any{map[string]any{"id": "rule.external-triage", "clauses": clauses}}
			findings := externalTriageClauseFindings(module, guide)
			found := false
			for _, finding := range findings {
				if finding == tc.want {
					found = true
				}
			}
			if !found {
				t.Fatalf("findings = %v, want %q", findings, tc.want)
			}
		})
	}
}

func TestTheExternalTriageRuleIsReplacedForSourceBaselineAdopters(t *testing.T) {
	t.Parallel()
	request, catalog := newExternalTriageAdopter(t)
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("declared replacement refused refresh: %+v", outcome.Result)
	}
	if len(outcome.Plan.FileChanges) == 0 {
		t.Fatal("drifted external triage adopter has no file changes")
	}
	if delta := outcome.Plan.ClauseDelta; delta == nil || delta.Dispositions["rule.external-triage"] != ClauseReplaced {
		t.Fatalf("replacement delta = %+v", delta)
	}
	found := false
	for _, evidence := range outcome.Plan.Retention {
		if evidence.FromClause != "rule.external-triage" {
			continue
		}
		found = true
		successor := externalTriageClauses[0].id
		if len(evidence.Targets) != 1 || evidence.Targets[0] != successor || !strings.Contains(evidence.Reason, successor) {
			t.Fatalf("replacement evidence = %+v", evidence)
		}
	}
	if !found {
		t.Fatal("old external triage rule has no retention evidence")
	}
}

func TestAnUndeclaredExternalTriageReplacementIsRefused(t *testing.T) {
	t.Parallel()
	request, catalog := newExternalTriageAdopter(t)
	declared, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if declared.Plan == nil || declared.Plan.ClauseDelta == nil {
		t.Fatalf("declared replacement has no plan delta: %+v", declared.Result)
	}
	mutateCatalogClause(t, catalog, "external-triage", externalTriageClauses[0].id, func(clause document) bool { delete(clause, "replaces"); return true })
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" || outcome.Result.Category != "classification" {
		t.Fatalf("undeclared replacement outcome = %+v", outcome)
	}
	if delta := outcome.Result.ClauseDelta; delta == nil || delta.Dispositions["rule.external-triage"] != ClauseUnaccounted {
		t.Fatalf("unaccounted delta = %+v", delta)
	}
	if !strings.Contains(outcome.Result.Message, "rule.external-triage") {
		t.Fatalf("refusal does not name rule: %q", outcome.Result.Message)
	}
	if len(outcome.Result.ClauseDelta.Dispositions) != len(declared.Plan.ClauseDelta.Dispositions) {
		t.Fatal("replacement declaration changed the accounted entry set")
	}
	for id, disposition := range declared.Plan.ClauseDelta.Dispositions {
		if id != "rule.external-triage" && outcome.Result.ClauseDelta.Dispositions[id] != disposition {
			t.Errorf("replacement declaration changed %s disposition", id)
		}
	}
}

func newExternalTriageAdopter(t *testing.T) (PlanRequest, *Catalog) {
	t.Helper()
	request, catalog := newClauseReplacementAdopter(t)
	enabled := false
	for index := range request.Decisions {
		if request.Decisions[index].ID == "triage.external" {
			request.Decisions[index].Value = true
			enabled = true
		}
	}
	if !enabled {
		t.Fatal("adopter has no triage.external decision")
	}
	adoption, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if adoption.Plan == nil {
		t.Fatalf("enable external triage: %+v", adoption.Result)
	}
	if _, err := applyPlanWithCatalog(context.Background(), request.Repository, *adoption.Plan, adoption.Plan.PlanDigest, catalog); err != nil {
		t.Fatal(err)
	}
	manifest := adoption.Plan.SetupManifest
	for index := range manifest.ManagedArtifacts {
		manifest.ManagedArtifacts[index].Digest = strings.Repeat("0", 64)
	}
	data, err := marshalSetupManifestBytes(manifest)
	if err != nil {
		t.Fatal(err)
	}
	writeTransactionFile(t, request.Repository, manifestPath, string(data), 0o644)
	commitInspectionRepository(t, request.Repository, "age external triage adopter digests")
	return request, catalog
}
