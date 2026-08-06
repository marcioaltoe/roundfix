package speccheck

import (
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

const (
	// CodeVerifyWorkIndependent identifies a Verification that cannot distinguish Task work from no work.
	CodeVerifyWorkIndependent = "SC-VERIFY-WORK-INDEPENDENT"
)

var (
	makeVerifyPattern = regexp.MustCompile(`^(?:rtk\s+)?(?:env\s+(?:\S+=\S+\s+)*)?make\s+verify$`)
	goWideGatePattern = regexp.MustCompile(`\bgo\s+(?:build|test|vet)\b`)
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

func workingTreeCleanlinessCheck(command string) bool {
	for _, check := range []string{
		"git status --porcelain",
		"git status --short",
		"git diff --check",
		"git diff --quiet",
		"git diff --exit-code",
	} {
		if strings.Contains(command, check) {
			return true
		}
	}
	return false
}
