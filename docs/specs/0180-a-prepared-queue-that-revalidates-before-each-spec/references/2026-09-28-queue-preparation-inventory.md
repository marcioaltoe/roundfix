---
type: feat
status: promoted
created: 2026-09-28
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
reason: null
---

# Queue preparation inventories what is approved to run

## Opportunity

Spec 0127 Core Feature 1: `deliver start` takes slugs; nothing inventories active Specs, Backlog intent, Findings and the Inbox, records dependencies and pending decisions, and names which Specs are actually approved.

## Value

The operator (or Supervisor) sees one prepared queue with its prerequisites and open decisions before anything runs; an inventory is never implementation authority.

## Shape

A read-only `deliver plan` (or `deliver start --dry-run`) over the lifecycle roots, reusing `spec check` and authorization resolution.

Evidence: carried from Spec 0127 (`docs/history/specs/0127-durable-unattended-spec-workflow/_prd.md`) when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
