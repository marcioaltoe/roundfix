---
name: write-techspec
description: Write the technical spec for an approved PRD at docs/specs/<slug>/_techspec.md — explore the architecture, settle technical decisions one question at a time, record them as ADRs, and produce a build order that write-tasks can decompose.
disable-model-invocation: true
argument-hint: "<spec slug or path to docs/specs/<slug>/_prd.md>"
metadata:
  category: engineering-design
  tags: [architecture, documentation, workflow]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Write TechSpec

Produce `docs/specs/<slug>/_techspec.md` — the technical answer to an approved `_prd.md`. The PRD said _what_ and _why_; this document decides _how_, _where_, and _with which_ — and its Build Order is what `write-tasks` turns into the task graph, so sequencing quality here becomes execution quality later.

## Preconditions

`$ARGUMENTS` names the spec (slug or path). `docs/specs/<slug>/_prd.md` must exist — if it doesn't, stop and point the user at `write-prd`. If the PRD contains low-level technical decisions that belong here, surface that as a finding and propose relocating them rather than silently duplicating.

## Ground rules

- **Read first**: `CONTEXT.md` (use the glossary vocabulary in every interface and component name), existing `docs/adr/` (respect prior decisions or challenge them explicitly — never contradict one silently), and the full PRD.
- **Decide, then write.** Settle decisions during clarification, record the significant ones as ADRs, then write the finished document. No outline-approval rounds.
- **YAGNI.** Do not propose new packages, layers, or directories when the feature fits as an addition to an existing module. Every new structural element must earn its place against a stated alternative.
- **Reference, don't duplicate.** Point at PRD sections by name ("covers User Stories 3–5"); repeating business context here creates two copies that drift.

## Process

### 1. Explore the architecture

Map the code the feature will live in: the modules it extends, the seams it can attach to, existing patterns for the same kind of work (an existing adapter, an existing job runner), and the test infrastructure available. Prefer existing seams to new ones — the ideal number of new seams is zero. Delegate the exploration to a subagent when one is available; the spec needs conclusions, not file dumps.

### 2. Clarify technical decisions

Same protocol as the PRD stage — **one question per message, multiple-choice with a suggested default and a one-line rationale** — but scoped to how/where/which: architecture placement, data model, API shape, migration strategy, testing seams, performance and security constraints. Only ask what exploration couldn't settle.

### 3. Record ADRs

Each significant technical decision (hard to reverse + surprising without context + real trade-off) becomes `docs/adr/NNNN-slug.md`, **continuing the same numbering the PRD stage used** — one decision log spans the whole feature. 1–3 sentences each.

### 4. Write

Write `_techspec.md` from the template in [references/techspec-template.md](references/techspec-template.md). Two sections carry the most downstream weight:

- **Every PRD goal and user story maps to a named technical component.** An unmapped story is a design hole; find it now, not during task execution.
- **Build Order** — numbered steps where every step after the first states which previous steps it depends on. `write-tasks` derives the task graph edges from this.

Keep interface sketches under 20 lines each; they document shape, not implementation. Target 1,500–2,500 words — a spec nobody reads protects nobody.

### 5. Report

Reply with the file path, the ADRs created, any decisions still open, and the next step: `write-tasks`.

## Anti-patterns

- Restating PRD business context instead of referencing it.
- Speculative generality — abstractions for requirements the PRD doesn't contain.
- Interface sketches that grow into implementations.
- A Build Order without dependency statements — that forces `write-tasks` to guess the graph.
- Contradicting an existing ADR without naming it and proposing to supersede it.

## References

- [references/techspec-template.md](references/techspec-template.md) — the full `_techspec.md` template. Read it before writing the file.
