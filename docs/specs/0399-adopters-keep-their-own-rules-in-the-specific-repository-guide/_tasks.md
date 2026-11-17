---
schema: spec-tasks/v1
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
qa: task_04
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
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — Adopters keep their own rules in the specific-repository guide

| id      | title                                                                                                   | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Baseline text sends every repository rule to the specific-repository guide, with an example         | docs    | medium     | —                |
| task_02 | The Baseline lists each Misplaced Repository Rule from the Setup Manifest's carriers                    | backend | high       | —                |
| task_03 | Baseline update and Doctor name each Misplaced Repository Rule, and the guides and skills describe them | backend | high       | task_01, task_02 |
| task_04 | Run the final QA gate                                                                                   | qa      | high       | task_03          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

The graph follows `_techspec.md` → Build Order. Tasks 01 and 02 share no
file, so they run together: task 01 changes the catalog text, its derived
files and this repository's managed guides; task 02 adds a new file pair to
`internal/baseline` and one glossary entry. Task 03 needs both: the detector
of task 02, and the managed refresh of task 01, without which its Surface
Transcript does not report this repository current. Each implementation Task
runs the suites that consume its changed contracts in its own Verification,
so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, where a
repository's own rules go, in three implementation Tasks and the QA gate.
Each Task earns its verification cycle: a reworded catalog with its
regenerated records and refreshed guides, a pure detector with a table of
the inventory's shapes, and two command surfaces with their guides and
skills. The command surfaces and their documentation stay in one Task
because a Task that changes a command names its guide itself or through a
Task it needs, and the guide Task records the Behavior Surface Record.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on the planning branch, based on `10ecd250`.

- **The derived set of task 01 was measured** on a scratch worktree of the
  planning branch, removed afterwards. With the three clauses reworded, the
  four templates changed, each changed rule, guide, root block and template
  version raised, the module record command, `make baseline-digests` and
  the sanctioned managed refresh changed exactly these files, and a second
  refresh reported `current` with `File changes: 0`:
  - the governed sources: the three modules, the four templates and
    `templates/index.json`;
  - the governed derived files: the Standard TypeScript profile pin and its
    three formatter goldens `agent-instructions.md`, `docs-layout.md` and
    `skill-dispatch.md`;
  - this repository's governed guides `docs/agents/agent-instructions.md`,
    `docs/agents/docs-layout.md`, `docs/agents/skill-dispatch.md` and the
    Setup Manifest;
  - files that are not governed: `AGENTS.md`,
    `internal/baseline/module-versions.json`, the three catalog records under
    `internal/baseline/testdata/` and the four plan characterization goldens.

  No Setup Snapshot, Source Baseline, other profile, formatter golden
  `AGENTS.md` or `docs/agents/specific-repository.md` changed, so
  `_authorization.md` needs no amendment. The same run showed that the
  only existing tests the new text breaks are the two of
  `vendored_skill_clauses_test.go`, which pin
  `clause.core.activate-matching-skills` word for word; task 01 changes
  that literal.
- **This repository reports no Misplaced Repository Rule.** Every Markdown
  carrier its Setup Manifest lists holds only Managed Regions, and its root
  `CLAUDE.md` is a symbolic link to `AGENTS.md`, as the inventory recorded.
  The Surface Transcript of task 03 depends on it.
- **The place of the Doctor line.** The guide of `roundfix doctor` already
  promises that `user-skills:` comes right after `skills:` and that
  `authoring-models:` follows them. `baseline-rules:` therefore comes after
  `skills:` and after those two lines when they are present, and before the
  residue lines, which keeps both promises and the TechSpec's "after the
  `skills` line".
- **Doctor tests that count lines.** Two table tests of
  `internal/cli/doctor_test.go` count fifteen lines. With no Setup Manifest
  under their fake repository root, Doctor prints one more line,
  `baseline-rules: skipped (...)`, so task 03 changes those counts and no
  other assertion of theirs.
- **Exact text.** Task 01 fixes the words of the three clauses and the four
  templates from `_techspec.md` → API Contract 3, and task 03 the sentences
  of the skills, so the QA gate searches for fixed strings.

## Time in tests

No test of this Spec reads the clock, compares a time it seeded with the
real clock, sleeps or waits on a timer. The detector and both commands read
files only.
