package judge

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"roundfix/internal/secretscan"
)

type FollowUpItem struct{ Ref, Kind, Type, Origin, Title, Text string }
type FollowUpRequest struct {
	RunID, RepoRoot, Project string
	Items                    []FollowUpItem
	Keys                     map[string]string
	HomeDir                  string
	Transport                http.RoundTripper
	Now                      func() time.Time
}
type FollowUpAnswer struct {
	Ref, Outcome string
	Reason       *string
	Nouls        map[string]float64
	Verdicts     map[string]string
	Model        *string
}
type FollowUpReport struct {
	Model                                  string
	Skipped, Stopped                       *string
	Answers                                []FollowUpAnswer
	Calls                                  int
	CostUSD, MonthCostUSD, MonthCeilingUSD float64
}

type followUpStateItem struct {
	Kind   string `json:"kind"`
	Type   string `json:"type"`
	Origin string `json:"origin"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

// PrioritizeFollowUp asks four independent questions per item. It is advisory:
// credential, service, and log failures are represented in the report and do
// not prevent callers from ordering items by their declared type.
func PrioritizeFollowUp(ctx context.Context, q Questions, req FollowUpRequest) (FollowUpReport, error) {
	report := FollowUpReport{Model: q.PinnedModel, Answers: []FollowUpAnswer{}, MonthCeilingUSD: q.MonthlyCeilingUSD}
	if req.RunID != "" {
		report.Skipped = reasonPointer("inside Run " + req.RunID + ", named by " + RunMarkerVariable)
		for _, item := range req.Items {
			report.Answers = append(report.Answers, FollowUpAnswer{Ref: item.Ref, Outcome: "skipped", Reason: report.Skipped})
		}
		return report, nil
	}
	now := req.Now
	if now == nil {
		now = time.Now
	}
	t, variable, key, ok := selectTransport(q, req.Keys)
	log := newJudgeLog(req.HomeDir, now())
	ceiling := func() string {
		return fmt.Sprintf("monthly ceiling reached (US$%.4f of US$%.2f)", report.MonthCostUSD, q.MonthlyCeilingUSD)
	}
	if !ok {
		vars := q.KeyVariables()
		reason := "no key variable is set"
		if len(vars) > 0 {
			reason = vars[0] + " is not set"
		}
		report.Skipped = reasonPointer(reason)
	} else if cost, err := log.monthCost(); err != nil {
		report.Skipped = reasonPointer("judge log unreadable: " + err.Error())
	} else {
		report.MonthCostUSD = cost
		if cost >= q.MonthlyCeilingUSD {
			report.Skipped = reasonPointer(ceiling())
		}
		report.Model = q.PinnedModel
	}
	secrets := []string{}
	for _, name := range q.KeyVariables() {
		if secret := req.Keys[name]; secret != "" {
			secrets = append(secrets, secret)
		}
	}
	redact := func(s string) string {
		for _, secret := range secrets {
			s = strings.ReplaceAll(s, secret, "[redacted]")
		}
		return s
	}
	if report.Skipped != nil {
		report.Skipped = reasonPointer(redact(*report.Skipped))
	}
	for index, item := range req.Items {
		itemAnswer := FollowUpAnswer{Ref: item.Ref, Outcome: "skipped"}
		if index >= q.FollowUp.MaxItems {
			itemAnswer.Reason = reasonPointer(fmt.Sprintf("item limit reached (%d)", q.FollowUp.MaxItems))
			report.Answers = append(report.Answers, itemAnswer)
			continue
		}
		if report.Skipped != nil || report.Stopped != nil {
			itemAnswer.Reason = report.Skipped
			if report.Stopped != nil {
				itemAnswer.Reason = report.Stopped
			}
			report.Answers = append(report.Answers, itemAnswer)
			continue
		}
		text := []rune(item.Text)
		if len(text) > q.FollowUp.ItemMaxChars {
			text = text[:q.FollowUp.ItemMaxChars]
		}
		excerpt := string(text)
		scanText := strings.Join([]string{item.Title, item.Origin, excerpt, item.Type, item.Kind}, "\n")
		if matches := secretscan.ScanContent(item.Ref, []byte(scanText)); len(matches) > 0 {
			itemAnswer.Reason = reasonPointer("not sent: secret pattern " + matches[0].Kind)
			report.Answers = append(report.Answers, itemAnswer)
			continue
		}
		state, err := json.Marshal(struct {
			Project string            `json:"project"`
			Item    followUpStateItem `json:"item"`
		}{req.Project, followUpStateItem{item.Kind, item.Type, item.Origin, item.Title, excerpt}})
		if err != nil {
			return report, fmt.Errorf("encode follow-up item: %w", err)
		}
		pending := PendingJudgment{Kind: "follow-up-priority", Artifact: item.Ref, Target: item.Kind, state: state}
		c := client{http: &http.Client{Transport: req.Transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, transport: t, key: key, q: q}
		for attempt := 1; attempt <= 3; attempt++ {
			callTime := now().UTC()
			month, err := newJudgeLog(req.HomeDir, callTime).monthCost()
			if err != nil {
				report.Stopped = reasonPointer("judge log unreadable: " + redact(err.Error()))
				itemAnswer.Reason = report.Stopped
				break
			}
			report.MonthCostUSD = month
			if month >= q.MonthlyCeilingUSD {
				report.Stopped = reasonPointer(ceiling())
				itemAnswer.Reason = report.Stopped
				break
			}
			if err := ctx.Err(); err != nil {
				report.Stopped = reasonPointer("service unavailable (" + err.Error() + ")")
				itemAnswer.Reason = report.Stopped
				break
			}
			result := c.ask(ctx, pending, attempt)
			result.Error, result.Model, result.ResponseID, result.Provider = redact(result.Error), redact(result.Model), redact(result.ResponseID), redact(result.Provider)
			a, outcome, reason, stop := evaluate(q, pending, result)
			itemAnswer.Outcome, itemAnswer.Nouls, itemAnswer.Reason, itemAnswer.Model = outcome, a.Probabilities, reasonPointer(reason), reasonPointer(result.Model)
			if outcome == "answered" {
				itemAnswer.Verdicts = make(map[string]string, len(a.Probabilities))
				for id, value := range a.Probabilities {
					verdict := "unsure"
					if value >= q.FollowUp.YesAtLeast {
						verdict = "yes"
					} else if value <= q.FollowUp.NoAtMost {
						verdict = "no"
					}
					itemAnswer.Verdicts[id] = verdict
				}
			}
			cost, source := result.cost(q)
			report.Calls++
			report.CostUSD += cost
			report.MonthCostUSD += cost
			hash := sha256.Sum256(state)
			row := logLine{Schema: "roundfix/judge-log/v1", Time: callTime, Repository: req.RepoRoot, Judgment: pending.Kind, Artifact: item.Ref, Target: item.Kind, StateHash: fmt.Sprintf("%x", hash)[:16], QuestionID: q.FollowUp.QuestionID, Transport: t.Name, KeyVariable: variable, ResponseID: result.ResponseID, Provider: result.Provider, RequestedModel: t.RequestModel, Model: result.Model, Probabilities: a.Probabilities, CostUSD: cost, CostSource: source, LatencyMS: result.LatencyMS, InputTokens: result.Usage.InputTokens, OutputTokens: result.Usage.OutputTokens, Status: result.Status, Attempts: attempt, Error: result.Error, Outcome: outcome}
			if err := newJudgeLog(req.HomeDir, callTime).append(row); err != nil {
				report.Stopped = reasonPointer("judge log not writable: " + redact(err.Error()))
				itemAnswer.Reason = report.Stopped
				break
			}
			if stop && (attempt == 3 || (result.Status != 429 && result.Status != 529)) {
				report.Stopped = reasonPointer(reason)
				itemAnswer.Reason = report.Stopped
				break
			}
			if result.Status != 429 && result.Status != 529 {
				break
			}
			if err := waitRetry(ctx, retryDelay(result.RetryAfter, attempt)); err != nil {
				report.Stopped = reasonPointer("service unavailable (" + err.Error() + ")")
				itemAnswer.Reason = report.Stopped
				break
			}
		}
		report.Answers = append(report.Answers, itemAnswer)
	}
	return report, nil
}
