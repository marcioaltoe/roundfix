// Suite: advisory Flow Command List reports
// Invariant: Doctor reports every present harness, preserves exit status and writes nothing.
// Boundary IN: CLI dispatch, temporary Git repositories, real permission planner.
// Boundary OUT: forge, machine readiness and Agent probes replaced by Doctor fakes.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/flowpermissions"
)

// Each behind fixture drops exactly two allow entries and retains every guard.
func writeFlowReportFixture(t *testing.T, repository, fixture string) {
	t.Helper()
	switch fixture {
	case "claude_not_applied", "codex_not_applied":
		directory := ".claude"
		if fixture == "codex_not_applied" {
			directory = ".codex"
		}
		if err := os.MkdirAll(filepath.Join(repository, directory), 0755); err != nil {
			t.Fatal(err)
		}
	case "claude_current", "claude_behind", "claude_missing_guard":
		var allow, ask []string
		for _, entry := range flowpermissions.Entries() {
			allow = append(allow, flowpermissions.ClaudeRule(entry))
		}
		for _, guard := range flowpermissions.Guards() {
			ask = append(ask, guard.ClaudeRule)
		}
		if fixture == "claude_behind" {
			allow = allow[:len(allow)-2]
		}
		if fixture == "claude_missing_guard" {
			ask = ask[:len(ask)-1]
		}
		data, err := json.Marshal(map[string]any{"permissions": map[string]any{"allow": allow, "ask": ask}})
		if err != nil {
			t.Fatal(err)
		}
		writeBaselinePlanTestFile(t, repository, ".claude/settings.json", string(data))
	case "claude_unreadable":
		writeBaselinePlanTestFile(t, repository, ".claude/settings.json", "{")
	case "codex_current", "codex_behind":
		data := string(flowpermissions.RenderCodexRules())
		if fixture == "codex_behind" {
			// A valid older rendering with one retired command and one missing command.
			data = strings.Replace(data, `"archive"`, `"retired-command"`, 1)
		}
		writeBaselinePlanTestFile(t, repository, ".codex/rules/roundfix-flow.rules", data)
	case "codex_unreadable":
		writeBaselinePlanTestFile(t, repository, ".codex/rules/roundfix-flow.rules", "not rules")
	case "opencode_not_needed":
		writeBaselinePlanTestFile(t, repository, "opencode.json", "{}")
	case "opencode_own_rules":
		writeBaselinePlanTestFile(t, repository, "opencode.jsonc", `{/* project policy */"permission":{"bash":"ask"},}`)
	case "opencode_unreadable":
		writeBaselinePlanTestFile(t, repository, "opencode.json", "{")
	default:
		t.Fatalf("unknown fixture %s", fixture)
	}
}

func TestDoctorReportsFlowPermissions(t *testing.T) {
	t.Parallel()
	const skipped = "flow-permissions: skipped (no harness files in this repository)"
	const claudeBehind = "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, 2 to add, 0 to remove; next: roundfix permissions --harness claude)"
	const claudeUnreadable = "flow-permissions: warn (DR-FLOW-PERMISSIONS-UNREADABLE: claude: not readable; next: roundfix permissions --harness claude)"
	for _, tc := range []struct {
		name     string
		fixtures []string
		want     string
	}{
		{"no_harness_files", nil, skipped},
		{"outside_a_git_repository", nil, "flow-permissions: skipped (outside a Git repository)"},
		{"claude_not_applied", []string{"claude_not_applied"}, "flow-permissions: skipped (claude: not applied; to see the list: roundfix permissions)"},
		{"claude_current", []string{"claude_current"}, "flow-permissions: ok (claude: current)"},
		{"missing_guard_counts_as_addition", []string{"claude_missing_guard"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, 1 to add, 0 to remove; next: roundfix permissions --harness claude)"},
		{"claude_behind", []string{"claude_behind", "codex_not_applied"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, 2 to add, 0 to remove; codex: not applied; next: roundfix permissions --harness claude)"},
		{"claude_unreadable", []string{"claude_unreadable"}, claudeUnreadable},
		{"codex_not_applied", []string{"codex_not_applied"}, "flow-permissions: skipped (codex: not applied; to see the list: roundfix permissions)"},
		{"codex_behind", []string{"codex_behind"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: codex: behind, 1 to add, 1 to remove; next: roundfix permissions --harness codex)"},
		{"opencode_not_needed", []string{"opencode_not_needed"}, "flow-permissions: ok (opencode: not needed)"},
		{"opencode_own_rules", []string{"opencode_own_rules"}, "flow-permissions: ok (opencode: own rules)"},
		{"behind_and_not_applied", []string{"claude_behind", "codex_not_applied"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, 2 to add, 0 to remove; codex: not applied; next: roundfix permissions --harness claude)"},
		{"behind_first_harness_named", []string{"claude_not_applied", "codex_behind"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: not applied; codex: behind, 1 to add, 1 to remove; next: roundfix permissions --harness codex)"},
		{"unreadable_without_behind", []string{"claude_not_applied", "codex_unreadable"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-UNREADABLE: claude: not applied; codex: not readable; next: roundfix permissions --harness codex)"},
		{"never_fails", []string{"claude_unreadable", "codex_behind", "opencode_unreadable"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: not readable; codex: behind, 1 to add, 1 to remove; opencode: not readable; next: roundfix permissions --harness claude)"},
		{"exit_unchanged", []string{"claude_behind"}, claudeBehind},
		{"after_the_baseline_checks_before_residue", nil, skipped},
		{"writes_nothing", []string{"claude_behind", "codex_behind", "opencode_own_rules"}, "flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, 2 to add, 0 to remove; codex: behind, 1 to add, 1 to remove; opencode: own rules; next: roundfix permissions --harness claude)"},
		{"codex_current", []string{"codex_current"}, "flow-permissions: ok (codex: current)"},
		{"mixed_current_and_not_applied", []string{"claude_not_applied", "codex_current"}, "flow-permissions: ok (claude: not applied; codex: current)"},
		{"outside_target_link", nil, claudeUnreadable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, repository := permissionsWorkspace(t)
			for _, fixture := range tc.fixtures {
				writeFlowReportFixture(t, repository, fixture)
			}
			if tc.name == "outside_target_link" {
				outside := t.TempDir()
				writeBaselinePlanTestFile(t, outside, "settings.json", `{}`)
				writeFlowReportFixture(t, repository, "claude_not_applied")
				if err := os.Symlink(filepath.Join(outside, "settings.json"), filepath.Join(repository, ".claude/settings.json")); err != nil {
					t.Fatal(err)
				}
			}
			gitRoot := repository
			if tc.name == "outside_a_git_repository" {
				gitRoot = ""
			}
			nodeStatus, wantCode := CheckStatusOK, exitOK
			// The existing skills check fails outside Git; the advisory stays skipped.
			if tc.name == "outside_a_git_repository" {
				wantCode = exitRunFailed
			}
			if tc.name == "exit_unchanged" {
				nodeStatus, wantCode = CheckStatusFailed, exitRunFailed
			}
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: nodeStatus}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: gitRoot, HomeDir: home}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			before, homeBefore := snapshotDoctorPath(t, repository), snapshotDoctorPath(t, home)
			var out, diagnostics bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &out, &diagnostics)
			if code != wantCode || diagnostics.Len() != 0 {
				t.Fatalf("exit=%d want=%d stderr=%s stdout=%s", code, wantCode, &diagnostics, &out)
			}
			lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
			lines = slices.DeleteFunc(lines, func(line string) bool { return strings.HasPrefix(line, "  DR-") })
			position := slices.Index(lines, tc.want)
			if position < 2 || position+1 >= len(lines) || !strings.HasPrefix(lines[position-1], "baseline-docs:") || !strings.HasPrefix(lines[position-2], "baseline-export:") || !strings.HasPrefix(lines[position+1], "residue:") || strings.Count(out.String(), "flow-permissions:") != 1 {
				t.Fatalf("missing single ordered line %q:\n%s", tc.want, &out)
			}
			if !reflect.DeepEqual(before, snapshotDoctorPath(t, repository)) || !reflect.DeepEqual(homeBefore, snapshotDoctorPath(t, home)) {
				t.Fatal("Doctor wrote files")
			}
		})
	}
}
