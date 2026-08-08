# Tooling authorization — the roundfix skill ships with the CLI change (2026-08-08)

On 2026-08-08 the maintainer directed the managed-refresh defect to be fixed
before the merge, by a Spec that replaces the behavior and the process:

> Vamos consertar antes do merge, refazendo o trabalho com uma nova spec que
> substitua e ajuste o que deveria estar funcionando e ser seguido.

## What this covers

`docs/agents/specific-repository.md` carries a HARD RULE: a pull request that
changes CLI behavior ships the roundfix skill update with it. Spec 0084 changes
observable `roundfix baseline update` behavior — a repository whose managed
regions diverge from the recorded digests now reaches a ready plan instead of an
action-required state, and the presented plan names each unrecorded region and
the lines the refresh removes.

The skill describes that command in detail today, including its flags, its exit
codes, and the sentence "Managed refresh never invokes semantic classification
and preserves every non-managed byte exactly". That sentence stays true; what is
missing is the classification and the removed-line report, which an Agent reading
the skill would otherwise not know exists.

The gap is measured. Task 03 of Spec 0084 declared the skill update as a
requirement and the implementing Agent refused it, correctly, because no
authorization named the path — it recorded the blocker and changed nothing rather
than widening its own boundary. This record removes the blocker without widening
anything else.

## Authorized paths

- `.agents/skills/roundfix/SKILL.md`, limited to describing the managed-refresh
  classification, the unrecorded-region report, and the removed-line report in
  the section that already documents `roundfix baseline update`.

`.claude/skills/` is a symbolic link to `.agents/skills/`, which is the
authoritative source. The generated copies under `skills/` are rewritten by
`make skills-sync` and are sanctioned fallout, not separate targets.

## Bounded by purpose

This grant covers the update command's own description. It does not authorize
changing the skill's Run, resolve, watch, implement, settle, reconcile, archive,
or release guidance, its agent bundles under `.agents/skills/roundfix/agents/`,
or any other skill in the repository.

## Sanctioned fallout — no separate grant

Skill digests and generated copies rewritten by `make skills-sync` and
`make baseline-digests` are deterministic consequences of the authorized source
edit, per ADR-0081. A hand-edited digest remains an unauthorized mutation.

## Consuming Spec

This authorization is consumed by Spec `0084-an-update-that-can-run`.

## Commit choreography

This record lands as its own commit, before the commit that changes the skill.
