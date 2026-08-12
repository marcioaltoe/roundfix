---
task: task_10
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: low
---

# Task 10: Repoint the pinned archive-path assertion

## Overview

The consequent fix of Task 07. That Task's authorized skill edit moved the
archive path, and the skill contract test pins the old path as a literal, so the
assertion now names a location the contract no longer has and the gate is red.
This is its own Task rather than a reopening of Task 07, because Task 07's work
is complete and its gates would measure an already-finished state.

## Requirements

1. MUST update the pinned archive-path literal in the skill contract test so it
   names the location the resolver answers.
2. MUST keep asserting the same reference-lifecycle contract; only where the
   assertion says the archive lives may change.
3. MUST NOT weaken, delete, or skip the assertion, and MUST NOT change any other
   expectation in that test.
4. MUST NOT change any repository path outside the bounded scope below plus this
   Task file; stop and fail the Task if a changed-file check finds another path.

## Subtasks

- [ ] Repoint the pinned archive-path literal.
- [ ] Confirm the contract test passes and still asserts the contract.

## Acceptance Criteria

- [ ] The skill contract test passes.
- [ ] The assertion still requires the reference-lifecycle contract, with only
      the archive location updated.
- [ ] No expectation was deleted or made conditional.
- [ ] The changed-file set is the bounded scope plus this Task file.

## Bounded scope

This Task may create or modify only:

- `skills/baseline_skill_contract_test.go`
- `docs/specs/0094-one-history-root-under-docs/task_10.md`

Express maintainer authorization:
`docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`,
extended on 2026-08-12 for exactly this file and bounded to the assertion that
pins the archive path.

## Verification

- `go test -count=1 ./skills -run 'TestSpecReferenceLifecycleSkillContracts' -v > /tmp/0094-task-10.log 2>&1; s=$?; grep -q '^--- PASS: TestSpecReferenceLifecycleSkillContracts' /tmp/0094-task-10.log || { cat /tmp/0094-task-10.log; exit 1; }; exit $s` — expected: exits 0. Fails today, where the assertion pins a path the skills no longer name.
- `grep -q 'from automatic link rewrites' skills/baseline_skill_contract_test.go && grep -q 'docs/history/specs' skills/baseline_skill_contract_test.go` — expected: exits 0, proving the assertion still requires the reference-lifecycle contract and now names the relocated path. Both clauses are one command because either alone is satisfiable without the other: the contract text is present today with the wrong path, and a bare path match would pass on an assertion that had been gutted. A `grep -c` stood here and was vacuous — it exits 0 whenever it prints a count, including a count that proves nothing changed.
- `git diff --name-only HEAD > /tmp/0094-task-10-all.txt; test -s /tmp/0094-task-10-all.txt || { echo 'no file changed'; exit 1; }; grep -v -e '^skills/baseline_skill_contract_test\.go$' -e '^docs/specs/0094-one-history-root-under-docs/task_10\.md$' /tmp/0094-task-10-all.txt > /tmp/0094-task-10-scope.txt; test ! -s /tmp/0094-task-10-scope.txt || { cat /tmp/0094-task-10-scope.txt; exit 1; }` — expected: exits 0, proving work happened and every changed path is in bounds.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`
- interface: `skills/baseline_skill_contract_test.go`

## References

`_prd.md` → Core Feature 10; Project Constraints: Tooling authority.
`_techspec.md` → Build Order 7. This Task is the consequent fix the tooling
chronology rule requires to land after the authorized change that made it
necessary, never folded into it.
