//go:build unix

package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

type lifelineProcessRecord struct {
	PID      int
	Identity string
}

func TestCLITestHelpersEndWithTheirTestBinary(t *testing.T) {
	t.Parallel()
	const recordEnv = "ROUNDFIX_CLI_LIFELINE_RECORD"
	if record := os.Getenv(recordEnv); record != "" {
		owner, _ := startCLIForceStopOwnerProcess(t)
		pids := []int{owner}
		for _, mode := range []string{detachTestChildRunCreated, detachTestChildIgnoreSentinel} {
			pids = append(pids, startCLILifelineDetachChild(t, mode))
		}

		var records []lifelineProcessRecord
		for _, pid := range pids {
			identity, err := store.OwnerProcessIdentity(t.Context(), pid)
			if err != nil {
				t.Fatalf("record helper %d identity: %v", pid, err)
			}
			records = append(records, lifelineProcessRecord{pid, identity})
		}
		data, err := json.Marshal(records)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(record+".tmp", data, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(record+".tmp", record); err != nil {
			t.Fatal(err)
		}
		testwait.Poll(t, "outer binary to kill this inner binary", (<-chan error)(nil), func() (bool, string) {
			return false, "all helpers recorded and waiting"
		})
		return
	}
	record := filepath.Join(t.TempDir(), "helpers.json")
	log, err := os.Create(filepath.Join(t.TempDir(), "inner.log"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = log.Close()
		if t.Failed() {
			data, err := os.ReadFile(log.Name())
			t.Logf("inner binary log (read error %v):\n%s", err, data)
		}
	})
	command := exec.Command(os.Args[0], "-test.run=^TestCLITestHelpersEndWithTheirTestBinary$", "-test.timeout="+(testwait.Bound(t)+time.Minute).String())
	command.Env = append(os.Environ(), recordEnv+"="+record)
	command.Stdout, command.Stderr = log, log
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	ended := make(chan error, 1)
	go func() { ended <- command.Wait() }()
	reaped := false
	var records []lifelineProcessRecord
	t.Cleanup(func() {
		if !reaped {
			_ = command.Process.Kill()
			_ = testwait.Until(t, "inner binary cleanup", ended, (<-chan error)(nil))
		}
		// Success never signals a helper; only a failed assertion permits this
		// rescue, and each recorded identity is proven again by the signal API.
		if t.Failed() {
			for _, process := range records {
				if _, err := store.SignalProvenProcess(context.Background(), process.PID, process.Identity, syscall.SIGKILL); err != nil {
					t.Errorf("helper cleanup: %v", err)
				}
			}
		}
	})
	testwait.Poll(t, "all helper identity records", ended, func() (bool, string) {
		data, err := os.ReadFile(record)
		if errors.Is(err, os.ErrNotExist) {
			return false, "inner binary starting helpers"
		}
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, &records); err != nil {
			t.Fatal(err)
		}
		if len(records) != 3 {
			t.Fatalf("records = %v, want 3 helper identities", records)
		}
		for _, process := range records {
			identity, err := store.OwnerProcessIdentity(t.Context(), process.PID)
			if err != nil || identity != process.Identity || !store.ProcessAlive(process.PID) {
				t.Fatalf("helper not alive at rendezvous: %v, %v", process, err)
			}
		}
		return true, "every helper alive"
	})
	if err := command.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	err = testwait.Until(t, "SIGKILL inner binary exit", ended, (<-chan error)(nil))
	reaped = true
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("inner exit = %v, want SIGKILL", err)
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() || status.Signal() != syscall.SIGKILL {
		t.Fatalf("inner exit = %v, want SIGKILL", err)
	}
	testwait.Poll(t, "every recorded helper identity to end", (<-chan error)(nil), func() (bool, string) {
		for _, process := range records {
			if !store.ProcessAlive(process.PID) {
				continue
			}
			identity, err := store.OwnerProcessIdentity(t.Context(), process.PID)
			if err != nil {
				if !store.ProcessAlive(process.PID) {
					continue
				}
				t.Fatalf("read still-live helper %d identity: %v", process.PID, err)
			}
			if identity == process.Identity {
				return false, fmt.Sprintf("helper %d still has identity %s", process.PID, identity)
			}
		}
		return true, "all recorded identities gone"
	})
}

func cliHelperLifeline() <-chan struct{} {
	ended := make(chan struct{})
	go func() { _, _ = io.Copy(io.Discard, os.Stdin); close(ended) }()
	return ended
}

func keepCLIHelperLifeline(t *testing.T, command *exec.Cmd) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	command.Stdin = reader
	t.Cleanup(func() { _ = writer.Close(); _ = reader.Close() })
}

func startCLILifelineDetachChild(t *testing.T, mode string) int {
	t.Helper()
	root := t.TempDir()
	command := exec.Command(os.Args[0], "-test.run=^$")
	command.Env = withEnvValue(os.Environ(), detachTestChildModeEnv, mode)
	command.Env = detachTestOwnerEnvironment(t, command.Env)
	command.Env = withEnvValue(command.Env, detachTestPIDPathEnv, filepath.Join(root, "pid"))
	command.Env = withEnvValue(command.Env, detachTestSentinelPathEnv, filepath.Join(root, "release"))
	command.Env = withEnvValue(command.Env, detachConsoleTempEnv, filepath.Join(root, "console.log"))
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = reader.Close(); _ = writer.Close() })
	command.ExtraFiles = []*os.File{writer}
	command.Env = withEnvValue(command.Env, detachHandshakeFDEnv, strconv.Itoa(3))
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	_ = writer.Close()
	t.Cleanup(func() { _ = command.Process.Kill(); _ = command.Wait() })
	waitForFile(t, filepath.Join(root, "pid"), nil)
	if mode == detachTestChildRunCreated {
		buffered := bufio.NewReader(reader)
		marker, err := buffered.ReadByte()
		if err != nil || marker != detachLivenessMarker {
			t.Fatalf("detach liveness = %v, %v", marker, err)
		}
		if _, err := buffered.ReadString('\n'); err != nil {
			t.Fatal(err)
		}
	}
	return command.Process.Pid
}
