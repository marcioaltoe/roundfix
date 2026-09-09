---
spec: 0131-cleanup-delivery-contracts
status: active
created: 2026-09-09
surfaces: [infra, docs]
---

# Finish the cleanup delivery contracts

The combined cleanup reaches two stale verification assumptions: the loop-order
mutation fixture copies a retired sentence, and CI's shallow checkout omits the
permanent Git ancestor used by the repaired historical authorization tests.
Repair those assumptions while preserving what the gates reject.

## Project Constraints

- Identifier strategy: applicable — preserve existing clause IDs, Spec slugs and repository-relative paths; no new domain identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no application interface or credential access; the existing read-only CI checkout retains its credentials policy. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0153 defines explicit pre-PR review policy, including none; the fixture must track the approved clause without restoring the superseded Watch requirement. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on 2026-09-09: "Considere autorizada as mudanças e ajustes necessários para finalizar o trabalho nessa branch e fazemos o pull request, squash merge e sync main." The bounded repair is recorded in [_authorization.md](_authorization.md); bounded files: `internal/docscontract/corpus_test.go`, `.github/workflows/ci-verify.yml`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

1. The divergence detector is exercised against the current canonical clause.
2. CI obtains the historical evidence its unchanged test contract requires.

## Core Features

1. **Current-clause mutation.** All three source cases still produce a real
   ordering disagreement and assert error severity, all source labels/paths and
   the divergent declaration. Missing or ambiguous mutation targets fail.
2. **Complete CI history.** The existing checkout fetches full history while
   retaining its action version, read-only permissions, credentials policy,
   verification commands and suite budget.

## Non-Goals / Out of Scope

No production detector change, canonical policy change, skipped assertion,
changed timeout/budget, source reformat outside the Task, or broader pending Spec
implementation. Spec 0130 remains archived with its original evidence.

## Outside Evidence

OE-1: the existing PR #180 CI job at
https://github.com/marcioaltoe/roundfix/actions/runs/34346394563/job/102448882825
failed `TestOutputsForCommand` because Git could not read the retained ancestor.
The local `make verify-docs` failure names all three stale loop-order mutation
cases. These failures predate this Spec and establish its acceptance boundary.
