---
task: task_05
spec: 0078-roundfix-asks-for-the-review
status: pending
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill

## Overview

This Spec adds a Review Source mutation, two configuration keys, and a
Preflight refusal, so the Skill must teach them before the Spec can close. This
is the authorized tooling Task.

## Requirements

1. MUST document that `watch` and `resolve` publish one review request per Round
   after the Final Push when `review_source.request_review` is enabled, and that
   `fetch` never publishes one.
2. MUST document both configuration keys with their defaults, and state that
   asking is never Evidence — a published request does not advance a Round.
3. MUST document the Preflight refusal and both refused combinations, so a
   reader meeting the exit `2` knows which pair produced it.
4. MUST state that no automatic retry, backoff, or capacity wait exists, and
   that a refused request ends the Run.
5. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
6. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the request, the two keys, the refusal, and the no-retry rule.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill states that a Round's Final Push is followed by one review
      request, and that `fetch` publishes none.
- [ ] The Skill documents `review_source.request_review` and
      `review_source.request_command` with their defaults.
- [ ] The Skill states that a published request is not Evidence.
- [ ] The Skill names both refused configuration combinations.
- [ ] The Skill states that no automatic retry exists.
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
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0078-roundfix-asks-for-the-review/task_05\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 3, 4
  and 5.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
