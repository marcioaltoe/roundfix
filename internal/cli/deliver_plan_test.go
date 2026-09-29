// Suite: Delivery Plan and queue authorization preflight.
// Invariant: the public CLI reports every approval fact without writing state, and start records only fully authorized Specs.
// Boundary IN: public CLI dispatch, committed Spec artifacts, Git, and the real SQLite Delivery Queue store.
// Boundary OUT: detached owner launch, injected as an operating-system boundary.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const deliverPlanSecondSlug = "0002-second-flow"

var allDeliveryOperations = []string{"implement", "commit", "push", "pull_request", "merge"}

type deliverPlanFixtureSpec struct {
	slug       string
	operations []string
	premises   []string
}

func TestDeliverPlanReportsAnApprovedAndABlockedSpec(t *testing.T) {
	t.Parallel()
	_, _ = newDeliverPlanWorkspace(t,
		deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations},
		deliverPlanFixtureSpec{slug: deliverPlanSecondSlug, operations: allDeliveryOperations[:4]},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug, deliverPlanSecondSlug}, &stdout, &stderr)

	if code != exitRunFailed || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, want := range []string{
		"spec\t" + implementTestSlug + "\tapproved\t1/1\t-",
		"spec\t" + deliverPlanSecondSlug + "\tblocked\t1/1\tauthorization lacks merge",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("deliver plan output does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestDeliverPlanExitsZeroWhenEverySpecIsApproved(t *testing.T) {
	t.Parallel()
	_, repoDir := newDeliverPlanWorkspace(t,
		deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations},
		deliverPlanFixtureSpec{slug: deliverPlanSecondSlug, operations: allDeliveryOperations},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug, deliverPlanSecondSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		findings, _ := strictSpecFindings(filepath.Join(repoDir, "docs", "specs"), repoDir, implementTestSlug)
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q findings=%+v", code, stderr.String(), stdout.String(), findings)
	}
	if strings.Contains(stdout.String(), "\tblocked\t") {
		t.Fatalf("approved plan contains blocked verdict:\n%s", stdout.String())
	}
}

func TestDeliverPlanJSONCarriesTheSameFacts(t *testing.T) {
	t.Parallel()
	_, repoDir := newDeliverPlanWorkspace(t,
		deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations, premises: []string{"internal/shared.go"}},
		deliverPlanFixtureSpec{slug: deliverPlanSecondSlug, operations: allDeliveryOperations[:4], premises: []string{"internal/shared.go"}},
	)
	writeDeliverPlanIntent(t, repoDir, "docs/backlog/2026-09-28-plan-json.md", "open")
	gitImplement(t, repoDir, "add", "docs/backlog/2026-09-28-plan-json.md")
	gitImplement(t, repoDir, "commit", "-m", "seed plan JSON intent")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", "--json", implementTestSlug, deliverPlanSecondSlug}, &stdout, &stderr)

	if code != exitRunFailed || stderr.Len() != 0 {
		t.Fatalf("deliver plan --json exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var document struct {
		Schema string `json:"schema"`
		Specs  []struct {
			Slug           string   `json:"slug"`
			Verdict        string   `json:"verdict"`
			Tasks          int      `json:"tasks"`
			Unfinished     int      `json:"unfinishedTasks"`
			Reasons        []string `json:"reasons"`
			SharedPremises []struct {
				With  string   `json:"with"`
				Paths []string `json:"paths"`
			} `json:"sharedPremises"`
		} `json:"specs"`
		Intent []struct {
			Kind   string `json:"kind"`
			Path   string `json:"path"`
			Status string `json:"status"`
		} `json:"intent"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatalf("decode deliver plan JSON: %v\n%s", err, stdout.String())
	}
	if document.Schema != "roundfix-deliver-plan/v1" || len(document.Specs) != 2 {
		t.Fatalf("deliver plan JSON header/specs = %+v", document)
	}
	if first := document.Specs[0]; first.Slug != implementTestSlug || first.Verdict != "approved" || first.Tasks != 1 || first.Unfinished != 1 || len(first.Reasons) != 0 {
		t.Fatalf("first plan Spec = %+v", first)
	}
	second := document.Specs[1]
	if second.Slug != deliverPlanSecondSlug || second.Verdict != "blocked" || len(second.Reasons) != 1 || second.Reasons[0] != "authorization lacks merge" {
		t.Fatalf("second plan Spec = %+v", second)
	}
	if len(second.SharedPremises) != 1 || second.SharedPremises[0].With != implementTestSlug || len(second.SharedPremises[0].Paths) != 1 || second.SharedPremises[0].Paths[0] != "internal/shared.go" {
		t.Fatalf("second plan shared premises = %+v", second.SharedPremises)
	}
	if len(document.Intent) != 1 || document.Intent[0].Kind != "backlog" || document.Intent[0].Status != "open" {
		t.Fatalf("deliver plan JSON intent = %+v", document.Intent)
	}
}

func TestDeliverPlanNamesSharedProductionPremises(t *testing.T) {
	t.Parallel()
	_, repoDir := newDeliverPlanWorkspace(t,
		deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations, premises: []string{"internal/shared.go"}},
		deliverPlanFixtureSpec{slug: deliverPlanSecondSlug, operations: allDeliveryOperations, premises: []string{"internal/shared.go", "internal/other.go"}},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug, deliverPlanSecondSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		findings, _ := strictSpecFindings(filepath.Join(repoDir, "docs", "specs"), repoDir, implementTestSlug)
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q findings=%+v", code, stderr.String(), stdout.String(), findings)
	}
	want := "shared\t" + deliverPlanSecondSlug + "\t" + implementTestSlug + "\tinternal/shared.go\n"
	if !strings.Contains(stdout.String(), want) {
		t.Fatalf("deliver plan shared row = %q, want %q", stdout.String(), want)
	}
	if strings.Count(stdout.String(), "\tapproved\t") != 2 {
		t.Fatalf("shared premise changed verdicts:\n%s", stdout.String())
	}
}

func TestDeliverPlanIgnoresSharedTestsAndGuides(t *testing.T) {
	t.Parallel()
	ignored := []string{"internal/shared_test.go", "docs/user-guide/shared.md"}
	_, _ = newDeliverPlanWorkspace(t,
		deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations, premises: ignored},
		deliverPlanFixtureSpec{slug: deliverPlanSecondSlug, operations: allDeliveryOperations, premises: ignored},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug, deliverPlanSecondSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if strings.Contains(stdout.String(), "shared\t") {
		t.Fatalf("deliver plan reported non-production shared premises:\n%s", stdout.String())
	}
}

func TestDeliverPlanListsIntentThatIsNotApprovedToRun(t *testing.T) {
	t.Parallel()
	_, repoDir := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations})
	writeDeliverPlanIntent(t, repoDir, "docs/backlog/2026-09-28-next.md", "open")
	writeDeliverPlanIntent(t, repoDir, "docs/findings/2026-09-28-observed.md", "pending")
	mustMkdir(t, filepath.Join(repoDir, "docs", "_inbox"))
	mustWrite(t, filepath.Join(repoDir, "docs", "_inbox", "note.txt"), "raw note\n")
	gitImplement(t, repoDir, "add", "docs/backlog", "docs/findings", "docs/_inbox")
	gitImplement(t, repoDir, "commit", "-m", "seed open intent")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, want := range []string{
		"backlog\tdocs/backlog/2026-09-28-next.md\topen",
		"finding\tdocs/findings/2026-09-28-observed.md\tpending",
		"inbox\tdocs/_inbox/note.txt\t-",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("deliver plan intent does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestDeliverPlanWritesNothing(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations})
	headBefore := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	statusBefore := gitImplementOutput(t, repoDir, "status", "--short")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if headAfter := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD")); headAfter != headBefore {
		t.Fatalf("deliver plan moved HEAD from %s to %s", headBefore, headAfter)
	}
	if statusAfter := gitImplementOutput(t, repoDir, "status", "--short"); statusAfter != statusBefore {
		t.Fatalf("deliver plan changed status from %q to %q", statusBefore, statusAfter)
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !os.IsNotExist(err) {
		t.Fatalf("deliver plan created Run Database: err=%v", err)
	}
}

func TestDeliverPlanDefaultsToEveryActiveSpec(t *testing.T) {
	t.Parallel()
	_, _ = newDeliverPlanWorkspace(t,
		deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations},
		deliverPlanFixtureSpec{slug: deliverPlanSecondSlug, operations: allDeliveryOperations},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, slug := range []string{implementTestSlug, deliverPlanSecondSlug} {
		if !strings.Contains(stdout.String(), "spec\t"+slug+"\tapproved") {
			t.Fatalf("default plan omits %s:\n%s", slug, stdout.String())
		}
	}
}

func TestDeliverPlanRefusesAnUnknownSlug(t *testing.T) {
	t.Parallel()
	homeDir, _ := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", "9999-unknown-spec"}, &stdout, &stderr)

	if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), "9999-unknown-spec") {
		t.Fatalf("unknown plan slug exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !os.IsNotExist(err) {
		t.Fatalf("unknown plan slug created Run Database: err=%v", err)
	}
}

func TestDeliverPlanRefusesAnUnknownFlag(t *testing.T) {
	t.Parallel()
	homeDir, _ := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", "--unknown", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("unknown plan flag exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !os.IsNotExist(err) {
		t.Fatalf("unknown plan flag created Run Database: err=%v", err)
	}
}

func TestDeliverStartRefusesASpecWithoutDeliveryAuthority(t *testing.T) {
	t.Parallel()
	homeDir, _ := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations[:4]})
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight || stdout.Len() != 0 || started != 0 {
		t.Fatalf("unauthorized start exit=%d starts=%d stdout=%q stderr=%q", code, started, stdout.String(), stderr.String())
	}
	for _, want := range []string{implementTestSlug, "authorization lacks merge", "roundfix deliver plan"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("unauthorized start diagnostic %q does not contain %q", stderr.String(), want)
		}
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !os.IsNotExist(err) {
		t.Fatalf("unauthorized start created Run Database: err=%v", err)
	}
}

func TestDeliverStartAcceptsASpecWithEveryDeliveryOperation(t *testing.T) {
	t.Parallel()
	_, _ = newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{slug: implementTestSlug, operations: allDeliveryOperations})
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 || started != 1 {
		t.Fatalf("authorized start exit=%d starts=%d stderr=%q stdout=%q", code, started, stderr.String(), stdout.String())
	}
}

func TestDeliverHelpNamesThePlanCommand(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "--help"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver help exit=%d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{"roundfix deliver plan [--json] [<slug>...]", "plan    Report which Specs are approved to run"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("deliver help does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func newDeliverPlanWorkspace(t *testing.T, fixtures ...deliverPlanFixtureSpec) (string, string) {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	for _, fixture := range fixtures {
		if fixture.slug != implementTestSlug {
			writeImplementSpec(t, repoDir, fixture.slug, []implementSeed{{id: "task_01"}})
		}
		authorizationPath := filepath.Join(repoDir, "docs", "specs", fixture.slug, "_authorization.md")
		mustWrite(t, authorizationPath, implementFixtureAuthorization(fixture.slug, fixture.operations...))
		if len(fixture.premises) > 0 {
			taskPath := filepath.Join(repoDir, "docs", "specs", fixture.slug, "task_01.md")
			content := mustRead(t, taskPath) + "\n## Context\n\n"
			for _, premise := range fixture.premises {
				content += "- interface: `" + premise + "`\n"
				premisePath := filepath.Join(repoDir, filepath.FromSlash(premise))
				mustMkdir(t, filepath.Dir(premisePath))
				mustWrite(t, premisePath, "package fixture\n")
			}
			mustWrite(t, taskPath, content)
		}
	}
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "prepare Delivery Plan fixtures")
	return homeDir, repoDir
}

func writeDeliverPlanIntent(t *testing.T, repoDir, relativePath, status string) {
	t.Helper()
	path := filepath.Join(repoDir, filepath.FromSlash(relativePath))
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, "---\nstatus: "+status+"\n---\n\n# Fixture intent\n")
}
