---
type: fix
status: promoted
created: 2026-09-25
spec: 0172-a-qa-gate-that-tells-the-truth
reason: null
---

# A QA Report the Agent never filled settles as pass

## Symptom

The Daemon writes every QA Report starting as `verdict: pass` with an empty Results table. When the QA Agent stops without editing it, settlement, `archive`, `settle` and `qa-report accept` all accept that pass: a gate that measured nothing completes the Spec. Fiscus archived a Spec on a 669-byte, 30-second report.

## Where

`internal/daemon/task_engine.go` (`writeMechanicalQAReport` starts `verdict := spec.VerdictPass`; `settleQAVerdict` reads only `QAReportEligibility`), `internal/spec/qa.go` (pass branch checks only blocked-row counts). The empty-Results check in `internal/speccheck/mechanical.go` runs only before the Agent.

## Expected

The starting report carries a non-accepting verdict (e.g. `pending`), and settlement refuses a pass whose Results table has no rows.

## Evidence

secondbrain `inbox/roundfix/2026-09-18-um-relatorio-de-qa-sem-matriz-passa-na-verificacao-canonica.md`.
