---
status: completed
type: backend
---

# Task: Bootstrap serialized across sibling Task Worktrees

Sibling worktrees bootstrap against one shared Git directory and one shared
package cache, and collide on the lock — reporting failure after having done
every byte of the work. Measured in a repository this Spec did not build: a
`prepare` step writing `git config` failed with `could not lock config file`.

## Work

- Serialize the bootstrap step across sibling Task Worktrees. Only that step:
  Agent turns and Verification stay parallel, because the collision is in the
  shared Git directory and the shared cache, not in the work.
- A bootstrap that fails after completing its work reports that state
  distinctly from one that failed before starting. Telling a maintainer nothing
  happened when everything did is the part that cost the Run.
- Change no default concurrency and recommend none; the Spec's Non-Goals forbid
  both.
- Cover bootstrap at capacity above one across repeated attempts with no lock
  collision, and the two failure classifications kept apart.

## References

- `_prd.md` → Goal 2, User Story 2, Core Features 2 and 3
- `_techspec.md` → Build Order 4
- ADR-0056 separates Task Capacity from Verification Capacity, and this Task
  changes neither

## Verification
- `grep -q "TestBootstrapSerializesAcrossSiblings" internal/worktree/worktree_test.go || exit 1; grep -q "TestBootstrapFailureAfterWorkIsClassifiedApart" internal/worktree/worktree_test.go || exit 1; go test -count=1 ./internal/worktree`

## Result

Implementation-ready behavior:

- Worktree Bootstrap commands now acquire one process-wide,
  cancellation-aware permit. The permit covers only command execution;
  worktree creation and copying, Agent work, and Verification retain their
  existing concurrency.
- `BootstrapError` now carries `BootstrapFailureBeforeStart` or
  `BootstrapFailureAfterStart`. Its stable error prefix remains intact, while
  the suffix states either that bootstrap work did not start or that bootstrap
  work may have been applied.
- No concurrency default, capacity setting, or concurrency recommendation was
  changed.

Acceptance evidence:

- `TestBootstrapSerializesAcrossSiblings` launches four sibling bootstrap
  commands together for six rounds against a real exclusive filesystem lock.
  A FIFO releases each admitted command without sleeps; the race-enabled
  focused check observed all 24 commands exit without a lock collision.
- `TestBootstrapFailureAfterWorkIsClassifiedApart` distinguishes validation
  failure before process start from a command that writes
  `bootstrap.complete` and then exits non-zero. It checks both the typed stage
  and the maintainer-facing message.
- Existing non-zero-exit and timeout tests now assert the post-start stage,
  while focused CLI and Daemon checks retain the prior failure prefix and Task
  isolation behavior.

Focused checks run after the final implementation edit:

- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -race -count=1 -timeout=60s -run 'Bootstrap' ./internal/worktree` — exit 0.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -timeout=30s -run '^TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks$' ./internal/daemon` — exit 0.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -timeout=60s -run '^TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork$' ./internal/cli` — exit 0.
- `git diff --check` — exit 0 after the code and Result formatting edits.

The Daemon-owned `## Verification` command was not run during this Agent turn.
