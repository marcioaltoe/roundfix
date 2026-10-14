---
task: task_01
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
status: completed
type: backend
complexity: medium
---

# Task 01: A partial blocked only by the Spec's own network prohibition qualifies

## Overview

The shared QA policy gains the Network-Prohibited Row of `_techspec.md` →
Data Models and API Contract 1: an outside-evidence row recorded exactly as
`blocked (environment: network prohibited by Spec authorization)` is an
exempt row like the Network-Denied Row. Every caller reads the policy from
`internal/spec`, so settlement, archive, `qa-report accept` and the Delivery
Queue's park decision change together. It is verifiable on its own through
the policy tests and the shape tests that run each command and the Daemon.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-13,
   [outside evidence blocked by the Spec's own authorization does not qualify](references/2026-10-13-outside-evidence-blocked-by-the-spec-authorization-does-not-qualify.md),
   by implementing API Contract 1 of `_techspec.md` in `internal/spec/qa.go`:
   the exported `QANetworkProhibitedStatus` constant, the
   `RowsBlockedNetworkProhibited` count derived from the Results rows in the
   same pass as the other exempt kinds, its subtraction in
   `EnvironmentRowsNeedingOverride`, and its place among the unmet exempt rows
   in `QAReportEligibility` (ADR-0276, ADR-0240).
2. MUST count a row only when its trimmed status equals the constant exactly
   and its provenance names the outside-evidence row as
   `qaReportProvenanceNames` reads it. A host suffix, a changed case, another
   source or a missing Provenance column counts nothing, and the row stays an
   ordinary environment row.
3. MUST keep every existing refusal text unchanged for a report with no
   Network-Prohibited Row, and MUST use the refusal text of API Contract 1 when
   a non-exempt environment row stands beside one. The frontmatter counts keep
   their meaning and validation.
4. MUST add the three `internal/spec` tests that `_techspec.md` → Testing
   Approach names, the five shapes it names to `qaPartialShapes` in both the
   command and the Daemon test files, and the two park cases to
   `TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly`. The
   `Spec 0258 network prohibited with PR and declaration` shape is four
   network-prohibited outside-evidence rows, one Pull Request row and one
   declared row covered by one declaration, with
   `rows_blocked_environment: 5`, and it qualifies.
5. Every new top-level test MUST call `t.Parallel()` first, and every new case
   MUST use the exported constants rather than a copied status literal, except
   the cases whose point is a wrong literal.
6. MUST NOT change any other caller of the policy, any skill, guide or user
   guide page; those are task_03 and task_04.

## Subtasks

- [ ] Add the constant, the count and its use in the policy.
- [ ] Add the three policy tests.
- [ ] Add the shapes and park cases in the command and Daemon tests.

## Acceptance Criteria

- [ ] A partial whose only unmet rows are Network-Prohibited Rows, Pull
      Request rows, Network-Denied Rows and covered declarations is accepted by
      `qa-report accept`, `archive`, `settle` and Daemon settlement, and does not
      park `qa-environment-partial`.
- [ ] A network-prohibited status without the outside-evidence provenance or
      with a host suffix refuses with `rows_blocked_environment is 1; expected 0`
      and parks.
- [ ] A non-exempt environment row beside a Network-Prohibited Row refuses
      with the new text naming both exempt outside-evidence kinds.
- [ ] Every existing policy test passes unchanged.

## Consumers

- The shared QA eligibility policy (`QAReportEligibility`, `EnvironmentRowsNeedingOverride`) — consumers: `./internal/spec`, `./internal/cli`, `./internal/daemon`

## Context

- instruction: `docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`
- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- interface: `internal/spec/qa.go`
- interface: `internal/spec/qa_partial_policy_test.go`
- interface: `internal/cli/qa_partial_policy_test.go`
- interface: `internal/daemon/qa_partial_policy_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestANetworkProhibitedOutsideEvidenceRowNeverDecidesAQualifyingPartial|TestANetworkProhibitedRowNeedsItsExactStatusAndOutsideEvidence|TestAnotherEnvironmentRowBesideANetworkProhibitedRowNamesBoth)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestANetworkProhibitedOutsideEvidenceRowNeverDecidesAQualifyingPartial TestANetworkProhibitedRowNeedsItsExactStatusAndOutsideEvidence TestAnotherEnvironmentRowBesideANetworkProhibitedRowNamesBoth; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape|TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly|TestSettleNamesTheReportItRefused)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape/network_prohibited_only TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape/Spec_0258_network_prohibited_with_PR_and_declaration TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape/network_prohibited_without_outside_evidence TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape/network_prohibited_with_a_host TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape/another_environment_beside_network_prohibited TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly/network_prohibited_outside_evidence TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly/network_prohibited_without_outside_evidence; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in network_prohibited_only Spec_0258_network_prohibited_with_PR_and_declaration network_prohibited_without_outside_evidence network_prohibited_with_a_host another_environment_beside_network_prohibited; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo/$name " || { printf 'missing Daemon pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/spec` — expected: exit 0; before this Task none of the three policy tests exists, so the command fails at its first pass check.

## References

- [_prd.md](_prd.md) — Goals 1; User Stories 1-2; Core Features 1; Success Metric 1; Success Metric 2
- [_techspec.md](_techspec.md) — Data Models; API Contract 1; Invariants 1-3; Testing Approach; Build Order 1
- ADR-0276; ADR-0240; ADR-0104; ADR-0167; ADR-0259; ADR-0271

## Result

Implemented API Contract 1 in the shared QA policy. `QANetworkProhibitedStatus`
identifies only the exact trimmed status; `RowsBlockedNetworkProhibited` is
computed in the existing Results-table pass and requires an outside-evidence
provenance item. The policy subtracts this count with the other exempt kinds,
clamps the subtraction to the frontmatter environment count, and includes the
new kind among the unmet rows that qualify a partial. A remaining environment
row beside a recognized prohibited row uses the exact new refusal. Existing
refusal branches and frontmatter validation are unchanged.

Added the three named policy tests, the five named shapes in both command and
Daemon fixtures, and the two named Delivery Queue park cases. The Spec 0258
shape contains four prohibited outside-evidence rows, one Pull Request row,
one declared row, one declaration, and five environment-blocked rows. New
policy tests call `t.Parallel()` first and use the exported row constants.
Additional policy cases exercise trimmed whitespace, annotated provenance,
exemption clamping, changed case, source boundaries, and missing provenance.

### Focused checks

- Before wiring the policy, `rtk proxy go test ./internal/spec -count=1 -run
  '^TestANetworkProhibitedOutsideEvidenceRowNeverDecidesAQualifyingPartial$/^prohibited_only$'`
  exited 1: the canonical row produced `RowsBlockedNetworkProhibited: 0`
  instead of 1. This established the regression before the behavior change.
- After the implementation, `rtk proxy go test ./internal/spec -count=1 -run
  'Test.*(Partial|PullRequest|NetworkProhibited|NetworkDenied|EnvironmentRowsNeedingOverride|QAReport|QAVerdict)'`
  exited 0 (`ok roundfix/internal/spec`). This includes every existing test
  in `qa_partial_policy_test.go`, unchanged, and the new policy tests.
- `rtk proxy go test ./internal/cli -count=1 -run
  'TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape|TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly'`
  exited 0 (`ok roundfix/internal/cli`). All command shapes and park cases
  passed through isolated local fixtures.
- `rtk proxy go test ./internal/daemon -count=1 -run
  'TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo/(network_prohibited|Spec_0258|another_environment_beside_network_prohibited)'`
  exited 0 (`ok roundfix/internal/daemon`). The five new shapes passed at
  the TaskCycle settlement boundary.
- Ran `rtk proxy gofmt -w internal/spec/qa.go
  internal/spec/qa_partial_policy_test.go internal/cli/qa_partial_policy_test.go
  internal/daemon/qa_partial_policy_test.go`; `rtk proxy git diff --check`
  exited 0.

### Acceptance evidence

1. The qualifying policy test accepts prohibited rows alone, beside PR or
   network-denied rows, with all exempt kinds and a covered declaration, and
   in the Spec 0258 shape. Both new positive command shapes are accepted by
   `qa-report accept`, `archive`, and `settle`; Daemon settlement accepts both.
   The `network prohibited outside evidence` park case remains unparked.
2. The exact-status/provenance policy test derives zero exemptions and one
   override row for missing outside-evidence provenance and host suffixes,
   refusing with `rows_blocked_environment is 1; expected 0`. Both shapes
   refuse across the commands and Daemon. The missing-provenance park case
   parks; host-suffixed rows retain the same one-row override count that the
   park decision reads.
3. The mixed-environment policy test and matching command/Daemon shape
   assert exactly: `rows_blocked_environment is 2, 1 outside the pre-PR Pull
   Request row and the network-denied and network-prohibited outside-evidence
   rows; expected 0 outside them`.
4. All existing tests in the shared partial-policy file ran unchanged in the
   focused spec check; all existing command shapes and park cases ran
   unchanged in the focused command check. Broader spec-package validation
   remains with the Daemon's declared Verification.

No follow-up work was added to this diff. Task status and checkboxes were
left unchanged. The authored Verification command, repository Verification,
and incremental Verification were not run; no commit, push, or Pull Request
was made. Implementation and focused evidence are handed back for Daemon
Verification and settlement.
