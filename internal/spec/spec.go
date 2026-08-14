// Package spec owns the on-disk Spec contract: discovering active Specs
// under a Spec Root, loading and validating a Task Graph into deterministic
// topological order, parsing task files, rewriting task status, and reading
// the QA Report verdict. Nothing outside this package touches spec markdown.
package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Task status values; the task file is the sole owner of its status. The
// accepted synonym set is intentionally exact after trimming: "done" maps to
// completed, and "in-progress"/"in progress" map to in_progress.
const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

const (
	manifestSchema  = "spec-tasks/v1"
	prdStatusActive = "active"
	archivedDirName = "_archived"
)

// Status is the lifecycle state a task file carries in its frontmatter.
type Status string

var statusSynonyms = map[string]Status{
	"done":        StatusCompleted,
	"in-progress": StatusInProgress,
	"in progress": StatusInProgress,
}

var canonicalStatuses = []Status{StatusPending, StatusInProgress, StatusCompleted, StatusFailed}

// NormalizeStatus maps documented synonyms to canonical statuses. Unknown
// values pass through trimmed and fail AllowedStatus as today.
func NormalizeStatus(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if status, ok := statusSynonyms[trimmed]; ok {
		return string(status)
	}
	return trimmed
}

func allowedStatusValues() string {
	values := make([]string, 0, len(canonicalStatuses))
	for _, status := range canonicalStatuses {
		values = append(values, string(status))
	}
	return strings.Join(values, ", ")
}

// AllowedStatus reports whether the status is one of the four task lifecycle
// states.
func AllowedStatus(status Status) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

// Spec identifies one Spec directory under a Spec Root.
type Spec struct {
	Slug string
	Dir  string
}

// SkippedSpec reports a Spec directory that discovery intentionally did not
// offer to the Implement picker, with the reason a user can fix.
type SkippedSpec struct {
	Dir    string
	Reason string
}

// Task is one Task Graph node joined with its parsed task file. File is the
// task file path relative to the Spec Root.
type Task struct {
	ID               string
	File             string
	Title            string
	TitleLine        int
	Needs            []string
	Status           Status
	StatusNormalized bool
	Type             TaskType
	Context          []TaskContextRef
	Requirements     []TaskDeclaration
	RehearsalCases   []TaskDeclaration
	Verification     []string
	NegativeControl  []string
}

// CarryForward is one completed Task commit a stopped Run can hand back to
// the checkout. MovedInputs names every declared input whose bytes no longer
// match the tree the Task saw before its settlement commit.
type CarryForward struct {
	TaskID        string   `json:"taskId"`
	RunID         string   `json:"runId"`
	Commit        string   `json:"commit"`
	TaskFile      string   `json:"taskFile"`
	InputsMoved   bool     `json:"inputsMoved"`
	MovedInputs   []string `json:"movedInputs"`
	Action        string   `json:"action"`
	RefusalReason string   `json:"refusalReason"`
}

// TaskDeclaration is one author-written Task declaration and its 1-based
// Markdown source line.
type TaskDeclaration struct {
	Text string
	Line int
}

// ContextKind classifies a Task-authored context path.
type ContextKind string

const (
	ContextKindInstruction ContextKind = "instruction"
	ContextKindInterface   ContextKind = "interface"
	ContextKindCreates     ContextKind = "creates"
)

// TaskContextRef is one labeled repository-relative path from a Task's
// optional ## Context section.
type TaskContextRef struct {
	Kind ContextKind
	Path string
}

// Graph is a validated Task Graph with Tasks in deterministic topological
// order (Kahn's algorithm, manifest node order as tiebreak).
type Graph struct {
	Spec       Spec
	Tasks      []Task
	QATaskID   string
	QADeclined bool
	QAReason   string
}

// UnreachableDeclaration is one acceptance the Spec's author declared beyond
// the reach of any hermetic Verification before the gate ran.
type UnreachableDeclaration struct {
	Criterion   string
	Reason      string
	SatisfiedBy string
	Line        int
}

type unreachableDeclarationBuilder struct {
	declaration UnreachableDeclaration
	activeField string
}

type manifestNode struct {
	ID    string   `yaml:"id"`
	File  string   `yaml:"file"`
	Needs []string `yaml:"needs"`
}

type manifestFrontmatter struct {
	Schema   string                 `yaml:"schema"`
	QA       manifestOptionalString `yaml:"qa"`
	QAReason manifestOptionalString `yaml:"qa_reason"`
	Graph    struct {
		Nodes []manifestNode `yaml:"nodes"`
	} `yaml:"graph"`
}

func (manifest *manifestFrontmatter) UnmarshalYAML(node *yaml.Node) error {
	type plainManifest manifestFrontmatter
	var decoded plainManifest
	if err := node.Decode(&decoded); err != nil {
		return fmt.Errorf("decode manifest frontmatter: %w", err)
	}
	*manifest = manifestFrontmatter(decoded)
	for index := 0; index+1 < len(node.Content); index += 2 {
		switch node.Content[index].Value {
		case "qa":
			manifest.QA.Present = true
		case "qa_reason":
			manifest.QAReason.Present = true
		}
	}
	return nil
}

type manifestOptionalString struct {
	Value   string
	Present bool
}

func (value *manifestOptionalString) UnmarshalYAML(node *yaml.Node) error {
	value.Present = true
	if node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
		return nil
	}
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		return fmt.Errorf("must be a string")
	}
	value.Value = node.Value
	return nil
}

type qaDeclaration struct {
	Present  bool
	TaskID   string
	Declined bool
	Reason   string
}

// Unreachable reads declarations from the Spec PRD's optional
// ## Unreachable Acceptance section. A Spec with no such section returns no
// declarations and no error.
func Unreachable(specDir string) ([]UnreachableDeclaration, error) {
	prdPath := filepath.Join(specDir, "_prd.md")
	content, err := os.ReadFile(prdPath)
	if err != nil {
		return nil, fmt.Errorf("read Spec PRD %q: %w", prdPath, err)
	}
	return parseUnreachableDeclarations(content, prdPath)
}

func parseUnreachableDeclarations(content []byte, prdPath string) ([]UnreachableDeclaration, error) {
	var declarations []UnreachableDeclaration
	var current *unreachableDeclarationBuilder
	inSection := false
	flush := func() error {
		if current == nil {
			return nil
		}
		declaration, err := current.finish(prdPath)
		if err != nil {
			return err
		}
		declarations = append(declarations, declaration)
		current = nil
		return nil
	}

	for index, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") || strings.HasPrefix(trimmed, "## ") {
			if inSection {
				if err := flush(); err != nil {
					return nil, err
				}
			}
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Unreachable Acceptance"
			continue
		}
		if !inSection {
			continue
		}

		if strings.HasPrefix(line, "- ") {
			if err := flush(); err != nil {
				return nil, err
			}
			current = &unreachableDeclarationBuilder{
				declaration: UnreachableDeclaration{Line: index + 1},
			}
			current.readField(strings.TrimSpace(strings.TrimPrefix(line, "- ")))
			continue
		}
		if current == nil || trimmed == "" {
			continue
		}
		if current.readField(trimmed) {
			continue
		}
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			current.appendContinuation(trimmed)
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	return declarations, nil
}

func (builder *unreachableDeclarationBuilder) readField(line string) bool {
	label, value, ok := strings.Cut(line, ":")
	if !ok {
		return false
	}
	field := strings.TrimSpace(label)
	switch field {
	case "criterion":
		builder.declaration.Criterion = strings.TrimSpace(value)
	case "reason":
		builder.declaration.Reason = strings.TrimSpace(value)
	case "satisfied-by":
		builder.declaration.SatisfiedBy = strings.TrimSpace(value)
	default:
		return false
	}
	builder.activeField = field
	return true
}

func (builder *unreachableDeclarationBuilder) appendContinuation(line string) {
	appendValue := func(value *string) {
		if *value == "" {
			*value = line
			return
		}
		*value += " " + line
	}
	switch builder.activeField {
	case "criterion":
		appendValue(&builder.declaration.Criterion)
	case "reason":
		appendValue(&builder.declaration.Reason)
	case "satisfied-by":
		appendValue(&builder.declaration.SatisfiedBy)
	}
}

func (builder *unreachableDeclarationBuilder) finish(prdPath string) (UnreachableDeclaration, error) {
	missing := ""
	switch {
	case builder.declaration.Criterion == "":
		missing = "criterion"
	case builder.declaration.Reason == "":
		missing = "reason"
	case builder.declaration.SatisfiedBy == "":
		missing = "satisfied-by"
	}
	if missing != "" {
		return UnreachableDeclaration{}, UnreachableDeclarationError{
			Path:  prdPath,
			Line:  builder.declaration.Line,
			Field: missing,
		}
	}
	return builder.declaration, nil
}

// ListActive discovers the Specs eligible for the Implement Command:
// directories under the Spec Root (excluding _archived/) whose _prd.md
// frontmatter carries status active. Directories without a readable active
// PRD are skipped; Load names the exact problem when a slug is requested
// explicitly. The result is sorted by slug.
func ListActive(specsRoot string) ([]Spec, error) {
	specs, _, err := ListActiveDetailed(specsRoot)
	return specs, err
}

// ListActiveDetailed discovers active Specs and reports non-active Spec
// directories skipped because their PRD is missing, unreadable, or inactive.
// _archived/ is outside active discovery and is not reported as skipped.
func ListActiveDetailed(specsRoot string) ([]Spec, []SkippedSpec, error) {
	root := filepath.Clean(specsRoot)
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read Spec Root %q: %w", root, err)
	}
	var specs []Spec
	var skipped []SkippedSpec
	// os.ReadDir returns entries sorted by filename, so the result is
	// already sorted by slug.
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == archivedDirName {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		content, err := os.ReadFile(filepath.Join(dir, "_prd.md"))
		if err != nil {
			skipped = append(skipped, SkippedSpec{
				Dir:    specDisplayDir(root, entry.Name()),
				Reason: prdReadSkipReason(err),
			})
			continue
		}
		status, err := prdStatus(content)
		if err != nil {
			skipped = append(skipped, SkippedSpec{
				Dir:    specDisplayDir(root, entry.Name()),
				Reason: fmt.Sprintf("unreadable _prd.md frontmatter: %v", err),
			})
			continue
		}
		if status != prdStatusActive {
			skipped = append(skipped, SkippedSpec{
				Dir:    specDisplayDir(root, entry.Name()),
				Reason: fmt.Sprintf("status %q is not active", status),
			})
			continue
		}
		specs = append(specs, Spec{Slug: entry.Name(), Dir: dir})
	}
	return specs, skipped, nil
}

// Load parses and validates one Spec's Task Graph and task files from the
// Spec Root, returning the Tasks in deterministic topological order. Every
// validation failure is a typed error naming the offending Task or check.
func Load(specsRoot string, slug string) (*Graph, error) {
	root := filepath.Clean(specsRoot)
	dir := filepath.Join(root, slug)
	if err := requireActive(slug, dir); err != nil {
		return nil, err
	}

	manifestPath := filepath.Join(dir, "_tasks.md")
	nodes, projections, hasProjections, qa, err := loadManifestNodes(manifestPath)
	if err != nil {
		return nil, err
	}

	order, err := topologicalOrder(nodes)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, len(nodes))
	for _, index := range order {
		node := nodes[index]
		task, err := loadTask(dir, slug, node)
		if err != nil {
			return nil, err
		}
		if hasProjections {
			projected, ok := projections[task.ID]
			if !ok {
				return nil, ManifestError{
					Path:   manifestPath,
					Reason: fmt.Sprintf("projection table has no row for Task %q; add a type cell matching the task frontmatter", task.ID),
				}
			}
			if projected != task.Type {
				return nil, TaskTypeProjectionError{
					TaskID:       task.ID,
					ManifestPath: manifestPath,
					TaskPath:     filepath.Join(dir, node.File),
					ManifestType: projected,
					FileType:     task.Type,
				}
			}
		}
		tasks = append(tasks, task)
	}
	if err := validateQAGate(manifestPath, nodes, tasks, qa); err != nil {
		return nil, fmt.Errorf("validate qa gate: %w", err)
	}
	return &Graph{
		Spec:       Spec{Slug: slug, Dir: dir},
		Tasks:      tasks,
		QATaskID:   qa.TaskID,
		QADeclined: qa.Declined,
		QAReason:   qa.Reason,
	}, nil
}

func specDisplayDir(specsRoot string, slug string) string {
	if filepath.Base(specsRoot) == "specs" && filepath.Base(filepath.Dir(specsRoot)) == "docs" {
		return filepath.ToSlash(filepath.Join("docs", "specs", slug))
	}
	if filepath.IsAbs(specsRoot) {
		return filepath.ToSlash(filepath.Join(specsRoot, slug))
	}
	return filepath.ToSlash(filepath.Join("docs", "specs", slug))
}

func prdReadSkipReason(err error) string {
	if errors.Is(err, os.ErrNotExist) {
		return "missing _prd.md"
	}
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return fmt.Sprintf("read _prd.md: %v", pathErr.Err)
	}
	return fmt.Sprintf("read _prd.md: %v", err)
}

func requireActive(slug string, dir string) error {
	prdPath := filepath.Join(dir, "_prd.md")
	content, err := os.ReadFile(prdPath)
	if errors.Is(err, os.ErrNotExist) {
		return SpecNotFoundError{Slug: slug, Dir: dir}
	}
	if err != nil {
		return fmt.Errorf("read Spec PRD %q: %w", prdPath, err)
	}
	status, err := prdStatus(content)
	if err != nil {
		return fmt.Errorf("parse Spec PRD %q: %w", prdPath, err)
	}
	if status != prdStatusActive {
		return InactiveSpecError{Slug: slug, Status: status}
	}
	return nil
}

func prdStatus(content []byte) (string, error) {
	frontmatterBytes, _, err := splitFrontmatter(content)
	if err != nil {
		return "", err
	}
	var prd struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal(frontmatterBytes, &prd); err != nil {
		return "", err
	}
	return prd.Status, nil
}

func loadManifestNodes(manifestPath string) ([]manifestNode, map[string]TaskType, bool, qaDeclaration, error) {
	content, err := os.ReadFile(manifestPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, false, qaDeclaration{}, ManifestError{Path: manifestPath, Reason: "file does not exist; run the write-tasks workflow to create the Task Graph"}
	}
	if err != nil {
		return nil, nil, false, qaDeclaration{}, fmt.Errorf("read Task Graph manifest %q: %w", manifestPath, err)
	}
	frontmatterBytes, body, err := splitFrontmatter(content)
	if err != nil {
		return nil, nil, false, qaDeclaration{}, ManifestError{Path: manifestPath, Reason: "invalid frontmatter", Err: err}
	}
	var manifest manifestFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &manifest); err != nil {
		return nil, nil, false, qaDeclaration{}, ManifestError{Path: manifestPath, Reason: "invalid frontmatter", Err: err}
	}
	if manifest.Schema != manifestSchema {
		return nil, nil, false, qaDeclaration{}, ManifestSchemaError{Path: manifestPath, Schema: manifest.Schema}
	}
	qa, err := parseQADeclaration(manifestPath, manifest)
	if err != nil {
		return nil, nil, false, qaDeclaration{}, err
	}
	nodes := manifest.Graph.Nodes
	if len(nodes) == 0 {
		return nil, nil, false, qaDeclaration{}, ManifestError{Path: manifestPath, Reason: "graph has no nodes"}
	}

	known := make(map[string]bool, len(nodes))
	for index, node := range nodes {
		if node.ID == "" || node.File == "" {
			return nil, nil, false, qaDeclaration{}, ManifestError{Path: manifestPath, Reason: fmt.Sprintf("graph node %d is missing id or file", index+1)}
		}
		if known[node.ID] {
			return nil, nil, false, qaDeclaration{}, ManifestError{Path: manifestPath, Reason: fmt.Sprintf("duplicate Task id %q", node.ID)}
		}
		known[node.ID] = true
	}
	for _, node := range nodes {
		for _, need := range node.Needs {
			if !known[need] {
				return nil, nil, false, qaDeclaration{}, UnknownNeedError{TaskID: node.ID, Need: need}
			}
		}
	}
	projections, hasProjections, err := parseTaskTypeProjections(manifestPath, body)
	if err != nil {
		return nil, nil, false, qaDeclaration{}, err
	}
	for taskID := range projections {
		if !known[taskID] {
			return nil, nil, false, qaDeclaration{}, ManifestError{
				Path:   manifestPath,
				Reason: fmt.Sprintf("projection table row names unknown Task %q; remove it or add a matching graph node", taskID),
			}
		}
	}
	return nodes, projections, hasProjections, qa, nil
}

func parseQADeclaration(manifestPath string, manifest manifestFrontmatter) (qaDeclaration, error) {
	if !manifest.QA.Present {
		if manifest.QAReason.Present {
			return qaDeclaration{}, QAGateError{ManifestPath: manifestPath, Reason: "qa_reason requires a qa: declaration"}
		}
		return qaDeclaration{}, nil
	}

	qa := strings.TrimSpace(manifest.QA.Value)
	if qa == "" {
		return qaDeclaration{}, QAGateError{ManifestPath: manifestPath, Reason: "qa: must name a Task id or be declined"}
	}
	reason := strings.TrimSpace(manifest.QAReason.Value)
	if qa == "declined" {
		if !manifest.QAReason.Present || reason == "" {
			return qaDeclaration{}, QAGateError{ManifestPath: manifestPath, Reason: "qa: declined requires a non-empty qa_reason"}
		}
		return qaDeclaration{Present: true, Declined: true, Reason: reason}, nil
	}
	if manifest.QAReason.Present {
		return qaDeclaration{}, QAGateError{ManifestPath: manifestPath, Reason: "qa_reason is allowed only with qa: declined"}
	}
	return qaDeclaration{Present: true, TaskID: qa}, nil
}

func validateQAGate(manifestPath string, nodes []manifestNode, tasks []Task, qa qaDeclaration) error {
	taskByID := make(map[string]Task, len(tasks))
	var qaTaskIDs []string
	for _, task := range tasks {
		taskByID[task.ID] = task
		if task.Type == TaskTypeQA {
			qaTaskIDs = append(qaTaskIDs, task.ID)
		}
	}

	if !qa.Present {
		if len(qaTaskIDs) > 0 {
			return QAGateError{
				ManifestPath: manifestPath,
				Reason:       fmt.Sprintf("Task %q has type %q but the manifest has no qa: declaration", qaTaskIDs[0], TaskTypeQA),
			}
		}
		return nil
	}
	if qa.Declined {
		if len(qaTaskIDs) > 0 {
			return QAGateError{
				ManifestPath: manifestPath,
				Reason:       fmt.Sprintf("qa: declined cannot be combined with QA Task %q", qaTaskIDs[0]),
			}
		}
		return nil
	}

	gate, ok := taskByID[qa.TaskID]
	if !ok {
		return QAGateError{ManifestPath: manifestPath, Reason: fmt.Sprintf("qa: names Task %q, which is not a Task Graph node", qa.TaskID)}
	}
	if gate.Type != TaskTypeQA {
		return QAGateError{
			ManifestPath: manifestPath,
			Reason:       fmt.Sprintf("qa: names Task %q with type %q; expected %q", qa.TaskID, gate.Type, TaskTypeQA),
		}
	}
	if len(qaTaskIDs) != 1 {
		return QAGateError{
			ManifestPath: manifestPath,
			Reason:       fmt.Sprintf("qa: names Task %q, but QA Tasks must be unique; found: %s", qa.TaskID, strings.Join(qaTaskIDs, ", ")),
		}
	}

	nodeByID := make(map[string]manifestNode, len(nodes))
	var gateDependents []string
	for _, node := range nodes {
		nodeByID[node.ID] = node
		for _, need := range node.Needs {
			if need == qa.TaskID {
				gateDependents = append(gateDependents, node.ID)
			}
		}
	}
	if len(gateDependents) > 0 {
		return QAGateError{
			ManifestPath: manifestPath,
			Reason:       fmt.Sprintf("QA Task %q is not terminal; these Tasks depend on it: %s", qa.TaskID, strings.Join(gateDependents, ", ")),
		}
	}

	dependencyClosure := make(map[string]bool, len(nodes))
	stack := append([]string(nil), nodeByID[qa.TaskID].Needs...)
	for len(stack) > 0 {
		last := len(stack) - 1
		id := stack[last]
		stack = stack[:last]
		if dependencyClosure[id] {
			continue
		}
		dependencyClosure[id] = true
		stack = append(stack, nodeByID[id].Needs...)
	}

	hasNonGateDependent := make(map[string]bool, len(nodes))
	for _, node := range nodes {
		if node.ID == qa.TaskID {
			continue
		}
		for _, need := range node.Needs {
			hasNonGateDependent[need] = true
		}
	}
	var uncoveredLeaves []string
	for _, node := range nodes {
		if node.ID != qa.TaskID && !hasNonGateDependent[node.ID] && !dependencyClosure[node.ID] {
			uncoveredLeaves = append(uncoveredLeaves, node.ID)
		}
	}
	if len(uncoveredLeaves) > 0 {
		return QAGateError{
			ManifestPath: manifestPath,
			Reason:       fmt.Sprintf("QA Task %q does not depend on every leaf; uncovered Tasks: %s", qa.TaskID, strings.Join(uncoveredLeaves, ", ")),
		}
	}

	if gate.Status != StatusCompleted && gate.Status != StatusFailed {
		return nil
	}
	var staleDependencies []string
	for _, node := range nodes {
		if dependencyClosure[node.ID] && taskByID[node.ID].Status != StatusCompleted {
			staleDependencies = append(staleDependencies, node.ID)
		}
	}
	if len(staleDependencies) > 0 {
		return StaleGateError{QATaskID: qa.TaskID, TaskIDs: staleDependencies}
	}
	return nil
}

func parseTaskTypeProjections(manifestPath string, body []byte) (map[string]TaskType, bool, error) {
	projections := make(map[string]TaskType)
	inProjectionTable := false
	for _, line := range strings.Split(string(body), "\n") {
		cells := markdownTableCells(line)
		if !inProjectionTable {
			if taskTypeProjectionHeader(cells) {
				inProjectionTable = true
			}
			continue
		}
		if len(cells) == 0 {
			break
		}
		if markdownTableSeparator(cells) {
			continue
		}
		if len(cells) < 5 || !strings.HasPrefix(cells[0], "task_") {
			return nil, true, ManifestError{
				Path:   manifestPath,
				Reason: "projection table has a malformed Task row; each row must start with a canonical task_NN id and include title, type, complexity, and needs cells",
			}
		}
		if _, exists := projections[cells[0]]; exists {
			return nil, true, ManifestError{
				Path:   manifestPath,
				Reason: fmt.Sprintf("projection table defines Task %q more than once; remove duplicate rows", cells[0]),
			}
		}
		taskType, err := ParseTaskType(manifestPath, cells[2])
		if err != nil {
			return nil, true, ManifestError{
				Path:   manifestPath,
				Reason: fmt.Sprintf("projection row for Task %q has invalid type %q (allowed: %s); update the _tasks.md type cell to one allowed value", cells[0], cells[2], allowedTaskTypeValues()),
				Err:    err,
			}
		}
		projections[cells[0]] = taskType
	}
	return projections, inProjectionTable, nil
}

func taskTypeProjectionHeader(cells []string) bool {
	if len(cells) < 5 {
		return false
	}
	return strings.EqualFold(cells[0], "id") &&
		strings.EqualFold(cells[1], "title") &&
		strings.EqualFold(cells[2], "type") &&
		strings.EqualFold(cells[3], "complexity") &&
		strings.EqualFold(cells[4], "needs")
}

func markdownTableSeparator(cells []string) bool {
	if len(cells) == 0 {
		return false
	}
	for _, cell := range cells {
		trimmed := strings.TrimSpace(cell)
		if trimmed == "" {
			return false
		}
		for _, char := range trimmed {
			if char != '-' && char != ':' {
				return false
			}
		}
	}
	return true
}

func markdownTableCells(line string) []string {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
		return nil
	}
	trimmed = strings.TrimPrefix(strings.TrimSuffix(trimmed, "|"), "|")
	parts := strings.Split(trimmed, "|")
	cells := make([]string, 0, len(parts))
	for _, part := range parts {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells
}

// topologicalOrder runs Kahn's algorithm over the manifest nodes, always
// picking the first ready node in manifest order so the same manifest yields
// the same order on every run.
func topologicalOrder(nodes []manifestNode) ([]int, error) {
	indexByID := make(map[string]int, len(nodes))
	for index, node := range nodes {
		indexByID[node.ID] = index
	}
	indegree := make([]int, len(nodes))
	dependents := make([][]int, len(nodes))
	for index, node := range nodes {
		for _, need := range node.Needs {
			needIndex := indexByID[need]
			dependents[needIndex] = append(dependents[needIndex], index)
			indegree[index]++
		}
	}

	placed := make([]bool, len(nodes))
	order := make([]int, 0, len(nodes))
	for len(order) < len(nodes) {
		next := -1
		for index := range nodes {
			if !placed[index] && indegree[index] == 0 {
				next = index
				break
			}
		}
		if next < 0 {
			var remaining []string
			for index, node := range nodes {
				if !placed[index] {
					remaining = append(remaining, node.ID)
				}
			}
			return nil, CycleError{TaskIDs: remaining}
		}
		placed[next] = true
		order = append(order, next)
		for _, dependent := range dependents[next] {
			indegree[dependent]--
		}
	}
	return order, nil
}

func loadTask(dir string, slug string, node manifestNode) (Task, error) {
	path := filepath.Join(dir, node.File)
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Task{}, MissingTaskFileError{TaskID: node.ID, Path: path}
	}
	if err != nil {
		return Task{}, fmt.Errorf("read Task %q file %q: %w", node.ID, path, err)
	}
	document, err := parseTaskDocument(content, path)
	if err != nil {
		return Task{}, TaskFileError{TaskID: node.ID, Path: path, Err: err}
	}
	if len(document.Verification) == 0 {
		return Task{}, MissingVerificationError{TaskID: node.ID, Path: path}
	}
	return Task{
		ID:               node.ID,
		File:             filepath.Join(slug, node.File),
		Title:            document.Title,
		TitleLine:        document.TitleLine,
		Needs:            append([]string(nil), node.Needs...),
		Status:           Status(document.Frontmatter.Status),
		StatusNormalized: document.StatusNormalized,
		Type:             document.Type,
		Context:          append([]TaskContextRef(nil), document.Context...),
		Requirements:     append([]TaskDeclaration(nil), document.Requirements...),
		RehearsalCases:   append([]TaskDeclaration(nil), document.RehearsalCases...),
		Verification:     document.Verification,
		NegativeControl:  append([]string(nil), document.NegativeControl...),
	}, nil
}

// CarryForwardInputs returns the repository paths whose bytes declare the
// Task Roundfix handed to the Agent: the Task and Spec contracts, the root
// instructions, the canonical domain context, the implementation skill, and
// every Task-authored Context reference.
func CarryForwardInputs(specDir string, taskFile string, taskContent []byte) ([]string, error) {
	cleanSpecDir, err := cleanCarryForwardPath(specDir)
	if err != nil {
		return nil, fmt.Errorf("carry-forward Spec directory: %w", err)
	}
	cleanTaskFile, err := cleanCarryForwardPath(taskFile)
	if err != nil {
		return nil, fmt.Errorf("carry-forward Task file: %w", err)
	}
	document, err := parseTaskDocument(taskContent, cleanTaskFile)
	if err != nil {
		return nil, fmt.Errorf("parse carry-forward Task input %q: %w", cleanTaskFile, err)
	}

	inputs := []string{
		cleanTaskFile,
		filepath.ToSlash(filepath.Join(cleanSpecDir, "_prd.md")),
		filepath.ToSlash(filepath.Join(cleanSpecDir, "_techspec.md")),
		filepath.ToSlash(filepath.Join(cleanSpecDir, "_tasks.md")),
		"AGENTS.md",
		"CONTEXT.md",
		".agents/skills/implement-task/SKILL.md",
	}
	seen := make(map[string]bool, len(inputs)+len(document.Context))
	unique := make([]string, 0, len(inputs)+len(document.Context))
	for _, path := range inputs {
		if !seen[path] {
			seen[path] = true
			unique = append(unique, path)
		}
	}
	for _, ref := range document.Context {
		clean, cleanErr := cleanCarryForwardPath(ref.Path)
		if cleanErr != nil {
			return nil, fmt.Errorf("carry-forward Context path %q: %w", ref.Path, cleanErr)
		}
		if !seen[clean] {
			seen[clean] = true
			unique = append(unique, clean)
		}
	}
	return unique, nil
}

// CarryForwardStatus reads the Task status from committed Task bytes.
func CarryForwardStatus(taskFile string, content []byte) (Status, error) {
	document, err := parseTaskDocument(content, taskFile)
	if err != nil {
		return "", fmt.Errorf("parse carry-forward Task status %q: %w", taskFile, err)
	}
	return Status(document.Frontmatter.Status), nil
}

// RecordCarryForward marks a Task completed and appends the source Run and
// settlement commit to the Task file. Every pre-existing byte other than the
// status value is preserved exactly.
func RecordCarryForward(taskPath string, runID string, commit string) error {
	runID = strings.TrimSpace(runID)
	commit = strings.TrimSpace(commit)
	if err := validateCarryForwardRecordValue("Run ID", runID); err != nil {
		return err
	}
	if err := validateCarryForwardRecordValue("commit", commit); err != nil {
		return err
	}
	info, err := os.Stat(taskPath)
	if err != nil {
		return fmt.Errorf("stat carry-forward Task file %q: %w", taskPath, err)
	}
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return fmt.Errorf("read carry-forward Task file %q: %w", taskPath, err)
	}
	const heading = "## Carry-forward provenance"
	if strings.Contains(string(content), "\n"+heading+"\n") {
		return fmt.Errorf("carry-forward Task file %q already records provenance", taskPath)
	}
	updated, err := rewriteStatus(content, StatusCompleted)
	if err != nil {
		return fmt.Errorf("rewrite carry-forward status in Task file %q: %w", taskPath, err)
	}
	updated = append(updated, []byte(fmt.Sprintf(
		"\n%s\n\n- Source Run: `%s`\n- Source commit: `%s`\n",
		heading,
		runID,
		commit,
	))...)
	if err := os.WriteFile(taskPath, updated, info.Mode().Perm()); err != nil {
		return fmt.Errorf("write carry-forward Task file %q: %w", taskPath, err)
	}
	return nil
}

func cleanCarryForwardPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("path is empty")
	}
	clean := filepath.Clean(path)
	if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("path must stay repository-relative")
	}
	return filepath.ToSlash(clean), nil
}

func validateCarryForwardRecordValue(label string, value string) error {
	if value == "" {
		return fmt.Errorf("carry-forward %s is required", label)
	}
	if strings.ContainsAny(value, "\r\n`") {
		return fmt.Errorf("carry-forward %s contains an unsupported character", label)
	}
	return nil
}

// splitFrontmatter mirrors the Round artifact frontmatter contract: an
// opening --- line, YAML, and a closing --- line ahead of the markdown body.
func splitFrontmatter(content []byte) ([]byte, []byte, error) {
	text := string(content)
	const opening = "---\n"
	if !strings.HasPrefix(text, opening) {
		return nil, nil, errors.New("missing YAML frontmatter opening marker")
	}
	rest := text[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return nil, nil, errors.New("missing YAML frontmatter closing marker")
	}
	bodyStart := end + len("\n---")
	if strings.HasPrefix(rest[bodyStart:], "\n") {
		bodyStart++
	}
	return []byte(rest[:end]), []byte(rest[bodyStart:]), nil
}
