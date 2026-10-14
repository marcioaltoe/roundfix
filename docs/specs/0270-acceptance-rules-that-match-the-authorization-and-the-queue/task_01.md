---
task: task_01
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
status: pending
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
