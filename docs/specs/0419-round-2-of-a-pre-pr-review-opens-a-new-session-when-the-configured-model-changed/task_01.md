---
task: task_01
spec: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
status: pending
type: backend
complexity: medium
---

# Task 01: Round 2 continues the reviewer session only on the recorded model, and the rule is documented

## Overview

This Task delivers `_techspec.md` → Build Order 1. It answers the finding reported by fiscus on 2026-11-23 (ai-memory message 01a20a55-66b6-7b31-bd03-ed4f3d6a1fbc): after the review profile's model changed from `gpt-6-luna` to `gpt-6.1-sol`, round 2 of a Reviewer Lineage continued the round-1 Agent Session, Codex put a notice about the changed model before `No findings.`, and the item parked as `review-blocked: unclassifiable agent output: no-findings verdict does not account for other content`. After this Task each round that sends a prompt records the model of its session in `lineage.model`, and round 2 continues the round-1 session only when the selection at the recorded index still names that model; otherwise it ends the round-1 session and opens a new one with the round-2 delta prompt. ADR-0451 records the rule, ADR-0197 gains its amendment notice, and the glossary, the review guide and the Roundfix skill's review reference state it. It is verifiable on its own through three new tests of `internal/cli` over a fake runner, three new documentation contract tests, and the existing review, skills and documentation contracts.

## Requirements

1. MUST implement `_techspec.md` → Interfaces and API Contracts 1 and 2 in `internal/cli/review_lineage.go` and `internal/cli/review.go`:
   - `reviewLineage` gains `Model string` with the JSON name `model` and `omitempty`, declared after `Selection`; no other field moves;
   - the block that records the session of a round whose prompt was sent also records the model of the runtime the prompt went to; a record written without a prompt, a `ceiling-closed` record among them, gains no model;
   - the round-2 branch that prepares the recorded session name runs only when, in addition to today's condition, the prior lineage holds a model and it equals, after trimming, the model of the Agent Selection at the recorded index (invariant 1). Every other case reaches the two steps that follow the branch today, unchanged: the recorded session is ended, then the selection loop starts from the preferred selection with a new random name and the round-2 delta prompt (invariants 2 and 3). This follows the Delegated Decision of 2026-11-23: What does round 2 compare before it continues the round-1 session? (the model only); the Delegated Decision of 2026-11-23: A round-1 record written before this fix holds no model. Does round 2 continue its session or open a new one? (a new session); and the Delegated Decision of 2026-11-23: Which prompt does the new session of round 2 receive? (the round-2 delta prompt).
   The classifier of the reviewer's answer keeps its code (invariant 5, ADR-0174), and no output line, flag, exit code, help text or key of Project Config is added (invariant 6; Delegated Decision of 2026-11-23: Does the review print why it did not continue the round-1 session?, answered with no new output line). No other Go source file changes.
2. MUST add `internal/cli/review_model_change_test.go` with the fake runner and the three tests of `_techspec.md` → Testing Approach, built on `newReviewCommandFixture`, `runLineageReview`, `commitLineageFile`, `disposeLineageFinding` and `sessionReviewRunner`, as `references/sources.md`, section 6, probed them:
   - `modelNoticeReviewRunner`, which remembers the model each session name was first prepared with and puts the notice of `references/sources.md`, section 1, with the two model names, before the answer of a session name prepared again with another model;
   - `TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged`, asserting everything the Testing Approach lists for it, with the subtest `the_classifier_still_refuses_the_reported_answer`, which gives the reported answer to a round 1 as the reviewer's only message and asserts exit `2`, the outcome `blocked` and the reason `unclassifiable agent output: no-findings verdict does not account for other content`;
   - `TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel`;
   - `TestReviewRecordsTheModelOfTheSessionOfEachRound`, with the subtests `preferred_then_continued` and `fallback`.
   Each test reads `lineage.model` from the JSON the command prints or from the record file, so the file compiles on the unchanged tree and fails there by assertion. No test reaches the network, reads the real home directory, starts an Agent session or waits on elapsed time.
3. MUST keep these existing tests unchanged and passing: `TestReviewRoundTwoContinuesTheRecordedSession`, `TestReviewRoundTwoContinuesTheRecordedFallbackSelection`, `TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails`, `TestReviewLineageChangeEndsTheOpenSession` and `TestACeilingRecordNamesNoSelection`, whose pinned `ceiling-closed` lineage keeps its bytes. No existing test, golden or fixture of `internal/cli` or `internal/delivery` changes.
4. MUST create `docs/adr/0451-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed.md`, the decision record of `_techspec.md` → API Contract 3, with the title `Round 2 of a pre-PR review opens a new session when the configured model changed`, the lifecycle header of `docs/agents/docs-layout.md` with `status: accepted` and, on the next line, `knowledge: true`, the report of the Fiscus of 2026-11-23 as its context, API Contracts 1 and 2 as its rules, the five options `_techspec.md` → Decisions gives up, and the statement that it amends ADR-0197 in part. It names no client other than the reporting project and cites no Spec folder path. When 0451 is no longer the next free number on the tree this Task runs on, stop and report it; do not take another number.
5. MUST edit by hand, as API Contract 3 states: ADR-0197 gains `ADR-0451` at the end of its `amended_by` list, takes ADR-0451's `created_at` as its `updated_at`, and gains one paragraph after the notice of ADR-0391 that begins with the bold lead "Amended by", ADR-0451 as a link to its file and the date of its `created_at` in parentheses, in the form of the notices above it; every other line of ADR-0197 keeps its bytes. `.secondbrain-export` gains the ADR's path once, in ascending order, inside the Standard Export List.
6. MUST revise the glossary entry **Reviewer Lineage** in `CONTEXT.md` through `domain-modeling`, as `_techspec.md` → API Contract 4 states: the first sentence ends with the clause it gives word for word, and the entry's list of decision records gains a link to ADR-0451 after the link to ADR-0391, written as its neighbours are. No other entry changes.
7. MUST replace, in `docs/user-guide/commands/review.md` and `.agents/skills/roundfix/references/review.md`, the paragraph that begins "When round 1 leaves standing findings" and the first sentence of the paragraph after it with the text of `_techspec.md` → API Contract 5, wrapped as each file wraps its neighbours; the rest of each file keeps its bytes, and no `### QA settlement` section of any skill changes (`TestSettlementGuidanceIsOneTable`).
8. MUST then run, in this order, the commands `_authorization.md` → Sanctioned regeneration declares, and hand-edit none of their outputs:
   1. `make skills-sync`, which copies the canonical reference to `skills/roundfix/references/review.md`;
   2. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, the `-record-skill-versions` step, which raises both version fields of both `SKILL.md` copies of the Roundfix skill to the next version and records it in `skills/testdata/owned-skill-versions.json`;
   3. `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, the `-record-skill-coverage` step, which re-records the fingerprints of the review guide page and of the `command review` surfaces in `docs/references/behavior-surfaces.json`.
   A version is "the next version", never a number written by hand. Spec 0418 changes the same skill, records and glossary first; these commands regenerate on the tree this Task runs on. `docs/references/skill-coverage.json` does not change; when `TestTheSkillCoverageMapIsCurrent` asks for a change to it, stop and report it.
9. MUST add `internal/docscontract/review_model_decision_test.go`, with the `docscontract` build tag and the `verify:always` marker the package's other files carry, holding the three tests of `_techspec.md` → Testing Approach:
   - `TestTheReviewModelDecisionAmendsTheLineageRule`, which finds the decision record by its title, never by its number (invariant 7), as `TestTheIndexusDecisionAmendsTheConsultationRule` does, and asserts with `spec.ReadADRHeader` that it is accepted with a true Knowledge Mark; that ADR-0197's `amended_by` ends with its identifier, its `updated_at` equals the record's `created_at`, the new notice follows the notice of ADR-0391 with the record's `created_at` date, and the lines of ADR-0197 before the new notice keep the text of its decision and of its three earlier notices; and that `.secondbrain-export` holds the record's path once between the Standard Export List markers, with its ADR lines in ascending order;
   - `TestTheReviewerLineageEntryStatesTheModelRule`, which asserts the clause of API Contract 4 in the entry Reviewer Lineage and the entry's link to the record the first test found;
   - `TestTheReviewGuideAndReferenceDescribeTheModelRule`, with the subtests `guide`, `canonical_reference` and `mirror_reference`, each asserting, after whitespace is collapsed, the four phrases of the Testing Approach and the absence of "Round 2 prepares that name on that selection through acpx. If preparation".
10. MUST change in `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` only the version fields the record command raises (`_authorization.md` → What the grant covers).
11. MUST NOT change any other reference of the Roundfix skill, any other owned or Vendored Skill, any guide under `docs/agents/`, `docs/agents/setup-context.json`, `docs/references/skill-coverage.json`, any accepted ADR other than ADR-0197, any Go source other than the two named in requirement 1 and the two new test files, or any file `_authorization.md` → Limits names.

## Subtasks

- [ ] Record the model in the lineage and require it at the round-2 branch.
- [ ] Add the fake runner and the three review tests.
- [ ] Write ADR-0451, its notice in ADR-0197 and its export line.
- [ ] Revise the glossary entry Reviewer Lineage.
- [ ] Rewrite the two paragraphs of the guide page and of the canonical review reference, then run the skill sync, the skill record command and the Behavior Surface record command, in that order.
- [ ] Add the three documentation contract tests.

## Acceptance Criteria

- [ ] After the review profile's model changes from `gpt-6-luna` to `gpt-6.1-sol`, round 2 exits `0` with the outcome `reviewed`, ends the round-1 session before it prepares a new name, sends the round-2 delta prompt, records `continued: false` and `lineage.model` of `gpt-6.1-sol`, and prints nothing on standard error; the classifier still refuses the reported answer.
- [ ] A round-1 record without `lineage.model` makes round 2 open a new session the same way; with the profile unchanged, round 2 continues the round-1 session on the preferred or the fallback selection exactly as before.
- [ ] Every round that sent a prompt records the model of its session, in the printed record and the record file; a `ceiling-closed` lineage keeps its bytes.
- [ ] ADR-0451 is accepted with a true Knowledge Mark; ADR-0197 lists it in `amended_by` and carries its dated notice after the notice of ADR-0391; the Standard Export List holds its path once, in order; the entry Reviewer Lineage states the rule and links ADR-0451.
- [ ] The guide page, the canonical reference and its mirror state when round 2 continues the session and when it opens a new one; both copies of the Roundfix skill are byte-identical and carry the next recorded version; the Behavior Surface Record is current; no `### QA settlement` section changed.

## Context

- instruction: `internal/cli/review_session_test.go`
- instruction: `internal/cli/review_lineage_test.go`
- instruction: `internal/cli/review_lineage_selection_test.go`
- instruction: `internal/cli/review_test.go`
- instruction: `internal/docscontract/indexus_decision_record_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `internal/spec/adr_header.go`
- instruction: `skills/settlement_guidance_repocontract_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_lineage.go`
- interface: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- interface: `.secondbrain-export`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/cli/review_model_change_test.go`
- creates: `internal/docscontract/review_model_decision_test.go`
- creates: `docs/adr/0451-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed.md`

## Consumers

- the Pre-PR Review Record and its lineage — consumers: `./internal/cli`, `./internal/delivery`
- the owned Roundfix skill, its mirror and the owned-skill version record — consumers: `./skills`
- the review guide page, the glossary, the decision records, the Standard Export List and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged|TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel|TestReviewRecordsTheModelOfTheSessionOfEachRound|TestReviewRoundTwoContinuesTheRecordedSession|TestReviewRoundTwoContinuesTheRecordedFallbackSelection|TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails|TestReviewLineageChangeEndsTheOpenSession|TestACeilingRecordNamesNoSelection)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged/the_classifier_still_refuses_the_reported_answer TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel TestReviewRecordsTheModelOfTheSessionOfEachRound/preferred_then_continued TestReviewRecordsTheModelOfTheSessionOfEachRound/fallback TestReviewRoundTwoContinuesTheRecordedSession TestReviewRoundTwoContinuesTheRecordedFallbackSelection TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails TestReviewLineageChangeEndsTheOpenSession TestACeilingRecordNamesNoSelection; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the three new review tests and their subtests pass, and the five existing session and lineage tests still pass.
- `grep -q 'json:"model,omitempty"' internal/cli/review_lineage.go && go test -count=1 -timeout 20m -run 'Review|Lineage' ./internal/cli && go test -count=1 -timeout 20m ./internal/delivery` — expected: exit 0, the lineage records the model, and the review tests of the command and the whole Delivery Engine package, which read the Pre-PR Review Record, pass.
- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestTheReviewModelDecisionAmendsTheLineageRule|TestTheReviewerLineageEntryStatesTheModelRule|TestTheReviewGuideAndReferenceDescribeTheModelRule|TestTheSkillCoverageMapIsCurrent|TestGlossaryLinksNameAcceptedADRs)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheReviewModelDecisionAmendsTheLineageRule TestTheReviewerLineageEntryStatesTheModelRule TestTheReviewGuideAndReferenceDescribeTheModelRule/guide TestTheReviewGuideAndReferenceDescribeTheModelRule/canonical_reference TestTheReviewGuideAndReferenceDescribeTheModelRule/mirror_reference TestTheSkillCoverageMapIsCurrent TestGlossaryLinksNameAcceptedADRs; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the decision record amends ADR-0197, the glossary entry and the three review documents state the rule, and the Behavior Surface Record is current.
- `f=docs/adr/0451-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed.md; test -f "$f" && grep -qx 'status:.accepted' "$f" && grep -qx 'knowledge:.true' "$f" && grep -q '^amended_by:.*ADR-0451' docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md && grep -qxF -- "$f" .secondbrain-export && grep -F -- '**Reviewer Lineage**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Reviewer Lineage\*\*(:| —( |$))' && { tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "still names the Agent Model that session was opened with, and otherwise opening a new one with the same prompt" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "still names the Agent Model that session was opened with" >&2; exit 1; }; }` — expected: exit 0, ADR-0451 is accepted and marked, ADR-0197 and the export list name it, and the glossary still defines Reviewer Lineage on a line of its own with the new rule.
- `for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'only when the selection at that index still names the recorded model' || { printf 'missing phrase in %s: %s\n' "$f" 'only when the selection at that index still names the recorded model' >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'or when the record holds no model because an earlier Roundfix wrote it' || { printf 'missing phrase in %s: %s\n' "$f" 'or when the record holds no model because an earlier Roundfix wrote it' >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Round 2 prepares that name on that selection through acpx. If preparation'; then printf 'replaced sentence still in %s\n' "$f" >&2; exit 1; fi; done; cmp -s .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0, the guide page and both copies of the reference state the rule and drop the replaced sentence, and both skill copies are byte-identical.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/review.md | grep -qF -- "still names the recorded model" || { printf 'missing phrase in %s: %s\n' skills/roundfix/references/review.md "still names the recorded model" >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestSettlementGuidanceIsOneTable)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestSettlementGuidanceIsOneTable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the changed skill carries its recorded version, its mirror is in sync and the QA settlement guidance is untouched.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- "still names the recorded model" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/review.md "still names the recorded model" >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0, the consumers of the skill, the guide page, the glossary and the decision records pass whole.

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Non-Goals / Out of Scope; Glossary; Decisions
- `_techspec.md` → System Architecture; Interfaces; Invariants; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Coverage Map; Testing Approach; Build Order 1; Decisions
- `_authorization.md` → What the grant covers; Sanctioned regeneration; Limits
- [`references/sources.md`](references/sources.md), sections 1, 5, 6 and 7 (reported by fiscus, 2026-11-23)
- ADR-0197, ADR-0174, ADR-0256, ADR-0233, ADR-0286
