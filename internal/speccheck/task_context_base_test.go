// Suite: Task Context paths resolved through the checked tree and base revisions.
// Invariant: only HEAD or the Delivery Base can resolve a missing Task Context path.
// Boundary IN: temporary local Git repositories and one pending Task Context entry.
// Boundary OUT: public Spec Consistency Check findings, including their exact diagnostics.
package speccheck_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/speccheck"
)

func TestTaskContextPathTrackedAtTheBaseResolves(t *testing.T) {
	t.Parallel()
	const path = "src/gone.txt"
	for _, tc := range []struct {
		name string
		want int
	}{
		{"present", 0},
		{"working tree removal", 0},
		{"staged removal", 0},
		{"committed removal", 0},
		{"instruction committed removal", 0},
		{"branch file removed in the working tree", 0},
		{"never tracked", 1},
		{"created and removed on the branch", 1},
		{"removed on main before the branch", 1},
		{"no Git repository", 1},
		{"nested below another work tree", 1},
		{"no default branch committed removal", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture := newUndeclaredFixture(t)
			if tc.name == "nested below another work tree" {
				outer := t.TempDir()
				gittest.InitRepo(t, outer, "-b", "main")
				writeCitationFixtureFile(t, outer, path, "outer source\n")
				gittest.Run(t, outer, "add", ".")
				gittest.Run(t, outer, "commit", "-m", "test: seed outer source")
				gittest.Run(t, outer, "update-ref", "refs/remotes/origin/main", "HEAD")
				checked := filepath.Join(outer, "checked")
				if err := os.Rename(fixture.repoRoot, checked); err != nil {
					t.Fatal(err)
				}
				fixture.repoRoot = checked
				fixture.specsRoot = filepath.Join(checked, "docs", "specs")
			}
			kind := "interface"
			if tc.name == "instruction committed removal" {
				kind = "instruction"
			}
			fixture.writeTask("pending", "backend", "- "+kind+": `"+path+"`", "true")
			writeCitationFixtureFile(t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/_prd.md", "---\nstatus: active\n---\n\n# Fixture\n\n## Success Metrics\n\nNone. Local path resolution only.\n")
			if tc.name != "no Git repository" && tc.name != "nested below another work tree" {
				branch := "main"
				if tc.name == "no default branch committed removal" {
					branch = "work"
				}
				gittest.InitRepo(t, fixture.repoRoot, "-b", branch)
				switch tc.name {
				case "never tracked", "created and removed on the branch", "branch file removed in the working tree":
				default:
					writeCitationFixtureFile(t, fixture.repoRoot, path, "source\n")
				}
				gittest.Run(t, fixture.repoRoot, "add", ".")
				gittest.Run(t, fixture.repoRoot, "commit", "-m", "test: seed context")
				if tc.name == "removed on main before the branch" {
					gittest.Run(t, fixture.repoRoot, "rm", path)
					gittest.Run(t, fixture.repoRoot, "commit", "-m", "test: remove on main")
				}
				if branch == "main" {
					gittest.Run(t, fixture.repoRoot, "update-ref", "refs/remotes/origin/main", "HEAD")
					gittest.Run(t, fixture.repoRoot, "checkout", "-b", "work")
				}
				switch tc.name {
				case "created and removed on the branch", "branch file removed in the working tree":
					writeCitationFixtureFile(t, fixture.repoRoot, path, "branch source\n")
					gittest.Run(t, fixture.repoRoot, "add", path)
					gittest.Run(t, fixture.repoRoot, "commit", "-m", "test: add branch source")
				}
				switch tc.name {
				case "working tree removal", "branch file removed in the working tree":
					if err := os.Remove(filepath.Join(fixture.repoRoot, path)); err != nil {
						t.Fatal(err)
					}
				case "staged removal":
					gittest.Run(t, fixture.repoRoot, "rm", path)
				case "committed removal", "instruction committed removal", "created and removed on the branch", "no default branch committed removal":
					gittest.Run(t, fixture.repoRoot, "rm", path)
					gittest.Run(t, fixture.repoRoot, "commit", "-m", "test: remove source")
				}
			}
			result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
			if err != nil {
				t.Fatal(err)
			}
			findings := findingsWithCode(result, speccheck.CodeReferenceUnresolved)
			if len(findings) != tc.want {
				t.Fatalf("findings = %+v, want %d", findings, tc.want)
			}
			task := "docs/specs/" + fixture.slug + "/task_01.md"
			for _, finding := range findings {
				want := speccheck.Finding{
					Code:     speccheck.CodeReferenceUnresolved,
					Severity: speccheck.SeverityError,
					Summary:  task + " declares unresolved Task Context path " + path,
					Where:    []speccheck.Location{{Path: task, Line: 12}, {Path: path, Line: 1}},
					Fix:      "Remove the Task Context entry for " + path + " from " + task + ", or declare it under deletes: when the Task removes it; replace it only with a path the Task itself reads or changes.",
				}
				if !reflect.DeepEqual(finding, want) {
					t.Errorf("finding = %+v, want %+v", finding, want)
				}
			}
		})
	}
}
