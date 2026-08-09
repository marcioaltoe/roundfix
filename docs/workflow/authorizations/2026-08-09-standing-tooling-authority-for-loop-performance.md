# Standing tooling authority — loop performance with reliability (2026-08-09)

On 2026-08-09, after seeing that Spec 0090's QA gate spent 461 tool calls and
two context compactions to find a defect two file reads would have caught, the
maintainer granted standing authority over protected tooling:

> Sobre autorização de editar o tooling, está totalmente autorizado para fazer
> isso com o foco de melhoria na performance garantindo a confiabilidade.

This record exists because a broad grant needs a *more* precise boundary than a
narrow one, not a looser one. Nothing below widens what was said; it writes down
what it means so a later reader can audit an edit against it.

## What this covers

Protected tooling — authored skills, `Makefile`, CI workflow files, ignore
files, plugin declarations, and version pins — may be edited when the change
serves **loop performance without reducing reliability**. Both halves are load-
bearing:

- *Performance* means fewer Agent turns, fewer tokens, less wall clock, or
  earlier discovery of the same defect. Moving a check from the QA gate to
  authoring qualifies. So does deleting a step that proves nothing.
- *Reliability guaranteed* means the change must not reduce what the loop can
  detect. A check may move, be made cheaper, or be replaced by a stronger one.
  A check may not be dropped because it is inconvenient, and a gate may not be
  made cheaper by making it blinder.

When the two conflict, reliability wins and the change stops.

## What it does not cover

- Any edit whose purpose is not performance-with-reliability. A feature, a
  preference, or a stylistic rewrite of a skill still needs its own grant.
- Weakening a Normative Clause, a HARD RULE, or a gate's fail-closed behaviour.
  Making a rule cheaper to satisfy is in scope; making it optional is not.
- Release, publication, or anything irreversible and outward-facing, which
  `docs/agents/autonomous-work.md` keeps with the maintainer regardless.
- Another repository. This grant is scoped to Roundfix.

## Obligations that travel with it

Every edit under this grant must, in its own commit message, name which of the
two halves it serves and how reliability was preserved — the check that moved
rather than vanished, the evidence that the new form catches what the old one
caught. An edit that cannot state that is outside the grant.

Where an edit removes a check from one place, the commit must show it running in
another, measured. Spec 0093 is the worked example: the QA gate's governance
rows leave only because `roundfix spec check` runs the same rules during
authoring, in 0.04 seconds, and the Spec's gate proves the corpus verdict did not
move.

## Consuming Specs

Open-ended. Consumed first by Spec 0093, and available to every later Spec whose
purpose matches. A Spec relying on it records that reliance in its Tooling
authority row and cites this record.
