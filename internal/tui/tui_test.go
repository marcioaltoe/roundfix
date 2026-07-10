package tui

import (
	"context"
	"io"
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
	input := strings.NewReader("\nclaude\nsonnet\nmaximum\n2\n\n")
	var output strings.Builder

	values, err := CollectInput(context.Background(), InputRequest{
		Command: "resolve",
		Values: CommandValues{
			Agent:       "codex",
			Round:       "all",
			ArtifactDir: ".roundfix",
		},
		PRSuggestion:      Suggestion{Value: "123", Source: "remembered"},
		SelectionDefaults: testSelectionDefaults(),
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
	if values.ReasoningEffort != "maximum" {
		t.Fatalf("expected reasoning override, got %q", values.ReasoningEffort)
	}
	if !strings.Contains(output.String(), "Open Pull Request [123]:") {
		t.Fatalf("expected prompted PR default, got %q", output.String())
	}
}

func TestCollectInputRecomputesSelectionDefaultsWhenAgentChanges(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantAgent        string
		wantModel        string
		wantReasoning    string
		wantModelPrompt  string
		wantEffortPrompt string
	}{
		{
			name:             "same agent keeps seeded selection",
			input:            "\n\n\n\n\n\n",
			wantAgent:        "codex",
			wantModel:        "configured-codex",
			wantReasoning:    "configured-xhigh",
			wantModelPrompt:  "Agent Model [configured-codex]:",
			wantEffortPrompt: "Default Reasoning Effort [configured-xhigh]:",
		},
		{
			name:             "changed agent uses selected runtime defaults",
			input:            "\nclaude\n\n\n\n\n",
			wantAgent:        "claude",
			wantModel:        "opus",
			wantReasoning:    "high",
			wantModelPrompt:  "Agent Model [opus]:",
			wantEffortPrompt: "Default Reasoning Effort [high]:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder

			values, err := CollectInput(context.Background(), InputRequest{
				Command: "resolve",
				Values: CommandValues{
					PRNumber:        "123",
					Agent:           "codex",
					Model:           "configured-codex",
					ReasoningEffort: "configured-xhigh",
					Round:           "all",
					ArtifactDir:     ".roundfix",
				},
				SelectionDefaults: testSelectionDefaults(),
			}, strings.NewReader(tt.input), &output)
			if err != nil {
				t.Fatalf("collect input: %v", err)
			}

			if values.Agent != tt.wantAgent {
				t.Fatalf("expected agent %q, got %q", tt.wantAgent, values.Agent)
			}
			if values.Model != tt.wantModel {
				t.Fatalf("expected model %q, got %q", tt.wantModel, values.Model)
			}
			if values.ReasoningEffort != tt.wantReasoning {
				t.Fatalf("expected reasoning effort %q, got %q", tt.wantReasoning, values.ReasoningEffort)
			}
			if !strings.Contains(output.String(), tt.wantModelPrompt) {
				t.Fatalf("expected model prompt %q, got:\n%s", tt.wantModelPrompt, output.String())
			}
			if !strings.Contains(output.String(), tt.wantEffortPrompt) {
				t.Fatalf("expected reasoning prompt %q, got:\n%s", tt.wantEffortPrompt, output.String())
			}
		})
	}
}

func TestCollectInputSpecPickerSelectsListedSpec(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantSpec string
	}{
		{name: "by number", input: "2\nclaude\n\n\n", wantSpec: "0002-other-flow"},
		{name: "by slug", input: "0001-widget-flow\nclaude\n\n\n", wantSpec: "0001-widget-flow"},
		{name: "out-of-range number passes through", input: "9\nclaude\n\n\n", wantSpec: "9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder

			values, err := CollectInput(context.Background(), InputRequest{
				Command:           "implement",
				Values:            CommandValues{Agent: "codex"},
				SelectionDefaults: testSelectionDefaults(),
				SpecOptions:       []string{"0001-widget-flow", "0002-other-flow"},
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
				"Agent Model [opus]:",
				"Default Reasoning Effort [high]:",
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
		{name: "yes enables QA", input: "\n\n\n\ny\n", wantQA: true, wantPrompt: "QA gate [y/N]:"},
		{name: "empty keeps QA disabled", input: "\n\n\n\n\n", wantQA: false, wantPrompt: "QA gate [y/N]:"},
		{name: "empty keeps QA flag default", input: "\n\n\n\n\n", defaultQA: true, wantQA: true, wantPrompt: "QA gate [Y/n]:"},
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
				SelectionDefaults: testSelectionDefaults(),
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
		Command:           "implement",
		Values:            CommandValues{Spec: "0001-widget-flow", Agent: "codex"},
		SelectionDefaults: testSelectionDefaults(),
	}, strings.NewReader("\n\n\n\nmaybe\nlater\n"), &output)

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

func TestCollectInputDisplaysCodexCatalogAndMapsNumbers(t *testing.T) {
	var output strings.Builder

	values, err := CollectInput(context.Background(), InputRequest{
		Command: "resolve",
		Values: CommandValues{
			PRNumber:    "123",
			Agent:       "codex",
			Round:       "all",
			ArtifactDir: ".roundfix",
		},
		SelectionDefaults: testSelectionDefaults(),
	}, strings.NewReader("\n\n1\n4\n\n\n"), &output)
	if err != nil {
		t.Fatalf("collect input: %v", err)
	}

	if values.Model != "gpt-5.6-sol" {
		t.Fatalf("expected first Codex catalog model, got %q", values.Model)
	}
	if values.ReasoningEffort != "xhigh" {
		t.Fatalf("expected fourth Codex reasoning choice, got %q", values.ReasoningEffort)
	}
	assertContainsInOrder(t, output.String(),
		"Agent Model Choices (codex):",
		"1. gpt-5.6-sol",
		"2. gpt-5.6-terra",
		"3. gpt-5.6-luna",
		"4. gpt-5.5",
		"5. gpt-5.4",
		"6. gpt-5.4-mini",
		"7. gpt-5.3-codex-spark",
	)
	if strings.Contains(output.String(), "Custom") {
		t.Fatalf("expected no synthetic Custom catalog entry, got:\n%s", output.String())
	}
}

func TestCollectInputDisplaysClaudeCatalogDefaultAsConcreteModel(t *testing.T) {
	defaults := testSelectionDefaults()
	claude := defaults["claude"]
	claude.Model = "claude-project-default"
	defaults["claude"] = claude
	var output strings.Builder

	values, err := CollectInput(context.Background(), InputRequest{
		Command: "resolve",
		Values: CommandValues{
			PRNumber:    "123",
			Agent:       "claude",
			Round:       "all",
			ArtifactDir: ".roundfix",
		},
		SelectionDefaults: defaults,
	}, strings.NewReader("\n\n1\n2\n\n\n"), &output)
	if err != nil {
		t.Fatalf("collect input: %v", err)
	}

	if values.Model != "claude-project-default" {
		t.Fatalf("expected Claude Default to resolve to configured model, got %q", values.Model)
	}
	if values.ReasoningEffort != "high" {
		t.Fatalf("expected Claude reasoning choice, got %q", values.ReasoningEffort)
	}
	assertContainsInOrder(t, output.String(),
		"Agent Model Choices (claude):",
		"1. Default -> claude-project-default",
		"2. Opus -> opus",
		"3. Fable -> fable",
		"4. Sonnet -> sonnet",
		"5. Haiku -> haiku",
	)
	if strings.Contains(output.String(), "Custom") {
		t.Fatalf("expected no synthetic Custom catalog entry, got:\n%s", output.String())
	}
}

func TestCollectInputPreservesCustomModelAndReasoningValues(t *testing.T) {
	values, err := CollectInput(context.Background(), InputRequest{
		Command: "resolve",
		Values: CommandValues{
			PRNumber:    "123",
			Agent:       "codex",
			Round:       "all",
			ArtifactDir: ".roundfix",
		},
		SelectionDefaults: testSelectionDefaults(),
	}, strings.NewReader("\n\nfuture-model\nexperimental-reasoning\n\n\n"), io.Discard)
	if err != nil {
		t.Fatalf("collect input: %v", err)
	}
	if values.Model != "future-model" || values.ReasoningEffort != "experimental-reasoning" {
		t.Fatalf("expected custom values to survive, got %#v", values)
	}
}

func TestCollectInputOpenCodeRequiresTypedOrConfiguredSelectionValues(t *testing.T) {
	tests := []struct {
		name     string
		defaults RuntimeSelectionDefaults
		input    string
		want     string
	}{
		{
			name:  "missing model",
			input: "\n\n\n",
			want:  "Agent Model",
		},
		{
			name:     "missing reasoning",
			defaults: RuntimeSelectionDefaults{Model: "opencode-model"},
			input:    "\n\n\n\n",
			want:     "Default Reasoning Effort",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults := testSelectionDefaults()
			defaults["opencode"] = tt.defaults
			var output strings.Builder

			_, err := CollectInput(context.Background(), InputRequest{
				Command:           "implement",
				Values:            CommandValues{Spec: "0001-widget-flow", Agent: "opencode"},
				SelectionDefaults: defaults,
			}, strings.NewReader(tt.input), &output)

			if err == nil {
				t.Fatal("expected missing OpenCode selection to fail")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error to name %q, got %v", tt.want, err)
			}
			if strings.Contains(output.String(), "Agent Model Choices (opencode):") {
				t.Fatalf("expected no fabricated OpenCode catalog, got:\n%s", output.String())
			}
		})
	}
}

func TestRenderLiveRunViewGroupsIssuesAndShowsStatusStrips(t *testing.T) {
	view := RenderLiveRunView(LiveRunView{
		Command:         "resolve",
		Repository:      "owner/project",
		PRNumber:        "123",
		HeadBranch:      "feature/review",
		ReviewSource:    "CodeRabbit",
		Agent:           "Codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "xhigh",
		HEAD:            "abc123",
		RunID:           "run_123",
		PipelineState:   "ResolvingWithAgent",
		BudgetState:     "38m / 2h",
		GitState:        "clean, 1 unpushed commit",
		CurrentRound:    2,
		MaxRounds:       6,
		AutoCommit:      true,
		AutoPush:        true,
		LastPush:        "pending",
		Width:           100,
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
		"Agent Model: gpt-5.5",
		"Default Reasoning Effort: xhigh",
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
		Command:         "implement",
		RunKind:         store.KindImplement,
		SpecSlug:        "0001-widget-flow",
		GitRoot:         "/repo",
		WorkDir:         "/home/user/.roundfix/worktrees/repo/run_9",
		HeadBranch:      "ma/widget-flow",
		Agent:           "Codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "xhigh",
		HEAD:            "abc123",
		RunID:           "run_9",
		PipelineState:   "ResolvingWithAgent",
		Concurrency:     2,
		BudgetState:     "38m / 2h",
		GitState:        "clean, 1 unpushed commit",
		CurrentRound:    2,
		MaxRounds:       6,
		AutoCommit:      true,
		AutoPush:        false,
		LastPush:        "disabled",
		Width:           100,
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
		"Agent Model: gpt-5.5",
		"Default Reasoning Effort: xhigh",
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

func TestRenderLiveRunViewShowsLegacyEmptySelectionAsDash(t *testing.T) {
	view := RenderLiveRunView(LiveRunView{
		Command:       "attach",
		Repository:    "owner/project",
		PRNumber:      "123",
		HeadBranch:    "feature/review",
		ReviewSource:  "CodeRabbit",
		Agent:         "Codex",
		RunID:         "run_legacy",
		PipelineState: "Clean",
	})

	for _, expected := range []string{
		"Agent: Codex",
		"Agent Model: -",
		"Default Reasoning Effort: -",
	} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected legacy selection display %q, got:\n%s", expected, view)
		}
	}
}

func TestRenderAgentTimelineDoesNotRenderAutoModelPlaceholder(t *testing.T) {
	timeline := stripANSI(renderAgentTimeline(LiveRunView{
		Agent:           "Codex",
		Model:           "",
		ReasoningEffort: "",
	}, nil, 80, 8))

	if strings.Contains(timeline, "auto") {
		t.Fatalf("expected empty selection to render as dash, got:\n%s", timeline)
	}
	if !strings.Contains(timeline, "0 entries · Codex · - · -") {
		t.Fatalf("expected model and reasoning dash placeholders, got:\n%s", timeline)
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

func testSelectionDefaults() map[string]RuntimeSelectionDefaults {
	return map[string]RuntimeSelectionDefaults{
		"codex": {
			Model:           "gpt-5.5",
			ReasoningEffort: "xhigh",
			ModelCatalog: []ModelChoice{
				{Label: "gpt-5.6-sol", Value: "gpt-5.6-sol"},
				{Label: "gpt-5.6-terra", Value: "gpt-5.6-terra"},
				{Label: "gpt-5.6-luna", Value: "gpt-5.6-luna"},
				{Label: "gpt-5.5", Value: "gpt-5.5"},
				{Label: "gpt-5.4", Value: "gpt-5.4"},
				{Label: "gpt-5.4-mini", Value: "gpt-5.4-mini"},
				{Label: "gpt-5.3-codex-spark", Value: "gpt-5.3-codex-spark"},
			},
			ReasoningChoices: []string{"low", "medium", "high", "xhigh"},
		},
		"claude": {
			Model:           "opus",
			ReasoningEffort: "high",
			ModelCatalog: []ModelChoice{
				{Label: "Default", Value: "default"},
				{Label: "Opus", Value: "opus"},
				{Label: "Fable", Value: "fable"},
				{Label: "Sonnet", Value: "sonnet"},
				{Label: "Haiku", Value: "haiku"},
			},
			ReasoningChoices: []string{"default", "high", "maximum"},
		},
		"opencode": {},
	}
}
