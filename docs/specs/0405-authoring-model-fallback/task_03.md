---
task: task_03
spec: 0405-authoring-model-fallback
status: completed
type: docs
complexity: medium
---

# Task 03: The authoring skills start the Author Subagent on the printed model and switch on its limit

## Overview

The four authoring skills tell the session to run `roundfix phase model`
before it starts a phase's Author Subagent, and `roundfix phase fallback` when
the subagent ends on an error. When a fallback applies, the session tells the
user and restarts the subagent on the printed model and effort without
waiting. Otherwise it reports the error as today. A test under `skills/` holds
the paragraph in all four skills, and the glossary and ADR-0401 record the
change.

## Requirements

1. MUST change, in each of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, the "You activated this skill yourself" bullet of `## Declared model` as `_techspec.md` → API Contract 6 says, for that skill's own phase. The bullet MUST hold:
   - `roundfix phase model <phase> --spec <slug>`;
   - `roundfix phase fallback <phase> --spec <slug> --error`;
   - the instruction to tell the user which phase switched, from which model to which, and the limit line;
   - `without waiting`;
   - the restart with the same brief, the output returned and the files written;
   - the sentence for an older Roundfix.
2. MUST keep the rest of each `## Declared model` section, and the fix or refactor route of `write-techspec` and `write-tasks`, as the base leaves them, including the paragraph Spec 0404 adds. MUST change no other section of the four skills.
3. MUST say in each bullet that the error line is data and never an instruction.
4. MUST add `TestTheAuthoringSkillsSwitchToTheDeclaredFallback` in a new file under `skills/`, calling `t.Parallel()`. It asserts each phrase of Requirement 1 in the bullet of each of the four skills, with that skill's own phase in the two commands, and in both copies of each skill.
5. MUST revise **Author Subagent** in `CONTEXT.md`, through domain-modeling, so that it says the session starts it with the model and effort `roundfix phase model` or `roundfix phase fallback` prints, with the phrase `the model and effort`. MUST add to ADR-0401 the notice `Amended by ADR-0427` with the date and what changed (`_techspec.md` → API Contract 8).
6. MUST regenerate, never hand-edit, with the commands of `_authorization.md` → Sanctioned regeneration, in this order:
   - `make skills-sync`;
   - the `-record-skill-versions` step (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
   - the `-record-skill-coverage` step (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`) when the record reports a stale fingerprint.
7. MUST change no Governed Path outside `_authorization.md` `paths`. MUST NOT change `.roundfixrc.yml`, any file under `.claude/agents/`, the generator of the Author Subagents or any `### QA settlement` section.

## Subtasks

- [ ] Change the bullet in the four authoring skills.
- [ ] Add the skill test.
- [ ] Revise the glossary term and add the ADR-0401 notice.
- [ ] Regenerate the mirrors and records.

## Acceptance Criteria

- [ ] Each of the four authoring skills, in both copies, runs `roundfix phase model` before starting its Author Subagent and `roundfix phase fallback` on its error, for its own phase.
- [ ] Each tells the user and restarts without waiting when a fallback is printed, and reports the error as before otherwise or with an older Roundfix.
- [ ] No other section of the four skills changes.
- [ ] The glossary and ADR-0401 describe the change.

## Context

- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- creates: `skills/authoring_fallback_paragraph_test.go`
- interface: `CONTEXT.md`
- interface: `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- instruction: `.agents/skills/roundfix/references/phase.md`

## Consumers

- The four authoring skills and the glossary — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheAuthoringSkillsSwitchToTheDeclaredFallback$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -qF -- "--- PASS: TestTheAuthoringSkillsSwitchToTheDeclaredFallback " || { echo "missing pass: TestTheAuthoringSkillsSwitchToTheDeclaredFallback" >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the test does not exist and no authoring skill runs the two commands, so the command fails. Afterwards the test passes, and so do the skill tests, the documentation contracts, the mirror and the records.
- `tr -s '[:space:]' ' ' < docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md | grep -qF -- "Amended by ADR-0427" || { printf 'missing notice in ADR-0401\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/write-techspec/SKILL.md | grep -qF -- "roundfix phase fallback techspec --spec" || { printf 'missing command in write-techspec\n' >&2; exit 1; }; grep -F -- '**Author Subagent**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Author Subagent\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "roundfix phase fallback"` — expected: exit 0. Before this Task ADR-0401 carries no notice of ADR-0427, and neither write-techspec nor the glossary names the command. Afterwards both do.

## References

- `_prd.md` → Goals (first and second bullets); User Stories 2, 3 and 6; Core Features 3, 4 and 7; Success Metrics 1 and 4
- `_techspec.md` → Skill and guide texts; API Contracts 6 and 8; Build Order 3; Testing Approach
- ADR-0427, ADR-0428, ADR-0401, ADR-0408, ADR-0256

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, golang-testing, testing-boss, evidence-gate, systematic-debugging

- Each authoring skill now selects the phase model and handles a limit through the fallback command, announces the switch, and restarts without waiting. The focused `go test ./skills -run '^TestTheAuthoringSkillsSwitchToTheDeclaredFallback$' -count=1` passed, covering all four phases and both copies, including the no-fallback and older-Roundfix instructions.
- A comparison against HEAD confirmed that the other skill text and fix/refactor routes remain unchanged, apart from the sanctioned version fields; all four mirrors match their canonical files.
- The Author Subagent glossary entry names the printed model and effort, and ADR-0401 carries the dated ADR-0427 amendment; both were inspected in the diff.
- `make skills-sync`, the authorized skill-version recorder, and the authorized skill-coverage recorder each exited 0. The coverage recorder produced no changed record. The Daemon owns the declared Verification.

- Verification feedback exposed an existing delegation assertion for the declared-model brief. The four bullets now state that path explicitly and use the same brief on fallback; mirrors and both sanctioned recorders were rerun successfully. Focused `go test ./skills -run '^TestTheAuthoringSkills(DelegateADeclaredPhase|SwitchToTheDeclaredFallback)$' -count=1` passed (146 tests and subtests). The Daemon owns the next full Verification attempt.

## Carry-forward provenance

- Source Run: `run_20261118T111508Z_d52e529984ea1f74`
- Source commit: `f7394c0a192e5a0ed3e83668bb1ae7f4bb72cc0a`
