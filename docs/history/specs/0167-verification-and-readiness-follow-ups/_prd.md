---
spec: 0167-verification-and-readiness-follow-ups
status: archived
created: 2026-09-24
surfaces: [backend, cli]
archived: "2026-09-25"
source_slug: 0167-verification-and-readiness-follow-ups
---


# Verification and readiness follow-ups

Spec 0158 shipped independent Verification, precondition repair and access
readiness with three minor defects recorded by its second pre-PR review:

- In independent mode, a command that fails in the first run and again in the
  exclusive retry reaches the repair turn twice, and a command the retry shows
  passing is still carried as failed.
- A named repair Task that completed after the Agent removed or padded the
  repository command in its Task file makes the next `implement` of the Spec
  refuse at planning, because the verbatim check also inspects completed Tasks.
- A degraded full-access policy appears only in `profiles validate --json`; the
  text output and Doctor print `passed`.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes and ACP adapters
  only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0038 allows one Verification repair,
  ADR-0056 separates Task and Verification Capacity, which the retry still uses,
  ADR-0159 runs Verification past a failure only when declared independent,
  ADR-0160 lets only the frozen authorization open a red repository gate,
  ADR-0093 checks Spec consistency by citation, ADR-0104 accepts on evidence a
  Spec did not author, ADR-0130 keeps a path governed once bounded, ADR-0155
  makes the `qa` Task declare the matrix and ADR-0156 makes a declared promise
  name a consuming Task. This Spec's gate is bound by ADR-0080, ADR-0091,
  ADR-0096, ADR-0097 and ADR-0117. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the standing grant of 2026-09-21 for governed
  paths a slice needs, recorded in [_authorization.md](_authorization.md);
  bounded files: `internal/cli/cli_test.go`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

Spec 0165 carried two more defects here: a project-scope `roundfix init` makes
every later command warn about `runs.max_active`, and the Implement budget test
races a real 500 ms clock.

## Goals

- The repair turn sees each current failure exactly once.
- A completed repair never blocks the Spec's next Run.
- A degraded access policy is visible wherever readiness is reported.

## Core Features

1. **The retry's verdict wins.** For every command the exclusive retry reached,
   its verdict replaces the first run's; first-run failures are carried only for
   commands the retry never reached.
2. **Completed repairs are history.** The verbatim planning check applies only
   to named repair Tasks that are not yet completed.
3. **Degraded access is printed.** The text output of `profiles validate` and
   Doctor's profile readiness line name a degraded effective access policy.

4. **A quiet project init and a deterministic budget test.** Project Config
   templates omit `runs.max_active`, and the budget test's clock is injected.

## Non-Goals / Out of Scope

- Changing when a repair Task may enter or how it settles.

## Success Metrics

1. A command failing in both runs appears once in the repair request; a command
   passing on retry does not appear.
2. `implement` proceeds for a Spec whose completed repair Task no longer carries
   the command verbatim.
3. The text readiness output names a degraded policy.
4. A project-scope `init` followed by any command warns nothing, and the budget
   test passes without depending on wall-clock timing.

## Decisions

- **Latest evidence wins.** The exclusive retry exists to replace an unreliable
  verdict, so its results supersede the first run's for the commands it ran.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph: a duplicate in the repair request, a refused re-run and a silent
degraded policy are the failures under test.

## Research basis

The three defects and their reproductions are recorded in the archived PRD of
Spec 0158, from its second pre-PR review of 2026-09-24.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
