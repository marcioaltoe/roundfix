package tui

import (
	"context"
	"strings"
	"testing"

	"roundfix/internal/rounds"
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
		Issues: []rounds.Issue{
			{Round: 2, Severity: "minor", Status: rounds.StatusPending, File: "README.md", Line: 12},
			{Round: 1, Severity: "major", Status: rounds.StatusResolved, File: "api/auth.go", Line: 88},
			{Round: 2, Severity: "major", Status: rounds.StatusValid, File: "src/cache.ts", Line: 41},
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
		"Progress:",
		"codex resolving batch 1/2",
		"running make verify",
		"Review Issues:",
		"Round 001",
		"major    resolved   api/auth.go:88",
		"Round 002",
		"major    valid      src/cache.ts:41",
		"minor    pending    README.md:12",
	}
	for _, text := range expected {
		if !strings.Contains(view, text) {
			t.Fatalf("expected live view to contain %q, got:\n%s", text, view)
		}
	}
	for _, removed := range []string{"[tab] focus", "[s] stop", "Console"} {
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
