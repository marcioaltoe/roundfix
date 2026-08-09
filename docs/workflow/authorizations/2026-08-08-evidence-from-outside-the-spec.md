# Tooling authorization — acceptance evidence from outside the Spec (2026-08-08)

On 2026-08-08 the maintainer directed that the managed-refresh defect be fixed
by a Spec that repairs both the behavior and the process that shipped it:

> Vamos consertar antes do merge, refazendo o trabalho com uma nova spec que
> substitua e ajuste o que deveria estar funcionando e ser seguido.

## What this covers

Spec 0082 shipped a requirement that made its own command unrunnable on the
first update of every repository that existed before the command did. Its QA
passed with zero blocked rows.

It passed because the rubric and the requirement had the same author and the
same premise. Task 02 required that a hand-edited managed marker block; its
Rehearsal Case read *Case: refresh a copy with a hand-edited managed marker;
Observation: the command blocks*. The gate observed that the code did what the
requirement said, which it did. Nothing in the Spec asked whether the
requirement was right, and nothing forced a measurement against a repository the
Spec had not built.

The Secondbrain names this failure mode directly.
`wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md` states that a
gate is trustworthy only when there is evidence it *observed the right property*
and can *fail a known negative case*, and reports that in 3,730 replayed events
46.0% of comparable positive results carried no information discriminating the
bug from an unrelated pass. A rehearsal that confirms the design is one of those
non-discriminating positives.

The correction that would have caught it is cheap: at least one acceptance row
per Spec must rest on evidence the Spec did not author — a real repository, a
measurement, or published literature. Running the finished 0082 command against
eight adopted repositories took minutes and produced the diagnosis the entire QA
missed.

## Authorized paths

- `internal/baseline/assets/modules/spec-workflow.json`, limited to adding one
  clause to `rule.spec.project-constraints` and bumping that rule's version.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`
  and the same Source Baseline's `manifest.json`, limited to seating that
  clause.
- `.agents/skills/write-tasks/SKILL.md` and `.agents/skills/qa-gate/SKILL.md`,
  limited to stating the same obligation where each skill already contracts
  Rehearsal Cases and QA rows. `.claude/skills/` is a symbolic link to
  `.agents/skills/`, which is the authoritative source; the generated copies
  under `skills/` are rewritten by `make skills-sync` and are sanctioned
  fallout, not separate targets.

The Source Baseline paths are included from the start: the catalog refuses a
mandatory clause that no Source Baseline row carries, and the regenerator
maintains rows but never creates them.

The two skills are repo-owned authorial workflow skills, which
`docs/agents/skill-dispatch.md` permits adapting locally; they are named here
because their text is the operative contract a Task author reads.

## Bounded by purpose

The clause requires that a Spec's acceptance rest, in at least one named row, on
evidence originating outside the Spec's own artifacts, and that the row record
where that evidence came from. It must not require a specific number of rows,
mandate human interaction, or turn an unavailable external source into a Spec
blocker — a row whose external evidence cannot be obtained is recorded as
blocked with its reason, which the QA contract already supports.

It must not restate the existing QA audit obligation in
`clause.spec.project-constraints-04-qa-audit`, whose concern is that the audit
happens at all.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edits, per ADR-0081. The maintained Source Baseline
entry count moves with the new corpus entry and its fixture expectation moves
with it. Skill digests rewritten by the same command follow the same rule.

## Consuming Spec

This authorization is consumed by Spec `0084-an-update-that-can-run`.

## Commit choreography

This record lands as its own commit, before the commit that changes the module
or the skills.
