package judge

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Request struct {
	RepoRoot, SpecDir, Spec string
	Stage                   Stage
	Keys                    map[string]string
	HomeDir                 string
	Transport               http.RoundTripper
	Now                     func() time.Time
}
type Judgment struct {
	Kind          string             `json:"kind"`
	Artifact      string             `json:"artifact"`
	Line          int                `json:"line"`
	Target        string             `json:"target"`
	Text          string             `json:"text"`
	SectionTitle  *string            `json:"section_title"`
	Outcome       string             `json:"outcome"`
	Reason        *string            `json:"reason"`
	Answer        *string            `json:"answer"`
	Probabilities map[string]float64 `json:"probabilities"`
	Confidence    *float64           `json:"confidence"`
	Noul          *float64           `json:"noul"`
	Model         *string            `json:"model"`
}
type ArtifactSkip struct {
	Artifact string `json:"artifact"`
	Reason   string `json:"reason"`
}
type Report struct {
	Schema           string         `json:"schema"`
	Spec             string         `json:"spec"`
	Model            string         `json:"model"`
	Transport        *string        `json:"transport"`
	Skipped          *string        `json:"skipped"`
	Stopped          *string        `json:"stopped"`
	ArtifactsSkipped []ArtifactSkip `json:"artifacts_skipped"`
	Judgments        []Judgment     `json:"judgments"`
	Calls            int            `json:"calls"`
	InputTokens      int            `json:"input_tokens"`
	CostUSD          float64        `json:"cost_usd"`
	MonthCostUSD     float64        `json:"month_cost_usd"`
	MonthCeilingUSD  float64        `json:"month_ceiling_usd"`
}

func reasonPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func relativeArtifact(root, artifact string) string {
	if !filepath.IsAbs(artifact) {
		return filepath.ToSlash(artifact)
	}
	relative, err := filepath.Rel(root, artifact)
	if err != nil {
		return filepath.ToSlash(artifact)
	}
	return filepath.ToSlash(relative)
}

// Run fails open for judgments and service/log errors. An unreadable required
// artifact is an error; environment and filesystem roots are explicit.
func Run(ctx context.Context, q Questions, req Request) (Report, error) {
	report := Report{Schema: "roundfix/spec-judge/v1", Spec: req.Spec, Model: q.PinnedModel, MonthCeilingUSD: q.MonthlyCeilingUSD, Judgments: []Judgment{}, ArtifactsSkipped: []ArtifactSkip{}}
	// Invalid stages are a command validation concern; retain the default plan.
	stage := req.Stage
	if stage != "prd" && stage != "techspec" && stage != "tasks" {
		stage = ""
	}
	plan, err := PlanSpec(q, req.RepoRoot, req.SpecDir, stage)
	if err != nil {
		return report, err
	}
	for _, a := range plan.SkippedArtifacts {
		report.ArtifactsSkipped = append(report.ArtifactsSkipped, ArtifactSkip{relativeArtifact(req.RepoRoot, artifactSkipPath(req.SpecDir, a.Artifact)), a.Reason})
	}
	for _, s := range plan.Skipped {
		report.Judgments = append(report.Judgments, Judgment{Kind: s.Kind, Artifact: relativeArtifact(req.RepoRoot, s.Artifact), Line: s.Line, Target: s.Target, Outcome: "skipped", Reason: reasonPointer(s.Reason)})
	}
	t, key, ok := selectTransport(q, req.Keys)
	now := req.Now
	if now == nil {
		now = time.Now
	}
	log := newJudgeLog(req.HomeDir, now())
	ceilingReason := func() string {
		return fmt.Sprintf("monthly ceiling reached (US$%.4f of US$%.2f)", report.MonthCostUSD, q.MonthlyCeilingUSD)
	}
	if !ok {
		report.Skipped = reasonPointer("ROUNDFIX_OPENROUTER_API_KEY is not set (nor ROUNDFIX_TYPESAFE_API_KEY)")
	} else {
		report.Transport = &t.Name
		cost, err := log.monthCost()
		if err != nil {
			report.Skipped = reasonPointer("judge log unreadable: " + err.Error())
		} else {
			report.MonthCostUSD = cost
			if cost >= q.MonthlyCeilingUSD {
				report.Skipped = reasonPointer(ceilingReason())
			}
		}
	}
	// Snapshot credentials with the transport so later changes cannot weaken redaction.
	secrets := make([]string, 0, len(q.Transports))
	for _, candidate := range q.Transports {
		if secret := req.Keys[candidate.KeyVariable]; secret != "" {
			secrets = append(secrets, secret)
		}
	}
	redact := func(s string) string {
		for _, secret := range secrets {
			if secret != "" {
				s = strings.ReplaceAll(s, secret, "[redacted]")
			}
		}
		return s
	}
	httpClient := &http.Client{Transport: req.Transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	c := client{http: httpClient, transport: t, key: key, q: q}
	seen := map[string]Judgment{}
	for _, p := range plan.Pending {
		j := Judgment{Kind: p.Kind, Artifact: relativeArtifact(req.RepoRoot, p.Artifact), Line: p.Line, Target: relativeArtifact(req.RepoRoot, p.Target), Text: p.Text, Outcome: "skipped"}
		if p.Kind == "goal-mechanism" {
			_, title, _ := strings.Cut(p.Target, " → ")
			j.SectionTitle = reasonPointer(title)
		}
		if report.Skipped != nil || report.Stopped != nil {
			j.Reason = report.Skipped
			if report.Stopped != nil {
				j.Reason = report.Stopped
			}
			report.Judgments = append(report.Judgments, j)
			continue
		}
		if previous, exists := seen[string(p.state)]; exists {
			j.Outcome, j.Reason, j.Answer, j.Probabilities, j.Confidence, j.Noul, j.Model = previous.Outcome, previous.Reason, previous.Answer, previous.Probabilities, previous.Confidence, previous.Noul, previous.Model
			report.Judgments = append(report.Judgments, j)
			continue
		}
		for attempt := 1; attempt <= 3; attempt++ {
			callTime := now().UTC()
			callLog := newJudgeLog(req.HomeDir, callTime)
			monthCost, err := callLog.monthCost()
			if err != nil {
				report.Stopped = reasonPointer("judge log unreadable: " + redact(err.Error()))
				j.Reason = report.Stopped
				break
			}
			report.MonthCostUSD = monthCost
			if report.MonthCostUSD >= q.MonthlyCeilingUSD {
				report.Stopped = reasonPointer(ceilingReason())
				j.Reason = report.Stopped
				break
			}
			if err := ctx.Err(); err != nil {
				report.Stopped = reasonPointer("service unavailable (" + err.Error() + ")")
				j.Reason = report.Stopped
				break
			}
			result := c.ask(ctx, p, attempt)
			result.Error, result.Model, result.ResponseID, result.Provider = redact(result.Error), redact(result.Model), redact(result.ResponseID), redact(result.Provider)
			a, outcome, reason, stop := evaluate(q, p, result)
			if a.Probabilities != nil {
				probabilities := make(map[string]float64, len(a.Probabilities))
				for name, value := range a.Probabilities {
					probabilities[redact(name)] = value
				}
				a.Probabilities = probabilities
			}
			if a.Choice != nil {
				choice := redact(*a.Choice)
				a.Choice = &choice
			}
			retry := (result.Status == 429 || result.Status == 529) && attempt < 3 && result.Error == ""
			if retry {
				stop = false
			}
			j.Outcome, j.Reason, j.Answer, j.Probabilities, j.Confidence, j.Noul, j.Model = outcome, reasonPointer(reason), a.Choice, a.Probabilities, a.Confidence, a.Noul, reasonPointer(result.Model)
			cost, costSource := result.cost(q)
			report.Calls++
			report.InputTokens += result.Usage.InputTokens
			report.CostUSD += cost
			report.MonthCostUSD += cost
			id := q.Citation.QuestionID
			if p.Kind == "goal-mechanism" {
				id = q.Goal.QuestionID
			}
			if p.Kind == "source-grouping" {
				id = q.Grouping.QuestionID
			}
			if p.Kind == "model-tier" {
				id = q.ModelTier.QuestionID
			}
			hash := sha256.Sum256(p.state)
			recordError := result.Error
			if recordError == "" && outcome == "skipped" {
				recordError = reason
			}
			row := logLine{Schema: "roundfix/judge-log/v1", Time: callTime, Repository: req.RepoRoot, Spec: req.Spec, Judgment: p.Kind, Artifact: j.Artifact, Line: p.Line, Target: j.Target, StateHash: fmt.Sprintf("%x", hash)[:16], QuestionID: id, Transport: t.Name, ResponseID: result.ResponseID, Provider: result.Provider, RequestedModel: t.RequestModel, Model: result.Model, Answer: a.Choice, Probabilities: a.Probabilities, Confidence: a.Confidence, Noul: a.Noul, LatencyMS: result.LatencyMS, InputTokens: result.Usage.InputTokens, OutputTokens: result.Usage.OutputTokens, CostUSD: cost, CostSource: costSource, Status: result.Status, Attempts: attempt, Error: recordError, Outcome: outcome}
			for _, field := range []*string{&row.Schema, &row.Repository, &row.Spec, &row.Judgment, &row.Artifact, &row.Target, &row.StateHash, &row.QuestionID, &row.Transport, &row.RequestedModel, &row.CostSource, &row.Outcome} {
				*field = redact(*field)
			}
			if err := callLog.append(row); err != nil {
				report.Stopped = reasonPointer("judge log not writable: " + redact(err.Error()))
				break
			}
			if stop {
				report.Stopped = reasonPointer(reason)
				break
			}
			if !retry {
				break
			}
			if err := waitRetry(ctx, retryDelay(result.RetryAfter, attempt)); err != nil {
				report.Stopped = reasonPointer("service unavailable (" + err.Error() + ")")
				j.Reason = report.Stopped
				break
			}
		}
		seen[string(p.state)] = j
		report.Judgments = append(report.Judgments, j)
	}
	return report, nil
}

func artifactSkipPath(specDir, artifact string) string {
	if filepath.IsAbs(artifact) {
		return artifact
	}
	return filepath.Join(specDir, artifact)
}

// Grouping has no line or excerpt; retain existing Go fields for callers.
func (j Judgment) MarshalJSON() ([]byte, error) {
	type plain Judgment
	if j.Kind != "source-grouping" {
		return json.Marshal(plain(j))
	}
	return json.Marshal(struct {
		plain
		Line *int    `json:"line"`
		Text *string `json:"text"`
	}{plain: plain(j)})
}
