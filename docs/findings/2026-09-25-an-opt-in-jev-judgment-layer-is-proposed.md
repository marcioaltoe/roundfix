---
status: pending
created_at: 2026-09-25
updated_at: 2026-09-25
---

# Triage — An opt-in Jev judgment layer is proposed for gates (2026-09-25)

Triaged from secondbrain `inbox/roundfix/2026-09-18-jev-typesafe-toggle-roundfixrc.md` on 2026-09-25 and revalidated against main 7a9b6ec6.
A proposal to let a per-project `.roundfixrc.yml` toggle send gate state (QA verdict, evidence sufficiency, archive decision, agent selection) to Jev. It would put probabilistic judgments where ADRs require mechanical proof, and the vendor numbers are unmeasured on Roundfix cases.

## 1. An opt-in Jev judgment layer is proposed for gates

- Symptom / evidence: No Go source mentions TypeSafe or Jev; Jev has only been used offline to rank the queue.
- Root cause: The proposed gates would put probabilistic judgments where ADR-0091, ADR-0140 and ADR-0148 require mechanical proof, and would send gate state to a remote service.
- Action / suggestion: Suggestion: benchmark Jev on real gate cases before any Backlog Entry; no implementation intent yet.
