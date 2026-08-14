// Suite: work-independent Task Verification classification
// Invariant: Tasks whose declared commands are all repository gates or clean-tree checks are classified independently of prose or status.
// Boundary IN: parsed Task Verification commands and the public detector
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestWorkIndependentVerificationRefusesOnlyWorkIndependentCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       spec.Status
		verification []string
		wantFinding  bool
	}{
		{
			name:         "repository gate and clean tree",
			status:       spec.StatusPending,
			verification: []string{"make verify", "git status --porcelain"},
			wantFinding:  true,
		},
		{
			name:         "repository wide Go gates and diff cleanliness",
			status:       spec.StatusInProgress,
			verification: []string{"go build -buildvcs=false ./...", "go test -parallel 16 ./...", "git diff --check"},
			wantFinding:  true,
		},
		{
			name:         "repository gate plus focused effect assertion",
			status:       spec.StatusPending,
			verification: []string{"make verify", "go test ./internal/speccheck -run '^TestTaskEffect$'"},
		},
		{
			name:         "repository wide command with effect selector",
			status:       spec.StatusPending,
			verification: []string{"go test ./... -run '^TestTaskEffect$'"},
		},
		{
			name:         "only effect assertions",
			status:       spec.StatusPending,
			verification: []string{"grep -q 'expected effect' internal/speccheck/verification.go", "go test ./internal/speccheck"},
		},
		{
			name:         "status does not change the declared command property",
			status:       spec.StatusCompleted,
			verification: []string{"make verify", "git status --porcelain"},
			wantFinding:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			finding, got := speccheck.WorkIndependentVerification(spec.Task{
				File:         "fixture/task_01.md",
				Status:       tt.status,
				Verification: tt.verification,
			})
			if got != tt.wantFinding {
				t.Fatalf("WorkIndependentVerification() finding = %v, want %v; finding: %#v", got, tt.wantFinding, finding)
			}
			if !got {
				return
			}
			if finding.Code != speccheck.CodeVerifyWorkIndependent || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeVerifyWorkIndependent, speccheck.SeverityError)
			}
			if len(finding.Where) != 1 || finding.Where[0].Path != "fixture/task_01.md" || finding.Where[0].Line < 1 {
				t.Errorf("finding locations = %#v, want the declaring Task", finding.Where)
			}
			if !strings.Contains(finding.Fix, "Task's own effect") {
				t.Errorf("finding fix = %q, want an effect-asserting repair", finding.Fix)
			}
		})
	}
}

func TestVerifyInvertedExit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		command         string
		wantForm        string
		wantReplacement string
	}{
		{
			name:            "grep count uses match status instead of the printed count",
			command:         "grep -c 'PASS' /tmp/task.log",
			wantForm:        "grep -c count-and-exit",
			wantReplacement: "grep -q",
		},
		{
			name:            "filtered count uses wc pipeline status",
			command:         "grep -v '^PASS' /tmp/task.log | wc -l",
			wantForm:        "grep -v ... | wc -l filtered count",
			wantReplacement: `test "$(grep -v ... | wc -l)" -eq 0`,
		},
		{
			name:            "test command substitution has no comparison",
			command:         `test "$(grep -v '^PASS' /tmp/task.log)"`,
			wantForm:        "test $(...) without a comparison",
			wantReplacement: `test -z "$(cmd)"`,
		},
		{
			name:    "grep q uses match status intentionally",
			command: "grep -q 'PASS' /tmp/task.log",
		},
		{
			name:    "test command substitution compares the count",
			command: `test "$(grep -c 'PASS' /tmp/task.log)" -ge 4`,
		},
		{
			name:    "redirected grep count is compared by the terminal test",
			command: `grep -c 'PASS' /tmp/task.log > /tmp/pass-count; test "$(cat /tmp/pass-count)" -ge 4`,
		},
		{
			name:    "filtered count is compared by the enclosing test",
			command: `test "$(grep -v '^PASS' /tmp/task.log | wc -l)" -eq 0`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			findings := speccheck.InvertedExitVerification(spec.Task{
				File:         "fixture/task_01.md",
				Verification: []string{tt.command},
			})
			if tt.wantForm == "" {
				if len(findings) != 0 {
					t.Fatalf("InvertedExitVerification(%q) = %#v, want no finding", tt.command, findings)
				}
				return
			}
			if len(findings) != 1 {
				t.Fatalf("InvertedExitVerification(%q) = %#v, want one finding", tt.command, findings)
			}
			finding := findings[0]
			if finding.Code != speccheck.CodeVerifyInvertedExit || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeVerifyInvertedExit, speccheck.SeverityError)
			}
			if !strings.Contains(finding.Summary, tt.wantForm) {
				t.Errorf("finding summary = %q, want form %q", finding.Summary, tt.wantForm)
			}
			if !strings.Contains(finding.Fix, tt.wantReplacement) {
				t.Errorf("finding fix = %q, want replacement %q", finding.Fix, tt.wantReplacement)
			}
			if len(finding.Where) != 1 || finding.Where[0].Path != "fixture/task_01.md" {
				t.Errorf("finding locations = %#v, want the declaring Task", finding.Where)
			}
		})
	}
}

func TestVerifyInvertedExitSkipsWithoutTaskGraph(t *testing.T) {
	t.Parallel()

	result, err := speccheck.Check(fixtureSpecRoot, "testdata/repo", "no-taskgraph")
	if err != nil {
		t.Fatalf("Check(no-taskgraph): %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeVerifyInvertedExit); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want detector skipped", speccheck.CodeVerifyInvertedExit, findings)
	}
	if !hasSkip(result, speccheck.CodeVerifyInvertedExit, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeVerifyInvertedExit)
	}
}

func TestVacuousVerificationCommandIsCaughtBesideHonestSiblings(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name    string
		command string
		vacuous bool
	}{
		{
			name:    "an assertion over which paths changed",
			command: "git diff --name-only HEAD | grep -v -E '^(kept)$' | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			name:    "a porcelain cleanliness check",
			command: "git status --porcelain | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			// Verified on a clean tree on 2026-08-10: git writes nothing, grep
			// matches nothing, and the command exits 1. It fails on an
			// unchanged tree rather than passing, so it is not vacuous.
			name:    "a nonempty-diff assertion fails on an unchanged tree",
			command: "git diff --name-status HEAD | grep -q .",
			vacuous: false,
		},
		{
			name:    "the same paths with a predicate that passes when empty",
			command: "git diff --name-status HEAD | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			// The equality predicate compares two snapshots captured by earlier
			// chain segments. An intervening command (`make regen`) can change
			// the second snapshot, so `[ "$s1" = "$s2" ]` is not a guaranteed
			// success: its operands are variables, not statically proven
			// values. The matcher must not report it vacuous.
			name:    "a two-snapshot comparison with an intervening command",
			command: `s1="$(git status --porcelain)"; make regen >/dev/null 2>&1; s2="$(git status --porcelain)"; [ "$s1" = "$s2" ]`,
			vacuous: false,
		},
		{
			name:    "a named test that does not exist yet",
			command: "go test ./internal/agent -run '^TestX$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestX'",
			vacuous: false,
		},
		{
			name:    "a grep over declared source",
			command: "grep -c 'Declared break' internal/agent/selection_test.go",
			vacuous: false,
		},
		{
			// git --exit-code does not make the terminal `grep -q .` pass on an
			// unchanged tree; grep still receives no input and exits 1.
			name:    "an exit-code diff piped into grep still fails when empty",
			command: "git diff --exit-code HEAD | grep -q .",
			vacuous: false,
		},
		{
			// A terminal that consumes empty output successfully: cat exits 0
			// on an unchanged tree, so the command passes before any work.
			name:    "a diff piped into cat passes when empty",
			command: "git diff --name-only HEAD | cat",
			vacuous: true,
		},
		{
			// An inverted empty-output check: it asserts changes exist, so on
			// an unchanged tree the terminal `exit 1` fires and fails honestly.
			name:    "an inverted empty-output check fails when empty",
			command: "git diff --name-status HEAD | grep -q . && exit 1",
			vacuous: false,
		},
		{
			// A success form inside another command's quoted argument is not a
			// terminal predicate: grep finds no lines in an unchanged tree's
			// empty output and exits 1, whatever text it was asked to match.
			name:    "a quoted success form inside grep still fails when empty",
			command: "git diff --name-only HEAD | grep -q 'exit 0'",
			vacuous: false,
		},
		{
			// `&&` short-circuits: grep -q . fails on an unchanged tree, so
			// cat is skipped and never decides the exit status, which stays 1.
			name:    "an and-chain whose tail is skipped when empty fails honestly",
			command: "git diff --name-only HEAD | grep -q . && cat",
			vacuous: false,
		},
		{
			// `||` runs its right side when the left fails, so the terminal cat
			// does execute on an unchanged tree and is genuinely vacuous.
			name:    "an or-chain whose tail runs when the left fails passes when empty",
			command: "git diff --name-only HEAD | grep -q . || cat",
			vacuous: true,
		},
		{
			// A quoted test -z form is likewise a grep of ordinary text, not an
			// empty-string test, so it fails on an unchanged tree.
			name:    "a quoted empty-test inside grep still fails when empty",
			command: `git diff --name-only HEAD | grep -q 'test -z'`,
			vacuous: false,
		},
		{
			// `test -z "$(printf x)"` reads a substitution of fixed non-empty
			// output, not the working tree, so it fails on an unchanged tree.
			name:    "an empty-test over non-empty command substitution fails",
			command: `git diff --name-only HEAD | test -z "$(printf x)"`,
			vacuous: false,
		},
		{
			// Single quotes make the operand literal text, not a substitution;
			// `'$(printf x)'` is a fixed non-empty string that fails the test.
			name:    "a single-quoted command substitution is literal and fails",
			command: `git diff --name-only HEAD | test -z '$(printf x)'`,
			vacuous: false,
		},
		{
			// A working-tree substitution is genuinely vacuous: an unchanged
			// tree yields an empty value from it and the test passes.
			name:    "an empty-test over a working-tree substitution passes when empty",
			command: `git diff --name-only HEAD | test -z "$(git diff --name-only HEAD)"`,
			vacuous: true,
		},
		{
			// The bracket form behaves like the test form.
			name:    "a bracket empty-test over a working-tree substitution passes when empty",
			command: `git diff --name-only HEAD | [ -z "$(git diff --name-only HEAD)" ]`,
			vacuous: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			task := spec.Task{Verification: []string{testCase.command}}
			vacuous := speccheck.VacuousVerificationCommands(task)
			if (len(vacuous) != 0) != testCase.vacuous {
				t.Fatalf("VacuousVerificationCommands(%q) = %#v, want vacuous %v", testCase.command, vacuous, testCase.vacuous)
			}
		})
	}
}

// TestOneHonestCommandDoesNotAbsolveAVacuousSibling pins the unit of judgement:
// the Daemon's pre-work probe refuses a Task for any single vacuous command, so
// the static check must not stay silent because a sibling is honest.
func TestOneHonestCommandDoesNotAbsolveAVacuousSibling(t *testing.T) {
	t.Parallel()

	task := spec.Task{Verification: []string{
		"go test ./internal/agent -run '^TestY$' -count=1 -v 2>&1 | grep -q '^--- PASS: TestY'",
		"test -z \"$(git diff --name-only HEAD)\"",
	}}
	vacuous := speccheck.VacuousVerificationCommands(task)
	if len(vacuous) != 1 || !strings.Contains(vacuous[0], "--name-only") {
		t.Fatalf("VacuousVerificationCommands() = %#v, want only the working-tree assertion", vacuous)
	}
	if _, reported := speccheck.WorkIndependentVerification(task); reported {
		t.Fatal("WorkIndependentVerification judges the whole Task and must stay silent here")
	}
}
