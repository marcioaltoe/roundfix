---
type: refactor
status: open
created: 2026-09-30
spec: null
reason: null
---

# Baseline wording left after the stack wave

## Opportunity

Spec 0206 adopted the Backlog Entry "Stack modules need rules with force", and Spec 0207 covers the composed profile and the frontend layout. Five items from that entry are covered by neither:

- the backend and frontend guides do not render the declared workspace paths;
- the production-code and debugging skill triggers belong to the TypeScript module but are worded as language-neutral, so no module clearly owns them;
- the generic-layers prohibition has no scope, and it reads as forbidding domain services;
- the backend guide names a fixed auth owner;
- the TypeScript setup requires skills that no dispatch uses (`zustand`, `ai-sdk`, `motion`, `tanstack-table`).

## Value

Each item puts a sentence in an adopter's guide that is either untrue for that repository or owned by the wrong module.

## Shape

One Spec that settles the owner of the two triggers, scopes the generic-layers clause, renders the declared workspace paths, makes the auth owner a recorded decision or removes it, and drops the unused setup skills. Each change states its Source Baseline retention disposition. Evidence: `docs/specs/0206-stack-rules-with-force-and-the-rules-adopters-repeat/references/2026-09-30-stack-modules-need-rules-with-force.md`.
