---
type: fix
status: promoted
created: 2026-09-25
spec: 0170-authoring-that-fails-before-dispatch
reason: null
---

# A section declared `None.` fails when its reason wraps

## Symptom

Spec check accepts a `None.` declaration only on a single line, so a normally wrapped reason triggers `SC-METRIC-UNDECLARED` and authors invent metrics.

## Where

`internal/speccheck/citations.go` (`nonEmptyLines == 1` predicate).

## Expected

Accept a wrapped `None.` paragraph; keep refusing missing, empty or mixed sections.

## Evidence

secondbrain `inbox/roundfix/2026-09-18-o-none-declarado-so-vale-numa-linha.md`.
