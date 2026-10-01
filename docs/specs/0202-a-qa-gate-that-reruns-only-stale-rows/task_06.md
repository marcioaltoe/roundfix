---
task: task_06
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: pending
type: test
complexity: medium
---

# Task 06: Two real gate passes carry a row into the committed report

## Overview

The QA gate of Run `run_20261001T133454Z_7983d7d312c8ec5c` blocked row 5 with finding F1. `TestQAGateCarriesARowFromAnUnintegratedFailedPass` builds the first pass's Evidence Snapshot by hand, and its fake Agent replaces the second pass's seeded report with a hollow passing report. So no test proves that a first real gate records a row and that the next real gate keeps it in its committed report. QA found no production defect: this Task adds the missing proof.

## Requirements

1. MUST add `internal/daemon/qa_two_pass_carry_test.go` with `TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport`. Through the existing task-cycle fixture it MUST:
   - drive two `TaskCycle` QA passes over a disposable repository;
   - in the first pass, use a fake Agent that writes a report with one passing carriable row declaring its inputs and one failing row, and let the Daemon record the Evidence Snapshot itself, never by hand;
   - leave the first row's inputs unmoved and change only a file the failing row depends on;
   - in the second pass, use a fake Agent that honors the carried rows the Daemon seeds, keeping them unchanged, and re-runs only the rest.
2. MUST assert, for both passes, the `prior_report`, `mechanical` and `evidence_snapshots` `daemon.qa` events with their payloads. The second pass's `prior_report` must name the first pass, and its `evidence_snapshots` must cover the re-run row.
3. MUST read the second pass's committed QA Report independently from Git (`git show <commit>:<path>`). It MUST assert that the carried row keeps the first pass's provenance cell and is listed as carried in `## Row carry-forward`, and that the changed row is re-run with its disposition.
4. MUST add `TestTwoGatePassesReRunARowWhoseInputMoved`, the mirror case: when the carried row's input moves between the passes, the second pass re-runs it as `input moved:` and does not keep it.
5. MUST NOT change production code or any existing test.

## Subtasks

- [ ] Build the two-pass fixture with honoring fake Agents.
- [ ] Assert the events and read the committed report from Git.
- [ ] Add the moved-input mirror case.

## Acceptance Criteria

- [ ] Both new tests pass. Sabotaging the carry, so that the seed drops the carried row, makes the first test fail.

## Context

- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/daemon/qa_two_pass_carry_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport|TestTwoGatePassesReRunARowWhoseInputMoved)$' ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport TestTwoGatePassesReRunARowWhoseInputMoved; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither test exists.

## References

- task_01, task_02, task_03
- QA Report of Run `run_20261001T133454Z_7983d7d312c8ec5c`, finding F1
