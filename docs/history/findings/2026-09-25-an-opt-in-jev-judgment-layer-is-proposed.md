---
status: deferred
created_at: 2026-09-25
updated_at: 2026-09-29
closure_reason: "Not implemented: probabilistic judgments would sit where ADR-0091, ADR-0096, ADR-0140 and ADR-0148 require mechanical proof, the nearest independent evidence is negative, and the maintainer closed the Jev front on 2026-09-29."
closure_evidence: docs/history/backlog/2026-09-25-jev-gateway-baseline-before-routing.md
---

# Triage — An opt-in Jev judgment layer is proposed for gates (2026-09-25)

Triaged from secondbrain `inbox/roundfix/2026-09-18-jev-typesafe-toggle-roundfixrc.md` on 2026-09-25 and revalidated against main 7a9b6ec6.
A proposal to let a per-project `.roundfixrc.yml` toggle send gate state (QA verdict, evidence sufficiency, archive decision, agent selection) to Jev. It would put probabilistic judgments where ADRs require mechanical proof, and the vendor numbers are unmeasured on Roundfix cases.

## 1. An opt-in Jev judgment layer is proposed for gates

- Symptom / evidence: No Go source mentions TypeSafe or Jev; Jev has only been used offline to rank the queue.
- Root cause: The proposed gates would put probabilistic judgments where ADR-0091, ADR-0140 and ADR-0148 require mechanical proof, and would send gate state to a remote service.
- Action / suggestion: Suggestion: benchmark Jev on real gate cases before any Backlog Entry; no implementation intent yet.

## Addendum — 2026-09-29 — Deferred: the Jev front is closed

The maintainer closed the Jev front on 2026-09-29 after a read-only analysis of the Secondbrain evidence and of Roundfix's own gateway logs. No gate judgment will be implemented.

- Conflict: the proposed gates would replace the hermetic mechanical stage (ADR-0096), the prober that refuses vacuous Verification (ADR-0148) and the proved selection tuple (ADR-0091, ADR-0140) with probabilistic answers.
- Evidence: the closest independent measurement is negative (trajectory fault attribution at AUROC 0.56 and difficulty routing at 51%, secondbrain `wiki/sources/jev-ecosistema-verificadores-routers-2026-09-25.md`). Jev itself ranked this item last in both queue rankings (0.31, secondbrain `raw/roundfix/2026-09-25-jev_triage_priority.json`).
- Exposure: gate state such as diffs and QA reports would leave the machine.
- Related: the gateway measurement closes in `docs/history/backlog/2026-09-25-jev-gateway-baseline-before-routing.md`, and the 2026-09-28 jev-guard inbox entry was discarded with the same analysis.

Reopen only when a gate appears where mechanical proof is impossible, and at least 50 labelled historical outcomes show Jev improving on the current verdict.

