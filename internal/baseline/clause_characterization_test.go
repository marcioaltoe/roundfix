package baseline

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type characterizedBaselineClause struct {
	ID          string `json:"id"`
	Guidance    string `json:"guidance"`
	Enforcement string `json:"enforcement"`
	RuleLevel   bool   `json:"-"`
}

func TestNoTwoBaselineClausesShareText(t *testing.T) {
	if pairs := repeatedBaselineGuidance(embeddedBaselineClauses(t)); len(pairs) != 0 {
		t.Fatalf("repeated Baseline guidance: %v", pairs)
	}
}

func TestADuplicatedClauseTextIsReported(t *testing.T) {
	clauses := embeddedBaselineClauses(t)
	// Change a copy, including case and whitespace, without mutating the catalog.
	clauses[1].Guidance = " \n" + strings.ToUpper(strings.ReplaceAll(clauses[0].Guidance, " ", "\t  ")) + " \n"
	want := [][2]string{{clauses[0].ID, clauses[1].ID}}
	if got := repeatedBaselineGuidance(clauses); !reflect.DeepEqual(got, want) {
		t.Fatalf("repeated guidance pairs = %v, want %v", got, want)
	}
}

func TestBaselineClauseForceIsCharacterized(t *testing.T) {
	if issues := baselineClauseForceDifferences(embeddedBaselineClauses(t), characterizedBaselineForce()); len(issues) != 0 {
		t.Fatalf("Baseline clause force changed: %v", issues)
	}
}

func TestAMissingBaselineClauseIsReported(t *testing.T) {
	clauses := embeddedBaselineClauses(t)
	want := []string{"missing clause " + clauses[0].ID}
	if got := baselineClauseForceDifferences(clauses[1:], characterizedBaselineForce()); !reflect.DeepEqual(got, want) {
		t.Fatalf("force differences = %v, want %v", got, want)
	}
}

func TestAnUnlistedBaselineClauseIsReported(t *testing.T) {
	clauses := append(embeddedBaselineClauses(t), characterizedBaselineClause{ID: "clause.test.unlisted", Enforcement: "mandatory"})
	want := []string{"unlisted clause clause.test.unlisted"}
	if got := baselineClauseForceDifferences(clauses, characterizedBaselineForce()); !reflect.DeepEqual(got, want) {
		t.Fatalf("force differences = %v, want %v", got, want)
	}
}

func TestAChangedBaselineClauseForceIsReported(t *testing.T) {
	clauses := embeddedBaselineClauses(t)
	prior := clauses[0].Enforcement
	clauses[0].Enforcement = "prohibited"
	if prior == clauses[0].Enforcement {
		clauses[0].Enforcement = "mandatory"
	}
	want := []string{fmt.Sprintf("clause %s force = %s, want %s", clauses[0].ID, clauses[0].Enforcement, prior)}
	if got := baselineClauseForceDifferences(clauses, characterizedBaselineForce()); !reflect.DeepEqual(got, want) {
		t.Fatalf("force differences = %v, want %v", got, want)
	}
}

func TestTheBackendBoundaryParagraphRendersOnce(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	const assetPath = "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md"
	asset, ok := catalog.Asset(assetPath)
	if !ok {
		t.Fatalf("embedded catalog has no asset %q", assetPath)
	}
	const sentence = "Keep blocking, network, process, database, and daemon boundaries explicit about ownership, cancellation, timeouts, and error reporting."
	if count := strings.Count(string(asset.Data), sentence); count != 1 {
		t.Fatalf("backend boundary sentence occurs %d times, want 1", count)
	}
}

func embeddedBaselineClauses(t *testing.T) []characterizedBaselineClause {
	t.Helper()
	catalog := mustEmbeddedCatalog(t)
	var clauses []characterizedBaselineClause
	for _, moduleID := range catalog.ModuleIDs() {
		module, ok := catalog.Module(moduleID)
		if !ok {
			t.Fatalf("embedded catalog has no module %q", moduleID)
		}
		var document struct {
			Rules []struct {
				characterizedBaselineClause
				Clauses []characterizedBaselineClause `json:"clauses"`
			} `json:"rules"`
		}
		if err := json.Unmarshal(module.Data, &document); err != nil {
			t.Fatalf("decode module %s: %v", moduleID, err)
		}
		for _, rule := range document.Rules {
			if rule.Guidance != "" {
				rule.RuleLevel = true
				clauses = append(clauses, rule.characterizedBaselineClause)
			}
			clauses = append(clauses, rule.Clauses...)
		}
	}
	return clauses
}

func repeatedBaselineGuidance(clauses []characterizedBaselineClause) [][2]string {
	seen := make(map[string]string)
	var pairs [][2]string
	for _, clause := range clauses {
		text := strings.Join(strings.Fields(strings.ToLower(clause.Guidance)), " ")
		if prior, ok := seen[text]; ok {
			pairs = append(pairs, [2]string{prior, clause.ID})
		} else {
			seen[text] = clause.ID
		}
	}
	return pairs
}

func baselineClauseForceDifferences(clauses []characterizedBaselineClause, want map[string]string) []string {
	seen := make(map[string]bool)
	var issues []string
	for _, clause := range clauses {
		// Rule-level guidance has no enforcement field; only clauses carry force.
		if clause.RuleLevel {
			continue
		}
		seen[clause.ID] = true
		force, ok := want[clause.ID]
		if !ok {
			issues = append(issues, "unlisted clause "+clause.ID)
		} else if clause.Enforcement != force {
			issues = append(issues, fmt.Sprintf("clause %s force = %s, want %s", clause.ID, clause.Enforcement, force))
		}
	}
	for id := range want {
		if !seen[id] {
			issues = append(issues, "missing clause "+id)
		}
	}
	sort.Strings(issues)
	return issues
}

// This literal records every clause at 9e439dbb except the duplicate backend
// entry. Wording may change without changing identifiers or enforcement.
func characterizedBaselineForce() map[string]string {
	return map[string]string{
		"clause.go.change-module-files-through-go":                    "prohibited",
		"clause.go.keep-entry-points-thin":                            "mandatory",
		"clause.go.record-the-reason-for-a-module":                    "mandatory",
		"clause.go.own-every-goroutine":                               "mandatory",
		"clause.go.pass-context-first":                                "mandatory",
		"clause.go.wrap-errors-with-the-operation":                    "mandatory",
		"clause.go.build-every-constrained-platform":                  "mandatory",
		"clause.go.test-observable-behavior":                          "mandatory",
		"clause.go.test-through-go-test":                              "mandatory",
		"clause.cli.public-command-contract":                          "mandatory",
		"clause.cli.separate-output-streams":                          "mandatory",
		"clause.cli.ship-the-skill-with-the-behavior":                 "mandatory",
		"clause.cli.deterministic-non-interactive":                    "mandatory",
		"clause.cli.explicit-safe-writes":                             "mandatory",
		"clause.rust.change-manifests-through-cargo":                  "mandatory",
		"clause.rust.read-current-docs":                               "mandatory",
		"clause.rust.iterate-with-focused-cargo-commands":             "mandatory",
		"clause.rust.test-the-command-through-execution":              "mandatory",
		"clause.rust.keep-the-binary-thin":                            "mandatory",
		"clause.rust.keep-type-erased-errors-at-the-entry-point":      "prohibited",
		"clause.rust.prohibit-panics-on-user-paths":                   "prohibited",
		"clause.rust.return-typed-errors":                             "mandatory",
		"clause.tui.drive-models-synchronously":                       "mandatory",
		"clause.tui.emulate-the-terminal-last":                        "mandatory",
		"clause.tui.keep-design-policy-in-repository-guidance":        "mandatory",
		"clause.core.regenerate-generated-files":                      "prohibited",
		"clause.core.flaky-tests-block":                               "mandatory",
		"clause.core.lint-warnings-block":                             "mandatory",
		"clause.core.prohibit-test-only-production-hooks":             "prohibited",
		"clause.core.prohibit-editing-vendored-skills":                "prohibited",
		"clause.core.ask-the-person-to-start-a-person-only-skill":     "mandatory",
		"clause.core.ask-before-database-mutation":                    "stop-and-ask",
		"clause.core.prove-the-mutation-predicate":                    "mandatory",
		"clause.core.prohibit-disguised-database-mutation":            "prohibited",
		"clause.autonomous.daemon-verifies-task":                      "mandatory",
		"clause.autonomous.delegate-through-roundfix":                 "mandatory",
		"clause.autonomous.hook-strictness":                           "mandatory",
		"clause.autonomous.loop-01-qa-once":                           "mandatory",
		"clause.autonomous.loop-02-autonomous-decomposition":          "mandatory",
		"clause.autonomous.loop-03-specify-the-hazard":                "mandatory",
		"clause.autonomous.loop-04-verify-the-class":                  "mandatory",
		"clause.autonomous.loop-05-clean-is-not-evidence":             "mandatory",
		"clause.autonomous.loop-06-branch-hygiene":                    "mandatory",
		"clause.autonomous.loop-07-outlive-the-turn":                  "mandatory",
		"clause.autonomous.supervisor-prohibition":                    "prohibited",
		"clause.backend.boundary-contracts":                           "mandatory",
		"clause.backend.http-independent-use-cases":                   "mandatory",
		"clause.backend.layered-architecture":                         "mandatory",
		"clause.backend.persistence-owner":                            "mandatory",
		"clause.backend.prohibit-generic-buckets":                     "prohibited",
		"clause.backend.thin-http-handlers":                           "mandatory",
		"clause.bun.add-from-owning-workspace":                        "mandatory",
		"clause.bun.prohibit-other-package-managers":                  "prohibited",
		"clause.bun.use-bun-owned-commands":                           "mandatory",
		"clause.bun.verify-dependency-before-add":                     "mandatory",
		"clause.context.adr-01-template":                              "mandatory",
		"clause.context.adr-02-active-status":                         "mandatory",
		"clause.context.adr-03-legacy-compatibility":                  "mandatory",
		"clause.context.backlog-01-operational-contract":              "mandatory",
		"clause.context.backlog-02-finding-boundary":                  "mandatory",
		"clause.context.docs-one-job-per-directory":                   "mandatory",
		"clause.context.docs-upstream-flow":                           "mandatory",
		"clause.context.findings-01-frontmatter":                      "mandatory",
		"clause.context.findings-02-pending":                          "mandatory",
		"clause.context.findings-03-partial":                          "mandatory",
		"clause.context.findings-04-deferred":                         "mandatory",
		"clause.context.findings-05-done":                             "mandatory",
		"clause.context.findings-06-append-evidence":                  "mandatory",
		"clause.context.findings-07-update-timestamp":                 "mandatory",
		"clause.context.findings-08-rollup":                           "mandatory",
		"clause.context.findings-09-archive":                          "mandatory",
		"clause.context.findings-10-live-work-health":                 "mandatory",
		"clause.context.findings-11-rollup-closure":                   "mandatory",
		"clause.context.inbox-01-triage":                              "mandatory",
		"clause.context.inbox-02-fleet-flow":                          "mandatory",
		"clause.context.inbox-03-extend-before-minting":               "mandatory",
		"clause.context.read-domain-contract":                         "mandatory",
		"clause.core.activate-matching-skills":                        "mandatory",
		"clause.core.ask-before-delivery":                             "stop-and-ask",
		"clause.core.ask-before-destructive-git":                      "stop-and-ask",
		"clause.core.ask-before-verification-configuration-change":    "stop-and-ask",
		"clause.core.ask-user-answerable-decisions":                   "stop-and-ask",
		"clause.core.assertion-reads-the-constant":                    "mandatory",
		"clause.core.conventional-commit-titles":                      "mandatory",
		"clause.core.fix-root-causes":                                 "mandatory",
		"clause.core.follow-dependency-workflow":                      "mandatory",
		"clause.core.keep-follow-ups-outside-slice":                   "mandatory",
		"clause.core.keep-root-compact":                               "mandatory",
		"clause.core.never-let-a-pipe-hide-a-gate":                    "mandatory",
		"clause.core.plan-the-release-first":                          "mandatory",
		"clause.core.preserve-git-scope":                              "mandatory",
		"clause.core.prohibit-external-research-for-local-code":       "prohibited",
		"clause.core.prohibit-secret-exposure":                        "prohibited",
		"clause.core.prohibit-verification-contract-bypass":           "prohibited",
		"clause.core.prohibit-verification-workarounds":               "prohibited",
		"clause.core.record-acceptance-evidence":                      "mandatory",
		"clause.core.request-pull-request-review":                     "mandatory",
		"clause.core.request-review-explicitly":                       "mandatory",
		"clause.core.require-fresh-evidence":                          "mandatory",
		"clause.core.require-tooling-authorization":                   "prohibited",
		"clause.core.research-authoritative-external-sources":         "mandatory",
		"clause.core.research-local-code-locally":                     "mandatory",
		"clause.core.review-new-dependencies":                         "mandatory",
		"clause.core.run-selected-verification":                       "mandatory",
		"clause.core.tooling-commit-choreography":                     "mandatory",
		"clause.core.use-conventional-commits":                        "mandatory",
		"clause.core.use-declared-external-research-fallback":         "mandatory",
		"clause.core.use-github-pr-workflow":                          "mandatory",
		"clause.core.verification-two-tiers":                          "mandatory",
		"clause.core.write-generated-guidance-in-english":             "mandatory",
		"clause.domain.canonical-language":                            "mandatory",
		"clause.domain.glossary-currency":                             "mandatory",
		"clause.domain.layout-decision":                               "mandatory",
		"clause.external-triage.classify-before-labelling":            "mandatory",
		"clause.external-triage.move-through-mapped-states":           "mandatory",
		"clause.external-triage.ask-for-an-unmapped-label":            "stop-and-ask",
		"clause.external-triage.needs-info-asks-and-stops":            "mandatory",
		"clause.external-triage.disclose-ai-authorship":               "mandatory",
		"clause.external-triage.record-wontfix-decisions":             "mandatory",
		"clause.external-triage.pull-request-is-an-issue-with-code":   "mandatory",
		"clause.external-triage.route-accepted-work-to-specs":         "mandatory",
		"clause.frontend.inspect-runnable-ui":                         "mandatory",
		"clause.frontend.follow-recorded-layout":                      "mandatory",
		"clause.frontend.organize-by-system":                          "mandatory",
		"clause.frontend.prohibit-incidental-ui-assertions":           "prohibited",
		"clause.frontend.public-system-boundary":                      "mandatory",
		"clause.frontend.read-design-contract-before-ui-work":         "mandatory",
		"clause.frontend.test-user-visible-behavior":                  "mandatory",
		"clause.secondbrain.01-consult-triggers":                      "mandatory",
		"clause.secondbrain.02-query-order":                           "mandatory",
		"clause.secondbrain.03-decision-consultation":                 "mandatory",
		"clause.secondbrain.04-external-research":                     "mandatory",
		"clause.secondbrain.05-research-limitations":                  "mandatory",
		"clause.secondbrain.baseline-owned-guidance":                  "mandatory",
		"clause.secondbrain.capture-self-contained":                   "mandatory",
		"clause.secondbrain.capture-trigger":                          "mandatory",
		"clause.secondbrain.cite-used-files":                          "mandatory",
		"clause.secondbrain.escalate-durable-updates":                 "mandatory",
		"clause.secondbrain.inbox-auto-capture":                       "mandatory",
		"clause.secondbrain.inbox-empty-state":                        "mandatory",
		"clause.secondbrain.inbox-entry-contract":                     "mandatory",
		"clause.secondbrain.inbox-triage-order":                       "mandatory",
		"clause.secondbrain.inbox-write-permission":                   "mandatory",
		"clause.secondbrain.prohibit-external-local-discovery":        "prohibited",
		"clause.secondbrain.prohibit-secret-access":                   "prohibited",
		"clause.secondbrain.prohibit-writes":                          "prohibited",
		"clause.secondbrain.record-decision-sources":                  "mandatory",
		"clause.secondbrain.research-capture":                         "mandatory",
		"clause.spec.keep-artifacts-in-spec-folder":                   "mandatory",
		"clause.spec.local-task-tracker-only":                         "mandatory",
		"clause.spec.project-constraints-01-active-artifacts":         "mandatory",
		"clause.spec.project-constraints-02-tooling-authorization":    "mandatory",
		"clause.spec.project-constraints-03-bounded-execution":        "mandatory",
		"clause.spec.project-constraints-04-qa-audit":                 "mandatory",
		"clause.spec.project-constraints-05-legacy-and-ownership":     "mandatory",
		"clause.spec.project-constraints-06-outside-evidence":         "mandatory",
		"clause.spec.name-the-finding-a-task-answers":                 "mandatory",
		"clause.spec.verification-fails-before-the-change":            "mandatory",
		"clause.spec.prohibit-tests-that-read-specs":                  "prohibited",
		"clause.typescript.type-fixtures-from-the-schema":             "mandatory",
		"clause.spec.routing-01-large-initiative":                     "mandatory",
		"clause.spec.routing-02-standard-feature":                     "mandatory",
		"clause.spec.routing-03-refactor-bugfix":                      "mandatory",
		"clause.spec.routing-04-trivial-change":                       "mandatory",
		"clause.spec.routing-05-task-graph":                           "mandatory",
		"clause.spec.sources-01-group-by-shared-context":              "mandatory",
		"clause.spec.sources-02-bound-the-group":                      "mandatory",
		"clause.spec.specs-are-downstream-artifacts":                  "mandatory",
		"clause.spec.status-only-in-task":                             "mandatory",
		"clause.spec.tracker-artifacts":                               "mandatory",
		"clause.spec.verification-two-tiers":                          "mandatory",
		"clause.typescript.inspect-dependent-interfaces-before-tests": "mandatory",
		"clause.typescript.keep-type-errors-visible":                  "mandatory",
		"clause.typescript.prohibit-incidental-test-oracles":          "prohibited",
		"clause.typescript.read-current-authoritative-docs":           "mandatory",
		"clause.typescript.test-observable-failure-modes":             "mandatory",
		"rule.monorepo.context-boundaries":                            "mandatory",
	}
}
