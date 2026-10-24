---
task: task_03
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
status: pending
type: docs
complexity: medium
---

# Task 03: The follow-up and unattended modes gather every decision and record each answer where it is read

## Overview

task_02 left the modes reference with how a mode is chosen, the optional
sources and `--status`. This Task writes the other two modes: `--follow-up`,
the eight steps from inventory to an armed continuation, and `--unattended`,
which gathers every decision and authorization in calls of up to three
independent questions, as the clause task_01 reworded allows, and records
each answer where the work reads it. Three tests hold the step order, the
agreement with the Baseline clause, the answer placement and the absence of
personal assistant names (`_prd.md` → Core Features 4 and 5; ADR-0320,
ADR-0095).

## Requirements

1. MUST answer the Backlog Entry Spec 0306 adopted,
   [`../0306-the-follow-up-command-lists-pending-work-in-priority-order/references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md`](../../references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md)
   of 2026-10-22, decision 4 and the `--follow-up` and `--unattended` steps.
2. MUST append to `.agents/skills/roundfix/references/modes.md`, after the
   `## --status` section, a `## --follow-up` section whose steps are one
   ordered list of exactly eight items, in this order:
   1. inventory: run the follow-up command in its JSON form, handing it the
      ai-memory message file when that source exists; when it reports
      nothing to do and nothing in progress, report what it checked and stop;
   2. plan: show a short plan in stages, grouping sources into Specs as
      `docs/agents/spec-routing.md` allows and following the command's
      priority order;
   3. run `--unattended` to gather every decision and authorization the plan
      needs;
   4. triage, after the answers: resolve each pending Inbox Entry of this
      project's namespace with `resolved_to` or `discarded_reason` and move it
      to `_triaged/`, as `docs/agents/secondbrain.md` and ADR-0095 describe;
   5. author each Spec through the repository's Spec routing, consulting the
      Secondbrain and Exa before the PRD and TechSpec, with the answers in its
      `_authorization.md` and PRD;
   6. run `roundfix spec check <slug> --strict --run-verification` on each
      committed plan and fix what it reports;
   7. queue the Specs with `roundfix deliver` in delivery order and follow
      them to merge, and to release only when the user granted release;
   8. before the turn ends: arm a detached Run, the Delivery Queue or a
      scheduled wake-up, or say plainly that work stops and what remains;
      capture each finding that belongs to another project as an Inbox Entry
      in that project's namespace; and report what was done, what is queued or
      armed, which sources were missing and what still waits on the user.
3. MUST append a `## --unattended` section that:
   - lists what to gather: missing information, decisions, and explicit
     authorizations for push, merge and release, and for a
     release-announcement step only when the user's own instructions define
     one;
   - contains, verbatim, the three sentences of `_techspec.md` → API Contract
     5;
   - says each question offers 2-3 options with the recommended one first and
     `(Recommended)` in its label, that without a structured tool the same
     questions come numbered in the chat and the mode waits, and that a
     recommendation, silence, a dismissal or a timeout is never an answer;
   - says an authorization is granted only to a consuming Spec, so a later
     session asks again when it writes that Spec;
   - ends with arming the continuation as `--follow-up` step 8 does when the
     mode is run alone.
4. MUST name the follow-up command and its options exactly as
   `_techspec.md` → API Contract 7 says and MUST NOT name a command or flag
   the tree lacks.
5. MUST add to `skills/roundfix_modes_test.go` the tests
   `TestTheFollowUpModeRunsItsStepsInOrder`,
   `TestTheUnattendedModeAsksWithinTheBaselineLimit` and
   `TestTheRoundfixSkillNamesNoPersonalAssistant` as `_techspec.md` → Testing
   Approach describes them; `TestTheUnattendedModeAsksWithinTheBaselineLimit`
   reads the core module source at
   `../internal/baseline/assets/modules/core.json`. Each top-level test and
   subtest calls `t.Parallel()` first. No existing test changes.
6. MUST run these commands, because this Task edits an owned skill. Before
   them, raise both version fields (`metadata.version` and the top-level
   `version`) of the canonical `SKILL.md` by one patch over the value the tree
   holds when this Task starts, because the record command refuses changed
   content under an unchanged version. The Agent runs them itself before
   handing back, in this order, and again after any later edit of the skill;
   no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
7. MUST NOT write `Jarvis`, `GT3` or any person's or channel's name into the
   skill, change the sections task_02 wrote other than to keep them
   consistent with these two, change the Baseline module, any other owned
   skill, any user guide page, or any text inside a `### QA settlement`
   section.

## Subtasks

- [ ] Write the `--follow-up` section with its eight steps.
- [ ] Write the `--unattended` section with the pinned sentences and answer placement.
- [ ] Add the three tests.
- [ ] Raise the skill version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The `--follow-up` steps run inventory, `--unattended`, triage, authoring, the strict check and the queue in that order, and end with an armed continuation and a report.
- [ ] The `--unattended` section states the Baseline clause's three-question limit and places each answer in the consuming Spec's `_authorization.md` and PRD, or in a Backlog Entry or Finding without authorizations.
- [ ] No Roundfix skill file names a personal assistant or announcement channel; the version is recorded and the mirrors are in sync.

## Consumers

- the modes reference's question limit and answer placement — consumers: `./skills`

## Context

task_02 creates the modes reference, its mirror and the modes test; this Task
extends them, so it declares them as files it writes and runs after task_02.

- instruction: `docs/adr/0320-a-question-call-carries-up-to-three-independent-questions.md`
- instruction: `docs/adr/0095-the-secondbrain-inbox-is-the-single-capture-door.md`
- instruction: `docs/agents/secondbrain.md`
- instruction: `docs/agents/spec-routing.md`
- instruction: `internal/baseline/assets/modules/core.json`
- creates: `.agents/skills/roundfix/references/modes.md`
- creates: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `skills/roundfix_modes_test.go`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for name in TestTheFollowUpModeRunsItsStepsInOrder TestTheUnattendedModeAsksWithinTheBaselineLimit TestTheRoundfixSkillNamesNoPersonalAssistant; do go test -list "^${name}\$" ./skills | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills` — expected: exit 0. Before this Task none of the three tests exists, so the command fails at its first check. Afterwards they exist, the mirrors are in sync, and the skills package passes with the step order, the Baseline limit, the answer placement, the product neutrality and the recorded version.

## References

- `_prd.md` → Goals; User Stories 1, 2, 4 and 5; Core Features 4 and 5; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contracts 5-7; Invariants 5-7; Testing Approach; Build Order 3
- ADR-0320
- ADR-0095
