---
type: fix
status: promoted
created: 2026-09-30
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
reason: null
---

# Authored Verification runs without a provenance check

## Symptom

The Baseline tells Agents to execute authored Verification only on committed provenance: the Spec artifacts that carry the commands must be tracked and byte-identical to their committed bytes. Until 2026-09-30 it also described an `execution_approvals` record that would let an untrusted source run. No Go source reads that field, and no command checks the provenance.

- `roundfix spec check --run-verification` loads the Task Graph from the working tree and executes its commands in a checkout of `HEAD`. It runs the commands of an uncommitted or modified Task file without a warning.
- `roundfix settle` reads the Task file of the work directory it settles.
- A Run is safe by construction, because its Run Worktree is created from a commit.

Spec 0119's Task 06 recorded the control as implemented, with a refusal code `SC-SOURCE-UNTRUSTED`. `git log -S` finds that code only in documentation.

## Where

`internal/cli/spec_check.go` (`probeSpecVerifications`), `internal/cli/settle.go`, and the clause `clause.core.verification-two-tiers` in `internal/baseline/assets/modules/core.json`.

## Expected

A command that executes authored Verification outside a Run compares the carrying artifact with its committed projection first. It refuses, or names, a command whose source is untracked or modified. Whether an approval record for untrusted sources should exist at all is a decision for the Spec.

## Evidence

`git grep -l "execution_approvals\|SOURCE-UNTRUSTED" -- '*.go'` on `9e439dbb` returns nothing. On 2026-09-29 and 2026-09-30 this session ran `spec check --strict --run-verification` on more than ten uncommitted Specs, and each run executed their commands. Spec 0193 rewrites the clause to say what is and is not enforced.

## Addendum 2026-10-01: a malformed command passes the honesty probe

During Spec 0205's delivery, an amended Verification command contained escaped backticks inside its inline-code span. The span ended at the first backtick, and the extracted command `grep -qF '| \` was a shell syntax error (`unexpected EOF while looking for matching`). It failed before the work, so `spec check --run-verification` reported it `honest`. It also failed after the work, so task_03 failed twice even though the implementation was complete (Run artifacts `run_20261001T214123Z_1a226c6ee6abc317/verification/batch-001-attempt-*.log`).

Expected, in the same command family: a Verification command that the shell cannot parse, or that `sh -n` rejects, is reported as malformed, not as honest. `spec check --strict` refuses an inline-code Verification whose span contains a backslash-escaped backtick.
