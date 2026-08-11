---
status: accepted
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Tasks and QA own Agent Sessions

Each Task executes in its own Agent Session selected by Task Type, and a requested QA phase executes in a separate Agent Session selected by the `qa` Agent Work Category; review Runs keep a review-selected Agent Session for their review work. The effective selection and every attempted fallback are persisted per Work Item or action. This refines ADR-0018 so a spec Run is no longer constrained to one Run-wide Agent Session when different work categories require different profiles.
