package daemon

import (
	"context"
	"fmt"
	"strings"

	"roundfix/internal/speccheck"
)

func (engine *Engine) specConsistencySettlementCheck(plan TaskPlan) verificationCheck {
	return verificationCheck{Label: "settlement check: spec consistency", Run: func(context.Context, string) (string, error) {
		findings, err := engine.deps.SettlementChecker.RefusingFindings(plan.SpecsRoot, plan.WorkDir, plan.Spec.Slug)
		if err != nil {
			return "", fmt.Errorf("read refusing Spec Consistency findings: %w", err)
		}
		var diagnostics strings.Builder
		for _, finding := range findings {
			reason := speccheck.RefusalReason(finding)
			if _, present := plan.specConsistencyBaseline[reason]; present {
				continue
			}
			fmt.Fprintln(&diagnostics, reason)
			for _, location := range finding.Where {
				fmt.Fprintf(&diagnostics, "%s:%d\n", location.Path, location.Line)
			}
			fmt.Fprintf(&diagnostics, "Fix: %s\n", finding.Fix)
		}
		return strings.TrimSpace(diagnostics.String()), nil
	}}
}
