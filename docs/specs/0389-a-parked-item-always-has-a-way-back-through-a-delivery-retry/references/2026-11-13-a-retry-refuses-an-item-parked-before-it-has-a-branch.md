---
type: fix
status: promoted
created: 2026-11-13
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
reason: null
---

# A retry refuses an item parked before it has a branch

## Symptom

Reported by fiscus in ai-memory message `01a1d42a-28f1-72d3-9957-76f106877d2e` of 2026-11-13 (06:07 BRT), with Roundfix 0.142.1. Its Delivery Queue item `0092-a-natureza-da-receita-da-skill-sai-da-tabela-oficial` parked with `delivery-error: read prerequisites: refresh prerequisite default branch "main": git fetch origin main: ... Recv failure: Operation timed out`, a passing network failure. The status said `resolve the blocker, then run roundfix deliver retry <item>`. With the network back (`git fetch` succeeding), the retry refused:

```text
Retry refused
Reason: retry Delivery Queue item "0092-...": use item branch "": use item branch: branch is required
```

The item parked before it had an item branch, and the retry requires one, so the item has no way back inside its queue. The only workaround is to wait for the owner to exit and start a new queue with the item, which replaces the earlier queue's record.

## Where

`roundfix deliver retry` for an item parked before its item branch exists: the refusal is raised at `internal/cli/deliver_workflow.go` (`use item branch: branch is required`), and the park comes from the prerequisite refresh (`refresh prerequisite default branch`).

## Expected

A retry of an item parked before its item branch exists returns it to `queued`, or creates the branch as `deliver start` would, and counts a retry. Optionally, a network failure while fetching the prerequisites' default branch is retried as transient before the item parks.

## Evidence

none yet — reported by the sender, not reproduced here. Sender's evidence: fiscus, `roundfix deliver status` and `retry` on 2026-11-13 06:07 BRT. Sender's repro: `deliver start` with two or more items, drop the network while an item without a branch fetches its prerequisites, then run `deliver retry` on it.
