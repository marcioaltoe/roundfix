// Suite: Full Contract Run wiring.
// Invariant: discovery drives ordered contract execution on main and before releases.
// Boundary IN: the repository Makefile with temporary command stubs and workflow YAML.
// Boundary OUT: real contract execution, workflow dispatch, publication, and network access.
package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestVerifyContractsRunsEveryContract(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	const lines = "docscontract ^(TestDocs)$ ./internal/docscontract\nrepocontract ^(TestOne|TestTwo)$ ./internal/baseline ./skills\n"
	const first = "test -count=1 -tags docscontract -run ^(TestDocs)$ ./internal/docscontract\n"
	const second = "test -count=1 -tags repocontract -run ^(TestOne|TestTwo)$ ./internal/baseline ./skills\n"
	for _, scenario := range []struct {
		name, selection, want             string
		selectorFail, goFail, wantFailure bool
	}{
		{name: "ordered invocations", selection: lines, want: first + second},
		{name: "first failure stops execution", selection: lines, goFail: true, wantFailure: true, want: first},
		{name: "selector failure prevents execution", selection: lines, selectorFail: true, wantFailure: true},
		{name: "empty selection"},
		{name: "line without packages is skipped", selection: "docscontract ^(TestEmpty)$\n" + lines, want: first + second},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "selection", scenario.selection)
			writeFile(t, dir, "calls", "")
			selector := "#!/bin/sh\n[ \"$#\" -eq 2 ] && [ \"$1\" = '-contracts' ] && [ \"$2\" = '-all' ] || exit 92\ncat selection\n"
			if scenario.selectorFail {
				selector += "exit 17\n"
			}
			goStub := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> calls\n"
			if scenario.goFail {
				goStub += "exit 19\n"
			}
			writeFile(t, dir, "selector", selector)
			writeFile(t, dir, "go", goStub)
			for _, name := range []string{"selector", "go"} {
				if err := os.Chmod(filepath.Join(dir, name), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "verify-contracts", "VERIFY_SELECT="+filepath.Join(dir, "selector"), "GO="+filepath.Join(dir, "go"))
			command.Dir = dir
			output, err := command.CombinedOutput()
			if (err != nil) != scenario.wantFailure {
				t.Fatalf("make error = %v, want failure %t\n%s", err, scenario.wantFailure, output)
			}
			calls, err := os.ReadFile(filepath.Join(dir, "calls"))
			if err != nil {
				t.Fatal(err)
			}
			if string(calls) != scenario.want {
				t.Fatalf("go invocations = %q, want %q\n%s", calls, scenario.want, output)
			}
		})
	}
}

func TestMainAndReleaseRunEveryContract(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	type step struct {
		Name string `yaml:"name"`
		Run  string `yaml:"run"`
		If   string `yaml:"if"`
	}
	type workflow struct {
		On struct {
			Push struct {
				Branches []string `yaml:"branches"`
			} `yaml:"push"`
		} `yaml:"on"`
		Jobs map[string]struct {
			Steps []step `yaml:"steps"`
		} `yaml:"jobs"`
	}
	read := func(name string) workflow {
		t.Helper()
		contents, err := os.ReadFile(filepath.Join(root, ".github", "workflows", name))
		if err != nil {
			t.Fatal(err)
		}
		var result workflow
		if err := yaml.Unmarshal(contents, &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	ci := read("ci-verify.yml")
	mainPush := false
	for _, branch := range ci.On.Push.Branches {
		mainPush = mainPush || branch == "main"
	}
	if !mainPush {
		t.Error("CI must trigger on pushes to main")
	}
	docsIndex, contractIndex, contractCount := -1, -1, 0
	for i, s := range ci.Jobs["verify"].Steps {
		if s.Name == "Verify docs" {
			docsIndex = i
		}
		if s.Run == "make verify-contracts" {
			contractIndex = i
			contractCount++
			if s.Name != "Verify every contract" || s.If != "github.event_name == 'push'" {
				t.Errorf("CI contract step must be named and restricted to push, got %+v", s)
			}
		}
	}
	if contractCount != 1 || docsIndex < 0 || contractIndex != docsIndex+1 {
		t.Errorf("CI contract steps = %d, indexes docs=%d contracts=%d; want one directly after docs", contractCount, docsIndex, contractIndex)
	}
	for jobName, job := range ci.Jobs {
		for _, s := range job.Steps {
			if s.Run == "make verify-contracts" && s.If != "github.event_name == 'push'" {
				t.Errorf("job %s has a contract step that is not restricted to push: %+v", jobName, s)
			}
		}
	}
	release := read("release.yml")
	gateIndex, preflightIndex := -1, -1
	contractIndex, contractCount = -1, 0
	for i, s := range release.Jobs["release"].Steps {
		if s.Run == "make verify" {
			gateIndex = i
		}
		if s.Name == "Publication preflight" {
			preflightIndex = i
		}
		if s.Run == "make verify-contracts" {
			contractIndex = i
			contractCount++
			if s.Name != "Verify every contract" || s.If != "" {
				t.Errorf("release contract step must be named and unconditional, got %+v", s)
			}
		}
	}
	if contractCount != 1 || gateIndex < 0 || contractIndex != gateIndex+1 || preflightIndex <= contractIndex {
		t.Errorf("release contract steps = %d, indexes gate=%d contracts=%d preflight=%d; want one after gate and before preflight", contractCount, gateIndex, contractIndex, preflightIndex)
	}
}
