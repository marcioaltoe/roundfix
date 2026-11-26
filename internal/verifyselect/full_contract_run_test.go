// Suite: Full Contract Run wiring.
// Invariant: discovery drives ordered contract execution once per pull request.
// Boundary IN: the repository Makefile with temporary command stubs and workflow YAML.
// Boundary OUT: real contract execution, workflow dispatch, publication, and network access.
package verifyselect_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"roundfix/internal/verifyselect"
	"strconv"
	"strings"
	"testing"
)

func TestVerifyContractsRunsEveryContract(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	const lines = "docscontract ^(TestDocs)$ ./internal/docscontract\nrepocontract ^(TestOne|TestTwo)$ ./internal/baseline ./skills\n"
	const first = "test -count=1 -timeout 20m -tags docscontract -run ^(TestDocs)$ ./internal/docscontract\n"
	const second = "test -count=1 -timeout 20m -tags repocontract -run ^(TestOne|TestTwo)$ ./internal/baseline ./skills\n"
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
			goStub := "#!/bin/sh\n[ \"$1\" = env ] && exit 0\nprintf '%s\\n' \"$*\" >> calls\n"
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
			// Exercise the default bound with only disposable cache and home paths.
			command.Env = []string{"PATH=" + os.Getenv("PATH"), "TMPDIR=" + dir, "HOME=" + dir}
			bound := "run -buildvcs=false ./cmd/gocache-bound -dir " + filepath.Join(dir, "roundfix-go-build-"+strconv.Itoa(os.Getuid())) + " -max-mb 10240 -idle 2h\n"
			output, err := command.CombinedOutput()
			if (err != nil) != scenario.wantFailure {
				t.Fatalf("make error = %v, want failure %t\n%s", err, scenario.wantFailure, output)
			}
			calls, err := os.ReadFile(filepath.Join(dir, "calls"))
			if err != nil {
				t.Fatal(err)
			}
			if want := bound + scenario.want; string(calls) != want {
				t.Fatalf("go invocations = %q, want %q\n%s", calls, want, output)
			}
		})
	}
}

func TestEveryPullRequestRunsEveryContractOnce(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	contracts, err := verifyselect.DiscoverContracts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(contracts) == 0 {
		t.Fatal("no contracts discovered")
	}
	workflow := readVerifyingWorkflow(t, filepath.Join(root, ".github/workflows/ci-verify.yml"))
	counts := make(map[string]int)
	for _, part := range workflow.Jobs["parts"].Strategy.Matrix.Part {
		if !strings.HasPrefix(part, "contracts-") {
			continue
		}
		tag := strings.Replace(part, "contracts-", "", 1) + "contract"
		var selected []verifyselect.ContractTest
		for _, c := range contracts {
			if c.Tag == tag {
				selected = append(selected, c)
			}
		}
		dir := t.TempDir()
		selector := "#!/bin/sh\n[ \"$*\" = '-contracts -all -tag " + tag + "' ] || exit 92\ncat '" + filepath.Join(dir, "selection") + "'\n"
		writeFile(t, dir, "selection", strings.Join(verifyselect.ContractInvocations(selected), "\n")+"\n")
		writeFile(t, dir, "selector", selector)
		writeFile(t, dir, "calls", "")
		writeFile(t, dir, "go", `#!/bin/sh
[ "$1" = env ] && exit 0
printf '%s\n' "$*" >> "$CALLS"
case "$1" in test) printf '%s\n' '{"Action":"pass","Package":"fixture","Elapsed":0}';; esac
`)
		for _, name := range []string{"selector", "go"} {
			if err := os.Chmod(filepath.Join(dir, name), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		makefile := filepath.Join(root, "Makefile")
		cmd := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "verify-part", "PART="+part, "VERIFY_SELECT="+filepath.Join(dir, "selector"), "GO="+filepath.Join(dir, "go"), "BIN="+filepath.Join(dir, "go"), "MAKE=make --no-print-directory -f "+makefile, "GOCACHE_BOUND_MB=0")
		cmd.Dir = dir
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + dir, "TMPDIR=" + dir, "GOCACHE=" + dir, "CALLS=" + filepath.Join(dir, "calls"), "GOFLAGS=-count=1"}
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s: %v\n%s", part, err, out)
		}
		calls, err := os.ReadFile(filepath.Join(dir, "calls"))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(calls), "\n") {
			fields := strings.Fields(line)
			if len(fields) == 0 || fields[0] != "test" {
				continue
			}
			runTag, pattern := "", ".*"
			for i, arg := range fields {
				if i+1 < len(fields) {
					switch arg {
					case "-tags":
						runTag = fields[i+1]
					case "-run":
						pattern = fields[i+1]
					}
				}
			}
			if runTag == "" {
				continue
			}
			match, err := regexp.Compile(pattern)
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range contracts {
				if c.Tag == runTag && match.MatchString(c.Name) {
					for _, arg := range fields {
						if arg == "./"+c.Package || arg == "./..." {
							counts[c.Tag+"/"+c.Package+"/"+c.Name]++
						}
					}
				}
			}
		}
	}
	for _, c := range contracts {
		key := fmt.Sprintf("%s/%s/%s", c.Tag, c.Package, c.Name)
		if counts[key] != 1 {
			t.Errorf("%s runs %d times, want once", key, counts[key])
		}
	}
}
