package tui

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
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
