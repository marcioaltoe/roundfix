package tui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"roundfix/internal/rounds"
)

type CommandValues struct {
	PRNumber     string
	ReviewSource string
	Agent        string
	Round        string
	ArtifactDir  string
	Model        string
	MaxRounds    int
	UntilClean   bool
}

type Suggestion struct {
	Value  string
	Source string
}

type InputRequest struct {
	Command         string
	Values          CommandValues
	PRSuggestion    Suggestion
	AgentSuggestion Suggestion
}

type LiveRunView struct {
	Repository    string
	PRNumber      string
	HeadBranch    string
	ReviewSource  string
	Agent         string
	HEAD          string
	RunID         string
	PipelineState string
	BudgetState   string
	GitState      string
	CurrentRound  int
	MaxRounds     int
	AutoCommit    bool
	AutoPush      bool
	LastPush      string
	Issues        []rounds.Issue
	Console       []string
	Keybindings   []string
}

type IssueGroup struct {
	Round  int
	Issues []rounds.Issue
}

func CollectInput(ctx context.Context, req InputRequest, input io.Reader, output io.Writer) (CommandValues, error) {
	if input == nil {
		input = strings.NewReader("")
	}
	if output == nil {
		output = io.Discard
	}
	values := req.Values
	defaults := DefaultsForInput(req)
	reader := bufio.NewReader(input)

	fmt.Fprint(output, RenderInteractiveInput(req))
	for _, field := range fieldsForCommand(req.Command) {
		if err := ctx.Err(); err != nil {
			return CommandValues{}, err
		}
		current := getValue(values, field)
		if current == "" || req.Values == values {
			current = defaults[field]
		}
		fmt.Fprintf(output, "%s [%s]: ", inputLabel(field), current)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return CommandValues{}, fmt.Errorf("read Interactive Input field %s: %w", field, err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			line = current
		}
		if line != "" {
			if err := setValue(&values, field, line); err != nil {
				return CommandValues{}, err
			}
		}
		if err == io.EOF {
			break
		}
	}
	return values, nil
}

func DefaultsForInput(req InputRequest) map[string]string {
	defaults := map[string]string{
		"pr":           req.Values.PRNumber,
		"source":       req.Values.ReviewSource,
		"agent":        req.Values.Agent,
		"round":        req.Values.Round,
		"artifact-dir": req.Values.ArtifactDir,
		"model":        req.Values.Model,
		"max-rounds":   "",
	}
	if req.Values.MaxRounds > 0 {
		defaults["max-rounds"] = strconv.Itoa(req.Values.MaxRounds)
	}
	if defaults["pr"] == "" {
		defaults["pr"] = req.PRSuggestion.Value
	}
	if defaults["agent"] == "" {
		defaults["agent"] = req.AgentSuggestion.Value
	}
	return defaults
}

func RenderInteractiveInput(req InputRequest) string {
	defaults := DefaultsForInput(req)
	var builder strings.Builder
	builder.WriteString("Roundfix Interactive Input\n")
	builder.WriteString(fmt.Sprintf("Command: %s\n", req.Command))
	if defaults["pr"] != "" {
		source := req.PRSuggestion.Source
		if source == "" && req.Values.PRNumber != "" {
			source = "provided"
		}
		builder.WriteString(fmt.Sprintf("Suggested Open Pull Request: #%s", defaults["pr"]))
		if source != "" {
			builder.WriteString(fmt.Sprintf(" (%s)", source))
		}
		builder.WriteByte('\n')
	}
	if defaults["agent"] != "" {
		source := req.AgentSuggestion.Source
		if req.Values.Agent != "" {
			source = "config"
		}
		builder.WriteString(fmt.Sprintf("Suggested Agent: %s", defaults["agent"]))
		if source != "" {
			builder.WriteString(fmt.Sprintf(" (%s)", source))
		}
		builder.WriteByte('\n')
	}
	builder.WriteString("Press Enter to accept a suggestion.\n")
	return builder.String()
}

func RenderLiveRunView(view LiveRunView) string {
	keybindings := view.Keybindings
	if len(keybindings) == 0 {
		keybindings = []string{"[tab] focus", "[j/k] issue", "[f] fetch", "[r] resolve", "[p] push", "[t] trigger", "[d] detach", "[s] stop", "[q] quit"}
	}
	var builder strings.Builder
	builder.WriteString("Roundfix\n")
	builder.WriteString(fmt.Sprintf("repo: %s     pr: #%s     head: %s     source: %s     agent: %s     head_sha: %s\n", emptyDash(view.Repository), emptyDash(view.PRNumber), emptyDash(view.HeadBranch), emptyDash(view.ReviewSource), emptyDash(view.Agent), emptyDash(view.HEAD)))
	builder.WriteString(fmt.Sprintf("run: %s     state: %s     round: %s     budget: %s\n", emptyDash(view.RunID), emptyDash(view.PipelineState), formatRound(view.CurrentRound, view.MaxRounds), emptyDash(view.BudgetState)))
	builder.WriteString(fmt.Sprintf("git: %s     auto-commit: %s     auto-push: %s     last push: %s\n", emptyDash(view.GitState), onOff(view.AutoCommit), onOff(view.AutoPush), emptyDash(view.LastPush)))
	builder.WriteString(strings.Repeat("-", 80))
	builder.WriteByte('\n')
	builder.WriteString("Review Issues\n")
	for _, group := range GroupIssuesByRound(view.Issues) {
		builder.WriteString(fmt.Sprintf("Round %03d\n", group.Round))
		for _, issue := range group.Issues {
			builder.WriteString(fmt.Sprintf("  %-8s %-10s %s:%d\n", emptyDash(issue.Severity), emptyDash(issue.Status), emptyDash(issue.File), issue.Line))
		}
	}
	if len(view.Issues) == 0 {
		builder.WriteString("  none\n")
	}
	builder.WriteString("\nConsole\n")
	for _, line := range view.Console {
		builder.WriteString("  ")
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	if len(view.Console) == 0 {
		builder.WriteString("  waiting for output\n")
	}
	builder.WriteString("\n")
	builder.WriteString(strings.Join(keybindings, "   "))
	builder.WriteByte('\n')
	return builder.String()
}

func GroupIssuesByRound(issues []rounds.Issue) []IssueGroup {
	sorted := append([]rounds.Issue{}, issues...)
	sort.SliceStable(sorted, func(i, j int) bool {
		left := sorted[i]
		right := sorted[j]
		if left.Round != right.Round {
			return left.Round < right.Round
		}
		if left.Severity != right.Severity {
			return left.Severity < right.Severity
		}
		if left.Status != right.Status {
			return left.Status < right.Status
		}
		if left.File != right.File {
			return left.File < right.File
		}
		return left.Line < right.Line
	})
	groups := make([]IssueGroup, 0)
	for _, issue := range sorted {
		if len(groups) == 0 || groups[len(groups)-1].Round != issue.Round {
			groups = append(groups, IssueGroup{Round: issue.Round})
		}
		groups[len(groups)-1].Issues = append(groups[len(groups)-1].Issues, issue)
	}
	return groups
}

type StreamBuffer struct {
	MaxLines int
	lines    []string
	pending  string
}

func (buffer *StreamBuffer) Write(payload []byte) (int, error) {
	text := buffer.pending + string(payload)
	parts := strings.Split(text, "\n")
	buffer.pending = parts[len(parts)-1]
	for _, line := range parts[:len(parts)-1] {
		buffer.append(line)
	}
	return len(payload), nil
}

func (buffer *StreamBuffer) Lines() []string {
	lines := append([]string{}, buffer.lines...)
	if buffer.pending != "" {
		lines = append(lines, buffer.pending)
	}
	return lines
}

func (buffer *StreamBuffer) append(line string) {
	buffer.lines = append(buffer.lines, line)
	max := buffer.MaxLines
	if max <= 0 {
		max = 200
	}
	if len(buffer.lines) > max {
		buffer.lines = append([]string{}, buffer.lines[len(buffer.lines)-max:]...)
	}
}

func fieldsForCommand(command string) []string {
	switch command {
	case "fetch":
		return []string{"pr", "source", "round", "artifact-dir"}
	case "resolve":
		return []string{"pr", "agent", "round", "artifact-dir", "model"}
	case "watch":
		return []string{"pr", "source", "agent", "artifact-dir", "model", "max-rounds"}
	default:
		return []string{"pr"}
	}
}

func inputLabel(field string) string {
	switch field {
	case "pr":
		return "Open Pull Request"
	case "source":
		return "Review Source"
	case "agent":
		return "Agent"
	case "round":
		return "Round"
	case "artifact-dir":
		return "Artifact Directory"
	case "model":
		return "Model"
	case "max-rounds":
		return "Max Rounds"
	default:
		return field
	}
}

func getValue(values CommandValues, field string) string {
	switch field {
	case "pr":
		return values.PRNumber
	case "source":
		return values.ReviewSource
	case "agent":
		return values.Agent
	case "round":
		return values.Round
	case "artifact-dir":
		return values.ArtifactDir
	case "model":
		return values.Model
	case "max-rounds":
		if values.MaxRounds > 0 {
			return strconv.Itoa(values.MaxRounds)
		}
	}
	return ""
}

func setValue(values *CommandValues, field string, value string) error {
	switch field {
	case "pr":
		values.PRNumber = value
	case "source":
		values.ReviewSource = value
	case "agent":
		values.Agent = value
	case "round":
		values.Round = value
	case "artifact-dir":
		values.ArtifactDir = value
	case "model":
		values.Model = value
	case "max-rounds":
		number, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Max Rounds must be a number: %w", err)
		}
		values.MaxRounds = number
	}
	return nil
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func onOff(value bool) string {
	if value {
		return "on"
	}
	return "off"
}

func formatRound(current int, max int) string {
	if current <= 0 && max <= 0 {
		return "-"
	}
	if max <= 0 {
		return strconv.Itoa(current)
	}
	return fmt.Sprintf("%d / %d", current, max)
}
