package jevrouter

// Suite: local HTTP boundary only, explicit sentinel credentials and temporary Home.
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestReadCreditsSendsOnlyTheBearerHeader(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/credits" || r.URL.RawQuery != "" {
			t.Errorf("request=%s %s", r.Method, r.URL)
		}
		if len(r.Header) != 1 || r.Header.Get("Authorization") != "Bearer "+testKey {
			t.Errorf("unexpected request headers: %v", r.Header)
		}
		fmt.Fprint(w, `{"data":{"total_credits":1,"total_usage":0.5}}`)
	}))
	defer server.Close()
	client := server.Client()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	jar.SetCookies(endpoint, []*http.Cookie{{Name: "private", Value: testKey}})
	client.Jar = jar
	usage, err := ReadCredits(context.Background(), client, server.URL+"/api/v1/", testKey)
	if err != nil || usage.Balance() != 0.5 || calls != 1 {
		t.Fatalf("usage=%v calls=%d err=%v", usage, calls, err)
	}
	if client.CheckRedirect != nil || client.Timeout != 0 || client.Jar != jar {
		t.Fatal("caller client changed")
	}
}

func TestReadCreditsRefusesMalformedAnswersAndHTTPFailures(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bodies := []string{`{}`, `null`, `{"data":null}`, `{"data":{}}`, testKey, strings.Repeat("x", (1<<20)+1), `{"data":{"total_credits":1,"total_usage":0}} trailing`}
	for _, field := range []string{"total_credits", "total_usage"} {
		other := "total_usage"
		if field == other {
			other = "total_credits"
		}
		for _, value := range []string{"null", "-1", "1e999", `"NaN"`, `"Infinity"`, `"0.5"`, "true", "[]", "{}"} {
			bodies = append(bodies, fmt.Sprintf(`{"data":{"%s":%s,"%s":0}}`, field, value, other))
		}
		bodies = append(bodies, fmt.Sprintf(`{"data":{"%s":0}}`, other))
	}
	for i, body := range bodies {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, body) }))
			defer server.Close()
			_, err := ReadCredits(context.Background(), server.Client(), server.URL, testKey)
			if err == nil || strings.Contains(err.Error(), testKey) || !strings.HasPrefix(err.Error(), "read account credits: ") {
				t.Fatalf("unsafe or absent error: %v", err)
			}
		})
	}
	testReadCreditsHTTPFailures(t)
}

func testReadCreditsHTTPFailures(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	for _, status := range []int{302, 401, 402, 403, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls++
				w.Header().Set("Location", "/redirected")
				w.WriteHeader(status)
				fmt.Fprint(w, testKey)
			}))
			defer server.Close()
			_, err := ReadCredits(context.Background(), server.Client(), server.URL, testKey)
			if err == nil || strings.Contains(err.Error(), testKey) || calls != 1 {
				t.Fatalf("calls=%d error=%v", calls, err)
			}
		})
	}
}

func TestReadCreditsHonorsCancellationAndNetworkFailure(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"data":{"total_credits":0,"total_usage":0}}`)
	}))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ReadCredits(ctx, server.Client(), server.URL, testKey); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error=%v", err)
	}
	server.Close()
	if _, err := ReadCredits(context.Background(), server.Client(), server.URL+"/"+testKey, testKey); err == nil || strings.Contains(err.Error(), testKey) {
		t.Fatalf("network error=%v", err)
	}
	// A response that never completes must stop at the caller's earlier deadline.
	blocked := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer blocked.Close()
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := ReadCredits(ctx, blocked.Client(), blocked.URL, testKey); !errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), testKey) {
		t.Fatalf("deadline error=%v", err)
	}
}

func TestReadCreditsCapsTheRequestAtTenSeconds(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := ReadCredits(ctx, server.Client(), server.URL, testKey)
	if !errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil || strings.Contains(err.Error(), testKey) {
		t.Fatalf("request did not enforce its own deadline: %v (parent=%v)", err, ctx.Err())
	}
}

func TestCreditLeftTakesTheLowerOfBalanceAndKeyLimit(t *testing.T) {
	for _, tc := range []struct {
		name      string
		remaining *float64
		want      float64
	}{
		{"absent", nil, 20}, {"lower", creditPointer(3), 3}, {"higher", creditPointer(40.31), 20}, {"zero", creditPointer(0), 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := CreditLeft(Credits{TotalCredits: 30, TotalUsage: 10}, KeyStatus{LimitRemaining: tc.remaining}); got != tc.want {
				t.Fatalf("left=%v want=%v", got, tc.want)
			}
		})
	}
	if got := (Credits{TotalCredits: 1, TotalUsage: 2}).Balance(); got != -1 {
		t.Fatalf("balance=%v", got)
	}
}
func creditPointer(value float64) *float64 { return &value }
func TestMinCreditDefaultsToFifteenDollars(t *testing.T) {
	if DefaultMinCreditUSD != 15 {
		t.Fatal("default floor changed")
	}
	for _, value := range []float64{0, -1, 4, 20} {
		want := DefaultMinCreditUSD
		if value > 0 {
			want = value
		}
		if got := MinCredit(value); got != want {
			t.Fatalf("MinCredit(%v)=%v want=%v", value, got, want)
		}
	}
}
