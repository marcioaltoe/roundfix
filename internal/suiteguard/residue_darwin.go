//go:build darwin

package suiteguard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

func platformProcessVanished(err error) bool { return errors.Is(err, unix.ESRCH) }

// kern.procargs2 holds argc, the executable, padding, argc arguments, then
// NUL-separated environment entries. Arguments must never count as marks.
func darwinProcessEnvironment(pid int) ([]byte, string, error) {
	raw, err := unix.SysctlRaw("kern.procargs2", pid)
	if err != nil {
		return nil, "", err
	}
	return parseDarwinProcessArguments(raw, pid)
}

func parseDarwinProcessArguments(raw []byte, pid int) ([]byte, string, error) {
	if len(raw) < 4 {
		return nil, "", fmt.Errorf("short process arguments for pid %d", pid)
	}
	argc := int(binary.NativeEndian.Uint32(raw[:4]))
	raw = raw[4:]
	end := bytes.IndexByte(raw, 0)
	if end < 0 {
		return nil, "", fmt.Errorf("missing executable for pid %d", pid)
	}
	executable := string(raw[:end])
	raw = bytes.TrimLeft(raw[end+1:], "\x00")
	var arguments []string
	for range argc {
		end = bytes.IndexByte(raw, 0)
		if end < 0 {
			return nil, "", fmt.Errorf("short argument list for pid %d", pid)
		}
		arguments = append(arguments, string(raw[:end]))
		raw = raw[end+1:]
	}
	command := strings.Join(arguments, " ")
	if command == "" {
		command = executable
	}
	return raw, command, nil
}

func helperResidue(mark string, self int) ([]residueProcess, error) {
	processes, err := unix.SysctlKinfoProcSlice("kern.proc.all")
	if err != nil {
		return nil, fmt.Errorf("read process table: %w", err)
	}
	var residue []residueProcess
	for _, process := range processes {
		pid := int(process.Proc.P_pid)
		if pid <= 0 || pid == self || process.Eproc.Ucred.Uid != uint32(os.Geteuid()) || process.Proc.P_stat == 5 || process.Proc.P_flag&0x2000 != 0 {
			continue
		}
		environment, command, err := darwinProcessEnvironment(pid)
		if err != nil {
			// An exiting process can lose procargs before its table entry disappears.
			absent, readErr := darwinResidueProcessAbsent(pid)
			if residueProcessVanished(err) || (readErr == nil && absent) {
				continue
			}
			if readErr == nil && (errors.Is(err, unix.EIO) || errors.Is(err, unix.EINVAL)) {
				// XNU's sysctl_procargsx may lose the referenced argument memory to
				// exit/exec before kern.proc.pid reflects that transition.
				return nil, fmt.Errorf("%w: read process %d environment: %w", errResidueProcessChanging, pid, err)
			}
			return nil, fmt.Errorf("read process %d environment: %w", pid, err)
		}
		if environmentHasMark(environment, mark) {
			residue = append(residue, residueProcess{pid, command})
		}
	}
	return residue, nil
}

func killResidue(process residueProcess, mark string) error {
	if process.PID <= 0 || process.PID == os.Getpid() || mark == "" {
		return nil
	}
	environment, _, err := darwinProcessEnvironment(process.PID)
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

// The typed SysctlKinfoProc reports EIO when an already-reaped PID returns
// zero bytes. Read the raw response to distinguish that absence from a real
// process-table failure; never classify an arbitrary EIO as an exit.
func darwinResidueProcessAbsent(pid int) (bool, error) {
	raw, err := unix.SysctlRaw("kern.proc.pid", pid)
	if residueProcessVanished(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	if len(raw) == 0 {
		return true, nil
	}
	if len(raw) != unix.SizeofKinfoProc {
		return false, fmt.Errorf("process %d returned %d status bytes", pid, len(raw))
	}
	process := *(*unix.KinfoProc)(unsafe.Pointer(&raw[0]))
	return process.Proc.P_stat == 5 || process.Proc.P_flag&0x2000 != 0, nil
}
