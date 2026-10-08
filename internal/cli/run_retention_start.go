package cli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func runRetentionAtStart(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, stderr io.Writer) {
	last, found, err := runStore.LastRunRetentionSweep(ctx)
	if err == nil && found && last.RetentionDays == loaded.Config.Store.RunRetentionDays && commandDependenciesForContext(ctx).gc.now().Sub(last.CompletedAt) < 24*time.Hour {
		return
	}
	var report runRetentionReport
	if err == nil {
		report, err = sweepRunDatabase(ctx, runStore, loaded, runRetentionOptions{budget: 2 * time.Second})
	}
	switch {
	case err != nil:
		fmt.Fprintf(stderr, "roundfix: warning: Run Retention failed: %v\n", err)
	case report.paused:
		fmt.Fprintln(stderr, "roundfix: Run Retention paused after its 2s budget; it continues at the next Run start")
	case report.rows.Runs > 0:
		rows := report.rows
		fmt.Fprintf(stderr, "roundfix: Run Retention removed runs=%d rows=%d database_bytes_reclaimed=%d\n", rows.Runs, rows.Runs+rows.RunEvents+rows.AgentSelections+rows.TokenUsage+rows.ActiveRunLocks, report.bytesBefore-report.bytesAfter)
	}
}

func unknownRunMessage(ctx context.Context, id string, retentionDays int) string {
	message := fmt.Sprintf("Run %q does not exist", id)
	timestamp, suffix, ok := strings.Cut(strings.TrimPrefix(id, "run_"), "_")
	if !strings.HasPrefix(id, "run_") || !ok || suffix == "" {
		return message
	}
	for _, char := range suffix {
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f' || char >= 'A' && char <= 'F') {
			return message
		}
	}
	created, err := time.Parse("20060102T150405Z", timestamp)
	if err == nil && created.Before(commandDependenciesForContext(ctx).gc.now().Add(-time.Duration(retentionDays)*24*time.Hour)) {
		message += fmt.Sprintf("; Run Retention may have removed it, because it removes terminal Runs that completed more than %d days ago", retentionDays)
	}
	return message
}
