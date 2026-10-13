---
task: task_03
spec: 0258-authoring-that-proves-what-it-claims
status: completed
type: docs
complexity: medium
---

# Task 03: The authoring skills state the four proof rules

## Overview

`write-techspec` asks for exploration but no evidence, and `write-tasks` forbids
a whole-package suite without naming the suites a contract change breaks. This
Task writes the texts of `_techspec.md` → Skill and document texts into both
skills and both templates, raises both skill versions, syncs the mirrors and
adds a skill contract test that reads the four rules. It is verifiable on its
own through phrase checks, the mirrors, the version record and the new test.

This is an authorized tooling Task. It may change only the files in its
Context and this Task file.

## Requirements

1. MUST answer the four gaps of the adopted Backlog Entry of 2026-10-12,
   [Spec authoring does not prove what it claims](references/2026-10-12-spec-authoring-does-not-prove-what-it-claims.md),
   by adding the `## Proof of what the Spec claims` section of
   `_techspec.md` → Skill and document texts to `write-techspec`, after
   `## Concrete contracts`, and the `## Consumers of a shared contract`
   section to `write-tasks`, after `## Surface Transcripts in Tasks`.
2. MUST add to `write-tasks` the exception sentence on whole-package suites
   and the Process step 1 sentence on the Coverage Cross-Read, both of
   `_techspec.md` → Skill and document texts.
3. MUST add `## Absence Evidence` and `## ADR Alternatives` sections after
   `## Coverage Map` in the TechSpec template, and the line
   `Deviations from the PRD: None.` to its Coverage Map guidance, and a
   `## Consumers` section before `## Verification` in the Task template.
4. MUST add `TestAuthoringSkillsStateTheProofRules` to the skill contract
   tests. It reads both canonical skills and both templates and fails when any
   rule, section heading or finding code of this Task is missing.
5. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` again, so that both skills record a raised version
   and every mirror matches its canonical copy.
6. MUST NOT change any `### QA settlement` section, any other skill, any
   Baseline module or any production Go file.

## Subtasks

- [ ] Add the two skill sections and the two sentences.
- [ ] Extend both templates.
- [ ] Add the skill contract test.
- [ ] Record both skill versions and sync the mirrors.

## Acceptance Criteria

- [ ] `write-techspec` names Absence Evidence, ADR Alternatives and the
      Coverage Cross-Read with their codes.
- [ ] `write-tasks` names the Consumer Declaration, the Integration
      Checkpoint, the whole-package exception and its code.
- [ ] Both templates carry the new sections.
- [ ] Both skills record a raised version, and the mirrors match.

## Consumers

- The write-techspec and write-tasks skill texts and templates — consumers: `./skills`

## Context

- instruction: `docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/baseline_skill_contract_test.go`

## Verification

- `for phrase in "## Proof of what the Spec claims" "SC-ABSENCE-UNPROVEN" "SC-ADR-ALTERNATIVES-UNRELATED" "SC-COVERAGE-UNREAD" "Deviations from the PRD:"; do tr -s '[:space:]' ' ' < .agents/skills/write-techspec/SKILL.md | grep -qF -- "$phrase" || { printf 'missing phrase in write-techspec: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "## Absence Evidence" "## ADR Alternatives" "Deviations from the PRD: None."; do tr -s '[:space:]' ' ' < .agents/skills/write-techspec/references/techspec-template.md | grep -qF -- "$phrase" || { printf 'missing phrase in the TechSpec template: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "## Consumers of a shared contract" "SC-CONSUMERS-UNPROVEN" "Integration Checkpoint" "A suite that consumes a Shared Contract the Task changes is the exception"; do tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "$phrase" || { printf 'missing phrase in write-tasks: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/references/task-template.md | grep -qF -- "## Consumers" || { printf 'the Task template lacks the Consumers section\n' >&2; exit 1; }; for file in write-tasks/SKILL.md write-tasks/references/task-template.md write-techspec/SKILL.md write-techspec/references/techspec-template.md; do cmp -s ".agents/skills/$file" "skills/$file" || { printf 'mirror differs: %s\n' "$file" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestAuthoringSkillsStateTheProofRules|TestWriteTasksSkillStatesTheDeclaredPathRules)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestAuthoringSkillsStateTheProofRules TestWriteTasksSkillStatesTheDeclaredPathRules; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task `write-techspec` lacks the proof section, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goals 1-4; User Stories 1-4; Core Features 1-4; Success Metric 2
- [_techspec.md](_techspec.md) — Skill and document texts; Invariants 1-7; Testing Approach; Build Order 3
- ADR-0265; ADR-0258

## Result

Implemented the authoring proof rules in the canonical `write-techspec` and
`write-tasks` skills and templates, added the
`TestAuthoringSkillsStateTheProofRules` contract test, and recorded the raised
skill versions. The canonical and distributed copies are synchronized.

Focused implementation evidence:

- `rtk proxy make skills-sync` — completed before version recording and after
  the final content/version update.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed and recorded `write-tasks` at `0.0.14` and `write-techspec` at `0.0.11`.
- `rtk proxy go test -count=1 ./skills -run '^TestAuthoringSkillsStateTheProofRules$'` — passed.
- `rtk proxy git diff --check` — passed.
- `cmp` checks for both skill files and both templates — passed for all four
  canonical/mirror pairs.

Acceptance evidence:

- `write-techspec` contains the Proof of what the Spec claims section with
  Absence Evidence, ADR Alternatives, Coverage Cross-Read, and all three
  finding codes.
- `write-tasks` contains the Consumer Declaration, Integration Checkpoint,
  whole-package exception, Coverage Cross-Read process rule, and
  `SC-CONSUMERS-UNPROVEN`.
- Both templates contain the requested sections and Coverage Map guidance.
- Both skill versions are raised and the final sync leaves every canonical
  copy byte-identical to its mirror.

## Carry-forward provenance

- Source Run: `run_20261013T122525Z_1c1258fdf2b03bd6`
- Source commit: `d148925f8ec07fa7df5b03b7d1713b30952ac957`
