---
task: task_01
spec: 0027-review-loop-integrity
status: pending
type: backend
complexity: low
---

# Task 01: Add the CleanUnverified terminal state and exit code 3

## Overview

Groundwork for the Merge-Ready confirmation change: register Clean Unverified as a terminal Run state and reserve process exit code 3 for it. No command produces the state yet — this task makes the state and its exit mapping exist and provably correct so later work can emit it.

## Requirements

1. MUST add a Clean Unverified state to the Run state vocabulary, the terminal-state predicate, and the journal's recognized state list.
2. MUST map the new state to process exit code 3 in the watch outcome-to-exit-code mapping, leaving the existing contract untouched (0 success, 1 run failed, 2 preflight, 130 SIGINT).
3. MUST leave every existing outcome's exit code unchanged, proven by tests.
4. SHOULD name the state consistently with the glossary term Clean Unverified.

## Subtasks

- [ ] Add the state constant and include it in the terminal-state set
- [ ] Register the state wherever terminal states are enumerated (journal, state listings)
- [ ] Extend the watch exit-code mapping with the new state → 3
- [ ] Table-test the full outcome → exit-code matrix, including the new state and all pre-existing states

## Acceptance Criteria

- [ ] The terminal-state predicate returns true for the new state
- [ ] A test asserts the new state maps to exit code 3 and every pre-existing state keeps its current code
- [ ] The full test suite passes

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/journal.go`
- interface: `internal/cli/cli.go`

## Verification

- `rg -q "CleanUnverified" internal/store/store.go` — expected: exit 0 (state constant exists)
- `go test ./internal/store/... ./internal/cli/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goals 2, User Story 5, Core Feature 5; `_techspec.md` → Build Order 1, Data Models (Run state), API Contracts (exit code 3); ADR-0043.
