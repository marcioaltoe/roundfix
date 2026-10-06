package agent

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/config"
)

// Called by the existing fake acpx process; capture the inline option and the
// generic variable's presence, never the implementation key's value.
func recordFakeLightEnvironment() error {
	path := os.Getenv("ROUNDFIX_FAKE_LIGHT_ENV")
	if path == "" {
		return nil
	}
	_, generic := os.LookupEnv("OPENROUTER_API_KEY")
	row, err := json.Marshal(struct {
		Config  string
		Generic bool
	}{os.Getenv("OPENCODE_CONFIG_CONTENT"), generic})
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(append(row, '\n'))
	return errors.Join(writeErr, file.Close())
}

func lightRuntimeForTest() RuntimeSpec {
	variable, _ := config.OpenRouterImplementKey(func(string) string { return "fixture" })
	return RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "openrouter/" + config.DefaultLightModel, OpenRouterKeyVariable: variable}
}

func TestLightSessionReadsTheImplementKeyByName(t *testing.T) {
	t.Parallel()
	for _, inherited := range []string{"", `{"theme":"dark","provider":{"other":{"name":"kept"},"openrouter":{"options":{"timeout":42,"apiKey":"old"}}}}`} {
		t.Run(inherited, func(t *testing.T) {
			h := newFakeACPXHarness(t)
			runtime := lightRuntimeForTest()
			path := filepath.Join(h.gitRoot, "light-env.jsonl")
			h.setEnv("ROUNDFIX_FAKE_LIGHT_ENV", path)
			h.setEnv(runtime.OpenRouterKeyVariable, "fixture-secret")
			h.setEnv("OPENROUTER_API_KEY", "generic-secret")
			h.setEnv("OPENCODE_CONFIG_CONTENT", inherited)
			if _, err := h.run(t.Context(), runtime, "light-key"); err != nil {
				t.Fatal(err)
			}
			session := SessionRef{Name: "light-key", WorkDir: h.gitRoot}
			if err := h.runner.CancelSession(t.Context(), runtime, session); err != nil {
				t.Fatal(err)
			}
			if err := h.runner.EndSession(t.Context(), runtime, session); err != nil {
				t.Fatal(err)
			}
			rows := strings.Split(strings.TrimSpace(readFile(t, path)), "\n")
			if len(rows) != len(readJSONInvocations(t, h.invocationsPath)) {
				t.Fatal("not every acpx command captured its environment")
			}
			for _, row := range rows {
				var got struct {
					Config  string
					Generic bool
				}
				if err := json.Unmarshal([]byte(row), &got); err != nil {
					t.Fatal(err)
				}
				if got.Generic || strings.Contains(got.Config, "secret") {
					t.Fatalf("credential escaped: %s", row)
				}
				var object map[string]any
				if err := json.Unmarshal([]byte(got.Config), &object); err != nil {
					t.Fatal(err)
				}
				provider := object["provider"].(map[string]any)
				options := provider["openrouter"].(map[string]any)["options"].(map[string]any)
				if options["apiKey"] != "{env:"+runtime.OpenRouterKeyVariable+"}" {
					t.Fatalf("option = %v", options)
				}
				if inherited != "" && (object["theme"] != "dark" || provider["other"] == nil || options["timeout"] != float64(42)) {
					t.Fatal("inherited options lost")
				}
			}
			if environmentValue(h.runner.Environment, "OPENCODE_CONFIG_CONTENT") != inherited {
				t.Fatal("base environment mutated")
			}
			log := readFile(t, filepath.Join(h.gitRoot, "runs", "run-acpx", "agent", "batch-007.log"))
			if strings.Contains(log, "secret") {
				t.Fatal("key leaked to Agent log")
			}
		})
	}
	for _, inherited := range []string{"null", "[]", `"fixture-secret"`, "invalid", "42"} {
		t.Run("nonobject/"+inherited, func(t *testing.T) {
			h := newFakeACPXHarness(t)
			h.setEnv("OPENCODE_CONFIG_CONTENT", inherited)
			_, err := h.run(t.Context(), lightRuntimeForTest(), "bad-config")
			var failure *SelectionPreflightError
			if !errors.As(err, &failure) || strings.Contains(err.Error(), "fixture-secret") {
				t.Fatalf("selection failure = %v", err)
			}
			assertNoFile(t, "work prompt", h.promptsPath)
			assertNoFile(t, "acpx commands", h.invocationsPath)
		})
	}
}

func TestLightSessionSendsNoWarmup(t *testing.T) {
	t.Parallel()
	h := newFakeACPXHarness(t)
	runtime := lightRuntimeForTest()
	h.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{"sessions show": sessionCapabilitySnapshotFixture(t, runtime.Model, []string{runtime.Model}, "effort", "low", []string{"low", "high"})}))
	if _, err := h.run(t.Context(), runtime, "light-no-warmup"); err != nil {
		t.Fatal(err)
	}
	prompts := readJSONStrings(t, h.promptsPath)
	if len(prompts) != 1 || prompts[0] != "prompt" {
		t.Fatalf("prompts = %v", prompts)
	}
}
