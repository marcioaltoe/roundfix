package judge

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"roundfix/internal/secretscan"
)

type IdeaStageRequest struct {
	RunID, RepoRoot, Project string
	Item                     FollowUpItem
	Keys                     map[string]string
	HomeDir                  string
	Transport                http.RoundTripper
	Now                      func() time.Time
}

type IdeaStageAnswer struct {
	ID          string
	Probability float64
	Answer      string
}

type IdeaStageReport struct {
	State                                  string
	Reason                                 *string
	Verdict                                string
	Answers                                []IdeaStageAnswer
	Model, Transport, KeyVariable          *string
	Calls, InputTokens                     int
	CostUSD, MonthCostUSD, MonthCeilingUSD float64
}

// IdeaStageVerdict combines the independent answers without hiding uncertainty.
func IdeaStageVerdict(answers []string) string {
	allNo := len(answers) > 0
	for _, answer := range answers {
		if answer == "yes" {
			return "indicated"
		}
		if answer != "no" {
			allNo = false
		}
	}
	if allNo {
		return "not indicated"
	}
	return "undecided"
}

// JudgeIdeaStage sends only one bounded, secret-scanned item. Service and log
// failures are advisory reports, never errors or invented answers.
func JudgeIdeaStage(ctx context.Context, q Questions, req IdeaStageRequest) (IdeaStageReport, error) {
	report := IdeaStageReport{State: "judge_unavailable", Answers: []IdeaStageAnswer{}, MonthCeilingUSD: q.MonthlyCeilingUSD}
	item := req.Item
	if slices.Contains(q.IdeaStage.NotJudgedTypes, item.Type) {
		report.State = "not_applicable"
		report.Reason = reasonPointer(fmt.Sprintf("the item is typed %s; only a feature is judged for the Idea stage", item.Type))
		return report, nil
	}
	if req.RunID != "" {
		report.Reason = reasonPointer("inside Run " + req.RunID + ", named by " + RunMarkerVariable)
		return report, nil
	}
	transport, variable, key, ok := selectTransport(q, req.Keys)
	if !ok {
		variables := q.KeyVariables()
		reason := "no key variable is set"
		if len(variables) > 0 {
			reason = variables[0] + " is not set"
			if len(variables) > 1 {
				reason += " (nor " + strings.Join(variables[1:], ", nor ") + ")"
			}
		}
		report.Reason = reasonPointer(reason)
		return report, nil
	}
	redact := func(value string) string {
		for _, name := range q.KeyVariables() {
			if secret := req.Keys[name]; secret != "" {
				value = strings.ReplaceAll(value, secret, "[redacted]")
			}
		}
		return value
	}
	unavailable := func(reason string) (IdeaStageReport, error) {
		report.State, report.Reason, report.Verdict = "judge_unavailable", reasonPointer(redact(reason)), ""
		report.Answers = []IdeaStageAnswer{}
		return report, nil
	}
	now := req.Now
	if now == nil {
		now = time.Now
	}
	ceiling := func() string {
		return fmt.Sprintf("monthly ceiling reached (US$%.4f of US$%.2f)", report.MonthCostUSD, q.MonthlyCeilingUSD)
	}
	month, err := newJudgeLog(req.HomeDir, now()).monthCost()
	if err != nil {
		return unavailable("judge log unreadable: " + err.Error())
	}
	report.MonthCostUSD = month
	if month >= q.MonthlyCeilingUSD {
		return unavailable(ceiling())
	}
	text := []rune(item.Text)
	if len(text) > q.IdeaStage.ItemMaxChars {
		text = text[:q.IdeaStage.ItemMaxChars]
	}
	excerpt := string(text)
	scan := strings.Join([]string{item.Title, item.Origin, excerpt, item.Type, item.Kind}, "\n")
	if matches := secretscan.ScanContent(item.Ref, []byte(scan)); len(matches) > 0 {
		return unavailable("not sent: secret pattern " + matches[0].Kind)
	}
	state, err := json.Marshal(struct {
		Project string            `json:"project"`
		Item    followUpStateItem `json:"item"`
	}{req.Project, followUpStateItem{item.Kind, item.Type, item.Origin, item.Title, excerpt}})
	if err != nil {
		return unavailable("unreadable answer")
	}
	pending := PendingJudgment{Kind: "idea-stage", Artifact: item.Ref, Target: item.Kind, state: state}
	c := client{http: &http.Client{Transport: req.Transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, transport: transport, key: key, q: q}
	for attempt := 1; attempt <= 3; attempt++ {
		callTime := now().UTC()
		month, err := newJudgeLog(req.HomeDir, callTime).monthCost()
		if err != nil {
			return unavailable("judge log unreadable: " + err.Error())
		}
		report.MonthCostUSD = month
		if month >= q.MonthlyCeilingUSD {
			return unavailable(ceiling())
		}
		if err := ctx.Err(); err != nil {
			return unavailable("service unavailable (" + err.Error() + ")")
		}
		result := c.ask(ctx, pending, attempt)
		result.Error, result.Model, result.ResponseID, result.Provider = redact(result.Error), redact(result.Model), redact(result.ResponseID), redact(result.Provider)
		a, outcome, reason, _ := evaluate(q, pending, result)
		cost, source := result.cost(q)
		report.Calls++
		report.InputTokens += result.Usage.InputTokens
		report.CostUSD += cost
		report.MonthCostUSD += cost
		report.Model, report.Transport, report.KeyVariable = reasonPointer(result.Model), &transport.Name, &variable
		hash := sha256.Sum256(state)
		row := logLine{Schema: "roundfix/judge-log/v1", Time: callTime, Repository: req.RepoRoot, Judgment: pending.Kind, Artifact: item.Ref, Target: item.Kind, StateHash: fmt.Sprintf("%x", hash)[:16], QuestionID: q.IdeaStage.QuestionID, Transport: transport.Name, KeyVariable: variable, ResponseID: result.ResponseID, Provider: result.Provider, RequestedModel: transport.RequestModel, Model: result.Model, Probabilities: a.Probabilities, CostUSD: cost, CostSource: source, LatencyMS: result.LatencyMS, InputTokens: result.Usage.InputTokens, OutputTokens: result.Usage.OutputTokens, Status: result.Status, Attempts: attempt, Error: result.Error, Outcome: outcome}
		for _, field := range []*string{&row.Repository, &row.Artifact, &row.Target, &row.QuestionID, &row.RequestedModel} {
			*field = redact(*field)
		}
		if err := newJudgeLog(req.HomeDir, callTime).append(row); err != nil {
			return unavailable("judge log not writable: " + err.Error())
		}
		if (result.Status == 429 || result.Status == 529) && attempt < 3 {
			if err := waitRetry(ctx, retryDelay(result.RetryAfter, attempt)); err != nil {
				return unavailable("service unavailable (" + err.Error() + ")")
			}
			continue
		}
		if outcome != "answered" {
			return unavailable(reason)
		}
		answers := make([]string, 0, len(q.IdeaStage.Criteria))
		for _, criterion := range q.IdeaStage.Criteria {
			probability := a.Probabilities[criterion.ID]
			answer := "unsure"
			if probability >= q.IdeaStage.YesAtLeast {
				answer = "yes"
			} else if probability <= q.IdeaStage.NoAtMost {
				answer = "no"
			}
			report.Answers = append(report.Answers, IdeaStageAnswer{criterion.ID, probability, answer})
			answers = append(answers, answer)
		}
		report.State, report.Reason, report.Verdict = "answered", nil, IdeaStageVerdict(answers)
		return report, nil
	}
	return report, nil
}
