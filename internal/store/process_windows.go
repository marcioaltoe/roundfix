//go:build windows

package store

import "syscall"

const processQueryLimitedInformation = 0x1000

// ProcessAlive reports whether pid can be opened by the current process.
func ProcessAlive(pid int) bool {
	if pid <= 0 {
		return true
	}
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		if err == syscall.ERROR_ACCESS_DENIED {
			return true
		}
		return false
	}
	_ = syscall.CloseHandle(handle)
	return true
}
