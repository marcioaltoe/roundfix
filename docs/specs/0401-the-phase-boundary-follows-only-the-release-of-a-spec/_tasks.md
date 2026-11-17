---
schema: spec-tasks/v1
spec: 0401-the-phase-boundary-follows-only-the-release-of-a-spec
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
---

# Tasks — The Phase Boundary follows only the release of a Spec

| id      | title                                                                                                       | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Boundary Notice follows only a release                                                                  | backend | medium     | —                |
| task_02 | The autonomous-work clause states the boundary after a Spec's delivery and release                           | docs    | medium     | —                |
| task_03 | The authoring skills and the Roundfix modes state the boundary after a Spec's delivery and release           | docs    | medium     | task_01          |
| task_04 | Run the final QA gate                                                                                       | qa      | high       | task_02, task_03 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

The graph follows `_techspec.md` → Build Order. Task 01 changes the
boundary rule with its tests and the texts that describe the notice. Task 02
rewords the Baseline clause and regenerates its derived files; it shares no
file with the other Tasks, so it runs beside task 01. Task 03 needs task 01:
both raise the Roundfix skill's version through
`skills/testdata/owned-skill-versions.json` and both revise `CONTEXT.md`.
Each implementation Task runs the suites that consume its changed contracts
in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, where a
Supervisor session ends, in three implementation Tasks and the QA gate.
Each Task earns its verification cycle: a rule with its table, transcript
and user-guide tests; a Baseline clause with the whole Baseline suite and a
managed refresh; and five skills with their mirror and version records.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on the planning branch, based on `379803c8`.

- **The derived set was measured** on a scratch worktree of `379803c8`,
  removed afterwards. With the clause reworded,
  `-record-module-versions`, `make baseline-digests` and the sanctioned
  `baseline update` changed exactly the files task 02 declares, and a
  second refresh reported `current`; `AGENTS.md` did not change. With the
  four sections, the modes and the phase reference reworded,
  `make skills-sync` and `-record-skill-versions` changed the mirrors, the
  five `SKILL.md` pairs and `skills/testdata/owned-skill-versions.json`;
  `-record-skill-coverage` changed `docs/references/behavior-surfaces.json`
  only when a user-guide page changed, so only task 01 declares it. The
  governed paths among them are the 17 in `_authorization.md` `paths`, so
  the record needs no amendment.
- **The current tests keep passing after the reword.** On that worktree,
  with every text changed and no test changed, `./skills` and
  `-tags docscontract ./internal/docscontract` passed, and in
  `./internal/baseline` only `TestTheOutliveTheTurnClauseStatesThePhaseBoundary`
  failed. So each Task adds checks of the new wording and of the absence of
  the old one, which fail before its work.
- **The tests the boundary rule moves were measured** with the condition
  changed on the same worktree: in `./internal/phaserecord`,
  `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession`,
  `TestStartPhaseSkipsUnrelatedAndMalformedHistory` and
  `TestStartPhaseReadsAnOversizedLineWithoutBufferingIt`; in
  `./internal/cli`, `TestPhaseStartSurfaceTranscripts/transcript_3`. Task 01
  declares both test files.
- **The other open Spec plans share only derived files.** Specs 0399 and
  0400 raise the Roundfix skill's version, and 0399 records Baseline module
  versions and the profile pin. Every shared file has a derived-path
  declaration in `.roundfixrc.yml`, so no `requires` is declared.
