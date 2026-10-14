//go:build !darwin && !linux

package childexit

import "syscall"

func signalName(signal syscall.Signal) string {
	return signal.String()
}
