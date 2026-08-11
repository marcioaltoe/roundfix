# Constraint and ancestry audit

Build: `a51a94cb7773639b96fd4b081a1b78584faab0a5`.

## Graph and prerequisite state

- `_tasks.md` names `qa: task_06`.
- `task_06` has `type: qa`, is the terminal node, and depends on `task_05`,
  the graph's only non-QA leaf.
- Tasks 01–05 are `completed`; task_06 remains Daemon-owned `pending` while
  this gate runs.
- The PRD has no `## Unreachable Acceptance` declarations.

## Project Constraints

- Identifier strategy: the PRD correctly treats project-owned Internal
  Identifiers as not applicable. The TechSpec explicitly accounts for the new
  diagnostic `SC-*` names and explains that they are not lifecycle-bearing
  Internal Identifiers. Both cite `docs/agents/domain.md`.
- Authentication and HTTP: both artifacts record this as not applicable and
  cite `docs/agents/agent-instructions.md`; the feature reads local artifacts
  and opens no transport.
- Active ADR obligations: both artifacts account for accepted ADR-0080,
  ADR-0081, ADR-0091, and ADR-0093 with their roles. Fresh reads confirmed all
  four are accepted and unsuperseded.
- Tooling authority: the TechSpec names the protected Skill paths exactly and
  cites the 2026-08-02 and 2026-08-04 authorization records. The PRD does not
  satisfy the stricter artifact contract; see F-001 below.

## Authorization chronology and changed paths

Fresh `git merge-base --is-ancestor` checks exited 0 for the authorization
commits before the protected-tooling Task commits. The Spec-authoring commit
`4d796ed2` also precedes both tooling-bearing commits.

- Task 01 commit `bd38544d` changes the repository guide, the Baseline module
  product asset, its assigned Task file, and deterministic Baseline fallout.
  The 2026-08-05 authorization addendum now classifies module assets as product
  content, but the recorded 2026-08-04 boundary also predates this commit.
- Task 05 commit `a51a94cb` changes only
  `.agents/skills/{roundfix,write-tasks}/SKILL.md`, their corresponding
  `skills/` mirrors, its assigned Task file, and deterministic Baseline and
  characterization fallout. `git diff-tree --no-commit-id --name-only -r`
  supplied the inventory; no unlisted protected-tooling path appears.
- No prerequisite repair or consequent fix is folded into either Task commit,
  and no later Spec 0065 implementation commit follows task_05.
- The named sanctioned regeneration and characterization commands are present
  in both Task Results. Fresh `make verify`, `make skills-sync-check`, Skill
  check, and mirror comparisons pass on the generated state.

## F-001 — PRD does not name the exact protected Skill paths

Impact: Trust-Damage.

Actor and step: a maintainer audits the active PRD before permitting Task 05's
protected-tooling mutation.

Expected: `docs/agents/spec-routing.md` requires the active PRD and present
TechSpec to record express authorization and the exact bounded
repository-relative files.

Actual: `_prd.md` says only `bounded files: the owned skill pair`. It neither
names `write-tasks` and `roundfix` nor their canonical and mirror paths. The
TechSpec supplies those exact paths, but the rule requires both artifacts.

Reproduction:

1. Read `_prd.md` Project Constraints, Tooling authority.
2. Compare its phrase `the owned skill pair` with
   `.agents/skills/write-tasks/**`, `skills/write-tasks/**`,
   `.agents/skills/roundfix/**`, and `skills/roundfix/**` in the TechSpec.
3. Read `docs/agents/spec-routing.md:62`, which requires the exact bounded
   paths in both artifacts.

The public `spec check` returns no finding because its current
`SC-TOOLING-UNBOUNDED` parser accepts any non-empty text after `bounded files:`;
that mechanical silence does not satisfy the final QA audit.

Affected row: R02.
