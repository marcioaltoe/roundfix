//go:build unix

package store

import "syscall"

// ProcessAlive reports whether pid provably exists. Only ESRCH proves death.
func ProcessAlive(pid int) bool {
	if pid <= 0 {
		return true
	}
	err := syscall.Kill(pid, syscall.Signal(0))
	return err != syscall.ESRCH
}
