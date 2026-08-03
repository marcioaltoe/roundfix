---
task: task_05
spec: 0072-qa-is-a-task-not-a-flag
status: pending
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
