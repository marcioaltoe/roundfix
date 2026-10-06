package judge

// Suite: fail-open execution. Every environment and writable root is injected.
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func runFixture(t *testing.T) (Questions, Request) {
	t.Helper()
	root := t.TempDir()
	spec := filepath.Join(root, "docs/specs/0300-example")
	writeFixture(t, spec, "_prd.md", englishContext+measuredClaim+"\n\n## Goals\n- The author reads the evidence for the work.\n")
	writeFixture(t, spec, "_techspec.md", englishContext+"## Coverage Map\n- Goal 1 → Reader mechanism\n## Reader mechanism\n"+strings.Repeat(englishContext, 5))
	writeFixture(t, root, "docs/adr/0123-decision.md", "---\nstatus: accepted\n---\n"+englishContext)
	req := Request{RepoRoot: root, SpecDir: spec, Spec: "0300-example", HomeDir: t.TempDir(), Keys: map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "fake-openrouter-secret"}, Now: func() time.Time { return time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC) }, Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
	})}
	return loadQuestions(t), req
}
func runChecked(t *testing.T, q Questions, req Request) Report {
	t.Helper()
	got, err := Run(context.Background(), q, req)
	if err != nil {
		t.Fatal(err)
	}
	return got
}
func TestRunRaisesAtTheMeasuredThresholds(t *testing.T) {
	for _, tc := range []struct {
		name             string
		confidence, noul float64
		citation, goal   string
	}{{"at confidence floor and below probability ceiling", 0.8, 0.29, "advisory", "advisory"}, {"below confidence floor and at probability ceiling", 0.79, 0.3, "clear", "clear"}} {
		t.Run(tc.name, func(t *testing.T) {
			q, req := runFixture(t)
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return fixtureResponse(200, answerBody("jev-1.13.0", tc.confidence, tc.noul, "")), nil
			})
			got := runChecked(t, q, req)
			if len(got.Judgments) != 2 || got.Judgments[0].Outcome != tc.citation || got.Judgments[1].Outcome != tc.goal {
				t.Fatalf("judgments=%+v", got.Judgments)
			}
		})
	}
}
func TestRunNeverComparesAnotherModelsAnswer(t *testing.T) {
	for _, model := range []string{"jev-1.13.0", "typesafe/jev-1.13-20260917", "jev-1.14.0"} {
		t.Run(model, func(t *testing.T) {
			q, req := runFixture(t)
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return fixtureResponse(200, answerBody(model, 0.8, 0.29, "")), nil
			})
			got := runChecked(t, q, req)
			for _, j := range got.Judgments {
				if model == "jev-1.14.0" {
					if j.Outcome != "skipped" || j.Reason == nil || *j.Reason != "answered by jev-1.14.0, thresholds belong to jev-1.13" {
						t.Fatalf("judgment=%+v", j)
					}
				} else if j.Outcome != "advisory" {
					t.Fatalf("judgment=%+v", j)
				}
			}
			for _, line := range readLog(t, req) {
				if line.Model != model {
					t.Fatalf("logged model=%q", line.Model)
				}
			}
		})
	}
}
func TestRunSendsNothingWithoutAKey(t *testing.T) {
	q, req := runFixture(t)
	req.Keys = nil
	req.Transport = neverRequest(t)
	got := runChecked(t, q, req)
	assertNoKey(t, got)
}
func neverRequest(t *testing.T) http.RoundTripper {
	t.Helper()
	return roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unexpected request")
		return nil, errors.New("unexpected request")
	})
}
func assertNoKey(t *testing.T, got Report) {
	t.Helper()
	if got.Skipped == nil || *got.Skipped != "ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set (nor ROUNDFIX_OPENROUTER_API_KEY, nor ROUNDFIX_TYPESAFE_API_KEY)" || got.Calls != 0 || len(got.Judgments) != 2 {
		t.Fatalf("report=%+v", got)
	}
	for _, j := range got.Judgments {
		if j.Outcome != "skipped" {
			t.Fatalf("judgment=%+v", j)
		}
	}
}
func TestRunStopsAtTheMonthlyCeiling(t *testing.T) {
	t.Run("already exactly at ceiling", func(t *testing.T) {
		q, req := runFixture(t)
		writeFixture(t, req.HomeDir, ".roundfix/judge/2026-10.jsonl", fmt.Sprintf("{\"cost_usd\":%g}\n", q.MonthlyCeilingUSD))
		req.Transport = neverRequest(t)
		got := runChecked(t, q, req)
		if got.Skipped == nil || *got.Skipped != fmt.Sprintf("monthly ceiling reached (US$%.4f of US$%.2f)", q.MonthlyCeilingUSD, q.MonthlyCeilingUSD) {
			t.Fatalf("report=%+v", got)
		}
	})
	t.Run("run reaches ceiling", func(t *testing.T) {
		q, req := runFixture(t)
		calls := 0
		req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, fmt.Sprint(q.MonthlyCeilingUSD))), nil
		})
		got := runChecked(t, q, req)
		if calls != 1 || got.Stopped == nil || got.Judgments[1].Outcome != "skipped" {
			t.Fatalf("report=%+v calls=%d", got, calls)
		}
	})
}
func TestRunSendsNothingWithAnUnreadableLog(t *testing.T) {
	q, req := runFixture(t)
	writeFixture(t, req.HomeDir, ".roundfix/judge/2026-10.jsonl", "broken\n")
	req.Transport = neverRequest(t)
	got := runChecked(t, q, req)
	if got.Skipped == nil || !strings.HasPrefix(*got.Skipped, "judge log unreadable: ") {
		t.Fatalf("report=%+v", got)
	}
}
func TestRunStopsWhenTheServiceFails(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		reason string
	}{{"unavailable", 503, "service unavailable (HTTP 503)"}, {"unauthorized", 401, "key refused (HTTP 401)"}, {"no credits", 402, "key refused (HTTP 402)"}, {"forbidden", 403, "key refused (HTTP 403)"}, {"rate limit", 429, "service unavailable (HTTP 429)"}, {"overloaded", 529, "service unavailable (HTTP 529)"}} {
		t.Run(tc.name, func(t *testing.T) {
			q, req := runFixture(t)
			calls := 0
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return fixtureResponse(tc.status, ""), nil })
			got := runChecked(t, q, req)
			want := 1
			if tc.status == 429 || tc.status == 529 {
				want = 3
			}
			if calls != want || got.Stopped == nil || *got.Stopped != tc.reason || got.Judgments[1].Outcome != "skipped" {
				t.Fatalf("report=%+v calls=%d", got, calls)
			}
		})
	}
}
func TestRunSkipsARefusedRequestAndContinues(t *testing.T) {
	q, req := runFixture(t)
	calls := 0
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return fixtureResponse(422, ""), nil
		}
		return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
	})
	got := runChecked(t, q, req)
	if calls != 2 || got.Stopped != nil || got.Judgments[0].Reason == nil || *got.Judgments[0].Reason != "request refused (HTTP 422)" || got.Judgments[1].Outcome != "advisory" {
		t.Fatalf("report=%+v", got)
	}
}
func TestRunSelectsTheTransportByKey(t *testing.T) {
	for _, tc := range []struct {
		name  string
		keys  map[string]string
		index int
	}{{"openrouter", map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "router-secret"}, 0}, {"typesafe", map[string]string{"ROUNDFIX_TYPESAFE_API_KEY": "direct-secret"}, 1}, {"both", map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "router-secret", "ROUNDFIX_TYPESAFE_API_KEY": "direct-secret"}, 0}} {
		t.Run(tc.name, func(t *testing.T) {
			q, req := runFixture(t)
			req.Keys = tc.keys
			calls := 0
			req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				expected := q.Transports[tc.index]
				if r.URL.String() != expected.Endpoint || r.Header.Get("Authorization") != "Bearer "+tc.keys[expected.KeyVariable] {
					t.Fatalf("wrong transport: %s %v", r.URL, r.Header)
				}
				body, _ := io.ReadAll(r.Body)
				if tc.index == 0 && (strings.Contains(string(body), "direct-secret") || strings.Contains(r.Header.Get("Authorization"), "direct-secret")) {
					t.Fatal("direct key leaked")
				}
				return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
			})
			got := runChecked(t, q, req)
			if calls != 2 || got.Transport == nil || *got.Transport != tc.name && tc.name != "both" {
				t.Fatalf("report=%+v calls=%d", got, calls)
			}
		})
	}
}
func TestRunIgnoresTheGenericOpenRouterKey(t *testing.T) {
	for _, tc := range []struct {
		name string
		keys map[string]string
	}{{"only OPENROUTER_API_KEY", map[string]string{"OPENROUTER_API_KEY": "generic-router"}}, {"only TYPESAFE_API_KEY", map[string]string{"TYPESAFE_API_KEY": "generic-direct"}}, {"both generic keys", map[string]string{"OPENROUTER_API_KEY": "generic-router", "TYPESAFE_API_KEY": "generic-direct"}}} {
		t.Run(tc.name, func(t *testing.T) {
			q, req := runFixture(t)
			req.Keys = tc.keys
			req.Transport = neverRequest(t)
			assertNoKey(t, runChecked(t, q, req))
		})
	}
}
func TestRequestsCarryOnlySpecArtifactText(t *testing.T) {
	for _, index := range []int{0, 1} {
		t.Run([]string{"openrouter", "typesafe"}[index], func(t *testing.T) {
			q, req := runFixture(t)
			writeFixture(t, req.RepoRoot, "source.go", "SOURCE_SENTINEL")
			writeFixture(t, req.RepoRoot, "docs/findings/finding.md", "FINDING_SENTINEL")
			// Canonical source retains the literal phrase substituted by the planner.
			writeFixture(t, req.SpecDir, "_prd.md", englishContext+"The cited decision keeps the verification gate in the repository for every author.\n\nADR-0123 keeps the verification gate in the repository for every author.\n\n## Goals\n- The author reads the evidence for the work.\n")
			var sources []string
			for _, p := range []string{filepath.Join(req.SpecDir, "_prd.md"), filepath.Join(req.SpecDir, "_techspec.md"), filepath.Join(req.RepoRoot, "docs/adr/0123-decision.md")} {
				b, err := os.ReadFile(p)
				if err != nil {
					t.Fatal(err)
				}
				sources = append(sources, string(b))
			}
			corpus := strings.Join(sources, "\n")
			req.Keys = map[string]string{q.Transports[index].KeyVariable: "boundary-secret"}
			calls := 0
			req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				b, _ := io.ReadAll(r.Body)
				for _, forbidden := range []string{"SOURCE_SENTINEL", "FINDING_SENTINEL", "boundary-secret"} {
					if strings.Contains(string(b), forbidden) {
						t.Fatalf("forbidden request text %s", forbidden)
					}
				}
				if len(r.Header) != 2 || r.Header.Get("Authorization") == "" || r.Header.Get("Content-Type") != "application/json" {
					t.Fatalf("headers=%v", r.Header)
				}
				var body struct {
					State map[string]string `json:"state"`
				}
				if err := json.Unmarshal(b, &body); err != nil {
					t.Fatal(err)
				}
				if len(body.State) == 0 {
					t.Fatal("missing state")
				}
				for _, text := range body.State {
					for _, part := range strings.Split(strings.ReplaceAll(text, "The cited decision", "the cited decision"), "the cited decision") {
						if !strings.Contains(corpus, strings.TrimSpace(part)) {
							t.Fatalf("state outside sources: %q", part)
						}
					}
				}
				return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
			})
			runChecked(t, q, req)
			if calls != 2 {
				t.Fatalf("calls=%d", calls)
			}
		})
	}
}
func TestRunHandlesMalformedAnswersAndWriteFailures(t *testing.T) {
	for _, body := range []string{"invalid", `{"model":"jev-1.13.0","answers":{}}`, `{"model":"jev-1.13.0","answers":{"relation":{"type":"choice","choice":"says_nothing","confidence":null}}}`} {
		t.Run(body, func(t *testing.T) {
			q, req := runFixture(t)
			req.Stage = "prd"
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { return fixtureResponse(200, body), nil })
			got := runChecked(t, q, req)
			if got.Judgments[0].Reason == nil || *got.Judgments[0].Reason != "unreadable answer" || got.Stopped != nil {
				t.Fatalf("report=%+v", got)
			}
		})
	}
	t.Run("append fails after answer", func(t *testing.T) {
		q, req := runFixture(t)
		// Induce the write failure after the initial log read.
		calls := 0
		req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			writeFixture(t, req.HomeDir, ".roundfix", "regular file blocks append")
			return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
		})
		got := runChecked(t, q, req)
		if calls != 1 || got.Stopped == nil || !strings.HasPrefix(*got.Stopped, "judge log not writable: ") || got.Judgments[0].Outcome != "advisory" {
			t.Fatalf("report=%+v", got)
		}
	})
	t.Run("only PRD read is an error", func(t *testing.T) {
		q, req := runFixture(t)
		if err := os.Remove(filepath.Join(req.SpecDir, "_prd.md")); err != nil {
			t.Fatal(err)
		}
		req.Transport = neverRequest(t)
		if _, err := Run(context.Background(), q, req); err == nil {
			t.Fatal("missing PRD accepted")
		}
	})
	t.Run("network error redacted", func(t *testing.T) {
		q, req := runFixture(t)
		req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("failed fake-openrouter-secret") })
		got := runChecked(t, q, req)
		if got.Stopped == nil || strings.Contains(*got.Stopped, "fake-openrouter-secret") {
			t.Fatalf("report=%+v", got)
		}
		b, err := os.ReadFile(newJudgeLog(req.HomeDir, req.Now()).path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), "fake-openrouter-secret") {
			t.Fatal("key in log")
		}
	})
}

func TestRunRechecksSpendBeforeRetryAndHonorsCancellation(t *testing.T) {
	t.Run("retry reaches ceiling", func(t *testing.T) {
		q, req := runFixture(t)
		calls := 0
		req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			writeFixture(t, req.HomeDir, ".roundfix/judge/2026-10.jsonl", fmt.Sprintf("{\"cost_usd\":%g}\n", q.MonthlyCeilingUSD))
			return fixtureResponse(429, ""), nil
		})
		got := runChecked(t, q, req)
		if calls != 1 || got.Stopped == nil || *got.Stopped != fmt.Sprintf("monthly ceiling reached (US$%.4f of US$%.2f)", q.MonthlyCeilingUSD, q.MonthlyCeilingUSD) {
			t.Fatalf("report=%+v calls=%d", got, calls)
		}
	})
	t.Run("cancelled retry wait", func(t *testing.T) {
		q, req := runFixture(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		calls := 0
		req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			cancel()
			response := fixtureResponse(529, "")
			response.Header.Set("Retry-After", "10")
			return response, nil
		})
		got, err := Run(ctx, q, req)
		if err != nil {
			t.Fatal(err)
		}
		if calls != 1 || got.Stopped == nil || *got.Stopped != "service unavailable (context canceled)" {
			t.Fatalf("report=%+v calls=%d", got, calls)
		}
	})
	t.Run("already cancelled", func(t *testing.T) {
		q, req := runFixture(t)
		req.Transport = neverRequest(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		got, err := Run(ctx, q, req)
		if err != nil {
			t.Fatal(err)
		}
		if got.Calls != 0 || got.Stopped == nil {
			t.Fatalf("report=%+v", got)
		}
	})
}
func TestRunAsksIdenticalStatesOnceAndKeepsTheSelectedTransport(t *testing.T) {
	q, req := runFixture(t)
	writeFixture(t, req.SpecDir, "_techspec.md", englishContext+measuredClaim+"\n\n"+measuredClaim)
	calls := 0
	req.Keys["ROUNDFIX_TYPESAFE_API_KEY"] = "direct-secret"
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != q.Transports[0].Endpoint || r.Header.Get("Authorization") != "Bearer fake-openrouter-secret" {
			t.Fatalf("transport changed: %s %v", r.URL, r.Header)
		}
		delete(req.Keys, "ROUNDFIX_OPENROUTER_API_KEY")
		return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
	})
	got := runChecked(t, q, req)
	if calls != 1 || got.Calls != 1 || len(got.Judgments) != 1 {
		t.Fatalf("deduplicated report=%+v calls=%d", got, calls)
	}
	// Distinct states continue using the selection made before the first call.
	q, req = runFixture(t)
	calls = 0
	req.Keys["ROUNDFIX_TYPESAFE_API_KEY"] = "direct-secret"
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != q.Transports[0].Endpoint || r.Header.Get("Authorization") != "Bearer fake-openrouter-secret" {
			t.Fatal("transport switched during run")
		}
		delete(req.Keys, "ROUNDFIX_OPENROUTER_API_KEY")
		return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
	})
	got = runChecked(t, q, req)
	if calls != 2 || got.Calls != 2 {
		t.Fatalf("report=%+v calls=%d", got, calls)
	}
}

func TestRunReportsSkippedArtifactPathsWithoutSendingText(t *testing.T) {
	q, req := runFixture(t)
	writeFixture(t, req.SpecDir, "_prd.md", "A decisão é uma regra para os autores e não está na sua documentação.\n")
	req.Stage = "prd"
	req.Transport = neverRequest(t)
	got := runChecked(t, q, req)
	if got.Calls != 0 || len(got.ArtifactsSkipped) != 1 || got.ArtifactsSkipped[0].Artifact != "docs/specs/0300-example/_prd.md" || got.ArtifactsSkipped[0].Reason != "not English" {
		t.Fatalf("report=%+v", got)
	}
}

func TestWithMonthlyCeilingReplacesOnlyAPositiveCeiling(t *testing.T) {
	q := loadQuestions(t)
	for _, value := range []float64{50, 0, -1} {
		got := q.WithMonthlyCeiling(value)
		want := q
		if value > 0 {
			want.MonthlyCeilingUSD = value
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("override %v changed other settings: %+v", value, got)
		}
	}
	if q.MonthlyCeilingUSD != loadQuestions(t).MonthlyCeilingUSD {
		t.Fatal("original questions changed")
	}
}

func TestRunStopsAtAConfiguredCeiling(t *testing.T) {
	for _, spend := range []float64{49, 50} {
		t.Run(fmt.Sprint(spend), func(t *testing.T) {
			q, req := runFixture(t)
			q = q.WithMonthlyCeiling(50)
			writeFixture(t, req.HomeDir, ".roundfix/judge/2026-10.jsonl", fmt.Sprintf("{\"cost_usd\":%g}\n", spend))
			calls := 0
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				return fixtureResponse(200, answerBody("jev-1.13.0", 0.8, 0.29, "")), nil
			})
			got := runChecked(t, q, req)
			if spend < 50 {
				if calls != 2 || got.Skipped != nil {
					t.Fatalf("below ceiling: %+v calls=%d", got, calls)
				}
			} else if calls != 0 || got.Skipped == nil || *got.Skipped != "monthly ceiling reached (US$50.0000 of US$50.00)" {
				t.Fatalf("at ceiling: %+v calls=%d", got, calls)
			}
		})
	}
}
