//go:build darwin

package store

import (
	"context"
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

func processStartIdentity(_ context.Context, pid int) (string, error) {
	raw, err := unix.SysctlRaw("kern.proc.pid", pid)
	if err != nil {
		return "", fmt.Errorf("read start time for process %d: %w", pid, err)
	}
	if len(raw) == 0 {
		return "", fmt.Errorf("read start time for process %d: %w", pid, unix.ESRCH)
	}
	if len(raw) != unix.SizeofKinfoProc {
		return "", fmt.Errorf("read start time for process %d: sysctl returned %d bytes, want %d", pid, len(raw), unix.SizeofKinfoProc)
	}
	process := *(*unix.KinfoProc)(unsafe.Pointer(&raw[0]))
	started := process.Proc.P_starttime
	return fmt.Sprintf("darwin:%d.%d", started.Sec, started.Usec), nil
}
