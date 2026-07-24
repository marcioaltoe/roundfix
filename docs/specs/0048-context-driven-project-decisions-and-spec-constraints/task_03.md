---
task: task_03
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
type: backend
complexity: medium
---

# Task 03: Render confirmed project guidance

## Overview

Render the confirmed identifier strategy in the domain guide and the complete
HTTP and Better Auth contract in the backend guide. Both outputs remain
operative without an ADR and stable across formatting and empty reapply.

## Requirements

1. MUST render UUID version 7 only for new project-owned Internal Identifiers
   and state every confirmed source-contract exception.
2. MUST render repository-defined identifier guidance exactly after safe
   marker validation.
3. MUST render HTTP mode and every ordered typed exception in the backend
   guide.
4. MUST render the complete Better Auth ownership and protocol rationale.
5. MUST reject marker-shaped or non-canonical structured render content.
6. MUST keep operative guidance independent of ADR existence.

## Subtasks

- [ ] Add the identifier-strategy domain renderer.
- [ ] Add the structured HTTP and Better Auth backend renderer.
- [ ] Validate safe Markdown content and ordering.
- [ ] Add formatter and no-ADR fixtures.

## Acceptance Criteria

- [ ] Domain guidance distinguishes project-owned IDs from provider, protocol,
  natural, and business identifiers.
- [ ] Backend guidance contains the complete confirmed HTTP and Better Auth
  contract.
- [ ] Removing all repository ADR fixtures does not remove operative guidance.
- [ ] Unsafe marker content fails before postimage assembly.
- [ ] Formatter and empty reapply produce zero managed delta.

## Context

- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/templates/guides/domain-single-context.md`
- interface: `internal/baseline/assets/templates/guides/backend.md`

## Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestProjectDecisionRendering|TestIdentifierGuidance|TestBetterAuthGuidance|TestStructuredRenderSafety'` — expected: complete domain/backend output, exceptions, no-ADR, safety, formatter, and idempotency cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 2; User Stories 1–2 and 6; Core Features 4, 8–9.
- `_techspec.md` → Implementation Design: Data Models and API Contracts; Build Order 3.
- ADR-0076 → self-contained decision rendering.
