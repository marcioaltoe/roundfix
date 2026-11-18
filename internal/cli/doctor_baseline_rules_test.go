// Suite: Doctor repository-rule placement report
// Invariant: one ordered advisory line never changes Doctor exit behavior.
// Boundary IN: runCLI with stubbed readiness checks and temporary repositories and homes.
// Boundary OUT: live machine, forge, and adapter probes.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestDoctorReportsMisplacedRepositoryRules(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"ok", "warn_one_passage", "warn_two_passages", "warn_extension_off", "skipped", "unread_invalid_manifest", "after_skills_and_exit_code"} {
		t.Run(name, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			home := t.TempDir()
			config := roundconfig.Builtin()
			nodeStatus := CheckStatusOK
			if name == "after_skills_and_exit_code" {
				home, _, _ = staleUserSkillsFixture(t)
				config.Skills.Authoring.Claude = &roundconfig.AuthoringHarness{PRD: &roundconfig.AuthoringEntry{Model: "opus", ReasoningEffort: "high"}}
				nodeStatus = CheckStatusFailed
			}
			setCommandEnvironmentForTest(t, home, repository)
			if name != "warn_extension_off" {
				_, code := runMisplacedBaselineCommand(t, repository, "json", "--decision", "repository.extension.enabled=true", "--yes")
				if code != exitOK {
					t.Fatalf("enable extension exit=%d", code)
				}
			}
			want := "baseline-rules: ok (no repository rule outside docs/agents/specific-repository.md)"
			switch name {
			case "warn_one_passage", "warn_two_passages", "warn_extension_off":
				rules := misplacedBaselineFixture(t, repository, name == "warn_two_passages")
				locations := make([]string, len(rules))
				for i, rule := range rules {
					locations[i] = fmt.Sprintf("%s:%d-%d", rule.Path, rule.StartLine, rule.EndLine)
				}
				noun := "1 passage belongs"
				if len(rules) == 2 {
					noun = "2 passages belong"
				}
				next := "move each passage into docs/agents/specific-repository.md; roundfix baseline update lists them"
				if name == "warn_extension_off" {
					next = "roundfix baseline update --decision repository.extension.enabled=true, then move each passage"
				}
				want = "baseline-rules: warn (DR-BASELINE-RULE-MISPLACED: " + noun + " in docs/agents/specific-repository.md: " + strings.Join(locations, ", ") + "; next: " + next + ")"
			case "skipped":
				repository = t.TempDir()
				setCommandEnvironmentForTest(t, home, repository)
				want = "baseline-rules: skipped (no Setup Manifest at docs/agents/setup-context.json)"
			case "unread_invalid_manifest":
				writeBaselinePlanTestFile(t, repository, "docs/agents/setup-context.json", "{invalid\n")
				want = "baseline-rules: warn (DR-BASELINE-RULES-UNREAD: read misplaced-rule Setup Manifest: baseline.preservation.managed-refresh.manifest.invalid; next: re-run roundfix doctor after fixing the Setup Manifest)"
			}
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: nodeStatus}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: config, GitRoot: repository, HomeDir: home}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			before := snapshotDoctorPath(t, repository)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			wantCode := exitOK
			if name == "after_skills_and_exit_code" {
				wantCode = exitRunFailed
			}
			if code != wantCode || stderr.Len() != 0 {
				t.Fatalf("exit=%d want=%d stderr=%s stdout=%s", code, wantCode, stderr.String(), stdout.String())
			}
			lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			position := slices.Index(lines, want)
			if position < 0 {
				t.Fatalf("missing whole line %q:\n%s", want, stdout.String())
			}
			count := 0
			for i, line := range lines {
				if strings.HasPrefix(line, "baseline-rules:") {
					count++
				}
				for _, prefix := range []string{"skills:", "user-skills:", "authoring-models:"} {
					if strings.HasPrefix(line, prefix) && i >= position {
						t.Fatalf("%s is after baseline rules: %s", prefix, stdout.String())
					}
				}
			}
			if count != 1 || position+2 >= len(lines) || !strings.HasPrefix(lines[position+1], "baseline-claude-link:") || !strings.HasPrefix(lines[position+2], "residue:") {
				t.Fatalf("baseline rules position: %s", stdout.String())
			}
			if name == "after_skills_and_exit_code" {
				for _, prefix := range []string{"user-skills:", "authoring-models:"} {
					if !strings.Contains(stdout.String(), "\n"+prefix) {
						t.Fatalf("missing optional %s line: %s", prefix, stdout.String())
					}
				}
				// An existing failed check still controls the exit; the new line is advisory.
				if !strings.HasPrefix(lines[0], "node: failed") {
					t.Fatalf("missing independent failure: %s", stdout.String())
				}
			}
			if after := snapshotDoctorPath(t, repository); !reflect.DeepEqual(after, before) {
				t.Fatalf("Doctor changed repository %s", filepath.Base(repository))
			}
		})
	}
}
