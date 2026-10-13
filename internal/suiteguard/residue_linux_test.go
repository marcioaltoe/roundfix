//go:build linux

package suiteguard

import (
	"fmt"
	"io/fs"
	"syscall"
	"testing"
)

func TestHelperResidueTreatsAnUninspectableProcessAsUnmarked(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		err  error
		want bool
	}{
		{&fs.PathError{Op: "open", Path: "/proc/1002/environ", Err: syscall.EACCES}, true},
		{&fs.PathError{Op: "open", Path: "/proc/1002/environ", Err: syscall.EPERM}, true},
		{fmt.Errorf("read: %w", syscall.EACCES), true},
		{&fs.PathError{Op: "open", Path: "/proc/1002/environ", Err: syscall.ENOENT}, false},
		{&fs.PathError{Op: "read", Path: "/proc/1002/environ", Err: syscall.EIO}, false},
	} {
		if got := linuxProcessUninspectable(tc.err); got != tc.want {
			t.Errorf("linuxProcessUninspectable(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}
