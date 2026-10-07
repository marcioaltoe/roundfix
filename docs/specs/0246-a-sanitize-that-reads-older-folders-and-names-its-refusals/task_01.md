---
task: task_01
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: pending
type: docs
complexity: low
---

# Task 01: The glossary, the history reference and the Roundfix Skill describe Refused Units and failed-qa

## Overview

Describes the behavior that task_02, task_03 and task_04 build. The places are
`CONTEXT.md`, the history command reference, and the Roundfix Skill's archive
reference with its mirror, so the skill ships with the CLI behavior. This Task
answers the Backlog Entry "History sanitize refuses older archived Specs and
aborts the whole plan" of 2026-10-07. It is verifiable on its own through the
phrases, the mirrors and the recorded skill version.

## Requirements

1. MUST add to `CONTEXT.md`, following the `domain-modeling` skill, the terms
   **Refused Unit** and **Lenient Legacy Reading** in the glossary's
   `**Term**:` / definition / `_Avoid_:` shape, beside **Legacy Archive
   Folder**.
   - **Refused Unit** is a pending History Sanitize Command unit that cannot
     be converted. The plan and the batch name it with its reason and leave it
     untouched, and it does not count toward `--batch <n>` (ADR-0251).
   - **Lenient Legacy Reading** is the reading of a Legacy Archive Folder's
     Task Graph that tolerates and names projection rows outside the graph and
     retired Task types. It is never applied to an active Spec (ADR-0251).
2. MUST revise two existing `CONTEXT.md` entries in place.
   - **Archive Record**: a folder whose newest QA Report failed, with no
     override, receives the disposition `failed-qa`. It keeps that verdict and
     the report name and is never a pass or an override (ADR-0251).
   - **Sanitize Batch**: the batch is the next units that can be converted,
     and Refused Units do not count toward it (ADR-0251).
3. MUST update `docs/user-guide/commands/history.md` to describe:
   - the Lenient Legacy Reading and the `tolerates <folder>: <tolerance>` line;
   - `failed-qa` beside `no-qa`;
   - the `refused <unit>: <reason>` line, the `; <r> unit(s) refused` suffix of
     the plan's first line and the `; <r> unit(s) refused` part of the apply
     confirmation;
   - that `--batch <n>` counts only units it can convert;
   - that apply refuses with exit 2 when every examined unit is refused;
   - which refusals still stop the whole command (`_techspec.md` → Invariants 9
     to 13, Surface Transcripts 1 to 3).

   Every changed sentence MUST agree with the refusal list and the exit codes
   already in the reference.
4. MUST add to `.agents/skills/roundfix/references/archive.md`, under
   `## History Sanitize Command`, that a Refused Unit is listed with its reason
   and skipped, that a Legacy Archive Folder is read leniently while an active
   Spec is not, and that a failed QA without an override is recorded as
   `failed-qa`. The file must contain the exact phrases "Refused Unit" and
   "failed-qa".
5. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, from the value it holds when this Task
   starts to the next free version at record time. Then run `make skills-sync`
   and record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running those commands is an implementation step, never part of
   Verification, and no version is written by hand.
6. MUST NOT change:
   - any other command reference, or the skill's other text;
   - the `### QA settlement` section of any skill;
   - the Baseline guides under `docs/agents/`;
   - another glossary entry.

## Subtasks

- [ ] Add the two glossary terms and revise the two existing ones.
- [ ] Document the tolerated rows, `failed-qa` and the refused lines.
- [ ] Describe the behavior in the skill's archive reference.
- [ ] Raise and record the skill version and sync the mirrors.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Refused Unit** and **Lenient Legacy Reading** and
      names `failed-qa` in **Archive Record**.
- [ ] The history reference documents the `refused` and `tolerates` lines and
      the batch counting.
- [ ] Every skill mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `CONTEXT.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`

## Verification

- `f() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; f CONTEXT.md '**Refused Unit**:'; f CONTEXT.md '**Lenient Legacy Reading**:'; f CONTEXT.md '**Archive Record**:'; f CONTEXT.md '**Sanitize Batch**:'; f CONTEXT.md 'failed-qa'; f docs/user-guide/commands/history.md 'refused <unit>: <reason>'; f docs/user-guide/commands/history.md 'tolerates <folder>: <tolerance>'; f docs/user-guide/commands/history.md 'unit(s) refused'; f docs/user-guide/commands/history.md 'failed-qa'; f .agents/skills/roundfix/references/archive.md 'Refused Unit'; f .agents/skills/roundfix/references/archive.md 'failed-qa'; cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0. Before this Task no file carries the new terms, so the command fails at its first phrase check. After it, the raised version must be recorded for the skill test to pass.

## References

- `_prd.md` → Core Feature 4; Goals
- `_techspec.md` → Build Order 1; API Contract 1; API Contract 2; API Contract 3; Invariant 11; Invariant 12; Invariant 13
- ADR-0251
- ADR-0248
