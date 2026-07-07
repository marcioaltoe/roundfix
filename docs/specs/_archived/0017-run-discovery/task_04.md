---
task: task_04
spec: 0017-run-discovery
status: completed
type: docs
complexity: low
---

# Task 04: Docs and skill sync for the Run discovery surface

## Overview

Close the SKILL-matches-CLI gate for the new surface: document `runs list` and
the Attach picker everywhere the CLI surface is described, and re-sync the
embedded skill bundle. Verifiable by the skills drift check and by reading the
updated docs against the shipped behavior.

## Requirements

1. MUST document `runs list` (columns, `--all`, `--active`, empty result,
   exit codes) and the no-argument Attach behavior in the README Commands and
   Command Boundaries sections, truthfully against implemented behavior.
2. MUST update the operational usage guide's monitoring flow to discover Runs
   with `runs list` and the Attach picker instead of requiring a captured run
   id.
3. MUST update the canonical roundfix skill (SKILL.md) with the new command
   surface and re-sync the embedded bundle so the drift check passes.
4. MUST keep top-level `--help` and the `runs`/`attach` usage text consistent
   with the docs.

## Subtasks

- [x] README Commands and Command Boundaries entries for `runs list` and the
      Attach picker
- [x] Usage guide monitoring flow update
- [x] roundfix SKILL.md update and `make skills-sync`
- [x] Drift and skills checks pass

## Acceptance Criteria

- [x] README documents the listing columns, both flags, the empty-result
      contract, and the no-argument Attach behavior.
- [x] The usage guide shows Run discovery without a captured run id.
- [x] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; Success Metrics. `_techspec.md` → Build Order 4. CLAUDE.md
SKILL.md-matches-CLI HARD RULE.

## Result

Updated the Run discovery documentation and skill bundle to match the shipped
CLI behavior.

Acceptance evidence:

- `README.md` now documents `roundfix runs list`, the stable columns
  (`run-id`, state, kind, target), `--active`, `--all`, the `No Runs found.`
  empty-result line, exit-code behavior, and no-argument `roundfix attach`.
- `docs/usage.md` now shows monitoring a Detached Run from a fresh terminal by
  using `roundfix runs list --active` and the no-argument Attach picker instead
  of requiring a previously captured Run id.
- `.agents/skills/roundfix/SKILL.md` now includes the Run discovery and Attach
  picker surface, and `rtk make skills-sync` regenerated the embedded
  `skills/roundfix/SKILL.md` copy.
- CLI help consistency was checked with `rtk go run -buildvcs=false
  ./cmd/roundfix --help`, `rtk go run -buildvcs=false ./cmd/roundfix runs
  --help`, and `rtk go run -buildvcs=false ./cmd/roundfix attach --help`.

Verification:

- `rtk make skills-sync-check` passed with no drift output.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed and reported
  all shipped Roundfix skills valid.
- `rtk make verify` passed: `rtk go test ./...` reported 838 tests across 18
  packages, `roundfix skills check` passed, and `go build` completed.
