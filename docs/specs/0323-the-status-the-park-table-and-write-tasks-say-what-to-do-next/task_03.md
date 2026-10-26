---
task: task_03
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
status: pending
type: docs
complexity: low
---

# Task 03: write-tasks tells authors to declare a shared test file in every Task

## Overview

Two Tasks of one Wave that each change an undeclared test fake are refused
by the QA gate's precondition with `SC-WAVE-COLLISION` after both have run,
because their settlement commits are collision evidence, while the strict
check was clean before the Run. This Task adds the `write-tasks` rule that a
test helper, fake, fixture or golden two Tasks may each change is declared in
every such Task, revises the **Recorded Path** glossary entry, and records the
skill version.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-two-tasks-of-one-wave-that-change-an-undeclared-test-fake-lose-the-run-at-the-qa-gate.md`](references/2026-10-26-two-tasks-of-one-wave-that-change-an-undeclared-test-fake-lose-the-run-at-the-qa-gate.md)
   of 2026-10-26, from the Fluxus report of 2026-10-23 (`_prd.md` →
   Sources 3), with the option `_techspec.md` → Decisions chose.
2. MUST add the `write-tasks` paragraph of `_techspec.md` → Shared test
   files in the Task Graph to the
   canonical `.agents/skills/write-tasks/SKILL.md` directly after the
   paragraph ending "Declaring every foreseeable path stays the rule.",
   outside any `### QA settlement` section.
3. MUST revise the **Recorded Path** entry of `CONTEXT.md` through the
   `domain-modeling` skill by adding the sentence of that section after its
   last definition sentence, keeping its earlier sentences and its `_Avoid_`
   line.
4. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises both version fields of the `write-tasks` skill and records
   them), then `make skills-sync` again so the mirror matches.
5. MUST NOT change any Go source, the Wave collision rule, any other skill,
   any ADR, `docs/agents/`, the Makefile or the CI workflows.

## Subtasks

- [ ] Add the `write-tasks` paragraph.
- [ ] Revise the **Recorded Path** glossary entry.
- [ ] Sync and record the skill.

## Acceptance Criteria

- [ ] Both `write-tasks` copies carry the shared-test-file rule and say a settlement commit is collision evidence.
- [ ] `CONTEXT.md` **Recorded Path** carries the revised sentence.
- [ ] The `write-tasks` version is recorded and the mirror matches.

## Consumers

None. This Task changes skill text, the glossary and the version record only;
no Go package reads them as a contract except the skill tests it runs.

## Context

- instruction: `internal/spec/collision.go`
- instruction: `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
- instruction: `skills/skills_test.go`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'Its settlement commit is still collision evidence' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Recorded Path**' && for skill in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md; do tr -s '[:space:]' ' ' < "$skill" | grep -qF -- 'A test helper, fake, fixture or golden that two Tasks may each change is a foreseeable path' && tr -s '[:space:]' ' ' < "$skill" | grep -qF -- 'precondition refuses the Run with' || { printf 'missing shared test file rule: %s\n' "$skill" >&2; exit 1; }; done && go test -timeout 10m -count=1 ./skills` — expected: exit 0. Before this Task the glossary sentence is absent, so the first check fails; afterwards both skill copies carry the rule and the owned-skill version is recorded.

## References

- `_prd.md` → Sources 3; Goal 4; Core Feature 3; Success Metric 4; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior 5-6; Shared test files in the Task Graph; Invariant 5; Decisions; Build Order 3
