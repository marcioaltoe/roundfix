---
task: task_04
spec: 0140-a-spec-traces-the-promises-it-makes
status: pending
type: docs
complexity: medium
---

# Task 04: State the declaration rule where artifacts are written

## Overview

The checker now enforces the rule; this slice states it where authors read. The
PRD and TechSpec skills gain the declaration form and its `None.` variant, the
Task skill gains the mapping and References obligation, and each template shows
the shape it prescribes. This Task is an authorized tooling mutation and may
change only the bounded files below plus its own Task file.

## Requirements

1. MUST state, in the PRD skill and its template, that Success Metrics are
   numbered items or a single `None.` with the reason none applies.
2. MUST state, in the TechSpec skill and its template, the same form for API
   Contracts.
3. MUST state, in the Task skill, that every declared promise is mapped and that
   each appears in some Task's References, beside the existing user story and
   Core Feature obligations.
4. MUST show, in the Task template's References example, a Success Metric and an
   API Contract named the way the checker reads them.
5. MUST regenerate the distributed mirror so both copies of every bounded file
   stay identical, and MUST rewrite a derived pin only through the sanctioned
   regeneration command.
6. MUST NOT change any path outside the bounded list, and MUST NOT change a
   skill's version, the Task Type contract, the QA decision rules, or any
   Verification rule.

## Subtasks

- [ ] Write the declaration form into the PRD skill and its template.
- [ ] Write the same form into the TechSpec skill and its template.
- [ ] Write the mapping and References obligation into the Task skill and show it
      in the Task template's example.
- [ ] Regenerate the mirror and confirm both copies are identical.

## Acceptance Criteria

- [ ] Each canonical skill and template carries its new clause, and the mirror
      copy of each is byte-identical.
- [ ] The Task template's References example names a Success Metric and an API
      Contract.
- [ ] No skill version changes, so no setup minimum moves.
- [ ] The changed-path set is inside the bounded list of the Spec's approved
      authority record.

## Context

- instruction: `.agents/skills/write-prd/SKILL.md`
- instruction: `.agents/skills/write-prd/references/prd-template.md`
- instruction: `.agents/skills/write-techspec/SKILL.md`
- instruction: `.agents/skills/write-techspec/references/techspec-template.md`
- instruction: `.agents/skills/write-tasks/SKILL.md`
- instruction: `.agents/skills/write-tasks/references/task-template.md`
- instruction: `skills/write-prd/SKILL.md`
- instruction: `skills/write-prd/references/prd-template.md`
- instruction: `skills/write-techspec/SKILL.md`
- instruction: `skills/write-techspec/references/techspec-template.md`
- instruction: `skills/write-tasks/SKILL.md`
- instruction: `skills/write-tasks/references/task-template.md`

## Verification

- `grep -q "or the single entry" .agents/skills/write-prd/references/prd-template.md` — expected: exit 0; the PRD template prescribes the numbered form and its explicit-none alternative.
- `grep -q "or the single entry" .agents/skills/write-techspec/references/techspec-template.md` — expected: exit 0; the TechSpec template prescribes the same form.
- `grep -q "numbered Success Metric" .agents/skills/write-prd/SKILL.md` — expected: exit 0; the PRD skill states the rule.
- `grep -q "numbered API Contract" .agents/skills/write-techspec/SKILL.md` — expected: exit 0; the TechSpec skill states the rule.
- `grep -q "Success Metric" .agents/skills/write-tasks/SKILL.md` — expected: exit 0; the Task skill carries the mapping and References obligation.
- `grep -q "Success Metric" .agents/skills/write-tasks/references/task-template.md` — expected: exit 0; the References example names a metric.
- `grep -q "API Contract" .agents/skills/write-tasks/references/task-template.md` — expected: exit 0; the References example names a contract.
- `grep -q "or the single entry" skills/write-prd/references/prd-template.md` — expected: exit 0; the mirror carries the same clause.
- `grep -q "or the single entry" skills/write-techspec/references/techspec-template.md` — expected: exit 0.
- `grep -q "numbered Success Metric" skills/write-prd/SKILL.md` — expected: exit 0.
- `grep -q "numbered API Contract" skills/write-techspec/SKILL.md` — expected: exit 0.
- `grep -q "Success Metric" skills/write-tasks/SKILL.md` — expected: exit 0.
- `grep -q "Success Metric" skills/write-tasks/references/task-template.md` — expected: exit 0.
- `grep -q "Success Metric" skills/write-tasks/SKILL.md && diff -q .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md` — expected: exit 0; the mirror carries the new clause and stays identical to its canonical copy. Before this Task the clause is absent, so the command fails.

## References

`_prd.md` → Core Features 1, 2; User Stories 1-3; Success Metric 2;
Project Constraints: Tooling authority;
`_techspec.md` → System Architecture: Authoring rules; Testing Approach 5;
Build Order 4; `_authorization.md`; ADR-0156.
