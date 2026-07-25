---
spec: project-constraint-journey-fixture
status: active
---

# Project Constraint journey fixture

This disposable Spec fixture exercises bounded tooling authorization without
granting authority to change the Roundfix repository.

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

## Goals

- Prove that exact tooling bounds permit only the listed disposable files.

## Non-goals

- Changing any Roundfix repository tooling file.
