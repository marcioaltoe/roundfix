// Suite: Doctor Claude Instruction Link report
// Invariant: one ordered advisory line never changes Doctor's exit code or files.
// Boundary IN: runCLI, temporary Git repositories and homes, stubbed readiness.
// Boundary OUT: live machine, forge and Agent probes.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestDoctorReportsTheClaudeInstructionLink(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"ok", "absent", "regular_file", "link_elsewhere", "file_with_passages", "link_elsewhere_with_passages", "agents_not_regular", "unlinked_checkout", "unread", "skipped", "after_baseline_rules_and_exit_code"} {
		t.Run(name, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			home := t.TempDir()
			setCommandEnvironmentForTest(t, home, repository)
			want := "baseline-claude-link: ok (CLAUDE.md is a symbolic link to AGENTS.md)"
			switch name {
			case "absent":
				removeClaudeInstructionFixture(t, repository)
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: CLAUDE.md is absent; next: run roundfix baseline update to create the symbolic link to AGENTS.md)"
			case "regular_file":
				writeClaudeInstructionFixture(t, repository, "@AGENTS.md\n")
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: CLAUDE.md is a regular file; next: run roundfix baseline update to replace it with a symbolic link to AGENTS.md)"
			case "link_elsewhere", "link_elsewhere_with_passages":
				removeClaudeInstructionFixture(t, repository)
				content := "@AGENTS.md\n"
				if name == "link_elsewhere_with_passages" {
					// Two tool-delimited spans separate the repository's passages.
					content = "# Local rules\n- Keep this rule.\n<!-- BEGIN:tool-rules -->\nTool guidance.\n<!-- END:tool-rules -->\n# More rules\n- Keep this too.\n"
				}
				writeBaselinePlanTestFile(t, repository, "docs/claude-notes.md", content)
				if err := os.Symlink("docs/claude-notes.md", filepath.Join(repository, "CLAUDE.md")); err != nil {
					t.Fatal(err)
				}
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: CLAUDE.md is a symbolic link to docs/claude-notes.md; next: run roundfix baseline update to replace it with a symbolic link to AGENTS.md)"
				if name == "link_elsewhere_with_passages" {
					want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: CLAUDE.md is a symbolic link to docs/claude-notes.md with 2 passages of its own: docs/claude-notes.md:1-2, docs/claude-notes.md:6-7; next: move each passage into docs/agents/specific-repository.md, then run roundfix baseline update)"
				}
			case "file_with_passages":
				writeClaudeInstructionFixture(t, repository, "# Claude notes\n- Keep this rule.\n")
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: CLAUDE.md is a regular file with 1 passage of its own: CLAUDE.md:1-2; next: move each passage into docs/agents/specific-repository.md, then run roundfix baseline update)"
			case "agents_not_regular":
				writeClaudeInstructionFixture(t, repository, "@AGENTS.md\n")
				if err := os.Remove(filepath.Join(repository, "AGENTS.md")); err != nil {
					t.Fatal(err)
				}
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: AGENTS.md is not a regular file; next: make AGENTS.md the root instruction file, then run roundfix baseline update)"
			case "unlinked_checkout":
				unlinkedClaudeInstructionFixture(t, repository)
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK-UNLINKED: CLAUDE.md is the symbolic link to AGENTS.md checked out as a regular file; next: run git config core.symlinks true, then check CLAUDE.md out again)"
			case "unread":
				removeClaudeInstructionFixture(t, repository)
				if err := os.Mkdir(filepath.Join(repository, "CLAUDE.md"), 0o755); err != nil {
					t.Fatal(err)
				}
				want = "baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK-UNREAD: inspect CLAUDE.md: expected a regular file or safe symbolic link; next: re-run roundfix doctor after fixing CLAUDE.md)"
			case "skipped":
				if err := os.Remove(filepath.Join(repository, baselineSetupManifestPath)); err != nil {
					t.Fatal(err)
				}
				want = "baseline-claude-link: skipped (no Setup Manifest at docs/agents/setup-context.json)"
			}
			nodeStatus, wantCode := CheckStatusOK, exitOK
			if name == "after_baseline_rules_and_exit_code" {
				nodeStatus, wantCode = CheckStatusFailed, exitRunFailed
			}
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: nodeStatus}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: repository, HomeDir: home}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			before := snapshotDoctorPath(t, repository)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			if code != wantCode || stderr.Len() != 0 {
				t.Fatalf("exit=%d want=%d stderr=%s stdout=%s", code, wantCode, stderr.String(), stdout.String())
			}
			lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			lines = slices.DeleteFunc(lines, func(line string) bool { return strings.HasPrefix(line, "  DR-") })
			position := slices.Index(lines, want)
			if position < 1 || position+3 >= len(lines) || !strings.HasPrefix(lines[position-1], "baseline-rules:") || !strings.HasPrefix(lines[position+1], "baseline-export:") || !strings.HasPrefix(lines[position+2], "baseline-docs:") || !strings.HasPrefix(lines[position+3], "residue:") || strings.Count(stdout.String(), "baseline-claude-link:") != 1 {
				t.Fatalf("missing single ordered whole line %q:\n%s", want, stdout.String())
			}
			if after := snapshotDoctorPath(t, repository); !reflect.DeepEqual(after, before) {
				t.Fatal("Doctor changed repository")
			}
		})
	}
}

func removeClaudeInstructionFixture(t *testing.T, repository string) {
	t.Helper()
	if err := os.Remove(filepath.Join(repository, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
}

func writeClaudeInstructionFixture(t *testing.T, repository, content string) {
	t.Helper()
	removeClaudeInstructionFixture(t, repository)
	writeBaselinePlanTestFile(t, repository, "CLAUDE.md", content)
}

func unlinkedClaudeInstructionFixture(t *testing.T, repository string) {
	t.Helper()
	runBaselinePlanTestCommand(t, repository, "git", "config", "core.symlinks", "false")
	removeClaudeInstructionFixture(t, repository)
	runBaselinePlanTestCommand(t, repository, "git", "checkout", "--", "CLAUDE.md")
	info, err := os.Lstat(filepath.Join(repository, "CLAUDE.md"))
	if err != nil || !info.Mode().IsRegular() {
		t.Fatalf("unlinked checkout info=%v err=%v", info, err)
	}
	if got := string(mustReadBytes(t, filepath.Join(repository, "CLAUDE.md"))); got != "AGENTS.md" {
		t.Fatalf("unlinked bytes=%q", got)
	}
}
