package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

// deliveryRunLines reads Active Runs independently of the Queue Owner record.
func deliveryRunLines(ctx context.Context, runStore *store.Store, gitRoot string, queue store.DeliveryQueue) ([]string, error) {
	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{GitRoot: gitRoot, States: store.StatesActive})
	if err != nil {
		return nil, fmt.Errorf("read delivery Active Runs: %w", err)
	}
	var lines []string
	for _, item := range queue.Items {
		if item.Stage != store.DeliveryStageRunning {
			continue
		}
		var active *store.Run
		for i := range runs {
			run := &runs[i]
			if run.Kind == store.KindImplement && run.SpecSlug == item.SpecSlug && run.LocalBranch == item.Branch {
				active = run
				break
			}
		}
		if active == nil {
			lines = append(lines, fmt.Sprintf("Run: %s no Active Run", item.SpecSlug))
			continue
		}
		events, err := runStore.RunEventsOfKinds(ctx, active.ID, runevent.KindDaemonTask)
		if err != nil {
			return nil, fmt.Errorf("read delivery Tasks in progress: %w", err)
		}
		started := map[string]int{}
		for i, event := range events {
			var task struct {
				Task  string `json:"task"`
				Phase string `json:"phase"`
			}
			if err := json.Unmarshal(event.Event.Payload, &task); err != nil {
				return nil, fmt.Errorf("read Task event of Run %s: %w", active.ID, err)
			}
			if task.Task == "" {
				continue
			}
			switch task.Phase {
			case "started":
				started[task.Task] = i
			case "settled", "skipped":
				delete(started, task.Task)
			}
		}
		tasks := make([]string, 0, len(started))
		for task := range started {
			tasks = append(tasks, task)
		}
		sort.Slice(tasks, func(i, j int) bool { return started[tasks[i]] < started[tasks[j]] })
		progress := "no Task in progress"
		if len(tasks) > 0 {
			progress = "Tasks in progress: " + strings.Join(tasks, ", ")
		}
		lines = append(lines, fmt.Sprintf("Run: %s %s; %s", item.SpecSlug, active.ID, progress))
	}
	return lines, nil
}
