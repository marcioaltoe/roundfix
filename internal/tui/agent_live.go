package tui

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/rounds"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type AgentLiveStream struct {
	output io.Writer
	live   bool
	prog   *tea.Program
	done   chan error

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

func (stream *AgentLiveStream) Write(payload []byte) (int, error) {
	if !stream.live || stream.prog == nil {
		return stream.output.Write(payload)
	}
	stream.prog.Send(agentLiveRawMsg{text: string(payload)})
	return len(payload), nil
}

func (stream *AgentLiveStream) HandleAgentUpdate(update agent.StreamUpdate) {
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
	text := agentTextFromFormat(update)
	if text == "" {
		return ""
	}
	return text
}

func agentTextFromFormat(update agent.StreamUpdate) string {
	switch update.Kind {
	case agent.StreamUpdateMessage:
		return update.Text
	case agent.StreamUpdateThought:
		if strings.TrimSpace(update.Text) == "" {
			return ""
		}
		return "THINK " + strings.TrimRight(update.Text, "\r\n") + "\n"
	case agent.StreamUpdateToolStarted:
		title := strings.TrimSpace(update.Title)
		if title == "" {
			title = "tool call"
		}
		if update.ToolID != "" {
			return fmt.Sprintf("[TOOL] %s (%s)\n", title, update.ToolID)
		}
		return fmt.Sprintf("[TOOL] %s\n", title)
	case agent.StreamUpdateToolUpdated:
		lines := []string{}
		if strings.TrimSpace(update.Title) != "" {
			lines = append(lines, update.Title)
		}
		if strings.TrimSpace(update.ToolState) != "" {
			lines = append(lines, update.ToolState)
		}
		if strings.TrimSpace(update.Text) != "" {
			lines = append(lines, strings.TrimRight(update.Text, "\r\n"))
		}
		if len(lines) == 0 {
			return ""
		}
		return strings.Join(lines, "\n") + "\n"
	case agent.StreamUpdatePlan:
		if strings.TrimSpace(update.Text) == "" {
			return ""
		}
		return "PLAN\n" + strings.TrimRight(update.Text, "\r\n") + "\n"
	case agent.StreamUpdateStatus:
		if update.Status == "" {
			return ""
		}
		return "SESSION " + strings.ToUpper(update.Status) + "\n"
	case agent.StreamUpdateRaw:
		return update.Text
	default:
		return update.Text
	}
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
	lines = append(lines, styleLime.Render("[ : batch_001")+padRightDisplay("", max(width-22, 0))+styleLime.Render(elapsed.String()))
	files := countIssueFiles(view.Issues)
	lines = append(lines, styleMuted.Render(fmt.Sprintf("  RUNNING · FILES %d · ISSUES %d", files, len(view.Issues))))
	lines = append(lines, "")
	for _, group := range GroupIssuesByRound(view.Issues) {
		lines = append(lines, styleSection.Render(fmt.Sprintf("ROUND %03d", group.Round)))
		for _, issue := range group.Issues {
			lines = append(lines, issueLine(issue, width-4))
			if title := strings.TrimSpace(issue.Title); title != "" {
				lines = append(lines, "  "+styleMuted.Render(truncateDisplay(title, width-6)))
			}
		}
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

func issueLine(issue rounds.Issue, width int) string {
	location := emptyDash(issue.File)
	if issue.Line > 0 {
		location = fmt.Sprintf("%s:%d", location, issue.Line)
	}
	line := fmt.Sprintf("%-7s %-9s %s", emptyDash(issue.Severity), emptyDash(issue.Status), location)
	return truncateDisplay(line, width)
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
	styleSection      = lipgloss.NewStyle().Foreground(lipgloss.Color("70")).Bold(true)
	styleTool         = lipgloss.NewStyle().Foreground(lipgloss.Color("48"))
	styleBar          = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238"))
	styleFooter       = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Background(lipgloss.Color("234"))
	styleBorder       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("238")).Padding(0, 1)
	styleActiveBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("10")).Padding(0, 1)
)
