---
task: task_04
spec: 0077-a-green-check-is-not-a-review
status: completed
type: chore
complexity: low
---

# Task 04: Synchronise the Roundfix Skill

## Overview

This Spec changes when `watch --until-clean` will and will not merge, so the
Skill must teach it before the Spec can close. This is the authorized tooling
Task.

## Requirements

1. MUST document that a Review Source refusal resolves as skipped evidence and
   never verifies a head, and that `watch --until-clean` will not merge it.
2. MUST document that an unrecognised signal resolves pending, and state plainly
   that a green check is not evidence a review ran.
3. MUST document that no automatic retry exists yet, so a reader does not wait
   for one — the follow-on Spec owns it.
4. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
5. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

6. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
7. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the refusal, pending, and no-retry contract.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill states a refusal never verifies a head and blocks the merge.
- [ ] The Skill states an unrecognised signal is pending, and that a green check
      is not evidence a review ran.
- [ ] The Skill states no automatic retry exists yet.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0077-a-green-check-is-not-a-review/task_04\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 2 and 4.
- `_techspec.md` → Integration Points; Build Order 4.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.

## Result

Updated the Roundfix Skill's Review Source Evidence contract so only a
recognised review-completed current-head signal can verify a head. The Skill
now states that a refusal resolves to `skipped`, never verifies the head, and
prevents `watch --until-clean` from clearing it for merge; an unrecognised
signal resolves to `pending` even when its check is green; and a refused head
has no automatic retrigger or retry yet.

Regeneration and focused-check evidence:

- Initial inspection:
  `rtk grep -n 'green check is not evidence\|does not automatically retrigger\|unrecognised.*pending\|will not merge' .agents/skills/roundfix/SKILL.md`
  exited 1 before the edit, confirming the required contract was absent.
- `rtk make skills-sync` exited 0 and regenerated the embedded Roundfix Skill.
- `rtk make baseline-digests` exited 0 and reported deterministic changes only
  under the ADR-0081 `DERIVED_DIGEST_PATHS` roots.
- `rtk go test ./internal/baseline -count=1 -run TestBaselinePlanCharacterization -update-baseline-plan-characterization`
  exited 0 with 7 passing tests.
- `rtk go test ./internal/baseline -count=1 -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics`
  exited 0 with 2 passing tests.
- `rtk go test ./internal/baseline -count=1 -run 'Test(BaselinePlanCharacterization|CatalogDiagnosticCharacterization)'`
  exited 0 with 9 passing tests after the final regeneration edit.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  exited 0, proving the mirror is byte-identical.
- ``rtk grep -n 'never verifies a head\|will not merge that head\|unrecognised signal resolves to `pending`\|green check is not evidence\|does not automatically retrigger' .agents/skills/roundfix/SKILL.md``
  exited 0 for the refusal, pending, green-check, and no-retry wording.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.
- The final changed-path postflight contains only the canonical and mirrored
  Roundfix Skill, this Task file, and regenerated files under
  `internal/baseline/assets/setups` or `internal/baseline/testdata`. No Go
  source file changed.

Acceptance evidence:

1. The canonical Skill says a Review Source refusal resolves to `skipped`,
   never verifies a head, and `watch --until-clean` will not merge or clear the
   head for merge.
2. The canonical Skill says an unrecognised signal resolves to `pending` even
   for a successful check, and states plainly that a green check is not
   evidence a review ran.
3. The canonical Skill distinguishes refusals from transient failures and says
   Roundfix does not automatically retrigger or retry a refused head; follow-on
   work owns that policy.
4. The byte comparison above proves `skills/roundfix/` matches
   `.agents/skills/roundfix/` after `make skills-sync`.
5. `make verify` was not run because it is a declared `## Verification`
   command owned by the Daemon in this Implement turn.
6. The changed-path postflight contains no `.go` file and no path outside the
   authorized Skill pair, this Task file, and ADR-0081 derived-digest roots.
