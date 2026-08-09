---
status: done
created_at: 2026-08-08
updated_at: 2026-08-08
kind: finding
spec: 0084-an-update-that-can-run
---

# The update refuses six of the eight copies it exists to update (2026-08-08)

`roundfix baseline update`, the command Spec 0082 shipped so a maintainer could
refresh a repository's Context-Driven Baseline without re-answering setup
questions, was run in read-only mode against every local checkout carrying
`docs/agents/setup-context.json`. Two reach a plan. Six stop before planning
anything, and none of the six refusals is recoverable by the maintainer, because
each fires before the plan that would repair it exists.

Binary `0.4.0` built from `92867d41` on branch `ma/specs-0082-0083`, run as
`roundfix baseline update --repo <checkout> --format json` with no `--yes`, so no
checkout was mutated.

| repository | exit | state | blocker |
| --- | --- | --- | --- |
| gss | 3 | `plan_ready` | none; plan ready for approval |
| oraculum | 3 | `plan_ready` | none; plan ready for approval |
| roundfix | 3 | `action_required` | `managed-marker.modified` on four files |
| fiscus | 3 | `action_required` | `managed-marker.modified` on `docs/agents/autonomous-work.md` |
| conexus | 3 | `action_required` | fourteen unaccounted structural clauses |
| tax-poc | 3 | `action_required` | the same fourteen |
| vortex | 3 | `action_required` | the same fourteen |
| fluxus | 2 | `failed` | Baseline Profile `oraculum-backend` does not resolve; `.roundfix` absent |

The fourteen are identical across the three repositories that report them: six
`clause.backend.*`, two `clause.domain.*`, two `clause.frontend.*`, two
`clause.spec.*`, `rule.backend.boundary-contracts`, and
`rule.monorepo.context-boundaries`.

## The marker block is testing the wrong thing

Detection compares each managed region's bytes against the digest the Setup
Manifest recorded **on adoption day**. In this repository
`docs/agents/setup-context.json` was written once, in pull request #36, and never
again; `docs/agents/autonomous-work.md` changed in seven later pull requests and
`docs/agents/docs-layout.md` in four. The manifest describes bytes that stopped
existing, and the command reads that as damage.

The experiment that separates the two hypotheses: with the marker check disabled
behind an environment variable in a probe build, the same repository produced a
ready plan, applied it, and verified. The entire diff across six files was the
three capture clauses and the rewritten permission clause in
`docs/agents/secondbrain.md`, the `loop-07` clause in
`docs/agents/autonomous-work.md`, one rewritten trigger line in
`docs/agents/skill-dispatch.md`, and the republished manifest. No authored line
was lost: every removal in the diff was the previous version of a rewritten
line. The bytes on disk were a legitimate rendering of an earlier catalog, not a
human edit.

## The block is cold-start, not steady-state

After that apply, a second run — including one with the shipped binary carrying
the unmodified check — answered:

```
state: current | the repository already matches the current Baseline catalog
fileChanges: 0
```

The check is therefore wrong exactly once: on the first update of any repository
whose managed regions moved before the update command existed, which is every
repository that predates Spec 0082. One successful apply republishes the manifest
and the detection is correct from then on.

## Why the gate that shipped it passed

Spec 0082's Task 02 declared *Case: refresh a copy with a hand-edited managed
marker; Observation: the command blocks*. The gate observed that the code did
what the requirement asked, and it did. The requirement was wrong, and nothing in
the Spec measured against a repository the Spec had not built. The QA report
recorded zero blocked rows.

## Suggested action

Classify rather than block: a managed region whose bytes differ from the recorded
digest is unrecorded, planning proceeds, and the presented plan names the region
and every line the refresh removes. Keep blocking only the case with no
defensible target — the same managed identity appearing twice in one file.
