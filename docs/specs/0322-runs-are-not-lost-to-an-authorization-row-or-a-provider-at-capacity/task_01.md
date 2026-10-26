---
task: task_01
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: pending
type: backend
complexity: medium
---

# Task 01: A Network-Prohibited Row names what the authorization forbids

## Overview

The shared QA partial policy exempts an outside-evidence row the Spec's
authorization keeps offline only when its status equals one network
sentence, so Specs 0317 and 0309, whose rows named what their authorizations
actually forbade, needed the QA Archive Override. This Task reads the status
as the form `blocked (environment: <subject> prohibited by Spec
authorization)` with a bounded subject, counts each Results row toward one
exempt kind at most, and proves that settlement, `qa-report accept`,
`archive` and the Daemon agree on the archived shapes.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md`](references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md)
   of 2026-10-26 (`_prd.md` → Sources 1) by implementing API Contracts 1 and
   2 of `_techspec.md`.
2. MUST add the constants and `qaProhibitedRowSubject` of `_techspec.md` →
   Interfaces to `internal/spec/qa.go`, counting the subject with
   `utf8.RuneCountInString`, and make `qaReportResultsCounts` count a
   Network-Prohibited Row through it, with the Network-Denied Row check moved
   into an `else if` so a row counts toward one exempt kind at most.
3. MUST NOT change `QANetworkProhibitedStatus`, `QANetworkDeniedStatusPrefix`,
   `QAOutsideEvidenceRowSource`, `qaReportProvenanceNames`, the frontmatter
   counts, skipped-row counting, `EnvironmentRowsNeedingOverride`,
   `QAReportEligibility` or any refusal message.
4. MUST create `internal/spec/qa_prohibited_subject_test.go` with
   `TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids` and
   `TestTheArchivedAuthorizationRowsQualify` and their subtests exactly as
   `_techspec.md` → Testing Approach names them, using the archived Results
   rows quoted in `_prd.md` → Validation and the package's existing helpers.
   Each top-level test and subtest calls `t.Parallel()` first.
5. MUST add the three shapes of `_techspec.md` → Testing Approach to
   `qaPartialShapes` in both `internal/cli/qa_partial_policy_test.go` and
   `internal/daemon/qa_partial_policy_test.go`, with identical rows and
   expectations, and change no existing shape.
6. MUST NOT change any document, skill, Baseline module, the Daemon's
   production code, the Makefile or the CI workflows.

## Subtasks

- [ ] Read the subject form and count each row once.
- [ ] Add the subject and archived-report tests.
- [ ] Add the three shapes to the CLI and Daemon tables.

## Acceptance Criteria

- [ ] Specs 0309's and 0317's archived tables qualify and Spec 0258's still does.
- [ ] The subject bounds, the parenthesis, the suffix case and the provenance rule refuse as specified, and a row in both forms counts once.
- [ ] `settle`, `qa-report accept`, `archive` and Daemon settlement agree on the three new shapes, and every existing shape keeps its outcome.

## Consumers

- `internal/spec` QA partial eligibility (`QAReportEligibility`, `ReadQAReportFile`) — consumers: `internal/cli`, `internal/daemon`

## Context

- instruction: `internal/spec/qa_partial_policy_test.go`
- instruction: `internal/spec/qa_prepr_row_test.go`
- interface: `internal/spec/qa.go`
- creates: `internal/spec/qa_prohibited_subject_test.go`
- interface: `internal/cli/qa_partial_policy_test.go`
- interface: `internal/daemon/qa_partial_policy_test.go`

## Verification

- `out="$(go test -timeout 15m -count=1 -v -run '^(TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids|TestTheArchivedAuthorizationRowsQualify|TestANetworkProhibitedOutsideEvidenceRowNeverDecidesAQualifyingPartial|TestANetworkProhibitedRowNeedsItsExactStatusAndOutsideEvidence|TestAnotherEnvironmentRowBesideANetworkProhibitedRowNamesBoth|TestANetworkDeniedOutsideEvidenceRowNeverDecidesAQualifyingPartial|TestANetworkDeniedRowWithoutOutsideEvidenceProvenanceStillRefuses|TestAPartialWithASkippedRowNeverQualifies|TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in accepts_Spec_0309_R07_another_repository accepts_Spec_0317_QA-08_upstream_checkout_and_network accepts_the_canonical_network_form accepts_a_subject_of_80_multibyte_characters refuses_an_empty_subject refuses_a_blank_subject refuses_a_subject_of_81_characters refuses_a_subject_with_a_parenthesis refuses_a_changed_suffix_case refuses_a_row_without_outside_evidence a_row_in_both_forms_counts_once a_skipped_row_still_refuses; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids/$name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for name in Spec_0258 Spec_0309 Spec_0317; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheArchivedAuthorizationRowsQualify/$name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -timeout 15m -count=1 -v -run '^TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in Spec_0309_another_repository_prohibited_with_PR_and_declaration/qa-report Spec_0309_another_repository_prohibited_with_PR_and_declaration/archive Spec_0309_another_repository_prohibited_with_PR_and_declaration/settle Spec_0317_upstream_checkout_and_network_prohibited_with_PR_and_declaration/qa-report Spec_0317_upstream_checkout_and_network_prohibited_with_PR_and_declaration/archive Spec_0317_upstream_checkout_and_network_prohibited_with_PR_and_declaration/settle authorization_subject_with_a_parenthesis/qa-report authorization_subject_with_a_parenthesis/archive authorization_subject_with_a_parenthesis/settle; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape/$name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -timeout 15m -count=1 -v -run '^TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for shape in Spec_0309_another_repository_prohibited_with_PR_and_declaration Spec_0317_upstream_checkout_and_network_prohibited_with_PR_and_declaration authorization_subject_with_a_parenthesis; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo/$shape" || { printf 'missing pass: %s\n' "$shape" >&2; exit 1; }; done` — expected: exit 0. Before this Task the new tests and shapes do not exist, so the first pass check fails. Afterwards the subject form qualifies within its bounds, each row counts once, the archived tables of Specs 0258, 0309 and 0317 qualify, the existing partial-policy tests pass unchanged, and the commands and the Daemon agree on the new shapes.

## References

- `_prd.md` → Sources 1; Goals 1-2; Core Feature 1; Success Metrics 1-3; Validation
- `_techspec.md` → Current behavior 1-3; Interfaces; API Contracts 1-2; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0335; ADR-0276
