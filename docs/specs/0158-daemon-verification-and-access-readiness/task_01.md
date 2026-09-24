---
task: task_01
spec: 0158-daemon-verification-and-access-readiness
status: completed
type: backend
complexity: medium
---

# Task 01: Collect every independent Verification failure

## Overview

A Task's Verification attempt returns at its first failed command, so the one repair turn ADR-0038 allows sees a single failure even when later commands were already failing for their own reasons.

## Requirements

1. MUST accept Task frontmatter `verification: independent`, record it on `spec.Task`, and refuse any other value by name.
2. MUST, for a Task so declared, run every Verification command even after one fails and publish each failure as today.
3. MUST hand every failed command, with its diagnostics, to the single repair turn, and run every command again on the retry.
4. MUST keep today's behaviour for a Task without the declaration: commands run in order and the attempt stops at the first failure.
5. MUST keep the retry ceiling and the handling of a temporary failure unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A declared Task with two failing commands of three hands both failures to its repair turn, and the third command ran.
- [ ] An undeclared Task stops at its first failure and later commands never run.
- [ ] An unknown `verification` value is refused by name.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskParsesIndependentVerification|TestIndependentVerificationHandsEveryFailureToRepair|TestUndeclaredVerificationStopsAtFirstFailure)$" ./internal/spec ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestTaskParsesIndependentVerification TestIndependentVerificationHandsEveryFailureToRepair TestUndeclaredVerificationStopsAtFirstFailure; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the three cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Independent Verification

## Result

Implemented Task-declared independent Verification. The Task parser records
`verification: independent` on `spec.Task` and rejects any other non-empty
value by name. Independent attempts retain command-specific diagnostics,
publish every deterministic command failure, send the ordered failure set to
the single Verification Feedback turn, and rerun the complete command list on
attempt 2. Undeclared attempts still stop at the first failure; temporary and
unobserved failures keep their existing short-circuit paths.

Focused-check evidence:

- Pre-change: `GOCACHE=/private/tmp/roundfix-0158-task01-gocache go test -count=1 -run '^TestTaskParsesIndependentVerification$' ./internal/spec` failed to compile because `Task.VerificationMode` and `VerificationModeIndependent` did not exist.
- `GOCACHE=/private/tmp/roundfix-0158-task01-gocache go test -count=1 -run '^TestTaskParsesIndependentVerification$' ./internal/spec` passed. This covers recording the declaration and refusing `verification: sequential-groups` with that value in the error.
- `GOCACHE=/private/tmp/roundfix-0158-task01-gocache go test -count=1 -run '^TestIndependentVerificationHandsEveryFailureToRepair$' ./internal/daemon` passed. A three-command declared Task published failures for commands one and three, retained distinct diagnostic paths for both, included both in one repair prompt, and ran all three commands again on attempt 2.
- `GOCACHE=/private/tmp/roundfix-0158-task01-gocache go test -count=1 -run '^TestIndependentVerificationKeepsTemporaryRetryHandling$' ./internal/daemon` passed. A temporary failure stopped the initial command sequence, restarted the complete sequence under the existing exclusive retry, and consumed no Agent repair turn.
- `GOCACHE=/private/tmp/roundfix-0158-task01-gocache go test -count=1 -run '^TestUndeclaredVerificationStopsAtFirstFailure$' ./internal/daemon` passed. Both bounded attempts stopped on command one and never ran commands two or three.
- `GOCACHE=/private/tmp/roundfix-0158-task01-gocache go test -count=1 ./internal/spec ./internal/agent ./internal/daemon` passed, including the existing retry-ceiling, temporary-failure, unknown-verdict, repeated-failure, and single-failure prompt coverage.

The Daemon-owned `## Verification` command was not run in this Agent turn.
