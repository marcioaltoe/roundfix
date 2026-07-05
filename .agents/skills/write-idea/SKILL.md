---
name: write-idea
description: Expand a raw product idea into a research-backed spec at docs/specs/<slug>/_idea.md — targeted questions, parallel codebase + market research, business-viability scoring, council debate, and an opportunity scan before drafting. The _idea.md feeds write-prd.
disable-model-invocation: true
argument-hint: "<feature idea or problem description>"
metadata:
  category: discovery
  tags: [product, research, requirements, workflow]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Write Idea

Turn a raw idea into `docs/specs/<slug>/_idea.md` — the Stage 0 artifact that answers "should we build this, and in what shape?" before `write-prd` answers "what exactly are we building?". The idea stage exists because unexamined business assumptions are the cheapest thing to fix here and the most expensive thing to fix after tasks are cut.

## When this stage earns its cost

Use it for product-level ideas: greenfield features, ideas with real market/competitive questions, or anything where the right _shape_ of the solution is still open. For a well-understood feature in a known area, say so and go straight to `grilling`/`grill-with-docs` → `write-prd` — a full idea study of an obvious feature is ceremony, not insight.

## Ground rules

- **Read `CONTEXT.md` first** and keep its vocabulary; new terms this idea coins get sharpened and added to the glossary as they're resolved.
- **One question per message**, multiple-choice with a suggested default whenever options are enumerable (`D) Other — describe` as the escape). 3–6 questions total; stop early if answers show the idea is already well-defined.
- **WHAT, WHY, WHO only.** Databases, APIs, frameworks, and architecture are forbidden topics here — they belong to `write-techspec`.
- **Research before writing.** Never draft an idea unbacked by codebase findings and (for user-facing ideas) market data.
- **Do not write the file until the user approves the draft** — unlike `write-prd`, this stage is inherently interactive: its whole product is shared judgment.

## Process

### 1. Resolve the slug and check for prior art

**HARD RULE — spec folders are numbered `docs/specs/NNNN-<kebab-slug>/`** (zero-padded 4 digits, e.g. `0001-implement-command`). Determine `NNNN` by scanning **both** `docs/specs/` and `docs/specs/_archived/` for the highest existing prefix and adding 1; use `0001` when no specs exist anywhere. Numbers are never reused and travel with the spec when archived. Never create an unnumbered spec folder.

Derive the kebab-case name (2–5 words) and prepend the number. If a spec folder for this idea already exists (any number), read its `_idea.md` and operate in update mode (preserve sections the user hasn't asked to change) — do not mint a new number. Check `docs/specs/` and `docs/specs/_archived/` for overlapping specs — overlap is a finding to surface, not absorb.

### 2. Understand the idea

Walk the question phases in order, one question at a time: problem/pain point → target user and current workarounds → V1 scope (offer: minimal MVP / complete feature / platform — recommend one of the first two) → ambition (quick win / strategic bet / compounding feature, and "what would make this 10x more valuable?") → dependencies on existing features → success criteria. YAGNI ruthlessly during the scope phase.

### 3. Research in parallel

Spawn two parallel explorations: one over the codebase (existing patterns, integration points, adjacent features) and one over the web (3–7 searches varying angle: competitive landscape, market data/statistics, user-expectation patterns, pricing/cost when relevant). If web tools are unavailable, note the limitation and proceed with codebase findings only. Present a merged research summary — codebase findings, competitors and what they do, potential differentiator, relevant data — before any analysis.

### 4. Score viability

Apply the `business-analyst` skill in feature-scoring mode: 3–6 KPIs with numeric targets, and the six-criteria assessment (Impact, Reach, Frequency, Differentiation, Defensibility, Feasibility — each `Must do / Strong / Maybe / Pass`) plus a leverage type (Quick Win / Strategic Bet / Compounding). Present the analysis to the user before proceeding.

### 5. Debate trade-offs

Run the `council` skill in embedded mode on the real dilemmas: V1 scope, priority vs other work, simpler alternatives, risks and hidden dependencies, and the 10x challenge. Extract the recommended approach, key trade-offs, out-of-scope items, and an optional V2+ stretch goal. If the scope decision is hard to reverse, surprising, and a real trade-off, record it as `docs/adr/NNNN-slug.md` (continue the repository numbering).

### 6. Opportunity scan

Before committing to the original shape, run the scan in [references/opportunity-scan.md](references/opportunity-scan.md): assess the original's ceiling, then propose up to three scored alternatives — one more ambitious, one simpler, one adjacent — and recommend a direction. Ask the user to pick: original / alternative N / hybrid / other. Fold the chosen direction into the draft.

### 7. Draft, review, save

Draft from [references/idea-template.md](references/idea-template.md), present it, and iterate (`approved / adjust sections / rewrite section X / discard`). On approval, write `docs/specs/<slug>/_idea.md` with frontmatter `status: proposed`, and point the user at `write-prd` as the next step — it reads `_idea.md` as its exploration input and flips the status to `promoted`. Ideas consciously shelved get `status: parked` and stay as the record that answers "didn't we already consider this?".

## Anti-patterns

- Asking about implementation ("which database?") — the fastest way to poison a product conversation.
- Skipping research because the idea "is obviously good" — that's precisely the idea that needs the market check.
- Padding the draft with features nobody validated — every feature traces to an answer or a finding.
- Smoothing over the council's dissent — preserved disagreement is input for the PRD's risk section.
- Running the full factory on a small, well-understood feature instead of saying "skip to write-prd".

## References

- [references/idea-template.md](references/idea-template.md) — the `_idea.md` template. Read before drafting.
- [references/opportunity-scan.md](references/opportunity-scan.md) — the strategist lens for step 6. Read when the scan starts.
