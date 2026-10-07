---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
absorbed_by: 0129-spec-authoring-and-gate-recovery
---

# Spec authoring — rules a release night made checkable (2026-08-06)

Eight Specs went from queue to released v0.4.0 in one supervised session. Every QA gate that failed did so on an authoring defect, not an implementation one, and each defect is mechanically checkable at authoring time — the same conclusion `2026-08-06-every-run-that-failed-tonight-failed-on-a-contract.md` reached for Daemon contracts, now reached for Spec artifacts. Spec 0065 shipped the `SC-*` enforcement surface the same night; these are its next rules.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-06-authoring-rules-a-release-night-made-checkable.md`.
