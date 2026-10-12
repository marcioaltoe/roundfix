//go:build linux

package suiteguard

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func platformProcessVanished(err error) bool { return errors.Is(err, unix.ESRCH) }

func linuxProcessEnvironment(pid int) ([]byte, error) {
	root := fmt.Sprintf("/proc/%d", pid)
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("process %d has no uid", pid)
	}
	if stat.Uid != uint32(os.Geteuid()) {
		return nil, nil
	}
	return os.ReadFile(root + "/environ")
}

func helperResidue(mark string, self int) ([]residueProcess, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read process table: %w", err)
	}
	var residue []residueProcess
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid <= 0 || pid == self {
			continue
		}
		environment, err := linuxProcessEnvironment(pid)
		if residueProcessVanished(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read process %d environment: %w", pid, err)
		}
		if !environmentHasMark(environment, mark) {
			continue
		}
		command, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
		if residueProcessVanished(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read process %d command: %w", pid, err)
		}
		// Zombies expose neither a live command nor a usable environment.
		if len(command) != 0 {
			residue = append(residue, residueProcess{pid, strings.TrimSpace(string(bytes.ReplaceAll(command, []byte{0}, []byte{' '})))})
		}
	}
	return residue, nil
}

func killResidue(process residueProcess, mark string) error {
	if process.PID <= 0 || process.PID == os.Getpid() || mark == "" {
		return nil
	}
	environment, err := linuxProcessEnvironment(process.PID)
	if residueProcessVanished(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reread process mark: %w", err)
	}
	if !environmentHasMark(environment, mark) {
		return nil
	}
	if err := unix.Kill(process.PID, unix.SIGKILL); err != nil && !residueProcessVanished(err) {
		return err
	}
	return nil
}
