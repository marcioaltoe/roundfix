---
status: done
created_at: 2026-08-11
updated_at: 2026-09-08
spec: 0124-verification-capacity-and-measured-economics
---

# A git worktree that fails only under load

`TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork` failed
twice on 2026-08-11 inside a full `make verify`, and passed on every focused
rerun. The failure is always the same:

```
fatal: failed to read .git/worktrees/<run>.task_01/commondir: Result too large
```

`Result too large` is `ERANGE`. Git raises it while creating the *second*
concurrent Task Worktree, reading the `commondir` file the *first* one wrote.
The Run then ends `Unresolved` and the test sees exit 1 where it expects Clean.

Two independent observations put this in the same family rather than making it
a curiosity:

- Spec 0080's Task 11 hit a sibling failure in the same session: this
  repository's fsmonitor IPC diagnostic entered `git diff-tree` combined output,
  and the mechanical detector read that diagnostic text as a changed path. The
  same check passed with `core.fsmonitor=false`.
- `internal/gittest` already sets `core.fsmonitor=false` and
  `GIT_OPTIONAL_LOCKS=0` for isolated test repositories, added on 2026-08-10 as
  hygiene rather than in response to a measured failure. Production worktree
  creation does not.

## Why it is worth fixing

It is invisible until it is expensive. Under `make verify` the suite runs
packages in parallel and the failure appears; run the same test alone and it
passes, so the natural reaction is to call it flaky and rerun. That is what
happened here — twice — and each time the real cost was a full gate cycle spent
deciding whether a Spec had regressed.

It also lands on the least legible surface. The message names `commondir` and an
errno, not the Task, the Run, or the concurrency that produced it, so a reader
has no path from the symptom to "two Task Worktrees were being created at once
while a filesystem monitor was watching the same directory".

## Shape

Non-binding.

**Find out whether the monitor is the cause before removing it.** The fsmonitor
correlation is suggestive, not proved: nothing here ran the concurrent-worktree
path with the monitor disabled and measured the failure rate. That experiment is
cheap — repeat the failing test under `make verify` load with and without
`core.fsmonitor=false` in the production worktree path — and it decides whether
the fix is a config default or something else entirely.

**Whatever the cause, the error deserves a Roundfix sentence.** A Task Worktree
that cannot be created is a Run-level event with a Run, a Task, and a
concurrency level; the git text is evidence for that sentence, not a substitute
for it. Today `implement failed after Run start: create Task Worktree: ...`
passes the raw fatal through and stops.

## Addendum — 2026-09-08 — Current triage

Spec `0097-a-wave-that-cannot-collide` supplied the contextual error in
`internal/worktree/worktree.go`, naming Run, Task, and configured concurrency.
The current production Git runner also passes `core.fsmonitor=false`.

Those source changes address the diagnostic and the suspected interference,
but this triage did not reproduce the original ERANGE or run the controlled
experiment needed to prove its cause. Root cause therefore remains unknown.
Keep the Finding pending and route the remaining characterization to
provisional P5, Verification capacity and measured economics. Do not treat an
isolated passing rerun as proof that the load failure is fixed.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0124-verification-capacity-and-measured-economics](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
