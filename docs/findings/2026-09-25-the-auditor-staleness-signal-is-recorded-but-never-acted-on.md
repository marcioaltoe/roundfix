---
status: pending
created_at: 2026-09-25
updated_at: 2026-09-25
---

# Triage — The auditor staleness signal is recorded but never acted on (2026-09-25)

Triaged from secondbrain `inbox/roundfix/2026-09-19-qa-gate-passes-on-self-declared-stale-binary.md` on 2026-09-25 and revalidated against main 7a9b6ec6.
QA Reports record `auditor_staleness`, but settlement never reads it; in Roundfix self-audit it reads `stale` by construction, so refusing it outright would block every gate. The QA Agent runs public-CLI rows with whatever `roundfix` is on PATH.

## 1. The auditor staleness signal is recorded but never acted on

- Symptom / evidence: `internal/spec/qa.go` parses `AuditorStaleness` but `QAReportEligibility` and `settleQAVerdict` never read it; `internal/spec/auditor_evidence.go` returns stale for any build commit that is a strict ancestor of HEAD.
- Root cause: The Homebrew shadow named in the original capture is gone and the derived QA command no longer calls a binary (Spec 0152); the QA Agent still runs public-CLI rows with whatever `roundfix` is on PATH.
- Action / suggestion: Suggestion: have the gate build `./cmd/roundfix` from the audited revision for public rows, or redefine staleness; needs a design decision.
