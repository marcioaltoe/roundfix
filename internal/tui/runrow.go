package tui

import (
	"fmt"
	"strings"
	"time"

	"roundfix/internal/store"
)

// runRowEmptyField stands in for a missing Run field so every row keeps the
// stable column count on both discovery surfaces.
const runRowEmptyField = "-"

// FormatRunRow renders one Run's discovery fields, newest-first listing
// order preserved by the caller: run id, state, kind, target, agent, start
// time, duration, and local branch, plus the repository when withRepo is
// set. Relative start time is TUI-only; the text surface always uses
// absolute UTC (RFC 3339). Run ids are never truncated here.
func FormatRunRow(run store.Run, now time.Time, relative bool, withRepo bool) []string {
	fields := []string{
		run.ID,
		run.State,
		run.Kind,
		runRowField(runRowTarget(run)),
		runRowField(run.Agent),
		runRowStart(run, now, relative),
		runRowDuration(run, now),
		runRowField(run.LocalBranch),
	}
	if withRepo {
		fields = append(fields, runRowField(run.GitRoot))
	}
	return fields
}

func runRowField(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return runRowEmptyField
	}
	return value
}

func runRowTarget(run store.Run) string {
	if run.Kind == store.KindImplement {
		if strings.TrimSpace(run.SpecSlug) == "" {
			return ""
		}
		return "spec:" + run.SpecSlug
	}
	if strings.TrimSpace(run.PRNumber) == "" {
		return ""
	}
	return "pr:" + run.PRNumber
}

func runRowStart(run store.Run, now time.Time, relative bool) string {
	if relative {
		return formatRunElapsed(now.Sub(run.CreatedAt)) + " ago"
	}
	return run.CreatedAt.UTC().Format(time.RFC3339)
}

// runRowDuration renders how long the Run ran: `running <elapsed>` against
// the wall clock for non-terminal Runs, the completion span for terminal
// ones. A terminal Run without a completion time renders the empty field.
func runRowDuration(run store.Run, now time.Time) string {
	if !store.IsTerminalState(run.State) {
		return "running " + formatRunElapsed(now.Sub(run.CreatedAt))
	}
	if run.CompletedAt == nil {
		return runRowEmptyField
	}
	return formatRunElapsed(run.CompletedAt.Sub(run.CreatedAt))
}

// formatRunElapsed truncates to the largest useful unit: `53s` under a
// minute, `42m` under an hour, `1h12m` beyond. Negative spans (clock skew)
// clamp to zero; skew is cosmetic only.
func formatRunElapsed(elapsed time.Duration) string {
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed < time.Minute {
		return fmt.Sprintf("%ds", int(elapsed.Seconds()))
	}
	minutes := int(elapsed.Minutes())
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%dh%dm", minutes/60, minutes%60)
}
