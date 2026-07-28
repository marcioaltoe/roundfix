---
task: task_03
spec: 0052-claude-adapter-standardization
status: pending
type: backend
complexity: low
---

# Task 03: Default the frontend profile to the proven Claude tuple

## Overview

Change the built-in `frontend` Agent Selection Profile's Preferred Selection
from the unadvertised `claude / claude-opus-5 / xhigh` to the proven
`claude / opus / xhigh`, so a fresh installation's frontend routing proves
against the official adapter without configuration edits. The Fallback Chain
and every other built-in profile are untouched.

## Requirements

1. MUST change the built-in `frontend` Preferred Selection model to `opus`,
   keeping runtime `claude` and reasoning effort `xhigh`.
2. MUST keep the `frontend` Fallback Chain and all other built-in profiles
   byte-identical in behavior.
3. MUST keep the rendered default Project Config in agreement with the
   built-in value.
4. MUST NOT add adapter aliases to the Model Catalog; `opus` is proven
   through Exact Agent Selection Proof, not catalog membership.

## Subtasks

- [ ] Update the built-in `frontend` profile value.
- [ ] Update the default Project Config rendering that mirrors it.
- [ ] Update every test that pins the old `claude-opus-5` frontend default.

## Acceptance Criteria

- [ ] Resolving the `frontend` profile with no user or project configuration
      yields preferred `claude / opus / xhigh` with the unchanged fallback.
- [ ] The generated default Project Config names `opus` for the frontend
      preferred model.
- [ ] The Model Catalog's Claude identifiers are unchanged.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`
- interface: `internal/agent/catalog.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/config/` — expected: pass.
- `grep -n 'claude-opus-5' internal/config/profiles.go ; test $? -eq 1` — expected: no matches (exit 1).
- `grep -c 'claude-opus-5' internal/agent/catalog.go | grep -x 1` — expected: `1`; the catalog identifier is preserved.

## References

`_prd.md` → User Story 5, Core Feature 6; `_techspec.md` → Build Order 3,
Data Models; ADR-0049.
