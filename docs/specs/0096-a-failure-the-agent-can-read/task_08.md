---
task: task_08
spec: 0096-a-failure-the-agent-can-read
status: failed # pending | in_progress | completed | failed — only implement-task changes this
type: qa
complexity: high
---

# Task 08: Run the final QA gate

## Overview

The authored terminal gate. It walks every user story through its surface, carries
the outside-evidence row, and closes the Spec's vocabulary.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the task evidence,
   and execute each row through a declared surface rather than by reading
   artifacts.
2. MUST rest one named row on evidence originating outside this Spec's artifacts:
   the repeated failure measured on 2026-08-08 in a repository this Spec did not
   build, recorded in `docs/findings/`. Record the row as blocked with its reason
   if that record cannot be read.
3. MUST classify every finding by user impact and record auditable evidence per
   row.
4. MUST check whether this Spec introduced, changed, or retired a term the
   glossary should carry, and update the domain context when it found something.
   Repeated Failure is declared in the Vocabulary Contract and MUST reach the
   glossary.
5. MUST record a row as blocked with its reason rather than passing it on
   equivalent artifacts when its surface cannot be reached.
6. MUST write the dated QA report into this Spec's evidence directory.
7. MUST NOT commit scratch repositories, built binaries, or any file it did not
   author as evidence.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Execute every row through its declared surface.
- [ ] Read the 2026-08-08 record for the outside-evidence row.
- [ ] Run the glossary check and update the domain context.
- [ ] Write the dated QA report with per-row evidence.

## Acceptance Criteria

- [ ] Every PRD user story has at least one executed row with recorded evidence.
- [ ] The outside-evidence row names the record it read and what the repetition
      check does with that case now, or is blocked with its reason.
- [ ] The absent diagnostic, the repetition, the vacuity event, the budget
      statement, the surface naming, and the ceiling exits are each proven through
      their own surface rather than through unit tests alone.
- [ ] Repeated Failure is in the glossary with the semantics the Daemon emits.
- [ ] The QA report carries a verdict, per-row classification, and evidence paths.
- [ ] No scratch repository, built binary, or unauthored file is committed as
      evidence.

## Verification

- `ls docs/specs/0096-a-failure-the-agent-can-read/qa/qa-report-*.md > /dev/null 2>&1 || { echo 'no QA report written'; exit 1; }` — expected: exits 0, proving the gate wrote its report.
- `r=$(ls docs/specs/0096-a-failure-the-agent-can-read/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to read'; exit 1; }; grep -q '^verdict:' "$r" || { echo "no verdict in $r"; exit 1; }; grep -qE '^\| ' "$r" || { echo "the Results table in $r has no rows"; exit 1; }; grep -q 'Repeated Failure' CONTEXT.md || { echo 'Repeated Failure never reached the glossary'; exit 1; }` — expected: exits 0, proving the report carries a verdict and rows, and that the coined term landed. Fails today on the last clause.
- `r=$(ls docs/specs/0096-a-failure-the-agent-can-read/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to audit'; exit 1; }; git ls-files -s docs/specs/0096-a-failure-the-agent-can-read/qa | grep '^160000' && { echo 'a scratch repository was committed as a gitlink'; exit 1; }; exit 0` — expected: exits 0, proving no scratch repository was committed as a gitlink, anchored to an existing report.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## References

`_prd.md` → every User Story, Success Metrics. `_techspec.md` → Coverage Map;
Risks & Considerations. ADR-0091, ADR-0104, ADR-0135, ADR-0136, ADR-0137.
