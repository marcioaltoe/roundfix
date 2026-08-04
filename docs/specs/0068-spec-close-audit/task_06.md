---
task: task_06
spec: 0068-spec-close-audit
status: pending
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

- [ ] Document the command and its contract in the canonical Skill.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.
- [ ] Confirm no Go source changed.

## Acceptance Criteria

- [ ] The Skill documents the command, both formats, the four kinds, and the
      exit codes.
- [ ] The Skill states the audit never reclaims.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

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
