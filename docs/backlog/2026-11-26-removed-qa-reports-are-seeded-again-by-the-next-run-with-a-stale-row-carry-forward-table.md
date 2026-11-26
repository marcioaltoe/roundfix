---
type: fix
status: open
created: 2026-11-26
spec: null
reason: null
---

# Removed QA reports are seeded again by the next Run, with a stale Row carry-forward table

## Symptom

In the delivery of Vortex Spec 0222, the whole `qa/` folder was removed from the item's branch, in a commit that also reopened the QA Task as pending. The next Run seeded the malformed report `-01` and the earlier failed report again, beside the new `-02`.

The `Row carry-forward` table seeded into `-02` lists the mechanical rows "QA Report structure" and "QA evidence paths", which do not exist in the Results of `-02`. The pre-PR review flagged `-01` twice and flagged that inconsistency.

## Open question

Where the Run takes the reports to seed from (the branch, an earlier Run's artifacts, or the Spec's history), and whether a deliberate removal on the branch should win. Not measured here: how the carry-forward table is built from the earlier Results.

## Source

Cross-project message `01a2168c-7b8e-7580-9d8e-8f08971b2f10` from the Vortex project, 2026-11-26. Nature: reported (observed in a Run audited with Roundfix 0.153.3, no response expected). Recorded by triage on 2026-11-26.
