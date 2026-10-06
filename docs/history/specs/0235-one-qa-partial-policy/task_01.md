---
task: task_01
spec: 0235-one-qa-partial-policy
status: completed
type: backend
complexity: medium
---

# Task 01: A partial whose only unmet rows nobody in the Run can reach qualifies under the one eligibility policy

## Overview

The Backlog Entry of 2026-10-06
([settlement refuses a partial that archive accepts](references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md))
records QA reports refused although their only unmet rows were the pre-PR
Pull Request row, sometimes with declared rows. Every QA settlement caller
reads `spec.QAReportEligibility`, so this Task changes that one function and
the QA Report reader beneath it. It adds the network-denied outside-evidence
row of ADR-0240, recognizes a Pull Request row whose provenance carries a
note, refuses a partial with a skipped row, and adds
`QAReport.EnvironmentRowsNeedingOverride`.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-06 named in the Overview by
   implementing the TechSpec's Interfaces and Invariants 1 to 7 in
   `internal/spec/qa.go`: the exported constants
   `QANetworkDeniedStatusPrefix` and `QAOutsideEvidenceRowSource`, the
   derived `QAReport.RowsBlockedNetworkDenied` and `QAReport.RowsSkipped`, the
   method `QAReport.EnvironmentRowsNeedingOverride`, and the eligibility
   order of Invariant 6. `QAReportEligibility` keeps its signature.
2. MUST derive the two new counts from the Results rows in the same pass and
   under the same rules as `RowsBlockedPrePullRequest`: only tables under
   `## Results` with a Status column (and a Provenance column, for the two
   exempt kinds) count, and fenced blocks never count. The frontmatter
   counts keep their meaning and validation.
3. MUST keep every existing refusal message byte-identical, and use the
   messages of API Contract 2 and API Contract 3 only for the two new
   refusals. Declarations are read only when `rows_blocked_declared` is
   greater than zero, so a qualifying partial with no declared row does not
   need `_prd.md`.
4. MUST add `internal/spec/qa_partial_policy_test.go` with the seven tests
   named in Verification. Each test builds reports through
   `ReadQAReportFile` in a temporary Spec directory, using these measured
   shapes: Pull Request rows only (archived Specs 0213 and 0220 and the
   Oraculum report), a provenance of
   `Pull Request row (Requirement 10 constrains execution)` (archived Spec
   0227), two covered declared rows with the Pull Request row (the Fluxus
   report), a `blocked (environment: network denied: docs.example.com)` row
   with provenance `Requirement 7; outside-evidence row (published guide)`,
   the same status with provenance `Requirement 8`, another environment row
   beside a network-denied outside-evidence row, a skipped row beside the
   Pull Request row, and a status that only looks like the marker, such as
   `blocked (environment: network denied: )`.
5. MUST change the one existing expectation the new policy reverses, and no
   other. That is the subtest "partial without a declared row keeps its
   refusal" of `TestAPassIsUnchangedByThePrePullRequestRow` in
   `internal/spec/qa_prepr_row_test.go`. Its report now qualifies, and the
   subtest is renamed to say so. Every other case in that file and in
   `internal/spec/qa_test.go` and `internal/spec/archive_test.go` keeps its
   expectation.
6. MUST NOT change Daemon settlement, `internal/spec/archive.go`, the
   command callers, the mechanical stage, any skill or guide, `CONTEXT.md`
   or `CHANGELOG.md`.

## Subtasks

- [ ] Derive network-denied outside-evidence rows and skipped rows from the Results rows.
- [ ] Accept a note after the Pull Request row's provenance item.
- [ ] Apply the eligibility order of Invariant 6 and add the override count.
- [ ] Add the policy tests and flip the one reversed expectation.

## Acceptance Criteria

- [ ] A partial whose unmet rows are only pre-PR Pull Request rows,
      network-denied outside-evidence rows and covered declared rows
      qualifies, whatever the mix of those kinds.
- [ ] A partial with any other environment row, a network-denied row
      without outside-evidence provenance, a skipped row, a finding row or
      an uncovered declared row is refused, and an existing refusal keeps
      its message.
- [ ] `EnvironmentRowsNeedingOverride` counts exactly the environment rows
      outside the two exempt kinds, and is never negative.

## Context

- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- instruction: `docs/adr/0167-the-pre-pr-pull-request-row-never-decides-a-qualifying-partial.md`
- interface: `internal/spec/qa.go`
- interface: `internal/spec/qa_prepr_row_test.go`
- creates: `internal/spec/qa_partial_policy_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestAPartialWhoseOnlyUnmetRowsArePullRequestRowsQualifies|TestAPullRequestRowWithANoteInItsProvenanceIsRecognized|TestANetworkDeniedOutsideEvidenceRowNeverDecidesAQualifyingPartial|TestANetworkDeniedRowWithoutOutsideEvidenceProvenanceStillRefuses|TestAPartialWithAnotherEnvironmentRowBesideANetworkDeniedRowNamesBoth|TestAPartialWithASkippedRowNeverQualifies|TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows|TestAPassIsUnchangedByThePrePullRequestRow|TestANoPullRequestStatusWithoutThePullRequestSourceStillRefuses|TestAPartialWithAnotherEnvironmentRowStillRefuses|TestAPartialQualifiesWhenItsOnlyEnvironmentRowIsThePrePullRequestRow)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAPartialWhoseOnlyUnmetRowsArePullRequestRowsQualifies TestAPullRequestRowWithANoteInItsProvenanceIsRecognized TestANetworkDeniedOutsideEvidenceRowNeverDecidesAQualifyingPartial TestANetworkDeniedRowWithoutOutsideEvidenceProvenanceStillRefuses TestAPartialWithAnotherEnvironmentRowBesideANetworkDeniedRowNamesBoth TestAPartialWithASkippedRowNeverQualifies TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows TestAPassIsUnchangedByThePrePullRequestRow; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the seven new tests exists, so their pass lines are missing and the command fails; after it they pass beside the existing pre-PR row tests.

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Acceptance evidence
- `_techspec.md` → Interfaces; Invariants 1 to 7; API Contract 2; API Contract 3; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 1
- ADR-0240; ADR-0167; ADR-0080


## Result

Implemented the Task 01 policy in `internal/spec/qa.go`. The Results reader
now derives pre-PR Pull Request, network-denied outside-evidence and skipped
row counts in one pass. Provenance items recognize the specified note
boundaries, and the network marker requires a non-empty host. Frontmatter
counts and their validation retain their existing meaning.

`QAReportEligibility` applies Invariant 6's refusal order and reads
Unreachable Acceptance declarations only when declared rows exist. Existing
refusal text is preserved; the two new refusals use API Contracts 2 and 3.
`EnvironmentRowsNeedingOverride` subtracts the two exempt kinds, clamped to
the environment count. The pass rule is unchanged.

Focused checks:

- Before implementation,
  `rtk proxy env GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/spec -run 'TestAPartialWhoseOnlyUnmetRowsArePullRequestRowsQualifies|TestAPullRequestRowWithANoteInItsProvenanceIsRecognized'`
  exited 1: the PR-only partials were refused with `expected "pass"`, and
  annotated PR provenance derived zero exempt rows.
- After the final code edit,
  `rtk proxy env GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/spec`
  exited 0 (`ok roundfix/internal/spec`, 22.060s). Output was captured in
  `/tmp/roundfix-task01-spec-check.log`. An earlier package-check attempt
  was invalidated by a formatting edit during execution; its suiteguard
  failure is not used as evidence.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

Acceptance evidence:

1. Eligible mixes: `TestAPartialWhoseOnlyUnmetRowsArePullRequestRowsQualifies`
   covers the 0213, 0220 and Oraculum PR-only shapes without `_prd.md`.
   `TestAPullRequestRowWithANoteInItsProvenanceIsRecognized` covers 0227's
   annotated provenance. `TestANetworkDeniedOutsideEvidenceRowNeverDecidesAQualifyingPartial`
   covers network-only, PR/network, Fluxus's PR/two-declaration shape,
   all three kinds, network/declared and declared-only reports.
2. Refusals: `TestANetworkDeniedRowWithoutOutsideEvidenceProvenanceStillRefuses`
   covers Requirement 8 provenance, empty/blank hosts, missing closing
   parenthesis and invalid source/status matches.
   `TestAPartialWithAnotherEnvironmentRowBesideANetworkDeniedRowNamesBoth`
   asserts API Contract 2 exactly. `TestAPartialWithASkippedRowNeverQualifies`
   asserts API Contract 3, finding/environment/skipped/declaration precedence,
   uncovered declarations and no-unmet-row refusal. Existing QA, pre-PR and
   archive tests passed unchanged except for the one prescribed PR-only
   partial expectation.
3. Override count: `TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows`
   covers both exempt kinds, other environment rows, zero and undersized
   frontmatter counts, missing columns, fenced and outside-Results tables,
   multiple tables and deeper headings. Counts stay non-negative, and skipped
   rows need no Provenance column.

Only this Task's implementation/test paths and this Result section were
edited. The incoming Daemon-owned `status: in_progress` is preserved.
Authored Verification was not run; Verification and settlement remain with
the Daemon. No commit, push or Pull Request was made. No follow-up was found
inside this slice; caller and guidance work remains with the later Tasks.

## Carry-forward provenance

- Source Run: `run_20261006T120012Z_3cae1c6ab9ca3edb`
- Source commit: `fa25b7267dfbc31833197effc3ea9448ee17b4bb`
