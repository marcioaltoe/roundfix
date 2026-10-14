//go:build darwin || linux

package childexit

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func signalName(signal syscall.Signal) string {
	return unix.SignalName(signal)
}
