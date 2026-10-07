---
status: done
absorbed_by: 0124-verification-capacity-and-measured-economics
created_at: 2026-08-14
updated_at: 2026-09-08
kind: finding
---

# Preflight starves when the machine is busy (2026-08-14)

`profile proof failed … adapter error: context deadline exceeded` was already on record from 2026-08-08 as tuple unavailability. This adds the missing variable: it is not random intermittence. It appears when other Roundfix Runs are active on the same machine and disappears when they finish. Minted from the Inbox Entry `inbox/roundfix/2026-08-12-preflight-e-inanicao-por-runs-concorrentes.md` in the Secondbrain, observed in `conexus` on 2026-08-11 and 2026-08-12.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-14-preflight-starves-when-the-machine-is-busy.md`.
