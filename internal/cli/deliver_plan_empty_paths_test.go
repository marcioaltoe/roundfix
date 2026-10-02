// Suite: Delivery Plan explicit empty authorization paths.
// Invariant: an explicit empty path list grants delivery operations, while a null field grants nothing.
// Boundary IN: public Delivery Plan/start dispatch, committed Spec fixtures, Git, and the real queue store.
// Boundary OUT: detached owner launch, injected as an operating-system boundary.
package cli

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestDeliverPlanApprovesAnExplicitEmptyPathsGrant(t *testing.T) {
	t.Parallel()

	_, repoDir := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{
		slug: implementTestSlug, operations: allDeliveryOperations,
	})
	setDeliverPlanFixturePaths(t, repoDir, "paths: []\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "spec\t"+implementTestSlug+"\tapproved\t1/1\t-") {
		t.Fatalf("deliver plan did not approve the explicit empty paths grant:\n%s", stdout.String())
	}
}

func TestDeliverPlanRefusesANullPathsGrant(t *testing.T) {
	t.Parallel()

	_, repoDir := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{
		slug: implementTestSlug, operations: allDeliveryOperations,
	})
	setDeliverPlanFixturePaths(t, repoDir, "paths:\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "plan", implementTestSlug}, &stdout, &stderr)

	if code != exitRunFailed || stderr.Len() != 0 {
		t.Fatalf("deliver plan exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, want := range []string{"spec\t" + implementTestSlug + "\tblocked", "authorization refused: paths"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("deliver plan output does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestDeliverStartAcceptsAnExplicitEmptyPathsGrant(t *testing.T) {
	t.Parallel()

	_, repoDir := newDeliverPlanWorkspace(t, deliverPlanFixtureSpec{
		slug: implementTestSlug, operations: allDeliveryOperations,
	})
	setDeliverPlanFixturePaths(t, repoDir, "paths: []\n")
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.deliveryReadiness = readyDeliveryReadiness
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 || started != 1 {
		t.Fatalf("deliver start exit=%d starts=%d stderr=%q stdout=%q", code, started, stderr.String(), stdout.String())
	}
}

func setDeliverPlanFixturePaths(t *testing.T, repoDir, declaration string) {
	t.Helper()

	recordPath := filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_authorization.md")
	record := mustRead(t, recordPath)
	updated := strings.Replace(record, "paths:\n  - docs/agents/domain.md\n", declaration, 1)
	if updated == record {
		t.Fatal("Delivery Plan fixture does not contain its expected paths declaration")
	}
	record = updated
	mustWrite(t, recordPath, record)
	gitImplement(t, repoDir, "add", filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_authorization.md")))
	gitImplement(t, repoDir, "commit", "-m", "record empty paths fixture")
}
