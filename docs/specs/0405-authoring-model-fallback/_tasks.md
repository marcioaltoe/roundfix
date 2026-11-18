---
schema: spec-tasks/v1
spec: 0405-authoring-model-fallback
requires: [0401-a-supervisor-session-ends-only-when-the-maintainer-asks, 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration, 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation, 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model]
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Authoring Model Declaration accepts a fallback for each Authoring Phase, and the recommendation carries one | backend | medium     | —                |
| task_02 | Roundfix decides and records an Authoring Fallback and prints the model a phase starts on      | backend | high       | task_01          |
| task_03 | The authoring skills start the Author Subagent on the printed model and switch on its limit    | docs    | medium     | task_01, task_02 |
| task_04 | The Baseline tells each root instruction file how an Authoring Phase keeps its model           | backend | medium     | task_03          |
| task_05 | Run the final QA gate                                                                          | qa      | high       | task_04          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The graph follows `_techspec.md` → Build Order, one Task per step. The steps
form one chain. Each of tasks 01 to 03 changes the Roundfix skill's `SKILL.md` pair,
`skills/testdata/owned-skill-versions.json` and `CONTEXT.md`. Tasks 01 and 02
also change `docs/references/behavior-surfaces.json`. So no two Tasks can
share a Wave.

- Task 02 needs task 01 for the `fallback` key it reads.
- Task 03 needs task 02 for the two commands its skill text runs.
- Task 04 needs task 03: its statement names `roundfix phase fallback`, which
  task 02 adds, and its regeneration (`make baseline-digests` and the Managed
  Refresh of `AGENTS.md` and `docs/agents/setup-context.json`) must not run
  beside task 03's skill regeneration, so the two do not share a Wave.

Each implementation Task runs the suites that consume its changed contracts in
its own Verification, so no separate Integration Checkpoint is needed.

The manifest requires Specs 0401 to 0404. They change the phase command, the
phase entry, the authoring skills' `## Declared model` and phase sections, the
Roundfix skill, its phase reference, the glossary, the surface record and the
configuration guide before this Spec does, as `_prd.md` → Prerequisites
orders. Spec 0404 also adds the repository identity that task 02 reads.

The Spec holds one context, how an Authoring Phase keeps its model and falls
back, in four implementation Tasks and the QA gate, at the reference size.
Each Task earns its verification cycle:

- a configuration key with its refusals, its guide and the recommended block with its contract;
- a pure reader, a record kind and two commands;
- four skills and their test;
- a Baseline root block with its versions, its derived files, this repository's Managed Refresh and its test.

The maintainer joined the root statement to this Spec on 2026-11-18 ("Juntar
na 0405 (Recommended)"); it does not fit task 03, whose skill regeneration
and documentation contracts it would double.

## What the Tasks stage settled beyond the TechSpec

[`references/2026-11-17-harness-and-test-measurements.md`](references/2026-11-17-harness-and-test-measurements.md)
records each item, measured on 2026-11-17 on a scratch worktree of
`a8236a87` that was removed afterwards. Specs 0401 to 0404 were not on the
tree, so each Task runs its packages again on its own base.

- **The tests the change moves.** The prototype moved
  `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/phase_field` in
  `./internal/config` and `TestPhaseStartSurfaceTranscripts/transcript_1` in
  `./internal/cli`; task 01 declares the first and task 02 the second. Spec
  0404 adds its own test of the full `roundfix phase --help`, which task 02
  updates on its base.
- **The documentation contracts.** The four documentation contracts that
  failed (`TestEveryCommandIsNamedInTheRoundfixSkill`,
  `TestEveryCommandIsNamedInTheUserGuide`,
  `TestAnUndocumentedCommandIsReported` and
  `TestTheSkillCoverageMapIsCurrent`) pass once task 02 writes the guide, the
  reference and the coverage entries and runs the regeneration.
- **Untouched packages.** `./skills` passed.
- **Files of Spec 0402.** Task 01 declares under `creates:` the four files Spec 0402 adds: the recommendation, its test, its command test and its repository contract. They were not on this plan's tree, so `interface:` could not resolve. On the delivery base, which holds Spec 0402, the Task edits them.
- **Governed Paths.** The record bounds the Roundfix skill's setup reference,
  which describes the Authoring Model Declaration and which task 01 changes.

- **The root statement.**
  [`references/2026-11-18-root-statement-decision-and-measurement.md`](references/2026-11-18-root-statement-decision-and-measurement.md)
  records the prototype on a scratch worktree of `eb888135`, removed
  afterwards: the files its regeneration rewrites, which task 04 declares,
  the six Governed Paths it adds to the record, and the packages that passed
  unchanged (`./internal/baseline`, `./internal/cli`, `./internal/docscontract`
  with its tag, `./skills`, `./internal/speccheck`, `./internal/spec`,
  `./internal/baselineacp`). It moves no existing test.

QA gate: included. The Spec adds a configuration key, two command surfaces,
skill text that a session follows and a Baseline root statement, and the gate
reproduces the eight Surface Transcripts through the built binary and adopts a
fresh fixture to read the statement in both root files.
