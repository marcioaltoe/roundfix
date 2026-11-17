---
schema: spec-tasks/v1
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
requires: [0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration]
qa: task_05
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
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The Update Check names each profile that differs from a refreshed recommendation

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Model Catalog follows the installed Codex and Claude adapters                                   | backend | medium     | —       |
| task_02 | The Recommended Profile becomes the 2026-11-17 snapshot, and this repository's profiles follow it   | backend | high       | task_01 |
| task_03 | The guides and the --update Skill Mode describe the profiles the Update Check names                 | docs    | medium     | task_02 |
| task_04 | The Update Check names each configured profile that differs from the Recommended Profile            | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                                              | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The graph follows `_techspec.md` → Build Order. Task 01 adds the models the
snapshot recommends to the catalog, because `TestRecommendedModelsAreInTheModelCatalog`
fails on a snapshot naming a model the catalog lacks. Task 02 changes the
snapshot, the reference that states it, the built-in lists of the guides and
this repository's Project Config together, because
`TestModelSelectionReferenceStatesTheShippedSnapshot` and the docs contract
fail on any of them alone. Task 03 writes the guides and the skill text of the
new report, and task 04 the command, because a Task that changes a command
names its guide through a Task it needs. Tasks 01 to 03 each change the
Roundfix skill's `SKILL.md` pair, its version record and the Behavior Surface
Record, so they run in series.

Sized by the grouping criterion, the Spec holds one context, how a repository
learns that its Agent Selection Profiles fell behind the recommendation, in
four implementation Tasks and the QA gate. Each earns its verification cycle:
a picker table with its tests and guide sentences; a snapshot with its tests,
reference, guides and Project Config; the guides, glossary, skill and ADR
amendment of the report; and the report itself with a transcript.

## What the Tasks stage settled beyond the TechSpec

Each was measured on the planning branch, based on `63a9f493`.

- **The tests that pin the old values were found** by applying the new
  snapshot and catalog to a scratch worktree and running `./internal/config`,
  `./internal/agent`, `./internal/cli`, `./internal/tui`, `./skills` and
  `-tags docscontract ./internal/docscontract`. They are
  `TestRecommendedProfileCoversEveryWorkCategory` (`internal/config`),
  `TestModelCatalogOpensWithTheCurrentModels` and
  `TestModelCatalogsExposeOrderedPickerData` (`internal/agent`),
  `TestProfilesShowJSONIsSchemaV2WithoutBenchmarkFields`
  (`internal/cli/profiles_show_recommended_test.go`),
  `TestModelSelectionReferenceStatesTheShippedSnapshot` and
  `TestProfilesDocumentationContractMatchesPublicGuidance`
  (`internal/docscontract`, the second pinning `claude-fable-5-1` and the
  snapshot date in the usage page, the configuration page and the Roundfix
  skill). Tasks 01 and 02 declare each. `internal/cli/cli_test.go` and
  `internal/tui` passed; every other test reads the snapshot through
  `config.RecommendedProfile` or `config.ModelRecommendationSnapshotVersion`.
- **The derived set was measured** on the same worktree, removed afterwards:
  with a sentence added to each guide and skill reference the Tasks change,
  `make skills-sync`, the `-record-skill-versions` step and the
  `-record-skill-coverage` step changed exactly the `SKILL.md` pair (version
  `0.1.248` to `0.1.249`), the mirrors of `profiles`, `setup` and `modes`,
  `skills/testdata/owned-skill-versions.json`, and the usage, configuration
  and update guide entries of `docs/references/behavior-surfaces.json`. The
  governed paths among them are in `_authorization.md` `paths`.
- **The Skill Coverage Map** gives the surface `command update check`
  (`internal/cli/update_check.go`) to the `setup` reference, which task 03
  changes before task 04 changes the source, and the surface `config keys`
  (`internal/config/*.go`) to the `profiles` reference, which task 02 changes
  with the snapshot.
- **Spec 0402 changes the same Update Check, guide and skill references.**
  The manifest requires it, so task 04 is written against its output and
  places the `Profiles:` lines after whatever 0402 prints after the Baseline
  line and before `Next:`.
