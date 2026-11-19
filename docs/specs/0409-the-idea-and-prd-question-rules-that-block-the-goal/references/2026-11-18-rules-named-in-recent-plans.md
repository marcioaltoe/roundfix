# Rules named in the plans of Specs 0364 to 0406

Measured by the author of the Idea on 2026-11-18, in this repository at `5d3237bd5f6db3fde9cfeaf74de65b23d8f42f65`, from the Git history alone. It answers one question: how often does a plan draw its scope at a recorded decision, and how often does it ask the maintainer about that decision?

## Method

1. List every `_prd.md` and `_idea.md` that a commit added under `docs/specs/` for Specs 0364 to 0406 (`git log --diff-filter=A --name-only`). When a Spec was renamed, keep the latest path. Numbers 0394 and 0395 have no plan in the history.
2. Read each file as it was in the commit that added it.
3. In a PRD, take the bullets of `## Non-Goals / Out of Scope`. In an Idea, take the bullets of `## Out of Scope (V1)`. Count the bullets that name an ADR (`ADR-NNNN`) or a guide under `docs/agents/`.
4. In a PRD, take the bullets of `## Decisions`. Count those that name an ADR and say it is amended or superseded.
5. Read every counted line, and say whether it records a question put to the maintainer about the decision it names.

## Result

The 41 Specs were planned between 2026-11-06 and 2026-11-18. They hold 41 PRDs and 6 Ideas.

| measure | count |
| --- | --- |
| Scope exclusions in the 41 PRDs | 267 |
| Of those, exclusions that name an ADR or a guide | 19, in 15 PRDs |
| Scope exclusions in the 6 Ideas | 44 |
| Of those, exclusions that name an ADR or a guide | 2, in 2 Ideas |
| Decisions in the 41 PRDs | 400 |
| Of those, decisions that amend or supersede an ADR | 13, in 12 PRDs |
| Of the 13, the author's own statement ("ADR-NNNN records the decision and amends ADR-MMMM") | 11 |
| Of the 13, a dated answer of the maintainer to another question, whose chosen option included the amendment | 2 |
| Lines, of the 34 counted, that record a question about whether a recorded decision is kept, changed or retired | 0 |

## What the count shows and what it does not

- A plan in this repository changes an ADR often: 12 of 41 PRDs amend or supersede one. The author decides the amendment and records it. No plan records that it asked the maintainer about a decision that stood in the way.
- A plan draws its scope at a recorded decision in 15 of 41 PRDs. The count cannot tell a boundary the author chose from a goal the author gave up: an exclusion such as "Changing the ADR horizon of" a named ADR reads as a boundary, and nothing in the file says whether a better solution was dropped.
- The count covers this repository only. The maintainer reported the problem from the projects that adopt the Baseline, and no count of it exists there.
