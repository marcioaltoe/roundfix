// Suite: Verification runner and release dependency contracts.
// Invariant: every verifying job uses macos-26 and publication depends on verification.
// Boundary IN: this repository's workflow YAML.
// Boundary OUT: Baseline assets, GitHub dispatch and network access.
package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// verifyingRunner is the label of every job that runs a test, a vet or a contract.
const verifyingRunner = "macos-26"

type verifyingWorkflowStep struct {
	Name string            `yaml:"name"`
	Uses string            `yaml:"uses"`
	Run  string            `yaml:"run"`
	If   string            `yaml:"if"`
	Env  map[string]string `yaml:"env"`
	With map[string]string `yaml:"with"`
}

type verifyingWorkflow struct {
	Name string               `yaml:"name"`
	On   map[string]yaml.Node `yaml:"on"`
	Jobs map[string]struct {
		Name        string                  `yaml:"name"`
		Runner      yaml.Node               `yaml:"runs-on"`
		Needs       string                  `yaml:"needs"`
		If          string                  `yaml:"if"`
		Timeout     int                     `yaml:"timeout-minutes"`
		Permissions map[string]string       `yaml:"permissions"`
		Steps       []verifyingWorkflowStep `yaml:"steps"`
		Env         map[string]string       `yaml:"env"`
		Strategy    struct {
			FailFast *bool `yaml:"fail-fast"`
			Matrix   struct {
				Part []string `yaml:"part"`
			} `yaml:"matrix"`
		} `yaml:"strategy"`
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

func TestEveryVerifyingJobRunsOnTheStandardMacOSRunner(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(findRepositoryRoot(t), ".github", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	verification := regexp.MustCompile(`\b(?:make[ \t]+verify(?:-changed|-docs|-contracts|-part)?|go[ \t]+(?:test|vet))(?:[\s;|&]|$)`)
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
				if want := verifyingRunner; job.Runner.Kind != yaml.ScalarNode || job.Runner.Value != want {
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
	if verify.Runner.Value != "ubuntu-latest" || !reflect.DeepEqual(verify.Permissions, map[string]string{"contents": "read", "checks": "read", "pull-requests": "read"}) {
		t.Errorf("verify runner/permissions = %s/%v", verify.Runner.Value, verify.Permissions)
	}
	if len(verify.Steps) != 2 {
		t.Fatalf("verify steps = %d, want checkout and provenance", len(verify.Steps))
	}
	checkout := verify.Steps[0]
	if !strings.HasPrefix(checkout.Uses, "actions/checkout@") || !reflect.DeepEqual(checkout.With, map[string]string{"fetch-depth": "0", "persist-credentials": "false"}) {
		t.Errorf("verify checkout = %+v", checkout)
	}
	proof := verify.Steps[1]
	if proof.Name != "Release provenance" || !strings.Contains(proof.Run, "gh api") || proof.Env["GH_TOKEN"] != "${{ github.token }}" {
		t.Errorf("missing token-bound provenance: %+v", proof)
	}
	for name, job := range workflow.Jobs {
		for _, step := range job.Steps {
			if regexp.MustCompile(`\b(make verify(?:-[a-z]+)?|go test|go vet)\b`).MatchString(step.Run) {
				t.Errorf("release job %s runs a test: %s", name, step.Run)
			}
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
	if !exists || workflow.Name != "CI | Verify" || job.Name != "Verification gate" || job.Timeout != 60 || job.Runner.Value != verifyingRunner || job.Needs != "parts" || job.If != "${{ always() }}" {
		t.Fatalf("gathering job = %+v", job)
	}
	if len(job.Steps) != 1 {
		t.Fatalf("gathering steps = %d", len(job.Steps))
	}
	// Execute the real condition for all GitHub dependency outcomes.
	for _, result := range []string{"success", "failure", "cancelled", "skipped", ""} {
		script := strings.ReplaceAll(job.Steps[0].Run, "${{ needs.parts.result }}", result)
		cmd := exec.CommandContext(t.Context(), "sh", "-c", script)
		out, err := cmd.CombinedOutput()
		if (err == nil) != (result == "success") {
			t.Errorf("parts %q: %v, %s", result, err, out)
		}
	}
	parts, exists := workflow.Jobs["parts"]
	if !exists || parts.Runner.Value != verifyingRunner || parts.Timeout != 60 || parts.Strategy.FailFast == nil || *parts.Strategy.FailFast || parts.Env["GOFLAGS"] != "-count=1" {
		t.Fatalf("parts job = %+v", parts)
	}
	wantParts := []string{"packages", "cli-1", "cli-2", "contracts-repo", "contracts-docs"}
	if !reflect.DeepEqual(parts.Strategy.Matrix.Part, wantParts) {
		t.Errorf("matrix = %v", parts.Strategy.Matrix.Part)
	}
	want := []string{"actions/checkout", "actions/setup-go", "Align GOCACHE with the restored cache", "Merge the current base branch", "Verify part"}
	var got []string
	for _, step := range parts.Steps {
		name := step.Name
		if name == "" {
			name = strings.Split(step.Uses, "@")[0]
		}
		got = append(got, name)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parts steps = %v, want %v", got, want)
	}
	if parts.Steps[len(parts.Steps)-1].Run != "make verify-part PART=${{ matrix.part }}" {
		t.Error("matrix does not execute its part")
	}
}

func TestContinuousIntegrationRunsOnlyForAPullRequest(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	workflow := readVerifyingWorkflow(t, filepath.Join(root, ".github/workflows/ci-verify.yml"))
	if len(workflow.On) != 1 {
		t.Fatalf("triggers = %v", workflow.On)
	}
	var trigger struct {
		Branches []string `yaml:"branches"`
		Types    []string `yaml:"types"`
	}
	node, ok := workflow.On["pull_request"]
	if !ok {
		t.Fatal("missing pull_request trigger")
	}
	if err := node.Decode(&trigger); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(trigger.Branches, []string{"main"}) || !reflect.DeepEqual(trigger.Types, []string{"opened", "reopened", "synchronize"}) {
		t.Fatalf("trigger = %+v", trigger)
	}
	// The executable gate surfaces retain only the existing 20m contract timeout.
	timeout := regexp.MustCompile(`-timeout(?:[ =]+)([0-9]+[a-z]+)`)
	for _, path := range []string{"Makefile", ".github/workflows/ci-verify.yml", ".github/workflows/release.yml"} {
		b, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		text := string(b)
		for _, retired := range []string{"SUITE_BUDGET", "test-budget", "suite-budget-check", "docs-test-skipped", "VERIFY_DOCS_TEST"} {
			if strings.Contains(text, retired) {
				t.Errorf("%s keeps %s", path, retired)
			}
		}
		for _, match := range timeout.FindAllStringSubmatch(text, -1) {
			d, err := time.ParseDuration(match[1])
			if err != nil || d > 20*time.Minute || path != "Makefile" {
				t.Errorf("%s: forbidden timeout %s", path, match[0])
			}
		}
	}
}
