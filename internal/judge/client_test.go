package judge

// Suite: measured transport contract. Boundary: injected HTTP only; no sockets.
import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func fixtureResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Retry-After": []string{"0"}}, Body: io.NopCloser(strings.NewReader(body))}
}
func answerBody(model string, confidence, noul float64, cost string) string {
	usage := `"input_tokens":842,"output_tokens":20`
	if cost != "" {
		usage += `,"cost":` + cost
	}
	return `{"model":"` + model + `","id":"fixture-id","provider":"TypeSafe","answers":{"relation":{"type":"choice","choice":"says_nothing","probabilities":{"supports":0.02,"contradicts":0.05,"says_nothing":0.93},"confidence":` + floatText(confidence) + `},"delivers_goal":{"type":"noul","noul":` + floatText(noul) + `}},"usage":{` + usage + `}}`
}
func floatText(v float64) string { b, _ := json.Marshal(v); return string(b) }
func pendingFixture(t *testing.T) PendingJudgment {
	t.Helper()
	return fixturePlan(t, loadQuestions(t), englishContext+measuredClaim, "", englishContext, "prd").Pending[0]
}
func TestAskSendsTheMeasuredRequest(t *testing.T) {
	q := loadQuestions(t)
	p := pendingFixture(t)
	c := client{q: q, transport: q.Transports[0], key: "fixture-secret", http: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != "POST" || r.URL.String() != q.Transports[0].Endpoint || len(r.Header) != 2 || r.Header.Get("Authorization") != "Bearer fixture-secret" || r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("request contract: %s %s %v", r.Method, r.URL, r.Header)
		}
		remaining := time.Until(mustDeadline(t, r.Context()))
		if remaining > 30*time.Second || remaining < 29*time.Second {
			t.Fatalf("attempt deadline=%s", remaining)
		}
		var body struct {
			State     json.RawMessage     `json:"state"`
			Model     string              `json:"model"`
			Questions map[string]Question `json:"questions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if string(body.State) != string(p.state) || body.Model != q.Transports[0].RequestModel || len(body.Questions) != 1 || body.Questions[q.Citation.QuestionID].Instructions != q.Citation.Question.Instructions {
			t.Fatalf("body=%+v", body)
		}
		return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
	})}}
	got := c.ask(context.Background(), p, 1)
	if got.Status != 200 || got.Usage.InputTokens != 842 || got.Error != "" {
		t.Fatalf("call=%+v", got)
	}
}
func mustDeadline(t *testing.T, ctx context.Context) time.Time {
	t.Helper()
	d, ok := ctx.Deadline()
	if !ok {
		t.Fatal("no deadline")
	}
	return d
}
func TestAskSendsEachTransportsOwnModelID(t *testing.T) {
	q := loadQuestions(t)
	for _, tc := range []struct {
		name, model string
		index       int
	}{{"openrouter", "jev-1.13", 0}, {"typesafe", "jev-1.13.0", 1}} {
		t.Run(tc.name, func(t *testing.T) {
			c := client{q: q, transport: q.Transports[tc.index], key: "fake", http: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(r.Body)
				if !strings.Contains(string(body), `"model":"`+tc.model+`"`) {
					t.Fatalf("wrong request model: %s", body)
				}
				return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
			})}}
			c.ask(context.Background(), pendingFixture(t), 1)
		})
	}
}
func TestModelPinAcceptsOnlyJev113(t *testing.T) {
	q := loadQuestions(t)
	for _, tc := range []struct {
		model  string
		accept bool
	}{{"jev-1.13.0", true}, {"typesafe/jev-1.13", true}, {"typesafe/jev-1.13-20260917", true}, {"jev-1.14.0", false}, {"typesafe/jev-1.14-20261101", false}, {"~typesafe/jev-latest", false}, {"", false}} {
		t.Run(tc.model, func(t *testing.T) {
			if q.pinned(tc.model) != tc.accept {
				t.Fatalf("pin accepted=%v want=%v", q.pinned(tc.model), tc.accept)
			}
		})
	}
}
func TestAskRetriesARateLimitThenAnswers(t *testing.T) {
	q, req := runFixture(t)
	calls := 0
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls++
		if calls < 3 {
			return fixtureResponse(429, ""), nil
		}
		return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
	})
	req.Stage = "prd"
	got := runChecked(t, q, req)
	if calls != 3 || got.Calls != 3 || got.Judgments[0].Outcome != "advisory" || got.Stopped != nil {
		t.Fatalf("report=%+v calls=%d", got, calls)
	}
	if rows := readLog(t, req); len(rows) != 3 || rows[2].Attempts != 3 {
		t.Fatalf("retry records=%+v", rows)
	}
}
func TestRetryWaitsAndAttemptsHonorCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		if retryDelay("99", 1) != 10*time.Second || retryDelay("", 1) != 1500*time.Millisecond || retryDelay("", 2) != 3*time.Second {
			t.Fatal("retry delays")
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if waitRetry(ctx, 10*time.Second) == nil {
			t.Fatal("wait ignored cancellation")
		}
		q, req := runFixture(t)
		req.Stage = "prd"
		calls := 0
		req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
			calls++
			<-r.Context().Done()
			return nil, r.Context().Err()
		})
		got := runChecked(t, q, req)
		if calls != 1 || got.Stopped == nil || !strings.Contains(*got.Stopped, "deadline exceeded") {
			t.Fatalf("timeout=%+v", got)
		}
	})
}

func TestAskIgnoresUnknownFieldsAndRejectsMalformedValues(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		reason     string
	}{
		{"unknown metadata cannot override attempt", `{"model":"jev-1.13.0","status":503,"attempts":99,"error":"injected","answers":{"relation":{"type":"choice","choice":"says_nothing","probabilities":{"supports":0.02,"contradicts":0.05,"says_nothing":0.93},"confidence":0.8}}}`, ""},
		{"missing noul", `{"model":"jev-1.13.0","answers":{"delivers_goal":{"type":"noul"}}}`, "unreadable answer"},
		{"wrong primitive", `{"model":"jev-1.13.0","answers":{"relation":{"type":"noul","noul":0.1}}}`, "unreadable answer"},
		{"unknown choice", `{"model":"jev-1.13.0","answers":{"relation":{"type":"choice","choice":"unknown","confidence":0.8,"probabilities":{"supports":0.02,"contradicts":0.05,"says_nothing":0.93}}}}`, "unreadable answer"},
		{"negative tokens", `{"model":"jev-1.13.0","usage":{"input_tokens":-10}}`, "unreadable answer"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q, req := runFixture(t)
			req.Stage = "prd"
			if tc.name == "missing noul" {
				req.Stage = "techspec"
			}
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { return fixtureResponse(200, tc.body), nil })
			got := runChecked(t, q, req)
			if tc.reason == "" {
				if got.Judgments[0].Outcome != "advisory" || got.Stopped != nil || readLog(t, req)[0].Status != 200 || readLog(t, req)[0].Attempts != 1 {
					t.Fatalf("report=%+v", got)
				}
			} else if got.Judgments[0].Reason == nil || *got.Judgments[0].Reason != tc.reason || got.Stopped != nil {
				t.Fatalf("report=%+v", got)
			}
		})
	}
}
func TestRunNeverFollowsARedirectWithAKey(t *testing.T) {
	q, req := runFixture(t)
	calls := 0
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != q.Transports[0].Endpoint {
			t.Fatal("key left its endpoint")
		}
		resp := fixtureResponse(307, "")
		resp.Header.Set("Location", "https://example.org/steal")
		return resp, nil
	})
	got := runChecked(t, q, req)
	if calls != 2 || got.Calls != 2 {
		t.Fatalf("redirect calls=%d report=%+v", calls, got)
	}
}
