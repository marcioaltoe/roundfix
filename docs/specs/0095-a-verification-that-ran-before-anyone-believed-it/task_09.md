---
task: task_09
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: qa
complexity: high
---

# Task 09: Run the final QA gate

## Overview

The authored terminal gate. It walks every user story through the surfaces the
Spec declares, rests one row on evidence this Spec did not author, and carries the
glossary check that closes the Spec's vocabulary.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the task evidence,
   and execute each row through a declared surface rather than by reading
   artifacts.
2. MUST rest one named row on evidence originating outside this Spec's artifacts:
   the authored graphs of Specs this Spec did not build, checked with
   `--run-verification` to see whether it finds what a clean checker missed.
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
- [ ] Run the outside-evidence row against Specs this Spec did not author.
- [ ] Run the glossary check and update the domain context if it found something.
- [ ] Write the dated QA report with per-row evidence.

## Acceptance Criteria

- [ ] Every PRD user story has at least one executed row with recorded evidence.
- [ ] The outside-evidence row names the Specs it checked and what
      `--run-verification` reported for each. Source: the authored graphs of the
      Spec set in this repository, none of which this Spec wrote. Blocked with
      that reason recorded if none can be checked.
- [ ] The two new refusals and the restored one each refuse a measured form and
      spare its correct counterpart, proven through the command rather than
      through unit tests alone.
- [ ] The glossary check ran, and its outcome is recorded either way.
- [ ] The QA report carries a verdict, per-row classification, and evidence paths.
- [ ] No scratch repository, built binary, or unauthored file is committed as
      evidence.

## Verification

- `ls docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/qa/qa-report-*.md > /dev/null 2>&1 || { echo 'no QA report written'; exit 1; }` — expected: exits 0, proving the gate wrote its report.
- `grep -q '^verdict:' $(ls -t docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/qa/qa-report-*.md | head -1)` — expected: exits 0, proving the newest report carries a verdict.
- `r=$(ls -t docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/qa/qa-report-*.md 2>/dev/null | head -1); test -n "$r" || { echo 'no QA report to audit'; exit 1; }; git ls-files -s docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/qa | grep '^160000' && { echo 'a scratch repository was committed as a gitlink'; exit 1; }; exit 0` — expected: exits 0, proving no scratch repository was committed as a gitlink. The audit is anchored to an existing report so an absent evidence directory fails rather than passes.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## References

`_prd.md` → every User Story, Success Metrics. `_techspec.md` → Coverage Map;
Risks & Considerations. ADR-0091, ADR-0104, ADR-0124.
