// Suite: Park Class policy.
// Boundary IN: delivery source constants and persisted blocker strings.
// Boundary OUT: status rendering is covered by the CLI suite.
package delivery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func unclassifiedBlockers(t *testing.T, files map[string]*ast.File) []string {
	t.Helper()
	var missing []string
	for _, file := range files {
		for _, decl := range file.Decls {
			constants, ok := decl.(*ast.GenDecl)
			if !ok || constants.Tok != token.CONST {
				continue
			}
			for _, spec := range constants.Specs {
				value := spec.(*ast.ValueSpec)
				for i, name := range value.Names {
					if !ast.IsExported(name.Name) || !strings.HasPrefix(name.Name, "Blocker") {
						continue
					}
					if i >= len(value.Values) {
						t.Fatalf("blocker %s needs an explicit value for the class sweep", name.Name)
					}
					literal, ok := value.Values[i].(*ast.BasicLit)
					if !ok || literal.Kind != token.STRING {
						t.Fatalf("blocker %s is not a string literal", name.Name)
					}
					blocker, err := strconv.Unquote(literal.Value)
					if err != nil {
						t.Fatal(err)
					}
					for _, suffix := range []string{"", ": details"} {
						if ClassifyPark(store.DeliveryQueue{}, store.DeliveryQueueItem{Blocker: blocker + suffix}).Class == ParkClassUnclassified {
							missing = append(missing, name.Name+suffix)
						}
					}
				}
			}
		}
	}
	return missing
}

func TestEveryBlockerHasAParkClass(t *testing.T) {
	packages, err := parser.ParseDir(token.NewFileSet(), ".", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	files := packages["delivery"].Files
	if missing := unclassifiedBlockers(t, files); len(missing) != 0 {
		t.Fatalf("unclassified exported blockers: %v", missing)
	}
	injected, err := parser.ParseFile(token.NewFileSet(), "new_blocker.go", `package delivery; const BlockerFuture = "future-blocker"`, 0)
	if err != nil {
		t.Fatal(err)
	}
	files["new_blocker.go"] = injected
	if missing := unclassifiedBlockers(t, files); len(missing) != 2 || missing[0] != "BlockerFuture" {
		t.Fatalf("sweep did not reject a new blocker: %v", missing)
	}
}

func TestParkClassesAndNextActions(t *testing.T) {
	for _, tt := range []struct {
		blocker string
		class   ParkClass
	}{
		{BlockerChecksTimeout, ParkClassEnvironment}, {BlockerItemWorktreeMissing, ParkClassEnvironment}, {BlockerDeliveryError, ParkClassEnvironment},
		{BlockerRunUnresolved, ParkClassFinding}, {BlockerReviewFindings, ParkClassFinding}, {BlockerCorrectiveSpecRequired, ParkClassFinding}, {BlockerGateFailed, ParkClassFinding}, {BlockerChecksFailed, ParkClassFinding}, {BlockerRevalidationFailed, ParkClassFinding},
		{BlockerRunBudgetExceeded, ParkClassBudget}, {BlockerQueueDeadline, ParkClassBudget}, {BlockerReviewBlocked, ParkClassReview}, {BlockerReviewStale, ParkClassReview}, {BlockerUnauthorized, ParkClassAuthorization}, {BlockerFlakyCheck, ParkClassFlakyCheck}, {"unknown", ParkClassUnclassified},
	} {
		t.Run(tt.blocker, func(t *testing.T) {
			for _, suffix := range []string{"", ": pkg"} {
				item := store.DeliveryQueueItem{SpecSlug: "example", Worktree: "/worktrees/example", Stage: store.DeliveryStageParked, Blocker: tt.blocker + suffix}
				queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{item}}
				park := ClassifyPark(queue, item)
				question, found := PendingQuestionFor(queue)
				if park.Class != tt.class || !found || question.Answer != park.Next {
					t.Fatalf("classification=%+v question=%+v", park, question)
				}
			}
		})
	}
	item := store.DeliveryQueueItem{SpecSlug: "example", Blocker: BlockerFlakyCheck + ": roundfix/internal/other"}
	if got := ClassifyPark(store.DeliveryQueue{}, item).Next; got != "the check failed twice in roundfix/internal/other, which example did not change; fix or re-run it, then run roundfix deliver retry example" {
		t.Fatalf("flaky next action = %q", got)
	}
}

func TestParkClassRequiresAnExactBlockerOrColonPrefix(t *testing.T) {
	for _, blocker := range []string{BlockerReviewStale + "-extra", BlockerRevalidationFailed + "extra", " " + BlockerChecksFailed} {
		if got := ClassifyPark(store.DeliveryQueue{}, store.DeliveryQueueItem{Blocker: blocker}).Class; got != ParkClassUnclassified {
			t.Fatalf("%q classified as %s", blocker, got)
		}
	}
}

func TestExistingGenericParkAnswersAreUnchanged(t *testing.T) {
	for _, blocker := range []string{
		BlockerRunUnresolved, BlockerReviewFindings, BlockerCorrectiveSpecRequired,
		BlockerGateFailed, BlockerChecksFailed, BlockerChecksTimeout,
		BlockerUnauthorized, BlockerDeliveryError, BlockerItemWorktreeMissing,
		BlockerRunBudgetExceeded, BlockerReviewBlocked, BlockerReviewStale,
	} {
		t.Run(blocker, func(t *testing.T) {
			item := store.DeliveryQueueItem{SpecSlug: "example", Blocker: blocker, Stage: store.DeliveryStageParked}
			question, found := PendingQuestionFor(store.DeliveryQueue{Items: []store.DeliveryQueueItem{item}})
			if !found || question.Answer != "resolve the blocker, then run roundfix deliver retry example" {
				t.Fatalf("question=%+v found=%t", question, found)
			}
		})
	}
}
