---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-14T17:50:12Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: ADR-0148
---

# Authoring and Run time share one Verification prober

The Daemon already runs every authored Verification command against the unchanged
tree before opening an Agent Session and classifies each as vacuous, failing, or
unknown; the authoring check this Spec adds asks the identical question earlier,
so the classification is extracted into one prober that both callers use rather
than reimplemented beside it. A second implementation would be free to disagree,
and a checker that approves what the probe later refuses is the defect this Spec
exists to remove, reappearing one layer up. The Daemon keeps its Run bookkeeping —
run state, verification capacity, artifact paths — around the shared loop, and the
authoring caller supplies a working directory and nothing else.
