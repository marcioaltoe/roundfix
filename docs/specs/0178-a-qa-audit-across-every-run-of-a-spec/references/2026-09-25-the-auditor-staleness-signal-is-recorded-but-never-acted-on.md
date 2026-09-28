---
status: done
created_at: 2026-09-25
updated_at: 2026-09-28
absorbed_by: 0178-a-qa-audit-across-every-run-of-a-spec
---

# Triage — The auditor staleness signal is recorded but never acted on (2026-09-25)

Triaged from secondbrain `inbox/roundfix/2026-09-19-qa-gate-passes-on-self-declared-stale-binary.md` on 2026-09-25 and revalidated against main 7a9b6ec6.
QA Reports record `auditor_staleness`, but settlement never reads it; in Roundfix self-audit it reads `stale` by construction, so refusing it outright would block every gate. The QA Agent runs public-CLI rows with whatever `roundfix` is on PATH.

## 1. The auditor staleness signal is recorded but never acted on

- Symptom / evidence: `internal/spec/qa.go` parses `AuditorStaleness` but `QAReportEligibility` and `settleQAVerdict` never read it; `internal/spec/auditor_evidence.go` returns stale for any build commit that is a strict ancestor of HEAD.
- Root cause: The Homebrew shadow named in the original capture is gone and the derived QA command no longer calls a binary (Spec 0152); the QA Agent still runs public-CLI rows with whatever `roundfix` is on PATH.
- Action / suggestion: Suggestion: have the gate build `./cmd/roundfix` from the audited revision for public rows, or redefine staleness; needs a design decision.

## Addendum — 2026-09-28 — Adopted by Spec 0178

[0178-a-qa-audit-across-every-run-of-a-spec](../_prd.md) adopts this Finding.
Fresh evidence from 2026-09-25 to 2026-09-28: the QA Reports of Specs 0171,
0172 and 0174 (first report) record `auditing_binary` as the Daemon's
`bin/roundfix` built from main (`256ad156`, `f0faa780`) and
`auditor_staleness: stale`, while the QA Agent's user-flow rows ran a
`bin/roundfix` built from the audited tree. The reports of Specs 0170, 0173,
0174 (second report) and 0176 instead record `auditing_binary` as a build of
the audited head, and three of them carry a staleness reason no Roundfix writer
produces, so the field named two different binaries. The Codex pre-PR review of
Spec 0172 flagged the stale field as a defect.

The Spec makes the design decision this Finding asked for. `auditing_binary`
names the Daemon's binary, which runs the mechanical stage and settles the
gate, and `auditor_staleness` compares its build commit with the Delivery Base
(the merge base of the audited head and the repository default branch) instead
of the audited head, so commits the candidate adds never make it stale. A
`stale` auditor is published as a `daemon.qa` warning and recorded in the
seeded report, and the gate proceeds, because a multi-item `roundfix deliver`
queue keeps one binary while later items start from a newer main. The QA Agent records
the binary its public rows ran as `user_flow_binary`, and in a Roundfix
self-audit settlement requires that binary to be built from the audited head.
The source moves once into that Spec's reference index; its lifecycle status
records adoption, not implementation or QA completion.
