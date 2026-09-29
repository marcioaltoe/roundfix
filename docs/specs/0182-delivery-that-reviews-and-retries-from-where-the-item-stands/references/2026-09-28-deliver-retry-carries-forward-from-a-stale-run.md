---
type: fix
status: promoted
created: 2026-09-28
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
reason: null
---

# `deliver retry` carries forward from the Run recorded on the item, not the newest Run

## Symptom

Spec 0175's item was retried twice after `BudgetExceeded`. The second retry tried to carry forward from the first Run (`run_20260928T154312Z_63679c6540b34603`, already carried: "carry-forward refused the whole set: ... declared input(s) moved") instead of the newest Run (`run_20260928T174529Z_c0cad0da5734a85f`), so the item branch did not receive Tasks 03–05 and `implement`'s preflight refused ("holds proved Tasks available for Task Carry-Forward"). The item parked again as `delivery-error`; the operator ran `roundfix reconcile <newest-run> --carry-forward` by hand and retried.

## Where

`internal/cli/deliver_workflow.go` / `internal/delivery` retry recovery (Spec 0173): the Run chosen for carry-forward.

## Expected

Retry carries forward from every terminal Run of the item's Spec on the item branch that still holds proved Tasks, newest first, and treats a Run whose Tasks are already on the item branch as nothing to carry rather than a refusal.

## Evidence

Delivery owner console logs `delivery-339f8dac2b68-53886` and `-61946`, 2026-09-28.
