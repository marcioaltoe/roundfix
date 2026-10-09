---
type: fix
status: open
created: 2026-10-09
spec: null
---

# Each worktree keeps its own multi-gigabyte Go build cache

## Problem

The Makefile sets `GOCACHE ?= $(CURDIR)/.gocache`, so every Delivery Queue
item worktree and every Run worktree builds its own Go cache. On 2026-10-09,
during Spec 0254's QA, the disk of this 228 GB Mac went from 14 GB free to
2 GB free three times. One item worktree's `.gocache` held 8.0 GB.
`go clean -cache` does not reach these caches, and removing a terminal Run
worktree only frees its own copy. The QA gate also builds scratch clones with
separate caches to compare timings.

## Direction

- Share one build cache between worktrees of the same repository, for example
  under Roundfix Home, keyed so that concurrent Runs stay safe. Alternatively,
  remove a worktree's `.gocache` when its Run settles.
- Report cache bytes in `roundfix doctor` storage, and reclaim them in
  `roundfix gc`.

## Sources

Operator log, 2026-10-09, entries 251, 254 and 256: disk at 5.3 GB, then 2 GB
twice; the cleanup freed the item `.gocache` (8.0 GB), the default cache
(9.9 GB) and terminal Run worktrees.
