---
status: pending
created_at: 2026-07-29
updated_at: 2026-07-29
---

# QA cycles — the cost is a cold Run Worktree and the Agent's turn count, not the gate itself (2026-07-29)

A QA cycle on Spec 0061 took roughly twenty minutes of wall-clock. Measuring
where that time goes contradicted the obvious hypothesis: the repository gate
accounts for about an eighth of it, and compilation for a fraction of that.
This report separates what is specific to this repository from what every
repository using Roundfix pays, because the two need different fixes.

Companion to
[the round-economics finding](2026-07-29-qa-gate-round-economics.md), which
covers how many cycles happen. This one covers what one cycle costs.

## Measurements

All timings on one macOS machine, moderate load, no competing build.

| scenario | wall-clock |
| --- | --- |
| `make verify`, cold Go cache — what every QA cycle pays today | **148s** |
| `make verify`, warm cache, nothing changed | **5s** |
| `make verify`, warm cache, one file changed in `internal/cli` | **119s** |

Per-package test execution, forced with `-count=1`:

| package | time |
| --- | --- |
| `internal/cli` | 133.5s |
| `internal/baseline` | 95.6s |
| `internal/worktree` | 15.8s |
| `internal/agent` | 12.9s |
| remaining seven packages | ~15s combined |

Packages run in parallel, so the 148s wall-clock is essentially
`internal/cli` plus compilation. Inside it, eight journey and macro tests sum
to roughly 72s, the largest at 12.4s. In `internal/baseline`, a single test —
`TestTransactionFailureMatrix` — costs 28.9s, and the package performs 44
catalog loads across its tests.

The QA Run itself made **333 tool invocations** and **95 reasoning blocks**.
Measured compute across the gate, two extra full `internal/baseline` runs,
focused suites, two disposable clones, and one binary build accounts for
roughly 450–500s of the ~1200s cycle.

## 1. Every Run Worktree is a cold build environment — this generalizes

- **Symptom / evidence**: the Go build cache resolves to
  `.../worktrees/<repo>/<run-id>/.gocache`, inside the Run Worktree, which is
  created fresh per Run. Every QA cycle therefore recompiles the module from
  nothing. The measurements above put the avoidable portion between **29s**,
  when the change lands in the slowest package, and **143s**, when the change
  is documentation or Spec text only — which is the common case for a
  corrective cycle.

  The same shape appears in every consumer repository, with a different and
  usually larger price. `fluxus` configures
  `bootstrap: "bun install"` with `bootstrap_timeout: 10m` and `copy: []`, so
  each Run Worktree performs a complete dependency install before any Agent
  work. A Vortex session on 2026-07-29 reported roughly fifteen minutes per
  gate cycle because its gate rebuilds images and starts containers.

- **Root cause**: Worktree Bootstrap (ADR-0034) prepares each worktree
  independently. It is a per-worktree preparation step, not a shared cache, so
  identical work is repeated for every Run of the same repository.

- **Action / suggestion**: give a repository one durable build-cache location
  that Run Worktrees share, keyed by repository, while keeping the
  deterministic behavior Spec 0054 delivered — for Go that is `GOCACHE`, and
  the same principle covers a package manager's store, a bundler cache, or a
  container layer cache. Bootstrap stays for genuine per-worktree setup;
  cache-shaped state stops being rebuilt. Sharing must remain safe under
  concurrency, which Go's cache already is.

## 2. Over half the cycle is Agent turn latency — this generalizes

- **Symptom / evidence**: ~1200s of wall-clock against ~450–500s of measured
  compute leaves roughly 700s, about 58%, in the Agent's own turns: 333 tool
  invocations, including 12 `git status`, 10 `git -c …`, 8 `rg`, 5 `git log`,
  and 5 `git diff-tree`, plus 95 reasoning blocks.
- **Root cause**: the QA Agent discovers repository state one command at a
  time. Nothing hands it the facts a gate always needs — the changed-path set
  for the Spec's commits, the tooling-authorization comparison, the built
  binary, the Spec's own acceptance rows.
- **Action / suggestion**: precompute the invariant evidence and give it to
  the Agent at Batch start. The Daemon already knows the Run's target branch,
  the Spec's commit range, and the Task graph; the changed-path audit and the
  authorization comparison are deterministic and cheap to compute once. This
  is the largest single lever, and unlike the gate it costs nothing in
  coverage. It is also stack-independent, so it helps every repository.

## 3. Two slow test packages set this repository's floor — this does not generalize

- **Symptom / evidence**: `internal/cli` at 133.5s bounds the gate's
  wall-clock, and `internal/baseline` at 95.6s becomes the bound as soon as
  the first improves. Eight journey and macro tests account for ~72s of the
  former; one test accounts for 28.9s of the latter.
- **Root cause**: those tests build binaries, create disposable repositories,
  and drive real CLI flows. That is why they catch what unit tests do not, and
  it is also why they are slow.
- **Action / suggestion**: keep the coverage; change its shape. Consider
  separating journey suites so the fast gate stays fast and the journeys run
  as their own target, and reducing the 44 catalog loads in
  `internal/baseline` to a shared fixture. Do not trade this coverage for
  speed: the adapter-echo defect in Spec 0052 and the regeneration-ordering
  defect in Spec 0054 were both found by exactly this kind of test.

## What generalizes, and what does not

The measured numbers are this repository's. The first two causes are not:
every repository pays a cold Run Worktree and every repository pays the Agent's
turn count, and both are worse in a stack whose preparation is a dependency
install or a container build than in one whose preparation is a compiler cache.
Only the third cause — which two test packages are slow — is local, and it is
the one with the smallest reach.

## What worked — keep

- Spec 0054's deterministic cache default is what made this measurable at all:
  before it, cache location varied with the sandbox and the numbers would not
  have been comparable between runs.
- The gate's expensive parts are its valuable parts. Both product defects found
  by QA today came from building the real binary and driving it, not from any
  unit test.

## Routing — 2026-08-01

Routed to [Spec 0063](../specs/0063-qa-cycle-economics/_prd.md) on 2026-08-01.
