---
spec: 0166-docs-changes-run-the-tests-that-read-them
status: active
created: 2026-09-24
surfaces: [backend]
---

# Docs changes run the tests that read them

## Executive Summary

Remove the no-set class from `ClassifyPath`, so documentation and root Markdown
fall through to both sets.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local classification only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — empty intersection with `GovernedPath`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The classifier

In `internal/verifyselect/verifyselect.go`, the `docs`/Markdown case that
returns `NoSet` is removed; such paths reach the default `BothSets`. Markdown
inside a package directory keeps selecting its owning set. The table in
`TestClassifyPath` and the fixture change sets that expected `NoSet` for
documentation are updated to the new contract.

## API Contracts

1. `verify-select` prints both sets for a documentation-only change.

## Coverage Map

- Goal 1 → The classifier; API Contract 1.
- Core Feature 1 → The classifier.
- Success Metric 1 → Testing Approach 1.
- API Contract 1 → The classifier.

## Integration Points

- **Spec 0155.** Owns the selector this Spec corrects.

## Testing Approach

1. **Documentation.** `docs/user-guide/run-database-lifecycle.md` and
   `README.md` alone select both sets; `internal/app/README.md` still selects
   core.
2. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The classifier (depends on: none).
2. Terminal QA (depends on: 1).

## Risks & Considerations

- **A slower docs-only pull request.** Accepted: correctness of the gate over
  its speed for documentation.
