---
task: task_03
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
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

- [x] Add the identifier-strategy domain renderer.
- [x] Add the structured HTTP and Better Auth backend renderer.
- [x] Validate safe Markdown content and ordering.
- [x] Add formatter and no-ADR fixtures.

## Acceptance Criteria

- [x] Domain guidance distinguishes project-owned IDs from provider, protocol,
  natural, and business identifiers.
- [x] Backend guidance contains the complete confirmed HTTP and Better Auth
  contract.
- [x] Removing all repository ADR fixtures does not remove operative guidance.
- [x] Unsafe marker content fails before postimage assembly.
- [x] Formatter and empty reapply produce zero managed delta.

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

## Result

Confirmed project decisions now render through typed Markdown renderers before
postimage assembly. The domain guide scopes UUID version 7 or exact
repository-defined guidance to new project-owned Internal Identifiers and
preserves provider, protocol, natural, and business source contracts. Both
single-context and multi-context domain templates carry the rule.

The backend guide now renders the confirmed HTTP mode, every normalized
exception in stable order, and a self-contained Better Auth ownership and
provider-protocol rationale. Dynamic inline fields are Markdown-escaped.
Structured render text rejects managed-marker shapes, template-token shapes,
control characters, non-canonical whitespace, and invalid inline newlines.

### Acceptance evidence

- `TestIdentifierGuidance` proves UUID version 7 remains limited to new
  project-owned Internal Identifiers in both domain layouts, all four
  source-contract exceptions remain operative, and repository-defined guidance
  renders exactly once after validation.
- `TestBetterAuthGuidance` proves the backend output contains the selected HTTP
  mode, every ordered typed exception, escaped repository content, Better Auth
  ownership, methods, scope, and the complete session, OAuth redirect,
  callback, and provider-protocol rationale.
- `TestProjectDecisionRendering` builds from a repository with no `docs/adr`
  fixture and still produces the complete domain and backend rules.
- `TestStructuredRenderSafety` proves marker-shaped, template-shaped, and
  non-canonical structured content returns an error with no Plan before
  postimages can be assembled.
- `TestProjectDecisionRendering` compares both affected postimages byte-for-byte
  with the selected formatter fixtures, applies the Plan, and proves an exact
  empty reapply reports the already-applied state with the full verified
  postimage set.

### Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestProjectDecisionRendering|TestIdentifierGuidance|TestBetterAuthGuidance|TestStructuredRenderSafety'`
  passed: 10 tests in 1 package.
- `rtk make verify` passed: 2,244 Go tests in 22 packages, 4 skill-contract
  tests, Roundfix skill synchronization, and the Roundfix build.

### Follow-up

None.
