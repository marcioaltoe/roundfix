package delivery

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestTheLoopClauseOrderMatchesTheDeliveryQueue(t *testing.T) {
	t.Parallel()
	if err := compareLoopClauseOrder(readLoopClause(t), reviewedDeliveryActions(t)); err != nil {
		t.Fatal(err)
	}
}

func TestALoopClauseThatArchivesBeforeReviewIsRefused(t *testing.T) {
	t.Parallel()
	const oldOrder = "Follow one order per Spec: implement the graph including its authored gate, archive and commit the candidate, apply the configured pre-PR review policy, open the Pull Request, verify required checks for the current head, and merge."
	if err := compareLoopClauseOrder(oldOrder, reviewedDeliveryActions(t)); err == nil || !strings.Contains(err.Error(), "order mismatch") {
		t.Fatalf("old order comparison = %v, want order mismatch", err)
	}
}

func TestArchiveBeforeReviewWithAllDeliveryActionsIsRefused(t *testing.T) {
	t.Parallel()
	guidance := strings.Replace(readLoopClause(t),
		"apply the configured pre-PR review policy, archive and commit the candidate",
		"archive and commit the candidate, apply the configured pre-PR review policy", 1)
	if err := compareLoopClauseOrder(guidance, reviewedDeliveryActions(t)); err == nil || !strings.Contains(err.Error(), "order mismatch") {
		t.Fatalf("reordered clause comparison = %v, want order mismatch", err)
	}
}

func TestAnUnmappedDeliveryActionIsRefused(t *testing.T) {
	t.Parallel()
	events := append(reviewedDeliveryActions(t), "new-stage")
	if err := compareLoopClauseOrder(readLoopClause(t), events); err == nil || !strings.Contains(err.Error(), "unmapped Delivery Queue action") {
		t.Fatalf("unknown action comparison = %v, want unmapped action", err)
	}
}

func readLoopClause(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile("../baseline/assets/modules/autonomous-work.json")
	if err != nil {
		t.Fatal(err)
	}
	var module struct {
		Rules []struct {
			Clauses []struct {
				ID       string `json:"id"`
				Guidance string `json:"guidance"`
			} `json:"clauses"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(content, &module); err != nil {
		t.Fatal(err)
	}
	for _, rule := range module.Rules {
		for _, clause := range rule.Clauses {
			if clause.ID == "clause.autonomous.loop-01-qa-once" {
				return clause.Guidance
			}
		}
	}
	t.Fatal("module lacks the loop clause")
	return ""
}

func reviewedDeliveryActions(t *testing.T) []string {
	t.Helper()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot, slug = "/repo", "loop-order"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug}); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	if _, err := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary()).Run(ctx, gitRoot); err != nil {
		t.Fatal(err)
	}
	if item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; item.Stage != store.DeliveryStageMerged {
		t.Fatalf("reviewed item stage = %s, want merged", item.Stage)
	}
	return workflow.events[slug]
}

func compareLoopClauseOrder(guidance string, events []string) error {
	const marker = "Follow one order per Spec:"
	if !strings.HasPrefix(guidance, marker) {
		return fmt.Errorf("loop clause lacks order marker")
	}
	sentence, _, terminated := strings.Cut(strings.TrimPrefix(guidance, marker), ".")
	if !terminated {
		return fmt.Errorf("loop order sentence lacks a full stop")
	}
	phrases := map[string]string{
		"implement the graph including its authored gate": "run",
		"apply the configured pre-PR review policy":       "review",
		"archive and commit the candidate":                "archive",
		"run the repository gate":                         "gate",
		"push the candidate":                              "push",
		"open the Pull Request":                           "pull-request",
		"verify required checks for the current head":     "checks",
		"merge": "merge",
	}
	var declared, observed []string
	for _, phrase := range strings.Split(sentence, ",") {
		phrase = strings.TrimPrefix(strings.TrimSpace(phrase), "and ")
		action, ok := phrases[phrase]
		if !ok {
			return fmt.Errorf("unmapped loop phrase %q", phrase)
		}
		declared = append(declared, action)
	}
	for _, event := range events {
		switch event {
		case "create-branch", "policy", "authorization", "publication":
			// Preparation and policy reads are bookkeeping within the declared steps.
		case "run", "review", "archive", "gate", "push", "pull-request", "checks", "merge":
			observed = append(observed, event)
		default:
			return fmt.Errorf("unmapped Delivery Queue action %q", event)
		}
	}
	if !reflect.DeepEqual(declared, observed) {
		return fmt.Errorf("loop order mismatch: declared %v, Delivery Queue ran %v", declared, observed)
	}
	return nil
}
