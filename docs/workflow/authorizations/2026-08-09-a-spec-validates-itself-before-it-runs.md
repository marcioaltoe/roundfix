# Tooling authorization — a Spec validates itself before it runs (2026-08-09)

On 2026-08-09, after the Spec 0090 QA gate returned `fail` on two authoring
defects, the maintainer proposed moving artifact validation out of the gate and
into authoring:

> Acredito que precisamos de uma skill que valide o adr, prd, techspec e tasks
> antes de iniciar o run de implementação, o ideal seria ate no momento de
> criação de cada etapa. Validamos o spec antes de finalizar e Assim mantemos o
> qa-gate focado no resultado esperado.

Shown that `roundfix spec check` already carries nineteen of those checks and
runs in 0.04 seconds for one Spec, the maintainer scoped the immediate step:

> Se o roundfix spec check ja faz isso, que seja rodado ao final da construção
> das tasks para validar a techspec. Assim já resolve uma parte do problema

## What this covers

`write-tasks` already instructs the author to "verify mechanically — parse,
don't eyeball" before reporting, and then lists checks an agent performs by
reading. The repository owns a command that performs them: `roundfix spec check
<slug>`. This grant adds that command as the mechanical half of that step, and
makes an unresolved finding block the report.

Measured on 2026-08-09, running it during authoring caught `SC-ADR-UNLISTED` on
Spec 0090 and `SC-ADR-RELATED` on Specs 0091 and 0092 — three defects that would
otherwise have reached a gate. It does not catch every class: Spec 0090's F-001,
a PRD citing an ADR that does not say what the PRD claims, has no check yet.

## Authorized paths

- `.agents/skills/write-tasks/SKILL.md`, limited to adding `roundfix spec check
  <slug>` to the graph-validation step and requiring a clean result before the
  breakdown is reported.

The generated copy under `skills/write-tasks/SKILL.md`, rewritten by
`make skills-sync`, is sanctioned fallout under ADR-0081 rather than a separate
target.

## Bounded by purpose

This grant covers running the existing checker at the end of Task authoring. It
does not authorize changes to any other skill, to the checker's rules, or to
what `write-tasks` decides about decomposition. The wider redesign — validation
at every authoring stage, the missing checks, and the QA gate shedding its
governance rows — is Spec 0093 and needs its own authorization.

## Consuming Spec

Applied directly rather than through a Spec: it is one paragraph in one file,
and it makes the authoring of Spec 0093's own Task graph safer.

## Commit choreography

This record lands as its own commit, before the commit that changes the skill.
