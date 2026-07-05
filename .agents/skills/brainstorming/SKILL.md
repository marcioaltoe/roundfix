---
name: brainstorming
description: "You MUST use this before any creative work - creating features, building components, adding functionality, or modifying behavior. Explores user intent, requirements and design through one question at a time, then routes the outcome to write-idea, write-prd, write-techspec (refactors/bug fixes), or a direct task."
metadata:
  category: discovery
  tags: [product, requirements, research, workflow]
  version: 0.3.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Brainstorming

Turn a raw intention into a shared understanding through collaborative dialogue — before any spec is written or any code is touched. This skill is the questioning discipline the discovery stage runs on; `write-idea` and `write-prd` both assume a conversation shaped like this has happened.

<HARD-GATE>
Do NOT invoke any implementation skill, write any code, or scaffold anything until the exploration below has run and the user has confirmed the direction. This applies to EVERY request regardless of perceived simplicity.
</HARD-GATE>

## Anti-pattern: "this is too simple to need exploration"

A todo list, a single-function utility, a config change — all of them go through this. "Simple" requests are where unexamined assumptions cause the most wasted work. The exploration can be three questions and two minutes for genuinely simple work, but it must happen.

## The questioning discipline

- **One question per message.** A topic that needs depth becomes a sequence of questions, never a questionnaire.
- **Multiple choice preferred** when the options are enumerable — `A/B/C` plus `D) Other — describe` — always with a suggested default and a one-line reason. Open-ended is fine when options genuinely can't be predetermined.
- **Explore before asking.** Check project state first — files, docs, `CONTEXT.md`, `docs/specs/` (including `_archived/`), recent commits. A question the codebase can answer is wasted user time.
- **WHAT, WHY, WHO — not HOW.** Databases, APIs, frameworks, and architecture belong to the tech-spec stage; raising them here derails the product conversation.
- **YAGNI ruthlessly.** Challenge every feature against the smallest version that delivers the value.

## Question phases

Walk these in order, skipping what the conversation or the codebase already answered. Minimum 3 questions, maximum ~6 — stop early once problem, user, and scope are clear.

1. **Problem / pain point** — what concrete problem does this solve? What prompted it? Offer 2–3 interpretations if ambiguous.
2. **Target user and context** — who is it for, where in their workflow, what's their current workaround?
3. **Scope** — ideal V1 size: minimal MVP / complete feature / platform. Recommend one of the first two.
4. **Ambition** — quick win / strategic bet / compounding feature; and "what would make this 10x more valuable instead of incremental?"
5. **Dependencies** — does it touch or extend existing features? Name candidates from the codebase.
6. **Success criteria** — how will we know it worked?

## Converging

When understanding feels complete:

1. **Propose 2–3 approaches** with trade-offs, leading with your recommendation and why.
2. **Present the direction in sections** scaled to their complexity (a few sentences when straightforward, 200–300 words when nuanced), validating each with the user before moving on.
3. **Capture language as it resolves** — a term the conversation sharpened goes into `CONTEXT.md` right then; a hard-to-reverse decision that survives the three-part gate becomes an ADR.

## Where the outcome goes

The conversation itself is the artifact — route it, don't duplicate it:

- **Product-level idea** (market questions, open solution shape) → `write-idea` (produces `docs/specs/<slug>/_idea.md`).
- **Feature ready to specify** → `write-prd` (synthesizes this conversation; it will not re-interview).
- **Refactor or bug fix** (no product behavior change) → `write-techspec` directly — it mints the spec folder with a minimal `_prd.md`, no product interview.
- **Trivial, well-understood change** (one-line fix, typo, config tweak) → direct implementation, no spec folder; say so explicitly.

When a repo carries `docs/agents/spec-routing.md`, follow its routing table — it is the canonical version of this list for that project.

Do not write a separate design document — the spec pipeline's artifacts (`_idea.md`, `_prd.md`) are where the outcome lives.

## Key principles

One question at a time · multiple choice preferred · explore before asking · YAGNI ruthlessly · always propose alternatives before settling · validate incrementally · be ready to loop back when an answer breaks an earlier assumption.
