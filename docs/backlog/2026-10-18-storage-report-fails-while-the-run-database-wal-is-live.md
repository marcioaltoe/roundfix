---
type: fix
status: open
created: 2026-10-18
spec: null
reason: null
---

# Storage report fails while the Run Database WAL is live

## Symptom

On 2026-10-18, while the Delivery Queue ran, `roundfix storage report`
failed with `reconcile Run Database storage groups: file bytes=1113714688
grouped bytes=1112350720 difference=1363968 tolerance=4096`. This was
observed by the Spec 0279 authoring session.

## Where

The storage report's reconciliation of the Run Database file size against
its grouped bytes.

## Expected

The report accounts for the live write-ahead log, or reads a consistent
snapshot, so it works while a queue runs. Its tolerance is not the fix.
