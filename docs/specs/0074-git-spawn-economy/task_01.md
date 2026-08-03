---
task: task_01
spec: 0074-git-spawn-economy
status: completed
type: docs
complexity: medium
---

# Task 01: Commit the spawn census baseline

## Overview

The comparison target every later Task is measured against. A PATH shim
counts every git invocation the fresh suite makes; `/usr/bin/time` attributes
wall clock against user and kernel CPU. Both procedures and their numbers are
committed under `baseline/`, so the after-measurement in task_06 re-runs the
same procedure instead of trusting memory. The motivating numbers from
design time (13,926 spawns; 36% utilization; sys 3× user) are re-derived
here on this Task's own machine and commit — the design figures are context,
not the baseline.

## Requirements

1. MUST record the census procedure exactly as reproducible commands: the
   PATH shim script, the invocation, and the parsing of per-subcommand
   counts.
2. MUST run the census against the current tree and commit the results: the
   total, the per-subcommand table, and the wall/user/sys attribution of
   `go test ./... -count=1 -parallel 16`.
3. MUST separate production-issued spawns from fixture-issued ones as far as
   the census allows, naming the method and its limits honestly.
4. MUST change nothing outside `docs/specs/0074-git-spawn-economy/` — the
   census observes; it does not optimize.

## Subtasks

- [ ] Write the shim and attribution procedure into `baseline/`.
- [ ] Run both and commit the measured numbers.
- [ ] Attribute production versus fixture spawns with stated limits.

## Acceptance Criteria

- [ ] `baseline/` contains the procedure and the measured census: total
      spawns, per-subcommand counts, wall/user/sys attribution.
- [ ] Every number in the file is measured by the committed procedure, not
      quoted from the design documents.
- [ ] `git status --porcelain` shows no path outside
      `docs/specs/0074-git-spawn-economy/`.

## Verification

- `ls docs/specs/0074-git-spawn-economy/baseline/ | grep -q .` — expected:
  exit 0; the baseline exists.
- `grep -rq "GOCACHE\|PATH" docs/specs/0074-git-spawn-economy/baseline/ || grep -rq "shim" docs/specs/0074-git-spawn-economy/baseline/`
  — expected: exit 0; the procedure is recorded, not just the numbers.
- `grep -rqE "[0-9]{3,} " docs/specs/0074-git-spawn-economy/baseline/` —
  expected: exit 0; measured counts are present.

## References

- `_prd.md` → Problem (the census that motivates the Spec).
- `_techspec.md` → Build Order 1; Testing Approach (the census, before and
  after, with the same shim).

## Result

Implemented the reproducible baseline under `baseline/`: an executable PATH
shim, a command-shape attribution parser, the exact census and timing
commands, the measured per-subcommand table, and the attribution limits.

Focused implementation evidence:

- `rtk proxy sh -n docs/specs/0074-git-spawn-economy/baseline/git-census-shim`
  exited 0.
- `rtk proxy awk -f docs/specs/0074-git-spawn-economy/baseline/attribution.awk /dev/null`
  exited 0 and emitted zero for every bucket.
- `rtk proxy test -x docs/specs/0074-git-spawn-economy/baseline/git-census-shim`
  exited 0.
- The first manual census attempt used a narrowed PATH and failed because
  `rustfmt` was unavailable. Its output and counts were discarded. A fresh
  replacement run preserved the caller PATH and exited 0.
- The replacement census command recorded 12,099 invocations with zero
  malformed records. Its subcommand parser attributed 7,405 as
  production-read-shaped, 2,736 as fixture-setup-shaped, and 1,958 as
  ambiguous or other.
- The separate fresh timing command exited 0 and reported `real 78.45`,
  `user 127.76`, and `sys 268.46` seconds.
- `rtk git -c core.fsmonitor=false diff --check` exited 0 after the final
  documentation update.

Acceptance-criterion evidence:

1. `baseline/README.md` records the successful census total, every
   per-subcommand count, command-shape attribution, and wall/user/sys output;
   `git-census-shim` and `attribution.awk` record the executable procedure.
2. The baseline uses only outputs from the fresh census and timing runs on
   commit `dbdad8ac1b8a2335ab88c65a0a47f50d86ef6c4e`; it does not copy the
   design-time census figures.
3. `rtk git -c core.fsmonitor=false status --short` showed changes only in
   this Task file and its `baseline/` directory after the final documentation
   update.

The Task's authored `## Verification` commands were not run; the Daemon owns
that Verification and the terminal Task status.
