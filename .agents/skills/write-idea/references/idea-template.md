# Idea template

Write `docs/specs/<slug>/_idea.md` with this structure. Guidance appears as `<!-- comments -->`; delete every comment from the final file. Insufficient information goes to Open Questions, never guessed into a section.

```markdown
---
spec: <NNNN-feature-slug> # numbered folder slug, e.g. 0001-implement-command
created: YYYY-MM-DD
status: proposed # proposed | promoted (write-prd ran) | parked (consciously shelved)
leverage: quick-win # quick-win | strategic-bet | compounding
---

# <Idea name, in glossary vocabulary>

## Overview

<!-- What problem it solves, who it's for, why it's valuable, how ambitious V1 should be. One tight paragraph each. -->

## Problem

<!-- 2-4 paragraphs with concrete scenarios; why current solutions/workarounds fall short.
     Add a "### Market Data" subsection when research produced numbers. -->

## Core Features

<!-- 3-10 rows, ordered by priority. Behavior in 1-2 lines each — no implementation. -->

| #   | Feature | Priority                 | Description |
| --- | ------- | ------------------------ | ----------- |
| F1  | <name>  | Critical / High / Medium | <behavior>  |

## KPIs

<!-- 3-6 rows. Targets numeric ("> 30%", "< 200ms", "-80%"); measurement method concrete and implementable. -->

| KPI | Target | How to measure |
| --- | ------ | -------------- |

## Feature Assessment

<!-- The business-analyst scoring. Each score: Must do / Strong / Maybe / Pass. -->

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  |       |
| Reach           | What % of users would this affect?                  |       |
| Frequency       | How often would users encounter this value?         |       |
| Differentiation | Does this set us apart or just match competitors?   |       |
| Defensibility   | Is this easy to copy or does it compound over time? |       |
| Feasibility     | Can we actually build this?                         |       |

## Council Insights

<!-- From the embedded council session. Keep the dissent — it feeds the PRD's risks. -->

- **Recommended approach:** ...
- **Key trade-offs:** ...
- **Risks identified:** ...
- **Dissenting view:** ...
- **Stretch goal (V2+):** ...

## Opportunity Scan

<!-- The alternatives considered (ambitious / simpler / adjacent), each with a one-line what/why/score,
     and the chosen direction with rationale. This is the record of the shapes NOT taken. -->

**Chosen direction:** original / alternative N / hybrid — <rationale>

## Out of Scope (V1)

<!-- Minimum 3 exclusions, each with a justification. This section is the PRD's scope firewall. -->

- **<excluded>** — <why it's out of V1>

## Decisions

<!-- Decisions made during this session, one line each; ADR-gated ones link: `See ADR-0014.` -->

## Open Questions

<!-- What still needs an answer, and the default that applies until it gets one. -->
```

## Optional sections

Add between Core Features and Out of Scope when content justifies: **Differentiator** (the competitive angle, when there is one), **Integration with Existing Features** (table: integration point → how), **Sub-Features** (when the idea should split into multiple specs), **Cost Estimate** (paid APIs, storage — volume and monthly cost).
