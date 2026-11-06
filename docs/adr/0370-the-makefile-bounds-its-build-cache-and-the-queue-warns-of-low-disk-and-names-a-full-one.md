---
status: accepted
created_at: 2026-11-06T12:00:00Z
updated_at: 2026-11-06T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The Makefile bounds its build cache, and the queue warns of low disk and names a full one

On 2026-11-05 the Mac mini's system volume filled twice. At 15:58 BRT
oraculum's Postgres stopped in the middle of a Run, and at 19:07 BRT the
Delivery Queue parked three items as `delivery-error` because `make build`
failed with `no space left on device`. Each time the largest temporary
directory was this repository's Makefile cache,
`$TMPDIR/roundfix-go-build-<uid>`, at 25 GB, and on 2026-11-06 it was 26 GB
again. The Shared Build Cache trim of ADR-0327 and ADR-0334 ran, but only at
Run starts and gates, with a 20 GiB limit and four idle hours, so a busy day
kept every entry young enough to stay. On 2026-11-06 the maintainer chose
"Spec: limite + aviso (Recommended)" (option text: "O Makefile limpa ou poda
o cache acima de um tamanho antes de compilar. O doctor e a fila avisam
quando o disco do cache fica com pouco espaço, e um park por disco cheio diz
isso em vez de 'delivery-error'.").

The go command has no size limit for its cache. It deletes an entry only after
five days without use, and it refreshes a used entry's modification time at
most once an hour. `go clean -cache` deletes every entry, and a concurrent
build can then fail: golang/go issue 31948, "cmd/go: concurrent build and
cache clean is unsafe", shows `can't open import: ... -d: no such file or
directory`.

## Decision

**This repository's Makefile holds the cache it chooses to a Build Cache
Bound, removing only idle entries; `roundfix doctor` and the Delivery Queue
warn below a Free-Space Floor; and a park caused by a full filesystem is a
Disk-Full Park.**

- The Build Cache Bound. Before each make target that builds or tests, the
  top-level make runs `go run ./cmd/gocache-bound` over the cache the
  Makefile chose (the shared `roundfix-go-build-<uid>` or the `.gocache`
  fallback). Above `GOCACHE_BOUND_MB` (default 10240 MiB) it removes entries
  idle for `GOCACHE_BOUND_IDLE` (default and minimum two hours), oldest first,
  until the cache is within the bound, with the rules ADR-0327 gave the GC
  Command, under one thirty-second budget. A sub-make does not repeat it. An
  explicit `GOCACHE` is bounded only when `GOCACHE_BOUND_MB` is also set, and
  `GOCACHE_BOUND_MB=0` turns it off. A refusal, an error, the end of the
  budget or a failed `go run` prints a warning and never fails the build.
  `.roundfixrc.yml` declares the same cache with the same values, so the GC
  Command, Run starts and gates trim it alike.
- The Free-Space Floor is 10,240 MiB. Doctor's `disk-space:` line inspects the
  temporary directory beside the worktree location and each declared Shared
  Build Cache. Before an item's Implement Run and before its repository gate,
  the Delivery Queue inspects the same locations and, below the floor,
  records a `low-disk` segment in the item's warning and writes one line to
  the Queue Owner's log; at or above the floor the segment goes. It never
  fails, parks or delays an item.
- A Disk-Full Park. When the item binary build or the repository gate fails
  and its error is `ENOSPC` or its retained log says `no space left on
  device`, the item parks as `disk-full: <path>: <N> MiB free` in the
  environment Park Class. `<path>` is the last absolute path the log names
  before that text, or else the inspected location with the least free
  space, and `<N>` is its free space at the park. The next action names what
  to free, and a Delivery Retry from `disk-full` is not counted against the
  retry limit.

The Build Cache Bound is this repository's own: adopters keep their own build
caches. The Free-Space Floor and the Disk-Full Park reach every repository.

## Considered Options

- Clearing the cache above the bound (`go clean -cache`). Rejected: three Runs
  share it, and a clean during a build can fail that build (golang/go issue
  31948), which ADR-0327 already refused.
- Removing young entries when nothing idle is left. Rejected for the same
  reason: an entry used in the last hour may still be read by a running build;
  the bound can stay above its size until entries age, and the floor warns
  meanwhile.
- `-trimpath`, so worktrees share entries. Rejected: thirteen test files
  locate their fixtures through `runtime.Caller`, which reports trimmed paths
  under `-trimpath`.
- A shell trim with `find -mmin`. Rejected: it cannot remove oldest first
  until within the bound, and the Go trim already proves its rules.
- Refusing to start a Run below the floor. Rejected again, as ADR-0334
  deferred it: the failures happen while a Run writes, and a warning that
  leaves the decision to the operator stops no work.
- Moving the cache to another volume. Left to each machine: `GOCACHE` and
  `TMPDIR` choose it, and the user guide names the option.

## Consequences

- A Run that shares the cache may recompile an entry removed after two idle
  hours.
- Each top-level make spends about a second compiling and running the bound.
- A Run whose Task Verification fills the disk still parks `run-unresolved`;
  `roundfix runs causes` names its `disk-full` signature.
