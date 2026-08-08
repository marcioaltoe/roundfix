---
status: accepted
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An unrecorded managed region is refreshed and named

A managed region whose bytes differ from the recorded digest carries no
information about why, so the choice is between refusing every such repository
and replacing content a human may have written. The Baseline already resolves
the ownership question: guidance inside `setup-context-driven` markers belongs to
the Baseline, and proposing a change to it is an Inbox Entry to the Baseline's
owner rather than a local edit. What was wrong was not the replacement but the
silence around it. An unrecorded managed region is therefore refreshed, and the
presented plan names it by path and managed identity and lists every line on
disk that the refreshed rendering does not reproduce, so approval is given with
the replacement in view; the existing Plan Digest confirmation remains the only
gate. Blocking was rejected because it makes the first update of every
pre-existing repository unrecoverable, and refreshing silently was rejected
because divergence a maintainer cannot see is indistinguishable from data loss.
