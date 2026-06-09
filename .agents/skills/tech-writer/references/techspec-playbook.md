# Tech spec playbook

Tech specs, design docs, RFCs, and ADRs. A spec answers: **given that we're building this, what exactly does it do and how?** The reader is writing code, tests, or rollout plans — the spec makes the behavior unambiguous.

## Rules for every spec

- **One proposal per document.** Lead with the problem; the design earns its place by solving it.
- **Don't repeat the PRD.** Reference its sections by name ("see PRD → Goals"). The spec adds the how; restating the why creates drift.
- **Precision matches the source.** If the source names an endpoint, name it; if it only says "a queue", do not pick a vendor. `TBD` markers beat invented details.
- **Trade-offs are mandatory.** A spec with no downsides is not finished. Name what the design gives up, the alternatives rejected and why, and the remaining risks.
- **Contracts beat prose.** Interfaces, schemas, and diagrams over paragraphs. Keep code examples under 20 lines each.
- **Numbers for every non-functional claim.** Latency, throughput, error budget, retention — measurable targets only, and only the ones the source states or the user confirms.
- **Decide rollout and rollback.** Behind which flag, in what order, and how to back out. A migration without a rollback path is a one-way door — say so explicitly.

## Full tech spec skeleton

For implementation-ready specifications. Common convention: lives beside the PRD as `_techspec.md`, with ADRs in a sibling `adrs/` folder.

```markdown
# <feature> — Technical Specification

## Executive Summary

<1–2 paragraphs: key architectural decisions, implementation strategy, primary trade-offs>

## System Architecture

### Component Overview

<each component: name, purpose, boundaries, data flow, external interactions>

## Implementation Design

### Core Interfaces

<key service interfaces as code blocks, ≤20 lines each, in the project's language>

### Data Models

<domain entities, relationships, field types, request/response types, schemas>

### API Endpoints

<by resource: method, path, purpose, request format, response format, status codes>

## Integration Points

<only when the design touches outside systems: service, purpose, auth, error handling, retry strategy>

## Impact Analysis

| Component | Impact Type | Description & Risk | Required Action |
| --------- | ----------- | ------------------ | --------------- |

## Testing Approach

### Unit Tests

<strategy, key components, mock requirements, critical scenarios and edge cases>

### Integration Tests

<what gets tested together, test data, environment dependencies>

## Development Sequencing

### Build Order

1. <step>
2. <step — depends on step 1>

### Technical Dependencies

<blocking dependencies: infrastructure, external services, other teams' deliverables>

## Monitoring and Observability

<key metrics, log events with structured fields, alerting thresholds>

## Technical Considerations

### Key Decisions

<for each: the decision, the rationale, what it trades away, alternatives rejected and why>

### Known Risks

<risk, likelihood, mitigation, open research areas>

## Architecture Decision Records

- [ADR-001: <title>](adrs/adr-001.md) — <one-line summary>
```

Gates: every step in Build Order after the first names the steps it depends on; Core Interfaces shows at least one real interface or type definition; the spec ends with at least one ADR documenting the primary approach; every PRD goal maps to a component.

## RFC / design-doc skeleton

For proposals seeking a decision rather than full implementation detail. Copy literally — plain section text, no bold-prefix labels.

```markdown
# <proposal as a short noun phrase>

Status: <Draft | In review | Approved>
Decider: <named person or group> · Decision needed by: <date, or omit the line>

## Problem

<what hurts today, with numbers and cited signals>

## Goals and non-goals

Goals: <from the source>
Non-goals: <explicit exclusions>

## Proposed design

<the how, at exactly the precision the source provides — contracts, schemas, and diagrams beat prose>

## Alternatives considered

<name> — rejected: <the stated reason, restated plainly. Stop there.>

## Risks and trade-offs

<honest negatives the source states or that follow directly from the design>

## Open questions

- <question> — owner: <name>
```

## ADR

One architecturally significant decision per record. Numbered file in source control (`adrs/adr-001.md` or `docs/adr/0007-use-webhooks.md` — follow the project's convention). Immutable after acceptance: supersede with a new ADR and link both, never edit. One to two pages; if it needs 4,000 words it is a design doc, not an ADR.

```markdown
# ADR-NNN: <decision as a short noun phrase>

## Status

<Proposed | Accepted | Deprecated | Superseded by ADR-XXX>

## Date

<YYYY-MM-DD>

## Context

<the forces in tension, as neutral facts with their numbers and sources>

## Decision

We will <the decision, in full sentences, active voice>.

## Alternatives Considered

### <alternative name>

- Description: <what it is>
- Pros: <stated upsides>
- Cons: <stated downsides>
- Why rejected: <the stated reason, plainly. Stop there.>

## Consequences

- Positive: <stated outcomes>
- Negative: <stated trade-offs — at least one; never soften or omit them>
- Risks: <each with its mitigation, when the source gives one>

## References

<links to related docs, specs, external resources — omit if none>
```

State each rejection reason and consequence only as given. "Stateful connections complicate the load balancer setup" must not grow into session affinity, connection pooling, or peak-season scenarios the source never mentioned.

## Quality gates

1. Problem stated before design; no PRD content restated.
2. Every interface, schema, and endpoint at source precision; `TBD` where the source is silent.
3. Alternatives and trade-offs present; at least one honest negative.
4. Build order is numbered with explicit dependencies.
5. Testing, monitoring, and rollout addressed (or explicitly marked out of the decision's scope).
6. At least one ADR records the primary approach.
