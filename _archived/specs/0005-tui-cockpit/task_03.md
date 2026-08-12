---
task: task_03
spec: 0005-tui-cockpit
status: completed
type: frontend
complexity: medium
---

# Task 03: Upgrade the Work Queue rows

## Overview

Bring the queue pane to the mocked density: Batch separators with elapsed
time, per-row status marker, severity when the Work Item carries one,
ordinal, title, and location, selected-row emphasis, and the compact totals
footer. Verifiable through row-rendering table tests for both Run kinds.

## Requirements

1. MUST group queue rows under Batch separators showing the batch ordinal and
   elapsed time for the executing Batch, per the mockup.
2. MUST render each row as: status marker (`[run]`, `[done]`, `[wait]`,
   `[locked]`, plus the failure/invalid states the pane already knows),
   severity label only when the underlying Work Item carries one (verify what
   Review Issue artifacts actually provide; render nothing when absent —
   never invent a level), ordinal, title, and file location when known.
3. MUST keep the selected row visually distinct and stable while statuses
   refresh (selection semantics unchanged).
4. MUST render the queue footer totals line (`N issues total · X resolved ·
   Y unresolved` — vocabulary adjusted per Run kind for Tasks).
5. MUST keep spec-Run rows (`task_NN`) flowing through the same renderer.

## Subtasks

- [x] Batch separators with elapsed placement
- [x] Row renderer: marker, optional severity, ordinal, title, location
- [x] Totals footer per Run kind
- [x] Table tests across all Work Item states and both kinds

## Acceptance Criteria

- [x] Row tests cover every status label the pane supports today plus the
      severity-present and severity-absent cases.
- [x] Separator tests pin batch ordinal and elapsed rendering.
- [x] Totals footer asserts for a review and a spec fixture.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 3. `_techspec.md` → Build Order 3,
Risks (severity display). `design/ui-redesign-plan.md` → Required Changes
(queue rows); `design/roundfix-01.png`.

## Result

- Row evidence: `TestCockpitWorkQueueRowsRenderMarkersMetadataAndOptionalSeverity`
  covers `Executing`, `Resolved`, `Completed`, `Waiting`, `Paused`, `Invalid`,
  `Duplicated`, and `Failed`, including severity-present Review Issue rows and
  severity-absent Task rows.
- Separator evidence: `TestCockpitWorkQueueBatchSeparatorShowsOrdinalAndElapsed`
  pins `BATCH 001/002`, elapsed placement for the executing Batch, and no
  elapsed time on a waiting Batch.
- Totals evidence: `TestCockpitWorkQueueFooterTotalsForRunKinds` asserts the
  review footer and the spec Task footer.
- Existing behavior evidence: updated cockpit tests keep selection movement,
  spec Task refresh, and mid-write last-good status behavior passing with the
  new row markers.
- Snapshot evidence: cockpit goldens under
  `internal/tui/testdata/cockpit_snapshots/` were deliberately regenerated for
  this row-density visual change.
- Verification passed: `rtk go test ./internal/tui/` (70 tests),
  `rtk go test ./...` (550 tests), and `rtk make verify`.
