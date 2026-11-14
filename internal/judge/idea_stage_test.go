package judge

import (
	"context"
	"crypto/sha256"
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

func ideaStageFixture(t *testing.T) (Questions, IdeaStageRequest) {
	t.Helper()
	q := loadQuestions(t)
	return q, IdeaStageRequest{Project: "roundfix", RepoRoot: t.TempDir(), HomeDir: t.TempDir(),
		Keys: map[string]string{q.Transports[0].KeyVariable: strings.Join([]string{"fake", "judge", "key"}, "-")},
		Now:  func() time.Time { return time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC) },
		Item: FollowUpItem{Ref: "docs/backlog/feature.md", Kind: "backlog", Type: "feat", Origin: "roundfix", Title: "Share Run results", Text: "Talvez uma página, talvez um bot."},
	}
}

func ideaStageBody(model string, values ...float64) string {
	return fmt.Sprintf(`{"model":%q,"answers":{"shape_open":{"type":"noul","noul":%g},"unverified_assumption":{"type":"noul","noul":%g},"market_question":{"type":"noul","noul":%g}},"usage":{"input_tokens":512,"output_tokens":10,"cost":0.000021504}}`, model, values[0], values[1], values[2])
}

func ideaStageRun(t *testing.T, q Questions, req IdeaStageRequest) IdeaStageReport {
	t.Helper()
	got, err := JudgeIdeaStage(context.Background(), q, req)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestIdeaStageAsksThreeCriteriaInOneRequest(t *testing.T) {
	t.Parallel()
	for _, index := range []int{0, 1} {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			q, req := ideaStageFixture(t)
			req.Keys = map[string]string{q.Transports[index].KeyVariable: strings.Join([]string{"fake", "transport", "key"}, "-")}
			calls := 0
			req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				var body struct {
					Model     string                     `json:"model"`
					Questions map[string]Question        `json:"questions"`
					State     map[string]json.RawMessage `json:"state"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if len(body.State) != 2 || string(body.State["project"]) != `"roundfix"` || body.Model != q.Transports[index].RequestModel {
					t.Fatalf("body=%+v", body)
				}
				var item map[string]string
				if err := json.Unmarshal(body.State["item"], &item); err != nil {
					t.Fatal(err)
				}
				want := map[string]string{"kind": req.Item.Kind, "type": req.Item.Type, "origin": req.Item.Origin, "title": req.Item.Title, "text": req.Item.Text}
				if !reflect.DeepEqual(item, want) || len(body.Questions) != 3 {
					t.Fatalf("item=%v questions=%v", item, body.Questions)
				}
				for _, criterion := range q.IdeaStage.Criteria {
					if !reflect.DeepEqual(body.Questions[criterion.ID], criterion.Question) {
						t.Fatalf("question=%+v", body.Questions)
					}
				}
				if r.Header.Get("Authorization") != "Bearer "+req.Keys[q.Transports[index].KeyVariable] {
					t.Fatal("wrong command key")
				}
				return fixtureResponse(200, ideaStageBody("jev-1.13.0", .98, .84, .17)), nil
			})
			got := ideaStageRun(t, q, req)
			if calls != 1 || got.Calls != 1 || got.State != "answered" || got.Verdict != "indicated" || got.InputTokens != 512 {
				t.Fatalf("calls=%d report=%+v", calls, got)
			}
		})
	}
}

func TestIdeaStageVerdictUsesTheThresholds(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		values  []float64
		answers []string
		verdict string
	}{
		{[]float64{.70, .30, .30}, []string{"yes", "no", "no"}, "indicated"},
		{[]float64{.69, .30, .30}, []string{"unsure", "no", "no"}, "undecided"},
		{[]float64{.30, .30, .30}, []string{"no", "no", "no"}, "not indicated"},
		{[]float64{.31, .30, .30}, []string{"unsure", "no", "no"}, "undecided"},
		{[]float64{1, 1, 1}, []string{"yes", "yes", "yes"}, "indicated"},
		{[]float64{0, 0, .70}, []string{"no", "no", "yes"}, "indicated"},
	} {
		t.Run(fmt.Sprint(tc.values), func(t *testing.T) {
			q, req := ideaStageFixture(t)
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return fixtureResponse(200, ideaStageBody("jev-1.13.0", tc.values...)), nil
			})
			got := ideaStageRun(t, q, req)
			if got.State != "answered" || got.Verdict != tc.verdict || IdeaStageVerdict(tc.answers) != tc.verdict {
				t.Fatalf("report=%+v", got)
			}
			for i, a := range got.Answers {
				if a.ID != q.IdeaStage.Criteria[i].ID || a.Answer != tc.answers[i] || a.Probability != tc.values[i] {
					t.Fatalf("answer=%+v", a)
				}
			}
		})
	}
	// The thresholds are data, rather than constants hidden in the evaluator.
	q, req := ideaStageFixture(t)
	q.IdeaStage.YesAtLeast, q.IdeaStage.NoAtMost = .8, .2
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return fixtureResponse(200, ideaStageBody("jev-1.13.0", .7, .3, .1)), nil
	})
	got := ideaStageRun(t, q, req)
	if got.Verdict != "undecided" || got.Answers[0].Answer != "unsure" || got.Answers[1].Answer != "unsure" {
		t.Fatalf("report=%+v", got)
	}
}

func TestIdeaStageSendsNothingItMustNot(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"fix", "perf", "refactor", "finding", "Run", "no key", "ceiling", "unreadable log", "secret text", "secret title", "secret origin", "secret type", "secret kind"} {
		t.Run(name, func(t *testing.T) {
			q, req := ideaStageFixture(t)
			calls := 0
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return fixtureResponse(500, ""), nil })
			wantState, wantReason := "judge_unavailable", ""
			switch name {
			case "fix", "perf", "refactor", "finding":
				req.Item.Type = name
				wantState, wantReason = "not_applicable", "the item is typed "+name+"; only a feature is judged for the Idea stage"
			case "Run":
				req.RunID = "fixture-run"
				wantReason = "inside Run fixture-run, named by ROUNDFIX_RUN_ID"
			case "no key":
				req.Keys = nil
				wantReason = "ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set (nor ROUNDFIX_OPENROUTER_API_KEY, nor ROUNDFIX_TYPESAFE_API_KEY)"
			case "ceiling", "unreadable log":
				path := newJudgeLog(req.HomeDir, req.Now()).path
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				data := fmt.Sprintf("{\"cost_usd\":%g}\n", q.MonthlyCeilingUSD)
				if name == "unreadable log" {
					data = "broken json\n"
				}
				if err := os.WriteFile(path, []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
				if name == "ceiling" {
					wantReason = fmt.Sprintf("monthly ceiling reached (US$%.4f of US$%.2f)", q.MonthlyCeilingUSD, q.MonthlyCeilingUSD)
				} else {
					_, err := newJudgeLog(req.HomeDir, req.Now()).monthCost()
					wantReason = "judge log unreadable: " + err.Error()
				}
			default:
				secret := "ghp_" + strings.Repeat("a1", 18)
				matches := secretscan.ScanContent("fixture", []byte(secret))
				if len(matches) == 0 {
					t.Fatal("bad secret fixture")
				}
				wantReason = "not sent: secret pattern " + matches[0].Kind
				switch name {
				case "secret text":
					req.Item.Text = secret
				case "secret title":
					req.Item.Title = secret
				case "secret origin":
					req.Item.Origin = secret
				case "secret type":
					req.Item.Type = secret
				case "secret kind":
					req.Item.Kind = secret
				}
			}
			got := ideaStageRun(t, q, req)
			if calls != 0 || got.Calls != 0 || got.State != wantState || got.Reason == nil || *got.Reason != wantReason || got.Verdict != "" || len(got.Answers) != 0 || got.Model != nil || got.Transport != nil || got.KeyVariable != nil {
				t.Fatalf("calls=%d report=%+v reason=%v want=%s", calls, got, got.Reason, wantReason)
			}
		})
	}
}

func TestIdeaStageFailsOpen(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"outside pin", "broken json", "missing", "null", "negative", "above one", "wrong type", "500", "401", "refused", "redaction", "log not writable"} {
		t.Run(name, func(t *testing.T) {
			q, req := ideaStageFixture(t)
			calls := 0
			want := "unreadable answer"
			body := ideaStageBody("jev-1.13.0", .9, .2, .1)
			status := 200
			switch name {
			case "outside pin":
				body = ideaStageBody("jev-1.14.0", .9, .2, .1)
				want = "answered by jev-1.14.0, thresholds belong to " + q.PinnedModel
			case "broken json":
				body = "{"
			case "missing":
				body = strings.Replace(body, `"shape_open":{"type":"noul","noul":0.9},`, "", 1)
			case "null":
				body = strings.Replace(body, `"noul":0.9`, `"noul":null`, 1)
			case "negative":
				body = strings.Replace(body, `"noul":0.9`, `"noul":-0.1`, 1)
			case "above one":
				body = strings.Replace(body, `"noul":0.9`, `"noul":1.1`, 1)
			case "wrong type":
				body = strings.Replace(body, `"type":"noul"`, `"type":"choice"`, 1)
			case "500":
				status = 500
				want = "service unavailable (HTTP 500)"
			case "401":
				status = 401
				want = "key refused (HTTP 401)"
			case "refused":
				status = 400
				want = "request refused (HTTP 400)"
			case "redaction":
				want = "service unavailable (Post \"" + q.Transports[0].Endpoint + "\": failed with [redacted])"
			case "log not writable":
				want = "judge log not writable: "
			}
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				if name == "redaction" {
					return nil, fmt.Errorf("failed with %s", req.Keys[q.Transports[0].KeyVariable])
				}
				if name == "log not writable" {
					path := newJudgeLog(req.HomeDir, req.Now()).path
					if err := os.MkdirAll(path, 0700); err != nil {
						t.Fatal(err)
					}
				}
				return fixtureResponse(status, body), nil
			})
			got := ideaStageRun(t, q, req)
			if calls != 1 || got.Calls != 1 || got.State != "judge_unavailable" || got.Verdict != "" || len(got.Answers) != 0 || got.Reason == nil {
				t.Fatalf("report=%+v", got)
			}
			if name == "log not writable" {
				if !strings.HasPrefix(*got.Reason, want) {
					t.Fatalf("reason=%s", *got.Reason)
				}
				return
			}
			if *got.Reason != want {
				t.Fatalf("reason=%s want=%s", *got.Reason, want)
			}
			rows, err := ReadMonth(context.Background(), req.HomeDir, req.Now())
			if err != nil {
				t.Fatal(err)
			}
			if len(rows) != 1 || rows[0].Outcome != "skipped" || rows[0].Judgment != "idea-stage" {
				t.Fatalf("rows=%+v", rows)
			}
			data, _ := json.Marshal(struct {
				Report IdeaStageReport
				Rows   []LogLine
			}{got, rows})
			if strings.Contains(string(data), req.Keys[q.Transports[0].KeyVariable]) {
				t.Fatal("key in report or log")
			}
		})
	}
}

func TestIdeaStageCutsTextOnARuneBoundary(t *testing.T) {
	t.Parallel()
	q, req := ideaStageFixture(t)
	req.Item.Text = strings.Repeat("a", 3999) + "ção"
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body struct {
			State struct {
				Item followUpStateItem `json:"item"`
			} `json:"state"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		text := body.State.Item.Text
		if utf8.RuneCountInString(text) != 4000 || !strings.HasSuffix(text, "ç") || !utf8.ValidString(text) {
			t.Fatal("invalid excerpt")
		}
		return fixtureResponse(200, ideaStageBody("jev-1.13.0", .9, .2, .1)), nil
	})
	if got := ideaStageRun(t, q, req); got.State != "answered" {
		t.Fatalf("report=%+v", got)
	}
}

func TestIdeaStageLogsItsJudgment(t *testing.T) {
	t.Parallel()
	q, req := ideaStageFixture(t)
	var sent []byte
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body struct {
			State json.RawMessage `json:"state"`
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		sent = body.State
		return fixtureResponse(200, ideaStageBody("jev-1.13.0", .98, .84, .17)), nil
	})
	got := ideaStageRun(t, q, req)
	rows, err := ReadMonth(context.Background(), req.HomeDir, req.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows=%+v", rows)
	}
	row := rows[0]
	hash := sha256.Sum256(sent)
	want := map[string]float64{"shape_open": .98, "unverified_assumption": .84, "market_question": .17}
	if row.Schema != "roundfix/judge-log/v1" || row.Judgment != "idea-stage" || row.Artifact != req.Item.Ref || row.Target != req.Item.Kind || row.QuestionID != "idea_stage" || row.Spec != "" || row.Model != "jev-1.13.0" || row.RequestedModel != q.Transports[0].RequestModel || row.StateHash != fmt.Sprintf("%x", hash)[:16] || !reflect.DeepEqual(row.Probabilities, want) || row.Outcome != "answered" || !row.Time.Equal(req.Now()) || row.Attempts != 1 || row.InputTokens != 512 || got.CostUSD != .000021504 || got.MonthCostUSD != got.CostUSD {
		t.Fatalf("row=%+v report=%+v", row, got)
	}
	data, _ := json.Marshal(row)
	if strings.Contains(string(data), req.Item.Text) || strings.Contains(string(data), req.Item.Title) {
		t.Fatal("item contents in log")
	}
}

func TestTheIdeaStageQuestionsAreTheMeasuredOnes(t *testing.T) {
	t.Parallel()
	q := loadQuestions(t)
	var want IdeaStageJudgment
	if err := json.Unmarshal([]byte("{\n     \"question_id\": \"idea_stage\",\n     \"criteria\": [\n       {\n         \"id\": \"shape_open\",\n         \"summary\": \"the right shape of the solution is still open\",\n         \"question\": {\n           \"type\": \"noul\",\n           \"instructions\": \"`item` is one pending request addressed to the software project `project`. Is the right shape of the solution for `item` still open?\",\n           \"criteria\": {\n             \"true\": \"At least two materially different designs could satisfy the item, and the item does not settle which one to build.\",\n             \"false\": \"The item names one concrete design, or only one reasonable design exists.\"\n           }\n         }\n       },\n       {\n         \"id\": \"unverified_assumption\",\n         \"summary\": \"it rests on an assumption about its users, the business or the market that it gives no evidence for\",\n         \"question\": {\n           \"type\": \"noul\",\n           \"instructions\": \"`item` is one pending request addressed to the software project `project`. Does `item` rest on an assumption about its users, the business or the market that it gives no evidence for?\",\n           \"criteria\": {\n             \"true\": \"A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong.\",\n             \"false\": \"The need is shown by evidence in the item itself, such as a measured defect, a reported failure or a decision the product owner already made.\"\n           }\n         }\n       },\n       {\n         \"id\": \"market_question\",\n         \"summary\": \"deciding to build it depends on a market or competitive question\",\n         \"question\": {\n           \"type\": \"noul\",\n           \"instructions\": \"`item` is one pending request addressed to the software project `project`. Does deciding whether to build `item` depend on a market or competitive question?\",\n           \"criteria\": {\n             \"true\": \"The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it.\",\n             \"false\": \"The decision depends only on the product's own behavior and its current users.\"\n           }\n         }\n       }\n     ],\n     \"yes_at_least\": 0.7,\n     \"no_at_most\": 0.3,\n     \"item_max_chars\": 4000,\n     \"not_judged_types\": [\n       \"fix\",\n       \"perf\",\n       \"refactor\",\n       \"finding\"\n     ]\n   }"), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(q.IdeaStage, want) {
		t.Fatalf("measured judgment changed: %+v", q.IdeaStage)
	}
}

func TestIdeaStageRetriesUnderTheCeiling(t *testing.T) {
	t.Parallel()
	for _, status := range []int{429, 529} {
		for _, mode := range []string{"answers", "exhausted", "ceiling", "cancelled"} {
			t.Run(fmt.Sprintf("%d/%s", status, mode), func(t *testing.T) {
				q, req := ideaStageFixture(t)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				calls := 0
				req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
					calls++
					if mode == "answers" && calls == 3 {
						return fixtureResponse(200, ideaStageBody("jev-1.13.0", .8, .2, .1)), nil
					}
					if mode == "ceiling" {
						if err := AppendLogLine(req.HomeDir, req.Now(), LogLine{CostUSD: q.MonthlyCeilingUSD}); err != nil {
							t.Fatal(err)
						}
					}
					if mode == "cancelled" {
						cancel()
					}
					response := fixtureResponse(status, "")
					response.Header.Set("Retry-After", "0")
					return response, nil
				})
				got, err := JudgeIdeaStage(ctx, q, req)
				if err != nil {
					t.Fatal(err)
				}
				wantCalls := 3
				wantState := "judge_unavailable"
				if mode == "answers" {
					wantState = "answered"
				}
				if mode == "ceiling" || mode == "cancelled" {
					wantCalls = 1
				}
				if got.State != wantState || got.Calls != wantCalls || calls != wantCalls {
					t.Fatalf("calls=%d report=%+v", calls, got)
				}
				if mode == "exhausted" && (got.Reason == nil || *got.Reason != fmt.Sprintf("service unavailable (HTTP %d)", status)) {
					t.Fatalf("report=%+v", got)
				}
				if mode == "ceiling" && (got.Reason == nil || !strings.HasPrefix(*got.Reason, "monthly ceiling reached")) {
					t.Fatalf("report=%+v", got)
				}
				rows, err := ReadMonth(context.Background(), req.HomeDir, req.Now())
				if err != nil {
					t.Fatal(err)
				}
				if mode != "ceiling" {
					if len(rows) != wantCalls {
						t.Fatalf("rows=%+v", rows)
					}
					for i, row := range rows {
						if row.Attempts != i+1 {
							t.Fatalf("row=%+v", row)
						}
					}
				}
			})
		}
	}
}
