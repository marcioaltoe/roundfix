---
task: task_03
spec: 0070-declared-unreachable-acceptance
status: pending
type: backend
complexity: high
---

# Task 03: Archive the declared case and refuse the rest

## Overview

The archive boundary gains exactly one accepting case: a `partial` report whose
only unmet rows are declared unreachable, with no finding-blocked and no
environment-blocked row, and whose declared count is covered by the Spec's own
declarations. It also stamps what remains unproven, so the archive record
carries the debt instead of hiding it.

This Spec widens an archive gate, so this slice is written as refusals with one
exception. The fixtures that must still refuse are worth more than the one that
must now accept.

## Requirements

1. MUST archive a Spec whose newest report is `partial` with
   `rows_blocked_declared` greater than zero, `rows_blocked_finding` zero,
   `rows_blocked_environment` zero, and at least as many declarations as
   declared rows.
2. MUST refuse, with a diagnostic naming the cause, when: any row is
   finding-blocked; any row is environment-blocked; the declared count exceeds
   the Spec's declarations; or the verdict is `fail`.
3. MUST keep environment-blocked rows blocking. Circumstance is not
   unreachability, and this is the refusal most likely to be relaxed by
   accident.
4. MUST stamp the archived Spec with the satisfying action of every declaration
   that covered an unmet row, so a reader of the archive learns what was never
   verified.
5. MUST leave `qa_override` unchanged in meaning, reach, and wording.
6. MUST archive every Spec that archives today under `pass` identically, proven
   over the archived corpus rather than one fixture.
7. SHOULD make each refusal diagnostic name the count that caused it, so the
   operator does not have to open the report to learn which rule fired.

## Subtasks

- [ ] Add the accepting case to the archive precondition.
- [ ] Stamp the unproven actions into the archived Spec.
- [ ] Add the refusal diagnostics naming each cause.
- [ ] Build the fixture matrix over the whole API Contracts table.
- [ ] Assert non-regression over the archived corpus.

## Acceptance Criteria

- [ ] A `partial` report with declared rows only, fully covered by
      declarations, archives and stamps every satisfying action.
- [ ] A `partial` report with any finding-blocked row refuses, naming the
      finding count.
- [ ] A `partial` report with any environment-blocked row refuses, naming the
      environment count.
- [ ] A report declaring more blocked rows than the Spec declares refuses,
      naming the shortfall.
- [ ] A `fail` report refuses exactly as today.
- [ ] A `pass` report archives exactly as today, with no stamp added.
- [ ] `qa_override` archives a genuinely failed Spec, unchanged.
- [ ] Every Spec in the archived corpus still satisfies the archive
      precondition it satisfied before this Task.

## Context

- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -count=1 -run 'Archive' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the archive tests ran and passed.
- `go test ./internal/cli -count=1 -run 'Archive' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the command tests ran and passed.
- `go test ./internal/spec ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; this Task touches no skill; the gate contract is task_04's
  bounded scope.

## References

- `_prd.md` → Core Features 3, 4, 5 and 6; Decisions; Success Metrics 2 and 4.
- `_techspec.md` → API Contracts; Build Order 3; Risks & Considerations.
- ADR-0080.
