package judge

// Suite: real temporary Judge Log IO, shared spend across transports, UTC month.
import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func readLog(t *testing.T, req Request) []logLine {
	t.Helper()
	b, err := os.ReadFile(newJudgeLog(req.HomeDir, req.Now()).path)
	if err != nil {
		t.Fatal(err)
	}
	var rows []logLine
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var row logLine
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}
	return rows
}
func TestJudgeLogRecordsEveryCall(t *testing.T) {
	q, req := runFixture(t)
	got := runChecked(t, q, req)
	rows := readLog(t, req)
	if len(rows) != got.Calls || len(rows) != 2 {
		t.Fatalf("rows=%d calls=%d", len(rows), got.Calls)
	}
	path := newJudgeLog(req.HomeDir, req.Now()).path
	for _, tc := range []struct {
		path string
		mode os.FileMode
	}{{path, 0600}, {filepath.Dir(path), 0700}} {
		info, err := os.Stat(tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != tc.mode {
			t.Fatalf("mode=%o want=%o", info.Mode().Perm(), tc.mode)
		}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), req.Keys[q.Transports[0].KeyVariable]) {
		t.Fatal("secret in log")
	}
	required := []string{"schema", "time", "repository", "spec", "judgment", "artifact", "line", "target", "state_hash", "question_id", "transport", "response_id", "provider", "requested_model", "model", "answer", "probabilities", "confidence", "noul", "latency_ms", "input_tokens", "output_tokens", "cost_usd", "cost_source", "status", "attempts", "error", "outcome"}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			t.Fatal(err)
		}
		if len(fields) != len(required) {
			t.Fatalf("fields=%v", fields)
		}
		for _, name := range required {
			if _, ok := fields[name]; !ok {
				t.Fatalf("missing %s", name)
			}
		}
	}
	for _, row := range rows {
		if row.Schema != "roundfix/judge-log/v1" || !row.Time.Equal(req.Now()) || row.Repository != req.RepoRoot || row.Spec != req.Spec || len(row.StateHash) != 16 || row.Status != 200 || row.Attempts != 1 || row.InputTokens != 842 || row.OutputTokens != 20 || row.Outcome != "advisory" || row.CostSource != "computed" || row.Transport != "openrouter" || row.RequestedModel != q.Transports[0].RequestModel {
			t.Fatalf("row=%+v", row)
		}
	}
	if rows[0].Answer == nil || rows[0].Confidence == nil || rows[0].Noul != nil || rows[1].Answer != nil || rows[1].Probabilities != nil || rows[1].Confidence != nil || rows[1].Noul == nil {
		t.Fatalf("answer fields=%+v", rows)
	}
}
func TestRunRecordsTheReportedCost(t *testing.T) {
	q, req := runFixture(t)
	req.Stage = "prd"
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return fixtureResponse(200, answerBody("typesafe/jev-1.13-20260917", 0.8, 0.29, "0.75")), nil
	})
	first := runChecked(t, q, req)
	rows := readLog(t, req)
	row := rows[0]
	if row.CostUSD != 0.75 || row.CostSource != "reported" || row.ResponseID != "fixture-id" || row.Provider != "TypeSafe" || row.Model != "typesafe/jev-1.13-20260917" || first.CostUSD != 0.75 {
		t.Fatalf("row=%+v report=%+v", row, first)
	}
	req.Keys = map[string]string{q.Transports[1].KeyVariable: "direct-secret"}
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		body := strings.ReplaceAll(answerBody("jev-1.13.0", 0.8, 0.29, ""), `,"id":"fixture-id","provider":"TypeSafe"`, "")
		return fixtureResponse(200, body), nil
	})
	second := runChecked(t, q, req)
	rows = readLog(t, req)
	computed := 842 * q.USDPerMillionInputTokens / 1e6
	if len(rows) != 2 || rows[1].CostUSD != computed || rows[1].CostSource != "computed" || rows[1].Transport != "typesafe" || rows[1].ResponseID != "" || rows[1].Provider != "" || second.MonthCostUSD != 0.75+computed {
		t.Fatalf("rows=%+v report=%+v", rows, second)
	}
	q.MonthlyCeilingUSD = second.MonthCostUSD
	req.Transport = neverRequest(t)
	third := runChecked(t, q, req)
	if third.Skipped == nil || third.MonthCostUSD != second.MonthCostUSD {
		t.Fatalf("summed ceiling=%+v", third)
	}
}
func TestJudgeLogUsesTheUTCMonthAndRejectsBadCosts(t *testing.T) {
	home := t.TempDir()
	now := time.Date(2026, 9, 30, 22, 0, 0, 0, time.FixedZone("UTC-3", -3*3600))
	log := newJudgeLog(home, now)
	if filepath.Base(log.path) != "2026-10.jsonl" {
		t.Fatalf("path=%s", log.path)
	}
	if cost, err := log.monthCost(); err != nil || cost != 0 {
		t.Fatalf("missing cost=%v err=%v", cost, err)
	}
	for _, tc := range []struct{ name, line string }{{"invalid JSON", "bad"}, {"null", "null"}, {"negative", `{"cost_usd":-1}`}, {"absent", `{}`}} {
		t.Run(tc.name, func(t *testing.T) {
			writeFixture(t, home, ".roundfix/judge/2026-10.jsonl", tc.line+"\n")
			if _, err := log.monthCost(); err == nil {
				t.Fatal("invalid log accepted")
			}
		})
	}
	for _, tc := range []struct {
		name, cost, source string
		want               float64
	}{{"negative reported", `-1`, "computed", 0.042}, {"null reported", `null`, "computed", 0.042}, {"string reported", `"0.1"`, "computed", 0.042}, {"zero reported", `0`, "reported", 0}} {
		t.Run(tc.name, func(t *testing.T) {
			c := call{}
			c.Usage.InputTokens = 1000000
			c.Usage.Cost = json.RawMessage(tc.cost)
			cost, source := c.cost(loadQuestions(t))
			if cost != tc.want || source != tc.source {
				t.Fatalf("cost=%v source=%s", cost, source)
			}
		})
	}
}
