---
schema: spec-tasks/v1
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
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

# Tasks — A Supervisor session ends only when the maintainer asks

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Phase Record marks no boundary and the phase command prints no notice                     | backend | medium     | —                |
| task_02 | The autonomous-work clause says no phase of a Spec ends the session                            | docs    | medium     | —                |
| task_03 | The authoring skills, the Roundfix modes and the glossary say no phase of a Spec ends the session | docs    | medium     | task_01          |
| task_04 | Run the final QA gate                                                                         | qa      | high       | task_02, task_03 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

The graph follows `_techspec.md` → Build Order. Task 01 changes the
Phase Record's writer and the phase command with their tests and the texts
that describe them. Task 02 rewords the Baseline clause and regenerates its
derived files; it shares no file with the other Tasks, so it runs beside
task 01. Task 03 needs task 01: both raise the Roundfix skill's version
through `skills/testdata/owned-skill-versions.json`. Each implementation
Task runs the suites that consume its changed contracts in its own
Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, when a
Supervisor session ends, in three implementation Tasks and the QA gate.
Each Task earns its verification cycle: a writer and a command with their
table, transcript, reader and user-guide tests; a Baseline clause with the
whole Baseline suite and a managed refresh; and five skills with their
mirror and version records, the glossary and an ADR amendment.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on 2026-11-17 on a scratch worktree of
`bcd716df`, removed afterwards.

- **The derived set was measured.** With the clause reworded,
  `-record-module-versions`, `make baseline-digests` and the sanctioned
  `baseline update` changed exactly the files task 02 declares, and a
  second refresh reported `current`; `AGENTS.md` did not change. With the
  four sections, the modes and the phase reference reworded,
  `make skills-sync` and `-record-skill-versions` changed the mirrors, the
  five `SKILL.md` pairs and `skills/testdata/owned-skill-versions.json`.
  With both user-guide pages rewritten, `-record-skill-coverage` changed
  `docs/references/behavior-surfaces.json`, so only task 01 declares it.
  `speccheck.GovernedPath`, run through a `go test -overlay` probe that
  wrote nothing to the tree, found 17 governed paths among the changed and
  planned files, exactly the 17 in `_authorization.md` `paths`, so the
  record needs no amendment.
- **The tests the change moves were measured.** With the writer and the
  command changed and the texts reworded, and no test changed:
  `./internal/phaserecord` no longer compiles until `record_test.go` and
  `read_test.go` follow the new `StartPhase`, and neither does
  `./internal/cli` until `follow_up_queue_test.go` does; in `./internal/cli`
  `TestPhaseStartSurfaceTranscripts/transcript_3` and
  `TestPhaseBoundaryNoticeOmitsAnAbsentEarlierSpec` fail; in `./skills`
  `TestThePhaseReferenceDocumentsTheCommandAndTheRecord/*/the_notice_refuses_nothing`,
  the `no_further_phase` checks of
  `TestTheAuthoringSkillsStateThePhaseBoundary`, and the `no_further_phase`,
  `a_wake_up_never_crosses`, `unattended_ends_at_the_boundary` and
  `unattended_wake_up` checks of
  `TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary` fail; in
  `./internal/baseline` the `phase_boundary_is_recorded`,
  `no_further_phase_starts` and `wake_up_stays_in_the_session` subtests of
  the clause test fail; `-tags docscontract ./internal/docscontract`
  passes. Each Task declares the tests it moves.
- **The Coverage Review of the usage page is needed.** The Behavior Surface
  `guide docs/user-guide/usage.md` is not covered by the Roundfix skill's
  phase reference, so task 01 records its Coverage Review.
- **The glossary moves in one Task.** The current **Phase Boundary**
  definition names the notice, so removing the notice's entry in task 01
  would leave an undefined phrase in the glossary until task 03; task 03
  makes all three glossary changes together.
- **The other open Spec plan shares only derived files.** Spec 0400 raises
  the Roundfix skill's version. Every shared file has a derived-path
  declaration in `.roundfixrc.yml`, so no `requires` is declared.
