---
spec: 0064-spec-artifact-consistency-gate
status: active
created: 2026-08-01
surfaces: [backend, cli, docs]
---

# Spec artifact consistency gate

Half of the QA findings this repository paid twenty-minute cycles for were not
code defects — they were contradictions between a Spec's own artifacts, each
detectable by reading files. Spec 0058's QA-001 found a PRD promising a
preflight check npm makes impossible while its own ADR acknowledged the limit;
QA-004 found the workflow emitting five failure prefixes while the runbook
documented four. Spec 0056 repeated both shapes: F-001 found Project
Constraints omitting an ADR that governs the changed behavior, and F-002 found
a Core Feature contradicting the TechSpec decision that superseded it. Four
findings, four full gate cycles, zero of them needing a running binary.

The gate is correct to catch these. It is the wrong place to catch them: it
runs last, costs twenty minutes, and serializes discovery behind everything
else. A static consistency check over the Spec folder would find the same
contradictions in seconds, before a Run exists.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier is
  created; Spec slugs, ADR numbers, Task identifiers, and diagnostic codes keep
  their existing contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; this Spec
  reads local Markdown and writes diagnostics. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0080 owns QA verdict semantics, and
  nothing here may weaken or duplicate a verdict; this gate reports before QA
  rather than inside it. ADR-0088 and ADR-0091 own the authored QA gate as a
  typed Task node, under which this Spec's own graph is authored. ADR-0093 and
  ADR-0094 are minted by this Spec and govern its detection boundary and its
  presence-awareness. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: on 2026-08-02
  the maintainer authorized tooling adjustment for the queued Specs, recorded at
  `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`; bounded files:
  `Makefile`, plus the owned skill pair if the check ships as a skill. Deterministic digest fallout is sanctioned by ADR-0081. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- A contradiction between a Spec's PRD, TechSpec, ADRs, and Task Graph is
  reported in seconds, before a Run is created.
- Every Project Constraint row cites an operative source, and every ADR that
  governs the changed behavior appears in the Active ADR row.
- Vocabulary a Spec's implementation emits is present in the documentation the
  Spec claims documents it.
- The check reports; it never edits a Spec and never substitutes for QA.

## Core Features

1. A read-only consistency check over one Spec folder, runnable before
   `implement` and fast enough to run on every authoring change.
2. Contradiction detection between artifacts: a PRD requirement the TechSpec's
   Decisions supersede without the PRD being amended is reported, naming both
   locations.
3. Constraint completeness: every Project Constraints row states applicability
   with a reason and cites an operative source path that exists, in both the
   PRD and every present TechSpec.
4. ADR coverage: an accepted ADR whose subject the Spec's changed surface
   touches, and which the Active ADR row omits, is reported as a gap rather
   than inferred silently.
5. Coverage completeness: every PRD user story and core feature appears in at
   least one Task's References, and every Task's declared references resolve.
6. Emitted-vocabulary coverage: an identifier the Spec's own artifacts declare
   as user-facing vocabulary — failure prefixes, exit codes, diagnostic codes —
   is present in the documentation the Spec names as documenting it.
7. Findings are reported with the file and line that carries each side of a
   contradiction, so the fix is unambiguous.

## Non-Goals / Out of Scope

- Replacing, weakening, or duplicating the QA gate. This check is cheap and
  static; QA remains the behavioral authority.
- Editing Spec artifacts automatically. The gate reports; the Supervisor
  decides.
- Judging whether a decision is correct — only whether the artifacts agree
  about it.
- Enforcing prose style, wording, or document length.

## Success Metrics

- Replayed against Specs 0056 and 0058 as they stood when their gates ran, the
  check reports QA-001, QA-004, F-001, and F-002 without executing anything.
- The check completes in seconds on a Spec folder, fast enough to run on every
  authoring change.
- The check reports no finding against a Spec whose artifacts agree, measured
  across the archived Spec corpus, so its false-positive rate is observable.

## Decisions

- The check is read-only and reports before a Run exists; it never mutates a
  Spec and never emits a QA verdict.
- Findings name both sides of a contradiction with file and line. A finding a
  reader must go hunting for costs more than it saves.
- This Spec evolves the authoring loop and never regresses it: no existing
  command changes behavior, and a Spec that passes today's gates must not be
  blocked by a new check that merely disagrees with its style.

## Open Questions

None.
