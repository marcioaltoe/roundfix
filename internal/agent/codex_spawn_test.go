package agent

import (
	"strings"
	"testing"
)

func TestACPXCommandEnvStripsGuardAndHygieneVariables(t *testing.T) {
	t.Parallel()

	base := []string{
		claudeNestedGuardEnv + "=1",
		codexPathEnv + "=/tmp/dirty-codex",
		"ROUNDFIX_TEST_KEEP=kept",
	}
	env := acpxCommandEnv(base, []string{codexPathEnv + "=/tmp/clean-codex"})

	var keep, override bool
	for _, entry := range env {
		key, value, _ := strings.Cut(entry, "=")
		switch key {
		case claudeNestedGuardEnv:
			t.Fatalf("expected %s stripped from acpx child env, got %q", claudeNestedGuardEnv, entry)
		case codexPathEnv:
			if value != "/tmp/clean-codex" {
				t.Fatalf("expected only the override %s, got %q", codexPathEnv, entry)
			}
			override = true
		case "ROUNDFIX_TEST_KEEP":
			keep = true
		}
	}
	if !keep {
		t.Fatal("expected unrelated environment to be inherited")
	}
	if !override {
		t.Fatal("expected the codex path override to be appended")
	}
}

func TestACPXCommandEnvStripsGuardWithoutOverrides(t *testing.T) {
	t.Parallel()

	for _, entry := range acpxCommandEnv([]string{claudeNestedGuardEnv + "=1"}, nil) {
		if strings.HasPrefix(entry, claudeNestedGuardEnv+"=") {
			t.Fatalf("expected %s stripped even without overrides, got %q", claudeNestedGuardEnv, entry)
		}
	}
}

func TestACPXRunnerCommandEnvDefaultsToProcessEnvironment(t *testing.T) {
	// Sequential: verifies the zero-value runner reads the process environment.
	t.Setenv("ROUNDFIX_TEST_PROCESS_ENV", "process-value")

	environment := (ACPXRunner{}).commandEnv(nil)
	if got := environmentValue(environment, "ROUNDFIX_TEST_PROCESS_ENV"); got != "process-value" {
		t.Fatalf("process environment value = %q, want process-value", got)
	}
}
