---
spec: failed-recovery
prd: _prd.md
created: 2026-09-10
---

# A failed gate accepts its repair — Technical Spec

## Executive Summary

Narrow the Task Graph's gate-staleness rule so it invalidates only a gate whose
recorded verdict is a pass. The trade-off is deliberate and small: one condition
plus its characterization, against the alternative of a general recovery
contract that Spec 0129 already owns and that this repair cannot wait for.

## Project Constraints

- Identifier strategy: applicable — preserve existing Spec slugs, Task IDs and error identities; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0091 keeps the authored QA gate terminal, and this Spec preserves that; ADR-0057 keeps the Daemon the exclusive writer of Task status, so recovery must not ask a human or an Agent to edit a status; ADR-0117 places a check at the stage that can establish it, and gate staleness is established only where a passing verdict exists; ADR-0096 requires mechanical facts before the QA Agent turn and is preserved unchanged, because this Spec alters when a graph loads and not which facts the gate proves. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The change is ordinary source in `internal/spec` with its test. Source: `docs/agents/agent-instructions.md`.

## System Architecture

| Component | File | Responsibility |
| --- | --- | --- |
| Task Graph loader | `internal/spec/spec.go` | Validate the authored QA declaration, including whether a settled gate stands above incomplete dependencies. |
| Loader characterization | `internal/spec/spec_test.go` | Record today's answer for each gate verdict and hold the change to it. |

No other component participates. The error type, its message and its exported
identity are unchanged, so every caller that matches on them keeps working.

## Implementation Design

The loader currently treats a gate as settled when its status is completed or
failed, then reports every dependency in the gate's closure that is not
completed. Restrict that first condition to a completed gate.

The reasoning is what the rule protects. A completed gate carries a passing
verdict, and a pass above work the gate never saw is a false certificate — the
refusal is the whole point. A failed gate carries no certificate. Refusing to
load it protects nothing and removes the only documented recovery: author
corrective Tasks, let them run, let the gate run again. Recovery must not
require editing a Task status, because the Daemon owns status.

The narrowing is exact. A gate whose status is absent, malformed, or any value
other than completed continues down the existing path; only the literal failed
status gains acceptance. That keeps the change fail-closed on evidence rather
than on the absence of evidence.

## Coverage Map

- PRD Goal 1 → Task Graph loader.
- PRD Goal 2 → Task Graph loader, Loader characterization.
- User Story 1 → Task Graph loader.
- User Story 2 → Loader characterization.
- Core Feature 1 → Task Graph loader.
- Core Feature 2 → Task Graph loader.
- Core Feature 3 → Loader characterization.
- Core Feature 4 → Task Graph loader, by not touching status at all.

## Testing Approach

Focused loader tests over real Spec fixtures on disk. Required observations:

1. A completed gate above one incomplete dependency still refuses, with the
   same error type and the dependency named.
2. A failed gate above one incomplete dependency loads.
3. A failed gate still must cover every non-QA leaf; an uncovered leaf refuses
   as it does today.
4. Loading changes no Task file, proven by comparing the Spec directory before
   and after.
5. The real repository's Spec corpus still loads, so no active Spec regresses.

Observation 5 rests on evidence this Spec did not author: the repository's own
eleven active Specs and its archived corpus, written by earlier work under a
contract this Spec is changing. If that corpus cannot be read, the row records
blocked with that reason.

## Build Order

1. Characterize the loader's present answer for each gate verdict (depends on: none).
2. Narrow the staleness condition to a completed gate (depends on: 1).
3. Terminal QA (depends on: 2).

## Risks & Considerations

The risk is widening acceptance past the failed verdict and letting a malformed
or absent status load a graph that should refuse. The design answers it by
matching the literal failed status and leaving every other value on the
existing path, and the characterization holds that boundary.

## Decisions

- One condition, not a recovery framework. Spec 0129 owns the general contract.
- The error type and message stay byte-identical so no caller or fixture that
  matches on them changes.

## Vocabulary Contract

No token is coined. The change reuses the existing stale-gate error identity,
whose glossary owner is unchanged.
