# Issue tracker: Local (canonical)

Planning artifacts for this repo live as local markdown under `docs/specs/` — there is no
external tracker. This is the default and canonical mode of the CONTEXT-driven spec workflow.

## Conventions

- One feature per directory: `docs/specs/<feature-slug>/`
- Artifacts: `_idea.md` (optional), `_prd.md`, `_techspec.md` (optional), `_tasks.md` (the
  dependency graph — dependencies live **only** here), and one `task_NN.md` per task
- Task status lives **only** in each `task_NN.md` frontmatter: `pending | in_progress | completed | failed`
- QA evidence lives in `docs/specs/<feature-slug>/qa/`
- Shipped specs move to `docs/specs/_archived/<feature-slug>/` (via `archive-spec`)

## When a skill says "publish to the issue tracker"

There is no external tracker: the task files written by `write-tasks` **are** the published
issues. Nothing further to do.

## When a skill says "fetch the relevant ticket"

Read the `task_NN.md` file in the spec folder. The user will normally pass the spec slug or the
task file path directly.

## Knowledge workspace

In this repo `docs` is a symlink into `.knowledge/` — spec artifacts physically live in the
central knowledge repository. Commit them with `git -C .knowledge …` per the
`knowledge-workspace` skill, never in the code repository.

## Legacy planning artifacts

`docs/plans/`, `docs/handoffs/`, and `docs/_inbox/` predate the spec workflow and are read-only
history. New work always goes to `docs/specs/<feature-slug>/`.
