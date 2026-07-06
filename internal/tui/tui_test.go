package tui

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestRenderInteractiveInputShowsCurrentAndConfiguredDefaults(t *testing.T) {
	view := RenderInteractiveInput(InputRequest{
		Command: "resolve",
		Values: CommandValues{
			Agent: "codex",
			Round: "all",
		},
		PRSuggestion: Suggestion{
			Value:  "123",
			Source: "current",
		},
		AgentSuggestion: Suggestion{
			Value:  "claude",
			Source: "remembered",
		},
	})

	for _, expected := range []string{
		"Roundfix Interactive Input",
		"Command: resolve",
		"Suggested Open Pull Request: #123 (current)",
		"Suggested Agent: codex (config)",
		"Press Enter to accept a suggestion.",
	} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected view to contain %q, got:\n%s", expected, view)
		}
	}
}

func TestCollectInputAppliesDefaultsAndUserOverrides(t *testing.T) {
	input := strings.NewReader("\nclaude\n2\n\nsonnet\n")
	var output strings.Builder

	values, err := CollectInput(context.Background(), InputRequest{
		Command: "resolve",
		Values: CommandValues{
			Agent:       "codex",
			Round:       "all",
			ArtifactDir: ".roundfix",
		},
		PRSuggestion: Suggestion{Value: "123", Source: "remembered"},
	}, input, &output)
	if err != nil {
		t.Fatalf("collect input: %v", err)
	}

	if values.PRNumber != "123" {
		t.Fatalf("expected default PR 123, got %q", values.PRNumber)
	}
	if values.Agent != "claude" {
		t.Fatalf("expected agent override, got %q", values.Agent)
	}
	if values.Round != "2" {
		t.Fatalf("expected round override, got %q", values.Round)
	}
	if values.ArtifactDir != ".roundfix" {
		t.Fatalf("expected artifact default, got %q", values.ArtifactDir)
	}
	if values.Model != "sonnet" {
		t.Fatalf("expected model override, got %q", values.Model)
	}
	if !strings.Contains(output.String(), "Open Pull Request [123]:") {
		t.Fatalf("expected prompted PR default, got %q", output.String())
	}
}

func TestCollectInputSpecPickerSelectsListedSpec(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantSpec string
	}{
		{name: "by number", input: "2\nclaude\n", wantSpec: "0002-other-flow"},
		{name: "by slug", input: "0001-widget-flow\nclaude\n", wantSpec: "0001-widget-flow"},
		{name: "out-of-range number passes through", input: "9\nclaude\n", wantSpec: "9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder

			values, err := CollectInput(context.Background(), InputRequest{
				Command:     "implement",
				Values:      CommandValues{Agent: "codex"},
				SpecOptions: []string{"0001-widget-flow", "0002-other-flow"},
			}, strings.NewReader(tt.input), &output)
			if err != nil {
				t.Fatalf("collect input: %v", err)
			}

			if values.Spec != tt.wantSpec {
				t.Fatalf("expected Spec %q, got %q", tt.wantSpec, values.Spec)
			}
			if values.Agent != "claude" {
				t.Fatalf("expected agent override, got %q", values.Agent)
			}
			for _, expected := range []string{
				"Active Specs:",
				"1. 0001-widget-flow",
				"2. 0002-other-flow",
				"Pick a Spec by number or slug.",
				"Spec []:",
				"Agent [codex]:",
			} {
				if !strings.Contains(output.String(), expected) {
					t.Fatalf("expected the Spec picker to show %q, got:\n%s", expected, output.String())
				}
			}
		})
	}
}

func TestCollectInputImplementQAGate(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultQA  bool
		wantQA     bool
		wantPrompt string
	}{
		{name: "yes enables QA", input: "\n\ny\n", wantQA: true, wantPrompt: "QA gate [y/N]:"},
		{name: "empty keeps QA disabled", input: "\n\n\n", wantQA: false, wantPrompt: "QA gate [y/N]:"},
		{name: "empty keeps QA flag default", input: "\n\n\n", defaultQA: true, wantQA: true, wantPrompt: "QA gate [Y/n]:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder

			values, err := CollectInput(context.Background(), InputRequest{
				Command: "implement",
				Values: CommandValues{
					Spec:  "0001-widget-flow",
					Agent: "codex",
					QA:    tt.defaultQA,
				},
			}, strings.NewReader(tt.input), &output)
			if err != nil {
				t.Fatalf("collect input: %v", err)
			}
			if values.QA != tt.wantQA {
				t.Fatalf("expected QA %v, got %v", tt.wantQA, values.QA)
			}
			if !strings.Contains(output.String(), tt.wantPrompt) {
				t.Fatalf("expected QA prompt %q, got:\n%s", tt.wantPrompt, output.String())
			}
		})
	}
}

func TestCollectInputImplementQAGateInvalidInputRepromptsOnce(t *testing.T) {
	var output strings.Builder

	_, err := CollectInput(context.Background(), InputRequest{
		Command: "implement",
		Values:  CommandValues{Spec: "0001-widget-flow", Agent: "codex"},
	}, strings.NewReader("\n\nmaybe\nlater\n"), &output)

	if err == nil {
		t.Fatal("expected invalid QA input to fail")
	}
	if !strings.Contains(err.Error(), "QA gate") {
		t.Fatalf("expected error to name QA gate, got %v", err)
	}
	if count := strings.Count(output.String(), "QA gate [y/N]:"); count != 2 {
		t.Fatalf("expected one QA re-prompt, got %d prompts:\n%s", count, output.String())
	}
}

func TestRenderLiveRunViewGroupsIssuesAndShowsStatusStrips(t *testing.T) {
	view := RenderLiveRunView(LiveRunView{
		Command:       "resolve",
		Repository:    "owner/project",
		PRNumber:      "123",
		HeadBranch:    "feature/review",
		ReviewSource:  "CodeRabbit",
		Agent:         "Codex",
		HEAD:          "abc123",
		RunID:         "run_123",
		PipelineState: "ResolvingWithAgent",
		BudgetState:   "38m / 2h",
		GitState:      "clean, 1 unpushed commit",
		CurrentRound:  2,
		MaxRounds:     6,
		AutoCommit:    true,
		AutoPush:      true,
		LastPush:      "pending",
		Width:         100,
		Issues: []rounds.Issue{
			{Round: 2, Title: "fix stale readme", Severity: "minor", Status: rounds.StatusPending, File: "README.md", Line: 12},
			{Round: 1, Title: "guard auth cache", Severity: "major", Status: rounds.StatusResolved, File: "api/auth.go", Line: 88},
			{Round: 2, Title: "invalidate cache", Severity: "major", Status: rounds.StatusValid, File: "src/cache.ts", Line: 41},
		},
		Console: []string{
			"codex resolving batch 1/2",
			"running make verify",
		},
	})

	expected := []string{
		"Roundfix resolve",
		"Target:",
		"PR: #123 owner/project",
		"Branch: feature/review",
		"Source: CodeRabbit",
		"Agent: Codex",
		"HEAD: abc123",
		"Run:",
		"ID: run_123",
		"State: ResolvingWithAgent",
		"Round: 2 / 6",
		"Budget: 38m / 2h",
		"Git: clean, 1 unpushed commit",
		"Auto-commit: on",
		"Auto-push: on",
		"Last push: pending",
		"Review Issues",
		"Agent Console",
		"codex resolving batch 1/2",
		"running make verify",
		"Round 001",
		"major    resolved   api/auth.go:88",
		"guard auth cache",
		"Round 002",
		"major    valid      src/cache.ts:41",
		"invalidate cache",
		"minor    pending    README.md:12",
		"fix stale readme",
		"Keys: Ctrl-C stop",
	}
	for _, text := range expected {
		if !strings.Contains(view, text) {
			t.Fatalf("expected live view to contain %q, got:\n%s", text, view)
		}
	}
	for _, removed := range []string{"[tab] focus", "[s] stop"} {
		if strings.Contains(view, removed) {
			t.Fatalf("did not expect non-interactive hint %q, got:\n%s", removed, view)
		}
	}
	if strings.Contains(view, "Concurrency:") {
		t.Fatalf("expected review Run header not to include concurrency, got:\n%s", view)
	}
}

func TestRenderLiveRunViewSpecRunRendersTasksAsWorkItems(t *testing.T) {
	view := RenderLiveRunView(LiveRunView{
		Command:       "implement",
		RunKind:       store.KindImplement,
		SpecSlug:      "0001-widget-flow",
		GitRoot:       "/repo",
		WorkDir:       "/home/user/.roundfix/worktrees/repo/run_9",
		HeadBranch:    "ma/widget-flow",
		Agent:         "Codex",
		HEAD:          "abc123",
		RunID:         "run_9",
		PipelineState: "ResolvingWithAgent",
		Concurrency:   2,
		BudgetState:   "38m / 2h",
		GitState:      "clean, 1 unpushed commit",
		CurrentRound:  2,
		MaxRounds:     6,
		AutoCommit:    true,
		AutoPush:      false,
		LastPush:      "disabled",
		Width:         100,
		Tasks: []spec.Task{
			{ID: "task_01", Title: "Build core", Status: spec.StatusCompleted},
			{ID: "task_02", Title: "Wire API", Status: spec.StatusInProgress},
			{ID: "task_03", Title: "Write docs", Status: spec.StatusPending},
		},
		Console: []string{"Task task_01 settled completed."},
	})

	for _, expected := range []string{
		"Roundfix implement",
		"Spec: 0001-widget-flow",
		"Branch: ma/widget-flow",
		"Agent: Codex",
		"Run:",
		"ID: run_9",
		"State: ResolvingWithAgent",
		"Concurrency: 2",
		"Run Worktree: /home/user/.roundfix/worktrees/repo/run_9",
		"Git: clean, 1 unpushed commit",
		"Auto-commit: on",
		"Auto-push: off",
		"Last push: disabled",
		"Tasks",
		"Agent Console",
		"task_01 completed — Build core",
		"task_02 in_progress — Wire API",
		"task_03 pending — Write docs",
		"Task task_01 settled completed.",
	} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected spec live view to contain %q, got:\n%s", expected, view)
		}
	}
	for _, absent := range []string{"Review Issues", "PR: #", "Source:"} {
		if strings.Contains(view, absent) {
			t.Fatalf("expected review vocabulary %q absent from a spec Run view, got:\n%s", absent, view)
		}
	}
	for _, absent := range []string{"\n  Round:", "\n  Budget:"} {
		if strings.Contains(view, absent) {
			t.Fatalf("expected spec Run header line %q absent, got:\n%s", strings.TrimSpace(absent), view)
		}
	}
	if strings.Index(view, "task_01") > strings.Index(view, "task_02") || strings.Index(view, "task_02") > strings.Index(view, "task_03") {
		t.Fatalf("expected Tasks rendered in Task Graph order, got:\n%s", view)
	}
}

func TestRunTimelineRendersTaskAndQAEventSummaries(t *testing.T) {
	timeline := NewRunTimeline(10)
	timeline.Append(runevent.RunEvent{Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, Summary: "Task task_01 started as Batch 001: Build core"})
	timeline.Append(runevent.RunEvent{Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, Summary: "Task task_01 settled completed."})
	timeline.Append(runevent.RunEvent{Source: runevent.SourceDaemon, Kind: runevent.KindDaemonQA, Summary: "QA verdict pass for Spec 0001-widget-flow."})

	lines := timeline.Lines()
	expected := []string{
		"Task task_01 started as Batch 001: Build core",
		"Task task_01 settled completed.",
		"QA verdict pass for Spec 0001-widget-flow.",
	}
	if len(lines) != len(expected) {
		t.Fatalf("expected one timeline line per event, got %v", lines)
	}
	for index := range expected {
		if lines[index] != expected[index] {
			t.Fatalf("expected timeline line %q, got %q", expected[index], lines[index])
		}
	}
}

func TestStreamBufferKeepsRecentConsoleOutput(t *testing.T) {
	buffer := &StreamBuffer{MaxLines: 2}
	if _, err := buffer.Write([]byte("first\nsecond\nthi")); err != nil {
		t.Fatalf("write stream: %v", err)
	}
	if _, err := buffer.Write([]byte("rd\n")); err != nil {
		t.Fatalf("write stream: %v", err)
	}

	lines := buffer.Lines()
	if got := strings.Join(lines, "|"); got != "second|third" {
		t.Fatalf("expected bounded stream lines, got %q", got)
	}
}

func TestRenderAgentSidebarShowsBatchProgressAndTotalIssues(t *testing.T) {
	view := LiveRunView{
		BatchNumber: 1,
		BatchTotal:  3,
		TotalIssues: 3,
		Issues: []rounds.Issue{
			{Path: "/repo/.roundfix/reviews/pr-20/round-001/issue_001.md", Round: 1, Title: "first", Severity: "minor", Status: rounds.StatusPending, File: "apps/api/test.ts", Line: 7},
			{Path: "/repo/.roundfix/reviews/pr-20/round-001/issue_002.md", Round: 1, Title: "second", Severity: "major", Status: rounds.StatusPending, File: "docker-compose.yml", Line: 22},
			{Path: "/repo/.roundfix/reviews/pr-20/round-001/issue_003.md", Round: 1, Title: "third", Severity: "major", Status: rounds.StatusFailed, File: "Makefile", Line: 52},
		},
	}

	sidebar := stripANSI(renderAgentSidebar(view, time.Now().Add(-90*time.Second), 42, 14))
	for _, expected := range []string{
		"batch_001/003",
		"FILES 3 · ISSUES 3",
		"Issue 001 • minor",
		"RUNNING •",
		"Issue 002 • major",
		"PENDING • --",
		"Issue 003 • major",
		"FAILED • --",
	} {
		if !strings.Contains(sidebar, expected) {
			t.Fatalf("expected sidebar to contain %q, got:\n%s", expected, sidebar)
		}
	}
	for _, hidden := range []string{"first", "apps/api/test.ts", "Makefile:52"} {
		if strings.Contains(sidebar, hidden) {
			t.Fatalf("expected sidebar to hide %q, got:\n%s", hidden, sidebar)
		}
	}
}

func rawAgentEvent(text string) runevent.RunEvent {
	return runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Summary: text,
		Payload: []byte(`{"text":` + strconv.Quote(text) + `}`),
	}
}

func TestRunTimelineCoalescesMessageChunksIntoOneLine(t *testing.T) {
	timeline := NewRunTimeline(10)
	for _, chunk := range []string{"Hel", "lo ", "world\n"} {
		timeline.Append(runevent.RunEvent{
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentMessage,
			Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":` + strconv.Quote(chunk) + `}}}`),
		})
	}

	lines := timeline.Lines()
	if len(lines) != 1 || lines[0] != "Hello world" {
		t.Fatalf("expected chunks coalesced into one line, got %v", lines)
	}
}

func TestRunTimelineBoundsConsoleMemory(t *testing.T) {
	timeline := NewRunTimeline(5)
	for index := 0; index < 50; index++ {
		timeline.Append(rawAgentEvent("line " + strconv.Itoa(index) + "\n"))
	}

	lines := timeline.Lines()
	if len(lines) != 5 {
		t.Fatalf("expected ring-bounded console of 5 lines, got %d", len(lines))
	}
	if lines[4] != "line 49" {
		t.Fatalf("expected newest lines kept, got %v", lines)
	}
}

func TestRunTimelineSkipsUnknownEventKinds(t *testing.T) {
	timeline := NewRunTimeline(5)
	timeline.Append(runevent.RunEvent{Source: runevent.SourceAgent, Kind: "future.unknown", Payload: []byte(`{}`)})
	timeline.Append(rawAgentEvent("kept\n"))

	lines := timeline.Lines()
	if len(lines) != 1 || lines[0] != "kept" {
		t.Fatalf("expected unknown kinds skipped, got %v", lines)
	}
}

func TestRunTimelineRendersToolEventsFromRawPayloads(t *testing.T) {
	timeline := NewRunTimeline(20)
	timeline.Append(runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentToolStarted,
		Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"tool_call","toolCallId":"call_1","title":"rtk go test","status":"pending"}}`),
	})
	timeline.Append(runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Payload: []byte(`{"status":"completed"}`),
	})

	lines := strings.Join(timeline.Lines(), "\n")
	if !strings.Contains(lines, "[TOOL] rtk go test (call_1)") {
		t.Fatalf("expected tool marker rendered from raw payload, got %q", lines)
	}
	if !strings.Contains(lines, "SESSION COMPLETED") {
		t.Fatalf("expected session status rendered, got %q", lines)
	}
}
