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

<!-- setup-context-driven:begin id=guide.issue-tracker version=0.0.1 -->

# Issue tracker

- **mandatory**: Keep each Spec under `docs/specs/<feature-slug>/`. Dependencies live only in `_tasks.md`; status lives only in each Task file, and the local Task files are the published planning issues.

<!-- setup-context-driven:end id=guide.issue-tracker -->

<!-- roundfix:repository-rule:begin id=rule.70eb8ff15c03ef56d724a61c44c3d479f91912f693b68f832c0f2fea799924dd -->
## Agent skills

### Issue tracker

Tasks live as local markdown under `docs/specs/<feature-slug>/` (the canonical
source — no external tracker). See `docs/agents/issue-tracker.md`.


<!-- roundfix:repository-rule:end id=rule.70eb8ff15c03ef56d724a61c44c3d479f91912f693b68f832c0f2fea799924dd -->

<!-- roundfix:repository-rule:begin id=rule.9754a69cf216a520e424e4705dfebf6dcf22fd8d021b243c967598765b9622f8 -->
2. Marking a spec task `completed` without fresh verification evidence

<!-- roundfix:repository-rule:end id=rule.9754a69cf216a520e424e4705dfebf6dcf22fd8d021b243c967598765b9622f8 -->

<!-- roundfix:repository-rule:begin id=rule.64fe5b9bda0a948e5a23375064fcdf9158651a35d555b973a2fc92224622a7e5 -->
3. Tracking progress in `_tasks.md` — status lives only in each `task_NN.md`

<!-- roundfix:repository-rule:end id=rule.64fe5b9bda0a948e5a23375064fcdf9158651a35d555b973a2fc92224622a7e5 -->
