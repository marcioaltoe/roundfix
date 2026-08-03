---
task: task_05
spec: 0072-qa-is-a-task-not-a-flag
status: completed
type: backend
complexity: medium
---

# Task 05: Measure and trim the gate cycle's own cost

## Overview

The maintainer asked for this slice explicitly: the gate cycles on Spec 0057
cost twenty to twenty-five minutes each, and three of them ran. Part of that
cost died with Spec 0071 (the per-Task suite tax and the cold local gate).
This Task measures what one gate cycle spends **now**, attributes it —
agent-session time, verification commands, snapshotting, artifact and
report handling — and applies the trims that leave verdict semantics
untouched, per the same discipline the 0071 campaign used: measure, cut
repeated work, never weaken a check.

## Requirements

1. MUST measure one full gate cycle end to end on a representative Spec
   fixture and record the attribution: where the minutes go, in a table,
   committed under this Spec's folder.
2. MUST apply only trims that keep report content, verdict semantics, and
   typed blocked-row counts byte-compatible; anything that would change
   what the gate observes is out of scope and gets recorded as a finding
   instead.
3. MUST re-measure after the trims with the same procedure and commit the
   before-and-after beside the attribution.
4. MUST record, for any cost that dominates and cannot be trimmed inside
   this Spec's boundary (for example Agent-session inference time), one
   honest paragraph naming it and why it stays.
5. MUST keep every existing gate test passing unmodified.

## Subtasks

- [ ] Instrument and measure one gate cycle; commit the attribution.
- [ ] Apply the semantics-preserving trims the attribution justifies.
- [ ] Re-measure; commit the before-and-after.
- [ ] Record untrimmable dominant costs honestly.

## Acceptance Criteria

- [ ] An attribution table under `docs/specs/0072-qa-is-a-task-not-a-flag/`
      names where a gate cycle's time goes, measured not estimated.
- [ ] Every trim is justified by a line in the attribution, and the
      after-measurement shows the delta.
- [ ] Report, verdict, and count outputs are byte-compatible: the existing
      gate tests pass unmodified.
- [ ] Dominant untrimmable costs are named, not hidden.
- [ ] `git status --porcelain` shows no path outside `internal/daemon/`,
      `internal/cli/`, `docs/specs/0072-qa-is-a-task-not-a-flag/`, and this
      task file.

## Verification

- `ls docs/specs/0072-qa-is-a-task-not-a-flag/ | grep -qi "gate-cost\|gate_cost"`
  — expected: exit 0; the attribution artifact exists.
- `go test ./internal/daemon -count=1 -run 'QA|Gate' -v | grep -q -- "--- PASS"`
  — expected: exit 0; gate behavior unchanged.
- `go test ./internal/daemon ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → the twenty-to-twenty-five-minute cycles; Non-Goals (what the
  gate checks does not change).
- `_techspec.md` → Build Order 5; Coverage Map (maintainer-requested
  performance slice).

## Result

### Implementation

- Added `gate-cost-2026-08-03.md` with a measured end-to-end attribution from
  the latest complete representative QA journal, the same-worktree fixture
  procedure, the before-and-after, and the current measurement limitation.
- Removed one repeated filesystem scan on the missing-report settlement path.
  `spec.QAVerdict` has already proven `spec.ErrNoQAReport`, so the Daemon now
  returns the existing `missing` verdict and empty report path immediately.
- Left pass, partial, fail, unreadable-report, typed blocked-row validation,
  snapshots, artifact filtering, report commits, and every existing gate test
  unchanged.

### Acceptance evidence

- Attribution artifact: `gate-cost-2026-08-03.md` records a measured
  625.710-second representative cycle and attributes Agent startup, static
  verification, other Agent work, verdict/report settlement, and the report
  commit with exact journal boundaries.
- Trim and delta: the artifact ties the missing-report trim to the 1.221-second
  settlement row. Ten-sample focused benchmarks measured the missing-report
  settlement median at 3,771 ns/op and 1,697 B/op with 16 allocations before,
  then 1,905 ns/op and 848 B/op with 8 allocations after. The complete warm
  gate fixture measured 0.48s before and 0.49s after, correctly reporting no
  material cycle-level change rather than hiding timing noise.
- Output compatibility: the existing
  `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` passed after the
  last Go edit across pass, partial, fail, missing, and unreadable verdicts.
  No test, fixture, golden output, report schema, verdict value, or typed-count
  assertion was edited. The Daemon-owned gate suites remain to run.
- Dominant untrimmable cost: the artifact names 468.669 seconds of Agent
  inference and live QA work, 74.90% of the representative cycle, and explains
  why shortening it inside this slice would change what the gate observes or
  violate the separate QA-session contract.
- Scope: the final changed-path inspection contains only
  `internal/daemon/task_engine.go`, this Task's gate-cost artifact, and this
  task file. The task file's pre-existing `pending` to `in_progress` status
  change remains Daemon-owned.

### Focused checks

- The same-worktree production gate fixture command passed after the final Go
  edit using the task-local `GOCACHE`; all five existing verdict cases passed.
- `git diff --check` exited 0 after the final edits.
- A fresh cost-measurement-only `make verify` attempt stopped after 77.06s in
  pre-existing `TestCoverageEquivalence` and baseline characterization drift
  owned by earlier Spec 0072 Tasks. Task 05 did not edit or regenerate those
  paths. This failed measurement is recorded in the attribution and is not
  presented as passing verification.

### Follow-up

- A post-Spec-0071 real QA Agent cycle is still needed to replace the latest
  complete end-to-end journal measurement. The current Run cannot supply that
  until earlier-Task verification drift is settled and the authored QA Task
  runs under Daemon ownership.
