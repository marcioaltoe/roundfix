---
spec: 0166-docs-changes-run-the-tests-that-read-them
status: active
created: 2026-09-24
surfaces: [backend]
---

# Docs changes run the tests that read them

Spec 0155 gave pull requests and Runs a selective gate, `make verify-changed`.
Its classifier sends a path under `docs/**`, and Markdown at the repository
root, to no test set. Tests read those files:
`TestDurableTableLifecyclePolicyCoversEveryTable` reads
`docs/user-guide/run-database-lifecycle.md`, `TestBaselineExamplesParse` reads
`README.md`, and the archive corpus tests read `docs/history/specs`. A change
to one of them passes the selective gate and fails only in the complete
`make verify` that pushes to `main` still run. Spec 0155 recorded this limit and
carried it here.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local classification only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- No documentation or root Markdown change passes the selective gate untested.

## Core Features

1. **No path selects nothing.** Paths under `docs/**` and Markdown outside a
   package directory select both test sets, like any other path the classifier
   cannot attribute to one set.

## Non-Goals / Out of Scope

- An allowlist of documents no test reads.
- Changing the Makefile, CI or `make verify`.

## Success Metrics

1. A change to `docs/user-guide/run-database-lifecycle.md` or `README.md` alone
   selects both sets.

## Recorded limits

- A core-only change still selects only the core set, although Baseline
  packages import core packages (`internal/baseline` imports `internal/spec`).
  Spec 0155 accepted this coupling to keep the selective gate selective; the
  complete `make verify` on pushes to `main` covers it. Found by the pre-PR
  review of 2026-09-25.

## Decisions

- **Fail safe over clever.** Every Spec delivery already changes code, so
  running both sets for documentation costs little; an allowlist would need its
  own proof that no test reads the listed files.

## Acceptance evidence

The Core Feature requires positive and negative public-contract evidence in the
Task Graph: the two reproductions recorded by Spec 0155 must select both sets.

## Research basis

Spec 0155's second pre-PR review of 2026-09-24 reproduced both escapes with
`go run ./cmd/verify-select -base HEAD~1`, recorded in its archived PRD.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
