package baseline

import (
	"context"
	"os"
	"strings"
	"testing"
)

const askUserClauseID = "clause.core.ask-user-answerable-decisions"
const askUserClauseLimit = "Ask at most three questions per call or message, never more than the tool accepts, and keep only that call pending."
const askUserClauseIndependence = "Put questions in one call only when they are independent of each other; a question whose options depend on another answer waits for that answer."
const askUserClauseRetired = "Ask exactly one question per call or message"

func TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall(t *testing.T) {
	t.Parallel()
	for _, clause := range embeddedBaselineClauses(t) {
		if clause.ID != askUserClauseID {
			continue
		}
		if clause.Enforcement != "stop-and-ask" {
			t.Errorf("clause enforcement = %q, want stop-and-ask", clause.Enforcement)
		}
		if len(clause.Replaces) != 0 {
			t.Errorf("clause replaces = %v, want none", clause.Replaces)
		}
		assertQuestionBatchGuidance(t, clause.Guidance)
		return
	}
	t.Fatalf("embedded catalog lacks %s", askUserClauseID)
}

func TestTheAskUserClauseRendersInTheAgentInstructionsGuide(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path string }{
		{"Standard TypeScript golden", "assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md"},
		{"repository", "../../docs/agents/agent-instructions.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			assertQuestionBatchGuidance(t, string(data))
		})
	}
}

func TestAnAdopterRetainsTheAskUserClause(t *testing.T) {
	t.Parallel()
	request, catalog := newClauseReplacementAdopter(t)
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("adopter refresh = %+v, want ready plan", outcome.Result)
	}
	if delta := outcome.Plan.ClauseDelta; delta == nil || delta.Dispositions[askUserClauseID] != ClauseRetained {
		t.Fatalf("clause delta = %+v, want %s retained", delta, askUserClauseID)
	}
	for _, evidence := range outcome.Plan.Retention {
		if evidence.FromClause != askUserClauseID {
			continue
		}
		if len(evidence.Targets) != 1 || evidence.Targets[0] != askUserClauseID {
			t.Fatalf("retention evidence = %+v, want self target", evidence)
		}
		return
	}
	t.Fatalf("no retention evidence for %s", askUserClauseID)
}

func assertQuestionBatchGuidance(t *testing.T, text string) {
	t.Helper()
	normalized := strings.Join(strings.Fields(text), " ")
	for _, sentence := range []string{askUserClauseLimit, askUserClauseIndependence} {
		if count := strings.Count(normalized, sentence); count != 1 {
			t.Errorf("sentence %q occurs %d times, want 1", sentence, count)
		}
	}
	if strings.Contains(normalized, askUserClauseRetired) {
		t.Errorf("guidance contains retired sentence %q", askUserClauseRetired)
	}
}
