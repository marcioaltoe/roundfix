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

Asked to be precise, the maintainer clarified the direction:

> O que disse foi para manter o write-tasks com a alteração e ajustar o qa-gate
> e adr

So the resolve Batch's authoring text stands: decomposition never stalls and
never asks a human, it records the blocked row and proceeds, and the row is
carried into the QA gate. The `qa-gate` skill and ADR-0104 move to say the same
thing in the same words — the gate holds pull request preparation until the row
is satisfied or carried forward on declared unmoved evidence under ADR-0097.

The obligation therefore lands at the gate rather than during authoring. That is
what keeps ADR-0104 from being satisfiable by silence: a Spec may reach its gate
with the row open, but it may not ship past one.

## Authorized paths

- `.agents/skills/qa-gate/SKILL.md`, limited to the outside-evidence row's
  effect on the verdict and on pull request preparation.
- `docs/adr/0104-a-spec-accepts-on-evidence-it-did-not-author.md`, limited to
  stating where the obligation lands and what it blocks.
- `.agents/skills/write-tasks/SKILL.md`, limited to the outside-evidence clause
  the resolve Batch rewrote.

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
