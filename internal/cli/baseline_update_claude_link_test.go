// Suite: Baseline update Claude Instruction Link warnings and surface transcript
// Invariant: advisory warnings preserve bytes and exits; this repository plans no link change.
// Boundary IN: runCLI, adopted temporary Git repositories and homes, read-only repository preview.
// Boundary OUT: confirmation, live home, network and external skill acquisition.
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineUpdateWarnsAboutTheClaudeInstructionLink(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"withheld_text", "withheld_json", "unlinked_checkout_text", "unlinked_checkout_json", "exit_code_unchanged"} {
		t.Run(name, func(t *testing.T) {
			shapes := []string{"withheld"}
			formats := []string{"text"}
			if strings.HasPrefix(name, "unlinked_checkout") {
				shapes = []string{"unlinked_checkout"}
			}
			if strings.HasSuffix(name, "_json") {
				formats = []string{"json"}
			}
			if name == "exit_code_unchanged" {
				shapes, formats = []string{"withheld", "unlinked_checkout"}, []string{"text", "json"}
			}
			for _, shape := range shapes {
				for _, format := range formats {
					repository := newBaselineUpdateRepository(t)
					if name == "exit_code_unchanged" {
						// An unrelated managed-region change keeps the plan at exit 3.
						content := string(mustReadBytes(t, filepath.Join(repository, misplacedTestCarrier)))
						writeBaselinePlanTestFile(t, repository, misplacedTestCarrier, strings.Replace(content, "# Agent instructions", "# Older instructions", 1))
					}
					setCommandEnvironmentForTest(t, t.TempDir(), repository)
					warning := baseline.Finding{Code: "baseline.root.claude-link.withheld", Path: "CLAUDE.md", Message: "CLAUDE.md holds 1 passage of its own (CLAUDE.md:1-2); move each into docs/agents/specific-repository.md, then run roundfix baseline update to replace it with a symbolic link to AGENTS.md"}
					if shape == "unlinked_checkout" {
						unlinkedClaudeInstructionFixture(t, repository)
						warning.Code = "baseline.root.claude-link.unlinked-checkout"
						warning.Message = "this checkout holds the symbolic link CLAUDE.md -> AGENTS.md as a regular file; run git config core.symlinks true, then check CLAUDE.md out again"
					} else {
						writeClaudeInstructionFixture(t, repository, "# Claude notes\n- Keep this rule.\n")
					}
					before := snapshotDoctorPath(t, repository)
					claudeBefore := mustReadBytes(t, filepath.Join(repository, "CLAUDE.md"))
					output, code := runMisplacedBaselineCommand(t, repository, format)
					if !bytes.Equal(mustReadBytes(t, filepath.Join(repository, "CLAUDE.md")), claudeBefore) {
						t.Fatal("preview changed CLAUDE.md bytes")
					}
					if after := snapshotDoctorPath(t, repository); !reflect.DeepEqual(after, before) {
						t.Fatal("preview changed repository")
					}
					if format == "json" {
						var result baselineUpdateResult
						if err := json.Unmarshal([]byte(output), &result); err != nil {
							t.Fatal(err)
						}
						found := 0
						for _, got := range result.Warnings {
							if got.Code == warning.Code {
								found++
								if got.Path != warning.Path || got.Message != warning.Message {
									t.Fatalf("warning=%+v, want %+v", got, warning)
								}
							}
						}
						if found != 1 {
							t.Fatalf("want one warning %+v: %s", warning, output)
						}
						for _, change := range result.FileChanges {
							if change.Path == "CLAUDE.md" {
								t.Fatalf("withheld file planned: %+v", change)
							}
						}
					} else {
						line := "Warning: " + warning.Code + ": " + warning.Path + ": " + warning.Message
						if !strings.Contains(output, "\n"+line+"\n") || strings.Count(output, "Warning: "+warning.Code+":") != 1 {
							t.Fatalf("missing single whole warning %q:\n%s", line, output)
						}
					}
					removeClaudeInstructionFixture(t, repository)
					if err := os.Symlink("AGENTS.md", filepath.Join(repository, "CLAUDE.md")); err != nil {
						t.Fatal(err)
					}
					linkedOutput, linkedCode := runMisplacedBaselineCommand(t, repository, format)
					wantCode := exitOK
					if name == "exit_code_unchanged" {
						wantCode = exitUnverified
					}
					if code != linkedCode || code != wantCode {
						t.Fatalf("%s %s: exit=%d, linked=%d, want=%d; warned=%s linked=%s", shape, format, code, linkedCode, wantCode, output, linkedOutput)
					}
				}
			}
		})
	}
}

func TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	setCommandEnvironmentForTest(t, t.TempDir(), root)
	assertLink := func() {
		t.Helper()
		target, err := os.Readlink(filepath.Join(root, "CLAUDE.md"))
		if err != nil || target != "AGENTS.md" {
			t.Fatalf("CLAUDE.md target=%q error=%v", target, err)
		}
	}
	assertLink()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"baseline", "update", "--repo", root, "--no-skills", "--no-history", "--format", "text"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Baseline update: current", "Result: the repository already matches the current Baseline catalog", "Prior catalog: " + catalog.Digest(), "Current catalog: " + catalog.Digest(), "File changes: 0", "Misplaced repository rules: 0", "Retention evidence: 0", "Declared decisions: 0"}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) <= len(want) || !reflect.DeepEqual(lines[:len(want)], want) {
		t.Fatalf("transcript=%q, want prefix=%q followed by remaining output", lines, want)
	}
	for _, line := range lines {
		if strings.Contains(line, "baseline.root.claude-link") || strings.Contains(line, "CLAUDE.md") {
			t.Fatalf("unexpected link plan or warning: %s", line)
		}
	}
	assertLink()
}
