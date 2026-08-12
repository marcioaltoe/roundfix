---
status: accepted
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Release planning is read-only and confirmation-gated

Roundfix exposes release analysis as a read-only Release Plan Command that classifies committed changes, proposes the next semantic version, cites the evidence, and names any manual classification or approval still required. It never edits release files, creates or pushes tags, publishes packages, or creates a GitHub Release; those mutations remain separate user-directed actions after the Release Plan is accepted. Combining analysis and release execution was rejected because a generic release request must not silently authorize a minor, major, or breaking version decision, and a deterministic plan can be rerun and audited before any irreversible action.
