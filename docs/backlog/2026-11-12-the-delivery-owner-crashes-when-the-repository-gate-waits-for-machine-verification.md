---
type: fix
status: open
created: 2026-11-12
spec: null
reason: null
---

# The delivery owner crashes when the repository gate waits for Machine Verification

## Symptom

On 2026-11-12 at 10:20 UTC, Roundfix 0.141.1 on the Mac mini reached the repository gate for Spec 0382 while a Quaestor Run held the only Machine Verification Capacity slot. The detached Delivery Owner then failed with `panic: runtime error: invalid memory address or nil pointer dereference` in `fmt.Fprintf` and exited. `roundfix deliver status` kept reporting the item as `gating` with no owner for about 8.5 hours, and no status line showed that the owner was gone. `roundfix deliver resume` started a new owner, and the item continued.

## Where

`commandDeliveryEngineDependencies` in `internal/cli/deliver_workflow.go` builds `commandDeliveryWorkflow` without setting `log`; only tests and `deliver lane` set it. `Gate` writes to `workflow.log` with no nil check in three places:

- the Machine Verification wait callback, which panicked here;
- the shared-cache trim report;
- the lease-release warning.

The item-binary probe already falls back to `os.Stderr` when `log` is nil.

## Expected

- The production workflow always has a log writer, the owner's console log, so every item whose gate waits for capacity logs the wait instead of crashing.
- `deliver status` reports a queue whose recorded owner is no longer running instead of showing the item as still active.

## Evidence

- On the mini, the owner's console log is at `~/.roundfix/artifacts/339f8dac2b687a04/delivery/delivery-339f8dac2b68-49802/console.log`. Its stack trace names `deliver_workflow.go:1713` (`Gate.func1`) called from `store.AcquireMachineVerification`.
- The defect has existed since Spec 0374 added the capacity wait at the gate. Until now it went unseen because the gate rarely waited.
