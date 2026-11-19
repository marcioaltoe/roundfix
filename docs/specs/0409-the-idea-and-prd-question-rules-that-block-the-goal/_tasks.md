---
schema: spec-tasks/v1
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — The Idea and the PRD question the rules that block the goal

| id      | title                                                                           | type | complexity | needs   |
| ------- | ------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The Spec workflow puts a Rule Challenge to the user                             | docs | medium     | —       |
| task_02 | The authoring skills put a Rule Challenge to the user                           | docs | medium     | task_01 |
| task_03 | The Idea and the PRD record each Rule Challenge, and a Spec waits for the answer | docs | medium     | task_02 |
| task_04 | Run the final QA gate                                                           | qa   | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The graph follows `_techspec.md` → Build Order. task_01 is step 1: the Baseline rule with its two clauses, the closed list of clause forces, the clause tests and the regeneration through the Managed Refresh. Steps 2 and 3 become task_02 and task_03, cut by what a reader meets rather than by file kind: task_02 gives the four authoring skills the step that asks, and task_03 gives the artifacts the record and the wait (the two templates, the `modes` reference) with the glossary and the user guide that name them. Step 3 alone, two glossary entries and one subsection, would not earn its own verification cycle, so it joins the Task that adds the record those words describe.

The Tasks run in series. task_02 follows task_01 as the Build Order says, so the skills restate a clause already on the branch and `make baseline-digests` runs before any skill changes. task_02 and task_03 both change the `SKILL.md` of `write-idea` and `write-prd`, through the version record, and the owned-skill version record; the Wave collision rule puts them in series. Each implementation Task runs the suites that consume its changed contract in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec is one context, the question an author asks about a rule, in three implementation Tasks that each pay for their own cycle: a Baseline rule with its derived files and 33 pinned guarantees, four skills with 24 pinned guarantees in both copies, and two templates, a mode reference, the glossary and a guide page with their contracts. No Task changes a CLI surface, so none has a Surface Transcript.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on 2026-11-19 on a scratch worktree of `348962d2`, outside the planning clone and removed afterwards, where the texts of the TechSpec's API Contracts, the regeneration commands of `_authorization.md` and throwaway tests written to these Task files' subtest names and phrases were applied. Every Verification command of task_01, task_02 and task_03 exited 0 there; the regeneration rewrote exactly the files the Tasks declare; and each of the nine sabotages of task_04 failed its named subtest.

- **Subtest names and phrases.** The TechSpec names the tests and lists the guarantees; each Task names one subtest per guarantee and the phrase it pins, taken word for word from the API Contract, so the QA gate's sabotages name the subtest that must fail.
- **The managed refresh in Verification.** task_01 plans a Managed Refresh with a binary it builds into a temporary directory and a fresh `HOME`, so the Verification never reads the real home directory, as `_authorization.md` → Limits requires.
- **A third test file.** `_techspec.md` → Testing Approach puts the template and `modes` tests in `skills/rule_challenge_authoring_test.go` with the skill section test. A Task may declare an existing file or one it creates, never one an earlier Task creates, so task_03 creates `skills/rule_challenge_record_test.go` for its two tests, under the names the TechSpec gives. The path was measured not governed with the probe of `_authorization.md`, on the planning clone on 2026-11-19, and no grant changes.
- **Owned-skill versions.** task_02 raises `brainstorming`, `write-idea`, `write-prd` and `write-techspec`; task_03 raises `write-idea`, `write-prd` and `roundfix` again on its own tree, because the record command raises every skill whose files changed. One release carries both raises.
- **The Skill Coverage Map.** The user guide page `docs/user-guide/context-driven-development.md` is covered by `write-prd`, among others; task_03 changes the page and declares `.agents/skills/write-prd/SKILL.md`, so the map needs no Coverage Review and `docs/references/skill-coverage.json` does not change.

## Shared files

Spec 0411 delivers before this Spec and Spec 0410 after it (`_prd.md` → Prerequisites). Each Task edits only its own sentences and regenerates every shared record on the tree it runs on, with no version or digest written by hand.

- With Spec 0411: `.agents/skills/roundfix/references/modes.md` and its mirror (0411 changes step 5 of `--update`; task_03 changes the `--unattended` section and step 5 of `--follow-up`), both copies of the Roundfix `SKILL.md` (version fields), `skills/testdata/owned-skill-versions.json`, `CONTEXT.md` and `docs/references/behavior-surfaces.json`.
- With Spec 0410: `internal/baseline/module-versions.json`, the derived Baseline files under `internal/baseline/testdata/` and the Standard TypeScript profile pin, `docs/agents/setup-context.json`, the Roundfix skill's version and `modes` reference, `skills/testdata/owned-skill-versions.json` and `CONTEXT.md`. Spec 0410 changes the modules `core` and `context-workflow`, which this Spec does not touch.
- Not shared: `write-tasks`, `qa-gate`, `archive-spec` and `implement-task` stay as they are (`_prd.md` → Non-Goals).

## Time in tests

No test of this Spec reads the clock or waits on a timer; every new test reads files or builds a plan on a temporary repository.
