---
status: done
created_at: 2026-09-29
updated_at: 2026-09-30
absorbed_by: 0187-a-queue-that-recovers-without-a-supervisor
---

# Delivery — the first full queue trial still needed manual recovery (2026-09-29)

Wave 4 (Specs 0181 and 0182) was the first cycle driven through one `roundfix deliver start`, as the 2026-09-29 handoff asked, using the v0.20.0 plan, limits, revalidation and renewed budget. The queue started at 14:01 with two retries per item. Neither item reached merge on its own: 0182 merged at 17:03 and 0181 around 18:40, both finished by the Supervisor. This report records each manual intervention and what it exposed. Specs 0181, 0182 and 0185 cover most of them.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md`.
