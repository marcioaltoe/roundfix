package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func selectTransport(q Questions, keys map[string]string) (Transport, string, bool) {
	for _, t := range q.Transports {
		if key := keys[t.KeyVariable]; key != "" {
			return t, key, true
		}
	}
	return Transport{}, "", false
}
func (q Questions) pinned(model string) bool { return q.AcceptedModel.MatchString(model) }

type answer struct {
	Type          string             `json:"type"`
	Choice        *string            `json:"choice"`
	Probabilities map[string]float64 `json:"probabilities"`
	Confidence    *float64           `json:"confidence"`
	Noul          *float64           `json:"noul"`
}
type call struct {
	Model      string            `json:"model"`
	ResponseID string            `json:"id"`
	Provider   string            `json:"provider"`
	Answers    map[string]answer `json:"answers"`
	Usage      struct {
		InputTokens  int             `json:"input_tokens"`
		OutputTokens int             `json:"output_tokens"`
		Cost         json.RawMessage `json:"cost"`
	} `json:"usage"`
	Status, Attempts  int    `json:"-"`
	LatencyMS         int64  `json:"-"`
	Error, RetryAfter string `json:"-"`
}
type client struct {
	http      *http.Client
	transport Transport
	key       string
	q         Questions
}

// ask performs exactly one attempt; Run owns retries so every attempt is recorded
// and checked against the ceiling before another request can leave the machine.
func (c client) ask(ctx context.Context, pending PendingJudgment, attempt int) call {
	result := call{Attempts: attempt}
	id, question := c.q.Citation.QuestionID, c.q.Citation.Question
	if pending.Kind == "goal-mechanism" {
		id, question = c.q.Goal.QuestionID, c.q.Goal.Question
	}
	if pending.Kind == "source-grouping" {
		id, question = c.q.Grouping.QuestionID, c.q.Grouping.Question
	}
	var body bytes.Buffer
	encoder := json.NewEncoder(&body)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(struct {
		State     json.RawMessage     `json:"state"`
		Model     string              `json:"model"`
		Questions map[string]Question `json:"questions"`
	}{pending.state, c.transport.RequestModel, map[string]Question{id: question}}); err != nil {
		result.Error = err.Error()
		return result
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.transport.Endpoint, &body)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	request.Header.Set("Authorization", "Bearer "+c.key)
	request.Header.Set("Content-Type", "application/json")
	start := time.Now()
	response, err := c.http.Do(request)
	if err != nil {
		result.Error = err.Error()
		result.LatencyMS = time.Since(start).Milliseconds()
		return result
	}
	defer response.Body.Close()
	result.Status = response.StatusCode
	result.RetryAfter = response.Header.Get("Retry-After")
	if result.Status == http.StatusOK {
		data, readErr := io.ReadAll(io.LimitReader(response.Body, 1<<20+1))
		if readErr != nil {
			result.Error = readErr.Error()
		} else if len(data) > 1<<20 || json.Unmarshal(data, &result) != nil {
			result.Error = "unreadable answer"
		}
	}
	if result.Usage.InputTokens < 0 || result.Usage.OutputTokens < 0 {
		result.Usage.InputTokens, result.Usage.OutputTokens = 0, 0
		result.Error = "unreadable answer"
	}
	result.LatencyMS = time.Since(start).Milliseconds()
	return result
}

func retryDelay(header string, attempt int) time.Duration {
	if seconds, err := strconv.ParseFloat(header, 64); err == nil && seconds >= 0 {
		if seconds > 10 {
			seconds = 10
		}
		return time.Duration(seconds * float64(time.Second))
	}
	return time.Duration(attempt) * 1500 * time.Millisecond
}
func waitRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func (c call) cost(q Questions) (float64, string) {
	var cost float64
	if len(c.Usage.Cost) > 0 && string(c.Usage.Cost) != "null" && json.Unmarshal(c.Usage.Cost, &cost) == nil && cost >= 0 {
		return cost, "reported"
	}
	return float64(c.Usage.InputTokens) * q.USDPerMillionInputTokens / 1e6, "computed"
}
func evaluate(q Questions, p PendingJudgment, c call) (answer, string, string, bool) {
	if c.Error != "" {
		if c.Error == "unreadable answer" {
			return answer{}, "skipped", c.Error, false
		}
		return answer{}, "skipped", fmt.Sprintf("service unavailable (%s)", c.Error), true
	}
	if c.Status != 200 {
		reason, stop := fmt.Sprintf("request refused (HTTP %d)", c.Status), false
		if c.Status == 401 || c.Status == 402 || c.Status == 403 {
			reason, stop = fmt.Sprintf("key refused (HTTP %d)", c.Status), true
		} else if c.Status >= 500 || c.Status == 429 {
			reason, stop = fmt.Sprintf("service unavailable (HTTP %d)", c.Status), true
		}
		return answer{}, "skipped", reason, stop
	}
	id := q.Citation.QuestionID
	if p.Kind == "goal-mechanism" {
		id = q.Goal.QuestionID
	}
	if p.Kind == "source-grouping" {
		id = q.Grouping.QuestionID
	}
	a := c.Answers[id]
	if p.Kind == "source-grouping" {
		a = answer{Type: a.Type, Noul: a.Noul}
	}
	if !q.pinned(c.Model) {
		return a, "skipped", fmt.Sprintf("answered by %s, thresholds belong to %s", c.Model, q.PinnedModel), false
	}
	validProbability := func(v *float64) bool { return v != nil && *v >= 0 && *v <= 1 }
	advisory := false
	if p.Kind == "citation-support" {
		if a.Type != "choice" || a.Choice == nil || !validProbability(a.Confidence) || len(a.Probabilities) != len(q.Citation.Question.Criteria) {
			return answer{}, "skipped", "unreadable answer", false
		}
		if _, ok := q.Citation.Question.Criteria[*a.Choice]; !ok {
			return answer{}, "skipped", "unreadable answer", false
		}
		for name := range q.Citation.Question.Criteria {
			v, ok := a.Probabilities[name]
			if !ok || v < 0 || v > 1 {
				return answer{}, "skipped", "unreadable answer", false
			}
		}
		advisory = *a.Choice != q.Citation.RaiseWhenChoiceIsNot && *a.Confidence >= q.Citation.RaiseMinConfidence
	} else {
		if a.Type != "noul" || !validProbability(a.Noul) {
			return answer{}, "skipped", "unreadable answer", false
		}
		if p.Kind == "source-grouping" {
			if *a.Noul >= q.Grouping.SuggestWhenNoulAtLeast {
				return a, "suggested", "", false
			}
			return a, "clear", "", false
		}
		advisory = *a.Noul < q.Goal.RaiseWhenNoulBelow
	}
	if advisory {
		return a, "advisory", "", false
	}
	return a, "clear", "", false
}
