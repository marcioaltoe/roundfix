---
name: write-prd
description: Write a PRD as a durable local artifact at docs/specs/<slug>/_prd.md — research first, clarify with one multiple-choice question at a time, record decisions as ADRs, then write directly without draft-approval rounds.
disable-model-invocation: true
argument-hint: "<feature description, or nothing after a grilling/brainstorm session>"
metadata:
  category: planning
  tags: [prd, product, requirements, workflow, documentation]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Write PRD

Produce `docs/specs/<feature-slug>/_prd.md` — the product definition the rest of the pipeline (`write-techspec`, `write-tasks`, `implement-task`, `qa-gate`) implements from. The PRD owns _what_ and _why_; the tech spec owns _how_. Keeping that boundary is what lets each artifact stay small and stay true while the code changes underneath it.

## Inputs

`$ARGUMENTS` is a feature description, or empty when the current conversation already explored the feature (a grilling or brainstorm session). Everything already answered in the conversation counts as answered — do not re-interview.

If `docs/specs/<slug>/_idea.md` exists (produced by `write-idea`), it **is** the exploration: read it first, and treat its research, scoring, council insights, chosen direction, and Out of Scope list as answered ground truth. Clarify only what the idea left open.

## Size triage first

Not every change earns a PRD. If the work is a single vertical slice with no open product decisions (a bug fix, a small extension of existing behavior), say so and recommend skipping to a single task file or direct implementation. A PRD pays for itself when there are product decisions to record and multiple tasks to derive.

## Ground rules

- **Read `CONTEXT.md` and `docs/adr/` before anything else.** Use the glossary's vocabulary throughout — titles, user stories, feature names — and never drift to synonyms the glossary lists under `_Avoid_`. If a concept has no term yet, sharpen one with the user and add it to the glossary right then; a PRD written in fuzzy language produces fuzzy tasks. If either file is missing, proceed silently.
- **Decide, then write.** Once research and clarifications are done, write the file directly and let the user react to the finished artifact. Do not present outline drafts for approval — reviewing a real PRD is faster than reviewing a promise of one.
- **Durability.** No file paths, no code snippets, no line numbers. The PRD may sit in a queue for weeks while the codebase moves; describe behavior and interfaces, which survive refactors.

## Process

### 1. Research

Explore before asking anything — questions the codebase can answer are wasted user time:

- Existing behavior and adjacent features the change touches.
- Prior specs under `docs/specs/` and `docs/specs/_archived/` — overlap with something already built or planned is a finding to surface, not to silently absorb.
- Market/competitor context via web research when the feature is user-facing and positioning matters.

### 2. Clarify

Ask **one question per message**, multiple-choice whenever the options are enumerable:

```text
Which failure behavior should an expired import have?

A) Retry automatically up to 3 times  ← suggested: matches the sync retry ADR
B) Fail fast and notify the user
C) Park it for manual review
D) Other — describe
```

Always state a suggested default and the one-line reason. Cover, in order of importance: goals and success criteria, functional scope, non-goals, constraints, risks. Stop asking when the remaining unknowns don't change what gets built.

### 3. Record decisions

A product decision that is hard to reverse, surprising without context, and the result of a real trade-off becomes an ADR at `docs/adr/NNNN-slug.md`, continuing the repository's numbering. Keep it to 1–3 sentences: context, decision, why. Decisions that fail that three-part gate just live in the PRD body.

### 4. Write

**HARD RULE — spec folders are numbered `docs/specs/NNNN-<kebab-slug>/`** (zero-padded 4 digits, e.g. `0001-implement-command`). Determine `NNNN` by scanning **both** `docs/specs/` and `docs/specs/_archived/` for the highest existing prefix and adding 1; use `0001` when no specs exist anywhere. Numbers are never reused and travel with the spec when archived. Never create an unnumbered spec folder. When an `_idea.md` fed this PRD, its folder already carries the number — reuse it, don't mint a new one.

Create `docs/specs/NNNN-<kebab-slug>/` and write `_prd.md` from the template in [references/prd-template.md](references/prd-template.md). If an `_idea.md` fed this PRD, flip its frontmatter `status` to `promoted`. Set the PRD frontmatter carefully — downstream skills parse it:

- `spec` — the folder slug.
- `status: active` — flipped to `archived` by `archive-spec` after release.
- `surfaces` — every surface the feature touches (`frontend`, `backend`, `cli`, `data`, `infra`, `docs`). `qa-gate` routes browser-based QA from this list, so an omitted `frontend` means the feature ships without browser validation.

### 5. Report

Reply with the file path, any open questions that survived clarification, and the next step: `write-techspec` for features with architectural decisions to make, `write-tasks` directly when the technical approach is already obvious.

## Anti-patterns

- Technical design in the PRD — schemas, endpoints, package layout belong in `_techspec.md`.
- Re-asking what the conversation or the codebase already answered.
- Several questions in one message, or open-ended questions where options were enumerable.
- Inventing requirements the user never confirmed to make the document look complete.
- Padding user stories with variations nobody asked for — every story must trace to a confirmed need.

## References

- [references/prd-template.md](references/prd-template.md) — the full `_prd.md` template. Read it before writing the file.
