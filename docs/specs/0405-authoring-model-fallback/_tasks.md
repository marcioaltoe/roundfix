---
schema: spec-tasks/v1
spec: 0405-authoring-model-fallback
requires: [0401-a-supervisor-session-ends-only-when-the-maintainer-asks, 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration, 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation, 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model]
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

# Tasks — An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Authoring Model Declaration accepts a fallback for each Authoring Phase, and the recommendation carries one | backend | medium     | —                |
| task_02 | Roundfix decides and records an Authoring Fallback and prints the model a phase starts on      | backend | high       | task_01          |
| task_03 | The authoring skills start the Author Subagent on the printed model and switch on its limit    | docs    | medium     | task_01, task_02 |
| task_04 | Run the final QA gate                                                                          | qa      | high       | task_03          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The graph follows `_techspec.md` → Build Order, one Task per step. The steps
form one chain. Each Task changes the Roundfix skill's `SKILL.md` pair,
`skills/testdata/owned-skill-versions.json` and `CONTEXT.md`. Tasks 01 and 02
also change `docs/references/behavior-surfaces.json`. So no two Tasks can
share a Wave.

- Task 02 needs task 01 for the `fallback` key it reads.
- Task 03 needs task 02 for the two commands its skill text runs.

Each implementation Task runs the suites that consume its changed contracts in
its own Verification, so no separate Integration Checkpoint is needed.

The manifest requires Specs 0401 to 0404. They change the phase command, the
phase entry, the authoring skills' `## Declared model` and phase sections, the
Roundfix skill, its phase reference, the glossary, the surface record and the
configuration guide before this Spec does, as `_prd.md` → Prerequisites
orders. Spec 0404 also adds the repository identity that task 02 reads.

The Spec holds one context, the fallback of an Authoring Phase, in three
implementation Tasks and the QA gate, under the reference size. Each Task
earns its verification cycle:

- a configuration key with its refusals, its guide and the recommended block with its contract;
- a pure reader, a record kind and two commands;
- four skills and their test.

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

QA gate: included. The Spec adds a configuration key, two command surfaces
and skill text that a session follows, and the gate reproduces the eight
Surface Transcripts through the built binary.
