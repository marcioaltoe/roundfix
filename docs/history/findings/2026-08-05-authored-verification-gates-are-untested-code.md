---
status: done
created_at: 2026-08-05
updated_at: 2026-09-08
absorbed_by: 0129-spec-authoring-and-gate-recovery
---

# 2026-08-05 — Authored verification gates are untested code

Source: the second real Supervisor session in the `fiscus` repository, spec `0002-auth-staff-e-directory` (runs `run_20260805T131149Z_18c1483e9d4abee4` Stopped and `run_20260805T134406Z_41a3bfa9659f0917` in flight), one day after the spec-0001 night that produced `2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md`. The unifying observation: **a Task's `## Verification` commands are the only code in the pipeline that is never executed before it starts deciding outcomes.** Implementation passes through the gates; the gates themselves ship untested. One run burned three defective gates authored by a careful Supervisor following the current `write-tasks` contract — the contract itself is missing a step. This finding proposes improvements to the `write-*` skills (whose Baseline roundfix owns) and to Roundfix behavior.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-05-authored-verification-gates-are-untested-code.md`.
