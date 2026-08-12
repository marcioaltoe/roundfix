---
task: task_04
spec: 0062-baseline-digest-regeneration-bootstrap
status: completed
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

## Result

Added one ordinary `TestCatalogCompatibility` run after the post-regeneration
artifact comparison and immediately before the existing success branch. The
test omits `-update`, so it loads through the public strict catalog path after
all five regeneration steps have rewritten their artifacts. A failure exits
the recipe and reaches the existing structured failure trap with
`errorCode: strict_validation_failed`, `stage: strict-validation`, and
`retryable: false`.

Focused checks:

- Red signal: pre-change `rtk make -n baseline-digests` showed the five
  `-update` steps proceeding directly through the post-scan and comparison to
  the success output, with no strict validation invocation.
- Post-change `rtk make -n baseline-digests` exited 0 and showed the unchanged
  five-step list, pre/post derived-path scans, comparison, the new ordinary
  strict validation, and then the unchanged success output.
- `rtk go test ./internal/baseline -run TestCatalogCompatibility -count=1`
  exited 0, exercising `LoadEmbeddedCatalog` against the current consistent
  catalog without regeneration mode.
- `rtk git diff --check` exited 0.
- `rtk git -c core.fsmonitor=false status --short` and
  `rtk git -c core.fsmonitor=false diff --name-only` exited 0 and listed only
  `Makefile` and this task file.

Acceptance evidence:

1. The focused ordinary catalog-compatibility test accepted the current
   consistent catalog, and recipe inspection shows the existing no-change
   success branch remains after that check. The Daemon-owned end-to-end target
   invocation remains under `## Verification`.
2. The dry-run recipe places `stage: strict-validation` after every `-update`
   step, the post-scan, and the changed-artifact comparison; its command omits
   `-update`.
3. The strict command preserves its nonzero status and exits the recipe on any
   diagnostic. The existing EXIT trap then emits the structured failure with
   `errorCode: strict_validation_failed` and `stage: strict-validation` rather
   than reaching success output.
4. The `Makefile` diff adds only the three strict-stage recipe lines. The
   `BASELINE_DIGEST_STEPS`, `DERIVED_DIGEST_PATHS`, both scans, changed-artifact
   comparison, and success-output line are unchanged.
5. Final changed-path inspection is limited to `Makefile` and this task file;
   the task status line was the pre-existing Daemon-owned change.

The commands under `## Verification` were not run; the Daemon owns them.
