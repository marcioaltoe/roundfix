---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A concurrent `skills-lock.json` rewrite during the fetch can be overwritten

## Symptom

`baseline skills reconcile` and `baseline skills restore` read
`skills-lock.json`, then fetch the source, then snapshot the lock's preimage for
the transaction. A rewrite of the lock during the fetch (for example by
`npx skills add`) becomes the preimage, the transaction check passes, and the
result computed from the old bytes overwrites the new edit without error.

## Where

`internal/baseline/skills_reconcile.go` and `internal/baseline/skills_restore.go`
(lock read before fetch; preimage snapshot after).

## Expected

Fail with a stale-plan refusal unless the snapshot equals the bytes the plan was
computed from, or load the lock after the fetch.

## Evidence

Recorded limit of Spec 0163 (pre-PR review of 2026-09-25), found by reading the
code.
