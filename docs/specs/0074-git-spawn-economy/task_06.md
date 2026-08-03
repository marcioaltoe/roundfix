---
task: task_06
spec: 0074-git-spawn-economy
status: completed
type: docs
complexity: low
---

# Task 06: Publish the before-and-after

## Overview

The same procedure task_01 committed, run again on the finished tree, both
tables side by side, delta stated. The maintainer's target is on record —
the complete fresh suite under 60 seconds — and this report says plainly
whether the Spec reached it, by how much, and where the remaining time
lives if it did not. The 0071 campaign's rule applies: dominant costs that
cannot be trimmed get named, not hidden.

## Requirements

1. MUST re-run the census and the wall/user/sys attribution with the exact
   procedure `baseline/` records, on the same machine class, and publish
   both alongside the before numbers.
2. MUST state the fresh-suite wall clock against the 60-second target, and
   attribute the remaining time honestly if the target is missed.
3. MUST include the per-package deltas for the surfaces this Spec touched
   (`internal/baseline`, `internal/agent`) and the spawn-count deltas per
   subcommand.
4. MUST record what was deliberately not done (caching across mutations, a
   shared git client, test removal) with one line each on why.

## Subtasks

- [ ] Re-run both procedures and record the after-tables.
- [ ] Write the delta analysis against the target.
- [ ] Name the remaining costs and the rejected alternatives.

## Acceptance Criteria

- [ ] The report sits beside the baseline under the Spec and quotes both
      measurements from committed procedure runs.
- [ ] The 60-second verdict is stated in the first screen of the report.
- [ ] `git status --porcelain` shows no path outside
      `docs/specs/0074-git-spawn-economy/`.

## Verification

- `ls docs/specs/0074-git-spawn-economy/baseline/ | grep -c . | grep -qvx 1`
  — expected: exit 0; the after-measurement joined the baseline file.
- `grep -rqi "60" docs/specs/0074-git-spawn-economy/baseline/` — expected:
  exit 0; the target verdict is stated.

## References

- `_prd.md` → Goals 3; Success Metrics.
- `_techspec.md` → Build Order 6.

## Result

Published `baseline/after.md` beside the committed baseline. The report leads
with the missed-target verdict, quotes the successful before and after
full-suite measurements, includes every per-subcommand and attribution delta,
adds isolated deltas for `internal/baseline` and `internal/agent`, names the
remaining package floors, and records the three rejected alternatives.

Focused implementation evidence:

- The successful retry of the committed census procedure ran the fresh suite
  with `TEST_STATUS=0`, recorded 10,528 Git invocations, and reported zero
  malformed records. The first attempt was discarded after an intermittent
  `internal/agent` failure; the report records that limitation and follow-up.
- The separate committed timing procedure exited 0 and reported `real 83.38`,
  `user 121.55`, and `sys 262.24` seconds. The report compares those values
  with the committed `78.45 / 127.76 / 268.46` baseline and states that the
  60-second target was missed by 23.38 seconds.
- Supplemental fresh-cache package runs exited 0 at both
  `dbdad8ac1b8a2335ab88c65a0a47f50d86ef6c4e` and
  `a9d0097c590412dafd173fbcb4deaf1923bcae3a`: `internal/baseline` real time
  changed from 44.58s to 44.38s, and `internal/agent` changed from 11.15s to
  10.10s.
- `rtk git -c core.fsmonitor=false diff --check` exited 0. A focused
  `rtk git -c core.fsmonitor=false status --short` listed only this Task file
  and the new `baseline/after.md` report.

Acceptance-criterion evidence:

1. `baseline/after.md` sits beside `baseline/README.md` and quotes both
   measurements from the exact procedure recorded there, including revision
   and same-machine environment evidence.
2. The report's opening paragraph states `83.38` seconds against the
   under-60-second target and gives the 23.38-second miss before any detailed
   tables.
3. The focused status inspection listed no path outside
   `docs/specs/0074-git-spawn-economy/`.

The commands under `## Verification` were not run; the Daemon owns that
Verification and the terminal Task status.
