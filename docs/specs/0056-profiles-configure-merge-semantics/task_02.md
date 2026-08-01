---
task: task_02
spec: 0056-profiles-configure-merge-semantics
status: pending
type: backend
complexity: medium
---

# Task 02: Derive the effective change set

## Overview

The summary, the proof scope, and the merge must agree about what an operation
does; today nothing computes that once. This Task introduces the Effective
Change Set — per category, whether the operation adds, replaces, or removes it
— derived from the fragment and the new removal declaration. It is a pure value
with no write path attached yet, so it is verifiable entirely on its own.

## Requirements

1. MUST derive, from the fragment and the declared removals, one ordered set of
   per-category changes classified as added, replaced, or removed.
2. MUST classify a category the fragment names as replaced when the file
   already configures it and added when it does not.
3. MUST accept a repeatable removal declaration naming a category, and classify
   that category as removed.
4. MUST reject naming the same category in both the fragment and a removal
   declaration as a validation failure, before any proof or write.
5. MUST classify a removal of a category the file does not contain as a
   no-op-but-reported change rather than an error.
6. MUST order the set deterministically so the summary and any machine output
   are stable across runs.
7. MUST NOT change the write path, the confirmation, exit codes, or output in
   this Task.

## Subtasks

- [ ] Add the change-set value and its classification kinds.
- [ ] Derive added/replaced from the fragment against the existing file.
- [ ] Accept the repeatable removal declaration and derive removals.
- [ ] Reject the fragment/removal conflict.
- [ ] Make the ordering deterministic.

## Acceptance Criteria

- [ ] A fragment naming a category the file already has classifies as replaced;
      naming one it lacks classifies as added.
- [ ] A declared removal of a configured category classifies as removed.
- [ ] A declared removal of an absent category classifies as removed and is
      reported, without failing.
- [ ] Naming one category in both the fragment and a removal fails validation
      and names that category.
- [ ] The same inputs produce the same order on repeated derivation.
- [ ] The characterization corpus from task 01 is unchanged, proving the write
      path did not move.
- [ ] `git status --porcelain` shows no path outside `internal/config/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/config/profile_config.go`
- interface: `internal/cli/profiles_configure.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/config -run TestEffectiveChangeSet -count=1` — expected:
  exit 0; all five classification cases hold.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0; the writer's output is unchanged by this Task.
- `go test ./internal/config ./internal/cli -count=1` — expected: exit 0.
- `go vet ./internal/config ./internal/cli` — expected: exit 0.

## References

- `_prd.md` → Core Features 2; User Story 5.
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0086.
