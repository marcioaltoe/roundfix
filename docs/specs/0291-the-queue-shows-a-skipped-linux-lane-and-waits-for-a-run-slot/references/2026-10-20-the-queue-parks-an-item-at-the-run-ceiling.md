---
type: fix
status: promoted
created: 2026-10-20
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
reason: null
---

# The queue parks an item at the Active Implement Run ceiling

## Symptom

On 2026-10-20 the Delivery Queue item of Spec 0289 started its Implement Run
while other repositories' Runs held every slot of the machine-wide ceiling
(`runs.max_active`, User Config, `3` by default). The child `roundfix
implement` refused at preflight and the queue parked the item as
`delivery-error`, with the whole preflight report as its reason:

```text
delivery-error: run Implement executor: roundfix implement failed with exit code 2: Preflight failed

Reason:
  Active Implement Run ceiling reached: 3 Run(s) hold 3 slot(s)
- run_id=<id> repository=<repository> spec=<spec>; stop it with: roundfix stop <id>
...
```

The item stayed parked until the operator noticed and ran `roundfix deliver
retry`, although the only thing missing was a free slot that the other Runs
gave back a little later.

## Where

`internal/store/store.go`, `ActiveImplementRunCeilingError` (`:281`), is the
refusal; `internal/cli/implement.go` prints it at preflight (`:404`) and
`CreateRun` re-checks it inside its write transaction. The queue runs the
child through `commandDeliveryWorkflow.RunSpec`
(`internal/cli/deliver_workflow.go`); `deliveryRunResult` (`:1065`) turns every
preflight exit into the error `roundfix implement failed with exit code 2:
<stderr>` (`:1100`), and `Engine.Run` (`internal/delivery/engine.go:741`)
parks it as `delivery-error`.

## Expected

When the only preflight refusal is the ceiling, the item waits for a slot
instead of parking: `roundfix deliver status` shows it waiting for a Run slot
with the Runs holding the slots and the deadline, the queue polls for a free
slot and starts the Run when one frees. The wait is bounded by a configurable
duration (2h by default; `0` parks at once as today) and by the queue's own
deadline when that comes first; at the end of the wait the item parks exactly
as today. The store's ceiling check stays the authority, so two queues waiting
for one slot never both start a Run beyond the ceiling.

## Evidence

The operator's observation of 2026-10-20 (Spec 0289 waited for other
repositories' Runs), recorded in the cycle AM operator briefing for this Spec.
No Finding.
