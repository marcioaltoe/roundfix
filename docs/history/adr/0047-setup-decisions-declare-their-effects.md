---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null
superseded_by: ADR-0144
---

# ADR-0047: Setup decisions declare their effects

The setup asset catalog declares how each durable decision activates modules, includes or excludes managed artifacts, selects templates, introduces dependent decisions, or binds a value into generated guidance. Preview, audit, and apply resolve the same declarative Decision Plan; imperative per-decision branching and profile-per-combination variants were rejected because they would let the three paths drift or multiply the profile catalog.
