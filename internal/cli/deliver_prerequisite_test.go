// Suite: prerequisite delivery preflight and refreshed Git evidence.
// Boundary IN: public CLI, SQLite, manifest loading, and disposable local Git repos.
// Boundary OUT: owner launch and network remotes; no test reaches a network.
package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func addFixtureRequires(t *testing.T, repo, slug, declaration string) {
	t.Helper()
	path := filepath.Join(repo, "docs", "specs", slug, "_tasks.md")
	mustWrite(t, path, strings.Replace(mustRead(t, path), "schema:", "requires: "+declaration+"\nschema:", 1))
}

func TestDeliverStartRefusesAPrerequisiteCycle(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	const second = "0205-second"
	writeImplementSpec(t, repo, second, []implementSeed{{id: "task_01"}})
	addFixtureRequires(t, repo, implementTestSlug, "["+second+"]")
	addFixtureRequires(t, repo, second, "["+implementTestSlug+"]")
	reason := "Delivery Queue Specs require each other in a cycle: " + implementTestSlug + " -> " + second + " -> " + implementTestSlug
	assertPrerequisiteStartRefusal(t, home, repo, []string{implementTestSlug, second}, reason)
}

func TestDeliverStartRefusesAnUnknownPrerequisite(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	addFixtureRequires(t, repo, implementTestSlug, "[0205-unknown]")
	assertPrerequisiteStartRefusal(t, home, repo, []string{implementTestSlug}, "Delivery Queue Spec "+implementTestSlug+" requires unknown Spec 0205-unknown")
}

func assertPrerequisiteStartRefusal(t *testing.T, home, repo string, slugs []string, reason string) {
	t.Helper()
	launched := false
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			launched = true
			return exitOK
		}
	})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, append([]string{"deliver", "start"}, slugs...), &stdout, &stderr)
	want := "Preflight failed\n\nReason:\n  " + reason + "\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix deliver start --help' for usage.\n"
	if code != 2 || stdout.Len() != 0 || stderr.String() != want || launched {
		t.Fatalf("exit=%d stdout=%q stderr=%q launched=%v; want stderr=%q", code, stdout.String(), stderr.String(), launched, want)
	}
	// Fixtures have no delivery grant, so the reason also proves prerequisites
	// are checked before authorization.
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, found, err := db.DeliveryQueue(t.Context(), repo); err != nil || found {
		t.Fatalf("queue found=%v err=%v", found, err)
	}
}

func TestAPrerequisiteIsMetByItsArchiveOnTheRefreshedDefaultBranch(t *testing.T) {
	t.Parallel()
	for _, remote := range []string{"origin", "delivery"} {
		t.Run(remote, func(t *testing.T) {
			origin := filepath.Join(t.TempDir(), "origin")
			gittest.InitRepo(t, origin, "-b", "main")
			manifest := "---\nschema: spec-tasks/v1\nrequires: [met, missing, local-only]\n---\n"
			path := filepath.Join(origin, "docs", "specs", "dependent", "_tasks.md")
			mustMkdir(t, filepath.Dir(path))
			mustWrite(t, path, manifest)
			gittest.Run(t, origin, "add", ".")
			gittest.Run(t, origin, "commit", "-m", "seed")
			checkout := filepath.Join(t.TempDir(), "checkout")
			gittest.Run(t, "", "clone", origin, checkout)
			gittest.Harden(t, checkout)
			if remote != "origin" {
				gittest.Run(t, checkout, "remote", "add", remote, origin)
			}
			// Advance the remote after cloning, while leaving the owner checkout stale.
			archived := filepath.Join(origin, "docs", "history", "specs", "met", "_prd.md")
			mustMkdir(t, filepath.Dir(archived))
			mustWrite(t, archived, "---\nstatus: archived\n---\n")
			gittest.Run(t, origin, "add", ".")
			gittest.Run(t, origin, "commit", "-m", "archive prerequisite")
			old := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "refs/remotes/origin/main"))
			newest := strings.TrimSpace(gittest.Run(t, origin, "rev-parse", "HEAD"))
			if old == newest {
				t.Fatal("fixture is already fetched")
			}
			local := filepath.Join(checkout, "docs", "history", "specs", "local-only", "_prd.md")
			mustMkdir(t, filepath.Dir(local))
			mustWrite(t, local, "local archive")
			// Change the checkout declaration too: only the default branch is trusted.
			mustWrite(t, filepath.Join(checkout, "docs", "specs", "dependent", "_tasks.md"), "---\nrequires: []\n---\n")
			config := roundconfig.Config{}
			config.Specs.Root = "docs/specs"
			if remote != "origin" {
				config.Watch.PushRemote = remote
			}
			workflow := &commandDeliveryWorkflow{loaded: roundconfig.Loaded{GitRoot: checkout, Config: config}, git: preflight.ExecGitRunner{}}
			unmet, err := workflow.UnmetPrerequisites(t.Context(), checkout, "dependent")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(unmet, []string{"missing", "local-only"}) {
				t.Fatalf("unmet=%v", unmet)
			}
			got := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "refs/remotes/"+remote+"/main"))
			if got != newest {
				t.Fatalf("ref=%s want=%s", got, newest)
			}
			if _, err := os.Stat(filepath.Join(checkout, "docs", "history", "specs", "met", "_prd.md")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("archive reached checkout: %v", err)
			}
		})
	}
}

type prerequisiteGitFailure struct {
	commands []string
	fail     string
}

func (runner *prerequisiteGitFailure) RunGit(_ context.Context, _ string, args ...string) (string, error) {
	command := strings.Join(args, " ")
	runner.commands = append(runner.commands, command)
	if args[0] == "symbolic-ref" {
		return "refs/remotes/origin/main", nil
	}
	if args[0] == runner.fail {
		return "", errors.New("read failed")
	}
	if args[0] == "show" {
		return "---\nrequires: [missing]\n---\n", nil
	}
	return "", nil
}
func TestPrerequisiteGitFailuresRemainErrors(t *testing.T) {
	t.Parallel()
	for _, fail := range []string{"fetch", "show", "ls-tree"} {
		_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
		runner := &prerequisiteGitFailure{fail: fail}
		config := roundconfig.Config{}
		config.Specs.Root = "docs/specs"
		workflow := &commandDeliveryWorkflow{loaded: roundconfig.Loaded{GitRoot: repo, Config: config}, git: runner}
		if _, err := workflow.UnmetPrerequisites(t.Context(), repo, implementTestSlug); err == nil {
			t.Fatalf("%s failure was treated as absence", fail)
		}
	}
}

func TestDeliverPrerequisitePreflightAcceptsActiveAndArchivedSpecs(t *testing.T) {
	t.Parallel()
	for _, builtIn := range []bool{false, true} {
		root := roundconfig.SpecsRoot{Path: filepath.Join(t.TempDir(), "docs", "specs"), BuiltInRoot: builtIn}
		archive := filepath.Join(root.Path, "_archived")
		if builtIn {
			archive = filepath.Join(filepath.Dir(root.Path), "history", "specs")
		}
		for _, directory := range []string{root.Path, archive} {
			if err := os.MkdirAll(filepath.Join(directory, "required"), 0755); err != nil {
				t.Fatal(err)
			}
			mustWrite(t, filepath.Join(directory, "required", "_prd.md"), "Spec")
		}
		graphs := []*spec.Graph{{Spec: spec.Spec{Slug: "dependent"}, Requires: []string{"required"}}}
		if err := validateDeliveryPrerequisites(root, graphs); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(filepath.Join(root.Path, "required")); err != nil {
			t.Fatal(err)
		}
		if err := validateDeliveryPrerequisites(root, graphs); err != nil {
			t.Fatalf("archived prerequisite: %v", err)
		}
	}
}

func TestDeliverPrerequisitePreflightReportsLongerCyclesDeterministically(t *testing.T) {
	t.Parallel()
	root := roundconfig.SpecsRoot{Path: t.TempDir()}
	var graphs []*spec.Graph
	for index, slug := range []string{"first", "second", "third"} {
		if err := os.MkdirAll(filepath.Join(root.Path, slug), 0755); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(root.Path, slug, "_prd.md"), "Spec")
		required := []string{"second", "third", "first"}[index]
		graphs = append(graphs, &spec.Graph{Spec: spec.Spec{Slug: slug}, Requires: []string{required}})
	}
	want := "Delivery Queue Specs require each other in a cycle: first -> second -> third -> first"
	if err := validateDeliveryPrerequisites(root, graphs); err == nil || err.Error() != want {
		t.Fatalf("cycle error=%v", err)
	}
}

func TestPrerequisiteReaderWithoutRequiresSkipsArchiveRead(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	config := roundconfig.Config{}
	config.Specs.Root = "docs/specs"
	runner := &noRequiresGitRunner{}
	workflow := &commandDeliveryWorkflow{loaded: roundconfig.Loaded{GitRoot: repo, Config: config}, git: runner}
	unmet, err := workflow.UnmetPrerequisites(t.Context(), repo, implementTestSlug)
	if err != nil || len(unmet) != 0 {
		t.Fatalf("unmet=%v err=%v", unmet, err)
	}
	if runner.calls != 3 {
		t.Fatalf("Git calls=%d want default detection, fetch, manifest", runner.calls)
	}
}

type noRequiresGitRunner struct{ calls int }

func (runner *noRequiresGitRunner) RunGit(_ context.Context, _ string, args ...string) (string, error) {
	runner.calls++
	switch args[0] {
	case "symbolic-ref":
		return "refs/remotes/origin/main", nil
	case "fetch":
		return "", nil
	case "show":
		return "---\nschema: spec-tasks/v1\n---\n", nil
	default:
		return "", errors.New("unexpected archive read")
	}
}
