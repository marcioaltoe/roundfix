package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"roundfix/internal/rounds"
	"roundfix/internal/spec"
	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// CockpitMode selects key ownership: the owning resolve/watch process keeps
// Ctrl-C as Stop Request and has no detach; attach detaches with q/Ctrl-C
// and has no stop key.
type CockpitMode int

const (
	CockpitAttach CockpitMode = iota
	CockpitOwning
)

// CockpitSource is the read-only journal surface the cockpit polls.
// *store.Store satisfies it; per ADR 0009 the cockpit never consumes the
// live sink.
type CockpitSource interface {
	TimelineSource
	DataVersion(ctx context.Context) (int64, error)
	Run(ctx context.Context, runID string) (store.Run, bool, error)
}

// CockpitConfig wires one interactive Live Run View.
type CockpitConfig struct {
	Mode         CockpitMode
	View         LiveRunView
	RunID        string
	Source       CockpitSource
	PollInterval time.Duration
	// OnStop handles Ctrl-C in owning mode (Stop Request). Nil in attach.
	OnStop func()
	// Now overrides the clock for elapsed-time rendering. Nil means time.Now.
	Now func() time.Time
}

const defaultCockpitPollInterval = 250 * time.Millisecond

type cockpitFocus int

const (
	focusTimeline cockpitFocus = iota
	focusIssues
)

type cockpitTickMsg struct{}

type issueDetailView struct {
	issue   rounds.Issue
	missing bool
	lines   []string
	scroll  int
}

type cockpitModel struct {
	ctx      context.Context
	cfg      CockpitConfig
	viewport *TimelineViewport
	now      func() time.Time

	focus    cockpitFocus
	selected int
	issueTop int
	detail   *issueDetailView

	issueStatuses  []string
	taskStatuses   []string
	currentBatch   int
	batchStartedAt time.Time

	runState    string
	terminal    bool
	lastVersion int64

	width  int
	height int
}

// newCockpitModel replays the backlog and primes the model. It is the test
// seam: tests drive Update/View directly with synthetic messages.
func newCockpitModel(ctx context.Context, cfg CockpitConfig) (*cockpitModel, error) {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultCockpitPollInterval
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	model := &cockpitModel{
		ctx:         ctx,
		cfg:         cfg,
		viewport:    NewTimelineViewport(cfg.Source, cfg.RunID, 0, 0),
		now:         now,
		runState:    cfg.View.PipelineState,
		lastVersion: -1,
		width:       maxInt(cfg.View.Width, 88),
		height:      32,
	}
	if store.IsTerminalState(model.runState) {
		model.terminal = true
		model.viewport.SetTerminal()
	}
	model.refreshWorkItems()
	model.viewport.SetHeight(model.bodyHeight())
	if err := model.viewport.Replay(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

// RunCockpit opens the interactive Live Run View and blocks until the user
// leaves it (detach in attach mode, Run end or stop in owning mode).
func RunCockpit(ctx context.Context, output io.Writer, cfg CockpitConfig) error {
	model, err := newCockpitModel(ctx, cfg)
	if err != nil {
		return err
	}
	prog := tea.NewProgram(model, tea.WithOutput(output), tea.WithoutSignalHandler())
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			prog.Quit()
		case <-done:
		}
	}()
	_, err = prog.Run()
	return err
}

func (model *cockpitModel) Init() tea.Cmd {
	return model.scheduleTick()
}

func (model *cockpitModel) scheduleTick() tea.Cmd {
	// The owning cockpit lingers after the Run ends and must still drain
	// late events (the terminal outcome); attach opens after they exist.
	if model.terminal && model.cfg.Mode == CockpitAttach {
		return nil
	}
	return tea.Tick(model.cfg.PollInterval, func(time.Time) tea.Msg {
		return cockpitTickMsg{}
	})
}

func (model *cockpitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = value.Width
		model.height = value.Height
		model.viewport.SetHeight(model.bodyHeight())
		return model, nil
	case cockpitTickMsg:
		model.poll()
		return model, model.scheduleTick()
	case tea.KeyPressMsg:
		return model.handleKey(tea.Key(value))
	}
	return model, nil
}

func (model *cockpitModel) poll() {
	version, err := model.cfg.Source.DataVersion(model.ctx)
	if err != nil || version == model.lastVersion {
		return
	}
	model.lastVersion = version
	_ = model.viewport.Poll(model.ctx)
	run, found, err := model.cfg.Source.Run(model.ctx, model.cfg.RunID)
	if err != nil || !found {
		return
	}
	model.runState = run.State
	if store.IsTerminalState(run.State) {
		model.terminal = true
		model.viewport.SetTerminal()
	}
	model.refreshWorkItems()
}

// specRun reports whether this cockpit renders a spec Run's Tasks; every
// other Run Kind keeps the Review Issue pane unchanged.
func (model *cockpitModel) specRun() bool {
	return specRunView(model.cfg.View)
}

// refreshWorkItems refreshes the work-item pane keyed on the Run Kind:
// Review Issue artifacts for review Runs, task files for spec Runs.
func (model *cockpitModel) refreshWorkItems() {
	if model.specRun() {
		model.refreshTasks()
		return
	}
	model.refreshIssues()
}

// refreshTasks re-reads Task statuses from the task files located through
// the Run row's git root. A parse failure keeps the last good status: the
// Agent rewrites task files while the pane polls, and a mid-write read must
// never fail the view (ADR 0009 keeps the cockpit on the journal plus these
// files; it never consumes the live sink).
func (model *cockpitModel) refreshTasks() {
	tasks := model.cfg.View.Tasks
	if len(model.taskStatuses) != len(tasks) {
		model.taskStatuses = make([]string, len(tasks))
		for index, task := range tasks {
			model.taskStatuses[index] = string(task.Status)
		}
	}
	for index := range tasks {
		current := tasks[index]
		if err := spec.ReloadTask(model.cfg.View.GitRoot, &current); err == nil {
			model.taskStatuses[index] = string(current.Status)
		}
	}
}

// refreshIssues re-reads Review Issue artifact statuses and derives which
// Batch is executing, so the sidebar and the progress bar track the Run.
func (model *cockpitModel) refreshIssues() {
	issues := model.cfg.View.Issues
	if len(model.issueStatuses) != len(issues) {
		model.issueStatuses = make([]string, len(issues))
	}
	for index, issue := range issues {
		status := issue.Status
		if parsed, err := rounds.ParseIssue(issue.Path); err == nil {
			status = parsed.Status
		}
		model.issueStatuses[index] = status
	}
	current := 0
	for index := range issues {
		if model.issueStatuses[index] == rounds.StatusPending || model.issueStatuses[index] == "" {
			current = model.batchOf(index)
			break
		}
	}
	if current != model.currentBatch {
		model.currentBatch = current
		model.batchStartedAt = model.now()
	}
}

// batchOf maps an issue index to its 1-based Batch number; 0 means the plan
// is unknown (attach without batch info).
func (model *cockpitModel) batchOf(index int) int {
	consumed := 0
	for batch, size := range model.cfg.View.BatchSizes {
		consumed += size
		if index < consumed {
			return batch + 1
		}
	}
	return 0
}

// issueStatusLabel renders the execution state per recorded design:
// terminal artifact statuses verbatim, Executing for the current Batch,
// Waiting ahead of it, Paused once the Run itself has ended.
func (model *cockpitModel) issueStatusLabel(index int) string {
	status := model.issueStatuses[index]
	switch status {
	case rounds.StatusResolved, rounds.StatusInvalid, rounds.StatusDuplicated, rounds.StatusFailed:
		return strings.ToUpper(status[:1]) + status[1:]
	}
	if model.terminal || store.IsTerminalState(model.runState) {
		return "Paused"
	}
	if model.currentBatch > 0 && model.batchOf(index) == model.currentBatch {
		return "Executing"
	}
	return "Waiting"
}

// taskStatusLabel mirrors the Review Issue labels: terminal task statuses
// verbatim, Executing for the Task the cycle is on, Waiting ahead of it,
// Paused once the Run itself has ended.
func (model *cockpitModel) taskStatusLabel(index int) string {
	status := spec.Status(model.taskStatuses[index])
	switch status {
	case spec.StatusCompleted:
		return "Completed"
	case spec.StatusFailed:
		return "Failed"
	}
	if model.terminal || store.IsTerminalState(model.runState) {
		return "Paused"
	}
	if status == spec.StatusInProgress || index == model.currentTaskIndex() {
		return "Executing"
	}
	return "Waiting"
}

// currentTaskIndex approximates the executing Task the same way the review
// pane derives the executing Batch: the cycle runs Tasks in Task Graph
// order, so the first unsettled Task is the one in flight.
func (model *cockpitModel) currentTaskIndex() int {
	for index, status := range model.taskStatuses {
		if spec.Status(status) != spec.StatusCompleted && spec.Status(status) != spec.StatusFailed {
			return index
		}
	}
	return -1
}

// workItemCount sizes selection over the pane's Work Items, keyed on the
// Run Kind.
func (model *cockpitModel) workItemCount() int {
	if model.specRun() {
		return len(model.cfg.View.Tasks)
	}
	return len(model.cfg.View.Issues)
}

func (model *cockpitModel) progressCounts() (int, int) {
	if model.specRun() {
		done := 0
		for _, status := range model.taskStatuses {
			if spec.Status(status) == spec.StatusCompleted {
				done++
			}
		}
		return done, len(model.taskStatuses)
	}
	done := 0
	for _, status := range model.issueStatuses {
		switch status {
		case rounds.StatusResolved, rounds.StatusInvalid, rounds.StatusDuplicated:
			done++
		}
	}
	return done, len(model.issueStatuses)
}

func (model *cockpitModel) handleKey(key tea.Key) (tea.Model, tea.Cmd) {
	keystroke := key.String()
	switch keystroke {
	case "q":
		if model.cfg.Mode == CockpitAttach || model.terminal {
			return model, tea.Quit
		}
		return model, nil
	case "ctrl+c":
		if model.cfg.Mode == CockpitAttach || model.terminal {
			return model, tea.Quit
		}
		if model.cfg.OnStop != nil {
			model.cfg.OnStop()
		}
		return model, nil
	case "esc":
		model.detail = nil
		return model, nil
	case "tab":
		if model.detail == nil {
			if model.focus == focusTimeline {
				model.focus = focusIssues
			} else {
				model.focus = focusTimeline
			}
		}
		return model, nil
	}
	if model.detail != nil {
		model.handleDetailKey(keystroke)
		return model, nil
	}
	if model.focus == focusIssues {
		model.handleIssueKey(keystroke)
		return model, nil
	}
	model.handleTimelineKey(keystroke)
	return model, nil
}

func (model *cockpitModel) handleTimelineKey(keystroke string) {
	switch keystroke {
	case "up", "k":
		_ = model.viewport.ScrollUp(model.ctx, 1)
	case "down", "j":
		_ = model.viewport.ScrollDown(model.ctx, 1)
	case "pgup":
		_ = model.viewport.PageUp(model.ctx)
	case "pgdown":
		_ = model.viewport.PageDown(model.ctx)
	case "home":
		model.viewport.JumpToTop()
	case "end", "G":
		_ = model.viewport.JumpToTail(model.ctx)
	}
}

func (model *cockpitModel) handleIssueKey(keystroke string) {
	switch keystroke {
	case "up", "k":
		if model.selected > 0 {
			model.selected--
		}
	case "down", "j":
		if model.selected < model.workItemCount()-1 {
			model.selected++
		}
	case "enter":
		model.openDetail()
	}
}

func (model *cockpitModel) handleDetailKey(keystroke string) {
	detail := model.detail
	switch keystroke {
	case "up", "k":
		if detail.scroll > 0 {
			detail.scroll--
		}
	case "down", "j":
		detail.scroll++
	case "pgup":
		detail.scroll = maxInt(detail.scroll-(model.bodyHeight()-1), 0)
	case "pgdown":
		detail.scroll += model.bodyHeight() - 1
	case "home":
		detail.scroll = 0
	}
	if limit := maxInt(len(detail.lines)-1, 0); detail.scroll > limit {
		detail.scroll = limit
	}
}

// openDetail loads the selected Review Issue artifact read-only. A missing
// or cleaned artifact degrades to a notice, never a failure. The detail
// pane applies to Review Issues only; a spec Run's Task pane has no detail.
func (model *cockpitModel) openDetail() {
	if model.specRun() {
		return
	}
	if model.selected < 0 || model.selected >= len(model.cfg.View.Issues) {
		return
	}
	listed := model.cfg.View.Issues[model.selected]
	detail := &issueDetailView{issue: listed}
	parsed, err := rounds.ParseIssue(listed.Path)
	if err == nil {
		detail.issue = parsed
	}
	content, readErr := os.ReadFile(listed.Path)
	if readErr != nil {
		detail.missing = true
		detail.lines = []string{"artifact not available", listed.Path}
	} else {
		detail.lines = artifactBodyLines(string(content))
	}
	model.detail = detail
}

// artifactBodyLines drops the YAML frontmatter: the detail header already
// shows the structured fields, so the pane renders the markdown body.
func artifactBodyLines(content string) []string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for index := 1; index < len(lines); index++ {
			if strings.TrimSpace(lines[index]) == "---" {
				body := lines[index+1:]
				for len(body) > 0 && strings.TrimSpace(body[0]) == "" {
					body = body[1:]
				}
				return body
			}
		}
	}
	return lines
}

func (model *cockpitModel) bodyHeight() int {
	height := model.height - 5
	if height < 8 {
		height = 8
	}
	return height
}

type cockpitLayout struct {
	width        int
	bodyHeight   int
	sidebarWidth int
	rightWidth   int
}

func cockpitLayoutFor(model *cockpitModel) cockpitLayout {
	width := maxInt(model.width, 88)
	bodyHeight := model.bodyHeight()
	innerWidth := width - 2
	sidebarWidth := innerWidth * 28 / 100
	if sidebarWidth < 30 {
		sidebarWidth = 30
	}
	if sidebarWidth > 46 {
		sidebarWidth = 46
	}
	return cockpitLayout{
		width:        width,
		bodyHeight:   bodyHeight,
		sidebarWidth: sidebarWidth,
		rightWidth:   innerWidth - sidebarWidth - 1,
	}
}

func (model *cockpitModel) View() tea.View {
	layout := cockpitLayoutFor(model)
	view := tea.NewView(renderCockpitLayout(model, layout))
	view.AltScreen = true
	return view
}

func renderCockpitLayout(model *cockpitModel, layout cockpitLayout) string {
	return strings.Join([]string{
		renderCockpitHeaderArea(model, layout.width),
		renderCockpitBody(model, layout),
		renderCockpitFooter(model, layout.width),
	}, "\n")
}

func renderCockpitHeaderArea(model *cockpitModel, width int) string {
	return strings.Join([]string{
		renderAgentHeader(model.cfg.View, width),
		model.renderStatusBar(width),
		"",
	}, "\n")
}

func renderCockpitBody(model *cockpitModel, layout cockpitLayout) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderCockpitWorkQueue(model, layout),
		renderCockpitRightPane(model, layout),
	)
}

func renderCockpitWorkQueue(model *cockpitModel, layout cockpitLayout) string {
	return panel(
		layout.sidebarWidth,
		layout.bodyHeight,
		model.renderWorkItemPane(layout.sidebarWidth, layout.bodyHeight),
		model.focus == focusIssues && model.detail == nil,
	)
}

func renderCockpitRightPane(model *cockpitModel, layout cockpitLayout) string {
	if model.detail != nil {
		return renderCockpitDetailPane(model, layout.rightWidth, layout.bodyHeight)
	}
	return renderCockpitTimelinePane(model, layout.rightWidth, layout.bodyHeight)
}

func renderCockpitTimelinePane(model *cockpitModel, width int, height int) string {
	lines := model.viewport.VisibleLines()
	content := []string{styleAccent.Bold(true).Render("SESSION.TIMELINE"), ""}
	if len(lines) == 0 {
		content = append(content, styleMuted.Render("No Run Events yet..."))
	} else {
		content = append(content, colorTimelineLines(lines, width-4)...)
	}
	return panel(width, height, strings.Join(limitTail(content, height-2), "\n"), model.focus == focusTimeline)
}

func renderCockpitDetailPane(model *cockpitModel, width int, height int) string {
	return panel(width, height, model.renderDetail(width, height), true)
}

func renderCockpitFooter(model *cockpitModel, width int) string {
	return model.renderFooter(width)
}

func (model *cockpitModel) renderDetail(width int, height int) string {
	detail := model.detail
	issue := detail.issue
	header := []string{
		styleAccent.Bold(true).Render("REVIEW.ISSUE"),
		styleBright.Render(emptyDash(issue.Title)),
		styleMuted.Render(fmt.Sprintf("%s · %s · %s:%d", emptyDash(issue.Severity), emptyDash(issue.Status), emptyDash(issue.File), issue.Line)),
		styleMuted.Render("source: " + emptyDash(issue.SourceRef)),
		"",
	}
	visible := detail.lines
	if detail.scroll < len(visible) {
		visible = visible[detail.scroll:]
	} else {
		visible = nil
	}
	body := []string{}
	for _, line := range visible {
		body = append(body, truncateDisplay(line, width-4))
	}
	return strings.Join(limitLines(append(header, body...), height-2), "\n")
}

// renderWorkItemPane renders the sidebar's Work Items keyed on the Run
// Kind: Tasks for spec Runs, Review Issues for everything else.
func (model *cockpitModel) renderWorkItemPane(width int, height int) string {
	if model.specRun() {
		return model.renderTaskPane(width, height)
	}
	return model.renderIssuePane(width, height)
}

// renderTaskPane renders the spec Run's Tasks in Task Graph order with the
// statuses of the latest task-file refresh, reusing the issue pane's
// windowing so selection scrolls identically.
func (model *cockpitModel) renderTaskPane(width int, height int) string {
	items := model.taskWorkItems()
	lines := []string{styleAccent.Bold(true).Render("SPEC.TASKS"), styleMuted.Render(fmt.Sprintf("%d Task(s)", len(items))), ""}
	if len(items) == 0 {
		lines = append(lines, styleMuted.Render("No Tasks"))
		return strings.Join(limitLines(lines, height-2), "\n")
	}
	visible := maxInt((height-5)/3, 1)
	if model.selected < model.issueTop {
		model.issueTop = model.selected
	}
	if model.selected >= model.issueTop+visible {
		model.issueTop = model.selected - visible + 1
	}
	end := minInt(model.issueTop+visible, len(items))
	for index := model.issueTop; index < end; index++ {
		lines = append(lines, model.workItemBlock(items[index], model.taskStatusLabel(index), index, width)...)
	}
	return strings.Join(limitLines(lines, height-2), "\n")
}

// taskWorkItems maps the spec Run's Tasks into Work Items carrying the
// statuses of the latest refresh instead of the load-time ones.
func (model *cockpitModel) taskWorkItems() []WorkItem {
	items := TaskWorkItems(model.cfg.View.Tasks)
	for index := range items {
		if index < len(model.taskStatuses) {
			items[index].Status = model.taskStatuses[index]
		}
	}
	return items
}

func (model *cockpitModel) renderIssuePane(width int, height int) string {
	items := model.issueWorkItems()
	lines := []string{styleAccent.Bold(true).Render("REVIEW.ISSUES"), styleMuted.Render(fmt.Sprintf("%d issue(s)", len(items))), ""}
	if len(items) == 0 {
		lines = append(lines, styleMuted.Render("No Review Issues"))
		return strings.Join(limitLines(lines, height-2), "\n")
	}
	// Each issue renders as a two-line block plus a blank spacer, with a
	// Batch separator ahead of each Batch's first issue.
	visible := maxInt((height-5)/3, 1)
	if model.selected < model.issueTop {
		model.issueTop = model.selected
	}
	if model.selected >= model.issueTop+visible {
		model.issueTop = model.selected - visible + 1
	}
	end := minInt(model.issueTop+visible, len(items))
	for index := model.issueTop; index < end; index++ {
		lines = append(lines, model.issueBlock(items[index], index, width)...)
	}
	return strings.Join(limitLines(lines, height-2), "\n")
}

// issueWorkItems maps the Run's Review Issues into Work Items: positional
// display names plus the statuses of the latest artifact refresh.
func (model *cockpitModel) issueWorkItems() []WorkItem {
	items := make([]WorkItem, len(model.cfg.View.Issues))
	for index, issue := range model.cfg.View.Issues {
		status := ""
		if index < len(model.issueStatuses) {
			status = model.issueStatuses[index]
		}
		items[index] = WorkItem{Name: fmt.Sprintf("Issue #%03d", index+1), Title: issue.Title, Status: status}
	}
	return items
}

func (model *cockpitModel) issueBlock(item WorkItem, index int, width int) []string {
	lines := []string{}
	if separator := model.batchSeparator(index, width); separator != "" {
		lines = append(lines, styleAccent.Render(separator))
	}
	return append(lines, model.workItemBlock(item, model.issueStatusLabel(index), index, width)...)
}

// workItemBlock renders one Work Item as the two-line sidebar block both
// Run Kinds share: the marker plus display name, then the status label.
func (model *cockpitModel) workItemBlock(item WorkItem, label string, index int, width int) []string {
	marker := "  "
	nameStyle := styleMuted
	if index == model.selected {
		marker = "> "
		nameStyle = styleBright
	}
	return []string{
		nameStyle.Render(truncateDisplay(marker+item.Name, width-4)),
		model.statusStyle(label).Render(truncateDisplay("  "+label, width-4)),
		"",
	}
}

// batchSeparator labels the first issue of each Batch when the plan is
// known. One Agent executes the whole Batch, so the elapsed clock lives on
// the executing Batch's separator, not on individual issues.
func (model *cockpitModel) batchSeparator(index int, width int) string {
	total := len(model.cfg.View.BatchSizes)
	if total == 0 {
		return ""
	}
	batch := model.batchOf(index)
	if batch == 0 {
		return ""
	}
	if index > 0 && model.batchOf(index-1) == batch {
		return ""
	}
	separator := fmt.Sprintf("─── Batch %03d/%03d", batch, total)
	if batch == model.currentBatch && !model.terminal && !store.IsTerminalState(model.runState) {
		elapsed := formatElapsed(model.now().Sub(model.batchStartedAt))
		pad := maxInt(width-4-displayWidth(separator)-displayWidth(elapsed), 1)
		separator += strings.Repeat(" ", pad) + elapsed
	}
	return truncateDisplay(separator, width-4)
}

func (model *cockpitModel) statusStyle(label string) lipgloss.Style {
	switch label {
	case "Executing":
		return styleAccent
	case "Resolved", "Invalid", "Duplicated", "Completed":
		return styleTool
	case "Failed":
		return styleError
	default:
		return styleMuted
	}
}

func formatElapsed(elapsed time.Duration) string {
	if elapsed < 0 {
		elapsed = 0
	}
	total := int(elapsed.Seconds())
	if total >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", total/3600, (total%3600)/60, total%60)
	}
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}

// renderStatusBar is the Run progress bar: resolved issues over total,
// filled solid blue to the completed percentage. Scrollback and read-only
// hints surface only when they apply.
func (model *cockpitModel) renderStatusBar(width int) string {
	label := styleAccent.Bold(true).Render("RUN.PROGRESS")
	barWidth := maxInt(width-displayWidth("RUN.PROGRESS ")-2, 8)
	done, total := model.progressCounts()
	text := strings.ToUpper(emptyDash(model.runState))
	percent := 0
	if total > 0 {
		percent = done * 100 / total
		text = fmt.Sprintf(" %d of %d issue(s) resolved · %d%%", done, total, percent)
		if model.specRun() {
			text = fmt.Sprintf(" %d of %d Task(s) completed · %d%%", done, total, percent)
		}
	}
	if suffix := model.statusSuffix(); suffix != "" {
		text += " · " + suffix
	}
	padded := padRightDisplay(truncateDisplay(text, barWidth), barWidth)
	fill := barWidth * percent / 100
	runes := []rune(padded)
	if fill > len(runes) {
		fill = len(runes)
	}
	return label + " " + styleBarFill.Render(string(runes[:fill])) + styleBarRest.Render(string(runes[fill:]))
}

// statusSuffix narrates only the states the user must notice: a frozen
// scrolled viewport, replay in progress, or a finished read-only Run.
func (model *cockpitModel) statusSuffix() string {
	state, below := model.viewport.State()
	switch state {
	case FollowReplaying:
		return "REPLAYING BACKLOG..."
	case FollowScrolled:
		if below > 0 {
			return fmt.Sprintf("SCROLLED · %d new event(s) below — End to follow", below)
		}
		return "SCROLLED — End to follow"
	case FollowTerminal:
		return "READ-ONLY"
	default:
		return ""
	}
}

func (model *cockpitModel) renderFooter(width int) string {
	keys := "Tab focus · ↑↓ scroll · PgUp/PgDn page · End follow · Enter issue · Esc back"
	switch {
	case model.cfg.Mode == CockpitAttach:
		keys += " · q detach"
	case model.terminal:
		keys += " · q close"
	default:
		keys += " · Ctrl-C stop"
	}
	return styleMuted.Render(padRightDisplay("Keys: "+keys, width))
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
