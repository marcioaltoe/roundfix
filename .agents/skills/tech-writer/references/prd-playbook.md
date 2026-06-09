# PRD playbook

PRDs, idea/discovery briefs, user stories, and acceptance criteria. A PRD answers one question: **should we build this, and what exactly counts as done?** The reader is deciding whether the problem is real, whether the shape is right, and what success means — not how to code it.

## Rules for every PRD

- **What and why only.** No databases, frameworks, endpoints, or architecture. If a technical constraint genuinely binds the product (compliance, an integration that must exist, a performance floor), state it as a constraint, not a design.
- **Problem before solution.** Lead with the pain and a cited signal — ticket count, support volume, churn number, named report. "Users want X" without a source is a banned claim.
- **Non-goals are mandatory.** An explicit out-of-scope list, with the source's exclusions verbatim. In co-author mode, propose at least three exclusions and confirm them.
- **One owner, dated.** Name who owns the doc and when it was last updated.
- **Ruthless priority.** If everything is P0, nothing is. Group features by priority and keep the MVP honest.
- **Right-size the document.** A feature touching one team gets the lightweight shape; a product or multi-phase feature gets the full shape. If an engineer won't finish reading it, it is too long for its decision.
- **Link, don't restate.** Context that lives elsewhere (research, designs, prior docs) gets a link.

## Full PRD skeleton

For products and multi-phase features. Common convention: lives at the feature's planning folder as `_prd.md`, with ADRs in a sibling `adrs/` folder.

```markdown
# PRD: <feature>

## Overview

<the problem, who has it, the value of solving it, and why now>

## Goals

<business objectives and success metrics, each with a numeric target and timeline where the source gives one>

## User Stories

<grouped by persona>
- As a <user type>, I want <action>, so that <benefit>.

## Core Features

<grouped by priority; for each feature: name, what it does, why it matters, and its functional requirements>

## User Experience

<personas, the primary user flows step by step, UX/accessibility considerations>

## Non-Goals (Out of Scope)

- <excluded feature> — <why, or "deferred to Phase N">

## Phased Rollout Plan

<Phase 1 (MVP), Phase 2, Phase 3 — each with its success criteria>

## Success Metrics

<engagement, performance, business impact, quality — numeric targets only>

## Risks and Mitigations

<adoption, competitive, timeline, dependency risks — product-level, not technical>

## High-Level Technical Constraints

<optional — required integrations, compliance, performance floors, data privacy. Constraints, never designs.>

## Architecture Decision Records

- [ADR-001: <title>](adrs/adr-001.md) — <one-line summary>

## Open Questions

- <question> — owner: <name>
```

## Lightweight PRD skeleton

For a feature scoped to one team, often published directly as a tracker issue. Decisions already made during exploration go in, file paths and code stay out (they go stale) — except a snippet that encodes a decision more precisely than prose can (a state machine, reducer, schema, or type shape).

```markdown
# <feature>

## Problem Statement

<the problem, from the user's perspective, with the cited signal>

## Solution

<the solution, from the user's perspective>

## User Stories

1. As a <actor>, I want <feature>, so that <benefit>.

## Implementation Decisions

<decisions already made: modules touched, interface changes, schema changes, API contracts, architectural choices — as decisions, not designs>

## Testing Decisions

<what a good test looks like (external behavior only), which modules get tested, prior art for similar tests>

## Out of Scope

- <explicit exclusion>

## Further Notes

<anything that fits nowhere else>
```

## Idea / discovery brief

The document before the PRD: a raw idea expanded into something assessable. Tables force precision.

```markdown
# Idea: <name>

## Overview

<problem it solves, target users, value, how ambitious V1 should be>

## Problem

<2–4 paragraphs with concrete scenarios; include a Market Data subsection when research exists>

## Core Features

| #  | Feature | Priority | Description |
| -- | ------- | -------- | ----------- |
| F1 | <name>  | Critical | <1–2 lines of concrete behavior> |

## KPIs

| KPI | Target | How to Measure |
| --- | ------ | -------------- |
| <observable metric> | <numeric value with unit> | <concrete, implementable method> |

## Out of Scope (V1)

- <excluded feature> — <short justification>

## Open Questions

- <question> — owner: <name>
```

Gates: 3–10 features ordered Critical > High > Medium; 3–6 KPIs, every target numeric; at least 3 exclusions, each justified.

## User stories and acceptance criteria

- Story shape: "As a <user>, I want <task>, so that <goal>." If the "so that" is hollow, the story is a task wearing a costume — rewrite or drop it.
- Each criterion is an observable someone who did not build it can check: a behavior, an output, a passing command. "Works correctly" checks nothing.
- Use Given-When-Then when state matters: "Given an expired session, when the user submits the form, then the draft is preserved and the login prompt appears."
- Cover the negative and edge cases the source raises (invalid input, concurrency, empty states). QA must be able to write test cases without asking.
- Replace adjectives with constraints: "fast" → "loads within 0.5s"; "scalable" → "handles 10k concurrent sessions" — only with numbers the source provides or the user confirms.

## Quality gates

1. Problem cites a real signal; no unattributed "users want".
2. Zero implementation details outside the constraints section.
3. Non-goals present and explicit.
4. Every metric numeric; every criterion checkable.
5. Each feature traces to a goal; each phase has success criteria.
6. Transcribe mode: nothing invented — missing sections carry `TBD — needs <owner>`. Co-author mode: every unconfirmed proposal still marked.
