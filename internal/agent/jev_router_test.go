package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/jevrouter"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Session environment tests cross the acpx process boundary using the compiled
// test binary. OpenCode and network calls remain outside this boundary.
const jevRouterFixtureConfigPath = "ROUNDFIX_FAKE_ACPX_JEV_CONFIG_PATH"
const jevRouterSentinel = "task01-sentinel-key-never-in-config"

type jevRouterFixtureEnvironment struct {
	Config     string
	KeyPresent bool
}

func recordJevRouterFixtureEnvironment() error {
	path := os.Getenv(jevRouterFixtureConfigPath)
	if path == "" {
		return nil
	}
	record := jevRouterFixtureEnvironment{Config: os.Getenv("OPENCODE_CONFIG_CONTENT"), KeyPresent: os.Getenv(JevRouterKeyEnv) == jevRouterSentinel}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal router environment: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("record router environment: %w", err)
	}
	return nil
}

func runJevRouterPrompt(t *testing.T, harness *fakeACPXHarness, model string) error {
	t.Helper()
	_, err := harness.runner.RunPrompt(t.Context(), ACPXPromptRequest{
		ExecuteRequest: ExecuteRequest{Runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: model}, GitRoot: harness.gitRoot, Prompt: "prompt"},
		Session:        "router-test",
	}, nil)
	return err
}

func TestJevRouterSessionCarriesThePlaceholderConfig(t *testing.T) {
	t.Parallel()
	harness := newFakeACPXHarness(t)
	t.Cleanup(func() { harness.runner.clearSessionState("router-test") })
	path := filepath.Join(harness.gitRoot, "config.json")
	harness.setEnv(jevRouterFixtureConfigPath, path)
	harness.setEnv(JevRouterKeyEnv, jevRouterSentinel)
	harness.setEnv("OPENCODE_CONFIG_CONTENT", `{"inherited":true}`)
	if err := runJevRouterPrompt(t, harness, JevRouterModel); err != nil {
		t.Fatal(err)
	}
	var record jevRouterFixtureEnvironment
	data := readFile(t, path)
	if err := json.Unmarshal([]byte(data), &record); err != nil {
		t.Fatal(err)
	}
	if record.Config != jevRouterProviderConfig(routerFixtureBaseURL(t, record.Config)) || !record.KeyPresent {
		t.Fatalf("routed child did not receive placeholder config and environment key")
	}
	if !json.Valid([]byte(record.Config)) {
		t.Fatal("invalid provider JSON")
	}
	if strings.Contains(record.Config, jevRouterSentinel) {
		t.Fatal("key in configuration")
	}
	acpxConfig := filepath.Join(environmentValue(harness.runner.Environment, "HOME"), ".acpx", "config.json")
	if strings.Contains(readFile(t, acpxConfig), jevRouterSentinel) {
		t.Fatal("key in acpx configuration")
	}
	calls := readJSONInvocations(t, harness.invocationsPath)
	if len(calls) == 0 {
		t.Fatal("no acpx calls")
	}
	for _, args := range calls {
		for _, arg := range args {
			if strings.Contains(arg, jevRouterSentinel) {
				t.Fatal("key in acpx argument")
			}
		}
	}
	if got := environmentValue(harness.runner.Environment, "OPENCODE_CONFIG_CONTENT"); got != `{"inherited":true}` {
		t.Fatal("runner base environment mutated")
	}
	// The same runner retains inherited configuration for an ordinary session.
	if err := runJevRouterPrompt(t, harness, "opencode-go/kimi-k3"); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(readFile(t, path)), &record); err != nil {
		t.Fatal(err)
	}
	if record.Config != `{"inherited":true}` {
		t.Fatal("router override leaked to another model")
	}
}

func TestJevRouterRefusedWithoutTheKey(t *testing.T) {
	t.Parallel()
	for _, present := range []bool{false, true} {
		t.Run(fmt.Sprintf("empty-variable-present=%t", present), func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			base := harness.runner.Environment[:0]
			for _, entry := range harness.runner.Environment {
				key, _, _ := strings.Cut(entry, "=")
				if key != JevRouterKeyEnv {
					base = append(base, entry)
				}
			}
			harness.runner.Environment = base
			if present {
				harness.setEnv(JevRouterKeyEnv, "")
			}
			err := runJevRouterPrompt(t, harness, JevRouterModel)
			var failure *SelectionFailureError
			if !errors.As(err, &failure) || failure.Runtime != "opencode" || failure.Reason != JevRouterKeyMissing+": "+JevRouterKeyEnv+" is not set" {
				t.Fatalf("unexpected refusal: %v", err)
			}
			assertNoFile(t, "invocations", harness.invocationsPath)
		})
	}
}

func TestOtherOpenCodeModelsGetNoInlineConfig(t *testing.T) {
	t.Parallel()
	for _, inherited := range []string{"", `{"user":true}`} {
		t.Run(fmt.Sprintf("inherited=%s", inherited), func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			path := filepath.Join(harness.gitRoot, "config.json")
			harness.setEnv(jevRouterFixtureConfigPath, path)
			harness.setEnv("OPENCODE_CONFIG_CONTENT", inherited)
			harness.setEnv(JevRouterKeyEnv, "")
			if err := runJevRouterPrompt(t, harness, "opencode-go/kimi-k3"); err != nil {
				t.Fatal(err)
			}
			var record jevRouterFixtureEnvironment
			if err := json.Unmarshal([]byte(readFile(t, path)), &record); err != nil {
				t.Fatal(err)
			}
			if record.Config != inherited {
				t.Fatalf("config = %q, want inherited config", record.Config)
			}
		})
	}
}

func TestIsJevRouterSelection(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, runtime, model string
		want                 bool
	}{
		{"router", "opencode", JevRouterModel, true},
		{"custom", "opencode-custom", JevRouterModel, true},
		{"trimmed", " opencode ", " " + JevRouterModel + " ", true},
		{"other-runtime", "codex", JevRouterModel, false},
		{"other-model", "opencode", "other", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := IsJevRouterSelection(test.runtime, test.model); got != test.want {
				t.Fatalf("selection = %t, want %t", got, test.want)
			}
		})
	}
}

const jevRouterFixturePost = "ROUNDFIX_FAKE_ACPX_JEV_POST"

func routerFixtureBaseURL(t *testing.T, content string) string {
	t.Helper()
	var config struct {
		Provider map[string]struct {
			Options struct {
				BaseURL string `json:"baseURL"`
			} `json:"options"`
		} `json:"provider"`
	}
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		t.Fatal(err)
	}
	base := config.Provider["roundfix-openrouter"].Options.BaseURL
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme != "http" || parsed.Hostname() != "127.0.0.1" || parsed.Port() == "" || parsed.Path == "" {
		t.Fatalf("not a session loopback URL: %q", base)
	}
	return base
}

// The compiled acpx fixture posts only to the relay URL supplied by its parent.
func postJevRouterFixture() error {
	if os.Getenv(jevRouterFixturePost) != "1" {
		return nil
	}
	var config struct {
		Provider map[string]struct {
			Options struct {
				BaseURL string `json:"baseURL"`
			} `json:"options"`
		} `json:"provider"`
	}
	if err := json.Unmarshal([]byte(os.Getenv("OPENCODE_CONFIG_CONTENT")), &config); err != nil {
		return err
	}
	base := config.Provider["roundfix-openrouter"].Options.BaseURL
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme != "http" || parsed.Hostname() != "127.0.0.1" {
		return errors.New("fixture refuses non-loopback endpoint")
	}
	req, err := http.NewRequest(http.MethodPost, base+"/chat/completions", strings.NewReader(`{"model":"typesafe/jev-router"}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv(JevRouterKeyEnv))
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	return err
}

func TestJevRouterSessionPointsOpenCodeAtTheRelay(t *testing.T) {
	t.Parallel()
	harness := newFakeACPXHarness(t)
	t.Cleanup(func() { harness.runner.clearSessionState("router-test") })
	path := filepath.Join(harness.gitRoot, "config.json")
	harness.setEnv(jevRouterFixtureConfigPath, path)
	harness.setEnv(JevRouterKeyEnv, jevRouterSentinel)
	var first string
	for i := 0; i < 2; i++ {
		if err := runJevRouterPrompt(t, harness, JevRouterModel); err != nil {
			t.Fatal(err)
		}
		var record jevRouterFixtureEnvironment
		if err := json.Unmarshal([]byte(readFile(t, path)), &record); err != nil {
			t.Fatal(err)
		}
		base := routerFixtureBaseURL(t, record.Config)
		if record.Config != jevRouterProviderConfig(base) || !strings.Contains(record.Config, "{env:ROUNDFIX_OPENROUTER_API_KEY}") || strings.Contains(record.Config, jevRouterSentinel) {
			t.Fatal("inline config violated credential boundary")
		}
		if i == 0 {
			first = base
		} else if base != first {
			t.Fatal("session endpoint changed")
		}
	}
}

func TestJevRouterPromptReturnsTheRelayObservation(t *testing.T) {
	t.Parallel()
	for _, failed := range []bool{false, true} {
		t.Run(fmt.Sprintf("failed=%t", failed), func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/chat/completions" || r.Header.Get("Authorization") != "Bearer "+jevRouterSentinel {
					t.Error("changed routed request")
				}
				w.Header().Set("Content-Type", "text/event-stream")
				fmt.Fprint(w, "data: {\"model\":\"a/one\",\"provider\":\"p\",\"id\":\"first\"}\n\ndata: {\"model\":\"b/two\",\"provider\":\"q\",\"id\":\"last\"}\n\n")
			}))
			defer upstream.Close()
			harness := newFakeACPXHarness(t)
			defer harness.runner.clearSessionState("router-test")
			harness.runner.RouterEndpoint = upstream.URL
			harness.setEnv(JevRouterKeyEnv, jevRouterSentinel)
			harness.setEnv(jevRouterFixturePost, "1")
			if failed {
				harness.setEnv(fakeACPXStdoutBy, `{}`)
				harness.setEnv(fakeACPXExitCode, "1")
			}
			result, err := harness.runner.RunPrompt(t.Context(), ACPXPromptRequest{ExecuteRequest: ExecuteRequest{Runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: JevRouterModel}, GitRoot: harness.gitRoot, Prompt: "fixture"}, Session: "router-test"}, nil)
			if (err != nil) != failed {
				t.Fatalf("err=%v failed=%t", err, failed)
			}
			want := jevrouter.Observation{Models: []string{"a/one", "b/two"}, Providers: []string{"p", "q"}, ResponseID: "last"}
			if !reflect.DeepEqual(result.Router, want) {
				t.Fatalf("router=%+v want=%+v", result.Router, want)
			}
			if got := harness.runner.routerRelay.Take(harness.runner.routerSessions["router-test"].token); !reflect.DeepEqual(got, jevrouter.Observation{}) {
				t.Fatal("prompt did not take observation")
			}
		})
	}
}

func TestJevRouterClearedSessionReleasesItsRelayToken(t *testing.T) {
	t.Parallel()
	harness := newFakeACPXHarness(t)
	harness.setEnv(JevRouterKeyEnv, jevRouterSentinel)
	runtime := RuntimeSpec{ID: "opencode", Model: JevRouterModel}
	for _, name := range []string{"first", "second"} {
		if err := harness.runner.PrepareSession(t.Context(), ExecuteRequest{Runtime: runtime, Session: SessionRef{Name: name, WorkDir: harness.gitRoot}, GitRoot: harness.gitRoot}, nil); err != nil {
			t.Fatal(err)
		}
	}
	defer harness.runner.clearSessionState("second")
	base := harness.runner.routerSessions["first"].baseURL
	if err := harness.runner.EndSession(t.Context(), runtime, SessionRef{Name: "first", WorkDir: harness.gitRoot}); err != nil {
		t.Fatal(err)
	}
	response, err := http.Get(base + "/chat/completions")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 404 {
		t.Fatalf("released token status=%d", response.StatusCode)
	}
	harness.runner.clearSessionState("second")
	if harness.runner.routerRelay != nil {
		t.Fatal("last token did not close relay")
	}
	response, err = http.Get(base)
	if err == nil {
		response.Body.Close()
		t.Fatal("relay still serves")
	}
}

func TestJevRouterCreditRefusalBeforeWorkIsAFailedSelection(t *testing.T) {
	testRouterCreditRefusal(t, false)
}

func TestJevRouterCreditRefusalAfterWorkNamesItsReason(t *testing.T) {
	testRouterCreditRefusal(t, true)
}

func testRouterCreditRefusal(t *testing.T, output bool) {
	t.Helper()
	for _, source := range []string{"openrouter_credits", ""} {
		t.Run(fmt.Sprintf("source=%s", source), func(t *testing.T) {
			var requests atomic.Int32
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusPaymentRequired)
				fmt.Fprintf(w, `{"error":{"metadata":{"limit_source":%q}}}`, source)
			}))
			defer upstream.Close()
			harness := newFakeACPXHarness(t)
			defer harness.runner.clearSessionState("router-test")
			harness.runner.RouterEndpoint = upstream.URL
			harness.setEnv(JevRouterKeyEnv, jevRouterSentinel)
			harness.setEnv(jevRouterFixturePost, "1")
			stdout := `{}`
			if output {
				stdout = acpxUpdateLine(`{"sessionId":"fake","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"started"}}}`)
			}
			harness.setEnv(fakeACPXStdout, stdout)
			harness.setEnv(fakeACPXStdoutBy, `{}`)
			harness.setEnv(fakeACPXExitCode, "1")
			harness.setEnv(fakeACPXStderr, "credit refusal detail")
			sink := newCaptureSink("")
			result, err := harness.runner.RunPrompt(t.Context(), ACPXPromptRequest{ExecuteRequest: ExecuteRequest{Runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: JevRouterModel}, GitRoot: harness.gitRoot, Prompt: "fixture"}, Session: "router-test"}, sink)
			limit := source
			if limit == "" {
				limit = "unspecified"
			}
			want := "openrouter_credit_refused: OpenRouter refused a routed request for credit (" + limit + ")"
			if result.Router.Refusal == nil || result.Router.Refusal.Status != 402 || result.Router.Refusal.LimitSource != source {
				t.Fatalf("observation=%+v", result.Router)
			}
			if requests.Load() != 1 {
				t.Fatalf("requests=%d, refused request must not be retried", requests.Load())
			}
			var selection *SelectionFailureError
			var batch *BatchFailureError
			if output {
				if !errors.As(err, &batch) || errors.As(err, &selection) || batch.Reason != want || batch.ExitCode != 1 || batch.Stderr != "credit refusal detail" {
					t.Fatalf("batch refusal=%v", err)
				}
			} else {
				if !errors.As(err, &selection) || selection.Runtime != "opencode" || selection.Reason != want {
					t.Fatalf("selection refusal=%v", err)
				}
				if countStatusEventsForTest(sink.Events(), AgentSelectionFailedStatus) != 1 {
					t.Fatal("selection failure was not reported")
				}
			}
			if err == nil || !strings.Contains(err.Error(), want) || strings.Contains(err.Error(), "agent/protocol error") {
				t.Fatalf("message=%v", err)
			}
		})
	}
}
