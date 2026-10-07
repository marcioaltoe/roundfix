---
status: done
created_at: 2026-07-29
updated_at: 2026-09-08
absorbed_by: 0124-verification-capacity-and-measured-economics
---

# QA cycles — the gate serializes findings and sits far from the defects it must catch (2026-07-29)

Spec 0052 needed four QA cycles at roughly twenty minutes each, and a parallel Vortex session needed seven on one spec. Neither number came from bad implementation: in both repositories the Agents built what the Spec asked for. The cost came from two structural properties of the gate — it stops at the first static failure, so one stale assertion hides every flow finding and buys a whole cycle; and it is the only detector that exercises real external surfaces, so a defect a thirty-second probe would expose waits for a twenty-minute gate. This report separates what Roundfix must change from what the supervisor must change, because only the first belongs in a Spec.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-29-qa-cycle-latency-and-detector-placement.md`.
