// Suite: documentation contract Makefile and CI wiring.
// Invariant: pull requests run each discovered documentation contract once.
// Boundary IN: the repository Makefile with temporary command stubs, workflow YAML, and contract declarations.
// Boundary OUT: real contract execution, workflow dispatch, publication, and network access.
package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestVerifyDocsRunsTheDocumentationContractsOnce(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	dir := t.TempDir()
	writeFile(t, dir, "calls", "")
	for _, name := range []string{"go", "roundfix-stub"} {
		writeFile(t, dir, name, "#!/bin/sh\n[ \"$1\" = env ] && exit 0\nprintf '%s\\n' \"$*\" >> calls\n")
		if err := os.Chmod(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "verify-docs", "GO="+filepath.Join(dir, "go"), "BIN="+filepath.Join(dir, "roundfix-stub"))
	command.Dir = dir
	// Isolate make variables and its cache probe from the caller's environment.
	command.Env = []string{"PATH=" + os.Getenv("PATH"), "TMPDIR=" + dir}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("make verify-docs: %v\n%s", err, output)
	}
	calls, err := os.ReadFile(filepath.Join(dir, "calls"))
	if err != nil {
		t.Fatal(err)
	}
	var docs, builds, repo, budget, spec, bound int
	for _, line := range strings.Split(strings.TrimSpace(string(calls)), "\n") {
		switch {
		case strings.HasPrefix(line, "run -buildvcs=false ./cmd/gocache-bound "):
			bound++
		case strings.HasPrefix(line, "build "):
			builds++
		case strings.HasPrefix(line, "test ") && strings.Contains(line, "./internal/docscontract"):
			if strings.Contains(line, "TestCheckCorpusBudget") {
				budget++
			} else {
				docs++
			}
		case strings.HasPrefix(line, "test ") && strings.Contains(line, "-tags repocontract "):
			repo++
		case strings.Contains(line, "./cmd/gate-time report -gate before-pull-request"):
		case line == "spec check":
			spec++
		default:
			t.Errorf("unexpected invocation %q", line)
		}
	}
	if docs != 1 || builds != 1 || repo != 1 || budget != 1 || spec != 1 || bound != 1 {
		t.Fatalf("invocations: docs=%d build=%d repo=%d budget=%d spec=%d bound=%d; want %d/1/1/1/1/1\n%s", docs, builds, repo, budget, spec, bound, 1, calls)
	}
}

func TestTheCIPullRequestRunsEachDocumentationContractOnce(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	workflow := readVerifyingWorkflow(t, filepath.Join(root, ".github/workflows/ci-verify.yml"))
	count := 0
	for _, part := range workflow.Jobs["parts"].Strategy.Matrix.Part {
		if part == "contracts-docs" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("docs parts = %d, want 1", count)
	}
	for name, job := range workflow.Jobs {
		for _, step := range job.Steps {
			if strings.Contains(step.Run, "make verify-docs") || step.Env["VERIFY_DOCS_TEST"] != "" {
				t.Errorf("%s repeats documentation contracts: %+v", name, step)
			}
		}
	}
	makefile, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	recipe := regexp.MustCompile(`(?ms)^verify-part:.*?(?:\n\n|\z)`).FindString(string(makefile))
	if !strings.Contains(recipe, "-contracts -all -tag") || strings.Contains(recipe, "spec-budget") || strings.Contains(recipe, "docs-test") {
		t.Fatalf("docs part must run discovery without repeating a contract: %s", recipe)
	}
}

func TestEveryDocumentationContractFileIsAlwaysSelected(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(findRepositoryRoot(t), "internal", "docscontract")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		contents, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(contents), "//go:build docscontract") {
			continue
		}
		count++
		always := false
		for _, line := range strings.Split(string(contents), "\n") {
			always = always || strings.TrimSpace(line) == "//verify:always"
		}
		if !always {
			t.Errorf("%s must declare //verify:always so settlements select it regardless of changed paths", entry.Name())
		}
	}
	if count == 0 {
		t.Fatal("no docscontract-tagged test files found")
	}
}
