---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Run outcome is derived from unresolved work, not from a failed Batch

`MarkBatchFailed` overwrote every Review Issue in a Batch when the Batch failed.
On pull request 143 round 001 that turned twenty recorded resolutions into
`failed`; each carried triage notes describing work that a spot check confirmed
had landed in the tree.

The obvious repair is to skip issues already in a Terminal status, exactly as the
sibling `SettleAssignedIssues` already does. It was tried on 2026-08-09 and
reverted within the hour, because six tests in `internal/cli` and
`internal/daemon` encode the opposite contract: a Batch whose Agent failed must
leave its issues `failed` so the Run ends `Unresolved` with exit 1 and a later
Round retries. Preserve the resolutions and the same Run reaches `Clean` — it
reports success for a Run whose Agent crashed.

So the defect is not the marking. It is that two different questions share one
answer. "Did this Batch finish?" and "is there unresolved work left?" are not the
same question, and deriving the second from the first is what forces the
overwrite.

Roundfix therefore records Batch settlement and Run outcome separately. A Batch
records what its Agent achieved, including on a Batch that failed, and an issue
the Agent resolved keeps that outcome. The Run's outcome is computed from whether
unresolved Review Issues remain, so a Batch failure with seventeen issues
resolved and four still open ends `Unresolved` on the four, not on the failure.
Both statements stay true and both stay visible.

Keeping the overwrite was rejected because it destroys work and invites a later
Round to redo it, and because it destroys the evidence that the Batch failed
partway — the twenty issues are the only record that this happened at all.
Preserving outcomes while leaving the Run outcome derived from Batch failure was
rejected because it is the version that reports `Clean` on a crashed Agent.
