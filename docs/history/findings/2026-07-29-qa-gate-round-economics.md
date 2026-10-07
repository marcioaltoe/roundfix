---
status: done
created_at: 2026-07-29
updated_at: 2026-09-08
absorbed_by: 0124-verification-capacity-and-measured-economics
---

# QA gate — discovery stops at the first hard gate, so one Spec pays for five expensive rounds (2026-07-29)

Evidence from implementing Spec `0010-contrato-de-health` in the `vortex` repository, a Bun/TypeScript monorepo. The Spec shipped correct behavior, but it needed **8 Tasks and 7 QA rounds**. Reading the reports back, most of that cost is not defect density — it is the order in which the gate discovers things, plus an under-provisioned gate environment. The three findings below are what a Supervisor cannot work around from the outside.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-29-qa-gate-round-economics.md`.
