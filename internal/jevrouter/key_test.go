package jevrouter

// Suite: local HTTP boundary only, explicit sentinel credentials and temporary Home.
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testKey = "router-test-secret"

func TestKeyUsageSendsOnlyTheBearerHeader(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/key" || r.URL.RawQuery != "" {
			t.Errorf("request=%s %s", r.Method, r.URL)
		}
		if len(r.Header) != 1 || r.Header.Get("Authorization") != "Bearer "+testKey {
			t.Errorf("unexpected request headers: %v", r.Header)
		}
		fmt.Fprint(w, `{"data":{"usage_monthly":0.5}}`)
	}))
	defer server.Close()
	client := server.Client()
	usage, err := KeyUsage(context.Background(), client, server.URL+"/api/v1/", testKey)
	if err != nil || usage != 0.5 || calls != 1 {
		t.Fatalf("usage=%v calls=%d err=%v", usage, calls, err)
	}
	if client.CheckRedirect != nil || client.Timeout != 0 {
		t.Fatal("caller client changed")
	}
}

func TestKeyUsageRefusesAMalformedAnswer(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	for _, body := range []string{`{}`, `null`, `{"data":null}`, `{"data":{"usage_monthly":null}}`, `{"data":{"usage_monthly":-1}}`, `{"data":{"usage_monthly":"0.5"}}`, `{"data":{"usage_monthly":true}}`, `{"data":{"usage_monthly":1e999}}`, `{"data":{"usage_monthly":0.5}} trailing`, testKey} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, body) }))
			defer server.Close()
			_, err := KeyUsage(context.Background(), server.Client(), server.URL, testKey)
			if err == nil || strings.Contains(err.Error(), testKey) {
				t.Fatalf("unsafe or absent error: %v", err)
			}
		})
	}
}

func TestKeyUsageRefusesHTTPFailuresAndRedirects(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	for _, status := range []int{302, 401, 402, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls++
				w.Header().Set("Location", "/redirected")
				w.WriteHeader(status)
				fmt.Fprint(w, testKey)
			}))
			defer server.Close()
			_, err := KeyUsage(context.Background(), server.Client(), server.URL, testKey)
			if err == nil || strings.Contains(err.Error(), testKey) || calls != 1 {
				t.Fatalf("calls=%d error=%v", calls, err)
			}
		})
	}
}

func TestKeyUsageHonorsCancellationAndNetworkFailure(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{"data":{"usage_monthly":0}}`) }))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := KeyUsage(ctx, server.Client(), server.URL, testKey); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error=%v", err)
	}
	server.Close()
	if _, err := KeyUsage(context.Background(), server.Client(), server.URL+"/"+testKey, testKey); err == nil || strings.Contains(err.Error(), testKey) {
		t.Fatalf("network error=%v", err)
	}
	// A response that never completes must stop at the caller's earlier deadline.
	blocked := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer blocked.Close()
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := KeyUsage(ctx, blocked.Client(), blocked.URL, testKey); !errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), testKey) {
		t.Fatalf("deadline error=%v", err)
	}
}

func TestKeyUsageCapsTheRequestAtTenSeconds(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := KeyUsage(ctx, server.Client(), server.URL, testKey)
	if !errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil || strings.Contains(err.Error(), testKey) {
		t.Fatalf("request did not enforce its own deadline: %v (parent=%v)", err, ctx.Err())
	}
}
