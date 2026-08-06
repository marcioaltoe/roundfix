---
status: done
created_at: 2026-08-04
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-04 — Fail-fast Verification spends the single repair turn on the first of N defects

status: pending

## What was observed

Task 07 of fluxus Spec 0012 failed with two independent defects in the working
tree, and the repair protocol can only ever address one of them.

The repository Verification is `make verify`, which runs fmt, lint, onioncry,
typecheck, and tests in sequence and stops at the first non-zero exit.

1. **Attempt 1** failed at the architecture step. `batch-007-attempt-1.log`
   ends at `@fluxus/backend#onioncry: status: fail`, 458 files checked. The
   later steps never ran, so nothing downstream was known.
2. The Daemon sent its one Verification Feedback prompt. The Agent repaired
   exactly what the log showed — renamed a helper to satisfy the layer rule —
   and reported honest focused evidence: onioncry clean, typecheck passing,
   `git diff --check` clean.
3. **Attempt 2** ran the full sequence for the first time and failed at the
   test step: a ZodError in `openapi-schemas.test.ts`, three required fields
   missing from a pre-existing fixture the Task had not updated.
4. There is no third attempt and no second repair prompt, so the Task settled
   `failed` with a defect the Agent had never been shown.

The Agent did nothing wrong. It fixed the only failure it was given, and the
second defect was invisible until the first was gone.

Recovering cost a Supervisor intervention and a complete new Run.

## Root cause

The repair budget is one turn, but the Verification that feeds it reveals
defects one layer at a time. When a Task leaves N independent problems, the
loop needs N-1 turns it does not have, and the Agent's single turn is spent
against a partial picture by construction.

The two mechanisms are individually reasonable and jointly broken: fail-fast
is correct for a developer at a terminal who reruns freely, and one repair
turn is a correct budget for a bounded loop. Combining them caps the loop at
"one defect per Run" whenever the gate is a sequential pipeline.

Nothing in the contract tells the Agent the picture is partial, so it cannot
compensate by looking wider on its own.

## What would settle it

Give the repair turn the complete failure set rather than the first failure.
Options, cheapest first:

- Let a repository declare its Verification as independent steps that all run
  before the verdict, so attempt 1 reports every failing step at once. This
  costs one full gate execution instead of stopping early, and buys a repair
  turn that sees everything.
- Failing that, state in the Verification Feedback prompt that the sequence is
  fail-fast and later steps did not run, so the Agent knows to sweep for
  adjacent breakage rather than trusting the log's silence as evidence of
  health.
- Consider making the retry budget a function of distinct failing steps rather
  than a constant, so a Task with two unrelated defects is not fated to fail.

## Related

The same Spec produced
[[2026-08-04-a-static-gate-row-reported-one-instance-per-cycle]]. Both are the
same shape at different altitudes: a checker that stops at the first problem,
feeding a loop with a bounded number of attempts.

## Spec pointer

None yet.
