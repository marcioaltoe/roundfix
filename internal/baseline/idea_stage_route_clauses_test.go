package baseline

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

var ideaStageRouteClausePhrases = []struct {
	name   string
	clause string
	phrase string
}{
	{
		name:   "large initiative route",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`.",
	},
	{
		name:   "feature takes indicated route",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "A feature takes the same route when its Idea stage is indicated.",
	},
	{
		name:   "judge three criteria",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Before writing the PRD of a feature that has no `_idea.md` and that the user did not send to `write-idea`, judge three criteria about the request: whether the right shape of the solution is still open, whether the request rests on an assumption about its users, the business or the market that it gives no evidence for, and whether deciding to build it depends on a market or competitive question.",
	},
	{
		name:   "judge recorded item or session request",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "When the request is recorded as an item that is not yet a Spec, which is an open Backlog Entry, a note in `docs/_inbox/` or a message handed in as a file, run `roundfix spec judge --stage idea --item <reference>`; when it exists only in the session, write it to a file outside the repository and run `roundfix spec judge --stage idea --request <path>`, which sends the judge only the project's name, the request's first line and its first 4,000 characters, after the secret scan.",
	},
	{
		name:   "agent answers when judge does not answer",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Take the judge's answers; when the command reports that the judge did not answer, answer the same three criteria yourself, each as yes, no or unsure, and say that the judge did not answer.",
	},
	{
		name:   "verdict rules",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "The Idea stage is indicated when any criterion is yes, not indicated when every criterion is no, and undecided otherwise.",
	},
	{
		name:   "indicated starts idea",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Indicated: start `write-idea` without asking, show the verdict, each criterion's answer with its probability and who judged, and let the user dismiss the stage at its first question.",
	},
	{
		name:   "undecided is user's choice",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Undecided, or in doubt: the user decides, seeing the same answers.",
	},
	{
		name:   "not indicated continues to prd",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Not indicated: continue to `write-prd`.",
	},
	{
		name:   "record judgment",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Record who judged, each answer, the verdict, whether the stage ran or was dismissed and why, and what the Idea did to the shape of the request, in the Idea when one is written and always in the PRD.",
	},
	{
		name:   "judgment never refuses",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "The judgment informs and never refuses: no command fails for it, and a dismissal needs only its recorded reason.",
	},
	{
		name:   "unattended indicated or undecided waits",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "When the stage is indicated or undecided and no user is present to answer, write no PRD for that feature: keep it as a Backlog Entry that records the judgment and waits, with `waiting_on` and `since` in its frontmatter, and never answer the stage's questions or the user's decision yourself.",
	},
	{
		name:   "fix or refactor never judged",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "A fix or a refactor is never judged for the Idea stage.",
	},
	{
		name:   "standard feature route",
		clause: "clause.spec.routing-02-standard-feature",
		phrase: "Standard feature that changes product behavior: run `write-prd` → `write-techspec` → `write-tasks` once its Idea stage is judged not indicated or the user dismissed it; skip the TechSpec only when the feature has no architectural surface.",
	},
	{
		name:   "idea judgment chooses route",
		clause: "clause.spec.routing-05-task-graph",
		phrase: "Use `brainstorming` before creative or feature work, start with the smaller sufficient route when two routes fit, except that the Idea-stage judgment, never the smaller route, decides whether a feature starts at `write-idea`, and execute implementation from the Task Graph.",
	},
	{
		name:   "surface transcript assertion",
		clause: "clause.spec.routing-05-task-graph",
		phrase: "The Task that implements a Surface Transcript asserts the whole transcript in one named test: every line of its standard output and standard error, in order, and its exit code, copied from the transcript with only its declared controlled variations; the Task's Verification runs that test by name.",
	},
	{
		name:   "invariant first sentence retained",
		clause: "clause.spec.routing-01-large-initiative",
		phrase: "Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`.",
	},
	{
		name:   "invariant techspec exception retained",
		clause: "clause.spec.routing-02-standard-feature",
		phrase: "skip the TechSpec only when the feature has no architectural surface.",
	},
	{
		name:   "invariant surface transcript retained",
		clause: "clause.spec.routing-05-task-graph",
		phrase: "The Task that implements a Surface Transcript asserts the whole transcript in one named test: every line of its standard output and standard error, in order, and its exit code, copied from the transcript with only its declared controlled variations; the Task's Verification runs that test by name.",
	},
}

func TestTheRouteClausesJudgeAFeatureForTheIdeaStage(t *testing.T) {
	t.Parallel()

	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("spec-workflow")
	if !ok {
		t.Fatal("missing module spec-workflow")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}

	clauses := make(map[string]document)
	for _, rule := range objectsOrEmpty(module["rules"]) {
		if rule["id"] != "rule.spec.routing" {
			continue
		}
		if rule["version"] != float64(12) {
			t.Fatalf("rule.spec.routing version = %v, want 12", rule["version"])
		}
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			clauses[clause["id"].(string)] = clause
		}
	}

	golden, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md")
	if !ok {
		t.Fatal("missing Standard TypeScript Spec routing golden")
	}
	local, err := os.ReadFile("../../docs/agents/spec-routing.md")
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range ideaStageRouteClausePhrases {
		want := want
		t.Run(want.name, func(t *testing.T) {
			clause, ok := clauses[want.clause]
			if !ok {
				t.Fatalf("missing clause %s", want.clause)
			}
			if clause["enforcement"] != "mandatory" {
				t.Fatalf("clause %s enforcement = %v, want mandatory", want.clause, clause["enforcement"])
			}
			if _, exists := clause["replaces"]; exists {
				t.Fatalf("clause %s has replaces", want.clause)
			}
			guidance, ok := clause["guidance"].(string)
			if !ok || strings.Count(guidance, want.phrase) != 1 {
				t.Fatalf("clause %s contains phrase %d times, want exactly once", want.clause, strings.Count(guidance, want.phrase))
			}

			literal := strings.Join(strings.Fields(want.phrase), " ")
			for name, data := range map[string][]byte{
				"Standard TypeScript golden": golden.Data,
				"repository guide":           local,
			} {
				t.Run(name, func(t *testing.T) {
					text := strings.Join(strings.Fields(string(data)), " ")
					if count := strings.Count(text, literal); count != 1 {
						t.Fatalf("guide contains literal %d times, want exactly once", count)
					}
				})
			}
		})
	}
}
