package suiteguard

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"roundfix/internal/testwait"
)

func TestResidueMarkRequiresTheExactEnvironmentAssignment(t *testing.T) {
	t.Parallel()
	const mark = "123.0123456789abcdef"
	for _, tc := range []struct {
		name, environment string
		want              bool
	}{
		{"exact", testBinaryMarkEnv + "=" + mark + "\x00OTHER=value\x00", true},
		{"another mark", testBinaryMarkEnv + "=" + mark + "0\x00", false},
		{"another key", "OTHER_" + testBinaryMarkEnv + "=" + mark + "\x00", false},
		{"value containing assignment", "OTHER=" + testBinaryMarkEnv + "=" + mark + "\x00", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := environmentHasMark([]byte(tc.environment), mark); got != tc.want {
				t.Fatalf("mark match = %v, want %v", got, tc.want)
			}
		})
	}
	if environmentHasMark([]byte(testBinaryMarkEnv+"=\x00"), "") {
		t.Fatal("empty mark matched")
	}
}

func TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("process environment unavailable on this platform")
	}
	mark := t.TempDir()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	command := exec.Command(os.Args[0], "-test.run=^TestGuardFixture$")
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case testBinaryMarkEnv, "ROUNDFIX_SUITEGUARD_FIXTURE_MODE", "ROUNDFIX_SUITEGUARD_FIXTURE_ROOT":
			continue
		}
		command.Env = append(command.Env, entry)
	}
	repository := t.TempDir()
	git := exec.CommandContext(t.Context(), "git", "-C", repository, "init")
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(entry), "GIT_") {
			git.Env = append(git.Env, entry)
		}
	}
	git.Env = append(git.Env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull)
	if output, err := git.CombinedOutput(); err != nil {
		t.Fatalf("initialize helper repository: %v: %s", err, output)
	}
	command.Env = append(command.Env, testBinaryMarkEnv+"="+mark,
		"ROUNDFIX_SUITEGUARD_FIXTURE_MODE=helper",
		"ROUNDFIX_SUITEGUARD_FIXTURE_ROOT="+repository)
	command.Stdin = reader
	var helperErrors bytes.Buffer
	command.Stderr = &helperErrors
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	ended := make(chan error, 1)
	go func() {
		ended <- command.Wait()
		close(ended)
	}()
	waited := false
	ready := make(chan string, 1)
	readDone := make(chan struct{})
	t.Cleanup(func() {
		// EOF is the helper's lifeline, including when an assertion fails.
		if err := writer.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			t.Error(err)
		}
		if !waited {
			testwait.Until(t, "residue fixture cleanup", ended, (<-chan error)(nil))
		}
		testwait.Until(t, "helper readiness reader cleanup", readDone, (<-chan error)(nil))
	})
	go func() {
		defer close(readDone)
		scanner := bufio.NewScanner(stdout)
		scanner.Scan()
		ready <- scanner.Text()
	}()
	if got := testwait.Until(t, "residue helper readiness", ready, ended); got != "ready" {
		t.Fatalf("helper readiness = %q, want ready", got)
	}

	unreadablePID := os.Getpid()
	unreadable := fmt.Errorf("%w: short argument list for pid %d", errResidueProcessChanging, unreadablePID)
	reads := 0
	scan := func(gotMark string, self int) ([]residueProcess, error) {
		if gotMark != mark || self != os.Getpid() {
			t.Fatalf("scan identity = %q, %d", gotMark, self)
		}
		reads++
		return []residueProcess{{PID: command.Process.Pid, Command: "marked fixture"}},
			residueReadErrors{{PID: unreadablePID, Err: unreadable}}
	}
	var output bytes.Buffer
	found := checkHelperResidueWithScan(mark, &output, scan)
	t.Log(output.String())
	if !found {
		t.Errorf("residue check returned false beside an unreadable process")
	}
	if reads < 2 {
		t.Errorf("process scanned %d times, want retries inside the grace", reads)
	}
	for _, want := range []string{
		fmt.Sprintf("suiteguard: helper residue check skipped an unreadable process: pid %d: %v\n", unreadablePID, unreadable),
		fmt.Sprintf("suiteguard: helper residue: pid %d: marked fixture\n", command.Process.Pid),
	} {
		if strings.Count(output.String(), want) != 1 {
			t.Errorf("want one diagnostic %q; output:\n%s", want, output.String())
		}
	}
	if strings.Contains(output.String(), "check unavailable") ||
		strings.Contains(output.String(), fmt.Sprintf("suiteguard: helper residue: pid %d:", unreadablePID)) {
		t.Errorf("unreadable process misclassified; output:\n%s", output.String())
	}
	if t.Failed() {
		return
	}
	err = testwait.Until(t, "marked helper death", ended, (<-chan error)(nil))
	waited = true
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("helper death = %v, want killed process; stderr:\n%s", err, helperErrors.String())
	}
	state, ok := exitError.ProcessState.Sys().(syscall.WaitStatus)
	if !ok || !state.Signaled() || state.Signal() != syscall.SIGKILL {
		t.Fatalf("helper exit = %v, want SIGKILL", exitError.ProcessState)
	}
}

func TestHelperResidueRetriesAnUnreadableProcessThatRecovers(t *testing.T) {
	t.Parallel()
	reads := 0
	var output bytes.Buffer
	found := checkHelperResidueWithScan("fixture", &output, func(_ string, self int) ([]residueProcess, error) {
		reads++
		if reads == 1 {
			return nil, residueReadErrors{{PID: self, Err: errResidueProcessChanging}}
		}
		return nil, nil
	})
	if found || reads != 2 || output.Len() != 0 {
		t.Fatalf("recovering scan: residue = %v, reads = %d, output = %q", found, reads, output.String())
	}
}

func TestHelperResidueKeepsAProcessTableFailureUnavailable(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	reads := 0
	found := checkHelperResidueWithScan("fixture", &output, func(string, int) ([]residueProcess, error) {
		reads++
		return nil, errors.New("read process table: fixture unavailable")
	})
	const want = "suiteguard: helper residue check unavailable: read process table: fixture unavailable\n"
	if found || reads != 1 || output.String() != want {
		t.Fatalf("unavailable scan: residue = %v, reads = %d, output = %q", found, reads, output.String())
	}
}

func TestHelperResidueSkipsUnreadableProcessesWithoutResidue(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	reads := 0
	found := checkHelperResidueWithScan("fixture", &output, func(string, int) ([]residueProcess, error) {
		reads++
		return nil, residueReadErrors{
			{PID: 2, Err: errors.New("second unreadable process")},
			{PID: 1, Err: errResidueProcessChanging},
		}
	})
	want := fmt.Sprintf("suiteguard: helper residue check skipped an unreadable process: pid 1: %v\n", errResidueProcessChanging) +
		"suiteguard: helper residue check skipped an unreadable process: pid 2: second unreadable process\n"
	if found || reads < 2 || output.String() != want {
		t.Fatalf("unreadable scan: residue = %v, reads = %d, output = %q", found, reads, output.String())
	}
}
