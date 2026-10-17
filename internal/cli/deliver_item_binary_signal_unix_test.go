//go:build darwin || linux

// Boundary: item binaries signal only their own PID in disposable repositories.
package cli

import (
	"os"
	"testing"

	"roundfix/internal/testfixture"
)

const deliverySignalItemFixtureSource = `package main
import ("encoding/json"; "fmt"; "os"; "syscall")
func main() {
 dir, err := os.Getwd(); if err != nil { panic(err) }
 file, err := os.OpenFile(os.Getenv("DELIVERY_ITEM_RECORD"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); if err != nil { panic(err) }
 err = json.NewEncoder(file).Encode(struct { Args []string; Dir, Home string }{os.Args[1:], dir, os.Getenv("HOME")}); if err != nil { panic(err) }; if err = file.Close(); err != nil { panic(err) }
 if os.Args[1] == os.Getenv("DELIVERY_ITEM_SIGNAL_STEP") {
  fmt.Fprint(os.Stderr, os.Getenv("DELIVERY_ITEM_STDERR")); fmt.Fprint(os.Stdout, os.Getenv("DELIVERY_ITEM_STDOUT"))
  if err := syscall.Kill(os.Getpid(), syscall.SIGKILL); err != nil { panic(err) }
  panic("self signal returned")
 }
}`

func newDeliverySignalItemFixture(t *testing.T, step string) *deliveryItemFixture {
	t.Helper()
	f := newDeliveryItemFixture(t)
	binary := testfixture.FixtureBinary(t, "signal-item-roundfix", deliverySignalItemFixtureSource)
	f.workflow.loaded.Config.Delivery.ItemBinary.Build = deliveryItemBuildCommand(binary)
	setCommandEnvForTest(t, "DELIVERY_ITEM_SIGNAL_STEP", step)
	return f
}

func TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe(t *testing.T) {
	t.Parallel()
	f := newDeliverySignalItemFixture(t, "migrate")
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	executable, err := f.workflow.stepExecutable(commandContextForTest(t, t.Context()), f.repo, "widget", "review")
	owner, ownerErr := os.Executable()
	if err != nil || ownerErr != nil || executable != owner {
		t.Fatalf("executable=%q err=%v ownerErr=%v", executable, err, ownerErr)
	}
	want := "roundfix: notice: Delivery Queue item widget: review runs the owner's binary; the item binary's migrate --check was killed by signal SIGKILL (9), no core dump: \n"
	if f.log.String() != want {
		t.Fatalf("log=%q want=%q", f.log.String(), want)
	}
	result, err := f.workflow.runRoundfix(commandContextForTest(t, t.Context()), executable, f.repo, "--version")
	if err != nil || result.exitCode != exitOK || result.stdout == "" {
		t.Fatalf("owner result=%+v err=%v", result, err)
	}
	f.assertStarts(t, []string{"migrate", "--check"})
}

func TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary(t *testing.T) {
	t.Parallel()
	f := newDeliverySignalItemFixture(t, "review")
	for _, tc := range []struct {
		name, stderr, stdout, detail string
	}{
		{name: "empty"},
		{name: "stderr", stderr: "\n reason from stderr\nsecond line\n", stdout: "other reason\n", detail: ": reason from stderr"},
		{name: "stdout", stderr: "\n \n", stdout: "\n reason from stdout\nsecond line\n", detail: ": reason from stdout"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Remove(f.record); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			f.log.Reset()
			setCommandEnvForTest(t, "DELIVERY_ITEM_STDERR", tc.stderr)
			setCommandEnvForTest(t, "DELIVERY_ITEM_STDOUT", tc.stdout)
			_, err := f.workflow.Review(commandContextForTest(t, t.Context()), f.repo, "widget", f.head)
			want := "roundfix review was killed by signal SIGKILL (9), no core dump" + tc.detail
			if err == nil || err.Error() != want {
				t.Fatalf("review error=%v want=%q", err, want)
			}
			f.assertStarts(t, []string{"migrate", "--check"}, []string{"review"})
			if f.log.String() != f.selectionLine("review") {
				t.Fatalf("log=%q", f.log.String())
			}
		})
	}
}
