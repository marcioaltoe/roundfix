---
type: fix
status: promoted
created: 2026-10-26
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
reason: null
---

# Two Tasks of one Wave that change an undeclared test fake lose the Run at the QA gate

## Symptom

In a Fluxus Run, task_01 and task_02 of the same Wave each edited test fakes
that neither Task declared; the fakes appeared only in each Task's
`## Recorded paths`. Both Tasks settled and integrated, then the QA gate's
precondition, `roundfix spec check <slug> --strict`, refused the Run with
`SC-WAVE-COLLISION` on the two settlement commits. The same strict check
was clean before the Run, because no settlement commit existed yet, so the
author had nothing to act on until the Run was lost.

## Where

`internal/spec/collision.go` `Collisions` adds each Task's newest
settlement commit (`priorRunTaskTouches`, source `prior Run settlement
commit`) to its declared and Verification paths;
`internal/speccheck/coherence.go` `detectWaveCollisions` reports the pair;
`internal/daemon/task_engine.go` `specGatePrecondition` runs the check as
the QA gate's precondition. The Spec Consistency settlement check of each
Task compares with a baseline taken when that Task starts, so the pair is
first refused at the gate. The `write-tasks` skill says a Recorded Path
"reserves nothing" and does not say that its settlement commit is still
collision evidence.

## Expected

Prevent the lost Run with the smallest structural change: either the
`write-tasks` skill tells authors to declare a shared test helper or fake
in every Task that may change it, so the strict check reports the collision
while the Spec is authored, or a settlement-time diagnostic names the
undeclared shared path and the fix.

## Evidence

Fluxus, ai-memory message `01a16a54` to `default/marcioaltoe-roundfix`,
2026-10-23, after Fluxus proved Spec 0301 in Run
`run_20261023T134436Z_a7dbc620803968fe`; consumed by the operator on
2026-10-24 at 11:25 BRT and recorded as item 10 of the cycle AO triage
inputs: "QA gate recusou SC-WAVE-COLLISION quando task_01 e task_02 da mesma
onda editaram fakes de teste não declarados (só em Recorded paths); spec
check --strict não prevê. Avaliar: aviso no settlement ou orientação no
write-tasks."
