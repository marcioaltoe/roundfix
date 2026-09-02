---
status: pending
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
