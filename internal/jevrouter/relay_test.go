package jevrouter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func localRelay(t *testing.T, handler http.HandlerFunc) (*Relay, string, string) {
	t.Helper()
	upstream := httptest.NewServer(handler)
	t.Cleanup(upstream.Close)
	relay, err := StartRelay(upstream.URL+"/api/v1", upstream.Client())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := relay.Close(context.Background()); err != nil {
			t.Error(err)
		}
	})
	token, base := relay.Open()
	return relay, token, base
}
func relayPost(t *testing.T, base string) (int, string) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, base+"/chat/completions?raw=a%2Fb&x=1&x=2", strings.NewReader(`{"messages":["private fixture"]}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer relay-sentinel")
	req.Header.Set("Forwarded", "for=private")
	req.Header.Set("X-Forwarded-For", "private")
	req.Header.Set("X-Forwarded-Host", "private")
	req.Header.Set("X-Forwarded-Proto", "private")
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, string(body)
}

func TestRelayForwardsStreamedAndWholeAnswersUnchanged(t *testing.T) {
	t.Parallel()
	for _, stream := range []bool{false, true} {
		t.Run(fmt.Sprintf("stream=%t", stream), func(t *testing.T) {
			answer := "{\"id\":\"last\",\"model\":\"a/one\",\"provider\":\"p\",\"extra\":\"private\"}\n"
			if stream {
				answer = ": comment\r\ndata: " + strings.TrimSpace(answer) + "\r\n\r\ndata: [DONE]\n\n"
			}
			relay, token, base := localRelay(t, func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil || string(body) != `{"messages":["private fixture"]}` || r.Header.Get("Authorization") != "Bearer relay-sentinel" {
					t.Error("request body or authorization changed")
				}
				if r.URL.Path != "/api/v1/chat/completions" || r.URL.RawQuery != "raw=a%2Fb&x=1&x=2" {
					t.Error("path or query changed")
				}
				for _, header := range []string{"Forwarded", "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto"} {
					if r.Header.Get(header) != "" {
						t.Errorf("forwarding header %s", header)
					}
				}
				w.Header().Set("X-Upstream", "unchanged")
				if stream {
					w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
				} else {
					w.Header().Set("Content-Type", "application/json")
				}
				for _, b := range []byte(answer) {
					_, _ = w.Write([]byte{b})
					if stream {
						w.(http.Flusher).Flush()
					}
				}
			})
			status, body := relayPost(t, base)
			if status != 200 || body != answer {
				t.Fatalf("changed answer: status=%d body=%q", status, body)
			}
			got := relay.Take(token)
			if !reflect.DeepEqual(got, Observation{Models: []string{"a/one"}, Providers: []string{"p"}, ResponseID: "last"}) {
				t.Fatalf("observation=%+v", got)
			}
		})
	}
}
func TestRelayNotesModelsProvidersAndTheLastResponseID(t *testing.T) {
	t.Parallel()
	relay, token, base := localRelay(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"id\":\"first\",\"model\":\"a/one\",\"provider\":\"p\"}\n\ndata: malformed\n\ndata: {\"id\":\"second\",\"model\":\"b/two\",\"provider\":\"q\"}\n\ndata: {\"id\":\"last\",\"model\":\"a/one\",\"provider\":\"p\"}\n\ndata: [DONE]\n\n")
	})
	relayPost(t, base)
	want := Observation{Models: []string{"a/one", "b/two"}, Providers: []string{"p", "q"}, ResponseID: "last"}
	if got := relay.Take(token); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%+v want=%+v", got, want)
	}
	if got := relay.Take(token); !reflect.DeepEqual(got, Observation{}) {
		t.Fatalf("Take did not clear: %+v", got)
	}
}
func TestRelayNotesACreditRefusal(t *testing.T) {
	t.Parallel()
	for _, answer := range []string{`{"error":{"metadata":{"limit_source":"account"}}}`, "unparseable"} {
		t.Run(answer, func(t *testing.T) {
			calls := 0
			relay, token, base := localRelay(t, func(w http.ResponseWriter, _ *http.Request) {
				calls++
				w.WriteHeader(402)
				if calls == 1 {
					fmt.Fprint(w, answer)
				} else {
					fmt.Fprint(w, `{"error":{"metadata":{"limit_source":"second"}}}`)
				}
			})
			status, body := relayPost(t, base)
			if status != 402 || body != answer {
				t.Fatal("refusal changed")
			}
			relayPost(t, base)
			want := ""
			if strings.HasPrefix(answer, "{") {
				want = "account"
			}
			if got := relay.Take(token).Refusal; got == nil || got.Status != 402 || got.LimitSource != want {
				t.Fatalf("refusal=%+v", got)
			}
		})
	}
}
func TestRelayRefusesAPathWithoutAnOpenToken(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	relay, token, base := localRelay(t, func(w http.ResponseWriter, _ *http.Request) { calls.Add(1); fmt.Fprint(w, `{}`) })
	for _, url := range []string{relay.baseURL + "/unknown", relay.baseURL} {
		if status, _ := relayPost(t, url); status != 404 {
			t.Fatalf("status=%d", status)
		}
	}
	if relay.Release(token) != 0 {
		t.Fatal("token retained")
	}
	if status, _ := relayPost(t, base); status != 404 || calls.Load() != 0 {
		t.Fatal("unknown token reached upstream")
	}
}
func TestRelayCloseStopsServing(t *testing.T) {
	t.Parallel()
	relay, _, base := localRelay(t, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{}`) })
	if err := relay.Close(t.Context()); err != nil {
		t.Fatal(err)
	}
	response, err := http.Get(base)
	if err == nil {
		response.Body.Close()
		t.Fatal("closed relay still serves")
	}
	select {
	case <-relay.done:
	default:
		t.Fatal("serving goroutine not awaited")
	}
}
func TestRelayIgnoresOversizedAndMalformedBodies(t *testing.T) {
	t.Parallel()
	for _, answer := range []string{"broken", `{"model":"` + strings.Repeat("x", relayBodyLimit) + `"}`} {
		relay, token, base := localRelay(t, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, answer) })
		_, body := relayPost(t, base)
		if body != answer || !reflect.DeepEqual(relay.Take(token), Observation{}) {
			t.Fatal("unparseable answer changed or observed")
		}
	}
}
func TestRelayUpstreamFailureIs502(t *testing.T) {
	t.Parallel()
	upstream := httptest.NewServer(http.NotFoundHandler())
	upstream.Close()
	relay, err := StartRelay(upstream.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer relay.Close(context.Background())
	_, base := relay.Open()
	if status, _ := relayPost(t, base); status != 502 {
		t.Fatalf("status=%d", status)
	}
}

func TestRelayFlushesBeforeUpstreamFinishes(t *testing.T) {
	t.Parallel()
	finish := make(chan struct{})
	var once sync.Once
	release := func() { once.Do(func() { close(finish) }) }
	// Release before cleanup closes the relay, even if a client read fails.
	defer release()
	relay, token, base := localRelay(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("X-Upstream", "unchanged")
		fmt.Fprint(w, "data: {\"model\":\"a/one\",\"id\":\"first\"}\n\n")
		w.(http.Flusher).Flush()
		select {
		case <-finish:
		case <-r.Context().Done():
			return
		}
		fmt.Fprint(w, "data: {\"model\":\"b/two\",\"id\":\"last\"}\n\n")
	})
	response, err := (&http.Client{Timeout: 5 * time.Second}).Get(base + "/chat/completions")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.Header.Get("X-Upstream") != "unchanged" {
		t.Fatal("response header changed")
	}
	first := "data: {\"model\":\"a/one\",\"id\":\"first\"}\n\n"
	data := make([]byte, len(first))
	if _, err := io.ReadFull(response.Body, data); err != nil || string(data) != first {
		t.Fatalf("stream did not flush: %q %v", data, err)
	}
	release()
	if _, err := io.Copy(io.Discard, response.Body); err != nil {
		t.Fatal(err)
	}
	if got := relay.Take(token); !reflect.DeepEqual(got.Models, []string{"a/one", "b/two"}) || got.ResponseID != "last" {
		t.Fatalf("observation=%+v", got)
	}
}

func TestRelayKeepsTokensAndTakesIndependent(t *testing.T) {
	t.Parallel()
	relay, first, base := localRelay(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/api/v1/escaped%2Fpath" {
			t.Errorf("escaped path changed: %s", r.URL.EscapedPath())
		}
		fmt.Fprint(w, `{"model":"a/one","error":"ignored extra shape"}`)
	})
	second, _ := relay.Open()
	response, err := http.Get(base + "/escaped%2Fpath")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if got := relay.Take(second); !reflect.DeepEqual(got, Observation{}) {
		t.Fatal("observation leaked between tokens")
	}
	if got := relay.Take(first); !reflect.DeepEqual(got.Models, []string{"a/one"}) {
		t.Fatalf("observation=%+v", got)
	}
	if relay.Release(first) != 1 {
		t.Fatal("release removed another token")
	}
}
