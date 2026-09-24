---
spec: 0165-a-measured-run-ceiling
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# A measured Run ceiling

## Executive Summary

Add a User Config ceiling on Active Implement Runs across repositories, enforced
by `implement` preflight, and make the budget test find its Run without the
header line.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Run Database only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0056, ADR-0080, ADR-0091, ADR-0093,
  ADR-0096, ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — empty intersection with `GovernedPath`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The ceiling

`internal/config/config.go` gains `runs.max_active` (integer, default 3, `0`
disables, negative refused). `internal/store/store.go` counts Active Runs of
kind implement across the whole Run Database. `implement` preflight in
`internal/cli/implement.go`, next to the per-checkout Active Run check, refuses
with exit 2 when the count has reached the ceiling, listing each holding Run's
id, repository and Spec, and creates no Run.

## The budget test

`TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` finds its Run
from the Run Database, or from any stderr line that names it, instead of the
header alone.

## API Contracts

1. `runs.max_active` bounds Active Implement Runs machine-wide; `implement`
   refuses at the ceiling with exit 2.

## Coverage Map

- Goal 1 → The ceiling; API Contract 1.
- Goal 2 → The budget test.
- Core Feature 1 → The ceiling.
- Core Feature 2 → The budget test.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- API Contract 1 → The ceiling.

## Integration Points

- **Spec 0156.** The delivery loop meets the refusal as a failed start and can
  retry later.

## Testing Approach

1. **Ceiling.** Two Active Runs in two repositories refuse a third at
   `max_active: 2`, naming both; `0` disables; a negative value is refused by
   config validation; below the ceiling `implement` proceeds.
2. **Budget test.** A stderr without the header still yields the Run.
3. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The ceiling (depends on: none).
2. The budget test (depends on: 1).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

- **A stale Active Run holding the ceiling.** The refusal names each holding Run
  so it can be stopped or reconciled.
