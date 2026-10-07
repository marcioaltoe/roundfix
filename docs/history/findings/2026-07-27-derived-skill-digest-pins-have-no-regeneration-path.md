---
status: done
created_at: 2026-07-27
updated_at: 2026-09-08
absorbed_by: 0067-derived-artifact-regeneration-boundary
---

# Skill edits — derived digest pins have no regeneration path and block their own Task three times per day (2026-07-27)

Editing a Roundfix-owned Skill changes its `contentDigest`. Five Baseline artifacts pin that digest, and nothing in the repository regenerates them, so every Skill edit fails `make verify` until a human hand-propagates five hashes. Because the pins are not part of any Spec's tooling authorization, the Task that edits the Skill also fails its own changed-path check. This blocked work three separate times on 2026-07-27 — Spec 0037 task_07, the Spec 0037 QA gate, and Spec 0038 task_07 — each time costing a failed Run, a manual digest chase, and a Spec amendment.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-27-derived-skill-digest-pins-have-no-regeneration-path.md`.
