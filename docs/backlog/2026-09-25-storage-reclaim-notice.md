---
type: feat
status: open
created: 2026-09-25
spec: null
reason: null
---

# Tell the operator when Run storage is reclaimable

## Opportunity

Run artifacts and the Run Database grow with every Run; only `gc` reports reclaimable bytes, and only when a human runs it (4.4 GB with 619 MB reclaimable after one day).

## Value

The operator learns about reclaimable space where they already look, without any automatic deletion.

## Shape

Reuse the storage measurement and print one threshold line in `runs list` or Doctor pointing to `roundfix gc`.

Evidence: secondbrain `inbox/roundfix/2026-09-17-os-artefatos-de-run-crescem-e-so-um-humano-recupera.md`.
