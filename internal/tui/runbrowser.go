package tui

import (
	"context"
	"io"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

// browserDefaultWidth sizes the Run Browser before the first WindowSizeMsg.
const browserDefaultWidth = 100

// browserChromeLines counts the non-row lines the Run Browser renders:
// header, blank, blank, footer.
const browserChromeLines = 4

// browserRowColumns indexes the FormatRunRow fields by meaning.
const (
	browserColumnID = iota
	browserColumnState
	browserColumnKind
	browserColumnTarget
	browserColumnAgent
	browserColumnStart
	browserColumnDuration
	browserColumnBranch
	browserColumnRepo
	browserColumnCount
)

// browserColumnDropOrder documents how the Run Browser degrades at small
// widths: columns drop in this order until the rows fit, so the short run
// id, state, and duration always survive.
var browserColumnDropOrder = []int{
	browserColumnBranch,
	browserColumnAgent,
	browserColumnStart,
	browserColumnRepo,
	browserColumnTarget,
	browserColumnKind,
}

// BrowserOutcome reports how the Run Browser closed. A zero value means
// the browser is still open.
type BrowserOutcome struct {
	RunID     string // selected Run id; empty on cancel
	Cancelled bool
}

// RunBrowser is the Run Browser model: every repository's Runs newest
// first, Active Runs only by default, with an all/active state toggle. It
// is read-only Run discovery — the model never attaches, mutates, or stops
// anything; it only reports the outcome the entry point acts on.
type RunBrowser struct {
	active  []store.Run
	all     []store.Run
	showAll bool
	cursor  int
	width   int
	height  int
	// now anchors relative start times and running durations for the
	// whole browser session; the list has no live auto-refresh.
	now     time.Time
	outcome BrowserOutcome
}

// NewRunBrowser builds the browser over pre-queried machine-wide listings,
// both newest first: the Active Runs shown by default and the full history
// behind the `a` toggle.
func NewRunBrowser(active []store.Run, all []store.Run) RunBrowser {
	return RunBrowser{active: active, all: all, now: time.Now()}
}

// Outcome reports how the browser closed; zero while still browsing.
func (browser RunBrowser) Outcome() BrowserOutcome {
	return browser.outcome
}

func (browser RunBrowser) rows() []store.Run {
	if browser.showAll {
		return browser.all
	}
	return browser.active
}

func (browser RunBrowser) Update(msg tea.Msg) (RunBrowser, tea.Cmd) {
	switch value := msg.(type) {
	case tea.WindowSizeMsg:
		browser.width = value.Width
		browser.height = value.Height
		return browser, nil
	case tea.KeyPressMsg:
		return browser.handleKey(tea.Key(value))
	}
	return browser, nil
}

func (browser RunBrowser) handleKey(key tea.Key) (RunBrowser, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		if browser.cursor > 0 {
			browser.cursor--
		}
	case "down", "j":
		if browser.cursor < len(browser.rows())-1 {
			browser.cursor++
		}
	case "a":
		browser.showAll = !browser.showAll
		browser.cursor = 0
	case "enter":
		rows := browser.rows()
		if len(rows) == 0 {
			return browser, nil
		}
		browser.outcome = BrowserOutcome{RunID: rows[browser.cursor].ID}
		return browser, tea.Quit
	case "q", "esc", "ctrl+c":
		browser.outcome = BrowserOutcome{Cancelled: true}
		return browser, tea.Quit
	}
	return browser, nil
}

func (browser RunBrowser) View() string {
	width := browser.width
	if width <= 0 {
		width = browserDefaultWidth
	}
	filter := "ACTIVE"
	if browser.showAll {
		filter = "ALL"
	}
	header := cockpitTokens.SectionLabel.Render(truncateDisplay("Run Browser — all repositories — "+filter, width))
	footer := cockpitTokens.Muted.Render(truncateDisplay("↑↓ move  enter attach  a all/active  q quit", width))

	body := browser.renderRows(width)
	if len(body) == 0 {
		body = []string{truncateDisplay(browser.emptyState(), width)}
	}

	lines := append([]string{header, ""}, body...)
	lines = append(lines, "", footer)
	return strings.Join(lines, "\n")
}

// emptyState names the current state filter. The browser is machine-wide,
// so an empty Active view really does mean nothing is running anywhere.
func (browser RunBrowser) emptyState() string {
	if browser.showAll {
		return "No Runs found."
	}
	return "No active Runs — press a to include terminal Runs."
}

// renderRows renders the visible window of Run rows: aligned columns, a
// `> ` cursor marker, and Active rows styled distinct when the all filter
// exposes terminal Runs next to them.
func (browser RunBrowser) renderRows(width int) []string {
	runs := browser.rows()
	if len(runs) == 0 {
		return nil
	}
	fields := make([][]string, 0, len(runs))
	for _, run := range runs {
		row := FormatRunRow(run, browser.now, true, true)
		row[browserColumnID] = browserShortRunID(run.ID)
		row[browserColumnRepo] = browserRepoName(run.GitRoot)
		fields = append(fields, row)
	}
	keep, widths := browser.fitColumns(fields, width)

	first, last := browser.visibleWindow(len(runs))
	lines := make([]string, 0, last-first)
	for index := first; index < last; index++ {
		marker := "  "
		if index == browser.cursor {
			marker = "> "
		}
		cells := make([]string, 0, len(keep))
		for _, column := range keep {
			cells = append(cells, padRightDisplay(fields[index][column], widths[column]))
		}
		line := truncateDisplay(strings.TrimRight(marker+strings.Join(cells, "  "), " "), width)
		if browser.showAll && !store.IsTerminalState(runs[index].State) {
			line = cockpitTokens.Running.Render(line)
		}
		lines = append(lines, line)
	}
	return lines
}

// fitColumns keeps every column that fits the width, dropping columns in
// browserColumnDropOrder until the widest row fits, and returns the kept
// column indexes with each column's alignment width.
func (browser RunBrowser) fitColumns(fields [][]string, width int) ([]int, []int) {
	widths := make([]int, browserColumnCount)
	for _, row := range fields {
		for column, value := range row {
			if size := displayWidth(value); size > widths[column] {
				widths[column] = size
			}
		}
	}
	kept := make(map[int]bool, browserColumnCount)
	for column := 0; column < browserColumnCount; column++ {
		kept[column] = true
	}
	rowWidth := func() int {
		total := 2 // cursor marker
		count := 0
		for column := 0; column < browserColumnCount; column++ {
			if kept[column] {
				total += widths[column]
				count++
			}
		}
		return total + 2*(count-1)
	}
	for _, column := range browserColumnDropOrder {
		if rowWidth() <= width {
			break
		}
		kept[column] = false
	}
	keep := make([]int, 0, browserColumnCount)
	for column := 0; column < browserColumnCount; column++ {
		if kept[column] {
			keep = append(keep, column)
		}
	}
	return keep, widths
}

// visibleWindow bounds the rendered rows to the terminal height, keeping
// the cursor row visible.
func (browser RunBrowser) visibleWindow(total int) (int, int) {
	if browser.height <= 0 {
		return 0, total
	}
	visible := browser.height - browserChromeLines
	if visible < 1 {
		visible = 1
	}
	if visible >= total {
		return 0, total
	}
	first := 0
	if browser.cursor >= visible {
		first = browser.cursor - visible + 1
	}
	return first, first + visible
}

// runBrowserProgram adapts RunBrowser to the Bubble Tea program surface so
// one full-screen program can drive the model the tests exercise directly.
type runBrowserProgram struct {
	browser RunBrowser
}

func (program runBrowserProgram) Init() tea.Cmd {
	return nil
}

func (program runBrowserProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	browser, cmd := program.browser.Update(msg)
	program.browser = browser
	return program, cmd
}

func (program runBrowserProgram) View() tea.View {
	view := tea.NewView(program.browser.View())
	view.AltScreen = true
	return view
}

// RunBrowserSession opens the Run Browser in the alternate screen and
// blocks until the user selects a Run or cancels. Context cancellation
// closes the session as a cancel: Run discovery never mutates anything.
func RunBrowserSession(ctx context.Context, output io.Writer, browser RunBrowser) (BrowserOutcome, error) {
	prog := tea.NewProgram(runBrowserProgram{browser: browser}, tea.WithOutput(output), tea.WithoutSignalHandler())
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			prog.Quit()
		case <-done:
		}
	}()
	final, err := prog.Run()
	if err != nil {
		return BrowserOutcome{}, err
	}
	outcome := final.(runBrowserProgram).browser.Outcome()
	if outcome.RunID == "" && !outcome.Cancelled {
		// The program ended without a key outcome: context cancellation
		// quit it, which is a cancel, never a selection.
		outcome.Cancelled = true
	}
	return outcome, nil
}

// browserRepoName is the readable repository column: the git root's base
// name; the full path stays on the text surface (`runs list --all`).
func browserRepoName(gitRoot string) string {
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return runRowEmptyField
	}
	return filepath.Base(gitRoot)
}

// browserShortRunID is the timestamp-less suffix of a run id
// (run_<timestamp>_<suffix>); ids in any other shape render in full. The
// text surface keeps full ids — the short form is TUI-only.
func browserShortRunID(id string) string {
	if !strings.HasPrefix(id, "run_") {
		return id
	}
	if index := strings.LastIndexByte(id, '_'); index >= len("run_") {
		return id[index+1:]
	}
	return id
}
