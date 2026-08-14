---
type: perf # feat | fix | perf | refactor
status: declined # open | promoted | declined
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: resolved directly on 2026-08-10 — the coupling was .git, not the docs tree; test git reads stopped writing the index
---

# One reader in cli still couples verify to the docs tree

## Opportunity

After the markdown contracts moved to `internal/docscontract`, touching a
markdown file should leave every code package `(cached)`. `internal/speccheck`
now behaves that way — proven on 2026-08-10: warm, append one line to a
`docs/backlog` entry, re-run, `(cached)`.

`internal/cli` does not, under the repository build cache. Measured the same
day, same file, `GOCACHE=$PWD/.gocache`:

```
go test ./internal/cli        → ok (54s, saves)
go test ./internal/cli        → (cached)
echo >> docs/backlog/<entry>  → one markdown line
go test ./internal/cli        → re-executes (54s)
```

So one input in the full cli run covers that file, and a docs-only commit
still pays 54 of the 61 seconds the post-docs `make verify` now costs.

Two measured oddities narrow the search and are worth keeping:

- Under the user-default `GOCACHE`, the same warm/touch/re-run sequence stayed
  `(cached)` — the coupling appeared only under the repository-local cache the
  Makefile exports, so an environment difference between the two runs is part
  of the mechanism.
- A `-run` bisection was attempted and produced an invalid result: an
  8KB alternation of 283 test names matched zero tests (`go test -list`
  confirmed), so both halves reported clean vacuously. A future bisection must
  verify its pattern matches before trusting a clean half.

## Value

The invalidation-domain split exists so a Spec-loop commit — which almost
always touches docs — re-runs seconds, not the 54-second cli package. This one
reader keeps the daily `make verify` at ~61s where ~10s is the design target.

## Shape

Find the reader with `GODEBUG=gocachetest=1` on a full package run and a
corrected bisection, then move that test's repository read behind a fixture or
into `internal/docscontract`. Worth settling in the same work: whether any cli
test walks the repository tree wholesale, since a directory listing records
every entry as an input and would explain a coupling no grep for the path
finds. This shape is non-binding.

## Resolution — 2026-08-10

The reader was `TestRunImplementRemovedQAFlagExplainsTaskGraph`, which ran the
CLI without an isolated work directory. Configuration resolution then started
from this package's own directory and walked up to the repository's `.git`,
recording it as a test input — so any later write under `.git` discarded the
whole package's cached result. Giving that one test a temporary home and work
directory removed `.git` from the package's inputs entirely (verified: zero
`.git` accesses in the package testlog).

The markdown was never the invalidator. Every reproduction that "touched a
doc" also ran `git checkout --` to revert it, and that write is what bumped
`.git`. Both measured oddities dissolve with the real cause: the
default-`GOCACHE` run stayed cached because no git command ran between its
warm-up and its probe.

Measured effect on two consecutive `make verify` runs with nothing changed:
61s before, **15.6s** after.
