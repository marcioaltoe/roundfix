---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A new QA Report fills a sequence gap and loses to an older report

## Symptom

The QA Report allocator takes the first free `-NN` while the selector picks the highest, so after a gap the fresh report is ignored and settlement reads a stale verdict.

## Where

`internal/daemon/task_engine.go` `writeMechanicalQAReport`; `.agents/skills/qa-gate/SKILL.md` wording.

## Expected

Allocate one above the highest existing sequence; fix the skill sentence.

## Evidence

secondbrain `inbox/roundfix/2026-09-11-o-qa-novo-reutiliza-um-numero-anterior-ao-relatorio-vigente.md`.
