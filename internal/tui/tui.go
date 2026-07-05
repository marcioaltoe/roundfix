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
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type CommandValues struct {
	PRNumber     string
	Spec         string
	ReviewSource string
	Agent        string
	Round        string
	ArtifactDir  string
	Model        string
	MaxRounds    int
	UntilClean   bool
	QA           bool
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
	// SpecOptions lists the active Spec slugs the Spec picker offers, in
	// slug order. The spec field accepts a 1-based number into this list
	// or a slug typed directly.
	SpecOptions []string
}

type LiveRunView struct {
	Command       string
	Repository    string
	PRNumber      string
	HeadBranch    string
	ReviewSource  string
	Agent         string
	Model         string
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
	BatchNumber   int
	BatchTotal    int
	TotalIssues   int
	Issues        []rounds.Issue
	// RunKind selects the Work Item vocabulary of the panes: implement Runs
	// render Tasks where review Runs render Review Issues. Empty means a
	// review Run, so existing callers keep their rendering unchanged.
	RunKind string
	// SpecSlug and GitRoot locate a spec Run's task files
	// (docs/specs/<slug>/ under the git root) so the cockpit can refresh
	// Task statuses by re-reading them.
	SpecSlug string
	GitRoot  string
	// Tasks lists the spec Run's Tasks in Task Graph order.
	Tasks []spec.Task
	// BatchSizes lists the planned Review Issue count per Batch, in Batch
	// order, when the caller knows the plan. The cockpit derives Batch
	// separators and Executing/Waiting states from it.
	BatchSizes []int
	Console    []string
	Width      int
}

type IssueGroup struct {
	Round  int
	Issues []rounds.Issue
}

// WorkItem is the Run-Kind-neutral view model of the work-item pane: the
// display name and current status of one Work Item. Review Runs map Review
// Issues into it and spec Runs map Tasks; the Run Kind keys the mapping, so
// each pane renders exactly one vocabulary.
type WorkItem struct {
	Name   string // "Issue #001" for Review Issues, the Task id for Tasks
	Title  string
	Status string // artifact status verbatim (pending, resolved, completed, ...)
}

// TaskWorkItems maps a spec Run's Tasks into Work Items, preserving Task
// Graph order.
func TaskWorkItems(tasks []spec.Task) []WorkItem {
	items := make([]WorkItem, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, WorkItem{Name: task.ID, Title: strings.TrimSpace(task.Title), Status: string(task.Status)})
	}
	return items
}

// specRunView reports whether the view renders a spec Run's Tasks; every
// other Run Kind keeps the Review Issue rendering byte-identical.
func specRunView(view LiveRunView) bool {
	return view.RunKind == store.KindImplement
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
		if field == "qa" {
			qa, done, err := collectQAGateInput(reader, output, current == "yes")
			if err != nil {
				return CommandValues{}, err
			}
			values.QA = qa
			if done {
				break
			}
			continue
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
		if field == "spec" {
			line = pickSpecOption(line, req.SpecOptions)
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
		"spec":         req.Values.Spec,
		"source":       req.Values.ReviewSource,
		"agent":        req.Values.Agent,
		"round":        req.Values.Round,
		"artifact-dir": req.Values.ArtifactDir,
		"model":        req.Values.Model,
		"max-rounds":   "",
		"qa":           "",
	}
	if req.Values.MaxRounds > 0 {
		defaults["max-rounds"] = strconv.Itoa(req.Values.MaxRounds)
	}
	if req.Values.QA {
		defaults["qa"] = "yes"
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
	if len(req.SpecOptions) > 0 {
		builder.WriteString("Active Specs:\n")
		for index, slug := range req.SpecOptions {
			builder.WriteString(fmt.Sprintf("  %d. %s\n", index+1, slug))
		}
		builder.WriteString("Pick a Spec by number or slug.\n")
	}
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
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Roundfix %s\n\n", strings.ToLower(emptyDash(view.Command))))
	builder.WriteString("Target:\n")
	if specRunView(view) {
		builder.WriteString(fmt.Sprintf("  Spec: %s\n", emptyDash(view.SpecSlug)))
		builder.WriteString(fmt.Sprintf("  Branch: %s\n", emptyDash(view.HeadBranch)))
	} else {
		builder.WriteString(fmt.Sprintf("  PR: #%s %s\n", emptyDash(view.PRNumber), emptyDash(view.Repository)))
		builder.WriteString(fmt.Sprintf("  Branch: %s\n", emptyDash(view.HeadBranch)))
		builder.WriteString(fmt.Sprintf("  Source: %s\n", emptyDash(view.ReviewSource)))
	}
	if view.Agent != "" {
		builder.WriteString(fmt.Sprintf("  Agent: %s\n", emptyDash(view.Agent)))
	}
	if view.HEAD != "" {
		builder.WriteString(fmt.Sprintf("  HEAD: %s\n", view.HEAD))
	}

	builder.WriteString("\nRun:\n")
	builder.WriteString(fmt.Sprintf("  ID: %s\n", emptyDash(view.RunID)))
	builder.WriteString(fmt.Sprintf("  State: %s\n", emptyDash(view.PipelineState)))
	if !specRunView(view) {
		builder.WriteString(fmt.Sprintf("  Round: %s\n", formatRound(view.CurrentRound, view.MaxRounds)))
		builder.WriteString(fmt.Sprintf("  Budget: %s\n", emptyDash(view.BudgetState)))
	}
	builder.WriteString(fmt.Sprintf("  Git: %s\n", emptyDash(view.GitState)))
	builder.WriteString(fmt.Sprintf("  Auto-commit: %s\n", onOff(view.AutoCommit)))
	builder.WriteString(fmt.Sprintf("  Auto-push: %s\n", onOff(view.AutoPush)))
	builder.WriteString(fmt.Sprintf("  Last push: %s\n", emptyDash(view.LastPush)))

	builder.WriteString("\n")
	builder.WriteString(renderSplitPanes(view))
	builder.WriteByte('\n')
	return builder.String()
}

func renderSplitPanes(view LiveRunView) string {
	width := view.Width
	if width <= 0 {
		width = 100
	}
	if width < 80 {
		width = 80
	}
	available := width - 7
	leftWidth := available * 42 / 100
	if leftWidth < 30 {
		leftWidth = 30
	}
	rightWidth := available - leftWidth
	if rightWidth < 30 {
		rightWidth = 30
		leftWidth = available - rightWidth
	}

	leftTitle := "Review Issues"
	leftLines := issuePaneLines(view.Issues)
	if specRunView(view) {
		leftTitle = "Tasks"
		leftLines = taskPaneLines(TaskWorkItems(view.Tasks))
	}
	rightLines := consolePaneLines(view.Console)
	rowCount := len(leftLines)
	if len(rightLines) > rowCount {
		rowCount = len(rightLines)
	}

	var builder strings.Builder
	builder.WriteString(paneBorder(leftWidth, rightWidth))
	builder.WriteString(paneRow(leftTitle, "Agent Console", leftWidth, rightWidth))
	builder.WriteString(paneBorder(leftWidth, rightWidth))
	for index := 0; index < rowCount; index++ {
		left := ""
		right := ""
		if index < len(leftLines) {
			left = leftLines[index]
		}
		if index < len(rightLines) {
			right = rightLines[index]
		}
		builder.WriteString(paneRow(left, right, leftWidth, rightWidth))
	}
	builder.WriteString(paneBorder(leftWidth, rightWidth))
	builder.WriteString("Keys: Ctrl-C stop\n")
	return builder.String()
}

func issuePaneLines(issues []rounds.Issue) []string {
	if len(issues) == 0 {
		return []string{"none"}
	}
	lines := []string{}
	for _, group := range GroupIssuesByRound(issues) {
		lines = append(lines, fmt.Sprintf("Round %03d", group.Round))
		for _, issue := range group.Issues {
			location := emptyDash(issue.File)
			if issue.Line > 0 {
				location = fmt.Sprintf("%s:%d", location, issue.Line)
			}
			lines = append(lines, fmt.Sprintf("%-8s %-10s %s", emptyDash(issue.Severity), emptyDash(issue.Status), location))
			if strings.TrimSpace(issue.Title) != "" {
				lines = append(lines, "  "+strings.TrimSpace(issue.Title))
			}
		}
	}
	return lines
}

// taskPaneLines renders Task Work Items one per line in Task Graph order,
// mirroring the implement stdout contract shape: `task_NN <status> — <title>`.
func taskPaneLines(items []WorkItem) []string {
	if len(items) == 0 {
		return []string{"none"}
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("%s %s — %s", item.Name, emptyDash(item.Status), emptyDash(item.Title)))
	}
	return lines
}

func consolePaneLines(console []string) []string {
	lines := make([]string, 0, len(console))
	for _, line := range console {
		line = strings.TrimRight(line, "\r\n")
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return []string{"waiting for output"}
	}
	return lines
}

func paneBorder(leftWidth int, rightWidth int) string {
	return "+" + strings.Repeat("-", leftWidth+2) + "+" + strings.Repeat("-", rightWidth+2) + "+\n"
}

func paneRow(left string, right string, leftWidth int, rightWidth int) string {
	return fmt.Sprintf("| %s | %s |\n", padRight(truncateText(left, leftWidth), leftWidth), padRight(truncateText(right, rightWidth), rightWidth))
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func truncateText(value string, width int) string {
	if width <= 0 || len(value) <= width {
		return value
	}
	if width <= 3 {
		return value[:width]
	}
	return value[:width-3] + "..."
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
	case "implement":
		return []string{"spec", "agent", "qa"}
	default:
		return []string{"pr"}
	}
}

// pickSpecOption maps a 1-based picker number onto the listed Spec slug and
// passes any other entry through unchanged, so a known slug can be typed
// directly.
func pickSpecOption(line string, options []string) string {
	index, err := strconv.Atoi(line)
	if err != nil || index < 1 || index > len(options) {
		return line
	}
	return options[index-1]
}

func inputLabel(field string) string {
	switch field {
	case "pr":
		return "Open Pull Request"
	case "spec":
		return "Spec"
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
	case "qa":
		return "QA gate"
	default:
		return field
	}
}

func getValue(values CommandValues, field string) string {
	switch field {
	case "pr":
		return values.PRNumber
	case "spec":
		return values.Spec
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
	case "qa":
		if values.QA {
			return "yes"
		}
	}
	return ""
}

func setValue(values *CommandValues, field string, value string) error {
	switch field {
	case "pr":
		values.PRNumber = value
	case "spec":
		values.Spec = value
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
	case "qa":
		qa, ok := parseQAGateChoice(value)
		if !ok {
			return fmt.Errorf("QA gate must be y, yes, n, no, or empty")
		}
		values.QA = qa
	}
	return nil
}

func collectQAGateInput(reader *bufio.Reader, output io.Writer, current bool) (bool, bool, error) {
	for attempt := 0; attempt < 2; attempt++ {
		fmt.Fprint(output, qaGatePrompt(current))
		line, err := reader.ReadString('\n')
		done := err == io.EOF
		if err != nil && err != io.EOF {
			return false, done, fmt.Errorf("read Interactive Input field QA gate: %w", err)
		}
		choice := strings.TrimSpace(line)
		if choice == "" {
			return current, done, nil
		}
		qa, ok := parseQAGateChoice(choice)
		if ok {
			return qa, done, nil
		}
		if attempt == 0 && !done {
			fmt.Fprintln(output, "QA gate must be y, yes, n, no, or empty.")
			continue
		}
		return false, done, fmt.Errorf("QA gate must be y, yes, n, no, or empty")
	}
	return false, false, fmt.Errorf("QA gate must be y, yes, n, no, or empty")
}

func qaGatePrompt(current bool) string {
	if current {
		return "QA gate [Y/n]: "
	}
	return "QA gate [y/N]: "
}

func parseQAGateChoice(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y", "yes":
		return true, true
	case "n", "no":
		return false, true
	default:
		return false, false
	}
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
