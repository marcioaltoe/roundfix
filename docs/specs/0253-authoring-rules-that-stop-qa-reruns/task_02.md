---
task: task_02
spec: 0253-authoring-rules-that-stop-qa-reruns
status: pending
type: docs
complexity: medium
---

# Task 02: The authoring skills and the repository guide state the rules that stop QA reruns

## Overview

The authoring skills tell an author to assert a transcript's "command, output,
and exit text" and to rest one row on outside evidence, but not to copy every
line or where that evidence must live. The rules that hold only for Roundfix
(the suite guard, the `### QA settlement` section, clause retention, skill and
module records, derived files and the Governed Path probe) live in the
operator's private briefing (findings B09 and B12 and the briefing section of
the Baseline audit of 2026-10-08,
`docs/references/2026-10-08-baseline-audit.md`). This Task writes the texts of
`_techspec.md` → Skill and document texts into `write-tasks`, the
`write-techspec` concrete contracts guide and
`docs/agents/specific-repository.md` (ADR-0258). It is verifiable on its own
through phrase checks, the mirrors and the skill version record.

This is an authorized tooling Task. It may change only the files in its
Context and this Task file. It starts from the merge of Spec 0252, which
deletes one bullet of `docs/agents/specific-repository.md`.

## Requirements

1. MUST answer finding B09 of the Baseline audit of 2026-10-08 by making the
   four `write-tasks` changes of `_techspec.md` → Skill and document texts in
   `.agents/skills/write-tasks/SKILL.md`. The first changes the sentence of
   `## Surface Transcripts in Tasks` that wraps across two lines.
2. MUST add the paragraph of `_techspec.md` → Skill and document texts to
   `## Surface transcript blocks` in
   `.agents/skills/write-techspec/references/concrete-contracts.md`. The guide
   MUST keep exactly one Claim Receipt and one transcript, which
   `TestTheGuideExamplesAreReadByTheCheck` proves.
3. MUST answer finding B12 of the Baseline audit of 2026-10-08 by inserting
   the seven bullets of `_techspec.md` → Skill and document texts into
   `docs/agents/specific-repository.md`, before the bullet that begins "The
   release plan the baseline requires". No other byte of that file changes.
4. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` again, so that the version fields of both skills
   rise and every mirror matches its canonical copy.
5. MUST NOT change any `### QA settlement` section, any other skill, any
   Baseline module or any production Go file.

## Subtasks

- [ ] Change `write-tasks` and the concrete contracts guide.
- [ ] Insert the repository rules.
- [ ] Record both skill versions and sync the mirrors.

## Acceptance Criteria

- [ ] `write-tasks` names the whole-transcript test, the `exit status <n>`
      line, reachable outside evidence, `requires` and the backtick rule.
- [ ] The concrete contracts guide tells authors to copy every line from a
      real run and still holds one receipt and one transcript.
- [ ] `docs/agents/specific-repository.md` carries the seven rules.
- [ ] Both skills record a raised version, and the mirrors match.

## Context

- instruction: `docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/concrete-contracts.md`
- interface: `skills/write-techspec/references/concrete-contracts.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/agents/specific-repository.md`

## Verification

- `for phrase in "asserts the whole transcript in one named test" "A line the reader might skip still counts" "Never cite an artifact that exists only on the operator's machine" "A cross-Spec prerequisite is declared." "A Verification command never contains a backtick"; do tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "$phrase" || { printf 'missing phrase in write-tasks: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/write-techspec/references/concrete-contracts.md | grep -qF -- "Copy each line from a run of the real command, not from memory." || { printf 'the concrete contracts guide lacks the copy rule\n' >&2; exit 1; }; for phrase in "spawning test packages install the suite guard" "New skill text goes under its own heading." "so prefer extending an existing clause" "measured on a scratch copy" "Suiteguard refuses the module and record writes of an undeclared command." "probe that writes nothing to the repository"; do tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- "$phrase" || { printf 'missing phrase in specific-repository.md: %s\n' "$phrase" >&2; exit 1; }; done; for file in write-tasks/SKILL.md write-techspec/SKILL.md write-techspec/references/concrete-contracts.md; do cmp -s ".agents/skills/$file" "skills/$file" || { printf 'mirror differs: %s\n' "$file" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestTheGuideExamplesAreReadByTheCheck|TestTheConcreteContractGuideShipsAtTheHorizonPath)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestTheGuideExamplesAreReadByTheCheck TestTheConcreteContractGuideShipsAtTheHorizonPath; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task `write-tasks` lacks the whole-transcript sentence, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goals 1-3 and 5; User Stories 1-3 and 5; Core Features 3 and 5; Success Metric 3
- [_techspec.md](_techspec.md) — API Contract 5; Skill and document texts; Version changes; Invariant 9; Testing Approach; Build Order 2
- ADR-0258
