# Issue tracker: local

Planning artifacts live as local markdown under `docs/specs/`. There is no
external planning tracker.

## Conventions

- One feature per directory: `docs/specs/<feature-slug>/`.
- Artifacts: `_idea.md` (optional), `_prd.md`, `_techspec.md` (optional),
  `_tasks.md`, and one `task_NN.md` per Task.
- Dependencies live only in `_tasks.md`.
- Task status lives only in each `task_NN.md` frontmatter:
  `pending | in_progress | completed | failed`.
- QA evidence lives in `docs/specs/<feature-slug>/qa/`.
- Completed Specs move to `docs/specs/_archived/<feature-slug>/`.
- Pipeline routing lives in `docs/agents/spec-routing.md`.

The `archive-spec` skill owns workflow archiving. Roundfix's
`roundfix archive <feature-slug>` command enforces the same Task and QA
eligibility contract when used directly.

## When a skill says "publish to the issue tracker"

The task files written by `write-tasks` are the published issues. Nothing else
must be created.

## When a skill says "fetch the relevant ticket"

Read the requested `task_NN.md` file in the Spec folder.

## Repository ownership

This repository stores `docs/` directly in the code repository. Commit
documentation through the normal Roundfix repository workflow.

`docs/handoffs/` and `docs/_inbox/` are supporting documentation, not planning
trackers. New planned work belongs under `docs/specs/<feature-slug>/`.
