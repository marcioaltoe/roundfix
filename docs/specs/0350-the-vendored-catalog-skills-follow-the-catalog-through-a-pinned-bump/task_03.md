---
task: task_03
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
status: pending
type: docs
complexity: low
---

# Task 03: The glossary and the repository rules name the Catalog Skill Bump

## Overview

Write the two terms this Spec declares into `CONTEXT.md` and add the
repository rule that the Catalog Skill Bump is the only writer of the
Vendored Skills, their lock entries and the upstream Setup Snapshots, so an
Agent that wants to correct a vendored skill is sent upstream
(`_prd.md` → Core Feature 4).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md`](references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md),
   its Expected section, for "nobody edits a vendored copy here".
2. MUST write, through the `domain-modeling` skill, two Glossary Definitions
   in `CONTEXT.md`, placed right after the **Vendored Skill** entry, each
   with an `_Avoid_:` line and a link to ADR-0360:
   **Catalog Revision** (the one full commit of `marcioaltoe/skills` that
   every catalog entry of `skills-lock.json` and every upstream Setup
   Snapshot name, never a branch name) and **Catalog Skill Bump** (`make
   skills-third-party-bump`, which moves the Vendored Skills, their lock
   entries and every upstream Setup Snapshot to one Catalog Revision, by
   default the catalog's current `main`; refuses a Roundfix-owned skill and a
   tree the catalog does not publish; a second run at the same commit changes
   nothing). MUST NOT change any other line of `CONTEXT.md`.
3. MUST add to `docs/agents/specific-repository.md`, right after the
   skill-ownership hard rule, a hard rule stating that the Vendored Skills
   under `.agents/skills`, their `skills-lock.json` entries and every upstream
   Setup Snapshot change only through `make skills-third-party-bump`, which
   pins them to one Catalog Revision, and that a correction to a Vendored
   Skill goes to `marcioaltoe/skills` and arrives with the next bump
   (ADR-0360). The sentence MUST contain the words `change only through`
   followed by `make skills-third-party-bump`, which may be set in code
   formatting.
4. MUST NOT change any other guide, skill, ADR or user guide page.

## Subtasks

- [ ] Write the two glossary entries.
- [ ] Add the repository rule.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Catalog Revision** and **Catalog Skill Bump** next to **Vendored Skill**.
- [ ] `docs/agents/specific-repository.md` names the bump as the only writer and sends corrections upstream.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/agents/domain.md`
- instruction: `docs/adr/0360-the-vendored-catalog-skills-follow-the-catalog-through-one-pinned-commit.md`
- interface: `CONTEXT.md`
- interface: `docs/agents/specific-repository.md`

## Consumers

None. The Task adds glossary entries and a repository rule; no code, test or record reads either text.

## Verification

- `grep -F -- '**Catalog Revision**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Catalog Revision\*\*(:| —( |$))'` — expected: exit 0; before this Task `CONTEXT.md` defines no **Catalog Revision**.
- `grep -F -- '**Catalog Skill Bump**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Catalog Skill Bump\*\*(:| —( |$))'` — expected: exit 0; before this Task `CONTEXT.md` defines no **Catalog Skill Bump**.
- `tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qE -- 'change only through .?make skills-third-party-bump' || { printf 'missing phrase in %s: %s\n' docs/agents/specific-repository.md 'change only through make skills-third-party-bump' >&2; exit 1; }` — expected: exit 0; before this Task the repository rules do not name the bump.

## References

- `_prd.md` → Core Feature 4; Glossary
- `_techspec.md` → Coverage Map; Glossary; Build Order 3
- ADR-0360
