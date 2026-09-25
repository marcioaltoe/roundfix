---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A QA Report whose front matter is empty settled the QA Task as passed

## Symptom

Spec 0170's second QA Report, `qa/qa-report-2026-09-25-01.md` at commit `e3042b40`, opened with two consecutive `---` lines, so its front matter was empty and `verdict: pass` sat outside it. The derived QA Verification's reader refuses that file (exit 1), yet the Run settled the QA Task `completed` and committed the report as `(pass)`. The pre-PR review caught it.

## Where

`internal/daemon/task_engine.go` QA settlement and report commit; `internal/spec/qa.go` `readQAReport`; whichever step wrote the extra delimiter (the QA Agent editing the seeded report, or the Daemon rewriting it).

## Expected

Settlement reads the report exactly as the derived Verification does and refuses a report whose front matter is empty or duplicated.

## Evidence

Run `run_20260925T195747Z_0ccff24444c45b55`; second pre-PR review of Spec 0170, 2026-09-25.
