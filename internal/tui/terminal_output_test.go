package tui

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

// embeddedWriter stands for any wrapper that embeds a writer, such as the
// CLI's command environment wrapper: it is not the terminal file.
type embeddedWriter struct {
	io.Writer
}

func nonFileWriters() []struct {
	name     string
	output   io.Writer
	typeName string
} {
	return []struct {
		name     string
		output   io.Writer
		typeName string
	}{
		{name: "buffer", output: &bytes.Buffer{}, typeName: "*bytes.Buffer"},
		{name: "embedded writer", output: embeddedWriter{Writer: &bytes.Buffer{}}, typeName: "tui.embeddedWriter"},
	}
}

func TestRunBrowserSessionRefusesAWriterThatIsNotAFile(t *testing.T) {
	t.Parallel()
	for _, tt := range nonFileWriters() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			outcome, err := RunBrowserSession(context.Background(), tt.output, NewRunBrowser(nil, nil))
			want := "run browser session: output " + tt.typeName + " is not a terminal file"
			if err == nil || err.Error() != want {
				t.Fatalf("RunBrowserSession error = %v; want %q", err, want)
			}
			if outcome != (BrowserOutcome{}) {
				t.Fatalf("expected no outcome from a refused writer, got %+v", outcome)
			}
		})
	}
}

func TestRunCockpitRefusesAWriterThatIsNotAFile(t *testing.T) {
	t.Parallel()
	for _, tt := range nonFileWriters() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// A nil Source would fail the model build; the guard runs first.
			err := RunCockpit(context.Background(), tt.output, CockpitConfig{Mode: CockpitAttach, RunID: "run_x"})
			want := "run cockpit: output " + tt.typeName + " is not a terminal file"
			if err == nil || err.Error() != want {
				t.Fatalf("RunCockpit error = %v; want %q", err, want)
			}
		})
	}
}

func TestRunBrowserCancelKeysNameTheirReason(t *testing.T) {
	t.Parallel()
	for _, keystroke := range []string{"q", "esc", "ctrl+c"} {
		t.Run(keystroke, func(t *testing.T) {
			t.Parallel()
			browser, cmd := pressBrowser(t, newTestRunBrowser(t), keystroke)

			assertQuit(t, cmd)
			want := BrowserOutcome{Cancelled: true, Reason: keystroke}
			if got := browser.Outcome(); got != want {
				t.Fatalf("outcome = %+v; want %+v", got, want)
			}
		})
	}
	t.Run("enter", func(t *testing.T) {
		t.Parallel()
		browser, cmd := pressBrowser(t, newTestRunBrowser(t), "enter")

		assertQuit(t, cmd)
		outcome := browser.Outcome()
		if outcome.Cancelled || outcome.RunID == "" || outcome.Reason != "" {
			t.Fatalf("expected a selection with no reason, got %+v", outcome)
		}
	})
}

func TestRunRepositoryNameNamesTheRepository(t *testing.T) {
	t.Parallel()
	queued := store.Run{
		ID:             "run_20261105T120000Z_2f7556b8aaaaaaaa",
		Kind:           store.KindImplement,
		State:          store.StateActive,
		SpecSlug:       "0063-x",
		RepositoryRoot: "/repos/fiscus",
		GitRoot:        "/worktrees/fiscus-2f7556b8/roundfix-deliver-0063-x-411942b2",
		CreatedAt:      time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC),
	}
	for _, tt := range []struct {
		name string
		run  store.Run
		want string
	}{
		{name: "repository root over the Run Worktree", run: queued, want: "fiscus"},
		{name: "only the Git root", run: store.Run{GitRoot: "/repos/quaestor"}, want: "quaestor"},
		{name: "neither", run: store.Run{RepositoryRoot: "  ", GitRoot: ""}, want: "-"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := RunRepositoryName(tt.run); got != tt.want {
				t.Fatalf("RunRepositoryName = %q; want %q", got, tt.want)
			}
		})
	}
	t.Run("the Run Browser names the repository", func(t *testing.T) {
		t.Parallel()
		browser := NewRunBrowserWithTokens([]store.Run{queued}, []store.Run{queued}, ResolveTokens(false))
		browser, _ = browser.Update(tea.WindowSizeMsg{Width: 160, Height: 20})
		view := browserViewText(browser)
		if !strings.Contains(view, "fiscus") || strings.Contains(view, "roundfix-deliver") {
			t.Fatalf("expected the REPO column to name fiscus, not the Run Worktree:\n%s", view)
		}
	})
}
