package baseline

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthorizationClausesStateTodaysRefusals(t *testing.T) {
	t.Parallel()
	clauses := embeddedBaselineClauses(t)
	for _, tt := range []struct{ id, guide, sentence string }{
		{"clause.core.verification-two-tiers", "agent-instructions.md", "A Run satisfies this by construction, because its Run Worktree is created from a commit. No other command checks it: `roundfix spec check --run-verification` executes the commands of the working-tree Spec in a checkout of `HEAD`, and `roundfix settle` reads the Task file of the directory it settles, so whoever runs them on an uncommitted or modified Spec first reads the commands they will execute. No approval record makes an untrusted source executable; commit the artifacts or do not execute them."},
		{"clause.core.tooling-commit-choreography", "agent-instructions.md", "a proposed, absent, contradictory, or withdrawn record grants no governed mutation and no delivery operation. `roundfix implement` still runs a Spec that changes no Governed Path, and `roundfix deliver start` refuses a Spec whose record does not grant `implement`, `commit`, `push`, `pull_request`, and `merge`. An executor cannot widen the grant it depends on."},
		{"clause.spec.project-constraints-02-tooling-authorization", "spec-routing.md", "A Spec that changes no Governed Path records `paths: []`, which grants its listed operations and bounds no path; a record without a `paths` key is refused."},
		{"clause.spec.project-constraints-03-bounded-execution", "spec-routing.md", "An authorized tooling Task may change a Governed Path, that is a protected tooling path, only when its grant bounds that path; stop before any other governed mutation, because the Task fails when its commit changes a Governed Path outside the grant. An ordinary path a Task changes without declaring it is not refused: the Daemon records it in the Task file under `## Recorded paths`, and the QA gate discloses it."},
		{"clause.spec.project-constraints-06-outside-evidence", "spec-routing.md", "When the outside source cannot be obtained during authoring, record the row as blocked with that reason and continue: decomposition never stalls and never asks a person. The QA gate then holds Pull Request preparation until the row is satisfied or carried forward on declared unmoved evidence."},
	} {
		t.Run(tt.id, func(t *testing.T) {
			found := false
			for _, clause := range clauses {
				if clause.ID != tt.id {
					continue
				}
				found = true
				if clause.Enforcement != "mandatory" {
					t.Errorf("enforcement = %q, want mandatory", clause.Enforcement)
				}
				if !strings.Contains(clause.Guidance, tt.sentence) {
					t.Errorf("embedded clause lacks %q", tt.sentence)
				}
			}
			if !found {
				t.Fatalf("missing embedded clause %s", tt.id)
			}
			for _, path := range []string{
				"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + tt.guide,
				"../../docs/agents/" + tt.guide,
			} {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(strings.Join(strings.Fields(string(content)), " "), tt.sentence) {
					t.Errorf("%s lacks %q", path, tt.sentence)
				}
			}
		})
	}
}

func TestAuthorizationClausesDropTheRemovedPhrases(t *testing.T) {
	t.Parallel()
	var modules, goldens int
	for _, root := range []string{"assets/modules", "assets/formatter-fixtures", "../../docs/agents"} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			if root == "assets/formatter-fixtures" && !strings.Contains(filepath.ToSlash(path), "/golden/") {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if root == "assets/modules" {
				modules++
			}
			if root == "assets/formatter-fixtures" {
				goldens++
			}
			text := strings.Join(strings.Fields(string(content)), " ")
			for _, phrase := range []string{"execution_approvals", "does not by itself refuse implementation", "may mutate only its bounded", "never blocks the Spec"} {
				if strings.Contains(text, phrase) {
					t.Errorf("%s contains removed phrase %q", path, phrase)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if modules == 0 || goldens == 0 {
		t.Fatalf("empty scan: %d modules, %d goldens", modules, goldens)
	}
}
