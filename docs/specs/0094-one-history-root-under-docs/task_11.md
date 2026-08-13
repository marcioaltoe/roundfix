---
task: task_11
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: low
---

# Task 11: Repoint the derived-artifact test's archive path

## Overview

The last pinned archive path. A repository-wide sweep after Task 10 found one
remaining literal: the owned-Skill derived-artifact test names `_archived/specs`
twice while proving that editing an owned Skill leaves archived artifacts
byte-identical, so the authorized relocation left it inspecting a directory that
no longer exists and the documentation gate is red.

This is a Task of its own rather than a widening of Task 10, whose work is
complete and merged. Widening a finished Task makes its satisfied gates measure
an already-done state, which is how three earlier attempts were refused as
vacuous.

## Requirements

1. MUST make the derived-artifact test inspect the archive at the location the
   resolver answers.
2. SHOULD read that location from the resolver rather than replacing one literal
   with another, so the next relocation does not break it again.
3. MUST keep proving the same property: that editing an owned Skill leaves
   archived artifacts byte-identical.
4. MUST NOT weaken, delete, or skip the assertion.
5. MUST NOT change any repository path outside the bounded scope below plus this
   Task file; stop and fail the Task if a changed-file check finds another path.

## Subtasks

- [ ] Point the derived-artifact test at the resolved archive location.
- [ ] Confirm it still proves archived artifacts stay byte-identical.

## Acceptance Criteria

- [ ] The derived-artifact test passes.
- [ ] Neither `_archived/specs` literal survives in that file.
- [ ] The test still compares archived artifact bytes before and after an owned
      Skill edit, with no expectation removed or made conditional.
- [ ] The changed-file set is the bounded scope plus this Task file.

## Bounded scope

This Task may create or modify only:

- `skills/owned_skill_edit_repocontract_test.go`
- `docs/specs/0094-one-history-root-under-docs/task_11.md`

Express maintainer authorization:
`docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`,
extended on 2026-08-13 for exactly this file and bounded to its two pinned
archive-path literals.

## Verification

- `go test -count=1 -tags repocontract ./skills -run 'TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical' -v > /tmp/0094-task-11.log 2>&1; s=$?; grep -q '^--- PASS: TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical' /tmp/0094-task-11.log || { cat /tmp/0094-task-11.log; exit 1; }; exit $s` — expected: exits 0. Fails today, where the test inspects a directory the authorized relocation removed.
- `! grep -q '_archived/specs' skills/owned_skill_edit_repocontract_test.go` — expected: exits 0, proving neither literal survives.
- `grep -q 'archived Spec artifact' skills/owned_skill_edit_repocontract_test.go` — expected: exits 0, proving the byte-identity assertion is still there rather than removed along with the path.
- `git diff --name-only HEAD > /tmp/0094-task-11-all.txt; test -s /tmp/0094-task-11-all.txt || { echo 'no file changed'; exit 1; }; grep -v -e '^skills/owned_skill_edit_repocontract_test\.go$' -e '^docs/specs/0094-one-history-root-under-docs/task_11\.md$' /tmp/0094-task-11-all.txt > /tmp/0094-task-11-scope.txt; test ! -s /tmp/0094-task-11-scope.txt || { cat /tmp/0094-task-11-scope.txt; exit 1; }` — expected: exits 0, proving work happened and every changed path is in bounds.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`
- interface: `skills/owned_skill_edit_repocontract_test.go`
- interface: `internal/spec/archive.go`

## References

`_prd.md` → Core Feature 10; Project Constraints: Tooling authority.
`_techspec.md` → Build Order 7. Like Task 10, this is a consequent fix landing
after the authorized change that made it necessary, in its own commit.
