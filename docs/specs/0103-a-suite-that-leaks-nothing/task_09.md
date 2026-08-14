---
task: task_09
spec: 0103-a-suite-that-leaks-nothing
status: failed # pending | in_progress | completed | failed — only implement-task changes this
type: qa
complexity: high
---

# Task 09: Run the final QA gate

## Overview

The authored terminal gate. It walks every user story through the surfaces the
Spec declares, carries the outside-evidence row this Spec rests on, and closes
the Spec's vocabulary.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the task evidence,
   and execute each row through a declared surface rather than by reading
   artifacts.
2. MUST carry the outside-evidence row from task_01: the upstream documentation
   of the write-then-execute hazard (golang/go#22315), and the 2026-08-06
   process-table measurement of four survivors consuming two hours and forty
   minutes of CPU — neither produced by this Spec.
3. MUST classify every finding by user impact and record auditable evidence per
   row.
4. MUST check whether this Spec introduced, changed, or retired a term the
   glossary should carry, and update the domain context when it found something.
5. MUST record a row as blocked with its reason rather than passing it on
   equivalent artifacts when its surface cannot be reached.
6. MUST write the dated QA report into this Spec's evidence directory.
7. MUST NOT commit scratch repositories, built binaries, or any file it did not
   author as evidence — which is also this Spec's own Core Feature 7, so the gate
   is its first subject.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Execute every row through its declared surface.
- [ ] Carry the outside-evidence row and name both sources.
- [ ] Run the glossary check and update the domain context if it found something.
- [ ] Write the dated QA report with per-row evidence.

## Acceptance Criteria

- [ ] Every PRD user story has at least one executed row with recorded evidence.
- [ ] The outside-evidence row names the upstream issue and the 2026-08-06
      measurement, and neither is an artifact this Spec authored.
- [ ] The compiled fixtures, the bounded wait, the guard, the teardown, the
      residue inventory, the tree proof, and the evidence refusal are each proven
      through their own surface rather than through unit tests alone.
- [ ] The glossary check ran, and its outcome is recorded either way.
- [ ] The QA report carries a verdict, per-row classification, and evidence paths.
- [ ] No scratch repository, built binary, or unauthored file is committed as
      evidence.

## Verification

- `ls docs/specs/0103-a-suite-that-leaks-nothing/qa/qa-report-*.md > /dev/null 2>&1 || { echo 'no QA report written'; exit 1; }` — expected: exits 0, proving the gate wrote its report.
- `r=$(ls docs/specs/0103-a-suite-that-leaks-nothing/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to read'; exit 1; }; grep -q '^verdict:' "$r" || { echo "no verdict in $r"; exit 1; }; grep -q '22315' "$r" || { echo "the outside-evidence row does not name its upstream source in $r"; exit 1; }` — expected: exits 0, proving the report carries a verdict and the outside evidence.
- `r=$(ls docs/specs/0103-a-suite-that-leaks-nothing/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to audit'; exit 1; }; git ls-files -s docs/specs/0103-a-suite-that-leaks-nothing/qa | grep '^160000' && { echo 'a scratch repository was committed as a gitlink'; exit 1; }; exit 0` — expected: exits 0, proving no scratch repository was committed as a gitlink. Anchored to an existing report so an absent evidence directory fails rather than passes.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## References

`_prd.md` → every User Story, Success Metrics. `_techspec.md` → Coverage Map;
Risks & Considerations. ADR-0091, ADR-0104, ADR-0125, ADR-0126, ADR-0127.
