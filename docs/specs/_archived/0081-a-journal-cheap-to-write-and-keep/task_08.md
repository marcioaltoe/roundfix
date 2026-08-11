---
task: task_08
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
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

## Result

### Implementation

- Re-ran the byte-identical task_01 harness twice against repaired production
  commit `27e3e21a5f0b7987c45dc5ccc57fe44a5c73c2ec` and recorded both after runs
  beside the before in `baseline/2026-08-11-repaired.md`.
- Recorded the decision to keep only the existing terminal-only, age-based
  Journal Retention shape. No payload is shed, rewritten, compressed, or moved
  to a second retention window.
- Left ADR-0008, ADR-0033, the reconcile replay probe, and
  `docs/user-guide/run-database-lifecycle.md` untouched. The measurement does
  not justify losing the journal's only durable raw agent payload copy, and no
  operator-observable behavior or durable table changed.

### Focused checks

- Pre-change signal: `rtk rg --files
  docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline` listed only
  `2026-08-11-before.md`; no after measurement or measured decision existed.
- `rtk git diff 04246137..HEAD --
  internal/store/journal_baseline_test.go` produced no output, proving the
  repaired-code run used the task_01 harness unchanged.
- `GOCACHE=/private/tmp/roundfix-task08-gocache rtk proxy go test -count=2
  ./internal/store -run '^TestJournalMeasurementHarness$' -v` passed both
  fresh runs. All 336 writes succeeded with zero typed `SQLITE_BUSY` errors.
  At 10,000 events, Run-start p50 fell from 11.449–11.551 ms before to
  0.041–0.047 ms after; after p50 stayed within 0.036–0.051 ms across every
  journal size.
- The same harness also exposed a cost rather than hiding it: four-writer
  direct-append p95 increased from 9.222–9.529 ms before to 33.104–33.428 ms
  after serialization. That after band is flat across 0, 1,000, and 10,000
  events, so reducing retained payloads would not address it.
- `rtk git diff --check` exited 0. `rtk git diff HEAD -- internal cmd` and
  `rtk git diff HEAD -- docs/adr
  docs/user-guide/run-database-lifecycle.md` both produced no output. The final
  status contains only this Task file and the new Spec-local report; the
  Task-file status transition was the pre-existing Daemon change.

### Acceptance evidence

1. **The after exists beside the before from the same harness.** The new
   `baseline/2026-08-11-repaired.md` records both complete after matrices,
   parameters, machine facts, measured commit, and before/after comparison;
   the exact harness-source diff from its task_01 commit is empty.
2. **The decision states whether payload is shed and explains why.** The report
   says no payload is shed. Run-start no longer scales with retained journal
   size, and remaining direct concurrent-write latency is also independent of
   retained size, so payload loss buys neither measured hot-path repair.
3. **ADR-0008 is untouched when payloads are retained.** The decision respects
   ADR-0008 without amending it. ADR-0030 confirms the journal is the only
   durable agent-payload copy, and ADR-0033's existing age boundary remains the
   only retention shape.
4. **No production code path changed in this Task.** This Task adds the
   Spec-local after report and this Result only; the frontmatter transition to
   `in_progress` was already present and remains Daemon-owned.

### Daemon handoff

- The commands under this Task's `## Verification` were not run. The Daemon
  owns those commands, terminal status, and the Task commit.
- No follow-up work was discovered within this Task's decision-only slice.
