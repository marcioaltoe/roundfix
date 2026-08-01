---
task: task_04
spec: 0062-baseline-digest-regeneration-bootstrap
status: pending
type: infra
complexity: medium
---

# Task 04: Re-validate strictly after regeneration

## Overview

Deferring a derived pin is only safe if something checks it afterwards.
Otherwise the previous slice traded a blocked run for a run that can finish
green on an inconsistent catalog. This Task closes the regeneration target with
one strict load, so a deferred pin is a pin checked later rather than a pin
abandoned. It is the authorized tooling change of this Spec.

## Requirements

1. MUST perform a strict catalog validation as the final stage of the
   regeneration target, after every regeneration step and after the derived
   artifacts have been rewritten.
2. MUST fail the target when that strict validation reports any diagnostic, so
   regeneration cannot report success on a catalog that is still inconsistent.
3. MUST report the failure in the target's existing structured failure shape,
   naming the stage, so the loop's tooling can consume it the way it consumes
   the other stages.
4. MUST leave the target's existing behavior intact: the step list, the derived
   path scan, the changed-artifact comparison, and the success output shape.
5. MUST change only the bounded authorized path `Makefile`, per the maintainer
   authorization recorded in this Spec's Tooling authority row. Any other path
   is out of scope and fails this Task.

## Subtasks

- [ ] Add the final strict validation stage to the regeneration target.
- [ ] Fail the target on any diagnostic from that stage.
- [ ] Report the failure in the existing structured shape with its stage name.
- [ ] Confirm the target's existing stages and outputs are unchanged.

## Acceptance Criteria

- [ ] On a consistent tree, the regeneration target completes successfully and
      reports no change.
- [ ] The target performs a strict validation after its regeneration steps, and
      that stage is named in the target.
- [ ] A catalog left inconsistent after regeneration fails the target rather
      than reporting success.
- [ ] The target's step list, derived path scan, and success output shape are
      unchanged.
- [ ] `git status --porcelain` shows no path outside `Makefile` and this task
      file.

## Context

- instruction: `docs/workflow/authorizations/2026-08-01-baseline-digest-bootstrap.md`
- interface: `Makefile`

## Verification

- `make baseline-digests` — expected: exit 0 on the unmodified tree, with
  `"changed":false` in its structured output, proving the target still
  completes and rewrites nothing when everything already agrees.
- `make baseline-digests` — expected: exit 0 on a second consecutive run with
  `"changed":false`, proving idempotence and that the new stage does not
  perturb derived artifacts.
- `grep -q "BASELINE_DIGEST_STEPS" Makefile` — expected: exit 0; the step list
  survived.
- `grep -q "DERIVED_DIGEST_PATHS" Makefile` — expected: exit 0; the derived path
  scan survived.
- `git diff --name-only HEAD | grep -v '^Makefile$' | grep -v '^docs/specs/0062' | grep -q . && exit 1 || exit 0`
  — expected: exit 0; nothing outside the bounded path changed.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 2; Project Constraints: Tooling authority.
- `_techspec.md` → Build Order 4; Risks (deferral could mask a real mismatch).
- ADR-0081, ADR-0085.
