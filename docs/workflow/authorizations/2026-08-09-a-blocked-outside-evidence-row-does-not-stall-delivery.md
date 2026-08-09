# Tooling authorization — a blocked outside-evidence row does not stall delivery (2026-08-09)

CodeRabbit's review of pull request #143 raised issue_027 (major): the
task-authoring contract and the QA contract disagree about what happens when a
Spec's outside evidence cannot be obtained.

- `.agents/skills/write-tasks/SKILL.md` said the obligation "never requires
  human interaction and never stalls the Spec".
- `.agents/skills/qa-gate/SKILL.md` marks such a row `partial` without
  equivalent evidence, and `partial` blocks pull request preparation.

Both cannot be true. A resolve Batch reconciled them on its own by making the
block explicit in both files, which reversed the authoring clause. That is a
normative decision about ADR-0104's territory, so it was surfaced rather than
kept. Asked which direction holds, the maintainer answered:

> Ajuste o qa-gate e a adr

## What this covers

The authoring clause stands: an unobtainable outside source never stalls the
Spec. The QA contract moves to match it — a blocked outside-evidence row is
recorded with its reason and does not by itself block pull request preparation —
and ADR-0104 states that consequence so the two skills read from one decision
rather than from each other.

This is the direction that keeps the obligation honest without giving a third
party a veto over delivery. ADR-0104 exists because a Spec that grades its own
homework passes; it does not exist to let an unreachable source hold a finished
Spec hostage.

## Authorized paths

- `.agents/skills/qa-gate/SKILL.md`, limited to the outside-evidence row's
  effect on the verdict and on pull request preparation.
- `docs/adr/0104-a-spec-accepts-on-evidence-it-did-not-author.md`, limited to
  stating that a blocked row is recorded and does not stall delivery.
- `.agents/skills/write-tasks/SKILL.md`, limited to restoring the clause the
  resolve Batch reversed.

`.claude/skills/` is a symbolic link to `.agents/skills/`, which is the
authoritative source. The generated copies under `skills/` are rewritten by
`make skills-sync`, and derived Baseline pins by `make baseline-digests`; both
are sanctioned fallout under ADR-0081, not separate targets.

## Bounded by purpose

This grant covers the outside-evidence contradiction and nothing else. It does
not authorize other changes to either skill, to any other skill, or to any other
ADR.

## Consuming work

Resolution of CodeRabbit issue_027 on pull request #143.

## Commit choreography

This record lands as its own commit, before the commit that changes the skills
and the ADR.
