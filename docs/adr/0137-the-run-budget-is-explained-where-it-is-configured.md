---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-16T09:01:00Z
updated_at: 2026-08-16T09:01:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The run budget is explained where it is configured

A maintainer setting a maximum Run duration cannot tell from the setting whether
it bounds wall-clock time from the Run's start or is evaluated at Work Item
boundaries, and the single measured overrun has no established cause — so
changing where the budget is evaluated would be a guess dressed as a fix. What it
bounds is therefore stated where it is set, in the configuration surface the tool
renders, rather than in a guide a maintainer reads separately or not at all. If a
reproduction later shows the evaluation point is wrong, that is a behaviour change
with its own evidence; documenting the current contract does not prejudge it.
