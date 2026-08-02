---
spec: 0070-declared-unreachable-acceptance
status: active
created: 2026-08-02
surfaces: [backend, cli, docs]
---

# Declared unreachable acceptance

Some acceptance genuinely cannot be reached by any hermetic Verification. Spec
0058's remaining QA row needs a real tagged release published against six live
npm trusted-publisher bindings — an irreversible act against a live registry
that no gate may perform. Its QA closed `partial` with **zero findings and zero
failed rows**, and the Spec still could not archive, because `roundfix archive`
requires `verdict: pass`.

The only available exit was `qa_override`, which stamps an archive as verified
despite unverified evidence. Using it for a Spec whose evidence is complete
except for something unreachable spends the one mechanism reserved for genuinely
failed evidence, and it records the wrong thing: the Spec was not waved through,
it was verified as far as verification can go.

ADR-0080 already types blocked rows by cause. Nothing consumes that typing at
the archive boundary.

## Project Constraints

- Identifier strategy: not applicable — QA row identifiers, verdict values, and
  Spec slugs keep their existing contracts; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; this Spec
  reads QA reports and Spec artifacts. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0080 owns verdict semantics and the
  typed blocked-row counts this Spec consumes; nothing here may make a verdict
  more permissive than the evidence it summarizes. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized; the work is product code, CLI surface, and documentation.
  Source: `docs/agents/agent-instructions.md`.

## Goals

- A Spec can declare, up front, an acceptance no hermetic Verification can
  reach, and say why.
- A QA report can confirm that a blocked row matches a declared unreachable
  acceptance, distinctly from a row blocked by circumstance.
- A Spec whose only unmet rows are declared unreachable archives without
  `qa_override`, so the override goes back to meaning failed evidence.
- Nothing becomes archivable that is merely inconvenient to verify.

## Core Features

1. A Spec declares unreachable acceptance in its own artifacts, up front and
   with a stated reason and the human action that would satisfy it — never at
   QA time, and never by the gate's own judgement.
2. The QA gate matches a blocked row against the Spec's declarations and
   records it as declared-unreachable, distinctly from environment-blocked by
   circumstance and from finding-blocked.
3. A verdict may reflect that every unmet row is declared unreachable, without
   claiming those rows passed.
4. Archive accepts a Spec whose only unmet rows are declared unreachable, and
   stamps what remains unproven with the action that would prove it, so the
   archive record carries the debt rather than hiding it.
5. `qa_override` keeps its current meaning and remains required for genuinely
   failed or missing evidence; declared unreachability is never a route around
   a real failure.
6. A blocked row with no matching declaration keeps blocking exactly as it does
   today.

## Non-Goals / Out of Scope

- Letting the QA gate decide on its own that something is unreachable.
- Weakening `pass`, or letting a failed row archive under any name.
- Removing `qa_override`.
- Changing which journeys a QA matrix derives, owned by the QA gate contract.

## Success Metrics

- Spec 0058's QA report and artifacts, replayed, archive without
  `qa_override`, and the archive record names the release action that remains
  unproven.
- A blocked row with no matching declaration still blocks the archive.
- A row declared unreachable that the environment could in fact reach is
  reported as a wrongly declared row rather than accepted.
- `qa_override` still archives a Spec with a genuinely failed row, unchanged.

## Decisions

- Unreachability is declared by the Spec's author before the gate runs, not
  discovered by the gate. A gate that can excuse itself is not a gate.
- The archive record carries the unproven action rather than dropping it, so a
  reader of the archive learns what was never verified.
- This Spec evolves the archive boundary and never regresses it: every Spec
  that archives today under `pass` still archives identically, and every
  refusal except the declared-unreachable case is preserved.

## Open Questions

None.
