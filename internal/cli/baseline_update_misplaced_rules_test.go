// Suite: Baseline update misplaced-rule result contract
// Invariant: reporting names passages without changing planning, writes, or exits.
// Boundary IN: runCLI, temporary adopted repositories, result JSON/text, and the current-repository transcript.
// Boundary OUT: detector internals and external skill acquisition.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

const misplacedTestCarrier = "docs/agents/agent-instructions.md"
const misplacedTestPassage = "## Local rules\nUse the repository's release checklist.\n"

func TestBaselineUpdateNamesMisplacedRepositoryRules(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"plan_text", "plan_json", "current_text", "current_json", "applied_json", "plan_digest_unchanged", "exit_code_unchanged", "extension_decision_off", "no_manifest_omits_the_field"} {
		t.Run(name, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			setCommandEnvironmentForTest(t, t.TempDir(), repository)
			setCommandEnvForTest(t, "PATH", os.Getenv("PATH"))
			if name == "no_manifest_omits_the_field" {
				if err := os.Remove(filepath.Join(repository, "docs/agents/setup-context.json")); err != nil {
					t.Fatal(err)
				}
				for _, format := range []string{"text", "json"} {
					out, code := runMisplacedBaselineCommand(t, repository, format)
					if code != exitUnverified {
						t.Fatalf("exit=%d, want %d: %s", code, exitUnverified, out)
					}
					if strings.Contains(out, "misplacedRepositoryRules") || strings.Contains(out, "Misplaced repository rules:") {
						t.Fatalf("unexpected misplaced-rule output: %s", out)
					}
				}
				return
			}
			if name != "extension_decision_off" {
				_, code := runMisplacedBaselineCommand(t, repository, "json", "--decision", "repository.extension.enabled=true", "--yes")
				if code != exitOK {
					t.Fatalf("enable extension: exit=%d", code)
				}
			}
			current := strings.HasPrefix(name, "current_")
			want := misplacedBaselineFixture(t, repository, !current)
			before := baselinePlanTestTree(t, repository)
			format := "json"
			if strings.HasSuffix(name, "_text") || name == "extension_decision_off" {
				format = "text"
			}
			var extra []string
			if name == "applied_json" {
				extra = append(extra, "--yes")
			}
			out, code := runMisplacedBaselineCommand(t, repository, format, extra...)
			wantCode := exitUnverified
			if current || name == "applied_json" {
				wantCode = exitOK
			}
			if code != wantCode {
				t.Fatalf("exit=%d, want %d: %s", code, wantCode, out)
			}
			if name == "applied_json" {
				data := string(mustReadBytes(t, filepath.Join(repository, misplacedTestCarrier)))
				start := strings.Count(data[:strings.Index(data, misplacedTestPassage)], "\n") + 1
				if start <= want[1].StartLine {
					t.Fatalf("apply did not shift passage: %s", data)
				}
				want = []baseline.MisplacedRepositoryRule{{Path: misplacedTestCarrier, StartLine: start, EndLine: start + 1, FirstLine: "## Local rules", MoveTo: baseline.SpecificRepositoryPath}}
			} else if after := baselinePlanTestTree(t, repository); after != before {
				t.Fatal("preview changed repository bytes")
			}
			if format == "json" {
				assertMisplacedBaselineJSON(t, out, want)
				var result baselineUpdateResult
				if err := json.Unmarshal([]byte(out), &result); err != nil {
					t.Fatal(err)
				}
				wantState := "plan_ready"
				if current {
					wantState = "current"
				} else if name == "applied_json" {
					wantState = "verified"
				}
				if result.State != wantState {
					t.Fatalf("state=%s, want %s", result.State, wantState)
				}
				if name == "plan_digest_unchanged" {
					catalog, err := baseline.LoadEmbeddedCatalog()
					if err != nil {
						t.Fatal(err)
					}
					input, err := baseline.ResolveManifestInput(repository, catalog)
					if err != nil {
						t.Fatal(err)
					}
					directories, err := commandEnvironmentForTest(t).executableDirectories("test search path")
					if err != nil {
						t.Fatal(err)
					}
					outcome, err := baseline.BuildPlan(t.Context(), baseline.PlanRequest{WriterVersion: app.Version, Repository: repository, ProfileID: input.ProfileID, Decisions: input.Decisions, Preservation: baseline.RootPreservationRequest{Mode: baseline.PreservationModeManagedRefresh}, ExecutableDirectories: directories})
					if err != nil || outcome.Plan == nil {
						t.Fatalf("plan=%+v, error=%v", outcome, err)
					}
					if result.PlanDigest != outcome.Plan.PlanDigest {
						t.Fatalf("digest=%s, want %s", result.PlanDigest, outcome.Plan.PlanDigest)
					}
					data, err := json.Marshal(outcome.Plan)
					if err != nil {
						t.Fatal(err)
					}
					if bytes.Contains(data, []byte("misplacedRepositoryRules")) {
						t.Fatalf("passages leaked into Plan Document: %s", data)
					}
				}
			} else {
				lines := []string{fmt.Sprintf("Misplaced repository rules: %d", len(want))}
				for _, rule := range want {
					identity := ""
					if rule.ManagedID != "" {
						identity = " (managed identity " + rule.ManagedID + ")"
					}
					lines = append(lines, fmt.Sprintf("- %s:%d-%d%s: %s", rule.Path, rule.StartLine, rule.EndLine, identity, rule.FirstLine), "  Move to: "+baseline.SpecificRepositoryPath)
				}
				if name == "extension_decision_off" {
					lines = append(lines, "Next for misplaced rules: roundfix baseline update --decision repository.extension.enabled=true, then move each passage")
				}
				lines = append(lines, "Retention evidence: 0")
				section := strings.Join(lines, "\n") + "\n"
				if !strings.Contains(out, section) {
					t.Fatalf("missing whole section %q:\n%s", section, out)
				}
				if !current && strings.Index(out, "Unrecorded managed regions:") > strings.Index(out, lines[0]) {
					t.Fatalf("section before unrecorded regions: %s", out)
				}
				if name != "extension_decision_off" && strings.Contains(out, "Next for misplaced rules:") {
					t.Fatalf("unexpected extension advice: %s", out)
				}
			}
			if name == "exit_code_unchanged" {
				data := string(mustReadBytes(t, filepath.Join(repository, misplacedTestCarrier)))
				writeBaselinePlanTestFile(t, repository, misplacedTestCarrier, strings.TrimSuffix(data, "\n"+misplacedTestPassage))
				_, withoutCode := runMisplacedBaselineCommand(t, repository, "json")
				if code != withoutCode {
					t.Fatalf("passage changed exit from %d to %d", withoutCode, code)
				}
			}
		})
	}
}

func runMisplacedBaselineCommand(t *testing.T, repository, format string, extra ...string) (string, int) {
	t.Helper()
	args := []string{"baseline", "update", "--repo", repository, "--no-skills", "--no-history", "--format", format}
	args = append(args, extra...)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, args, &stdout, &stderr)
	if stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	return stdout.String(), code
}

func misplacedBaselineFixture(t *testing.T, repository string, mismatch bool) []baseline.MisplacedRepositoryRule {
	t.Helper()
	content := string(mustReadBytes(t, filepath.Join(repository, misplacedTestCarrier)))
	var want []baseline.MisplacedRepositoryRule
	if mismatch {
		manifest := ReadBaselineSetupManifest(t, repository)
		var version string
		for _, artifact := range manifest.ManagedArtifacts {
			if artifact.ID == "guide.agent-instructions" {
				version = artifact.Version
			}
		}
		if version == "" {
			t.Fatal("fixture has no agent-instructions managed region")
		}
		content = "<!-- setup-context-driven:begin id=guide.agent-instructions version=" + version + " -->\n\n# Local managed override\n\n<!-- setup-context-driven:end id=guide.agent-instructions -->\n"
		want = append(want, baseline.MisplacedRepositoryRule{Path: misplacedTestCarrier, StartLine: 1, EndLine: 5, FirstLine: "# Local managed override", ManagedID: "guide.agent-instructions", MoveTo: baseline.SpecificRepositoryPath})
	}
	start := strings.Count(content, "\n") + 2
	content += "\n" + misplacedTestPassage
	writeBaselinePlanTestFile(t, repository, misplacedTestCarrier, content)
	want = append(want, baseline.MisplacedRepositoryRule{Path: misplacedTestCarrier, StartLine: start, EndLine: start + 1, FirstLine: "## Local rules", MoveTo: baseline.SpecificRepositoryPath})
	return want
}

func assertMisplacedBaselineJSON(t *testing.T, output string, want []baseline.MisplacedRepositoryRule) {
	t.Helper()
	var result map[string]json.RawMessage
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatal(err)
	}
	expected, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result["misplacedRepositoryRules"], expected) {
		t.Fatalf("misplacedRepositoryRules=%s, want %s", result["misplacedRepositoryRules"], expected)
	}
	var got []baseline.MisplacedRepositoryRule
	if err := json.Unmarshal(result["misplacedRepositoryRules"], &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("rules=%+v, want %+v", got, want)
	}
}

func TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	setCommandEnvironmentForTest(t, t.TempDir(), root)
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
		t.Fatalf("transcript first eight lines=%q, want %q", lines, want)
	}
}

func TestBaselineUpdateMisplacedRuleReportAcrossFailures(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"empty_array", "invalid_manifest", "declared_decision_failure", "carrier_read_failure"} {
		t.Run(name, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			setCommandEnvironmentForTest(t, t.TempDir(), repository)
			var want []baseline.MisplacedRepositoryRule
			if name != "empty_array" {
				want = misplacedBaselineFixture(t, repository, true)
			}
			args := []string{"baseline", "update", "--repo", repository, "--no-skills", "--no-history", "--format=json"}
			wantCode := exitOK
			switch name {
			case "empty_array":
				want = []baseline.MisplacedRepositoryRule{}
			case "invalid_manifest":
				writeBaselinePlanTestFile(t, repository, "docs/agents/setup-context.json", "{invalid\n")
				wantCode = exitPreflight
			case "declared_decision_failure":
				args = append(args, "--decision", "preservation.mode=greenfield")
				wantCode = exitPreflight
			case "carrier_read_failure":
				if err := os.Remove(filepath.Join(repository, misplacedTestCarrier)); err != nil {
					t.Fatal(err)
				}
				wantCode = exitUnverified
				want = []baseline.MisplacedRepositoryRule{}
			}
			var stdout, stderr bytes.Buffer
			code := runCLI(t, args, &stdout, &stderr)
			if code != wantCode {
				t.Fatalf("exit=%d, want %d; stdout=%s stderr=%s", code, wantCode, stdout.String(), stderr.String())
			}
			if name == "invalid_manifest" {
				var fields map[string]json.RawMessage
				if err := json.Unmarshal(stdout.Bytes(), &fields); err != nil {
					t.Fatal(err)
				}
				if _, exists := fields["misplacedRepositoryRules"]; exists {
					t.Fatalf("field before valid manifest: %s", stdout.String())
				}
				return
			}
			assertMisplacedBaselineJSON(t, stdout.String(), want)
			if name == "carrier_read_failure" {
				var result baselineUpdateResult
				if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				count := 0
				for _, warning := range result.Warnings {
					if warning.Code == "baseline.misplaced-rules.unread" {
						count++
					}
				}
				if count != 1 {
					t.Fatalf("detector warnings=%+v", result.Warnings)
				}
				args[len(args)-1] = "--format=text"
				stdout.Reset()
				stderr.Reset()
				if textCode := runCLI(t, args, &stdout, &stderr); textCode != wantCode || stderr.Len() != 0 {
					t.Fatalf("text exit=%d stderr=%s stdout=%s", textCode, stderr.String(), stdout.String())
				}
				warning := "Warning: baseline.misplaced-rules.unread: docs/agents/setup-context.json: read misplaced-rule carrier \"docs/agents/agent-instructions.md\": file does not exist"
				if !strings.Contains(stdout.String(), "\n"+warning+"\n") || strings.Count(stdout.String(), "Warning: baseline.misplaced-rules.unread:") != 1 {
					t.Fatalf("missing single whole warning line %q: %s", warning, stdout.String())
				}
			}
		})
	}
}
