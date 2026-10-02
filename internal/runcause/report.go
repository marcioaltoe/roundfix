package runcause

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

// Reader is the report's read-only Run Database boundary.
type Reader interface {
	ListRuns(context.Context, store.ListRunsQuery) ([]store.Run, error)
	RunEventsOfKinds(context.Context, string, ...runevent.Kind) ([]store.JournalEvent, error)
}

type Request struct {
	Reader                 Reader
	RepositoryRoot         string
	SpecsRoot, ArchiveRoot string
	Since, Until           time.Time
}

type Window struct {
	Since *string `json:"since"`
	Until *string `json:"until"`
}

type Item struct {
	Kind      string  `json:"kind"`
	Spec      string  `json:"spec"`
	Task      string  `json:"task"`
	RunID     string  `json:"run_id"`
	Attempt   *int    `json:"attempt"`
	Check     string  `json:"check"`
	Class     string  `json:"class"`
	Signature *string `json:"signature"`
	Evidence  *string `json:"evidence"`
}

type TaskCounts struct {
	Spec           string `json:"spec"`
	Task           string `json:"task"`
	QA             bool   `json:"qa"`
	Corrective     *bool  `json:"corrective"`
	VerdictsPassed int    `json:"verdicts_passed"`
	VerdictsFailed int    `json:"verdicts_failed"`
	FeedbackRounds int    `json:"feedback_rounds"`
	Runs           int    `json:"runs"`
}

type Report struct {
	Schema           string         `json:"schema"`
	Window           Window         `json:"window"`
	SignaturesSHA256 string         `json:"signatures_sha256"`
	Runs             int            `json:"runs"`
	Items            []Item         `json:"items"`
	Tasks            []TaskCounts   `json:"tasks"`
	Summary          map[string]int `json:"summary"`
	SpecsNotFound    []string       `json:"specs_not_found"`
}

func date(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.UTC().Format(time.DateOnly)
	return &s
}

// Summarize recomputes every count from items; unclassified is never a class.
func (t Table) Summarize(items []Item) map[string]int {
	summary := map[string]int{"unclassified": 0, "items": len(items), "classified": 0, "repository_knowledge": 0}
	for _, class := range t.Classes {
		summary[class] = 0
	}
	for _, item := range items {
		summary[item.Class]++
		if slices.Contains(t.Classes, item.Class) {
			summary["classified"]++
		}
		if slices.Contains(t.RepositoryKnowledge, item.Class) {
			summary["repository_knowledge"]++
		}
	}
	return summary
}

// Build joins terminal Spec Runs, failed attempts and historical graph nodes.
// A nil Reader means an absent database, producing the same empty document.
func Build(ctx context.Context, t Table, req Request) (Report, error) {
	report := Report{Schema: "roundfix/runs-causes/v1", Window: Window{date(req.Since), date(req.Until)}, SignaturesSHA256: t.SHA256, Items: []Item{}, Tasks: []TaskCounts{}, SpecsNotFound: []string{}}
	report.Summary = t.Summarize(report.Items)
	if req.Reader == nil {
		return report, nil
	}
	runs, err := req.Reader.ListRuns(ctx, store.ListRunsQuery{RepositoryRoot: req.RepositoryRoot, States: store.StatesTerminal})
	if err != nil {
		return Report{}, fmt.Errorf("list cause Runs: %w", err)
	}
	sort.Slice(runs, func(i, j int) bool {
		if runs[i].CreatedAt.Equal(runs[j].CreatedAt) {
			return runs[i].ID < runs[j].ID
		}
		return runs[i].CreatedAt.Before(runs[j].CreatedAt)
	})
	tasks := map[string]*TaskCounts{}
	graphs := map[string]*spec.CauseGraph{}
	correctiveSeen := map[string]bool{}
	for _, run := range runs {
		if run.SpecSlug == "" || (!req.Since.IsZero() && run.CreatedAt.Before(req.Since)) || (!req.Until.IsZero() && !run.CreatedAt.Before(req.Until)) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return Report{}, err
		}
		report.Runs++
		graph, cached := graphs[run.SpecSlug]
		if !cached {
			graph, err = findGraph(req, run.SpecSlug)
			if err != nil {
				return Report{}, err
			}
			graphs[run.SpecSlug] = graph
			if graph == nil {
				report.SpecsNotFound = append(report.SpecsNotFound, run.SpecSlug)
			}
		}
		events, err := req.Reader.RunEventsOfKinds(ctx, run.ID, runevent.KindDaemonVerification, runevent.KindDaemonTask)
		if err != nil {
			return Report{}, fmt.Errorf("read cause events for %q: %w", run.ID, err)
		}
		payloads := make([]eventPayload, len(events))
		started := map[string]bool{}
		for i, entry := range events {
			if len(entry.Event.Payload) > 0 {
				if err := json.Unmarshal(entry.Event.Payload, &payloads[i]); err != nil {
					return Report{}, fmt.Errorf("decode cause event %q cursor %d: %w", run.ID, entry.Cursor, err)
				}
			}
			if payloads[i].Task == "" {
				payloads[i].Task = entry.Event.ReviewIssue
			}
			p := payloads[i]
			if entry.Event.Kind == runevent.KindDaemonTask && p.Phase == "started" && p.Task != "" {
				started[p.Task] = true
			}
		}
		for id := range started {
			key := run.SpecSlug + "/" + id
			if tasks[key] == nil {
				counts := &TaskCounts{Spec: run.SpecSlug, Task: id}
				if graph != nil {
					corrective := isCorrective(*graph, id)
					counts.Corrective = &corrective
					counts.QA = id == graph.QATaskID
				}
				tasks[key] = counts
			}
			tasks[key].Runs++
		}
		seenAttempts := map[string]bool{}
		for i, entry := range events {
			p := payloads[i]
			if started[p.Task] {
				counts := tasks[run.SpecSlug+"/"+p.Task]
				if entry.Event.Kind == runevent.KindDaemonVerification {
					if p.Verdict == "passed" {
						counts.VerdictsPassed++
					}
					if p.Verdict == "failed" {
						counts.VerdictsFailed++
					}
				}
				if entry.Event.Kind == runevent.KindDaemonTask && p.Phase == "verification_feedback" {
					counts.FeedbackRounds++
				}
			}
			if entry.Event.Kind != runevent.KindDaemonVerification || (p.Phase != "failed" && p.Verdict != "failed") {
				continue
			}
			key := p.Task + "/" + strconv.Itoa(p.Attempt)
			if seenAttempts[key] {
				continue
			}
			seenAttempts[key] = true
			check := p.Command
			if check == "" {
				check = strings.Join(p.Commands, "\n")
			}
			diagnostic, err := diagnosticTail(run, p.DiagnosticPath, t.DiagnosticTailBytes)
			if err != nil {
				return Report{}, err
			}
			item := Item{Kind: "verification", Spec: run.SpecSlug, Task: p.Task, RunID: run.ID, Attempt: &p.Attempt, Check: publicCheck(check)}
			item.Class, item.Signature, item.Evidence = classifyItem(t, Evidence{Diagnostic: diagnostic, Command: check})
			report.Items = append(report.Items, item)
		}
		ids := make([]string, 0, len(started))
		for id := range started {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool {
			if taskNumber(ids[i]) == taskNumber(ids[j]) {
				return ids[i] < ids[j]
			}
			return taskNumber(ids[i]) < taskNumber(ids[j])
		})
		for _, id := range ids {
			key := run.SpecSlug + "/" + id
			if graph == nil || !isCorrective(*graph, id) || correctiveSeen[key] {
				continue
			}
			correctiveSeen[key] = true
			text, err := graph.CauseTaskText(id, t.TaskTextBytes)
			if err != nil {
				return Report{}, err
			}
			item := Item{Kind: "corrective", Spec: run.SpecSlug, Task: id, RunID: run.ID, Check: t.Trigger(text)}
			item.Class, item.Signature, item.Evidence = classifyItem(t, Evidence{Task: text})
			report.Items = append(report.Items, item)
		}
	}
	keys := make([]string, 0, len(tasks))
	for key := range tasks {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		report.Tasks = append(report.Tasks, *tasks[key])
	}
	sort.Strings(report.SpecsNotFound)
	report.Summary = t.Summarize(report.Items)
	return report, nil
}

type eventPayload struct {
	Task           string   `json:"task"`
	Phase          string   `json:"phase"`
	Verdict        string   `json:"verdict"`
	Attempt        int      `json:"attempt"`
	Command        string   `json:"command"`
	Commands       []string `json:"commands"`
	DiagnosticPath string   `json:"diagnostic_path"`
}

func classifyItem(t Table, e Evidence) (string, *string, *string) {
	class, signature, source := t.Classify(e)
	if signature == "" {
		return class, nil, nil
	}
	return class, &signature, &source
}

func findGraph(req Request, slug string) (*spec.CauseGraph, error) {
	for _, root := range []string{req.SpecsRoot, req.ArchiveRoot} {
		if root == "" {
			continue
		}
		graph, err := spec.ReadCauseGraph(root, slug)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read cause graph %q: %w", slug, err)
		}
		return &graph, nil
	}
	return nil, nil
}

func taskNumber(id string) int {
	if !strings.HasPrefix(id, "task_") {
		return -1
	}
	n, err := strconv.Atoi(strings.TrimPrefix(id, "task_"))
	if err != nil {
		return -1
	}
	return n
}

func isCorrective(g spec.CauseGraph, id string) bool {
	_, member := g.Tasks[id]
	return member && g.QATaskID != "" && taskNumber(g.QATaskID) >= 0 && taskNumber(id) > taskNumber(g.QATaskID)
}

// diagnosticTail checks every path component with Lstat. Neither a leaf link
// nor a linked parent can turn external content into diagnostic evidence.
func diagnosticTail(run store.Run, path string, limit int64) (string, error) {
	if path == "" || run.ArtifactDir == "" || filepath.Base(run.ID) != run.ID {
		return "", nil
	}
	root := filepath.Join(run.ArtifactDir, "runs", run.ID)
	if !filepath.IsAbs(path) {
		return "", nil
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", nil
	}
	clean := filepath.Clean(path)
	current := string(filepath.Separator)
	for _, part := range strings.Split(strings.TrimPrefix(clean, string(filepath.Separator)), string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			return "", nil
		}
		if current == clean && !info.Mode().IsRegular() {
			return "", nil
		}
	}
	file, err := os.Open(clean)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("open Verification diagnostic: %w", err)
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("stat Verification diagnostic: %w", err)
	}
	if _, err = file.Seek(max(0, info.Size()-limit), io.SeekStart); err != nil {
		return "", fmt.Errorf("seek Verification diagnostic tail: %w", err)
	}
	data, err := io.ReadAll(io.LimitReader(file, limit))
	if err != nil {
		return "", fmt.Errorf("read Verification diagnostic tail: %w", err)
	}
	return string(data), nil
}

// The check remains useful without exporting absolute paths or credentials
// embedded in a Verification command. Classification uses the original text.
var checkCredentialAssignment = regexp.MustCompile(`(?i)([a-z0-9_]*(?:key|token|password|secret)\s*=\s*)(?:"[^"]*"|'[^']*'|[^\s;]+)`)
var checkCredentialFlag = regexp.MustCompile(`(?i)(--(?:api-key|token|password|secret)(?:=|\s+))(?:"[^"]*"|'[^']*'|[^\s;]+)`)
var checkAbsolutePath = regexp.MustCompile(`(^|[\s="'(])(/[^\s"'<>;]*)`)

func publicCheck(command string) string {
	command = checkCredentialAssignment.ReplaceAllString(command, "${1}[redacted]")
	command = checkCredentialFlag.ReplaceAllString(command, "${1}[redacted]")
	return checkAbsolutePath.ReplaceAllString(command, "${1}[path]")
}
