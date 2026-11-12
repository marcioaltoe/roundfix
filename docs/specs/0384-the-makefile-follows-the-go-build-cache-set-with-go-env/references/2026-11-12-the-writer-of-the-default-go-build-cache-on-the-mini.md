# The writer of the default Go build cache on the Mac mini

Measured by the author of Spec 0384 on the Mac mini on 2026-11-12, between
17:15 and 18:30 BRT, read only. Nothing was written, deleted or moved in any
cache directory. Times are BRT (UTC−3).

## The machine

- `go version`: go1.26.5 darwin/arm64, `/usr/local/go/bin/go`.
- `go env GOENV`: `~/Library/Application Support/go/env`, created and last
  modified 2026-11-11 20:51:04, holding one line,
  `GOCACHE=/Volumes/Brain-SSD/cache/go-build`.
- `go env GOCACHE`: `/Volumes/Brain-SSD/cache/go-build` (11 GB, 35,570 files,
  shared by every Go project of the user).
- `$TMPDIR/roundfix-go-build-501`: 10 GB, the cache this repository's
  `Makefile` chooses when `GOCACHE` is not in the environment.
- `~/Library/Caches/go-build`: 613 MB; the directory and its `README` were
  created 2026-11-11 22:00:19, `trim.txt` at 22:00:28.

## When the default cache was written

Entries by modification time, grouped by ten seconds:

| window | entries |
| --- | --- |
| 2026-11-11 22:00:1x–22:00:2x | 43 |
| 2026-11-11 22:02:4x–22:02:5x | 80 |
| 2026-11-11 23:50:0x | 41 |
| 2026-11-11 23:52:2x–23:52:3x | 72 |
| 2026-11-12 02:57:4x | 1,470 |
| 2026-11-12 02:59:5x–03:00:0x | 1,532 |

No entry is older than 22:00:19 on 2026-11-11 or newer than 03:00 on
2026-11-12. Every source path the entries cite is in the main checkout
`~/dev/roundfix/internal/...` or in `~/go/pkg/mod`; none is in a Run Worktree
(those live under `/Volumes/Brain-SSD/cache/.roundfix/worktrees`). The
`~/.vault/home/.codex/.../codex-path/rg` path the source message saw comes
from four committed fixtures under
`internal/baseline/testdata/parity-corpus/v1/fixtures/`, which a compiled
test package embeds; it is file content, not the environment of the writer.

## What ran at those times

Each burst is a pair around three short Codex sessions. The Codex state
database (`~/.codex/state_5.sqlite`, table `threads`) holds, with working
directory `~/dev/roundfix`, source `acpx` and no user message, three sessions
of eight seconds each at 22:00:38, 22:01:29 and 22:01:53; at 23:50:17,
23:51:09 and 23:51:33; and at 02:57:52, 02:58:43 and 02:59:06. They are the
profile proofs `roundfix doctor` opens. Each burst of cache entries starts
just before the first proof and resumes right after the third.

The Supervisor's Claude Code session transcript of that night records, at
the same minutes, after each Delivery Queue merge (0378 at 21:57, 0379 at
23:46, 0380 at 02:55), the call
`~/.roundfix-operator/release-checks.sh <tag> <reason>` (01:00:19, 02:50:05
and 05:57:40 UTC). That operator script, outside this repository, begins:

```sh
cd /Users/marcio/dev/roundfix || exit 1
export GOCACHE=$HOME/Library/Caches/go-build GIT_CONFIG_COUNT=1 ...
make build >/dev/null 2>&1 || { echo build failed; exit 1; }
env -u NODE_OPTIONS bin/roundfix doctor 2>&1 | grep -E '^skills'
go test -count=1 ./skills -run ...
go test -count=1 ./internal/baseline -run ...
...
env -u NODE_OPTIONS go run -buildvcs=false ./cmd/roundfix baseline update ...
```

`make build` fills the first burst, `roundfix doctor` opens the three
profile proofs, and the `go test` and `go run` lines fill the second burst.
The Codex shell snapshot of the proof session of 22:01:29 records
`export GOCACHE=/Users/marcio/Library/Caches/go-build`, inherited from the
script. The directory appeared at 22:00:19 although the script also ran at
18:04 and 20:27, which is consistent with the directory having been removed
when the Go environment file was written at 20:51.

## Why it grew without a bound

At `d79b30b4`, with `GOCACHE` in the environment, the `Makefile` exports it
unchanged and sets `GOCACHE_BOUND_MB` to 0 (ADR-0370: "An explicit `GOCACHE`
is bounded only when `GOCACHE_BOUND_MB` is also set"):

```text
$ GOCACHE=<dir>/explicit TMPDIR=<dir>/tmp/ make -f Makefile -f print.mk print-gocache
GOCACHE=<dir>/explicit BOUND=0
$ env -u GOCACHE GOENV=<dir>/env-ssd TMPDIR=<dir>/tmp/ make -f Makefile -f print.mk print-gocache
GOCACHE=<dir>/tmp/roundfix-go-build-501 BOUND=10240
```

The second line also shows the cause of the 10 GB: a cache set with
`go env -w GOCACHE` is not in the environment, so `make` ignores it.

## The candidates the brief named, ruled out

- A Run Agent's sandbox. The Codex sandbox of a Run Agent writes only to its
  Run Worktree, `$TMPDIR`, `/tmp` and `~/.agent-browser` (thread
  `sandbox_policy`), and a 2026-10-12 Secondbrain digest records that
  `codex sandbox -P :workspace` denied writing `~/Library/Caches/go-build`.
  The Run Agents of that night ran `go test` and `make` escalated, outside
  the sandbox, where the Go environment file is read and the cache set there
  is used. No Run Agent session works in the main checkout.
- The Daemon's Verification environment. Outside the Linux Lane, which sets
  `GOCACHE=/tmp/go-build` inside its container, no production code sets
  `GOCACHE` or `GOENV`; `worktree.shared_caches` only trims declared caches.
- The operator's own commands around the 0380 retry. The transcript shows
  them at 02:30–02:33 (`gofmt`, `go vet`, `roundfix deliver retry`), all in
  the item's worktree, and no default-cache entry is dated then; the next
  `go` command is the `go build` of 02:57:32 in the main checkout, outside
  the script, which reads the Go environment file and uses the cache set
  there.

## Conclusion

The writer is the operator's `~/.roundfix-operator/release-checks.sh`, which
exports `GOCACHE=$HOME/Library/Caches/go-build` before `make`, `go test` and
`go run` in the main checkout. It is outside this repository; removing that
export makes those commands use the cache the `Makefile` chooses (through
`make`) or the cache set with `go env` (direct `go` commands).
