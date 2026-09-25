---
task: task_02
spec: 0164-knowledge-lifecycle-and-capture
status: completed
type: backend
complexity: medium
---

# Task 02: Review retirement reads recorded evidence

## Overview

`ClassifyReview` in `internal/spec/review_liveness.go` answers from local Git, so the same orphan Review Artifact classifies three ways as refs are fetched and deleted, and a squash merge can never retire one by ancestry. Legacy `docs/specs/_reviews/` artifacts that are not finished are never renamed.

## Requirements

1. MUST make `ClassifyReview` read `outcome.md` at the Review Artifact root, with front matter `pull_request_state` (`merged`, `closed` or `open`), `merge_commit` and `recorded_at`, and run no Git command.
2. MUST answer `ReviewFinished` for `merged` with a hexadecimal `merge_commit` (a squash receipt) whether or not that commit or the recorded head exists locally, and for `closed`.
3. MUST answer `ReviewLive` for `open`, and `ReviewUndecidable` with a reason naming what is missing for an absent or malformed record or a `merged` record without a merge commit.
4. MUST make `internal/baseline/history_layout.go` relocate a live or undecidable artifact under `docs/specs/_reviews/` to `docs/specs/reviews/` and keep its retained report, while a finished one relocates to `docs/history/reviews/` as today, through the existing collision checks.
5. MUST replace the expectations of `TestClassifyReviewLocalGit` that ADR-0163 revokes, as a declared break, and keep every Review Artifact root resolver answer outside the history root.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, one artifact classifies the same with its recorded head present, absent, and reachable only from fetched refs.
- [ ] A recorded squash receipt retires an artifact whose merge commit is absent locally.
- [ ] An artifact with no recorded outcome is undecidable and retained.
- [ ] A legacy `_reviews` artifact relocates to `docs/specs/reviews/` when live or undecidable, and to history when finished.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/review_liveness.go`
- interface: `internal/baseline/history_layout.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestClassifyReviewIgnoresObjectStoreAvailability|TestClassifyReviewAcceptsARecordedSquashReceipt|TestClassifyReviewWithoutARecordedOutcomeIsUnknown|TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness)$" ./internal/spec ./internal/baseline 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestClassifyReviewIgnoresObjectStoreAvailability TestClassifyReviewAcceptsARecordedSquashReceipt TestClassifyReviewWithoutARecordedOutcomeIsUnknown TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Review retirement
- ADR-0163

## Result

Implemented recorded-evidence Review retirement. `ClassifyReview` now reads
`outcome.md`, validates the declared Pull Request state and required fields,
accepts a hexadecimal squash receipt without consulting Git, and retains
missing or malformed outcomes with a specific reason. History layout
discovery now routes every legacy `_reviews` artifact through the existing
relocation and collision ledger: live and undecidable artifacts move to
`docs/specs/reviews/`, while finished artifacts move to
`docs/history/reviews/`.

Focused checks and evidence:

- Before implementation, `GOCACHE=/tmp/roundfix-task02-gocache rtk go test
  -count=1 -run '^TestClassifyReviewIgnoresObjectStoreAvailability$'
  ./internal/spec` failed when the recorded head was absent, and the matching
  focused history-layout test failed because no legacy live or undecidable
  relocation was produced.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test -count=1 -run
  '^TestClassifyReview' ./internal/spec` passed 17 tests. This covers stable
  classification across present, absent and fetched-only Git objects; an
  absent-local squash receipt; open and closed outcomes; and missing or
  malformed records.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test -count=1 -run
  '^TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness$'
  ./internal/baseline` passed. The test checks live, undecidable and finished
  legacy artifacts, retained findings, report-file preservation and the three
  destination trees.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test -count=1 ./internal/spec
  ./internal/baseline` passed, including history-move transaction coverage.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test -count=1 -run
  '^TestReviewArtifactRootNeverResolvesIntoHistory$' ./internal/config`
  passed, preserving the resolver boundary required by ADR-0163.
- `rtk make verify-incremental` first reached two process-owner integration
  failures because the sandbox denied process-table access. The same command
  passed with host process-table permission; all packages, skill checks and the
  build passed.

Acceptance evidence:

- `TestClassifyReviewIgnoresObjectStoreAvailability` proves one recorded-open
  artifact returns the same answer and reason with its head present, absent or
  reachable only from a fetched ref.
- `TestClassifyReviewAcceptsARecordedSquashReceipt` proves a recorded
  hexadecimal merge commit retires the artifact without a matching local Git
  object.
- `TestClassifyReviewWithoutARecordedOutcomeIsUnknown` proves an absent outcome
  is undecidable and names the missing `outcome.md`; malformed-field cases are
  covered by `TestClassifyReviewRejectsMalformedRecordedOutcome`.
- `TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness` proves live
  and undecidable legacy artifacts retain their reports under the canonical
  live root, while a finished artifact relocates to history.

The Task's declared `## Verification` command was not run; the Daemon owns that
verification and status settlement.

## Carry-forward provenance

- Source Run: `run_20260924T223414Z_54e6f85ae98b6550`
- Source commit: `c941f8a9c39f435847f76a8bcb7f0110a741d336`
