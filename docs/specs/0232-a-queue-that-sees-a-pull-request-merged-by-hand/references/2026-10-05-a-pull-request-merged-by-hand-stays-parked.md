---
type: fix
status: promoted
created: 2026-10-05
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
reason: null
---

# A Pull Request merged by hand leaves its Delivery Queue item parked

## Symptom

When the operator merges a parked item's Pull Request by hand, after fixing
CI on the item branch or after a refused retry, the Delivery Queue never
learns it. On 2026-10-05 Spec 0231's Pull Request #404 was merged by hand
after a real CI fix; `roundfix deliver status` still prints
`0231-checks-that-hold-in-delivery parked checks-failed` with a Pending
Question asking for a retry, and `roundfix deliver retry` cannot answer it.
The same happened to Spec 0225 (Pull Request #382) and Spec 0229 (Pull
Request #396), whose retries were refused before the operator opened their
Pull Requests by hand.

## Where

- `Engine.Run` in `internal/delivery/engine.go` skips every parked item, so
  `deliver resume` never looks at one.
- `deliver status` reads only the Run Database.
- `Engine.Retry` in `internal/delivery/engine.go` never asks whether the
  item's Pull Request was merged or whether its Spec is already archived on
  the default branch. It first requires the item branch, and then applies the
  archived-head rules of ADR-0223 and ADR-0229.

## Expected

A retry of a parked item whose Pull Request is merged, or whose Spec the
default branch already archives through a delivery commit outside the item
branch, records the item `merged` and lets the owner run the normal
post-merge cleanup. A Pull Request closed without merging keeps the item
parked.

## Evidence

Read on 2026-10-05:

- Operator intervention log entries 175 to 177 (0231), 158 (0225) and 169
  (0229).
- `gh pr view 404`: state `MERGED`, head `dca3aa19` on branch
  `roundfix/deliver-0231-checks-that-hold-in-delivery-02e9e83bd5a5644b`,
  merge commit `4b5ea48b`. The queue item records Pull Request `404` and
  candidates `ea6135c3`, `09a30447`; `dca3aa19` is the operator's CI fix on
  top of `09a30447`.
- `gh pr view 382` and `gh pr view 396`: both `MERGED`, opened by hand, so
  their items recorded no Pull Request number.
- The local item branch and worktree of 0231 no longer exist, so a retry
  today stops at `item branch ... is missing`.
