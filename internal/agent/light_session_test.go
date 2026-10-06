package agent

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/config"
	"roundfix/internal/openrouterkey"
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
	variable, _ := config.OpenRouterImplementKey([]string{openrouterkey.Shared + "=fixture"})
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

func TestOpenModelImplementationPrefersTheImplementStageKey(t *testing.T) {
	t.Parallel()
	environ := []string{
		openrouterkey.Implement + "=implement-sentinel",
		openrouterkey.Shared + "=shared-sentinel",
		"OPENROUTER_API_KEY=generic-sentinel",
	}
	variable, present := config.OpenRouterImplementKey(environ)
	if !present || variable != openrouterkey.Implement {
		t.Fatalf("selected variable=%q present=%t", variable, present)
	}
	runtime := RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "openrouter/" + config.DefaultLightModel, OpenRouterKeyVariable: variable}
	overrides, err := lightSessionEnvironment(runtime, environ)
	if err != nil {
		t.Fatal(err)
	}
	inline := environmentValue(overrides, "OPENCODE_CONFIG_CONTENT")
	if !strings.Contains(inline, `"apiKey":"{env:`+openrouterkey.Implement+`}"`) {
		t.Fatalf("OpenCode config does not reference the implementation variable: %s", inline)
	}
	if strings.Contains(inline, openrouterkey.Shared) || strings.Contains(inline, "implement-sentinel") || strings.Contains(inline, "shared-sentinel") {
		t.Fatalf("OpenCode config carries another variable or a key value: %s", inline)
	}
	commandEnv := acpxCommandEnv(environ, overrides)
	if value := environmentValue(commandEnv, "OPENROUTER_API_KEY"); value != "" {
		t.Fatalf("generic key reached OpenCode: %q", value)
	}
	if value := environmentValue(commandEnv, openrouterkey.Implement); value != "implement-sentinel" {
		t.Fatalf("implementation key missing from OpenCode environment: %q", value)
	}
}

func TestOpenModelImplementationFallsBackToTheSharedKey(t *testing.T) {
	t.Parallel()
	environ := []string{
		openrouterkey.Shared + "=shared-sentinel",
		"OPENROUTER_API_KEY=generic-sentinel",
	}
	variable, present := config.OpenRouterImplementKey(environ)
	if !present || variable != openrouterkey.Shared {
		t.Fatalf("selected variable=%q present=%t", variable, present)
	}
	runtime := RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "openrouter/" + config.DefaultLightModel, OpenRouterKeyVariable: variable}
	overrides, err := lightSessionEnvironment(runtime, environ)
	if err != nil {
		t.Fatal(err)
	}
	inline := environmentValue(overrides, "OPENCODE_CONFIG_CONTENT")
	if !strings.Contains(inline, `"apiKey":"{env:`+openrouterkey.Shared+`}"`) {
		t.Fatalf("OpenCode config does not reference the shared variable: %s", inline)
	}
	if strings.Contains(inline, openrouterkey.Implement) {
		t.Fatalf("OpenCode config references the implementation variable: %s", inline)
	}
	commandEnv := acpxCommandEnv(environ, overrides)
	if value := environmentValue(commandEnv, "OPENROUTER_API_KEY"); value != "" {
		t.Fatalf("generic key reached OpenCode: %q", value)
	}
	if value := environmentValue(commandEnv, openrouterkey.Shared); value != "shared-sentinel" {
		t.Fatalf("shared key missing from OpenCode environment: %q", value)
	}
	variable, present = config.OpenRouterImplementKey([]string{"OPENROUTER_API_KEY=generic-sentinel"})
	if present || variable != "" {
		t.Fatalf("generic key selected variable=%q present=%t", variable, present)
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
