package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

// cockpitFakeSource extends the fake timeline source with run lookup and
// data-version change signaling.
type cockpitFakeSource struct {
	fakeTimelineSource
	run     store.Run
	version int64
}

func (source *cockpitFakeSource) DataVersion(context.Context) (int64, error) {
	return source.version, nil
}

func (source *cockpitFakeSource) Run(context.Context, string) (store.Run, bool, error) {
	return source.run, true, nil
}

func newTestCockpit(t *testing.T, source *cockpitFakeSource, view LiveRunView) *cockpitModel {
	t.Helper()
	view.Width = 100
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode:   CockpitAttach,
		View:   view,
		RunID:  "run-1",
		Source: source,
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	return model
}

func pressKey(t *testing.T, model *cockpitModel, keystroke string) tea.Cmd {
	t.Helper()
	key := tea.Key{}
	switch keystroke {
	case "tab":
		key.Code = tea.KeyTab
	case "enter":
		key.Code = tea.KeyEnter
	case "esc":
		key.Code = tea.KeyEscape
	case "up":
		key.Code = tea.KeyUp
	case "down":
		key.Code = tea.KeyDown
	case "pgup":
		key.Code = tea.KeyPgUp
	case "pgdown":
		key.Code = tea.KeyPgDown
	case "home":
		key.Code = tea.KeyHome
	case "end":
		key.Code = tea.KeyEnd
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
	_, cmd := model.Update(tea.KeyPressMsg(key))
	return cmd
}

func viewText(model *cockpitModel) string {
	return model.View().Content
}

func sampleIssues(count int) []rounds.Issue {
	issues := []rounds.Issue{}
	for index := 1; index <= count; index++ {
		issues = append(issues, rounds.Issue{
			Path:     fmt.Sprintf("/missing/issue_%03d.md", index),
			Title:    fmt.Sprintf("issue title %03d", index),
			Severity: "major",
			Status:   "pending",
		})
	}
	return issues
}

func TestCockpitTabSwitchesFocusAndArrowsMoveSelection(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	source.addLine("line one\n")
	model := newTestCockpit(t, source, LiveRunView{PipelineState: store.StateActive, Issues: sampleIssues(3)})

	if model.focus != focusTimeline {
		t.Fatalf("expected timeline focus by default, got %v", model.focus)
	}
	pressKey(t, model, "tab")
	if model.focus != focusIssues {
		t.Fatal("expected Tab to focus the Issues pane")
	}
	pressKey(t, model, "down")
	pressKey(t, model, "down")
	pressKey(t, model, "up")
	if model.selected != 1 {
		t.Fatalf("expected selection at index 1, got %d", model.selected)
	}
	pressKey(t, model, "down")
	pressKey(t, model, "down")
	pressKey(t, model, "down")
	if model.selected != 2 {
		t.Fatalf("expected selection clamped at the last issue, got %d", model.selected)
	}
	if !strings.Contains(viewText(model), "> major") {
		t.Fatalf("expected selection marker rendered, got:\n%s", viewText(model))
	}
	pressKey(t, model, "tab")
	if model.focus != focusTimeline {
		t.Fatal("expected Tab to return focus to the timeline")
	}
}

func TestCockpitEnterOpensIssueDetailAndEscCloses(t *testing.T) {
	artifactDir := t.TempDir()
	persisted, err := rounds.PersistRound(context.Background(), rounds.PersistRequest{
		ArtifactDir:    artifactDir,
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		HeadSHA:        "abc123",
		Round:          1,
		CreatedAt:      time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Items: []reviewsource.ReviewItem{{
			Title:                   "major: handle nil cache",
			File:                    "internal/cache/cache.go",
			Line:                    42,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Guard the map lookup before dereferencing.",
			SourceRef:               "thread:PRRT_1,comment:PRRC_1",
			ReviewHash:              "abc",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		}},
	})
	if err != nil {
		t.Fatalf("persist round: %v", err)
	}
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateClean}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateClean,
		Issues:        []rounds.Issue{{Path: persisted.IssuePaths[0], Title: "major: handle nil cache", Severity: "major", Status: "pending"}},
	})

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")

	if model.detail == nil {
		t.Fatal("expected Enter to open the issue detail pane")
	}
	rendered := viewText(model)
	for _, expected := range []string{"REVIEW.ISSUE", "major: handle nil cache", "thread:PRRT_1,comment:PRRC_1", "Guard the map lookup"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected detail to contain %q, got:\n%s", expected, rendered)
		}
	}

	pressKey(t, model, "esc")
	if model.detail != nil {
		t.Fatal("expected Esc to close the detail pane")
	}
}

func TestCockpitMissingArtifactDegradesWithoutFailing(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateClean}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateClean,
		Issues:        []rounds.Issue{{Path: filepath.Join(t.TempDir(), "gone.md"), Title: "cleaned", Severity: "minor", Status: "resolved"}},
	})

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")

	if model.detail == nil || !model.detail.missing {
		t.Fatal("expected missing artifact to open a degraded detail pane")
	}
	if !strings.Contains(viewText(model), "artifact not available") {
		t.Fatalf("expected degraded notice, got:\n%s", viewText(model))
	}
}

func TestCockpitDetachKeysQuitInAttachModeAndRunIsUntouched(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{PipelineState: store.StateActive})

	for _, keystroke := range []string{"q", "ctrl+c"} {
		cmd := pressKey(t, model, keystroke)
		if cmd == nil {
			t.Fatalf("expected %q to detach in attach mode", keystroke)
		}
		if msg := cmd(); fmt.Sprintf("%T", msg) != "tea.QuitMsg" {
			t.Fatalf("expected quit command for %q, got %T", keystroke, msg)
		}
	}
	if source.run.State != store.StateActive {
		t.Fatal("expected detach to leave the Run untouched")
	}
}

func TestCockpitScrollFreezesFollowAndStatusBarNarratesStates(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	for index := 1; index <= 60; index++ {
		source.addLine(fmt.Sprintf("line %04d\n", index))
	}
	model := newTestCockpit(t, source, LiveRunView{PipelineState: store.StateActive})

	if !strings.Contains(viewText(model), "FOLLOWING") {
		t.Fatalf("expected FOLLOWING status after replay, got:\n%s", viewText(model))
	}

	pressKey(t, model, "pgup")
	source.addLine("line 0061\n")
	source.version++
	model.Update(cockpitTickMsg{})

	rendered := viewText(model)
	if !strings.Contains(rendered, "SCROLLED · 1 new event(s) below — End to follow") {
		t.Fatalf("expected scrolled status with arrival count, got:\n%s", rendered)
	}

	pressKey(t, model, "end")
	rendered = viewText(model)
	if !strings.Contains(rendered, "FOLLOWING") || !strings.Contains(rendered, "line 0061") {
		t.Fatalf("expected End to resume following at the tail, got:\n%s", rendered)
	}
}

func TestCockpitTickPollsViewportOnlyOnDataVersionChange(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 7}
	source.addLine("line one\n")
	model := newTestCockpit(t, source, LiveRunView{PipelineState: store.StateActive})
	model.Update(cockpitTickMsg{}) // first tick drains the replay gap
	readsAfterFirst := source.reads

	model.Update(cockpitTickMsg{})
	model.Update(cockpitTickMsg{})

	if source.reads != readsAfterFirst {
		t.Fatalf("expected idle ticks to skip event reads, got %d extra", source.reads-readsAfterFirst)
	}
}

func TestCockpitTerminalRunShowsReadOnlyAndStopsTicking(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateClean}, version: 1}
	source.addLine("line one\n")
	model := newTestCockpit(t, source, LiveRunView{PipelineState: store.StateClean})

	rendered := viewText(model)
	if !strings.Contains(rendered, "READ-ONLY") || !strings.Contains(rendered, "CLEAN") {
		t.Fatalf("expected terminal read-only status, got:\n%s", rendered)
	}
	if cmd := model.Init(); cmd != nil {
		t.Fatal("expected no follow ticking for terminal Runs")
	}
	if !strings.Contains(rendered, "q detach") {
		t.Fatalf("expected attach footer keys, got:\n%s", rendered)
	}
}

func TestCockpitOwningModeKeysDifferFromAttach(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	stopped := false
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode:   CockpitOwning,
		View:   LiveRunView{PipelineState: store.StateActive, Width: 100},
		RunID:  "run-1",
		Source: source,
		OnStop: func() { stopped = true },
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}

	if cmd := pressKey(t, model, "q"); cmd != nil {
		t.Fatal("expected q to do nothing in owning mode")
	}
	if cmd := pressKey(t, model, "ctrl+c"); cmd != nil {
		t.Fatal("expected ctrl+c to not quit the owning cockpit directly")
	}
	if !stopped {
		t.Fatal("expected ctrl+c to trigger the Stop Request callback in owning mode")
	}
	if !strings.Contains(viewText(model), "Ctrl-C stop") {
		t.Fatalf("expected owning footer keys, got:\n%s", viewText(model))
	}
}

func TestOwningCockpitPollsJournalWhileOwnProcessWrites(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	writer, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	t.Cleanup(func() { _ = writer.Close() })
	run, err := writer.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        t.TempDir(),
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(t.TempDir(), ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	stopRequested := false
	model, err := newCockpitModel(ctx, CockpitConfig{
		Mode:   CockpitOwning,
		View:   LiveRunView{PipelineState: store.StateActive, Width: 100},
		RunID:  run.ID,
		Source: reader,
		OnStop: func() { stopRequested = true },
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}

	const total = 25
	written := make(chan error, 1)
	go func() {
		for index := 0; index < total; index++ {
			if _, err := writer.AppendRunEvent(ctx, runevent.RunEvent{
				RunID:   run.ID,
				Source:  runevent.SourceAgent,
				Kind:    runevent.KindAgentRaw,
				Summary: fmt.Sprintf("live line %02d\n", index),
				Payload: []byte(fmt.Sprintf(`{"text":"live line %02d\n"}`, index)),
			}); err != nil {
				written <- err
				return
			}
		}
		_, err := writer.CompleteRun(ctx, run.ID, store.StateStopped)
		written <- err
	}()

	// The owning cockpit polls its read-only connection while the same
	// process writes — a Stop Request mid-poll must keep working.
	sawAll := false
	for tick := 0; tick < 10_000 && !sawAll; tick++ {
		model.Update(cockpitTickMsg{})
		if tick == 5 {
			pressKey(t, model, "ctrl+c")
		}
		rendered := viewText(model)
		sawAll = strings.Contains(rendered, fmt.Sprintf("live line %02d", total-1)) && strings.Contains(rendered, "READ-ONLY")
	}
	if err := <-written; err != nil {
		t.Fatalf("writer: %v", err)
	}
	model.Update(cockpitTickMsg{})

	if !sawAll {
		t.Fatalf("expected the owning cockpit to render journal writes and the terminal state, got:\n%s", viewText(model))
	}
	if !stopRequested {
		t.Fatal("expected ctrl+c during active polling to trigger the Stop Request callback")
	}
}
