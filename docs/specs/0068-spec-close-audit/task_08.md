---
task: task_08
spec: 0068-spec-close-audit
status: pending
type: backend
complexity: medium
---

# Task 08: Reclaim the scratch worktree whose work already survives

## Overview

QA finding F-001. A Supervisor scratch worktree whose branch is pushed and
whose content is merged classifies `preserved` and carries no reclaim command,
so the case Core Feature 4 exists for cannot be closed. The emitted branch
deletion also cannot succeed while that branch is still checked out in the
worktree.

The root cause is in the Task that specified it: task_02's Requirement 3
demanded `preserved` for every worktree without a Run, while its own References
note said a pushed-and-merged one should be `residue`. The implementation
followed the MUST, correctly. Requirement 3 is corrected; this Task delivers
the distinction.

## Requirements

1. MUST classify a worktree with no matching Run as `residue` when its branch
   is pushed and its content is merged into the default branch, and
   `preserved` otherwise.
2. MUST prove the merge by content before offering reclamation, using the same
   prove-only-integration asymmetry task_01 established. Ambiguity resolves to
   `preserved`.
3. MUST emit a reclaim command whose steps are ordered so they can actually
   run: remove the worktree first, then delete the branch. A branch checked out
   in a worktree cannot be deleted, and an unrunnable command is worse than
   none.
4. MUST still classify as `preserved` a worktree whose branch is unpushed, whose
   content is unmerged, or whose state cannot be determined.
5. MUST NOT execute the reclaim. The audit reports; the operator reclaims.
6. MUST NOT weaken the Active Run guard: a worktree belonging to an Active Run
   is never `residue`, whatever its branch state.

## Subtasks

- [ ] Add the pushed-and-merged determination for Run-less worktrees.
- [ ] Emit the ordered reclaim command.
- [ ] Extend the fixtures: merged, unpushed, unmerged, indeterminate, active.

## Acceptance Criteria

- [ ] A scratch worktree whose branch is pushed and merged classifies
      `residue` and carries a reclaim command that removes the worktree before
      deleting the branch.
- [ ] Running that emitted command in a fixture succeeds, proving it is
      runnable rather than merely printed.
- [ ] A scratch worktree whose branch is unpushed classifies `preserved`.
- [ ] A scratch worktree whose content is unmerged classifies `preserved`.
- [ ] A worktree belonging to an Active Run is never `residue`.
- [ ] The audit still mutates nothing, proven by Git state identical before and
      after.

## Context

- interface: `internal/specaudit`
- interface: `internal/gittest`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/specaudit -count=1 -run 'Scratch|Reclaim|Preserved|Active' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the scratch-worktree tests ran and passed.
- `go test ./internal/specaudit ./internal/worktree ./internal/cli -count=1`
  — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if grep -rn "os.RemoveAll\|worktree remove" internal/specaudit --include="*.go" | grep -v "_test.go" | grep -v "fmt.Sprintf\|const \|reclaim" | grep -q .; then exit 1; fi`
  — expected: exit 0; the reclaim stays a string, never a call.

## References

- `_prd.md` → Core Feature 4; Decisions.
- `_techspec.md` → Interfaces; Risks & Considerations.
- `qa/qa-report-2026-08-04.md` → F-001 and its required repair.
