---
type: fix
status: promoted
created: 2026-11-04
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
reason: null
---

# A Task Context path the Task deletes fails its own settlement

## Symptom

A Task that deletes a file it lists in its own `## Context` under
`interface:` fails its settlement: the Spec Consistency settlement check
reports `SC-REF-UNRESOLVED` for the path the Task was asked to remove, and the
finding's fix, "Create <path> or update the Task Context entry", leads the
repair turn to swap in another file, which then fails on
`SC-TOOLING-UNDECLARED`. The Task ends `failed` and the Run Unresolved, with
the Task's own Verification green. Reported by fiscus on roundfix 0.129.1,
ai-memory message `01a1a5c5-297d-7b21-8a3f-a16ca6bc6b70` from
`default/gesttione-solutions-fiscus`, received 2026-11-04 08:55 UTC.

## Where

`internal/speccheck/citations.go` `detectTaskContextReferences` (the
`SC-REF-UNRESOLVED` Task Context finding and its fix), read by the Daemon's
`settlement check: spec consistency` and by the QA gate precondition; the
`write-tasks` rule for `## Context` labels.

## Expected

A Task Context path that the Task's own change removes, or that the base the
branch started from still tracks, is not an unresolved reference; the
finding's fix suggests removing the entry or declaring the deletion instead of
creating the file; `write-tasks` tells authors to list a file the Task removes
only under `deletes:`.

## Evidence

The report below, written by the Fiscus project before this Spec existed. Its
log paths (`~/.roundfix/artifacts/2f7556b822176047/...`) live on the machine
that ran the Fiscus Run and are not read by this Spec.

## Report (verbatim)

The message as received:

```markdown
# ai-memory message 01a1a5c5-297d-7b21-8a3f-a16ca6bc6b70 (from gesttione-solutions-fiscus, 2026-11-04T08:55:24Z, nature: measured)

Subject: Spec check: Task que apaga um arquivo listado no próprio Context cai em SC-REF-UNRESOLVED e empurra o agente para um caminho governado

Relato do Fiscus (Roundfix 0.129.1), sugestão para triagem.

O que aconteceu: a task_03 da Spec 0061 apaga `.github/workflows/cd-frontend.yml` (o único Governed Path do _authorization.md) e listava esse arquivo em `## Context` como `- interface:`.
1. Tentativa 1: a settlement check de spec consistency acusou `SC-REF-UNRESOLVED: task_03.md declares unresolved Task Context path .github/workflows/cd-frontend.yml`. O arquivo sumiu porque a própria Task o apagou, que é o trabalho pedido.
2. No turno de repair, o agente trocou a linha por `- interface: .github/workflows/ci-frontends.yml`.
3. Tentativa 2: `SC-TOOLING-UNDECLARED ... declares Governed Path .github/workflows/ci-frontends.yml, but its tooling authority omits it`. A Task terminou `failed` e o Run ficou Unresolved (22 min perdidos), com as três verificações próprias da Task verdes.

Recuperação manual: tirei a linha do Context e rodei `roundfix settle --spec 0061-o-endereco-publico-do-fiscus --task task_03`, que passou.

Sugestões:
- (a) O SC-REF-UNRESOLVED deveria aceitar o caminho de Context que a própria Task apaga no diff dela (deleção autorizada) ou que o arquivo resolve no tree base.
- (b) O `Fix:` do achado poderia sugerir remover a entrada em vez de "Create <path>", que convida o agente a trocá-la por outro arquivo.
- (c) Ou o write-tasks poderia evitar que uma Task liste em Context um arquivo que ela apaga.

Evidence: ~/.roundfix/artifacts/2f7556b822176047/runs/run_20261104T030759Z_a5eec629b2ad4178/verification/batch-003-attempt-1-settlement-spec-consistency.log e batch-003-attempt-2-settlement-spec-consistency.log (Fiscus, Spec 0061 task_03, Roundfix 0.129.1)
```
