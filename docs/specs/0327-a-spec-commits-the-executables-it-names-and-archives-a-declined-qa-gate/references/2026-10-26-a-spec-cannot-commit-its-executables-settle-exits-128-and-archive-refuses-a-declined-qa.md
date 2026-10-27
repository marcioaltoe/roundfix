---
type: fix
status: promoted
created: 2026-10-26
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
reason: null
---

# A Spec cannot commit the executables it asks for, settle exits 128 after settling, and archive refuses a declared QA decline

## Symptom

Quaestor reported three defects of Roundfix 0.123.0 in part 2 of its
ai-memory message, quoted verbatim (Portuguese, as received):

```text
2. Defeitos na v0.123.0:
(a) Run run_20261026T190411Z_4760f2955b8df218, Spec 0019-auditor-continuo, task_02. A Verification passou, mas o commit da Task recusou scripts/auditor-continuo.sh e scripts/auditor-continuo-instalar.sh por serem 0755 ("deliberately executable repository files are not valid Work Item output"). A Task virou failed e o Run ficou Unresolved. O esperado: a Spec pedia scripts para o launchd; ou o write-tasks/prompt avisa da regra, ou o _authorization.md pode autorizar executáveis pelo nome. Contornamos com scripts 0644 chamados por /bin/sh e `roundfix settle`. O settle apagou o Run Worktree e saiu com exit 128 depois de imprimir "settled completed".
(b) `roundfix archive 0020-atualizar-dependencias-skills-tooling` (e 0021), com qa: declined, qa_reason e todas as Tasks concluídas, recusou com "no QA Report found". O docs-layout do Baseline diz que um QA declinado declarado tem contrato próprio e não é override. Arquivamos com --qa-override e aprovação de Marcio.
```

Part 1 of the same message (Baseline profile-change frictions) belongs to
Spec 0328 and is out of this entry's scope.

## Where

- The Task commit's stage filter, `FilterStageablePaths` in
  `internal/daemon/task_engine.go`, and its refusal line in
  `publishDroppedStagePath`.
- The Archive Command's eligibility in `internal/spec/archive.go` `Archive`,
  which reads a QA Report for every Spec with a Task Graph.
- The Settle Command's Run-level integration in `internal/cli/settle.go`
  `integrateSettledRun`, which removes the kept Run Worktree without saying so.

## Expected

- A Task commit keeps an untracked executable file at a path the Spec's
  `_authorization.md` names, with its mode; every other executable stays
  refused, and the refusal and the `write-tasks` skill say how to name one.
- `roundfix archive` accepts a Spec whose Task Graph declares `qa: declined`
  with a `qa_reason` once every Task is completed, without an override.
- `roundfix settle` exits with the settled outcome and names each worktree it
  removes, so a caller whose shell stood in it knows why a later command there
  fails.

## Evidence

ai-memory message `01a17abd-20fe-7213-ab5c-f1a1bdb043fe` from
`default/gesttione-solutions-quaestor`, created 2026-10-26 21:22 BRT; it
consolidates the earlier messages `01a17994` and `01a179a8` of the same day.
Run `run_20261026T190411Z_4760f2955b8df218` of Quaestor Spec
`0019-auditor-continuo`; Quaestor Specs `0020-atualizar-dependencias-skills-tooling`
and `0021`.
