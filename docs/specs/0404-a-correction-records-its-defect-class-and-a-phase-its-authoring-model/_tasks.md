---
schema: spec-tasks/v1
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
requires: [0401-a-supervisor-session-ends-only-when-the-maintainer-asks, 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration, 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation]
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A Correction records its Defect Class and a phase entry its authoring model

| id      | title                                                                                                                 | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | A Correction is declared, with a fixed finding or on its own, into one record per Spec                                | backend | high       | —                         |
| task_02 | A phase entry of an Authoring Phase records its model, effort, runner and repository                                  | backend | high       | task_01                   |
| task_03 | The Delivery Queue reads the Correction Signs at its merge and parks an undeclared one where the repository requires it | backend | high       | task_01, task_02          |
| task_04 | The Corrections Report counts Corrections per Spec and plan defects per authoring model, beside its sample            | backend | high       | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                                                                 | qa      | high       | task_04                   |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The graph follows `_techspec.md` → Build Order, one Task per step. The steps
form one chain: each Task changes the Roundfix skill's `SKILL.md` pair,
`skills/testdata/owned-skill-versions.json`,
`docs/references/behavior-surfaces.json` and `CONTEXT.md`, and tasks 01, 02
and 04 change the root help in `internal/cli/cli.go`, so no two can share a
Wave. Task 02 needs task 01 for `baseline.InspectRepositoryIdentity`; task 03
needs task 01 for the record and task 02 for the order of the shared files;
task 04 reads what the three wrote. Each implementation Task runs the suites
that consume its changed contracts in its own Verification, so no separate
Integration Checkpoint is needed. The manifest requires Specs 0401, 0402 and
0403, which change the phase command, the authoring skills, the phase
reference, the Roundfix skill, the glossary and the guides before this Spec
does, as `_prd.md` → Prerequisites orders.

Sized by the grouping criterion, the Spec holds one context, the record of
why a delivery was corrected and which model wrote the plan, in four
implementation Tasks and the QA gate, the reference size. Each Task earns
its verification cycle: a new package, two writers and a retention target;
a writer, a command and four skills; a reader, a configuration key and a
step of the Delivery Queue; and a report with its JSON.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on 2026-11-17 on scratch worktrees of
`d65bae06`, removed afterwards, and is recorded in
[`references/2026-11-17-tasks-measurements.md`](references/2026-11-17-tasks-measurements.md).
Specs 0401, 0402 and 0403 were not on the tree, so each Task runs its
packages again on its own base.

- **The tests the change moves were measured.** A prototype of every
  surface an existing test observes moved, in `./internal/cli`, the GC
  report of `TestGCCountsPhaseRecordDayFilesWithoutARunDatabase` and the
  pattern three Worktree Build Cache tests share, which task 01 declares;
  and `TestPhaseStartSurfaceTranscripts`, four subtests of
  `TestPhaseStartRecordsEachOfTheSevenPhases`,
  `TestPhaseStartLoadsNoConfiguration` and, in `./internal/phaserecord`,
  `TestAPhaseEntryIsOneLineOfItsDeclaredKeys`, which task 02 declares.
  `TestEveryCommandPathPrintsItsHelp` and three documentation contracts read
  every command path of the root help; they keep their expectations, and
  tasks 01 and 04 each name the help and the documentation their command
  paths need. `./internal/delivery`, `./internal/config`,
  `./internal/runcause`, `./internal/baseline` and `./skills` passed, and
  `TestReviewDisposeHelpPrintsUsage` passed with the new `--fixed-by` line.
  The merge step and the sign reader were not prototyped; task 03 runs its
  packages before and after, as `_techspec.md` → Testing Approach says.
- **The derived set was measured.** With the canonical references, the
  five `SKILL.md` files and the six guide pages changed and the
  `corrections` reference and page created, `make skills-sync`,
  `-record-skill-versions` and `-record-skill-coverage` wrote exactly the
  mirrors, the five `SKILL.md` pairs,
  `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`, and nothing under
  `.claude/agents/`. The last step needs the hand-written Skill Coverage Map
  entry of the new guide page, which task 01 writes.
- **The authorization record needs no amendment.** `speccheck.GovernedPath`
  over the 92 paths the Tasks declare found 16 governed, exactly the 16 of
  `_authorization.md` `paths`.
- **A Task that changes a file a needed Task creates declares it under
  `creates:`.** Task 04 extends the `corrections` command file, reference,
  guide page and skill test that task 01 creates; the checker resolves an
  `interface:` path only on the tree, so those paths are declared under
  `creates:`, and the chain keeps the two Tasks in series.
- **`not-a-correction` reaches every Task that names the classes.** Task 01
  validates and records it, task 03 names it in the park's next action and
  covers a sign with it, task 04 counts it apart, and the QA gate sabotages
  it, as the TechSpec amendment of 2026-11-17 requires.
- **The QA gate keeps the host's Rust toolchain.** The gate runs the whole
  `./internal/cli` package under a temporary `HOME` with `RUSTUP_HOME` and
  `CARGO_HOME` taken from the host and used read only, because that
  package's Baseline tests resolve the Rust toolchain through the home
  directory, which stopped the QA of Specs 0399, 0400 and 0401.
- **The release is no Task.** The authorization record's Release Ceiling
  `minor` answers the release plan after the merge; this repository's
  Project Config gains the key in the operator's own commit after that
  release, as `_authorization.md` says.
