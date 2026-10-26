---
task: task_02
spec: 0315-a-task-that-deletes-a-tracked-file-settles
status: completed
type: test
complexity: medium
---

# Task 02: A gated Task whose declared deletion is only on disk settles

## Overview

Oraculum's Task failed in the Daemon's Spec Consistency settlement check,
not in a direct call to the readers. This Task proves the fix through the
real settlement path: a gated `TaskCycle` with the default settlement checker,
a repository that retired a history file, and a Task whose fake Agent deletes
its declared `deletes:` path from disk only. It also proves the deletion is
still judged: a broken live link is still a finding, and an unauthorized
Governed Path deletion is still refused.

## Requirements

1. MUST answer the oraculum Inbox Entry of 2026-10-22,
   [`references/2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md`](references/2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md)
   (`_prd.md` → Sources 1), by proving Success Metrics 1-4 through the
   Daemon's settlement checks (API Contract 4, Invariant 3 of
   `_techspec.md`).
2. MUST add `internal/daemon/task_deletes_settlement_test.go` with a
   `Suite:`/`Invariant:`/`Boundary` header, one fixture built on
   `settlementRepositoryFixture` exactly as `_techspec.md` → Testing Approach
   describes (the recorded retirement rule with
   `spec.RetirementRuleSentinel`, one history Finding committed and then
   removed with `git rm` in a second commit, `docs/user-guide/index.md`
   linking `guide.md`, `docs/user-guide/guide.md`,
   `packages/backend/test/gone.test.ts`, and task_01 of the fixture Spec
   declaring `deletes: packages/backend/test/gone.test.ts`), and these tests,
   each top-level test and subtest calling `t.Parallel()` first:
   - `TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles` with subtests
     `disk-only` (the fake Agent's `afterTask` removes the file from disk)
     and `staged` (it runs `git rm -q --ignore-unmatch` on it); in both,
     `TaskCycle` returns no error, one completed Task and no failed Task;
   - `TestSettlementStillReportsALinkToADeletedFile`: after removing
     `docs/user-guide/guide.md` from disk,
     `SpecCheckSettlementChecker{}.RefusingFindings` returns no error and an
     `SC-HISTORY-RETAINED` finding whose summary names
     `docs/user-guide/index.md` and `guide.md`;
   - `TestSettlementRefusesAnUnauthorizedGovernedDeletion`: in a temporary
     repository whose authorization record lists only `Makefile` and which
     tracks `.agents/skills/escaped/SKILL.md`, after removing that file from
     disk, `SpecCheckSettlementChecker{}.AuditCommit` for a prospective
     commit changing it returns no error and exactly one `QA-AUTH-PATHS`
     finding naming it.
3. MUST make the `afterTask` callback safe on a repair turn: it reports with
   `t.Error`, never `t.Fatal`, and treats a file already gone as done, so the
   test fails instead of hanging when settlement refuses the first turn.
4. MUST NOT change any production file, existing test, golden, fixture,
   guide, skill or glossary entry.

## Subtasks

- [ ] Build the retired-history deletion fixture on the settlement fixture.
- [ ] Add the settlement, link-finding and authorization tests.

## Acceptance Criteria

- [ ] A gated Task whose declared deletion is only on disk completes through the real settlement checker, and so does a staged one.
- [ ] A deleted file's broken live link is a finding, not an error.
- [ ] An unauthorized Governed Path deletion is refused with `QA-AUTH-PATHS`.

## Consumers

None. The Task adds tests only; it changes no contract another package reads.

## Context

- instruction: `internal/daemon/settlement_spec_consistency.go`
- instruction: `internal/daemon/settlement_checks.go`
- instruction: `internal/daemon/settlement_checks_test.go`
- instruction: `internal/daemon/settlement_repository_verification_test.go`
- instruction: `internal/daemon/task_deletes.go`
- instruction: `internal/daemon/task_deletes_test.go`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/speccheck/history_retained.go`
- creates: `internal/daemon/task_deletes_settlement_test.go`

## Verification

- `out="$(go test -timeout 10m -count=1 -v -run '^(TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles|TestSettlementStillReportsALinkToADeletedFile|TestSettlementRefusesAnUnauthorizedGovernedDeletion|TestATaskThatLeavesItsDeletesPathDoesNotSettle|TestATaskThatRemovesItsDeletesPathSettles|TestDefaultSettlementCheckerAuditsAProspectiveCommit|TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles/disk-only TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles/staged TestSettlementStillReportsALinkToADeletedFile TestSettlementRefusesAnUnauthorizedGovernedDeletion TestATaskThatLeavesItsDeletesPathDoesNotSettle TestATaskThatRemovesItsDeletesPathSettles TestDefaultSettlementCheckerAuditsAProspectiveCommit TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/daemon` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check; without task_01's readers the `disk-only` subtest fails with oraculum's `inspect citing file packages/backend/test/gone.test.ts` refusal and the link test fails with an error. Afterwards the new tests pass with the unchanged deletes and settlement checker tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Feature 3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior; API Contract 4; Invariant 3; Testing Approach; Build Order 2

## Result

Added `internal/daemon/task_deletes_settlement_test.go` only; production code
and existing tests, fixtures, guides, skills and glossary remain unchanged.
The shared fixture extends `settlementRepositoryFixture` with the recorded
retirement rule using `spec.RetirementRuleSentinel`, a tracked Finding removed
with `git rm` in a subsequent commit, a linked user guide, and task_01's
declared TypeScript deletion. All repositories and commits created by the
tests are disposable fixtures.

Acceptance evidence from the focused implementation check:

- `TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles/disk-only` and
  `/staged` exercise the gated `TaskCycle` with the default settlement
  checker and assert no error, one completed Task and zero failed Tasks.
  The Agent callback uses nonfatal error reporting and tolerates repeated
  deletion on a repair turn (`os.ErrNotExist` or `git rm --ignore-unmatch`).
- `TestSettlementStillReportsALinkToADeletedFile` removes the tracked guide
  from disk and asserts no checker error plus an `SC-HISTORY-RETAINED`
  finding naming `docs/user-guide/index.md` and `guide.md`.
- `TestSettlementRefusesAnUnauthorizedGovernedDeletion` removes a tracked
  skill from disk under a grant listing only `Makefile`; the prospective
  commit audit returns no error and exactly one `QA-AUTH-PATHS` finding
  naming `.agents/skills/escaped/SKILL.md`.

Focused check:

```text
go test -count=1 -timeout 3m -run '^(TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles|TestSettlementStillReportsALinkToADeletedFile|TestSettlementRefusesAnUnauthorizedGovernedDeletion)$' ./internal/daemon
ok  roundfix/internal/daemon  4.208s
```

The new test file was absent before this change. `gofmt` formatted it, and
`git diff --check` exited 0. The only pre-existing change was the Daemon's
`status: in_progress` update in this Task file; that field was preserved.
The executor rejected the symlinked writable root; execution succeeded using
the resolved checkout with approved escalation. No configuration was changed.

Not run: the declared Verification commands (including Windows vet), selected
repository Verification and incremental Verification; these remain
Daemon-owned. No follow-up work identified.
