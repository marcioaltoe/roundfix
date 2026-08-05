---
task: task_04
spec: 0066-run-teardown-reclaims-what-it-created
status: pending
type: backend
complexity: medium
---

# Task 04: Offer both debris kinds through reconcile

## Overview

`reconcile` is already dry-run first, already enumerates without a Run ID,
already carries proof per candidate, already preserves ambiguity, and already
refuses to apply what it did not inspect. This slice adds the two new debris
kinds to it rather than minting a command that would have to re-learn all four
properties.

## Requirements

1. MUST offer orphaned process candidates and releasable Run Branch candidates
   beside the existing worktree candidates.
2. MUST keep dry-run the default and `--apply` the only acting mode.
3. MUST carry, for every new candidate, the proof that makes it reclaimable.
4. MUST preserve anything ambiguous and report it as preserved.
5. MUST never offer or act on anything belonging to an Active Run, per ADR-0052.
6. MUST be idempotent: a second pass after applying is a no-op.
7. MUST extend the `--format json` schema additively, leaving existing fields
   and their meanings unchanged.

## Subtasks

- [ ] Add the two candidate kinds with their proofs.
- [ ] Wire them through dry-run and apply.
- [ ] Extend the JSON schema additively.
- [ ] Assert idempotence and the Active Run guard.

## Acceptance Criteria

- [ ] A dry-run names every candidate of both kinds with its proof and changes
      nothing.
- [ ] `--apply` reclaims exactly the named candidates.
- [ ] A second pass after applying is a no-op.
- [ ] An Active Run's process and branch are never offered.
- [ ] An ambiguous candidate is reported preserved, never offered.
- [ ] Existing `--format json` fields keep their names and meanings, asserted
      against the current schema.

## Context

- interface: `internal/cli/reconcile.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1 -run 'Reconcile' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the reconcile tests ran and passed.
- `go test ./internal/cli -count=1` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix reconcile --format json | grep -q "runs\|Run"`
  — expected: exit 0; the command still runs and emits its report.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; the Skill is task_05's bounded scope.

## References

- `_prd.md` → Core Features 5 and 6; Success Metrics 3 and 4.
- `_techspec.md` → API Contracts; Build Order 4.
- ADR-0052.
