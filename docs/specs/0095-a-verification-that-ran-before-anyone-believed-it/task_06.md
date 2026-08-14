---
task: task_06
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 06: Restore the vacuity refusal and account what it finds

## Overview

The vacuity detector exists, works, and is disabled: a maintainer decision on
2026-08-10 shelved it until the loop-performance work happened, and the pre-work
probe carried the discovery at Run time instead. With the probe now reachable at
authoring, the refusal belongs back in the staged registry — and turning it on
reports against active Specs that predate it, which this slice accounts for
deliberately rather than regenerating a golden on reflex.

## Requirements

1. MUST restore the vacuity refusal to the staged registry at the tasks stage.
2. MUST account every finding it reports across the active Spec corpus, either by
   correcting the command or by recording why the command stands.
3. MUST re-record the characterization corpus golden in the same change with its
   reason stated, and MUST NOT regenerate it silently.
4. MUST NOT weaken the detector to reduce its findings.

## Subtasks

- [ ] Restore the registry line.
- [ ] Account every finding across the active corpus.
- [ ] Re-record the corpus golden with its reason.

## Acceptance Criteria

- [ ] The refusal runs at the tasks stage.
- [ ] The active corpus reports no unaccounted vacuity finding.
- [ ] The corpus golden records why it was re-recorded, naming this Spec.
- [ ] The detector's own logic is unchanged.

## Verification

- `grep -q '^	{code: CodeVerifyVacuousCommand, stage: StageTasks},' internal/speccheck/coherence.go` — expected: exits 0, proving the registry line is restored and uncommented. Fails today, where the line is commented out.
- `grep -q 'Re-recorded because Spec 0095' internal/docscontract/testdata/corpus-golden.json` — expected: exits 0. The golden already carries earlier re-record reasons, so the sentence must name this Spec for the check to mean anything.
- `go build -buildvcs=false -o /tmp/0095-t06-roundfix ./cmd/roundfix && /tmp/0095-t06-roundfix spec check > /tmp/0095-t06.log 2>&1; grep -q 'SC-VERIFY-VACUOUS-COMMAND' /tmp/0095-t06.log && { echo 'unaccounted vacuity findings remain:'; grep 'SC-VERIFY-VACUOUS-COMMAND' /tmp/0095-t06.log; exit 1; }; grep -q '^	{code: CodeVerifyVacuousCommand, stage: StageTasks},' internal/speccheck/coherence.go || { echo 'FAIL: the sweep is clean only because the detector is still commented out'; exit 1; }` — expected: exits 0, proving the restored detector reports nothing unaccounted across the active corpus. It prints the offending findings on failure, and the clean sweep is anchored to the registry line so a still-disabled detector cannot pass it.

## Context

- interface: `internal/speccheck/coherence.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`

## References

`_techspec.md` → Build Order 6; Risks: restoring the vacuity detector moves the
corpus golden. `_prd.md` → Core Feature 4. Related:
`docs/backlog/2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md`,
which records the decision that disabled it.
