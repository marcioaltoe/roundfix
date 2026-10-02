package cli

import (
	"context"
	"strings"
	"time"
)

type readinessFinding struct {
	Code   string
	Status CheckStatus
	Text   string
	Next   string
}

func readinessResult(name, okDetail string, findings []readinessFinding) CheckResult {
	result := CheckResult{Name: name, Status: CheckStatusOK}
	var details, actions []string
	if okDetail != "" {
		details = append(details, okDetail)
	}
	for _, finding := range findings {
		if finding.Status == CheckStatusFailed || (finding.Status == CheckStatusWarn && result.Status == CheckStatusOK) {
			result.Status = finding.Status
		}
		details = append(details, finding.Code+": "+finding.Text)
		if finding.Next != "" {
			actions = append(actions, finding.Next)
		}
	}
	result.Detail = strings.Join(details, "; ")
	result.NextAction = strings.Join(actions, " && ")
	return result
}

type readinessRunner func(ctx context.Context, dir string, env []string, name string, args ...string) (stdout, stderr string, exitCode int, err error)

type readinessDependencies struct {
	run     readinessRunner
	resolve func(name string) (string, error)
	environ []string
	exists  func(path string) bool
	timeout time.Duration
}
