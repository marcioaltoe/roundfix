// Suite: repository verification planning and timing.
// Invariant: selection and reporting preserve the documented boundaries.
// Boundary IN: in-process functions and isolated fixtures.
// Boundary OUT: live delivery, GitHub and the Run Database.
package verifyselect_test

import (
	"bytes"
	"reflect"
	"regexp"
	"roundfix/internal/verifyselect"
	"strings"
	"testing"
)

func TestTheCLIShardPatternsAreDerivedFromTheSource(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeFile(t, root, "internal/cli/fixture_test.go", "package cli\nimport \"testing\"\nfunc TestOne(t *testing.T){}\nfunc TestTwo(t *testing.T){}\n")
	names, err := verifyselect.CLITestNames(root)
	if err != nil || !reflect.DeepEqual(names, []string{"TestOne", "TestTwo"}) {
		t.Fatal(names, err)
	}
	flag, pattern, err := verifyselect.CLIShard(root, 1, 2)
	if err != nil || flag != "-run" {
		t.Fatal(flag, pattern, err)
	}
	_, skip, err := verifyselect.CLIShard(root, 2, 2)
	if err != nil || skip != pattern {
		t.Fatal(skip, err)
	}
	if regexp.MustCompile(pattern).MatchString("TestNotDeclared") {
		t.Fatal("undeclared test selected by first shard")
	}
}
func TestTheJobsPartitionTheSelectedPackages(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeFile(t, root, "internal/cli/fixture_test.go", "package cli\nimport \"testing\"\nfunc TestOne(t *testing.T){}\n")
	plan := verifyselect.Plan{Packages: []string{"./core/direct", "./core/testonly", "./internal/baseline", "./internal/baselineacp", "./internal/cli", "./skills/check"}}
	args, err := verifyselect.JobArguments(root, plan, "packages")
	if err != nil || len(args) != 5 {
		t.Fatal(args, err)
	}
	for _, name := range []string{"cli-1", "cli-2"} {
		args, err := verifyselect.JobArguments(root, plan, name)
		if err != nil || len(args) != 3 || args[2] != "./internal/cli" {
			t.Fatal(args, err)
		}
	}
}

// Sequential: sets the Run Marker used by the command.
func TestSurfaceTranscriptThePackagePlanAtASettlement(t *testing.T) {
	repo := newGitRepository(t)
	runGit(t, repo, "checkout", "-b", "roundfix/run-fixture-task_03")
	t.Setenv("ROUNDFIX_RUN_ID", "fixture")
	var out, err bytes.Buffer
	code := verifyselect.Run(t.Context(), []string{"-repo", repo, "-plan"}, &out, &err)
	want := "verify-select: packages: none at a settlement; the Complete Suite Run runs them once per Spec\n"
	if code != 0 || out.Len() != 0 || err.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), err.String())
	}
}
func TestChangedPathsNameTheAddedArchiveRecord(t *testing.T) {
	t.Parallel()
	repo := newGitRepository(t)
	writeFile(t, repo, "base.txt", "base")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "base")
	base := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	writeFile(t, repo, "docs/history/specs/example.md", "archive")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "archive")
	paths, err := verifyselect.AddedPaths(t.Context(), repo, base)
	if err != nil || !reflect.DeepEqual(paths, []string{"docs/history/specs/example.md"}) {
		t.Fatal(paths, err)
	}
}
func TestSelectListsCommittedUnstagedAndUntrackedPaths(t *testing.T) {
	t.Parallel()
	repo := newGitRepository(t)
	writeFile(t, repo, "base.txt", "base")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "base")
	base := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	writeFile(t, repo, "internal/app/committed.go", "package app")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "change")
	writeFile(t, repo, "base.txt", "unstaged")
	writeFile(t, repo, "skills/example/SKILL.md", "untracked")
	paths, err := verifyselect.ChangedPaths(t.Context(), repo, base)
	want := []string{"base.txt", "internal/app/committed.go", "skills/example/SKILL.md"}
	if err != nil || !reflect.DeepEqual(paths, want) {
		t.Fatal(paths, err)
	}
}
func TestSelectListsStagedPaths(t *testing.T) {
	t.Parallel()
	repo := newGitRepository(t)
	writeFile(t, repo, "base.txt", "base")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "base")
	base := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	writeFile(t, repo, "internal/app/staged.go", "package app")
	runGit(t, repo, "add", ".")
	paths, err := verifyselect.ChangedPaths(t.Context(), repo, base)
	if err != nil || !reflect.DeepEqual(paths, []string{"internal/app/staged.go"}) {
		t.Fatal(paths, err)
	}
}
func TestPlanFailsSafeToEveryPackage(t *testing.T) {
	t.Parallel()
	for _, git := range []bool{false, true} {
		repo := t.TempDir()
		if git {
			repo = newGitRepository(t)
		}
		writeFile(t, repo, "go.mod", "module example.test/fixture\n\ngo 1.26\n")
		writeFile(t, repo, "core/core.go", "package core\n")
		plan, err := verifyselect.Select(t.Context(), repo, "missing-base")
		if err != nil || !reflect.DeepEqual(plan.Packages, []string{"./..."}) || !strings.Contains(plan.Reason, "Git could not list") {
			t.Fatal(plan, err)
		}
	}
}
func TestRunPrintsRequestedDefinitions(t *testing.T) {
	t.Parallel()
	repo := newGitRepository(t)
	writeFile(t, repo, "go.mod", "module example.test/fixture\n\ngo 1.26\n")
	writeFile(t, repo, "core/core.go", "package core\n")
	writeFile(t, repo, "internal/baseline/base.go", "package baseline\n")
	writeFile(t, repo, "internal/cli/baseline_fixture_test.go", "package cli\nimport \"testing\"\nfunc TestBaselineFixture(t *testing.T){}\n")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "base")
	writeFile(t, repo, "README.md", "change")
	var out, err bytes.Buffer
	code := verifyselect.Run(t.Context(), []string{"-repo", repo, "-base", "HEAD", "-job", "packages"}, &out, &err)
	if code != 0 || out.String() != "./core\n./internal/baseline\n" {
		t.Fatalf("%d %s %s", code, &out, &err)
	}
	out.Reset()
	err.Reset()
	code = verifyselect.Run(t.Context(), []string{"-repo", repo, "-base", "HEAD", "-job", "cli-1"}, &out, &err)
	if code != 0 || !strings.HasSuffix(out.String(), "./internal/cli\n") {
		t.Fatalf("%d %s %s", code, &out, &err)
	}
}
