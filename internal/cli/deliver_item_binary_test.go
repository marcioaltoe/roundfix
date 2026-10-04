// Boundary: real Git ignore checks, builds and child processes in disposable repositories.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
	"roundfix/internal/testfixture"
)

const deliveryItemFixtureSource = `package main
import ("encoding/json"; "fmt"; "os"; "strconv")
func main() {
 dir, err := os.Getwd(); if err != nil { panic(err) }
 file, err := os.OpenFile(os.Getenv("DELIVERY_ITEM_RECORD"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); if err != nil { panic(err) }
 err = json.NewEncoder(file).Encode(struct { Args []string; Dir, Home string }{os.Args[1:], dir, os.Getenv("HOME")}); if err != nil { panic(err) }; if err = file.Close(); err != nil { panic(err) }
 if len(os.Args) == 3 && os.Args[1] == "migrate" && os.Args[2] == "--check" {
  fmt.Fprint(os.Stderr, os.Getenv("DELIVERY_ITEM_STDERR")); fmt.Fprint(os.Stdout, os.Getenv("DELIVERY_ITEM_STDOUT"))
  code, err := strconv.Atoi(os.Getenv("DELIVERY_ITEM_PROBE_EXIT")); if err != nil { panic(err) }; os.Exit(code)
 }
 fmt.Fprintln(os.Stderr, "fixture step stopped after recording"); os.Exit(7)
}`

type deliveryItemStart struct {
	Args      []string
	Dir, Home string
}
type deliveryItemFixture struct {
	workflow                 *commandDeliveryWorkflow
	repo, home, record, head string
	log                      bytes.Buffer
}

func newDeliveryItemFixture(t *testing.T) *deliveryItemFixture {
	t.Helper()
	binary := testfixture.FixtureBinary(t, "item-roundfix", deliveryItemFixtureSource)
	root := t.TempDir()
	f := &deliveryItemFixture{repo: filepath.Join(root, "item"), home: filepath.Join(root, "home"), record: filepath.Join(root, "starts.jsonl")}
	mustMkdir(t, f.home)
	gittest.InitRepo(t, f.repo, "-b", "main")
	gittest.PersistIdentity(t, f.repo)
	if err := os.WriteFile(filepath.Join(f.repo, ".gitignore"), []byte("bin/\n"), 0600); err != nil {
		t.Fatal(err)
	}
	// The item's config is deliberately unreadable to the owner; only its loaded declaration is used.
	if err := os.WriteFile(filepath.Join(f.repo, ".roundfixrc.yml"), []byte("unknown-item-key: true\n"), 0600); err != nil {
		t.Fatal(err)
	}
	mustMkdir(t, filepath.Join(f.repo, "docs", "specs", "widget"))
	if err := os.WriteFile(filepath.Join(f.repo, "docs", "specs", "widget", "_prd.md"), []byte("# Widget\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, f.repo, "add", ".")
	gittest.Run(t, f.repo, "commit", "-m", "fixture")
	f.head = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", "HEAD"))
	db, err := store.Open(t.Context(), f.home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	f.workflow = &commandDeliveryWorkflow{store: db, git: preflight.ExecGitRunner{}, log: &f.log, loaded: roundconfig.Loaded{GitRoot: f.repo, HomeDir: f.home}}
	f.workflow.loaded.Config.Specs.Root = "docs/specs"
	f.workflow.loaded.Config.Defaults.ArtifactDir = filepath.Join(root, "artifacts")
	f.workflow.loaded.Config.Delivery.ItemBinary = roundconfig.ItemBinaryDeclaration{Path: "bin/roundfix", Build: "mkdir -p bin && cp '" + strings.ReplaceAll(binary, "'", "'\"'\"'") + "' bin/roundfix"}
	t.Setenv("DELIVERY_ITEM_RECORD", f.record)
	t.Setenv("DELIVERY_ITEM_PROBE_EXIT", "0")
	t.Setenv("DELIVERY_ITEM_STDERR", "")
	t.Setenv("DELIVERY_ITEM_STDOUT", "")
	return f
}

func (f *deliveryItemFixture) starts(t *testing.T) []deliveryItemStart {
	t.Helper()
	data, err := os.ReadFile(f.record)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var starts []deliveryItemStart
	for _, line := range bytes.Split(bytes.TrimSpace(data), []byte("\n")) {
		var start deliveryItemStart
		if err := json.Unmarshal(line, &start); err != nil {
			t.Fatal(err)
		}
		starts = append(starts, start)
	}
	return starts
}
func (f *deliveryItemFixture) selectionLine(step string) string {
	return fmt.Sprintf("roundfix: Delivery Queue item widget: %s runs the item binary %s\n", step, filepath.Join(f.repo, "bin/roundfix"))
}
func (f *deliveryItemFixture) assertStarts(t *testing.T, args ...[]string) {
	t.Helper()
	starts := f.starts(t)
	if len(starts) != len(args) {
		t.Fatalf("starts = %+v, want %v", starts, args)
	}
	for i, start := range starts {
		// Getwd can canonicalize the macOS temporary directory symlink.
		wantDir, err := filepath.EvalSymlinks(f.repo)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(start.Args, args[i]) || start.Dir != wantDir || start.Home != f.home {
			t.Fatalf("start = %+v, want args %v dir %s HOME %s", start, args[i], wantDir, f.home)
		}
	}
}

func TestDeliveryStepRunsTheItemBinary(t *testing.T) {
	f := newDeliveryItemFixture(t)
	executable, err := f.workflow.stepExecutable(t.Context(), f.repo, "widget", "implement")
	if err != nil || executable != filepath.Join(f.repo, "bin/roundfix") {
		t.Fatalf("executable=%q err=%v", executable, err)
	}
	if f.log.String() != f.selectionLine("implement") {
		t.Fatalf("log=%q", f.log.String())
	}
	f.assertStarts(t, []string{"migrate", "--check"})
}

func TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments(t *testing.T) {
	f := newDeliveryItemFixture(t) // Compile once for all three steps.
	for _, step := range []string{"implement", "archive", "review"} {
		t.Run(step, func(t *testing.T) {
			if err := os.Remove(f.record); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			f.log.Reset()
			// Removing the previous output proves that every step rebuilds it.
			if err := os.Remove(filepath.Join(f.repo, "bin/roundfix")); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			var err error
			var args []string
			switch step {
			case "implement":
				_, err = f.workflow.RunSpec(t.Context(), f.repo, "widget")
				args = []string{"implement", "--spec", "widget"}
			case "archive":
				_, err = f.workflow.Archive(t.Context(), f.repo, "widget", f.head)
				args = []string{"archive", "widget"}
			case "review":
				_, err = f.workflow.Review(t.Context(), f.repo, "widget", f.head)
				args = []string{"review"}
			}
			if err == nil || !strings.Contains(err.Error(), "fixture step stopped after recording") {
				t.Fatalf("step error=%v", err)
			}
			f.assertStarts(t, []string{"migrate", "--check"}, args)
			if f.log.String() != f.selectionLine(step) {
				t.Fatalf("log=%q", f.log.String())
			}
		})
	}
}

func TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate(t *testing.T) {
	f := newDeliveryItemFixture(t)
	t.Setenv("DELIVERY_ITEM_PROBE_EXIT", "2")
	t.Setenv(cliTestHelperEnv, "1")
	for _, stream := range []string{"stderr", "stdout"} {
		t.Run(stream, func(t *testing.T) {
			if err := os.Remove(f.record); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			f.log.Reset()
			t.Setenv("DELIVERY_ITEM_STDERR", "\n  \n")
			t.Setenv("DELIVERY_ITEM_STDOUT", "\nstdout reason\nlater\n")
			detail := "stdout reason"
			if stream == "stderr" {
				t.Setenv("DELIVERY_ITEM_STDERR", "\nprobe schema differs\nsecond line\n")
				detail = "probe schema differs"
			}
			executable, err := f.workflow.stepExecutable(t.Context(), f.repo, "widget", "review")
			owner, ownerErr := os.Executable()
			if err != nil || ownerErr != nil || executable != owner {
				t.Fatalf("executable=%q err=%v ownerErr=%v", executable, err, ownerErr)
			}
			result, err := f.workflow.runRoundfix(t.Context(), executable, f.repo, "--version")
			if err != nil || result.exitCode != 0 || result.stdout == "" {
				t.Fatalf("owner result=%+v err=%v", result, err)
			}
			want := "roundfix: notice: Delivery Queue item widget: review runs the owner's binary; the item binary's migrate --check exited 2: " + detail + "\n"
			if f.log.String() != want {
				t.Fatalf("log=%q want=%q", f.log.String(), want)
			}
			f.assertStarts(t, []string{"migrate", "--check"})
		})
	}
}

func TestDeliveryStepParksWhenTheItemBuildFails(t *testing.T) {
	f := newDeliveryItemFixture(t)
	f.workflow.loaded.Config.Delivery.ItemBinary.Build = "printf 'build failed\\n' >&2; exit 7"
	_, err := f.workflow.RunSpec(t.Context(), f.repo, "widget")
	logPath := filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "item-binary-build.log")
	if err == nil || !strings.HasPrefix(err.Error(), "build item binary: ") || !strings.Contains(err.Error(), "exit status 7") || !strings.Contains(err.Error(), logPath) {
		t.Fatalf("error=%v", err)
	}
	data, readErr := os.ReadFile(logPath)
	if readErr != nil || string(data) != "build failed\n" {
		t.Fatalf("build log=%q err=%v", data, readErr)
	}
	f.assertStarts(t)
	if f.log.Len() != 0 {
		t.Fatalf("log=%q", f.log.String())
	}
}

func TestDeliveryStepParksWhenTheItemBinaryPathIsNotIgnored(t *testing.T) {
	f := newDeliveryItemFixture(t)
	f.workflow.loaded.Config.Delivery.ItemBinary.Path = "unignored/roundfix"
	f.workflow.loaded.Config.Delivery.ItemBinary.Build = "touch build-started"
	_, err := f.workflow.Review(t.Context(), f.repo, "widget", f.head)
	if err == nil || err.Error() != `item binary path "unignored/roundfix" is not ignored by Git` {
		t.Fatalf("error=%v", err)
	}
	if _, err := os.Stat(filepath.Join(f.repo, "build-started")); !os.IsNotExist(err) {
		t.Fatalf("build started: %v", err)
	}
	f.assertStarts(t)
	if f.log.Len() != 0 {
		t.Fatalf("log=%q", f.log.String())
	}
}

func TestDeliveryStepWithoutADeclarationRunsTheOwnerExecutable(t *testing.T) {
	f := newDeliveryItemFixture(t)
	f.workflow.loaded.Config.Delivery.ItemBinary = roundconfig.ItemBinaryDeclaration{}
	f.workflow.git = nil // Any Git access would panic.
	t.Setenv(cliTestHelperEnv, "1")
	executable, err := f.workflow.stepExecutable(t.Context(), f.repo, "widget", "implement")
	owner, ownerErr := os.Executable()
	if err != nil || ownerErr != nil || executable != owner {
		t.Fatalf("executable=%q err=%v ownerErr=%v", executable, err, ownerErr)
	}
	result, err := f.workflow.runRoundfix(t.Context(), executable, f.repo, "--version")
	if err != nil || result.exitCode != 0 || result.stdout == "" {
		t.Fatalf("owner result=%+v err=%v", result, err)
	}
	f.assertStarts(t)
	if _, err := os.Stat(filepath.Join(f.repo, "bin")); !os.IsNotExist(err) {
		t.Fatalf("build ran: %v", err)
	}
	if f.log.Len() != 0 {
		t.Fatalf("log=%q", f.log.String())
	}
}

func TestDeliveryItemBinaryThatCannotStartReturnsAParkError(t *testing.T) {
	f := newDeliveryItemFixture(t)
	f.workflow.loaded.Config.Delivery.ItemBinary.Build = "mkdir -p bin"
	_, err := f.workflow.Archive(t.Context(), f.repo, "widget", f.head)
	if err == nil || !strings.HasPrefix(err.Error(), fmt.Sprintf("run item binary %q: ", filepath.Join(f.repo, "bin/roundfix"))) {
		t.Fatalf("error=%v", err)
	}
	f.assertStarts(t)
	if f.log.Len() != 0 {
		t.Fatalf("log=%q", f.log.String())
	}
}
