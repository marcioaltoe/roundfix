---
task: task_03
spec: 0161-deliver-ready-for-real-repositories
status: completed
type: backend
complexity: medium
---

# Task 03: Clean parks, crash-safe archive, stale owners

## Overview

A non-exact archive park leaves the checkout dirty so every later item refuses to start; a crash after the archive commit resumes as `review-stale`; a reused owner PID makes both `resume` and `stop` refuse.

## Requirements

1. MUST, on park, discard uncommitted changes the item made and return the checkout to the branch it started from.
2. MUST accept on resume a HEAD whose parent is the reviewed head and whose change is exactly the Spec move into the archive root.
3. MUST release, on `resume`, an owner whose live process identity does not match the recorded identity.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, the next item starts after a non-exact archive park.
- [ ] A crash after the archive commit resumes past archiving.
- [ ] A reused PID with another identity is released.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/delivery/engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAParkLeavesACleanCheckout|TestResumeAcceptsTheArchiveCommit|TestResumeReleasesAStaleOwner)$" ./internal/cli ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAParkLeavesACleanCheckout TestResumeAcceptsTheArchiveCommit TestResumeReleasesAStaleOwner; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Park and resume

## Result

Implemented clean parking and crash-safe resume without changing the Task
status or running the Daemon-owned Verification command:

- Every park now asks the item workspace to discard tracked and untracked
  changes, restore the branch from which that item entered its delivery branch,
  and prove the restored checkout is clean before persisting `parked`.
- Archive reconciliation accepts an already-committed child of the reviewed
  head only when its complete Git diff stays inside the active and archived
  Spec paths, the source and destination tree identities match, the source is
  absent at the child, and the destination was absent at the parent. Dirty or
  broader changes remain stale.
- `deliver resume` now proves the identity behind a live recorded PID. A proven
  mismatch releases the stale owner record; an unreadable identity still fails
  closed, and a matching live owner still blocks the second owner.

Pre-change evidence:

- `TestAParkLeavesACleanCheckout` failed because only the first item started;
  the dirty checkout prevented the second item from entering its branch.
- `TestResumeAcceptsTheArchiveCommit` failed because the committed exact move
  returned `ExactSpecMove: false`.
- `TestResumeReleasesAStaleOwner` failed because PID liveness blocked resume
  before the recorded identity was checked.

Focused-check evidence:

- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^TestAParkLeavesACleanCheckout$' ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^TestResumeAcceptsTheArchiveCommit$' ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^TestResumeReleasesAStaleOwner$' ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 ./internal/cli ./internal/delivery` — passed with process-table permission; the sandboxed attempt was blocked only in existing force-stop integration tests that inspect the macOS process table.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go vet ./internal/cli ./internal/delivery` — passed.
- `rtk git diff --check` — passed.

Acceptance evidence:

1. `TestAParkLeavesACleanCheckout` uses real Git, dirties one tracked and one
   untracked path during a non-exact archive, observes the next queued Spec
   start, and finishes on a clean `main` checkout.
2. `TestResumeAcceptsTheArchiveCommit` creates a real archive commit, resumes
   an item persisted at `archiving`, and observes the archive head appended
   without `review-stale`; its negative case rejects the same move when the
   commit also changes an unrelated path.
3. `TestResumeReleasesAStaleOwner` observes a live PID with a mismatched process
   identity released before the new owner starts; its negative case proves an
   unreadable identity retains the owner and refuses resume.
