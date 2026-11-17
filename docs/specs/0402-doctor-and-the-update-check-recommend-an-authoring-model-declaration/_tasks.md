---
schema: spec-tasks/v1
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
qa: task_03
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
---

# Tasks — Doctor and the Update Check recommend an Authoring Model Declaration

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Doctor and the Update Check recommend an Authoring Model Declaration when none is declared               | backend | medium     | —       |
| task_02 | This repository's declaration and its guides follow the Recommended Authoring Declaration                | docs    | medium     | task_01 |
| task_03 | Run the final QA gate                                                                                   | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The graph follows `_techspec.md` → Build Order. Task 01 adds the
recommendation, the `undeclared` report and its two command guides,
because a Task that changes a command names its guide itself or through a
Task it needs. Task 02 needs it: its contract test calls the rendering
task_01 adds, and the configuration guide, the Roundfix skill and the
glossary describe the report task_01 builds. Both declare the Behavior
Surface Record, so the edge also runs them in series. Each implementation
Task runs the suites that consume its changed contracts in its own
Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, the report of
an undeclared authoring model, in two implementation Tasks and the QA gate.
Each Task earns its verification cycle: the command behavior with its
fixtures, transcript and the two changed tests, and the drift gate with the
skill, glossary and ADR text and their regenerated records. A separate Task
for the `internal/config` table would pay a whole cycle for one file whose
only consumer is task_01's report, so it joins task_01.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on the planning branch, based on `bcd716df`.

- **The tests whose expectations change were measured.** A probe on a
  scratch worktree of `bcd716df`, removed afterwards, made both checks
  report every installed, unpublished authoring skill when nothing is
  declared and ran the whole `./internal/cli` package with
  `-count=1`. Exactly two tests failed:
  `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted` and
  `TestRunDoctorRealRepositoryCheckDoesNotMutateState`. task_01 declares
  both files and names both changes.
- **The governed set was measured** with a `go test -overlay` probe of
  `speccheck.GovernedPath` over the 38 paths the Tasks declare or read.
  The governed ones the Tasks change are the four in `_authorization.md`
  `paths`. `internal/cli/cli_test.go`,
  `internal/docscontract/publicdocs_test.go` and `.roundfixrc.yml` are
  governed and no Task changes them; the contract test reads
  `.roundfixrc.yml` and no Task declares it.
- **The derived set of task_02** follows the pattern Spec 0400 measured for
  the Roundfix skill: `make skills-sync` rewrites the two mirrors under
  `skills/roundfix/references/`, and the `-record-skill-versions` step
  raises the `SKILL.md` pair and writes
  `skills/testdata/owned-skill-versions.json`. task_02 declares each, and
  the governed ones are in `paths`.
- **The Skill Coverage Map** lists `.agents/skills/roundfix/references/setup.md`
  as the covering skill of `command doctor`, `command update check`, their
  two guides, `config keys` and the configuration guide. task_01 changes
  `internal/cli/doctor.go` and the two command guides, and creates a file
  under `internal/config`; task_02 changes the covering setup reference.
- **The phrase test.** `skills/authoring_delegation_test.go` pins
  "Roundfix ships no model name" in the setup reference; task_02 declares
  it and changes the phrase. No other `_test.go` file pins the sentences the
  Tasks change.
- **The transcript runs on a development version.** The QA gate builds a
  binary stamped `0.0.0-dev`, so the Update Check makes no release lookup
  and needs no network.
