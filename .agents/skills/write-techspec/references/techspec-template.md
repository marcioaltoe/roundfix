# TechSpec template

Write `docs/specs/<slug>/_techspec.md` with this exact structure. Guidance appears as `<!-- comments -->`; delete every comment from the final file.

```markdown
---
spec: <feature-slug>
prd: _prd.md
created: YYYY-MM-DD
---

# <Feature name> — Technical Spec

## Executive Summary

<!-- 3-6 sentences. Must name the primary technical trade-off this design accepts and why. -->

## System Architecture

<!-- Which existing modules the feature extends, which components are new, and how they connect.
     Name components in glossary vocabulary. A small diagram (mermaid) is welcome when flow is non-obvious. -->

## Implementation Design

### Interfaces

<!-- The key contracts, each sketched in ≤20 lines of code. Shape, not implementation. -->

### Data Models

<!-- New/changed entities and their relationships. Schema changes described by shape and constraint. -->

### API Contracts

<!-- New/changed endpoints or commands: input, output, failure modes. Omit for features with none. -->

## Coverage Map

<!-- One line per PRD goal and user story → the component(s) that satisfy it.
     Format: `Story 3 → ImportScheduler, ImportStatusView`. An unmapped story is a design hole. -->

## Integration Points

<!-- External systems touched, and the boundary pattern used for each (adapter, webhook, queue). -->

## Testing Approach

<!-- The seams where tests attach — prefer existing seams; state each new seam and why it's needed.
     What is covered by unit vs integration tests. implement-task treats these seams as pre-agreed. -->

## Build Order

<!-- Numbered steps. Every step after the first states its dependencies explicitly:
     `3. Import status endpoint (depends on: 1, 2)`.
     write-tasks derives the task graph from this section — vague sequencing here becomes a wrong graph there. -->

1. ...
2. ... (depends on: 1)

## Risks & Considerations

<!-- Known risks with mitigations; performance/security/observability notes worth a task's attention. -->

## Decisions

<!-- Technical decisions made during clarification, one line each; ADR-gated ones link: `See ADR-0013.` -->
```
