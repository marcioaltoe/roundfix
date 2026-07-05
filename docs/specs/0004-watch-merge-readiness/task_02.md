---
task: task_02
spec: 0004-watch-merge-readiness
status: pending
type: backend
complexity: medium
---

# Task 02: Deterministic stdout reports for watch and resolve

## Overview

Watch's stdout is empty and resolve's is inconsistent with the Implement
Command's report; scripts get only the exit code. Both commands gain the
deterministic report: per-Review-Issue lines plus one outcome line, shaped
like implement's. Verifiable through buffer-captured CLI tests.

## Requirements

1. MUST print, on stdout only and only at Run end: one line per Review Issue
   in Round then fetch order — `issue <id> <status> — <title>` with the final
   local status (`resolved|invalid|failed|duplicated|unresolved`) — then one
   outcome line `<Outcome> after <N> Round(s): <X> resolved, <Y> invalid,
   <Z> failed, <W> unresolved.`
2. MUST produce the report for every terminal outcome (Clean, Unresolved,
   MaxRoundsReached, TimedOut, BudgetExceeded, Stopped where issues were
   fetched; a Run stopped before any fetch prints only the outcome line).
3. MUST keep resolve's shape identical with `1 Round(s)`; fetch is untouched.
4. MUST keep all diagnostics on stderr and every exit code unchanged.
5. MUST be byte-deterministic for equal inputs (stable ordering, no
   timestamps).

## Subtasks

- [ ] Report rendering from cycle results (watch)
- [ ] Resolve alignment to the same shape
- [ ] Terminal-outcome matrix tests with byte-exact stdout asserts

## Acceptance Criteria

- [ ] Buffer-captured watch runs assert byte-exact stdout for Clean,
      Unresolved, and MaxRoundsReached fixtures.
- [ ] Resolve prints the same shape; its existing exit-code tests pass
      unchanged.
- [ ] stdout contains nothing before Run end and nothing non-report.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 3; Core Feature 3. `_techspec.md` → API Contracts,
Build Order 2. Dogfood finding 19.
