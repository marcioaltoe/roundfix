# Issue tracker: local markdown

Work for this repo is tracked as markdown spec artifacts under `docs/specs/` — there is no external issue tracker. The local spec files are canonical.

## Conventions

- One feature per directory: `docs/specs/<feature-slug>/`
- The PRD is `docs/specs/<feature-slug>/_prd.md`
- The dependency graph is `_tasks.md`; dependencies live only there
- Work units are `task_NN.md` files, numbered from `01`; each task's status lives only in its own frontmatter
- QA evidence lives in `docs/specs/<feature-slug>/qa/`
- Shipped specs move to `docs/specs/_archived/<feature-slug>/`

## When a skill says "publish to the issue tracker"

Write the markdown file into the feature's `docs/specs/<feature-slug>/` directory, creating the directory if needed.

## When a skill says "fetch the relevant ticket"

Read the file at the referenced path. The user will normally pass the path or the task number directly.
