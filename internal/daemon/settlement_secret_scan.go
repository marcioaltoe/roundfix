package daemon

import (
	"context"
	"fmt"
	"strings"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func (engine *Engine) secretScanSettlementCheck(plan TaskPlan, task spec.Task, ordinal, attempt int, before []string) verificationCheck {
	return verificationCheck{Label: runevent.SettlementCheckSecretScan, HaltOnError: true, Run: func(ctx context.Context, _ string) (string, error) {
		prepared, err := engine.prepareTaskCommit(ctx, plan, task, before)
		if err != nil {
			return "", err
		}
		matches, err := scanTaskCommitPaths(ctx, plan.WorkDir, prepared.stageable)
		if err != nil {
			return "", err
		}
		if len(matches) == 0 {
			return "", nil
		}
		details := make([]string, len(matches))
		locations := make([]map[string]any, len(matches))
		for index, match := range matches {
			details[index] = match.String()
			locations[index] = map[string]any{"path": match.Path, "line": match.Line, "kind": match.Kind}
		}
		message := fmt.Sprintf("Task %s commit refused: %d Secret Match(es)", task.ID, len(matches))
		if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonTask,
			message, map[string]any{"phase": "secret_scan", "attempt": attempt, "matches": locations}); err != nil {
			return "", err
		}
		fmt.Fprintln(engine.deps.Progress, message)
		return "Task commit carries a secret: " + strings.Join(details, "; ") + "\n" +
			"Fix: replace each named value with one no Secret Pattern matches, such as letters only or fewer than 20 characters; copying a value the repository already tracks is refused the same way.", nil
	}}
}
