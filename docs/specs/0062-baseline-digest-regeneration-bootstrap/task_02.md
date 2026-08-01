---
task: task_02
spec: 0062-baseline-digest-regeneration-bootstrap
status: pending
type: backend
complexity: high
---

# Task 02: Give the catalog loader a regeneration mode

## Overview

Catalog validation is unconditionally strict, which is correct everywhere
except inside the run whose purpose is to rewrite the very pins it checks. This
Task adds an explicit regeneration mode to the loader and an enumerable
allowlist of the diagnostic codes that mode defers. Nothing consumes the mode
yet — it is verifiable on its own by proving that the mode defers exactly the
allowlisted codes and that the default path is unchanged.

## Requirements

1. MUST add an explicit regeneration mode to the catalog loader, defaulting to
   off, so every existing caller keeps today's strict behavior with no change
   at the call site.
2. MUST express the deferred set as an enumerable allowlist keyed by diagnostic
   code, containing exactly the formatter golden-digest mismatch code, so the
   blast radius can be read off the list rather than inferred from a severity.
3. MUST defer an allowlisted code only when regeneration mode is on; a code
   outside the allowlist MUST still be recorded in regeneration mode.
4. MUST expose the regeneration-mode load through an unexported entry point, so
   no package outside the Baseline package and no production, CLI, or CI caller
   can reach it.
5. MUST leave the public catalog-loading API's signature and behavior
   unchanged.
6. MUST NOT wire the update path to the new entry point in this Task; that is
   the next slice.

## Subtasks

- [ ] Add the regeneration flag to the loader, defaulting to off.
- [ ] Add the deferred-code allowlist with the single formatter digest code.
- [ ] Apply the deferral at the diagnostic sink, gated on the flag.
- [ ] Add the unexported regeneration-mode load entry point.

## Acceptance Criteria

- [ ] With regeneration mode off, a catalog whose formatter golden digest
      disagrees with its fixtures still records the mismatch diagnostic.
- [ ] With regeneration mode on, the same catalog records no mismatch
      diagnostic.
- [ ] With regeneration mode on, a diagnostic whose code is absent from the
      allowlist is still recorded.
- [ ] The regeneration entry point is unexported and unreachable from outside
      the package.
- [ ] The public catalog-loading API's signature is unchanged.
- [ ] The characterization corpus from task 01 is unchanged, proving the
      default path did not move.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/catalog_load.go`
- interface: `internal/baseline/catalog_validate.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestCatalogRegenerationMode -count=1` —
  expected: exit 0; deferral applies only in regeneration mode and only to the
  allowlisted code.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; the strict default path produces the same diagnostics as
  before this Task.
- `grep -q "goldenDigest.mismatch" internal/baseline/catalog_validate.go` —
  expected: exit 0; the code is still emitted, not deleted.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `go vet ./internal/baseline` — expected: exit 0.

## References

- `_prd.md` → Core Features 1; Goals (outside regeneration still fails closed).
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0085.
