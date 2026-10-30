---
type: fix
status: promoted
created: 2026-10-27
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
reason: null
---

# A Wave collision with a completed Task refuses after the Run starts

## Symptom

A corrective Task added after reopening (`task_06`, `needs [task_02]`) declares
a file that a completed Task (`task_05`) also declares. `roundfix spec check`
reports `SC-WAVE-COLLISION`, yet `roundfix implement` creates the Run and only
then refuses: `task cycle: refuse Wave before dispatch because Tasks collide`.
The completed Task can never run in parallel with the new one, so the
collision is not real, and the refusal comes after a Run exists.

## Where

`internal/daemon/task_engine.go` `refuseTaskPlanWaveCollisions` and
`spec.Collisions`, which consider every Task of the graph, including
`completed` ones; the Spec Check `SC-WAVE-COLLISION` rule.

## Expected

A Task already `completed` does not collide with a pending one, in the Spec
Check and in the daemon; or `implement` refuses a colliding graph before it
creates the Run.

## Evidence

Quaestor (ai-memory message `01a17e01-45fd-7a81-9052-54c660114746`,
2026-10-27 12:36 BRT), item 2:

```
2. Colisão de onda depois de reabrir (Spec 0027). Acrescentamos a task_06 corretiva (needs [task_02]), com a task_05 já concluída e ambas declarando tests/test_conferencia_estoque.py. O `spec check` acusou SC-WAVE-COLLISION, mas o implement começou mesmo assim e só recusou a onda depois do início do Run ("refuse Wave before dispatch because Tasks collide"), com a task_05 concluída e sem chance de rodar em paralelo. Esperado: não considerar colisão com Task concluída, ou recusar antes de criar o Run.
```
