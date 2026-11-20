---
status: accepted
knowledge: true
created_at: 2026-11-17T14:20:00Z
updated_at: 2026-11-17T14:20:00Z
deprecated_at: null
superseded_by: null
---

# A Supervisor session ends only when the maintainer asks or the work stops, and no phase marks a boundary

No phase of a Spec ends the Supervisor session, with or without a person
present: after an Authoring Phase, a delivery or a release the session goes
on with the next queued work, and it writes the Session Handoff and ends
only when the maintainer asks or the work stops because nothing remains or
what remains waits on the maintainer. `roundfix phase start` keeps recording
every phase, reads no earlier entry, marks none and prints no Boundary
Notice. This amends ADR-0407 and supersedes ADR-0413 of 2026-11-17, which
nothing had implemented.

## Context

ADR-0407 ended the Supervisor session after each Authoring Phase and after
a Spec's delivery and release, and `roundfix phase start` printed the
Boundary Notice when a session went on. On 2026-11-16 the maintainer wrote:
"Esse negocio de criar sessões a cada etapa está atrasando demais." and
chose a boundary only after a Spec's delivery and release, which ADR-0413
recorded on the morning of 2026-11-17. Before any Run implemented it, the
maintainer noticed that a session also stopped at the end of each Spec's
delivery asking to be restarted, which "vai prejudicar demais o autonomous
work", and answered on 2026-11-17 at about 10:00 BRT "Só quando você pedir
(Recommended)": no automatic stop at all, neither after a phase nor after a
Spec, with or without a person present; the handoff is written only when
the session really ends; the Phase Record keeps measuring; the notice of
`roundfix phase` goes.

The operator's Phase Record of 2026-11-16 to 2026-11-17 holds eight phase
entries of one session that authored three Specs, and Roundfix 0.147.0
printed the Boundary Notice for four of them. The Author Subagents already
give each feature phase its own context (ADR-0401), so a phase does not need
a new session to start small.

## Decision

**When a session ends.** A Supervisor session follows its queue: it hands an
Authoring Phase to its Author Subagent as the authoring skill's `Declared
model` section says, follows delivery and release, and goes on to the next
Spec. No phase of a Spec ends it. It writes the Session Handoff, with the
Spec, the phase that ended, the next phase, each open decision and what
stays armed, tells the maintainer what stays armed and ends only when the
maintainer asks or the work stops because nothing remains or what remains
waits on the maintainer. A turn that ends while work goes on still arms a
detached Run, the Delivery Queue or a scheduled wake-up, or says plainly
that work stops and what remains.

**The Phase Record.** `roundfix phase start` records each phase as before
and reads no earlier entry. It writes the `boundary` key of the
`roundfix/phase-record/v1` schema as `false` and never writes
`previous_phase` or `previous_spec`; an entry an earlier version marked
keeps its mark and stays readable. It prints nothing on standard error at a
phase change, and its standard output, refusals and exit codes are
unchanged.

**Where it is stated.** The `autonomous-work` Baseline clause that says
when a session ends, the four authoring skills' `Phase boundary` sections,
the Roundfix skill's `--follow-up` and `--unattended` modes and its phase
reference, the glossary and the user guide state the same rule.

## Considered Options

- **A boundary after a Spec's delivery and release** (ADR-0413).
  Superseded: it still stopped an unattended session at the end of every
  Spec, the stop the maintainer found harmful to autonomous work.
- **Keeping the notice for a session that goes on after a release.**
  Rejected: after this decision going on is the expected course, so the
  notice would only warn against the rule.
- **Dropping the `boundary` key from the record.** Rejected: an outside
  reader of `roundfix/phase-record/v1` may expect the key, and a constant
  `false` keeps the schema without a new version.
- **A stop after a number of Specs or a context size.** Not taken: the
  maintainer asked for no automatic stop, and the harness compacts or the
  maintainer asks when a session grows.

## Consequences

- One session can plan, deliver and release several Specs; its context
  grows with them, and the per-phase usage stays attributable from the
  Phase Record and the harness transcripts, as ADR-0407 says.
- A Session Handoff is written once, when the session ends, instead of once
  per phase or per Spec.
- Entries written by Roundfix 0.146.0 and 0.147.x keep `"boundary":true`
  where the earlier rules marked them; a reader that counts marks splits by
  date.
- An unattended cycle no longer waits at a phase for a session a person
  opens; it stops only when its queue is empty or it needs the maintainer.
