---
task: task_05
spec: 0136-a-rename-the-committer-can-stage
status: pending
type: backend
complexity: low
---

# Task 05: Report both sides of a rename in the prior-changed reader

## Overview

QA finding F-001: a governed rename auto-pushed under a record granting no
`push`, while a direct governed edit refused correctly. After integration the
Run asks for the paths changed between its initial head and HEAD with
`git diff --name-only`, and Git's rename detection answers with the destination
alone, so the governed source disappears before the final-push decision reads
it. The same reader feeds the changed-path audit against the grant, which is
also more correct with both sides present.

## Requirements

1. MUST report both the source and the destination of a rename, so a Run that
   renamed a Governed Path presents its governed side to the final-push
   decision.
2. MUST refuse the final push for a Run that renamed a Governed Path under a
   record that does not permit `push`.
3. MUST keep an ordinary Run passing its final push without `push` authority,
   which Task 02 established.
4. MUST keep a direct governed edit refusing its final push, which already
   works.
5. MUST keep the changed-path audit against the grant working, and MUST NOT
   change any refusal token.

## Subtasks

- [ ] Stop the prior-changed reader collapsing a rename.
- [ ] Prove the governed rename is refused at the final push.
- [ ] Prove the ordinary Run and the direct governed edit are unchanged.

## Acceptance Criteria

- [ ] The prior-changed reader reports both paths of a rename; today it reports
      the destination alone.
- [ ] A Run that renamed a Governed Path is refused at the final push without
      `push` authority; today it pushes.
- [ ] An ordinary Run still completes its final push without `push` authority.
- [ ] A direct governed edit is still refused at the final push.

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/cli/implement.go`

## Verification

- `grep -q 'func TestPriorChangedFilesReportsBothSidesOfARename' internal/worktree/worktree_test.go && go test -count=1 ./internal/worktree -run '^TestPriorChangedFilesReportsBothSidesOfARename$'` — both paths are reported; this fails today.
- `grep -q 'func TestFinalPushRefusesAGovernedRename' internal/cli/implement_test.go && go test -count=1 ./internal/cli -run '^TestFinalPushRefusesAGovernedRename$'` — the governed rename is refused at the final push; this fails today.
- `grep -q 'func TestFinalPushRefusesAGovernedRename' internal/cli/implement_test.go || exit 1; go test -count=1 ./internal/cli -run '^TestFinalPushAuthorityFollowsTheChangedPaths$'` — the ordinary Run and the direct governed edit keep Task 02's behavior.
- `grep -q 'func TestPriorChangedFilesReportsBothSidesOfARename' internal/worktree/worktree_test.go || exit 1; go test -count=1 ./internal/speccheck -run '^TestEnumeratedOutputsAreAuthoritative$'` — the changed-path audit that shares this reader still holds.

## References

- `_prd.md` → Goal 2; Core Feature 4; Regression locks.
- `_techspec.md` → Implementation Design: A rename survives the diff too;
  Testing Approach observation 7; Build Order 4.
