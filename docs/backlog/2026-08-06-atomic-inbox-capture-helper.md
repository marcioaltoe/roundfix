---
type: feat # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-06
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Add an atomic helper for durable inbox capture

## Opportunity

Fleet sessions could capture a contract-valid Inbox Entry through one
observable action instead of manually creating, staging, reviewing, committing,
and timing the entry. This intent came through
`inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md` in the
Secondbrain.

## Value

An atomic capture action would reduce the risk of missing the one-minute
durability target and make the path, brain commit, and elapsed time available
to the capturing session without reconstructing them from shell history.

## Shape

A future design could validate the Inbox Entry contract, create only the
destination entry, commit only that brain-relative path, and report the path,
commit hash, and capture-to-durability seconds. This shape is non-binding; the
owning Spec must choose whether it belongs in a skill, helper, or product
surface.
