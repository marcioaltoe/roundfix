package tui

import (
	"context"
	"encoding/json"
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
	// ColorEnabled is the effective color mode resolved by the caller
	// through the existing ROUNDFIX_COLOR / NO_COLOR contract. Disabled
	// color renders identity tokens: every distinction rides on markers.
	ColorEnabled bool
	// OnStop handles Ctrl-C in owning mode (Stop Request). Nil in attach.
	OnStop func()
	// Now overrides the clock for elapsed-time rendering. Nil means time.Now.
	Now func() time.Time
}

const defaultCockpitPollInterval = 250 * time.Millisecond

const minTwoPaneCockpitWidth = 88

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

// batchTimeSpan is the first-to-last Run Event timestamp range of one Batch.
type batchTimeSpan struct {
	first time.Time
	last  time.Time
}

type cockpitModel struct {
	ctx      context.Context
	cfg      CockpitConfig
	viewport *TimelineViewport
	now      func() time.Time
	tokens   Tokens

	focus    cockpitFocus
	selected int
	issueTop int
	detail   *issueDetailView

	issueStatuses       []string
	taskStatuses        []string
	taskJournalStatuses []string
	currentBatch        int
	batchTimes          map[int]batchTimeSpan
	batchTimeCursor     int64

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
		tokens:      ResolveTokens(cfg.ColorEnabled),
		batchTimes:  map[int]batchTimeSpan{},
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
// Review Issue artifacts for review Runs, Task journal events for spec Runs.
func (model *cockpitModel) refreshWorkItems() {
	if model.specRun() {
		model.refreshTasks()
		return
	}
	model.refreshIssues()
}

// refreshTasks keeps Task rows truthful under parallel execution. Task files
// remain the fallback/post-integration status source, while daemon.task events
// are the execution source for started, settled, and skipped states.
func (model *cockpitModel) refreshTasks() {
	tasks := model.cfg.View.Tasks
	if len(model.taskStatuses) != len(tasks) {
		model.taskStatuses = make([]string, len(tasks))
		model.taskJournalStatuses = make([]string, len(tasks))
		for index, task := range tasks {
			model.taskStatuses[index] = string(task.Status)
		}
	}
	model.refreshTaskFiles()
	model.refreshTaskJournalEvents()
	model.applyTaskJournalStatuses()
}

func (model *cockpitModel) refreshTaskFiles() {
	tasks := model.cfg.View.Tasks
	root := taskReadRoot(model.cfg.View)
	for index := range tasks {
		current := tasks[index]
		if err := spec.ReloadTask(root, &current); err == nil {
			model.taskStatuses[index] = string(current.Status)
		}
	}
}

func (model *cockpitModel) refreshTaskJournalEvents() {
	for index := range model.taskJournalStatuses {
		model.taskJournalStatuses[index] = ""
	}
	runID := model.cfg.RunID
	if strings.TrimSpace(runID) == "" {
		runID = model.cfg.View.RunID
	}
	if strings.TrimSpace(runID) == "" || model.cfg.Source == nil {
		return
	}
	cursor := int64(0)
	for {
		page, err := model.cfg.Source.RunEventsAfter(model.ctx, runID, cursor, 200)
		if err != nil {
			return
		}
		for _, entry := range page {
			if entry.Cursor > cursor {
				cursor = entry.Cursor
			}
			model.applyTaskJournalEvent(entry.Event)
		}
		if len(page) < 200 {
			return
		}
	}
}

type taskJournalEvent struct {
	Task   string `json:"task"`
	Phase  string `json:"phase"`
	Status string `json:"status"`
}

func (model *cockpitModel) applyTaskJournalEvent(event runevent.RunEvent) {
	if event.Kind != runevent.KindDaemonTask {
		return
	}
	taskEvent := parseTaskJournalEvent(event)
	if taskEvent.Task == "" {
		return
	}
	index := model.taskIndex(taskEvent.Task)
	if index < 0 {
		return
	}
	switch taskEvent.Phase {
	case "started":
		model.taskJournalStatuses[index] = string(spec.StatusInProgress)
	case "settled":
		switch spec.Status(taskEvent.Status) {
		case spec.StatusCompleted, spec.StatusFailed:
			model.taskJournalStatuses[index] = taskEvent.Status
		}
	case "skipped":
		model.taskJournalStatuses[index] = "skipped"
	}
}

func (model *cockpitModel) applyTaskJournalStatuses() {
	for index, journalStatus := range model.taskJournalStatuses {
		switch spec.Status(journalStatus) {
		case spec.StatusCompleted, spec.StatusFailed:
			model.taskStatuses[index] = journalStatus
		case spec.StatusInProgress:
			switch spec.Status(model.taskStatuses[index]) {
			case spec.StatusCompleted, spec.StatusFailed:
			default:
				model.taskStatuses[index] = journalStatus
			}
		default:
			if journalStatus == "skipped" {
				model.taskStatuses[index] = journalStatus
			}
		}
	}
}

func parseTaskJournalEvent(event runevent.RunEvent) taskJournalEvent {
	parsed := taskJournalEvent{Task: strings.TrimSpace(event.ReviewIssue)}
	if len(event.Payload) > 0 {
		var payload taskJournalEvent
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			if strings.TrimSpace(payload.Task) != "" {
				parsed.Task = strings.TrimSpace(payload.Task)
			}
			parsed.Phase = strings.TrimSpace(payload.Phase)
			parsed.Status = strings.TrimSpace(payload.Status)
		}
	}
	if parsed.Task == "" || parsed.Phase == "" || parsed.Status == "" {
		applyTaskJournalSummary(&parsed, event.Summary)
	}
	return parsed
}

func applyTaskJournalSummary(parsed *taskJournalEvent, summary string) {
	words := strings.Fields(strings.TrimSpace(summary))
	if parsed.Task == "" && len(words) >= 2 && strings.EqualFold(words[0], "Task") {
		parsed.Task = strings.TrimSpace(strings.TrimSuffix(words[1], ":"))
	}
	lower := strings.ToLower(summary)
	switch {
	case parsed.Phase == "" && strings.Contains(lower, " started "):
		parsed.Phase = "started"
	case parsed.Phase == "" && strings.Contains(lower, " settled "):
		parsed.Phase = "settled"
	case parsed.Phase == "" && strings.Contains(lower, " skipped"):
		parsed.Phase = "skipped"
	}
	if parsed.Status == "" {
		for _, status := range []spec.Status{spec.StatusCompleted, spec.StatusFailed} {
			if strings.Contains(lower, " "+string(status)) {
				parsed.Status = string(status)
				return
			}
		}
	}
}

func (model *cockpitModel) taskIndex(taskID string) int {
	for index, task := range model.cfg.View.Tasks {
		if task.ID == taskID {
			return index
		}
	}
	return -1
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
	model.currentBatch = current
	model.refreshBatchClocks()
}

// refreshBatchClocks folds new journal events into per-Batch timestamp
// spans. The scan is incremental: the cursor only moves forward, so each
// poll reads only the events that arrived since the last one.
func (model *cockpitModel) refreshBatchClocks() {
	if model.cfg.Source == nil {
		return
	}
	runID := strings.TrimSpace(model.cfg.RunID)
	if runID == "" {
		runID = strings.TrimSpace(model.cfg.View.RunID)
	}
	if runID == "" {
		return
	}
	if model.batchTimes == nil {
		model.batchTimes = map[int]batchTimeSpan{}
	}
	for {
		page, err := model.cfg.Source.RunEventsAfter(model.ctx, runID, model.batchTimeCursor, 200)
		if err != nil {
			return
		}
		for _, entry := range page {
			if entry.Cursor > model.batchTimeCursor {
				model.batchTimeCursor = entry.Cursor
			}
			event := entry.Event
			if event.Batch <= 0 || event.Time.IsZero() {
				continue
			}
			span, seen := model.batchTimes[event.Batch]
			if !seen || event.Time.Before(span.first) {
				span.first = event.Time
			}
			if event.Time.After(span.last) {
				span.last = event.Time
			}
			model.batchTimes[event.Batch] = span
		}
		if len(page) < 200 {
			return
		}
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
	case spec.StatusInProgress:
		if model.terminal || store.IsTerminalState(model.runState) {
			return "Paused"
		}
		return "Executing"
	}
	if model.taskStatuses[index] == "skipped" {
		return "Skipped"
	}
	if model.terminal || store.IsTerminalState(model.runState) {
		return "Paused"
	}
	return "Waiting"
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
	root := taskReadRoot(model.cfg.View)
	if err := spec.ReloadTask(root, &task); err != nil {
		return err
	}
	content, err := os.ReadFile(filepath.Join(root, task.File))
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
	collapsed    bool
}

func cockpitLayoutFor(model *cockpitModel) cockpitLayout {
	width := model.width
	if width <= 0 || model.height <= 0 {
		return cockpitLayout{}
	}
	bodyHeight := model.bodyHeight()
	if width < minTwoPaneCockpitWidth {
		return cockpitLayout{
			width:      width,
			bodyHeight: bodyHeight,
			rightWidth: width,
			collapsed:  true,
		}
	}
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
	if model.width <= 0 || model.height <= 0 {
		view := tea.NewView("")
		view.AltScreen = true
		return view
	}
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
		renderPhaseRow(model.tokens, width, runPhases(model)),
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
	runLabel := shortRunID(runID)
	chip := "[" + formatRunStateLabel(model.runState) + "]"
	suffix := model.statusSuffix()
	rightText := runLabel + " " + chip
	if suffix != "" {
		rightText += " " + suffix
	}
	right := model.tokens.Muted.Render(runLabel) + " " + model.runStateToken().Render(chip)
	if suffix != "" {
		right += " " + model.tokens.Muted.Render(suffix)
	}
	rightWidth := displayWidth(rightText)
	if rightWidth >= width {
		rightText = truncateDisplay(rightText, width-1)
		rightWidth = displayWidth(rightText)
		right = model.tokens.Muted.Render(rightText)
	}
	leftText = truncateDisplay(leftText, maxInt(width-rightWidth, 1))
	left := renderCockpitHeaderLeft(model.tokens, leftText)
	return padRightDisplay(left, maxInt(width-rightWidth, 1)) + right
}

func renderCockpitHeaderLeft(tokens Tokens, text string) string {
	if strings.HasPrefix(text, "ROUNDFIX") {
		return tokens.SectionLabel.Render("ROUNDFIX") + tokens.Muted.Render(strings.TrimPrefix(text, "ROUNDFIX"))
	}
	return tokens.Muted.Render(text)
}

// runStateToken maps the Run state to the state chip's token: green for a
// clean or integrated end, red for every other terminal outcome, amber
// while the Run is still moving.
func (model *cockpitModel) runStateToken() lipgloss.Style {
	switch {
	case model.runState == store.StateClean,
		model.runState == store.StateFetched,
		model.runState == store.StateIntegrationPending:
		return model.tokens.Done
	case store.IsTerminalState(model.runState):
		return model.tokens.Failed
	default:
		return model.tokens.Running
	}
}

func cockpitTargetLabel(view LiveRunView) string {
	target := ""
	if specRunView(view) {
		if strings.TrimSpace(view.SpecSlug) != "" {
			target = "SPEC " + strings.TrimSpace(view.SpecSlug)
		} else {
			target = "SPEC"
		}
		if view.Concurrency > 0 {
			target += fmt.Sprintf(" // Concurrency: %d", view.Concurrency)
		}
	} else if strings.TrimSpace(view.PRNumber) != "" {
		target = "PR #" + strings.TrimSpace(view.PRNumber)
	} else {
		target = strings.TrimSpace(view.HeadBranch)
	}
	if workDir := strings.TrimSpace(view.WorkDir); workDir != "" {
		if target != "" {
			target += " // "
		}
		if specRunView(view) {
			target += "WORKTREE " + workDir
		} else {
			target += "USER CHECKOUT " + workDir
		}
	}
	return target
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

func renderPhaseRow(tokens Tokens, width int, phases []cockpitPhase) string {
	rendered := make([]string, 0, len(phases))
	for _, phase := range phases {
		rendered = append(rendered, renderPhase(tokens, phase))
	}
	line := strings.Join(rendered, tokens.Muted.Render("  >  "))
	return panel(width, 3, line, false)
}

func renderPhase(tokens Tokens, phase cockpitPhase) string {
	return phase.name + " " + phaseMarkerStyle(tokens, phase.marker).Render("["+phase.marker+"]")
}

func phaseMarkerStyle(tokens Tokens, marker string) lipgloss.Style {
	switch marker {
	case phaseDone:
		return tokens.Done
	case phaseRun:
		return tokens.Running
	case phaseLocked:
		return tokens.Locked
	default:
		return tokens.Waiting
	}
}

func renderCockpitBody(model *cockpitModel, layout cockpitLayout) string {
	if layout.collapsed {
		return renderCockpitCollapsedTimelinePane(model, layout.width, layout.bodyHeight)
	}
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
	innerWidth := maxInt(width-4, 1)
	content := []string{model.timelinePaneHeaderLine(innerWidth), ""}
	body := model.timelineBodyLines(innerWidth)
	content = append(content, limitTail(body, maxInt(height-2-len(content), 1))...)
	return panel(width, height, strings.Join(content, "\n"), model.focus == focusTimeline)
}

func renderCockpitCollapsedTimelinePane(model *cockpitModel, width int, height int) string {
	innerWidth := maxInt(width-4, 1)
	content := []string{
		model.timelinePaneHeaderLine(innerWidth),
		model.tokens.Muted.Render(truncateDisplay("QUEUE.SUMMARY "+model.workQueueTotalsLine(), innerWidth)),
		"",
	}
	body := model.timelineBodyLines(innerWidth)
	content = append(content, limitTail(body, maxInt(height-2-len(content), 1))...)
	return panel(width, height, strings.Join(content, "\n"), true)
}

// timelinePaneHeaderLine pins the pane title and the Live · detail
// indicator to the top of the timeline pane; the indicator follows the
// Detail Modal state.
func (model *cockpitModel) timelinePaneHeaderLine(width int) string {
	label := "SESSION.TIMELINE"
	if width <= displayWidth(label) {
		return model.tokens.SectionLabel.Render(truncateDisplay(label, width))
	}
	indicator := "Live · detail hidden"
	if model.detail != nil {
		indicator = "Live · detail open"
	}
	available := width - displayWidth(label) - 1
	if displayWidth(indicator) > available {
		indicator = truncateDisplay(indicator, available)
	}
	if indicator == "" {
		return model.tokens.SectionLabel.Render(label)
	}
	padding := maxInt(width-displayWidth(label)-displayWidth(indicator), 1)
	return model.tokens.SectionLabel.Render(label) + strings.Repeat(" ", padding) + model.tokens.Muted.Render(indicator)
}

func (model *cockpitModel) timelineBodyLines(width int) []string {
	lines := model.viewport.VisibleLines()
	if len(lines) == 0 {
		return model.mutedEmptyLines(model.timelineEmptyCopy(), width)
	}
	return model.styleTimelineRows(lines, width)
}

// workQueueEmptyLines explains the empty Work Queue: which Run kind this
// is and what would populate the pane.
func (model *cockpitModel) workQueueEmptyLines(width int) []string {
	return model.mutedEmptyLines(model.workQueueEmptyCopy(), width)
}

func (model *cockpitModel) workQueueEmptyCopy() []string {
	switch {
	case model.specRun():
		return []string{
			"No Tasks yet.",
			"An implement Run queues its Spec's",
			"Task Graph here.",
		}
	case model.cfg.View.RunKind == store.KindFetch:
		return []string{
			"No Work Items.",
			"A Fetch Run writes Review artifacts",
			"to disk and starts no Agent.",
		}
	default:
		return []string{
			"No Work Items yet.",
			"Review Issues queue here once the",
			"Run records them.",
		}
	}
}

func (model *cockpitModel) timelineEmptyCopy() []string {
	switch {
	case model.cfg.View.RunKind == store.KindFetch:
		return []string{
			"No Run Events.",
			"A Fetch Run writes Review artifacts to disk and starts no Agent.",
		}
	case model.specRun():
		return []string{
			"No Run Events yet.",
			"Task execution and verification stream here.",
		}
	default:
		return []string{
			"No Run Events yet.",
			"Agent and Daemon activity streams here.",
		}
	}
}

func (model *cockpitModel) mutedEmptyLines(copyLines []string, width int) []string {
	lines := make([]string, 0, len(copyLines))
	for _, line := range copyLines {
		lines = append(lines, model.tokens.Muted.Render(truncateDisplay(line, width)))
	}
	return lines
}

// styleTimelineRows colors the plain viewport rows through the tokens:
// batch group headers, the muted timestamp gutter, and kind labels. Rows
// are truncated to the pane width first, so every event stays one bounded
// line.
func (model *cockpitModel) styleTimelineRows(lines []string, width int) []string {
	styled := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := truncateDisplay(strings.TrimRight(line, "\r\n"), width)
		styled = append(styled, model.styleTimelineRow(trimmed))
	}
	return styled
}

func (model *cockpitModel) styleTimelineRow(line string) string {
	if strings.HasPrefix(line, timelineExpandedMarker+" ") || strings.HasPrefix(line, timelineCollapsedMarker+" ") {
		return model.styleTimelineBatchHeader(line)
	}
	if gutter, rest, ok := splitTimelineGutter(line); ok {
		return model.tokens.Muted.Render(gutter) + model.styleTimelineKindRow(rest)
	}
	return model.styleTimelineKindRow(line)
}

// styleTimelineBatchHeader colors the group header: marker and Batch label
// as a section label, the state word in its state token, the elapsed stamp
// muted.
func (model *cockpitModel) styleTimelineBatchHeader(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return model.tokens.SectionLabel.Render(line)
	}
	label := fields[0] + " " + fields[1]
	rest := fields[2:]
	if len(rest) > 0 && strings.Contains(rest[0], "/") {
		label += " " + rest[0]
		rest = rest[1:]
	}
	styled := model.tokens.SectionLabel.Render(label)
	for _, field := range rest {
		if strings.Contains(field, ":") {
			styled += " " + model.tokens.Muted.Render(field)
			continue
		}
		styled += " " + model.timelineBatchStateToken(field).Render(field)
	}
	return styled
}

func (model *cockpitModel) timelineBatchStateToken(state string) lipgloss.Style {
	switch state {
	case "executing", "started", "running":
		return model.tokens.Running
	case "waiting":
		return model.tokens.Waiting
	case "completed":
		return model.tokens.Done
	case "failed", "stopped":
		return model.tokens.Failed
	default:
		return model.tokens.Muted
	}
}

// styleTimelineKindRow colors the row's kind label through the tokens:
// PLAN, [TOOL], and SESSION in the section-label cyan family, THINK rows
// muted, daemon milestone labels cyan. Text after the label stays plain so
// the no-color render is marker-only by construction.
func (model *cockpitModel) styleTimelineKindRow(line string) string {
	if label, rest, ok := cutTimelineKindLabel(line); ok {
		if label == "THINK" {
			return model.tokens.Muted.Render(line)
		}
		return model.tokens.SectionLabel.Render(label) + rest
	}
	return line
}

// cutTimelineKindLabel splits a known kind label off the row. It reports
// ok=false for plain streaming text rows.
func cutTimelineKindLabel(line string) (string, string, bool) {
	for _, label := range []string{
		"[TOOL]", "PLAN", "THINK", "SESSION",
		"STATUS", "FETCH", "VERIFY", "COMMIT", "PUSH", "SOURCE", "RETRY", "OUTCOME", "TASK", "QA", "SELECTION",
	} {
		if line == label {
			return label, "", true
		}
		if strings.HasPrefix(line, label+" ") {
			return label, line[len(label):], true
		}
	}
	return "", "", false
}

// splitTimelineGutter splits the aligned timestamp column off a rendered
// event row: either "HH:MM:SS " or a blank column of the same width.
func splitTimelineGutter(line string) (string, string, bool) {
	runes := []rune(line)
	if len(runes) < timelineGutterWidth {
		return "", "", false
	}
	gutter := string(runes[:timelineGutterWidth])
	rest := string(runes[timelineGutterWidth:])
	if gutter == strings.Repeat(" ", timelineGutterWidth) {
		return gutter, rest, true
	}
	if gutter[2] == ':' && gutter[5] == ':' && gutter[8] == ' ' {
		return gutter, rest, true
	}
	return "", "", false
}

// renderCockpitDetailPane draws the Detail Modal frame through the accent
// border token; the identity token keeps the frame structure without color.
func renderCockpitDetailPane(model *cockpitModel, width int, height int) string {
	return model.tokens.ActiveBorder.Width(width).Height(height).Render(model.renderDetail(width, height))
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
		lines[index] = model.tokens.Muted.Render(padRightDisplay(truncateDisplay(stripANSI(line), layout.width), layout.width))
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
		renderCockpitDetailPane(model, layout.width, layout.bodyHeight),
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
	header := model.detailHeaderLines(detail, innerWidth)
	bodyHeight := model.detailBodyHeight(detail, innerHeight)
	start, end := detailVisibleRange(detail, bodyHeight)
	body := []string{}
	for _, line := range detail.lines[start:end] {
		body = append(body, model.detailBodyLine(line, innerWidth))
	}
	for len(body) < bodyHeight {
		body = append(body, "")
	}
	footer := model.tokens.Muted.Render(truncateDisplay(detailScrollFooter(detail, start, end), innerWidth))
	lines := append(header, body...)
	lines = append(lines, footer)
	return strings.Join(limitLines(lines, innerHeight), "\n")
}

// detailBodyLine renders one body line; markdown headings become section
// labels so the body reads sectioned, and the heading markers survive the
// no-color render.
func (model *cockpitModel) detailBodyLine(line string, width int) string {
	rendered := truncateDisplay(line, width)
	if strings.HasPrefix(strings.TrimSpace(line), "#") {
		return model.tokens.SectionLabel.Render(rendered)
	}
	return rendered
}

// detailHeaderLines renders the modal's title block: the Work Item title
// line with the key hints, a rule, the subject, the severity/status/
// location line, and the read-only source line, all through the tokens.
func (model *cockpitModel) detailHeaderLines(detail *issueDetailView, width int) []string {
	lines := []string{
		model.detailTitleLine(detail, width),
		model.tokens.Muted.Render(strings.Repeat("-", width)),
		truncateDisplay(detailSubject(detail), width),
		model.detailMetaLine(detail, width),
		model.tokens.Muted.Render(truncateDisplay(detailSource(detail), width)),
	}
	if selectionLine := model.detailSelectionLine(detail); selectionLine != "" {
		lines = append(lines, model.tokens.Muted.Render(truncateDisplay(selectionLine, width)))
	}
	if detail.stale {
		lines = append(lines, model.tokens.Failed.Render(truncateDisplay("STALE: keeping last readable task file", width)))
	}
	return append(lines, "")
}

func (model *cockpitModel) detailTitleLine(detail *issueDetailView, width int) string {
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
	if displayWidth(left)+padding+displayWidth(hint) > width {
		return model.tokens.SectionLabel.Render(truncateDisplay(left+strings.Repeat(" ", padding)+hint, width))
	}
	return model.tokens.SectionLabel.Render(left) + strings.Repeat(" ", padding) + model.tokens.Muted.Render(hint)
}

func detailSubject(detail *issueDetailView) string {
	if detail.kind == detailTask {
		return emptyDash(detail.task.Title)
	}
	return emptyDash(detail.issue.Title)
}

// detailMetaLine renders severity and status in the artifact's state token
// and the location muted, per the modal design's title block.
func (model *cockpitModel) detailMetaLine(detail *issueDetailView, width int) string {
	const separator = " · "
	first, status, location := detailMetaParts(detail)
	plain := first + separator + status + separator + location
	if displayWidth(plain) > width {
		return model.tokens.Muted.Render(truncateDisplay(plain, width))
	}
	statusToken := model.artifactStatusToken(status)
	return statusToken.Render(first) +
		model.tokens.Muted.Render(separator) +
		statusToken.Render(status) +
		model.tokens.Muted.Render(separator+location)
}

func detailMetaParts(detail *issueDetailView) (string, string, string) {
	if detail.kind == detailTask {
		return emptyDash(string(detail.task.Type)), emptyDash(string(detail.task.Status)), emptyDash(detail.task.File)
	}
	issue := detail.issue
	return emptyDash(issue.Severity), emptyDash(issue.Status), emptyDash(issueLocation(issue))
}

// artifactStatusToken maps a Work Item artifact status (Review Issue or
// Task) to its state token.
func (model *cockpitModel) artifactStatusToken(status string) lipgloss.Style {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case rounds.StatusResolved, string(spec.StatusCompleted):
		return model.tokens.Done
	case rounds.StatusFailed:
		return model.tokens.Failed
	case string(spec.StatusInProgress):
		return model.tokens.Running
	case rounds.StatusInvalid, rounds.StatusDuplicated, "skipped":
		return model.tokens.Muted
	default:
		return model.tokens.Waiting
	}
}

func detailSource(detail *issueDetailView) string {
	if detail.kind == detailTask {
		return "source: " + emptyDash(detail.task.File) + " (read-only)"
	}
	return "source: " + emptyDash(detail.issue.SourceRef)
}

func (model *cockpitModel) detailSelectionLine(detail *issueDetailView) string {
	record, ok := model.selectionForDetail(detail)
	if !ok {
		return ""
	}
	return "selection: " + runevent.SelectionLifecycleSummary(record)
}

func (model *cockpitModel) selectionForDetail(detail *issueDetailView) (runevent.SelectionLifecycleRecord, bool) {
	if detail.kind == detailTask {
		return model.latestSelection("task", detail.task.ID)
	}
	if batch := model.batchOf(detail.ordinal - 1); batch > 0 {
		if record, ok := model.latestSelection("review", fmt.Sprintf("batch-%03d", batch)); ok {
			return record, true
		}
	}
	return model.latestReviewSelection()
}

func (model *cockpitModel) detailBodyHeight(detail *issueDetailView, innerHeight int) int {
	return maxInt(innerHeight-len(model.detailHeaderLines(detail, 1))-1, 1)
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
	return model.detailBodyHeight(model.detail, maxInt(height-2, 1))
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
	lines := []string{model.tokens.SectionLabel.Render(fmt.Sprintf("WORK QUEUE (%d)", len(items))), ""}
	if len(items) == 0 {
		lines = append(lines, model.workQueueEmptyLines(maxInt(width-4, 1))...)
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
				rowLines = append(rowLines, separator)
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
	footer := model.tokens.Muted.Render(truncateDisplay(model.workQueueTotalsLine(), maxInt(width-4, 1)))
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

// workItemBlock renders one Work Item as the card both Run Kinds share:
// marker and severity in the state color, ordinal and location muted, and
// the accent selection border around the selected card.
func (model *cockpitModel) workItemBlock(item WorkItem, label string, index int, width int) []string {
	rowWidth := maxInt(width-4, 1)
	selected := index == model.selected
	contentWidth := rowWidth
	if selected {
		contentWidth = maxInt(rowWidth-2, 1)
	}
	title := strings.TrimSpace(item.Title)
	if title == "" {
		title = strings.TrimSpace(item.Name)
	}
	if strings.HasPrefix(item.Name, "task_") && title != "" {
		title = item.Name + " " + title
	}
	lines := []string{
		model.workItemHeaderLine(item, label, selected, contentWidth),
		truncateDisplay("  "+title, contentWidth),
	}
	if location := strings.TrimSpace(item.Location); location != "" {
		lines = append(lines, model.tokens.Muted.Render(truncateDisplay("  "+location, contentWidth)))
	}
	if record, ok := model.selectionForWorkItem(index, item); ok {
		lines = append(lines, model.tokens.Muted.Render(truncateDisplay("  selection: "+runevent.SelectionLifecycleSummary(record), contentWidth)))
	}
	if selected {
		card := model.tokens.Selection.Width(rowWidth).Render(strings.Join(lines, "\n"))
		return append(strings.Split(card, "\n"), "")
	}
	return append(lines,
		"",
	)
}

func (model *cockpitModel) workItemHeaderLine(item WorkItem, label string, selected bool, width int) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}
	left := prefix + "[" + workItemStatusMarker(label) + "]"
	if severity := strings.TrimSpace(item.Severity); severity != "" {
		left += " " + strings.ToUpper(severity)
	}
	style := model.statusStyle(label)
	right := workItemOrdinal(item)
	if right == "" {
		return style.Render(truncateDisplay(left, width))
	}
	padding := maxInt(width-displayWidth(left)-displayWidth(right), 1)
	if displayWidth(left)+padding+displayWidth(right) > width {
		return style.Render(truncateDisplay(left+strings.Repeat(" ", padding)+right, width))
	}
	return style.Render(left) + strings.Repeat(" ", padding) + model.tokens.Muted.Render(right)
}

func workItemOrdinal(item WorkItem) string {
	if item.Ordinal > 0 {
		return fmt.Sprintf("#%d", item.Ordinal)
	}
	return strings.TrimSpace(item.Name)
}

func (model *cockpitModel) selectionForWorkItem(index int, item WorkItem) (runevent.SelectionLifecycleRecord, bool) {
	if model.specRun() {
		return model.latestSelection("task", item.Name)
	}
	if batch := model.batchOf(index); batch > 0 {
		return model.latestSelection("review", fmt.Sprintf("batch-%03d", batch))
	}
	return model.latestReviewSelection()
}

func (model *cockpitModel) latestSelection(scopeKind string, scopeID string) (runevent.SelectionLifecycleRecord, bool) {
	scopeKind = strings.TrimSpace(scopeKind)
	scopeID = strings.TrimSpace(scopeID)
	if scopeKind == "" || scopeID == "" {
		return runevent.SelectionLifecycleRecord{}, false
	}
	records := model.selectionRecords()
	for index := len(records) - 1; index >= 0; index-- {
		record := records[index]
		if record.ScopeKind == scopeKind && record.ScopeID == scopeID {
			return record, true
		}
	}
	return runevent.SelectionLifecycleRecord{}, false
}

func (model *cockpitModel) latestReviewSelection() (runevent.SelectionLifecycleRecord, bool) {
	records := model.selectionRecords()
	for index := len(records) - 1; index >= 0; index-- {
		record := records[index]
		if record.ScopeKind == "review" {
			return record, true
		}
	}
	return runevent.SelectionLifecycleRecord{}, false
}

func (model *cockpitModel) selectionRecords() []runevent.SelectionLifecycleRecord {
	records := append([]runevent.SelectionLifecycleRecord{}, model.cfg.View.Selections...)
	if model.viewport != nil {
		records = append(records, model.viewport.SelectionRecords()...)
	}
	return records
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
	case "Skipped":
		return "skip"
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
	lineWidth := maxInt(width-4, 1)
	label := fmt.Sprintf("BATCH %03d/%03d", batch, total)
	elapsed := model.batchElapsed(batch)
	if elapsed == "" {
		return model.tokens.SectionLabel.Render(truncateDisplay(label, lineWidth))
	}
	pad := maxInt(lineWidth-displayWidth(label)-displayWidth(elapsed), 1)
	if displayWidth(label)+pad+displayWidth(elapsed) > lineWidth {
		return model.tokens.SectionLabel.Render(truncateDisplay(label+strings.Repeat(" ", pad)+elapsed, lineWidth))
	}
	return model.tokens.SectionLabel.Render(label) + strings.Repeat(" ", pad) + model.tokens.Muted.Render(elapsed)
}

// batchElapsed derives the Batch's elapsed stamp from its Run Event
// timestamps: first event to last event, with the executing Batch of a live
// Run counting against now. Batches without timestamped events have no
// clock to show.
func (model *cockpitModel) batchElapsed(batch int) string {
	span, seen := model.batchTimes[batch]
	if !seen {
		return ""
	}
	end := span.last
	if batch == model.currentBatch && !model.terminal && !store.IsTerminalState(model.runState) {
		end = model.now()
	}
	return formatElapsed(end.Sub(span.first))
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

// statusStyle maps the Work Item's execution label to its state token: the
// same color family the item's text marker belongs to.
func (model *cockpitModel) statusStyle(label string) lipgloss.Style {
	switch label {
	case "Executing":
		return model.tokens.Running
	case "Resolved", "Completed":
		return model.tokens.Done
	case "Failed":
		return model.tokens.Failed
	case "Paused":
		return model.tokens.Locked
	case "Invalid", "Duplicated", "Skipped":
		return model.tokens.Muted
	default:
		return model.tokens.Waiting
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
	keys := model.footerKeys()
	line := model.fitFooterLine("Keys: "+keys, width)
	return styleMuted.Render(padRightDisplay(line, width))
}

func (model *cockpitModel) fitFooterLine(line string, width int) string {
	if displayWidth(line) <= width {
		return line
	}
	suffix := " · " + model.footerModeKey()
	if strings.HasSuffix(line, suffix) && displayWidth(suffix) < width {
		prefix := strings.TrimSuffix(line, suffix)
		return truncateDisplay(prefix, width-displayWidth(suffix)) + suffix
	}
	return truncateDisplay(line, width)
}

func (model *cockpitModel) footerKeys() string {
	if model.detail != nil {
		return strings.Join([]string{
			"Esc close",
			"j/k scroll",
			"PgUp/PgDn page",
			model.footerModeKey(),
		}, " · ")
	}
	if model.width > 0 && model.width < minTwoPaneCockpitWidth {
		return strings.Join([]string{
			"↑↓ scroll",
			"PgUp/PgDn page",
			"widen for Work Queue",
			model.footerModeKey(),
		}, " · ")
	}
	item := "issue"
	if model.specRun() {
		item = "Task"
	}
	return strings.Join([]string{
		"Tab focus",
		"↑↓ move/scroll",
		"PgUp/PgDn page",
		"Enter " + item,
		"D show detail",
		"End follow",
		model.footerModeKey(),
	}, " · ")
}

func (model *cockpitModel) footerModeKey() string {
	switch {
	case model.cfg.Mode == CockpitAttach:
		return "q detach"
	case model.terminal:
		return "q close"
	default:
		return "Ctrl-C stop"
	}
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
