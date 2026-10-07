---
status: done
created_at: 2026-08-04
updated_at: 2026-09-08
absorbed_by: 0129-spec-authoring-and-gate-recovery
---

# 2026-08-04 — Pre-contract Spec graphs run with no QA gate and say nothing

Roundfix 0.3.1 removed the `--qa` flag and moved the gate into the Task Graph as an authored `qa: task_NN` declaration. In a consumer repository (oraculum) with eleven queued Specs authored under the old flag, **none of the eleven declares a gate**:

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-04-pre-contract-spec-graphs-run-with-no-qa-gate-and-say-nothing.md`.
