package tui

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type AgentLiveStream struct {
	output io.Writer
	live   bool
	prog   *tea.Program
	done   chan error
	events *EventBuffer

	mu     sync.Mutex
	closed bool
}

type agentLiveModel struct {
	view      LiveRunView
	console   StreamBuffer
	startedAt time.Time
	width     int
	height    int
}

type agentLiveUpdateMsg struct {
	update agent.StreamUpdate
}

type agentLiveRawMsg struct {
	text string
}

func NewAgentLiveStream(output io.Writer, view LiveRunView, live bool) *AgentLiveStream {
	if output == nil {
		output = io.Discard
	}
	stream := &AgentLiveStream{
		output: output,
		live:   live,
		done:   make(chan error, 1),
	}
	stream.events = NewEventBuffer(agentEventBufferSize, stream.handleAgentUpdate)
	if !live {
		return stream
	}
	model := agentLiveModel{
		view:      view,
		console:   StreamBuffer{MaxLines: 300},
		startedAt: time.Now(),
		width:     view.Width,
		height:    32,
	}
	for _, line := range view.Console {
		_, _ = model.console.Write([]byte(line + "\n"))
	}
	prog := tea.NewProgram(model, tea.WithOutput(output), tea.WithInput(nil), tea.WithoutSignalHandler())
	stream.prog = prog
	go func() {
		_, err := prog.Run()
		stream.done <- err
		close(stream.done)
	}()
	return stream
}

// Live reports whether the stream renders through a terminal program rather
// than plain text output.
func (stream *AgentLiveStream) Live() bool {
	return stream.live && stream.prog != nil
}

func (stream *AgentLiveStream) Write(payload []byte) (int, error) {
	if !stream.live || stream.prog == nil {
		return stream.output.Write(payload)
	}
	stream.prog.Send(agentLiveRawMsg{text: string(payload)})
	return len(payload), nil
}

// Publish implements the Run Event sink as a best-effort consumer: events
// flow through the bounded buffer so rendering pressure never blocks or
// fails producers.
func (stream *AgentLiveStream) Publish(ctx context.Context, event runevent.RunEvent) error {
	return stream.events.Publish(ctx, event)
}

// DroppedEvents reports Run Events dropped under rendering pressure.
func (stream *AgentLiveStream) DroppedEvents() uint64 {
	return stream.events.DroppedEvents()
}

func (stream *AgentLiveStream) handleAgentUpdate(update agent.StreamUpdate) {
	if !stream.live || stream.prog == nil {
		_, _ = io.WriteString(stream.output, agentText(update))
		return
	}
	stream.prog.Send(agentLiveUpdateMsg{update: update})
}

func (stream *AgentLiveStream) Close() error {
	stream.mu.Lock()
	if stream.closed {
		stream.mu.Unlock()
		return nil
	}
	stream.closed = true
	prog := stream.prog
	done := stream.done
	stream.mu.Unlock()
	stream.events.Close()
	if prog == nil {
		return nil
	}
	prog.Quit()
	select {
	case err, ok := <-done:
		if !ok {
			_, _ = fmt.Fprintln(stream.output)
			return nil
		}
		_, _ = fmt.Fprintln(stream.output)
		return err
	case <-time.After(2 * time.Second):
		_, _ = fmt.Fprintln(stream.output)
		return nil
	}
}

func (model agentLiveModel) Init() tea.Cmd {
	return nil
}

func (model agentLiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = value.Width
		model.height = value.Height
	case agentLiveRawMsg:
		_, _ = model.console.Write([]byte(value.text))
	case agentLiveUpdateMsg:
		_, _ = model.console.Write([]byte(agentText(value.update)))
	}
	return model, nil
}

func (model agentLiveModel) View() tea.View {
	width := model.width
	if width <= 0 {
		width = model.view.Width
	}
	if width < 88 {
		width = 88
	}
	height := model.height
	if height < 24 {
		height = 24
	}
	innerWidth := width - 2
	sidebarWidth := innerWidth * 28 / 100
	if sidebarWidth < 30 {
		sidebarWidth = 30
	}
	if sidebarWidth > 46 {
		sidebarWidth = 46
	}
	timelineWidth := innerWidth - sidebarWidth - 1
	bodyHeight := height - 5
	if bodyHeight < 12 {
		bodyHeight = 12
	}

	content := strings.Join([]string{
		renderAgentHeader(model.view, width),
		renderPipelineBar(model.view, width),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			renderAgentSidebar(model.view, model.startedAt, sidebarWidth, bodyHeight),
			renderAgentTimeline(model.view, model.console.Lines(), timelineWidth, bodyHeight),
		),
		renderAgentFooter(width),
	}, "\n")
	view := tea.NewView(content)
	view.AltScreen = true
	return view
}

func agentText(update agent.StreamUpdate) string {
	return agent.ConsoleText(update)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func renderAgentHeader(view LiveRunView, width int) string {
	left := styleLime.Bold(true).Render("ROUNDFIX") + styleMuted.Render(" // ACP COCKPIT")
	right := styleBright.Render("RUN 0/1")
	return padRightDisplay(left, max(width-displayWidth(right), 1)) + right
}

func renderPipelineBar(view LiveRunView, width int) string {
	label := styleLime.Bold(true).Render("SYS.PIPELINE")
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
	lines = append(lines, styleLime.Render("[ : "+label)+padRightDisplay("", max(width-displayWidth(label)-11, 0))+styleLime.Render(elapsed.String()))
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
	header := styleLime.Bold(true).Render("SESSION.TIMELINE")
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
		label = styleLime.Render("▌ ") + truncateDisplay(label, width-2)
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
			colored = append(colored, styleLime.Render(truncateDisplay(trimmed, width)))
		case strings.HasPrefix(trimmed, "THINK"):
			colored = append(colored, styleMuted.Render(truncateDisplay(trimmed, width)))
		case strings.HasPrefix(trimmed, "SESSION"):
			colored = append(colored, styleLime.Render(truncateDisplay(trimmed, width)))
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

var (
	styleLime         = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleBright       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	styleMuted        = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleTool         = lipgloss.NewStyle().Foreground(lipgloss.Color("48"))
	styleBar          = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238"))
	styleFooter       = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Background(lipgloss.Color("234"))
	styleBorder       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("238")).Padding(0, 1)
	styleActiveBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("10")).Padding(0, 1)
)
