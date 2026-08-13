---
task: task_12
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 12: Admit a real fleet archive into the migration

## Overview

The QA gate's outside-evidence row proved the migration does not work on a
repository this Spec did not build: applying to a copy of an adopted repository
exits 2 before any mutation, because its 10.2 MiB of archived documentation
crosses the carrier inventory budget of 8 MiB. The budget bounds how much
instruction text the Baseline ingests to plan; an archive it only moves is
counted against it. This slice raises the budget so the measured fleet fits, with
the measurement recorded beside the constant.

## Requirements

1. MUST allow a repository whose archived documentation reaches the measured
   fleet maximum to plan and apply its migration without refusal.
2. MUST size the budget against a recorded measurement rather than against the
   one repository that failed, and MUST record that measurement and its date
   beside the constant so a later reader can tell what it was sized for.
3. MUST keep refusing an inventory that is implausible rather than merely large,
   so the budget still bounds a hostile or broken tree.
4. MUST keep the per-file and carrier-count bounds unchanged; only the cumulative
   byte budget moves.
5. MUST prove the new budget against a repository this Spec did not build, using
   a disposable copy and never a working checkout.

## Subtasks

- [ ] Raise the cumulative inventory byte budget with its measurement recorded.
- [ ] Cover the boundary with a test that fails at the old value.
- [ ] Prove a fleet-sized archive plans and applies.

## Acceptance Criteria

- [ ] An inventory at the measured fleet maximum plans and applies without
      refusal.
- [ ] The constant carries the measurement, its date, and the repositories it was
      sized against.
- [ ] An inventory beyond the new budget still refuses, naming the limit.
- [ ] The per-file and carrier-count bounds are unchanged.
- [ ] A test exercises the boundary and fails against the previous value.

## Verification

- `grep -q 'maxInventoryBytes' internal/baseline/repository.go && ! grep -q 'maxInventoryBytes     = 8 \* 1024 \* 1024' internal/baseline/repository.go` — expected: exits 0, proving the budget exists and no longer holds its previous value. Both clauses are one command because either alone passes on an unchanged tree.
- `grep -B4 -A1 'maxInventoryBytes' internal/baseline/repository.go > /tmp/0094-task-12-c.txt; grep -q '2026-08-13' /tmp/0094-task-12-c.txt && grep -qi 'MiB\|measured' /tmp/0094-task-12-c.txt` — expected: exits 0, proving the measurement and its date sit beside the constant rather than in a commit message nobody reads at the point of change.
- `go test -count=1 ./internal/baseline -run 'TestInventoryBudget' -v > /tmp/0094-task-12.log 2>&1; s=$?; grep -q '^--- PASS: TestInventoryBudget' /tmp/0094-task-12.log || { cat /tmp/0094-task-12.log; exit 1; }; exit $s` — expected: exits 0 and the log names the boundary test; fails today, where no such test exists.
- `grep -q 'maxInventoryFileBytes = 2 \* 1024 \* 1024' internal/baseline/repository.go && grep -q 'maxInventoryCarriers  = 256' internal/baseline/repository.go && ! grep -q 'maxInventoryBytes     = 8 \* 1024 \* 1024' internal/baseline/repository.go` — expected: exits 0, proving the two untouched bounds kept their values while the cumulative budget moved.

## Context

- interface: `internal/baseline/repository.go`

## References

`_prd.md` → Goal 2; User Story 2; Success Metrics, the fleet-migration row.
`_techspec.md` → Risks: moving files the tool does not own. QA report
`qa/qa-report-2026-08-13.md` → F-001.

Measured 2026-08-13 across four adopted repositories this Spec did not build:
vortex 413 files at 12.4 MiB, conexus 173 at 10.2 MiB, fluxus 458 at 8.1 MiB,
oraculum 735 at 3.4 MiB. File count does not predict bytes — oraculum holds the
most files and the fewest bytes — so the budget is sized on bytes, and an archive
only grows, so the value needs headroom rather than a fit.
