---
status: done
created_at: 2026-08-03
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# 2026-08-03 — Gate and review rounds need a convergence rule

status: pending

## What was observed

Closing Spec 0072 took four gate executions and two resolve Runs, and the
oscillation had no terminating rule of its own — it converged only because
the Supervisor started checking the review state before each gate attempt.

The loop, as it actually ran:

1. Gate #2 failed: PR #87 carried eight unresolved CodeRabbit threads.
2. `roundfix resolve` closed all eight and pushed — which gave the PR a new
   head, which triggered a fresh CodeRabbit review.
3. Gate #3 failed: four *new* threads on the new head.
4. Resolve round two closed those four and pushed — new head again.
5. This time the Supervisor waited for CodeRabbit to finish and checked
   `reviewDecision` before launching the gate; the PR was APPROVED with zero
   unresolved threads, and gate #4 found F-003 resolved.

Every fix round creates the next review round. The gate reads "unresolved
threads at current head" at whatever moment it runs, so gate-after-resolve
races the reviewer: launched too early it fails on threads the reviewer has
not finished opening, or passes on a head the reviewer has not finished
reading.

## Root cause

Neither the gate nor the resolve flow owns the ordering between "the head
settled" and "the reviewer's verdict on that head exists". The knowledge
lives only in operator discipline.

## What would settle it

A convergence precondition on the gate's review-readiness row: the gate
observes the reviewer's *decision state for the exact current head*
(complete review, zero unresolved, approved) rather than a thread count at
an arbitrary moment — and when the reviewer has not yet reported on the
current head, the row blocks as *environment: review pending* instead of
failing the Spec. `roundfix watch --until-clean` already embodies the right
loop shape for the resolve side.

## Spec pointer

None yet.
