// Suite: release provenance at the workflow script boundary.
// Invariant: only a verified squash of an up-to-date PR may reach publication.
// Boundary IN: local Git history and recorded GitHub API pages.
// Boundary OUT: network access, real GitHub credentials and publication.
package verifyselect_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestReleaseProvenanceAcceptsOnlyAVerifiedSquashMerge(t *testing.T) {
	t.Parallel()
	for _, scenario := range []struct{ name, reason string }{
		{"squash_merge_of_a_branch_that_held_main", ""},
		{"manual_dispatch_evaluates_the_release_commit", ""},
		{"merged_behind_main", "the pull request merged behind main"},
		{"failed_check", "the check did not pass or is missing"},
		{"missing_check", "the check did not pass or is missing"},
		{"commit_no_pull_request_merged", "the commit is not the squash merge of one pull request into main"},
		{"tag_on_another_commit", "the tag is not on the release commit"},
		{"github_unreadable", "GitHub could not be read"},
		// Adversarial answers must not borrow a different run's success.
		{"older_success_newer_failure", "the check did not pass or is missing"},
		{"older_failure_newer_success", ""},
		{"foreign_app_check", "the check did not pass or is missing"},
		{"wrong_head_check", "the check did not pass or is missing"},
		{"incomplete_check", "the check did not pass or is missing"},
		{"different_tree", "the pull request merged behind main"},
		{"wrong_base", "the commit is not the squash merge of one pull request into main"},
		{"wrong_merge_commit", "the commit is not the squash merge of one pull request into main"},
		{"unmerged_pull_request", "the commit is not the squash merge of one pull request into main"},
		{"ambiguous_pull_requests", "the commit is not the squash merge of one pull request into main"},
		{"malformed_pulls", "GitHub could not be read"},
		{"malformed_checks", "GitHub could not be read"},
		{"checks_unreadable", "GitHub could not be read"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			event := "push"
			if strings.HasPrefix(scenario.name, "manual_dispatch") {
				event = "workflow_dispatch"
			}
			h := newReleaseWorkflowHarness(t, releaseWorkflowCase{eventName: event, refName: "v" + releaseTestVersion})
			// Keep harness files outside the Git tree so ancestry and tree checks are real.
			gittest.InitRepo(t, h.root, "-b", "main")
			writeFile(t, h.root, "dist/npm/roundfix/package.json", `{"version":"0.0.0"}`)
			gittest.Run(t, h.root, "add", "-f", "dist/npm/roundfix/package.json")
			gittest.Run(t, h.root, "commit", "-m", "initial")
			base := strings.TrimSpace(gittest.Run(t, h.root, "rev-parse", "HEAD"))
			gittest.Run(t, h.root, "checkout", "-b", "chore/release-fixture")
			writeFile(t, h.root, "dist/npm/roundfix/package.json", `{"version":"`+releaseTestVersion+`"}`)
			gittest.Run(t, h.root, "add", "-f", "dist/npm/roundfix/package.json")
			gittest.Run(t, h.root, "commit", "-m", "release content")
			head := strings.TrimSpace(gittest.Run(t, h.root, "rev-parse", "HEAD"))
			gittest.Run(t, h.root, "checkout", "main")
			if scenario.name == "merged_behind_main" {
				// Same final tree, but main's new tip was never in the PR branch.
				gittest.Run(t, h.root, "commit", "--allow-empty", "-m", "main advanced")
			}
			gittest.Run(t, h.root, "merge", "--squash", head)
			if scenario.name == "different_tree" {
				writeFile(t, h.root, "extra", "merge-only content")
				gittest.Run(t, h.root, "add", "extra")
			}
			gittest.Run(t, h.root, "commit", "-m", "squash release")
			commit := strings.TrimSpace(gittest.Run(t, h.root, "rev-parse", "HEAD"))
			gittest.Run(t, h.root, "tag", "v"+releaseTestVersion)
			if scenario.name == "tag_on_another_commit" || event == "workflow_dispatch" {
				writeFile(t, h.root, "later", "later content")
				gittest.Run(t, h.root, "add", "later")
				gittest.Run(t, h.root, "commit", "-m", "after release")
				if event == "push" {
					gittest.Run(t, h.root, "tag", "-f", "v"+releaseTestVersion)
				}
			}
			pr := map[string]any{"number": 7, "merged_at": "2026-11-25T00:00:00Z", "merge_commit_sha": commit, "base": map[string]any{"ref": "main", "repo": map[string]any{"full_name": "roundfix/fixture"}}, "head": map[string]any{"sha": head}}
			prs := []any{pr}
			switch scenario.name {
			case "wrong_base":
				pr["base"] = map[string]any{"ref": "develop"}
			case "wrong_merge_commit":
				pr["merge_commit_sha"] = base
			case "unmerged_pull_request":
				pr["merged_at"] = nil
			case "ambiguous_pull_requests":
				prs = []any{pr, pr}
			case "commit_no_pull_request_merged":
				prs = []any{}
			}
			// --slurp returns one array per API page; place the match on a second page.
			writeAnswer := func(name string, v any) {
				t.Helper()
				b, err := json.Marshal(v)
				if err != nil {
					t.Fatal(err)
				}
				writeFile(t, h.fixtures, name, string(b))
			}
			writeAnswer("pulls.json", []any{[]any{}, prs})
			check := func(id int, result string) map[string]any {
				return map[string]any{"id": id, "name": "Verification gate", "app": map[string]any{"slug": "github-actions"}, "head_sha": head, "status": "completed", "conclusion": result}
			}
			checks := []any{check(10, "success")}
			switch scenario.name {
			case "failed_check":
				checks = []any{check(10, "failure")}
			case "missing_check":
				checks = []any{}
			case "older_success_newer_failure":
				checks = []any{check(20, "failure"), check(10, "success")}
			case "older_failure_newer_success":
				checks = []any{check(10, "failure"), check(20, "success")}
			case "foreign_app_check":
				checks[0].(map[string]any)["app"] = map[string]any{"slug": "foreign"}
			case "wrong_head_check":
				checks[0].(map[string]any)["head_sha"] = base
			case "incomplete_check":
				checks[0].(map[string]any)["status"] = "in_progress"
			}
			writeAnswer("checks.json", []any{map[string]any{"check_runs": []any{}}, map[string]any{"check_runs": checks}})
			if scenario.name == "malformed_pulls" {
				writeFile(t, h.fixtures, "pulls.json", "invalid")
			}
			if scenario.name == "malformed_checks" {
				writeFile(t, h.fixtures, "checks.json", `{"message":"unavailable"}`)
			}
			stub := fmt.Sprintf(`#!/bin/sh
set -eu
[ "$GH_TOKEN" = fixture-token ] || exit 91
[ "$1" = api ] && [ "$2" = -H ] && [ "$3" = 'X-GitHub-Api-Version: 2022-11-28' ] && [ "$4" = --paginate ] && [ "$5" = --slurp ] || exit 92
case "$6" in
 'repos/roundfix/fixture/commits/%s/pulls?per_page=100') [ '%s' != github_unreadable ] || exit 1; cat "$RELEASE_FIXTURES/pulls.json";;
 'repos/roundfix/fixture/commits/%s/check-runs?per_page=100&filter=all') [ '%s' != checks_unreadable ] || exit 1; cat "$RELEASE_FIXTURES/checks.json";;
 *) exit 93;;
esac
`, commit, scenario.name, head, scenario.name)
			writeFile(t, h.root, "bin/gh", stub)
			if err := os.Chmod(filepath.Join(h.root, "bin/gh"), 0o755); err != nil {
				t.Fatal(err)
			}
			h.env = append(h.env, "GH_TOKEN=fixture-token", "GITHUB_REPOSITORY=roundfix/fixture", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null", "GIT_OPTIONAL_LOCKS=0")
			output, code := h.jobStep(t, "verify", "Release provenance")
			want := fmt.Sprintf("provenance: %s is the squash merge of pull request #7, and Verification gate passed on its tree\n", commit)
			wantCode := 0
			if scenario.reason != "" {
				want = "::error::provenance: " + scenario.reason + "\n"
				wantCode = 1
			}
			if output != want || code != wantCode {
				t.Fatalf("exit=%d want=%d, output=%q want=%q", code, wantCode, output, want)
			}
		})
	}
}
