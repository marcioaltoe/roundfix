---
task: task_07
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: failed # pending | in_progress | completed | failed — only implement-task changes this
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate, and this Spec's first subject: every repair here
changes how this very node behaves, so the gate proves its own contract by
executing under it.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the task evidence,
   and execute each row through a declared surface rather than by reading
   artifacts.
2. MUST rest one named row on evidence originating outside this Spec's artifacts:
   the three refusals recorded on 2026-08-14 and 2026-08-15 in
   `docs/backlog/2026-08-14-a-spec-that-coins-a-term-cannot-pass-its-own-gate.md`
   and the reports under `docs/history/specs/0103-a-suite-that-leaks-nothing/qa/`,
   none of which this Spec authored. Record the row as blocked with its reason if
   those artifacts cannot be read.
3. MUST classify every finding by user impact and record auditable evidence per
   row.
4. MUST check whether this Spec introduced, changed, or retired a term the
   glossary should carry, and update the domain context when it found something.
5. MUST record a row as blocked with its reason rather than passing it on
   equivalent artifacts when its surface cannot be reached.
6. MUST write the dated QA report into this Spec's evidence directory.
7. MUST NOT commit scratch repositories, built binaries, or any file it did not
   author as evidence.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Execute every row through its declared surface.
- [ ] Read the three recorded refusals for the outside-evidence row.
- [ ] Run the glossary check and update the domain context if it found something.
- [ ] Write the dated QA report with per-row evidence.

## Acceptance Criteria

- [ ] Every PRD user story has at least one executed row with recorded evidence.
- [ ] The outside-evidence row names the recorded refusals it read and what the
      repaired node does with each now, or is blocked with its reason.
- [ ] The four repairs — the terminal refusal row, the literal-naming diagnostic,
      the single finding for one cause, and the performed repair — are each proven
      through the node's own surface rather than through unit tests alone.
- [ ] The glossary check ran, and its outcome is recorded either way.
- [ ] The QA report carries a verdict, per-row classification, and evidence paths.
- [ ] No scratch repository, built binary, or unauthored file is committed as
      evidence.

## Verification

- `ls docs/specs/0113-a-gate-report-that-does-not-block-its-successor/qa/qa-report-*.md > /dev/null 2>&1 || { echo 'no QA report written'; exit 1; }` — expected: exits 0, proving the gate wrote its report.
- `r=$(ls docs/specs/0113-a-gate-report-that-does-not-block-its-successor/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to read'; exit 1; }; grep -q '^verdict:' "$r" || { echo "no verdict in $r"; exit 1; }; grep -qE '^\| ' "$r" || { echo "the Results table in $r has no rows — the defect this Spec removes"; exit 1; }` — expected: exits 0, proving the report carries a verdict and at least one row. The row clause is this Spec's own first defect, asserted against its own gate.
- `r=$(ls docs/specs/0113-a-gate-report-that-does-not-block-its-successor/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to audit'; exit 1; }; git ls-files -s docs/specs/0113-a-gate-report-that-does-not-block-its-successor/qa | grep '^160000' && { echo 'a scratch repository was committed as a gitlink'; exit 1; }; exit 0` — expected: exits 0, proving no scratch repository was committed as a gitlink, anchored to an existing report.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## References

`_prd.md` → every User Story, Success Metrics. `_techspec.md` → Coverage Map;
Risks & Considerations. ADR-0091, ADR-0104, ADR-0132, ADR-0133, ADR-0134.
