// Suite: documentation contract Makefile and CI wiring.
// Invariant: pull requests skip only documentation contracts already selected as always relevant.
// Boundary IN: the repository Makefile with temporary command stubs, workflow YAML, and contract declarations.
// Boundary OUT: real contract execution, workflow dispatch, publication, and network access.
package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	const skipMessage = "docs-test: skipped; make verify-changed ran internal/docscontract earlier in this job"
	for _, scenario := range []struct {
		name     string
		skip     bool
		wantDocs int
	}{
		{name: "default", wantDocs: 1},
		{name: "already selected", skip: true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			writeFile(t, dir, "calls", "")
			for _, name := range []string{"go", "roundfix-stub"} {
				writeFile(t, dir, name, "#!/bin/sh\nprintf '%s\\n' \"$*\" >> calls\n")
				if err := os.Chmod(filepath.Join(dir, name), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "verify-docs", "GO="+filepath.Join(dir, "go"), "BIN="+filepath.Join(dir, "roundfix-stub"))
			command.Dir = dir
			// Isolate make variables and its cache probe from the caller's environment.
			command.Env = []string{"PATH=" + os.Getenv("PATH"), "TMPDIR=" + dir}
			if scenario.skip {
				command.Args = append(command.Args, "VERIFY_DOCS_TEST=docs-test-skipped")
			}
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("make verify-docs: %v\n%s", err, output)
			}
			calls, err := os.ReadFile(filepath.Join(dir, "calls"))
			if err != nil {
				t.Fatal(err)
			}
			var docs, builds, repo, budget, spec int
			for _, line := range strings.Split(strings.TrimSpace(string(calls)), "\n") {
				switch {
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
				case line == "spec check":
					spec++
				default:
					t.Errorf("unexpected invocation %q", line)
				}
			}
			if docs != scenario.wantDocs || builds != 1 || repo != 1 || budget != 1 || spec != 1 {
				t.Fatalf("invocations: docs=%d build=%d repo=%d budget=%d spec=%d; want %d/1/1/1/1\n%s", docs, builds, repo, budget, spec, scenario.wantDocs, calls)
			}
			if got := strings.Contains(string(output), skipMessage); got != scenario.skip {
				t.Errorf("skip message present = %t, want %t\n%s", got, scenario.skip, output)
			}
		})
	}
}

func TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts(t *testing.T) {
	t.Parallel()
	contents, err := os.ReadFile(filepath.Join(findRepositoryRoot(t), ".github", "workflows", "ci-verify.yml"))
	if err != nil {
		t.Fatal(err)
	}
	var workflow struct {
		Jobs map[string]struct {
			Steps []struct {
				Name string            `yaml:"name"`
				Run  string            `yaml:"run"`
				Env  map[string]string `yaml:"env"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(contents, &workflow); err != nil {
		t.Fatal(err)
	}
	const want = "${{ github.event_name == 'pull_request' && 'docs-test-skipped' || 'docs-test' }}"
	count := 0
	for _, step := range workflow.Jobs["verify"].Steps {
		if step.Name != "Verify docs" {
			continue
		}
		count++
		if step.Run != "make verify-docs" || step.Env["VERIFY_DOCS_TEST"] != want {
			t.Errorf("Verify docs must run make verify-docs with VERIFY_DOCS_TEST=%q; got %+v", want, step)
		}
	}
	if count != 1 {
		t.Errorf("Verify docs steps = %d, want 1", count)
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
			t.Errorf("%s must declare //verify:always so pull requests can skip its repeated docs-test run", entry.Name())
		}
	}
	if count == 0 {
		t.Fatal("no docscontract-tagged test files found")
	}
}
