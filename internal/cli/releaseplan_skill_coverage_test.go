package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/releaseplan"
	"roundfix/internal/skillcoverage"
)

// Suite: committed skill coverage at the release CLI boundary.
// Boundary IN: temporary Git repositories with authored maps and recorded surfaces.
// Boundary OUT: decision state, text, JSON, exit codes and repository immutability.
const coverageFixtureSkill = ".agents/skills/roundfix/references/implement.md"
const coverageFixtureSurface = "command implement"

func coverageFixtureMap(review string, uncovered bool) skillcoverage.Map {
	entry := skillcoverage.Surface{ID: coverageFixtureSurface, Skills: []string{coverageFixtureSkill}, Review: review}
	if uncovered {
		entry.Skills = nil
		entry.Uncovered = "fixture has no skill"
	}
	return skillcoverage.Map{SchemaVersion: skillcoverage.MapSchema, Surfaces: []skillcoverage.Surface{entry}}
}
func writeCoverageFixture(t *testing.T, repo string, m skillcoverage.Map, fingerprint string) {
	t.Helper()
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeReleasePlanFile(t, repo, skillcoverage.MapPath, string(data))
	surfaces := map[string]string{}
	if fingerprint != "" {
		surfaces[coverageFixtureSurface] = skillcoverage.Fingerprint([]byte(fingerprint))
	}
	writeReleasePlanFile(t, repo, skillcoverage.RecordPath, string(skillcoverage.EncodeRecord(skillcoverage.Record{SchemaVersion: skillcoverage.RecordSchema, Surfaces: surfaces})))
}
func newCoverageFixture(t *testing.T, review string, uncovered, declared bool) string {
	t.Helper()
	repo := newReleasePlanCommandRepo(t, "v0.4.0")
	writeReleasePlanFile(t, repo, coverageFixtureSkill, "initial skill\n")
	if declared {
		writeCoverageFixture(t, repo, coverageFixtureMap(review, uncovered), "before")
	}
	gitReleasePlan(t, repo, "add", "-A")
	commitReleasePlan(t, repo, "docs: describe fixture")
	gitReleasePlan(t, repo, "tag", "-f", "v0.4.0")
	return repo
}
func commitCoverageFixture(t *testing.T, repo, subject string) {
	t.Helper()
	gitReleasePlan(t, repo, "add", "-A")
	commitReleasePlan(t, repo, subject)
}
func checkCoverageFixture(t *testing.T, repo, status string, blocking bool, wantCode int, args ...string) releasePlanJSON {
	t.Helper()
	code, output, stderr := runReleasePlanCommandInRepo(t, repo, append(args, "--format=json")...)
	plan := decodeReleasePlanJSON(t, output)
	if code != wantCode || stderr != "" || plan.Checks.SkillCoverage.Status != status || plan.Checks.SkillCoverage.Blocking != blocking || plan.SchemaVersion != releaseplan.SchemaVersion {
		t.Fatalf("exit=%d stderr=%q plan=%+v check=%+v", code, stderr, plan, plan.Checks.SkillCoverage)
	}
	return plan
}
func assertCoverageText(t *testing.T, output, line string) {
	t.Helper()
	if strings.Count(output, "\n"+line+"\n") != 1 {
		t.Fatalf("missing exact transcript line %q in %s", line, output)
	}
}
func TestReleasePlanBlocksALaggingSurface(t *testing.T) {
	for _, subject := range []string{"fix: correct output", "feat: add output"} {
		t.Run(subject, func(t *testing.T) {
			repo := newCoverageFixture(t, "", false, true)
			writeCoverageFixture(t, repo, coverageFixtureMap("", false), "after")
			commitCoverageFixture(t, repo, subject)
			plan := checkCoverageFixture(t, repo, "behind", true, exitUnverified)
			wantState, wantVersion := releaseplan.StateReady, "v0.4.1"
			if strings.HasPrefix(subject, "feat:") {
				wantState, wantVersion = releaseplan.StateApprovalRequired, "v0.5.0"
			}
			if plan.State != wantState || plan.ProposedVersion != wantVersion {
				t.Fatalf("decision changed: %+v", plan)
			}
			lagging := plan.Checks.SkillCoverage.Lagging
			if len(lagging) != 1 || lagging[0].Surface != coverageFixtureSurface || lagging[0].Change != "changed" || len(lagging[0].Skills) != 1 || lagging[0].Skills[0] != coverageFixtureSkill {
				t.Fatalf("lagging=%+v", lagging)
			}
			code, output, stderr := runReleasePlanCommandInRepo(t, repo)
			if code != exitUnverified || stderr != "" {
				t.Fatalf("exit=%d stderr=%s", code, stderr)
			}
			assertCoverageText(t, output, "Release blocked: skill-coverage")
			assertCoverageText(t, output, "Next action: "+releasePlanBlockedNextAction)
			assertCoverageText(t, output, "skill-coverage: behind: 1 lagging surface(s) since v0.4.0; next: "+releasePlanSkillCoverageNextAction)
			assertCoverageText(t, output, "- lagging: command implement (changed; covering skills unchanged: .agents/skills/roundfix/references/implement.md)")
			if wantState == releaseplan.StateApprovalRequired && (!plan.Approval.Required || plan.Approval.Question == "" || !strings.Contains(output, "Approval question: "+plan.Approval.Question)) {
				t.Fatalf("approval lost: %+v", plan.Approval)
			}
			// A later skill edit clears HEAD but must not clear the earlier --to range.
			target := gitReleasePlanOutput(t, repo, "rev-parse", "HEAD")
			writeReleasePlanFile(t, repo, coverageFixtureSkill, "updated skill\n")
			commitCoverageFixture(t, repo, "docs: describe output")
			checkCoverageFixture(t, repo, "behind", true, exitUnverified, "--to="+strings.TrimSpace(target))
		})
	}
}
func TestReleasePlanSkillCoverageIsCurrentWhenItsSkillChanges(t *testing.T) {
	for _, tc := range []struct {
		name, beforeReview, afterReview, fingerprint, outcome string
		skill, uncovered                                      bool
	}{
		{name: "skill", fingerprint: "after", outcome: "described", skill: true},
		{name: "review", afterReview: "2026-10-08 — wording still applies", fingerprint: "after", outcome: "reviewed"},
		{name: "unchanged review", beforeReview: "2026-10-08 — wording still applies", afterReview: "2026-10-08 — wording still applies", fingerprint: "after", outcome: "lagging"},
		{name: "removed", outcome: "lagging"},
		{name: "uncovered", fingerprint: "after", outcome: "uncovered", uncovered: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := newCoverageFixture(t, tc.beforeReview, tc.uncovered, true)
			m := coverageFixtureMap(tc.afterReview, tc.uncovered)
			if tc.fingerprint == "" {
				m.Surfaces = []skillcoverage.Surface{}
			}
			writeCoverageFixture(t, repo, m, tc.fingerprint)
			if tc.skill {
				writeReleasePlanFile(t, repo, coverageFixtureSkill, "updated skill\n")
			}
			commitCoverageFixture(t, repo, "fix: correct output")
			status, code := "current", exitOK
			if tc.outcome == "lagging" {
				status, code = "behind", exitUnverified
			}
			plan := checkCoverageFixture(t, repo, status, tc.outcome == "lagging", code)
			if plan.State != releaseplan.StateReady || plan.ProposedVersion != "v0.4.1" {
				t.Fatalf("decision changed: %+v", plan)
			}
			textCode, output, stderr := runReleasePlanCommandInRepo(t, repo)
			if textCode != code || stderr != "" {
				t.Fatalf("exit=%d stderr=%s", textCode, stderr)
			}
			if tc.outcome != "lagging" {
				d, r, u := "0", "0", "0"
				switch tc.outcome {
				case "described":
					d = "1"
				case "reviewed":
					r = "1"
				case "uncovered":
					u = "1"
				}
				assertCoverageText(t, output, "skill-coverage: current: 1 changed surface(s) since v0.4.0: "+d+" described, "+r+" reviewed, "+u+" uncovered")
				assertCoverageText(t, output, "Next action: release may proceed for v0.4.1 after independent release verification.")
				if plan.Checks.SkillCoverage.NextAction != "" || len(plan.Checks.SkillCoverage.Lagging) != 0 {
					t.Fatalf("unexpected action: %+v", plan.Checks.SkillCoverage)
				}
			} else if len(plan.Checks.SkillCoverage.Lagging) != 1 {
				t.Fatal("missing lagging surface")
			}
		})
	}
}
func TestReleasePlanSkillCoverageNeverBlocksWithoutAMap(t *testing.T) {
	for _, tc := range []struct {
		name, subject, status string
		baseMap, targetMap    bool
	}{
		{"absent", "fix: correct output", "not_declared", false, false},
		{"introduced", "fix: correct output", "introduced", false, true},
		{"no release", "docs: describe output", "behind", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// The same range without coverage establishes its decision.
			args := []string{}
			if tc.status == "behind" {
				args = []string{"--impact=none", "--reason=maintenance only"}
			}
			control := newReleasePlanCommandRepo(t, "v0.4.0", releasePlanCommandCommit{subject: tc.subject, paths: []string{"docs/references/fixture.json"}})
			_, controlOutput, _ := runReleasePlanCommandInRepo(t, control, append(args, "--format=json")...)
			want := decodeReleasePlanJSON(t, controlOutput)
			repo := newCoverageFixture(t, "", false, tc.baseMap)
			if tc.targetMap {
				writeCoverageFixture(t, repo, coverageFixtureMap("", false), "after")
			} else {
				writeReleasePlanFile(t, repo, "docs/references/fixture.json", "changed")
			}
			commitCoverageFixture(t, repo, tc.subject)
			plan := checkCoverageFixture(t, repo, tc.status, false, exitOK, args...)
			if plan.State != want.State || plan.ProposedVersion != want.ProposedVersion {
				t.Fatalf("decision changed: %+v want %+v", plan, want)
			}
			code, output, stderr := runReleasePlanCommandInRepo(t, repo, args...)
			if code != exitOK || stderr != "" || strings.Contains(output, "Release blocked:") {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, output, stderr)
			}
			if tc.status == "not_declared" {
				assertCoverageText(t, output, "skill-coverage: not_declared: the repository has no Skill Coverage Map")
			}
			if tc.status == "introduced" {
				assertCoverageText(t, output, "skill-coverage: introduced: the Skill Coverage Map is new since v0.4.0; no surface is compared")
			}
		})
	}
}
func TestReleasePlanSkillCoverageFailsClosedOnAnUnreadableMap(t *testing.T) {
	for _, malformed := range []bool{true, false} {
		t.Run(map[bool]string{true: "malformed map", false: "missing record"}[malformed], func(t *testing.T) {
			repo := newCoverageFixture(t, "", false, true)
			if malformed {
				writeReleasePlanFile(t, repo, skillcoverage.MapPath, "{broken")
			} else if err := os.Remove(filepath.Join(repo, skillcoverage.RecordPath)); err != nil {
				t.Fatal(err)
			}
			commitCoverageFixture(t, repo, "fix: correct output")
			before := snapshotReleasePlanRepo(t, repo)
			plan := checkCoverageFixture(t, repo, "failed", true, exitUnverified)
			if plan.State != releaseplan.StateReady || plan.ProposedVersion != "v0.4.1" || plan.Checks.SkillCoverage.Detail == "" {
				t.Fatalf("plan=%+v", plan)
			}
			code, output, stderr := runReleasePlanCommandInRepo(t, repo)
			if code != exitUnverified || stderr != "" {
				t.Fatalf("exit=%d stderr=%s", code, stderr)
			}
			assertCoverageText(t, output, "Release blocked: skill-coverage")
			assertReleasePlanRepoUnchanged(t, repo, before)
		})
	}
}
