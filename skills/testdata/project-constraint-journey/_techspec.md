---
spec: project-constraint-journey-fixture
prd: _prd.md
---

# Project Constraint journey fixture — Technical Spec

The fixture models one authorized tooling change and rejects every path outside
its readable Project Constraint snapshot.

## Project Constraints

- Identifier strategy: not applicable because the fixture creates no Internal
  Identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable because the fixture changes no
  authentication provider, route, or HTTP Contract Decision. Source:
  `docs/agents/backend.md`.
- Active ADR obligations: applicable; the fixture must preserve confirmation
  gates and recoverable apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. For this disposable fixture, the maintainer
  expressly authorizes changes to exactly `.golangci.yml` and
  `scripts/verify.sh`; no other protected tooling file is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Verification

- Accept a changed-file set containing only `.golangci.yml` and
  `scripts/verify.sh`.
- Refuse a changed-file set containing `scripts/release.sh`.
