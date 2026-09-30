---
type: fix
status: open
created: 2026-09-30
spec: null
reason: null
---

# A retry refused for moved inputs does not say to carry forward first

## Symptom

On 2026-09-29 the Supervisor amended Spec 0181 on its item branch and ran `deliver retry`. It refused the whole carry-forward set because every Task's declared inputs had moved. The refusal's next action named only `reconcile <run> --carry-forward`. It did not say that the amendment itself caused the refusal and has to land after the carry-forward. The recovery took a reset, the carry-forward and a cherry-pick of the amendment.

## Where

The carry-forward refusal text in `internal/cli/carryforward.go` and `deliver retry` in `internal/cli/deliver_workflow.go`.

## Expected

When every refused Task's moved inputs are Spec artifacts changed by commits on the item branch after the Run started, the refusal names those commits. It says to carry forward from the Run's branch state first and reapply the amendment afterwards, and it prints the exact commands.

## Evidence

`docs/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md`, section 2.
