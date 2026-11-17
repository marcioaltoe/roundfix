---
task: task_03
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
status: pending
type: docs
complexity: medium
---

# Task 03: The authoring skills, the Roundfix modes and the glossary say no phase of a Spec ends the session

## Overview

The `## Phase boundary` sections of `write-idea`, `write-prd`,
`write-techspec` and `write-tasks`, and the `--follow-up` step 8 and
`--unattended` section of the Roundfix skill's modes, say that no phase of
a Spec ends the session, with or without a person present, and that the
session writes the Session Handoff and ends only when the user asks or the
work stops. The glossary redefines **Phase Boundary** and
**Session Handoff** and drops the entry of the notice the phase command no
longer prints, and ADR-0407 is amended by ADR-0414. The skill mirrors and
version records are regenerated.

## Requirements

1. MUST replace, in the canonical `.agents/skills/<name>/SKILL.md` of each
   of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, every
   line of `## Phase boundary` after its first paragraph with the text of
   `_techspec.md` → API Contract 4, word for word and the same in all four.
   The first paragraph, which runs `roundfix phase start <phase>`, and the
   `## Declared model` section stay as they are.
2. MUST make, in the canonical `.agents/skills/roundfix/references/modes.md`,
   only the two changes of `_techspec.md` → API Contract 5, word for word.
3. MUST then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
4. MUST change the tests in `skills/phase_boundary_test.go`, on both
   copies, letting each check name the occurrences it wants (one, or zero
   for an absence):
   - in `TestTheAuthoringSkillsStateThePhaseBoundary`, keep
     `records_its_own_phase`, `an_unknown_command_continues`,
     `the_author_subagent_returns`, `the_session_writes_the_handoff`,
     `the_user_is_told` and `the_next_session_follows_up`; drop
     `no_further_phase`; add, each once in the section,
     `no_phase_ends_the_session` ("When the phase ends, the session goes
     on: no phase of a Spec ends the session, whether a person is present
     or not."), `the_session_goes_on` ("go on in this session with the next
     phase or the next queued work, handing an Authoring Phase to its
     Author Subagent as its skill's `Declared model` section says, and
     write no Session Handoff for this phase"),
     `the_subagent_session_goes_on` ("the session that started you goes on
     with the next phase") and `the_session_ends_when_asked_or_stopped`
     ("End the session only when the user asks or the work stops because
     nothing remains or what remains waits on the user"); and add, each
     zero times, `no_old_boundary` ("When the phase ends, a Phase Boundary
     follows") and `no_stop_after_the_phase` ("start no further phase");
   - in `TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary`,
     keep every check but `no_further_phase`, `a_wake_up_never_crosses`,
     `unattended_ends_at_the_boundary` and `unattended_wake_up`, which it
     drops; add, each once in step 8, `no_phase_ends_the_session` ("No
     phase of a Spec ends the session, whether a person is present or
     not"), `the_session_follows_the_queue` ("after an Authoring Phase, a
     delivery or a release, go on with the next queued work"),
     `the_session_ends_when_asked_or_stopped` ("End the session only when
     the user asks or the work stops because nothing remains or what
     remains waits on the user") and `a_wake_up_returns_to_the_session` ("A
     scheduled wake-up returns to the same session with its whole
     context."); add, once in `--unattended`,
     `unattended_no_phase_ends_the_session` ("With nobody present no phase
     of a Spec ends the session either"); and add, each zero times,
     `no_boundary_after_each_phase` ("after each Authoring Phase") and
     `no_stop_after_the_phase` ("start no further phase") in step 8 and
     `unattended_no_boundary` ("ends at a Phase Boundary") in
     `--unattended`.
5. MUST revise **Phase Boundary** in `CONTEXT.md` through
   `domain-modeling` to this definition, keeping its `_Avoid_` line: "The
   change from one Supervisor Phase to the next, which `roundfix phase
   start` records. It ends no session, with or without a person present:
   the Supervisor goes on with the next phase or the next queued work, and
   an Author Subagent returns its report. A Supervisor session ends only
   when the user asks or the work stops because nothing remains or what
   remains waits on the user; there the Supervisor writes the Session
   Handoff and tells the user what stays armed", followed by the links to
   ADR-0407 and ADR-0414.
6. MUST revise **Session Handoff** in `CONTEXT.md` through
   `domain-modeling`, replacing only its last sentence, "At a Phase
   Boundary the Supervisor writes one …", with: "A Supervisor session
   writes one, with the Spec, the phase that ended, the next phase, each
   decision still open and what stays armed, only when it ends because the
   user asks or the work stops; no phase of a Spec ends it", followed by
   the links to ADR-0407 and ADR-0414. The rest of the entry and its
   `_Avoid_` line stay.
7. MUST remove the glossary entry of the notice from `CONTEXT.md` through
   `domain-modeling`: its definition line, its `_Avoid_` line and the blank
   line after them, so that no line of `CONTEXT.md` names it.
8. MUST add to ADR-0407, after its paragraph "Amended by ADR-0410", one
   paragraph that begins with the bold amendment line naming ADR-0414,
   linked by its file name in the `docs/adr/` directory, and the date
   2026-11-17, in the form of the ADR-0410 amendment, and says that no
   phase of a Spec ends the Supervisor session, which writes the Session
   Handoff and ends only when the maintainer asks or the work stops; that
   `roundfix phase start` reads no earlier entry, writes `boundary` as
   `false` without `previous_phase` or `previous_spec`, and prints no
   Boundary Notice; and that the Phase Record, its entries and the refusal
   of nothing stay. MUST raise its `updated_at` to the time of the change
   and change nothing else in it.
9. MUST NOT change any other section of the five skills, the
   `### QA settlement` section of any skill, the generated subagents under
   `.claude/agents/`, `.roundfixrc.yml`, any Vendored Skill, or any ADR but
   ADR-0407.

## Subtasks

- [ ] Rewrite the four `Phase boundary` sections.
- [ ] Make the two changes to the Roundfix modes.
- [ ] Sync the mirrors and record the skill versions.
- [ ] Change the two boundary tests with the new and absent phrases.
- [ ] Revise **Phase Boundary** and **Session Handoff**, remove the notice's entry, and amend ADR-0407.

## Acceptance Criteria

- [ ] No authoring skill and no Roundfix mode ends the session or asks for a Session Handoff after an Authoring Phase or a Spec's delivery and release, and each says the session ends only when the user asks or the work stops.
- [ ] The first paragraph of each `Phase boundary` section and every `Declared model` section are unchanged.
- [ ] The glossary and ADR-0407 state the same rule, and the glossary no longer defines the notice.
- [ ] Restoring any old sentence makes a named check fail.

## Context

- instruction: `docs/adr/0414-a-supervisor-session-ends-only-when-the-maintainer-asks.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/phase_boundary_test.go`
- interface: `CONTEXT.md`
- interface: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`

## Consumers

- The four authoring skills, the Roundfix modes and their version record — consumers: `./skills`, `./internal/docscontract`
- The glossary and ADR-0407 — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringSkillsStateThePhaseBoundary|TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in no_phase_ends_the_session the_session_goes_on the_subagent_session_goes_on the_session_ends_when_asked_or_stopped no_old_boundary no_stop_after_the_phase the_session_follows_the_queue a_wake_up_returns_to_the_session unattended_no_phase_ends_the_session no_boundary_after_each_phase unattended_no_boundary; do printf '%s\n' "$out" | grep -qE -- "--- PASS: .*/$name " || { echo "missing check: $name" >&2; exit 1; }; done; if printf '%s\n' "$out" | grep -q -- "--- FAIL"; then printf '%s\n' "$out"; exit 1; fi; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the named checks do not exist; afterwards each passes on both copies, and the skill, mirror, version and glossary records pass.
- `for f in .agents/skills/write-idea/SKILL.md .agents/skills/write-prd/SKILL.md .agents/skills/write-techspec/SKILL.md .agents/skills/write-tasks/SKILL.md .agents/skills/roundfix/references/modes.md; do if tr -s '[:space:]' ' ' < "$f" | grep -qE -- "When the phase ends, a Phase Boundary follows|after each Authoring Phase and after a Spec|start no further phase|ends at a Phase Boundary"; then echo "old boundary in $f" >&2; exit 1; fi; done; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "no phase of a Spec ends the session, whether a person is present or not" || { echo 'write-tasks lacks the new rule' >&2; exit 1; }` — expected: exit 0. Before this Task every authoring skill and the modes end the session at a Phase Boundary.
- `grep -F -- '**Phase Boundary**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Phase Boundary\*\*(:| —( |$))' && grep -F -- '**Session Handoff**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Session Handoff\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "It ends no session, with or without a person present" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "only when it ends because the user asks or the work stops" && ! { tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "Boundary Notice"; } && { tr -s '[:space:]' ' ' < docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md | grep -qF -- 'Amended by [ADR-0414](0414-a-supervisor-session-ends-only-when-the-maintainer-asks.md)' || { printf 'missing phrase in %s: %s\n' docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md 'Amended by ADR-0414' >&2; exit 1; }; }` — expected: exit 0. Before this Task the glossary ends the session at each Phase Boundary and defines the notice, and ADR-0407 carries no amendment by ADR-0414.

## References

- `_prd.md` → Goals (first, second and fourth bullets); Core Features 3 and 4; Success Metric 2; Glossary
- `_techspec.md` → API Contracts 4 and 5; Testing Approach (Task 3); Build Order 3; ADR Alternatives; Glossary
- [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), the Backlog Entry of 2026-11-16 this Spec adopts
- ADR-0414, ADR-0407, ADR-0401, ADR-0408
