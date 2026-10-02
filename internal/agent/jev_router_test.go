package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	if record.Config != jevRouterProviderConfig || !record.KeyPresent {
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
