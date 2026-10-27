---
task: task_01
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
status: pending
type: backend
complexity: high
---

# Task 01: QA acceptance matches each declared-blocked row to the criterion a declaration names

## Overview

`spec.QAReportEligibility` compares `rows_blocked_declared` with the number
of `## Unreachable Acceptance` declarations, and the declaration parser
keeps one criterion each, so one declaration covering two criteria cannot
cover two rows. This Task lets a declaration name several criteria, matches
each declared-blocked row of the Results table to a declaration naming its
criterion, keeps today's count as the fallback so no report accepted today
is refused, makes the
mechanical stage classify a declared row with the same function, and writes
the rule into the `qa-gate` and `write-prd` skills, the Roundfix `settle`
reference, the `qa-report` guide, the QA prompt and `CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md`](references/2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md)
   of 2026-10-27, its acceptance half, by implementing API Contracts 1 to 4
   and 8 of `_techspec.md` and ADR-0343.
2. MUST change `parseUnreachableDeclarations` and its builder in
   `internal/spec/spec.go` as API Contract 1 states, adding
   `Criteria []string` to `UnreachableDeclaration` and keeping `Criterion`
   as the first criterion; every other field and error is unchanged.
3. MUST change `internal/spec/qa.go` as API Contracts 2 and 3 state: export
   `QADeclaredRow`, `QADeclaredRowCriterion` and `CriterionNamed`, add
   `DeclaredRows` to `QAReport`, read the rows in the same Results pass as
   `qaReportResultsCounts`, and accept in `QAReportEligibility` when the
   listed rows all match or when today's count holds, exactly as API Contract
   3 states. When both fail, the refusal keeps today's shortfall words and
   appends the unmatched rows; every report accepted at `5e9475bf` stays
   accepted.
4. MUST make `detectMechanicalReportShape` in
   `internal/speccheck/mechanical.go` count a declared row exactly when
   `spec.QADeclaredRowCriterion` accepts its status (API Contract 4), with
   no change to any finding.
5. MUST add `internal/spec/qa_declared_rows_test.go` with
   `TestUnreachableDeclarationNamesEveryCriterion` (two `criterion:` lines,
   a continuation line, the first criterion kept as `Criterion`, an empty
   criterion still refused), `TestCriterionNamedMatchesTextOrLeadingLabel`
   (a table that includes `Success Metric 5` against `Success Metric 5, after
   release, …`, backticks and a trailing period ignored, and `Success Metric
   1` refused against `Success Metric 10, …`) and
   `TestQAReportEligibilityMatchesEachDeclaredRow` (read from report files:
   one declaration with two criteria and rows `Q10` and `Q11` accepted,
   which today's count refuses; rows `Q10` and `Q12`, where `Q12` names
   `Success Metric 12`, refused with `shortfall is 1` and `declared-blocked
   row Q12`, an input under which a count of criteria would pass;
   `rows_blocked_declared: 3` over two matched rows and one declaration
   refused) and `TestQAReportEligibilityKeepsTodaysCountForAParaphrasedRow`
   (a characterization that passes before and after this Task: one
   declaration `no Run delivered after this Spec fails the test`, one row
   `Q09` recorded `blocked (declared: no later Run fails the test)` and
   `rows_blocked_declared: 1`, accepted; and the same report without a
   Results table, accepted).
6. MUST add `internal/cli/qa_report_declared_rows_test.go` with
   `TestQAReportAcceptMatchesEachDeclaredRowToItsCriterion` through CLI
   dispatch: the Quaestor shape exits `0` with no stdout; adding row `Q12`
   in place of `Q11` exits `1` with stderr naming `declared-blocked row Q12`.
7. MUST update, and change nothing else in, these existing tests:
   `TestUnreachableReadsDeclaredAcceptance` in
   `internal/spec/spec_test.go` (compare fields, since the struct holds a
   slice) and `TestBuildQAPromptStatesTheDeclaredMatrix` in
   `internal/agent/spec_prompt_test.go` (the new sentence).
8. MUST change `qaGateContract` in `internal/agent/spec_prompt.go` to name
   "each criterion a PRD Unreachable Acceptance declaration names" as a
   non-waivable source.
9. MUST write API Contract 8's sentence, verbatim, into a new
   `### Declared rows name their criterion` section of
   `.agents/skills/qa-gate/SKILL.md` placed outside `### QA settlement`
   (with the row form and the no-paraphrase rule), the same skill's
   non-waivable list, the comment of the `## Unreachable Acceptance` section
   of `.agents/skills/write-prd/references/prd-template.md` (with a second
   `criterion:` line in its example), the `## QA Report acceptance` section
   of `.agents/skills/roundfix/references/settle.md`, and
   `docs/user-guide/commands/qa-report.md`; then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version pairs of `qa-gate`, `write-prd` and `roundfix`.
10. MUST revise **Unreachable Acceptance** in `CONTEXT.md` through
    `domain-modeling` with the sentence API Contract 8 gives, citing
    ADR-0343 and keeping its `_Avoid_` line.
11. MUST, as the last step after the last guide edit, re-record the Behavior
    Surface fingerprint with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
    Hand-edit no mirror, version field or record.
12. MUST NOT change any text inside a `### QA settlement` section, the
    Archive Record schema, `archiveUnprovenActions`, any existing test other
    than the two Requirement 7 names, settle, the Daemon, the Makefile or the
    CI workflows. Each new top-level test calls
    `t.Parallel()` first.

## Subtasks

- [ ] Read several criteria per declaration and the declared rows of the Results table; match them in the eligibility decision and the mechanical stage.
- [ ] Add the spec and CLI tests and update the two existing tests; keep every other eligibility, archive and replay test unchanged and passing.
- [ ] Write the rule into the skills, the reference, the guide, the prompt and the glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] One declaration with two criteria covers two declared rows; a report refused by both paths names the rows no declaration names.
- [ ] Every report today's count accepts is still accepted, a paraphrased row included; the four callers and the mechanical stage read one classifier.
- [ ] The skills, reference, guide, prompt and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- QA Report eligibility and the declared-row classifier — consumers: `./internal/cli`, `./internal/daemon`, `./internal/spec`, `./internal/speccheck`

## Context

- instruction: `docs/adr/0343-settle-follows-the-branch-plan-and-a-declared-row-matches-the-criterion-it-names.md`
- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/spec/qa_partial_policy_test.go`
- instruction: `internal/cli/qa_report_test.go`
- interface: `internal/spec/spec.go`
- interface: `internal/spec/qa.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/spec_prompt_test.go`
- interface: `internal/spec/spec_test.go`
- instruction: `internal/spec/archive_test.go`
- instruction: `internal/spec/qa_prohibited_subject_test.go`
- creates: `internal/spec/qa_declared_rows_test.go`
- creates: `internal/cli/qa_report_declared_rows_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestUnreachableDeclarationNamesEveryCriterion|TestCriterionNamedMatchesTextOrLeadingLabel|TestQAReportEligibilityMatchesEachDeclaredRow|TestQAReportEligibilityKeepsTodaysCountForAParaphrasedRow|TestSpec0058ReplayRefusesUnmatchedBlockedRow|TestTheArchivedAuthorizationRowsQualify|TestUnreachableReadsDeclaredAcceptance)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUnreachableDeclarationNamesEveryCriterion TestCriterionNamedMatchesTextOrLeadingLabel TestQAReportEligibilityMatchesEachDeclaredRow TestQAReportEligibilityKeepsTodaysCountForAParaphrasedRow TestSpec0058ReplayRefusesUnmatchedBlockedRow TestTheArchivedAuthorizationRowsQualify; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 15m -v -run '^TestQAReportAcceptMatchesEachDeclaredRowToItsCriterion$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestQAReportAcceptMatchesEachDeclaredRowToItsCriterion ' || { echo 'missing pass: TestQAReportAcceptMatchesEachDeclaredRowToItsCriterion' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/spec ./internal/speccheck ./internal/agent ./internal/daemon && go test -count=1 -timeout 20m -run 'QA|Archive|Settle|Partial|Declared' ./internal/cli` — expected: exit 0. Before this Task the four spec tests and the CLI test do not exist, so the first pass check fails. Afterwards one declaration covers two rows, an undeclared row is refused by name, a paraphrased row and the archived replays stay accepted or refused exactly as before, the prompt pin agrees with the new rule, the four packages pass whole, and the QA, archive and settle tests of `./internal/cli` pass.
- `for file in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md .agents/skills/write-prd/references/prd-template.md skills/write-prd/references/prd-template.md .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md docs/user-guide/commands/qa-report.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'may name several criteria, one criterion line each, and each criterion becomes its own QA row' || { printf '%s lacks the declared-row rule\n' "$file" >&2; exit 1; }; done; grep -q '^### Declared rows name their criterion' .agents/skills/qa-gate/SKILL.md || { echo 'qa-gate lacks its Declared rows section' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- 'each criterion a PRD Unreachable Acceptance declaration names' || { echo 'qa-gate non-waivable list not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'One declaration may name several criteria; each is its own QA row.' || { echo 'CONTEXT.md lacks the revised **Unreachable Acceptance**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Unreachable Acceptance**:' || { echo 'CONTEXT.md lacks **Unreachable Acceptance**' >&2; exit 1; }; diff -r .agents/skills/qa-gate skills/qa-gate && diff -r .agents/skills/write-prd skills/write-prd && diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no skill, reference or guide states the rule, so the first phrase check fails. Afterwards the skills and their mirrors agree, the owned-skill versions are recorded, the `### QA settlement` table is unchanged, and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 1, 2 and 6; Core Features 1, 2 and 6; Success Metrics 1-3 and 7; Glossary
- `_techspec.md` → Current behavior 1-4; API Contracts 1-4 and 8; Declared breaks 1-2; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0343, ADR-0240
