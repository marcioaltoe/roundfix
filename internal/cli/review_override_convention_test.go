// Suite: authorized QA Archive Override convention.
// Invariant: C5 requires a complete override record in the archived review head.
// Boundary IN: temporary Git repositories, convention regions, prompts and records.
// Boundary OUT: provider calls and semantic validator judgment.
package cli

import (
	"context"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

const overrideConventionPRD = "---\nstatus: archived\nqa_override: true\nqa_override_approval: Maintainer approval on 2026-10-02\nqa_override_reason: Environment-only partial QA\n---\n"

func overrideConventionRepository(t *testing.T, root, prd string) reviewRepository {
	t.Helper()
	repo := t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	slug := root + "/0001-example"
	if prd != "" {
		validatorWrite(t, repo, slug+"/_prd.md", prd)
	}
	validatorWrite(t, repo, slug+"/task_01.md", strings.Replace(validatorTaskBody, "status: completed", "status: failed", 1))
	validatorWrite(t, repo, slug+"/qa/report.md", "---\nverdict: partial\n---\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "test: archive override")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	return reviewRepository{Root: repo, Head: head, Base: head, SpecRoots: []string{"docs/specs", "docs/history/specs", "specs/_archived"}, Git: preflight.ExecGitRunner{}}
}

func TestOverrideArchiveIsEligibleForConventionC5(t *testing.T) {
	for _, root := range []string{"docs/history/specs", "specs/_archived"} {
		t.Run(root, func(t *testing.T) {
			repo := overrideConventionRepository(t, root, overrideConventionPRD)
			// Eligibility reads the pinned head even when the checkout removes authority.
			validatorWrite(t, repo.Root, root+"/0001-example/_prd.md", "---\nstatus: archived\n---\n")
			for _, file := range []string{"_prd.md", "task_01.md", "qa/report.md"} {
				t.Run(file, func(t *testing.T) {
					rules, err := conventionRegions(context.Background(), repo, reviewFindingAnchor{root + "/0001-example/" + file, 1, 1})
					if err != nil || !containsString(rules, "C5") || !containsString(rules, "C3") {
						t.Fatalf("rules=%v err=%v, want C5 and C3", rules, err)
					}
					if file == "task_01.md" && !containsString(rules, "C2") || file == "qa/report.md" && !containsString(rules, "C1") {
						t.Fatalf("existing region lost: %v", rules)
					}
				})
			}
		})
	}
}

func TestArchiveWithoutOverrideIsNotEligibleForC5(t *testing.T) {
	for _, test := range []struct{ name, root, prd string }{
		{"missing PRD", "docs/history/specs", ""},
		{"no override", "docs/history/specs", "---\nstatus: archived\n---\n"},
		{"false override", "docs/history/specs", strings.Replace(overrideConventionPRD, "qa_override: true", "qa_override: false", 1)},
		{"missing approval", "docs/history/specs", strings.Replace(overrideConventionPRD, "qa_override_approval: Maintainer approval on 2026-10-02\n", "", 1)},
		{"blank approval", "docs/history/specs", strings.Replace(overrideConventionPRD, "Maintainer approval on 2026-10-02", "'   '", 1)},
		{"missing reason", "docs/history/specs", strings.Replace(overrideConventionPRD, "qa_override_reason: Environment-only partial QA\n", "", 1)},
		{"blank reason", "docs/history/specs", strings.Replace(overrideConventionPRD, "Environment-only partial QA", "'   '", 1)},
		{"body only", "docs/history/specs", "---\nstatus: archived\n---\n" + overrideConventionPRD},
		{"unclosed front matter", "docs/history/specs", strings.TrimSuffix(overrideConventionPRD, "---\n")},
		{"active Spec", "docs/specs", overrideConventionPRD},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo := overrideConventionRepository(t, test.root, test.prd)
			// An uncommitted grant must never qualify the immutable candidate.
			validatorWrite(t, repo.Root, test.root+"/0001-example/_prd.md", overrideConventionPRD)
			rules, err := conventionRegions(context.Background(), repo, reviewFindingAnchor{test.root + "/0001-example/qa/report.md", 1, 1})
			if err != nil || containsString(rules, "C5") || !containsString(rules, "C1") {
				t.Fatalf("rules=%v err=%v, want C1 without C5", rules, err)
			}
		})
	}
}

func TestReviewPromptCarriesConventionC5(t *testing.T) {
	repo := overrideConventionRepository(t, "docs/history/specs", overrideConventionPRD)
	prompt := buildReviewPrompt(repo.Base, repo.Head, "diff")
	if !strings.Contains(prompt, "Delivery Conventions ("+deliveryConventionsVersion+")") {
		t.Fatal("prompt missing conventions version")
	}
	const c5 = "C5. A Spec archived through the QA Archive Override records `qa_override: true`, `qa_override_approval` and `qa_override_reason` in its archived `_prd.md` front matter; its QA Task keeps its observed status and its QA Report its observed verdict."
	if !strings.Contains(prompt, c5) {
		t.Fatal("prompt missing C5 contract")
	}
	record, code := validateReviewFindingAnchors(reviewRecord{Outcome: reviewOutcomeReviewed}, "diff")
	if code != exitOK || record.Validation == nil || record.Validation.Conventions != deliveryConventionsVersion {
		t.Fatalf("record validation=%+v exit=%d", record.Validation, code)
	}
}

func TestConventionC5MalformedRecordFailsClosed(t *testing.T) {
	for _, prd := range []string{
		"---\nqa_override: [\n---\n",
		overrideConventionPRD[:len(overrideConventionPRD)-4] + "qa_override: false\n---\n",
	} {
		t.Run(prd, func(t *testing.T) {
			repo := overrideConventionRepository(t, "docs/history/specs", prd)
			rules, err := conventionRegions(context.Background(), repo, reviewFindingAnchor{"docs/history/specs/0001-example/qa/report.md", 1, 1})
			if err == nil || containsString(rules, "C5") {
				t.Fatalf("malformed record rules=%v err=%v", rules, err)
			}
		})
	}
}

func TestConventionC5ValidatorHonorsEligibility(t *testing.T) {
	const reportPath = "docs/history/specs/0001-example/qa/report.md"
	for _, test := range []struct {
		name, prd, reason string
		wantDismissed     bool
	}{
		{"authorized restatement", overrideConventionPRD, "restates the authorized override", true},
		{"unauthorized rule", "---\nstatus: archived\n---\n", "restates the override", false},
		{"blank reason", overrideConventionPRD, " ", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo := overrideConventionRepository(t, "docs/history/specs", test.prd)
			record, _ := validateReviewFindingAnchors(reviewRecord{
				Outcome:      reviewOutcomeFindings,
				FindingItems: []reviewFinding{{ID: "F1", Text: reportPath + ":2 Failure: archived QA remains partial"}},
			}, "diff --git a/"+reportPath+" b/"+reportPath+"\n--- /dev/null\n+++ b/"+reportPath+"\n@@ -0,0 +1,3 @@\n+---\n+verdict: partial\n+---\n")
			runner := &conventionValidatorRunner{sealedResult: agent.SealedPromptResult{Output: validatorAnswer(validatorVerdict{"F1", "dismiss", "convention:C5", test.reason})}}
			record, code := validateReviewConventions(context.Background(), repo, record, runner, agent.RuntimeSpec{})
			if got := reviewFindingDismissedByValidation(record.FindingItems[0]); got != test.wantDismissed {
				t.Fatalf("dismissed=%v record=%+v", got, record)
			}
			if test.wantDismissed && (code != exitOK || record.Validation.Validator != "ran") || !test.wantDismissed && (code != exitRunFailed || record.Validation.Validator != "unavailable") {
				t.Fatalf("exit=%d validation=%+v", code, record.Validation)
			}
		})
	}
}
