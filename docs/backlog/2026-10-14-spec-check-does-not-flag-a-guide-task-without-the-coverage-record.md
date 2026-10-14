---
type: feat
status: open
created: 2026-10-14
spec: null
reason: null
---

# Spec check does not flag a guide Task without the coverage record

## Symptom

Spec 0268's task_01 edited three user guides, which changes their
fingerprints in `docs/references/behavior-surfaces.json`, but its
`_authorization.md` Limits forbade that file. Verification failed twice on
`TestTheSkillCoverageMapIsCurrent`, the Agent correctly refused to regenerate
the record, and the Run ended Unresolved. `roundfix spec check --strict` was
clean on the plan.

## Where

`internal/speccheck`, for a pending Task that declares a path under
`docs/user-guide/` or a command's help.

## Expected

`spec check` reports a gap when such a Task neither declares
`docs/references/behavior-surfaces.json` with the `-record-skill-coverage`
step nor is covered by another Task that does, and when the Spec's
authorization forbids that file.

## Evidence

0268 Run `run_20261014T102741Z_dd0d6c13d87ea4f9`, task_01 Verification logs;
the fix was plan commit `cc0ac876`. Specs 0266, 0269, 0270 and 0271 declared
the file and the step.
