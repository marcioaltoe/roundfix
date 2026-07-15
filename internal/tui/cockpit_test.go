package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
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

// addDaemonEvent appends one daemon-source journal event; the timeline
// renders daemon kinds from their bounded summaries.
func (source *cockpitFakeSource) addDaemonEvent(kind runevent.Kind, summary string) {
	cursor := int64(len(source.events) + 1)
	source.events = append(source.events, store.JournalEvent{
		Cursor: cursor,
		Event:  runevent.RunEvent{Source: runevent.SourceDaemon, Kind: kind, Summary: summary},
	})
}

func (source *cockpitFakeSource) addTaskEvent(t *testing.T, taskID string, phase string, status spec.Status, batch int) {
	t.Helper()
	payload := map[string]any{"task": taskID, "phase": phase}
	if status != "" {
		payload["status"] = string(status)
	}
	if batch > 0 {
		payload["batch"] = batch
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal task event: %v", err)
	}
	summary := fmt.Sprintf("Task %s %s.", taskID, phase)
	switch phase {
	case "started":
		summary = fmt.Sprintf("Task %s started as Batch %03d.", taskID, batch)
	case "settled":
		summary = fmt.Sprintf("Task %s settled %s.", taskID, status)
	case "skipped":
		summary = fmt.Sprintf("Task %s skipped.", taskID)
	}
	source.addEvent(runevent.RunEvent{
		Batch:       batch,
		Source:      runevent.SourceDaemon,
		Kind:        runevent.KindDaemonTask,
		ReviewIssue: taskID,
		Summary:     summary,
		Payload:     raw,
	})
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
	return stripANSI(model.View().Content)
}

func footerText(model *cockpitModel, width int) string {
	return strings.TrimRight(stripANSI(model.renderFooter(width)), " ")
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

// Suite: cockpit render snapshots
// Invariant: the shipped Live Run View visible render string stays byte-identical while render helpers are extracted.
// Boundary IN: synchronous cockpit model rendering for review and spec Runs.
// Boundary OUT: terminal emulation, ANSI color assertions, daemon/store behavior covered by lower suites.
func TestCockpitRenderSnapshots(t *testing.T) {
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{name: "88x24", width: 88, height: 24},
		{name: "120x40", width: 120, height: 40},
	}
	cases := []struct {
		name  string
		model func(t *testing.T) *cockpitModel
	}{
		{
			name: "review_normal",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				return newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
			},
		},
		{
			name: "review_detail_open",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				return newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, true)
			},
		},
		{
			name: "review_attach",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				return newReviewSnapshotCockpit(t, CockpitAttach, store.StateResolvingWithAgent, false)
			},
		},
		{
			name: "review_terminal",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				return newReviewSnapshotCockpit(t, CockpitOwning, store.StateClean, false)
			},
		},
		{
			name: "spec_run_pane",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				return newSpecSnapshotCockpit(t)
			},
		},
	}

	for _, tc := range cases {
		for _, size := range sizes {
			t.Run(tc.name+"/"+size.name, func(t *testing.T) {
				model := tc.model(t)
				model.Update(tea.WindowSizeMsg{Width: size.width, Height: size.height})
				assertCockpitSnapshot(t, tc.name+"_"+size.name, stripANSI(model.View().Content))
			})
		}
	}
}

func newReviewSnapshotCockpit(t *testing.T, mode CockpitMode, runState string, openDetail bool) *cockpitModel {
	t.Helper()
	artifactDir := t.TempDir()
	persisted, err := rounds.PersistRound(context.Background(), rounds.PersistRequest{
		ArtifactDir:    artifactDir,
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/cockpit",
		HeadSHA:        "abc123",
		Round:          1,
		CreatedAt:      time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
		Items: []reviewsource.ReviewItem{
			{
				Title:                   "major: guard nil cache",
				File:                    "internal/cache/cache.go",
				Line:                    42,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "Guard the map lookup before dereferencing.",
				SourceRef:               "thread:PRRT_1,comment:PRRC_1",
				ReviewHash:              "h1",
				SourceReviewID:          "1",
				SourceReviewSubmittedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
			},
			{
				Title:                   "minor: trim stale TODO",
				File:                    "internal/cache/cache_test.go",
				Line:                    17,
				Severity:                "minor",
				Author:                  "coderabbitai[bot]",
				Body:                    "Remove the stale TODO from the test fixture.",
				SourceRef:               "thread:PRRT_2,comment:PRRC_2",
				ReviewHash:              "h2",
				SourceReviewID:          "1",
				SourceReviewSubmittedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
			},
			{
				Title:                   "info: document cache eviction",
				File:                    "internal/cache/cache.go",
				Line:                    88,
				Severity:                "info",
				Author:                  "coderabbitai[bot]",
				Body:                    "Document why the eviction path ignores expired entries.",
				SourceRef:               "thread:PRRT_3,comment:PRRC_3",
				ReviewHash:              "h3",
				SourceReviewID:          "1",
				SourceReviewSubmittedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
			},
		},
	})
	if err != nil {
		t.Fatalf("persist round: %v", err)
	}
	if err := rounds.SetIssueStatus(persisted.IssuePaths[0], rounds.StatusResolved, "", ""); err != nil {
		t.Fatalf("set issue status: %v", err)
	}
	issues := make([]rounds.Issue, 0, len(persisted.IssuePaths))
	for _, path := range persisted.IssuePaths {
		issues = append(issues, rounds.Issue{Path: path, Status: rounds.StatusPending})
	}
	clock := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	source := &cockpitFakeSource{run: store.Run{ID: "run-review-00000001", State: runState}, version: 1}
	source.addLine("PLAN inspect current cockpit render\n")
	source.addLine("[TOOL] read_file * completed\n")
	source.addLine("THINK preserving the shipped layout\n")
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonBatch,
		Summary: "Batch 001 executing.",
		Time:    clock,
	})
	source.addLine("SESSION running\n")
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode: mode,
		View: LiveRunView{
			Command:       "resolve",
			Repository:    "owner/project",
			PRNumber:      "123",
			HeadBranch:    "feature/cockpit",
			ReviewSource:  string(reviewsource.SourceCodeRabbit),
			Agent:         "codex",
			Model:         "gpt-5",
			RunID:         "run-review-00000001",
			PipelineState: runState,
			Issues:        issues,
			BatchSizes:    []int{2, 1},
		},
		RunID:        "run-review-00000001",
		Source:       source,
		ColorEnabled: true,
		OnStop:       func() {},
		Now:          func() time.Time { return clock },
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	clock = clock.Add(83 * time.Second)
	if openDetail {
		pressKey(t, model, "tab")
		pressKey(t, model, "enter")
	}
	return model
}

func newSpecSnapshotCockpit(t *testing.T) *cockpitModel {
	t.Helper()
	gitRoot := t.TempDir()
	slug := "0005-tui-cockpit"
	fileOne := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Decompose renderers", spec.StatusCompleted)
	fileTwo := writeCockpitTaskFile(t, gitRoot, slug, "task_02", "Build phase row", spec.StatusInProgress)
	fileThree := writeCockpitTaskFile(t, gitRoot, slug, "task_03", "Upgrade work queue", spec.StatusPending)
	source := &cockpitFakeSource{run: store.Run{ID: "run-spec-00000001", State: store.StateResolvingWithAgent}, version: 1}
	source.addDaemonEvent(runevent.KindDaemonTask, "Task task_01 settled completed.")
	source.addLine("PLAN continue Task Graph execution\n")
	source.addLine("[TOOL] go test ./internal/tui/ * completed\n")
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode: CockpitOwning,
		View: LiveRunView{
			Command:       "implement",
			HeadBranch:    "feature/cockpit",
			Agent:         "codex",
			Model:         "gpt-5",
			RunID:         "run-spec-00000001",
			PipelineState: store.StateResolvingWithAgent,
			RunKind:       store.KindImplement,
			SpecSlug:      slug,
			GitRoot:       gitRoot,
			Tasks: []spec.Task{
				{ID: "task_01", File: fileOne, Title: "Decompose renderers", Status: spec.StatusPending},
				{ID: "task_02", File: fileTwo, Title: "Build phase row", Status: spec.StatusPending},
				{ID: "task_03", File: fileThree, Title: "Upgrade work queue", Status: spec.StatusPending},
			},
		},
		RunID:        "run-spec-00000001",
		Source:       source,
		ColorEnabled: true,
		OnStop:       func() {},
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	return model
}

func assertCockpitSnapshot(t *testing.T, name string, got string) {
	t.Helper()
	path := filepath.Join("testdata", "cockpit_snapshots", name+".golden")
	if os.Getenv("ROUNDFIX_UPDATE_COCKPIT_SNAPSHOTS") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create cockpit snapshot dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write cockpit snapshot: %v", err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read cockpit snapshot %s: %v", path, err)
	}
	if got != string(want) {
		t.Fatalf("cockpit snapshot %s mismatch\n%s", name, cockpitSnapshotDiff(string(want), got))
	}
}

func cockpitSnapshotDiff(want string, got string) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")
	limit := maxInt(len(wantLines), len(gotLines))
	for index := 0; index < limit; index++ {
		wantLine := "<missing>"
		gotLine := "<missing>"
		if index < len(wantLines) {
			wantLine = wantLines[index]
		}
		if index < len(gotLines) {
			gotLine = gotLines[index]
		}
		if wantLine != gotLine {
			return fmt.Sprintf("line %d:\nwant: %q\ngot:  %q\n\nfull got:\n%s", index+1, wantLine, gotLine, got)
		}
	}
	return "snapshot lengths differ"
}

// Suite: cockpit base layout and phases
// Invariant: review and spec Runs render one Work Queue, one dominant timeline, and the correct phase sequence.
// Boundary IN: synchronous cockpit model rendering and journal-backed phase derivation.
// Boundary OUT: rich Work Queue rows, timeline grouping, and modal detail covered by later cockpit tasks.
func TestCockpitReviewPhaseRowAndTwoPaneStructure(t *testing.T) {
	model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	for _, expected := range []string{
		"WORK QUEUE (3)",
		"SESSION.TIMELINE",
		"FETCH [done]",
		"TRIAGE [done]",
		"AGENT [run]",
		"VERIFY [wait]",
		"PUSH [locked]",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected cockpit render to contain %q, got:\n%s", expected, rendered)
		}
	}
	assertContainsInOrder(t, rendered,
		"FETCH [done]",
		"TRIAGE [done]",
		"AGENT [run]",
		"VERIFY [wait]",
		"PUSH [locked]",
	)
	for _, unexpected := range []string{"RUN.PROGRESS", "REVIEW.ISSUES"} {
		if strings.Contains(rendered, unexpected) {
			t.Fatalf("expected %q to be retired from the base cockpit, got:\n%s", unexpected, rendered)
		}
	}
}

func TestCockpitSpecPhaseRowCoversQAOmittedAndLocked(t *testing.T) {
	tests := []struct {
		name       string
		qaEvent    bool
		want       []string
		notWant    []string
		taskStatus []spec.Status
	}{
		{
			name:       "qa omitted without qa journal event",
			taskStatus: []spec.Status{spec.StatusCompleted, spec.StatusInProgress, spec.StatusPending},
			want:       []string{"AGENT [run]", "VERIFY [wait]", "COMMIT [wait]"},
			notWant:    []string{"QA ["},
		},
		{
			name:       "qa locked until every task completes",
			qaEvent:    true,
			taskStatus: []spec.Status{spec.StatusCompleted, spec.StatusInProgress, spec.StatusPending},
			want:       []string{"AGENT [run]", "VERIFY [wait]", "COMMIT [wait]", "QA [locked]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := newSpecPhaseCockpit(t, tt.qaEvent, store.StateResolvingWithAgent, tt.taskStatus...)
			model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			rendered := viewText(model)
			for _, expected := range append([]string{"WORK QUEUE (3)", "SESSION.TIMELINE"}, tt.want...) {
				if !strings.Contains(rendered, expected) {
					t.Fatalf("expected spec cockpit to contain %q, got:\n%s", expected, rendered)
				}
			}
			assertContainsInOrder(t, rendered, tt.want...)
			for _, unexpected := range tt.notWant {
				if strings.Contains(rendered, unexpected) {
					t.Fatalf("expected spec cockpit not to contain %q, got:\n%s", unexpected, rendered)
				}
			}
			if strings.Contains(rendered, "SPEC.TASKS") {
				t.Fatalf("expected spec Run to use the shared Work Queue pane, got:\n%s", rendered)
			}
		})
	}
}

func TestCockpitTimelinePaneIsDominantAtTestedSizes(t *testing.T) {
	for _, size := range []struct {
		name   string
		width  int
		height int
	}{
		{name: "88x24", width: 88, height: 24},
		{name: "120x40", width: 120, height: 40},
	} {
		t.Run(size.name, func(t *testing.T) {
			model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
			model.Update(tea.WindowSizeMsg{Width: size.width, Height: size.height})

			layout := cockpitLayoutFor(model)
			if layout.rightWidth <= layout.sidebarWidth {
				t.Fatalf("expected timeline width to dominate at %s, got queue=%d timeline=%d", size.name, layout.sidebarWidth, layout.rightWidth)
			}
		})
	}
}

func newSpecPhaseCockpit(t *testing.T, qaEvent bool, runState string, statuses ...spec.Status) *cockpitModel {
	t.Helper()
	gitRoot := t.TempDir()
	slug := "0005-tui-cockpit"
	tasks := make([]spec.Task, 0, len(statuses))
	for index, status := range statuses {
		id := fmt.Sprintf("task_%02d", index+1)
		title := fmt.Sprintf("Task %02d", index+1)
		file := writeCockpitTaskFile(t, gitRoot, slug, id, title, status)
		tasks = append(tasks, spec.Task{ID: id, File: file, Title: title, Status: spec.StatusPending})
	}
	source := &cockpitFakeSource{run: store.Run{ID: "run-spec-phase", State: runState}, version: 1}
	source.addLine("PLAN render spec cockpit phases\n")
	if qaEvent {
		source.addDaemonEvent(runevent.KindDaemonQA, "QA requested for Spec 0005-tui-cockpit.")
	}
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode: CockpitOwning,
		View: LiveRunView{
			Command:       "implement",
			RunKind:       store.KindImplement,
			SpecSlug:      slug,
			GitRoot:       gitRoot,
			Tasks:         tasks,
			HeadBranch:    "feature/cockpit",
			Agent:         "codex",
			Model:         "gpt-5",
			RunID:         "run-spec-phase",
			PipelineState: runState,
		},
		RunID:  "run-spec-phase",
		Source: source,
		OnStop: func() {},
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	return model
}

func assertContainsInOrder(t *testing.T, haystack string, needles ...string) {
	t.Helper()
	offset := 0
	for _, needle := range needles {
		index := strings.Index(haystack[offset:], needle)
		if index < 0 {
			t.Fatalf("expected %q after byte offset %d in:\n%s", needle, offset, haystack)
		}
		offset += index + len(needle)
	}
}

// assertTaskQueueRow checks the Work Queue card marker for one Task. Only
// the left pane is scanned: the timeline narrates the same Task ids, and
// pane rows do not stay vertically aligned across the two panels.
func assertTaskQueueRow(t *testing.T, rendered string, taskID string, marker string) {
	t.Helper()
	lines := strings.Split(rendered, "\n")
	for index := range lines {
		if cut := strings.Index(lines[index], "││"); cut >= 0 {
			lines[index] = lines[index][:cut]
		}
	}
	for index, line := range lines {
		if !strings.Contains(line, taskID+" ") {
			continue
		}
		if index == 0 {
			t.Fatalf("expected %s title line to have a marker line before it in:\n%s", taskID, rendered)
		}
		if !strings.Contains(lines[index-1], marker) {
			t.Fatalf("expected %s row marker %q, got marker line %q in:\n%s", taskID, marker, lines[index-1], rendered)
		}
		return
	}
	t.Fatalf("expected rendered Work Queue to contain %s, got:\n%s", taskID, rendered)
}

func TestCockpitWorkQueueRowsRenderMarkersMetadataAndOptionalSeverity(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		item     WorkItem
		want     []string
		notWant  []string
		selected bool
	}{
		{
			name:     "executing review issue with severity",
			label:    "Executing",
			selected: true,
			item: WorkItem{
				Name:     "Issue #007",
				Title:    "Guard nil cache",
				Severity: "major",
				Ordinal:  7,
				Location: "internal/cache/cache.go:42",
			},
			want: []string{"> [run]", "MAJOR", "#7", "Guard nil cache", "internal/cache/cache.go:42"},
		},
		{
			name:  "resolved review issue",
			label: "Resolved",
			item:  WorkItem{Name: "Issue #001", Title: "Resolved issue", Severity: "minor", Ordinal: 1, Location: "a.go:1"},
			want:  []string{"[done]", "MINOR", "#1", "Resolved issue", "a.go:1"},
		},
		{
			name:    "completed task without severity",
			label:   "Completed",
			item:    WorkItem{Name: "task_02", Title: "Wire API", Ordinal: 2, Location: "docs/specs/0001-widget-flow/task_02.md"},
			want:    []string{"[done]", "#2", "task_02", "Wire API", "docs/specs/0001-widget-flow/task_02.md"},
			notWant: []string{"MAJOR", "MINOR", "LOW", "HIGH"},
		},
		{
			name:    "waiting task without severity",
			label:   "Waiting",
			item:    WorkItem{Name: "task_03", Title: "Upgrade queue", Ordinal: 3},
			want:    []string{"[wait]", "#3", "task_03", "Upgrade queue"},
			notWant: []string{"MAJOR", "MINOR", "LOW", "HIGH"},
		},
		{
			name:  "paused row",
			label: "Paused",
			item:  WorkItem{Name: "Issue #004", Title: "Paused issue", Ordinal: 4},
			want:  []string{"[locked]", "#4", "Paused issue"},
		},
		{
			name:  "invalid review issue",
			label: "Invalid",
			item:  WorkItem{Name: "Issue #005", Title: "Invalid issue", Ordinal: 5},
			want:  []string{"[invalid]", "#5", "Invalid issue"},
		},
		{
			name:  "duplicated review issue",
			label: "Duplicated",
			item:  WorkItem{Name: "Issue #006", Title: "Duplicated issue", Ordinal: 6},
			want:  []string{"[dup]", "#6", "Duplicated issue"},
		},
		{
			name:  "failed row",
			label: "Failed",
			item:  WorkItem{Name: "Issue #008", Title: "Failed issue", Ordinal: 8},
			want:  []string{"[fail]", "#8", "Failed issue"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &cockpitModel{tokens: ResolveTokens(true)}
			if tt.selected {
				model.selected = 3
			} else {
				model.selected = -1
			}
			rendered := stripANSI(strings.Join(model.workItemBlock(tt.item, tt.label, 3, 46), "\n"))
			for _, expected := range tt.want {
				if !strings.Contains(rendered, expected) {
					t.Fatalf("expected work queue row to contain %q, got:\n%s", expected, rendered)
				}
			}
			for _, unexpected := range tt.notWant {
				if strings.Contains(rendered, unexpected) {
					t.Fatalf("expected work queue row not to contain %q, got:\n%s", unexpected, rendered)
				}
			}
		})
	}
}

func TestCockpitWorkQueueBatchSeparatorShowsOrdinalAndElapsed(t *testing.T) {
	startedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	clock := startedAt.Add(83 * time.Second)
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonBatch,
		Summary: "Batch 001 executing.",
		Time:    startedAt,
	})
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode: CockpitAttach,
		View: LiveRunView{
			PipelineState: store.StateResolvingWithAgent,
			Width:         100,
			Issues:        sampleIssues(3),
			BatchSizes:    []int{2, 1},
		},
		RunID:  "run-1",
		Source: source,
		Now:    func() time.Time { return clock },
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}

	executing := stripANSI(model.batchSeparator(0, 46))
	for _, expected := range []string{"BATCH 001/002", "01:23"} {
		if !strings.Contains(executing, expected) {
			t.Fatalf("expected executing separator to contain %q, got %q", expected, executing)
		}
	}
	if strings.Contains(executing, "───") {
		t.Fatalf("expected dense batch separator without decorative rule, got %q", executing)
	}

	waiting := stripANSI(model.batchSeparator(2, 46))
	if !strings.Contains(waiting, "BATCH 002/002") || strings.Contains(waiting, "01:23") {
		t.Fatalf("expected waiting separator without elapsed time, got %q", waiting)
	}
}

func TestCockpitWorkQueueFooterTotalsForRunKinds(t *testing.T) {
	review := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
	review.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	reviewRendered := viewText(review)
	if !strings.Contains(reviewRendered, "3 issues total · 1 resolved · 2 unresolved") {
		t.Fatalf("expected review totals footer, got:\n%s", reviewRendered)
	}

	specRun := newSpecPhaseCockpit(t, false, store.StateResolvingWithAgent, spec.StatusCompleted, spec.StatusInProgress, spec.StatusPending)
	specRun.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	specRendered := viewText(specRun)
	if !strings.Contains(specRendered, "3 Tasks total · 1 completed · 2 unresolved") {
		t.Fatalf("expected spec totals footer, got:\n%s", specRendered)
	}
}

func TestCockpitFooterHintsForStatesAndRunKinds(t *testing.T) {
	tests := []struct {
		name  string
		model func(t *testing.T) *cockpitModel
		want  string
	}{
		{
			name: "review normal owning",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter issue · D show detail · End follow · Ctrl-C stop",
		},
		{
			name: "spec normal owning",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newSpecPhaseCockpit(t, false, store.StateResolvingWithAgent, spec.StatusInProgress)
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter Task · D show detail · End follow · Ctrl-C stop",
		},
		{
			name: "review modal owning",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, true)
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Esc close · j/k scroll · PgUp/PgDn page · Ctrl-C stop",
		},
		{
			name: "spec modal owning",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newSpecPhaseCockpit(t, false, store.StateResolvingWithAgent, spec.StatusInProgress)
				model.detail = &issueDetailView{kind: detailTask, task: model.cfg.View.Tasks[0], lines: []string{"body"}}
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Esc close · j/k scroll · PgUp/PgDn page · Ctrl-C stop",
		},
		{
			name: "review attach",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newReviewSnapshotCockpit(t, CockpitAttach, store.StateResolvingWithAgent, false)
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter issue · D show detail · End follow · q detach",
		},
		{
			name: "spec attach",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newSpecPhaseCockpit(t, false, store.StateResolvingWithAgent, spec.StatusInProgress)
				model.cfg.Mode = CockpitAttach
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter Task · D show detail · End follow · q detach",
		},
		{
			name: "review terminal",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateClean, false)
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter issue · D show detail · End follow · q close",
		},
		{
			name: "spec terminal",
			model: func(t *testing.T) *cockpitModel {
				t.Helper()
				model := newSpecPhaseCockpit(t, false, store.StateClean, spec.StatusCompleted)
				model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				return model
			},
			want: "Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter Task · D show detail · End follow · q close",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.model(t)
			if got := footerText(model, 180); got != tt.want {
				t.Fatalf("expected footer %q, got %q", tt.want, got)
			}
		})
	}
}

func TestCockpitResponsiveFallbackAndStableSizes(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		collapsed  bool
		sidebar    int
		timeline   int
		bodyHeight int
		want       []string
		notWant    []string
		wantFooter string
	}{
		{
			name:       "80x24 collapsed",
			width:      80,
			height:     24,
			collapsed:  true,
			timeline:   80,
			bodyHeight: 19,
			want:       []string{"SESSION.TIMELINE", "QUEUE.SUMMARY 3 issues total · 1 resolved · 2 unresolved"},
			notWant:    []string{"WORK QUEUE ("},
			wantFooter: "widen for Work Queue",
		},
		{
			name:       "120x40 two pane",
			width:      120,
			height:     40,
			sidebar:    46,
			timeline:   71,
			bodyHeight: 35,
			want:       []string{"WORK QUEUE (3)", "SESSION.TIMELINE"},
			notWant:    []string{"QUEUE.SUMMARY"},
		},
		{
			name:       "200x50 two pane",
			width:      200,
			height:     50,
			sidebar:    46,
			timeline:   151,
			bodyHeight: 45,
			want:       []string{"WORK QUEUE (3)", "SESSION.TIMELINE"},
			notWant:    []string{"QUEUE.SUMMARY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
			model.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			layout := cockpitLayoutFor(model)
			if layout.collapsed != tt.collapsed || layout.sidebarWidth != tt.sidebar || layout.rightWidth != tt.timeline || layout.bodyHeight != tt.bodyHeight {
				t.Fatalf("unexpected layout for %s: %+v", tt.name, layout)
			}
			rendered := viewText(model)
			for _, expected := range tt.want {
				if !strings.Contains(rendered, expected) {
					t.Fatalf("expected %s render to contain %q, got:\n%s", tt.name, expected, rendered)
				}
			}
			for _, unexpected := range tt.notWant {
				if strings.Contains(rendered, unexpected) {
					t.Fatalf("expected %s render not to contain %q, got:\n%s", tt.name, unexpected, rendered)
				}
			}
			if tt.wantFooter != "" && !strings.Contains(rendered, tt.wantFooter) {
				t.Fatalf("expected %s footer to contain %q, got:\n%s", tt.name, tt.wantFooter, rendered)
			}
		})
	}
}

func TestCockpitDegenerateSizesRenderEmptyWithoutPanic(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{name: "zero width", width: 0, height: 24},
		{name: "negative width", width: -1, height: 24},
		{name: "zero height", width: 120, height: 0},
		{name: "negative height", width: 120, height: -1},
		{name: "zero both", width: 0, height: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := newReviewSnapshotCockpit(t, CockpitOwning, store.StateResolvingWithAgent, false)
			model.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			if got := viewText(model); got != "" {
				t.Fatalf("expected empty render for %s, got:\n%s", tt.name, got)
			}
		})
	}
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
	rendered := viewText(model)
	for _, expected := range []string{"> [wait]", "#3", "issue title 003"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected selection marker rendered with %q, got:\n%s", expected, rendered)
		}
	}
	pressKey(t, model, "tab")
	if model.focus != focusTimeline {
		t.Fatal("expected Tab to return focus to the timeline")
	}
}

func TestCockpitEnterOpensIssueDetailModalAndEscRestoresReviewContext(t *testing.T) {
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
	extraBody := []string{"", "## Debug notes"}
	for index := 1; index <= 30; index++ {
		extraBody = append(extraBody, fmt.Sprintf("detail body line %02d", index))
	}
	appendToFile(t, persisted.IssuePaths[0], "\n"+strings.Join(extraBody, "\n")+"\n")
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateClean}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateClean,
		Issues:        []rounds.Issue{{Path: persisted.IssuePaths[0], Title: "major: handle nil cache", Severity: "major", Status: "pending"}},
	})

	pressKey(t, model, "tab")
	before := viewText(model)
	pressKey(t, model, "enter")

	if model.detail == nil {
		t.Fatal("expected Enter to open the issue detail modal")
	}
	rendered := viewText(model)
	for _, expected := range []string{"SESSION.TIMELINE", "WORK QUEUE", "REVIEW.ISSUE  #001", "major: handle nil cache", "thread:PRRT_1,comment:PRRC_1", "Guard the map lookup", "Line 1-"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected detail to contain %q, got:\n%s", expected, rendered)
		}
	}

	selected := model.selected
	pressKey(t, model, "down")
	if model.selected != selected {
		t.Fatalf("expected queue selection to stay fixed while modal is open, got %d want %d", model.selected, selected)
	}
	if model.detail.scroll == 0 {
		t.Fatal("expected down to scroll the detail body while modal is open")
	}
	if !strings.Contains(viewText(model), "Line 2-") {
		t.Fatalf("expected scrolled modal footer to move to line 2, got:\n%s", viewText(model))
	}
	pressKey(t, model, "pgdown")
	if model.detail.scroll <= 1 {
		t.Fatalf("expected PgDn to page the modal body, got scroll %d", model.detail.scroll)
	}

	pressKey(t, model, "esc")
	if model.detail != nil {
		t.Fatal("expected Esc to close the detail modal")
	}
	if after := viewText(model); after != before {
		t.Fatalf("expected Esc to restore the exact prior review context\nbefore:\n%s\n\nafter:\n%s", before, after)
	}
}

func TestCockpitDetailTogglesWithDAndKeepsTimelineFollowing(t *testing.T) {
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
			Title:                   "minor: trim stale TODO",
			File:                    "internal/cache/cache.go",
			Line:                    42,
			Severity:                "minor",
			Author:                  "coderabbitai[bot]",
			Body:                    "Remove the stale TODO.",
			SourceRef:               "thread:PRRT_2,comment:PRRC_2",
			ReviewHash:              "def",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		}},
	})
	if err != nil {
		t.Fatalf("persist round: %v", err)
	}
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	source.addLine("line before modal\n")
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateActive,
		Issues:        []rounds.Issue{{Path: persisted.IssuePaths[0], Title: "minor: trim stale TODO", Severity: "minor", Status: "pending"}},
	})
	pressKey(t, model, "tab")
	before := viewText(model)

	pressKey(t, model, "d")
	if model.detail == nil {
		t.Fatal("expected D to open the detail modal")
	}
	source.addLine("line while modal is open\n")
	source.version++
	model.Update(cockpitTickMsg{})
	if !strings.Contains(strings.Join(model.viewport.VisibleLines(), "\n"), "line while modal is open") {
		t.Fatalf("expected Follow Mode to keep advancing under the modal, got %v", model.viewport.VisibleLines())
	}

	pressKey(t, model, "d")
	if model.detail != nil {
		t.Fatal("expected D to close the detail modal")
	}
	if rendered := viewText(model); !strings.Contains(rendered, "line while modal is open") {
		t.Fatalf("expected closing the modal to reveal the advanced timeline, got:\n%s", rendered)
	}

	pressKey(t, model, "d")
	pressKey(t, model, "d")
	withoutTimelineAdvance := viewText(model)
	if !strings.Contains(withoutTimelineAdvance, "line while modal is open") {
		t.Fatalf("expected detail toggle to preserve the current base context, got:\n%s", withoutTimelineAdvance)
	}
	if before == withoutTimelineAdvance {
		t.Fatal("expected the base context to differ after a real timeline advance")
	}
}

func TestCockpitSpecTaskDetailModalRestoresContextAndSurvivesStaleReload(t *testing.T) {
	gitRoot := t.TempDir()
	slug := "0005-tui-cockpit"
	file := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Open modal", spec.StatusInProgress)
	appendToFile(t, filepath.Join(gitRoot, file), "\n## Notes\n\nTask detail body line.\n")
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateResolvingWithAgent,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		Tasks:         []spec.Task{{ID: "task_01", File: file, Title: "Open modal", Status: spec.StatusPending}},
	})
	pressKey(t, model, "tab")
	before := viewText(model)

	pressKey(t, model, "enter")
	if model.detail == nil {
		t.Fatal("expected Enter to open the Task detail modal")
	}
	rendered := viewText(model)
	for _, expected := range []string{"SESSION.TIMELINE", "SPEC.TASK  task_01", "Open modal", "source: " + file + " (read-only)", "# Task 01: Open modal", "Line 1-"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected Task detail to contain %q, got:\n%s", expected, rendered)
		}
	}
	pressKey(t, model, "pgdown")
	if rendered := viewText(model); !strings.Contains(rendered, "Task detail body line.") {
		t.Fatalf("expected paged Task detail to show lower body text, got:\n%s", rendered)
	}

	if err := os.WriteFile(filepath.Join(gitRoot, file), []byte("---\ntask: task_01\nstatus: in_prog"), 0o644); err != nil {
		t.Fatalf("corrupt task file: %v", err)
	}
	source.version++
	model.Update(cockpitTickMsg{})
	rendered = viewText(model)
	for _, expected := range []string{"STALE: keeping last readable task file", "Task detail body line."} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected stale Task detail to contain %q, got:\n%s", expected, rendered)
		}
	}

	pressKey(t, model, "esc")
	if model.detail != nil {
		t.Fatal("expected Esc to close the Task detail modal")
	}
	if after := viewText(model); after != before {
		t.Fatalf("expected Esc to restore the exact prior spec context\nbefore:\n%s\n\nafter:\n%s", before, after)
	}
}

func TestCockpitDetailUsesFullSurfaceFallbackWhenTerminalIsTooShort(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateClean}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateClean,
		Issues:        []rounds.Issue{{Path: filepath.Join(t.TempDir(), "gone.md"), Title: "short terminal", Severity: "minor", Status: "resolved"}},
	})
	model.Update(tea.WindowSizeMsg{Width: 88, Height: 16})

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	rendered := viewText(model)
	if !strings.Contains(rendered, "REVIEW.ISSUE  #001") || !strings.Contains(rendered, "artifact not available") {
		t.Fatalf("expected full-surface detail fallback to render the detail, got:\n%s", rendered)
	}
	for _, background := range []string{"WORK QUEUE", "SESSION.TIMELINE"} {
		if strings.Contains(rendered, background) {
			t.Fatalf("expected too-short terminal fallback to hide background %q, got:\n%s", background, rendered)
		}
	}
}

func TestCockpitDetailKeepsAttachDetachKey(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateActive,
		Issues:        []rounds.Issue{{Path: filepath.Join(t.TempDir(), "gone.md"), Title: "attach detail", Severity: "minor", Status: "pending"}},
	})

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	if model.detail == nil {
		t.Fatal("expected detail modal open before detach")
	}
	cmd := pressKey(t, model, "q")
	if cmd == nil {
		t.Fatal("expected q to detach in attach mode while detail is open")
	}
	if msg := cmd(); fmt.Sprintf("%T", msg) != "tea.QuitMsg" {
		t.Fatalf("expected quit command while detail is open, got %T", msg)
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

// Suite: cockpit Attach parity
// Invariant: Attach replays the same read-only cockpit surfaces from the Run Event Journal for finished review and spec Runs.
// Boundary IN: synchronous cockpit model rendering, journal replay, Attach key routing, and Work Item detail modal rendering.
// Boundary OUT: CLI command parsing and store read-only opening, covered by internal/cli attach tests.
func TestCockpitAttachReplaysFinishedReviewRunThroughRedesignedCockpit(t *testing.T) {
	artifactDir := t.TempDir()
	persisted, err := rounds.PersistRound(context.Background(), rounds.PersistRequest{
		ArtifactDir:    artifactDir,
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/attach",
		HeadSHA:        "abc123",
		Round:          1,
		CreatedAt:      time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
		Items: []reviewsource.ReviewItem{
			{
				Title:                   "major: guard nil cache",
				File:                    "internal/cache/cache.go",
				Line:                    42,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "Guard the map lookup before dereferencing.",
				SourceRef:               "thread:PRRT_1,comment:PRRC_1",
				ReviewHash:              "h1",
				SourceReviewID:          "1",
				SourceReviewSubmittedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
			},
			{
				Title:                   "minor: trim stale TODO",
				File:                    "internal/cache/cache_test.go",
				Line:                    17,
				Severity:                "minor",
				Author:                  "coderabbitai[bot]",
				Body:                    "Remove the stale TODO from the test fixture.",
				SourceRef:               "thread:PRRT_2,comment:PRRC_2",
				ReviewHash:              "h2",
				SourceReviewID:          "1",
				SourceReviewSubmittedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
			},
		},
	})
	if err != nil {
		t.Fatalf("persist round: %v", err)
	}
	for _, path := range persisted.IssuePaths {
		if err := rounds.SetIssueStatus(path, rounds.StatusResolved, "", ""); err != nil {
			t.Fatalf("set issue status: %v", err)
		}
	}
	issues := []rounds.Issue{}
	for _, path := range persisted.IssuePaths {
		issues = append(issues, rounds.Issue{Path: path, Status: rounds.StatusPending})
	}

	startedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	source := &cockpitFakeSource{run: store.Run{ID: "run-review-attach", State: store.StateClean}, version: 1}
	source.addEvent(runevent.RunEvent{Batch: 1, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch, Summary: "Batch 001 executing.", Time: startedAt})
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentPlan,
		Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"plan","entries":[{"status":"pending","content":"Replay Attach cockpit"}]}}`),
		Time:    startedAt.Add(5 * time.Second),
	})
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentToolUpdated,
		Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"call_1","title":"read_file","status":"completed","rawInput":{"command":"rtk read internal/tui/cockpit.go"},"content":[{"content":{"type":"text","text":"loaded cockpit renderer"}}]}}`),
		Time:    startedAt.Add(20 * time.Second),
	})
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentThought,
		Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"checking Attach replay"}}}`),
		Time:    startedAt.Add(32 * time.Second),
	})
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Payload: []byte(`{"status":"completed"}`),
		Time:    startedAt.Add(38 * time.Second),
	})
	source.addEvent(runevent.RunEvent{
		Batch:   2,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonOutcome,
		Summary: "Clean after review Attach replay.",
		Time:    startedAt.Add(90 * time.Second),
	})
	model := newTestCockpit(t, source, LiveRunView{
		Command:       "attach",
		Repository:    "owner/project",
		PRNumber:      "123",
		HeadBranch:    "feature/attach",
		ReviewSource:  string(reviewsource.SourceCodeRabbit),
		Agent:         "codex",
		Model:         "gpt-5",
		RunID:         "run-review-attach",
		PipelineState: store.StateClean,
		Issues:        issues,
		BatchSizes:    []int{1, 1},
	})
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	for _, expected := range []string{
		"WORK QUEUE (2)",
		"SESSION.TIMELINE",
		"FETCH [done]",
		"TRIAGE [done]",
		"AGENT [done]",
		"VERIFY [done]",
		"PUSH [done]",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected review Attach cockpit to contain %q, got:\n%s", expected, rendered)
		}
	}
	assertContainsInOrder(t, rendered,
		"FETCH [done]",
		"TRIAGE [done]",
		"AGENT [done]",
		"VERIFY [done]",
		"PUSH [done]",
	)
	assertContainsInOrder(t, rendered,
		"▼ BATCH 001/002 executing 00:38",
		"PLAN",
		"[TOOL] read_file",
		"THINK checking Attach replay",
		"SESSION COMPLETED",
		"▼ BATCH 002/002 00:00",
		"OUTCOME",
		"Clean after review Attach replay.",
	)
	if strings.Contains(rendered, "pending  Replay Attach cockpit") {
		t.Fatalf("expected the plan payload to stay behind the PLAN summary row, got:\n%s", rendered)
	}
	for _, expected := range []string{
		"READ-ONLY",
		"2 issues total · 2 resolved · 0 unresolved",
		"Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter issue · D show detail · End follow · q detach",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected review Attach cockpit to contain %q, got:\n%s", expected, rendered)
		}
	}

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	rendered = viewText(model)
	for _, expected := range []string{
		"REVIEW.ISSUE  #001",
		"Guard the map lookup before dereferencing.",
		"source: thread:PRRT_1,comment:PRRC_1",
		"Keys: Esc close · j/k scroll · PgUp/PgDn page · q detach",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected review Attach modal to contain %q, got:\n%s", expected, rendered)
		}
	}
	pressKey(t, model, "esc")
	if strings.Contains(viewText(model), "REVIEW.ISSUE  #001") {
		t.Fatalf("expected Esc to close the review Attach modal, got:\n%s", viewText(model))
	}
	if source.run.State != store.StateClean {
		t.Fatalf("expected review Attach to leave Run state untouched, got %q", source.run.State)
	}
}

func TestCockpitAttachReplaysFinishedSpecRunThroughRedesignedCockpit(t *testing.T) {
	gitRoot := t.TempDir()
	slug := "0005-tui-cockpit"
	fileOne := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build modal detail", spec.StatusCompleted)
	fileTwo := writeCockpitTaskFile(t, gitRoot, slug, "task_02", "Sync skill docs", spec.StatusCompleted)

	startedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	source := &cockpitFakeSource{run: store.Run{ID: "run-spec-attach", State: store.StateClean}, version: 1}
	source.addEvent(runevent.RunEvent{
		Batch:       1,
		Source:      runevent.SourceDaemon,
		Kind:        runevent.KindDaemonTask,
		ReviewIssue: "task_01",
		Summary:     "Task task_01 started as Batch 001: Build modal detail",
		Time:        startedAt,
	})
	source.addEvent(runevent.RunEvent{
		Batch:       1,
		Source:      runevent.SourceDaemon,
		Kind:        runevent.KindDaemonVerification,
		ReviewIssue: "task_01",
		Summary:     "Verification command passed: rtk go test ./internal/tui/",
		Time:        startedAt.Add(12 * time.Second),
	})
	source.addEvent(runevent.RunEvent{
		Batch:       1,
		Source:      runevent.SourceDaemon,
		Kind:        runevent.KindDaemonTask,
		ReviewIssue: "task_01",
		Summary:     "Task task_01 settled completed.",
		Time:        startedAt.Add(20 * time.Second),
	})
	source.addEvent(runevent.RunEvent{
		Batch:   2,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonQA,
		Summary: "QA verdict pass for Spec 0005-tui-cockpit.",
		Time:    startedAt.Add(2 * time.Minute),
	})
	model := newTestCockpit(t, source, LiveRunView{
		Command:       "attach",
		HeadBranch:    "feature/cockpit",
		Agent:         "codex",
		Model:         "gpt-5",
		RunID:         "run-spec-attach",
		PipelineState: store.StateClean,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		Tasks: []spec.Task{
			{ID: "task_01", File: fileOne, Title: "Build modal detail", Status: spec.StatusPending},
			{ID: "task_02", File: fileTwo, Title: "Sync skill docs", Status: spec.StatusPending},
		},
	})
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	for _, expected := range []string{
		"WORK QUEUE (2)",
		"SESSION.TIMELINE",
		"AGENT [done]",
		"VERIFY [done]",
		"COMMIT [done]",
		"QA [done]",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected spec Attach cockpit to contain %q, got:\n%s", expected, rendered)
		}
	}
	assertContainsInOrder(t, rendered,
		"AGENT [done]",
		"VERIFY [done]",
		"COMMIT [done]",
		"QA [done]",
	)
	assertContainsInOrder(t, rendered,
		"▼ BATCH 001/002 00:20",
		"TASK",
		"Task task_01 started as Batch 001",
		"VERIFY",
		"Verification command passed: rtk go test",
		"TASK",
		"Task task_01 settled completed.",
		"▼ BATCH 002/002 00:00",
		"QA",
		"QA verdict pass for Spec 0005-tui-cockpit.",
	)
	for _, expected := range []string{
		"READ-ONLY",
		"2 Tasks total · 2 completed · 0 unresolved",
		"Keys: Tab focus · ↑↓ move/scroll · PgUp/PgDn page · Enter Task · D show detail · End follow · q detach",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected spec Attach cockpit to contain %q, got:\n%s", expected, rendered)
		}
	}

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	rendered = viewText(model)
	for _, expected := range []string{
		"SPEC.TASK  task_01",
		"Build modal detail",
		"source: docs/specs/0005-tui-cockpit/task_01.md (read-only)",
		"# Task 01: Build modal detail",
		"Keys: Esc close · j/k scroll · PgUp/PgDn page · q detach",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected spec Attach modal to contain %q, got:\n%s", expected, rendered)
		}
	}
	pressKey(t, model, "esc")
	if strings.Contains(viewText(model), "SPEC.TASK  task_01") {
		t.Fatalf("expected Esc to close the spec Attach modal, got:\n%s", viewText(model))
	}
	if source.run.State != store.StateClean {
		t.Fatalf("expected spec Attach to leave Run state untouched, got %q", source.run.State)
	}
}

func TestCockpitScrollFreezesFollowAndStatusBarNarratesStates(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	for index := 1; index <= 60; index++ {
		source.addLine(fmt.Sprintf("line %04d\n", index))
	}
	model := newTestCockpit(t, source, LiveRunView{PipelineState: store.StateActive})

	if strings.Contains(viewText(model), "SCROLLED") {
		t.Fatalf("expected no scroll hint while following, got:\n%s", viewText(model))
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
	if strings.Contains(rendered, "SCROLLED") || !strings.Contains(rendered, "line 0061") {
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
	if !strings.Contains(rendered, "READ-ONLY") {
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

func TestCockpitSidebarShowsBatchesStatusAndElapsed(t *testing.T) {
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
		Items: []reviewsource.ReviewItem{
			{Title: "one", File: "a.go", Line: 1, Severity: "major", Author: "bot", Body: "b", SourceRef: "thread:1,comment:1", ReviewHash: "h1", SourceReviewID: "1", SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)},
			{Title: "two", File: "b.go", Line: 2, Severity: "minor", Author: "bot", Body: "b", SourceRef: "thread:2,comment:2", ReviewHash: "h2", SourceReviewID: "1", SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)},
			{Title: "three", File: "c.go", Line: 3, Severity: "minor", Author: "bot", Body: "b", SourceRef: "thread:3,comment:3", ReviewHash: "h3", SourceReviewID: "1", SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)},
		},
	})
	if err != nil {
		t.Fatalf("persist round: %v", err)
	}
	// First issue resolved on disk; the other two stay pending.
	if err := rounds.SetIssueStatus(persisted.IssuePaths[0], rounds.StatusResolved, "", ""); err != nil {
		t.Fatalf("set status: %v", err)
	}
	issues := []rounds.Issue{}
	for _, path := range persisted.IssuePaths {
		issues = append(issues, rounds.Issue{Path: path, Status: rounds.StatusPending})
	}
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	clock := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	source.addEvent(runevent.RunEvent{
		Batch:   1,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonBatch,
		Summary: "Batch 001 executing.",
		Time:    clock,
	})
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode:   CockpitOwning,
		View:   LiveRunView{PipelineState: store.StateResolvingWithAgent, Width: 100, Issues: issues, BatchSizes: []int{2, 1}},
		RunID:  "run-1",
		Source: source,
		Now:    func() time.Time { return clock },
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	clock = clock.Add(83 * time.Second)

	rendered := viewText(model)
	for _, expected := range []string{
		"WORK QUEUE (3)",
		"FETCH [done]",
		"TRIAGE [done]",
		"AGENT [run]",
		"VERIFY [wait]",
		"PUSH [locked]",
		"BATCH 001/002",
		"BATCH 002/002",
		"[done] MAJOR",
		"#1",
		"one",
		"a.go:1",
		"[run] MINOR",
		"#2",
		"two",
		"b.go:2",
		"01:23",
		"[wait] MINOR",
		"#3",
		"three",
		"c.go:3",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected sidebar to contain %q, got:\n%s", expected, rendered)
		}
	}
	// The clock belongs to the executing Batch separator, never to issues.
	for _, line := range strings.Split(rendered, "\n") {
		if !strings.Contains(line, "01:23") {
			continue
		}
		if !strings.Contains(line, "BATCH 001/002") {
			t.Fatalf("expected elapsed clock on the executing Batch separator, got %q", line)
		}
	}
}

// writeCockpitTaskFile writes a parseable task file under the git root's
// docs/specs/<slug>/ and returns its path relative to the git root, the way
// spec.Load records Task files.
func writeCockpitTaskFile(t *testing.T, gitRoot string, slug string, id string, title string, status spec.Status) string {
	t.Helper()
	dir := filepath.Join(gitRoot, "docs", "specs", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir spec dir: %v", err)
	}
	content := fmt.Sprintf("---\ntask: %s\nspec: %s\nstatus: %s\ntype: backend\n---\n\n# Task 01: %s\n\n## Verification\n\n- `true` — expected: passes.\n", id, slug, status, title)
	relative := filepath.Join("docs", "specs", slug, id+".md")
	if err := os.WriteFile(filepath.Join(gitRoot, relative), []byte(content), 0o644); err != nil {
		t.Fatalf("write task file: %v", err)
	}
	return relative
}

func appendToFile(t *testing.T, path string, content string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open file for append: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("close appended file: %v", err)
		}
	}()
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("append file content: %v", err)
	}
}

func TestCockpitSpecRunShowsTasksInGraphOrderAndRefreshesStatuses(t *testing.T) {
	gitRoot := t.TempDir()
	slug := "0001-widget-flow"
	fileOne := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build core", spec.StatusPending)
	fileTwo := writeCockpitTaskFile(t, gitRoot, slug, "task_02", "Wire API", spec.StatusPending)
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateResolvingWithAgent,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		Tasks: []spec.Task{
			{ID: "task_01", File: fileOne, Title: "Build core", Status: spec.StatusPending},
			{ID: "task_02", File: fileTwo, Title: "Wire API", Status: spec.StatusPending},
		},
	})

	rendered := viewText(model)
	for _, expected := range []string{
		"WORK QUEUE (2)",
		"AGENT [run]",
		"VERIFY [wait]",
		"COMMIT [wait]",
		"task_01",
		"[run]",
		"task_02",
		"[wait]",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected the Task pane to contain %q, got:\n%s", expected, rendered)
		}
	}
	if strings.Index(rendered, "task_01") > strings.Index(rendered, "task_02") {
		t.Fatalf("expected Tasks in Task Graph order, got:\n%s", rendered)
	}

	// The Agent settles task_01 in its file and the Daemon journals the
	// settlement; the poll tick re-reads the task files (never the sink).
	if err := spec.SetStatus(filepath.Join(gitRoot, fileOne), spec.StatusCompleted); err != nil {
		t.Fatalf("settle task file: %v", err)
	}
	source.addDaemonEvent(runevent.KindDaemonTask, "Task task_01 settled completed.")
	source.version++
	model.Update(cockpitTickMsg{})

	rendered = viewText(model)
	for _, expected := range []string{
		"[done]",
		"WORK QUEUE (2)",
		"Task task_01 settled completed.",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected the refreshed view to contain %q, got:\n%s", expected, rendered)
		}
	}
}

func TestCockpitSpecRunDerivesConcurrentTaskStateFromJournal(t *testing.T) {
	gitRoot := t.TempDir()
	slug := "0009-parallel-scheduling"
	files := []string{
		writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build scheduler", spec.StatusPending),
		writeCockpitTaskFile(t, gitRoot, slug, "task_02", "Wire queue", spec.StatusPending),
		writeCockpitTaskFile(t, gitRoot, slug, "task_03", "Write docs", spec.StatusPending),
	}
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	source.addTaskEvent(t, "task_01", "started", "", 1)
	source.addTaskEvent(t, "task_02", "started", "", 2)
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateResolvingWithAgent,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		Concurrency:   2,
		Tasks: []spec.Task{
			{ID: "task_01", File: files[0], Title: "Build scheduler", Status: spec.StatusPending},
			{ID: "task_02", File: files[1], Title: "Wire queue", Status: spec.StatusPending},
			{ID: "task_03", File: files[2], Title: "Write docs", Status: spec.StatusPending},
		},
	})
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	assertTaskQueueRow(t, rendered, "task_01", "[run]")
	assertTaskQueueRow(t, rendered, "task_02", "[run]")
	assertTaskQueueRow(t, rendered, "task_03", "[wait]")
	assertContainsInOrder(t, rendered, "task_01 Build scheduler", "task_02 Wire queue", "task_03 Write docs")
	if !strings.Contains(rendered, "3 Tasks total · 0 completed · 3 unresolved") {
		t.Fatalf("expected initial concurrent totals, got:\n%s", rendered)
	}

	source.addTaskEvent(t, "task_02", "settled", spec.StatusCompleted, 2)
	source.version++
	model.Update(cockpitTickMsg{})

	rendered = viewText(model)
	assertTaskQueueRow(t, rendered, "task_01", "[run]")
	assertTaskQueueRow(t, rendered, "task_02", "[done]")
	assertTaskQueueRow(t, rendered, "task_03", "[wait]")
	if !strings.Contains(rendered, "3 Tasks total · 1 completed · 2 unresolved") {
		t.Fatalf("expected out-of-order settlement totals, got:\n%s", rendered)
	}

	source.addTaskEvent(t, "task_01", "settled", spec.StatusCompleted, 1)
	source.version++
	model.Update(cockpitTickMsg{})

	rendered = viewText(model)
	assertTaskQueueRow(t, rendered, "task_01", "[done]")
	assertTaskQueueRow(t, rendered, "task_02", "[done]")
	assertTaskQueueRow(t, rendered, "task_03", "[wait]")
	if !strings.Contains(rendered, "3 Tasks total · 2 completed · 1 unresolved") {
		t.Fatalf("expected final concurrent totals, got:\n%s", rendered)
	}
}

func TestCockpitSpecRunInterleavedTaskReplayMatchesLivePolling(t *testing.T) {
	gitRoot := t.TempDir()
	slug := "0009-parallel-scheduling"
	tasks := []spec.Task{
		{ID: "task_01", File: writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build scheduler", spec.StatusPending), Title: "Build scheduler", Status: spec.StatusPending},
		{ID: "task_02", File: writeCockpitTaskFile(t, gitRoot, slug, "task_02", "Wire queue", spec.StatusPending), Title: "Wire queue", Status: spec.StatusPending},
	}
	view := LiveRunView{
		PipelineState: store.StateResolvingWithAgent,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		Concurrency:   2,
		Tasks:         tasks,
	}

	liveSource := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	live := newTestCockpit(t, liveSource, view)
	live.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	liveSource.addTaskEvent(t, "task_01", "started", "", 1)
	liveSource.addTaskEvent(t, "task_02", "started", "", 2)
	liveSource.version++
	live.Update(cockpitTickMsg{})
	liveSource.addTaskEvent(t, "task_02", "settled", spec.StatusCompleted, 2)
	liveSource.version++
	live.Update(cockpitTickMsg{})
	liveRendered := viewText(live)

	replaySource := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	replaySource.addTaskEvent(t, "task_01", "started", "", 1)
	replaySource.addTaskEvent(t, "task_02", "started", "", 2)
	replaySource.addTaskEvent(t, "task_02", "settled", spec.StatusCompleted, 2)
	replay := newTestCockpit(t, replaySource, view)
	replay.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if replayRendered := viewText(replay); replayRendered != liveRendered {
		t.Fatalf("expected Attach-style replay to match live polling\nlive:\n%s\n\nreplay:\n%s", liveRendered, replayRendered)
	}
}

func TestCockpitSpecRunHeaderShowsConcurrency(t *testing.T) {
	model := newSpecPhaseCockpit(t, false, store.StateResolvingWithAgent, spec.StatusPending)
	model.cfg.View.Concurrency = 3
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	if !strings.Contains(rendered, "Concurrency: 3") {
		t.Fatalf("expected spec Run header to show concurrency, got:\n%s", rendered)
	}
}

func TestCockpitSpecRunReadsTaskStatusAndDetailFromWorkDir(t *testing.T) {
	gitRoot := t.TempDir()
	workDir := t.TempDir()
	slug := "0001-widget-flow"
	file := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build core", spec.StatusPending)
	workFile := writeCockpitTaskFile(t, workDir, slug, "task_01", "Build core", spec.StatusCompleted)
	appendToFile(t, filepath.Join(gitRoot, file), "\n## Notes\n\nuser-root stale detail\n")
	appendToFile(t, filepath.Join(workDir, workFile), "\n## Notes\n\nworktree truth detail\n")
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateResolvingWithAgent,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		WorkDir:       workDir,
		Tasks:         []spec.Task{{ID: "task_01", File: file, Title: "Build core", Status: spec.StatusPending}},
	})

	rendered := viewText(model)
	if !strings.Contains(rendered, "[done]") {
		t.Fatalf("expected WorkDir task status to render completed, got:\n%s", rendered)
	}
	if strings.Contains(rendered, "[run]") {
		t.Fatalf("expected stale user-root pending status ignored, got:\n%s", rendered)
	}

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	pressKey(t, model, "pgdown")
	rendered = viewText(model)
	if !strings.Contains(rendered, "worktree truth detail") {
		t.Fatalf("expected Task detail to read WorkDir body, got:\n%s", rendered)
	}
	if strings.Contains(rendered, "user-root stale detail") {
		t.Fatalf("expected Task detail not to read stale GitRoot body, got:\n%s", rendered)
	}
}

func TestCockpitSpecRunFallsBackToGitRootWhenWorkDirIsGone(t *testing.T) {
	gitRoot := t.TempDir()
	workDir := filepath.Join(t.TempDir(), "pruned-worktree")
	slug := "0001-widget-flow"
	file := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build core", spec.StatusCompleted)
	appendToFile(t, filepath.Join(gitRoot, file), "\n## Notes\n\nintegrated user-root detail\n")
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateClean}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateClean,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		WorkDir:       workDir,
		Tasks:         []spec.Task{{ID: "task_01", File: file, Title: "Build core", Status: spec.StatusPending}},
	})

	rendered := viewText(model)
	if !strings.Contains(rendered, "[done]") {
		t.Fatalf("expected missing WorkDir to fall back to GitRoot status, got:\n%s", rendered)
	}

	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	pressKey(t, model, "pgdown")
	if rendered = viewText(model); !strings.Contains(rendered, "integrated user-root detail") {
		t.Fatalf("expected missing WorkDir detail fallback to GitRoot, got:\n%s", rendered)
	}
}

func TestCockpitSpecRunKeepsLastGoodStatusOnMidWriteTaskFile(t *testing.T) {
	gitRoot := t.TempDir()
	slug := "0001-widget-flow"
	file := writeCockpitTaskFile(t, gitRoot, slug, "task_01", "Build core", spec.StatusCompleted)
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateResolvingWithAgent}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		PipelineState: store.StateResolvingWithAgent,
		RunKind:       store.KindImplement,
		SpecSlug:      slug,
		GitRoot:       gitRoot,
		// The seeded status is stale on purpose: the initial refresh must
		// read the task file, not trust the loaded snapshot.
		Tasks: []spec.Task{{ID: "task_01", File: file, Title: "Build core", Status: spec.StatusPending}},
	})

	if !strings.Contains(viewText(model), "[done]") {
		t.Fatalf("expected the initial refresh to read the task file, got:\n%s", viewText(model))
	}

	// A mid-write read sees a truncated file: the pane keeps the last good
	// status and no error surfaces in the view.
	if err := os.WriteFile(filepath.Join(gitRoot, file), []byte("---\ntask: task_01\nstatus: comp"), 0o644); err != nil {
		t.Fatalf("truncate task file: %v", err)
	}
	source.version++
	model.Update(cockpitTickMsg{})

	rendered := viewText(model)
	if !strings.Contains(rendered, "[done]") {
		t.Fatalf("expected the last good status kept through the mid-write read, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "COMMIT [done]") {
		t.Fatalf("expected the phase row to keep the last good completed state, got:\n%s", rendered)
	}
}

func TestOwningCockpitClosesOnQOnlyAfterTerminalState(t *testing.T) {
	source := &cockpitFakeSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode:   CockpitOwning,
		View:   LiveRunView{PipelineState: store.StateActive, Width: 100},
		RunID:  "run-1",
		Source: source,
		OnStop: func() {},
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}

	if cmd := pressKey(t, model, "q"); cmd != nil {
		t.Fatal("expected q to do nothing while the Run is active")
	}

	source.run.State = store.StateClean
	source.version++
	model.Update(cockpitTickMsg{})

	if !strings.Contains(viewText(model), "q close") {
		t.Fatalf("expected lingering footer after the Run ended, got:\n%s", viewText(model))
	}
	cmd := pressKey(t, model, "q")
	if cmd == nil {
		t.Fatal("expected q to close the lingering cockpit")
	}
	if msg := cmd(); fmt.Sprintf("%T", msg) != "tea.QuitMsg" {
		t.Fatalf("expected quit, got %T", msg)
	}
	if cmd := pressKey(t, model, "ctrl+c"); cmd == nil {
		t.Fatal("expected ctrl+c to also close after terminal state")
	}
}
