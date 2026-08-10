---
task: task_01
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: test
complexity: medium
---

# Task 01: Measure what the journal actually costs

## Overview

The first artifact, and the one every later claim in this Spec cites. The
three-gigabyte measurement that opened the Spec is a symptom; the costs worth
fixing — a write lock held across a full-table aggregate at Run start, one
transaction and one fsync per agent output line — were only visible once the
write path was read rather than the disk.

Nothing here changes behaviour. It builds the harness and records the before.

## Requirements

1. MUST measure event-write latency, lock-wait time, and `SQLITE_BUSY`
   frequency as functions of journal size and concurrent writer count.
2. MUST measure Run-start cost against journal size, since the retention
   eligibility scan is the path this Spec claims is size-dependent.
3. MUST record the results as a committed artifact in the Spec folder, with
   the harness parameters, the machine facts that affect the numbers, and the
   commit measured.
4. MUST make the harness rerunnable, so an after can be produced from the same
   code and compared honestly against the before.
5. MUST seed its own journal fixtures rather than measuring the operator's live
   database, so the measurement is hermetic and repeatable.
6. MUST NOT change any production code path, configuration default, or schema.

## Subtasks

- [ ] Build the seeded-journal harness with its parameters.
- [ ] Record write latency, lock wait, `SQLITE_BUSY`, and Run-start cost.
- [ ] Commit the baseline artifact with the commit and machine facts.

## Acceptance Criteria

- [ ] The baseline artifact exists in the Spec folder and names the commit it
      measured.
- [ ] Rerunning the harness on the same commit reproduces comparable numbers.
- [ ] No production code path or default changed.

## Verification

- `output="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/*.md 2>/dev/null)"; [ -n "$output" ]`
  — expected: exit 0; the committed baseline artifact exists.
- `report="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/*.md | tail -1)"; grep -qiE 'latency' "$report" && grep -qi 'SQLITE_BUSY' "$report" && grep -qiE 'run start|run-start' "$report"`
  — expected: exit 0; the three measured dimensions are recorded.
- `output="$(go test -count=1 ./internal/store -run 'Baseline|Harness' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the harness is selected and runs.
  — expected: exit 0; nothing outside the Spec folder and store test files
  changed, so no production path moved.

## References

- `_prd.md` → Core Feature 1; Success Metrics; Decisions (measurement precedes
  design).
- `_techspec.md` → Testing Approach; Build Order 1.
- `references/2026-08-06-event-journal-payload-economics.md` → the intent this
  measurement serves.
