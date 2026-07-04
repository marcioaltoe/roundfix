---
name: archive-spec
description: Archive a shipped spec — verify every task completed, QA passed, and the work merged/released, then stamp the release metadata and move docs/specs/<slug>/ to docs/specs/_archived/<slug>/.
disable-model-invocation: true
argument-hint: "<spec slug> [--release <tag or PR URL>]"
metadata:
  category: delivery
  tags: [workflow, documentation, process]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Archive Spec

Move a shipped spec out of the active set: `docs/specs/<slug>/` → `docs/specs/_archived/<slug>/`, with the release stamped in its frontmatter. Archived means _implemented and shipped_ — after this, one `ls docs/specs/` separates live work from history, and the archive stays greppable as the record of what was built and why.

## Preconditions — verify, don't trust

Check all three with fresh evidence before touching anything:

1. **Every task completed** — read each `task_NN.md` listed in `_tasks.md`; every `status` is `completed`.
2. **QA passed** — the newest report in `qa/` has `verdict: pass`. A missing `qa/` directory or a failing latest report blocks the archive; proceed only if the user explicitly says "archive anyway", and record that override in the stamped frontmatter (`qa_override: true`).
3. **Shipped** — the work is merged or released. Confirm via `--release`, or find it: the merged PR (`gh pr list --state merged --search <slug>`) or the release tag containing the spec's commits. Archiving unmerged work rewrites history that hasn't happened.

If any check fails, stop and report exactly which — the spec stays active.

## Steps

1. **Stamp** `_prd.md` frontmatter:

   ```yaml
   status: archived
   release: <tag or PR URL>
   archived: YYYY-MM-DD
   ```

2. **Move** with history preserved:

   ```bash
   mkdir -p docs/specs/_archived
   git mv docs/specs/<slug> docs/specs/_archived/<slug>
   ```

3. **Commit** — `chore(specs): archive <slug>` (Conventional Commits). Do not push unless asked.

4. **Report** — the new path, the release reference, and anything carried over (open follow-ups from task `## Result` sections belong in new specs, not in the archive).

## Unarchive

Rare, explicit, reversed: `git mv` back, set `status: active`, remove `release`/`archived`. Reopening usually means new work — prefer a fresh spec that references the archived one.

## Anti-patterns

- Archiving with failing or missing QA silently — the override must be the user's word, on the record.
- Archiving unmerged work because "it's done locally".
- Editing an archived spec — it is a record; new requirements get a new spec that links back.
- Deleting instead of archiving — the graveyard is where "didn't we already try this?" gets answered.
