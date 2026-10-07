---
status: done
created_at: 2026-07-30
updated_at: 2026-09-08
absorbed_by: 0067-derived-artifact-regeneration-boundary
---

# `make baseline-digests` cannot bootstrap a Baseline module edit (2026-07-30)

Delivering `rule.autonomous.loop` through the Baseline module hit three distinct failures in the derived-artifact contract. All three surface on any module edit that changes a generated guide, which is the ordinary case, so they block the autonomous Spec loop rather than being incidental to one change.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md`.
