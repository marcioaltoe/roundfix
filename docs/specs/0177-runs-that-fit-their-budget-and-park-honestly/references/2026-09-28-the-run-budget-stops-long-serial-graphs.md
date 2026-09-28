---
type: fix
status: promoted
created: 2026-09-28
spec: 0177-runs-that-fit-their-budget-and-park-honestly
reason: null
---

# The two-hour Run budget stops Specs whose Task Graph is a long serial chain

## Symptom

On 2026-09-28 two Onda 2 Runs reached `BudgetExceeded` with all implementation Tasks completed and only the QA Task (0176) or one Task plus QA (0173) pending: `budget.max_run_duration: 2h` in `.roundfixrc.yml` counts the whole Run, and graphs that `SC-WAVE-COLLISION` forces into a chain (every Task shares the instruction file or a guide) run their Tasks one after another. Each stop cost a relaunch and a fresh QA session.

## Where

`.roundfixrc.yml` `budget.max_run_duration`; budget enforcement in `internal/daemon`; `internal/speccheck` wave-collision rule counting the shared `instruction:` file as a collision.

## Expected

The budget scales with the graph (per Task or per critical-path length), or the QA Task is never started when the remaining budget cannot cover it; the instruction file every Task names does not count as a collision.

## Evidence

Runs `run_20260928T111930Z_13d79b0227e06da5` (0176) and `run_20260928T125114Z_30db2f43ddd8a501` (0173).
