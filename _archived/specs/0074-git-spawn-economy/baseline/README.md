# Git spawn census baseline

This baseline measures every `git` executable resolved through `PATH` during
one fresh `go test ./... -count=1 -parallel 16` run. It also records a separate
fresh run's wall, user, and system CPU time. Task 06 must use the same commands
when it publishes the after-measurement.

## Census procedure

Run from the repository root on macOS:

```sh
repo_root=$(rtk proxy pwd -P)
work_dir=$(rtk mktemp -d "${TMPDIR:-/tmp}/roundfix-git-census.XXXXXX")
rtk mkdir "$work_dir/bin" "$work_dir/gocache"
rtk ln -s "$repo_root/docs/specs/0074-git-spawn-economy/baseline/git-census-shim" "$work_dir/bin/git"
real_git=$(rtk proxy which git)
invocation_log="$work_dir/git-invocations.tsv"
: >"$invocation_log"
rtk proxy env \
  GIT_CENSUS_REAL="$real_git" \
  GIT_CENSUS_LOG="$invocation_log" \
  GOCACHE="$work_dir/gocache" \
  PATH="$work_dir/bin:$PATH" \
  go test ./... -count=1 -parallel 16
rtk proxy wc -l "$invocation_log"
rtk proxy cut -f1 "$invocation_log" | rtk proxy env LC_ALL=C sort | rtk proxy uniq -c | rtk proxy sort -nr
rtk proxy awk -f "$repo_root/docs/specs/0074-git-spawn-economy/baseline/attribution.awk" "$invocation_log"
```

The empty task-local `GOCACHE` makes compilation fresh, and `-count=1`
disables Go's test-result cache. Keep `work_dir` until the three parsing
commands have run; the raw log is temporary and is not a stable artifact.

## Timing procedure

Run separately from the census so the shim's own shell process and log write
do not affect the attribution:

```sh
timing_dir=$(rtk mktemp -d "${TMPDIR:-/tmp}/roundfix-git-timing.XXXXXX")
rtk mkdir "$timing_dir/gocache"
rtk proxy env GOCACHE="$timing_dir/gocache" /usr/bin/time -p go test ./... -count=1 -parallel 16
```

## Measured results

The production and test code measured by both runs was at HEAD
`dbdad8ac1b8a2335ab88c65a0a47f50d86ef6c4e`; the only worktree changes were
this Task's documentation artifacts. The measurement environment was macOS
26.5.2 (`25F84`) on arm64, Go 1.26.5, and Git 2.55.0, recorded with:

```sh
rtk git -c core.fsmonitor=false rev-parse HEAD
rtk proxy sw_vers
rtk proxy go version
rtk proxy git --version
```

The census run exited 0 and recorded **12,099 total Git spawns**. The committed
parsing commands produced this per-subcommand table:

```text
3859 rev-parse
 998 commit
 997 add
 985 ls-tree
 974 cat-file
 843 rev-list
 646 init
 452 status
 362 worktree
 304 branch
 292 for-each-ref
 196 checkout
 187 check-ref-format
 138 remote
 113 merge-base
 107 show-ref
  95 tag
  95 config
  86 archive
  80 symbolic-ref
  71 show
  52 merge
  51 log
  36 diff
  35 diff-tree
  16 fetch
   8 push
   7 cherry-pick
   6 clone
   5 reflog
   1 stash
   1 ls-remote
   1 ls-files
```

The attribution parser reported:

```text
7405 production-read-shaped
2736 fixture-setup-shaped
1958 ambiguous-or-other
0 malformed-records
```

The separate timing run exited 0. `/usr/bin/time -p` reported:

```text
real 78.45
user 127.76
sys 268.46
```

## Attribution method and limits

`attribution.awk` separates command shapes into three disjoint buckets:

- `production-read-shaped`: `rev-parse`, `ls-tree`, `cat-file`, `rev-list`,
  `status`, and `for-each-ref`, the read operations identified in the PRD.
- `fixture-setup-shaped`: `init`, `config`, `add`, and `commit`, which primarily
  construct temporary repositories.
- `ambiguous-or-other`: every remaining subcommand.

These are shape-based proxies, not exact caller attribution. Fixture helpers
also query repositories, and tests exercise production mutation paths. The Go
test binary is the parent for both kinds of call, and temporary-repository
working directories are used by both, so neither parent process nor working
directory can identify the call site. The shim also cannot observe an
absolute `/path/to/git` that bypasses `PATH`. The recorded split therefore
describes recognizable command shapes and states its uncertainty; it must not
be read as an exact production-versus-fixture count.
