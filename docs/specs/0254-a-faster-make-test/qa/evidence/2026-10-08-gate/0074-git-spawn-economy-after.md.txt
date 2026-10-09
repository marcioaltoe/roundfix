# Git spawn census — before and after

**Verdict: the fresh suite missed the under-60-second target.** The after run
took **83.38 seconds**, 23.38 seconds over the target and 4.93 seconds slower
than the committed 78.45-second baseline. Git spawns fell by 1,571 (13.0%),
but that reduction did not lower the measured wall clock.

Both full-suite columns use the exact census and timing procedures in
[`README.md`](README.md): fresh task-local `GOCACHE` directories,
`go test ./... -count=1 -parallel 16`, and a separate timing run without the
PATH shim. The before run measured
`dbdad8ac1b8a2335ab88c65a0a47f50d86ef6c4e`; the after run measured
`a9d0097c590412dafd173fbcb4deaf1923bcae3a`. Both ran on macOS 26.5.2
(`25F84`) arm64 with Go 1.26.5 and Git 2.55.0 on the same host class.

## Fresh-suite timing

Both timing runs exited 0.

| Measurement (seconds) | Before | After | Delta |
| --- | ---: | ---: | ---: |
| Wall (`real`) | 78.45 | 83.38 | +4.93 (+6.3%) |
| User CPU (`user`) | 127.76 | 121.55 | -6.21 (-4.9%) |
| System CPU (`sys`) | 268.46 | 262.24 | -6.22 (-2.3%) |

The after run used less user and system CPU but more elapsed time. Its system
CPU remained 2.16 times its user CPU, so process and filesystem work still
dominate compute. The procedure does not isolate scheduler noise or attribute
CPU to individual call sites, so the 4.93-second wall regression cannot be
assigned to one implementation change from these totals.

The after run's longest Go-reported package times show where the remaining
critical path lives:

| Package | After elapsed |
| --- | ---: |
| `internal/cli` | 72.913s |
| `internal/baseline` | 67.505s |
| `internal/worktree` | 34.911s |
| `skills` | 21.761s |
| `internal/spec` | 20.076s |

Package tests run concurrently, so these values are not additive. They do
show that `internal/cli` alone remained 12.913 seconds beyond the target and
that `internal/baseline` also remained beyond it. The Spec reduced repository
read spawns, but it did not remove those two package floors.

## Git spawn census

The successful after census exited 0 and recorded **10,528** Git spawns with
zero malformed records, down from **12,099**.

| Git subcommand | Before | After | Delta |
| --- | ---: | ---: | ---: |
| `rev-parse` | 3,859 | 2,283 | -1,576 |
| `commit` | 998 | 1,002 | +4 |
| `add` | 997 | 1,001 | +4 |
| `ls-tree` | 985 | 985 | 0 |
| `cat-file` | 974 | 967 | -7 |
| `rev-list` | 843 | 843 | 0 |
| `init` | 646 | 650 | +4 |
| `status` | 452 | 452 | 0 |
| `worktree` | 362 | 362 | 0 |
| `branch` | 304 | 304 | 0 |
| `for-each-ref` | 292 | 292 | 0 |
| `checkout` | 196 | 196 | 0 |
| `check-ref-format` | 187 | 187 | 0 |
| `remote` | 138 | 138 | 0 |
| `merge-base` | 113 | 113 | 0 |
| `show-ref` | 107 | 107 | 0 |
| `tag` | 95 | 95 | 0 |
| `config` | 95 | 95 | 0 |
| `archive` | 86 | 86 | 0 |
| `symbolic-ref` | 80 | 80 | 0 |
| `show` | 71 | 71 | 0 |
| `merge` | 52 | 52 | 0 |
| `log` | 51 | 51 | 0 |
| `diff` | 36 | 36 | 0 |
| `diff-tree` | 35 | 35 | 0 |
| `fetch` | 16 | 16 | 0 |
| `push` | 8 | 8 | 0 |
| `cherry-pick` | 7 | 7 | 0 |
| `clone` | 6 | 6 | 0 |
| `reflog` | 5 | 5 | 0 |
| `stash` | 1 | 1 | 0 |
| `ls-remote` | 1 | 1 | 0 |
| `ls-files` | 1 | 1 | 0 |
| **Total** | **12,099** | **10,528** | **-1,571 (-13.0%)** |

The combined repository-resolution queries account for the visible result:
`rev-parse` fell by 1,576 (40.8%). The two bounded object-read loops reduced
`cat-file` by seven in this suite; the remaining 967 calls come from other
read scopes and fixtures outside those loops.

The committed shape-based attribution parser produced:

| Attribution bucket | Before | After | Delta |
| --- | ---: | ---: | ---: |
| Production-read-shaped | 7,405 | 5,822 | -1,583 (-21.4%) |
| Fixture-setup-shaped | 2,736 | 2,748 | +12 (+0.4%) |
| Ambiguous or other | 1,958 | 1,958 | 0 |
| Malformed records | 0 | 0 | 0 |

These remain command-shape proxies, not exact caller attribution; the limits
recorded in `README.md` apply unchanged.

## Touched-package deltas

The original full-suite output did not preserve per-package lines, so this
table does not invent them. It records a supplemental comparison run on the
same host: each committed tree was extracted from Git, each package ran alone
with a fresh `GOCACHE`, and each command used
`go test <package> -count=1 -parallel 16` under `/usr/bin/time -p`.

| Package | Metric | Before | After | Delta |
| --- | --- | ---: | ---: | ---: |
| `internal/baseline` | Go package elapsed | 41.372s | 41.126s | -0.246s (-0.6%) |
| `internal/baseline` | Real | 44.58s | 44.38s | -0.20s (-0.4%) |
| `internal/baseline` | User | 39.78s | 39.26s | -0.52s |
| `internal/baseline` | System | 80.15s | 70.05s | -10.10s |
| `internal/agent` | Go package elapsed | 8.754s | 7.259s | -1.495s (-17.1%) |
| `internal/agent` | Real | 11.15s | 10.10s | -1.05s (-9.4%) |
| `internal/agent` | User | 12.04s | 12.24s | +0.20s |
| `internal/agent` | System | 5.20s | 5.53s | +0.33s |

The package comparison passed at both revisions. It isolates the touched
surfaces but is supplemental: its elapsed values cannot replace the
full-suite timing because package competition and build scheduling differ.

## Deliberately not done

- Caching across mutations: rejected because a stale repository fact can
  corrupt a Run; batching stays within provably immutable read scopes.
- Shared Git client: rejected because process count, not the seven local
  runners, was the measured cost, and packages keep distinct error contracts.
- Test removal: rejected because the suite is the behavior characterization,
  and prior measurement showed that removing the heaviest tests could not
  close the target gap.

## Measurement notes

One earlier after-census attempt was discarded after `internal/agent` failed:
the locally resolved ACP session did not advertise `sandbox_mode`, followed by
a cancellation-test milestone timeout. The successful retry used the same
committed procedure and produced the same 10,528-count census. This report
does not treat the discarded run as timing or pass evidence; the intermittent
agent failure is follow-up reliability work outside Task 06's documentation
slice.
