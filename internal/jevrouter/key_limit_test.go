package jevrouter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"roundfix/internal/judge"
)

func checkKeyLimit(t *testing.T, fields, reason string) {
	t.Helper()
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/key" || r.Header.Get("Authorization") != "Bearer "+testKey {
			t.Errorf("unexpected key request: %s %s", r.Method, r.URL.Path)
		}
		fmt.Fprintf(w, `{"data":{"usage_monthly":0.5,%s}}`, fields)
	}))
	defer server.Close()
	spend, err := MonthSpend(context.Background(), Deps{HomeDir: t.TempDir(), Env: []string{"ROUNDFIX_OPENROUTER_API_KEY=" + testKey}, Client: server.Client(), Endpoint: server.URL}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	err = spend.CheckKeyLimit()
	if reason == "" {
		if err != nil {
			t.Fatalf("bounded key refused: %v", err)
		}
	} else if err == nil || !strings.HasPrefix(err.Error(), reason+":") || strings.Contains(err.Error(), testKey) {
		t.Fatalf("wanted %s, got %v", reason, err)
	} else if reason == "jev_router_key_unbounded" && !strings.Contains(err.Error(), fmt.Sprintf("monthly credit limit of at most US$%.4f", questions.MonthlyCeilingUSD)) {
		t.Fatalf("missing maintainer remedy: %v", err)
	}
}

func TestARoutedPromptNeedsAMonthlyKeyLimit(t *testing.T) {
	for _, fields := range []string{`"limit":null,"limit_remaining":null,"limit_reset":null`, `"limit_reset":"monthly"`, `"limit":null,"limit_reset":"monthly"`} {
		t.Run(fields, func(t *testing.T) { checkKeyLimit(t, fields, "jev_router_key_unbounded") })
	}
}
func TestALifetimeKeyLimitIsRefused(t *testing.T) {
	for _, reset := range []string{`null`, `"daily"`, `"weekly"`, `"Monthly"`} {
		t.Run(reset, func(t *testing.T) { checkKeyLimit(t, `"limit":1,"limit_reset":`+reset, "jev_router_key_unbounded") })
	}
}
func TestAMonthlyKeyLimitAboveTheCeilingIsRefused(t *testing.T) {
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	checkKeyLimit(t, fmt.Sprintf(`"limit":%v,"limit_reset":"monthly"`, questions.MonthlyCeilingUSD+0.01), "jev_router_key_unbounded")
}
func TestAMonthlyKeyLimitWithinTheCeilingIsAccepted(t *testing.T) {
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, limit := range []float64{1, questions.MonthlyCeilingUSD} {
		t.Run(fmt.Sprint(limit), func(t *testing.T) {
			checkKeyLimit(t, fmt.Sprintf(`"limit":%v,"limit_remaining":0.5,"limit_reset":"monthly"`, limit), "")
		})
	}
	checkKeyLimit(t, `"limit":1,"limit_remaining":null,"limit_reset":"monthly"`, "")
}
func TestAnExhaustedKeyLimitIsRefused(t *testing.T) {
	for _, remaining := range []float64{0, -0.1} {
		t.Run(fmt.Sprint(remaining), func(t *testing.T) {
			checkKeyLimit(t, fmt.Sprintf(`"limit":1,"limit_remaining":%v,"limit_reset":"monthly"`, remaining), "jev_ceiling_reached")
		})
	}
}
