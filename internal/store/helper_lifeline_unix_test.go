//go:build unix

package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/testwait"
)

type lifelineProcessRecord struct {
	PID      int
	Identity string
}

func TestOwnerProcessHelpersEndWithTheirTestBinary(t *testing.T) {
	t.Parallel()
	const recordEnv = "ROUNDFIX_STORE_LIFELINE_RECORD"
	if record := os.Getenv(recordEnv); record != "" {
		pids := []int{}
		for _, mode := range []string{"graceful", "ignore"} {
			pid, _ := startOwnerProcessHelper(t, mode)
			pids = append(pids, pid)
		}
		root, grandchild, _ := startOwnerProcessTreeHelper(t)
		pids = append(pids, root, grandchild)

		var records []lifelineProcessRecord
		for _, pid := range pids {
			identity, err := OwnerProcessIdentity(t.Context(), pid)
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
	command := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelpersEndWithTheirTestBinary$", "-test.timeout="+(testwait.Bound(t)+time.Minute).String())
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
				if _, err := SignalProvenProcess(context.Background(), process.PID, process.Identity, syscall.SIGKILL); err != nil {
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
		if len(records) != 4 {
			t.Fatalf("records = %v, want 4 helper identities", records)
		}
		for _, process := range records {
			identity, err := OwnerProcessIdentity(t.Context(), process.PID)
			if err != nil || identity != process.Identity || !ProcessAlive(process.PID) {
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
			if !ProcessAlive(process.PID) {
				continue
			}
			identity, err := OwnerProcessIdentity(t.Context(), process.PID)
			if err != nil {
				if !ProcessAlive(process.PID) {
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
