---
task: task_05
spec: 0065-loop-order-and-verification-honesty
status: completed
type: chore
complexity: low
---

# Task 05: Synchronise the authoring and CLI Skills

## Overview

This Spec adds four `SC-*` rules that refuse Task Graphs which author cleanly
today, and it settles the loop order. Two Skills must teach that before the
Spec can close: `write-tasks`, which authors the graphs the rules now refuse,
and `roundfix`, whose `spec check` contract gained the rules.

This is the authorized tooling Task.

## Requirements

1. MUST document, in `write-tasks`, that a Task's Verification must be able to
   fail when no work was done, with the refused shape named: only
   repository-wide gates plus a clean-tree check.
2. MUST document, in `write-tasks`, that mutually unsatisfiable requirements
   are refused, and the section a rehearsal Task uses to declare its cases and
   their observation.
3. MUST document, in `roundfix`, the four new `SC-*` identifiers as part of the
   `spec check` contract.
4. MUST state the settled loop order in both Skills identically to the sources
   task_01 corrected, so task_04's rule passes over them.
5. MUST regenerate both mirrors with `make skills-sync`.
6. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/write-tasks/**`,
   `.agents/skills/roundfix/**`, their `skills/**` mirrors, this Task file, and
   the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [x] Document the three authoring rules in `write-tasks`.
- [x] Document the four `SC-*` identifiers in `roundfix`.
- [x] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [x] `write-tasks` names the refused Verification shape.
- [x] `write-tasks` names the contradiction rule and the rehearsal declaration.
- [x] `roundfix` lists the four new `SC-*` identifiers.
- [x] Both Skills state the settled loop order identically.
- [x] Both mirrors are byte-identical to their canonical Skills.
- [ ] `make verify` exits 0 after the regeneration chain.
- [x] No Go source file changed.

## Context

- instruction: `.agents/skills/write-tasks/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; both mirrors match.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/(write-tasks|roundfix)/|skills/(write-tasks|roundfix)/|docs/specs/0065-loop-order-and-verification-honesty/task_05\.md$|internal/baseline/(assets/(modules|setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 3, 4
  and 5.
- `_techspec.md` → Integration Points; Build Order 4.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.

## Result

### Implementation

- `write-tasks` now refuses the work-independent Verification shape composed
  only of repository-wide gates plus working-tree cleanliness checks. It also
  refuses declared `MUST`/`MUST NOT` clauses that require and forbid the same
  named state.
- `write-tasks` now requires a gate rehearsal or proof Task to declare every
  case under `## Rehearsal Cases` as
  `- Case: <case>; Observation: <observation>`.
- `roundfix` now documents the read-only `spec check` command contract and the
  four stable error identifiers: `SC-VERIFY-WORK-INDEPENDENT`,
  `SC-REQUIREMENT-CONTRADICTORY`, `SC-REHEARSAL-UNDECLARED`, and
  `SC-LOOP-ORDER-DIVERGENT`.
- Both Skills now carry the same settled order as task_01's sources: implement
  the graph including its authored gate, archive, open the Pull Request, watch
  until Clean, and merge.
- Regenerated the two embedded Skill mirrors and the deterministic ADR-0081
  Baseline digest and characterization fallout. No behavior or Go source was
  changed.

### Focused checks

- Red signal: focused searches before editing found none of the four new
  `SC-*` identifiers in the Roundfix Skill, and `write-tasks` did not name the
  refused repository-wide-gates-plus-clean-tree shape, the contradiction rule,
  or the `## Rehearsal Cases` declaration.
- `rtk make skills-sync` — exit 0; regenerated the embedded Skill bundle.
- `rtk make baseline-digests` — exit 0; reported `ok: true` and regenerated
  only deterministic files under `DERIVED_DIGEST_PATHS`.
- The first direct plan-characterization re-record could not open the
  sandbox-external default Go build cache. The same selector rerun with the
  repository-local `.gocache` exited 0 with 7 tests passed:
  `rtk env GOCACHE=<worktree>/.gocache rtk go test ./internal/baseline
  -count=1 -run TestBaselinePlanCharacterization
  -update-baseline-plan-characterization`.
- `rtk env GOCACHE=<worktree>/.gocache rtk go test ./internal/baseline
  -count=1 -run TestCatalogDiagnosticCharacterization
  -update-catalog-diagnostics` — exit 0 with 2 tests passed.
- `rtk env GOCACHE=<worktree>/.gocache rtk go test ./internal/cli -count=1
  -run 'Skills'` — exit 0 with 10 tests passed.
- Direct `cmp -s` checks for both canonical/mirror pairs exited 0 after the
  regeneration chain. `rtk git diff --check` also exited 0.

### Acceptance-criterion evidence

- Refused Verification shape: both `write-tasks` copies contain the explicit
  work-independent shape composed only of repository-wide gates plus
  working-tree cleanliness checks.
- Contradiction and rehearsal declaration: both copies contain the
  mutually-satisfiable `MUST`/`MUST NOT` rule and the exact
  `## Rehearsal Cases` entry syntax.
- Four identifiers: focused searches found each required `SC-*` identifier in
  both Roundfix Skill copies.
- Settled order: focused searches found the identical canonical order sentence
  in both Skills and their mirrors, matching the repository guide, shipped
  clause, and Baseline module asset.
- Mirror equality: direct byte comparisons exited 0 for `write-tasks` and
  `roundfix`.
- No Go source: the post-regeneration changed-path inventory contains no
  `.go` path. Every changed non-Task path is one of the two authorized Skill
  trees or deterministic fallout under `DERIVED_DIGEST_PATHS`.

### Daemon-owned verification

- `make skills-sync-check`, `go run -buildvcs=false ./cmd/roundfix skills
  check`, `make verify`, and both declared changed-path shell gates were not
  run in this Agent turn. The Daemon owns those commands, the remaining
  `make verify` acceptance criterion, and the terminal Task verdict.
