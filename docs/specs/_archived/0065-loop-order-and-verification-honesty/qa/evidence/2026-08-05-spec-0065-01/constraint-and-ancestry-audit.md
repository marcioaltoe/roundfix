# Constraint and ancestry audit

Build: `d603031e808e3c64a539c4875f00d62ed778f522`.

## Graph and prerequisite state

- `_tasks.md` names `qa: task_06`; `task_06` has `type: qa`, is terminal,
  depends on `task_07`, and task_07 is the graph's only non-QA leaf.
- Tasks 01–05 and corrective task_07 are `completed`; task_06 remains
  Daemon-owned `pending` during this gate.
- The PRD has no `## Unreachable Acceptance` declaration.

## Project Constraints

- Identifier strategy: the PRD records no project-owned Internal Identifier;
  the TechSpec separately accounts for `SC-*` diagnostic names and their
  existing non-lifecycle convention. Both cite `docs/agents/domain.md`.
- Authentication and HTTP: both active artifacts record this as not
  applicable because the feature reads local Spec artifacts and opens no
  transport. Both cite `docs/agents/agent-instructions.md`.
- Active ADR obligations: both artifacts account for accepted, unsuperseded
  ADR-0080, ADR-0081, ADR-0091, and ADR-0093 with their operative roles and
  cite `docs/agents/domain.md`.
- Tooling authority: both active artifacts now name the protected paths
  exactly: `.agents/skills/write-tasks/**`, `skills/write-tasks/**`,
  `.agents/skills/roundfix/**`, and `skills/roundfix/**`. They cite the
  2026-08-02 grant and 2026-08-04 confirmation. The Baseline module asset is
  also named; the 2026-08-05 authorization addendum classifies module assets
  as product content, while the earlier confirmation independently records
  the requested boundary.

The prior report's F-001 is closed: commit `b4347ac6` replaced the PRD's vague
`owned skill pair` phrase with those four exact repository-relative path
families before corrective Task 07 ran.

## Authorization chronology and changed paths

Fresh `git merge-base --is-ancestor` checks all exited 0:

- authorization commits `2e560cea` and `9bdaedbe` and Spec-authoring commit
  `4d796ed2` precede Task 01 commit `bd38544d`;
- the same three commits precede Task 05 commit `a51a94cb`;
- corrective authoring commit `b4347ac6` precedes Task 07 commit `d603031e`.

Fresh `git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r`
inventories show:

- Task 01 changes the repository guide, Baseline module product asset, its
  assigned Task file, and regenerated formatter, profile, digest, diagnostic,
  and characterization outputs.
- Task 05 changes only the two authorized canonical Skills, their mirrors, its
  assigned Task file, and sanctioned Baseline digest and characterization
  fallout.
- `b4347ac6` changes only the PRD, Task Graph, QA Task, and new corrective Task.
- Task 07 changes the two authorized canonical Skills and mirrors, the guide,
  Baseline module product asset, TechSpec and assigned Task, plus sanctioned
  formatter, profile, digest, diagnostic, parity, and characterization
  fallout. It changes no Go source.

No prerequisite repair or consequent fix is folded into a tooling Task
commit. The Task Results name `make skills-sync`, `make baseline-digests`, and
the two sanctioned characterization re-records that produced the derived
files. At audit time the only worktree delta was this in-progress QA report;
no implementation or protected-tooling delta was present.
