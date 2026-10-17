// Boundary: file metadata and recorded starts around a real delivery item probe.
package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"
)

type deliveryProbeReport struct {
	InodeBefore, InodeAfter uint64
	Start                   *deliveryItemStart
	ContextErr              error
}

func (r deliveryProbeReport) String() string {
	form := "replaced"
	switch {
	case r.InodeAfter == 0:
		form = "absent"
	case r.InodeBefore == 0:
		form = "created"
	case r.InodeBefore == r.InodeAfter:
		form = "rewritten in place"
	}
	start := "not recorded"
	if r.Start != nil {
		start = fmt.Sprintf("pid %d ppid %d pgid %d", r.Start.Pid, r.Start.Ppid, r.Start.Pgid)
	}
	state := "live"
	if r.ContextErr != nil {
		state = r.ContextErr.Error()
	}
	return fmt.Sprintf("probe report: item binary inode %d -> %d (%s); probe start %s; context %s", r.InodeBefore, r.InodeAfter, form, start, state)
}

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

func TestDeliveryProbeReportFormsEachEnd(t *testing.T) {
	t.Parallel()
	for _, form := range []struct {
		name          string
		before, after uint64
		want          string
	}{
		{"absent", 41, 0, "probe report: item binary inode 41 -> 0 (absent)"},
		{"created", 0, 42, "probe report: item binary inode 0 -> 42 (created)"},
		{"rewritten in place", 41, 41, "probe report: item binary inode 41 -> 41 (rewritten in place)"},
		{"replaced", 41, 42, "probe report: item binary inode 41 -> 42 (replaced)"},
	} {
		t.Run(form.name, func(t *testing.T) {
			for _, start := range []struct {
				name  string
				value *deliveryItemStart
				want  string
			}{
				{"recorded", &deliveryItemStart{Pid: 101, Ppid: 102, Pgid: 103}, "pid 101 ppid 102 pgid 103"},
				{"unrecorded", nil, "not recorded"},
			} {
				t.Run(start.name, func(t *testing.T) {
					for _, state := range []struct {
						name string
						err  error
						want string
					}{
						{"live", nil, "live"},
						{"canceled", context.Canceled, "context canceled"},
					} {
						t.Run(state.name, func(t *testing.T) {
							report := deliveryProbeReport{InodeBefore: form.before, InodeAfter: form.after, Start: start.value, ContextErr: state.err}
							want := form.want + "; probe start " + start.want + "; context " + state.want
							if got := report.String(); got != want {
								t.Fatalf("report=%q want=%q", got, want)
							}
						})
					}
				})
			}
		})
	}
}

func TestDeliveryProbeReportRecordsTheProbeStart(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	executable, err, report := f.step(t, commandContextForTest(t, t.Context()), "review")
	path := filepath.Join(f.repo, "bin/roundfix")
	if err != nil || executable != path {
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
