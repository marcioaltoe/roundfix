//go:build windows

package store

import (
	"context"
	"fmt"
	"math"
	"time"

	"golang.org/x/sys/windows"
)

func readOwnedProcess(_ context.Context, pid int) (OwnedProcess, error) {
	handle, err := windows.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		return OwnedProcess{}, fmt.Errorf("open process: %w", err)
	}
	defer func() {
		_ = windows.CloseHandle(handle)
	}()

	var created, exited, kernel, user windows.Filetime
	if err := windows.GetProcessTimes(handle, &created, &exited, &kernel, &user); err != nil {
		return OwnedProcess{}, fmt.Errorf("read process times: %w", err)
	}
	kernelTime, err := windowsFiletimeDuration(kernel)
	if err != nil {
		return OwnedProcess{}, fmt.Errorf("read process kernel time: %w", err)
	}
	userTime, err := windowsFiletimeDuration(user)
	if err != nil {
		return OwnedProcess{}, fmt.Errorf("read process user time: %w", err)
	}
	if kernelTime > time.Duration(math.MaxInt64)-userTime {
		return OwnedProcess{}, fmt.Errorf("read process CPU time: duration overflows")
	}
	return OwnedProcess{
		Started: time.Unix(0, created.Nanoseconds()),
		CPUTime: kernelTime + userTime,
	}, nil
}

func windowsFiletimeDuration(value windows.Filetime) (time.Duration, error) {
	ticks := uint64(value.HighDateTime)<<32 | uint64(value.LowDateTime)
	if ticks > uint64(math.MaxInt64/100) {
		return 0, fmt.Errorf("duration overflows")
	}
	return time.Duration(ticks * 100), nil
}
