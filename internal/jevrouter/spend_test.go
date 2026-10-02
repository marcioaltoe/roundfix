package jevrouter

// Suite: shared Judge Log, UTC month boundaries and a local key endpoint.
import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/judge"
)

func TestMonthSpendAddsTypeSafeToTheLargerOpenRouterFigure(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	now := time.Date(2026, 9, 30, 22, 0, 0, 0, time.FixedZone("UTC-3", -3*3600))
	for _, row := range []judge.LogLine{{Transport: "typesafe", CostUSD: 0.10}, {Transport: "openrouter", CostUSD: 0.20}} {
		if err := judge.AppendLogLine(home, now, row); err != nil {
			t.Fatal(err)
		}
	}
	if err := judge.AppendLogLine(home, now.AddDate(0, -1, 0), judge.LogLine{Transport: "typesafe", CostUSD: 100}); err != nil {
		t.Fatal(err)
	}
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ usage, total float64 }{{0.50, 0.60}, {0.05, 0.30}, {0, 0.30}} {
		t.Run(fmt.Sprint(tc.usage), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprintf(w, `{"data":{"usage_monthly":%v}}`, tc.usage)
			}))
			defer server.Close()
			spend, err := MonthSpend(context.Background(), Deps{HomeDir: home, Env: []string{"ROUNDFIX_OPENROUTER_API_KEY=" + testKey}, Client: server.Client(), Endpoint: server.URL}, now)
			if err != nil || math.Abs(spend.Total-tc.total) > 1e-12 || spend.TypeSafeLogged != 0.10 || spend.OpenRouterLogged != 0.20 || spend.KeyUsageMonthly != tc.usage || spend.Ceiling != questions.MonthlyCeilingUSD {
				t.Fatalf("spend=%+v error=%v", spend, err)
			}
		})
	}
}

func TestMonthSpendFailsOnAnUnreadableKeyEndpoint(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, status := range []int{200, 500} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			fmt.Fprint(w, testKey)
		}))
		_, err := MonthSpend(context.Background(), Deps{HomeDir: home, Env: []string{"ROUNDFIX_OPENROUTER_API_KEY=" + testKey}, Client: server.Client(), Endpoint: server.URL}, time.Now())
		server.Close()
		if err == nil || strings.Contains(err.Error(), testKey) {
			t.Fatalf("error=%v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	server.Close()
	_, err := MonthSpend(context.Background(), Deps{HomeDir: home, Env: []string{"ROUNDFIX_OPENROUTER_API_KEY=" + testKey}, Client: server.Client(), Endpoint: server.URL}, time.Now())
	if err == nil || strings.Contains(err.Error(), testKey) {
		t.Fatalf("network error=%v", err)
	}
}

func TestMonthSpendRejectsUnreadableLogsAndMissingKeys(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	now := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		fmt.Fprint(w, `{"data":{"usage_monthly":0.5}}`)
	}))
	defer server.Close()
	deps := Deps{HomeDir: home, Env: []string{"ROUNDFIX_OPENROUTER_API_KEY=" + testKey}, Client: server.Client(), Endpoint: server.URL}
	spend, err := MonthSpend(context.Background(), deps, now)
	if err != nil || spend.Total != 0.5 {
		t.Fatalf("missing log spend=%+v err=%v", spend, err)
	}
	deps.Env = nil
	if _, err := MonthSpend(context.Background(), deps, now); err == nil || calls != 1 {
		t.Fatalf("missing key calls=%d err=%v", calls, err)
	}
	deps.Env = []string{"ROUNDFIX_OPENROUTER_API_KEY=" + testKey}
	path := filepath.Join(home, ".roundfix", "judge", "2026-10.jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{"bad", `{"cost_usd":-1}`, `{"cost_usd":null}`} {
		if err := os.WriteFile(path, []byte(line+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := MonthSpend(context.Background(), deps, now); err == nil || calls != 1 {
			t.Fatalf("invalid log calls=%d err=%v", calls, err)
		}
	}
}
