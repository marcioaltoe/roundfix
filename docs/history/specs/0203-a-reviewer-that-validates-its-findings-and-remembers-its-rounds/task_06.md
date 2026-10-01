---
task: task_06
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: completed
type: backend
complexity: low
---

# Task 06: A ceiling record names no reviewer selection

## Overview

The QA gate of Run `run_20261001T145917Z_95f0b38476c6ebcc` found QA-0203-02. The ceiling outputs of Surface Transcripts 5 and 6 include `"selection":0` in their lineage object, although no reviewer ran. `reviewLineage.Selection` in `internal/cli/review_lineage.go` is an `int` with the unconditional tag `json:"selection"`, and `closeAtCeiling` leaves it at zero. In rounds 1 and 2, `selection` is the index of the Agent Selection that reviewed, and `0` means the Preferred Selection. So at the ceiling, a zero reads as though the preferred reviewer ran.

## Requirements

1. MUST change `reviewLineage.Selection` to `*int` with the tag `json:"selection,omitempty"`, set it where `internal/cli/review.go` records a round's selection, and leave it nil on the ceiling path. Rounds 1 and 2 MUST still emit `selection`, including the value `0`.
2. MUST add `internal/cli/review_lineage_selection_test.go` with:
   - `TestARoundRecordNamesItsSelectionEvenWhenItIsThePreferred`: a round reviewed by selection 0 emits `"selection":0`;
   - `TestACeilingRecordNamesNoSelection`: the ceiling record's lineage has no `selection` key, in both the blocked and the closed outcome.
3. MUST reproduce Surface Transcripts 5 and 6 as authored, with no `selection` in their lineage object.
4. MUST NOT change any other field, record or behavior, or any existing test except one whose expected bytes include the ceiling's `"selection":0`.

## Subtasks

- [ ] Make the selection optional, and set it only for rounds.
- [ ] Add the two tests.

## Acceptance Criteria

- [ ] Both tests pass, and the existing `TestReview*` suite stays green.

## Context

- interface: `internal/cli/review_lineage.go`
- interface: `internal/cli/review.go`
- creates: `internal/cli/review_lineage_selection_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestARoundRecordNamesItsSelectionEvenWhenItIsThePreferred|TestACeilingRecordNamesNoSelection)$' ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestARoundRecordNamesItsSelectionEvenWhenItIsThePreferred TestACeilingRecordNamesNoSelection; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 -run '^TestReview' ./internal/cli` — expected: exit 0; before this Task neither new test exists.

## References

- task_03
- `_techspec.md` → Surface Transcripts 5 and 6
- QA Report of Run `run_20261001T145917Z_95f0b38476c6ebcc`, finding QA-0203-02

## Result

Added the two requested regression tests in
`internal/cli/review_lineage_selection_test.go`. The round test checks public
JSON from both rounds and requires an explicit numeric selection of zero.
The ceiling test replays Surface Transcripts 5 and 6 through the public
review command with a temporary repository and fake reviewer, checks their
exit codes and stderr, and requires the exact authored lineage fields with
no selection. It also checks that only the first two rounds call the
reviewer. Its verbose output records both surface transcripts.

Focused evidence:

- `rtk proxy go test ./internal/cli -run 'Test(A.*RecordNames)' -count=1 -v`
  could not access the default Go build cache under the sandbox.
- `GOCACHE=/private/tmp/roundfix-task06-go-cache rtk proxy go test ./internal/cli -run 'Test(A.*RecordNames)' -count=1 -v`
  exited 1: the preferred-selection test passed for rounds 1 and 2; both
  ceiling subtests reproduced the unwanted `"selection":0`. The blocked
  transcript exited 2 with its authored diagnostic, and the closed
  transcript exited 0 with empty stderr.
- `rtk proxy gofmt -w internal/cli/review_lineage_selection_test.go`
  exited 0.

### Verification Feedback repair — attempt 1

Inspected the Daemon's diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261001T165954Z_dec5e0e74ad42b1b/verification/batch-001-attempt-1.log`.
It confirmed the same ceiling serialization defect observed in the focused
check above.

Following the instruction to repair this slice, changed `Selection` to
`*int` with `json:"selection,omitempty"` and stored its address only where a
reviewed round records its selection. The ceiling constructor remains
unchanged and leaves it nil. The continued-session reader dereferences a
recorded selection and preserves the previous zero default for older
records without the field. No other record field or outcome changed.

The required type change also requires mechanical updates to three existing
integer comparisons in `review_session_test.go`. Those assertions now check
presence and dereference the value, retaining their preferred-selection and
fallback-selection expectations. The affected failure diagnostic prints the
lineage rather than formatting a pointer as an integer. This is the narrow
compile adaptation identified in the initial scope conflict; no test
expectation was weakened or otherwise changed.

Focused repair evidence:

- `GOCACHE=/private/tmp/roundfix-task06-go-cache rtk proxy go test ./internal/cli -run 'Test(A.*RecordNames|ReviewRound.*|ReviewCeiling.*)' -count=1 -v`
  exited 0. Both requested tests passed, including both ceiling outcomes.
  Surface Transcripts 5 and 6 have exactly the authored lineage fields and
  no selection, with exit 2 and its diagnostic for blocked and exit 0 with
  empty stderr for closed. Existing round, ceiling, preferred-session and
  continued-fallback tests in this focused set also passed.
- `rtk proxy gofmt -w internal/cli/review_lineage.go internal/cli/review.go internal/cli/review_session_test.go internal/cli/review_lineage_selection_test.go`
  exited 0.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- `rtk make verify-incremental GOCACHE=/private/tmp/roundfix-task06-go-cache`
  initially exited 2 because sandbox restrictions blocked process-table
  integration tests and an Agent Result edit during the run triggered the
  suite's repository-change guard. Repeated with the required process access
  and no concurrent edits; it exited 0. Formatting, vet, the full Go suite
  (including existing `TestReview*` tests), skill checks, and build passed.

Acceptance evidence: both requested tests passed in the focused check;
existing review behavior also passed in the full incremental suite. The
Daemon still owns the declared Verification rerun and Task settlement.

Declared Verification remains the Daemon's responsibility and was not
rerun. Task status, checkboxes, the Task Graph, and other Task files were not
edited; no commit, push, or pull request was created.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/review_session_test.go`

## Carry-forward provenance

- Source Run: `run_20261001T165954Z_dec5e0e74ad42b1b`
- Source commit: `925f354a24c4392835d66465fdd31ea9d63c2095`
