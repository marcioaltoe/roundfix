package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
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
	task    spec.Task
	kind    detailKind
	ordinal int
	missing bool
	stale   bool
	lines   []string
	scroll  int
}

type detailKind int

const (
	detailReviewIssue detailKind = iota
	detailTask
)

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
	model.refreshOpenDetail()
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
			model.cfg.View.Issues[index] = parsed
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

func (model *cockpitModel) allReviewIssuesSettled() bool {
	if len(model.issueStatuses) == 0 {
		return true
	}
	for _, status := range model.issueStatuses {
		switch status {
		case rounds.StatusResolved, rounds.StatusInvalid, rounds.StatusDuplicated:
		default:
			return false
		}
	}
	return true
}

func (model *cockpitModel) allTasksCompleted() bool {
	if len(model.taskStatuses) == 0 {
		return false
	}
	for _, status := range model.taskStatuses {
		if spec.Status(status) != spec.StatusCompleted {
			return false
		}
	}
	return true
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
	case "d":
		if model.detail != nil {
			model.detail = nil
			return model, nil
		}
		model.openDetail()
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
	pageSize := model.detailPageSize()
	switch keystroke {
	case "up", "k":
		if detail.scroll > 0 {
			detail.scroll--
		}
	case "down", "j":
		if detail.scroll < model.detailMaxScroll() {
			detail.scroll++
		}
	case "pgup":
		detail.scroll = maxInt(detail.scroll-pageSize, 0)
	case "pgdown":
		detail.scroll += pageSize
	case "home":
		detail.scroll = 0
	}
	if limit := model.detailMaxScroll(); detail.scroll > limit {
		detail.scroll = limit
	}
}

// openDetail loads the selected Work Item read-only. Review Runs show the
// Review Issue artifact; spec Runs show the Task file body.
func (model *cockpitModel) openDetail() {
	if model.specRun() {
		model.openTaskDetail()
		return
	}
	model.openReviewIssueDetail()
}

// openReviewIssueDetail loads the selected Review Issue artifact read-only.
// A missing or cleaned artifact degrades to a notice, never a failure.
func (model *cockpitModel) openReviewIssueDetail() {
	if model.selected < 0 || model.selected >= len(model.cfg.View.Issues) {
		return
	}
	listed := model.cfg.View.Issues[model.selected]
	detail := &issueDetailView{kind: detailReviewIssue, issue: listed, ordinal: model.selected + 1}
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

func (model *cockpitModel) openTaskDetail() {
	if model.selected < 0 || model.selected >= len(model.cfg.View.Tasks) {
		return
	}
	task := model.cfg.View.Tasks[model.selected]
	if model.selected < len(model.taskStatuses) {
		task.Status = spec.Status(model.taskStatuses[model.selected])
	}
	detail := &issueDetailView{kind: detailTask, task: task, ordinal: model.selected + 1}
	if err := model.loadTaskDetail(detail); err != nil {
		detail.stale = true
		detail.lines = []string{"task file temporarily unreadable", task.File}
	}
	model.detail = detail
}

func (model *cockpitModel) refreshOpenDetail() {
	if model.detail == nil || model.detail.kind != detailTask {
		return
	}
	if err := model.loadTaskDetail(model.detail); err != nil {
		model.detail.stale = true
	}
	if model.detail.scroll > model.detailMaxScroll() {
		model.detail.scroll = model.detailMaxScroll()
	}
}

func (model *cockpitModel) loadTaskDetail(detail *issueDetailView) error {
	task := detail.task
	if err := spec.ReloadTask(model.cfg.View.GitRoot, &task); err != nil {
		return err
	}
	content, err := os.ReadFile(filepath.Join(model.cfg.View.GitRoot, task.File))
	if err != nil {
		return err
	}
	detail.task = task
	detail.lines = artifactBodyLines(string(content))
	detail.stale = false
	return nil
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
	sidebarWidth := innerWidth * 40 / 100
	if sidebarWidth < 30 {
		sidebarWidth = 30
	}
	if sidebarWidth > 46 {
		sidebarWidth = 46
	}
	if innerWidth-sidebarWidth-1 <= sidebarWidth {
		sidebarWidth = maxInt((innerWidth-2)/2, 1)
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
	base := renderCockpitBaseLayout(model, layout)
	if model.detail != nil {
		return renderCockpitDetailOverlay(model, layout, base)
	}
	return base
}

func renderCockpitBaseLayout(model *cockpitModel, layout cockpitLayout) string {
	return strings.Join([]string{
		renderCockpitHeaderArea(model, layout.width),
		renderCockpitBody(model, layout),
		renderCockpitFooter(model, layout.width),
	}, "\n")
}

func renderCockpitHeaderArea(model *cockpitModel, width int) string {
	return strings.Join([]string{
		renderCockpitHeaderLine(model, width),
		renderPhaseRow(width, runPhases(model)),
	}, "\n")
}

func renderCockpitHeaderLine(model *cockpitModel, width int) string {
	target := cockpitTargetLabel(model.cfg.View)
	leftText := "ROUNDFIX // LIVE RUN VIEW"
	if target != "" {
		leftText += " // " + target
	}
	runID := model.cfg.View.RunID
	if runID == "" {
		runID = model.cfg.RunID
	}
	rightText := shortRunID(runID) + " [" + formatRunStateLabel(model.runState) + "]"
	if suffix := model.statusSuffix(); suffix != "" {
		rightText += " " + suffix
	}
	rightWidth := displayWidth(rightText)
	if rightWidth >= width {
		rightText = truncateDisplay(rightText, width-1)
		rightWidth = displayWidth(rightText)
	}
	leftText = truncateDisplay(leftText, maxInt(width-rightWidth, 1))
	left := renderCockpitHeaderLeft(leftText)
	right := styleBright.Render(rightText)
	return padRightDisplay(left, maxInt(width-displayWidth(right), 1)) + right
}

func renderCockpitHeaderLeft(text string) string {
	if strings.HasPrefix(text, "ROUNDFIX") {
		return styleAccent.Bold(true).Render("ROUNDFIX") + styleMuted.Render(strings.TrimPrefix(text, "ROUNDFIX"))
	}
	return styleMuted.Render(text)
}

func cockpitTargetLabel(view LiveRunView) string {
	if specRunView(view) {
		if strings.TrimSpace(view.SpecSlug) != "" {
			return "SPEC " + strings.TrimSpace(view.SpecSlug)
		}
		return "SPEC"
	}
	if strings.TrimSpace(view.PRNumber) != "" {
		return "PR #" + strings.TrimSpace(view.PRNumber)
	}
	return strings.TrimSpace(view.HeadBranch)
}

func formatRunStateLabel(state string) string {
	if state == "" {
		return "UNKNOWN"
	}
	var builder strings.Builder
	previousLowerOrDigit := false
	for index, r := range state {
		upper := r >= 'A' && r <= 'Z'
		if index > 0 && upper && previousLowerOrDigit {
			builder.WriteByte(' ')
		}
		builder.WriteRune(r)
		previousLowerOrDigit = (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
	}
	return strings.ToUpper(builder.String())
}

type cockpitPhase struct {
	name   string
	marker string
}

const (
	phaseDone   = "done"
	phaseRun    = "run"
	phaseWait   = "wait"
	phaseLocked = "locked"
)

func runPhases(model *cockpitModel) []cockpitPhase {
	if model.specRun() {
		return specRunPhases(model)
	}
	return reviewRunPhases(model)
}

func reviewRunPhases(model *cockpitModel) []cockpitPhase {
	pushMarker := phaseLocked
	if model.allReviewIssuesSettled() {
		pushMarker = phaseWait
	}
	phases := []cockpitPhase{
		{name: "FETCH", marker: phaseDone},
		{name: "TRIAGE", marker: phaseDone},
		{name: "AGENT", marker: phaseWait},
		{name: "VERIFY", marker: phaseWait},
		{name: "PUSH", marker: pushMarker},
	}
	switch model.runState {
	case store.StatePushing:
		phases[2].marker = phaseDone
		phases[3].marker = phaseDone
		phases[4].marker = phaseRun
	case store.StateVerifying:
		phases[2].marker = phaseDone
		phases[3].marker = phaseRun
	case store.StateClean:
		for index := range phases {
			phases[index].marker = phaseDone
		}
	case store.StateUnresolved:
		phases[2].marker = phaseDone
		phases[3].marker = phaseDone
		phases[4].marker = phaseLocked
	default:
		if store.IsTerminalState(model.runState) {
			phases[2].marker = phaseDone
			phases[3].marker = phaseDone
			return phases
		}
		phases[2].marker = phaseRun
	}
	return phases
}

func specRunPhases(model *cockpitModel) []cockpitPhase {
	phases := []cockpitPhase{
		{name: "AGENT", marker: phaseWait},
		{name: "VERIFY", marker: phaseWait},
		{name: "COMMIT", marker: phaseWait},
	}
	allCompleted := model.allTasksCompleted()
	switch {
	case allCompleted:
		for index := range phases {
			phases[index].marker = phaseDone
		}
	case model.runState == store.StateVerifying:
		phases[0].marker = phaseDone
		phases[1].marker = phaseRun
	case store.IsTerminalState(model.runState):
		phases[0].marker = phaseDone
		phases[1].marker = phaseDone
	default:
		phases[0].marker = phaseRun
	}
	if !model.viewport.HasKind(runevent.KindDaemonQA) {
		return phases
	}
	qaMarker := phaseLocked
	if allCompleted {
		qaMarker = phaseRun
		if store.IsTerminalState(model.runState) {
			qaMarker = phaseDone
		}
	}
	return append(phases, cockpitPhase{name: "QA", marker: qaMarker})
}

func renderPhaseRow(width int, phases []cockpitPhase) string {
	rendered := make([]string, 0, len(phases))
	for _, phase := range phases {
		rendered = append(rendered, renderPhase(phase))
	}
	line := strings.Join(rendered, styleMuted.Render("  >  "))
	return panel(width, 3, line, false)
}

func renderPhase(phase cockpitPhase) string {
	return styleBright.Render(phase.name) + " " + phaseMarkerStyle(phase.marker).Render("["+phase.marker+"]")
}

func phaseMarkerStyle(marker string) lipgloss.Style {
	switch marker {
	case phaseDone:
		return styleTool
	case phaseRun:
		return styleAccent
	case phaseLocked:
		return styleError
	default:
		return styleMuted
	}
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

const (
	detailModalMinWidth  = 76
	detailModalMinHeight = 20
)

func renderCockpitDetailOverlay(model *cockpitModel, layout cockpitLayout, base string) string {
	baseLines := strings.Split(base, "\n")
	if layout.width < detailModalMinWidth || len(baseLines) < detailModalMinHeight {
		return renderCockpitFullSurfaceDetail(model, layout)
	}
	modalWidth, modalHeight := detailModalSize(layout.width, len(baseLines))
	modal := renderCockpitDetailPane(model, modalWidth, modalHeight)
	modalLines := strings.Split(modal, "\n")
	top := maxInt((len(baseLines)-len(modalLines))/2, 0)
	left := maxInt((layout.width-modalWidth)/2, 0)
	lines := make([]string, len(baseLines))
	for index, line := range baseLines {
		lines[index] = styleMuted.Render(padRightDisplay(truncateDisplay(stripANSI(line), layout.width), layout.width))
	}
	for index, line := range modalLines {
		target := top + index
		if target >= len(lines) {
			break
		}
		lines[target] = padRightDisplay(strings.Repeat(" ", left)+line, layout.width)
	}
	return strings.Join(lines, "\n")
}

func renderCockpitFullSurfaceDetail(model *cockpitModel, layout cockpitLayout) string {
	return strings.Join([]string{
		renderCockpitHeaderArea(model, layout.width),
		panel(layout.width, layout.bodyHeight, model.renderDetail(layout.width, layout.bodyHeight), true),
		renderCockpitFooter(model, layout.width),
	}, "\n")
}

func detailModalSize(width int, height int) (int, int) {
	modalWidth := minInt(maxInt(width*78/100, 62), maxInt(width-4, 1))
	modalHeight := minInt(maxInt(height*50/100, 12), maxInt(height-2, 1))
	return modalWidth, modalHeight
}

func renderCockpitFooter(model *cockpitModel, width int) string {
	return model.renderFooter(width)
}

func (model *cockpitModel) renderDetail(width int, height int) string {
	detail := model.detail
	innerWidth := maxInt(width-4, 1)
	innerHeight := maxInt(height-2, 1)
	header := detailHeaderLines(detail, innerWidth)
	bodyHeight := detailBodyHeight(detail, innerHeight)
	start, end := detailVisibleRange(detail, bodyHeight)
	body := []string{}
	for _, line := range detail.lines[start:end] {
		body = append(body, truncateDisplay(line, innerWidth))
	}
	for len(body) < bodyHeight {
		body = append(body, "")
	}
	footer := styleMuted.Render(truncateDisplay(detailScrollFooter(detail, start, end), innerWidth))
	lines := append(header, body...)
	lines = append(lines, footer)
	return strings.Join(limitLines(lines, innerHeight), "\n")
}

func detailHeaderLines(detail *issueDetailView, width int) []string {
	lines := []string{
		styleAccent.Bold(true).Render(detailTitleLine(detail, width)),
		styleMuted.Render(strings.Repeat("-", width)),
		styleBright.Render(truncateDisplay(detailSubject(detail), width)),
		styleMuted.Render(truncateDisplay(detailMeta(detail), width)),
		styleMuted.Render(truncateDisplay(detailSource(detail), width)),
	}
	if detail.stale {
		lines = append(lines, styleError.Render(truncateDisplay("STALE: keeping last readable task file", width)))
	}
	return append(lines, "")
}

func detailTitleLine(detail *issueDetailView, width int) string {
	left := "REVIEW.ISSUE"
	if detail.kind == detailTask {
		left = "SPEC.TASK"
	}
	if detail.kind == detailReviewIssue && detail.ordinal > 0 {
		left += fmt.Sprintf("  #%03d", detail.ordinal)
	}
	if detail.kind == detailTask && strings.TrimSpace(detail.task.ID) != "" {
		left += "  " + strings.TrimSpace(detail.task.ID)
	}
	hint := "Esc close · j/k scroll"
	padding := maxInt(width-displayWidth(left)-displayWidth(hint), 1)
	return truncateDisplay(left+strings.Repeat(" ", padding)+hint, width)
}

func detailSubject(detail *issueDetailView) string {
	if detail.kind == detailTask {
		return emptyDash(detail.task.Title)
	}
	return emptyDash(detail.issue.Title)
}

func detailMeta(detail *issueDetailView) string {
	if detail.kind == detailTask {
		parts := []string{emptyDash(detail.task.Type), emptyDash(string(detail.task.Status)), emptyDash(detail.task.File)}
		return strings.Join(parts, " · ")
	}
	issue := detail.issue
	return fmt.Sprintf("%s · %s · %s", emptyDash(issue.Severity), emptyDash(issue.Status), emptyDash(issueLocation(issue)))
}

func detailSource(detail *issueDetailView) string {
	if detail.kind == detailTask {
		return "source: " + emptyDash(detail.task.File) + " (read-only)"
	}
	return "source: " + emptyDash(detail.issue.SourceRef)
}

func detailBodyHeight(detail *issueDetailView, innerHeight int) int {
	return maxInt(innerHeight-len(detailHeaderLines(detail, 1))-1, 1)
}

func detailVisibleRange(detail *issueDetailView, bodyHeight int) (int, int) {
	total := len(detail.lines)
	if total == 0 {
		return 0, 0
	}
	start := detail.scroll
	if start > total {
		start = total
	}
	end := minInt(start+bodyHeight, total)
	return start, end
}

func detailScrollFooter(detail *issueDetailView, start int, end int) string {
	total := len(detail.lines)
	if total == 0 {
		return "Line 0-0 of 0 · PgUp/PgDn page"
	}
	return fmt.Sprintf("Line %d-%d of %d · PgUp/PgDn page", start+1, end, total)
}

func (model *cockpitModel) detailPageSize() int {
	if model.detail == nil {
		return 1
	}
	layout := cockpitLayoutFor(model)
	height := layout.bodyHeight
	if layout.width >= detailModalMinWidth {
		_, totalHeight := detailModalSize(layout.width, layout.bodyHeight+4)
		height = totalHeight
	}
	return detailBodyHeight(model.detail, maxInt(height-2, 1))
}

func (model *cockpitModel) detailMaxScroll() int {
	if model.detail == nil {
		return 0
	}
	return maxInt(len(model.detail.lines)-model.detailPageSize(), 0)
}

// renderWorkItemPane renders the shared Work Queue surface. Run Kind only
// decides which Work Items feed it; the layout path stays identical.
func (model *cockpitModel) renderWorkItemPane(width int, height int) string {
	items := model.workItems()
	innerHeight := maxInt(height-2, 1)
	lines := []string{styleAccent.Bold(true).Render(fmt.Sprintf("WORK QUEUE (%d)", len(items))), ""}
	if len(items) == 0 {
		lines = append(lines, styleMuted.Render("No Work Items"))
		return model.renderWorkQueueWithFooter(lines, innerHeight, width)
	}
	if model.selected >= len(items) {
		model.selected = len(items) - 1
	}
	rowBudget := maxInt(innerHeight-len(lines)-1, 0)
	visible := maxInt(rowBudget/4, 1)
	if model.selected < model.issueTop {
		model.issueTop = model.selected
	}
	if model.selected >= model.issueTop+visible {
		model.issueTop = model.selected - visible + 1
	}
	end := minInt(model.issueTop+visible, len(items))
	rowLines := []string{}
	for index := model.issueTop; index < end; index++ {
		if !model.specRun() {
			if separator := model.batchSeparator(index, width); separator != "" {
				rowLines = append(rowLines, styleAccent.Render(separator))
			}
		}
		block := model.workItemBlock(items[index], model.workItemStatusLabel(index), index, width)
		if len(rowLines)+len(block) > rowBudget && len(rowLines) > 0 {
			break
		}
		rowLines = append(rowLines, block...)
	}
	lines = append(lines, limitLines(rowLines, rowBudget)...)
	return model.renderWorkQueueWithFooter(lines, innerHeight, width)
}

func (model *cockpitModel) renderWorkQueueWithFooter(lines []string, innerHeight int, width int) string {
	footer := styleMuted.Render(truncateDisplay(model.workQueueTotalsLine(), maxInt(width-4, 1)))
	if innerHeight <= 1 {
		return strings.Join(limitLines(lines, innerHeight), "\n")
	}
	for len(lines) < innerHeight-1 {
		lines = append(lines, "")
	}
	lines = append(limitLines(lines, innerHeight-1), footer)
	return strings.Join(limitLines(lines, innerHeight), "\n")
}

func (model *cockpitModel) workItems() []WorkItem {
	if model.specRun() {
		return model.taskWorkItems()
	}
	return model.issueWorkItems()
}

func (model *cockpitModel) workItemStatusLabel(index int) string {
	if model.specRun() {
		return model.taskStatusLabel(index)
	}
	return model.issueStatusLabel(index)
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

// issueWorkItems maps the Run's Review Issues into Work Items: positional
// display names plus the statuses of the latest artifact refresh.
func (model *cockpitModel) issueWorkItems() []WorkItem {
	items := make([]WorkItem, len(model.cfg.View.Issues))
	for index, issue := range model.cfg.View.Issues {
		status := ""
		if index < len(model.issueStatuses) {
			status = model.issueStatuses[index]
		}
		title := strings.TrimSpace(issue.Title)
		if title == "" {
			title = fmt.Sprintf("Issue #%03d", index+1)
		}
		items[index] = WorkItem{
			Name:     fmt.Sprintf("Issue #%03d", index+1),
			Title:    title,
			Status:   status,
			Severity: strings.TrimSpace(issue.Severity),
			Ordinal:  index + 1,
			Location: issueLocation(issue),
		}
	}
	return items
}

func issueLocation(issue rounds.Issue) string {
	file := strings.TrimSpace(issue.File)
	if file == "" {
		return ""
	}
	if issue.Line > 0 {
		return fmt.Sprintf("%s:%d", file, issue.Line)
	}
	return file
}

// workItemBlock renders one Work Item as the compact queue row both Run
// Kinds share: marker/severity/ordinal, title, and optional location.
func (model *cockpitModel) workItemBlock(item WorkItem, label string, index int, width int) []string {
	rowWidth := maxInt(width-4, 1)
	rowStyle := model.statusStyle(label)
	if index == model.selected {
		rowStyle = styleBright
	}
	title := strings.TrimSpace(item.Title)
	if title == "" {
		title = strings.TrimSpace(item.Name)
	}
	if strings.HasPrefix(item.Name, "task_") && title != "" {
		title = item.Name + " " + title
	}
	lines := []string{
		rowStyle.Render(workItemHeaderLine(item, label, index == model.selected, rowWidth)),
		styleBright.Render(truncateDisplay("  "+title, rowWidth)),
	}
	if location := strings.TrimSpace(item.Location); location != "" {
		lines = append(lines, styleMuted.Render(truncateDisplay("  "+location, rowWidth)))
	}
	return append(lines,
		"",
	)
}

func workItemHeaderLine(item WorkItem, label string, selected bool, width int) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}
	left := prefix + "[" + workItemStatusMarker(label) + "]"
	if severity := strings.TrimSpace(item.Severity); severity != "" {
		left += " " + strings.ToUpper(severity)
	}
	right := workItemOrdinal(item)
	if right == "" {
		return truncateDisplay(left, width)
	}
	padding := maxInt(width-displayWidth(left)-displayWidth(right), 1)
	return truncateDisplay(left+strings.Repeat(" ", padding)+right, width)
}

func workItemOrdinal(item WorkItem) string {
	if item.Ordinal > 0 {
		return fmt.Sprintf("#%d", item.Ordinal)
	}
	return strings.TrimSpace(item.Name)
}

func workItemStatusMarker(label string) string {
	switch label {
	case "Executing":
		return phaseRun
	case "Resolved", "Completed":
		return phaseDone
	case "Paused":
		return phaseLocked
	case "Invalid":
		return "invalid"
	case "Duplicated":
		return "dup"
	case "Failed":
		return "fail"
	default:
		return phaseWait
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
	separator := fmt.Sprintf("BATCH %03d/%03d", batch, total)
	if batch == model.currentBatch && !model.terminal && !store.IsTerminalState(model.runState) {
		elapsed := formatElapsed(model.now().Sub(model.batchStartedAt))
		pad := maxInt(width-4-displayWidth(separator)-displayWidth(elapsed), 1)
		separator += strings.Repeat(" ", pad) + elapsed
	}
	return truncateDisplay(separator, width-4)
}

func (model *cockpitModel) workQueueTotalsLine() string {
	done, total := model.progressCounts()
	unresolved := total - done
	if unresolved < 0 {
		unresolved = 0
	}
	if model.specRun() {
		return fmt.Sprintf("%d Tasks total · %d completed · %d unresolved", total, done, unresolved)
	}
	return fmt.Sprintf("%d issues total · %d resolved · %d unresolved", total, done, unresolved)
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
