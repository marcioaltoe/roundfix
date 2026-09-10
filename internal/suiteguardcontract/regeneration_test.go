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

	"roundfix/internal/authorization"
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
			"docs/specs/current-grant/references/approved-authorization.md",
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
			"docs/workflow/authorizations/2026-08-12-legacy.md",
			cleanupLegacyGrant("0099-legacy-consumer", "legacy-generator", "generated/legacy.txt"),
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

func TestSanctionedRegenerationReadsOnlyCandidateRecords(t *testing.T) {
	repository := t.TempDir()
	authorizationRoot := filepath.Join(repository, filepath.FromSlash(specAuthorizationRoot))
	writeCleanupRegenerationFile(t, repository, "docs/specs/current/_authorization.md", "canonical record\n")
	writeCleanupRegenerationFile(
		t,
		repository,
		"docs/specs/current/references/2026-09-08-historical-authorization.md",
		"preserved record\n",
	)

	countCandidateReads := func() int {
		t.Helper()
		reads := 0
		err := walkAuthorizationRecords(
			authorizationRoot,
			authorization.AuthorizationRoleSpec,
			func(relative string) error {
				if _, err := os.ReadFile(filepath.Join(authorizationRoot, relative)); err != nil {
					return err
				}
				reads++
				return nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		return reads
	}

	before := countCandidateReads()
	writeCleanupRegenerationFile(t, repository, "docs/specs/current/task_01.md", "non-record\n")
	after := countCandidateReads()

	if before != 2 {
		t.Fatalf("candidate reads before unrelated Markdown = %d, want 2", before)
	}
	if after != before {
		t.Fatalf("candidate reads after unrelated Markdown = %d, want unchanged %d", after, before)
	}
}

func TestSanctionedRegenerationResolvesOncePerProcess(t *testing.T) {
	repository := t.TempDir()
	recordPath := "docs/specs/current/_authorization.md"
	writeCleanupRegenerationFile(
		t,
		repository,
		recordPath,
		cleanupSpecGrantWithOutput("current", "first-command", "generated/first.txt"),
	)

	want := []SanctionedRegeneration{{
		Command: "first-command",
		Outputs: []string{"generated/first.txt"},
	}}
	got, err := ReadSanctionedRegenerations(repository)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("first resolution = %#v, want %#v", got, want)
	}

	writeCleanupRegenerationFile(
		t,
		repository,
		recordPath,
		cleanupSpecGrantWithOutput("current", "second-command", "generated/second.txt"),
	)
	again, err := ReadSanctionedRegenerations(repository)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(again, want) {
		t.Fatalf("second resolution = %#v, want cached %#v", again, want)
	}
}

func TestSanctionedRegenerationResolvesLegacyRecordsWithoutFrontmatter(t *testing.T) {
	repository := t.TempDir()
	writeCleanupRegenerationFile(
		t,
		repository,
		"docs/workflow/authorizations/legacy.md",
		cleanupLegacyRegenerationOnly(
			"legacy-generator",
			[]string{"generated/legacy-z.txt", "generated/legacy-a.txt"},
		),
	)

	want := []SanctionedRegeneration{{
		Command: "legacy-generator",
		Outputs: []string{"generated/legacy-a.txt", "generated/legacy-z.txt"},
	}}
	got, err := ReadSanctionedRegenerations(repository)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("legacy declarations = %#v, want %#v", got, want)
	}
}

func TestSanctionedRegenerationReadsArchivedSpecGrants(t *testing.T) {
	repository := t.TempDir()
	writeCleanupRegenerationFile(t, repository, "Makefile", "DERIVED_DIGEST_PATHS := internal/baseline/derived\n")
	writeCleanupRegenerationFile(t, repository, "internal/baseline/derived/_ownership.yml", "owner: sanctioned\nreason: fixture\n")
	writeCleanupRegenerationFile(t, repository, "internal/baseline/derived/generated.txt", "generated\n")
	writeCleanupRegenerationFile(
		t,
		repository,
		"docs/history/specs/archived-grant/_authorization.md",
		cleanupMultiSpecGrant(
			"archived-grant",
			"another-consumer",
			"internal/baseline/derived/_ownership.yml",
			"make baseline-digests",
		),
	)

	want := []SanctionedRegeneration{{
		Command: "make baseline-digests",
		Outputs: []string{"internal/baseline/derived/generated.txt"},
	}}
	got, err := ReadSanctionedRegenerations(repository)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("archived multi-Spec declarations = %#v, want %#v", got, want)
	}
}

func TestSanctionedRegenerationRejectsNonOperativeRecords(t *testing.T) {
	tests := []struct {
		name       string
		recordPath string
		content    string
	}{
		{
			name:       "proposed",
			recordPath: "docs/specs/proposed/_authorization.md",
			content: cleanupSpecGrant(
				"proposed", "2026-09-09", "not granted", "proposed", "Makefile", "proposed-command",
			),
		},
		{
			name:       "null dated",
			recordPath: "docs/specs/null-date/_authorization.md",
			content: cleanupSpecGrant(
				"approved", "null", "missing date", "null-date", "Makefile", "null-date-command",
			),
		},
		{
			name:       "malformed",
			recordPath: "docs/specs/malformed/_authorization.md",
			content:    "---\nstatus: approved\npaths: [\n---\n\n## Sanctioned regeneration\n\n```yaml\ncommand: malformed-command\n```\n",
		},
		{
			name:       "unrelated",
			recordPath: "docs/specs/unrelated/_authorization.md",
			content: cleanupSpecGrant(
				"approved", "2026-09-09", "different consumer", "another-spec", "Makefile", "unrelated-command",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := t.TempDir()
			writeCleanupRegenerationFile(t, repository, tt.recordPath, tt.content)

			got, err := ReadSanctionedRegenerations(repository)
			if err != nil {
				t.Fatal(err)
			}
			if got != nil {
				t.Fatalf("non-operative declarations = %#v, want nil", got)
			}
		})
	}
}

func TestSanctionedRegenerationSetOnlyGrows(t *testing.T) {
	repository := t.TempDir()
	writeCleanupRegenerationFile(
		t,
		repository,
		"docs/specs/present-spec/_authorization.md",
		cleanupSpecGrantWithOutputs(
			"present-spec",
			"present-spec-command",
			[]string{"generated/present-spec-z.txt", "generated/present-spec-a.txt"},
		),
	)
	writeCleanupRegenerationFile(
		t,
		repository,
		"docs/workflow/authorizations/present-legacy.md",
		cleanupLegacyRegenerationOnly("present-legacy-command", []string{"generated/present-legacy.txt"}),
	)
	writeCleanupRegenerationFile(
		t,
		repository,
		"docs/history/specs/new-archive/_authorization.md",
		cleanupSpecGrantWithOutput("new-archive", "new-archive-command", "generated/new-archive.txt"),
	)

	recordedPresent := []SanctionedRegeneration{
		{Command: "present-legacy-command", Outputs: []string{"generated/present-legacy.txt"}},
		{
			Command: "present-spec-command",
			Outputs: []string{"generated/present-spec-a.txt", "generated/present-spec-z.txt"},
		},
	}
	got, err := ReadSanctionedRegenerations(repository)
	if err != nil {
		t.Fatal(err)
	}
	for _, present := range recordedPresent {
		if !containsSanctionedRegeneration(got, present) {
			t.Errorf("widened declarations = %#v, missing recorded present declaration %#v", got, present)
		}
	}
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

func cleanupMultiSpecGrant(consuming, additionalConsumer, boundedPath, command string) string {
	return "---\n" +
		"status: approved\n" +
		"granted: 2026-09-09\n" +
		"action: regenerate archived outputs\n" +
		"consuming:\n" +
		"  - " + additionalConsumer + "\n" +
		"  - " + consuming + "\n" +
		"paths:\n" +
		"  - " + boundedPath + "\n" +
		"---\n\n" +
		"# Grant\n\n" +
		"## Sanctioned regeneration\n\n" +
		"```yaml\n" +
		"command: " + command + "\n" +
		"```\n"
}

func cleanupSpecGrantWithOutput(consuming, command, output string) string {
	return cleanupSpecGrantWithOutputs(consuming, command, []string{output})
}

func cleanupSpecGrantWithOutputs(consuming, command string, outputs []string) string {
	var declaredOutputs strings.Builder
	for _, output := range outputs {
		declaredOutputs.WriteString("  - ")
		declaredOutputs.WriteString(output)
		declaredOutputs.WriteByte('\n')
	}
	return "---\n" +
		"status: approved\n" +
		"granted: 2026-09-09\n" +
		"action: regenerate recorded output\n" +
		"consuming: " + consuming + "\n" +
		"paths:\n" +
		"  - " + outputs[0] + "\n" +
		"---\n\n" +
		"# Grant\n\n" +
		"## Sanctioned regeneration\n\n" +
		"```yaml\n" +
		"command: " + command + "\n" +
		"outputs:\n" +
		declaredOutputs.String() +
		"```\n"
}

func cleanupLegacyGrant(consuming, command, output string) string {
	return "# Legacy grant — regenerate recorded output\n\n" +
		"## Consuming Spec\n\n" +
		"- " + consuming + "\n\n" +
		"## Authorized paths\n\n" +
		"- `" + output + "`\n\n" +
		"## Sanctioned regeneration\n\n" +
		"```yaml\n" +
		"command: " + command + "\n" +
		"outputs:\n" +
		"  - " + output + "\n" +
		"```\n"
}

func cleanupLegacyRegenerationOnly(command string, outputs []string) string {
	var declaredOutputs strings.Builder
	for _, output := range outputs {
		declaredOutputs.WriteString("  - ")
		declaredOutputs.WriteString(output)
		declaredOutputs.WriteByte('\n')
	}
	return "# Legacy authorization\n\n" +
		"## Sanctioned regeneration\n\n" +
		"```yaml\n" +
		"command: " + command + "\n" +
		"outputs:\n" +
		declaredOutputs.String() +
		"```\n"
}

func containsSanctionedRegeneration(
	declarations []SanctionedRegeneration,
	want SanctionedRegeneration,
) bool {
	for _, declaration := range declarations {
		if reflect.DeepEqual(declaration, want) {
			return true
		}
	}
	return false
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
