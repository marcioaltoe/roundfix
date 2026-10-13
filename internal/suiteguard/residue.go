package suiteguard

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

const testBinaryMarkEnv = "SUITEGUARD_TEST_BINARY"
const helperResidueGrace = 10 * time.Second

var errResidueProcessChanging = errors.New("process argument memory changed during read")

type residueProcess struct {
	PID     int
	Command string
}

func markTestBinary() (string, error) {
	var nonce [8]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", fmt.Errorf("create test binary mark: %w", err)
	}
	mark := fmt.Sprintf("%d.%x", os.Getpid(), nonce)
	if err := os.Setenv(testBinaryMarkEnv, mark); err != nil {
		return "", fmt.Errorf("set test binary mark: %w", err)
	}
	return mark, nil
}

func environmentHasMark(environment []byte, mark string) bool {
	if mark == "" {
		return false
	}
	assignment := []byte(testBinaryMarkEnv + "=" + mark)
	for entry := range bytes.SplitSeq(environment, []byte{0}) {
		if bytes.Equal(entry, assignment) {
			return true
		}
	}
	return false
}

type residueReadError struct {
	PID int
	Err error
}

// residueReadErrors retains per-process failures alongside readable residue.
// A process-table failure is returned separately and still makes the scan unavailable.
type residueReadErrors []residueReadError

func (errs residueReadErrors) Error() string {
	return fmt.Sprintf("read process %d environment: %v", errs[0].PID, errs[0].Err)
}

func (errs residueReadErrors) Unwrap() []error {
	causes := make([]error, len(errs))
	for i, err := range errs {
		causes[i] = err.Err
	}
	return causes
}

func checkHelperResidue(mark string, output io.Writer) bool {
	return checkHelperResidueWithScan(mark, output, helperResidue)
}

func checkHelperResidueWithScan(mark string, output io.Writer, scan func(string, int) ([]residueProcess, error)) bool {
	deadline := time.Now().Add(helperResidueGrace)
	for {
		processes, err := scan(mark, os.Getpid())
		if err != nil {
			var unreadable residueReadErrors
			if !errors.As(err, &unreadable) {
				fmt.Fprintf(output, "suiteguard: helper residue check unavailable: %v\n", err)
				return false
			}
			// Argument memory can change before the process table records exit.
			// Retry unreadable processes inside the same grace while retaining
			// readable residue, then report each skip without signalling it.
			if time.Until(deadline) > 0 {
				time.Sleep(min(50*time.Millisecond, time.Until(deadline)))
				continue
			}
			sort.Slice(unreadable, func(i, j int) bool { return unreadable[i].PID < unreadable[j].PID })
			for _, process := range unreadable {
				fmt.Fprintf(output, "suiteguard: helper residue check skipped an unreadable process: pid %d: %v\n", process.PID, process.Err)
			}
		}
		if len(processes) == 0 {
			return false
		}
		if time.Until(deadline) > 0 {
			time.Sleep(min(50*time.Millisecond, time.Until(deadline)))
			continue
		}
		sort.Slice(processes, func(i, j int) bool { return processes[i].PID < processes[j].PID })
		for _, process := range processes {
			fmt.Fprintf(output, "suiteguard: helper residue: pid %d: %s\n", process.PID, process.Command)
			if err := killResidue(process, mark); err != nil {
				fmt.Fprintf(output, "suiteguard: kill helper residue: pid %d: %v\n", process.PID, err)
			}
		}
		return true
	}
}

func residueProcessVanished(err error) bool {
	return errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrProcessDone) || platformProcessVanished(err)
}
