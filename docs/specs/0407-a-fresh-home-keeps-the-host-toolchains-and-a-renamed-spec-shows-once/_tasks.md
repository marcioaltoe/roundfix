---
schema: spec-tasks/v1
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
qa: task_03
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
---

# Tasks — A fresh home keeps the host toolchains, and a renamed Spec shows once

| id      | title                                                     | type    | complexity | needs            |
| ------- | --------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The status view shows a renamed Spec once                  | backend | medium     | —                |
| task_02 | A fresh home keeps the host toolchains in the QA skills    | docs    | medium     | —                |
| task_03 | Run the final QA gate                                      | qa      | high       | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03

The graph follows `_techspec.md` → Build Order. The two implementation Tasks share no declared file, so they run in the same Wave: task_01 changes the status view's row reader with its guide, glossary entry, ADR-0410 amendment, Coverage Review and Behavior Surface Record; task_02 changes the `qa-gate` and `write-tasks` skills with their mirrors, versions and contract test. Each runs the suites that consume its changed contracts in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds two defects of one Backlog Entry that share no component, kept together by the maintainer's answer "Uma Spec só (Recommended)" (`_prd.md` → Decisions). Each half is one Task that earns its verification cycle: a reader change with eight cases and four documents, and two skill sections with a contract test of sixteen guarantees per copy. Splitting the source into two Specs would give each a one-Task Spec paying a whole Run, QA gate, review and release.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on a prototype on a scratch worktree of the planning branch, based on `1a8e08b4`, removed afterwards; the measurements are in [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md).

- **The skill edits do not move the Behavior Surface Record**; only the `follow-up` guide edit does, so only task_01 re-records it. The Skill Coverage Map covers `internal/cli/follow_up_queue.go` and the `follow-up` guide by the Roundfix `follow-up` reference, which needs no change, so task_01 records the Coverage Review that `_prd.md` → Skills declares instead of changing a governed Roundfix skill file.
- **Existing tests.** With the prototype, `TestASpecIsBeingAuthored`, `TestTheQueueRowsComeInTableOrder`, the whole `./skills` package and `-tags docscontract ./internal/docscontract` passed unchanged; the whole `./internal/cli` package is in task_01's Verification.
- **Owned-skill versions.** The record command raised `qa-gate` and `write-tasks` one patch each on that base. Spec 0408 also appends to `skills/testdata/owned-skill-versions.json` and edits `CONTEXT.md` and `docs/references/behavior-surfaces.json`; it delivers after this Spec and regenerates on its own tree.
- **The fresh home of a test.** A `HOME` or `TMPDIR` inside the checkout is refused by the suite guard, so every fresh home of the QA gate lives outside it.
- **Exact text.** Both implementation Tasks copy their text word for word from `_techspec.md` → API Contracts 1 and 2, so the QA gate searches for fixed strings.

## Time in tests

No test of this Spec reads the clock or waits on a timer. `TestARenamedSpecLeavesNoAuthoringRow` records phase entries at fixed times relative to the fixed time its input carries.
