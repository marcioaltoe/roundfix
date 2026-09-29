// Suite: concurrent pre-PR review finding dispositions.
// Invariant: one finding accepts at most one disposition across Roundfix processes.
// Boundary IN: public review dispose command, Git checks, and Artifact Directory files.
// Boundary OUT: review execution and disposition reuse by later review commands.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const concurrentDisposeProcessCount = 32
const reviewDisposeStartBarrierEnv = "ROUNDFIX_REVIEW_DISPOSE_START_BARRIER"

func init() {
	readyPath := os.Getenv(reviewDisposeStartBarrierEnv)
	if readyPath == "" {
		return
	}
	if err := os.WriteFile(readyPath, nil, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "signal review dispose start barrier: %v\n", err)
		os.Exit(97)
	}
	releasePath := filepath.Join(filepath.Dir(readyPath), "release")
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(releasePath); err == nil {
			return
		} else if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "wait for review dispose start barrier: %v\n", err)
			os.Exit(97)
		}
		if time.Now().After(deadline) {
			fmt.Fprintln(os.Stderr, "wait for review dispose start barrier: deadline exceeded")
			os.Exit(97)
		}
		time.Sleep(time.Millisecond)
	}
}

type reviewDisposeProcessResult struct {
	stdout string
	stderr string
	code   int
}

func TestConcurrentDisposeOfOneFindingAppendsOnce(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)

	barrier := newReviewDisposeStartBarrier(t)
	commands := make([]*exec.Cmd, 0, concurrentDisposeProcessCount)
	stdout := make([]bytes.Buffer, concurrentDisposeProcessCount)
	stderr := make([]bytes.Buffer, concurrentDisposeProcessCount)
	for index := 0; index < concurrentDisposeProcessCount; index++ {
		command := newReviewDisposeProcess(
			t,
			fixture,
			barrier.readyPath(index),
			&stdout[index],
			&stderr[index],
			"F1",
			"--dismiss",
			"--evidence",
			"concurrent evidence",
		)
		if err := command.Start(); err != nil {
			t.Fatalf("start review dispose process %d: %v", index, err)
		}
		commands = append(commands, command)
	}
	barrier.release(t, len(commands))

	results := waitForReviewDisposeProcesses(commands, stdout, stderr)
	successes := 0
	refusals := 0
	const alreadyDisposed = "roundfix: review dispose refused: finding \"F1\" already has a disposition\n"
	for index, result := range results {
		switch result.code {
		case exitOK:
			successes++
			if result.stdout == "" || result.stderr != "" {
				t.Fatalf("successful process %d stdout=%q stderr=%q", index, result.stdout, result.stderr)
			}
		case exitPreflight:
			refusals++
			if result.stdout != "" || result.stderr != alreadyDisposed {
				t.Fatalf("refused process %d stdout=%q stderr=%q", index, result.stdout, result.stderr)
			}
		default:
			t.Fatalf("review dispose process %d exit=%d stdout=%q stderr=%q", index, result.code, result.stdout, result.stderr)
		}
	}
	if successes != 1 || refusals != concurrentDisposeProcessCount-1 {
		t.Fatalf("review dispose results: successes=%d refusals=%d, want 1 and %d", successes, refusals, concurrentDisposeProcessCount-1)
	}

	content, err := os.ReadFile(filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName))
	if err != nil {
		t.Fatalf("read disposition ledger: %v", err)
	}
	if lines := nonBlankDispositionLines(content); len(lines) != 1 {
		t.Fatalf("disposition ledger has %d entries, want 1:\n%s", len(lines), content)
	}
}

func TestConcurrentDisposeOfTwoFindingsBothSucceed(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(
		t,
		fixture,
		"- internal/cli/review.go:10: first\n- internal/cli/review.go:20: second",
		true,
	)

	barrier := newReviewDisposeStartBarrier(t)
	stdout := make([]bytes.Buffer, 2)
	stderr := make([]bytes.Buffer, 2)
	commands := []*exec.Cmd{
		newReviewDisposeProcess(t, fixture, barrier.readyPath(0), &stdout[0], &stderr[0], "F1", "--dismiss", "--evidence", "first evidence"),
		newReviewDisposeProcess(t, fixture, barrier.readyPath(1), &stdout[1], &stderr[1], "F2", "--dismiss", "--evidence", "second evidence"),
	}
	for index, command := range commands {
		if err := command.Start(); err != nil {
			t.Fatalf("start review dispose process %d: %v", index, err)
		}
	}
	barrier.release(t, len(commands))
	for index, result := range waitForReviewDisposeProcesses(commands, stdout, stderr) {
		if result.code != exitOK || result.stdout == "" || result.stderr != "" {
			t.Fatalf("review dispose process %d exit=%d stdout=%q stderr=%q", index, result.code, result.stdout, result.stderr)
		}
	}

	content, err := os.ReadFile(filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName))
	if err != nil {
		t.Fatalf("read disposition ledger: %v", err)
	}
	lines := nonBlankDispositionLines(content)
	if len(lines) != 2 {
		t.Fatalf("disposition ledger has %d entries, want 2:\n%s", len(lines), content)
	}
	seen := make(map[string]bool, 2)
	for _, line := range lines {
		seen[decodeDispositionLine(t, line+"\n").Finding] = true
	}
	if !seen["F1"] || !seen["F2"] {
		t.Fatalf("disposition ledger findings = %v, want F1 and F2", seen)
	}
}

func newReviewDisposeProcess(
	t *testing.T,
	fixture reviewCommandFixture,
	barrierReadyPath string,
	stdout *bytes.Buffer,
	stderr *bytes.Buffer,
	args ...string,
) *exec.Cmd {
	t.Helper()
	commandArgs := append([]string{"review", "dispose"}, args...)
	command := exec.CommandContext(t.Context(), os.Args[0], commandArgs...)
	command.Dir = fixture.repository
	command.Env = cliHelperEnv(t, "", map[string]string{reviewDisposeStartBarrierEnv: barrierReadyPath})
	command.Stdout = stdout
	command.Stderr = stderr
	return command
}

type reviewDisposeStartBarrier struct {
	directory string
}

func newReviewDisposeStartBarrier(t *testing.T) reviewDisposeStartBarrier {
	t.Helper()
	return reviewDisposeStartBarrier{directory: t.TempDir()}
}

func (barrier reviewDisposeStartBarrier) readyPath(index int) string {
	return filepath.Join(barrier.directory, fmt.Sprintf("ready-%d", index))
}

func (barrier reviewDisposeStartBarrier) release(t *testing.T, processCount int) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		entries, err := os.ReadDir(barrier.directory)
		if err != nil {
			t.Fatalf("read review dispose start barrier: %v", err)
		}
		ready := 0
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "ready-") {
				ready++
			}
		}
		if ready == processCount {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("review dispose start barrier has %d ready processes, want %d", ready, processCount)
		}
		time.Sleep(time.Millisecond)
	}
	if err := os.WriteFile(filepath.Join(barrier.directory, "release"), nil, 0o600); err != nil {
		t.Fatalf("release review dispose processes: %v", err)
	}
}

func waitForReviewDisposeProcesses(
	commands []*exec.Cmd,
	stdout []bytes.Buffer,
	stderr []bytes.Buffer,
) []reviewDisposeProcessResult {
	results := make([]reviewDisposeProcessResult, len(commands))
	for index, command := range commands {
		err := command.Wait()
		results[index] = reviewDisposeProcessResult{
			stdout: stdout[index].String(),
			stderr: stderr[index].String(),
			code:   exitCodeFromWait(err),
		}
	}
	return results
}

func nonBlankDispositionLines(content []byte) []string {
	lines := make([]string, 0)
	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
