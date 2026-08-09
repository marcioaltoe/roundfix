---
task: task_01
spec: 0091-a-proof-that-can-refuse
status: pending
type: test
complexity: medium
---

# Task 01: Record that a nonexistent claude model proves passed

## Overview

The corpus this Spec exists to break. It records today's asymmetry: a
nonexistent model proves `passed` on `claude` and is refused on `codex`, from
the same profile shape. It also records the contaminated read that causes it —
an adapter reporting a requested model as current and appending it to its own
advertised options.

## Requirements

1. MUST record that a `claude` Agent Selection naming a model the runtime does
   not offer proves successfully today.
2. MUST record that the equivalent `codex` selection is refused with
   `model_not_advertised`, so the asymmetry is the captured fact rather than the
   `claude` behaviour alone.
3. MUST record the contaminated capability payload as a fixture: a response
   whose `currentValue` and options both contain a model the runtime does not
   offer.
4. MUST declare each break in a documentation comment naming the Task that
   changes it.
5. MUST NOT change any production behaviour.

## Subtasks

- [ ] Capture the passing `claude` proof.
- [ ] Capture the refused `codex` proof.
- [ ] Add the contaminated payload fixture.

## Acceptance Criteria

- [ ] A test proves the `claude` selection succeeds today.
- [ ] A test proves the `codex` selection is refused today.
- [ ] A fixture carries a capability payload whose advertised options include a
      model the runtime does not offer.
- [ ] Each declared break names its Task.

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_catalogue_characterization_test.go`
- `docs/specs/0091-a-proof-that-can-refuse/task_01.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestSelectionCatalogueCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSelectionCatalogueCharacterizationClaudeProvesAnUnofferedModel'` — expected: exits 0. A `-run` pattern selecting no cases exits 0, so this asserts the named case ran.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestSelectionCatalogueCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel'` — expected: exits 0.
- `grep -c 'Declared break: task_0' internal/agent/selection_catalogue_characterization_test.go | grep -qE '^[2-9]'` — expected: exits 0, proving both breaks are declared.

## References

- `_prd.md` → Goal 1.
- `_techspec.md` → Build Order 1; Testing Approach.
