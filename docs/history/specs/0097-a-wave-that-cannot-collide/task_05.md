---
status: completed
type: backend
---

# Task: A worktree failure in Roundfix's own words

Creating a second Task Worktree under suite load fails with a raw filesystem
errno that names neither the Run, the Task, nor the concurrency that produced
it. A Supervisor reading it has to translate a filesystem error into the loop's
terms before they can act.

## Work

- A Task Worktree that cannot be created reports the Run, the Task, and the
  configured concurrency.
- The underlying error travels as evidence, not as the message. A reader needs
  both: what the filesystem said, and what Roundfix was doing when it said it.
- Follow the repository's existing shape for an actionable failure rather than
  inventing one — the loop already has refusals that name their next action.
- Cover the message naming all three facts and still carrying the underlying
  error, asserted through the error rather than through printed text.

## References

- `_prd.md` → Goal 3, User Story 3, Core Feature 4
- `_techspec.md` → Build Order 5
- ADR-0135 makes an absent diagnostic a reported state rather than an empty
  message, which is the principle this Task applies to a raw errno

## Verification
- `grep -q "TestTaskWorktreeCreationFailureNamesItsContext" internal/worktree/worktree_test.go || exit 1; go test -count=1 ./internal/worktree ./internal/daemon`

## Result

- Task Worktree creation options now carry the configured concurrency from the
  Daemon. The sequential convenience path records concurrency `1`.
- Every Task Worktree creation failure now names its Run, Task, and configured
  concurrency, states the next action, and wraps the original error with `%w`.
- `TestTaskWorktreeCreationFailureNamesItsContext` forces a real filesystem
  failure, checks all three context values and the next-action label, and uses
  `errors.As` plus `errors.Is` to prove the filesystem error remains in the
  chain.
- Focused evidence: `rtk proxy env
  GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run
  '^(TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot|TestTaskWorktreeCreationFailureNamesItsContext|TestTaskWorktreesIntegrateFirstByFastForwardThenCherryPick)$'
  ./internal/worktree` passed.
- Focused evidence: `rtk proxy env
  GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run
  '^(TestTaskCycleCreatesTaskWorktreesWithBootstrapBeforeAgentWork|TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks)$'
  ./internal/daemon` passed; the first test asserts that the Task Plan's
  concurrency reaches each Task Worktree creation request.
- Focused evidence: `rtk git diff --check` passed.
- The declared Verification command was not run; the Daemon owns that gate.
