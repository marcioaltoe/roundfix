package judge

// Suite: test-only acceptance measurement. Boundary: local Task files and
// Judge Log are real; external HTTP is injected. Only the explicit live flag
// permits the process environment and default network transport.
import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/config"
	"roundfix/internal/runcause"
	"roundfix/internal/spec"
)

var (
	measureAcceptance = flag.Bool("measure-task-acceptance", false, "send the acceptance measurement requests")
	measureLabels     = flag.String("measure-labels", "", "cause record used as labels")
	measureRepo       = flag.String("measure-repo", "", "repository root (default: this repository)")
	measureOut        = flag.String("measure-out", "", "output acceptance record")
)

type acceptanceQuestion struct {
	Schema      string   `json:"schema"`
	Judgment    string   `json:"judgment"`
	QuestionID  string   `json:"question_id"`
	Question    Question `json:"question"`
	TitleMax    int      `json:"task_title_max_chars"`
	CriteriaMax int      `json:"acceptance_criteria_max_chars"`
	FlagBelow   float64  `json:"flag_when_noul_below"`
	Bootstrap   struct {
		Resamples    int     `json:"resamples"`
		Seed         uint64  `json:"seed"`
		Cluster      string  `json:"cluster"`
		MaxDiscarded float64 `json:"max_discarded_share"`
	} `json:"bootstrap"`
	RandomSeed uint64 `json:"random_baseline_seed"`
}

func loadAcceptanceQuestion() (acceptanceQuestion, []byte, error) {
	data, err := os.ReadFile("testdata/task-acceptance-question.json")
	if err != nil {
		return acceptanceQuestion{}, nil, err
	}
	var q acceptanceQuestion
	if err := json.Unmarshal(data, &q); err != nil {
		return q, nil, err
	}
	if q.Schema != "roundfix/task-acceptance-question/v1" || q.Judgment != "task-acceptance" || q.Question.Type != "noul" || q.QuestionID == "" || q.TitleMax <= 0 || q.CriteriaMax <= 0 || q.Bootstrap.Resamples <= 0 || q.Bootstrap.Cluster != "spec" {
		return q, nil, errors.New("invalid acceptance question")
	}
	return q, data, nil
}
func acceptanceHash(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }

type acceptanceState struct {
	Title    string `json:"task_title"`
	Criteria string `json:"acceptance_criteria"`
}
type acceptanceTask struct {
	Spec          string   `json:"spec"`
	Task          string   `json:"task"`
	Repaired      bool     `json:"repaired"`
	Failed        int      `json:"verdicts_failed"`
	CriteriaChars int      `json:"criteria_chars"`
	StateHash     string   `json:"state_sha256"`
	Noul          *float64 `json:"noul"`
	Model         string   `json:"model"`
	Outcome       string   `json:"outcome"`
	Reason        string   `json:"reason"`
}
type acceptanceExcluded struct {
	Spec   string `json:"spec"`
	Task   string `json:"task"`
	Reason string `json:"reason"`
}
type acceptanceRecord struct {
	Schema       string               `json:"schema"`
	Status       string               `json:"status"`
	Reason       string               `json:"reason"`
	LabelsHash   string               `json:"labels_sha256"`
	QuestionHash string               `json:"question_sha256"`
	PinnedModel  string               `json:"pinned_model"`
	Transport    string               `json:"transport"`
	Calls        int                  `json:"calls"`
	InputTokens  int                  `json:"input_tokens"`
	Cost         float64              `json:"cost_usd"`
	Tasks        []acceptanceTask     `json:"tasks"`
	Excluded     []acceptanceExcluded `json:"excluded"`
	Statistics   acceptanceStatistics `json:"statistics"`
	Verdict      string               `json:"verdict"`
}

func acceptanceRoots(repo, home string) ([]string, error) {
	loaded, err := config.Load(config.LoadOptions{WorkDir: repo, HomeDir: home, Stderr: io.Discard})
	if err != nil {
		return nil, fmt.Errorf("load Spec Root: %w", err)
	}
	root, err := config.ResolveSpecsRoot(loaded, repo)
	if err != nil {
		return nil, err
	}
	return []string{root.Path, spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot)}, nil
}

var acceptanceTaskName = regexp.MustCompile(`^task_[0-9]{2}$`)

// Check every component, including the configured root, before opening a
// direct Task child. A symlink is never a missing file to fall back around.
func acceptanceTaskFile(roots []string, slug, task string) (string, error) {
	if slug == "" || slug == "." || slug == ".." || strings.ContainsAny(slug, `/\`) || !acceptanceTaskName.MatchString(task) {
		return "", errors.New("invalid Spec or Task name")
	}
	for _, root := range roots {
		path := filepath.Join(root, slug, task+".md")
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		// Inspect from the filesystem root so a symlink above Spec Root is refused too.
		current := string(filepath.Separator)
		missing := false
		for _, part := range strings.Split(strings.TrimPrefix(abs, string(filepath.Separator)), string(filepath.Separator)) {
			current = filepath.Join(current, part)
			info, err := os.Lstat(current)
			if errors.Is(err, os.ErrNotExist) {
				missing = true
				break
			}
			if err != nil {
				return "", err
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return "", errors.New("symbolic-link Task path")
			}
			if current == abs {
				if !info.Mode().IsRegular() {
					return "", errors.New("Task is not a regular file")
				}
			} else if !info.IsDir() {
				return "", errors.New("Task parent is not a directory")
			}
		}
		if !missing {
			return abs, nil
		}
	}
	return "", errors.New("Task file missing from active and archived Specs")
}
func acceptanceCut(s string, limit int) string {
	r := []rune(s)
	if len(r) > limit {
		r = r[:limit]
	}
	return string(r)
}
func buildAcceptanceState(roots []string, slug, task string, q acceptanceQuestion) (acceptanceState, []byte, error) {
	path, err := acceptanceTaskFile(roots, slug, task)
	if err != nil {
		return acceptanceState{}, nil, err
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return acceptanceState{}, nil, errors.New("Task changed before opening")
	}
	f, err := os.Open(path)
	if err != nil {
		return acceptanceState{}, nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return acceptanceState{}, nil, errors.New("Task changed while opening")
	}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	var state acceptanceState
	var criteria strings.Builder
	found, inCriteria, titleFound := false, false, false
	fence := ""
	lineNumber, frontMatter := 0, false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		lineNumber++
		if lineNumber == 1 && line == "---" {
			frontMatter = true
			continue
		}
		if frontMatter {
			if line == "---" {
				frontMatter = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			marker := trimmed[:3]
			if fence == "" {
				fence = marker
			} else if fence == marker {
				fence = ""
			}
		}
		if fence == "" {
			if inCriteria && (strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "# ")) {
				break
			}
			if !titleFound && strings.HasPrefix(line, "# ") {
				titleFound = true
				state.Title = acceptanceCut(strings.TrimSpace(strings.TrimPrefix(line, "# ")), q.TitleMax)
			}
			if line == "## Acceptance Criteria" {
				found, inCriteria = true, true
				continue
			}
		}
		if inCriteria {
			criteria.WriteString(line)
			criteria.WriteByte('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		return state, nil, err
	}
	if !found {
		return state, nil, errors.New("no Acceptance Criteria section")
	}
	state.Criteria = acceptanceCut(strings.TrimSpace(criteria.String()), q.CriteriaMax)
	data, err := json.Marshal(state)
	return state, data, err
}

// The delivered client already accepts a PendingJudgment and Questions.
// Override its in-memory question slot only; the embedded judge is unchanged.
func measureTaskAcceptance(ctx context.Context, repo string, labels []byte, keys map[string]string, home string, transport http.RoundTripper, now func() time.Time) (acceptanceRecord, error) {
	aq, questionBytes, err := loadAcceptanceQuestion()
	if err != nil {
		return acceptanceRecord{}, err
	}
	q, err := Load()
	if err != nil {
		return acceptanceRecord{}, err
	}
	q.Citation.QuestionID, q.Citation.Question = aq.QuestionID, aq.Question
	record := acceptanceRecord{Schema: "roundfix/task-acceptance-remeasurement/v1", Status: "measured", LabelsHash: acceptanceHash(labels), QuestionHash: acceptanceHash(questionBytes), PinnedModel: q.PinnedModel, Tasks: []acceptanceTask{}, Excluded: []acceptanceExcluded{}}
	finish := func() acceptanceRecord {
		record.Statistics = computeAcceptanceStatistics(record.Tasks, aq)
		record.Verdict = acceptanceVerdict(record.Status, record.Statistics, aq)
		return record
	}
	// Only Roundfix-scoped names enter transport selection or redaction.
	scoped := map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": keys["ROUNDFIX_OPENROUTER_API_KEY"], "ROUNDFIX_TYPESAFE_API_KEY": keys["ROUNDFIX_TYPESAFE_API_KEY"]}
	recipient, key, ok := selectTransport(q, scoped)
	if !ok {
		record.Status, record.Reason = "blocked", "neither Roundfix transport key is set"
		return finish(), nil
	}
	if transport == nil {
		return record, errors.New("measurement transport must be explicit")
	}
	if now == nil {
		return record, errors.New("measurement clock must be explicit")
	}
	record.Transport = recipient.Name
	var cause runcause.Report
	if err := json.Unmarshal(labels, &cause); err != nil {
		return record, fmt.Errorf("read labels: %w", err)
	}
	if cause.Schema != "roundfix/runs-causes/v1" {
		return record, errors.New("invalid cause record schema")
	}
	roots, err := acceptanceRoots(repo, home)
	if err != nil {
		return record, err
	}
	redact := func(text string) string {
		for _, secret := range scoped {
			if secret != "" {
				text = strings.ReplaceAll(text, secret, "[redacted]")
			}
		}
		return text
	}
	c := client{q: q, transport: recipient, key: key, http: &http.Client{Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	stopped := ""
	for _, label := range cause.Tasks {
		if label.QA || label.VerdictsPassed+label.VerdictsFailed == 0 {
			continue
		}
		state, data, err := buildAcceptanceState(roots, label.Spec, label.Task, aq)
		if err == nil && !q.Language.isEnglish(state.Title+"\n"+state.Criteria) {
			err = errors.New("state fails judge language gate")
		}
		if err != nil {
			record.Excluded = append(record.Excluded, acceptanceExcluded{label.Spec, label.Task, redact(err.Error())})
			continue
		}
		task := acceptanceTask{Spec: label.Spec, Task: label.Task, Repaired: label.VerdictsFailed >= 1, Failed: label.VerdictsFailed, CriteriaChars: utf8.RuneCountInString(state.Criteria), StateHash: acceptanceHash(data), Outcome: "skipped"}
		for attempt := 1; attempt <= 3; attempt++ {
			if stopped != "" {
				task.Reason = stopped
				break
			}
			callTime := now().UTC()
			log := newJudgeLog(home, callTime)
			monthCost, err := log.monthCost()
			if err != nil {
				stopped = "judge log unreadable: " + redact(err.Error())
				task.Reason = stopped
				break
			}
			if monthCost >= q.MonthlyCeilingUSD {
				stopped = "monthly ceiling reached"
				task.Reason = stopped
				break
			}
			if err := ctx.Err(); err != nil {
				stopped = err.Error()
				task.Reason = stopped
				break
			}
			pending := PendingJudgment{Kind: aq.Judgment, state: data}
			result := c.ask(ctx, pending, attempt)
			result.Error, result.Model, result.ResponseID, result.Provider = redact(result.Error), redact(result.Model), redact(result.ResponseID), redact(result.Provider)
			a, outcome, reason, stop := evaluate(q, pending, result)
			retry := (result.Status == 429 || result.Status == 529) && attempt < 3 && result.Error == ""
			if retry {
				stop = false
			}
			task.Model, task.Reason = result.Model, redact(reason)
			if outcome != "skipped" {
				task.Outcome, task.Noul = "answered", a.Noul
			} else {
				task.Outcome, task.Noul = "skipped", nil
			}
			cost, source := result.cost(q)
			record.Calls++
			record.InputTokens += result.Usage.InputTokens
			record.Cost += cost
			row := logLine{Schema: "roundfix/judge-log/v1", Time: callTime, Repository: redact(repo), Spec: redact(label.Spec), Judgment: aq.Judgment, Artifact: redact(label.Task + ".md"), StateHash: task.StateHash[:16], QuestionID: aq.QuestionID, Transport: recipient.Name, RequestedModel: recipient.RequestModel, Model: result.Model, ResponseID: result.ResponseID, Provider: result.Provider, Noul: a.Noul, LatencyMS: result.LatencyMS, InputTokens: result.Usage.InputTokens, OutputTokens: result.Usage.OutputTokens, CostUSD: cost, CostSource: source, Status: result.Status, Attempts: attempt, Error: task.Reason, Outcome: task.Outcome}
			if err := log.append(row); err != nil {
				stopped = "judge log not writable: " + redact(err.Error())
				break
			}
			if stop {
				stopped = task.Reason
				break
			}
			if !retry {
				break
			}
			if err := waitRetry(ctx, retryDelay(result.RetryAfter, attempt)); err != nil {
				stopped = err.Error()
				break
			}
		}
		record.Tasks = append(record.Tasks, task)
	}
	if stopped != "" {
		record.Status, record.Reason = "blocked", stopped
	}
	return finish(), nil
}
func writeAcceptanceRecord(path string, record acceptanceRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}
func TestMeasureTaskAcceptance(t *testing.T) {
	if !*measureAcceptance {
		t.Skip("requires -measure-task-acceptance; no keys, files or home read")
	}
	if *measureLabels == "" || *measureOut == "" {
		t.Fatal("-measure-labels and -measure-out are required")
	}
	repo := *measureRepo
	if repo == "" {
		var err error
		repo, err = filepath.Abs("../..")
		if err != nil {
			t.Fatal(err)
		}
	}
	labels, err := os.ReadFile(*measureLabels)
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": os.Getenv("ROUNDFIX_OPENROUTER_API_KEY"), "ROUNDFIX_TYPESAFE_API_KEY": os.Getenv("ROUNDFIX_TYPESAFE_API_KEY")}
	home := os.Getenv("HOME")
	if home == "" {
		t.Fatal("HOME is required")
	}
	record, err := measureTaskAcceptance(t.Context(), repo, labels, keys, home, http.DefaultTransport, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeAcceptanceRecord(*measureOut, record); err != nil {
		t.Fatal(err)
	}
}

const acceptanceFixtureTitle = "Task 01: The command prints the expected output"
const acceptanceFixtureCriteria = "- [ ] The command exits with code 0 and prints the expected output.\n- [ ] The test reads the file and checks the named field."
const acceptanceResultSentinel = "RESULT_PRIVATE_SENTINEL_0214"
const acceptanceSourceSentinel = "SOURCE_PRIVATE_SENTINEL_0214"

func acceptanceFixtureText() string {
	return "---\ntask: task_01\nstatus: pending\n---\n\n# " + acceptanceFixtureTitle + "\n\n## Overview\nThe implementation follows the command contract.\n\n## Acceptance Criteria\n\n" + acceptanceFixtureCriteria + "\n\n## Result\n" + acceptanceResultSentinel + "\n"
}
func acceptanceWrite(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
}
func acceptanceTestRepo(t *testing.T) (string, string, []byte) {
	t.Helper()
	repo, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	acceptanceWrite(t, filepath.Join(repo, "docs/specs/0300-example/task_01.md"), acceptanceFixtureText())
	acceptanceWrite(t, filepath.Join(repo, "main.go"), "package example\n// "+acceptanceSourceSentinel+"\n")
	labels, err := json.Marshal(runcause.Report{Schema: "roundfix/runs-causes/v1", Tasks: []runcause.TaskCounts{{Spec: "0300-example", Task: "task_01", VerdictsPassed: 1, VerdictsFailed: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	return repo, home, labels
}
func acceptanceClock() time.Time { return time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC) }
func acceptanceResponse(model string, cost float64) *http.Response {
	return fixtureResponse(200, fmt.Sprintf(`{"model":%q,"id":"fixture","provider":"TypeSafe","answers":{"acceptance_observable":{"type":"noul","noul":0.2}},"usage":{"input_tokens":100,"output_tokens":1,"cost":%v}}`, model, cost))
}
func acceptanceFake() http.RoundTripper {
	return roundTripFunc(func(*http.Request) (*http.Response, error) { return acceptanceResponse("jev-1.13.0", 0.001), nil })
}
func acceptanceMeasure(t *testing.T, repo, home string, labels []byte, keys map[string]string, transport http.RoundTripper) acceptanceRecord {
	t.Helper()
	record, err := measureTaskAcceptance(t.Context(), repo, labels, keys, home, transport, acceptanceClock)
	if err != nil {
		t.Fatal(err)
	}
	return record
}
func TestTaskAcceptanceQuestionFileLoads(t *testing.T) {
	q, _, err := loadAcceptanceQuestion()
	if err != nil {
		t.Fatal(err)
	}
	if q.Question.Criteria["true"] == "" || q.Question.Criteria["false"] == "" || q.FlagBelow <= 0 || q.FlagBelow >= 1 {
		t.Fatalf("invalid question: %+v", q)
	}
}
func TestTaskAcceptanceStateCarriesOnlyTitleAndCriteria(t *testing.T) {
	repo, _, _ := acceptanceTestRepo(t)
	q, _, err := loadAcceptanceQuestion()
	if err != nil {
		t.Fatal(err)
	}
	roots := []string{filepath.Join(repo, "docs/specs"), filepath.Join(repo, "docs/history/specs")}
	state, data, err := buildAcceptanceState(roots, "0300-example", "task_01", q)
	if err != nil {
		t.Fatal(err)
	}
	if state.Title != acceptanceFixtureTitle || state.Criteria != acceptanceFixtureCriteria {
		t.Fatalf("state=%+v", state)
	}
	var fields map[string]string
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 || fields["task_title"] != state.Title || fields["acceptance_criteria"] != state.Criteria {
		t.Fatalf("state fields=%v", fields)
	}
	t.Run("rune limits and section boundaries", func(t *testing.T) {
		title := strings.Repeat("é", q.TitleMax+10)
		criteria := strings.Repeat("界", q.CriteriaMax+10)
		acceptanceWrite(t, filepath.Join(repo, "docs/specs/0300-example/task_01.md"), "# "+title+"\n# Second title\n## Acceptance Criteria\n"+criteria+"\n## Result\n"+acceptanceResultSentinel)
		state, _, err := buildAcceptanceState(roots, "0300-example", "task_01", q)
		if err != nil {
			t.Fatal(err)
		}
		if utf8.RuneCountInString(state.Title) != q.TitleMax || utf8.RuneCountInString(state.Criteria) != q.CriteriaMax || strings.Contains(state.Criteria, acceptanceResultSentinel) {
			t.Fatalf("incorrect limits or boundary: %+v", state)
		}
	})
}
func TestTaskAcceptanceRequestsCarryNoResultOrSourceText(t *testing.T) {
	repo, home, labels := acceptanceTestRepo(t)
	requests := 0
	record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "openrouter-secret"}, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		for _, sentinel := range []string{acceptanceResultSentinel, acceptanceSourceSentinel, "openrouter-secret"} {
			if strings.Contains(string(data), sentinel) {
				t.Fatalf("request leaked %s", sentinel)
			}
		}
		var body struct {
			State     map[string]string   `json:"state"`
			Questions map[string]Question `json:"questions"`
		}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		if len(body.State) != 2 || body.State["task_title"] != acceptanceFixtureTitle || body.State["acceptance_criteria"] != acceptanceFixtureCriteria || len(body.Questions) != 1 || body.Questions["acceptance_observable"].Type != "noul" {
			t.Fatalf("body=%s", data)
		}
		return acceptanceResponse("typesafe/jev-1.13-20261001", 0.001), nil
	}))
	if requests != 1 || record.Calls != 1 || len(record.Tasks) != 1 || record.Tasks[0].Outcome != "answered" {
		t.Fatalf("requests=%d, record=%+v", requests, record)
	}
	path := filepath.Join(t.TempDir(), "record.json")
	if err := writeAcceptanceRecord(path, record); err != nil {
		t.Fatal(err)
	}
	if err := checkAcceptanceRecord(repo, home, acceptanceRead(t, path), []byte("Verdict: inconclusive\n"), labels); err != nil {
		t.Fatalf("fake measurement record failed audit: %v", err)
	}
}
func TestTaskAcceptanceSendsEachKeyOnlyToItsEndpoint(t *testing.T) {
	q, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		keys  map[string]string
		index int
		key   string
	}{
		{"both", map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "or-secret", "ROUNDFIX_TYPESAFE_API_KEY": "ts-secret"}, 0, "or-secret"},
		{"TypeSafe only", map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts-secret"}, 1, "ts-secret"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, home, labels := acceptanceTestRepo(t)
			count := 0
			record := acceptanceMeasure(t, repo, home, labels, tc.keys, roundTripFunc(func(r *http.Request) (*http.Response, error) {
				count++
				if r.URL.String() != q.Transports[tc.index].Endpoint || r.Header.Get("Authorization") != "Bearer "+tc.key || len(r.Header) != 2 {
					t.Fatalf("wrong recipient/header: %s %v", r.URL, r.Header)
				}
				data, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				for _, secret := range tc.keys {
					if strings.Contains(string(data), secret) {
						t.Fatal("key in request body")
					}
				}
				return acceptanceResponse("jev-1.13.0", 0.001), nil
			}))
			if count != 1 || record.Transport != q.Transports[tc.index].Name {
				t.Fatalf("recipient=%s, count=%d", record.Transport, count)
			}
			serialized, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			log := acceptanceRead(t, newJudgeLog(home, acceptanceClock()).path)
			for _, secret := range tc.keys {
				if strings.Contains(string(serialized), secret) || strings.Contains(string(log), secret) {
					t.Fatal("key in record or Judge Log")
				}
			}
		})
	}
}
func TestTaskAcceptanceWithoutARoundfixKeyIsBlocked(t *testing.T) {
	for _, tc := range []struct {
		name string
		keys map[string]string
	}{
		{"no keys", nil},
		{"generic keys only", map[string]string{"OPENROUTER_API_KEY": "generic-or", "TYPESAFE_API_KEY": "generic-ts"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, home, labels := acceptanceTestRepo(t)
			calls := 0
			record, err := measureTaskAcceptance(t.Context(), repo, labels, tc.keys, home, roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return acceptanceResponse("jev-1.13.0", 0), nil }), acceptanceClock)
			if err != nil {
				t.Fatal(err)
			}
			if calls != 0 || record.Calls != 0 || record.Status != "blocked" || record.Verdict != "inconclusive" {
				t.Fatalf("generic-key run=%+v, requests=%d", record, calls)
			}
		})
	}
}
func TestTaskAcceptanceSkipsAnUnpinnedModel(t *testing.T) {
	repo, home, labels := acceptanceTestRepo(t)
	record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts"}, roundTripFunc(func(*http.Request) (*http.Response, error) { return acceptanceResponse("jev-1.14.0", 0.001), nil }))
	if len(record.Tasks) != 1 || record.Tasks[0].Outcome != "skipped" || record.Tasks[0].Noul != nil || record.Statistics.Answered != 0 || record.Calls != 1 {
		t.Fatalf("unpinned record=%+v", record)
	}
	var row logLine
	if err := json.Unmarshal(acceptanceRead(t, newJudgeLog(home, acceptanceClock()).path), &row); err != nil {
		t.Fatal(err)
	}
	if row.Model != "jev-1.14.0" || row.Outcome != "skipped" {
		t.Fatalf("log=%+v", row)
	}
}
func TestTaskAcceptanceLogsEveryRequest(t *testing.T) {
	repo, home, labels := acceptanceTestRepo(t)
	attempts := 0
	record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts"}, roundTripFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return fixtureResponse(429, ""), nil
		}
		return acceptanceResponse("jev-1.13.0", 0.001), nil
	}))
	data := acceptanceRead(t, newJudgeLog(home, acceptanceClock()).path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if attempts != 2 || record.Calls != attempts || len(lines) != attempts {
		t.Fatalf("calls=%d, log=%s", record.Calls, data)
	}
	for i, line := range lines {
		var row logLine
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatal(err)
		}
		if row.Judgment != "task-acceptance" || row.QuestionID != "acceptance_observable" || row.Attempts != i+1 || row.Transport != record.Transport {
			t.Fatalf("row=%+v", row)
		}
	}
}
func TestTaskAcceptanceStopsAtTheCeiling(t *testing.T) {
	q, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"already reached", "reached by first request", "unreadable log"} {
		t.Run(name, func(t *testing.T) {
			repo, home, labels := acceptanceTestRepo(t)
			log := newJudgeLog(home, acceptanceClock())
			switch name {
			case "already reached":
				if err := log.append(logLine{CostUSD: q.MonthlyCeilingUSD}); err != nil {
					t.Fatal(err)
				}
			case "unreadable log":
				acceptanceWrite(t, log.path, "invalid JSON\n")
			case "reached by first request":
				acceptanceWrite(t, filepath.Join(repo, "docs/specs/0300-example/task_02.md"), strings.ReplaceAll(acceptanceFixtureText(), "task_01", "task_02"))
				var cause runcause.Report
				if err := json.Unmarshal(labels, &cause); err != nil {
					t.Fatal(err)
				}
				cause.Tasks = append(cause.Tasks, runcause.TaskCounts{Spec: "0300-example", Task: "task_02", VerdictsPassed: 1})
				labels, err = json.Marshal(cause)
				if err != nil {
					t.Fatal(err)
				}
			}
			requests := 0
			record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts"}, roundTripFunc(func(*http.Request) (*http.Response, error) {
				requests++
				return acceptanceResponse("jev-1.13.0", q.MonthlyCeilingUSD), nil
			}))
			if name == "reached by first request" {
				if requests != 1 || len(record.Tasks) != 2 || record.Tasks[1].Outcome != "skipped" {
					t.Fatalf("ceiling crossed: %+v", record)
				}
			} else if requests != 0 {
				t.Fatalf("sent %d requests past ceiling/unreadable log", requests)
			}
			if record.Status != "blocked" {
				t.Fatalf("status=%s", record.Status)
			}
		})
	}
}
func TestTaskAcceptanceExcludesUnsafeAndMissingStates(t *testing.T) {
	for _, name := range []string{"missing file", "missing criteria", "Portuguese", "Task symlink", "Spec symlink", "root symlink", "nested Task", "directory Task", "archived Task"} {
		t.Run(name, func(t *testing.T) {
			repo, home, labels := acceptanceTestRepo(t)
			active := filepath.Join(repo, "docs/specs/0300-example/task_01.md")
			switch name {
			case "missing file":
				if err := os.Remove(active); err != nil {
					t.Fatal(err)
				}
			case "missing criteria":
				acceptanceWrite(t, active, "# The task title\n## Result\n"+acceptanceResultSentinel)
			case "Portuguese":
				acceptanceWrite(t, active, "# O título da tarefa\n## Acceptance Criteria\n- O comando deve mostrar a saída e o resultado da tarefa.\n")
			case "Task symlink":
				target := filepath.Join(repo, "other.md")
				acceptanceWrite(t, target, acceptanceFixtureText())
				if err := os.Remove(active); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, active); err != nil {
					t.Fatal(err)
				}
			case "Spec symlink":
				dir := filepath.Dir(active)
				moved := filepath.Join(repo, "other-spec")
				if err := os.Rename(dir, moved); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(moved, dir); err != nil {
					t.Fatal(err)
				}
			case "root symlink":
				root := filepath.Join(repo, "docs/specs")
				moved := filepath.Join(repo, "other-root")
				if err := os.Rename(root, moved); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(moved, root); err != nil {
					t.Fatal(err)
				}
			case "nested Task":
				var causes runcause.Report
				if err := json.Unmarshal(labels, &causes); err != nil {
					t.Fatal(err)
				}
				causes.Tasks[0].Task = "nested/task_01"
				labels, _ = json.Marshal(causes)
			case "directory Task":
				if err := os.Remove(active); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(active, 0700); err != nil {
					t.Fatal(err)
				}
			case "archived Task":
				archived := filepath.Join(repo, "docs/history/specs/0300-example/task_01.md")
				acceptanceWrite(t, archived, acceptanceFixtureText())
				if err := os.Remove(active); err != nil {
					t.Fatal(err)
				}
			}
			calls := 0
			record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts"}, roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return acceptanceResponse("jev-1.13.0", 0), nil }))
			if name == "archived Task" {
				if calls != 1 || len(record.Excluded) != 0 {
					t.Fatalf("archive refused: %+v", record)
				}
				return
			}
			if calls != 0 || len(record.Excluded) != 1 || record.Excluded[0].Reason == "" {
				t.Fatalf("unsafe/missing state sent: %+v", record)
			}
		})
	}
}

func TestTaskAcceptanceUsesOnlyEligibleLabels(t *testing.T) {
	repo, home, labels := acceptanceTestRepo(t)
	var cause runcause.Report
	if err := json.Unmarshal(labels, &cause); err != nil {
		t.Fatal(err)
	}
	// Neither entry has a file: filtering must precede even the state reader.
	cause.Tasks = append(cause.Tasks, runcause.TaskCounts{Spec: "0300-example", Task: "task_02", QA: true, VerdictsPassed: 1}, runcause.TaskCounts{Spec: "0300-example", Task: "task_03"})
	labels, err := json.Marshal(cause)
	if err != nil {
		t.Fatal(err)
	}
	record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts"}, acceptanceFake())
	if record.Calls != 1 || len(record.Tasks) != 1 || len(record.Excluded) != 0 || !record.Tasks[0].Repaired {
		t.Fatalf("eligibility=%+v", record)
	}
}
func TestTaskAcceptanceResolvesConfiguredAndArchivedRoot(t *testing.T) {
	for _, archived := range []bool{false, true} {
		name := "active"
		if archived {
			name = "archived"
		}
		t.Run(name, func(t *testing.T) {
			repo, home, labels := acceptanceTestRepo(t)
			if err := os.Mkdir(filepath.Join(repo, ".git"), 0700); err != nil {
				t.Fatal(err)
			}
			external, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			configText := fmt.Sprintf("specs:\n  root: %q\n", external)
			acceptanceWrite(t, filepath.Join(repo, ".roundfixrc.yml"), configText)
			root := external
			if archived {
				root = filepath.Join(external, "_archived")
			}
			acceptanceWrite(t, filepath.Join(root, "0300-example/task_01.md"), acceptanceFixtureText())
			if err := os.Remove(filepath.Join(repo, "docs/specs/0300-example/task_01.md")); err != nil {
				t.Fatal(err)
			}
			record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "ts"}, acceptanceFake())
			if record.Calls != 1 || len(record.Excluded) != 0 {
				t.Fatalf("configured Spec Root not used: %+v", record)
			}
		})
	}
}
func TestTaskAcceptanceStopsOnRefusalAndRedactsErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		failure error
	}{
		{"refused", 401, nil},
		{"transport failure", 0, errors.New("transport reported fixture-key")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, home, labels := acceptanceTestRepo(t)
			record := acceptanceMeasure(t, repo, home, labels, map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "fixture-key"}, roundTripFunc(func(*http.Request) (*http.Response, error) { return fixtureResponse(tc.status, ""), tc.failure }))
			if record.Status != "blocked" || record.Calls != 1 || record.Statistics.Answered != 0 {
				t.Fatalf("refusal=%+v", record)
			}
			data, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			log := acceptanceRead(t, newJudgeLog(home, acceptanceClock()).path)
			if strings.Contains(string(data), "fixture-key") || strings.Contains(string(log), "fixture-key") {
				t.Fatal("error exposed key")
			}
		})
	}
}
