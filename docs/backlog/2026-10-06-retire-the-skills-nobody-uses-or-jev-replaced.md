---
type: chore
status: open
created: 2026-10-06
spec: null
---

# Retire the skills nobody uses, or that the Jev judge replaced

## Problem

This repository installs 45 skills, 14 Roundfix-owned and 31 external. The
Baseline modules (`internal/baseline/assets/modules/*.json`) require most of
them from every adopter. The maintainer asked on 2026-10-06 to retire from the
Roundfix set and the Baseline the skills that no longer have a use, or whose
function the Jev judge took over, and named `council` and `the-fool`.

Usage was measured read-only on 2026-10-06 from the Run Database:
`agent.tool_started` events whose payload names `skills/<name>/`, across 220
Runs. Runs cover implementation; interactive authoring sessions are not in the
Run Database, so these counts understate skills meant for a human
conversation.

| Skill | Owner | Baseline module | Reads (Runs) | Role |
| --- | --- | --- | ---: | --- |
| `the-fool` | external | context-workflow | 0 (0) | devil's advocate, pre-mortem, red team |
| `grilling` | external | context-workflow | 0 (0) | interview the user to stress-test a plan |
| `grill-with-docs` | external | context-workflow | 0 (0) | interview that also writes ADRs and glossary entries |
| `autoresearch` | external | none (this repository only) | 0 (0) | tune a skill by repeated evals |
| `handoff` | external | core | 2 (1) | compact a conversation for another agent |
| `business-analyst` | owned | context-workflow | 11 (3) | score ideas, KPI targets |
| `council` | owned | context-workflow | 23 (3) | multi-advisor debate on high-impact decisions |
| `write-idea` | owned | context-workflow | 25 (4) | raw idea to `_idea.md` |
| `brainstorming` | owned | context-workflow | 87 (39) | explore before creative work |
| `domain-modeling` | external | context-workflow | 97 (25) | glossary and domain model |

For comparison, workflow skills in daily use: `implement-task` (201 Runs),
`evidence-gate` (196), `tech-writer` (193), `testing-boss` (159),
`systematic-debugging` (159), `qa-gate` (153), `no-workarounds` (145),
`knowledge-workspace` (136).

## What the Jev judge replaced, and what it did not

The Jev judge (`roundfix spec judge`, the `the_judge` skill) asks fixed, typed
questions about a Spec and returns probabilities. Examples: does a goal have
a mechanism, does a citation support its claim, which open item belongs with
which source, which tier a Task needs. That covers the "check this artifact"
part of `the-fool` (stress-test assumptions) and part of `business-analyst`
(scored viability). It does not generate alternatives, argue positions, or
interview a person, which is what `council`, `grilling` and `grill-with-docs`
do. Since the autonomous authoring route (write-prd → write-techspec →
write-tasks, decided autonomously and recorded as ADRs) took over Spec
authoring, those interactive skills have had no reader in the measured Runs.

## Proposal (non-binding; the maintainer decides per skill)

- **Retire:** `the-fool` (no use; its check role is covered by the judge and
  the pre-PR review), `autoresearch` (no use; not in any module), and
  `handoff` from `core` (one Run; Roundfix has its own handoff and history
  conventions).
- **Decide:**
  - `council`, `business-analyst` and `write-idea`: owned, rarely used, part
    of the human-led route before a PRD. Either retire them or keep them as
    optional skills outside the required modules.
  - `grilling` and `grill-with-docs`: unused by Runs. They are part of the
    CONTEXT-driven method credited to Matt Pocock, so retiring them changes
    the method's interactive entry point.
- **Keep:** the skills in daily use above, plus `brainstorming` and
  `domain-modeling` (regular readers).
- **Mechanics:**
  - Remove the retired skills from their Baseline modules, from
    `skills-lock.json` and the Setup Snapshot, and from the owned set (owned
    ones).
  - Update guides and skill-dispatch text that names them.
  - Adopters keep any installed copy until they remove it, as Spec 0208 did
    for `review` and `triage`.
  - Vendored upstream skills are not edited, only dropped from the set.
  - Baseline clause retention (rule 11) applies to any clause that names a
    retired skill.

## Notes

Resolve after the current cycle (0235–0237), alongside or after the history
clean-up entry, `2026-10-06-history-keeps-only-what-the-secondbrain-needs.md`.

## Maintainer decision — 2026-10-06

Asked through AskUserQuestion, the maintainer answered: "remover: the-fool,
autoresearch, council // manter: grilling, grill-with-docs, write-idea,
business-analyst e handoff". Scope therefore:

- **Retire:** `the-fool`, `autoresearch` and `council`. `council` is
  Roundfix-owned, so it leaves the owned set and its mirror too.
- **Keep:** `grilling`, `grill-with-docs`, `write-idea`, `business-analyst`
  and `handoff`. They stay in their modules.

The maintainer also asked to validate in the Baseline why `grilling` and
`domain-modeling` are not being executed and why `CONTEXT.md` is not kept
current during implementations. That is recorded separately in
`2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md`.

