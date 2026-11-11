// Suite: suite-time budget recording and CI push wiring.
// Invariant: an overrun records a verdict that CI judges after the contracts.
// Boundary IN: the real Makefile with a temporary Go stub and workflow YAML.
// Boundary OUT: real Go test execution, workflow dispatch, and network access.
package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTheSuiteTimeBudgetIsRecordedThenJudged(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	const hint = "the budget is SUITE_BUDGET_SECONDS in the Makefile; raise it deliberately or find the regression"
	const exceeded = "suite-time budget exceeded: the suite took 500s against a budget of 420s"
	const within = "suite-time budget: 300s of 420s"
	for _, scenario := range []struct {
		name, target, budget, initialVerdict    string
		filePattern, outputPattern, exactOutput string
		record, failSuite, wantFailure, noHint  bool
	}{
		{
			name: "an overrun is recorded without failing", target: "test-budget", budget: "-1", record: true, noHint: true,
			filePattern: `^suite-time budget exceeded: the suite took [0-9]+s against a budget of -1s\n$`,
		},
		{
			name: "without a verdict file an overrun fails at once", target: "test-budget", budget: "-1", wantFailure: true,
			outputPattern: `(?m)^suite-time budget exceeded: the suite took [0-9]+s against a budget of -1s\n` + regexp.QuoteMeta(hint) + `\n`,
		},
		{
			name: "a failing suite records no verdict", target: "test-budget", record: true, failSuite: true, wantFailure: true,
			initialVerdict: "suite-time budget: 1s of 420s\n",
		},
		{
			name: "the default budget is 420 seconds", target: "test-budget", record: true,
			filePattern: `^suite-time budget: [0-9]+s of 420s\n$`,
		},
		{
			name: "the check fails an exceeded verdict", target: "suite-budget-check", record: true, wantFailure: true,
			initialVerdict: exceeded + "\n", filePattern: "^" + regexp.QuoteMeta(exceeded+"\n") + "$",
			outputPattern: "^" + regexp.QuoteMeta(exceeded+"\n"+hint+"\n"),
		},
		{
			name: "the check passes a verdict within the budget", target: "suite-budget-check", record: true,
			initialVerdict: within + "\n", filePattern: "^" + regexp.QuoteMeta(within+"\n") + "$", exactOutput: within + "\n",
		},
		{
			name: "the check fails without a verdict", target: "suite-budget-check", record: true, wantFailure: true,
			outputPattern: `suite-time budget: no verdict recorded at `,
		},
		{
			name: "the check fails an unreadable verdict", target: "suite-budget-check", record: true, wantFailure: true,
			initialVerdict: "garbage\n", filePattern: "^garbage\n$", outputPattern: `suite-time budget: unreadable verdict at `,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			verdictPath := filepath.Join(dir, "suite-budget-verdict")
			stub := "#!/bin/sh\nexit 0\n"
			if scenario.failSuite {
				stub = "#!/bin/sh\nif [ \"$1\" = test ]; then exit 1; fi\nexit 0\n"
			}
			writeFile(t, dir, "go", stub)
			if err := os.Chmod(filepath.Join(dir, "go"), 0o755); err != nil {
				t.Fatal(err)
			}
			if scenario.initialVerdict != "" {
				writeFile(t, dir, "suite-budget-verdict", scenario.initialVerdict)
			}
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, scenario.target, "GO="+filepath.Join(dir, "go"))
			command.Dir = dir
			// Do not inherit caller overrides or let make touch the real home or cache.
			command.Env = []string{"PATH=" + os.Getenv("PATH"), "TMPDIR=" + dir, "HOME=" + dir}
			if scenario.record {
				command.Env = append(command.Env, "SUITE_BUDGET_VERDICT="+verdictPath)
			}
			if scenario.budget != "" {
				command.Args = append(command.Args, "SUITE_BUDGET_SECONDS="+scenario.budget)
			}
			output, err := command.CombinedOutput()
			if (err != nil) != scenario.wantFailure {
				t.Fatalf("make %s error = %v, want failure %t\n%s", scenario.target, err, scenario.wantFailure, output)
			}
			if scenario.outputPattern != "" && !regexp.MustCompile(scenario.outputPattern).Match(output) {
				t.Errorf("output = %q, want pattern %q", output, scenario.outputPattern)
			}
			if scenario.exactOutput != "" && string(output) != scenario.exactOutput {
				t.Errorf("output = %q, want %q", output, scenario.exactOutput)
			}
			if scenario.noHint && strings.Contains(string(output), hint) {
				t.Errorf("recording the overrun printed the immediate-failure hint: %s", output)
			}
			verdict, err := os.ReadFile(verdictPath)
			if scenario.filePattern == "" {
				if !os.IsNotExist(err) {
					t.Fatalf("verdict = %q, error = %v, want no file", verdict, err)
				}
			} else if err != nil {
				t.Fatal(err)
			} else if !regexp.MustCompile(scenario.filePattern).Match(verdict) {
				t.Errorf("verdict = %q, want pattern %q", verdict, scenario.filePattern)
			}
		})
	}
}

func TestTheCIPushJudgesTheSuiteTimeBudgetLast(t *testing.T) {
	t.Parallel()
	contents, err := os.ReadFile(filepath.Join(findRepositoryRoot(t), ".github", "workflows", "ci-verify.yml"))
	if err != nil {
		t.Fatal(err)
	}
	var workflow struct {
		Jobs map[string]struct {
			Steps []struct {
				Name string            `yaml:"name"`
				ID   string            `yaml:"id"`
				If   string            `yaml:"if"`
				Run  string            `yaml:"run"`
				Env  map[string]string `yaml:"env"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(contents, &workflow); err != nil {
		t.Fatal(err)
	}
	const verdictPath = "${{ runner.temp }}/suite-budget-verdict"
	const condition = "${{ !cancelled() && github.event_name == 'push' && steps.verify.outcome == 'success' }}"
	count := 0
	for jobName, job := range workflow.Jobs {
		for _, step := range job.Steps {
			if _, exists := step.Env["SUITE_BUDGET_SECONDS"]; exists {
				t.Errorf("job %s step %q overrides the Makefile's budget", jobName, step.Name)
			}
			if step.Env["VERIFY_TEST_TARGET"] == "test-budget" {
				count++
				if jobName != "verify" || step.Name != "Verify" || step.ID != "verify" || step.If != "github.event_name == 'push'" || step.Run != "make verify" || step.Env["SUITE_BUDGET_VERDICT"] != verdictPath {
					t.Errorf("timed Verify step must record the push verdict, got job %s: %+v", jobName, step)
				}
			}
		}
	}
	if count != 1 {
		t.Errorf("steps selecting test-budget = %d, want exactly 1", count)
	}
	steps := workflow.Jobs["verify"].Steps
	if len(steps) < 2 {
		t.Fatalf("verify job has %d steps, want the contracts followed by the budget check", len(steps))
	}
	last := steps[len(steps)-1]
	if steps[len(steps)-2].Name != "Verify every contract" || steps[len(steps)-2].Run != "make verify-contracts" {
		t.Errorf("step before the budget check must run every contract, got %+v", steps[len(steps)-2])
	}
	if last.Name != "Check the suite-time budget" || last.If != condition || last.Run != "make suite-budget-check" || last.Env["SUITE_BUDGET_VERDICT"] != verdictPath {
		t.Errorf("last step must judge the recorded verdict whenever push Verify succeeded, got %+v", last)
	}
}
