# PRD template

Write `docs/specs/<feature-slug>/_prd.md` with this exact structure. Guidance appears as `<!-- comments -->`; delete every comment from the final file. Omit a section only when it is genuinely empty — write "None." rather than deleting Non-Goals or Open Questions, so readers know it was considered.

```markdown
---
spec: <feature-slug>
status: active
created: YYYY-MM-DD
surfaces: [backend] # every surface touched: frontend | backend | cli | data | infra | docs
---

# <Feature name, in glossary vocabulary>

<!-- One paragraph: the problem from the user's perspective and the outcome this feature buys. -->

## Goals

<!-- 2-5 bullets. Each goal is an observable outcome, not an activity. -->

## User Stories

<!-- Numbered list, `As a <actor>, I want <capability>, so that <benefit>`.
     Every story traces to a confirmed need; qa-gate walks each one against the running app. -->

1. As a ..., I want ..., so that ...

## Core Features

<!-- Numbered functional requirements. Behavior only — no schemas, endpoints, or file paths.
     Each requirement must be verifiable: a reader can tell "done" from "not done". -->

1. ...

## User Experience

<!-- Only for user-facing surfaces: key flows, states (empty/loading/error), breakpoints if relevant.
     Describe what the user sees and does, not which components render it. -->

## Non-Goals / Out of Scope

<!-- Explicit exclusions. This section protects the tasks from scope creep — be generous here. -->

## Success Metrics

<!-- How we know the feature worked after shipping. Measurable where possible. -->

## Decisions

<!-- Product decisions made during clarification, one line each.
     Decisions that passed the ADR gate link to their ADR: `See ADR-0012.` -->

## Open Questions

<!-- Anything that survived clarification, with the default that applies until answered. -->
```
