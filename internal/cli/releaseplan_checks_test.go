package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/releaseplan"
)

// Suite: advisory release plan checks
// Boundary IN: release plan CLI and real local repositories.
// Boundary OUT: release decisions, repository writes, and external services.
func TestReleasePlanReportsTheSkillsAndBaselineChecks(t *testing.T) {
	repo := newReleasePlanCommandRepo(t, "v0.4.0", releasePlanCommandCommit{subject: "fix: correct output", paths: []string{"internal/cli/output.go"}})
	code, output, stderr := runReleasePlanCommandInRepo(t, repo)
	if code != exitOK || stderr != "" {
		t.Fatalf("exit=%d stderr=%s", code, stderr)
	}
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "Next action:") {
			if i+2 >= len(lines) || !strings.HasPrefix(lines[i+1], "skills: failed: missing:") || !strings.HasPrefix(lines[i+2], "baseline: action_required: the repository has no Setup Manifest; next: ") {
				t.Fatalf("Surface Transcript 1: %s", output)
			}
		}
	}
	for _, prefix := range []string{"skills:", "baseline:"} {
		if strings.Count(output, "\n"+prefix) != 1 {
			t.Fatalf("expected one %s line: %s", prefix, output)
		}
	}
	_, output, _ = runReleasePlanCommandInRepo(t, repo, "--format=json")
	checks := decodeReleaseChecks(t, output)
	if checks.Skills.Status != "failed" || checks.Baseline.Status != "action_required" || checks.Baseline.Detail != "the repository has no Setup Manifest" {
		t.Fatalf("checks=%+v", checks)
	}
	for _, check := range []releaseChecksTestValue{checks.Skills, checks.Baseline} {
		if check.NextAction != "complete the skills and guides check in the release runbook before the release Pull Request" {
			t.Fatalf("next action=%q", check.NextAction)
		}
	}
	if strings.Index(output, `"checks"`) < strings.Index(output, `"approval"`) || strings.Index(output, `"checks"`) > strings.Index(output, `"changes"`) {
		t.Fatal("checks JSON order")
	}
	help := commandUsage("release plan")
	if !strings.Contains(help, "A range plan also reports the skills and baseline checks read-only; they never change the decision state, the proposed version, or the exit code.") {
		t.Fatal("missing help sentence")
	}
}

func TestReleasePlanChecksNeverChangeTheDecision(t *testing.T) {
	for _, tc := range []struct {
		subject, path, base, state, version string
		code                                int
		args                                []string
	}{
		{"fix: correct output", "internal/cli/output.go", "v0.4.0", "ready", "v0.4.1", exitOK, nil},
		{"feat: add output", "internal/cli/output.go", "v0.4.0", "approval_required", "v0.5.0", exitUnverified, nil},
		{"feat!: replace output", "internal/cli/output.go", "v1.4.0", "approval_required", "v2.0.0", exitUnverified, nil},
		{"fix!: replace output", "internal/cli/output.go", "v0.4.0", "approval_required", "v0.5.0", exitUnverified, nil},
		{"docs: describe output", "docs/specs/fixture/output.md", "v0.4.0", "no_release", "", exitOK, nil},
		{"chore: tune output", "internal/cli/output.go", "v0.4.0", "manual_classification_required", "", exitUnverified, nil},
		{"chore: tune output", "internal/cli/output.go", "v0.4.0", "approval_required", "v0.5.0", exitUnverified, []string{"--impact=minor", "--reason=public behavior changed"}},
	} {
		t.Run(tc.subject+tc.base+tc.state, func(t *testing.T) {
			repo := newReleasePlanCommandRepo(t, tc.base, releasePlanCommandCommit{subject: tc.subject, paths: []string{tc.path}})
			code, output, stderr := runReleasePlanCommandInRepo(t, repo, append([]string{"--format=json"}, tc.args...)...)
			plan := decodeReleasePlanJSON(t, output)
			checks := decodeReleaseChecks(t, output)
			if code != tc.code || string(plan.State) != tc.state || plan.ProposedVersion != tc.version || stderr != "" || checks.Skills.Status != "failed" {
				t.Fatalf("exit=%d plan=%+v checks=%+v stderr=%s", code, plan, checks, stderr)
			}
			if plan.SchemaVersion != releaseplan.SchemaVersion {
				t.Fatal("schema changed")
			}
		})
	}
}

func TestReleasePlanChecksOnAnAdoptedRepository(t *testing.T) {
	for _, stale := range []bool{false, true} {
		t.Run(fmt.Sprint(stale), func(t *testing.T) {
			repo := newBaselineUpdateRepository(t)
			if stale {
				repo = staleBaselineUpdateRepository(t)
			}
			copyReleaseCheckSkills(t, repo)
			commitBaselinePlanTestRepository(t, repo)
			addReleaseCheckRange(t, repo)
			code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--format=json")
			if code != exitOK || stderr != "" {
				t.Fatalf("exit=%d output=%s stderr=%s", code, output, stderr)
			}
			checks := decodeReleaseChecks(t, output)
			if checks.Skills.Status != "ok" || checks.Skills.NextAction != "" {
				t.Fatalf("skills=%+v", checks.Skills)
			}
			if !stale {
				if checks.Baseline.Status != "current" || checks.Baseline.NextAction != "" {
					t.Fatalf("Surface Transcript 2: %+v", checks)
				}
				if strings.Contains(output, `"nextAction"`) {
					t.Fatal("empty nextAction must be omitted")
				}
			} else {
				result, _, stderr, _ := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repo, "--no-skills", "--format=json")
				want := fmt.Sprintf("%d file change(s)", len(result.FileChanges))
				if len(result.HistoryMoves) > 0 {
					want += fmt.Sprintf(", %d history move(s)", len(result.HistoryMoves))
				}
				if stderr != "" || checks.Baseline.Status != "plan_ready" || checks.Baseline.Detail != want {
					t.Fatalf("baseline=%+v update=%+v stderr=%s", checks.Baseline, result, stderr)
				}
			}
		})
	}
}

func TestReleasePlanChecksWriteNothing(t *testing.T) {
	repo := staleBaselineUpdateRepository(t)
	copyReleaseCheckSkills(t, repo)
	commitBaselinePlanTestRepository(t, repo)
	addReleaseCheckRange(t, repo)
	before := snapshotReleasePlanRepo(t, repo)
	for _, format := range []string{"text", "json"} {
		code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--format="+format)
		if code != exitOK || stderr != "" || !strings.Contains(output, "plan_ready") {
			t.Fatalf("exit=%d output=%s stderr=%s", code, output, stderr)
		}
		assertReleasePlanRepoUnchanged(t, repo, before)
	}
}

type releaseChecksTestValue struct {
	Status     string
	Detail     string
	NextAction string
}
type releaseChecksTestOutput struct {
	Skills   releaseChecksTestValue
	Baseline releaseChecksTestValue
}

func decodeReleaseChecks(t *testing.T, output string) releaseChecksTestOutput {
	t.Helper()
	var value struct{ Checks *releaseChecksTestOutput }
	if err := json.Unmarshal([]byte(output), &value); err != nil {
		t.Fatal(err)
	}
	if value.Checks == nil {
		t.Fatalf("missing checks: %s", output)
	}
	return *value.Checks
}
func addReleaseCheckRange(t *testing.T, repo string) {
	t.Helper()
	gitReleasePlan(t, repo, "tag", "v0.4.0")
	writeReleasePlanFile(t, repo, "internal/output.go", "package output\n")
	gitReleasePlan(t, repo, "add", "-A")
	commitReleasePlan(t, repo, "fix: correct output")
}
func copyReleaseCheckSkills(t *testing.T, repo string) {
	t.Helper()
	lock, err := os.ReadFile(filepath.Join(baselineDocumentationRepoRoot(), "skills-lock.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "skills-lock.json"), lock, 0644); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(baselineDocumentationRepoRoot(), ".agents", "skills")
	err = filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(repo, ".agents", "skills", relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
