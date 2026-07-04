---
name: business-analyst
description: "Turns exploration into decision-ready analysis: scores feature ideas on six viability criteria with KPI targets, converts open questions into A/B/C options with a recommended default, and produces executive-grade deliverables (KPI frameworks, forecasts, cohort/LTV-CAC models, A/B readouts, market sizing). Use during idea/PRD discovery to score or structure a decision, and whenever the task is quantitative decision support. Don't use for software implementation, UI work, or database migration."
metadata:
  category: discovery
  tags: [product, research, requirements]
  version: 0.2.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Business Analyst

Convert exploration into decisions. Three modes, picked by what the caller needs; `write-idea` uses the first two, executive deliverables use the third.

## Mode 1 — Feature scoring (for idea/PRD discovery)

Score a feature idea so its priority is an argument, not a feeling:

1. **Six-criteria assessment** — each scored `Must do / Strong / Maybe / Pass`, each score justified in one line citing research or codebase findings:

   | Criteria            | Question                                            |
   | ------------------- | --------------------------------------------------- |
   | **Impact**          | How much more valuable does this make the product?  |
   | **Reach**           | What % of users would this affect?                  |
   | **Frequency**       | How often would users encounter this value?         |
   | **Differentiation** | Does this set us apart or just match competitors?   |
   | **Defensibility**   | Is this easy to copy or does it compound over time? |
   | **Feasibility**     | Can we actually build this?                         |

2. **Leverage type** — Quick Win (small effort, disproportionate value), Strategic Bet (larger effort, potentially transformative), or Compounding (gets more valuable over time: data effects, habit formation, network effects).
3. **KPIs** — 3 to 6, each with a numeric target ("> 30%", "< 200ms", "-80%") and a concrete, implementable measurement method. A KPI nobody can measure is a wish.
4. **Viability verdict** — one paragraph grounded in the research: proceed, reshape, or pass, and what would change the answer.

## Mode 2 — Decision support (structure an open question)

When a discovery conversation hits an open decision, convert it into a decision-ready block instead of prose:

```text
Which retention lever should V1 optimize for?

A) Weekly digest email  ← suggested: cheapest to ship, measurable in one cycle
B) In-app streak mechanics
C) Usage-based notifications
D) Other — describe

Assumptions to confirm: users check email weekly; digest infra exists.
```

Rules: 2–4 options per decision; the suggestion always carries a one-line rationale; assumptions that would invalidate the suggestion are listed, not hidden. One decision per block — batching decisions produces rubber-stamping.

## Mode 3 — Executive deliverables

For full quantitative work (KPI frameworks, dashboards, forecasts, cohort/LTV-CAC models, A/B readouts, market sizing), follow the Output Contract:

**Operating rules.** Every deliverable names a decision, not a topic. Every number has a source, formula, or assumption-table row. State the worst case explicitly (CI lower bound, downside sensitivity, stressed LTV). Flag data-quality gaps before conclusions — never hedge with "TBD". If the contract cannot be met, say so in the opening paragraph.

**Structure.** Open with `## Business Objective & Success Criteria` (the decision this enables + 3–5 numeric success criteria) and close with `## Recommendations` (3–7 numbered items, each with **Decision / Owner** (a named role, never "the team") **/ By when / Expected impact**). Between them, use executive headings — never `## Context`/`## Analysis`/`## Summary`:

- **Action titles** stating the insight in ≤15 words: `## Enterprise cohort NRR fell 8pts; ship usage-based pricing in Q3`.
- **SCR** for decision memos: `## Situation` → `## Complication` → `## Resolution` (the resolution is 60–70% of the content).
- Standard vocabulary where the type calls for it: `## Input Metrics` / `## Output Metrics` (KPI frameworks), `## Decision Required` (before Recommendations on approval memos), `## So What` (page-level takeaway), `## Assumptions & Sensitivities`, `## Reconciliation` (market sizing: where top-down meets bottom-up).

## Anti-patterns

- Scoring without justification lines — a bare "Strong" transfers no information.
- Presenting one option with a question mark instead of genuine alternatives.
- KPI targets without measurement methods, or "improve engagement" as a KPI.
- Generic persona prose ("as an expert analyst I…") — this skill is its outputs, not its résumé.
