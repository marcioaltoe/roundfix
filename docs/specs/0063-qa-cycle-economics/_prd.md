---
spec: 0063-qa-cycle-economics
status: active
created: 2026-08-01
surfaces: [backend, docs]
---

# QA cycle economics

A QA cycle costs roughly twenty minutes, and Specs routinely need four to seven
of them. Measurement shows the cost is not defect density and not the gate's
own work: `make verify` takes 148s on the cold Go cache every Run Worktree
starts with, against 5s warm. Discovery order compounds it — the gate stops at
the first static failure, so one stale assertion pinning an outdated literal
recorded `SG-01 fail` and marked eight behavioral journeys `skipped`, buying a
whole cycle to learn one line. And the gate is the only detector that exercises
real external surfaces, so a defect a thirty-second probe would expose waits
twenty minutes to surface. Evidence:
[what one cycle costs](../../findings/2026-07-29-qa-cycle-cost-is-cold-environments-and-agent-turns.md),
[the gate serializes findings and sits far from the defects](../../findings/2026-07-29-qa-cycle-latency-and-detector-placement.md),
and
[discovery stops at the first hard gate](../../findings/2026-07-29-qa-gate-round-economics.md).

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier is
  created; QA report identifiers, row identifiers, and verdict values keep
  their existing contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; this Spec
  handles no credential and opens no transport. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0080 keeps QA verdicts
  distinguishing environment-blocked rows, and the typed blocked-cause counts
  it introduced must survive any change to discovery order. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: on 2026-08-02
  the maintainer authorized tooling adjustment for the queued Specs, recorded at
  `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`; bounded files:
  `Makefile`. Deterministic digest fallout is sanctioned by ADR-0081. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- A static-gate failure stops reporting behavioral journeys as `skipped` when
  they could have run, so one stale assertion no longer hides eight findings.
- Cheap detectors run before the expensive gate, so a defect a short probe can
  expose does not wait for a full cycle.
- A Run Worktree does not pay a cold compilation cost on every QA cycle.
- Discovery order is measurable: a report says what it could not observe and
  why, not merely that it stopped.

## Core Features

1. A static-gate failure no longer hard-stops every behavioral journey. Rows
   that do not depend on the failed gate still run and report, and rows that
   genuinely cannot run record their cause under the existing typed
   blocked-count contract.
2. Fast detectors run ahead of the full gate, so defects they can expose are
   reported in their timeframe rather than the gate's.
3. The Run Worktree reuses a warm build cache across QA cycles, with the
   reuse proven not to leak state between Runs.
4. The report distinguishes "did not run because its subject was absent" from
   "did not run because an earlier gate failed", so a reader can tell wasted
   discovery from genuine blocking.
5. The change is bounded by a characterization corpus captured before it: every
   verdict a report reaches today, it still reaches, and no verdict becomes
   more permissive.

## Non-Goals / Out of Scope

- Weakening what a `pass` verdict means, or letting a failed static gate
  produce a passing report.
- Changing the QA matrix's derivation from the PRD, or how journeys are chosen.
- Re-litigating verdict semantics or blocked-row typing, which ADR-0080 owns.
- Supervisor-side workflow changes, which the loop discipline owns.

## Success Metrics

- A Spec whose static gate fails on one stale assertion reports the behavioral
  findings that were reachable, instead of eight `skipped` rows.
- A cold Run Worktree's QA cycle no longer pays the full cold-cache
  compilation, measured against the 148s baseline the finding recorded.
- No verdict in the characterization corpus becomes more permissive after the
  change.

## Decisions

- Discovery order is the target, not the gate's thoroughness: nothing here
  makes the gate check less.
- This Spec evolves QA and never regresses it. The declared break is that a
  static-gate failure no longer marks independent journeys `skipped`; any other
  change to verdicts, counts, or report shape is a defect.

## Open Questions

None.
