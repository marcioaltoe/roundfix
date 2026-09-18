---
task: task_01
spec: 0145-a-runner-that-owns-its-lock
status: pending
type: backend
complexity: medium
---

# Task 01: Give the runner's state one owner

## Overview

Sixteen methods take the acpx runner by value while the only holder keeps a
pointer. Each call copies a struct whose maps survive the copy and whose mutex
does not, so the guard is per-call over shared maps. This slice converts those
receivers to pointers and leaves behavior alone.

## Requirements

1. MUST convert every method of the runner that takes it by value to a pointer
   receiver.
2. MUST leave session ensure, warm, work-start, selection assignment, probing,
   codex spawn resolution and cancellation behaving exactly as they do today.
3. MUST NOT rename, add or remove any exported identifier of the package.
4. MUST NOT edit an existing assertion to make a test pass; a test helper may
   gain an address-of where the compiler requires one.
5. MUST NOT add a linter, an analyzer configuration, or a repository gate that
   runs one.

## Subtasks

- [ ] Convert the value receivers.
- [ ] Fix whatever the compiler names at call sites and in test helpers.
- [ ] Run the analyzer over the package and read an empty report.
- [ ] Run the package's tests with no assertion edited.

## Acceptance Criteria

- [ ] The analyzer reports no copied-lock diagnostic for this package.
- [ ] The package's tests pass with their assertions unchanged.
- [ ] No exported identifier changes.
- [ ] No tooling configuration or repository gate is added.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/agent.go`

## Verification

- `report="$(go vet ./internal/agent 2>&1 | grep "passes lock by value")"; test -z "$report"` — expected: exit 0; the analyzer names no copied lock. Before this Task it names thirty-one, so the command fails.
- `count="$(grep -c "^func (runner ACPXRunner)" internal/agent/acpx_runner.go)"; test "$count" = "0"` — expected: exit 0; no method takes the runner by value. Before this Task sixteen do.
- `count="$(grep -c "^func (runner ACPXRunner)" internal/agent/acpx_runner.go)"; test "$count" = "0" && go test -count=1 ./internal/agent` — expected: exit 0; every receiver is a pointer and the package's behavior is unchanged. Before this Task sixteen methods take the runner by value, so the command fails.

## References

`_prd.md` → Core Features 1-3; User Stories 1-3; Goals 1-3; Success Metrics 1-3;
Declared intentional breaks; Regression locks;
`_techspec.md` → Implementation Design: One receiver kind, What must not move;
API Contracts 1-3; Build Order 1; ADR-0020.
