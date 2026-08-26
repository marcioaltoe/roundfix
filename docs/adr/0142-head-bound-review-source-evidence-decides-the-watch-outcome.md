---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Head-bound Review Source Evidence decides the watch outcome

Watch classifies Review Source signals as evidence bound to the pushed head: an explicit skip ends Review Skipped, a current-head approval with no unresolved threads proves Merge-Ready, and an exact Daemon-created artifact-only descendant may inherit its verified parent's evidence without another review request. When the Review Source's status check never appears within the grace period, the Run ends Clean Unverified — a distinct terminal outcome with its own exit code — so Clean keeps meaning affirmatively Merge-Ready and a caller can distinguish "verified merge-ready" from "pushed, but the Review Source never looked". Refines ADR-0019 and preserves ADR-0036's separate review-artifact commit.

Consolidates ADR-0043 and ADR-0054 (2026-08-26); both are archived under docs/history/adr/.
