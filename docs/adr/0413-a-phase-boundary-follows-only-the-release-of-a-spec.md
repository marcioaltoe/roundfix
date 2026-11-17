---
status: accepted
created_at: 2026-11-17T10:00:00Z
updated_at: 2026-11-17T10:00:00Z
deprecated_at: null
superseded_by: null
---

# A Phase Boundary follows only the release of a Spec, and the Boundary Notice prints only after a release

A Phase Boundary follows the delivery and release of a Spec, or its merge
when it releases nothing, and never an Authoring Phase. Between the Idea,
PRD, TechSpec and Tasks phases the Supervisor session goes on to the next
phase, and each authoring skill still hands its phase to the Author
Subagent where its `Declared model` section says so. `roundfix phase start`
prints the Boundary Notice only when the same session's latest phase entry
is `release`. This amends ADR-0407, which put a boundary after every
Authoring Phase.

## Context

ADR-0407 ended the Supervisor session after each Authoring Phase so that a
long orchestration session stopped carrying a large context into every
call. On 2026-11-16, the evening Roundfix 0.146.0 shipped it, the
maintainer wrote: "Esse negocio de criar sessões a cada etapa está
atrasando demais." He then chose "Só por Spec (Recommended)": the session
ends only after the Spec is delivered and released, the Authoring Phases
run in agents with their own context, and the Phase Record keeps
measuring.

The operator's Phase Record of 2026-11-16 to 2026-11-17 holds eight phase
entries of one session that authored three Specs. Roundfix 0.147.0 printed
the Boundary Notice for four of them, each after an Authoring Phase and none
after a release. The Author Subagents already give each feature phase its
own context (ADR-0401), and a subagent's command reads its parent's session
identifier, so the record cannot tell the phases a subagent ran from those
the session ran.

## Decision

**The boundary.** The Phase Boundary of a Spec follows its delivery and
release, or its merge when the Spec releases nothing. There the Supervisor
writes the Session Handoff, says what stays armed and starts no further
phase in the session, as ADR-0407 decided. An Authoring Phase is followed by
none: the Author Subagent returns its report, and the session goes on to the
next phase. The fix or refactor route keeps authoring in the session
(ADR-0408).

**The notice.** `roundfix phase start` marks the new entry and prints the
Boundary Notice when the same session's latest phase entry of the last seven
days is `release` and the new entry is not a `release` of the same Spec. Any
other earlier phase marks nothing. The notice keeps its text, its stream and
its refusal of nothing, and the record keeps its schema and its fields.

**Where it is stated.** The `autonomous-work` Baseline clause that says
when a session ends, the four authoring skills' `Phase boundary` sections,
the Roundfix skill's `--follow-up` and `--unattended` modes and its phase
reference, the glossary and the user guide state the same boundary.

## Considered Options

- **A boundary when the session changes Spec.** Rejected: one Supervisor
  session plans several Specs through Author Subagents, which record under
  its session, so the notice would fire on work the maintainer accepted.
- **No notice at all.** Rejected: a session that goes on after a release
  is the case the boundary still names, and the notice is its only
  reminder.
- **A boundary after `delivery`.** Rejected: one release often carries
  several delivered Specs, and the session that follows them is the one
  that cuts that release.
- **Delegating the fix route's phases too.** Not taken: ADR-0408 keeps a fix
  on the session's model, and the maintainer's answer did not reopen it.

## Consequences

- A session that plans and delivers a Spec writes one Session Handoff, at
  the end, instead of one per phase.
- A session that goes on after a release still gets the notice. A Spec
  that releases nothing ends its session at its merge, which the record
  does not see, so that boundary stays an instruction without a notice.
- Entries written before this decision keep the `boundary` mark the earlier
  rule gave them; a reader of the record tells them apart by the version
  that wrote them or by their date.
- The usage of the Authoring Phases is still attributable from the Phase
  Record and the harness transcripts, as ADR-0407 says.
