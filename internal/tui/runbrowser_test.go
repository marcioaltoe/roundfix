package tui

import (
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

var browserTestNow = time.Date(2026, 7, 7, 15, 0, 0, 0, time.UTC)

func browserTestRuns() (active []store.Run, all []store.Run) {
	completed := time.Date(2026, 7, 7, 12, 42, 0, 0, time.UTC)
	newest := store.Run{
		ID:          "run_20260707T144800Z_aaaaaaaaaaaaaaaa",
		Kind:        store.KindResolve,
		State:       store.StateActive,
		PRNumber:    "123",
		Agent:       "codex",
		LocalBranch: "feature/one",
		CreatedAt:   time.Date(2026, 7, 7, 14, 48, 0, 0, time.UTC),
	}
	older := store.Run{
		ID:          "run_20260707T143000Z_bbbbbbbbbbbbbbbb",
		Kind:        store.KindImplement,
		State:       store.StateResolvingWithAgent,
		SpecSlug:    "0020-run-browser",
		Agent:       "codex",
		LocalBranch: "ma/two",
		CreatedAt:   time.Date(2026, 7, 7, 14, 30, 0, 0, time.UTC),
	}
	terminal := store.Run{
		ID:          "run_20260707T120000Z_cccccccccccccccc",
		Kind:        store.KindResolve,
		State:       store.StateClean,
		PRNumber:    "99",
		Agent:       "codex",
		LocalBranch: "feature/old",
		CreatedAt:   time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC),
		CompletedAt: &completed,
	}
	active = []store.Run{newest, older}
	all = []store.Run{newest, older, terminal}
	return active, all
}

func newTestRunBrowser(t *testing.T) RunBrowser {
	t.Helper()
	active, all := browserTestRuns()
	browser := NewRunBrowser(active, all)
	browser.now = browserTestNow
	return browser
}

func pressBrowser(t *testing.T, browser RunBrowser, keystroke string) (RunBrowser, tea.Cmd) {
	t.Helper()
	key := tea.Key{}
	switch keystroke {
	case "enter":
		key.Code = tea.KeyEnter
	case "esc":
		key.Code = tea.KeyEscape
	case "up":
		key.Code = tea.KeyUp
	case "down":
		key.Code = tea.KeyDown
	case "ctrl+c":
		key.Code = 'c'
		key.Mod = tea.ModCtrl
	default:
		key.Code = []rune(keystroke)[0]
		key.Text = keystroke
	}
	if got := key.String(); got != keystroke {
		t.Fatalf("synthetic key mismatch: built %q, wanted %q", got, keystroke)
	}
	next, cmd := browser.Update(tea.KeyPressMsg(key))
	return next, cmd
}

func browserViewText(browser RunBrowser) string {
	return stripANSI(browser.View())
}

func browserCursorLine(t *testing.T, browser RunBrowser) string {
	t.Helper()
	for _, line := range strings.Split(browserViewText(browser), "\n") {
		if strings.HasPrefix(line, "> ") {
			return line
		}
	}
	t.Fatalf("no cursor row in view:\n%s", browserViewText(browser))
	return ""
}

func assertQuit(t *testing.T, cmd tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected a quit command, got nil")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", cmd())
	}
}

func TestRunBrowserDefaultShowsActiveRunsNewestFirst(t *testing.T) {
	t.Parallel()
	browser := newTestRunBrowser(t)

	view := browserViewText(browser)

	if !strings.Contains(view, "Run Browser — all repositories — ACTIVE") {
		t.Fatalf("expected header with repository and ACTIVE filter, got:\n%s", view)
	}
	if !strings.Contains(view, "aaaaaaaaaaaaaaaa") || !strings.Contains(view, "bbbbbbbbbbbbbbbb") {
		t.Fatalf("expected both Active Runs, got:\n%s", view)
	}
	if strings.Contains(view, "cccccccccccccccc") {
		t.Fatalf("expected terminal Run hidden by default, got:\n%s", view)
	}
	if strings.Contains(view, "run_20260707T144800Z") {
		t.Fatalf("expected short run ids only, got:\n%s", view)
	}
	if !strings.Contains(view, "12m ago") || !strings.Contains(view, "running 12m") {
		t.Fatalf("expected relative start and running duration, got:\n%s", view)
	}
	if !strings.Contains(view, "↑↓ move  enter attach  a all/active  q quit") {
		t.Fatalf("expected footer keys, got:\n%s", view)
	}
	cursor := browserCursorLine(t, browser)
	if !strings.Contains(cursor, "aaaaaaaaaaaaaaaa") {
		t.Fatalf("expected cursor on the newest Run, got %q", cursor)
	}
}

func TestRunBrowserToggleShowsAllAndDistinguishesActive(t *testing.T) {
	t.Parallel()
	browser := newTestRunBrowser(t)

	browser, _ = pressBrowser(t, browser, "a")

	view := browserViewText(browser)
	if !strings.Contains(view, "Run Browser — all repositories — ALL") {
		t.Fatalf("expected ALL filter in header, got:\n%s", view)
	}
	if !strings.Contains(view, "cccccccccccccccc") || !strings.Contains(view, "42m") {
		t.Fatalf("expected terminal Run with duration after toggle, got:\n%s", view)
	}
	styledActive := false
	for _, line := range strings.Split(browser.View(), "\n") {
		stripped := stripANSI(line)
		switch {
		case strings.Contains(stripped, "aaaaaaaaaaaaaaaa"):
			if line != stripped {
				styledActive = true
			}
		case strings.Contains(stripped, "cccccccccccccccc"):
			if line != stripped {
				t.Fatalf("expected terminal row unstyled, got %q", line)
			}
		}
	}
	if !styledActive {
		t.Fatalf("expected Active row visually distinct under the all filter, got:\n%s", browser.View())
	}

	browser, _ = pressBrowser(t, browser, "a")
	view = browserViewText(browser)
	if !strings.Contains(view, "— ACTIVE") || strings.Contains(view, "cccccccccccccccc") {
		t.Fatalf("expected toggle back to Active-only, got:\n%s", view)
	}
}

func TestRunBrowserUsesPerInstanceIdentityTokens(t *testing.T) {
	t.Parallel()
	active, all := browserTestRuns()
	browser := NewRunBrowserWithTokens(active, all, ResolveTokens(false))
	browser.now = browserTestNow
	browser, _ = pressBrowser(t, browser, "a")

	view := browser.View()

	if strings.Contains(view, "\x1b") {
		t.Fatalf("expected identity tokens to render without ANSI, got:\n%q", view)
	}
	if !strings.Contains(view, "Run Browser — all repositories — ALL") {
		t.Fatalf("expected unstyled Run Browser header, got:\n%s", view)
	}
	if !strings.Contains(view, "cccccccccccccccc") {
		t.Fatalf("expected terminal Run in all view, got:\n%s", view)
	}
}

func TestRunBrowserNavigationMovesAndClamps(t *testing.T) {
	t.Parallel()
	browser := newTestRunBrowser(t)

	browser, _ = pressBrowser(t, browser, "up")
	if cursor := browserCursorLine(t, browser); !strings.Contains(cursor, "aaaaaaaaaaaaaaaa") {
		t.Fatalf("expected top clamp to keep the newest Run, got %q", cursor)
	}

	browser, _ = pressBrowser(t, browser, "down")
	if cursor := browserCursorLine(t, browser); !strings.Contains(cursor, "bbbbbbbbbbbbbbbb") {
		t.Fatalf("expected cursor on the second Run, got %q", cursor)
	}

	browser, _ = pressBrowser(t, browser, "down")
	if cursor := browserCursorLine(t, browser); !strings.Contains(cursor, "bbbbbbbbbbbbbbbb") {
		t.Fatalf("expected bottom clamp on the oldest Active Run, got %q", cursor)
	}

	browser, _ = pressBrowser(t, browser, "up")
	if cursor := browserCursorLine(t, browser); !strings.Contains(cursor, "aaaaaaaaaaaaaaaa") {
		t.Fatalf("expected cursor back on the newest Run, got %q", cursor)
	}
}

func TestRunBrowserEnterReportsSelectedRunID(t *testing.T) {
	t.Parallel()
	browser := newTestRunBrowser(t)

	browser, _ = pressBrowser(t, browser, "down")
	browser, cmd := pressBrowser(t, browser, "enter")

	assertQuit(t, cmd)
	outcome := browser.Outcome()
	if outcome.Cancelled {
		t.Fatal("expected selection, got cancel")
	}
	if outcome.RunID != "run_20260707T143000Z_bbbbbbbbbbbbbbbb" {
		t.Fatalf("expected the selected full run id, got %q", outcome.RunID)
	}
}

func TestRunBrowserCancelKeysReportCancel(t *testing.T) {
	t.Parallel()
	for _, keystroke := range []string{"q", "esc", "ctrl+c"} {
		t.Run(keystroke, func(t *testing.T) {
			browser := newTestRunBrowser(t)

			browser, cmd := pressBrowser(t, browser, keystroke)

			assertQuit(t, cmd)
			outcome := browser.Outcome()
			if !outcome.Cancelled || outcome.RunID != "" {
				t.Fatalf("expected cancel outcome with no selection, got %+v", outcome)
			}
		})
	}
}

func TestRunBrowserEmptyStatesNameTheFilter(t *testing.T) {
	t.Parallel()
	_, all := browserTestRuns()
	browser := NewRunBrowser(nil, all)
	browser.now = browserTestNow

	view := browserViewText(browser)
	if !strings.Contains(view, "No active Runs — press a to include terminal Runs.") {
		t.Fatalf("expected the toggle invitation, got:\n%s", view)
	}

	browser, cmd := pressBrowser(t, browser, "enter")
	if cmd != nil {
		t.Fatal("expected Enter on an empty list to do nothing")
	}
	if outcome := browser.Outcome(); outcome.RunID != "" || outcome.Cancelled {
		t.Fatalf("expected no outcome on empty Enter, got %+v", outcome)
	}

	browser, _ = pressBrowser(t, browser, "a")
	if view := browserViewText(browser); !strings.Contains(view, "cccccccccccccccc") {
		t.Fatalf("expected history after toggle, got:\n%s", view)
	}

	empty := NewRunBrowser(nil, nil)
	empty, _ = pressBrowser(t, empty, "a")
	if view := browserViewText(empty); !strings.Contains(view, "No Runs found.") {
		t.Fatalf("expected all-filter empty state, got:\n%s", view)
	}
}

func TestRunBrowserNarrowWidthDropsColumnsInOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		width    int
		contains []string
		excludes []string
	}{
		{
			name:     "full width keeps every column",
			width:    120,
			contains: []string{"codex", "feature/one", "spec:0020-run-browser", "12m ago", "running 12m"},
		},
		{
			name:     "medium width drops branch, agent, then start",
			width:    90,
			contains: []string{"spec:0020-run-browser", "implement", "running 12m"},
			excludes: []string{"feature/one", "codex", "12m ago"},
		},
		{
			name:     "minimal width keeps id, state, and duration",
			width:    60,
			contains: []string{"aaaaaaaaaaaaaaaa", "ResolvingWithAgent", "running 12m"},
			excludes: []string{"feature/one", "codex", "12m ago", "spec:0020-run-browser", "implement"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			browser := newTestRunBrowser(t)
			browser, _ = browser.Update(tea.WindowSizeMsg{Width: tt.width, Height: 40})

			view := browserViewText(browser)

			for _, line := range strings.Split(view, "\n") {
				if got := displayWidth(line); got > tt.width {
					t.Fatalf("line overflows width %d (%d): %q", tt.width, got, line)
				}
			}
			for _, want := range tt.contains {
				if !strings.Contains(view, want) {
					t.Fatalf("expected %q at width %d, got:\n%s", want, tt.width, view)
				}
			}
			for _, exclude := range tt.excludes {
				if strings.Contains(view, exclude) {
					t.Fatalf("expected %q dropped at width %d, got:\n%s", exclude, tt.width, view)
				}
			}
		})
	}
}

func TestRunBrowserSmallHeightKeepsCursorVisible(t *testing.T) {
	t.Parallel()
	browser := newTestRunBrowser(t)
	browser, _ = pressBrowser(t, browser, "a")
	browser, _ = browser.Update(tea.WindowSizeMsg{Width: 120, Height: 6})

	browser, _ = pressBrowser(t, browser, "down")
	browser, _ = pressBrowser(t, browser, "down")

	view := browserViewText(browser)
	if !strings.Contains(view, "cccccccccccccccc") {
		t.Fatalf("expected cursor row visible in the height window, got:\n%s", view)
	}
	if strings.Contains(view, "aaaaaaaaaaaaaaaa") {
		t.Fatalf("expected the newest row scrolled out of a two-row window, got:\n%s", view)
	}
	if lines := strings.Count(view, "\n") + 1; lines > 6 {
		t.Fatalf("expected at most 6 lines, got %d:\n%s", lines, view)
	}
}

func TestRunBrowserProgramDelegatesToTheModel(t *testing.T) {
	t.Parallel()
	browser := newTestRunBrowser(t)
	var program tea.Model = runBrowserProgram{browser: browser}

	view := program.(runBrowserProgram).View()
	if !view.AltScreen {
		t.Fatal("expected the browser program to render in the alternate screen")
	}
	if !strings.Contains(stripANSI(view.Content), "Run Browser — all repositories — ACTIVE") {
		t.Fatalf("expected the program view to carry the model view, got:\n%s", view.Content)
	}

	key := tea.Key{Code: tea.KeyEnter}
	next, cmd := program.Update(tea.KeyPressMsg(key))
	assertQuit(t, cmd)
	outcome := next.(runBrowserProgram).browser.Outcome()
	if outcome.Cancelled || outcome.RunID == "" {
		t.Fatalf("expected the program to report the model selection, got %+v", outcome)
	}
}

type unexpectedBrowserModel struct{}

func (unexpectedBrowserModel) Init() tea.Cmd {
	return nil
}

func (unexpectedBrowserModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return unexpectedBrowserModel{}, nil
}

func (unexpectedBrowserModel) View() tea.View {
	return tea.NewView("")
}

func TestBrowserOutcomeFromModelRejectsUnexpectedModel(t *testing.T) {
	t.Parallel()
	_, err := browserOutcomeFromModel(unexpectedBrowserModel{})

	if err == nil {
		t.Fatal("expected unexpected model type error")
	}
	if !strings.Contains(err.Error(), "unexpected model type") {
		t.Fatalf("expected model type in error, got %v", err)
	}
}

func TestFormatRunRowSharedByBothSurfaces(t *testing.T) {
	t.Parallel()
	active, all := browserTestRuns()
	now := browserTestNow

	absolute := FormatRunRow(active[0], now, false, false)
	if absolute[browserColumnStart] != "2026-07-07T14:48:00Z" {
		t.Fatalf("expected absolute UTC start, got %q", absolute[browserColumnStart])
	}
	relative := FormatRunRow(active[0], now, true, false)
	if relative[browserColumnStart] != "12m ago" {
		t.Fatalf("expected relative start, got %q", relative[browserColumnStart])
	}
	if absolute[browserColumnDuration] != relative[browserColumnDuration] {
		t.Fatalf("expected identical durations, got %q vs %q", absolute[browserColumnDuration], relative[browserColumnDuration])
	}
	if absolute[browserColumnID] != active[0].ID {
		t.Fatalf("expected the untruncated run id, got %q", absolute[browserColumnID])
	}

	terminal := FormatRunRow(all[2], now, false, true)
	if terminal[browserColumnDuration] != "42m" {
		t.Fatalf("expected completion duration, got %q", terminal[browserColumnDuration])
	}
	if len(terminal) != browserColumnCount {
		t.Fatalf("expected the repository column included, got %d fields", len(terminal))
	}
}
