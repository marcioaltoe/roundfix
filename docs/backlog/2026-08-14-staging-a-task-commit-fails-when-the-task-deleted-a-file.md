---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Staging a Task commit fails when the Task deleted a file

## Symptom

An `implement` Run ended `Failed` because the Daemon could not commit a Task that
had finished its work correctly. The staging command names every path the Task
touched, including the ones it **deleted**, and `git add -f -- <path>` does not
match a path absent from both the working tree and the index.

```text
roundfix: implement failed after Run start: git add -f -- ... \
  packages/frontend/src/systems/catalog-archival/lib/remaining-window.ts ... \
  failed: exit status 128: \
  fatal: pathspec 'packages/frontend/.../__tests__/remaining-window.test.ts' did not match any files
```

The Task replaced a countdown with a pending age, so deleting
`remaining-window.ts` and its suite was exactly what it was supposed to do. The
work was correct and verified; only the staging call failed.

This is the same family as a hook refusal killing a Run whose work was already
verified: the Daemon is the authority that proved the work, and it loses it at
the commit step.

## Where

The Daemon's Task staging call, where it builds a `git add -f` pathspec from the
Task's changed-file set.

## Expected

Staging covers deletions. `git add -A -- <paths>` records a removal where
`git add -f -- <path>` refuses it, and the Task's changed-file set already
distinguishes what was deleted from what was written.

Worth settling in the same work: whether a staging failure should be recoverable
at all rather than terminal, since the work it loses has already passed the
authoritative Verification.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-12-git-add-f-falha-quando-a-task-apaga-arquivo.md` in the
Secondbrain. Related:
`docs/backlog/2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md`
records the same shape reached through a commit hook instead of a pathspec.
