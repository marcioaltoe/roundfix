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

func checkHelperResidue(mark string, output io.Writer) bool {
	deadline := time.Now().Add(helperResidueGrace)
	for {
		processes, err := helperResidue(mark, os.Getpid())
		if err != nil {
			// Argument memory can disappear before the process table records exit.
			// Retry that transition inside the same grace, never extend the grace
			// and never silently omit an unreadable live process from a scan.
			if errors.Is(err, errResidueProcessChanging) && time.Until(deadline) > 0 {
				time.Sleep(min(50*time.Millisecond, time.Until(deadline)))
				continue
			}
			fmt.Fprintf(output, "suiteguard: helper residue check unavailable: %v\n", err)
			return false
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
