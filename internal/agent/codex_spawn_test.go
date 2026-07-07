package agent

import (
	"strings"
	"testing"
)

func TestACPXCommandEnvStripsGuardAndHygieneVariables(t *testing.T) {
	t.Setenv(claudeNestedGuardEnv, "1")
	t.Setenv(codexPathEnv, "/tmp/dirty-codex")
	t.Setenv("ROUNDFIX_TEST_KEEP", "kept")

	env := acpxCommandEnv([]string{codexPathEnv + "=/tmp/clean-codex"})

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
	t.Setenv(claudeNestedGuardEnv, "1")

	for _, entry := range acpxCommandEnv(nil) {
		if strings.HasPrefix(entry, claudeNestedGuardEnv+"=") {
			t.Fatalf("expected %s stripped even without overrides, got %q", claudeNestedGuardEnv, entry)
		}
	}
}
