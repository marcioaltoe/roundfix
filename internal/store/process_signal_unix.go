//go:build unix

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// SignalProvenProcess sends only while pid has the start identity recorded
// when it was learned. A different identity proves the recorded process ended.
func SignalProvenProcess(ctx context.Context, pid int, identity string, signal syscall.Signal) (bool, error) {
	if pid <= 0 || pid == os.Getpid() || strings.TrimSpace(identity) == "" || signal == 0 {
		return false, ErrOwnerProcessIdentityUnproven
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	live, err := processStartIdentity(ctx, pid)
	if errors.Is(err, syscall.ESRCH) || errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if live != identity {
		return false, nil
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	err = syscall.Kill(pid, signal)
	if errors.Is(err, syscall.ESRCH) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("signal proven process %d: %w", pid, err)
	}
	return true, nil
}

// SignalProvenProcessGroup proves the leader before signalling its group.
// After the leader is gone, live members reserve the group ID against reuse.
func SignalProvenProcessGroup(ctx context.Context, pgid int, leaderIdentity string, signal syscall.Signal) (bool, error) {
	if pgid <= 0 || pgid == os.Getpid() || strings.TrimSpace(leaderIdentity) == "" || signal == 0 {
		return false, ErrOwnerProcessIdentityUnproven
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	live, err := processStartIdentity(ctx, pgid)
	switch {
	case errors.Is(err, syscall.ESRCH), errors.Is(err, os.ErrNotExist):
		members, readErr := liveProcessGroupMembers(ctx, pgid)
		if readErr != nil {
			return false, readErr
		}
		if len(members) == 0 {
			return false, nil
		}
	case err != nil:
		return false, err
	case live != leaderIdentity:
		return false, nil
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	err = syscall.Kill(-pgid, signal)
	if errors.Is(err, syscall.ESRCH) {
		return false, nil
	}
	if errors.Is(err, syscall.EPERM) {
		members, readErr := liveProcessGroupMembers(ctx, pgid)
		if readErr != nil {
			return false, readErr
		}
		if len(members) == 0 {
			return false, nil
		}
	}
	if err != nil {
		return false, fmt.Errorf("signal proven process group %d: %w", pgid, err)
	}
	return true, nil
}
