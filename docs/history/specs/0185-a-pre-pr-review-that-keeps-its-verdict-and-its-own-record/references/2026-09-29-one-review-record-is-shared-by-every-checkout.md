---
type: fix
status: promoted
created: 2026-09-29
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
reason: null
---

# One pre-PR review record is shared by every checkout of a repository

## Symptom

`roundfix review` writes `pre-pr-review.json` and `pre-pr-review-answer.txt` at the repository's Artifact Root. Since the durable repository key (Spec 0162), every checkout of a repository shares that root: the main checkout, extra worktrees, and each Delivery Queue item worktree. On 2026-09-29, reviews run from `~/dev/roundfix` and from `~/dev/roundfix-wave5` overwrote the same file. Its `repository` field changed from one checkout to the other.

While a Delivery Queue item is in its `reviewing` stage, an operator review from any other checkout can therefore replace the record the delivery engine is about to read. The dispositions file is shared too; it is keyed by head, which keeps it consistent.

## Where

The review record path under the Artifact Root in `internal/cli/review.go`, and the delivery engine's read of that record (`internal/delivery/engine.go` `reviewCandidate` → `internal/cli/deliver_workflow.go`).

## Expected

A review record belongs to the checkout, or the head, it reviewed. A concurrent review from another checkout never replaces a record another reader depends on. A reader that finds a record for a different checkout or head refuses it instead of acting on it.

## Evidence

`~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review.json`, 2026-09-29. At 13:50 the `repository` field read `/Users/marcio/dev/roundfix` (head `fe517d2b`). At 14:20 it read `/Users/marcio/dev/roundfix-wave5` (head `3e77ed70`). No delivery item was in `reviewing` at those times, so no queue decision was affected.
