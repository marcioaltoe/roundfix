---
granted: 2026-08-27
action: Describe the Unresolved-Run carry-forward, the implement Preflight refusal that names it, and the two reconcile disposition flags the skill never documented.
consuming: 0118-a-task-proved-once-does-not-run-twice
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
---

# Tooling authorization — carry-forward reaches an Unresolved Run (2026-08-27)

On 2026-08-27 the maintainer was asked which scope Spec 0118 should record,
given that the Spec changes CLI behavior the roundfix skill documents while the
skill's `reconcile` section is already behind the shipped command. The
maintainer chose to include the skill under authorization and to close the
pre-existing gap in the same Spec.

## What this covers

`docs/agents/specific-repository.md` carries a HARD RULE: a pull request that
changes CLI behavior ships the roundfix skill update with it. Spec 0118 changes
two observable behaviors an Agent reading the skill would otherwise not know
exist.

The first is `roundfix reconcile --carry-forward`, which today refuses every Run
that is not `Stopped` and after this Spec also accepts an `Unresolved` one. The
second is a new `implement` Preflight refusal: a Run that would re-execute Tasks
already completed and verified on a prior Run Branch stops before creating a Run
and names the carry-forward command instead.

The gap is measured and it predates this Spec. The skill documents `reconcile`
with `--apply` alone at three call sites and mentions neither `--carry-forward`
nor `--discard-superseded`, both shipped by Spec 0092 on 2026-08-11. An Agent
following the skill therefore cannot reach the command this Spec exists to make
reachable, which would leave the Spec's own remedy undiscoverable through the
documented surface. Closing that gap is what makes the rest of the Spec usable,
so it is bounded here rather than deferred.

## Authorized paths

- `.agents/skills/roundfix/SKILL.md`, limited to the sections that already
  document `roundfix reconcile` and `roundfix implement`: adding
  `--carry-forward` and `--discard-superseded` to the reconcile guidance,
  recording the accepted Run states for carry-forward, and describing the
  implement Preflight refusal with the command that clears it.

`.claude/skills/` is a symbolic link to `.agents/skills/`, which is the
authoritative source. `skills/roundfix/SKILL.md` is the generated copy that
`make skills-sync` rewrites from it. It is sanctioned fallout of the authorized
source edit rather than a separate target, and it is named in the bounded paths
above because the changed-path audit reads that list and not this sentence — a
grant whose prose says more than its frontmatter is a grant the gate refuses.

## Bounded by purpose

This grant covers the reconcile and implement guidance named above. It does not
authorize changing the skill's resolve, watch, settle, archive, baseline, or
release guidance, its agent bundles under `.agents/skills/roundfix/agents/`, or
any other skill in the repository. It does not authorize changing what the
tooling-authority rule requires, and it does not widen any detector.

## Sanctioned fallout — no separate grant

Skill digests and generated copies rewritten by `make skills-sync` and
`make baseline-digests` are deterministic consequences of the authorized source
edit, per ADR-0081. A hand-edited digest remains an unauthorized mutation.

## Consuming Spec

This authorization is consumed by Spec
`0118-a-task-proved-once-does-not-run-twice`.

## Commit choreography

This record lands as its own commit, before the commit that changes the skill.
