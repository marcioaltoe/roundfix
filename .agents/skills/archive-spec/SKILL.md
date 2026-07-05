---
name: archive-spec
description: Archive a completed spec — verify every task completed and QA passed, then stamp the archive metadata and move docs/specs/<slug>/ to docs/specs/_archived/<slug>/. Runs automatically at the end of the implement-spec loop after a QA pass, or whenever the user asks to archive a spec.
argument-hint: "<spec slug> [--release <tag or PR URL>]"
metadata:
  category: delivery
  tags: [workflow, documentation, process]
  version: 0.2.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Archive Spec

Move a completed spec out of the active set: `docs/specs/<slug>/` → `docs/specs/_archived/<slug>/`, with the completion stamped in its frontmatter. Archived means _implemented and verified_ — every task done and QA passed — after this, one `ls docs/specs/` separates live work from history, and the archive stays greppable as the record of what was built and why.

The trigger is spec completion, not publication: run this automatically at the end of the `implement-spec` loop once the QA gate passes, or whenever the user asks. Merge and release are separate, user-driven steps — the archive commit simply travels with the branch and ships inside the feature's own PR.

## Preconditions — verify, don't trust

Check both with fresh evidence before touching anything:

1. **Every task completed** — read each `task_NN.md` listed in `_tasks.md`; every `status` is `completed`.
2. **QA passed** — the newest report in `qa/` has `verdict: pass`. A missing `qa/` directory or a failing latest report blocks the archive; proceed only if the user explicitly says "archive anyway", and record that override in the stamped frontmatter (`qa_override: true`).

A merged PR or release tag is **not** a precondition. If the user passes `--release`, or a merged PR/tag is already known, stamp it as metadata — but never block the archive waiting for one.

If any check fails, stop and report exactly which — the spec stays active.

## Steps

1. **Stamp** `_prd.md` frontmatter:

   ```yaml
   status: archived
   archived: YYYY-MM-DD
   release: <tag or PR URL> # only when known — from --release or an already-merged PR/tag
   ```

2. **Move** with history preserved:

   ```bash
   mkdir -p docs/specs/_archived
   git mv docs/specs/<slug> docs/specs/_archived/<slug>
   ```

3. **Commit** — `chore(specs): archive <slug>` (Conventional Commits). Do not push unless asked.

4. **Report** — the new path, the release reference when one was stamped, and anything carried over (open follow-ups from task `## Result` sections belong in new specs, not in the archive).

5. **Suggest the publish step** — when the work isn't merged yet, close by suggesting the PR (via `github-pr-workflow`). This is where that suggestion lives in the workflow — the implement loop ends at the archive and doesn't offer it. Suggest only: opening the PR is the user's call.

## Unarchive

Rare, explicit, reversed: `git mv` back, set `status: active`, remove `release`/`archived`. Reopening usually means new work — prefer a fresh spec that references the archived one.

## Anti-patterns

- Archiving with failing or missing QA silently — the override must be the user's word, on the record.
- Blocking the archive on a merged PR or release tag — completion (tasks + QA) is the gate; publishing is a separate, user-driven step.
- Leaving a completed spec in the active set "until the PR merges" — the active folder is for live work only.
- Editing an archived spec — it is a record; new requirements get a new spec that links back.
- Deleting instead of archiving — the graveyard is where "didn't we already try this?" gets answered.
