---
task: task_06
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: pending
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
