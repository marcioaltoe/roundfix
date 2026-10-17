// Boundary: real builds and running executable copies in disposable repositories.
package cli

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func (f *deliveryItemFixture) copyPath(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.repo, "bin/roundfix"))
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	return filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "item-binary", fmt.Sprintf("%x", hash)[:16], "roundfix")
}

func TestDeliveryStepRunsACopyOfTheItemBinary(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	ctx := commandContextForTest(t, t.Context())
	executable, err := f.workflow.stepExecutable(ctx, f.repo, "widget", "implement")
	if err != nil {
		t.Fatal(err)
	}
	if executable != f.copyPath(t) {
		t.Fatalf("executable=%q want=%q", executable, f.copyPath(t))
	}
	built := mustRead(t, filepath.Join(f.repo, "bin/roundfix"))
	if got := mustRead(t, executable); got != built {
		t.Fatal("copy differs from build output")
	}
	info, err := os.Stat(executable)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0755 {
		t.Fatalf("mode=%v want 0755", info.Mode())
	}
	wantLine := fmt.Sprintf("roundfix: Delivery Queue item widget: implement runs the item binary %s\n", executable)
	if f.log.String() != wantLine {
		t.Fatalf("log=%q want=%q", f.log.String(), wantLine)
	}
	result, err := f.workflow.runRoundfix(ctx, executable, f.repo, "implement", "--spec", "widget")
	if err != nil || result.exitCode != 7 {
		t.Fatalf("step=%+v err=%v", result, err)
	}
	f.assertStarts(t, []string{"migrate", "--check"}, []string{"implement", "--spec", "widget"})
	for _, start := range f.starts(t) {
		if start.Executable != executable {
			t.Fatalf("executed=%q want copy=%q", start.Executable, executable)
		}
	}
}

func TestDeliveryRebuildInPlaceNeverRewritesARunningCopy(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	ctx := commandContextForTest(t, t.Context())
	first, err := f.workflow.stepExecutable(ctx, f.repo, "widget", "review")
	if err != nil {
		t.Fatal(err)
	}
	if first != f.copyPath(t) {
		t.Fatalf("first executable=%q want copy=%q", first, f.copyPath(t))
	}
	original := []byte(mustRead(t, first))
	originalInfo, err := os.Stat(first)
	if err != nil {
		t.Fatal(err)
	}
	buildInfo, err := os.Stat(filepath.Join(f.repo, "bin/roundfix"))
	if err != nil {
		t.Fatal(err)
	}
	held := exec.CommandContext(t.Context(), first, "hold")
	held.Dir = f.repo
	// Only this test holds the input writer; test-binary exit also releases the child.
	input, err := held.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	ready, err := held.StdoutPipe()
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
			t.Errorf("held copy exit=%v, want exit 0", err)
		}
	}()
	if line, err := bufio.NewReader(ready).ReadString('\n'); err != nil || line != "ready\n" {
		t.Fatalf("holder ready=%q err=%v", line, err)
	}
	rebuilt := testfixture.FixtureBinary(t, "rebuilt-item-roundfix", strings.ReplaceAll(deliveryItemFixtureSource, "fixture step stopped", "rebuilt step stopped"))
	// Deliberately rebuild in place: only the artifact copy may ever have run.
	f.workflow.loaded.Config.Delivery.ItemBinary.Build = "mkdir -p bin && cp '" + strings.ReplaceAll(rebuilt, "'", "'\"'\"'") + "' bin/roundfix"
	f.log.Reset()
	second, err := f.workflow.stepExecutable(ctx, f.repo, "widget", "review")
	if err != nil {
		t.Fatal(err)
	}
	if second != f.copyPath(t) || second == first || second == filepath.Join(f.repo, "bin/roundfix") {
		t.Fatalf("first=%q second=%q want=%q", first, second, f.copyPath(t))
	}
	afterBuild, err := os.Stat(filepath.Join(f.repo, "bin/roundfix"))
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(buildInfo, afterBuild) {
		t.Fatal("fixture did not rebuild the build output in place")
	}
	afterHeld, err := os.Stat(first)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(originalInfo, afterHeld) || !bytes.Equal(original, []byte(mustRead(t, first))) {
		t.Fatal("in-place rebuild changed the held copy's bytes or identity")
	}
	if f.log.String() != f.selectionLine(t, "review") {
		t.Fatalf("log=%q", f.log.String())
	}
	f.assertStarts(t, []string{"migrate", "--check"}, []string{"migrate", "--check"})
	starts := f.starts(t)
	if starts[0].Executable != first || starts[1].Executable != second {
		t.Fatalf("executed paths=%+v", starts)
	}
	// Repeating identical bytes must also replace the copy, never write it in place.
	beforeRepeat, err := os.Stat(second)
	if err != nil {
		t.Fatal(err)
	}
	repeated, err := f.workflow.stepExecutable(ctx, f.repo, "widget", "review")
	if err != nil || repeated != second {
		t.Fatalf("repeated=%q err=%v", repeated, err)
	}
	afterRepeat, err := os.Stat(second)
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(beforeRepeat, afterRepeat) {
		t.Fatal("identical build rewrote the copy in place")
	}
	root := filepath.Dir(filepath.Dir(second))
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && entry.Name() != "roundfix" {
			t.Errorf("temporary file remains: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeliveryItemBinaryCopyFailureNeverRunsTheBuildOutput(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"directory", "blocked destination"} {
		t.Run(kind, func(t *testing.T) {
			f := newDeliveryItemFixture(t)
			ctx := commandContextForTest(t, t.Context())
			var directory string
			if kind == "directory" {
				f.workflow.loaded.Config.Delivery.ItemBinary.Build = "mkdir -p bin/roundfix"
			} else {
				executable, err := f.workflow.stepExecutable(ctx, f.repo, "widget", "review")
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(executable); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(executable, 0755); err != nil {
					t.Fatal(err)
				}
				directory = filepath.Dir(executable)
				if err := os.Remove(f.record); err != nil {
					t.Fatal(err)
				}
				f.log.Reset()
			}
			_, err := f.workflow.Review(ctx, f.repo, "widget", f.head)
			if err == nil || !strings.HasPrefix(err.Error(), "copy item binary: ") {
				t.Fatalf("error=%v", err)
			}
			f.assertStarts(t)
			if f.log.Len() != 0 {
				t.Fatalf("selection logged after copy failure: %q", f.log.String())
			}
			if directory != "" {
				entries, err := os.ReadDir(directory)
				if err != nil {
					t.Fatal(err)
				}
				if len(entries) != 1 || entries[0].Name() != "roundfix" {
					t.Fatalf("copy failure left temporary files: %v", entries)
				}
			}
		})
	}
}
