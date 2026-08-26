---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Release planning is read-only, including its reset mode

The Release Plan Command classifies committed changes, proposes the next semantic version, cites its evidence, and names any manual classification or approval still required — and never edits release files, creates or pushes tags, publishes packages, or creates a GitHub Release; those mutations remain separate user-directed actions after the plan is accepted. `--reset-to v0.0.1` is a mutually exclusive mode of the same boundary: it inventories prior tags and GitHub Releases, produces a digest-bound reset plan requiring approval, and exposes no deletion path — removal needs a separate, explicitly authorized post-QA action. A generic release request must never silently authorize a version decision.

Consolidates ADR-0048 and ADR-0065 (2026-08-26); both are archived under docs/history/adr/.
