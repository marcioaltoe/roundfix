---
task: task_01
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: completed
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

## Result

Implemented the bounded authorization-subject form in the shared Results-row
reader. `qaProhibitedRowSubject` trims the subject, requires 1–80 Unicode
code points using `utf8.RuneCountInString`, rejects either parenthesis, and
keeps the prefix and suffix case-sensitive. Prohibited classification now
precedes network-denied classification through `else if`, so one row cannot
supply both exemptions. The existing canonical constants, provenance reader,
frontmatter handling, skipped counting, eligibility policy and refusal
messages retain their bytes.

Acceptance evidence:

- Archived shapes: `TestTheArchivedAuthorizationRowsQualify` embeds the
  Results tables read at `078869bf` (0258), `06a9b7ce` (0309) and `ef1a4bab`
  (0317), with environment counts 5, 4 and 3 and one covered declaration
  each. It checks prohibited counts 4, 3 and 2, one PR exemption, no denied
  exemptions and no remaining override rows. Tests read temporary fixtures,
  never repository history or active Spec paths at runtime.
- Subject and refusal boundaries: all twelve prescribed subtests of
  `TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids` exercise the
  named subjects, canonical form, trimmed 80 multibyte code points, empty,
  blank, 81-character, parenthesis, suffix-case and provenance refusals,
  exclusive counting and the skipped-row refusal. Both top-level tests and
  every subtest call `t.Parallel()` first.
- Consumer agreement: the three new shapes passed focused checks through
  CLI `qa-report accept`, `archive` and `settle`, and Daemon settlement.
  Additional focused checks covered Spec 0258, missing provenance, malformed
  canonical statuses, another environment row, skipped rows and uncovered
  declarations. A Python comparison against HEAD confirmed that the three
  additions have identical rows and expectations in both tables and that
  removing only those additions reproduces each prior table byte-for-byte.
  Full declared Verification remains for the Daemon.

Focused checks (all final commands exited 0):

- `go test ./internal/spec -run 'Test(ANetworkProhibited|TheArchivedAuthorization|AnotherEnvironmentRowBesideANetworkProhibited|EnvironmentRowsNeedingOverride)' -count=1`
- `go test ./internal/spec -run 'Test(ANetworkDenied|APartialWithASkipped|APartialWhoseOnly|APullRequestRowWith)' -count=1`
- `go test ./internal/cli -run '^TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape$/(Spec_0309|Spec_0317|authorization_subject)' -count=1`
- `go test ./internal/daemon -run '^TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo$/(Spec_0309|Spec_0317|authorization_subject)' -count=1`
- `go test ./internal/cli -run '^TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape$/(Spec_0258|network_prohibited|another_environment|PR_and_skipped|uncovered_declaration)' -count=1`
- `go test ./internal/daemon -run '^TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo$/(Spec_0258|network_prohibited|another_environment|PR_and_skipped|uncovered_declaration)' -count=1`
- `git diff --check`

Before the parser edit, the focused Spec 0309 subtest exited 1 with
`prohibited/denied counts = 0/0, want 1/0`, establishing the regression.
An intermediate CLI build exposed an insertion into the wrong table; the
three rows were moved into `qaPartialShapes` and both final CLI checks above
passed. No existing shape was changed.

The initial worktree change was the Daemon's `pending` → `in_progress`
status transition in this Task file; it is preserved. Only this Result and
the four code/test paths in this Task were edited. No declared Verification,
repository-wide Verification, commit, push or PR operation was run. No
follow-up outside this slice was identified.
