// Suite: repository boundary guard
// Invariant: package tests leave every non-ignored repository path unchanged.
// Boundary IN: suiteguard.Main, filesystem fingerprints, and repository ignore rules.
// Boundary OUT: installing the guard across packages, owned by task_04.
package suiteguard_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/suiteguard"
)

const (
	fixtureModeEnv = "ROUNDFIX_SUITEGUARD_FIXTURE_MODE"
	fixtureRootEnv = "ROUNDFIX_SUITEGUARD_FIXTURE_ROOT"
	fixturePathEnv = "ROUNDFIX_SUITEGUARD_FIXTURE_PATH"
)

func TestMain(m *testing.M) {
	root := os.Getenv(fixtureRootEnv)
	if root == "" {
		root = filepath.Join("..", "..")
	}
	os.Exit(suiteguard.Main(m, root))
}

func TestGuardNamesTheViolatingPath(t *testing.T) {
	repository := t.TempDir()
	mustWrite(t, filepath.Join(repository, "modified.txt"), "before")
	mustWrite(t, filepath.Join(repository, "removed.txt"), "before")

	output, exitCode := runGuardedFixture(t, repository, "violating")
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

	output, exitCode := runGuardedFixture(t, repository, "isolated")
	t.Log(output)
	if exitCode != 0 {
		t.Fatalf("guarded package exit code = %d, want success; output:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "guard cost: measured") || !strings.Contains(output, "fingerprint took") {
		t.Fatalf("guarded package did not record its measured cost:\n%s", output)
	}
}

func TestGuardIgnoresRepositoryRules(t *testing.T) {
	repository := t.TempDir()
	mustWrite(t, filepath.Join(repository, ".gitignore"), "/ignored/\n")

	output, exitCode := runGuardedFixture(t, repository, "ignored")
	t.Log(output)
	if exitCode != 0 {
		t.Fatalf("guarded package exit code = %d, want success; output:\n%s", exitCode, output)
	}
	if strings.Contains(output, "ignored/build.log") {
		t.Fatalf("guarded package reported an ignored path:\n%s", output)
	}
}

func TestSanctionedRegenerationIsNotAViolation(t *testing.T) {
	repository := t.TempDir()
	writeSanctionedAuthorization(t, repository, guardedFixtureCommand(), "declared.txt")

	output, exitCode := runGuardedFixture(t, repository, "regeneration", "declared.txt")
	t.Log(output)
	if exitCode != 0 {
		t.Fatalf("guarded package exit code = %d, want success; output:\n%s", exitCode, output)
	}
}

func TestSanctionedRegenerationIsNotAViolationWrongCommandIsRefused(t *testing.T) {
	repository := t.TempDir()
	writeSanctionedAuthorization(t, repository, "roundfix-not-the-running-command", "declared.txt")

	output, exitCode := runGuardedFixture(t, repository, "regeneration", "declared.txt")
	t.Log(output)
	if exitCode == 0 {
		t.Fatalf("guarded package exit code = 0, want failure; output:\n%s", output)
	}
	if !strings.Contains(output, "created: declared.txt") {
		t.Fatalf("guarded package did not name the wrong-command write:\n%s", output)
	}
}

func TestSanctionedRegenerationIsNotAViolationUndeclaredPathIsRefused(t *testing.T) {
	repository := t.TempDir()
	writeSanctionedAuthorization(t, repository, guardedFixtureCommand(), "declared.txt")

	output, exitCode := runGuardedFixture(t, repository, "regeneration", "undeclared.txt")
	t.Log(output)
	if exitCode == 0 {
		t.Fatalf("guarded package exit code = 0, want failure; output:\n%s", output)
	}
	if !strings.Contains(output, "created: undeclared.txt") {
		t.Fatalf("guarded package did not name the undeclared write:\n%s", output)
	}
}

func TestGuardFixture(t *testing.T) {
	mode := os.Getenv(fixtureModeEnv)
	root := os.Getenv(fixtureRootEnv)
	if mode == "" || root == "" {
		t.Skip("runs only in a guarded subprocess")
	}

	switch mode {
	case "violating":
		mustWrite(t, filepath.Join(root, "created.txt"), "created")
		mustWrite(t, filepath.Join(root, "modified.txt"), "after")
		if err := os.Remove(filepath.Join(root, "removed.txt")); err != nil {
			t.Fatal(err)
		}
	case "isolated":
		mustWrite(t, filepath.Join(t.TempDir(), "isolated.txt"), "isolated")
	case "ignored":
		mustWrite(t, filepath.Join(root, "ignored", "build.log"), "ignored")
	case "regeneration":
		filePath := os.Getenv(fixturePathEnv)
		if filePath == "" {
			t.Fatal("regeneration fixture path is empty")
		}
		mustWrite(t, filepath.Join(root, filePath), "regenerated")
	default:
		t.Fatalf("unknown fixture mode %q", mode)
	}
}

func runGuardedFixture(t *testing.T, repository, mode string, filePath ...string) (string, int) {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=^TestGuardFixture$", "-test.v=true")
	command.Env = append(environmentWithout(os.Environ(), fixtureModeEnv, fixtureRootEnv, fixturePathEnv),
		fixtureModeEnv+"="+mode,
		fixtureRootEnv+"="+repository,
	)
	if len(filePath) != 0 {
		command.Env = append(command.Env, fixturePathEnv+"="+filePath[0])
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

func guardedFixtureCommand() string {
	return strings.Join([]string{
		filepath.Base(os.Args[0]),
		"-test.run=^TestGuardFixture$",
		"-test.v=true",
	}, " ")
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
