---
task: task_04
spec: 0014-run-store-retention
status: completed
type: docs
complexity: low
---

# Task 04: Docs and skill sync

## Overview

Document Run store retention and the GC Command, and sync the Roundfix Skill for
the new command surface. Runs after the behavior tasks so the documented
contract matches shipped behavior and the SKILL.md-matches-CLI gate closes.

## Requirements

1. MUST document the `store.journal_retention` config key (default, `0`
   semantics), the `roundfix gc` command with `--dry-run`, and that retention
   never touches Active Runs, `runs` rows, or locks.
2. MUST note the self-healing preflight-sweep prune and that review artifacts
   under the spec tree are out of retention scope.
3. MUST update `.agents/skills/roundfix/SKILL.md` (and manifest) for the new
   `gc` command surface and regenerate the embedded copy with `make
   skills-sync`.
4. MUST leave no skill drift and no doc/behavior mismatch; CONTEXT.md already
   carries GC Command and Journal Retention.

## Subtasks

- [x] Docs for `store.journal_retention` + `roundfix gc [--dry-run]` + safety scope
- [x] Note self-healing sweep prune and out-of-scope review artifacts
- [x] Update SKILL.md/manifest for `gc` + `make skills-sync`
- [x] Verify no skill drift

## Acceptance Criteria

- [x] Docs accurately describe retention, `roundfix gc`/`--dry-run`, and the never-prune-Active-Runs guarantee.
- [x] SKILL.md matches the shipped CLI surface; the embedded copy is regenerated.
- [x] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories (documentation). `_techspec.md` → Build Order 4. ADR-0033.
CONTEXT.md → GC Command, Journal Retention. CLAUDE.md SKILL.md-matches-CLI gate.

## Result

- Documented `store.journal_retention` in `README.md` with default `336h`, Go
  duration semantics, `0` keep-everything behavior, and config context. The
  command docs now describe `roundfix gc`, `--dry-run`, reclaimed
  Run/journal/artifact counts, orphaned run artifact directories, and the
  guarantee that retention never deletes Active Runs, `runs` rows, or
  active-run locks.
- Documented the self-healing `implement`/`resolve`/`watch` preflight-sweep
  prune as best-effort and noted that Review artifacts under the Spec tree are
  outside retention scope.
- Updated `.agents/skills/roundfix/SKILL.md` and
  `.agents/skills/roundfix/agents/openai.yaml` for the GC Command, then ran
  `rtk make skills-sync` to regenerate `skills/roundfix/SKILL.md` and
  `skills/roundfix/agents/openai.yaml`.
- Verified behavior/docs alignment with `rtk go run ./cmd/roundfix gc --help`;
  it reported `roundfix gc [--dry-run]`, terminal-Run journal pruning,
  orphaned `runs/<id>` cleanup under the run artifact root, and the safety
  guarantee for Run rows, Active Run locks, and Review artifacts outside the
  run artifact root.
- Verification passed:
  - `rtk make skills-sync`
  - `rtk go run ./cmd/roundfix skills check` (`Roundfix skill check passed:
    roundfix`)
  - `rtk make skills-sync-check` (no drift output)
  - `rtk make verify` (`go test ./...` reported 801 passing tests in 18
    packages, `roundfix skills check` passed, and `go build` passed)
