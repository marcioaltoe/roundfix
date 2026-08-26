---
verdict: fail
rows_blocked_precondition: 1
rows_blocked_environment: 0
rows_blocked_finding: 0
rows_blocked_declared: 0
precondition_check: "spec check --strict"
precondition_reason: "SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Identifier strategy required by docs/agents/spec-routing.md; SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Authentication and HTTP required by docs/agents/spec-routing.md; SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Active ADR obligations required by docs/agents/spec-routing.md; SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Tooling authority required by docs/agents/spec-routing.md"
---

# QA Report

## Results

| # | Status | Provenance |
| - | --- | --- |
| 0 | blocked | precondition |

## Precondition refusal

- check: spec check --strict
- reason: SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Identifier strategy required by docs/agents/spec-routing.md; SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Authentication and HTTP required by docs/agents/spec-routing.md; SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Active ADR obligations required by docs/agents/spec-routing.md; SC-CONSTRAINT-MISSING: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Tooling authority required by docs/agents/spec-routing.md

The gate stopped at this check before it built its QA matrix, so no requirement was measured and the row above records the refusal itself.

## Performed repairs

None.

## Assigned repair failures

None.

## Mechanical findings

### SC-CONSTRAINT-MISSING

- location: `docs/specs/9001-a-spec-that-contradicts-itself/_prd.md:1`
- detail: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Identifier strategy required by docs/agents/spec-routing.md
- fix: Add the Identifier strategy row to docs/specs/9001-a-spec-that-contradicts-itself/_prd.md with applicability, a reason, and an operative Source path.

### SC-CONSTRAINT-MISSING

- location: `docs/specs/9001-a-spec-that-contradicts-itself/_prd.md:1`
- detail: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Authentication and HTTP required by docs/agents/spec-routing.md
- fix: Add the Authentication and HTTP row to docs/specs/9001-a-spec-that-contradicts-itself/_prd.md with applicability, a reason, and an operative Source path.

### SC-CONSTRAINT-MISSING

- location: `docs/specs/9001-a-spec-that-contradicts-itself/_prd.md:1`
- detail: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Active ADR obligations required by docs/agents/spec-routing.md
- fix: Add the Active ADR obligations row to docs/specs/9001-a-spec-that-contradicts-itself/_prd.md with applicability, a reason, and an operative Source path.

### SC-CONSTRAINT-MISSING

- location: `docs/specs/9001-a-spec-that-contradicts-itself/_prd.md:1`
- detail: docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Tooling authority required by docs/agents/spec-routing.md
- fix: Add the Tooling authority row to docs/specs/9001-a-spec-that-contradicts-itself/_prd.md with applicability, a reason, and an operative Source path.

## Mechanical rows

| # | Status | Provenance |
| - | --- | --- |
| SC-CONSTRAINT-MISSING | blocked (finding: SC-CONSTRAINT-MISSING — waits on docs/specs/9001-a-spec-that-contradicts-itself/_prd.md omits Identifier strategy required by docs/agents/spec-routing.md) | mechanical finding |

## Mechanical skips

| Detector | Missing artifact |
| --- | --- |
| authorization bounded paths | tooling authorization |
| consequent-fix commit order | consequent-fix declaration |
| QA Report structure | QA Report |
| QA evidence paths | QA Report |
