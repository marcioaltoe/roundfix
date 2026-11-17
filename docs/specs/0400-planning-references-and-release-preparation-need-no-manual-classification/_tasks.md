---
schema: spec-tasks/v1
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — Planning references and release preparation need no manual classification

| id      | title                                                                                                        | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | The Release Plan counts planning references and a Release-Preparation Commit as maintenance-only              | backend | medium     | —                |
| task_02 | The release guides describe both rules and keep the repository's own publication trigger                      | docs    | medium     | task_01          |
| task_03 | The release plan command reports a Release-Preparation Commit and reaches the same decision after it          | backend | medium     | task_01, task_02 |
| task_04 | Run the final QA gate                                                                                        | qa      | high       | task_03          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The graph follows `_techspec.md` → Build Order. Task 01 changes the release
plan domain with table tests and no Git. Task 02 needs it: the usage text,
the release reference, the runbook, the usage page and the glossary
describe the rules task 01 builds, and the JSON field whose name
`_techspec.md` → API Contract 2 fixes. Task 03 needs both: the command's
JSON evidence reads task 01's field, and a Task that changes a command
names its guide itself or through a Task it needs, which task 02 is. Each
implementation Task runs the suites that consume its changed contracts in
its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, how a release
is classified and published, in three implementation Tasks and the QA
gate. Each Task earns its verification cycle: a domain rule with its
tables, a command surface with fixture repositories and a transcript on
this repository's history, and the guides and skill with their regenerated
records. The release-guide half of the Quaestor's entry joins task 02
because it rewrites the same section of the same reference.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on the planning branch, based on `c2d1f77c`.

- **The derived set of task 02 was measured** on a scratch worktree of
  `c2d1f77c`, removed afterwards. With a sentence added to the canonical
  release reference, a paragraph to the usage text in `internal/cli/cli.go`
  and a sentence to the runbook and the usage page, `make skills-sync`, the
  `-record-skill-versions` step and the `-record-skill-coverage` step
  changed exactly these files besides the edited ones:
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` (the
  version, `0.1.246` to `0.1.247`), `skills/roundfix/references/release.md`,
  `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json` (the `command release plan`,
  release runbook and usage page surfaces). The whole `./skills` package
  and `-tags docscontract ./internal/docscontract` passed afterwards. The
  governed paths among them are the three in `_authorization.md` `paths`,
  so the record needs no amendment.
- **No test pins the text this Spec changes** outside the files the Tasks
  declare: a search of the `_test.go` files for the usage paragraph, the
  planning-evidence sentence and "Preserve the existing tag-triggered
  workflow" found none. `internal/docscontract/publicdocs_test.go` pins
  "tag-triggered" and "publication workflow" in the runbook, which task 02
  keeps.
- **The Skill Coverage Map** lists the release runbook as the source of the
  `command release plan` and runbook surfaces, covered by the canonical
  release reference, which task 02 changes. `releaseplan_command.go` is the
  source of no surface, so task 03 needs no skills Task of its own.
- **The transcript test reads this repository's own history.** CI checks
  out with `fetch-depth: 0`, so `v0.143.3` and `v0.143.4` are present
  there and in every Run worktree; the test clones with `--shared` into a
  temporary directory because the plan refuses a dirty worktree.
