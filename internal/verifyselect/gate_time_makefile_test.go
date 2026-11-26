// Suite: verification recipes.
// Invariant: the real recipes execute the plan, preserve failures and report each gate.
// Boundary IN: Makefile and disposable go, git and gate-time stand-ins.
// Boundary OUT: nested Go test execution and production verification.
package verifyselect_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func runGateRecipe(t *testing.T, mode, target, part string, mutations ...string) (string, string, error) {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "calls", "")
	script := `#!/bin/sh
printf '%s\n' "$*" >> "$CALLS"
case "$1" in
 env) exit 0;;
 run) case "$*" in
 *cmd/verify-select*) case "$*" in
 *-contracts*) printf 'docscontract ^(TestDocs)$ ./internal/docscontract\n';;
 *-job*) case "$*" in *cli-1*) printf '%s\n' '-run' '^(TestOne)$' './internal/cli';; *cli-2*) printf '%s\n' '-skip' '^(TestOne)$' './internal/cli';; *) printf './internal/app\n';; esac;;
 *) if [ "$MODE" = settlement ]; then echo 'verify-select: packages: none at a settlement; the Complete Suite Run runs them once per Spec' >&2; else echo 'verify-select: packages: every package with fresh results for the Complete Suite Run' >&2; echo './...'; fi;; esac;;
 *cmd/gate-time*) case "$*" in
 *'tests '*) echo 'gate time: test binary 0s (target: 300s)';;
 *task-settlement*) echo 'gate time: task settlement 31s (target: median 120s, p95 240s)';;
 *complete-suite-run*) echo 'gate time: complete suite run 352s (target: 360s)';;
 *before-pull-request*) echo 'gate time: before pull request 60s (target: 90s)';;
 *) echo "gate time: pull request, job $PART_NAME 412s (target: median 480s, p95 600s)";; esac;; esac;;
 test) case "$*" in *-json*) echo '{"Action":"pass","Package":"fixture","Elapsed":0}';; esac;;
esac
exit 0
`
	writeFile(t, dir, "go", script)
	writeFile(t, dir, "git", "#!/bin/sh\nexit 0\n")
	for _, name := range []string{"go", "git"} {
		if err := os.Chmod(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	if len(mutations) > 0 {
		data, err := os.ReadFile(makefile)
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		for i := 0; i < len(mutations); i += 2 {
			if !strings.Contains(text, mutations[i]) {
				t.Fatal("mutation did not match")
			}
			text = strings.ReplaceAll(text, mutations[i], mutations[i+1])
		}
		writeFile(t, dir, "Makefile", text)
		makefile = filepath.Join(dir, "Makefile")
	}

	cmd := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, target, "PART="+part, "MAKE=make --no-print-directory -f "+makefile, "OWNED_SKILLS=", "GO="+filepath.Join(dir, "go"), "GOFMT=true", "BIN="+filepath.Join(dir, "go"), "GOCACHE_BOUND_MB=0", "RTK=", "GO_HUMAN="+filepath.Join(dir, "go"))
	cmd.Dir = dir
	cmd.Env = append([]string{"PATH=" + dir + ":" + os.Getenv("PATH"), "HOME=" + dir, "GOCACHE=" + dir, "CALLS=" + filepath.Join(dir, "calls"), "MODE=" + mode, "PART_NAME=" + part}, "TMPDIR="+dir)
	out, err := cmd.CombinedOutput()
	calls, readErr := os.ReadFile(filepath.Join(dir, "calls"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	return string(out), string(calls), err
}
func TestTheSettlementStartsNoPackageTest(t *testing.T) {
	t.Parallel()
	out, calls, err := runGateRecipe(t, "settlement", "verify-changed", "")
	if err != nil {
		t.Fatal(err, out)
	}
	for _, line := range strings.Split(calls, "\n") {
		if strings.HasPrefix(line, "test ") && !strings.Contains(line, "-tags docscontract") && !strings.Contains(line, "TestCoverageEquivalence") {
			t.Fatal(line)
		}
	}
	if !strings.Contains(calls, "TestCoverageRecordCountsNoPackageUnderDocs") {
		t.Fatal(calls)
	}
}
func TestTheCompleteSuiteRunRunsEveryPackageOnceWithFreshResults(t *testing.T) {
	t.Parallel()
	out, calls, err := runGateRecipe(t, "complete", "verify-changed", "")
	if err != nil || strings.Count(calls, "test -count=1 -parallel 16 ./...\n") != 1 {
		t.Fatal(err, out, calls)
	}
}
func TestTheSettlementPrintsItsGateTimeReport(t *testing.T) {
	t.Parallel()
	out, _, err := runGateRecipe(t, "settlement", "verify-changed", "")
	if err != nil || !strings.HasSuffix(out, "gate time: task settlement 31s (target: median 120s, p95 240s)\n") {
		t.Fatal(err, out)
	}
}
func TestTheCompleteSuiteRunPrintsItsGateTimeReport(t *testing.T) {
	t.Parallel()
	out, _, err := runGateRecipe(t, "complete", "verify-changed", "")
	if err != nil || !strings.HasSuffix(out, "gate time: complete suite run 352s (target: 360s)\n") {
		t.Fatal(err, out)
	}
}
func TestVerifyDocsPrintsItsGateTimeReport(t *testing.T) {
	t.Parallel()
	out, _, err := runGateRecipe(t, "complete", "verify-docs", "")
	if err != nil || !strings.HasSuffix(out, "gate time: before pull request 60s (target: 90s)\n") {
		t.Fatal(err, out)
	}
}
func TestEveryVerificationPartPrintsItsGateTimeReport(t *testing.T) {
	t.Parallel()
	for _, part := range []string{"packages", "cli-1", "cli-2", "contracts-repo", "contracts-docs"} {
		t.Run(part, func(t *testing.T) {
			out, calls, err := runGateRecipe(t, "reach", "verify-part", part)
			want := fmt.Sprintf("gate time: pull request, job %s 412s (target: median 480s, p95 600s)\n", part)
			if err != nil || !strings.HasSuffix(out, want) || !strings.Contains(out, "gate time: test binary 0s (target: 300s)\n") || !strings.Contains(calls, "cmd/gate-time tests -events") {
				t.Fatal(err, out, calls)
			}
		})
	}
}
func TestVerifyPartRunsThePartsThePlanPrints(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ part, want string }{{"packages", "test -json -parallel 16 ./internal/app"}, {"cli-1", "test -json -parallel 16 -run ^(TestOne)$ ./internal/cli"}, {"cli-2", "test -json -parallel 16 -skip ^(TestOne)$ ./internal/cli"}} {
		out, calls, err := runGateRecipe(t, "reach", "verify-part", c.part)
		if err != nil || !strings.Contains(calls, c.want+"\n") {
			t.Fatal(err, out, calls)
		}
	}
	for _, seed := range []struct{ name, replacement string }{{"missing", "cli-2"}, {"duplicate", "cli-1"}} {
		t.Run(seed.name, func(t *testing.T) {
			_, first, e := runGateRecipe(t, "reach", "verify-part", "cli-1")
			if e != nil {
				t.Fatal(e)
			}
			_, second, e := runGateRecipe(t, "reach", "verify-part", "cli-2", `-job "$(PART)"`, `-job `+seed.replacement)
			if e != nil {
				t.Fatal(e)
			}
			if seed.name == "missing" {
				_, first, e = runGateRecipe(t, "reach", "verify-part", "cli-1", `-job "$(PART)"`, `-job cli-2`)
				if e != nil {
					t.Fatal(e)
				}
			}
			counts := recipeShardMemberships(first + second)
			if counts["TestOne"] == 1 && counts["TestTwo"] == 1 {
				t.Fatal("seeded recipe still partitions tests")
			}
		})
	}

}

func recipeShardMemberships(calls string) map[string]int {
	counts := map[string]int{"TestOne": 0, "TestTwo": 0}
	for _, line := range strings.Split(calls, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "test" || !strings.Contains(line, "./internal/cli") {
			continue
		}
		for i, arg := range fields {
			if (arg == "-run" || arg == "-skip") && i+1 < len(fields) {
				p := regexp.MustCompile(fields[i+1])
				for name := range counts {
					match := p.MatchString(name)
					if arg == "-skip" {
						match = !match
					}
					if match {
						counts[name]++
					}
				}
			}
		}
	}
	return counts
}
