//go:build darwin || linux

package suiteguard

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

const residueFixtureEnv = "SUITEGUARD_RESIDUE_FIXTURE"

func TestResidueProcessHelper(t *testing.T) {
	t.Parallel()
	if os.Getenv(residueFixtureEnv) != "1" {
		return
	}
	os.Stdout.WriteString("ready\n")
	if _, err := io.Copy(io.Discard, os.Stdin); err != nil {
		t.Fatal(err)
	}
}

func TestKillResidueRereadsTheExactMark(t *testing.T) {
	t.Parallel()
	mark := os.Getenv(testBinaryMarkEnv)
	if mark == "" {
		t.Fatal("test binary is unmarked")
	}
	command := exec.Command(os.Args[0], "-test.run=^TestResidueProcessHelper$")
	command.Env = append(os.Environ(), residueFixtureEnv+"=1")
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	command.Stdin = reader
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	_ = reader.Close()
	waited := false
	t.Cleanup(func() {
		_ = writer.Close()
		if !waited {
			_ = command.Process.Kill()
			_ = command.Wait()
		}
	})
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() || scanner.Text() != "ready" {
		t.Fatalf("helper readiness = %q, %v", scanner.Text(), scanner.Err())
	}
	process := residueProcess{PID: command.Process.Pid}
	if err := killResidue(process, mark+"different"); err != nil {
		t.Fatal(err)
	}
	if err := command.Process.Signal(syscall.Signal(0)); err != nil {
		t.Fatalf("mismatched mark signaled helper: %v", err)
	}
	if err := killResidue(process, mark); err != nil {
		t.Fatal(err)
	}
	err = command.Wait()
	waited = true
	if err == nil {
		t.Fatal("matching marked helper exited without SIGKILL")
	}
	status, ok := command.ProcessState.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() || status.Signal() != syscall.SIGKILL {
		t.Fatalf("helper status = %v, want SIGKILL", command.ProcessState)
	}
}
