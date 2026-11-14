package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/judge"
)

type ideaStageDocumentItem struct {
	Ref   string `json:"ref"`
	Kind  string `json:"kind"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type ideaStageDocumentCriterion struct {
	ID          string   `json:"id"`
	Summary     string   `json:"summary"`
	Yes         string   `json:"yes"`
	No          string   `json:"no"`
	Probability *float64 `json:"probability"`
	Answer      *string  `json:"answer"`
}

type ideaStageDocument struct {
	Schema     string                       `json:"schema"`
	Item       ideaStageDocumentItem        `json:"item"`
	State      string                       `json:"state"`
	Reason     *string                      `json:"reason"`
	Verdict    *string                      `json:"verdict"`
	JudgedBy   *string                      `json:"judged_by"`
	Criteria   []ideaStageDocumentCriterion `json:"criteria"`
	Thresholds struct {
		YesAtLeast float64 `json:"yes_at_least"`
		NoAtMost   float64 `json:"no_at_most"`
	} `json:"thresholds"`
	Model           *string `json:"model"`
	Transport       *string `json:"transport"`
	KeyVariable     *string `json:"key_variable"`
	Calls           int     `json:"calls"`
	InputTokens     int     `json:"input_tokens"`
	CostUSD         float64 `json:"cost_usd"`
	MonthCostUSD    float64 `json:"month_cost_usd"`
	MonthCeilingUSD float64 `json:"month_ceiling_usd"`
}

// ideaStageItem is the sole reference lookup; disk paths are never a fallback.
func ideaStageItem(items []followUpItem, reference string) (followUpItem, error) {
	for _, item := range items {
		if item.Ref != reference {
			continue
		}
		if item.Kind == "spec" {
			return followUpItem{}, fmt.Errorf("%q is an active Spec; the Idea stage is judged before a PRD exists", reference)
		}
		return item, nil
	}
	return followUpItem{}, fmt.Errorf("no pending item %q; name a reference that 'roundfix follow-up' lists", reference)
}

// ideaStageRequestItem reads the caller's file without changing it. Its path
// never becomes a reference or a field of the state sent to the judge.
func ideaStageRequestItem(path, project string) (followUpItem, error) {
	fail := func(err error) (followUpItem, error) {
		return followUpItem{}, fmt.Errorf("request file %q: %w", path, err)
	}
	file, err := os.Open(path)
	if err != nil {
		return fail(err)
	}
	defer file.Close()
	const limit = 1 << 20
	content, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return fail(err)
	}
	if len(content) > limit {
		return fail(fmt.Errorf("exceeds 1 MiB"))
	}
	text := strings.TrimSpace(string(content))
	if text == "" {
		return fail(fmt.Errorf("holds no text"))
	}
	title, _, _ := strings.Cut(text, "\n")
	title = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(title), "# "))
	runes := []rune(title)
	if len(runes) > 200 {
		title = string(runes[:200])
	}
	return followUpItem{Ref: "the request in the session", Kind: "request", Type: "request", Title: title, text: text, origin: project}, nil
}

func runIdeaStageCommand(ctx context.Context, req specJudgeRequest, loaded roundconfig.Loaded, environment commandEnvironment, stdout io.Writer) error {
	project := followUpNamespace(ctx, loaded)
	var item followUpItem
	var err error
	if req.request != "" {
		path := req.request
		if !filepath.IsAbs(path) {
			path = filepath.Join(environment.workDir, path)
		}
		item, err = ideaStageRequestItem(path, project)
	} else {
		var items []followUpItem
		items, _, err = followUpInventory(loaded, project, req.messages, environment.workDir)
		if err == nil {
			item, err = ideaStageItem(items, req.item)
		}
	}
	if err != nil {
		return err
	}
	q, err := judge.Load()
	if err != nil {
		return err
	}
	q = q.WithMonthlyCeiling(loaded.Config.Jev.MonthlyCeilingUSD)
	runID := commandRunID(environment.environ)
	keys := make(map[string]string)
	if runID == "" {
		allowed := make(map[string]bool)
		for _, variable := range q.KeyVariables() {
			allowed[variable] = true
		}
		for _, entry := range environment.environ {
			name, value, _ := strings.Cut(entry, "=")
			if allowed[name] {
				keys[name] = value
			}
		}
	}
	report := judge.IdeaStageReport{State: "not_applicable", MonthCeilingUSD: q.MonthlyCeilingUSD}
	if item.WaitingOn != "" {
		reason := fmt.Sprintf("the item waits on %s since %s", item.WaitingOn, item.Since)
		report.Reason = &reason
	} else {
		report, err = judge.JudgeIdeaStage(ctx, q, judge.IdeaStageRequest{
			RunID: runID, RepoRoot: loaded.GitRoot, Project: project,
			Item: judge.FollowUpItem{Ref: item.Ref, Kind: item.Kind, Type: item.Type, Origin: item.origin, Title: item.Title, Text: item.text},
			Keys: keys, HomeDir: environment.homeDir, Transport: environment.dependencies.judgeTransport, Now: environment.dependencies.judgeNow,
		})
		if err != nil {
			return err
		}
	}
	doc := ideaStageDocument{Schema: "roundfix/idea-stage-judgment/v1", Item: ideaStageDocumentItem{item.Ref, item.Kind, item.Type, item.Title}, State: report.State, Reason: report.Reason, Criteria: []ideaStageDocumentCriterion{}, Model: report.Model, Transport: report.Transport, KeyVariable: report.KeyVariable, Calls: report.Calls, InputTokens: report.InputTokens, CostUSD: report.CostUSD, MonthCostUSD: report.MonthCostUSD, MonthCeilingUSD: report.MonthCeilingUSD}
	doc.Thresholds.YesAtLeast, doc.Thresholds.NoAtMost = q.IdeaStage.YesAtLeast, q.IdeaStage.NoAtMost
	if report.State == "answered" {
		by := "judge"
		doc.Verdict, doc.JudgedBy = &report.Verdict, &by
	}
	for i, criterion := range q.IdeaStage.Criteria {
		field := ideaStageDocumentCriterion{ID: criterion.ID, Summary: criterion.Summary, Yes: criterion.Question.Criteria["true"], No: criterion.Question.Criteria["false"]}
		if report.State == "answered" {
			field.Probability, field.Answer = &report.Answers[i].Probability, &report.Answers[i].Answer
		}
		doc.Criteria = append(doc.Criteria, field)
	}
	if req.format == "json" {
		if err := json.NewEncoder(stdout).Encode(doc); err != nil {
			return fmt.Errorf("write judge JSON: %w", err)
		}
	} else {
		fmt.Fprint(stdout, renderIdeaStageText(doc))
	}
	return nil
}

func renderIdeaStageText(doc ideaStageDocument) string {
	var out strings.Builder
	const rule = "Rule: indicated when any criterion is yes, not indicated when every criterion is no, undecided otherwise"
	if doc.State == "not_applicable" {
		return fmt.Sprintf("Idea stage for %s: not judged: %s\n", doc.Item.Ref, *doc.Reason)
	}
	if doc.State == "answered" {
		fmt.Fprintf(&out, "Idea stage for %s: %s\n", doc.Item.Ref, *doc.Verdict)
		for _, criterion := range doc.Criteria {
			fmt.Fprintf(&out, "  %s: %s (%.2f) — %s\n", criterion.ID, *criterion.Answer, *criterion.Probability, criterion.Summary)
		}
		fmt.Fprintf(&out, "%s; yes at %.2f or more, no at %.2f or less.\nRecord:\n  - item: %s\n  - judged-by: judge (%s)\n", rule, doc.Thresholds.YesAtLeast, doc.Thresholds.NoAtMost, doc.Item.Ref, *doc.Model)
		for _, criterion := range doc.Criteria {
			fmt.Fprintf(&out, "  - %s: %s (%.2f)\n", criterion.ID, *criterion.Answer, *criterion.Probability)
		}
		fmt.Fprintf(&out, "  - verdict: %s\n", *doc.Verdict)
	} else {
		fmt.Fprintf(&out, "Idea stage for %s: not judged: %s\nAnswer the three criteria yourself, each as yes, no or unsure, and say that the judge did not answer:\n", doc.Item.Ref, *doc.Reason)
		for _, criterion := range doc.Criteria {
			fmt.Fprintf(&out, "  %s: %s\n    yes: %s\n    no: %s\n", criterion.ID, criterion.Summary, criterion.Yes, criterion.No)
		}
		fmt.Fprintf(&out, "%s.\nRecord:\n  - item: %s\n  - judged-by: agent (%s)\n", rule, doc.Item.Ref, *doc.Reason)
	}
	if doc.Calls > 0 {
		model := ""
		if doc.Model != nil {
			model = *doc.Model
		}
		fmt.Fprintf(&out, "Judge: %d call(s), %d input tokens, US$%.4f; month US$%.4f of US$%.2f; model %s via %s on %s\n", doc.Calls, doc.InputTokens, doc.CostUSD, doc.MonthCostUSD, doc.MonthCeilingUSD, model, *doc.Transport, *doc.KeyVariable)
	}
	return out.String()
}
