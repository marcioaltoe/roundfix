---
title: "Resolve: deduplicate Review Issues and assign Batches"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 22
  - 23
  - 24
  - 25
  - 26
  - 27
blocked_by:
  - 06-resolve-compatible-artifacts.md
---

# Resolve: deduplicate Review Issues and assign Batches

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Before invoking any Agent, make `resolve` deduplicate repeated unresolved Review
Issues across selected Compatible Artifact Rounds and assign only the newest
occurrences into bounded Batches. Older duplicate occurrences must stay local
artifact bookkeeping until the newest occurrence reaches a terminal outcome.

## Acceptance criteria

- [x] Review Issue Fingerprints prefer source thread identity and fall back to a
      provider-specific review hash.
- [x] Newest duplicate selection orders by Round, source review submission time,
      and Round creation time.
- [x] Ambiguous newest duplicate selection fails during Preflight Validation and
      does not create a Run.
- [x] Only newest unresolved occurrences are assigned to Agent Batches.
- [x] Older duplicate occurrences are associated with the newest occurrence and
      marked `duplicated` only after the newest occurrence reaches `resolved` or
      `invalid`.
- [x] Tests prove duplicate groups produce one Agent assignment and older
      duplicate occurrences do not resolve source threads separately.

## Blocked by

- 06-resolve-compatible-artifacts.md
