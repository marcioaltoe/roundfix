// Suite: repository verification planning and timing.
// Invariant: selection and reporting preserve the documented boundaries.
// Boundary IN: in-process functions and isolated fixtures.
// Boundary OUT: live delivery, GitHub and the Run Database.
package verifyselect

import (
	"bytes"
	"path/filepath"
	"reflect"
	"roundfix/internal/store"
	"strings"
	"testing"
)

func fixtureModule() []goListPackage {
	var result []goListPackage
	for _, dir := range []string{"internal/app", "internal/baseline", "internal/baselineacp", "internal/cli", "skills", "cmd/roundfix"} {
		result = append(result, goListPackage{Dir: filepath.Join("/repo", dir), ImportPath: dir})
	}
	result[5].Imports = []string{"internal/app", "internal/baseline", "internal/baselineacp", "internal/cli", "skills"}
	return result
}
func checkPaths(t *testing.T, paths, want []string) {
	t.Helper()
	p, e := PlanPackages(Facts{Changed: paths}, fixtureModule(), "/repo")
	if e != nil || !reflect.DeepEqual(p.Packages, want) {
		t.Fatalf("%v: %+v %v, want %v", paths, p, e, want)
	}
}
func TestPlanMapsEachChangedPathToItsPackage(t *testing.T) {
	t.Parallel()
	for _, p := range []string{"internal/baseline/catalog.go", "internal/baselineacp/session.go", "skills/roundfix/SKILL.md", "internal/cli/baseline_plan.go", "internal/cli/baseline_plan_test.go", "internal/app/version.go", "internal/app/version_test.go", "internal/app/README.md"} {
		owner := strings.Split(p, "/")[0] + "/" + strings.Split(p, "/")[1]
		if strings.HasPrefix(p, "skills/") {
			owner = "skills"
		}
		checkPaths(t, []string{p}, []string{"./cmd/roundfix", "./" + owner})
	}
	checkPaths(t, []string{".agents/skills/roundfix/SKILL.md"}, []string{"./cmd/roundfix", "./internal/baseline", "./internal/baselineacp", "./skills"})
	for _, p := range []string{"go.mod", "go.sum", "Makefile", "README.md", "testdata/input.json"} {
		checkAll(t, p)
	}
	checkPaths(t, []string{"docs/user-guide/run-database-lifecycle.md"}, nil)
}
func checkAll(t *testing.T, p string) {
	t.Helper()
	plan, e := PlanPackages(Facts{Changed: []string{p}}, fixtureModule(), "/repo")
	if e != nil || len(plan.Packages) != len(fixtureModule()) {
		t.Fatalf("%s: %+v %v", p, plan, e)
	}
}
func TestTheLiveDocumentationTreesSelectNoPackage(t *testing.T) {
	t.Parallel()
	for _, p := range []string{"docs/history/specs/x.md", "docs/adr/x.md", "docs/specs/x.md", "docs/user-guide/run-database-lifecycle.md"} {
		checkPaths(t, []string{p}, nil)
	}
}
func TestRootMarkdownSelectsEveryPackage(t *testing.T) { t.Parallel(); checkAll(t, "README.md") }
func TestUnknownNonDocumentationPathsFailSafe(t *testing.T) {
	t.Parallel()
	checkAll(t, ".github/workflows/unknown.yml")
}
func TestAFixtureChangeSelectsItsPackageAndItsImporters(t *testing.T) {
	t.Parallel()
	for _, p := range []string{"internal/app/testdata/version.golden", "internal/baseline/testdata/profile.golden"} {
		owner := strings.Split(p, "/")[1]
		checkPaths(t, []string{p}, []string{"./cmd/roundfix", "./internal/" + owner})
	}
}
func TestPlanSelectsTheReachOfAChangeSet(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ paths, want []string }{
		{[]string{"internal/app/version.go"}, []string{"./cmd/roundfix", "./internal/app"}},
		{[]string{"skills/roundfix/SKILL.md"}, []string{"./cmd/roundfix", "./skills"}},
		{[]string{"internal/app/version.go", "internal/baseline/catalog.go"}, []string{"./cmd/roundfix", "./internal/app", "./internal/baseline"}},
		{[]string{"docs/user-guide/x.md"}, nil},
	} {
		checkPaths(t, c.paths, c.want)
	}
	checkAll(t, "go.mod")
}
func TestASkillChangeSelectsTheBaselinePackagesAndTheirImporters(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"import", "test", "external", "none"} {
		module := fixtureModule()
		module[5].Imports = nil
		switch kind {
		case "import":
			module[5].Imports = []string{"internal/baseline"}
		case "test":
			module[5].TestImports = []string{"internal/baseline"}
		case "external":
			module[5].XTestImports = []string{"internal/baseline"}
		}
		p, e := PlanPackages(Facts{Changed: []string{".agents/skills/x.md"}}, module, "/repo")
		want := 3
		if kind != "none" {
			want++
		}
		if e != nil || len(p.Packages) != want {
			t.Fatalf("%s: %+v %v", kind, p, e)
		}
	}
}
func TestPlanPackagesSettlesOnlyOnTheRunsBranch(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		marker, branch string
		want           Mode
	}{{"a", "roundfix/run-a", Settlement}, {"a", "roundfix/run-a-task_03", Settlement}, {"a", "", Reach}, {"a", "feat/other", Reach}, {"a", "roundfix/batch-a", Reach}, {"a", "roundfix/run-ab", Reach}, {"", "roundfix/run-a", Reach}} {
		p, e := PlanPackages(Facts{RunMarker: c.marker, Branch: c.branch}, nil, "/repo")
		if e != nil || p.Mode != c.want {
			t.Fatalf("%+v: %+v %v", c, p, e)
		}
	}
}
func TestPlanPackagesRunsTheCompleteSuiteForAnArchiveRecord(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"docs/history/specs/x.md", "docs/history/specs/nested/x.md", "docs/history/specs/x.json"} {
		p, e := PlanPackages(Facts{Added: []string{name}}, nil, "/repo")
		if e != nil {
			t.Fatal(e)
		}
		if name == "docs/history/specs/x.md" {
			if p.Mode != Complete || !p.Fresh || !reflect.DeepEqual(p.Packages, []string{"./..."}) {
				t.Fatal(p)
			}
		} else if p.Mode != Reach {
			t.Fatal(p)
		}
	}
}
func TestTheSettlementBranchPrefixIsTheStores(t *testing.T) {
	t.Parallel()
	if runBranchPrefix != store.RunBranchPrefix {
		t.Fatal("Run prefix drift")
	}
}
func TestTheVersionFilesOfAReleaseSelectEveryPackage(t *testing.T) {
	t.Parallel()
	for _, p := range []string{"CHANGELOG.md", "dist/npm/roundfix/package.json"} {
		checkAll(t, p)
	}
}
func TestRunPrintsTheContractsOfOneBuildTag(t *testing.T) {
	t.Parallel()
	for _, tag := range []string{"docscontract", "repocontract"} {
		var out, err bytes.Buffer
		code := Run(t.Context(), []string{"-repo", "../..", "-contracts", "-all", "-tag", tag}, &out, &err)
		if code != 0 || out.Len() == 0 {
			t.Fatalf("%d %s", code, err.String())
		}
		for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
			if !strings.HasPrefix(line, tag+" ") {
				t.Fatal(line)
			}
		}
	}
}
