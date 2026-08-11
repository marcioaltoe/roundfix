---
task: task_06
spec: 0068-spec-close-audit
status: completed
type: chore
complexity: low
---

# Task 06: Synchronise the Roundfix Skill with the new command

## Overview

The repository HARD RULE requires a Pull Request that changes CLI behaviour to
ship the Roundfix Skill update with it. This Spec adds a command, so the Skill
must teach it before the Spec can close. This is the authorized tooling Task.

## Requirements

1. MUST document `roundfix spec audit` in the canonical Roundfix Skill: its
   argument, both formats, the four survivor kinds, and the exit-code contract.
2. MUST state that the audit reports and never reclaims, so a reader does not
   expect it to remove anything.
3. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
4. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

   Neither test is in `BASELINE_DIGEST_STEPS`, so the sanctioned command
   reports "no changes" while the gate stays red. The flag names are verbatim
   because they do not match their test names.
5. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`. The
   Skill paths are authorized by the 2026-08-04 standing grant for CLI-contract
   synchronisation; any other path is out of scope — stop rather than widen it.
6. MUST NOT change the command's behaviour. This Task documents what shipped.

## Subtasks

- [x] Document the command and its contract in the canonical Skill.
- [x] Run `make skills-sync`, then `make baseline-digests` and both re-records.
- [x] Confirm no Go source changed.

## Acceptance Criteria

- [x] The Skill documents the command, both formats, the four kinds, and the
      exit codes.
- [x] The Skill states the audit never reclaims.
- [x] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [x] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `docs/agents/skill-dispatch.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `grep -q "spec audit" .agents/skills/roundfix/SKILL.md` — expected: exit 0.
- `grep -q "spec audit" skills/roundfix/SKILL.md` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0068-spec-close-audit/task_06\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority).
- `_techspec.md` → Integration Points; Build Order 6.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` → standing
  Roundfix Skill CLI synchronisation grant.
- ADR-0081.

## Result

### Implementation

- Added a `Spec close audit` reference section to the canonical Roundfix Skill.
  It documents the required slug, text and JSON formats, the
  `roundfix-specaudit/v1` schema, all four survivor kinds, the `0`/`1`/`2`
  exit-code contract, and the undelivered-artifact report.
- Stated that the command reports without changing Git state, the Run Database,
  or Spec artifacts. Residue reclaim commands remain operator actions.
- Regenerated the `skills/roundfix/` mirror and the ADR-0081 derived digest and
  characterization fallout through the sanctioned commands.

### Focused checks

- Pre-change signal: the canonical Skill had no `spec audit` documentation.
- `rtk make skills-sync`: exited 0.
- `rtk make baseline-digests`: exited 0 and reported five regenerated derived
  files.
- `rtk go test ./internal/baseline -count=1 -run TestBaselinePlanCharacterization -update-baseline-plan-characterization`:
  exited 0; 7 tests passed.
- `rtk go test ./internal/baseline -count=1 -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics`:
  exited 0; 2 tests passed.
- `rtk proxy diff -qr .agents/skills/roundfix skills/roundfix`: exited 0 with
  no differences.
- `rtk git diff --name-only HEAD -- '*.go'`: exited 0 with no paths.
- `make verify`: not run because it is declared in `## Verification` and is
  Daemon-owned for this Implement Task.

### Acceptance criterion evidence

1. The canonical and mirrored Skill sections name the required
   `roundfix spec audit <slug>` command, both formats, `pull-request`,
   `pending`, `residue`, `preserved`, and exits `0`, `1`, and `2`.
2. The section explicitly says the audit reports and never reclaims, and that
   printed reclaim commands are operator actions.
3. The focused recursive comparison found no difference between the canonical
   and mirrored Skill directories.
4. Repository Verification remains for the Daemon; no result is claimed here.
5. The focused Go-path diff inspection returned no path.

### Follow-ups

None.
