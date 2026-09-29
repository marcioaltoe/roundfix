---
task: task_04
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: pending
type: docs
complexity: medium
---

# Task 04: The skills, the guides and the glossary describe the three gates as they now behave

## Overview

task_01, task_02 and task_03 change three gates. After them, the owned skills, the user guide and `CONTEXT.md` would describe the old rules. `qa-gate` would tell the QA Agent that the Pull Request row needs a declaration and that every undeclared path fails scope. `write-tasks` would state that every edited path is declared. The Roundfix skill would say an environment-blocked row refuses every `partial`. This Task aligns them: QA Agents and Spec authors read the skills, the maintainer reads the guides, and the Baseline ships the skill mirrors to adopters. The edits to the three skills are bounded by the Spec's `_authorization.md`, and no other governed path may change.

## Requirements

1. MUST edit `.agents/skills/qa-gate/SKILL.md`, the canonical copy, as follows:
   - its table row for a qualifying declared `partial` and its Pull Request journey rule MUST state that the pre-PR Pull Request row never decides a qualifying partial and needs no Unreachable Acceptance declaration. That row is recorded as `blocked (environment: no open Pull Request)` with the Pull Request row named in its provenance.
   - it MUST add a scope rule: when a Requirement audits a Task's changed files against its declarations, a path the Task file lists under `## Recorded paths` counts as declared and is named in that row, while a Governed Path still needs its authorization.
   - it MUST keep every phrase the skill contract tests in `skills/baseline_skill_contract_test.go` require.
2. MUST edit `.agents/skills/write-tasks/SKILL.md` after its declared-path rule to say that the Daemon records a path a Task changed without declaring it under `## Recorded paths` at commit, that the QA scope audit counts it as declared, and that recording discloses a change and reserves nothing. Declaring every foreseeable path stays the rule. It MUST keep every phrase `TestWriteTasksSkillStatesTheDeclaredPathRules` requires.
3. MUST move both version declarations of each of those two skills to `0.0.3`.
4. MUST edit `.agents/skills/roundfix/SKILL.md` so that:
   - its settlement and `archive` eligibility text names the pre-PR Pull Request row exception with the phrase `pre-PR Pull Request row`;
   - its Task commit and `roundfix settle` text names the `## Recorded paths` section;
   - no other command description changes.
5. MUST regenerate the `skills/` mirrors with `make skills-sync` and run `make baseline-digests`. Any path either command rewrites MUST be named in the Result. No other governed path may change.
6. MUST edit `docs/user-guide/commands.md` so that the `roundfix archive` refusal list and the `roundfix qa-report accept` description name the pre-PR Pull Request row exception with the phrase `pre-PR Pull Request row`. It MUST also edit `docs/user-guide/context-driven-development.md` so that its blocked-row paragraph states that exception and its consistency-check text states the ADR horizon with the word `horizon`.
7. MUST add **Recorded Path** to `CONTEXT.md`, and extend **QA Report** with the pre-PR Pull Request row exception and **Spec Consistency Check** with the ADR horizon, using the phrases `pre-PR Pull Request row` and `horizon`.

## Subtasks

- [ ] Align the `qa-gate` and `write-tasks` skills and move their versions.
- [ ] Align the Roundfix skill's eligibility and Task commit text.
- [ ] Regenerate the mirrors and the Baseline digests.
- [ ] Align the two user guides and the glossary.

## Acceptance Criteria

- [ ] The canonical `qa-gate` skill names the pre-PR Pull Request row exception and the `## Recorded paths` scope rule, and declares version `0.0.3`.
- [ ] The canonical `write-tasks` skill names the `## Recorded paths` record and keeps its declared-path phrases, and declares version `0.0.3`.
- [ ] The canonical Roundfix skill names the `pre-PR Pull Request row` and the `## Recorded paths` section.
- [ ] Each mirror is byte-identical to its canonical skill, and `make skills-sync-check` exits `0`.
- [ ] The commands guide, the context-driven development guide and `CONTEXT.md` carry the exception, and the development guide and `CONTEXT.md` carry the horizon; `CONTEXT.md` defines **Recorded Path**.

## Context

- instruction: `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
- instruction: `docs/adr/0167-the-pre-pr-pull-request-row-never-decides-a-qualifying-partial.md`
- instruction: `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `CONTEXT.md`

## Verification

- `for pair in ".agents/skills/qa-gate/SKILL.md|never decides a qualifying" ".agents/skills/qa-gate/SKILL.md|## Recorded paths" ".agents/skills/write-tasks/SKILL.md|## Recorded paths" ".agents/skills/write-tasks/SKILL.md|reserves nothing" ".agents/skills/roundfix/SKILL.md|pre-PR Pull Request row" ".agents/skills/roundfix/SKILL.md|## Recorded paths" "docs/user-guide/commands.md|pre-PR Pull Request row" "docs/user-guide/context-driven-development.md|pre-PR Pull Request row" "docs/user-guide/context-driven-development.md|horizon" "CONTEXT.md|**Recorded Path**" "CONTEXT.md|pre-PR Pull Request row" "CONTEXT.md|horizon"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && grep -q "^version: 0.0.3$" .agents/skills/qa-gate/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/qa-gate/SKILL.md && grep -q "^version: 0.0.3$" .agents/skills/write-tasks/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/write-tasks/SKILL.md && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null && diff -r .agents/skills/write-tasks skills/write-tasks >/dev/null && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check && out="$(go test -count=1 -v -run "^(TestWriteTasksSkillStatesTheDeclaredPathRules|TestProjectConstraintQAGate|TestProjectConstraintJourney|TestToolingAuthorizationJourney|TestLegacySpecConstraintExemption|TestAuthorialSkillSync)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestWriteTasksSkillStatesTheDeclaredPathRules TestProjectConstraintQAGate TestProjectConstraintJourney TestToolingAuthorizationJourney TestLegacySpecConstraintExemption TestAuthorialSkillSync; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task no skill, guide or glossary names the recorded section, the pre-PR row or the horizon, and both skills declare `0.0.2`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Skills, guides and glossary
- `_prd.md` → Goal 4; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 4
- `_techspec.md` → Testing Approach 4

## Result
