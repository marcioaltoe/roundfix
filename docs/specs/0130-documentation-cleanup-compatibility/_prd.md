---
spec: 0130-documentation-cleanup-compatibility
status: active
created: 2026-09-09
surfaces: [infra, docs]
---

# Verification survives the documentation cleanup

PR #180 removes obsolete workflow documents as requested, but test fixtures and
the regeneration reader still consume their paths. Restore those contracts before
integrating the cleanup, preserving real authorization evidence and strict output
ownership. This is the five-file repair approved by the maintainer, not the full
implementation of Spec 0120.

## Project Constraints

- Identifier strategy: applicable — preserve Spec slugs, repository-relative paths,
  immutable Git object identities and the existing ownership tokens; introduce no
  new domain term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — filesystem and local Git reads only;
  no credential, HTTP interface, external mutation or paid call is introduced.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0149 states that the grant names the regeneration command and the tree names its outputs; the suite guard and audit read the same ownership declaration. This repair preserves that command/ownership boundary. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: "Autorizo os reparos", 2026-09-09, recorded in [the approved repair grant](_authorization.md); bounded files: `internal/baseline/derived_ownership_test.go`, `internal/speccheck/mechanical_test.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/suiteguardcontract/regeneration.go`, `internal/suiteguardcontract/regeneration_test.go`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

1. Approved regeneration works with Spec-contained grants and no workflow directory.
2. Historical authorization tests retain meaningful, reproducible evidence.
3. Unauthorized outputs and invalid grants remain blocked.

## Core Features

1. **Operative grant discovery.** Accept approved Spec-contained regeneration grants
   and compatible legacy inputs; reject unrelated or unapproved records.
2. **Ownership-based regeneration.** Resolve command-only grants from the existing
   ownership declarations, preserving frozen paths, dedicated commands and errors.
3. **Durable historical coverage.** Recover real grants from immutable main history,
   explicitly classify missing historical objects, and retain controlled assertions.
4. **Cleanup-compatible fixtures.** Exercise the regeneration contract with current
   grants while the legacy directory is absent, without changing canonical rules.

## Non-Goals / Out of Scope

- Implementing the broader Specs 0119–0129, changing public CLI/configuration,
  restoring docs/workflow in the delivered checkout, or adding an exemption list.
- Editing any implementation/test path outside the five-file grant, hand-editing
  generated expectations, weakening a check, or approving paid runtime use.

## Outside Evidence

OE-1 uses the original PR #180 CI failure, not a fixture created by this Spec:
https://github.com/marcioaltoe/roundfix/actions/runs/34276050740/job/102229198572
It names the deleted proof-cost authorization and stale derived artifacts.
The repair must exercise the corresponding repository boundary without the legacy
path, preserving strict refusal cases alongside the successful regeneration.

## Open Questions

None within the approved repair scope.
