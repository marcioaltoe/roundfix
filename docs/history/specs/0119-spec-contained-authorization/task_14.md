---
task: task_14
spec: 0119-spec-contained-authorization
status: completed
type: backend
complexity: high
---

# Task 14: Let the audit read the reference the checker resolved

## Overview

The changed-path audit is this Spec's central control, and it fails open for
the citation form this Spec makes canonical. The constraint reader accepts a
Spec-relative Markdown reference; the audit re-derives the reference with a
backtick-only extractor, gets nothing, finds no Task commits, and records a
skip. A Spec that follows the new templates therefore has no changed-path audit
at all. This slice makes both consumers read one resolved reference.

## Requirements

1. MUST resolve a Tooling authority row's authorization reference once and give
   the audit that resolved value, so the audit never re-derives it with a
   narrower rule than the reader used.
2. MUST audit a Spec that cites its record in the Spec-relative Markdown form,
   producing the same Task-commit comparison it produces for a backticked
   reference today.
3. MUST resolve a row citing several records to the one the row names as its
   express authorization, and MUST refuse rather than resolve to none when the
   claim cannot be matched to exactly one record. An unresolved reference must
   not read as a skip: a skip says there was nothing to audit, and that is the
   failure this Task removes.
4. MUST keep every audit outcome that works today working, including the
   out-of-grant refusal, the grant-edited-in-the-consuming-commit refusal, and
   the presence-aware skip for a Spec that genuinely declares no authorization.
5. MUST prove the control cannot silently do nothing: a Spec with a governed
   change outside its grant must refuse under both citation forms.

## Subtasks

- [ ] Carry the reader's resolved reference into the audit request.
- [ ] Audit the Spec-relative Markdown form end to end.
- [ ] Resolve a multi-record row and refuse an unmatchable claim.
- [ ] Prove the preserved outcomes, including the genuine no-authorization skip.

## Acceptance Criteria

- [ ] A Spec citing its record as a Spec-relative Markdown reference produces a
      changed-path audit with Task commits; today it produces a skip.
- [ ] A governed change outside the grant refuses under both the backticked and
      the Markdown citation forms.
- [ ] A row citing an approved narrow grant beside a proposed record audits
      against the approved one; a row whose claim matches no single record
      refuses instead of skipping.
- [ ] A Spec that declares no authorization still records the presence-aware
      skip it records today.
- [ ] A test fails if the audit reports zero Task commits for a Spec whose
      reference resolved, so a fail-open cannot pass as a skip again.

## Context

- interface: `internal/speccheck/mechanical.go`
- interface: `internal/speccheck/constraints.go`

## Verification

- `grep -q 'func TestAuditReadsTheResolvedReference' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestAuditReadsTheResolvedReference$'` — the Markdown citation form produces a real audit; today it produces a skip.
- `grep -q 'func TestAuditRefusesOutOfGrantUnderBothCitationForms' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestAuditRefusesOutOfGrantUnderBothCitationForms$'` — an out-of-grant governed change refuses whichever form the row uses.
- `grep -q 'func TestUnresolvedReferenceIsNotASkip' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestUnresolvedReferenceIsNotASkip$'` — an unmatchable claim refuses, and a genuine no-authorization Spec still skips.
- `grep -q 'func TestAuditReadsTheResolvedReference' internal/speccheck/mechanical_test.go || exit 1; go test -count=1 ./internal/speccheck` — the package suite passes with the change present.

## References

- `_prd.md` → User Stories 3; Core Features 2; Goals 2, 4.
- `_techspec.md` → Implementation Design: One resolved reference, read by every consumer.
- `qa/qa-report-2026-09-09.md` → the audit skip behind F-003.

## Result

The constraint artifact reader now resolves the Tooling authority row's
role-labelled authorization references once and stores the selection on the
row. Authoring validation and `MechanicalAuthorization` consume that same
selection, so the Daemon receives the Spec-relative Markdown path that the
reader accepted instead of a second backtick-only projection. An operative
claim that does not identify exactly one record returns an error before the
mechanical stage can turn the missing path into a skip; a row that genuinely
declares no authorization still returns the empty presence-aware input.

Focused red evidence before the production change:

- `rtk go test -count=1 ./internal/speccheck -run 'Test(AuditReadsTheResolvedReference|AuditRefusesOutOfGrantUnderBothCitationForms|UnresolvedReferenceIsNotASkip)'`
  reached the expected regressions: the Markdown path and the approved record
  beside a proposal both resolved to empty, the unmatched operative claim
  returned no error, and the Markdown out-of-grant case produced no finding.

Focused post-change checks:

- The same focused command passed 8 tests after the shared selection change.
- `rtk go test ./internal/speccheck` passed 375 tests, including the existing
  out-of-grant, grant-edited-in-the-consuming-commit, ancestor-grant,
  regeneration, and presence-aware skip coverage.
- `rtk gofmt -d internal/speccheck/constraints.go internal/speccheck/mechanical.go internal/speccheck/mechanical_test.go`
  produced no diff.
- `rtk make verify-incremental` stopped at `fmt-check` because the unchanged
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. Neither path is
  changed in this Task worktree, so this Task did not alter them.

Acceptance evidence:

1. `TestAuditReadsTheResolvedReference` resolves
   `[_authorization.md](_authorization.md)`, completes one authorization read
   for `task_01`, and rejects a zero-read result.
2. `TestAuditRefusesOutOfGrantUnderBothCitationForms` proves the backticked and
   Spec-relative Markdown forms each report `.golangci.yml` outside the exact
   `Makefile` grant.
3. `TestUnresolvedReferenceIsNotASkip/express_grant_is_selected_beside_a_proposal`
   proves the audit selects the express narrow grant instead of the proposed
   widening. Its `operative_claim_without_one_matching_record_refuses` case
   proves an unmatchable claim returns an exact-reference error.
4. `TestUnresolvedReferenceIsNotASkip/genuine_no-authorization_declaration_keeps_presence-aware_skip`
   proves an honest no-authorization row still records the existing
   `authorization bounded paths` skip for missing `tooling authorization`.
5. The explicit `AuthorizationReads` count in
   `TestAuditReadsTheResolvedReference` makes a resolved reference with zero
   audited Task commits fail instead of accepting a silent skip.

The Task's declared `## Verification` commands were not run; Daemon
Verification remains the settlement owner.
