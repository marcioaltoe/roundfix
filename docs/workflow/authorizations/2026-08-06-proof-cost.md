# Tooling authorization — the cost of proof

**Granted:** 2026-08-06 by the maintainer, in the write-prd authorization gate
for Spec 0080.

**Scope:** the two-tier verification contract and the two-stage QA gate.

## Bounded files

- `.agents/skills/qa-gate/**` and its `skills/qa-gate/**` mirror, regenerated
  by `make skills-sync`.
- `internal/baseline/assets/modules/core.json` and
  `internal/baseline/assets/modules/spec-workflow.json`, for the two-tier
  verification clause.
- `docs/agents/` guides carrying the regenerated managed-block postimages of
  those two modules, as adoption only.
- `Makefile`, when the incremental verification tier needs a target of its
  own.

Deterministic digest fallout of these edits is sanctioned by ADR-0081 and
needs no separate authorization. A hand-edited pin remains unauthorized.

## Outside this grant

- The Run Database, its schema, GC, and config keys — those belong to
  Spec 0081 and carry no protected-tooling surface.
- Any change to `internal/spec/qa.go` verdict semantics: ADR-0080 owns them,
  and a mechanical stage may compute the counts, never redefine them.
- `.github/**` workflows. If the complete fresh tier needs a CI change, it
  returns here for its own grant.

## Why the boundary is where it is

The gate's mechanical stage and the fast local verification tier are the same
split seen from two sides, so both must be authorable in one Spec or the
design gets made twice. Everything else the work touches is ordinary Go and
documentation.
