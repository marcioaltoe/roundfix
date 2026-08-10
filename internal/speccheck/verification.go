package speccheck

import (
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

const (
	// CodeVerifyWorkIndependent identifies a Verification that cannot distinguish Task work from no work.
	CodeVerifyWorkIndependent = "SC-VERIFY-WORK-INDEPENDENT"
	// CodeVerifyVacuousCommand identifies one Verification command that already
	// passes against the unchanged tree, whatever its siblings prove.
	CodeVerifyVacuousCommand = "SC-VERIFY-VACUOUS-COMMAND"
)

var (
	makeVerifyPattern = regexp.MustCompile(`^(?:rtk\s+)?(?:env\s+(?:\S+=\S+\s+)*)?make\s+verify$`)
	goWideGatePattern = regexp.MustCompile(`\bgo\s+(?:build|test|vet)\b`)
	// An assertion over which paths changed.
	workingTreeStatePattern = regexp.MustCompile(`git\s+(?:status|diff)\b[^|&;]*(?:--porcelain|--short|--check|--quiet|--exit-code|--name-only|--name-status|--stat)`)
	// A predicate that succeeds on empty output, which is what an unchanged
	// tree produces. Naming the git flags is not enough: `git diff --name-only
	// | grep -q .` reads the same paths and exits 1 on an unchanged tree, so it
	// fails rather than passing and is not vacuous.
	emptyOutputSucceedsPattern = regexp.MustCompile(`(?:\|\|\s*exit\s+0|\btest\s+-z\b|\[\s+-z\s|--quiet\b|--exit-code\b|"\$\{?\w+\}?"\s*=\s*"\$\{?\w+\}?")`)
)

// WorkIndependentVerification reports a Task whose declared Verification
// cannot distinguish the Task having run from the Task having done nothing.
func WorkIndependentVerification(task spec.Task) (Finding, bool) {
	if len(task.Verification) == 0 {
		return Finding{}, false
	}
	for _, command := range task.Verification {
		if !workIndependentVerificationCommand(command) {
			return Finding{}, false
		}
	}

	return Finding{
		Code:     CodeVerifyWorkIndependent,
		Severity: SeverityError,
		Summary:  task.File + " declares only repository-wide gates and working-tree cleanliness checks, so its Verification cannot distinguish Task work from no work",
		Where:    []Location{{Path: task.File, Line: 1}},
		Fix:      "Add a declared Verification command that asserts this Task's own effect.",
	}, true
}

func workIndependentVerificationCommand(command string) bool {
	normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
	return repositoryWideGate(normalized) || workingTreeCleanlinessCheck(normalized)
}

func repositoryWideGate(command string) bool {
	if makeVerifyPattern.MatchString(command) {
		return true
	}
	if !goWideGatePattern.MatchString(command) || !strings.Contains(command, "./...") {
		return false
	}
	return !strings.Contains(command, " -run ") && !strings.Contains(command, " -run=")
}

// workingTreeCleanlinessCheck reports a command that asserts over the set of
// changed paths and passes when that set is empty, which is what an unchanged
// tree produces. Naming the git flags is not enough: `git diff --name-only |
// grep -q .` reads the same paths and exits 1 on an unchanged tree, so it
// fails rather than passing and is honest Verification.
func workingTreeCleanlinessCheck(command string) bool {
	if !workingTreeStatePattern.MatchString(command) {
		return false
	}
	if !strings.ContainsAny(command, "|&;") {
		// Nothing consumes the output: git prints nothing and exits zero on an
		// unchanged tree, so the command passes before any work happens.
		return true
	}
	// Something does consume it, and whether the command passes depends on
	// what: `| grep -q .` fails on empty output, `|| exit 0` and `test -z`
	// succeed on it.
	return emptyOutputSucceedsPattern.MatchString(command)
}

// VacuousVerificationCommands reports each declared command that already passes
// against the unchanged tree. WorkIndependentVerification judges the Task as a
// whole and stays silent when one honest command sits beside a vacuous one;
// the Daemon's pre-work probe judges each command on its own and refuses the
// Task for any single one, so this check applies the Daemon's unit.
func VacuousVerificationCommands(task spec.Task) []string {
	var vacuous []string
	for _, command := range task.Verification {
		normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
		if workingTreeCleanlinessCheck(normalized) {
			vacuous = append(vacuous, command)
		}
	}
	return vacuous
}
