package baseline

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestTheComposedProfileTakesTheComposedSetup(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	composed, ok := catalog.profiles["go-cli-typescript-monorepo"]
	if !ok || composed["setup"] != "go-cli-typescript-bun" {
		t.Fatalf("composed profile = %v", composed)
	}
	standard := catalog.profiles["standard-typescript-monorepo"]
	for _, field := range []string{"entryDecisions", "workspaces", "optionalModules", "architecture", "httpContract", "capabilitySets", "activationBundles"} {
		if !reflect.DeepEqual(composed[field], standard[field]) {
			t.Errorf("composed %s differs from Standard TypeScript Monorepo", field)
		}
	}
}

func TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause(t *testing.T) {
	t.Parallel()
	_, plan := composedProfilePlan(t)
	for _, path := range []string{"go", "cli", "typescript-bun", "monorepo", "backend", "frontend"} {
		guide := planPostimage(t, plan, "docs/agents/"+path+".md")
		if !strings.Contains(string(guide.Content), "- ") {
			t.Errorf("guide %s has no clauses", path)
		}
	}
	if repeated := repeatedRenderedClause(plan.Postimages); repeated != "" {
		t.Fatalf("repeated clause: %s", repeated)
	}
}

func TestTheGoGuideNamesWhatItGoverns(t *testing.T) {
	t.Parallel()
	_, plan := composedProfilePlan(t)
	guide := string(planPostimage(t, plan, "docs/agents/go.md").Content)
	heading := strings.Index(guide, "# Go")
	scope := "# Go\n\nThese rules govern the repository's Go module: its commands, packages and\ntests. Code in another language follows its own guide."
	if heading < 0 || !strings.HasPrefix(guide[heading:], scope) {
		t.Fatalf("Go guide scope missing: %s", guide)
	}
}

func TestARepeatedRenderedClauseIsReported(t *testing.T) {
	t.Parallel()
	for _, clause := range []string{"- **mandatory**: Keep one owner.", "- Keep Go commands thin."} {
		t.Run(clause, func(t *testing.T) {
			guides := []Postimage{
				{Path: "docs/agents/first.md", Content: []byte("# First\n\n" + clause + "\n")},
				{Path: "docs/agents/second.md", Content: []byte("# Second\n\n" + clause + "\n")},
			}
			if got := repeatedRenderedClause(guides); got != clause {
				t.Fatalf("repeat = %q, want %q", got, clause)
			}
		})
	}
}

func repeatedRenderedClause(postimages []Postimage) string {
	seen := make(map[string]bool)
	for _, postimage := range postimages {
		if !strings.HasSuffix(postimage.Path, ".md") {
			continue
		}
		for _, line := range strings.Split(string(postimage.Content), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "- ") {
				continue
			}
			if seen[line] {
				return line
			}
			seen[line] = true
		}
	}
	return ""
}

func TestTheComposedProfilePlanConverges(t *testing.T) {
	t.Parallel()
	repository, plan := composedProfilePlan(t)
	if _, err := ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatal(err)
	}
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repository, ProfileID: "go-cli-typescript-monorepo",
		Decisions:    standardTypeScriptDecisions("make verify"),
		Preservation: RootPreservationRequest{Mode: PreservationModeManagedRefresh},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil {
		t.Fatalf("replan = %+v", outcome.Result)
	}
	if len(outcome.Plan.FileChanges) != 0 {
		t.Fatalf("replan changes = %+v", outcome.Plan.FileChanges)
	}
}

func composedProfilePlan(t *testing.T) (string, PlanDocument) {
	t.Helper()
	repository := newProjectDecisionPlanRepository(t)
	writeProfileAlignmentFile(t, repository, "go.mod", "module example.invalid/composed\n\ngo 1.26\n")
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repository, ProfileID: "go-cli-typescript-monorepo",
		Decisions:    standardTypeScriptDecisions("make verify"),
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil {
		t.Fatalf("plan = %+v", outcome.Result)
	}
	return repository, *outcome.Plan
}
