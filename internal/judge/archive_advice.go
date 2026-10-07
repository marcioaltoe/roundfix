package judge

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ArchiveAdviceRequest struct {
	RepoRoot, SpecDir, Spec string
	Files                   []string
	Keys                    map[string]string
	HomeDir                 string
	Transport               http.RoundTripper
	Now                     func() time.Time
}

type ArchiveAdvice struct {
	File          string
	Outcome       string
	Choice        *string
	Probabilities map[string]float64
	Confidence    *float64
	Reason        *string
	Model         *string
}

type ArchiveAdviceReport struct {
	Model                                  string
	Skipped, Stopped                       *string
	Advice                                 []ArchiveAdvice
	Calls                                  int
	CostUSD, MonthCostUSD, MonthCeilingUSD float64
}

// AdviseArchive classifies candidate files for an archive. It is advisory:
// credential, service, and log failures are represented in the report and do
// not prevent the archive command from proceeding.
func AdviseArchive(ctx context.Context, q Questions, req ArchiveAdviceRequest) (ArchiveAdviceReport, error) {
	report := ArchiveAdviceReport{Model: q.PinnedModel, Advice: []ArchiveAdvice{}, MonthCeilingUSD: q.MonthlyCeilingUSD}
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
	for _, name := range req.Files {
		advice := ArchiveAdvice{File: filepath.ToSlash(name), Outcome: "skipped"}
		if isArchiveCore(name) || strings.HasPrefix(filepath.ToSlash(name), "qa/evidence/") {
			continue
		}
		path := filepath.Join(req.SpecDir, filepath.Clean(name))
		data, err := os.ReadFile(path)
		if err != nil {
			advice.Reason = reasonPointer("unreadable file: " + err.Error())
			report.Advice = append(report.Advice, advice)
			continue
		}
		if isBinary(data) {
			advice.Reason = reasonPointer("binary")
			report.Advice = append(report.Advice, advice)
			continue
		}
		if report.Skipped != nil || report.Stopped != nil {
			advice.Reason = report.Skipped
			if report.Stopped != nil {
				advice.Reason = report.Stopped
			}
			report.Advice = append(report.Advice, advice)
			continue
		}
		state, _ := json.Marshal(struct {
			File    string `json:"file"`
			Content string `json:"content"`
		}{File: filepath.ToSlash(name), Content: string(data[:min(len(data), 6000)])})
		pending := PendingJudgment{Kind: "archive-value", Artifact: filepath.Join(req.SpecDir, name), Target: filepath.ToSlash(name), state: state}
		c := client{http: &http.Client{Transport: req.Transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, transport: t, key: key, q: q}
		for attempt := 1; attempt <= 3; attempt++ {
			callTime := now().UTC()
			month, err := newJudgeLog(req.HomeDir, callTime).monthCost()
			if err != nil {
				report.Stopped = reasonPointer("judge log unreadable: " + redact(err.Error()))
				advice.Reason = report.Stopped
				break
			}
			report.MonthCostUSD = month
			if month >= q.MonthlyCeilingUSD {
				report.Stopped = reasonPointer(ceiling())
				advice.Reason = report.Stopped
				break
			}
			if err := ctx.Err(); err != nil {
				report.Stopped = reasonPointer("service unavailable (" + err.Error() + ")")
				advice.Reason = report.Stopped
				break
			}
			result := c.ask(ctx, pending, attempt)
			result.Error, result.Model, result.ResponseID, result.Provider = redact(result.Error), redact(result.Model), redact(result.ResponseID), redact(result.Provider)
			a, outcome, reason, stop := evaluate(q, pending, result)
			advice.Outcome, advice.Choice, advice.Probabilities, advice.Confidence, advice.Reason, advice.Model = outcome, a.Choice, a.Probabilities, a.Confidence, reasonPointer(reason), reasonPointer(result.Model)
			cost, source := result.cost(q)
			report.Calls++
			report.CostUSD += cost
			report.MonthCostUSD += cost
			hash := sha256.Sum256(state)
			row := logLine{Schema: "roundfix/judge-log/v1", Time: callTime, Repository: req.RepoRoot, Spec: req.Spec, Judgment: pending.Kind, Artifact: advice.File, Target: advice.File, StateHash: fmt.Sprintf("%x", hash)[:16], QuestionID: q.Archive.QuestionID, Transport: t.Name, KeyVariable: variable, ResponseID: result.ResponseID, Provider: result.Provider, RequestedModel: t.RequestModel, Model: result.Model, Answer: a.Choice, Probabilities: a.Probabilities, Confidence: a.Confidence, CostUSD: cost, CostSource: source, Status: result.Status, Attempts: attempt, Error: result.Error, Outcome: outcome}
			if err := newJudgeLog(req.HomeDir, callTime).append(row); err != nil {
				report.Stopped = reasonPointer("judge log not writable: " + redact(err.Error()))
				advice.Reason = report.Stopped
				break
			}
			if stop {
				report.Stopped = reasonPointer(reason)
				advice.Reason = report.Stopped
				break
			}
			if result.Status != 429 && result.Status != 529 {
				break
			}
			if err := waitRetry(ctx, retryDelay(result.RetryAfter, attempt)); err != nil {
				report.Stopped = reasonPointer("service unavailable (" + err.Error() + ")")
				advice.Reason = report.Stopped
				break
			}
		}
		report.Advice = append(report.Advice, advice)
	}
	return report, nil
}

func isArchiveCore(name string) bool {
	name = filepath.ToSlash(filepath.Clean(name))
	base := filepath.Base(name)
	return name == "_prd.md" || name == "_techspec.md" || name == "_tasks.md" || name == "_authorization.md" || (strings.HasPrefix(base, "task_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "qa-report-") && strings.HasSuffix(base, ".md")) || name == "references/_index.md"
}
func isBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
