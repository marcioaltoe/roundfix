---
status: accepted
created_at: 2026-10-08T00:00:00Z
updated_at: 2026-10-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Inside a Run the Daemon is the only full gate, and the Baseline states each rule once

The Baseline audit of 2026-10-08 (`docs/references/2026-10-08-baseline-audit.md`)
measured Specs 0225–0248 in the Run Database. Task agents ran
`make verify-incremental`, the full `go test` suite, 74 times in 20 of 23
Specs: 284 agent-minutes, 231 s each on average. They did so because the
Baseline told them to run the incremental tier "for each Task ... before
handoff". The `implement-task` skill says the opposite for a Daemon-assigned
turn, and the Daemon already runs the Task's Verification and
`make verify-changed` (350 s on average) at every settlement. That Daemon gate
failed 6 times in 125, so the agent-side run caught almost nothing it would
not have caught. The audit also found guide text that contradicts the Project
Config (`rtk`-prefixed Verification values, model names that no profile uses)
and rules stated two or three times. On 2026-10-08 the maintainer asked to
"Atualizar e verificar o baseline de regras para que tenhamos o máximo de
performance e resultado com o roundfix e congruente com as mudanças".

**The tiers are scoped by session.** Outside a Run, an Agent runs the selected
incremental Verification for fast checks and the selected repository
Verification before a completion claim, as before. In a Daemon-assigned turn
the Agent runs focused tests of the packages it changed and neither selected
command, because the Daemon runs the Task's declared Verification and the
repository Verification at settlement. The core and Spec workflow clauses say
so, and `implement-task` says the same. Before a Run starts, the author runs
`roundfix spec check <slug> --strict --run-verification`, because only strict
mode fails on gap findings.

**Recorded values match the gates that run.** This repository records
`verification.gate: make verify`, which CI runs on every push to main, and
`verification.incremental: make verify-changed`, which CI runs on every Pull
Request and the Daemon runs at settlement. The catalog stops proposing an
`rtk`-prefixed command, because `rtk` is a local output filter, not part of any
repository's gate.

**The runtime decisions stay, and the guide says who decides.** Retiring
`runtime.backend` and `runtime.design` would make every adopter answer or drop
a recorded decision, so both stay (operator decision, 2026-10-08). Their
defaults become the tuples the Project Config prefers, `codex gpt-6.1-sol high`
and `claude opus high`, and the autonomous-work guide states that Agent
Selection Profiles choose each Agent Session's runtime, model and effort
(ADR-0049) and that a low-complexity Task runs on the Light Tier
(ADR-0238). A recorded value is a stated preference, never a selection.

**One statement per rule.** When two clauses say the same thing, one absorbs
the other and lists it in `replaces`, so an adopter's update records the
removed clause `replaced` and never `unaccounted`. A clause whose wording
changes keeps its identity, as Spec 0239 did, and this narrows the rename rule
of ADR-0222: a new identity needs a Source Baseline row that the regenerator
never creates. An authoring prototype of this Spec measured that cost, with
eight missing rows and one retention transition left without a target.

**The root names what to read and when.** The domain guide stays mandatory for
every session. The docs-layout guide is read before work that creates, changes,
moves or retires `CONTEXT.md` or a document under `docs/`, or before a test or
build step reads one. The Secondbrain guide is read before consulting or
writing the Secondbrain or authoring an Idea, PRD or TechSpec. A Run's Agent
Session cannot reach the Secondbrain. Together the two guides are about 23 KB
of the 69 KB the root points every session at.

## Considered options

- Keep the Baseline and change only `implement-task`. Rejected: agents followed
  the Baseline over the skill, and Gloaguen et al. (2026, arXiv:2602.11988)
  found that agents follow context-file instructions closely. Those
  instructions cause more testing and exploration, and inference cost rises
  20 % or more with no gain in success.
- Rename every reworded clause, as ADR-0222 asks. Rejected for the cost above.
- Remove the runtime decisions. Rejected because of the adopter impact.

## Consequences

An adopter's next `roundfix baseline update` rewrites its managed guides and
root blocks with these clauses. Its recorded decision values stay what it
recorded. The audit estimates 5–10 minutes of wall time saved per Spec from the
tier scoping alone. That number rests on the audit's 284 agent-minutes over 23
Specs, not on a measurement this ADR makes. The first Runs after release are
the place to confirm it.
