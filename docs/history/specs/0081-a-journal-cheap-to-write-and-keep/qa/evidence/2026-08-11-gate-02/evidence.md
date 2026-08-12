# QA evidence — 0081 gate 02

Build: `faa69a5a1da188b0955e7ff41b1dbfc7f195ae5e`.

All database journeys used temporary directories created by Go tests or the
disposable source archives named below. No operator Run Database was opened.

## Static gate

- `rtk roundfix spec check 0081-a-journal-cheap-to-write-and-keep --strict`
  exited 0 with no findings. It reported
  `SC-VOCABULARY-UNDOCUMENTED` skipped because the TechSpec has no Vocabulary
  Contract.
- The first `rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache
  GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp make verify` launch did not
  reach repository code because the selected disposable `GOTMPDIR` did not
  exist. After `rtk mkdir -p /private/tmp/roundfix-0081-qa02-gotmp`, the same
  command ran and reached `go test -parallel 16 ./...`; the sandbox then
  denied `api.github.com` because that domain is not on its network allowlist.
  The denial was not retried.
- Equivalent current-build evidence passed:
  `rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache
  GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp make fmt-check
  skills-sync-check skills-check build`, followed by `go test -count=1
  ./internal/store ./internal/tui ./internal/runevent`. The three packages
  exited 0 in 4.855 s, 0.942 s, and 0.568 s. The consumer and CLI-specific
  checks are recorded below.
- `rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache
  GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp make verify-docs` exited 0,
  including documentation contracts, measured sanctioned-ownership checks,
  corpus budget, and the full Spec consistency sweep.

## R01 — Run-start cost

Command:

```sh
rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp go test -count=2 ./internal/store -run '^TestJournalMeasurementHarness$' -v
```

The two fresh runs passed. Run-start p50 values for 0 / 1,000 / 10,000 events
were 51 / 113 / 41 us and 51 / 40 / 41 us. Database bytes were 61,440 /
2,162,688 / 21,135,360 in both runs. The 10,000-event p50 was lower than or
equal to the empty-journal control in both observations, so the former
journal-size curve was absent.

## R02 — Parallel Runs

Command:

```sh
rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp go test -count=1 ./internal/store -run '^TestParallelRuns$' -v
```

All three rehearsals passed. Six independent Runs persisted 1,536 of 1,536
events at `busy_timeout_ms=5000`, with `sqlite_busy=0` and
`completed_runs=6`. Wall time was 59,111 us (38,484 us per 1,000 events), and
the read-back and cancellation rehearsals proved contiguous cursors,
publisher order, and lock release.

## R03 — Same-boundary write cost

The pre-change source was extracted from
`a2a4c86b7570e4ef782ccc2ff390033f586d47fd` under
`/private/tmp/roundfix-0081-qa02-pre`; `task10-before-overlay.json` installed
the shared harness and old API adapter without changing that tree. The current
build used `task10-after-overlay.json` with the same harness and current API
adapter.

Commands:

```sh
rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp go test -overlay=<task10-before-overlay.json> -count=3 ./internal/store -run '^TestTask10SameBoundaryMeasurement$' -v
rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp go test -overlay=<task10-after-overlay.json> -count=3 ./internal/store -run '^TestTask10SameBoundaryMeasurement$' -v
```

Every sample persisted exactly 1,536 events at the publish-through-flush
boundary. Before values were 129,300 / 164,756 / 228,156 us per 1,000; current
values were 36,804 / 36,848 / 37,599. The medians fell from 164,756 to 36,848
us per 1,000, a 77.6% reduction with event count unchanged.

The focused batch suite also passed count, linger, immediate, explicit flush,
publisher order, contiguous cursors, raw payload bytes, publish-after-close,
failure preservation, shared Store writer, and all four ambiguous-commit
rehearsals.

## R04 — Bytes and retention shape

The R01 observation freshly confirmed unchanged database bytes at 0 / 1,000 /
10,000 events. Ten focused Store and CLI checks passed for bounded retention
eligibility, scan placement, exact raw-payload round-trip, terminal/age
selection, lifecycle preservation, zero-retention keep-everything, dry-run
non-mutation, actual pruning, and no implicit vacuum.

`rtk git diff --exit-code a2a4c86b7570e4ef782ccc2ff390033f586d47fd..HEAD --
.roundfixrc.yml go.mod go.sum
docs/adr/0008-run-event-payload-stores-raw-producer-json.md
docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md
docs/user-guide/run-database-lifecycle.md` exited 0 with no output. The
outside source was the pre-Spec maintainer Finding
`docs/findings/2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md`,
which records 2,915 MB of `run_events`, 1,645,457 events, and a largest Run of
42,000 events.

## R05 — Cockpit refresh and TUI sweep

Twelve current-build TUI checks passed for the forward cursor, header
projection, summary fallback, replay/live equivalence, stale event handling,
Attach replay, ten snapshots, 80x24 / 120x40 / 200x50 responsive states,
degenerate sizes, keyboard focus, Follow Mode, and color-independent markers.

`r05-10000-overlay.json` installed the QA-only
`r05-10000-event-harness_test.go.txt`. This command passed:

```sh
rtk env GOCACHE=/private/tmp/roundfix-0081-qa02-gocache GOTMPDIR=/private/tmp/roundfix-0081-qa02-gotmp go test -overlay=<r05-10000-overlay.json> -count=1 ./internal/tui -run '^TestQACockpitRefreshCostAtTenThousandEvents$' -v
```

After both 20- and 10,000-event backlogs, the same two new events cost exactly
2 header rows, 1 full row, and 67 payload bytes.

## R06 — Archive-safe pre-change consumer replay

The five relocated inputs matched the digests recorded before the move:

| Artifact | SHA-256 |
| --- | --- |
| `internal/store/testdata/2026-08-11-prechange-roundfix.db` | `1c38ed391769fb2c3ba5e48875ed868737bf0bc5659c12bdbf7ee8b19f5c5dec` |
| `internal/store/testdata/2026-08-11-prechange-consumer-expectations.json` | `3c4eb9511ee29dda3116ac8863e8b7f0b5a74fcaf594470ce6f0788b7f5ab04b` |
| `internal/store/testdata/2026-08-11-prechange-cockpit.golden` | `16eb1d01eb12a55728e0343c12768f04e40809359acc8214c368a42e5a8d9700` |
| `internal/store/testdata/task10-cli-consumer-harness_test.go.txt` | `204d5a2d9eab3d57b17bd5ef9050cf167b6112e1ca52b4eee3178f80a9439260` |
| `internal/store/testdata/task10-cockpit-consumer-harness_test.go.txt` | `d7c6d7da9220ad863b084564c100c5ad3012ddd40879c682e23a476fd02cb5f8` |

`git diff 1e5ed7e8..HEAD --name-status` classified all five as `R100` renames.
The current-worktree corpus replay passed both subtests: the real CLI harness
matched `events`, Attach, reconcile, and `gc`, and the real TUI harness matched
the complete Cockpit frame.

For the decisive archive simulation, `HEAD` was extracted under
`/private/tmp/roundfix-0081-qa02-archived`, the Spec directory was moved to
`docs/specs/_archived/0081-a-journal-cheap-to-write-and-keep`, and the same
`TestJournalConsumerCorpusReplaysEveryConsumer` command passed again. A
repository search for the Spec slug under `internal/store` returned no match.

The freshly built `bin/roundfix` returned exit 0 for `events --help`,
`attach --help`, `gc --help`, and `reconcile --help`. Their output retained the
canonical Run Database, Run Event Stream, Attach, Run Event Journal, and
Journal Retention vocabulary and public command shapes.
