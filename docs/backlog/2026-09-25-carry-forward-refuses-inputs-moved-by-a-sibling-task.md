---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# Carry-forward refuses a Run whose Tasks moved each other's inputs

## Symptom

`roundfix reconcile <run> --carry-forward` refused the whole settled set of Spec 0170's first Run: `Task task_04 declared input(s) moved: CONTEXT.md`. `CONTEXT.md` was changed by task_01 of the same Run, which is already integrated on the Run Branch, so the move is the Run's own work, not drift on the checkout. The operator had to fast-forward the item branch to the Run Branch by hand.

## Where

`internal/cli` reconcile carry-forward input-move check.

## Expected

An input moved only by commits of the same Run (sibling Tasks already on the Run Branch) does not refuse the carry-forward; only a move on the target branch does.

## Evidence

Run `run_20260925T182023Z_7e3e8c4b8bc65d93`, 2026-09-25.
