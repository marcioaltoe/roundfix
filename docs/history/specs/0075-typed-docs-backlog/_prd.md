---
spec: 0075-typed-docs-backlog
status: archived
created: 2026-08-03
surfaces: [backend, docs]
archived: "2026-08-05"
source_slug: 0075-typed-docs-backlog
---


# A typed backlog for triage intent

## Problem

The documentation layout gives observations a home and gives intent none.
`docs/findings/` holds evidence-backed field reports — what happened, dated
and immutable — and the maintainer confirmed that contract is right. But a
suggestion with no observed evidence has nowhere honest to live: it is not a
finding (nothing was observed), and it is not yet a Spec (nobody committed
to it). In practice such notes either get forced into findings, diluting
what "finding" means, or lost.

The maintainer's direction (2026-08-03): add a backlog directory to the
CONTEXT-driven layout. Deliberately named `docs/backlog/` rather than
`docs/ideas/`, because intent comes in more than one shape — a feature idea
today, a fix tomorrow, other kinds later — and a typed backlog can absorb
them all, opening a path to eventually deprecating `docs/findings/` into it.
A backlog `feat` entry is explicitly **not** the `write-idea` artifact: it
is upstream raw material — almost a feature backlog — that the spec pipeline
may later consume.

This is a Baseline product change, not a repository-local one: the layout
guide is a managed artifact every adopting repository receives, so the
backlog must be defined where the layout is defined.

## Project Constraints

- Identifier strategy: applicable — project-owned vocabulary is created:
  the Backlog Entry and its type values (`feat`, `fix`, `perf`,
  `refactor`), deliberately aligned with the Conventional Commits types this
  repository already enforces on commit messages and PR titles, added to
  the `CONTEXT.md` glossary and used verbatim in the layout contract.
  `refactor` is the canonical token, never an abbreviation. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — documentation layout and
  Baseline asset content only. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0081 sanctions the
  deterministic digest fallout of the authorized asset edit; ADR-0092
  (created by this Spec) records the backlog decision. The knowledge-flow
  principle stands: `CONTEXT.md` and guides never reference backlog
  content, which is downstream and deletable. ADR-0080 owns QA verdict
  semantics and ADR-0091 owns the authored QA gate as a typed Task node,
  under which this Spec's own graph is authored. ADR-0093 surfaces as a
  relation candidate because it cites ADR-0080; it does not apply — it
  governs the Spec Consistency Check's detection boundary. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no Makefile, workflow, skill, or pin
  changes; the surface is Baseline module assets, the generated guide, and
  the glossary. No owned skill is edited: the layout guide owns the routing
  contract, and the authoring skills already consume triage material
  generically. Source: `docs/agents/agent-instructions.md`.

## Goals

- `docs/backlog/` exists in the CONTEXT-driven layout with one job:
  **intent** — what to do next, typed, not yet committed to a Spec.
- `docs/findings/` keeps one job: **observation** — what happened, with
  evidence, immutable. The guide states the boundary so the two never blur
  again.
- Every backlog entry carries a typed frontmatter contract and a per-type
  body template, copyable from the guide exactly as the findings contract
  is today.
- Adopting repositories receive the new layout the same way they receive
  the rest of it: through the managed guide and its module clauses.

## Core Features

1. The one-job-per-directory clause gains `docs/backlog/` for typed intent
   entries; the guide explains what belongs there and what does not.
2. A Backlog Operational Contract in the guide: dated filenames
   (`YYYY-MM-DD-slug.md`), frontmatter with `type`
   (`feat` | `fix` | `perf` | `refactor`), `status`
   (`open` | `promoted` | `declined`), `created`, plus `spec` when
   promoted and `reason` when declined. The type vocabulary is the
   Conventional Commits vocabulary, so one word carries the intent from
   backlog entry to Spec to commit to PR title.
3. One body template per type: `feat` (opportunity, value, shape of a
   solution), `fix` (symptom, where, expected behavior, evidence pointer —
   a finding link when one exists), `perf` (what is slow, the measurement
   that says so, the target), and `refactor` (what is tangled, what it
   costs, the desired shape).
4. The finding-versus-backlog boundary stated in both directions: a
   finding may spawn a backlog entry; a backlog entry needs no finding; an
   entry is never evidence and a finding is never a commitment.
5. Promotion mirrors adoption: when a Spec adopts a backlog entry, the
   entry records `promoted` with the Spec link and moves to that Spec's
   `references/`, leaving the backlog; git history is the trail. An empty
   backlog is not the goal — an *honest* one is.
6. The `CONTEXT.md` glossary defines Backlog Entry and its types, and
   distinguishes a `feat` entry from the `write-idea` pipeline artifact.
7. The type set is deliberately open, and extension follows the same
   rule that seeded it: a new type must be a Conventional Commits type
   that expresses intent. The eventual deprecation of `docs/findings/`
   into a backlog type stays a future decision this structure can absorb
   without reshaping.

## Non-Goals / Out of Scope

- Deprecating `docs/findings/` now — recorded as direction, decided later.
- Changing `write-idea` or any owned skill; the guide owns the routing.
- Migrating existing findings or inbox notes into the backlog.
- Any Roundfix binary behavior change: no command reads the backlog.

## Success Metrics

- A fresh Baseline apply produces a layout guide that names
  `docs/backlog/`, carries the typed contract with both templates, and
  states the finding boundary.
- This repository's own `docs/agents/docs-layout.md` reflects the new
  contract after its baseline re-applies, and `docs/backlog/` accepts its
  first entry using the template unchanged.
- The compatibility corpus and digest chain regenerate deterministically
  per ADR-0081, with the breaks declared, not discovered.
- `CONTEXT.md` defines the vocabulary, and no guide references any backlog
  entry's content.

## Decisions

- `backlog`, not `ideas`: intent has more than one shape, and the folder
  name must not have to change when the second shape arrives. (ADR-0092)
- Findings stay observations and stay put; the deprecation path into the
  backlog is explicitly future work.
- The type vocabulary is the Conventional Commits vocabulary — one word
  from intent to PR title. Chosen by the maintainer (2026-08-03) from his
  own usage: findings kept absorbing problem reports that were really
  `fix` intent.
- A `feat` entry is upstream of `write-idea`, never its output.
- Four types now — `feat`, `fix`, `perf`, `refactor`; the contract is open
  to more.

## Open Questions

None.
