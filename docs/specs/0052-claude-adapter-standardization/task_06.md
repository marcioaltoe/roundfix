---
task: task_06
spec: 0052-claude-adapter-standardization
status: pending
type: backend
complexity: low
---

# Task 06: Explain the Preflight fallback boundary

## Overview

A Preflight refusal caused by a failed selection proof currently says nothing
about why the configured Fallback Chain did not rescue the Run, leaving
operators to suspect the fallback is broken. Append one sentence to the
profile proof failure explaining the ADR-0050 boundary: Fallback Chains
activate only after Run creation, and Preflight proves every configured tuple
without substituting any.

## Requirements

1. MUST append to the profile proof failure message a sentence stating that
   Fallback Chains activate only after Run creation per ADR-0050 and that
   Preflight substitutes none.
2. MUST keep every existing field and wording of the message (runtime,
   model, reasoning effort, affected categories, classification, adapter
   error, next action) byte-identical apart from the added sentence.
3. MUST apply to every consumer of the shared proof error — the Implement,
   Resolve, and Watch preflights and `profiles validate` — without
   per-command duplication.
4. SHOULD keep the added sentence out of the `profiles-validate` JSON
   schema's structured fields; it is message prose, not a new field.

## Subtasks

- [ ] Extend the proof-failure message construction with the fallback
      boundary sentence.
- [ ] Update the preflight-message assertions across the affected command
      tests.
- [ ] Confirm the JSON surfaces remain schema-stable.

## Acceptance Criteria

- [ ] A failed profile proof printed by an operational preflight contains
      the fallback boundary sentence naming ADR-0050.
- [ ] The pre-existing message fields keep their exact text and order.
- [ ] `profiles validate --json` output remains valid against its current
      schema with no new structured field.

## Context

- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/ -run 'TestProfilesValidate|TestImplement'` — expected: pass with the new message assertions.
- `grep -n 'ADR-0050' internal/cli/profiles_validate.go` — expected: at least one match; the boundary sentence names its ADR.

## References

`_prd.md` → User Story 6, Core Feature 8; `_techspec.md` → Build Order 6,
API Contracts; ADR-0050.
