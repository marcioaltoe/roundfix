package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/rounds"

	"charm.land/lipgloss/v2"
)

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func renderAgentHeader(view LiveRunView, width int) string {
	left := styleAccent.Bold(true).Render("ROUNDFIX") + styleMuted.Render(" // LIVE RUN VIEW")
	right := styleBright.Render(shortRunID(view.RunID))
	return padRightDisplay(left, max(width-displayWidth(right), 1)) + right
}

// shortRunID keeps the header compact: the trailing hash identifies the Run.
func shortRunID(runID string) string {
	if runID == "" {
		return "RUN"
	}
	if len(runID) > 8 {
		return "RUN " + runID[len(runID)-8:]
	}
	return "RUN " + runID
}

func renderPipelineBar(view LiveRunView, width int) string {
	label := styleAccent.Bold(true).Render("SYS.PIPELINE")
	state := strings.ToUpper(emptyDash(view.PipelineState))
	barWidth := max(width-displayWidth("SYS.PIPELINE ")-2, 8)
	bar := styleBar.Render(padRightDisplay(state, barWidth))
	return label + " " + bar
}

func renderAgentSidebar(view LiveRunView, startedAt time.Time, width int, height int) string {
	lines := []string{}
	elapsed := time.Since(startedAt).Truncate(time.Second)
	label := "run"
	if view.BatchNumber > 0 {
		label = fmt.Sprintf("batch_%03d", view.BatchNumber)
		if view.BatchTotal > 0 {
			label = fmt.Sprintf("batch_%03d/%03d", view.BatchNumber, view.BatchTotal)
		}
	}
	lines = append(lines, styleAccent.Render("[ : "+label)+padRightDisplay("", max(width-displayWidth(label)-11, 0))+styleAccent.Render(elapsed.String()))
	files := countIssueFiles(view.Issues)
	issueCount := len(view.Issues)
	if view.TotalIssues > 0 {
		issueCount = view.TotalIssues
	}
	lines = append(lines, styleMuted.Render(fmt.Sprintf("  RUNNING · FILES %d · ISSUES %d", files, issueCount)))
	lines = append(lines, "")
	for index, issue := range view.Issues {
		lines = append(lines, issueJobLines(issue, index+1, view.BatchNumber == index+1, elapsed, width-4)...)
	}
	if len(view.Issues) == 0 {
		lines = append(lines, styleMuted.Render("No Review Issues"))
	}
	return panel(width, height, strings.Join(limitLines(lines, height-2), "\n"), true)
}

func renderAgentTimeline(view LiveRunView, lines []string, width int, height int) string {
	header := styleAccent.Bold(true).Render("SESSION.TIMELINE")
	model := strings.TrimSpace(view.Model)
	if model == "" {
		model = "auto"
	}
	meta := styleMuted.Render(fmt.Sprintf("%d entries · %s · %s", len(lines), emptyDash(view.Agent), model))
	contentLines := []string{header, meta, ""}
	if len(lines) == 0 {
		contentLines = append(contentLines, styleMuted.Render("Waiting for Agent output..."))
	} else {
		contentLines = append(contentLines, colorTimelineLines(lines, width-4)...)
	}
	return panel(width, height, strings.Join(limitTail(contentLines, height-2), "\n"), false)
}

func renderAgentFooter(width int) string {
	text := "FOCUS JOBS   [J/K] ISSUE   [TAB] FOCUS   [CTRL-C] STOP"
	return styleFooter.Render(padRightDisplay(text, width))
}

func issueJobLines(issue rounds.Issue, fallbackNumber int, active bool, elapsed time.Duration, width int) []string {
	name := issueDisplayName(issue, fallbackNumber)
	severity := emptyDash(issue.Severity)
	status := strings.ToUpper(emptyDash(issue.Status))
	duration := "--"
	if active {
		duration = elapsed.String()
		if status == strings.ToUpper(rounds.StatusPending) {
			status = "RUNNING"
		}
	}
	label := fmt.Sprintf("%s • %s", name, severity)
	if active {
		label = styleAccent.Render("▌ ") + truncateDisplay(label, width-2)
	} else {
		label = "  " + truncateDisplay(label, width-2)
	}
	return []string{
		label,
		styleMuted.Render(truncateDisplay("  "+fmt.Sprintf("%s • %s", status, duration), width)),
	}
}

func issueDisplayName(issue rounds.Issue, fallbackNumber int) string {
	base := filepath.Base(issue.Path)
	if strings.HasPrefix(base, "issue_") && strings.HasSuffix(base, ".md") {
		number := strings.TrimSuffix(strings.TrimPrefix(base, "issue_"), ".md")
		if number != "" {
			return "Issue " + number
		}
	}
	if fallbackNumber > 0 {
		return fmt.Sprintf("Issue %03d", fallbackNumber)
	}
	return "Issue"
}

func colorTimelineLines(lines []string, width int) []string {
	colored := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r\n")
		switch {
		case strings.HasPrefix(trimmed, "[TOOL]"):
			colored = append(colored, styleTool.Render(truncateDisplay(trimmed, width)))
		case strings.HasPrefix(trimmed, "PLAN"):
			colored = append(colored, styleAccent.Render(truncateDisplay(trimmed, width)))
		case strings.HasPrefix(trimmed, "THINK"):
			colored = append(colored, styleMuted.Render(truncateDisplay(trimmed, width)))
		case strings.HasPrefix(trimmed, "SESSION"):
			colored = append(colored, styleAccent.Render(truncateDisplay(trimmed, width)))
		default:
			colored = append(colored, styleBright.Render(truncateDisplay(trimmed, width)))
		}
	}
	return colored
}

func panel(width int, height int, content string, active bool) string {
	border := styleBorder
	if active {
		border = styleActiveBorder
	}
	return border.Width(width).Height(height).Render(content)
}

func limitLines(lines []string, maxLines int) []string {
	if maxLines <= 0 || len(lines) <= maxLines {
		return lines
	}
	return lines[:maxLines]
}

func limitTail(lines []string, maxLines int) []string {
	if maxLines <= 0 || len(lines) <= maxLines {
		return lines
	}
	return lines[len(lines)-maxLines:]
}

func countIssueFiles(issues []rounds.Issue) int {
	seen := map[string]struct{}{}
	for _, issue := range issues {
		if strings.TrimSpace(issue.File) == "" {
			continue
		}
		seen[issue.File] = struct{}{}
	}
	return len(seen)
}

func padRightDisplay(value string, width int) string {
	if width <= 0 {
		return ""
	}
	size := displayWidth(value)
	if size >= width {
		return value
	}
	return value + strings.Repeat(" ", width-size)
}

func truncateDisplay(value string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= width {
		return value
	}
	if width <= 1 {
		return string(runes[:width])
	}
	return string(runes[:width-1]) + "…"
}

func displayWidth(value string) int {
	return len([]rune(stripANSI(value)))
}

func stripANSI(value string) string {
	var builder strings.Builder
	inEscape := false
	for _, r := range value {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

// Roundfix uses a blue accent palette so the cockpit reads distinctly from
// other terminal tools.
var (
	styleAccent       = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	styleBright       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	styleMuted        = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleTool         = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	styleBar          = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238"))
	styleBarFill      = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("27")).Bold(true)
	styleBarRest      = lipgloss.NewStyle().Foreground(lipgloss.Color("17")).Background(lipgloss.Color("153"))
	styleError        = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	styleFooter       = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Background(lipgloss.Color("234"))
	styleBorder       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("238")).Padding(0, 1)
	styleActiveBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("39")).Padding(0, 1)
)
