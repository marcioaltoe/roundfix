# Issue tracker: Local (canonical)

Planning artifacts for this repo live as local markdown under `docs/specs/` — there is no external tracker. This is the default and canonical mode of the CONTEXT-driven spec workflow.

## Conventions

- One feature per directory: `docs/specs/<feature-slug>/`
- Artifacts: `_idea.md` (optional), `_prd.md`, `_techspec.md` (optional), `_tasks.md` (the dependency graph — dependencies live **only** here), and one `task_NN.md` per task
- Task status lives **only** in each `task_NN.md` frontmatter: `pending | in_progress | completed | failed`
- QA evidence lives in `docs/specs/<feature-slug>/qa/`
- Completed specs (all tasks done, QA passed) move to `docs/specs/_archived/<feature-slug>/` (via `archive-spec`)
- Which pipeline stages a given change runs through (idea/PRD/techspec/tasks) is routed per `docs/agents/spec-routing.md`

## When a skill says "publish to the issue tracker"

There is no external tracker: the task files written by `write-tasks` **are** the published issues. Nothing further to do.

## When a skill says "fetch the relevant ticket"

Read the `task_NN.md` file in the spec folder. The user will normally pass the spec slug or the task file path directly.
