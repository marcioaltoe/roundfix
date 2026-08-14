---
task: task_11
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: backend
complexity: medium
---

# Task 11: Let the authorization detector read the fallout it already sanctions

## Overview

The mechanical stage this Spec built refused its own Spec, and the refusal is
wrong. `QA-AUTH-PATHS` blocked `task_06` five times because its commit touches
`docs/agents/setup-context.json` and four
`internal/baseline/testdata/plan-characterization/*.golden.json` files that the
`## Bounded files` list does not name.

The same authorization record says, three lines below that list:

> Deterministic digest fallout of these edits is sanctioned by ADR-0081 and
> needs no separate authorization. A hand-edited pin remains unauthorized.

The five paths are exactly that. Proved on 2026-08-11: running the sanctioned
`make baseline-digests` against the committed tree produces zero diff, so every
committed byte matches what regeneration emits. The detector parses `## Bounded
files` and nothing else, so it cannot see a sanction the record states in prose
one paragraph later.

This is the detector refusing legitimate work, which costs more than the round
it was built to save. It is not a reason to widen the grant: the grant already
covers these paths.

## Requirements

1. MUST let an authorization record declare its sanctioned deterministic
   regeneration in a form the detector reads, alongside the exact bounded files
   it already parses.
2. MUST keep the distinction the record draws: regenerated output is sanctioned,
   a hand-edited pin is not. A declared path is sanctioned only as the output of
   a declared command, never as a free path a Task may edit by hand.
3. MUST keep every other `QA-AUTH-PATHS` refusal intact. A changed path that is
   neither bounded nor declared regeneration output still blocks, and a record
   that declares no regeneration behaves exactly as today.
4. MUST add the declaration to `docs/workflow/authorizations/2026-08-06-proof-cost.md`
   naming `make baseline-digests` and the paths it owns. This grants nothing
   new: it makes machine-readable what the record's own prose and ADR-0081
   already sanction, and the maintainer's 2026-08-07 correction — that the list
   must name exact files rather than globs — applies to the declaration too.
5. MUST NOT let the detector run the regeneration itself. This Spec exists to
   make the pre-gate stage cheap; proving fallout by regenerating would put a
   full baseline rebuild in front of every QA gate.

## Subtasks

- [ ] Parse the sanctioned-regeneration declaration.
- [ ] Exempt declared outputs from the bounded-path refusal.
- [ ] Declare the regeneration in the 0080 authorization record.

## Acceptance Criteria

- [ ] A changed path declared as regeneration output does not raise `QA-AUTH-PATHS`.
- [ ] A changed path that is neither bounded nor declared still raises it.
- [ ] An authorization with no declaration behaves as it does today.
- [ ] The mechanical stage no longer blocks `task_06`.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/mechanical.go`
- `internal/speccheck/mechanical_test.go`
- `docs/workflow/authorizations/2026-08-06-proof-cost.md`
- `docs/specs/0080-cheap-detectors-run-before-the-gate/task_11.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestMechanicalAuthPathsAcceptsDeclaredRegenerationOutput$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestMechanicalAuthPathsAcceptsDeclaredRegenerationOutput'` — expected: exits 0. The case does not exist before this Task.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestMechanicalAuthPathsStillRefusesAnUndeclaredPath$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestMechanicalAuthPathsStillRefusesAnUndeclaredPath'` — expected: exits 0, proving the refusal survived rather than being traded away.
- `grep -q 'baseline-digests' docs/workflow/authorizations/2026-08-06-proof-cost.md` — expected: exits 0. The record does not name the command before this Task.

## References

- `_prd.md` → Goal 1.
- `docs/workflow/authorizations/2026-08-06-proof-cost.md` → the sanction in prose.
- ADR-0081.
- `qa/qa-report-2026-08-11-02.md` → the five `QA-AUTH-PATHS` refusals.

## Result

### Implementation

- Added a machine-readable `## Sanctioned regeneration` YAML declaration that
  binds `make baseline-digests` to exact output files. The declaration lists
  all sixteen deterministic outputs present in the actual `task_06` commit;
  the live failure report contains more fallout than the five paths summarized
  in this Task's overview.
- Kept bounded files and regeneration outputs as separate authorization
  classes. The detector accepts a regeneration output only when its declaration
  names a non-empty command and an exact repository-relative file; missing
  commands, globs, invalid paths, and undeclared paths remain unauthorized.
- The detector reads the declaration only. It does not execute the regeneration
  command or otherwise add baseline rebuild work to the mechanical stage.

### Focused checks

- Before the production edit,
  `GOCACHE=/private/tmp/roundfix-task11-gocache rtk proxy go test ./internal/speccheck -run '^TestMechanicalAuthPaths' -count=1`
  exited 1 because the new declared-output case still raised
  `QA-AUTH-PATHS` for `internal/baseline/testdata/catalog.digest`.
- After the production edit,
  `GOCACHE=/private/tmp/roundfix-task11-gocache rtk go test ./internal/speccheck -run '^TestMechanicalAuthPaths' -count=1`
  exited 0 with nine passing cases. The group includes the declared-output
  acceptance, the undeclared-path refusal, commandless and glob declaration
  refusals, and the pre-existing no-declaration refusal.
- `GOCACHE=/private/tmp/roundfix-task11-gocache rtk go test ./internal/speccheck -count=1`
  exited 0 with 176 passing package tests.
- A temporary Go overlay invoked `RunMechanicalStage` against the actual
  `task_06` commit `9d36349ac0060a4cff5ad1cbe1f782e4d6605e20` and the amended
  authorization. With repository-local fsmonitor disabled through Git's
  process environment, the focused check exited 0 and found no
  `QA-AUTH-PATHS` finding.

### Acceptance-criterion evidence

1. `TestMechanicalAuthPathsAcceptsDeclaredRegenerationOutput` exercises a real
   temporary Git commit and raises no `QA-AUTH-PATHS` finding for an exact
   declared output.
2. `TestMechanicalAuthPathsStillRefusesAnUndeclaredPath` keeps a valid
   regeneration declaration present and proves an unrelated changed path still
   raises exactly one `QA-AUTH-PATHS` finding.
3. The existing red subtest in `TestMechanicalAuthPaths`, included in the
   nine-case focused run, uses an authorization with no regeneration
   declaration and still refuses its unbounded path.
4. The overlay check exercised the real `task_06` SHA and authorization record;
   all sixteen formerly unbounded generated paths were accepted.

### Follow-up note

- The first live-commit overlay attempt exposed a separate existing issue:
  this worktree's fsmonitor IPC diagnostic entered `git diff-tree` combined
  output and the detector treated that diagnostic text as a changed path. The
  same actual-commit check passed with `core.fsmonitor=false`. Repairing Git
  stderr handling is outside this Task's bounded slice.

### Handoff boundary

- The Daemon-owned commands under `## Verification` were not run. Task status
  remains Daemon-owned; no commit, push, or pull request was created.
