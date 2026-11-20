---
status: accepted
knowledge: true
created_at: 2026-11-17T21:30:00Z
updated_at: 2026-11-17T21:30:00Z
deprecated_at: null
superseded_by: null
---

# A Spec whose phase fell back keeps the fallback for the limited model until the Spec ends

Once a phase of a Spec fell back because a model's usage limit ran out, every
later Authoring Phase of that Spec whose declared model is that model starts
on its own fallback, and the next Spec starts on the declared models. The
Phase Record holds the fallback, with the model that ran out and the limit
message, so the choice outlives the session.

## Context

ADR-0427 makes a phase fall back when its declared model's limit runs out.
Asked on 2026-11-17 whether the next phase tries the declared model again, the
maintainer answered "Fica no reserva na spec e volta ao modelo declarado na
próxima spec". A limit of this kind lasts days: the Fable limit observed that
day reset five days later. The maintainer also changes accounts during a
week, so a fallback must not outlive the Spec.

This repository declares Fable for the Idea, the PRD and the TechSpec and
Opus for the Tasks, with Sonnet as the Tasks' fallback "só para fallback ter
algo". Moving every later phase to its fallback would put the Tasks of a Spec
on Sonnet because Fable ran out.

## Decision

**The scope is the Spec and the model.** A later Authoring Phase of the same
Spec, in the same repository, starts on its fallback when its declared model is
the model whose limit ran out, and on its declared model otherwise. A phase
that declares no fallback starts on its declared model.

**The record.** The fallback is a new kind of entry in the Phase Record,
under the schema `roundfix/phase-record/v1`: the phase, the Spec, the
repository, the declared and fallback models and efforts, the name of the
limit that ran out and the limit message. The phase entry of the restarted
Author Subagent records the model that ran (ADR-0418). Readers of phase
entries skip the new kind.

**The next Spec.** Nothing carries a fallback to another Spec: its phases
start on the declared models.

## Considered Options

- **Try the declared model again at every phase.** Rejected by the
  maintainer. It would serve an account change in the middle of a Spec, and
  it costs a failed start while the limit lasts.
- **Every later phase on its fallback.** Rejected: a phase on a model that
  did not run out would be moved to a weaker fallback.
- **Until the reset time the message names.** Rejected: the time is text in
  the harness's language and zone, and an account change makes it wrong.
- **Kept in the session only.** Rejected: a session can end between phases,
  and the record is where a phase's model already is.

## Consequences

- A Spec authored across two machines falls back only on the machine whose
  record holds the fallback.
- A Spec whose limit cleared, or whose account changed, keeps its fallback
  until the Spec ends.
- The Run Retention Sweep removes old fallback entries with their day files.
