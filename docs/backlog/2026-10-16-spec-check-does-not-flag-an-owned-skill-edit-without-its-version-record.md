---
type: feat
status: open
created: 2026-10-16
spec: null
reason: null
---

# Spec check does not flag an owned skill edit without its version record

## Symptom

An owned skill changes only together with its version: the Task that edits
`.agents/skills/<name>/` must also declare its `skills/<name>/` mirror, both
`SKILL.md` files whose version fields rise, and
`skills/testdata/owned-skill-versions.json`, and name the
`-record-skill-versions` step. A plan that forgets one of them passes
`roundfix spec check --strict` and fails only at settlement, the same shape
Spec 0274 closes for the Behavior Surface Record of a guide Task.

## Where

`internal/speccheck`, beside the gap Spec 0274 adds for guide Tasks.

## Expected

`spec check` reports a gap when a pending Task declares a file under an owned
skill and does not declare the mirror, the version record and the record step.

## Evidence

Raised on 2026-10-16 while authoring Spec 0274
(`docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/`),
whose task_01, task_02 and task_04 each raise the `roundfix` skill version;
the rule is stated in `docs/agents/specific-repository.md`.
