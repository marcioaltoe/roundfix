// Suite: repository boundary guard
// Invariant: package tests leave every non-ignored repository path unchanged.
// Boundary IN: suiteguard.Main, filesystem fingerprints, and repository ignore rules.
// Boundary OUT: installing the guard across packages, owned by task_04.
package suiteguard_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/suiteguard"
	"roundfix/internal/testwait"
)

const (
	fixtureCaseEnv              = "ROUNDFIX_SUITEGUARD_FIXTURE_CASE"
	fixtureModeEnv              = "ROUNDFIX_SUITEGUARD_FIXTURE_MODE"
	fixtureRootEnv              = "ROUNDFIX_SUITEGUARD_FIXTURE_ROOT"
	fixturePathEnv              = "ROUNDFIX_SUITEGUARD_FIXTURE_PATH"
	fixtureSanctionedCommandEnv = "ROUNDFIX_SUITEGUARD_FIXTURE_SANCTIONED_COMMAND"
)

func TestMain(m *testing.M) {
	root := os.Getenv(fixtureRootEnv)
	if root == "" {
		root = filepath.Join("..", "..")
	}
	os.Exit(suiteguard.Main(m, root))
}

func TestGuardNamesTheViolatingPath(t *testing.T) {
	repository := newGitRepository(t)
	mustWrite(t, filepath.Join(repository, "modified.txt"), "before")
	mustWrite(t, filepath.Join(repository, "removed.txt"), "before")

	output, exitCode := runGuardedFixture(t, repository, "violating", "")
	t.Log(output)
	if exitCode == 0 {
		t.Fatalf("guarded package exit code = 0, want failure; output:\n%s", output)
	}
	for _, want := range []string{
		"created: created.txt",
		"modified: modified.txt",
		"removed: removed.txt",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("guarded package output does not contain %q:\n%s", want, output)
		}
	}
}

func TestGuardPassesOnAnIsolatedWrite(t *testing.T) {
	repository := repositoryRoot(t)

	output, exitCode := runGuardedFixture(t, repository, "isolated", "")
	t.Log(output)
	if exitCode != 0 {
		t.Fatalf("guarded package exit code = %d, want success; output:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "guard cost: measured") || !strings.Contains(output, "fingerprint took") {
		t.Fatalf("guarded package did not record its measured cost:\n%s", output)
	}
}

func TestGuardIgnoresRepositoryRules(t *testing.T) {
	repository := newGitRepository(t)
	mustWrite(t, filepath.Join(repository, ".gitignore"), "/ignored/\n")

	output, exitCode := runGuardedFixture(t, repository, "ignored", "")
	t.Log(output)
	if exitCode != 0 {
		t.Fatalf("guarded package exit code = %d, want success; output:\n%s", exitCode, output)
	}
	if strings.Contains(output, "ignored/build.log") {
		t.Fatalf("guarded package reported an ignored path:\n%s", output)
	}
}

func TestSanctionedRegenerationIsDeclaredInProcess(t *testing.T) {
	repository := newGitRepository(t)
	const declaredCommand = "roundfix-fixture-regeneration"
	writeSanctionedAuthorization(t, repository, declaredCommand, "declared.txt")

	output, exitCode := runGuardedFixture(t, repository, "regeneration", declaredCommand, "declared.txt")
	t.Log(output)
	if exitCode != 0 {
		t.Fatalf("guarded package exit code = %d, want success; output:\n%s", exitCode, output)
	}
}

func TestSanctionedRegenerationIsNotAViolationWrongCommandIsRefused(t *testing.T) {
	repository := newGitRepository(t)
	writeSanctionedAuthorization(t, repository, "roundfix-fixture-regeneration", "declared.txt")

	output, exitCode := runGuardedFixture(t, repository, "regeneration", "roundfix-wrong-regeneration", "declared.txt")
	t.Log(output)
	if exitCode == 0 {
		t.Fatalf("guarded package exit code = 0, want failure; output:\n%s", output)
	}
	if !strings.Contains(output, "created: declared.txt") {
		t.Fatalf("guarded package did not name the wrong-command write:\n%s", output)
	}
}

func TestSanctionedRegenerationIsNotAViolationUndeclaredCommandIsRefused(t *testing.T) {
	repository := newGitRepository(t)
	writeSanctionedAuthorization(t, repository, "roundfix-fixture-regeneration", "declared.txt")

	output, exitCode := runGuardedFixture(t, repository, "regeneration", "", "declared.txt")
	t.Log(output)
	if exitCode == 0 {
		t.Fatalf("guarded package exit code = 0, want failure; output:\n%s", output)
	}
	if !strings.Contains(output, "created: declared.txt") {
		t.Fatalf("guarded package did not name the undeclared-command write:\n%s", output)
	}
}

func TestGuardFixture(t *testing.T) {
	mode := os.Getenv(fixtureModeEnv)
	root := os.Getenv(fixtureRootEnv)
	if mode == "" || root == "" {
		t.Skip("runs only in a guarded subprocess")
	}

	switch mode {
	case "helper":
		fmt.Fprintln(os.Stdout, "ready")
		_, _ = io.Copy(io.Discard, os.Stdin)
		os.Exit(0)
	case "residue", "ended":
		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		command := exec.Command(os.Args[0], "-test.run=^TestGuardFixture$")
		command.Env = append(environmentWithout(os.Environ(), fixtureModeEnv), fixtureModeEnv+"=helper")
		command.Stdin = reader
		command.Stderr = os.Stderr
		stdout, err := command.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := command.Start(); err != nil {
			t.Fatal(err)
		}
		_ = reader.Close()
		scanner := bufio.NewScanner(stdout)
		if !scanner.Scan() || scanner.Text() != "ready" {
			t.Fatalf("helper readiness: %q, %v", scanner.Text(), scanner.Err())
		}
		fmt.Fprintf(os.Stdout, "fixture helper pid: %d\n", command.Process.Pid)
		if mode == "ended" {
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if err := command.Wait(); err != nil {
				t.Fatal(err)
			}
		} else {
			// Retain the writer until the guard runs. EOF at process exit cannot
			// rescue a missing guard: its ten-second scan happens before that exit.
			residueFixtureWriter = writer
		}
	case "mark":
		fmt.Fprintf(os.Stdout, "fixture mark: %d:%s\n", os.Getpid(), os.Getenv("SUITEGUARD_TEST_BINARY"))
	case "violating":
		mustWrite(t, filepath.Join(root, "created.txt"), "created")
		mustWrite(t, filepath.Join(root, "modified.txt"), "after")
		if err := os.Remove(filepath.Join(root, "removed.txt")); err != nil {
			t.Fatal(err)
		}
	case "leak":
		mutateLeak(t, root, os.Getenv(fixtureCaseEnv))
	case "isolated":
		mustWrite(t, filepath.Join(t.TempDir(), "isolated.txt"), "isolated")
	case "ignored":
		mustWrite(t, filepath.Join(root, "ignored", "build.log"), "ignored")
	case "regeneration":
		if command := os.Getenv(fixtureSanctionedCommandEnv); command != "" {
			suiteguard.DeclareSanctionedRegeneration(command)
		}
		filePath := os.Getenv(fixturePathEnv)
		if filePath == "" {
			t.Fatal("regeneration fixture path is empty")
		}
		mustWrite(t, filepath.Join(root, filePath), "regenerated")
	default:
		t.Fatalf("unknown fixture mode %q", mode)
	}
}

func runGuardedFixture(t *testing.T, repository, mode, sanctionedCommand string, filePath ...string) (string, int) {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=^TestGuardFixture$", "-test.v=true")
	if len(filePath) > 1 {
		command.Args = append(command.Args, filePath[1:]...)
	}
	command.Env = append(environmentWithout(
		os.Environ(),
		fixtureCaseEnv,
		fixtureModeEnv,
		fixtureRootEnv,
		fixturePathEnv,
		fixtureSanctionedCommandEnv,
	),
		fixtureModeEnv+"="+mode,
		fixtureRootEnv+"="+repository,
	)
	if sanctionedCommand != "" {
		command.Env = append(command.Env, fixtureSanctionedCommandEnv+"="+sanctionedCommand)
	}
	if len(filePath) != 0 {
		command.Env = append(command.Env, fixturePathEnv+"="+filePath[0])
		if mode == "leak" {
			command.Env = append(command.Env, fixtureCaseEnv+"="+filePath[0])
		}
	}
	output, err := command.CombinedOutput()
	if err == nil {
		return string(output), 0
	}
	exitError, ok := errors.AsType[*exec.ExitError](err)
	if !ok {
		t.Fatalf("run guarded fixture: %v", err)
	}
	return string(output), exitError.ExitCode()
}

func writeSanctionedAuthorization(t *testing.T, repository, command string, outputs ...string) {
	t.Helper()
	var declaration strings.Builder
	declaration.WriteString("# Fixture authorization\n\n## Sanctioned regeneration\n\n```yaml\ncommand: ")
	declaration.WriteString(command)
	declaration.WriteString("\noutputs:\n")
	for _, output := range outputs {
		declaration.WriteString("  - ")
		declaration.WriteString(output)
		declaration.WriteByte('\n')
	}
	declaration.WriteString("```\n")
	mustWrite(t, filepath.Join(repository, "docs", "workflow", "authorizations", "fixture.md"), declaration.String())
}

func environmentWithout(environment []string, keys ...string) []string {
	filtered := make([]string, 0, len(environment))
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		keep := true
		for _, excluded := range keys {
			if key == excluded {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func mustWrite(t *testing.T, filePath, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	return root
}

var residueFixtureWriter *os.File

func TestGuardReportsAndKillsHelperResidue(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("process table unavailable on this platform")
	}
	output, exitCode := runGuardedFixture(t, newGitRepository(t), "residue", "")
	t.Log(output)
	if exitCode != 1 {
		t.Fatalf("guard exit = %d, want 1; output:\n%s", exitCode, output)
	}
	match := regexp.MustCompile(`fixture helper pid: ([0-9]+)`).FindStringSubmatch(output)
	if len(match) != 2 {
		t.Fatalf("missing helper PID:\n%s", output)
	}
	pid, err := strconv.Atoi(match[1])
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("suiteguard: helper residue: pid %d: ", pid)
	if strings.Count(output, want) != 1 {
		t.Fatalf("want one residue diagnostic %q:\n%s", want, output)
	}
	testwait.Poll(t, "reported helper to be dead", (<-chan error)(nil), func() (bool, string) {
		// An unreaped orphan is dead too; it cannot retain the lifeline or work.
		command := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "stat=")
		state, err := command.Output()
		if err != nil {
			var exitError *exec.ExitError
			if !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
				t.Fatalf("read helper process state: %v", err)
			}
			return true, "process absent"
		}
		if len(strings.TrimSpace(string(state))) == 0 {
			return true, "process absent"
		}
		return strings.HasPrefix(strings.TrimSpace(string(state)), "Z"), string(state)
	})
}

func TestGuardPassesWhenEveryHelperEnded(t *testing.T) {
	t.Parallel()
	output, exitCode := runGuardedFixture(t, newGitRepository(t), "ended", "")
	if exitCode != 0 || strings.Contains(output, "suiteguard: helper residue:") {
		t.Fatalf("guard exit = %d, output:\n%s", exitCode, output)
	}
	t.Log(output)
}

func TestGuardMarksItsTestBinary(t *testing.T) {
	t.Parallel()
	output, exitCode := runGuardedFixture(t, newGitRepository(t), "mark", "")
	match := regexp.MustCompile(`fixture mark: ([0-9]+):([0-9]+)\.([0-9a-f]{16})\n`).FindStringSubmatch(output)
	if exitCode != 0 || len(match) != 4 || match[1] != match[2] {
		t.Fatalf("invalid binary mark, exit %d:\n%s", exitCode, output)
	}
}
