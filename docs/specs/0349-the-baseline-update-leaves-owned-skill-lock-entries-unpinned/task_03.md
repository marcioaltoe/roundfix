---
task: task_03
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
status: completed
type: docs
complexity: low
---

# Task 03: The Roundfix skill and the glossary state the unpinned owned lock entry

## Overview

The Roundfix skill's `baseline` and `setup` references still say the bundle
install keeps every lock field but `computedHash`, and nothing in the skill
or the glossary names the new Doctor warning. This Task writes the rule of
ADR-0359 into both references and into `CONTEXT.md`, raises the Roundfix
skill's version and re-records the derived records. It closes the marcioaltoe-skills report of
2026-11-04 (ai-memory message `01a1a5f7`, adopted as
[references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md](references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md)).

## Requirements

1. MUST state the rule of ADR-0359 in the canonical Roundfix references, as
   the guides task_01 and task_02 changed state it:
   - `.agents/skills/roundfix/references/baseline.md` (the lock sentence,
     `skills.lockUnpinned`, `- lock ref removed <name>`);
   - `.agents/skills/roundfix/references/setup.md` (the lock sentence,
     `Removed skills-lock.json ref: <name>`, the `DR-SKILL-LOCK-STALE`
     warning and its next action);
   - each new text under its own sentences, never inside a `### QA settlement`
     section. Then run `make skills-sync` and record the version with
     `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     which raises both version fields of the Roundfix `SKILL.md` pair; never
     hand-edit a mirror, a version field or the record.
2. MUST add `TestTheGuidesStateTheUnpinnedOwnedLockEntry` in
   `skills/owned_lock_guidance_test.go` as `_techspec.md` → Testing Approach
   describes, calling `t.Parallel()` first.
3. MUST add **Stale Owned Lock Entry** to `CONTEXT.md` through
   domain-modeling, next to **Setup Snapshot Difference**, as a
   `**Stale Owned Lock Entry**:` line. It is a `skills-lock.json` entry of an
   installed Roundfix-owned skill that carries a `ref` or a `computedHash`
   other than the installed copy's, so `skills update` would replace the
   installed copy with the one the entry pins. `roundfix doctor` warns
   `DR-SKILL-LOCK-STALE`, and the bundle install of `roundfix baseline update`
   or `roundfix skills install` rewrites the entry with the installed hash
   and no `ref`. Link ADR-0359. It keeps this repository's glossary contract.
4. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last skill edit, and never hand-edit
   it.
5. MUST NOT touch any code outside `skills/owned_lock_guidance_test.go`, any
   command guide, `docs/references/skill-coverage.json`, the release runbook,
   any skill other than the Roundfix skill, any guide under `docs/agents/`,
   the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Update the two canonical Roundfix references, sync the mirrors and record the version.
- [ ] Add the guidance test and the glossary entry.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the Roundfix references state the removed `ref`, the new report lines, `skills.lockUnpinned` and `DR-SKILL-LOCK-STALE`.
- [ ] `CONTEXT.md` defines **Stale Owned Lock Entry**.
- [ ] The Roundfix skill's version and the derived records are recorded by their commands.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/owned_lock_guidance_test.go`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix skill references, their mirrors and the owned-skill version record — consumers: `./skills`
- The glossary and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheGuidesStateTheUnpinnedOwnedLockEntry$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheGuidesStateTheUnpinnedOwnedLockEntry ' || { echo 'missing pass: TestTheGuidesStateTheUnpinnedOwnedLockEntry' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the test does not exist; afterwards it passes with the whole skills package, the owned-skill version record and mirror checks among them.
- `grep -F -- '**Stale Owned Lock Entry**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Stale Owned Lock Entry\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the glossary has no such definition; afterwards it does and the documentation contracts, the glossary contract and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 5; Glossary
- `_techspec.md` → API Contract 5; Testing Approach; Risks & Considerations; Glossary; Build Order 3

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, testing-boss, coding-guidelines, evidence-gate.

Implemented the Task slice for Daemon Verification:

- Updated the canonical Roundfix `baseline` and `setup` references with the unpinned owned-lock rule, the `skills.lockUnpinned` and removed-`ref` reports, the `DR-SKILL-LOCK-STALE` warning, and its `roundfix skills install --target project` next action.
- Ran `make skills-sync`; the Roundfix reference mirrors match their canonical copies. The owned-skill recorder raised both Roundfix `SKILL.md` version fields from `0.1.182` to `0.1.183` and appended the sanctioned digest record.
- Added `TestTheGuidesStateTheUnpinnedOwnedLockEntry`, with `t.Parallel()` as its first statement, covering both copies of each reference on whitespace-normalized text.
- Added **Stale Owned Lock Entry** beside **Setup Snapshot Difference** in `CONTEXT.md`, linked to ADR-0359.
- Re-recorded `docs/references/behavior-surfaces.json` with the required sanctioned command as the final regeneration step.

Acceptance evidence:

- Both canonical references contain the removed-`ref` wording and their required report or warning text; the two mirrors are byte-identical to their canonical sources. Static checks passed.
- `CONTEXT.md` contains the required `**Stale Owned Lock Entry**:` glossary definition and ADR-0359 link.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` passed; both Roundfix version fields are `0.1.183` and the version record contains its generated digest. The final Behavior Surface recorder also passed.

Focused checks:

- `gofmt -d skills/owned_lock_guidance_test.go` — no formatting changes.
- `git diff --check` — passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` — passed.
- `go test ./skills -run '^$'` — compiled the skills package; no tests were selected.

The declared Verification commands and repository gates remain Daemon-owned and were not run. Task status remains Daemon-owned `in_progress`; no commit, push, PR, Task Graph edit, or other Task file was changed.
