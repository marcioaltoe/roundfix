// Suite: Verification runner and release dependency contracts.
// Invariant: every verifying job uses macos-26 and publication depends on verification.
// Boundary IN: this repository's workflow YAML.
// Boundary OUT: Baseline assets, GitHub dispatch and network access.
package verifyselect_test

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// verifyingRunner is the label of every job that runs a test, a vet or a contract.
const verifyingRunner = "macos-26"

type verifyingWorkflowStep struct {
	Name string            `yaml:"name"`
	Uses string            `yaml:"uses"`
	Run  string            `yaml:"run"`
	If   string            `yaml:"if"`
	With map[string]string `yaml:"with"`
}

type verifyingWorkflow struct {
	Name string `yaml:"name"`
	Jobs map[string]struct {
		Name        string                  `yaml:"name"`
		Runner      yaml.Node               `yaml:"runs-on"`
		Needs       string                  `yaml:"needs"`
		If          string                  `yaml:"if"`
		Timeout     int                     `yaml:"timeout-minutes"`
		Permissions map[string]string       `yaml:"permissions"`
		Steps       []verifyingWorkflowStep `yaml:"steps"`
	} `yaml:"jobs"`
}

func readVerifyingWorkflow(t *testing.T, path string) verifyingWorkflow {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var workflow verifyingWorkflow
	if err := yaml.Unmarshal(contents, &workflow); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return workflow
}

// releaseVerifyRunnerOverride is the runner the release's verify job uses
// under the maintainer's override of 2026-11-25 (DECISÃO DE MARCIO, ai-memory
// 01a21237, 06:58 BRT): the full suite exceeds go test's 10-minute default on
// the standard macOS runner, so the release verifies on Linux until Spec 0421
// (test-suite performance) removes this exception and returns it to macOS.
const releaseVerifyRunnerOverride = "ubuntu-latest"

// wantVerifyingRunner is the runner a verifying job of workflow must name.
func wantVerifyingRunner(workflow, job string) string {
	if workflow == "release.yml" && job == "verify" {
		return releaseVerifyRunnerOverride
	}
	return verifyingRunner
}

func TestEveryVerifyingJobRunsOnTheStandardMacOSRunner(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(findRepositoryRoot(t), ".github", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	verification := regexp.MustCompile(`\b(?:make[ \t]+verify(?:-changed|-docs|-contracts)?|go[ \t]+(?:test|vet))(?:[\s;|&]|$)`)
	jobs := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		workflow := readVerifyingWorkflow(t, filepath.Join(dir, entry.Name()))
		for name, job := range workflow.Jobs {
			for _, step := range job.Steps {
				if !verification.MatchString(step.Run) {
					continue
				}
				jobs++
				if want := wantVerifyingRunner(entry.Name(), name); job.Runner.Kind != yaml.ScalarNode || job.Runner.Value != want {
					t.Errorf("%s job %s runs %q on %v, want %s", entry.Name(), name, step.Run, job.Runner, want)
				}
			}
		}
	}
	if jobs == 0 {
		t.Fatal("no verifying steps found")
	}
}

func TestTheReleaseVerifiesBeforeItBuildsOrPublishes(t *testing.T) {
	t.Parallel()
	workflow := readVerifyingWorkflow(t, filepath.Join(findRepositoryRoot(t), ".github/workflows/release.yml"))
	verify, exists := workflow.Jobs["verify"]
	if !exists {
		t.Fatal("release workflow has no verify job")
	}
	if verify.Runner.Value != wantVerifyingRunner("release.yml", "verify") || !reflect.DeepEqual(verify.Permissions, map[string]string{"contents": "read"}) {
		t.Errorf("verify runner/permissions = %s/%v", verify.Runner.Value, verify.Permissions)
	}
	if len(verify.Steps) != 4 {
		t.Fatalf("verify steps = %d, want checkout, Go setup and two gates", len(verify.Steps))
	}
	checkout, setup := verify.Steps[0], verify.Steps[1]
	if !strings.HasPrefix(checkout.Uses, "actions/checkout@") || !reflect.DeepEqual(checkout.With, map[string]string{"fetch-depth": "0", "persist-credentials": "false"}) {
		t.Errorf("verify checkout = %+v", checkout)
	}
	if !strings.HasPrefix(setup.Uses, "actions/setup-go@") || !reflect.DeepEqual(setup.With, map[string]string{"go-version": "1.26", "check-latest": "true", "cache": "false"}) {
		t.Errorf("verify Go setup = %+v", setup)
	}
	for i, want := range []verifyingWorkflowStep{{Name: "Verify gate", Run: "make verify"}, {Name: "Verify every contract", Run: "make verify-contracts"}} {
		if got := verify.Steps[i+2]; got.Name != want.Name || got.Run != want.Run || got.If != "" {
			t.Errorf("verification step = %+v, want %+v unconditionally", got, want)
		}
	}
	release, exists := workflow.Jobs["release"]
	if !exists || release.Needs != "verify" {
		t.Fatal("release must need verify")
	}
	for name, job := range workflow.Jobs {
		if job.If != "" {
			t.Errorf("job %s has condition %q; verification must succeed before publication", name, job.If)
		}
	}
	for _, step := range verify.Steps {
		if step.If != "" {
			t.Errorf("verify step must be unconditional: %+v", step)
		}
	}
	wantOrder := []string{"actions/checkout", "actions/setup-go", "actions/setup-node", "Guard Trusted Publishing runtime", "Validate tag", "Publication preflight", "Cross-compile and stage", "Compare published packages", "Publish to npm", "GitHub Release"}
	var gotOrder []string
	for _, step := range release.Steps {
		name := step.Name
		if name == "" {
			name = strings.Split(step.Uses, "@")[0]
		}
		gotOrder = append(gotOrder, name)
		if step.Name == "Verify gate" || step.Name == "Verify every contract" || step.Run == "make verify" || step.Run == "make verify-contracts" {
			t.Errorf("publishing job holds verification step: %+v", step)
		}
	}
	if !reflect.DeepEqual(gotOrder, wantOrder) {
		t.Errorf("release step order = %v, want %v", gotOrder, wantOrder)
	}
}

func TestTheVerificationGateKeepsItsNamesAndItsSteps(t *testing.T) {
	t.Parallel()
	workflow := readVerifyingWorkflow(t, filepath.Join(findRepositoryRoot(t), ".github/workflows/ci-verify.yml"))
	job, exists := workflow.Jobs["verify"]
	if !exists || workflow.Name != "CI | Verify" || job.Name != "Verification gate" || job.Timeout != 60 {
		t.Fatalf("CI identity = %q, job = %q, timeout = %d", workflow.Name, job.Name, job.Timeout)
	}
	want := []string{"actions/checkout", "actions/setup-go", "Align GOCACHE with the restored cache", "Merge the current base branch", "Verify changed", "Verify", "Verify docs", "Verify every contract", "Check the suite-time budget"}
	var got []string
	for _, step := range job.Steps {
		name := step.Name
		if name == "" {
			name = strings.Split(step.Uses, "@")[0]
		}
		got = append(got, name)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CI steps = %v, want %v", got, want)
	}
}
