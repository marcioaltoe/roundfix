---
status: accepted
created_at: 2026-10-24T12:00:00Z
updated_at: 2026-10-24T12:00:00Z
deprecated_at: null
superseded_by: null
---

# GC trims a declared shared build cache

This repository's Makefile gives every worktree one Go build cache under the
temporary directory. Each Run worktree compiles the module under its own
path, so the cache grows by a full build per worktree: 43.7 GB filled the
disk on 2026-10-22, and on 2026-10-24 it held 33 GB, none of it older than
fifteen hours. The go command deletes an entry only after five days without
use. The GC Command reclaims Worktree Build Caches inside ended worktrees but
leaves this shared cache alone. On 2026-10-24 the maintainer chose "roundfix
gc recolhe" (option text: "Vira uma mudança do produto: o comando gc passa a
recolher esse cache.").

The go command refreshes the modification time of a cache entry it uses at
most once an hour (`mtimeInterval` in `cmd/go/internal/cache`), and treats a
missing entry as a miss that it rebuilds. Its own trim deletes entries unused
for five days plus that hour.

## Decision

**A Shared Build Cache is declared, and GC trims only its idle entries.**
`worktree.shared_caches` in User Config or Project Config lists caches by
`path`, `kind`, `max_size_mb` and `idle`. The path is absolute or starts
with `{tmpdir}` or `{home}`, and `{uid}` stands for the user id, so a
repository can name the cache its build tooling uses without Roundfix knowing
that tooling. The only kind is `go-build`. `idle` is at least two hours,
twice the go command's refresh interval, so an entry any build used in the
last hour is never idle.

`roundfix gc --dry-run` reports each cache's size and the entries a trim
would remove; `roundfix gc` removes them. A trim acts only when the cache
exceeds `max_size_mb`, removes idle entries oldest first until the cache is
within its limit, and rereads each entry's modification time just before
removing it. It refuses a path that is a symbolic link or whose directory
lacks the go command's cache README, and it never removes a file that is not
a cache entry: a name of 64 hexadecimal digits and `-a` or `-d`, inside a
two-digit shard directory.

## Considered Options

- A Makefile target that trims the cache was rejected by the maintainer's
  answer: the reclamation belongs to the product, where adopters get it too.
- Emptying the cache above the limit (`go clean -cache`) was rejected: three
  concurrent Runs share it, and an empty cache makes each of them recompile.
- Discovering the cache from `go env GOCACHE` was rejected: the Makefile sets
  `GOCACHE` itself, and GC would see Go's default cache instead.

## Consequences

- This repository declares its Makefile cache in `.roundfixrc.yml`; a
  repository with no declaration sees no change in GC's output.
- A worktree whose entries were idle for the declared time recompiles them.
- The Worktree Build Cache rules and the Item Binary sweep stay as they are.
