---
status: accepted
created_at: 2026-11-12T21:00:00Z
updated_at: 2026-11-12T21:00:00Z
deprecated_at: null
superseded_by: null
---

# The Makefile follows the Go build cache set with go env, and falls back small

ADR-0370 left the choice of volume to each machine: "`GOCACHE` and `TMPDIR`
choose it". On 2026-11-11 the Mac mini, whose internal disk must keep 80 GB
free, set its user Go build cache on an external disk with
`go env -w GOCACHE=/Volumes/Brain-SSD/cache/go-build`. That setting lives in
the go command's environment file, not in the environment, so this
repository's Makefile kept choosing `$TMPDIR/roundfix-go-build-<uid>` on the
internal disk, which held 10 GB on 2026-11-12. The MacBook shares the shell
startup files with the mini and has no such disk. On 2026-11-12 the
maintainer decided that make uses the external cache with a 10 GB limit and,
when that disk is not mounted, "Autorizo; sem SSD, TMPDIR com 2 GB" (option
text: "O build segue no disco interno com limite pequeno e avisa.").

The go command reports a setting made with `go env -w` through
`go env -changed`, which "prints only those settings whose effective value
differs from the default value that would be obtained in an empty environment
with no prior uses of the -w flag". A missing or unreadable environment file
is ignored, and a relative `GOCACHE` value reads as `off`.

## Decision

**When the environment does not define `GOCACHE`, this repository's
Makefile uses the cache set with `go env -w GOCACHE` if it is an existing,
writable directory, under the Build Cache Bound of 10,240 MiB; when that
cache is set but unusable, it warns once and falls back to the cache it
chose before, under a default bound of 2,048 MiB; when nothing is set, it
chooses exactly as ADR-0370 decided.**

- Set. `go env -changed GOCACHE` prints a value. make reads the path with
  `go env GOCACHE`.
- Usable. The path is absolute, is an existing directory, and make can
  create and remove a probe file in it. make never creates the directory or
  any parent: on a disk that is not mounted, a missing path, or an empty
  mount-point directory left under `/Volumes`, must not become a cache on the
  internal disk.
- Used. `GOCACHE` is that path, and `GOCACHE_BOUND_MB` defaults to 10240, as
  for the cache make chose before: the cache is one make reads and chooses,
  not one the caller passed.
- Fallback. `GOCACHE` is `$TMPDIR/roundfix-go-build-<uid>`, or
  `$(CURDIR)/.gocache` when that cannot be created and written, and
  `GOCACHE_BOUND_MB` defaults to 2048. make prints one line to standard error
  and continues:
  `gocache: warning: GOCACHE set with go env (<set path>) is not a writable directory; using <fallback> with a 2048 MiB default bound`.
- Unchanged. A `GOCACHE` in the environment or on the command line still wins
  untouched and unbounded unless `GOCACHE_BOUND_MB` is set; `GOCACHE_BOUND_MB`
  and `GOCACHE_BOUND_IDLE` still override the defaults; a sub-make inherits the
  exported `GOCACHE` and neither queries nor warns again; with nothing set,
  the cache and the bound are byte-for-byte those of ADR-0370.

The repository never names a machine's path. It reads the go command's own
setting, so every Go command a user runs directly and every make target
share one cache.

## Considered Options

- Testing for `/Volumes/Brain-SSD/cache` in the Makefile, as the request
  suggested. Rejected: it names one machine's disk in a repository every
  machine builds, and a second machine or a renamed disk needs another edit.
- Failing make when the set cache is unusable. Rejected by the maintainer's
  answer: the build continues on the internal disk under a small bound and
  warns.
- Exporting `GOCACHE` from the shell startup files. Rejected: the MacBook
  shares them and has no such disk, and an exported `GOCACHE` is a caller's
  cache that make leaves unbounded.
- Creating the set directory when it is missing (`mkdir -p`), as make does for
  its temporary cache. Rejected: on an unmounted disk it creates the path
  under `/Volumes` on the internal disk, the silent fallback the request
  forbids.
- Reading the environment file directly. Rejected: its location depends on
  `GOENV` and on the operating system's configuration directory, and its
  format belongs to the go command; `go env -changed` is that command's own
  answer.
- Treating the set cache as an explicit one, unbounded unless
  `GOCACHE_BOUND_MB` is set. Rejected by the maintainer's answer, which bounds
  it to 10 GB.
- A placeholder in `worktree.shared_caches` that resolves the set cache, so
  the GC Command, Run starts and gates trim it too. Deferred: it changes the
  released binary for a cache the Makefile already bounds before every
  top-level build or test, including every Verification the Daemon runs.

## Consequences

- On the mini, every make target builds in the external cache and holds it to
  10,240 MiB. The cache is shared by every Go project of the user, so the
  bound counts their entries and may remove theirs, idle at least two hours,
  oldest first, never a young one and never the whole cache (ADR-0327).
- `.roundfixrc.yml` still declares `{tmpdir}/roundfix-go-build-{uid}` with
  10,240 MiB for the GC Command, Run starts and gates; on a machine with a
  usable set cache that declaration covers the fallback, and the set cache is
  trimmed only by make. Doctor's free-space line does not inspect it.
- A sandboxed Agent whose sandbox cannot write the set cache gets the fallback
  and the warning on each top-level make.
- Each top-level make runs `go env` two to four times, a few milliseconds
  each.
- A caller that exports `GOCACHE` keeps an unbounded cache, as before. The
  613 MB in `~/Library/Caches/go-build` on 2026-11-12 came from such a caller,
  an operator script outside this repository.
