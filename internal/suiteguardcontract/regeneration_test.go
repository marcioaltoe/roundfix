// Suite: sanctioned regeneration discovery.
// Invariant: only operative grants authorize the regular outputs owned by their declared command.
// Boundary IN: temporary repository files, Spec grants, legacy grants, and ownership declarations.
// Boundary OUT: suiteguard process isolation and Baseline's independent ownership reader.
package suiteguardcontract

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCleanupRegenerationDiscovery(t *testing.T) {
	t.Run("approved Spec grant works without the legacy directory", func(t *testing.T) {
		repository := t.TempDir()
		writeCleanupRegenerationFile(t, repository, "Makefile", "DERIVED_DIGEST_PATHS := internal/baseline/derived\n")
		writeCleanupRegenerationFile(t, repository, "internal/baseline/derived/_ownership.yml", "owner: sanctioned\nreason: fixture\n")
		writeCleanupRegenerationFile(t, repository, "internal/baseline/derived/generated.txt", "generated\n")

		writeCleanupRegenerationFile(
			t,
			repository,
			"docs/specs/current-grant/references/approved.md",
			cleanupSpecGrant("approved", "2026-09-09", "regenerate the fixture", "current-grant", "internal/baseline/derived/_ownership.yml", "make baseline-digests"),
		)
		invalid := map[string]string{
			"docs/specs/proposed/_authorization.md": cleanupSpecGrant(
				"proposed", "2026-09-09", "not granted", "proposed", "Makefile", "proposed-command",
			),
			"docs/specs/null-date/_authorization.md": cleanupSpecGrant(
				"approved", "null", "missing date", "null-date", "Makefile", "null-command",
			),
			"docs/specs/malformed/_authorization.md": "---\nstatus: approved\npaths: [\n---\n\n## Sanctioned regeneration\n\n```yaml\ncommand: malformed-command\n```\n",
			"docs/specs/unrelated/README.md":         "# Unrelated document\n\n## Sanctioned regeneration\n\n```yaml\ncommand: unrelated-command\n```\n",
			"docs/specs/mismatched/_authorization.md": cleanupSpecGrant(
				"approved", "2026-09-09", "wrong consumer", "some-other-spec", "Makefile", "mismatched-command",
			),
			"docs/specs/empty-action/_authorization.md": cleanupSpecGrant(
				"approved", "2026-09-09", "", "empty-action", "Makefile", "empty-action-command",
			),
			"docs/specs/unsafe-path/_authorization.md": cleanupSpecGrant(
				"approved", "2026-09-09", "unsafe path", "unsafe-path", "../Makefile", "unsafe-path-command",
			),
		}
		for path, content := range invalid {
			writeCleanupRegenerationFile(t, repository, path, content)
		}

		symlinkTarget := filepath.Join(repository, "symlinked-grant.md")
		if err := os.WriteFile(
			symlinkTarget,
			[]byte(cleanupSpecGrant("approved", "2026-09-09", "symlink", "symlinked", "Makefile", "symlink-command")),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
		symlinkPath := filepath.Join(repository, "docs", "specs", "symlinked", "_authorization.md")
		if err := os.MkdirAll(filepath.Dir(symlinkPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(symlinkTarget, symlinkPath); err != nil {
			t.Fatal(err)
		}

		want := []SanctionedRegeneration{{
			Command: "make baseline-digests",
			Outputs: []string{"internal/baseline/derived/generated.txt"},
		}}
		got, err := ReadSanctionedRegenerations(repository)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("operative declarations = %#v, want %#v", got, want)
		}
		again, err := ReadSanctionedRegenerations(repository)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(again, got) {
			t.Fatalf("second discovery = %#v, want deterministic %#v", again, got)
		}
	})

	t.Run("legacy unfrontmattered declaration remains compatible", func(t *testing.T) {
		repository := t.TempDir()
		writeCleanupRegenerationFile(
			t,
			repository,
			"docs/workflow/authorizations/legacy.md",
			"# Legacy grant\n\n## Sanctioned regeneration\n\n```yaml\ncommand: legacy-generator\noutputs:\n  - generated/legacy.txt\n```\n",
		)

		want := []SanctionedRegeneration{{
			Command: "legacy-generator",
			Outputs: []string{"generated/legacy.txt"},
		}}
		got, err := ReadSanctionedRegenerations(repository)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("legacy declarations = %#v, want %#v", got, want)
		}
	})

	t.Run("symlinked derived input fails closed", func(t *testing.T) {
		repository := t.TempDir()
		writeCleanupRegenerationFile(t, repository, "Makefile", "DERIVED_DIGEST_PATHS := internal/baseline/derived\n")
		writeCleanupRegenerationFile(t, repository, "internal/baseline/derived/_ownership.yml", "owner: sanctioned\nreason: fixture\n")
		writeCleanupRegenerationFile(
			t,
			repository,
			"docs/specs/symlink-output/_authorization.md",
			cleanupSpecGrant("approved", "2026-09-09", "reject symlink", "symlink-output", "Makefile", "make baseline-digests"),
		)
		outside := filepath.Join(repository, "outside.txt")
		if err := os.WriteFile(outside, []byte("outside\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output := filepath.Join(repository, "internal", "baseline", "derived", "output.txt")
		if err := os.Symlink(outside, output); err != nil {
			t.Fatal(err)
		}

		got, err := ReadSanctionedRegenerations(repository)
		if err == nil || !strings.Contains(err.Error(), "is a symlink") {
			t.Fatalf("symlinked output declarations = %#v, error = %v, want symlink refusal", got, err)
		}
		if got != nil {
			t.Fatalf("symlinked output declarations = %#v, want no partial authority", got)
		}
	})

	t.Run("missing optional roots grant nothing", func(t *testing.T) {
		got, err := ReadSanctionedRegenerations(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Fatalf("declarations = %#v, want nil", got)
		}
	})

	t.Run("invalid authorization root remains an error", func(t *testing.T) {
		repository := t.TempDir()
		writeCleanupRegenerationFile(t, repository, "docs/specs", "not a directory\n")
		got, err := ReadSanctionedRegenerations(repository)
		if err == nil || !strings.Contains(err.Error(), "is not a directory") {
			t.Fatalf("invalid-root declarations = %#v, error = %v, want directory error", got, err)
		}
		if got != nil {
			t.Fatalf("invalid-root declarations = %#v, want no partial authority", got)
		}
	})
}

func cleanupSpecGrant(status, granted, action, consuming, boundedPath, command string) string {
	return "---\n" +
		"status: " + status + "\n" +
		"granted: " + granted + "\n" +
		"action: " + action + "\n" +
		"consuming: " + consuming + "\n" +
		"paths:\n" +
		"  - " + boundedPath + "\n" +
		"---\n\n" +
		"# Grant\n\n" +
		"## Sanctioned regeneration\n\n" +
		"```yaml\n" +
		"command: " + command + "\n" +
		"```\n"
}

func writeCleanupRegenerationFile(t *testing.T, repository, relative, content string) {
	t.Helper()
	path := filepath.Join(repository, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory for %q: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture %q: %v", relative, err)
	}
}
