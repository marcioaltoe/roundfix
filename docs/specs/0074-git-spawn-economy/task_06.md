---
task: task_06
spec: 0074-git-spawn-economy
status: pending
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
