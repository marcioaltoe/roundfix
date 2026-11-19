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

## Trade-off Insights

<!-- From the step 6 trade-off debate. Keep the dissent — it feeds the PRD's risks. -->

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

## Idea-stage judgment

<!-- The judgment that started this stage, in this fielded form; write-prd copies the section into the PRD.
     judged-by: `judge (<model>)`, `agent (<why the judge did not answer>)` or `nobody (the user asked for the Idea stage)`.
     Each criterion: yes, no or unsure, with the probability in parentheses when the judge answered; omit the three lines when nobody judged.
     verdict: indicated | not indicated | undecided | not judged.
     stage: ran. shape: changed | split | shelved | unchanged, then a dash and one line that says what. -->

- item: <reference of the judged item, or `the request in the session`>
- judged-by: judge (<model>)
- shape_open: yes (0.85)
- unverified_assumption: unsure (0.41)
- market_question: no (0.09)
- verdict: indicated
- stage: ran
- shape: changed — <what the Idea changed>

## Rule challenges

<!-- One record per Rule Challenge, in this fielded form, or the single entry `None.` when no rule blocked the goal.
     rule: the ADR, clause, definition or document, and its path.
     prevents: up to three things, for this goal.
     asked: the date and the phase: brainstorming, idea, prd or techspec.
     answer: keep | change | retire | waiting, then a dash, the user's words in quotes and the date of the answer.
     effect: what the Spec does with the answer; for a rule the Baseline owns, the message sent and the scope left out.
     new fact: only when the rule already carried a Kept Line; what this question knows that the earlier one did not.
     Leave a blank line between two records. -->

- rule: ADR-NNNN, `docs/adr/NNNN-<slug>.md`
- imposes: <what the rule requires>
- prevents: <what it prevents for this goal>
- asked: YYYY-MM-DD, idea
- answer: change — "<the user's words>", YYYY-MM-DD
- effect: <the requirement that states the change, and the Task that makes it>

## Decisions

<!-- Decisions made during this session, one line each; ADR-gated ones link: `See ADR-0014.` -->

## Open Questions

<!-- What still needs an answer, and the default that applies until it gets one. -->
```

## Optional sections

Add between Core Features and Out of Scope when content justifies: **Differentiator** (the competitive angle, when there is one), **Integration with Existing Features** (table: integration point → how), **Sub-Features** (when the idea should split into multiple specs), **Cost Estimate** (paid APIs, storage — volume and monthly cost).
