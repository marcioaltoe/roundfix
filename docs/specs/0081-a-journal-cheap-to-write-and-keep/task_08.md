---
task: task_08
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: docs
complexity: medium
---

# Task 08: Decide the retention shape on the measurement

## Overview

The decision this whole Spec was arranged to defer. With the lock discipline,
the retention query, and the append path repaired and re-measured against the
task_01 baseline, the question finally has evidence behind it: does retention
still need a shape beyond age — and specifically, does anything need to shed
payloads at all?

The honest answer may be no, and that is a deliverable rather than a
shortfall. ADR-0008 makes the payload raw producer JSON, write-once and
read-as-blob, and ADR-0030 removed the per-Batch agent logs precisely because
the journal is the durable copy. Shedding payloads destroys the only copy, so
it is bought with a measurement or not at all.

## Requirements

1. MUST re-run the task_01 harness on the repaired code and record the after
   beside the before, so the comparison is same-harness rather than anecdotal.
2. MUST state, from that comparison, whether the remaining cost justifies any
   retention shape beyond age — and MUST record the reasoning either way.
3. MUST, when the answer is no, say so explicitly in the recorded decision and
   leave ADR-0008 and the retention policy untouched. A Spec that ends smaller
   than its title is the correct outcome when the measurement says so.
4. MUST, when the answer is yes, land it as an explicit ADR-0008 amendment
   naming the capability that becomes unrecoverable and the window in which it
   is still recoverable, and MUST NOT smuggle the change past the ADR.
5. MUST, in that case only, record that the reconcile replay probe's payload
   equality key is a precondition to be re-keyed before any payload rewrite
   becomes possible.
6. MUST update `docs/user-guide/run-database-lifecycle.md` when anything an
   operator can observe changes, and MUST leave it alone when nothing does.
   This includes its machine-checked durable-table row: if the chosen shape
   adds or changes a table, that row lands in the same change, which the
   guide's own test enforces.
7. MUST NOT change production behaviour in this Task. Its deliverable is a
   measured decision and, at most, an ADR.

## Subtasks

- [ ] Re-run the harness and record the after against the before.
- [ ] Write the decision with its reasoning, in whichever direction.
- [ ] Land the ADR amendment only if the measurement demands it.

## Acceptance Criteria

- [ ] The after measurement exists beside the before, from the same harness.
- [ ] The decision is recorded with reasoning, and states plainly whether any
      payload is shed.
- [ ] If payloads are shed, an ADR-0008 amendment exists naming the lost
      capability; if not, ADR-0008 is untouched.
- [ ] No production code path changed in this Task.

## Context

- instruction: docs/user-guide/run-database-lifecycle.md

## Verification

- `output="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/*.md 2>/dev/null | wc -l | tr -d ' ')"; [ "$output" -ge 2 ]`
  — expected: exit 0; a second measurement artifact sits beside the baseline.
- `report="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/*.md | tail -1)"; grep -qiE 'payload' "$report" && grep -qi 'ADR-0008' "$report"`
  — expected: exit 0; the decision addresses payload shedding and names the
  ADR it either respects or amends.
  — expected: exit 0; no production code path changed in this Task.
  — expected: exit 0; the store is still green after the re-measurement.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Features 5 and 7; Goal 4; User Story 6; Decisions
  (ADR-0008 is binding until explicitly amended).
- `_techspec.md` → Build Order 8; Risks (the measurement may disprove the
  premise).
- ADR-0008, ADR-0030, ADR-0033.
