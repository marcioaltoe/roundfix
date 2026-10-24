package judge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/secretscan"
)

func followUpFixture(t *testing.T) (Questions, FollowUpRequest) {
	t.Helper()
	q := loadQuestions(t)
	return q, FollowUpRequest{
		Project: "roundfix", RepoRoot: t.TempDir(), HomeDir: t.TempDir(),
		Keys: map[string]string{q.Transports[0].KeyVariable: "fake-" + "judge-key"},
		Now:  func() time.Time { return time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC) },
		Items: []FollowUpItem{
			{Ref: "inbox/one.md", Kind: "inbox", Type: "fix", Origin: "other-project", Title: "Fix dispatch", Text: "Corrigir o despacho da fila"},
			{Ref: "backlog/two.md", Kind: "backlog", Type: "feat", Origin: "roundfix", Title: "Add export", Text: "Export the report"},
		},
	}
}
func followUpBody(model string) string {
	return `{"model":"` + model + `","answers":{"is_fix":{"type":"noul","noul":0.7},"blocks_other_project":{"type":"noul","noul":0.69},"is_small":{"type":"noul","noul":0.31},"affects_production":{"type":"noul","noul":0.3}},"usage":{"input_tokens":100,"output_tokens":10,"cost":0.01}}`
}
func followUpRun(t *testing.T, q Questions, req FollowUpRequest) FollowUpReport {
	t.Helper()
	got, err := PrioritizeFollowUp(context.Background(), q, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Answers) != len(req.Items) {
		t.Fatalf("answers=%d items=%d", len(got.Answers), len(req.Items))
	}
	for i, a := range got.Answers {
		if a.Ref != req.Items[i].Ref {
			t.Fatalf("answer order at %d: %q", i, a.Ref)
		}
	}
	return got
}
func followUpRows(t *testing.T, req FollowUpRequest) []LogLine {
	t.Helper()
	rows, err := ReadMonth(context.Background(), req.HomeDir, req.Now())
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

func TestFollowUpAsksFourQuestionsInOneRequest(t *testing.T) {
	t.Parallel()
	q, req := followUpFixture(t)
	calls := 0
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body struct {
			Model     string              `json:"model"`
			Questions map[string]Question `json:"questions"`
			State     struct {
				Project string            `json:"project"`
				Item    map[string]string `json:"item"`
			} `json:"state"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		item := req.Items[calls]
		want := map[string]string{"kind": item.Kind, "type": item.Type, "origin": item.Origin, "title": item.Title, "text": item.Text}
		if body.Model != q.Transports[0].RequestModel || body.State.Project != req.Project || !reflect.DeepEqual(body.State.Item, want) || !reflect.DeepEqual(body.Questions, q.FollowUp.Questions) || len(body.Questions) != 4 {
			t.Fatalf("body=%+v", body)
		}
		for _, question := range body.Questions {
			if question.Type != "noul" || len(question.Criteria) != 0 {
				t.Fatalf("question=%+v", question)
			}
		}
		calls++
		return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
	})
	got := followUpRun(t, q, req)
	if calls != 2 || got.Calls != 2 {
		t.Fatalf("calls=%d report=%+v", calls, got)
	}
	for _, a := range got.Answers {
		if a.Outcome != "answered" {
			t.Fatalf("answer=%+v", a)
		}
	}
}
func TestFollowUpVerdictsUseTheThresholds(t *testing.T) {
	t.Parallel()
	q, req := followUpFixture(t)
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
	})
	got := followUpRun(t, q, req)
	want := map[string]string{"is_fix": "yes", "blocks_other_project": "unsure", "is_small": "unsure", "affects_production": "no"}
	for _, a := range got.Answers {
		if a.Outcome != "answered" || !reflect.DeepEqual(a.Verdicts, want) {
			t.Fatalf("answer=%+v", a)
		}
	}
}
func TestFollowUpLogsTheModelOfEachAnswer(t *testing.T) {
	t.Parallel()
	q, req := followUpFixture(t)
	calls := 0
	models := []string{"jev-1.13.0", "jev-1.13.2"}
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		model := models[calls]
		calls++
		return fixtureResponse(200, followUpBody(model)), nil
	})
	got := followUpRun(t, q, req)
	rows := followUpRows(t, req)
	if len(rows) != 2 || got.Calls != 2 {
		t.Fatalf("rows=%+v report=%+v", rows, got)
	}
	want := map[string]float64{"is_fix": 0.7, "blocks_other_project": 0.69, "is_small": 0.31, "affects_production": 0.3}
	for i, row := range rows {
		if row.Model != models[i] || row.RequestedModel != q.Transports[0].RequestModel || row.Judgment != "follow-up-priority" || row.QuestionID != "follow_up_priority" || row.Artifact != req.Items[i].Ref || row.Target != req.Items[i].Kind || row.Spec != "" || row.Noul != nil || !reflect.DeepEqual(row.Probabilities, want) {
			t.Fatalf("row=%+v", row)
		}
		if got.Answers[i].Model == nil || *got.Answers[i].Model != models[i] {
			t.Fatalf("answer=%+v", got.Answers[i])
		}
	}
	if got.CostUSD != 0.02 || got.MonthCostUSD != 0.02 {
		t.Fatalf("cost=%+v", got)
	}
}
func TestFollowUpNeverSendsASecretShapedItem(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"text", "title", "origin", "kind", "type"} {
		t.Run(field, func(t *testing.T) {
			q, req := followUpFixture(t)
			secret := "ghp_" + strings.Repeat("a1", 18)
			matches := secretscan.ScanContent("fixture", []byte(secret))
			if len(matches) == 0 {
				t.Fatal("fixture does not match scanner")
			}
			switch field {
			case "text":
				req.Items[0].Text = secret
			case "title":
				req.Items[0].Title = secret
			case "origin":
				req.Items[0].Origin = secret
			case "kind":
				req.Items[0].Kind = secret
			case "type":
				req.Items[0].Type = secret
			}
			calls := 0
			req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(string(body), secret) {
					t.Fatal("secret in request")
				}
				calls++
				return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
			})
			got := followUpRun(t, q, req)
			if calls != 1 || got.Calls != 1 || got.Answers[0].Outcome != "skipped" || got.Answers[0].Reason == nil || *got.Answers[0].Reason != "not sent: secret pattern "+matches[0].Kind || got.Answers[1].Outcome != "answered" {
				t.Fatalf("report=%+v calls=%d", got, calls)
			}
			data, err := json.Marshal(struct {
				Report FollowUpReport
				Rows   []LogLine
			}{got, followUpRows(t, req)})
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), secret) {
				t.Fatal("secret in report or log")
			}
		})
	}
}
func TestFollowUpCutsTextOnARuneBoundary(t *testing.T) {
	t.Parallel()
	q, req := followUpFixture(t)
	req.Items = req.Items[:1]
	req.Items[0].Text = strings.Repeat("a", 3999) + "ção"
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body struct {
			State struct {
				Item struct {
					Text string `json:"text"`
				} `json:"item"`
			} `json:"state"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		text := body.State.Item.Text
		if utf8.RuneCountInString(text) != 4000 || !strings.HasSuffix(text, "ç") || strings.ContainsRune(text, utf8.RuneError) {
			t.Fatalf("invalid excerpt: runes=%d", utf8.RuneCountInString(text))
		}
		return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
	})
	if got := followUpRun(t, q, req); got.Calls != 1 {
		t.Fatalf("report=%+v", got)
	}
}
func TestFollowUpFailsOpen(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"missing key", "Run", "ceiling", "unreadable log", "503", "outside pin"} {
		t.Run(name, func(t *testing.T) {
			q, req := followUpFixture(t)
			calls := 0
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				if name == "503" {
					return fixtureResponse(503, ""), nil
				}
				return fixtureResponse(200, followUpBody("jev-1.14.0")), nil
			})
			switch name {
			case "missing key":
				req.Keys = nil
			case "Run":
				req.RunID = "fixture-run"
			case "ceiling", "unreadable log":
				path := newJudgeLog(req.HomeDir, req.Now()).path
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				data := fmt.Sprintf("{\"cost_usd\":%g}\n", q.MonthlyCeilingUSD)
				if name == "unreadable log" {
					data = "not json\n"
				}
				if err := os.WriteFile(path, []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			got := followUpRun(t, q, req)
			wantCalls := 0
			if name == "503" {
				wantCalls = 1
				if got.Stopped == nil || *got.Stopped != "service unavailable (HTTP 503)" {
					t.Fatalf("report=%+v", got)
				}
			} else if name == "outside pin" {
				wantCalls = 2
				if got.Stopped != nil || got.Skipped != nil {
					t.Fatalf("report=%+v", got)
				}
			} else if got.Skipped == nil {
				t.Fatalf("report=%+v", got)
			}
			if calls != wantCalls || got.Calls != wantCalls {
				t.Fatalf("calls=%d report=%+v", calls, got)
			}
			for _, a := range got.Answers {
				if a.Outcome != "skipped" || a.Reason == nil {
					t.Fatalf("answer=%+v", a)
				}
				if name == "outside pin" && *a.Reason != "answered by jev-1.14.0, thresholds belong to "+q.PinnedModel {
					t.Fatalf("reason=%s", *a.Reason)
				}
			}
			if wantCalls > 0 && len(followUpRows(t, req)) != wantCalls {
				t.Fatal("missing log rows")
			}
		})
	}
}
func TestFollowUpStopsAtTheItemLimit(t *testing.T) {
	t.Parallel()
	q, req := followUpFixture(t)
	req.Items = make([]FollowUpItem, 101)
	for i := range req.Items {
		req.Items[i] = FollowUpItem{Ref: fmt.Sprintf("item/%03d", i), Kind: "backlog", Text: "small change"}
	}
	calls := 0
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls++
		return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
	})
	got := followUpRun(t, q, req)
	if calls != 100 || got.Calls != 100 || got.Answers[100].Outcome != "skipped" || got.Answers[100].Reason == nil || *got.Answers[100].Reason != "item limit reached (100)" {
		t.Fatalf("calls=%d last=%+v", calls, got.Answers[100])
	}
	if len(followUpRows(t, req)) != 100 {
		t.Fatal("missing log rows")
	}
}

func TestFollowUpRetriesAndChecksEachAttempt(t *testing.T) {
	t.Parallel()
	for _, status := range []int{429, 529} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			q, req := followUpFixture(t)
			calls := 0
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				if calls < 3 {
					return fixtureResponse(status, ""), nil
				}
				return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
			})
			got := followUpRun(t, q, req)
			rows := followUpRows(t, req)
			if got.Stopped != nil || got.Calls != 4 || len(rows) != 4 || got.Answers[0].Outcome != "answered" || rows[0].Attempts != 1 || rows[1].Attempts != 2 || rows[2].Attempts != 3 {
				t.Fatalf("report=%+v rows=%+v", got, rows)
			}
		})
	}
	for _, name := range []string{"ceiling before retry", "ceiling before next item", "retry exhausted", "key refused", "key redaction"} {
		t.Run(name, func(t *testing.T) {
			q, req := followUpFixture(t)
			calls := 0
			fakeKey := req.Keys[q.Transports[0].KeyVariable]
			if strings.HasPrefix(name, "ceiling") {
				q = q.WithMonthlyCeiling(0.01)
			}
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				switch name {
				case "ceiling before retry":
					return fixtureResponse(429, `{"usage":{"cost":0.01}}`), nil
				case "retry exhausted":
					return fixtureResponse(529, ""), nil
				case "key refused":
					return fixtureResponse(401, ""), nil
				case "key redaction":
					return nil, fmt.Errorf("transport failed with %s", fakeKey)
				default:
					return fixtureResponse(200, followUpBody("jev-1.13.0")), nil
				}
			})
			// Non-200 bodies are intentionally not parsed by the shared client.
			// Simulate another caller consuming the budget during a retry.
			if name == "ceiling before retry" {
				req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
					calls++
					if err := AppendLogLine(req.HomeDir, req.Now(), LogLine{CostUSD: 0.01}); err != nil {
						t.Fatal(err)
					}
					return fixtureResponse(429, ""), nil
				})
			}
			got := followUpRun(t, q, req)
			want := 1
			if name == "retry exhausted" {
				want = 3
			}
			if calls != want || got.Stopped == nil || got.Answers[1].Outcome != "skipped" || got.Answers[1].Reason == nil {
				t.Fatalf("calls=%d report=%+v", calls, got)
			}
			if strings.HasPrefix(name, "ceiling") && !strings.Contains(*got.Stopped, "monthly ceiling reached") {
				t.Fatalf("stopped=%s", *got.Stopped)
			}
			if name == "key redaction" {
				data, err := json.Marshal(struct {
					Report FollowUpReport
					Rows   []LogLine
				}{got, followUpRows(t, req)})
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(string(data), fakeKey) || !strings.Contains(string(data), "[redacted]") {
					t.Fatal("key redaction failed")
				}
			}
		})
	}
}

func TestFollowUpRejectsUnreadableNouls(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"missing", "null", "negative", "above one", "wrong type", "broken json"} {
		t.Run(name, func(t *testing.T) {
			q, req := followUpFixture(t)
			body := followUpBody("jev-1.13.0")
			switch name {
			case "missing":
				body = strings.Replace(body, `"is_fix":{"type":"noul","noul":0.7},`, "", 1)
			case "null":
				body = strings.Replace(body, `"noul":0.7`, `"noul":null`, 1)
			case "negative":
				body = strings.Replace(body, `"noul":0.7`, `"noul":-0.1`, 1)
			case "above one":
				body = strings.Replace(body, `"noul":0.7`, `"noul":1.1`, 1)
			case "wrong type":
				body = strings.Replace(body, `"type":"noul"`, `"type":"choice"`, 1)
			case "broken json":
				body = "{"
			}
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { return fixtureResponse(200, body), nil })
			got := followUpRun(t, q, req)
			if got.Calls != 2 || got.Stopped != nil {
				t.Fatalf("report=%+v", got)
			}
			for _, a := range got.Answers {
				if a.Outcome != "skipped" || a.Reason == nil || *a.Reason != "unreadable answer" || a.Verdicts != nil {
					t.Fatalf("answer=%+v", a)
				}
			}
		})
	}
}
