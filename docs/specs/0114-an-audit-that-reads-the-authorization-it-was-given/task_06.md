---
task: task_06
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: qa
complexity: high
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. It walks every user story through its surface, carries
the outside-evidence row from task_05, and closes the Spec's vocabulary.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the task evidence,
   and execute each row through a declared surface rather than by reading
   artifacts.
2. MUST carry task_05's outside-evidence row: the two refusals replayed from
   commits this Spec did not author — the Spec 0095 Task split to satisfy the
   audit, and the 2026-08-13 sanctioned-fallout refusal.
3. MUST rewrite this Spec's own Tooling authority rows to the wording ADR-0131
   settles, in both the PRD and the TechSpec, once the detector accepts it — the
   Spec was refused three times for writing what its own template teaches, and
   leaving the workaround in place would ship the defect as documentation.
4. MUST classify every finding by user impact and record auditable evidence per
   row.
5. MUST check whether this Spec introduced, changed, or retired a term the
   glossary should carry, and update the domain context when it found something.
   Governed Path is declared in the Vocabulary Contract and MUST reach the
   glossary; the gate is refused by its own static precondition otherwise.
6. MUST record a row as blocked with its reason rather than passing it on
   equivalent artifacts when its surface cannot be reached.
7. MUST write the dated QA report into this Spec's evidence directory.
8. MUST NOT commit scratch repositories, built binaries, or any file it did not
   author as evidence.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Execute every row through its declared surface.
- [ ] Replay the two historical refusals for the outside-evidence row.
- [ ] Restore this Spec's own tooling rows to the settled wording.
- [ ] Run the glossary check and update the domain context.
- [ ] Write the dated QA report with per-row evidence.

## Acceptance Criteria

- [ ] Every PRD user story has at least one executed row with recorded evidence.
- [ ] The outside-evidence row names both replayed commits and what the audit
      says about each now, or is blocked with its reason.
- [ ] This Spec's PRD and TechSpec record Tooling authority as applicable with no
      mutation proposed, and `spec check --strict` accepts both.
- [ ] Governed Path is in the glossary with the semantics the checker emits.
- [ ] The QA report carries a verdict, per-row classification, and evidence paths.
- [ ] No scratch repository, built binary, or unauthored file is committed as
      evidence.

## Verification

- `ls docs/specs/0114-an-audit-that-reads-the-authorization-it-was-given/qa/qa-report-*.md > /dev/null 2>&1 || { echo 'no QA report written'; exit 1; }` — expected: exits 0, proving the gate wrote its report.
- `r=$(ls docs/specs/0114-an-audit-that-reads-the-authorization-it-was-given/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to read'; exit 1; }; grep -q '^verdict:' "$r" || { echo "no verdict in $r"; exit 1; }` — expected: exits 0, proving the newest report carries a verdict.
- `for f in docs/specs/0114-an-audit-that-reads-the-authorization-it-was-given/_prd.md docs/specs/0114-an-audit-that-reads-the-authorization-it-was-given/_techspec.md; do grep -q 'Tooling authority: applicable' "$f" || { echo "FAIL: $f still carries the workaround wording"; exit 1; }; done; grep -q 'Governed Path' CONTEXT.md || { echo 'Governed Path never reached the glossary'; exit 1; }` — expected: exits 0, proving the Spec stopped documenting the defect it removed and that its coined term landed. Fails today on both clauses.
- `r=$(ls docs/specs/0114-an-audit-that-reads-the-authorization-it-was-given/qa/qa-report-*.md 2>/dev/null | tail -1); test -n "$r" || { echo 'no QA report to audit'; exit 1; }; git ls-files -s docs/specs/0114-an-audit-that-reads-the-authorization-it-was-given/qa | grep '^160000' && { echo 'a scratch repository was committed as a gitlink'; exit 1; }; exit 0` — expected: exits 0, proving no scratch repository was committed as a gitlink, anchored to an existing report.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## References

`_prd.md` → every User Story, Success Metrics. `_techspec.md` → Coverage Map;
Risks & Considerations. ADR-0091, ADR-0104, ADR-0131.
