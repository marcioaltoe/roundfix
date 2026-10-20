---
type: fix
status: promoted
created: 2026-10-20
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
reason: null
---

# The queue skips the Linux Lane without a trace

## Symptom

On 2026-10-20 the `golang:1.26` image had been pruned from the maintainer's
machine. With `delivery.linux_lane: true` on the default branch, the Delivery
Queue gate of Specs 0285, 0286 and 0287 skipped the Linux Lane, and each of the
three candidates was published and merged. The only trace was one notice on the
Queue Owner's console log and in the item's lane log:

```text
roundfix: notice: Delivery Queue item <slug>: Linux Lane skipped: image golang:1.26 is not present locally
```

`roundfix deliver status`, the Pull Request the queue opened and `roundfix
doctor` said nothing, so nobody noticed for three releases that the lane the
maintainer had turned on did not run.

## Where

`internal/cli/deliver_linux_lane.go`: `runLinuxLane`'s `skip` closure (`:56`)
writes the notice and returns a passing result for every skip reason, among
them `image <image> is not present locally` (`:158`); `recordLinuxLaneSkip`
(`:410`) does the same when the fetched default cannot be read.
`internal/cli/deliver_workflow.go`, `Gate` (`:1456`), returns
`delivery.GateResult{Passed: true}` for every skip, so the Delivery Engine
(`internal/delivery/engine.go`, `gateCandidate`) records nothing on the item,
and `Publication` (`:1572`) writes the body `Delivers Spec <slug>.` whatever
the gate saw. `roundfix doctor` (`internal/cli/doctor.go`) has no Linux Lane
check.

## Expected

A skip stays a skip and never parks the item. When the lane is enabled and is
skipped for an environment reason (Docker missing or unreachable, the image
absent, an object store it cannot mount, a fetched default it cannot read), the
skip is visible where the delivery is read: on the item in `roundfix deliver
status`, in the body of the Pull Request the queue opens, and in `roundfix
doctor`, which reports the configured lane image missing or Docker unreachable
while `delivery.linux_lane` is `true`. A skip because nothing needs the lane
(no `go.mod`, no changed Go package) may stay quiet.

## Evidence

The operator's observation of 2026-10-20, recorded in the cycle AM operator
briefing for this Spec; the lane logs of the three items held only the notice
above. No Finding.
