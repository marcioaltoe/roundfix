---
task: task_03
spec: 0077-a-green-check-is-not-a-review
status: pending
type: backend
complexity: low
---

# Task 03: Say what was observed instead of merging

## Overview

A stalled loop that does not say why is a stalled loop an operator will
override. This slice makes the refusal and the unclassified case legible where
the operator is already looking: the watch output and the Run Event Stream.

## Requirements

1. MUST make `watch --until-clean` name what it observed when it declines to
   proceed — the refusal and its reason, or that the signal was unrecognised.
2. MUST record the evidence state and reason in the Run Event Stream, so an
   operator can tell a refused review from an absent one without opening GitHub.
3. MUST NOT change what `watch` does on `verified` or `reviewed`.
4. MUST NOT push or merge on a `skipped` or `pending` head, asserted rather than
   assumed.

## Subtasks

- [ ] Name the observation in the watch diagnostic.
- [ ] Record state and reason in the Run Event Stream.
- [ ] Assert no push and no merge on refused or unclassified heads.

## Acceptance Criteria

- [ ] A refused head produces a diagnostic naming the refusal and its reason.
- [ ] An unclassified head produces a diagnostic saying the signal was not
      recognised.
- [ ] The Run Event Stream carries the evidence state and reason for both.
- [ ] Neither case pushes or merges, asserted by a fixture that would have
      merged before this Spec.
- [ ] A verified head behaves exactly as it does today.

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/runevent`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/watch -count=1 -run 'Skipped|Pending|Diagnostic|Clean' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the watch tests ran and passed.
- `go test ./internal/watch ./internal/runevent -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; the Skill is task_04's bounded scope.

## References

- `_prd.md` → Core Features 4 and 6; Goals.
- `_techspec.md` → Build Order 3.
