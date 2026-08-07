// Suite: authoritative repository gate
// Invariant: the Makefile test target returns the Go toolchain's exit status.
// Boundary IN: Makefile expansion, the Go test subprocess, and a throwaway module.
// Boundary OUT: make itself, this repository's test result, and RTK internals.

package spec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const gateFixtureModule = "example.com/roundfix-gate-fixture"

func TestAuthoritativeGateReportsFailure(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	makefileBytes, err := os.ReadFile(filepath.Join(repositoryRoot, "Makefile"))
	if err != nil {
		t.Fatalf("read repository Makefile: %v", err)
	}
	makefile := string(makefileBytes)
	maskingWrapper := writeMaskingWrapper(t)
	goCache := filepath.Join(t.TempDir(), "gocache")

	tests := []struct {
		name               string
		testSource         string
		wantFailure        bool
		wantFailingPackage bool
	}{
		{
			name: "high-volume failure",
			testSource: `package gatefixture

import "testing"

func TestFailsWithHighOutput(t *testing.T) {
	for line := range 300 {
		t.Logf("concealment output line %03d", line)
	}
	t.Error("deliberate high-volume failure")
}
`,
			wantFailure:        true,
			wantFailingPackage: true,
		},
		{
			name: "short failure",
			testSource: `package gatefixture

import "testing"

func TestFailsWithShortOutput(t *testing.T) {
	t.Fatal("deliberate short failure")
}
`,
			wantFailure:        true,
			wantFailingPackage: true,
		},
		{
			name: "passing package",
			testSource: `package gatefixture

import "testing"

func TestPasses(t *testing.T) {}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moduleRoot := writeGateFixtureModule(t, tt.testSource)
			output, runErr := runMakefileTestTarget(t, makefile, maskingWrapper, goCache, moduleRoot)
			if tt.wantFailure && runErr == nil {
				t.Fatalf("authoritative test target exited zero for a failing package:\n%s", output)
			}
			if !tt.wantFailure && runErr != nil {
				t.Fatalf("authoritative test target failed for a passing package: %v\n%s", runErr, output)
			}
			if tt.wantFailingPackage && !strings.Contains(output, gateFixtureModule) {
				t.Fatalf("authoritative test target did not name failing package %q:\n%s", gateFixtureModule, output)
			}
		})
	}

	t.Run("wrapper routing masks the failure signal", func(t *testing.T) {
		const authoritativeAssignment = "GO := go"
		if !strings.Contains(makefile, authoritativeAssignment) {
			t.Fatalf("Makefile does not contain authoritative assignment %q", authoritativeAssignment)
		}
		wrappedMakefile := strings.Replace(
			makefile,
			authoritativeAssignment,
			"GO := $(RTK) go",
			1,
		)
		moduleRoot := writeGateFixtureModule(t, tests[0].testSource)
		output, runErr := runMakefileTestTarget(t, wrappedMakefile, maskingWrapper, goCache, moduleRoot)
		if runErr != nil {
			t.Fatalf("synthetic output wrapper did not mask the failing toolchain status: %v\n%s", runErr, output)
		}
	})
}

func runMakefileTestTarget(t *testing.T, makefile, maskingWrapper, goCache, moduleRoot string) (string, error) {
	t.Helper()

	variables := makeVariables(makefile)
	variables["RTK"] = maskingWrapper
	recipe := makeTargetRecipe(t, makefile, "test")
	commandLine := expandMakeVariables(recipe, variables)
	if strings.Contains(commandLine, "$(") {
		t.Fatalf("test target contains an unresolved Make variable: %q", commandLine)
	}
	command := strings.Fields(commandLine)
	if len(command) == 0 {
		t.Fatal("test target expanded to an empty command")
	}

	cmd := exec.CommandContext(t.Context(), command[0], command[1:]...)
	cmd.Dir = moduleRoot
	cmd.Env = append(
		os.Environ(),
		"GOCACHE="+goCache,
		"GOWORK=off",
	)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func makeVariables(makefile string) map[string]string {
	variables := make(map[string]string)
	for _, line := range strings.Split(makefile, "\n") {
		if strings.HasPrefix(line, "\t") {
			continue
		}
		for _, operator := range []string{":=", "?=", "="} {
			before, after, found := strings.Cut(line, operator)
			if !found {
				continue
			}
			name := strings.TrimSpace(before)
			if name != "" && !strings.ContainsAny(name, " \t") {
				variables[name] = strings.TrimSpace(after)
			}
			break
		}
	}
	return variables
}

func makeTargetRecipe(t *testing.T, makefile, target string) string {
	t.Helper()

	lines := strings.Split(makefile, "\n")
	declaration := target + ":"
	for index, line := range lines {
		if !strings.HasPrefix(line, declaration) {
			continue
		}
		for _, recipe := range lines[index+1:] {
			if strings.HasPrefix(recipe, "\t") {
				return strings.TrimSpace(recipe)
			}
			if strings.TrimSpace(recipe) != "" {
				break
			}
		}
	}
	t.Fatalf("Makefile target %q has no recipe", target)
	return ""
}

func expandMakeVariables(value string, variables map[string]string) string {
	for range len(variables) + 1 {
		before := value
		for name, replacement := range variables {
			value = strings.ReplaceAll(value, "$("+name+")", replacement)
		}
		if value == before {
			return value
		}
	}
	return value
}

func writeGateFixtureModule(t *testing.T, testSource string) string {
	t.Helper()

	moduleRoot := t.TempDir()
	writeGateFixtureFile(t, filepath.Join(moduleRoot, "go.mod"), fmt.Sprintf("module %s\n\ngo 1.26\n", gateFixtureModule), 0o644)
	writeGateFixtureFile(t, filepath.Join(moduleRoot, "gate_test.go"), testSource, 0o644)
	return moduleRoot
}

func writeMaskingWrapper(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "rtk")
	writeGateFixtureFile(t, path, "#!/bin/sh\n\"$@\"\nexit 0\n", 0o755)
	return path
}

func writeGateFixtureFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
