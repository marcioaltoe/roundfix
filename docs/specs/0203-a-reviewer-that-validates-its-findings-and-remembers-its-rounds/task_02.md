---
task: task_02
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: pending
type: backend
complexity: high
---

# Task 02: A sealed validator dismisses a convention restatement only inside its region

## Overview

After task_01, an anchored finding always stands, including one that only
restates the delivery's designed order, such as the Daemon's status write. This
Task computes each convention's region from the head's tree. It asks a sealed,
tool-less validator prompt about the anchored findings that a region holds, or
that have no failure clause. It accepts a dismissal only for a rule the code
made eligible. Any validator failure leaves the findings standing.

The validator's input comes from the candidate diff, the head's tree and the
reviewer's answer. All three are untrusted and labelled so. The validator runs
on the local runtime of the selection that ran the review, through the
existing sealed path, which denies tools and closes its own session. Only
Roundfix reads its answer.

## Requirements

1. MUST add `conventionRegions` to `internal/cli/review_conventions.go`,
   implementing `_techspec.md` → Regions exactly. It reads files with
   `git show <head>:<path>` through the review's Git runner, uses the Spec
   roots `reviewCandidateSpecRoots` resolves, and requires every line of the
   anchor to lie in the region.
2. MUST compute eligibility as `_techspec.md` → Validation states: the
   `convention:Cn` rules whose region holds the anchor, plus `no-failure` when
   `reviewFindingHasFailureClause` is false. A finding with no eligible rule
   stands without a validator call. When no finding is eligible,
   `validation.validator` MUST stay `not-needed`.
3. MUST add `internal/cli/review_convention_validator.go` with
   `reviewSealedRunner` and `runConventionValidator`. It sends one sealed
   prompt per review with every eligible finding, on the runtime of the
   selection that ran the review. The prompt carries the input the TechSpec
   lists and labels it untrusted data. The answer is parsed strictly by the
   `roundfix/review-validation/v1` grammar.
4. MUST fail closed. A runner that does not implement `reviewSealedRunner`, a
   runner error, a timeout, `ToolUsed`, output that is not exactly the JSON
   object, a missing or duplicate ID, a `dismiss` with an ineligible rule, or
   a blank reason each make `validation.validator` `unavailable`, with reason
   `validator unavailable: <cause>`, and leave every eligible finding
   standing.
5. MUST record each accepted dismissal as `dismissed-by-validation` with the
   validator's rule and reason, and each `stands` verdict as `stands` with the
   validator's reason. Standard error and the outcome MUST follow task_01's
   rules, so Surface Transcript 1 is reproduced.
6. MUST NOT change the verdict grammar, the classifier or the reviewer call,
   and MUST NOT call the validator when the review outcome is `reviewed` or
   `blocked`.
7. MUST put the new tests in `internal/cli/review_convention_validator_test.go`.
   They use temporary repositories whose trees hold Spec directories, and a
   fake runner that implements `reviewSealedRunner` and records each sealed
   request. No test uses a real adapter.

## Subtasks

- [ ] Compute the four regions from the head's tree.
- [ ] Decide eligibility per finding.
- [ ] Build the sealed validator prompt and parse its answer strictly.
- [ ] Apply accepted dismissals and fail closed on everything else.
- [ ] Add a test for each acceptance criterion, with each negative case separate.

## Acceptance Criteria

- [ ] `TestReviewValidatorDismissesADaemonSettlementAsC2` reproduces Surface
      Transcript 1 without its `lineage` field, which task_03 and task_04 add. The sealed request names F1, its anchor and the eligible
      rule `convention:C2`, and marks the input untrusted.
- [ ] `TestReviewValidatorCannotDismissOutsideTheConventionRegion`: a finding
      with a failure clause, anchored on a Task's Requirements, stands without
      a validator call. A validator answer that dismisses it as `C2` when it is
      asked about another finding is refused as an ineligible rule.
- [ ] `TestReviewValidatorNeverAsksNoFailureForAFindingWithAFailureClause`:
      `no-failure` is never an eligible rule for a finding with a `Failure:`
      clause.
- [ ] `TestReviewValidatorFailureLeavesEveryAskedFindingStanding`: one subtest
      per cause in Requirement 4, each yielding `unavailable`, its reason, the
      findings standing, outcome `findings` and exit `1`.
- [ ] `TestReviewConventionRegionsAdmitAndRefuseAtTheirEdges`:
      - `C1` admits a path under `qa/` of an active and of an archived Spec,
        and refuses the Spec's `_prd.md`;
      - `C2` admits the front matter and a `## Result` section, and refuses a
        range that crosses into the next section;
      - `C3` admits the archive root only;
      - `C4` admits `_tasks.md` and front matter when every Task is pending
        and `qa/` is empty, and refuses both when one Task is `completed` or a
        QA Report exists.
- [ ] `TestReviewDisposeRefusesAConventionDismissal` reproduces Surface
      Transcript 7.
- [ ] task_01's tests and `TestReviewClassifiesVerdictVariants` still pass.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_conventions.go`
- creates: `internal/cli/review_validation.go`
- creates: `internal/cli/review_convention_validator.go`
- creates: `internal/cli/review_convention_validator_test.go`
- creates: `internal/cli/review_validation_test.go`
- instruction: `internal/agent/sealed.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewValidatorDismissesADaemonSettlementAsC2|TestReviewValidatorCannotDismissOutsideTheConventionRegion|TestReviewValidatorNeverAsksNoFailureForAFindingWithAFailureClause|TestReviewValidatorFailureLeavesEveryAskedFindingStanding|TestReviewConventionRegionsAdmitAndRefuseAtTheirEdges|TestReviewDisposeRefusesAConventionDismissal|TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne|TestReviewBlocksFindingsThatNameNoFileAndLine|TestReviewClassifiesVerdictVariants)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewValidatorDismissesADaemonSettlementAsC2 TestReviewValidatorCannotDismissOutsideTheConventionRegion TestReviewValidatorNeverAsksNoFailureForAFindingWithAFailureClause TestReviewValidatorFailureLeavesEveryAskedFindingStanding TestReviewConventionRegionsAdmitAndRefuseAtTheirEdges TestReviewDisposeRefusesAConventionDismissal TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne TestReviewBlocksFindingsThatNameNoFileAndLine TestReviewClassifiesVerdictVariants; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1–3; User Stories 1 and 3; Core Features 1, 4 and 5; Success Metrics 1 and 3
- [_techspec.md](_techspec.md) — Regions; Validation; Invariants 1 and 2; API Contracts 2–4; Surface Transcripts 1 and 7; Testing Approach 2; Build Order 2
- ADR-0196; ADR-0151; ADR-0153
