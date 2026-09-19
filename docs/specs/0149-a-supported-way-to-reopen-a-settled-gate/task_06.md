---
task: task_06
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: completed
type: backend
complexity: medium
---

# Task 06: A guard that guards nothing else

## Overview

Pre-PR review round 2 found two defects, both verified in source before this
Task was written. The first is self-inflicted: round 1 asked the command to
reuse `settle`'s active-Run guard, and it does — including the part that writes.

1. **The guard is not read-only.** `ensureNoSettleActiveRun` calls
   `reclaimOrphanedActiveRun`, which force-stops Agent sessions and calls
   `ReclaimOrphanedRun` to mark the Run failed, appending a daemon outcome
   event. The command's own usage text promises it "creates no Run, writes no
   Run Event Journal entry", and Core Feature 4 says the same. A `reopen` that
   reclaims an orphan keeps neither promise.
2. **The Spec Root check resolves before it validates.**
   `ensureReopenTaskInsideSpecsRoot` compares only `EvalSymlinks` results, so a
   manifest `file:` that lexically escapes the Spec directory but resolves
   through a symlink back inside passes. `replaceTaskFile` then builds its
   temporary file from the *lexical* parent and renames onto the lexical path,
   replacing an entry outside the root. The manifest's `file:` field is checked
   non-empty at `internal/spec/spec.go:588` and never validated further, so the
   escaping value is accepted.

This is the second corrective Task on this Spec, which is the accepted ceiling.
A further round of findings is a split or a reauthor, not a third correction.

## Requirements

1. MUST refuse when a Run is active for the Spec or the git root using a
   read-only check: no Run reclamation, no Agent session force-stop, no Run
   Event Journal write, no Run Database mutation of any kind.
2. MUST reject a manifest task path that is not confined to the Spec directory,
   judged lexically, before any symlink resolution.
3. MUST keep the resolved-path check as well; the lexical check is added to it,
   not substituted for it.
4. MUST keep every behavior Tasks 01 through 05 delivered.

## Subtasks

- [ ] Replace the reclaiming guard with a read-only active-Run refusal.
- [ ] Validate the manifest task path lexically before resolving it.
- [ ] Add a test for each.

## Acceptance Criteria

- [ ] With a dead-owner Active Run in the Run Database, the command refuses and
      the Run's state and its journal are unchanged.
- [ ] A manifest whose `file:` escapes the Spec directory is refused, naming
      that condition, and no file is written.
- [ ] The existing reopen tests still pass unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/reopen.go`
- interface: `internal/cli/cli.go`

## Verification

- `out="$(go test -count=1 -run "^TestReopenRefusesAnActiveRunWithoutReclaimingIt" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestReopenRejectsAManifestPathOutsideTheSpecDirectory" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

- [_techspec.md](_techspec.md) — Risks & Considerations

## Result

Implementation:

- `reopen` now checks the existing Run Database through a read-only Store and
  refuses any Active Run for the Spec target or git root without orphan
  reclamation, Agent Session force-stop, journal append, or Run mutation.
- The QA Task path must be lexically inside the selected Spec directory before
  the existing symlink-resolved Spec Root confinement check runs.
- Regression coverage exercises a dead-owner Active Run through the real Run
  Database and a lexical manifest escape that resolves through a symlink back
  inside the Spec Root.

Focused-check evidence:

- Acceptance criterion 1: `rtk env GOCACHE=/tmp/roundfix-task06-go-cache go
  test ./internal/cli -run
  'Test(ReopenRefusesAnActiveRunWithoutReclaimingIt|ReopenRejectsAManifestPathOutsideTheSpecDirectory)$'`
  passed. The Active Run test compares the complete Run record and Run Event
  Journal before and after refusal.
- Acceptance criterion 2: the same focused command passed. The manifest-path
  test requires the `outside Spec directory` refusal and compares the escaped
  QA Task bytes before and after the command.
- Acceptance criterion 3: `rtk env GOCACHE=/tmp/roundfix-task06-go-cache go
  test ./internal/cli -run '^TestReopen'` passed, covering the existing reopen
  command suite unchanged alongside the two regressions.
- Repository incremental check: `rtk env
  GOCACHE=/tmp/roundfix-task06-go-cache make verify-incremental` passed with
  process-table access. The first sandboxed attempt reached the suite but two
  existing force-stop integration tests could not enumerate the process table;
  the access-enabled rerun passed those tests and the full incremental gate.

Daemon Verification commands were not run; the Daemon retains that gate.
