package tui

import (
	"strings"
	"time"

	"roundfix/internal/store"
)

// ShortRunID displays the first eight characters of a Run id's random suffix.
// Older or shorter identifiers remain intact.
func ShortRunID(id string) string {
	suffix := id[strings.LastIndex(id, "_")+1:]
	if strings.HasPrefix(id, "run_") && len(suffix) >= 8 {
		return suffix[:8]
	}
	return id
}

// RunTableOptions controls the human listing without changing the stable rows.
type RunTableOptions struct {
	Now      time.Time
	Location *time.Location
	Repo     bool
	Wide     bool
	Width    int // Zero means no width budget; Wide always shows every cell in full.
}

// RenderRunTable renders Runs in the caller's order with local starts and bare
// durations. Only TARGET may shrink, keeping at least twelve display columns.
func RenderRunTable(runs []store.Run, options RunTableOptions) string {
	if len(runs) == 0 {
		return "No Runs found.\n"
	}
	location := options.Location
	if location == nil {
		location = time.Local
	}
	now := options.Now.In(location)
	header := []string{"RUN"}
	if options.Repo {
		header = append(header, "REPO")
	}
	header = append(header, "STATE")
	if options.Wide {
		header = append(header, "KIND")
	}
	targetColumn := len(header)
	header = append(header, "TARGET", "AGENT", "STARTED ("+now.Format("MST")+")", "DURATION")
	if options.Wide {
		header = append(header, "BRANCH", "WORKTREE")
	}
	note := false
	for _, run := range runs {
		if run.OwnerIdentityUnproven {
			note = true
			break
		}
	}
	if note {
		header = append(header, "NOTE")
	}
	rows := [][]string{header}
	for _, run := range runs {
		id := ShortRunID(run.ID)
		if options.Wide {
			id = run.ID
		}
		row := []string{id}
		if options.Repo {
			row = append(row, RunRepositoryName(run))
		}
		row = append(row, run.State)
		if options.Wide {
			row = append(row, run.Kind)
		}
		target := runRowTarget(run)
		if run.Kind == store.KindImplement {
			target = run.SpecSlug
		}
		started := run.CreatedAt.In(location)
		startFormat := "Jan 02 15:04"
		if started.Year() == now.Year() && started.YearDay() == now.YearDay() {
			startFormat = "15:04"
		}
		duration := runRowEmptyField
		if !store.IsTerminalState(run.State) {
			duration = formatRunElapsed(options.Now.Sub(run.CreatedAt))
		} else if run.CompletedAt != nil {
			duration = formatRunElapsed(run.CompletedAt.Sub(run.CreatedAt))
		}
		row = append(row, runRowField(target), runRowField(run.Agent), started.Format(startFormat), duration)
		if options.Wide {
			row = append(row, runRowField(run.LocalBranch), runRowField(run.GitRoot))
		}
		if note {
			value := ""
			if run.OwnerIdentityUnproven {
				value = "owner identity unproven"
			}
			row = append(row, value)
		}
		rows = append(rows, row)
	}
	widths := make([]int, len(header))
	for _, row := range rows {
		for column, cell := range row {
			widths[column] = max(widths[column], displayWidth(cell))
		}
	}
	if options.Width > 0 && !options.Wide {
		total := 2 * (len(widths) - 1)
		for _, width := range widths {
			total += width
		}
		if excess := total - options.Width; excess > 0 {
			widths[targetColumn] = max(12, widths[targetColumn]-excess)
			for _, row := range rows {
				row[targetColumn] = truncateDisplay(row[targetColumn], widths[targetColumn])
			}
		}
	}
	var table strings.Builder
	for _, row := range rows {
		cells := make([]string, len(row))
		for column, cell := range row {
			cells[column] = padRightDisplay(cell, widths[column])
		}
		table.WriteString(strings.TrimRight(strings.Join(cells, "  "), " "))
		table.WriteByte('\n')
	}
	return table.String()
}
