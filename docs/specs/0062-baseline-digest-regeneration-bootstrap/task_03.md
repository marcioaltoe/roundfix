---
task: task_03
spec: 0062-baseline-digest-regeneration-bootstrap
status: pending
type: backend
complexity: high
---

# Task 03: Break the regeneration cycle on the update path

## Overview

This is the slice that fixes the reported defect. The regeneration tests load
the catalog strictly, so a module edit that invalidates a derived pin makes the
load refuse and the refresh never runs — the command's remediation is itself.
This Task points the update path at the regeneration-mode entry point and adds
the fixture that reproduces the cycle, which fails before the change and passes
after it.

## Requirements

1. MUST make every regeneration step that runs under the update flag load the
   catalog through the regeneration-mode entry point, so a stale derived pin no
   longer blocks the run that rewrites it.
2. MUST leave every non-update test and every production path loading strictly.
3. MUST add a fixture that reproduces the reported cycle — a catalog whose
   generated guide changed, invalidating its formatter golden digest — and
   prove regeneration completes against it.
4. MUST prove the same fixture still fails strict validation when loaded
   outside regeneration mode, so the deferral is scoped and not a hole.
5. MUST leave the diagnostics of every other catalog unchanged.

## Subtasks

- [ ] Point the update-flag regeneration steps at the regeneration entry point.
- [ ] Add the cycle fixture reproducing the invalidated derived pin.
- [ ] Assert regeneration completes against the fixture.
- [ ] Assert the same fixture still fails strict validation.

## Acceptance Criteria

- [ ] The cycle fixture completes regeneration without the formatter digest
      mismatch stopping it.
- [ ] The same fixture, loaded outside regeneration mode, still produces the
      mismatch diagnostic.
- [ ] Regeneration steps that do not run under the update flag still load
      strictly.
- [ ] The characterization corpus from task 01 is unchanged.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/catalog_load.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestRegenerationBreaksGoldenDigestCycle -count=1`
  — expected: exit 0; the fixture that reproduces the reported cycle now
  regenerates instead of refusing.
- `go test ./internal/baseline -run TestCatalogRegenerationMode -count=1` —
  expected: exit 0; scoping from task 02 still holds.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; no other catalog's diagnostics moved.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 2; Success Metrics (adding a clause and running the
  regeneration once succeeds).
- `_techspec.md` → Testing Approach: cycle regression fixture; Build Order 3.
- ADR-0085.
