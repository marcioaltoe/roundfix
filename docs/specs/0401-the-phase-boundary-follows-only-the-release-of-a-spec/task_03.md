---
task: task_03
spec: 0401-the-phase-boundary-follows-only-the-release-of-a-spec
status: pending
type: docs
complexity: medium
---

# Task 03: The authoring skills and the Roundfix modes state the boundary after a Spec's delivery and release

## Overview

The `## Phase boundary` sections of `write-idea`, `write-prd`,
`write-techspec` and `write-tasks`, and the `--follow-up` step 8 and
`--unattended` section of the Roundfix skill's modes, say that an Authoring
Phase ends at no Phase Boundary and that a Spec's Phase Boundary follows its
delivery and release. The glossary's **Phase Boundary** and ADR-0407 say the
same. The skill mirrors and version records are regenerated.

## Requirements

1. MUST replace, in the canonical `.agents/skills/<name>/SKILL.md` of each
   of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, every
   line of `## Phase boundary` after its first paragraph with the text of
   `_techspec.md` → API Contract 4, word for word and the same in all four.
   The first paragraph, which runs `roundfix phase start <phase>`, and the
   `## Declared model` section stay as they are.
2. MUST make, in the canonical `.agents/skills/roundfix/references/modes.md`,
   only the three changes of `_techspec.md` → API Contract 5, word for
   word.
3. MUST then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
4. MUST extend the tests in `skills/phase_boundary_test.go`, on both
   copies, keeping every check they hold today:
   - in `TestTheAuthoringSkillsStateThePhaseBoundary`, the checks
     `no_boundary_after_the_phase` ("When the phase ends, no Phase Boundary
     follows: a Spec's Phase Boundary follows its delivery and release, or
     its merge when it releases nothing."), `the_session_goes_on` ("go on
     to the next phase in this session, handing an Authoring Phase to its
     Author Subagent as its skill's `Declared model` section says, and
     write no Session Handoff for this phase") and
     `the_subagent_session_goes_on` ("the session that started you goes on
     to the next phase"), each exactly once in the section, and
     `no_old_boundary`, which finds "When the phase ends, a Phase Boundary
     follows:" zero times;
   - in `TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary`,
     the checks `the_boundary_follows_delivery_and_release` ("At a Phase
     Boundary, after a Spec is delivered and released, or after it merges
     when it releases nothing,") and `an_authoring_phase_ends_at_no_boundary`
     ("An Authoring Phase ends at no Phase Boundary") in step 8,
     `unattended_boundary_follows_release` ("which follows a Spec's
     delivery and release and never an Authoring Phase") in
     `--unattended`, each exactly once, and `no_boundary_after_each_phase`,
     which finds "after each Authoring Phase" zero times in step 8.
5. MUST revise **Phase Boundary** in `CONTEXT.md` through `domain-modeling`
   to this definition, keeping its `_Avoid_` line: "The point that follows
   the delivery and release of a Spec, or its merge when it releases
   nothing; no Authoring Phase is followed by one. There the Supervisor
   writes the Session Handoff, tells the user what stays armed and starts
   no further phase in the session, and the next session starts with the
   follow-up Skill Mode. Between Authoring Phases the session goes on to
   the next phase, and an Author Subagent ends no session and returns its
   report. Roundfix cannot end a session, so the boundary is an instruction
   that the Boundary Notice recalls", followed by the links to ADR-0407 and
   ADR-0413.
6. MUST add to ADR-0407, after its paragraph "Amended by ADR-0410", one
   paragraph that begins with the bold amendment line naming ADR-0413,
   linked by its file name in the `docs/adr/` directory, and the date
   2026-11-17, in the form of the ADR-0410 amendment, and says that a Phase Boundary follows only the
   delivery and release of a Spec, or its merge when it releases nothing,
   never an Authoring Phase, and that `roundfix phase start` marks an entry
   and prints the Boundary Notice only when the same session's latest
   phase entry is `release` and the new entry is not a `release` of the
   same Spec. MUST raise its `updated_at` to the time of the change and
   change nothing else in it.
7. MUST NOT change any other section of the five skills, the
   `### QA settlement` section of any skill, the generated subagents under
   `.claude/agents/`, `.roundfixrc.yml` or any Vendored Skill.

## Subtasks

- [ ] Rewrite the four `Phase boundary` sections.
- [ ] Make the three changes to the Roundfix modes.
- [ ] Sync the mirrors and record the skill versions.
- [ ] Extend the two boundary tests with the new and absent phrases.
- [ ] Revise **Phase Boundary** and amend ADR-0407.

## Acceptance Criteria

- [ ] No authoring skill and no Roundfix mode asks for a Session Handoff after an Authoring Phase, and each states the boundary after a Spec's delivery and release.
- [ ] The first paragraph of each `Phase boundary` section and every `Declared model` section are unchanged.
- [ ] The glossary and ADR-0407 state the same boundary.
- [ ] Restoring any old sentence makes a named subtest fail.

## Context

- instruction: `docs/adr/0413-a-phase-boundary-follows-only-the-release-of-a-spec.md`
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

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringSkillsStateThePhaseBoundary|TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in no_boundary_after_the_phase the_session_goes_on the_subagent_session_goes_on no_old_boundary the_boundary_follows_delivery_and_release an_authoring_phase_ends_at_no_boundary unattended_boundary_follows_release no_boundary_after_each_phase; do printf '%s\n' "$out" | grep -qF -- "$name (" || { echo "missing check: $name" >&2; exit 1; }; done; if printf '%s\n' "$out" | grep -q -- "--- FAIL"; then printf '%s\n' "$out"; exit 1; fi; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the named checks do not exist; afterwards each passes on both copies, and the skill, mirror and version records pass.
- `for f in .agents/skills/write-idea/SKILL.md .agents/skills/write-prd/SKILL.md .agents/skills/write-techspec/SKILL.md .agents/skills/write-tasks/SKILL.md .agents/skills/roundfix/references/modes.md; do if tr -s '[:space:]' ' ' < "$f" | grep -qE -- "When the phase ends, a Phase Boundary follows:|after each Authoring Phase and after a Spec"; then echo "old boundary in $f" >&2; exit 1; fi; done; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "When the phase ends, no Phase Boundary follows" || { echo 'write-tasks lacks the new boundary' >&2; exit 1; }` — expected: exit 0. Before this Task every authoring skill and the modes state a boundary after each Authoring Phase.
- `grep -F -- '**Phase Boundary**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Phase Boundary\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "no Authoring Phase is followed by one" && { tr -s '[:space:]' ' ' < docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md | grep -qF -- 'Amended by [ADR-0413](0413-a-phase-boundary-follows-only-the-release-of-a-spec.md)' || { printf 'missing phrase in %s: %s\n' docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md 'Amended by ADR-0413' >&2; exit 1; }; }` — expected: exit 0. Before this Task the glossary puts a boundary after each Authoring Phase and ADR-0407 carries no amendment by ADR-0413.

## References

- `_prd.md` → Goals (first, second and fourth bullets); Core Features 3 and 4; Success Metric 2; Glossary
- `_techspec.md` → API Contracts 4 and 5; Testing Approach (Task 3); Build Order 3; ADR Alternatives; Glossary
- [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), the Backlog Entry of 2026-11-16 this Spec adopts
- ADR-0413, ADR-0407, ADR-0401, ADR-0408
