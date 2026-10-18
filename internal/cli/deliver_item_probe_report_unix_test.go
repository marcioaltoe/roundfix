//go:build unix

// Suite: delivery probe file metadata and Unix process starts.
// Invariant: probe reports retain file identity and the observed process start.
// Boundary IN: real item builds, executable metadata and probe processes in temporary repositories.
// Boundary OUT: delivery queue settlement and Windows filesystem metadata.
package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
)

func deliveryItemInode(t *testing.T, path string) uint64 {
	t.Helper()
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return info.Sys().(*syscall.Stat_t).Ino
}

func (f *deliveryItemFixture) step(t *testing.T, ctx context.Context, step string) (string, error, deliveryProbeReport) {
	t.Helper()
	path := filepath.Join(f.repo, "bin/roundfix")
	startCount := len(f.starts(t))
	report := deliveryProbeReport{InodeBefore: deliveryItemInode(t, path)}
	executable, err := f.workflow.stepExecutable(ctx, f.repo, "widget", step)
	report.ContextErr = ctx.Err()
	report.InodeAfter = deliveryItemInode(t, path)
	starts := f.starts(t)
	for i := startCount; i < len(starts); i++ {
		if slices.Equal(starts[i].Args, []string{"migrate", "--check"}) {
			report.Start = &starts[i]
		}
	}
	return executable, err, report
}

func TestDeliveryProbeReportRecordsTheProbeStart(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	executable, err, report := f.step(t, commandContextForTest(t, t.Context()), "review")
	path := filepath.Join(f.repo, "bin/roundfix")
	if err != nil || executable != f.copyPath(t) {
		t.Fatalf("executable=%q err=%v; %s", executable, err, report)
	}
	inode := deliveryItemInode(t, path)
	if inode == 0 || report.InodeBefore != 0 || report.InodeAfter != inode {
		t.Fatalf("inode=%d; %s", inode, report)
	}
	if report.Start == nil || report.Start.Pid <= 0 || report.Start.Ppid != os.Getpid() || report.Start.Pgid != syscall.Getpgrp() {
		t.Fatalf("want a probe parent %d and group %d; %s", os.Getpid(), syscall.Getpgrp(), report)
	}
	want := fmt.Sprintf("probe report: item binary inode 0 -> %d (created); probe start pid %d ppid %d pgid %d; context live", inode, report.Start.Pid, os.Getpid(), syscall.Getpgrp())
	if report.ContextErr != nil || report.String() != want {
		t.Fatalf("report=%q want=%q", report, want)
	}
	f.assertStarts(t, []string{"migrate", "--check"})

	ctx, cancel := context.WithCancel(commandContextForTest(t, t.Context()))
	cancel()
	_, err, report = f.step(t, ctx, "review")
	want = fmt.Sprintf("probe report: item binary inode %d -> %d (rewritten in place); probe start not recorded; context context canceled", inode, inode)
	if err == nil || report.Start != nil || report.ContextErr != context.Canceled || report.String() != want {
		t.Fatalf("step error=%v report=%q want=%q", err, report, want)
	}
	f.assertStarts(t, []string{"migrate", "--check"})
}

func TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	setCommandEnvForTest(t, "DELIVERY_ITEM_PROBE_EXIT", "2")
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	for _, stream := range []string{"stderr", "stdout"} {
		t.Run(stream, func(t *testing.T) {
			if err := os.Remove(f.record); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			f.log.Reset()
			setCommandEnvForTest(t, "DELIVERY_ITEM_STDERR", "\n  \n")
			setCommandEnvForTest(t, "DELIVERY_ITEM_STDOUT", "\nstdout reason\nlater\n")
			detail := "stdout reason"
			if stream == "stderr" {
				setCommandEnvForTest(t, "DELIVERY_ITEM_STDERR", "\nprobe schema differs\nsecond line\n")
				detail = "probe schema differs"
			}
			executable, err, report := f.step(t, commandContextForTest(t, t.Context()), "review")
			owner, ownerErr := os.Executable()
			if err != nil || ownerErr != nil || executable != owner {
				t.Fatalf("executable=%q err=%v ownerErr=%v", executable, err, ownerErr)
			}
			result, err := f.workflow.runRoundfix(commandContextForTest(t, t.Context()), executable, f.repo, "--version")
			if err != nil || result.exitCode != 0 || result.stdout == "" {
				t.Fatalf("owner result=%+v err=%v", result, err)
			}
			want := "roundfix: notice: Delivery Queue item widget: review runs the owner's binary; the item binary's migrate --check exited 2: " + detail + "\n"
			if f.log.String() != want {
				t.Fatalf("log=%q want=%q; %s", f.log.String(), want, report)
			}
			f.assertStarts(t, []string{"migrate", "--check"})
		})
	}
}

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
