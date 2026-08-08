---
task: task_06
spec: 0084-an-update-that-can-run
status: pending
type: backend
complexity: high
---

# Task 06: Emit the fourteen structural clauses again

## Overview

Three of the eight measured repositories stop on the same fourteen structural
Normative Clauses that the Source Baseline accounts for and the current catalog
no longer emits. The fail-closed retention gate is behaving exactly as ADR-0058
requires; the loss is upstream, in the catalog. This slice restores those clauses
to their owning modules so retention accounting on an unchanged repository stops
reporting them unaccounted.

## Requirements

1. MUST restore, to their owning catalog modules, the fourteen structural
   Normative Clauses the maintainer's authorization names, using the clause text
   the Source Baseline corpus already carries so retention accounts for them
   mechanically.
2. MUST bump the version of every rule that receives a restored clause.
3. MUST NOT weaken, reword to a weaker enforcement, or reject any restored clause;
   the maintainer chose restoration over rejection.
4. MUST NOT change any module's decisions, capabilities, required skills, skill
   dispatch, or template selection.
5. MUST regenerate derived pins with the sanctioned command rather than by hand,
   per ADR-0081.
6. MUST leave per-project customization as decisions, never clauses: identifier
   strategy, HTTP contract, authentication provider, and database choice stay
   maintainer-answered.
7. MUST prove, by test, that retention accounting over the standard profile
   reports zero unaccounted entries for the fourteen restored identities.
8. MUST change only the paths the authorization bounds.

## Subtasks

- [ ] Map each of the fourteen identities to its owning module.
- [ ] Restore each clause with the Source Baseline's text and bump its rule.
- [ ] Record accounting for any identity whose owner legitimately moved.
- [ ] Regenerate derived pins with the sanctioned command.
- [ ] Cover retention accounting reporting zero unaccounted for the fourteen.
- [ ] Confirm no decision, capability, skill, or template selection changed.

## Acceptance Criteria

- [ ] Each of the fourteen identities named in the authorization resolves to a
      clause or rule the catalog emits.
- [ ] Retention accounting over the standard profile reports none of the fourteen
      as unaccounted.
- [ ] Every rule that received a clause carries a higher version than before.
- [ ] The catalog's decision, capability, skill, and template selections are
      unchanged.
- [ ] Derived pins match what the sanctioned regeneration command produces.
- [ ] The change touches only paths the authorization bounds.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-restore-structural-clauses.md`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'Retention' -v > /tmp/0084-task-06-a.log 2>&1 && grep -q '^--- PASS: .*Retention' /tmp/0084-task-06-a.log` — expected: exits 0, proving the retention corpus ran and passed.
- `go test ./internal/baseline/ -run 'Catalog' -v > /tmp/0084-task-06-b.log 2>&1 && grep -q '^--- PASS: .*Catalog' /tmp/0084-task-06-b.log` — expected: exits 0, proving catalog validation accepts the restored clauses.
- `grep -c 'clause.backend.' internal/baseline/assets/modules/backend.json > /tmp/0084-task-06-c.log 2>&1 && grep -qv '^0$' /tmp/0084-task-06-c.log` — expected: exits 0, proving the module declares backend clauses where it previously declared none.
- `make baseline-digests > /tmp/0084-task-06-d.log 2>&1 && git diff --quiet -- internal/baseline/assets` — expected: exits 0, proving every derived pin already matches the sanctioned regeneration.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 6; Risks & Considerations.
- `_prd.md` → Core Feature 5; Goals 1 and 4; Open Questions.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the three repositories that stop on these fourteen.
- ADR-0058, ADR-0081, ADR-0099.
