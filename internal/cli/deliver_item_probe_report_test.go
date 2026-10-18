// Boundary: file metadata and recorded starts around a real delivery item probe.
package cli

import (
	"context"
	"fmt"
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
