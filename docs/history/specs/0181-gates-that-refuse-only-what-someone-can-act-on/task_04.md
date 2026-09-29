---
task: task_04
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: completed
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

Implementation evidence:

- `.agents/skills/qa-gate/SKILL.md` now declares version `0.0.3`, documents the
  pre-PR Pull Request row exception in settlement and Pull Request journeys,
  and counts `## Recorded paths` as declared for ordinary scope rows while
  retaining Governed Path authorization.
- `.agents/skills/write-tasks/SKILL.md` now declares version `0.0.3` and states
  that the Daemon records undeclared paths under `## Recorded paths`, that the
  QA scope audit counts them as declared, and that recording discloses a change
  and reserves nothing.
- `.agents/skills/roundfix/SKILL.md` now names the pre-PR Pull Request row in
  settlement and archive eligibility text and names `## Recorded paths` in
  Task-settlement commit guidance.
- `docs/user-guide/commands.md`,
  `docs/user-guide/context-driven-development.md`, and `CONTEXT.md` now carry
  the pre-PR Pull Request row exception; the latter two carry the ADR horizon,
  and `CONTEXT.md` defines **Recorded Path**.
- `make skills-sync` passed and rewrote only the three distributed mirrors:
  `skills/qa-gate/SKILL.md`, `skills/write-tasks/SKILL.md`, and
  `skills/roundfix/SKILL.md`. `make baseline-digests` passed with
  `changed:false` and rewrote no path.

Acceptance-criterion evidence:

- Canonical `qa-gate`: focused phrase/version scan passed for the pre-PR row,
  `## Recorded paths`, and both `0.0.3` declarations.
- Canonical `write-tasks`: focused phrase/version scan passed for
  `## Recorded paths`, `reserves nothing`, and both `0.0.3` declarations.
- Canonical Roundfix: focused phrase scan passed for the pre-PR row and
  `## Recorded paths`.
- Mirrors: `cmp` passed for all three canonical/mirror pairs; `git diff --check`
  passed. The declared `make skills-sync-check` command was not run because it
  is part of Daemon-owned Task Verification.
- Guides and glossary: focused phrase scan passed for the exception, horizon,
  and `**Recorded Path**` definition.

Focused checks:

- `make skills-sync`: passed.
- `make baseline-digests`: passed; `changed:false`.
- `cmp` for all three skill pairs: passed.
- `git diff --check`: passed.
- With a task-scoped `GOCACHE`, the focused process-table tests passed with
  host permission: `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`
  and `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner`.
- `make verify-incremental`: failed first on host process-table permission and
  on `TestSettlementGuidanceIsOneTable`; the process-table failures were
  cleared by the focused elevated rerun above. The settlement contract still
  compares the required changed `qa-gate` table with the unchanged
  upstream-managed `.agents/skills/archive-spec/SKILL.md`. Updating that path
  is outside this Task's authorization; follow up with a separately authorized
  alignment of that shared contract.
- The focused `TestSettlementGuidanceIsOneTable` rerun reproduced that same
  out-of-slice mismatch. No status field was changed, no other Task file or
  Task Graph manifest was edited, and no commit, push, or pull request was
  created.

## Carry-forward provenance

- Source Run: `run_20260929T170240Z_b8558d4fb9103028`
- Source commit: `65fd04b31ab776f093cd475131545413376c3965`
