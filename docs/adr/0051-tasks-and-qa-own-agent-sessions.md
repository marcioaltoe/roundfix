---
status: accepted
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-10-14T23:44:53Z
deprecated_at: null
superseded_by: null
---

# Tasks and QA own Agent Sessions

Each Task executes in its own Agent Session selected by Task Type, and a requested QA phase executes in a separate Agent Session selected by the `qa` Agent Work Category; review Runs keep a review-selected Agent Session for their review work. The effective selection and every attempted fallback are persisted per Work Item or action. This refines ADR-0018 so a spec Run is no longer constrained to one Run-wide Agent Session when different work categories require different profiles.

Consolidates ADR-0018 (2026-10-14): an Agent Session is still never reused
across Runs and every Run still starts from a clean session, so cross-Run
reproducibility holds; the single Run-wide session it decided survives only
as the review Run's review-selected Agent Session. Git keeps ADR-0018's text.
