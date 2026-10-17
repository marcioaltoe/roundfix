// Boundary: rebuilding a real item binary while a pipe holds its previous build running.
package cli

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDeliveryStepProbesARebuiltItemBinaryWhileThePreviousBuildStillRuns(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	executable, err, first := f.step(t, commandContextForTest(t, t.Context()), "review")
	if err != nil || executable != f.copyPath(t) {
		t.Fatalf("first executable=%q err=%v; %s", executable, err, first)
	}

	held := exec.Command(executable, "hold")
	held.Dir = f.repo
	// Only this test holds the write end, so test-binary exit also releases the holder.
	input, err := held.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := held.Start(); err != nil {
		_ = input.Close()
		t.Fatal(err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			t.Error(err)
		}
		if err := held.Wait(); err != nil {
			t.Errorf("held item binary exit=%v, want exit 0", err)
		}
	}()

	f.log.Reset()
	executable, err, report := f.step(t, commandContextForTest(t, t.Context()), "review")
	if err != nil || executable != f.copyPath(t) || f.log.String() != f.selectionLine(t, "review") {
		t.Fatalf("second executable=%q err=%v log=%q want=%q; %s", executable, err, f.log.String(), f.selectionLine(t, "review"), report)
	}
	if report.InodeBefore == 0 || report.InodeBefore != first.InodeAfter || report.InodeAfter == 0 || report.InodeAfter == report.InodeBefore || !strings.Contains(report.String(), "(replaced)") {
		t.Fatalf("want the previous build replaced by a new inode; %s", report)
	}
	if report.Start == nil || report.ContextErr != nil {
		t.Fatalf("want a recorded probe with a live context; %s", report)
	}
	f.assertStarts(t, []string{"migrate", "--check"}, []string{"migrate", "--check"})
}
